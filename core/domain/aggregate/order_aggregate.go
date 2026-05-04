package aggregate

import (
	"boilerplate-ddd/core/domain/entity"
	"boilerplate-ddd/pkg/abstract"
)

type OrderAggregate struct {
	abstract.AggregateRoot

	entity.Order
	orderItems []entity.OrderItem
	payments   []entity.Payment
	shipments  []entity.Shipment
}

func NewOrderAggregate(orderID string) (*OrderAggregate, error) {
	agg := &OrderAggregate{}
	if err := abstract.InitAggregate(&agg.AggregateRoot, agg, orderID); err != nil {
		return nil, err
	}

	return agg, nil
}

func (a *OrderAggregate) RegisterEvents(register abstract.RegisterEventsFunc) error {
	return register(
		&EventCreateShipment{},
	)
}

func (agg *OrderAggregate) Transition(e abstract.Event) error {
	switch data := e.EventData.(type) {
	case *EventCreateShipment:
		return agg.CreateShipment(data)
	default:
		return abstract.ErrUnsupportedEventType
	}
}
