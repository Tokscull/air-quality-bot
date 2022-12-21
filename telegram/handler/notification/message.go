package notification

import (
	"air-quality-bot/internal/locale"
	"air-quality-bot/pkg/utils"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"time"
)

type messageKey int

const (
	CreateTimeZoneInputMesKey messageKey = iota
	EditTimeZoneInputMesKey
)

type messageEntity struct {
	parentID       int
	messageKey     messageKey
	notificationID int64
	locationID     int64
}

type messageHandler struct {
	handler func(upd *tgbotapi.Update, entity messageEntity, langCode string) error
}

func (u *Handler) InitializeMessages() {
	u.messages[CreateTimeZoneInputMesKey] = messageHandler{handler: u.handleCreateTimeZoneInput}
	u.messages[EditTimeZoneInputMesKey] = messageHandler{handler: u.handleEditTimeZoneInput}
}

func (u *Handler) HandleMessage(upd *tgbotapi.Update, data interface{}) error {
	if entity, ok := data.(messageEntity); ok {
		if mes, k := u.messages[entity.messageKey]; k {
			langCode := upd.Message.From.LanguageCode
			return mes.handler(upd, entity, langCode)
		}
		u.log.Warn("Can't handler message, key not found", zap.Int("messageKey", int(entity.messageKey)))
		return fmt.Errorf("can't handler message, key not found: %d", entity.messageKey)
	}
	u.log.Error("Can't cast childData to messageEntity", zap.Any("childData", data))
	return fmt.Errorf("can't cast childData to messageEntity: %s", data)
}

func (u *Handler) HandleMenuMsg(upd *tgbotapi.Update, pKey int, langCode string) error {

	createCb := callbackEntity{parentID: pKey, callbackKey: CreateCbKey}
	getCb := callbackEntity{parentID: pKey, callbackKey: GetAllCbKey}

	text := locale.Get(locale.NotificationsMsq, langCode)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(locale.Get(locale.NotificationsCreateBtn, langCode), marshallCb(createCb)),
			tgbotapi.NewInlineKeyboardButtonData(locale.Get(locale.NotificationViewBtn, langCode), marshallCb(getCb)),
		))
	return u.bot.SendMessageWithKeyBoard(upd.Message.Chat.ID, text, keyboard)
}

func (u *Handler) handleCreateTimeZoneInput(upd *tgbotapi.Update, entity messageEntity, langCode string) error {

	_, err := time.LoadLocation(upd.Message.Text)
	if err != nil {
		u.log.Error("Error parsing timezone")
		return u.bot.SendMessage(upd.Message.From.ID, locale.Get(locale.ErrorMsg, langCode))
	}

	dbErr := u.locationRepo.UpdateTimeZoneById(context.Background(), entity.locationID, upd.Message.Text)
	if err != nil {
		u.log.Error("Error updating location time_zone by id", zap.Error(dbErr))
		return u.bot.SendMessage(upd.Message.From.ID, locale.Get(locale.ErrorMsg, langCode))
	}

	data := callbackEntity{
		parentID:    entity.parentID,
		callbackKey: CreateTimePickerCbKey,
		locationID:  entity.locationID,
	}

	text := locale.Get(locale.TimeZoneUpdatedMsg, langCode)
	var keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(locale.Get(locale.ContinueBtn, langCode), marshallCb(data)),
		))
	return u.bot.SendMessageWithKeyBoard(upd.Message.From.ID, text, keyboard)
}

func (u *Handler) handleEditTimeZoneInput(upd *tgbotapi.Update, entity messageEntity, langCode string) error {
	newLoc, err := time.LoadLocation(upd.Message.Text)
	if err != nil {
		u.log.Error("Error parsing timezone")
		return u.bot.SendMessage(upd.Message.From.ID, locale.Get(locale.ErrorMsg, langCode))
	}

	ntf, dbErr := u.notificationRepo.FindByIdWithLocation(context.Background(), entity.notificationID)
	if dbErr != nil {
		u.log.Error("Error find notification by id", zap.Error(dbErr))
		return u.bot.SendMessage(upd.Message.From.ID, locale.Get(locale.ErrorMsg, langCode))
	}

	oldLoc, _ := time.LoadLocation(ntf.Location.TimeZone)
	notifyAtOldUtc := utils.TruncateToDay(time.Now()).Add(time.Hour*time.Duration(ntf.NotifyAt.Hour()) +
		time.Minute*time.Duration(ntf.NotifyAt.Minute()) + time.Second*time.Duration(ntf.NotifyAt.Second()))
	notifyAtOldLocal := notifyAtOldUtc.In(oldLoc)

	notifyAtNew := utils.TruncateToDay(time.Now().In(newLoc)).Add(time.Hour*time.Duration(notifyAtOldLocal.Hour()) +
		time.Minute*time.Duration(notifyAtOldLocal.Minute()) + time.Second*time.Duration(notifyAtOldLocal.Second()))

	dbErr = u.locationRepo.UpdateTimeZoneAndNotifyAtById(context.Background(), ntf.Location.ID, upd.Message.Text, notifyAtNew.UTC())
	if dbErr != nil {
		u.log.Error("Error updating notifyAt and time_zone by id", zap.Error(dbErr))
		return u.bot.SendMessage(upd.Message.From.ID, locale.Get(locale.ErrorMsg, langCode))
	}

	data := callbackEntity{
		parentID:       entity.parentID,
		callbackKey:    DetailsCbKey,
		notificationID: entity.notificationID,
		locationID:     entity.locationID,
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(locale.Get(locale.NotificationEditBackToDetailsBtn, langCode), marshallCb(data)),
		))

	textFormat := locale.Get(locale.NotificationTimeZoneUpdatedMsg, langCode)
	text := fmt.Sprintf(textFormat, upd.Message.Text)
	return u.bot.SendMessageWithKeyBoard(upd.Message.From.ID, text, keyboard)
}
