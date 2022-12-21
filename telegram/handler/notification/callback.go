package notification

import (
	"air-quality-bot/internal/locale"
	"air-quality-bot/internal/location"
	"air-quality-bot/internal/notification"
	"air-quality-bot/internal/user"
	"air-quality-bot/pkg/utils"
	"air-quality-bot/telegram/cache"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

type callbackKey int

const (
	GetAllCbKey callbackKey = iota + 1
	CreateCbKey
	CreateTimePickerCbKey
	CreateTimeZoneInputCbKey
	SaveCbKey
	DetailsCbKey
	PauseUnpauseCbKey
	EditTimeCbKey
	EditTimeSaveCbKey
	EditTimeZoneCbKey
	EditLocationCbKey
	RemoveCbKey
)

type callbackEntity struct {
	parentID        int
	callbackKey     callbackKey
	notificationID  int64
	locationID      int64
	timePickerInput string
}

type callbackHandler struct {
	handler func(upd *tgbotapi.Update, entity callbackEntity, langCode string) error
}

func (u *Handler) InitializeCallbacks() {
	u.callbacks[GetAllCbKey] = callbackHandler{handler: u.handleGetAllCallback}
	u.callbacks[CreateCbKey] = callbackHandler{handler: u.handleCreateCallback}
	u.callbacks[CreateTimePickerCbKey] = callbackHandler{handler: u.handleCreateTimePickerCallback}
	u.callbacks[CreateTimeZoneInputCbKey] = callbackHandler{handler: u.handleCreateTimeZoneInputCallback}
	u.callbacks[SaveCbKey] = callbackHandler{handler: u.handleSaveCallback}
	u.callbacks[DetailsCbKey] = callbackHandler{handler: u.handleDetailsCallback}
	u.callbacks[PauseUnpauseCbKey] = callbackHandler{handler: u.handlePauseUnpauseCallback}
	u.callbacks[EditTimeCbKey] = callbackHandler{handler: u.handleEditTimeCallback}
	u.callbacks[EditTimeSaveCbKey] = callbackHandler{handler: u.handleEditTimeSaveCallback}
	u.callbacks[EditTimeZoneCbKey] = callbackHandler{handler: u.handleEditTimeZoneCallback}
	u.callbacks[EditLocationCbKey] = callbackHandler{handler: u.handleEditLocationCallback}
	u.callbacks[RemoveCbKey] = callbackHandler{handler: u.handleRemoveCallback}
}

func marshallCb(d callbackEntity) string {
	return fmt.Sprintf("%d@%d_%d;%d;%s", d.parentID, d.callbackKey, d.notificationID, d.locationID, d.timePickerInput)
}

func unmarshallCb(parentId int, data string) callbackEntity {
	d := strings.Split(data, "_")

	var cbType int
	if len(d) > 0 {
		cbType, _ = strconv.Atoi(d[0])
	}

	var params []string
	if len(d) > 1 {
		params = strings.Split(d[1], ";")
	}

	var notificationID int64
	if len(params) > 0 {
		notificationID, _ = strconv.ParseInt(params[0], 10, 64)
	}

	var locationID int64
	if len(params) > 1 {
		locationID, _ = strconv.ParseInt(params[1], 10, 64)
	}

	var timeInp string
	if len(params) > 2 {
		timeInp = params[2]
	}

	return callbackEntity{
		parentID:        parentId,
		callbackKey:     callbackKey(cbType),
		notificationID:  notificationID,
		locationID:      locationID,
		timePickerInput: timeInp,
	}
}

func (u *Handler) HandleCallback(upd *tgbotapi.Update, parentId int, data string) error {
	entity := unmarshallCb(parentId, data)

	if cb, ok := u.callbacks[entity.callbackKey]; ok {
		langCode := upd.CallbackQuery.From.LanguageCode
		return cb.handler(upd, entity, langCode)
	}

	u.log.Warn("Can't handler callback, handler not found", zap.String("cb", data))
	return fmt.Errorf("can't handler callback, handler not found %s", data)
}

func (u *Handler) handleGetAllCallback(upd *tgbotapi.Update, entity callbackEntity, langCode string) error {
	rows, dbErr := u.notificationRepo.FindAllByUserIdWithLocation(context.Background(), upd.CallbackQuery.From.ID)
	if dbErr != nil || len(rows) == 0 {
		u.log.Warn("User notifications not found", zap.Error(dbErr))

		text := locale.Get(locale.NotificationsUserListNotFoundMsg, langCode)
		data := callbackEntity{parentID: entity.parentID, callbackKey: CreateCbKey}
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(locale.Get(locale.NotificationsCreateBtn, langCode), marshallCb(data)),
			))
		return u.bot.SendMessageWithKeyBoard(upd.CallbackQuery.Message.Chat.ID, text, keyboard)
	}

	text := fmt.Sprintf(locale.Get(locale.NotificationsUserListMsg, langCode), len(rows))
	buttons := make([][]tgbotapi.InlineKeyboardButton, 0, len(rows)+1)

	for _, ntf := range rows {
		loc, _ := time.LoadLocation(ntf.Location.TimeZone)

		notifyAt := utils.TruncateToDay(time.Now()).Add(time.Hour*time.Duration(ntf.NotifyAt.Hour()) +
			time.Minute*time.Duration(ntf.NotifyAt.Minute()) + time.Second*time.Duration(ntf.NotifyAt.Second()))
		notifyAtFormatted := notifyAt.In(loc).Format("15:04")

		var bIsActive string
		if ntf.IsActive {
			bIsActive = locale.Get(locale.NotificationStatusActiveMsg, langCode)
		} else {
			bIsActive = locale.Get(locale.NotificationStatusPausedMsg, langCode)
		}

		btnText := fmt.Sprintf(locale.Get(locale.NotificationDetailsBtn, langCode), notifyAtFormatted, bIsActive)

		data := callbackEntity{
			parentID:       entity.parentID,
			callbackKey:    DetailsCbKey,
			notificationID: ntf.ID,
		}

		btn := tgbotapi.NewInlineKeyboardButtonData(btnText, marshallCb(data))
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(btn))
	}
	return u.bot.SendMessageWithKeyBoard(upd.CallbackQuery.Message.Chat.ID, text, tgbotapi.NewInlineKeyboardMarkup(buttons...))
}

func (u *Handler) handleCreateCallback(upd *tgbotapi.Update, entity callbackEntity, langCode string) error {
	text := locale.Get(locale.NotificationLocationMsg, langCode)

	keyboard := tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonLocation(locale.Get(locale.LocationBtn, langCode)),
			tgbotapi.NewKeyboardButton(locale.Get(locale.CancelBtn, langCode))),
	)

	u.bot.Cache[cache.Key(upd.CallbackQuery.From.ID)] = cache.Data{
		Type:       cache.Location,
		HandlerKey: cache.Notifications,
		ChildInfo: locationEntity{
			parentID:    entity.parentID,
			locationKey: CreateLocationKey,
		},
	}

	return u.bot.SendMessageWithKeyBoard(upd.CallbackQuery.Message.Chat.ID, text, keyboard)
}

func (u *Handler) handleCreateTimeZoneInputCallback(upd *tgbotapi.Update, entity callbackEntity, langCode string) error {
	u.bot.Cache[cache.Key(upd.CallbackQuery.From.ID)] = cache.Data{
		Type:       cache.Message,
		HandlerKey: cache.Notifications,
		ChildInfo: messageEntity{
			parentID:   entity.parentID,
			messageKey: CreateTimeZoneInputMesKey,
			locationID: entity.locationID,
		},
	}

	text := locale.Get(locale.TimeZoneManualMsg, langCode)
	return u.bot.SendMessage(upd.CallbackQuery.Message.Chat.ID, text)
}

func (u *Handler) handleCreateTimePickerCallback(upd *tgbotapi.Update, entity callbackEntity, langCode string) error {
	text := locale.Get(locale.TimePickerMsg, langCode)

	data := callbackEntity{
		parentID:    entity.parentID,
		callbackKey: SaveCbKey,
		locationID:  entity.locationID,
	}
	keyboard := getTimePickerKeyBoard(data)
	return u.bot.SendMessageWithKeyBoard(upd.CallbackQuery.Message.Chat.ID, text, keyboard)
}

func (u *Handler) handleDetailsCallback(upd *tgbotapi.Update, entity callbackEntity, langCode string) error {
	ntf, dbErr := u.notificationRepo.FindByIdWithLocation(context.Background(), entity.notificationID)
	if dbErr != nil {
		u.log.Warn("Notification not found", zap.Error(dbErr))
		return dbErr
	}

	loc, _ := time.LoadLocation(ntf.Location.TimeZone)
	notifyAt := utils.TruncateToDay(time.Now()).Add(time.Hour*time.Duration(ntf.NotifyAt.Hour()) +
		time.Minute*time.Duration(ntf.NotifyAt.Minute()) + time.Second*time.Duration(ntf.NotifyAt.Second()))
	notifyAtFormatted := notifyAt.In(loc).Format("15:04")
	var mIsActive string
	if ntf.IsActive {
		mIsActive = locale.Get(locale.NotificationStatusActiveMsg, langCode)
	} else {
		mIsActive = locale.Get(locale.NotificationStatusPausedMsg, langCode)
	}

	text := fmt.Sprintf(locale.Get(locale.NotificationDetailsMsg, langCode), notifyAtFormatted, ntf.Location.TimeZone, mIsActive)
	err := u.bot.SendMessage(upd.CallbackQuery.Message.Chat.ID, text)
	if err != nil {
		return err
	}

	_, lcErr := u.bot.API.Send(tgbotapi.NewLocation(upd.CallbackQuery.Message.Chat.ID, ntf.Location.Latitude, ntf.Location.Longitude))
	if lcErr != nil {
		_ = u.bot.SendMessage(upd.CallbackQuery.Message.Chat.ID, locale.Get(locale.ErrorLoadLocationMsg, langCode))
	}

	data := callbackEntity{
		parentID:       entity.parentID,
		notificationID: ntf.ID,
		locationID:     ntf.Location.ID,
	}

	keyboard := getDetailsKeyBoard(data, ntf.IsActive, langCode)
	return u.bot.SendMessageWithKeyBoard(upd.CallbackQuery.Message.Chat.ID, locale.Get(locale.NotificationDetailsManageMsg, langCode), keyboard)
}

func (u *Handler) handlePauseUnpauseCallback(upd *tgbotapi.Update, entity callbackEntity, langCode string) error {
	isActive, dbErr := u.notificationRepo.UpdateIsActiveById(context.Background(), entity.notificationID)
	if dbErr != nil {
		u.log.Error("error updating notification is_active", zap.Error(dbErr))
		return dbErr
	}

	var text string
	if isActive {
		text = locale.Get(locale.NotificationUnpausedMsg, langCode)
	} else {
		text = locale.Get(locale.NotificationPausedMsg, langCode)
	}
	_ = u.bot.SendMessage(upd.CallbackQuery.Message.Chat.ID, text)

	data := callbackEntity{
		parentID:       entity.parentID,
		notificationID: entity.notificationID,
		locationID:     entity.locationID,
	}

	keyboard := getDetailsKeyBoard(data, isActive, langCode)
	return u.bot.EditMessageKeyboard(upd.CallbackQuery.Message.Chat.ID, upd.CallbackQuery.Message.MessageID, keyboard)
}

func (u *Handler) handleEditTimeCallback(upd *tgbotapi.Update, entity callbackEntity, langCode string) error {
	text := locale.Get(locale.NotificationEditTimeMsg, langCode)

	data := callbackEntity{
		parentID:       entity.parentID,
		callbackKey:    EditTimeSaveCbKey,
		notificationID: entity.notificationID,
		locationID:     entity.locationID,
	}
	keyboard := getTimePickerKeyBoard(data)
	return u.bot.SendMessageWithKeyBoard(upd.CallbackQuery.Message.Chat.ID, text, keyboard)
}

func (u *Handler) handleSaveCallback(upd *tgbotapi.Update, entity callbackEntity, langCode string) error {
	loc, err := u.locationRepo.FindById(context.Background(), entity.locationID)
	if err != nil {
		u.log.Error("Error find location", zap.Error(err))
	}

	timeLoc, _ := time.LoadLocation(loc.TimeZone)
	ntfDuration, _ := time.ParseDuration(entity.timePickerInput)
	today := utils.TruncateToDay(time.Now().In(timeLoc))
	notifyAt := today.Add(ntfDuration).UTC()

	newNotification := notification.Notification{
		User: user.User{
			ID: upd.CallbackQuery.From.ID,
		},
		Location: location.Location{
			ID: loc.ID,
		},
		NotifyAt: &notifyAt,
		IsActive: true,
	}

	dbErr := u.notificationRepo.Save(context.Background(), &newNotification)
	if dbErr != nil {
		u.log.Info("Error creating new notification", zap.Error(dbErr))
		return dbErr
	}

	text := fmt.Sprintf(locale.Get(locale.NotificationCreatedMsg, langCode), entity.timePickerInput)
	return u.bot.SendMessage(upd.CallbackQuery.Message.Chat.ID, text)
}

func (u *Handler) handleEditTimeSaveCallback(upd *tgbotapi.Update, entity callbackEntity, langCode string) error {
	loc, err := u.locationRepo.FindById(context.Background(), entity.locationID)
	if err != nil {
		u.log.Error("Error find location", zap.Error(err))
		return err
	}

	tl, _ := time.LoadLocation(loc.TimeZone)
	ntfTime, _ := time.ParseDuration(entity.timePickerInput)
	today := utils.TruncateToDay(time.Now().In(tl))
	notifyAt := today.Add(ntfTime).UTC()

	fmt.Println(today.Add(ntfTime))
	fmt.Println(today.Add(ntfTime))
	fmt.Println(today.Add(ntfTime).Format("15:04"))

	dbErr := u.notificationRepo.UpdateNotifyAtById(context.Background(), notifyAt, entity.notificationID)
	if err != nil {
		u.log.Error("Error updating notification notify_at", zap.Error(dbErr))
		return err
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

	textFormat := locale.Get(locale.NotificationTimeUpdatedMsg, langCode)
	text := fmt.Sprintf(textFormat, today.Add(ntfTime).Format("15:04"))
	return u.bot.SendMessageWithKeyBoard(upd.CallbackQuery.Message.Chat.ID, text, keyboard)
}

func (u *Handler) handleEditTimeZoneCallback(upd *tgbotapi.Update, entity callbackEntity, langCode string) error {
	u.bot.Cache[cache.Key(upd.CallbackQuery.From.ID)] = cache.Data{
		Type:       cache.Message,
		HandlerKey: cache.Notifications,
		ChildInfo: messageEntity{
			parentID:       entity.parentID,
			messageKey:     EditTimeZoneInputMesKey,
			notificationID: entity.notificationID,
			locationID:     entity.locationID,
		},
	}

	text := locale.Get(locale.NotificationEditTimeZoneMsg, langCode)
	return u.bot.SendMessage(upd.CallbackQuery.Message.Chat.ID, text)
}

func (u *Handler) handleEditLocationCallback(upd *tgbotapi.Update, entity callbackEntity, langCode string) error {
	text := locale.Get(locale.NotificationEditLocationMsg, langCode)

	keyboard := tgbotapi.NewOneTimeReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonLocation(locale.Get(locale.LocationBtn, langCode)),
			tgbotapi.NewKeyboardButton(locale.Get(locale.CancelBtn, langCode))),
	)

	u.bot.Cache[cache.Key(upd.CallbackQuery.From.ID)] = cache.Data{
		Type:       cache.Location,
		HandlerKey: cache.Notifications,
		ChildInfo: locationEntity{
			parentID:       entity.parentID,
			locationKey:    EditLocationKey,
			notificationID: entity.notificationID,
			locationID:     entity.locationID,
		},
	}

	return u.bot.SendMessageWithKeyBoard(upd.CallbackQuery.Message.Chat.ID, text, keyboard)
}

func (u *Handler) handleRemoveCallback(upd *tgbotapi.Update, entity callbackEntity, langCode string) error {
	dbErr := u.notificationRepo.DeleteByIdWithLocation(context.Background(), entity.notificationID)
	if dbErr != nil {
		u.log.Error("error delete notification by id", zap.Error(dbErr))
		return dbErr
	}

	getCb := callbackEntity{parentID: entity.parentID, callbackKey: GetAllCbKey}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(locale.Get(locale.NotificationBackToListBtn, langCode), marshallCb(getCb)),
		))

	text := locale.Get(locale.NotificationDeletedMsg, langCode)
	return u.bot.SendMessageWithKeyBoard(upd.CallbackQuery.Message.Chat.ID, text, keyboard)
}
