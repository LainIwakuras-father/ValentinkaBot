package interfaces

type IBot interface {
	SendMessage(chatID int64, message string) error
}
