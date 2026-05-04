package entity

type Order struct {
	AbstractEntity
	Total         float64
	Discount      float64
	Subtotal      float64
	Status        string
	CustomerName  string
	CustomerPhone string
}
