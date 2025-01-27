package bot_settings

import (
	"context"
	"telegram-bots-gateway/domain"
)

type BotSettingsRepository interface {
	All(ctx context.Context) ([]domain.BotSettings, error)
	GetByID(ctx context.Context, id uint) (domain.BotSettings, error)
	Update(ctx context.Context, ar *domain.BotSettings) error
	Store(ctx context.Context, a *domain.BotSettings) error
	Delete(ctx context.Context, id uint) error
}

type Service struct {
	botSettingsRepo BotSettingsRepository
}

func NewService(a BotSettingsRepository) *Service {
	return &Service{
		botSettingsRepo: a,
	}
}

func (s Service) GetBotSettings(ctx context.Context) ([]domain.BotSettings, error) {
	botSettings, err := s.botSettingsRepo.All(ctx)
	if err != nil {
		return botSettings, err
	}
	return botSettings, nil
}
