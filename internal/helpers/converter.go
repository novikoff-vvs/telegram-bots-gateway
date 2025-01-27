package helpers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	pb "telegram-bots-gateway/internal/grpc"
)

func ConvertMessageToProto(m tgbotapi.Message) *pb.Message {
	return &pb.Message{
		MessageId: uint64(m.MessageID),
		Text:      m.Text,
	}
}

// TODO нужно вынести это все как-то не в отдельный конверт
func ConvertMessage2ToProto(m tgbotapi.Message) *pb.QueuingMessage {
	return &pb.QueuingMessage{
		MessageId: uint64(m.MessageID),
		Text:      m.Text,
	}
}
