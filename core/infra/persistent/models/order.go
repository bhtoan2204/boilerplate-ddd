package models

import "database/sql"

type OrderModel struct {
	ID            string           `gorm:"primaryKey;column:id"`
	Total         float64          `gorm:"column:total;not null"`
	Discount      float64          `gorm:"column:discount;not null"`
	Subtotal      float64          `gorm:"column:subtotal;not null"`
	Status        string           `gorm:"column:status;not null"`
	CustomerName  sql.NullString   `gorm:"column:customer_name;type:text"`
	CustomerPhone sql.NullString   `gorm:"column:customer_phone;type:text"`
	OrderItems    []OrderItemModel `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	Shipments     []ShipmentModel  `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	Payments      []PaymentModel   `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	CreatedAt     int64            `gorm:"column:created_at;not null"`
	UpdatedAt     int64            `gorm:"column:updated_at;not null"`
}

func (OrderModel) TableName() string {
	return "orders"
}
