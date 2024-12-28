package core

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"strconv"
	"strings"
)

func HandleAuthChallenge(ctx context.Context, b *bot.Bot, update *models.Update, config YmlConfigurationData, authmaps AuthMaps, data ImportedGICAuthData) {
	chatid := update.Message.Chat.ID
	if ok, _ := authmaps.Data[chatid]; ok != 0 {
		// The user has begun authentication process
		if authmaps.Steps[chatid] == 2 {
			for i := 0; i < len(authmaps.GroupIds[chatid]); i++ {
				// fmt.Println(data.Data[userAuthGroupIds[chatid][i]].Data[userAuthGroupMessages[chatid][i]].Text)
				encoded := strings.Split(data.Data[authmaps.GroupIds[chatid][i]].Data[authmaps.GroupMessages[chatid][i]], " ")[1]

				if check(update.Message.Text, encoded, config) {
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: chatid,
						Text:   "验证成功，喜欢您来，欢迎加入！",
					})
					b.ApproveChatJoinRequest(ctx, &bot.ApproveChatJoinRequestParams{
						ChatID: authmaps.Data[chatid],
						UserID: chatid,
					})
					if config.AdminUID != -1 {
						b.SendMessage(ctx, &bot.SendMessageParams{
							ChatID: config.AdminUID,
							Text:   update.Message.Chat.FirstName + update.Message.Chat.LastName + "(UID:" + strconv.FormatInt(update.Message.Chat.ID, 10) + ")已经通过GIC验证",
						})
					}
					delete(authmaps.Data, update.Message.Chat.ID)
					delete(authmaps.Steps, update.Message.Chat.ID)
					delete(authmaps.GroupIds, update.Message.Chat.ID)
					delete(authmaps.GroupMessages, update.Message.Chat.ID)
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
			delete(authmaps.Data, update.Message.Chat.ID)
			delete(authmaps.Steps, update.Message.Chat.ID)
			delete(authmaps.GroupIds, update.Message.Chat.ID)
			delete(authmaps.GroupMessages, update.Message.Chat.ID)
			return
		}
		if update.Message.Text == "1" {
			message := "您选择了自助验证。\n"
			message += "自助验证说明：接下来将会给您发送10个 t.me 的链接，其对应了10个私有频道/群组的聊天信息。\n"
			message += "您需要做的是从上到下依次点击这些链接，选择一个可以访问的链接地址，并将其包含的文本内容复制到这里。\n"
			message += "小提示: 你可以电脑鼠标右键/手机轻点选择复制文本/Copy Text以复制内容"

			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   message,
			})

			authmaps.Steps[update.Message.Chat.ID] = 2
			message = "10个链接如下，请逐个点击，直到寻找到您可以访问的频道或群组，并将其文本全文复制至此。\n"
			message += "请注意：如果您无法访问这些链接的内容，您可以随意输入一个内容以结束验证进程并稍后再次点击申请重试。"
			for i := 0; i <= 10; i++ {
				num, num2 := generate(data)
				authmaps.GroupMessages[chatid] = append(authmaps.GroupMessages[chatid], num2)
				authmaps.GroupIds[chatid] = append(authmaps.GroupIds[chatid], num)
				parts := strings.Split(data.Data[num].Data[num2], " ")
				message += "\nhttps://t.me/c/" + strconv.FormatInt(data.Data[num].ID, 10) + "/" + strconv.FormatUint(decode(parts[0]), 10)
			}
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   message,
			})
			//userAuthData[update.Message.Chat.ID] = 0
		} else {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "您输入了其他内容，因此您需等待管理员验证，您也可以再次点击加入申请来触发自助验证。",
			})
			delete(authmaps.Data, update.Message.Chat.ID)
			delete(authmaps.Steps, update.Message.Chat.ID)

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
