package handlers

import (
	"context"
	"github.com/chi-net/weiba/core/store"
	"github.com/chi-net/weiba/core/types"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"strconv"
)

func GetRankingHandler(ctx context.Context, b *bot.Bot, update *models.Update, config types.YmlConfigurationData) {
	message := update.Message
	chat := message.Chat
	senderid := message.From.ID
	title := message.From.FirstName + " " + message.From.LastName
	if message.SenderChat != nil {
		senderid = message.SenderChat.ID
		title = message.SenderChat.Title
	}
	store.UpdateUsername(senderid, title)
	for _, val := range config.Whitelists.Ranking {
		if val == chat.ID {
			message := ""
			if config.Ranking.Global {
				message += "水群排行榜\n"
				// message += "总消息数:" + strconv.FormatInt(store.GetGroupRecordedMessages(chat.ID, "global"), 10) + "\n"
				result := store.GetRanking("global")
				for i, val := range result {
					message += strconv.Itoa(i+1) + ". " + val.Name + ":" + strconv.FormatInt(val.Count, 10) + "次\n"
				}
				message += "\n"
			}
			if config.Ranking.Features.Cai {
				message += "卖菜排行榜\n"
				// message += "总消息数:" + strconv.FormatInt(store.GetGroupRecordedMessages(chat.ID, "cai"), 10) + "\n"
				result := store.GetRanking("cai")
				for i, val := range result {
					message += strconv.Itoa(i+1) + ". " + val.Name + ":" + strconv.FormatInt(val.Count, 10) + "次\n"
				}
				message += "\n"
			}
			if config.Ranking.Features.Xm {
				message += "羡慕排行榜\n"
				// message += "总消息数:" + strconv.FormatInt(store.GetGroupRecordedMessages(chat.ID, "xm"), 10) + "\n"
				result := store.GetRanking("xm")
				for i, val := range result {
					message += strconv.Itoa(i+1) + ". " + val.Name + ":" + strconv.FormatInt(val.Count, 10) + "次\n"
				}
			}
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chat.ID,
				Text:   message,
			})
		}
	}
}
