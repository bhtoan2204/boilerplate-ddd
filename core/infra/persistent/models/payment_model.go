package models

type PaymentModel struct {
	ID              string  `gorm:"primaryKey;column:id"`
	OrderID         string  `gorm:"column:order_id;not null;index"`
	TransactionCode string  `gorm:"column:transaction_code;not null"`
	Amount          float64 `gorm:"column:amount;not null"`
	Method          string  `gorm:"column:method;not null"`
	CreatedAt       int64   `gorm:"column:created_at;not null"`
	UpdatedAt       int64   `gorm:"column:updated_at;not null"`
}

func (PaymentModel) TableName() string {
	return "payments"
}
