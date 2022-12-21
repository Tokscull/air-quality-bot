package user

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

type User struct {
	ID         int64      `json:"id"`
	Username   string     `json:"user_name"`
	ChatID     int64      `json:"chat_id"`
	LangCode   string     `json:"lang_code"`
	IsActive   bool       `json:"is_active"`
	CreatedAt  *time.Time `json:"created_at"`
	LastSeenAt *time.Time `json:"last_seen_at"`
}

func FromTgUser(tgUser *tgbotapi.User, chat *tgbotapi.Chat) *User {
	now := time.Now()
	u := &User{
		ID:         tgUser.ID,
		Username:   tgUser.UserName,
		ChatID:     chat.ID,
		LangCode:   tgUser.LanguageCode,
		IsActive:   true,
		LastSeenAt: &now,
	}
	return u
}
