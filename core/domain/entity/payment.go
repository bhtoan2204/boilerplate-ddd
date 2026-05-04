package entity

type Payment struct {
	AbstractEntity
	TransactionCode string
	Amount          float64
	Method          string
}
