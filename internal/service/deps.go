package service

import (
	"context"

	"github.com/polonkoevv/wb-tech/internal/models"
)

type Repo interface {
	LoadCache(ctx context.Context) (map[string]models.Order, error)
	Save(ctx context.Context, order models.Order) error
}
