package utils

import (
	"context"
	"github.com/chi-net/weiba/core/types"
	"github.com/go-telegram/bot"
)

func SendDebugMessage(content string, ctx context.Context, b *bot.Bot, config types.YmlConfigurationData) {
	if config.Features.Debug {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: config.AdminUID,
			Text:   content,
		})
	}
}
