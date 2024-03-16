package user

import (
	"context"

	"github.com/phuwn/crawlie/src/model"
)

// Store - user store interface
type Store interface {
	Get(ctx context.Context, id string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Save(ctx context.Context, user *model.User) error
}
