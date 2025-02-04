package domain

import "gorm.io/gorm"

type File struct {
	gorm.Model
}

func NewFile(id int) File {
	return File{Model: gorm.Model{ID: uint(id)}}
}
