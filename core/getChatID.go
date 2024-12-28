package core

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"strconv"
)

func ChatIDHandler(ctx context.Context, b *bot.Bot, update *models.Update, config YmlConfigurationData) {
	msg := "Chat Information\n"
	msg += "ChatID: " + strconv.FormatInt(update.Message.Chat.ID, 10) + "\n"
	msg += "SenderID: " + strconv.FormatInt(update.Message.From.ID, 10) + "\n"
	if config.AdminUID == update.Message.From.ID {
		msg += "You are the admin of this bot."
	} else if config.AdminUID == -1 {
		msg += "Admin notification feature is not enabled on this bot."
	} else {
		msg += "You are not the admin of this bot."
	}
	for _, val := range config.Whitelists.UnpinChannelPosts {
		if val == update.Message.Chat.ID {
			msg += "\nUnpin feature enabled."
		}
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   msg,
	})
}
