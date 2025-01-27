package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"telegram-bots-gateway/bot"
	botSettings "telegram-bots-gateway/bot-settings"
	"telegram-bots-gateway/domain"
	"telegram-bots-gateway/internal/handlers"
	pgRep "telegram-bots-gateway/internal/repository/postgresql"

	"github.com/joho/godotenv"
)

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
	dbName := os.Getenv("DATABASE_NAME")

	dsn := "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Moscow"
	connection := fmt.Sprintf(dsn, dbHost, dbUser, dbPass, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(connection), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&domain.BotSettings{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	dbConn, err := initDatabase()
	if err != nil {
		log.Fatal(err)
	}
	botSettingsRepo := pgRep.NewBotSettingsRepository(dbConn)

	svc := botSettings.NewService(botSettingsRepo)
	botsvc := bot.NewService(*svc)
	bots, err := botsvc.GetBots()
	if err != nil {
		log.Fatal(err)
	}

	wg := &sync.WaitGroup{}
	var activeHandlers []handlers.BotHandler
	for _, b := range bots {
		wg.Add(1)
		h := handlers.NewBotHandler(b)
		activeHandlers = append(activeHandlers, *h)
		go func() {
			err := h.Handle(wg)
			if err != nil {
				log.Println(err)
			}
		}()
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	gracefulShutdown(signalChan, wg, activeHandlers)
}

func gracefulShutdown(signalChan chan os.Signal, wg *sync.WaitGroup, handlers []handlers.BotHandler) {
	sigReceived := <-signalChan
	log.Printf("Received signal: %v. Shutting down gracefully...", sigReceived)
	for _, h := range handlers {
		h.Close()
	}
	wg.Wait()

	log.Println("Server gracefully stopped")
}
