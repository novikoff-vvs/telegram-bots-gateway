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

type Bot struct {
	BotAPI      *tgbotapi.BotAPI
	BotSettings BotSettings
	BotClient   grpc.Service
}

func (b Bot) Handle(update *tgbotapi.Update) {
	switch v := b.BotClient.Client.(type) {
	case pb.RouteGuideClient:
		{
			_, err := v.Handle(context.Background(), helpers.ConvertMessageToProto(*update.Message))
			if err != nil {
				log.Println(err)
			}
		}
	case pb.QueueingBotClient:
		{
			if update.Message != nil {
				r, err := v.HandleMessage(context.Background(), helpers.ConvertMessage2ToProto(*update.Message))
				var m tgbotapi.MessageConfig
				if err != nil {
					log.Println(err)
				}
				err = json.Unmarshal([]byte(r.Result), &m)
				if err != nil {
					log.Println(err)
				}
				_, err = b.BotAPI.Send(m)
				if err != nil {
					log.Println(err)
				}
			} else if update.CallbackQuery != nil {
				callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")

				if _, err := b.BotAPI.Request(callback); err != nil {
					log.Println(err)
				}

				log.Println(update.CallbackQuery.Message.Text)
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
