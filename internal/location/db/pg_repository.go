package location

import (
	"air-quality-bot/internal/location"
	"air-quality-bot/pkg/postgres"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
)

type repository struct {
	client postgres.Client
}

func NewRepository(client postgres.Client) location.Repository {
	return &repository{
		client: client,
	}
}

func (r *repository) Save(ctx context.Context, location *location.Location) (int64, error) {
	var id int64
	if err := r.client.QueryRow(
		ctx,
		saveQuery,
		location.User.ID,
		location.Latitude,
		location.Longitude,
		location.TimeZone,
	).Scan(&id); err != nil {
		return id, err
	}
	return id, nil
}

func (r *repository) FindById(ctx context.Context, id int64) (*location.Location, error) {
	row := r.client.QueryRow(ctx, findByIdQuery, id)

	var u location.Location
	err := row.Scan(&u.ID, &u.User.ID, &u.Latitude, &u.Longitude, &u.TimeZone)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("location with id %d does not exist", id)
		}
		return nil, err
	}

	return &u, nil
}

func (r *repository) UpdateById(ctx context.Context, l *location.Location) error {
	tag, err := r.client.Exec(ctx, updateByIdQuery, l.User.ID, l.Latitude, l.Latitude, l.TimeZone, l.ID)

	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("row not found")
	}

	return nil
}

func (r *repository) UpdateTimeZoneById(ctx context.Context, id int64, timeZone string) error {
	tag, err := r.client.Exec(ctx, updateTimeZoneByIdQuery, timeZone, id)

	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("row not found")
	}

	return nil
}

func (r *repository) UpdateTimeZoneAndNotifyAtById(ctx context.Context, id int64, timeZone string, notifyAt time.Time) error {
	tag, err := r.client.Exec(ctx, updateTimeZoneAndNotifyAtByIdQuery, notifyAt, timeZone, id)

	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("row not found")
	}

	return nil
}
