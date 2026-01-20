package handlers

import (
	"fmt"
	"log"
	"strconv"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/LainIwakuras-father/ValentinkaBot/internal/adapter"
	"github.com/LainIwakuras-father/ValentinkaBot/internal/db"
	"github.com/LainIwakuras-father/ValentinkaBot/internal/models"
)

type Handler struct {
	adapter *adapter.TelegramAdapter
	db      *db.Db
}

func NewHandler(adapter *adapter.TelegramAdapter, db *db.Db) *Handler {
	return &Handler{
		adapter: adapter,
		db:      db,
	}
}

// HandleStart обрабатывает /start
func (h *Handler) HandleStart(userID int64, chatID int64) {
	// Сбрасываем состояние пользователя
	h.db.ResetState(userID)

	text := `💝 Добро пожаловать в бот Валентинки!

Нажмите кнопку ниже, чтобы отправить анонимную валентинку 💌`

	// Создаем Inline клавиатуру
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💌 Создать валентинку", "create_valentine"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(" Как это работает?", "show_help"),
		),
	)

	if err := h.adapter.SendMessageWithKeyboard(chatID, text, keyboard); err != nil {
		log.Printf("Ошибка отправки /start: %v", err)
	}
}

// HandleHelp обрабатывает /help или кнопку "Как это работает?"
func (h *Handler) HandleHelp(chatID int64) {
	text := ` Как работает бот:

1️⃣ Нажмите "Создать валентинку"
2️⃣ Напишите текст сообщения
3️⃣ Укажите ID получателя
4️⃣ Подтвердите отправку

⚠️ Важно:
• Получатель должен был написать боту /start
• Для получения своего ID используйте /myid
• Сообщения анонимны`

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💌 Создать валентинку", "create_valentine"),
		),
	)

	if err := h.adapter.SendMessageWithKeyboard(chatID, text, keyboard); err != nil {
		log.Printf("Ошибка отправки помощи: %v", err)
	}
}

// HandleMyID обрабатывает /myid
func (h *Handler) HandleMyID(userID int64, chatID int64) {
	text := fmt.Sprintf(`🆔 Ваш ID: %d

Поделитесь им с друзьями, чтобы получать валентинки!`, userID)

	if err := h.adapter.SendMessage(chatID, text); err != nil {
		log.Printf("Ошибка отправки ID: %v", err)
	}
}

// HandleCreateValentine начинает создание валентинки
func (h *Handler) HandleCreateValentine(userID int64, chatID int64) {
	// Получаем состояние пользователя
	state := h.db.GetUserState(userID)
	
	// Меняем шаг на "ввод сообщения"
	state.CurrentStep = "enter_message"
	state.Valentine = &models.ValentineMsg{
		FromID: userID,
	}
	
	// Сохраняем состояние
	h.db.SetUserState(userID, state)
	text := `💌 Напишите текст валентинки:

Это может быть признание, комплимент или что-то особенное!`

	// Клавиатура с кнопкой отмены
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(" Отменить", "cancel"),
		),
	)

	if err := h.adapter.SendMessageWithKeyboard(chatID, text, keyboard); err != nil {
		log.Printf("Ошибка создания валентинки: %v", err)
	}
}

// HandleTextMessage обрабатывает текстовые сообщения
// Паттерн: State Machine (обработка зависит от состояния)
func (h *Handler) HandleTextMessage(userID int64, chatID int64, text string) {
	state := h.db.GetUserState(userID)

	switch state.CurrentStep {
	case "enter_message":
		h.handleEnterMessage(userID, chatID, text, state)
	case "enter_recipient":
		h.handleEnterRecipient(userID, chatID, text, state)
	default:
		h.adapter.SendMessage(chatID, "Используйте /start для начала работы")
	}
}

// handleEnterMessage обрабатывает ввод текста валентинки
func (h *Handler) handleEnterMessage(userID int64, chatID int64, text string, state *models.UserState) {
	// Проверка длины
	if len(text) > 500 {
		h.adapter.SendMessage(chatID, " Текст слишком длинный! Максимум 500 символов.")
		return
	}

	// Сохраняем текст
	state.Valentine.Text = text
	state.CurrentStep = "enter_recipient"
	h.db.SetUserState(userID, state)

	// Запрашиваем ID получателя
	recipientText := ` Текст сохранен!

Теперь введите ID получателя:

Чтобы узнать свой ID, человек должен написать боту /myid`

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(" Отменить", "cancel"),
		),
	)

	if err := h.adapter.SendMessageWithKeyboard(chatID, recipientText, keyboard); err != nil {
		log.Printf("Ошибка запроса ID: %v", err)
	}
}

// handleEnterRecipient обрабатывает ввод ID получателя
func (h *Handler) handleEnterRecipient(userID int64, chatID int64, text string, state *models.UserState) {
	// Парсим ID
	recipientID, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		h.adapter.SendMessage(chatID, "❌ Неверный формат ID! Введите число.")
		return
	}

	// Проверка на самоотправку
	if recipientID == userID {
		h.adapter.SendMessage(chatID, "❌ Нельзя отправить валентинку самому себе!")
		return
	}

	// Сохраняем ID получателя
	state.Valentine.ToUser = recipientID
	h.db.SetUserState(userID, state)

	// Показываем предпросмотр с кнопками подтверждения
	h.showPreview(chatID, state.Valentine)
}

// showPreview показывает предпросмотр валентинки
func (h *Handler) showPreview(chatID int64, valentine *models.ValentineMsg) {
	previewText := fmt.Sprintf(`📋 Предпросмотр валентинки:

💌 Текст: %s

👤 Получатель: ID %d

⚠️ Сообщение будет отправлено анонимно!`, valentine.Text, valentine.ToUser)

	// Inline клавиатура для подтверждения
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Отправить", "confirm_send"),
			tgbotapi.NewInlineKeyboardButtonData("❌ Отменить", "cancel"),
		),
	)

	if err := h.adapter.SendMessageWithKeyboard(chatID, previewText, keyboard); err != nil {
		log.Printf("Ошибка показа предпросмотра: %v", err)
	}
}

// HandleCallback обрабатывает нажатия на inline-кнопки
// Паттерн: Command (каждая кнопка - команда)
func (h *Handler) HandleCallback(callbackQuery *tgbotapi.CallbackQuery) {
	userID := callbackQuery.From.ID
	chatID := callbackQuery.Message.Chat.ID
	data := callbackQuery.Data

	// Убираем "часики загрузки"
	h.adapter.AnswerCallbackQuery(callbackQuery.ID, "")

	// Обрабатываем нажатие
	switch data {
	case "create_valentine":
		h.HandleCreateValentine(userID, chatID)
	
	case "show_help":
		h.HandleHelp(chatID)
	
	case "cancel":
		h.handleCancel(userID, chatID)
	
	case "confirm_send":
		h.handleConfirmSend(userID, chatID)
	
	default:
		log.Printf("Неизвестный callback: %s", data)
	}
}

// handleCancel отменяет создание валентинки
func (h *Handler) handleCancel(userID int64, chatID int64) {
	h.db.ResetState(userID)
	
	text := " Создание валентинки отменено."
	
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💌 Создать новую", "create_valentine"),
		),
	)
	
	h.adapter.SendMessageWithKeyboard(chatID, text, keyboard)
}

// handleConfirmSend отправляет валентинку
func (h *Handler) handleConfirmSend(userID int64, chatID int64) {
	state := h.db.GetUserState(userID)

	if state.Valentine == nil {
		h.adapter.SendMessage(chatID, "❌ Ошибка: валентинка не найдена")
		return
	}

	// Отправляем валентинку получателю
	valentineText := fmt.Sprintf(`💝 Вы получили анонимную валентинку!

%s

🎭 От: Тайный поклонник`, state.Valentine.Text)
	err := h.adapter.SendMessage(state.Valentine.ToUser, valentineText)
	
	if err != nil {
		// Ошибка отправки
		errorText := ` Не удалось отправить валентинку.

Возможные причины:
• Получатель не начал диалог с ботом (/start)
• Неверный ID получателя

Попросите получателя написать боту /start`
		h.adapter.SendMessage(chatID, errorText)
		log.Printf("Ошибка отправки валентинки: %v", err)
	} else {
		// Успех!
		successText := " Валентинка отправлена! 💝"
		
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("💌 Отправить еще", "create_valentine"),
			),
		)
		
		h.adapter.SendMessageWithKeyboard(chatID, successText, keyboard)
	}
	
	// Очищаем состояние
	h.db.ResetState(userID)
}