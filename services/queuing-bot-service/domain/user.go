package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ChatId int64 `json:"chat_id"`
}
