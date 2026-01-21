package models

type Message struct {
	ChatID int64  `json:"chat_id"`
	UserID int64  `json:"from_id"`
	Text   string `json:"text"`
}

// Глобальное хранилище в памяти (потоконебезопасно, для демо)
var MessageStore = make(map[string]Message)
var messageCounter = 0

// // состояния для валидации и сохранения
// type UserState struct {
// 	CurrentStep string
// 	Valentine *Message
// }

// func NewUserState(userID int64) *UserState {
// 	return &UserState{
// 		CurrentStep: "idle",
// 		Valentine:   &Message{
// 			FromID: userID,
// 		},
// 	}
// }
