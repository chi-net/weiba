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
	authmaps.ChatOpened = make(map[int64]bool)

	// Read configuration.
	byteValue2, err := os.ReadFile("config.yml")

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(byteValue2, &config)
	if err != nil {
		panic(err)
	}

	if config.Mode == "env" {
		// in env mode, we automatically disabled whitelist mode, if you opened unpin channel posts feature, it will influence all groups joined in.
		config.AdminUID, _ = strconv.ParseInt(getEnv("ADMIN_UID", "-1"), 10, 64)
		config.BotToken = getEnv("BOT_TOKEN", "")
		if config.BotToken == "" {
			panic("BotToken is not defined in environmental variables! Stopping this application.")
		}
		config.TransformDigits, _ = strconv.Atoi(getEnv("TRANSFORM_DIGITS", "-1"))
		config.Features.MonitorMembers, _ = strconv.ParseBool(getEnv("MONITOR_MEMBERS", "false"))
		config.Features.UnpinChannelPosts, _ = strconv.ParseBool(getEnv("UNPIN_CHANNEL_POSTS", "false"))
		config.Features.GicAuth, _ = strconv.ParseBool(getEnv("GIC_AUTH", "false"))
		config.Features.Debug, _ = strconv.ParseBool(getEnv("DEBUG", "false"))
		config.Features.AnonymousChat, _ = strconv.ParseBool(getEnv("ANONYMOUS_CHAT", "false"))
	}

	if config.Features.GicAuth {
		byteValue, err := os.ReadFile("data.json")
		// Unmarshal the JSON into the struct
		err = json.Unmarshal(byteValue, &data)
		if err != nil {
			log.Fatal(err)
		}
	}

	if (config.Features.AnonymousChat || config.Features.Debug) && config.AdminUID == -1 {
		panic("You don't set any administrator UID for features that needs it!")
	}

	// initialize the bot.
	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
		bot.WithAllowedUpdates(bot.AllowedUpdates{
			"chat_member",
			"chat_join_request",
			"message",
			"channel_post",
		}),
	}

	// create the bot instance.
	b, err := bot.New(config.BotToken, opts...)
	if err != nil {
		panic(err)
	}

	// register commands
	b.RegisterHandler(bot.HandlerTypeMessageText, "/info", bot.MatchTypeExact,
		func(ctx context.Context, b *bot.Bot, update *models.Update) {
			core.InfoHandler(ctx, b, update, config)
		})
	b.RegisterHandler(bot.HandlerTypeMessageText, "/chatid", bot.MatchTypeExact,
		func(ctx context.Context, b *bot.Bot, update *models.Update) {
			core.ChatIDHandler(ctx, b, update, config)
		})
	b.RegisterHandler(bot.HandlerTypeMessageText, "/chat", bot.MatchTypeExact,
		func(ctx context.Context, b *bot.Bot, update *models.Update) {
			core.OpenChatHandler(ctx, b, update, config, authmaps)
		})
	b.RegisterHandler(bot.HandlerTypeMessageText, "/close", bot.MatchTypeExact,
		func(ctx context.Context, b *bot.Bot, update *models.Update) {
			core.CloseChatHandler(ctx, b, update, config, authmaps)
		})

	// start the bot
	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// checkout channel monitor status
	if config.AdminUID != -1 && config.Features.MonitorMembers && update.ChatMember != nil {
		core.MonitorStatus(ctx, b, update, config)
	}

	// handle ChatJoinRequests
	if update.ChatJoinRequest != nil {
		if config.Features.GicAuth {
			core.HandleJoinRequest(ctx, b, update, config, authmaps)
		} else {
			msg := "[Debug] Detected chat join request but you don't enable this feature yet.\n"
			msg += "If you want to enable it, please configure it in config.yml."
			core.SendDebugMessage(msg, ctx, b, config)
		}
	}
	if update.Message != nil && update.Message.Chat.Type == models.ChatTypePrivate {
		//fmt.Println(authmaps.Steps[update.Message.Chat.ID])
		if authmaps.Steps[update.Message.Chat.ID] == 1 || authmaps.Steps[update.Message.Chat.ID] == 2 {
			core.HandleAuthChallenge(ctx, b, update, config, authmaps, data)
		} else {
			core.HandleMessage(ctx, b, update, config, authmaps)
		}
	}

	if update.Message != nil && (update.Message.Chat.Type == models.ChatTypeGroup || update.Message.Chat.Type == models.ChatTypeSupergroup) {
		if config.Features.Tietie {
			core.HandleTietie(ctx, b, update)
		}
	}

	// handle Unpin Messages
	if update.Message != nil && update.Message.SenderChat != nil && config.Features.UnpinChannelPosts {
		if update.Message.Chat.Type == models.ChatTypeSupergroup && update.Message.SenderChat.Type == models.ChatTypeChannel {
			for _, val := range config.Whitelists.UnpinChannelPosts {
				if val == update.Message.Chat.ID {
					core.HandleChannelPosts(ctx, b, update, config)
				}
			}
		} else {
			msg := "[Debug] Detected channel posted in chat but you don't enable Unpin Messages in this group yet.\n"
			msg += "If you want to enable it, please configure it in config.yml."
			core.SendDebugMessage(msg, ctx, b, config)
		}
	}
}
