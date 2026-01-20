package main

import (
	"os"
	"log"
	
	"github.com/joho/godotenv"

	"github.com/LainIwakuras-father/ValentinkaBot/internal/adapter"
	"github.com/LainIwakuras-father/ValentinkaBot/internal/db"
	"github.com/LainIwakuras-father/ValentinkaBot/internal/handlers"
)

func main() {
	// load envirement
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Инициализация адаптеров
	bot, err := adapter.NewTelegramAdapter(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	db := db.NewDb()
	//DI 
	handler := handlers.NewHandler(bot, db)

	// Запуск бота
	updates := bot.ListenUpdates()
	log.Println(" Бот запущен и ожидает сообщения...")

	// 6. Обрабатываем обновления
	for update := range updates {
		// Обработка callback (нажатия на inline-кнопки)
		if update.CallbackQuery != nil {
			handler.HandleCallback(update.CallbackQuery)
			continue
		}

		// Обработка сообщений
		if update.Message == nil {
			continue
		}

		userID := update.Message.From.ID
		chatID := update.Message.Chat.ID

		// Обработка команд
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				handler.HandleStart(userID, chatID)
			case "help":
				handler.HandleHelp(chatID)
			case "myid":
				handler.HandleMyID(userID, chatID)
			}
			continue
		}

		// Обработка текстовых сообщений
		if update.Message.Text != "" {
			handler.HandleTextMessage(userID, chatID, update.Message.Text)
		}
	}
}