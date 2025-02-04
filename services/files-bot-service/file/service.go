package file

import "files-bot-service/internal/repository/postgresql"

type Service struct {
	FileRepository *postgresql.FileRepository
}

func NewService(fileRepository *postgresql.FileRepository) *Service {
	return &Service{FileRepository: fileRepository}
}
