package handlers

import (
	"context"
	"github.com/chi-net/weiba/core/types"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"strconv"
)

func MonitorStatus(ctx context.Context, b *bot.Bot, update *models.Update, config types.YmlConfigurationData) {
	// fmt.Println(update.ChatMember.NewChatMember.Member)
	// fmt.Println(update.ChatMember.NewChatMember.Banned)
	// fmt.Println(update.ChatMember.From)
	// fmt.Println(update.ChatMember.Chat)
	chatMember := update.ChatMember
	chatType := "频道"
	if chatMember.Chat.Type != "channel" {
		chatType = "群组"
	}
	msg := chatMember.Chat.Title + "(ID:" + strconv.FormatInt(chatMember.Chat.ID, 10) + ")的成员有新变动！\n"
	if chatMember.NewChatMember.Left != nil {
		userData := chatMember.NewChatMember.Left.User
		msg += userData.FirstName + "(UID:" + strconv.FormatInt(userData.ID, 10) + ")退出了" + chatType + "。"
	} else if chatMember.NewChatMember.Member != nil {
		userData := chatMember.NewChatMember.Member.User
		from := chatMember.From
		via := "直接"
		if chatMember.ViaJoinRequest {
			via = "通过管理员审批"
		} else if chatMember.ViaChatFolderInviteLink {
			via = "通过分享文件夹批量添加"
		}
		msg += userData.FirstName + "(UID:" + strconv.FormatInt(userData.ID, 10) + ")" + via + "加入了" + chatType + "。\n"
		msg += "邀请来自" + from.FirstName + from.LastName + "(UID:" + strconv.FormatInt(from.ID, 10) + ")"
	} else if chatMember.NewChatMember.Banned != nil {
		userData := chatMember.NewChatMember.Banned.User
		msg += userData.FirstName + userData.LastName + "(UID:" + strconv.FormatInt(userData.ID, 10) + ")" + "被封禁了。"
	}
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: config.AdminUID,
		Text:   msg,
	})
	if err != nil {
		return
	}
}
