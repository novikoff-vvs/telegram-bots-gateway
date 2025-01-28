package postgresql

import (
	"context"
	"gorm.io/gorm"
	"queuingbot/domain"
)

type UserRepository struct {
	Conn      *gorm.DB
	BaseModel domain.User
}

func NewUserRepository(conn *gorm.DB) *UserRepository {
	return &UserRepository{Conn: conn, BaseModel: domain.User{}}
}

func (r UserRepository) getQuery() *gorm.DB {
	return r.Conn.Model(&r.BaseModel)
}

func (r UserRepository) Store(ctx context.Context, item *domain.User) error {
	if err := r.getQuery().Save(item).Error; err != nil {
		return err

	}
	return nil
}
