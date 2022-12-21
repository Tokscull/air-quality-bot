package user

import (
	"context"
)

type Repository interface {
	SaveOrUpdateAndReturnIsActive(ctx context.Context, user *User) (bool, error)
	FindAll(ctx context.Context) ([]User, error)
	FindById(ctx context.Context, userID int64) (*User, error)
}
