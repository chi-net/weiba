package handlers

import (
	"context"
	"github.com/chi-net/weiba/core/store"
	"github.com/chi-net/weiba/core/types"
	"github.com/chi-net/weiba/core/utils"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func GetRankingHandler(ctx context.Context, b *bot.Bot, update *models.Update, config types.YmlConfigurationData, count int) {
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
				result, count := store.GetRanking("global", chat.ID, count)
				message += utils.GetRankingMessage(result, count)
			}
			if config.Ranking.Features.Cai {
				message += "卖菜排行榜\n"
				result, count := store.GetRanking("cai", chat.ID, count)
				message += utils.GetRankingMessage(result, count)
			}
			if config.Ranking.Features.Xm {
				message += "羡慕排行榜\n"
				result, count := store.GetRanking("xm", chat.ID, count)
				message += utils.GetRankingMessage(result, count)
			}
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    chat.ID,
				Text:      message,
				ParseMode: models.ParseModeMarkdown,
			})
		}
	}
}
