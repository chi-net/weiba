package handlers

import (
	"context"
	"github.com/chi-net/weiba/core/types"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"strconv"
)

func OpenChatHandler(ctx context.Context, b *bot.Bot, update *models.Update, config types.YmlConfigurationData, authmaps types.AuthMaps) {
	chatid := update.Message.Chat.ID
	user := update.Message.From
	if authmaps.ChatOpened[chatid] {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatid,
			Text:   "您的消息正在被传输到管理员处，无需再次开启聊天！",
		})
	} else {
		authmaps.ChatOpened[chatid] = true
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatid,
			Text:   "聊天连接已经建立，接下来您在本bot中输入的所有内容将会被转发至管理员处，管理员也会通过这个bot对你的消息进行实时回复。\n如果您想关闭聊天，请输入 /close\n如果长时间未响应或聊天一半后中断，请再次尝试输入 /chat 。",
		})
		adminMsg := "用户 " + user.FirstName + " " + user.LastName
		adminMsg += "(UID:" + strconv.FormatInt(user.ID, 10) + ";Username:" + user.Username + ")与你发起了对话！\n"
		adminMsg += "你可以回复这条消息来和他进行交流"
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: config.AdminUID,
			Text:   adminMsg,
		})
	}
}
