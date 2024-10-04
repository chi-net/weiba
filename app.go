package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var userAuthDataMap = make(map[int64]int64)
var userAuthGroupIdsMap = make(map[int64][]int64)
var userAuthGroupMessagesMap = make(map[int64][]int64)
var userAuthStepsMap = make(map[int64]int64)

var data importedGICAuthData

// Administrator UID
const adminUid int64 = 1

// Token from @BotFather
const botToken string = ""

type importedGICAuthData struct {
	LastUpdate int            `json:"lastupdate"`
	Data       []gicGroupData `json:"data"`
}

type gicGroupData struct {
	ID   int64         `json:"i"`
	Data []gicAuthData `json:"d"`
}

type gicAuthData struct {
	URL  int64  `json:"u"`
	Text string `json:"t"`
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Read the file contents
	byteValue, err := os.ReadFile("data.json")
	if err != nil {
		log.Fatal(err)
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

	b, err := bot.New(botToken, opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	print(update)
	if update.ChatJoinRequest != nil {
		fmt.Println(update.ChatJoinRequest)
		userAuthDataMap[update.ChatJoinRequest.UserChatID] = update.ChatJoinRequest.Chat.ID
		userAuthStepsMap[update.ChatJoinRequest.UserChatID] = 0
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.ChatJoinRequest.UserChatID,
			Text:   "您正在尝试加入" + update.ChatJoinRequest.Chat.Title + "！\n您可以选择自助验证(Group in Common验证)或者取消验证，等待管理员审核。\n如果需要自助验证，请输入1\n请注意，其他输入均视作取消自助验证操作，您也可以在晚些时候再次点击申请以再次触发自助验证。",
		})

		// admin notifier
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: adminUid,
			Text:   update.ChatJoinRequest.From.FirstName + update.ChatJoinRequest.From.LastName + "(UID:" + strconv.FormatInt(update.ChatJoinRequest.From.ID, 10) + ")正在尝试加入" + update.ChatJoinRequest.Chat.Title + "！",
		})

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
					if data.Data[userAuthGroupIdsMap[chatid][i]].Data[userAuthGroupMessagesMap[chatid][i]].Text == update.Message.Text {
						b.SendMessage(ctx, &bot.SendMessageParams{
							ChatID: chatid,
							Text:   "验证成功，喜欢您来，欢迎加入！",
						})
						b.ApproveChatJoinRequest(ctx, &bot.ApproveChatJoinRequestParams{
							ChatID: userAuthDataMap[chatid],
							UserID: chatid,
						})

						b.SendMessage(ctx, &bot.SendMessageParams{
							ChatID: adminUid,
							Text:   update.Message.Chat.FirstName + update.Message.Chat.LastName + "(UID:" + strconv.FormatInt(update.Message.Chat.ID, 10) + ")已经通过GIC验证",
						})

						return
					}
				}
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "对不起，验证失败，您可以在晚些时候再次点击加入申请来触发自助验证。",
				})
				b.DeclineChatJoinRequest(ctx, &bot.DeclineChatJoinRequestParams{
					ChatID: userAuthDataMap[update.Message.Chat.ID],
					UserID: update.Message.Chat.ID,
				})

				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: adminUid,
					Text:   update.Message.Chat.FirstName + update.Message.Chat.LastName + "(UID:" + strconv.FormatInt(update.Message.Chat.ID, 10) + ")并没有通过GIC验证",
				})

				return
			}
			if update.Message.Text == "1" {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "您选择了自助验证。\n自助验证说明：接下来将会给您发送10个 t.me 的链接，其对应了10个私有频道/群组，您需要做的是从上到下依次点击这些链接，选择一个可以访问的链接地址，并将其包含的文本内容复制到这里。(你可以电脑鼠标右键/手机轻点选择复制文本/Copy Text以复制内容)",
				})
				userAuthStepsMap[update.Message.Chat.ID] = 1
				message := "10个链接如下，请逐个点击，直到寻找到您可以访问的频道或群组，并将其文本全文复制至此。"
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

					message += "\nhttps://t.me/c/" + strconv.FormatInt(data.Data[num].ID, 10) + "/" + strconv.FormatInt(data.Data[num].Data[num2].URL, 10)
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
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: adminUid,
					Text:   update.Message.Chat.FirstName + update.Message.Chat.LastName + "(UID:" + strconv.FormatInt(update.Message.Chat.ID, 10) + ")选择了管理员审核，您可前往频道/群组页进行批准或拒绝",
				})

			}
		} else {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "非常抱歉，您暂未处于任何验证进程中。",
			})
		}
	}
}
