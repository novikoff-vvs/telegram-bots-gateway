package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"

	pb "queuingbot/internal/grpc"
)

const port = 50052

func main() {
	// Создаём слушатель для порта
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Инициализация gRPC сервера
	grpcServer := grpc.NewServer()
	pb.RegisterQueueingBotServer(grpcServer, newServer())

	// Запуск сервера
	log.Printf("Server is running at localhost:%d", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Создаём новый сервер
func newServer() *queueingBotServer {
	return &queueingBotServer{}
}

// Структура сервера
type queueingBotServer struct {
	pb.UnimplementedQueueingBotServer
}

// Реализация метода GetFeature
func (s *queueingBotServer) Handle(ctx context.Context, point *pb.QueuingMessage) (*pb.QueuingBoolResult, error) {
	log.Printf("Received point: %v", point.Text)
	// Возвращаем необработанную точку (без имени)
	return &pb.QueuingBoolResult{Result: true}, nil
}
