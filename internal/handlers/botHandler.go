package handlers

import (
	"sync"
	"telegram-bots-gateway/domain"
)

type BotHandler struct {
	Bot            domain.Bot
	ShutDownSignal chan bool
}

func (h BotHandler) Close() {
	h.ShutDownSignal <- true
}

func NewBotHandler(bot domain.Bot) *BotHandler {
	return &BotHandler{Bot: bot, ShutDownSignal: make(chan bool, 1)}
}

func (h BotHandler) Handle(wg *sync.WaitGroup) error {
	updates := h.Bot.GetUpdateChan()
	for {
		select {
		case <-h.ShutDownSignal:
			{
				wg.Done()
				return nil
			}
		case update, _ := <-updates:
			{
				h.Bot.Handle(&update)
			}
		}
	}
}
