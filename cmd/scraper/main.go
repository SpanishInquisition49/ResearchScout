package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"
	"unipi-research-crawler/internal/database"
	"unipi-research-crawler/internal/scraper"
	"unipi-research-crawler/internal/telegram"

	_ "github.com/mattn/go-sqlite3"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting Unipi Research Crawler...")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	botToken, validToken := os.LookupEnv("TELEGRAM_TOKEN")
	if !validToken {
		log.Fatalln("Invalid TELEGRAM_TOKEN")
	}
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	conn, err := sql.Open("sqlite3", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	// Spawn the bot goroutine
	requestsDB := database.New(conn)
	ctx := context.Background()
	log.Println("Connected to the database successfully.")
	// spawn a long lived goroutine to handle the bot
	go handleBot(bot, requestsDB)
	go handleCrawler(requestsDB)
	// Infinite loop to keep the service running
	for {
		users, err := requestsDB.GetUsers(ctx)
		if err != nil {
			log.Fatalf("Failed to get users: %v", err)
		}
		for _, user := range users {
			// Get all the calls that the user has not yet received
			calls, err := requestsDB.GetCallsToSend(context.Background(), user.ChatID)
			if err != nil {
				log.Printf("Failed to get calls for user %d: %v", user.ChatID, err)
				continue
			}
			// Send all the calls to the user
			for _, call := range calls {
				telegram.SendRichMessages(bot, user.ChatID, call)
				param := database.SendCallToUserParams{
					UserChatID: user.ChatID,
					CallID:     call.ID,
				}
				if _, err := requestsDB.SendCallToUser(ctx, param); err != nil {
					log.Printf("Failed to add call %d for user %d: %v", call.ID, user.ChatID, err)
				}
			}
		}
		time.Sleep(1 * time.Minute) // Wait for 30 minutes before the next scrape
	}
}

func handleCrawler(db *database.Queries) {
	for {
		log.Println("Starting crawler...")
		scraper.Scrape(db)
		log.Println("Crawling completed.")
		time.Sleep(30 * time.Minute) // Wait for 30 minutes before the next scrape
	}
}

func handleBot(bot *tgbotapi.BotAPI, db *database.Queries) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	ctx := context.Background()

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		if update.Message.IsCommand() {
			chatID := update.Message.Chat.ID
			switch update.Message.Command() {
			case "start":
				param := database.CreateUserParams{
					ChatID:           chatID,
					FirstInteraction: time.Now().Format(time.RFC3339),
				}
				if _, err := db.CreateUser(ctx, param); err != nil {
					log.Printf("Failed to create user: %v", err)
					continue
				}
				reply := tgbotapi.NewMessage(chatID, "👋 Welcome to Research Scout! You'll now receive updates.")
				if _, err := bot.Send(reply); err != nil {
					log.Printf("Could not send /start response to %d: %v", chatID, err)
				}
				log.Printf("User '%s' started the bot.", update.Message.Chat.FirstName)
			case "stop":
				db.DeactivateUser(ctx, chatID)
				reply := tgbotapi.NewMessage(chatID, "👋 You have been unsubscribed from Research Scout. You will no longer receive updates.")
				if _, err := bot.Send(reply); err != nil {
					log.Printf("Could not send /stop response to %d: %v", chatID, err)
				}
				log.Printf("User '%s' stopped the bot.", update.Message.Chat.FirstName)
			}
		}
	}
	log.Println("Bot handler stopped.")
}
