package handler

import (
	"air-quality-bot/internal/locale"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type callbackKey int

const (
	NotificationsCbKey callbackKey = iota + 1
	AirQualityFaqPollutionLvlCbKey
)

type callbackEntity struct {
	callbackKey callbackKey
	childInfo   string
}

type callbackHandler struct {
	handler func(upd *tgbotapi.Update, entity callbackEntity) error
}

func marshallCb(d callbackEntity) string {
	return fmt.Sprintf("%d@%s", d.callbackKey, d.childInfo)
}

func unmarshallCb(data string) callbackEntity {
	d := strings.Split(data, "@")

	var cbType int
	if len(d) > 0 {
		cbType, _ = strconv.Atoi(d[0])
	}

	var childInfo string
	if len(d) > 1 {
		childInfo = d[1]
	}

	return callbackEntity{
		callbackKey: callbackKey(cbType),
		childInfo:   childInfo,
	}
}

func (u *UpdateHandler) InitializeCallbacks() {
	u.callbacks[NotificationsCbKey] = callbackHandler{handler: u.handleNotificationsCallback}
	u.callbacks[AirQualityFaqPollutionLvlCbKey] = callbackHandler{handler: u.handleAirQualityPollutionLvlCallback}
}

func (u *UpdateHandler) handleCallback(upd *tgbotapi.Update) {
	u.log.Info("Handling callback", zap.Int64("chatID", upd.CallbackQuery.Message.Chat.ID), zap.String("callbackID", upd.CallbackQuery.Data))

	//remove loading alert
	callback := tgbotapi.NewCallback(upd.CallbackQuery.ID, "")
	_, _ = u.bot.API.Request(callback)

	data := upd.CallbackData()
	entity := unmarshallCb(data)

	if cb, ok := u.callbacks[entity.callbackKey]; ok {
		go func() {
			_ = cb.handler(upd, entity)
		}()
	} else {
		u.log.Warn("Can't handler callback, handler not found", zap.String("cb", data))
	}
}

func (u *UpdateHandler) handleNotificationsCallback(upd *tgbotapi.Update, entity callbackEntity) error {
	return u.notificationsHandler.HandleCallback(upd, int(entity.callbackKey), entity.childInfo)
}

func (u *UpdateHandler) handleAirQualityPollutionLvlCallback(upd *tgbotapi.Update, entity callbackEntity) error {
	langCode := upd.CallbackQuery.From.LanguageCode

	img := locale.GetImage(locale.AirQualityIndexScaleImg, langCode)
	text := locale.Get(locale.AirQualityPollutionLvlMsg, langCode)
	return u.bot.SendImage(upd.CallbackQuery.Message.Chat.ID, text, img)
}
