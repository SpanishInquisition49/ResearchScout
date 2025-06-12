# Research Scout 🔬

A Telegram bot that automatically crawls and notifies users about new research opportunities and calls for papers. Stay updated with the latest academic opportunities without manually checking multiple sources.

⚠️ Disclaimer: This is an unofficial, independent project and is not affiliated with or endorsed by the University of Pisa (Unipi) or any other academic institution. This bot is developed and maintained by me.

## 📱 How to Use

The bot is hosted and ready to use! Simply start a conversation with:

**[@researchscout_bot](https://t.me/researchscout_bot)**

### Commands

- `/start` - Subscribe to research updates
- `/stop` - Unsubscribe from notifications

Once subscribed, you'll automatically receive messages about new research calls.

## 🏗️ Architecture

This bot consists of several components running concurrently:

- **Web Scraper**: Continuously monitors research websites for new opportunities
- **Telegram Bot Handler**: Manages user subscriptions and commands
- **Notification System**: Delivers personalized updates to subscribed users
- **Database Manager**: Stores user data and research calls with automatic cleanup

## 🔧 Technical Details

### Built With

- **Language**: Go
- **Database**: SQLite
- **Bot Framework**: go-telegram-bot-api
- **Dependencies**:
  - `github.com/mattn/go-sqlite3` - SQLite driver
  - `github.com/go-telegram-bot-api/telegram-bot-api/v5` - Telegram Bot API
  - `github.com/joho/godotenv` - Environment configuration

### Key Features

- Concurrent processing with goroutines
- Automatic database cleanup (runs daily)
- Rich message formatting for research calls
- User subscription management
- Error handling and logging
- Configurable scraping intervals

### Processing Flow

1. **Crawler**: Runs every 30 minutes to scrape new research opportunities
2. **Notifications**: Checks every minute for new calls to send to users
3. **Cleanup**: Removes old data daily to maintain database performance

## 🚀 Self-Hosting

If you want to run your own instance:

### Prerequisites

- Go 1.19 or higher
- SQLite3
- sqlc
- sql-migrate

### Setup

1. Clone the repository
2. Create a `.env` file with:

   ```
   TELEGRAM_TOKEN=your_bot_token_here
   DATABASE_URL=path/to/your/database.db
   ```

3. Install dependencies:

   ```bash
   go mod tidy
   ```

4. Create the database and the database code:

   ```bash
   touch database.db
   sql-migrate up
   sqlc generate
   ```

5. Run the bot:

   ```bash
   go run main.go
   ```

### Project Structure

```
unipi-research-crawler/
├── main.go
├── internal/
│   ├── database/     # Database operations
│   ├── scraper/      # Web scraping logic
│   └── telegram/     # Telegram message formatting
└── .env             # Environment variables
```

## 📊 Database Schema

The bot uses SQLite with tables for:

- **Users**: Chat IDs and subscription status
- **Calls**: Research opportunities and metadata
- **Notifications**: Tracking sent messages to prevent duplicates

## 🤝 Contributing

This is a personal project, but feedback and suggestions are welcome! Feel free to open issues or reach out if you have ideas for improvements.

## 📄 License

This project is for educational and personal use. Please respect the terms of service of any websites being scraped.

---

**Bot**: [@researchscout_bot](https://t.me/researchscout_bot)  
**Status**: Active and maintained
