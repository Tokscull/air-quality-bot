package notification

import (
	"air-quality-bot/internal/location"
	"air-quality-bot/internal/notification"
	"air-quality-bot/pkg/logger"
	tgBot "air-quality-bot/telegram/bot"
	"github.com/ringsaturn/tzf"
)

type Handler struct {
	log              *logger.Logger
	bot              *tgBot.Bot
	tzFinder         *tzf.DefaultFinder
	notificationRepo notification.Repository
	locationRepo     location.Repository

	callbacks map[callbackKey]callbackHandler
	locations map[locationKey]locationHandler
	messages  map[messageKey]messageHandler
}

func NewHandler(log *logger.Logger, bot *tgBot.Bot, tz *tzf.DefaultFinder, ntfRepo notification.Repository,
	lcRepo location.Repository) *Handler {

	u := &Handler{
		log:              log,
		bot:              bot,
		tzFinder:         tz,
		notificationRepo: ntfRepo,
		locationRepo:     lcRepo,

		callbacks: make(map[callbackKey]callbackHandler),
		locations: make(map[locationKey]locationHandler),
		messages:  make(map[messageKey]messageHandler),
	}

	u.InitializeCallbacks()
	u.InitializeLocations()
	u.InitializeMessages()

	return u
}
