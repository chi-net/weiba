package handlers

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"regexp"
	"strconv"
	"strings"
)

func HandleTietie(ctx context.Context, b *bot.Bot, update *models.Update) {
	message := update.Message
	chat := message.Chat
	sender := message.From
	replytitle := "自己"
	replyid := sender.ID
	if message.ReplyToMessage != nil {
		replytitle = message.ReplyToMessage.From.FirstName + " " + message.ReplyToMessage.From.LastName
		replyid = message.ReplyToMessage.From.ID
		if message.ReplyToMessage.SenderChat != nil {
			replytitle = message.ReplyToMessage.SenderChat.Title
			replyid = message.ReplyToMessage.SenderChat.ID
		}
	}

	// group tietie feature
	// example:
	// user1: hi!
	// user2: /贴
	// bot: 'user1' 贴了 'user2'!
	// fmt.Println(message.Text)
	if message.Text[0] == '/' {
		pattern := `^[a-zA-Z0-9/\\@]+$`
		matched, _ := regexp.MatchString(pattern, message.Text)
		if matched && message.Text[1] != '/' {
			return
		}
		i := 0
		for message.Text[i] == '/' {
			i += 1
		}
		receive := strings.SplitN(message.Text[i:], " ", 2)
		msg := "[" + bot.EscapeMarkdown(sender.FirstName+" "+sender.LastName) + "](tg://user?id=" + strconv.FormatInt(sender.ID, 10) + ") "
		matched, _ = regexp.MatchString(pattern, receive[0])
		if matched && message.Text[1] != '/' {
			return
		}
		if len(receive) == 2 {
			msg += receive[0] + " "
		} else if len(receive) == 1 {
			msg += receive[0] + "了" + " "
		}
		if message.ReplyToMessage != nil {
			msg += "[" + bot.EscapeMarkdown(replytitle) + "](tg://user?id=" + strconv.FormatInt(replyid, 10) + ")"
			if len(receive) == 2 {
				msg += receive[1] + "\\!"
			} else {
				msg += "\\!"
			}
		} else {
			msg += "自己\\!"
		}
		// fmt.Println(msg)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chat.ID,
			Text:      msg,
			ParseMode: models.ParseModeMarkdown,
			ReplyParameters: &models.ReplyParameters{
				MessageID: message.ID,
				ChatID:    chat.ID,
			},
		})
	} else if message.Text[0] == '\\' {
		i := 0
		for message.Text[i] == '\\' {
			i += 1
		}
		receive := strings.SplitN(message.Text[i:], " ", 2)
		msg := "[" + bot.EscapeMarkdown(sender.FirstName+" "+sender.LastName) + "](tg://user?id=" + strconv.FormatInt(sender.ID, 10) + ") 被 "
		msg += "[" + bot.EscapeMarkdown(replytitle) + "](tg://user?id=" + strconv.FormatInt(replyid, 10) + ")"
		msg += receive[0] + "了" + " "

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chat.ID,
			Text:      msg,
			ParseMode: models.ParseModeMarkdown,
			ReplyParameters: &models.ReplyParameters{
				MessageID: message.ID,
				ChatID:    chat.ID,
			},
		})
	}
}
