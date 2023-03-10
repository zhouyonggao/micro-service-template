package repository

import (
	"context"
	"microServiceTemplate/internal/biz/do"
	"time"
)

// OrderRepo is a Greater repo.
type OrderRepo interface {
	FindByID(context.Context, int) *do.OrderDo
	UpdateTimeNow(context.Context, int, time.Time) error
	FindByIDByRedis(context.Context, int) (*do.OrderDo, error)
}
