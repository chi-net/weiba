package core

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func HandleChannelPosts(ctx context.Context, b *bot.Bot, update *models.Update, config YmlConfigurationData) {
	msg := "[Debug] Detected channel posted in chat and we unpinned it for you.\n"
	_, err := b.UnpinChatMessage(ctx, &bot.UnpinChatMessageParams{
		ChatID:    update.Message.Chat.ID,
		MessageID: update.Message.ID,
	})
	if err != nil {
		msg += "It returned an error:" + err.Error()
	}

	SendDebugMessage(msg, ctx, b, config)
}
