package models

import "database/sql"

type OrderItemModel struct {
	ID          string         `gorm:"primaryKey;column:id"`
	OrderID     string         `gorm:"column:order_id;not null;index"`
	Sku         string         `gorm:"column:sku;not null"`
	Name        sql.NullString `gorm:"column:name;type:text"`
	Description sql.NullString `gorm:"column:description;type:text"`
	Price       float64        `gorm:"column:price;not null"`
	Discount    float64        `gorm:"column:discount;not null"`
	Quantity    int            `gorm:"column:quantity;not null"`
	CreatedAt   int64          `gorm:"column:created_at;not null"`
	UpdatedAt   int64          `gorm:"column:updated_at;not null"`
}

func (OrderItemModel) TableName() string {
	return "order_items"
}
