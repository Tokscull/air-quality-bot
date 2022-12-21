package location

import (
	"air-quality-bot/internal/user"
)

type Location struct {
	ID        int64     `json:"id"`
	User      user.User `json:"bot_user"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	TimeZone  string    `json:"time_zone"`
}
