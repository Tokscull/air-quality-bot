package notification

import (
	"air-quality-bot/internal/locale"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func getTimePickerKeyBoard(entity callbackEntity) tgbotapi.InlineKeyboardMarkup {

	data := func(time string) callbackEntity {
		entity.timePickerInput = time
		return entity
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("05:00", marshallCb(data("05h00m"))),
			tgbotapi.NewInlineKeyboardButtonData("06:00", marshallCb(data("06h00m"))),
			tgbotapi.NewInlineKeyboardButtonData("07:00", marshallCb(data("07h00m"))),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("08:00", marshallCb(data("08h00m"))),
			tgbotapi.NewInlineKeyboardButtonData("09:00", marshallCb(data("09h00m"))),
			tgbotapi.NewInlineKeyboardButtonData("10:00", marshallCb(data("10h00m"))),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("11:00", marshallCb(data("11h00m"))),
			tgbotapi.NewInlineKeyboardButtonData("12:00", marshallCb(data("12h00m"))),
			tgbotapi.NewInlineKeyboardButtonData("13:00", marshallCb(data("13h00m"))),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("14:00", marshallCb(data("14h00m"))),
			tgbotapi.NewInlineKeyboardButtonData("15:00", marshallCb(data("15h00m"))),
			tgbotapi.NewInlineKeyboardButtonData("16:00", marshallCb(data("16h00m"))),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("17:00", marshallCb(data("17h00m"))),
			tgbotapi.NewInlineKeyboardButtonData("18:00", marshallCb(data("18h00m"))),
			tgbotapi.NewInlineKeyboardButtonData("19:00", marshallCb(data("19h00m"))),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("20:00", marshallCb(data("20h00m"))),
			tgbotapi.NewInlineKeyboardButtonData("21:00", marshallCb(data("21h00m"))),
			tgbotapi.NewInlineKeyboardButtonData("22:00", marshallCb(data("22h00m"))),
		),
	)
	return keyboard
}

func getDetailsKeyBoard(entity callbackEntity, isNtfActive bool, langCode string) tgbotapi.InlineKeyboardMarkup {

	data := func(callbackKey callbackKey) callbackEntity {
		entity.callbackKey = callbackKey
		return entity
	}

	var pauseText string
	if isNtfActive {
		pauseText = locale.Get(locale.NotificationDetailsPauseBtn, langCode)
	} else {
		pauseText = locale.Get(locale.NotificationDetailsUnpauseBtn, langCode)
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(pauseText, marshallCb(data(PauseUnpauseCbKey))),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(locale.Get(locale.NotificationDetailsEditTimeBtn, langCode), marshallCb(data(EditTimeCbKey))),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(locale.Get(locale.NotificationDetailsEditTimeZoneBtn, langCode), marshallCb(data(EditTimeZoneCbKey))),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(locale.Get(locale.NotificationDetailsEditLocationBtn, langCode), marshallCb(data(EditLocationCbKey))),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(locale.Get(locale.NotificationDetailsDeleteBtn, langCode), marshallCb(data(RemoveCbKey))),
		),
	)

	return keyboard
}
