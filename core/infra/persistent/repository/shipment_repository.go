package repository

import (
	"context"

	"boilerplate-ddd/core/infra/persistent/models"

	"gorm.io/gorm"
)

type ShipmentRepository struct {
	db *gorm.DB
}

func NewShipmentRepository(db *gorm.DB) *ShipmentRepository {
	return &ShipmentRepository{db: db}
}

func (r *ShipmentRepository) DeleteByOrderID(ctx context.Context, orderID string) error {
	return r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Delete(&models.ShipmentModel{}).
		Error
}

func (r *ShipmentRepository) SaveAll(ctx context.Context, shipments []models.ShipmentModel) error {
	if len(shipments) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&shipments).Error
}

func (r *ShipmentRepository) ReplaceByOrderID(ctx context.Context, orderID string, shipments []models.ShipmentModel) error {
	if err := r.DeleteByOrderID(ctx, orderID); err != nil {
		return err
	}
	return r.SaveAll(ctx, shipments)
}
