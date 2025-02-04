package handlers

import (
	"context"
	"files-bot-service/domain"
	"files-bot-service/file"
	"files-bot-service/user"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type MessageHandler struct {
	UserService *user.Service
	FileService *file.Service
}

func NewMessageHandler(userService *user.Service, fileService *file.Service) *MessageHandler {
	return &MessageHandler{
		UserService: userService,
		FileService: fileService,
	}
}

func (handler *MessageHandler) HandleMessage(message tgbotapi.Message) (msg tgbotapi.Chattable, err error) {
	var model domain.User
	model, err = handler.UserService.Repository.GetByID(context.Background(), uint(message.Chat.ID)) //TODO поправить на сервис
	if err != nil {
		return tgbotapi.MessageConfig{}, err
	}

	switch message.Text {
	case "/start":
		{
			if model.ID == 0 {
				model, err = handler.UserService.CreateUser(message.Chat.ID)
				if err != nil {
					return tgbotapi.MessageConfig{}, err
				}
				return tgbotapi.NewMessage(message.Chat.ID, "Вы успешно зарегистрированы!"), nil
			}

			return tgbotapi.NewMessage(message.Chat.ID, "Чем могу помочь?"), nil
		}
	case "/add":
		{
			if model.ID == 0 {
				log.Println("Пользователь не найден!")
			}
			model.Terminator = domain.WAIT_INPUT
			model, err = handler.UserService.Update(model)
			if err != nil {
				return tgbotapi.MessageConfig{}, err
			}

			return tgbotapi.NewMessage(int64(model.ID), "Пришлите файл, который вы хотите сохранить"), nil
		}
	case "/getID": //TODO сделать автосетап по команде
		{
			return tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Chat ID: %d", message.Chat.ID)), nil
		}
	default:
		{
			switch model.Terminator {
			case domain.WAIT_INPUT:
				{
					f := domain.NewFile(message.MessageID)
					err = handler.FileService.FileRepository.Store(context.Background(), &f)
					if err != nil {
						return tgbotapi.MessageConfig{}, err
					}
					return tgbotapi.NewForward(-4655598408, int64(model.ID), int(f.ID)), nil
				}
			}
		}
		return tgbotapi.MessageConfig{}, nil
	}
}
