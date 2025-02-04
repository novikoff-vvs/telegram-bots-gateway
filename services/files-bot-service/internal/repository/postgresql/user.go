package postgresql

import (
	"files-bot-service/domain"
	"gorm.io/gorm"

	"context"
)

type UserRepository struct {
	Conn      *gorm.DB
	BaseModel *domain.User
}

func NewUserRepository(conn *gorm.DB) *UserRepository {
	return &UserRepository{conn, &domain.User{
		Model: gorm.Model{},
	}}
}
func (m *UserRepository) All(ctx context.Context) (res []domain.User, err error) {
	m.getQuery().Find(&res)
	return res, nil
}
func (m *UserRepository) getQuery() *gorm.DB {
	return m.Conn.Model(m.BaseModel)
}

func (m *UserRepository) GetByID(ctx context.Context, id uint) (res domain.User, err error) {
	m.getQuery().Where("id = ?", id).First(&res)
	return res, nil
}

func (m *UserRepository) Store(ctx context.Context, a *domain.User) (err error) {
	if m.getQuery().Create(a).Error != nil {
		return m.getQuery().Create(a).Error
	}
	return nil
}

func (m *UserRepository) Delete(ctx context.Context, id uint) (err error) {
	return
}
func (m *UserRepository) Update(ctx context.Context, ar *domain.User) (err error) {
	m.getQuery().Save(&ar)
	return nil
}

func (m *UserRepository) Query() *gorm.DB {
	return m.getQuery()
}
