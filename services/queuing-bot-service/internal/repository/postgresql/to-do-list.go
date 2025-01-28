package postgresql

import (
	"context"
	"gorm.io/gorm"
	"queuingbot/domain"
)

type ToDoListItemRepository struct {
	Conn      *gorm.DB
	BaseModel domain.ToDoListItem
}

func (r ToDoListItemRepository) getQuery() *gorm.DB {
	return r.Conn.Model(&r.BaseModel)
}

func (r ToDoListItemRepository) Store(ctx context.Context, item domain.ToDoListItem) error {
	if err := r.getQuery().Save(item).Error; err != nil {
		return err

	}
	return nil
}
