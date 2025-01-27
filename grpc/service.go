package grpc

import (
	"google.golang.org/grpc"
)

type Service struct {
	ClientConn *grpc.ClientConn
	Client     interface{}
}

func NewService(conn *grpc.ClientConn, client interface{}) Service {
	return Service{ClientConn: conn, Client: client}
}

func (s Service) Close() {
	err := s.ClientConn.Close()
	if err != nil {
		return
	}
}
