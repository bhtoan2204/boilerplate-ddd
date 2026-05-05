package aggregate

import (
	"time"

	"boilerplate-ddd/core/domain/entity"
)

type EventShipmentCreated struct {
	ID          string
	TrackingNo  string
	Status      string
	FromAddress string
	ToAddress   string
	SKUs        []string
	Time        time.Time
}

type EventItemAdded struct {
	Item entity.OrderItem
	Time time.Time
}

type EventItemRemoved struct {
	ItemID string
	Time   time.Time
}

type EventPaymentMade struct {
	Payment entity.Payment
	Time    time.Time
}

type EventOrderCancelled struct {
	Reason string
	Time   time.Time
}

type EventOrderCompleted struct {
	Time time.Time
}
