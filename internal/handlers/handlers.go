package handlers

import (
	"log"
	"strconv"

	"github.com/LainIwakuras-father/ValentinkaBot/internal/adapter"
	"github.com/LainIwakuras-father/ValentinkaBot/internal/storage"
)

// поменять на интерфейсы
type Handler struct {
	adapter *adapter.TelegramAdapter
	db      *storage.MemoryStorage
}

func NewHandler(adapter *adapter.TelegramAdapter, db *storage.MemoryStorage) *Handler {
	return &Handler{
		adapter: adapter,
		db:      db,
	}
}

// HandleStart обрабатывает /start
func (h *Handler) HandleStart(userID int64, chatID int64) {
	// Сбрасываем состояние пользователя

	text := `Привествую Друг!
Я бот сохраняющий и отправляющий твои сообщения кому либо!
Напиши Юзернейм или ID пользователя и текст который хочешь ему отправить
Формат:
9038487587 Привет, Друг!
	`
	if err := h.adapter.SendMessage(chatID, text); err != nil {
		log.Printf("Ошибка отправки /start: %v", err)
	}
}

// HandleTextMessage обрабатывает текстовые сообщения любые кроме команд
func (h *Handler) HandleTextMessage(userID int64, chatID int64, text string) {

	// Сохранить сообщение в памяти
	h.SaveMessage(userID, chatID, text)
	//Отправить ответ подтверждение
	totalmsg := strconv.Itoa(h.db.Count())
	confirmation := "Cообщение сохранено! Всего Сообщений:" + totalmsg + "\n\n"

	if err := h.adapter.SendMessage(chatID, confirmation); err != nil {
		log.Printf("Ошибка сохранения сообщения %v", err)
		return
	}

}

func (h *Handler) SaveMessage(userID int64, chatID int64, text string) {
	h.db.Save(userID, chatID, text)
}
