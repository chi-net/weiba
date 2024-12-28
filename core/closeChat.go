package core

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"strconv"
)

func CloseChatHandler(ctx context.Context, b *bot.Bot, update *models.Update, config YmlConfigurationData, authmaps AuthMaps) {
	chatid := update.Message.Chat.ID
	user := update.Message.From
	if !authmaps.ChatOpened[chatid] {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatid,
			Text:   "您并未开启，无需关闭聊天！",
		})
	} else {
		authmaps.ChatOpened[chatid] = true
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatid,
			Text:   "聊天连接已经关闭。\n如果您想重新开启聊天，请再次输入 /chat",
		})
		adminMsg := "用户 " + user.FirstName + " " + user.LastName
		adminMsg += "(UID:" + strconv.FormatInt(user.ID, 10) + ";Username:" + user.Username + ")关闭了对话。\n"
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: config.AdminUID,
			Text:   adminMsg,
		})
	}
}
