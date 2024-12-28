package core

import (
	"context"
	"github.com/go-telegram/bot"
)

func SendDebugMessage(content string, ctx context.Context, b *bot.Bot, config YmlConfigurationData) {
	if config.Features.Debug {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: config.AdminUID,
			Text:   content,
		})
	}
}
