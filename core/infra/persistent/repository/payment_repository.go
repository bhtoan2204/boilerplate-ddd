package repository

import (
	"context"

	"boilerplate-ddd/core/infra/persistent/models"

	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) DeleteByOrderID(ctx context.Context, orderID string) error {
	return r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Delete(&models.PaymentModel{}).
		Error
}

func (r *PaymentRepository) SaveAll(ctx context.Context, payments []models.PaymentModel) error {
	if len(payments) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&payments).Error
}

func (r *PaymentRepository) ReplaceByOrderID(ctx context.Context, orderID string, payments []models.PaymentModel) error {
	if err := r.DeleteByOrderID(ctx, orderID); err != nil {
		return err
	}
	return r.SaveAll(ctx, payments)
}
