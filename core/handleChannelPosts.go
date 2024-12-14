package core

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func HandleChannelPosts(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.UnpinChatMessage(ctx, &bot.UnpinChatMessageParams{
		ChatID:    update.Message.Chat.ID,
		MessageID: update.Message.ID,
	})
}
