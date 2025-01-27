package bot

import (
	"context"
	botSettings "telegram-bots-gateway/bot-settings"
	"telegram-bots-gateway/domain"
	"time"
)

type Service struct {
	BotSettingsService botSettings.Service
	errors             []error
}

func NewService(botSettingsService botSettings.Service) *Service {
	return &Service{botSettingsService, make([]error, 0)}
}

func (s Service) GetErrors() []error {
	return s.errors
}

func (s Service) GetBots() (res []domain.Bot, err error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(100))
	settings, err := s.BotSettingsService.GetBotSettings(ctx)
	if err != nil {
		return nil, err
	}

	for _, setting := range settings {
		bot, err := domain.NewBot(setting)

		if err != nil {
			s.errors = append(s.errors, err)
		}

		res = append(res, bot)
	}

	return res, nil
}
