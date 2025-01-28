package message_handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"queuingbot/domain"
	"queuingbot/internal/grpc"
	userModel "queuingbot/user"
)

var MainMarkup = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Получить список задач", "get_todo_list"),
	),
)

type MessageHandler struct {
	Users       map[int64]domain.User
	UserService *userModel.Service
}

func NewMessageHandler(service *userModel.Service) MessageHandler {
	return MessageHandler{
		Users:       make(map[int64]domain.User),
		UserService: service,
	}
}

func (h MessageHandler) Handle(m *grpc.QueuingMessage) tgbotapi.MessageConfig {
	switch m.Text {
	case "/start":
		{
			user, ok := h.Users[m.FromChatId]
			var msg tgbotapi.MessageConfig
			if !ok {
				user = domain.User{
					ChatId: m.FromChatId,
				}
				err := h.UserService.StoreUser(&user)

				if err != nil {
					return tgbotapi.MessageConfig{}
				}

				h.Users[m.FromChatId] = user

				msg = tgbotapi.NewMessage(user.ChatId, "Вы зарегистрированы! Начинайте пользоваться ботом)")
			} else {
				msg = tgbotapi.NewMessage(user.ChatId, "Чем помочь?")
			}
			msg.ReplyMarkup = MainMarkup
			return msg
		}
	}

	return tgbotapi.MessageConfig{}
}
