// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.1
// source: services/queuingbot/queueingbot.proto

package grpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	QueueingBot_Handle_FullMethodName = "/QueueingBot.QueueingBot/Handle"
)

// QueueingBotClient is the client API for QueueingBot service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Сервис, который предоставляет метод для получения информации о точке
type QueueingBotClient interface {
	Handle(ctx context.Context, in *QueuingMessage, opts ...grpc.CallOption) (*QueuingBoolResult, error)
}

type queueingBotClient struct {
	cc grpc.ClientConnInterface
}

func NewQueueingBotClient(cc grpc.ClientConnInterface) QueueingBotClient {
	return &queueingBotClient{cc}
}

func (c *queueingBotClient) Handle(ctx context.Context, in *QueuingMessage, opts ...grpc.CallOption) (*QueuingBoolResult, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueuingBoolResult)
	err := c.cc.Invoke(ctx, QueueingBot_Handle_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueueingBotServer is the server API for QueueingBot service.
// All implementations must embed UnimplementedQueueingBotServer
// for forward compatibility.
//
// Сервис, который предоставляет метод для получения информации о точке
type QueueingBotServer interface {
	Handle(context.Context, *QueuingMessage) (*QueuingBoolResult, error)
	mustEmbedUnimplementedQueueingBotServer()
}

// UnimplementedQueueingBotServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedQueueingBotServer struct{}

func (UnimplementedQueueingBotServer) Handle(context.Context, *QueuingMessage) (*QueuingBoolResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Handle not implemented")
}
func (UnimplementedQueueingBotServer) mustEmbedUnimplementedQueueingBotServer() {}
func (UnimplementedQueueingBotServer) testEmbeddedByValue()                     {}

// UnsafeQueueingBotServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QueueingBotServer will
// result in compilation errors.
type UnsafeQueueingBotServer interface {
	mustEmbedUnimplementedQueueingBotServer()
}

func RegisterQueueingBotServer(s grpc.ServiceRegistrar, srv QueueingBotServer) {
	// If the following call pancis, it indicates UnimplementedQueueingBotServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&QueueingBot_ServiceDesc, srv)
}

func _QueueingBot_Handle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueuingMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueueingBotServer).Handle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: QueueingBot_Handle_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueueingBotServer).Handle(ctx, req.(*QueuingMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// QueueingBot_ServiceDesc is the grpc.ServiceDesc for QueueingBot service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var QueueingBot_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "QueueingBot.QueueingBot",
	HandlerType: (*QueueingBotServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Handle",
			Handler:    _QueueingBot_Handle_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "services/queuingbot/queueingbot.proto",
}
