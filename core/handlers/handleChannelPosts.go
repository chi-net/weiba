package handlers

import (
	"context"
	"github.com/chi-net/weiba/core/types"
	"github.com/chi-net/weiba/core/utils"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func HandleChannelPosts(ctx context.Context, b *bot.Bot, update *models.Update, config types.YmlConfigurationData) {
	msg := "[Debug] Detected channel posted in chat and we unpinned it for you.\n"
	_, err := b.UnpinChatMessage(ctx, &bot.UnpinChatMessageParams{
		ChatID:    update.Message.Chat.ID,
		MessageID: update.Message.ID,
	})
	if err != nil {
		msg += "It returned an error:" + err.Error()
	}

	utils.SendDebugMessage(msg, ctx, b, config)
}
