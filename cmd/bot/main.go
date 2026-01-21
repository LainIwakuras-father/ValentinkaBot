package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/LainIwakuras-father/ValentinkaBot/internal/adapter"
	"github.com/LainIwakuras-father/ValentinkaBot/internal/handlers"
	"github.com/LainIwakuras-father/ValentinkaBot/internal/storage"
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

	db := storage.NewMemoryStorage()
	//DI
	handler := handlers.NewHandler(bot, db)

	// Запуск бота
	updates := bot.ListenUpdates()
	log.Println(" Бот запущен и ожидает сообщения...")

	// 6. Обрабатываем обновления
	for update := range updates {

		// Обработка сообщений
		if update.Message == nil {
			continue
		}
		// сохранить chat_id пользователя (в памяти)
		userID := update.Message.From.ID
		chatID := update.Message.Chat.ID

		// Обработка команд
		if update.Message.IsCommand() {

			switch update.Message.Command() {
			case "start":
				handler.HandleStart(userID, chatID)
			default:
				// Можно добавить обработку неизвестных команд
				if err := bot.SendMessage(chatID, "Неизвестная команда. Используй /start"); err != nil {
					log.Printf("Ошибка отправки: %v", err)
				}
			}
			continue
		} // Обработка текстовых сообщений
        if update.Message.Text != "" {
            handler.HandleTextMessage(userID, chatID, update.Message.Text)
        }

	}
}
