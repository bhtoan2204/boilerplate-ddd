package repository

import "context"

type AggregateRepository interface {
	WithTransaction(ctx context.Context, fn func(AggregateRepository) error)
}

type OrderAggregateRepository interface {
	GetByID(ctx context.Context, id string)
	Save(ctx context.Context)
}
