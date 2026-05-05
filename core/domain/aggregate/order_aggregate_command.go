package aggregate

import "time"

type CreateShipmentCommand struct {
	ID          string
	TrackingNo  string
	Status      string
	FromAddress string
	ToAddress   string
	SKUs        []string
	Time        time.Time
}
