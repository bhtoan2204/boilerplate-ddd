package entity

type Shipment struct {
	AbstractEntity
	TrackingNo  string
	Status      string
	FromAddress string
	ToAddress   string
	Items       []OrderItem
}
