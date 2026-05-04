package abstract

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrEventStoreNil    = errors.New("event store can not be nil")
	ErrAggregateNil     = errors.New("aggregate can not be nil")
	ErrAggregateRootNil = errors.New("aggregate root can not be nil")
)

//go:generate mockgen -package=event -destination=publisher_mock.go -source=publisher.go
type Store interface {
	Append(ctx context.Context, event Event) error
}

//go:generate mockgen -package=event -destination=publisher_mock.go -source=publisher.go
type Publisher interface {
	Publish(ctx context.Context, events ...Event) error
	PublishAggregate(ctx context.Context, agg Aggregate) error
}

type publisher struct {
	store Store
}

func NewPublisher(store Store) Publisher {
	return &publisher{store: store}
}

func (p *publisher) Publish(ctx context.Context, events ...Event) error {
	if p == nil || p.store == nil {
		return ErrEventStoreNil
	}
	for _, ev := range events {
		if ev.AggregateID == "" {
			return ErrIDEmpty
		}
		if ev.EventName == "" {
			return ErrEventNameEmpty
		}
		if ev.CreatedAt <= 0 {
			ev.CreatedAt = time.Now().Unix()
		}
		if err := p.store.Append(ctx, ev); err != nil {
			return fmt.Errorf("append event=%s failed: %w", ev.EventName, err)
		}
	}
	return nil
}

func (p *publisher) PublishAggregate(ctx context.Context, agg Aggregate) error {
	if agg == nil {
		return ErrAggregateNil
	}
	root := agg.Root()
	if root == nil {
		return ErrAggregateRootNil
	}
	events := root.CloneEvents()
	if len(events) == 0 {
		return nil
	}
	if err := p.Publish(ctx, events...); err != nil {
		return err
	}
	root.MarkPersisted()
	return nil
}
