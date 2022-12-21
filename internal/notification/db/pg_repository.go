package notification

import (
	"air-quality-bot/internal/notification"
	"air-quality-bot/pkg/postgres"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
)

type repository struct {
	client postgres.Client
}

func NewRepository(client postgres.Client) notification.Repository {
	return &repository{
		client: client,
	}
}

func (r *repository) Save(ctx context.Context, ntf *notification.Notification) error {
	if err := r.client.QueryRow(
		ctx,
		saveQuery,
		ntf.User.ID,
		ntf.Location.ID,
		ntf.IsActive,
		ntf.NotifyAt,
	).Scan(&ntf.ID); err != nil {
		return err
	}
	return nil
}

func (r *repository) FindById(ctx context.Context, id int64) (*notification.Notification, error) {
	row := r.client.QueryRow(ctx, findById, id)

	var n notification.Notification
	err := row.Scan(
		&n.ID,
		&n.User.ID,
		&n.Location.ID,
		&n.NotifyAt,
		&n.IsActive,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("notification with id %d does not exist", id)
		}
		return nil, err
	}

	return &n, nil
}

func (r *repository) FindByIdWithLocation(ctx context.Context, id int64) (*notification.Notification, error) {
	row := r.client.QueryRow(ctx, findByIdWithLocationQuery, id)

	var n notification.Notification
	err := row.Scan(
		&n.ID,
		&n.NotifyAt,
		&n.IsActive,
		&n.User.ID,
		&n.Location.ID,
		&n.Location.User.ID,
		&n.Location.Latitude,
		&n.Location.Longitude,
		&n.Location.TimeZone,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("notification with id %d does not exist", id)
		}
		return nil, err
	}

	return &n, nil
}

func (r *repository) FindAllByUserIdWithLocation(ctx context.Context, userID int64) ([]notification.Notification, error) {
	rows, err := r.client.Query(ctx, findAllByUserIdWithLocationQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []notification.Notification
	for rows.Next() {
		var n notification.Notification
		err = rows.Scan(
			&n.ID,
			&n.NotifyAt,
			&n.IsActive,
			&n.User.ID,
			&n.Location.ID,
			&n.Location.User.ID,
			&n.Location.Latitude,
			&n.Location.Longitude,
			&n.Location.TimeZone,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, n)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *repository) FindAllByNotifyAtPerPage(ctx context.Context, notifyAt time.Time, size int, offset int) ([]notification.Notification, error) {
	rows, err := r.client.Query(ctx, findAllByNotifyAtPerPageQuery, notifyAt, size, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []notification.Notification
	for rows.Next() {
		var n notification.Notification
		err = rows.Scan(
			&n.ID,
			&n.NotifyAt,
			&n.IsActive,
			&n.User.ID,
			&n.User.Username,
			&n.User.ChatID,
			&n.User.LangCode,
			&n.User.IsActive,
			&n.Location.ID,
			&n.Location.User.ID,
			&n.Location.Latitude,
			&n.Location.Longitude,
			&n.Location.TimeZone,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, n)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *repository) UpdateLastTimeProcessedAt(ctx context.Context, processedAt time.Time, id int64) error {
	_, err := r.client.Exec(ctx, updateLastTimeProcessedAtQuery, processedAt, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateIsActiveById(ctx context.Context, id int64) (bool, error) {
	var isActive bool
	if err := r.client.QueryRow(ctx, updateIsActiveByIdQuery, id).Scan(&isActive); err != nil {
		return isActive, err
	}
	return isActive, nil
}

func (r *repository) UpdateNotifyAtById(ctx context.Context, notifyAt time.Time, id int64) error {
	tag, err := r.client.Exec(ctx, updateNotifyAtByIdQuery, notifyAt, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("no row found to update")
	}

	return nil
}

func (r *repository) DeleteByIdWithLocation(ctx context.Context, id int64) error {
	tag, err := r.client.Exec(ctx, deleteByIdWithLocationQuery, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("no row found to delete")
	}
	return nil
}
