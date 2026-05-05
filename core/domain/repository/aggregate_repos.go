package repository

import (
	"context"

	"boilerplate-ddd/core/domain/aggregate"
)

type AggregateRepository interface {
	WithTransaction(ctx context.Context, fn func(AggregateRepository) error) error
}

type OrderAggregateRepository interface {
	GetByID(ctx context.Context, id string) (*aggregate.OrderAggregate, error)
	Save(ctx context.Context, order *aggregate.OrderAggregate) error
}
