package domain

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gs "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"telegram-bots-gateway/grpc"
	"telegram-bots-gateway/internal/helpers"

	pb "telegram-bots-gateway/internal/grpc"
)

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

	return Bot{
		BotAPI:      botApi,
		BotSettings: settings,
		BotClient:   grpc.NewService(conn, pb.NewRouteGuideClient(conn)),
	}, nil
}

func (r Bot) GetUpdateChan() <-chan tgbotapi.Update {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 5
	return r.BotAPI.GetUpdatesChan(u)
}
