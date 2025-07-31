package handlers

import (
	"context"
	"github.com/chi-net/weiba/core/types"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"strconv"
)

func HandleJoinRequest(ctx context.Context, b *bot.Bot, update *models.Update, config types.YmlConfigurationData, authmaps types.AuthMaps) {

	// Checking whitelist...
	if len(config.Whitelists.GicAuth) != 0 {
		flag := false
		for _, val := range config.Whitelists.GicAuth {
			// fmt.Println(val)
			if val == update.ChatJoinRequest.Chat.ID {
				flag = true
			}
		}
		if !flag {
			userMessage := "我们遇到了一些无法修复的内部错误，请您稍后再试。\n"
			userMessage += "如果您是普通用户，请联系频道/群组的管理员；如果您是管理员，请使用管理员帐号start这个bot并查看报错信息。"
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.ChatJoinRequest.UserChatID,
				Text:   userMessage,
			})
			adminMessage := "发生错误：有非白名单之外的频道/群组邀请了您的bot并将它用作未经验证的GIC认证服务。\n"
			adminMessage += "相关信息：\n"
			adminMessage += "名称：" + update.ChatJoinRequest.Chat.Title + "\n"
			adminMessage += "ID: " + strconv.FormatInt(update.ChatJoinRequest.Chat.ID, 10) + "\n"
			adminMessage += "Username: " + update.ChatJoinRequest.Chat.Username + "\n"
			adminMessage += "如果您确定这是您自己的频道/群组，请在 config.yml 的 whitelist，gic_auth: 一栏下添加如下内容并重启程序：\n"
			adminMessage += "- " + strconv.FormatInt(update.ChatJoinRequest.Chat.ID, 10)
			if config.AdminUID != -1 {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: config.AdminUID,
					Text:   adminMessage,
				})
			}
			return
		}
	}

	//fmt.Println(update.ChatJoinRequest)
	authmaps.Data[update.ChatJoinRequest.UserChatID] = update.ChatJoinRequest.Chat.ID
	authmaps.Steps[update.ChatJoinRequest.UserChatID] = 1

	requestFrom := update.ChatJoinRequest.From
	joinChat := update.ChatJoinRequest.Chat

	userMessage := "您正在尝试加入" + update.ChatJoinRequest.Chat.Title + "！\n"
	userMessage += "您可以选择自助验证(Group in Common验证)或者取消验证，等待管理员审核。\n"
	userMessage += "如果需要自助验证，请输入1\n"
	userMessage += "请注意，其他输入均视作取消自助验证操作，您也可以在晚些时候再次点击申请以再次触发自助验证。\n"
	userMessage += "如果您没有回复任何内容，将默认视作为等待管理员手动批准。"

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.ChatJoinRequest.UserChatID,
		Text:   userMessage,
	})

	adminMessage := requestFrom.FirstName + requestFrom.LastName + "(UID:" + strconv.FormatInt(requestFrom.ID, 10) + ")"
	adminMessage += "正在尝试加入" + joinChat.Title + "(ID:" + strconv.FormatInt(joinChat.ID, 10) + ")！"

	// admin notifier
	if config.AdminUID != -1 {
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: config.AdminUID,
			Text:   adminMessage,
		})
	}

	if err != nil {
		return
	}
}
