package notification

import (
	"air-quality-bot/internal/locale"
	"air-quality-bot/internal/location"
	"air-quality-bot/internal/user"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type locationKey int

const (
	CreateLocationKey locationKey = iota
	EditLocationKey
)

type locationEntity struct {
	parentID       int
	locationKey    locationKey
	notificationID int64
	locationID     int64
}

type locationHandler struct {
	handler func(upd *tgbotapi.Update, entity locationEntity, langCode string) error
}

func (u *Handler) InitializeLocations() {
	u.locations[CreateLocationKey] = locationHandler{handler: u.handleCreateLocation}
	u.locations[EditLocationKey] = locationHandler{handler: u.handleNotificationEditLocation}
}

func (u *Handler) HandleLocation(upd *tgbotapi.Update, data interface{}) error {
	if entity, ok := data.(locationEntity); ok {
		if lc, ok := u.locations[entity.locationKey]; ok {
			langCode := upd.Message.From.LanguageCode
			return lc.handler(upd, entity, langCode)
		}
		u.log.Warn("Can't handler location, key not found", zap.Int("locationKey", int(entity.locationKey)))
		return fmt.Errorf("can't handler location, key not found %d", entity.locationKey)
	}
	u.log.Error("Can't cast childData to locationEntity", zap.Any("childData", data))
	return fmt.Errorf("can't cast childData to locationEntity%s", data)
}

func (u *Handler) handleCreateLocation(upd *tgbotapi.Update, entity locationEntity, langCode string) error {
	_ = u.bot.SendMessageWithMenu(upd.Message.Chat.ID, locale.Get(locale.LocationProcessingMsg, langCode), langCode)

	tz := u.tzFinder.GetTimezoneName(upd.Message.Location.Longitude, upd.Message.Location.Latitude)

	lc := location.Location{
		User: user.User{
			ID: upd.Message.From.ID,
		},
		Latitude:  upd.Message.Location.Latitude,
		Longitude: upd.Message.Location.Longitude,
		TimeZone:  tz,
	}

	id, dbErr := u.locationRepo.Save(context.Background(), &lc)
	if dbErr != nil {
		u.log.Error("Error saving location", zap.Error(dbErr))
		return dbErr
	}

	timepickerCb := callbackEntity{
		parentID:    entity.parentID,
		callbackKey: CreateTimePickerCbKey,
		locationID:  id,
	}

	manualCb := callbackEntity{
		parentID:    entity.parentID,
		callbackKey: CreateTimeZoneInputCbKey,
		locationID:  id,
	}

	text := fmt.Sprintf(locale.Get(locale.TimeZoneClarificationMsg, langCode), tz)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(locale.Get(locale.TimeZoneRightBtn, langCode), marshallCb(timepickerCb)),
			tgbotapi.NewInlineKeyboardButtonData(locale.Get(locale.TimeZoneWrongBtn, langCode), marshallCb(manualCb)),
		))

	return u.bot.SendMessageWithKeyBoard(upd.Message.Chat.ID, text, keyboard)
}

func (u *Handler) handleNotificationEditLocation(upd *tgbotapi.Update, entity locationEntity, langCode string) error {
	_ = u.bot.SendMessageWithMenu(upd.Message.Chat.ID, locale.Get(locale.LocationProcessingMsg, langCode), langCode)

	tz := u.tzFinder.GetTimezoneName(upd.Message.Location.Longitude, upd.Message.Location.Latitude)
	lc := location.Location{
		ID: entity.locationID,
		User: user.User{
			ID: upd.Message.From.ID,
		},
		Latitude:  upd.Message.Location.Latitude,
		Longitude: upd.Message.Location.Longitude,
		TimeZone:  tz,
	}

	dbErr := u.locationRepo.UpdateById(context.Background(), &lc)
	if dbErr != nil {
		u.log.Error("Error updating location", zap.Error(dbErr))
		text := locale.Get(locale.ErrorMsg, langCode)
		return u.bot.SendMessageWithMenu(upd.Message.Chat.ID, text, langCode)
	}

	data := callbackEntity{
		parentID:       entity.parentID,
		callbackKey:    DetailsCbKey,
		notificationID: entity.notificationID,
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(locale.Get(locale.NotificationEditBackToDetailsBtn, langCode), marshallCb(data)),
		))

	text := locale.Get(locale.NotificationLocationUpdatedMsg, langCode)
	return u.bot.SendMessageWithKeyBoard(upd.Message.Chat.ID, text, keyboard)
}
