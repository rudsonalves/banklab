package domain

import "context"

type Repository interface {
	Create(ctx context.Context, c *Customer) error
}
