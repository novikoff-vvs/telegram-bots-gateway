package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
)

var Bot *tgbotapi.BotAPI
var BotActivity = false
var botStarted = make(chan bool)
var QuitChan = make(chan os.Signal, 1)
var BotExited = make(chan bool)
var usersWaitRegister []int64

var ServiceMap = map[uint][]tgbotapi.InlineKeyboardButton{
	1: tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Запросить Wireguard", "get_wg"),
	),
	2: tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Запросить PPTP", "get_pptp"),
	),
}

var mainMarkup = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Получить конфиг VPN", "get_config"),
		tgbotapi.NewInlineKeyboardButtonData("Запросить услугу", "get_vpn"),
	),
)

var registrationMarkup = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Уже есть свой идентификатор", "register"),
		tgbotapi.NewInlineKeyboardButtonData("Получить VPN", "get_vpn"),
	),
)

type UpdateHandler struct{}

func (h UpdateHandler) HandleMessage(message string) {
	var msg tgbotapi.MessageConfig

	switch message {
	case "/start":
		{
			if user.ID != 0 {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "И снова дратути!")
				msg.ReplyMarkup = mainMarkup
			} else {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Чем могу помочь?")
				msg.ReplyMarkup = registrationMarkup
			}
			_, err := Bot.Send(msg)
			if err != nil {
				return
			}
		}
	case "/adduser":
		{
			if update.Message.Chat.ID != settings.AdminUserId {
				continue
			} else {
				name := RandomString(6)
				CreateNewCert(name)
				var vpnConfig VpnConfig
				var vpnUser VpnUser
				//TODO вынести нахуй в отдельную функцию
				vpnConfig.VpnConfigName = name
				vpnConfig.ConfigPath = "uploads/configs/wg0-client-" + name + ".conf"
				vpnConfig.QrCodeImagePath, _ = GenerateQRCode("uploads/configs/wg0-client-" + name + ".conf")
				vpnConfig.Save()

				vpnUser.VpnConfig = &vpnConfig
				vpnUser.Comment = "Создано из бота"
				vpnUser.VerificationCode = RandomString(6)

				vpnUser.Save()
				msg = tgbotapi.NewMessage(settings.AdminUserId, "Новый пользователь добавлен.\nИмя:"+name+"\nИдентификатор для регистрации:")
				Bot.Send(msg)
				msg = tgbotapi.NewMessage(settings.AdminUserId, vpnUser.VerificationCode)
			}
			_, err := Bot.Send(msg)
			if err != nil {
				return
			}
		}
	case "/notify":
		{

		}

	default:
		{
			isWait, index := userIsWaiting(update.Message.Chat.ID, usersWaitRegister)
			if isWait {
				var registeredUser VpnUser
				DB.Model(&VpnUser{}).Where("verification_code = ?", update.Message.Text).First(&registeredUser)
				if registeredUser.ID == 0 {
					Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Не правильный код! Проверьте и попробуйте ввести еще раз"))
					continue
				}
				registeredUser.ChatId = update.Message.Chat.ID
				registeredUser.Username = update.Message.Chat.UserName
				registeredUser.Save()

				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Регистрация прошла успешно!")
				msg.ReplyMarkup = mainMarkup
				Bot.Send(msg)

				usersWaitRegister = RemoveIndex(usersWaitRegister, index)

			}
		}
	}
}
