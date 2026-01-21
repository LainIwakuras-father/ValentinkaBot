package service
import(
	 "github.com/LainIwakuras-father/ValentinkaBot/internal/interfaces"
)

type MessageService struct{
	bot 	interfaces.IBot
	storage interfaces.IStorage
}

func NewMessageService(bot interfaces.IBot, storage interfaces.IStorage) *MessageService {
    return &MessageService{
        bot:     bot,
        storage: storage,
    }
}
