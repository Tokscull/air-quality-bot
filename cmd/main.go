package main

import (
	"air-quality-bot/internal/config"
	location "air-quality-bot/internal/location/db"
	notification "air-quality-bot/internal/notification/db"
	user "air-quality-bot/internal/user/db"
	"air-quality-bot/pkg/logger"
	"air-quality-bot/pkg/postgres"
	"air-quality-bot/pkg/waqi"
	"air-quality-bot/scheduler"
	telegramBot "air-quality-bot/telegram/bot"
	"air-quality-bot/telegram/handler"
	"context"
	"flag"
	"github.com/ringsaturn/tzf"
	"go.uber.org/zap"
)

func main() {
	log := logger.NewLogger()

	isProd := flag.Bool("prod", false, "running in a production environment")
	flag.Parse()

	cfg, err := config.GetConfig(isProd)
	if err != nil {
		log.Fatal("Error creating config", zap.Error(err))
	}

	dbClient, err := postgres.NewClient(log, context.Background(), cfg.Postgresql)
	if err != nil {
		log.Fatal("Error creating postgres client", zap.Error(err))
	}

	userRepo := user.NewRepository(dbClient)
	locationRepo := location.NewRepository(dbClient)
	notificationRepo := notification.NewRepository(dbClient)

	tzFinder, err := tzf.NewDefaultFinder()
	if err != nil {
		log.Fatal("Error creating timezone finder", zap.Error(err))
	}

	aqiService := waqi.NewService(cfg.WAQIApiToken)

	tgBot := telegramBot.NewBot(log, cfg.TgBot)

	sch := scheduler.NewScheduler(log, tgBot, aqiService, notificationRepo)
	sch.Start()

	tgBotHandler := handler.NewUpdateHandler(log, tgBot, aqiService, tzFinder, userRepo, locationRepo, notificationRepo)
	tgBotHandler.ServeWebhookHandler(cfg.TgBot)
}
