package handler

import (
	"air-quality-bot/internal/config"
	"air-quality-bot/internal/locale"
	"air-quality-bot/internal/location"
	"air-quality-bot/internal/notification"
	"air-quality-bot/internal/user"
	"air-quality-bot/pkg/logger"
	"air-quality-bot/pkg/waqi"
	tgBot "air-quality-bot/telegram/bot"
	notificationHandler "air-quality-bot/telegram/handler/notification"
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ringsaturn/tzf"
	"go.uber.org/zap"
	"net/http"
)

type UpdateHandler struct {
	log                  *logger.Logger
	bot                  *tgBot.Bot
	aqiService           *waqi.Service
	userRepo             user.Repository
	notificationsHandler *notificationHandler.Handler

	commands  map[commandKey]commandHandler
	callbacks map[callbackKey]callbackHandler
}

func NewUpdateHandler(log *logger.Logger, bot *tgBot.Bot, aqiService *waqi.Service, tzFinder *tzf.DefaultFinder,
	userRepo user.Repository, locationRepo location.Repository, notificationRepo notification.Repository) *UpdateHandler {

	ntfHandler := notificationHandler.NewHandler(log, bot, tzFinder, notificationRepo, locationRepo)

	u := &UpdateHandler{
		log:                  log,
		bot:                  bot,
		aqiService:           aqiService,
		userRepo:             userRepo,
		notificationsHandler: ntfHandler,
		commands:             make(map[commandKey]commandHandler),
		callbacks:            make(map[callbackKey]callbackHandler),
	}
	u.InitializeCommands()
	u.InitializeCallbacks()

	return u
}

func (u *UpdateHandler) ServeWebhookHandler(cfg *config.TgBotConfig) {
	wh, _ := tgbotapi.NewWebhook(cfg.Url + "/api/updates")
	_, errR := u.bot.API.Request(wh)
	if errR != nil {
		u.log.Error("Error setting webhook", zap.Error(errR))
	}

	http.HandleFunc("/api/updates", u.HandleWebHook)

	u.log.Info("Telegram webhook start listening for updates", zap.String("port", cfg.Port))
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		u.log.Error("Error serve webhook", zap.Error(err))
	}
}

func parseTelegramRequest(r *http.Request) (*tgbotapi.Update, error) {
	var update tgbotapi.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		return nil, err
	}
	return &update, nil
}

func (u *UpdateHandler) HandleWebHook(w http.ResponseWriter, r *http.Request) {
	var update, err = parseTelegramRequest(r)
	if err != nil {
		u.log.Error("Error parsing update", zap.Error(err))
		return
	}
	u.handleUpdate(update)
}

func (u *UpdateHandler) handleUpdate(upd *tgbotapi.Update) {
	if upd.Message != nil {
		if isAuth := u.validateUser(upd.Message.From, upd.Message.Chat); isAuth {
			switch {
			case upd.Message.IsCommand():
				u.handleCommand(upd)
			case upd.Message.Location != nil:
				_ = u.handleLocation(upd)
			default:
				_ = u.handleMessage(upd)
			}
		} else {
			_ = u.bot.SendMessage(upd.Message.Chat.ID, locale.Get(locale.AccessDenied, upd.Message.From.LanguageCode))
		}
	}

	if upd.CallbackQuery != nil {
		if isAuth := u.validateUser(upd.CallbackQuery.From, upd.CallbackQuery.Message.Chat); isAuth {
			u.handleCallback(upd)
		} else {
			_ = u.bot.SendMessage(upd.CallbackQuery.Message.Chat.ID, locale.Get(locale.AccessDenied, upd.CallbackQuery.From.LanguageCode))
		}
	}
}

func (u *UpdateHandler) validateUser(tgUser *tgbotapi.User, chat *tgbotapi.Chat) bool {
	currentUser := user.FromTgUser(tgUser, chat)
	isActive, dbErr := u.userRepo.SaveOrUpdateAndReturnIsActive(context.Background(), currentUser)
	if dbErr != nil {
		u.log.Info("Error creating or updating user", zap.Error(dbErr))
		return false
	}
	return isActive
}
