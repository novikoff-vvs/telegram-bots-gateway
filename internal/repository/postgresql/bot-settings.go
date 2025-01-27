package postgresql

import (
	"gorm.io/gorm"
	"telegram-bots-gateway/domain"

	"context"
)

type BotSettingsRepository struct {
	Conn      *gorm.DB
	BaseModel *domain.BotSettings
}

func NewBotSettingsRepository(conn *gorm.DB) *BotSettingsRepository {
	return &BotSettingsRepository{conn, &domain.BotSettings{
		Model:      gorm.Model{},
		ApiKey:     "",
		AdminID:    "",
		ServiceUrl: "",
	}}
}
func (m *BotSettingsRepository) All(ctx context.Context) (res []domain.BotSettings, err error) {
	m.getQuery().Find(&res)
	return res, nil
}
func (m *BotSettingsRepository) getQuery() *gorm.DB {
	return m.Conn.Model(m.BaseModel)
}

func (m *BotSettingsRepository) GetByID(ctx context.Context, id uint) (res domain.BotSettings, err error) {
	m.getQuery().Where("id = ?", id).First(&res)
	return res, nil
}

func (m *BotSettingsRepository) Store(ctx context.Context, a *domain.BotSettings) (err error) {
	if m.getQuery().Create(a).Error != nil {
		return m.getQuery().Create(a).Error
	}
	return nil
}

func (m *BotSettingsRepository) Delete(ctx context.Context, id uint) (err error) {
	return
}
func (m *BotSettingsRepository) Update(ctx context.Context, ar *domain.BotSettings) (err error) {
	return
}
