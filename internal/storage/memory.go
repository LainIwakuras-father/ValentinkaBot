package storage

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/LainIwakuras-father/ValentinkaBot/internal/models"
)

// MemoryStorage - временное хранилище в памяти
type MemoryStorage struct {
	messages map[string]*models.Message
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		messages: make(map[string]*models.Message),
	}
}

func (m *MemoryStorage) Save(userID int64, chatID int64, text string) (string, error) {
	// Генерируем уникальный ID
	generatedID := generateID()

	msg := &models.Message{
		ChatID: chatID,
		UserID: userID,
		Text:   text,
	}

	m.messages[generatedID] = msg
	return generatedID, nil

}

func (m *MemoryStorage) Count() int {
	return len(m.messages)
}

// generateID - генерирует уникальный ID
func generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
