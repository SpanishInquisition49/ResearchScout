package telegram

import (
	"fmt"
	"log"
	"strings"
	"unipi-research-crawler/internal/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func toBotStringHTML(c database.Call) string {
	msgText := fmt.Sprintf(
		"🎓 <b>%s</b>\n\n"+
			"• <b>Deadline:</b> %s\n",
		c.Title,
		c.Deadline,
	)
	return strings.ToValidUTF8(msgText, "")
}

func SendRichMessages(bot *tgbotapi.BotAPI, chatID int64, call database.Call) {
	// Option 1: Message with Inline Keyboard
	msg := tgbotapi.NewMessage(chatID, toBotStringHTML(call))
	msg.ParseMode = tgbotapi.ModeHTML

	// Add inline keyboard buttons
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("📖 Requirements", call.Requirements),
			tgbotapi.NewInlineKeyboardButtonURL("✍️ Apply Now", call.ApplyModule),
		),
	)
	msg.ReplyMarkup = keyboard

	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send message with keyboard: %v", err)
	}
}
