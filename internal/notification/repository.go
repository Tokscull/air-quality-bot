package notification

import (
	"context"
	"time"
)

type Repository interface {
	Save(ctx context.Context, notification *Notification) error
	FindById(ctx context.Context, id int64) (*Notification, error)
	FindByIdWithLocation(ctx context.Context, id int64) (*Notification, error)
	FindAllByUserIdWithLocation(ctx context.Context, userID int64) ([]Notification, error)
	FindAllByNotifyAtPerPage(ctx context.Context, notifyAt time.Time, size int, offset int) ([]Notification, error)
	UpdateIsActiveById(ctx context.Context, id int64) (bool, error)
	UpdateNotifyAtById(ctx context.Context, notifyAt time.Time, id int64) error
	UpdateLastTimeProcessedAt(ctx context.Context, processedAt time.Time, id int64) error
	DeleteByIdWithLocation(ctx context.Context, id int64) error
}
