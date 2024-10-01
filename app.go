package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Send any text message to the bot after the bot has been started
var userAuthMap = make(map[int64]int64)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	
	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}
	
	b, err := bot.New("[Your token]", opts...)
	if err != nil {
		panic(err)
	}
	
	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	print(update)
	if update.ChatJoinRequest != nil {
		fmt.Println(update.ChatJoinRequest)
		userAuthMap[update.ChatJoinRequest.UserChatID] = update.ChatJoinRequest.Chat.ID
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.ChatJoinRequest.UserChatID,
			Text:   "您正在尝试加入" + update.ChatJoinRequest.Chat.Title + "！\n您可以输入xxx以继续加入",
		})
		
		if err != nil {
			return
		}
	} else if update.Message != nil && update.Message.Chat.Type == models.ChatTypePrivate {
		if ok, _ := userAuthMap[update.Message.Chat.ID]; ok != 0 {
			if update.Message.Text == "xxx" {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "验证成功，欢迎加入",
				})
				b.ApproveChatJoinRequest(ctx, &bot.ApproveChatJoinRequestParams{
					ChatID: userAuthMap[update.Message.Chat.ID],
					UserID: update.Message.Chat.ID,
				})
				userAuthMap[update.Message.Chat.ID] = 0
			} else {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "验证失败，请稍后再申请",
				})
				b.DeclineChatJoinRequest(ctx, &bot.DeclineChatJoinRequestParams{
					ChatID: userAuthMap[update.Message.Chat.ID],
					UserID: update.Message.Chat.ID,
				})
				userAuthMap[update.Message.Chat.ID] = 0
			}
		} else {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "您暂未处于任何验证进程中",
			})
		}
	}
}
