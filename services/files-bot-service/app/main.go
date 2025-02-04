package main

import (
	"context"
	"encoding/json"
	"files-bot-service/domain"
	"files-bot-service/file"
	pb "files-bot-service/internal/grpc"
	"files-bot-service/internal/handlers"
	pgsql "files-bot-service/internal/repository/postgresql"
	"files-bot-service/user"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func initDatabase() (*gorm.DB, error) {
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPass := os.Getenv("DATABASE_PASS")
	dbName := os.Getenv("DATABASE_NAME")

	dsn := "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Moscow"
	connection := fmt.Sprintf(dsn, dbHost, dbUser, dbPass, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(connection), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&domain.User{}, &domain.File{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	db, err := initDatabase()
	if err != nil {
		return
	}
	// Создаём слушатель для порта
	port := 50053
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Инициализация gRPC сервера
	grpcServer := grpc.NewServer()
	pb.RegisterFilesBotServer(
		grpcServer,
		newServer(
			handlers.NewMessageHandler(
				user.NewUserService(
					pgsql.NewUserRepository(
						db,
					),
				),
				file.NewService(
					pgsql.NewFileRepository(
						db,
					),
				),
			),
		),
	)

	// Запуск сервера
	log.Printf("Server is running at localhost:%d", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func newServer(handler *handlers.MessageHandler) *queueingBotServer {
	return &queueingBotServer{
		MessageHandler: handler,
	}
}

type queueingBotServer struct {
	pb.UnimplementedFilesBotServer
	MessageHandler *handlers.MessageHandler
}

func (s *queueingBotServer) Handle(ctx context.Context, message *pb.FilesBotMessage) (*pb.FilesBotResult, error) {
	var msg tgbotapi.Message
	err := json.Unmarshal([]byte(message.Json), &msg)
	if err != nil {
		return nil, err
	}

	response, err := s.MessageHandler.HandleMessage(msg)
	if err != nil {
		return nil, err
	}

	jsonString, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	switch response.(type) {
	case tgbotapi.ForwardConfig:
		{
			return &pb.FilesBotResult{
				ForwardMessageJson: string(jsonString),
			}, nil
		}
	default:
		{
			return &pb.FilesBotResult{
				NewMessageJson: string(jsonString),
			}, nil
		}
	}
}
