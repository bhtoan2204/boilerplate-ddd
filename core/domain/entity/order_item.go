package entity

type OrderItem struct {
	AbstractEntity
	Sku         string
	Name        string
	Description string
	Price       float64
	Discount    float64
	Quantity    int64
}
