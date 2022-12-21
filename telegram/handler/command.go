package handler

import (
	"air-quality-bot/internal/locale"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type commandKey string

const (
	StartCmdKey = commandKey("start")
)

type commandHandler struct {
	handler func(upd *tgbotapi.Update) error
}

func (u *UpdateHandler) InitializeCommands() {
	u.commands[StartCmdKey] = commandHandler{
		handler: u.handleStartCommand,
	}
}

func (u *UpdateHandler) handleCommand(upd *tgbotapi.Update) {
	u.log.Info("Handling command", zap.Int64("chatID", upd.Message.Chat.ID), zap.String("cmd", upd.Message.Command()))
	key := upd.Message.Command()
	if cmd, ok := u.commands[commandKey(key)]; ok {
		go func() {
			_ = cmd.handler(upd)
		}()
	} else {
		u.log.Warn("Can't handle command, handler not found", zap.String("cmd", key))
	}
}

func (u *UpdateHandler) handleStartCommand(upd *tgbotapi.Update) error {
	name := upd.Message.From.FirstName
	if name == "" {
		name = upd.Message.From.UserName
	}

	langCode := upd.Message.From.LanguageCode
	text := fmt.Sprintf(locale.Get(locale.Greeting, langCode), name)

	return u.bot.SendMessageWithMenu(upd.Message.Chat.ID, text, langCode)
}
