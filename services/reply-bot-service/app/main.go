package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"

	pb "replybot/internal/grpc"
)

const port = 50051

func main() {
	// Создаём слушатель для порта
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Инициализация gRPC сервера
	grpcServer := grpc.NewServer()
	pb.RegisterRouteGuideServer(grpcServer, newServer())

	// Запуск сервера
	log.Printf("Server is running at localhost:%d", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Создаём новый сервер
func newServer() *routeGuideServer {
	return &routeGuideServer{}
}

// Структура сервера
type routeGuideServer struct {
	pb.UnimplementedRouteGuideServer
}

// Реализация метода GetFeature
func (s *routeGuideServer) Handle(ctx context.Context, point *pb.Message) (*pb.BoolResult, error) {
	log.Printf("Received point: %v", point.Text)
	// Возвращаем необработанную точку (без имени)
	return &pb.BoolResult{Result: true}, nil
}
