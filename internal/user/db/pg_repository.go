package repository

import (
	"air-quality-bot/internal/user"
	"air-quality-bot/pkg/postgres"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type repository struct {
	client postgres.Client
}

func NewRepository(client postgres.Client) user.Repository {
	return &repository{
		client: client,
	}
}

func (r *repository) SaveOrUpdateAndReturnIsActive(ctx context.Context, user *user.User) (bool, error) {
	var isActive bool
	if err := r.client.QueryRow(
		ctx,
		saveOrUpdateQuery,
		user.ID,
		user.Username,
		user.ChatID,
		user.LangCode,
		user.IsActive,
		user.LastSeenAt,
	).Scan(&isActive); err != nil {
		return false, err
	}
	return isActive, nil
}

func (r *repository) FindAll(ctx context.Context) ([]user.User, error) {
	rows, err := r.client.Query(ctx, findAllQuery)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var u user.User
		err = rows.Scan(&u.ID, &u.Username, &u.ChatID, &u.IsActive, &u.LangCode, &u.CreatedAt, &u.LastSeenAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *repository) FindById(ctx context.Context, userID int64) (*user.User, error) {
	row := r.client.QueryRow(ctx, findByIdQuery, userID)

	var u user.User
	err := row.Scan(&u.ID, &u.Username, &u.ChatID, &u.LangCode, &u.IsActive, &u.CreatedAt, &u.LastSeenAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user with userID %d does not exist", userID)
		}
		return nil, err
	}

	return &u, nil
}
