package handler

import (
	"air-quality-bot/internal/locale"
	"air-quality-bot/telegram/cache"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (u *UpdateHandler) handleMessage(upd *tgbotapi.Update) error {
	u.log.Info("Handling incoming message", zap.Int64("chatID", upd.Message.Chat.ID), zap.String("msg", upd.Message.Text))
	langCode := upd.Message.From.LanguageCode

	switch upd.Message.Text {
	case locale.Get(locale.AirQualityBtn, langCode):
		return u.handleAirQualityMsg(upd, langCode)
	case locale.Get(locale.NotificationsBtn, langCode):
		return u.notificationsHandler.HandleMenuMsg(upd, int(NotificationsCbKey), langCode)
	case locale.Get(locale.CancelBtn, langCode):
		return u.handleCancelMsg(upd, langCode)
	}

	cacheData, ok := u.bot.Cache[cache.Key(upd.Message.From.ID)]
	delete(u.bot.Cache, cache.Key(upd.Message.From.ID))

	if ok && cacheData.Type == cache.Message {
		switch cacheData.HandlerKey {
		case cache.Notifications:
			return u.notificationsHandler.HandleMessage(upd, cacheData.ChildInfo)
		}
	}

	text := locale.Get(locale.MessageNotRecognisedMsg, langCode)
	return u.bot.SendMessage(upd.Message.Chat.ID, text)
}

func (u *UpdateHandler) handleAirQualityMsg(upd *tgbotapi.Update, langCode string) error {
	text := locale.Get(locale.AirQualityLocationMsg, langCode)
	keyboard := tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonLocation(locale.Get(locale.LocationBtn, langCode)),
			tgbotapi.NewKeyboardButton(locale.Get(locale.CancelBtn, langCode))),
	)

	u.bot.Cache[cache.Key(upd.Message.From.ID)] = cache.Data{
		Type:       cache.Location,
		HandlerKey: cache.AirQuality,
	}

	return u.bot.SendMessageWithKeyBoard(upd.Message.Chat.ID, text, keyboard)
}

func (u *UpdateHandler) handleCancelMsg(upd *tgbotapi.Update, langCode string) error {
	text := locale.Get(locale.MenuMsg, langCode)
	return u.bot.SendMessageWithMenu(upd.Message.Chat.ID, text, langCode)
}
