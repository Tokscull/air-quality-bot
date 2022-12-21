package bot

import (
	"air-quality-bot/internal/config"
	"air-quality-bot/internal/locale"
	"air-quality-bot/pkg/logger"
	"air-quality-bot/telegram/cache"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Bot struct {
	API   *tgbotapi.BotAPI
	log   *logger.Logger
	Cache map[cache.Key]cache.Data
}

func NewBot(log *logger.Logger, cfg *config.TgBotConfig) *Bot {
	botApi, err := tgbotapi.NewBotAPI(cfg.Token)

	if err != nil {
		log.Error("Unable to start bot", zap.Error(err))
	}

	b := &Bot{
		API:   botApi,
		log:   log,
		Cache: make(map[cache.Key]cache.Data),
	}

	log.Info("Telegram bot connected", zap.String("username", botApi.Self.UserName))
	return b
}

func (b *Bot) SendMessage(chatID int64, message string) error {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = tgbotapi.ModeMarkdown

	b.log.Info("Sending message", zap.Int64("chatID", chatID))

	_, err := b.API.Send(msg)
	if err != nil {
		b.log.Error("Error sending message", zap.Error(err))
	}
	return err
}

func (b *Bot) SendMessageWithKeyBoard(chatID int64, message string, keyboard interface{}) error {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyMarkup = keyboard

	b.log.Info("Sending message with keyboard", zap.Int64("chatID", chatID))

	_, err := b.API.Send(msg)
	if err != nil {
		b.log.Error("Error sending message with keyboard", zap.Error(err))
	}
	return err
}

func (b *Bot) SendMessageWithMenu(chatID int64, message string, langCode string) error {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(locale.Get(locale.AirQualityBtn, langCode)),
			tgbotapi.NewKeyboardButton(locale.Get(locale.NotificationsBtn, langCode)),
		),
	)

	return b.SendMessageWithKeyBoard(chatID, message, keyboard)
}

func (b *Bot) EditMessageKeyboard(chatID int64, messageId int, keyboard tgbotapi.InlineKeyboardMarkup) error {
	b.log.Info("Edit message keyboard", zap.Int64("chatID", chatID))

	msg := tgbotapi.NewEditMessageReplyMarkup(chatID, messageId, keyboard)

	_, err := b.API.Send(msg)
	if err != nil {
		b.log.Error("Error editing message keyboard", zap.Error(err))
	}
	return err
}

func (b *Bot) DeleteMessage(chatID int64, messageId int) error {
	b.log.Info("Delete message", zap.Int64("chatID", chatID))

	msg := tgbotapi.NewDeleteMessage(chatID, messageId)

	_, err := b.API.Send(msg)
	if err != nil {
		b.log.Error("Error deleting message", zap.Error(err))
	}
	return err
}

func (b *Bot) SendImage(chatID int64, caption string, imgPath string) error {
	b.log.Info("Send image", zap.Int64("chatID", chatID))

	img := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(imgPath))
	img.ParseMode = tgbotapi.ModeMarkdown
	img.Caption = caption
	//img.ReplyMarkup

	_, err := b.API.Send(img)
	if err != nil {
		b.log.Error("Error sending image", zap.Error(err))
	}
	return err
}
