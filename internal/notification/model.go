package notification

import (
	"air-quality-bot/internal/location"
	"air-quality-bot/internal/user"
	"time"
)

type Notification struct {
	ID                  int64             `json:"id"`
	User                user.User         `json:"user_id"`
	Location            location.Location `json:"location_id"`
	IsActive            bool              `json:"is_active"`
	NotifyAt            *time.Time        `json:"notify_at"`
	LastTimeProcessedAt *time.Time        `json:"last_time_processed_at"`
}
