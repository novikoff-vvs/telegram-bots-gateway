package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
	"queuingbot/domain"
	message_handler "queuingbot/internal/handler/message-handler"
	"queuingbot/internal/repository/postgresql"
	"queuingbot/user"

	pb "queuingbot/internal/grpc"
)

const port = 50052

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

	dsn := "host=%s user=%s password=%s dbname=queuing port=%s sslmode=disable TimeZone=Europe/Moscow"
	connection := fmt.Sprintf(dsn, dbHost, dbUser, dbPass, dbPort)
	db, err := gorm.Open(postgres.Open(connection), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	db, err := initDatabase()
	if err != nil {
		panic(err)
	}
	repo := postgresql.NewUserRepository(db)
	userService := user.NewUserService(repo)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterQueueingBotServer(grpcServer, newServer(message_handler.NewMessageHandler(userService)))

	log.Printf("Server is running at localhost:%d", port)
	if err = grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}

func newServer(h message_handler.MessageHandler) *queueingBotServer {
	return &queueingBotServer{MessageHandler: h}
}

type queueingBotServer struct {
	pb.UnimplementedQueueingBotServer
	MessageHandler message_handler.MessageHandler
}

func (s *queueingBotServer) Handle(ctx context.Context, point *pb.QueuingMessage) (*pb.QueuingBoolResult, error) {
	log.Printf("Received point: %v", point.Text)
	return &pb.QueuingBoolResult{Result: true}, nil
}

func (s queueingBotServer) HandleMessage(ctx context.Context, message *pb.QueuingMessage) (*pb.QueuingJsonResult, error) {
	msg := s.MessageHandler.Handle(message)
	jsonString, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return &pb.QueuingJsonResult{Result: string(jsonString)}, nil
}
