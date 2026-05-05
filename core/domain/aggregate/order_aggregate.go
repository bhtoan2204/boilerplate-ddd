package aggregate

import (
	"boilerplate-ddd/core/domain/entity"
	"boilerplate-ddd/pkg/abstract"
	"fmt"
	"time"
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
		&EventShipmentCreated{},
		&EventItemAdded{},
		&EventItemRemoved{},
		&EventPaymentMade{},
		&EventOrderCancelled{},
		&EventOrderCompleted{},
	)
}

func (agg *OrderAggregate) Transition(e abstract.Event) error {
	switch data := e.EventData.(type) {
	case *EventShipmentCreated:
		return agg.applyShipmentCreated(data)
	case *EventItemAdded:
		return agg.applyItemAdded(data)
	case *EventItemRemoved:
		return agg.applyItemRemoved(data)
	case *EventPaymentMade:
		return agg.applyPaymentMade(data)
	case *EventOrderCancelled:
		return agg.applyOrderCancelled(data)
	case *EventOrderCompleted:
		return agg.applyOrderCompleted(data)
	default:
		return abstract.ErrUnsupportedEventType
	}
}

func (agg *OrderAggregate) CreateShipment(data *CreateShipmentCommand) error {
	if data == nil {
		return ErrNilData
	}

	if err := agg.checkShipmentInvariants(data.SKUs); err != nil {
		return err
	}

	if err := agg.ApplyChange(agg, &EventShipmentCreated{
		ID:          data.ID,
		TrackingNo:  data.TrackingNo,
		Status:      data.Status,
		FromAddress: data.FromAddress,
		ToAddress:   data.ToAddress,
		SKUs:        data.SKUs,
		Time:        data.Time,
	}); err != nil {
		return fmt.Errorf("apply event failed: %w", err)
	}

	return nil
}

func (agg *OrderAggregate) AddItem(item entity.OrderItem) error {
	if item.ID == "" {
		return ErrNilData
	}

	for _, existing := range agg.orderItems {
		if existing.ID == item.ID {
			return ErrItemAlreadyExists
		}
	}

	if err := agg.ApplyChange(agg, &EventItemAdded{
		Item: item,
		Time: item.CreatedAt,
	}); err != nil {
		return fmt.Errorf("apply event failed: %w", err)
	}

	return nil
}

func (agg *OrderAggregate) RemoveItem(itemID string) error {
	if itemID == "" {
		return ErrNilData
	}

	found := false
	for _, item := range agg.orderItems {
		if item.ID == itemID {
			found = true
			break
		}
	}
	if !found {
		return ErrItemNotFound
	}

	if agg.Status != OrderStatusPending {
		return ErrItemCannotBeRemoved
	}

	if err := agg.ApplyChange(agg, &EventItemRemoved{
		ItemID: itemID,
		Time:   time.Now(),
	}); err != nil {
		return fmt.Errorf("apply event failed: %w", err)
	}

	return nil
}

func (agg *OrderAggregate) Pay(payment entity.Payment) error {
	if payment.ID == "" {
		return ErrNilData
	}

	if agg.Status != OrderStatusPending {
		return ErrInvalidOrderStatus
	}

	for _, existing := range agg.payments {
		if existing.ID == payment.ID {
			return ErrPaymentAlreadyExists
		}
	}

	if err := agg.ApplyChange(agg, &EventPaymentMade{
		Payment: payment,
		Time:    payment.CreatedAt,
	}); err != nil {
		return fmt.Errorf("apply event failed: %w", err)
	}

	return nil
}

func (agg *OrderAggregate) Cancel() error {
	if agg.Status != OrderStatusPending && agg.Status != OrderStatusPaid {
		return ErrOrderCannotBeCancelled
	}

	if err := agg.ApplyChange(agg, &EventOrderCancelled{
		Reason: "",
		Time:   time.Now(),
	}); err != nil {
		return fmt.Errorf("apply event failed: %w", err)
	}

	return nil
}

func (agg *OrderAggregate) Complete() error {
	if agg.Status != OrderStatusPaid {
		return ErrOrderCannotBeCompleted
	}

	if err := agg.ApplyChange(agg, &EventOrderCompleted{
		Time: time.Now(),
	}); err != nil {
		return fmt.Errorf("apply event failed: %w", err)
	}

	return nil
}

func (agg *OrderAggregate) OrderItems() []entity.OrderItem {
	return agg.orderItems
}

func (agg *OrderAggregate) SetOrderItems(items []entity.OrderItem) {
	agg.orderItems = items
}

func (agg *OrderAggregate) Shipments() []entity.Shipment {
	return agg.shipments
}

func (agg *OrderAggregate) SetShipments(shipments []entity.Shipment) {
	agg.shipments = shipments
}

func (agg *OrderAggregate) Payments() []entity.Payment {
	return agg.payments
}

func (agg *OrderAggregate) SetPayments(payments []entity.Payment) {
	agg.payments = payments
}
