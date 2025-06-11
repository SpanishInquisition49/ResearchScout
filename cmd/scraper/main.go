package main

import (
	"log"
	"os"
	"strconv"
	"unipi-research-crawler/internal/scraper"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	botToken, validToken := os.LookupEnv("TELEGRAM_TOKEN")
	chatId, validChatId := os.LookupEnv("CHAT_ID")
	if !validToken || !validChatId {
		log.Fatal("TELEGRAM_TOKEN and CHAT_ID environment variables must be set.")
	}
	authorizedChatId, err := strconv.ParseInt(chatId, 10, 64)
	if err != nil {
		log.Fatalln("Invalid CHAT_ID")
	}
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	log.Println("Starting crawler...")
	calls := scraper.Scrape()
	for _, call := range *calls {
		sendRichMessages(bot, authorizedChatId, call)
	}
	log.Println("Crawling completed.")
}

func sendRichMessages(bot *tgbotapi.BotAPI, chatID int64, call scraper.Call) {
	// Option 1: Message with Inline Keyboard
	msg := tgbotapi.NewMessage(chatID, call.ToBotStringHTML())
	msg.ParseMode = tgbotapi.ModeHTML

	// Add inline keyboard buttons
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("📖 Requirements", call.CallUrl),
			tgbotapi.NewInlineKeyboardButtonURL("✍️ Apply Now", call.ModuleUrl),
		),
	)
	msg.ReplyMarkup = keyboard

	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send message with keyboard: %v", err)
	}
}
