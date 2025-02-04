package domain

import (
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gs "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"telegram-bots-gateway/grpc"
	"telegram-bots-gateway/internal/helpers"

	pb "telegram-bots-gateway/internal/grpc"
)

const ReplyBotType = "replybot"
const QueuingBotType = "queuingbot"
const FilesBotType = "filesbot"

type Bot struct {
	BotAPI      *tgbotapi.BotAPI
	BotSettings BotSettings
	BotClient   grpc.Service
}

func (r Bot) Handle(update *tgbotapi.Update) {
	switch v := r.BotClient.Client.(type) {
	case pb.RouteGuideClient:
		{
			_, err := v.Handle(context.Background(), helpers.ConvertMessageToProto(*update.Message))
			if err != nil {
				return
			}
		}
	case pb.QueueingBotClient:
		{
			_, err := v.Handle(context.Background(), helpers.ConvertMessage2ToProto(*update.Message))
			if err != nil {
				return
			}
		}
	case pb.FilesBotClient:
		{
			msg, _ := json.Marshal(update.Message)
			response, err := v.Handle(context.Background(), &pb.FilesBotMessage{Json: string(msg)})
			if err != nil {
				log.Println(err)
				return
			}

			if response.NewMessageJson != "" {
				var newMessage tgbotapi.MessageConfig
				err = json.Unmarshal([]byte(response.NewMessageJson), &newMessage)
				if err != nil {
					log.Println(err)
					return
				}
				_, err = r.BotAPI.Send(newMessage)
				if err != nil {
					log.Println(err)
					return
				}
			} else if response.ForwardMessageJson != "" {
				var forwardMessage tgbotapi.ForwardConfig
				err = json.Unmarshal([]byte(response.ForwardMessageJson), &forwardMessage)
				if err != nil {
					log.Println(err)
					return
				}
				_, err = r.BotAPI.Send(forwardMessage)
				if err != nil {
					log.Println(err)
					return
				}
			} else {
				log.Println("Не пришло сообщение!")
				return
			}
		}
	}
}

func NewBot(settings BotSettings) (bot Bot, err error) {
	botApi, err := tgbotapi.NewBotAPI(settings.ApiKey)
	if err != nil {
		return Bot{}, err
	}

	var opts = []gs.DialOption{
		gs.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := gs.NewClient(settings.ServiceUrl, opts...)
	if err != nil {
		return Bot{}, err
	}

	var client interface{}
	switch settings.ServiceType {
	case ReplyBotType:
		{
			client = pb.NewRouteGuideClient(conn)
		}
	case QueuingBotType:
		{
			client = pb.NewQueueingBotClient(conn)
		}
	case FilesBotType:
		{
			client = pb.NewFilesBotClient(conn)
		}
	}

	return Bot{
		BotAPI:      botApi,
		BotSettings: settings,
		BotClient:   grpc.NewService(conn, client),
	}, nil
}

func (r Bot) GetUpdateChan() <-chan tgbotapi.Update {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 5
	return r.BotAPI.GetUpdatesChan(u)
}
