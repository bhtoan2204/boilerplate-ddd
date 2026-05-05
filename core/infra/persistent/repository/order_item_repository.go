package repository

import (
	"context"

	"boilerplate-ddd/core/infra/persistent/models"

	"gorm.io/gorm"
)

type OrderItemRepository struct {
	db *gorm.DB
}

func NewOrderItemRepository(db *gorm.DB) *OrderItemRepository {
	return &OrderItemRepository{db: db}
}

func (r *OrderItemRepository) DeleteByOrderID(ctx context.Context, orderID string) error {
	return r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Delete(&models.OrderItemModel{}).
		Error
}

func (r *OrderItemRepository) SaveAll(ctx context.Context, orderItems []models.OrderItemModel) error {
	if len(orderItems) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&orderItems).Error
}

func (r *OrderItemRepository) ReplaceByOrderID(ctx context.Context, orderID string, orderItems []models.OrderItemModel) error {
	if err := r.DeleteByOrderID(ctx, orderID); err != nil {
		return err
	}
	return r.SaveAll(ctx, orderItems)
}
