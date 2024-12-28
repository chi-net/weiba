package core

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"strconv"
)

func HandleMessage(ctx context.Context, b *bot.Bot, update *models.Update, config YmlConfigurationData, authmaps AuthMaps) {
	message := update.Message
	chat := message.Chat
	senderid := message.From.ID
	sender := message.From
	if config.Features.AnonymousChat {
		if senderid == config.AdminUID && message.ReplyToMessage != nil {
			replyto := getUIDinMessage(message.ReplyToMessage.Text)
			if authmaps.ChatOpened[replyto] {
				b.CopyMessage(ctx, &bot.CopyMessageParams{
					ChatID:     replyto,
					FromChatID: message.Chat.ID,
					MessageID:  message.ID,
				})
			} else {
				msg := "非常抱歉，此用户关闭了聊天。"
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: chat.ID,
					Text:   msg,
				})
			}
		}
		if authmaps.ChatOpened[senderid] && senderid != config.AdminUID {
			msg := sender.FirstName + " " + sender.LastName + "(UID:" + strconv.FormatInt(senderid, 10) + ";username:" + sender.Username + ")发送了一条消息\n你可以回复本消息来回复他。"
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: config.AdminUID,
				Text:   msg,
			})
			b.ForwardMessage(ctx, &bot.ForwardMessageParams{
				ChatID:     config.AdminUID,
				MessageID:  message.ID,
				FromChatID: chat.ID,
			})
		}
	}
	if !authmaps.ChatOpened[senderid] && config.AdminUID != senderid {
		msg := "非常抱歉，您暂未处于任何验证进程中。"
		if config.Features.AnonymousChat {
			msg += "\n如果您想和管理员进行聊天，请输入 /chat 以开始聊天。"
		}
		if config.AdminUID == senderid {
			msg += "\n您是本bot管理员， /chat 和自己聊天可能有很创的表现，目前暂未修复，请谨慎输入 /chat ！"
		}
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chat.ID,
			Text:   msg,
		})
	}
}
