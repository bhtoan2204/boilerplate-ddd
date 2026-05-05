package repository

import (
	"context"

	"boilerplate-ddd/core/infra/persistent/models"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) FindByID(ctx context.Context, id string) (*models.OrderModel, error) {
	var order models.OrderModel
	if err := r.db.WithContext(ctx).
		Preload("OrderItems").
		Preload("Shipments.ShipmentItems").
		Preload("Payments").
		First(&order, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) Save(ctx context.Context, order *models.OrderModel) error {
	return r.db.WithContext(ctx).
		Save(order).Error
}

func (r *OrderRepository) Migrate(ctx context.Context) error {
	return r.db.WithContext(ctx).AutoMigrate(
		&models.OrderModel{},
		&models.OrderItemModel{},
		&models.ShipmentModel{},
		&models.ShipmentItemModel{},
		&models.PaymentModel{},
	)
}
