package models

type ShipmentModel struct {
	ID            string               `gorm:"primaryKey;column:id"`
	OrderID       string               `gorm:"column:order_id;not null;index"`
	TrackingNo    string               `gorm:"column:tracking_no;not null"`
	Status        string               `gorm:"column:status;not null"`
	FromAddress   string               `gorm:"column:from_address;not null"`
	ToAddress     string               `gorm:"column:to_address;not null"`
	ShipmentItems []ShipmentItemModel  `gorm:"foreignKey:ShipmentID;constraint:OnDelete:CASCADE"`
	CreatedAt     int64                `gorm:"column:created_at;not null"`
	UpdatedAt     int64                `gorm:"column:updated_at;not null"`
}

func (ShipmentModel) TableName() string {
	return "shipments"
}
