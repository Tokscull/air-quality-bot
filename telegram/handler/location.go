package handler

import (
	"air-quality-bot/internal/locale"
	"air-quality-bot/telegram/cache"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (u *UpdateHandler) handleLocation(upd *tgbotapi.Update) error {
	u.log.Info("Handling location", zap.Int64("chatID", upd.Message.Chat.ID))
	langCode := upd.Message.From.LanguageCode

	cacheData, ok := u.bot.Cache[cache.Key(upd.Message.From.ID)]
	delete(u.bot.Cache, cache.Key(upd.Message.From.ID))

	if ok && cacheData.Type == cache.Location {
		switch cacheData.HandlerKey {
		case cache.AirQuality:
			return u.handleAirQualityLocation(upd, langCode)
		case cache.Notifications:
			return u.notificationsHandler.HandleLocation(upd, cacheData.ChildInfo)
		}
	}

	u.log.Warn("Can't handle location, action not found", zap.Any("type", cacheData.Type),
		zap.Any("handlerKey", cacheData.HandlerKey))
	return u.handleAirQualityLocation(upd, langCode)
}

func (u *UpdateHandler) handleAirQualityLocation(upd *tgbotapi.Update, langCode string) error {
	_ = u.bot.SendMessageWithMenu(upd.Message.Chat.ID, locale.Get(locale.LocationProcessingMsg, langCode), langCode)

	resp, err := u.aqiService.GetByGeo(upd.Message.Location.Latitude, upd.Message.Location.Longitude)
	if err != nil {
		u.log.Error("Failed to load airQuality by geo", zap.Error(err))
		text := locale.Get(locale.ErrorMsg, langCode)
		return u.bot.SendMessageWithMenu(upd.Message.Chat.ID, text, langCode)
	}

	pollutionLvlCb := callbackEntity{callbackKey: AirQualityFaqPollutionLvlCbKey}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(locale.Get(locale.AirQualityPollutionLvlBtn, langCode), marshallCb(pollutionLvlCb)),
		),
	)

	textForm := locale.Get(locale.AirQualityInfoMsg, langCode)
	pollutionLvl := locale.GetPollutionLvl(resp.Level, langCode)

	text := fmt.Sprintf(textForm, resp.Station.Name, resp.AQI, pollutionLvl, resp.Station.URL)
	return u.bot.SendMessageWithKeyBoard(upd.Message.Chat.ID, text, keyboard)
}
