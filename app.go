package main

import (
	"context"
	"encoding/json"
	"github.com/chi-net/weiba/core"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/signal"
	"strconv"
)

var authmaps = core.AuthMaps{
	Data:          make(map[int64]int64),
	GroupIds:      make(map[int64][]int64),
	GroupMessages: make(map[int64][]int64),
	Steps:         make(map[int64]int64),
}

var data core.ImportedGICAuthData
var config core.YmlConfigurationData

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	authmaps.Data = make(map[int64]int64)
	authmaps.GroupIds = make(map[int64][]int64)

	// Read the file contents
	byteValue, err := os.ReadFile("data.json")
	byteValue2, err := os.ReadFile("config.yml")

	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(byteValue2, &config)
	if err != nil {
		panic("Can not read config.yml! If you are using a container, please ensure your application's folder has config.yml")
	}

	if config.Mode == "env" {
		// in env mode, we automatically disabled whitelist mode, if you opened unpin channel posts feature, it will influence all groups joined in.
		config.AdminUID, _ = strconv.ParseInt(getEnv("ADMIN_UID", "-1"), 10, 64)
		config.BotToken = getEnv("BOT_TOKEN", "")
		if config.BotToken == "" {
			panic("BotToken is not defined in environmental variables! Stopping this application.")
		}
		config.TransformDigits, _ = strconv.Atoi(getEnv("TRANSFORM_DIGITS", "-1"))
		config.EnhancedMonitorChannelMembers, _ = strconv.ParseBool(getEnv("ENHANCED_MONITOR_CHANNEL_MEMBERS", "false"))
		config.UnpinChannelPosts, _ = strconv.ParseBool(getEnv("UNPIN_CHANNEL_POSTS", "false"))
	}

	// Unmarshal the JSON into the struct
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		log.Fatal(err)
	}

	// Print the data
	//for _, d := range data.Data {
	//	fmt.Printf(strconv.FormatInt(d.ID, 10) + "\n")
	//}

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
		bot.WithAllowedUpdates(bot.AllowedUpdates{
			"chat_member",
			"chat_join_request",
			"message",
		}),
	}

	b, err := bot.New(config.BotToken, opts...)
	if err != nil {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/info", bot.MatchTypeExact,
		func(ctx context.Context, b *bot.Bot, update *models.Update) {
			core.InfoHandler(ctx, b, update, config)
		})
	b.RegisterHandler(bot.HandlerTypeMessageText, "/chatid", bot.MatchTypeExact,
		func(ctx context.Context, b *bot.Bot, update *models.Update) {
			core.ChatIDHandler(ctx, b, update, config)
		})

	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// checkout channel monintor status
	if config.AdminUID != -1 && config.EnhancedMonitorChannelMembers && update.ChatMember != nil {
		core.MonitorStatus(ctx, b, update, config)
	}
	// handle ChatJoinRequests
	if update.ChatJoinRequest != nil {
		core.HandleJoinRequest(ctx, b, update, config, authmaps)
	} else if update.Message != nil && update.Message.Chat.Type == models.ChatTypePrivate {
		core.HandleAuthChallenge(ctx, b, update, config, authmaps, data)
	}
}
