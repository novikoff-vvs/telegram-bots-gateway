package postgresql

import (
	"files-bot-service/domain"
	"gorm.io/gorm"

	"context"
)

type FileRepository struct {
	Conn      *gorm.DB
	BaseModel *domain.File
}

func NewFileRepository(conn *gorm.DB) *FileRepository {
	return &FileRepository{conn, &domain.File{
		Model: gorm.Model{},
	}}
}
func (m *FileRepository) All(ctx context.Context) (res []domain.File, err error) {
	m.getQuery().Find(&res)
	return res, nil
}
func (m *FileRepository) getQuery() *gorm.DB {
	return m.Conn.Model(m.BaseModel)
}

func (m *FileRepository) GetByID(ctx context.Context, id uint) (res domain.File, err error) {
	m.getQuery().Where("id = ?", id).First(&res)
	return res, nil
}

func (m *FileRepository) Store(ctx context.Context, a *domain.File) (err error) {
	if err = m.getQuery().Create(a).Error; err != nil {
		return m.getQuery().Create(a).Error
	}
	return nil
}

func (m *FileRepository) Delete(ctx context.Context, id uint) (err error) {
	return
}
func (m *FileRepository) Update(ctx context.Context, ar *domain.File) (err error) {
	m.getQuery().Save(&ar)
	return nil
}

func (m *FileRepository) Query() *gorm.DB {
	return m.getQuery()
}
