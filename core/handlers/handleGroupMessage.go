package handlers

import (
	"context"
	"fmt"
	"github.com/chi-net/weiba/core/store"
	"github.com/chi-net/weiba/core/types"
	"github.com/chi-net/weiba/core/utils"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"strings"
)

func HandleGroupMessage(ctx context.Context, b *bot.Bot, update *models.Update, config types.YmlConfigurationData) {
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
			if config.Ranking.Global {
				store.RecordRanking(senderid, chat.ID, "global")
			}
			if config.Ranking.Features.Cai {
				if message.Text != "" && strings.Contains(message.Text, "我菜") {
					store.RecordRanking(senderid, chat.ID, "cai")
				} else if message.Text != "" && utils.IsOnlyRepeats(message.Text, "您") {
					store.RecordRanking(senderid, chat.ID, "cai")
				}
				if message.Sticker != nil {
					if message.Sticker.FileUniqueID == "AgAD6QkAAqSCAAFX" || // yufox trashbin
						message.Sticker.FileUniqueID == "AgAD8gsAApMaeVc" || // suzume trashbin
						message.Sticker.FileUniqueID == "AgADtAADuDulNA" || // wo
						message.Sticker.FileUniqueID == "AgAD6AUAAgGeUVY" { // xiaoxinmiao wocai
						store.RecordRanking(senderid, chat.ID, "cai")
					}
				}
			}
			if config.Ranking.Features.Xm {
				if message.Text != "" && strings.Contains(strings.ToLower(message.Text), "xm") {
					store.RecordRanking(senderid, chat.ID, "xm")
				}
				if message.Sticker != nil {
					if message.Sticker.FileUniqueID == "AgADhhcAAs1rgFU" || // suzume xmsl
						message.Sticker.FileUniqueID == "AgADswADuDulNA" { // nin
						store.RecordRanking(senderid, chat.ID, "xm")
					}
				}
			}
			if message.Sticker != nil {
				fmt.Println(message.Sticker.FileUniqueID)
			}
		}
	}
}
