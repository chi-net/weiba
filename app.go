package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"gopkg.in/yaml.v3"
	"log"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var userAuthDataMap = make(map[int64]int64)
var userAuthGroupIdsMap = make(map[int64][]int64)
var userAuthGroupMessagesMap = make(map[int64][]int64)
var userAuthStepsMap = make(map[int64]int64)

var data importedGICAuthData
var config ymlConfigurationData

type importedGICAuthData struct {
	LastUpdate int            `json:"lastupdate"`
	Data       []gicGroupData `json:"data"`
}

type gicGroupData struct {
	ID   int64    `json:"i"`
	Data []string `json:"d"`
}

type ymlConfigurationData struct {
	AdminUID        int64  `yaml:"admin_uid"`
	BotToken        string `yaml:"bot_token"`
	Mode            string `yaml:"mode"`
	TransformDigits int    `yaml:"transform_digits"`
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func transback(from string) uint64 {
	result := uint64(0)
	strlist := "abcdefghijklmnopqrstuvwxyz!@#$ABCDEFGHIJKLMNOPQRSTUVWXYZ%^&*1234567890()-=_+[]{}|\\:;<>?,./`~"

	// Create a map for quick lookup of character indices
	indexMap := make(map[rune]int)
	for i, c := range strlist {
		indexMap[c] = i
	}

	for _, char := range from {
		index, exists := indexMap[char]
		if !exists {
			fmt.Printf("Character '%c' not found in strlist\n", char)
			return 0 // or handle the error appropriately
		}
		// fmt.Print(index, " ")                                // Print the index for debugging
		result = result*uint64(len(strlist)) + uint64(index) // Update the result
	}
	// fmt.Println(result)
	// fmt.Println() // Print a newline after indices
	return result
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Read the file contents
	byteValue, err := os.ReadFile("data.json")
	byteValue2, err := os.ReadFile("config.yml")

	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(byteValue2, &config)
	if err != nil {
		log.Fatal(err)
	}

	if config.Mode == "env" {
		config.AdminUID, _ = strconv.ParseInt(getEnv("ADMIN_UID", "-1"), 10, 64)
		config.BotToken = getEnv("BOT_TOKEN", "")
		if config.BotToken == "" {
			panic("BotToken is not defined in environmental variables! Stopping this application.")
		}
		config.TransformDigits, _ = strconv.Atoi(getEnv("TRANSFORM_DIGITS", "-1"))
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
	}

	b, err := bot.New(config.BotToken, opts...)
	if err != nil {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/info", bot.MatchTypeExact, infoHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/info@mrweibabot", bot.MatchTypeExact, infoHandler)

	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.ChatJoinRequest != nil {
		//fmt.Println(update.ChatJoinRequest)
		userAuthDataMap[update.ChatJoinRequest.UserChatID] = update.ChatJoinRequest.Chat.ID
		userAuthStepsMap[update.ChatJoinRequest.UserChatID] = 0
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.ChatJoinRequest.UserChatID,
			Text:   "您正在尝试加入" + update.ChatJoinRequest.Chat.Title + "！\n您可以选择自助验证(Group in Common验证)或者取消验证，等待管理员审核。\n如果需要自助验证，请输入1\n请注意，其他输入均视作取消自助验证操作，您也可以在晚些时候再次点击申请以再次触发自助验证。\n如果您没有回复任何内容，将默认视作为等待管理员手动批准。",
		})

		// admin notifier
		if config.AdminUID != -1 {
			_, err = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: config.AdminUID,
				Text:   update.ChatJoinRequest.From.FirstName + update.ChatJoinRequest.From.LastName + "(UID:" + strconv.FormatInt(update.ChatJoinRequest.From.ID, 10) + ")正在尝试加入" + update.ChatJoinRequest.Chat.Title + "！",
			})
		}

		if err != nil {
			return
		}
	} else if update.Message != nil && update.Message.Chat.Type == models.ChatTypePrivate {
		chatid := update.Message.Chat.ID
		if ok, _ := userAuthDataMap[chatid]; ok != 0 {
			// The user has begun authentication process
			if userAuthStepsMap[chatid] == 1 {
				for i := 0; i < len(userAuthGroupIdsMap[chatid]); i++ {
					// fmt.Println(data.Data[userAuthGroupIdsMap[chatid][i]].Data[userAuthGroupMessagesMap[chatid][i]].Text)
					splited := strings.Split(data.Data[userAuthGroupIdsMap[chatid][i]].Data[userAuthGroupMessagesMap[chatid][i]], " ")[1]

					// Convert string to a slice of runes
					var runes []rune
					for _, r := range update.Message.Text {
						runes = append(runes, r)
					}
					msg := string(runes)
					// sha1
					h := sha1.New()
					h.Write([]byte(msg))
					hash := h.Sum(nil)
					hashHex := hex.EncodeToString(hash)
					//fmt.Printf("SHA-1 Hash (hex): %s\n", hashHex)
					hashInt, _ := strconv.ParseUint(hashHex[:16], 16, 64)
					modulus := uint64(math.Pow(92, float64(config.TransformDigits)))
					result := hashInt % modulus
					//fmt.Println(hashInt)
					//fmt.Println(result)
					if transback(splited) == result {
						b.SendMessage(ctx, &bot.SendMessageParams{
							ChatID: chatid,
							Text:   "验证成功，喜欢您来，欢迎加入！",
						})
						b.ApproveChatJoinRequest(ctx, &bot.ApproveChatJoinRequestParams{
							ChatID: userAuthDataMap[chatid],
							UserID: chatid,
						})
						if config.AdminUID != -1 {
							b.SendMessage(ctx, &bot.SendMessageParams{
								ChatID: config.AdminUID,
								Text:   update.Message.Chat.FirstName + update.Message.Chat.LastName + "(UID:" + strconv.FormatInt(update.Message.Chat.ID, 10) + ")已经通过GIC验证",
							})
						}
						delete(userAuthDataMap, update.Message.Chat.ID)
						delete(userAuthStepsMap, update.Message.Chat.ID)
						delete(userAuthGroupIdsMap, update.Message.Chat.ID)
						delete(userAuthGroupMessagesMap, update.Message.Chat.ID)
						return
					}
				}
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "对不起，验证失败，您可以在晚些时候再次点击加入申请来触发自助验证。",
				})
				if config.AdminUID != -1 {
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: config.AdminUID,
						Text:   update.Message.Chat.FirstName + update.Message.Chat.LastName + "(UID:" + strconv.FormatInt(update.Message.Chat.ID, 10) + ")并没有通过GIC验证",
					})
				}
				delete(userAuthDataMap, update.Message.Chat.ID)
				delete(userAuthStepsMap, update.Message.Chat.ID)
				delete(userAuthGroupIdsMap, update.Message.Chat.ID)
				delete(userAuthGroupMessagesMap, update.Message.Chat.ID)
				return
			}
			if update.Message.Text == "1" {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "您选择了自助验证。\n自助验证说明：接下来将会给您发送10个 t.me 的链接，其对应了10个私有频道/群组，您需要做的是从上到下依次点击这些链接，选择一个可以访问的链接地址，并将其包含的文本内容复制到这里。(你可以电脑鼠标右键/手机轻点选择复制文本/Copy Text以复制内容)",
				})
				userAuthStepsMap[update.Message.Chat.ID] = 1
				message := "10个链接如下，请逐个点击，直到寻找到您可以访问的频道或群组，并将其文本全文复制至此。\n请注意：如果您无法访问这些链接的内容，您可以随意输入一个内容以结束验证进程并稍后再次点击申请重试。"
				for i := 0; i <= 10; i++ {
					source := rand.NewSource(time.Now().UnixNano())
					r := rand.New(source) // Create a new Rand instance
					num := r.Int63n(int64(len(data.Data)))

					userAuthGroupIdsMap[chatid] = append(userAuthGroupIdsMap[chatid], num)

					source = rand.NewSource(time.Now().UnixNano())
					r = rand.New(source) // Create a new Rand instance
					// fmt.Println(len(data.Data[num].Data))
					num2 := r.Int63n(int64(len(data.Data[num].Data)))

					userAuthGroupMessagesMap[chatid] = append(userAuthGroupMessagesMap[chatid], num2)

					parts := strings.Split(data.Data[num].Data[num2], " ")

					message += "\nhttps://t.me/c/" + strconv.FormatInt(data.Data[num].ID, 10) + "/" + strconv.FormatUint(transback(parts[0]), 10)
				}
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   message,
				})
				//userAuthDataMap[update.Message.Chat.ID] = 0
			} else {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "您输入了其他内容，因此您需等待管理员验证，您也可以再次点击加入申请来触发自助验证。",
				})
				delete(userAuthDataMap, update.Message.Chat.ID)
				delete(userAuthStepsMap, update.Message.Chat.ID)
				if config.AdminUID != -1 {
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: config.AdminUID,
						Text:   update.Message.Chat.FirstName + update.Message.Chat.LastName + "(UID:" + strconv.FormatInt(update.Message.Chat.ID, 10) + ")选择了管理员审核，您可前往频道/群组页进行批准或拒绝",
					})
				}
			}
		} else {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "非常抱歉，您暂未处于任何验证进程中。",
			})
		}
	}
}

func infoHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.From.ID != config.AdminUID && config.AdminUID != -1 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "weibabot@canary - A bot which can process authentication of private groups and channels\nOpenSource: https://github.com/chi-net/weiba",
		})
	} else {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		// Get CPU Info
		cpuInfo, err := cpu.Info()
		if err != nil {
			fmt.Println("Error getting CPU info:", err)
			return
		}
		// Get CPU usage percentage
		cpuPercent, err := cpu.Percent(2*time.Second, false)
		if err != nil {
			fmt.Println("Error getting CPU percent:", err)
			return
		}
		// Get Memory Info
		virtualMem, err := mem.VirtualMemory()
		if err != nil {
			fmt.Println("Error getting memory info:", err)
			return
		}
		// Get Disk Info
		diskInfo, err := disk.Usage("/")
		if err != nil {
			fmt.Println("Error getting disk info:", err)
			return
		}
		hostInfo, err := host.Info()
		if err != nil {
			log.Fatalf("Error getting host info: %v", err)
		}

		msg := "weibabot@canary - A bot which can process authentication of private groups and channels\nSystem Info:\n"
		msg += "Host:" + hostInfo.Platform + "(Kernel:" + hostInfo.OS + " " + hostInfo.KernelVersion + ")\n"
		msg += "CPU: " + cpuInfo[0].ModelName + "(" + strconv.Itoa(len(cpuInfo)) + ") " + strconv.FormatFloat(cpuInfo[0].Mhz/1000, 'f', 2, 64) + "GHz " + strconv.FormatFloat(cpuPercent[0], 'f', 2, 64) + "% Used\n"
		msg += "Memory: " + strconv.FormatFloat(float64(virtualMem.Total/1024/1024), 'f', 2, 64) + "MB Total;" + strconv.FormatFloat(float64(virtualMem.Used/1024/1024), 'f', 2, 64) + "MB Used\n"
		msg += "Disk: " + strconv.FormatFloat(float64(diskInfo.Total/1024/1024/1024), 'f', 2, 64) + "GB Total;" + strconv.FormatFloat(float64(diskInfo.Used/1024/1024/1024), 'f', 2, 64) + "GB Used\n"
		msg += "Application:\n"
		msg += "Current active goroutines:" + strconv.Itoa(runtime.NumGoroutine()) + "\n"
		msg += "Current occupied Memory:" + strconv.FormatFloat(float64(memStats.Alloc/1024/1024), 'f', 2, 64) + "MB"
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   msg,
		})
	}
}
