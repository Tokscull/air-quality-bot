package location

import (
	"context"
	"time"
)

type Repository interface {
	Save(ctx context.Context, location *Location) (int64, error)
	FindById(ctx context.Context, id int64) (*Location, error)
	UpdateById(ctx context.Context, location *Location) error
	UpdateTimeZoneById(ctx context.Context, id int64, timezone string) error
	UpdateTimeZoneAndNotifyAtById(ctx context.Context, id int64, timezone string, notifyAt time.Time) error
}
