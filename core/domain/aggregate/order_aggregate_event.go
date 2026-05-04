package aggregate

import "time"

type EventCreateShipment struct {
	ID          string
	TrackingNo  string
	Status      string
	FromAddress string
	ToAddress   string
	SKUs        []string
	Time        time.Time
}
