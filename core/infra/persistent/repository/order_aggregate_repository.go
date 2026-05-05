package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"boilerplate-ddd/core/domain/aggregate"
	"boilerplate-ddd/core/domain/entity"
	"boilerplate-ddd/core/domain/repository"
	"boilerplate-ddd/core/infra/persistent/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GormOrderAggregateRepository struct {
	orderRepo    *OrderRepository
	itemRepo     *OrderItemRepository
	shipmentRepo *ShipmentRepository
	paymentRepo  *PaymentRepository
}

func NewGormOrderAggregateRepository(db *gorm.DB) *GormOrderAggregateRepository {
	return &GormOrderAggregateRepository{
		orderRepo:    NewOrderRepository(db),
		itemRepo:     NewOrderItemRepository(db),
		shipmentRepo: NewShipmentRepository(db),
		paymentRepo:  NewPaymentRepository(db),
	}
}

func NewSqliteGormDB(dsn string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
}

func EnsureOrderAggregateSchema(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).AutoMigrate(
		&models.OrderModel{},
		&models.OrderItemModel{},
		&models.ShipmentModel{},
		&models.ShipmentItemModel{},
		&models.PaymentModel{},
	)
}

func (r *GormOrderAggregateRepository) WithTransaction(ctx context.Context, fn func(repository.AggregateRepository) error) error {
	return r.orderRepo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &GormOrderAggregateRepository{
			orderRepo:    NewOrderRepository(tx),
			itemRepo:     NewOrderItemRepository(tx),
			shipmentRepo: NewShipmentRepository(tx),
			paymentRepo:  NewPaymentRepository(tx),
		}
		return fn(txRepo)
	})
}

func (r *GormOrderAggregateRepository) GetByID(ctx context.Context, id string) (*aggregate.OrderAggregate, error) {
	if id == "" {
		return nil, errors.New("order id is empty")
	}

	orderModel, err := r.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	agg, err := aggregate.NewOrderAggregate(id)
	if err != nil {
		return nil, err
	}

	agg.Order = orderModelToEntity(orderModel)
	agg.SetOrderItems(orderItemModelsToEntities(orderModel.OrderItems))
	agg.SetShipments(shipmentModelsToEntities(orderModel.Shipments))
	agg.SetPayments(paymentModelsToEntities(orderModel.Payments))

	return agg, nil
}

func (r *GormOrderAggregateRepository) Save(ctx context.Context, agg *aggregate.OrderAggregate) error {
	if agg == nil {
		return errors.New("order aggregate is nil")
	}

	if agg.AggregateID() == "" {
		return errors.New("aggregate id is empty")
	}

	orderModel := orderEntityToModel(&agg.Order)
	orderModel.ID = agg.AggregateID()
	if orderModel.CreatedAt == 0 {
		orderModel.CreatedAt = time.Now().UTC().Unix()
	}
	orderModel.UpdatedAt = time.Now().UTC().Unix()

	if err := r.orderRepo.Save(ctx, orderModel); err != nil {
		return err
	}

	if err := r.itemRepo.ReplaceByOrderID(ctx, orderModel.ID, orderItemsToModels(orderModel.ID, agg.OrderItems())); err != nil {
		return err
	}

	if err := r.shipmentRepo.ReplaceByOrderID(ctx, orderModel.ID, shipmentEntitiesToModels(orderModel.ID, agg.Shipments())); err != nil {
		return err
	}

	if err := r.paymentRepo.ReplaceByOrderID(ctx, orderModel.ID, paymentEntitiesToModels(orderModel.ID, agg.Payments())); err != nil {
		return err
	}

	return nil
}

func orderModelToEntity(model *models.OrderModel) entity.Order {
	customerName := ""
	if model.CustomerName.Valid {
		customerName = model.CustomerName.String
	}
	customerPhone := ""
	if model.CustomerPhone.Valid {
		customerPhone = model.CustomerPhone.String
	}

	return entity.Order{
		AbstractEntity: entity.AbstractEntity{
			ID:        model.ID,
			CreatedAt: time.Unix(model.CreatedAt, 0),
			UpdatedAt: time.Unix(model.UpdatedAt, 0),
		},
		Total:         model.Total,
		Discount:      model.Discount,
		Subtotal:      model.Subtotal,
		Status:        model.Status,
		CustomerName:  customerName,
		CustomerPhone: customerPhone,
	}
}

func orderItemModelsToEntities(modelsSlice []models.OrderItemModel) []entity.OrderItem {
	items := make([]entity.OrderItem, 0, len(modelsSlice))
	for _, model := range modelsSlice {
		items = append(items, entity.OrderItem{
			AbstractEntity: entity.AbstractEntity{
				ID:        model.ID,
				CreatedAt: time.Unix(model.CreatedAt, 0),
				UpdatedAt: time.Unix(model.UpdatedAt, 0),
			},
			Sku:         model.Sku,
			Name:        model.Name.String,
			Description: model.Description.String,
			Price:       model.Price,
			Discount:    model.Discount,
			Quantity:    int64(model.Quantity),
		})
	}
	return items
}

func shipmentItemModelsToEntities(modelsSlice []models.ShipmentItemModel) []entity.OrderItem {
	items := make([]entity.OrderItem, 0, len(modelsSlice))
	for _, model := range modelsSlice {
		items = append(items, entity.OrderItem{
			AbstractEntity: entity.AbstractEntity{
				ID:        model.ID,
				CreatedAt: time.Unix(model.CreatedAt, 0),
				UpdatedAt: time.Unix(model.UpdatedAt, 0),
			},
			Sku:         model.Sku,
			Name:        model.Name.String,
			Description: model.Description.String,
			Price:       model.Price,
			Discount:    model.Discount,
			Quantity:    int64(model.Quantity),
		})
	}
	return items
}

func shipmentModelsToEntities(modelsSlice []models.ShipmentModel) []entity.Shipment {
	shipments := make([]entity.Shipment, 0, len(modelsSlice))
	for _, model := range modelsSlice {
		shipments = append(shipments, entity.Shipment{
			AbstractEntity: entity.AbstractEntity{
				ID:        model.ID,
				CreatedAt: time.Unix(model.CreatedAt, 0),
				UpdatedAt: time.Unix(model.UpdatedAt, 0),
			},
			TrackingNo:  model.TrackingNo,
			Status:      model.Status,
			FromAddress: model.FromAddress,
			ToAddress:   model.ToAddress,
			Items:       shipmentItemModelsToEntities(model.ShipmentItems),
		})
	}
	return shipments
}

func paymentModelsToEntities(modelsSlice []models.PaymentModel) []entity.Payment {
	payments := make([]entity.Payment, 0, len(modelsSlice))
	for _, model := range modelsSlice {
		payments = append(payments, entity.Payment{
			AbstractEntity: entity.AbstractEntity{
				ID:        model.ID,
				CreatedAt: time.Unix(model.CreatedAt, 0),
				UpdatedAt: time.Unix(model.UpdatedAt, 0),
			},
			TransactionCode: model.TransactionCode,
			Amount:          model.Amount,
			Method:          model.Method,
		})
	}
	return payments
}

func orderEntityToModel(order *entity.Order) *models.OrderModel {
	customerName := sqlNullString(order.CustomerName)
	customerPhone := sqlNullString(order.CustomerPhone)
	return &models.OrderModel{
		ID:            order.ID,
		Total:         order.Total,
		Discount:      order.Discount,
		Subtotal:      order.Subtotal,
		Status:        order.Status,
		CustomerName:  customerName,
		CustomerPhone: customerPhone,
	}
}

func orderItemsToModels(orderID string, items []entity.OrderItem) []models.OrderItemModel {
	modelsSlice := make([]models.OrderItemModel, 0, len(items))
	for idx, item := range items {
		id := item.ID
		if id == "" {
			id = generateOrderItemID(orderID, item.Sku, idx)
		}
		modelsSlice = append(modelsSlice, models.OrderItemModel{
			ID:          id,
			OrderID:     orderID,
			Sku:         item.Sku,
			Name:        sqlNullString(item.Name),
			Description: sqlNullString(item.Description),
			Price:       item.Price,
			Discount:    item.Discount,
			Quantity:    int(item.Quantity),
			CreatedAt:   time.Now().UTC().Unix(),
			UpdatedAt:   time.Now().UTC().Unix(),
		})
	}
	return modelsSlice
}

func shipmentEntitiesToModels(orderID string, shipments []entity.Shipment) []models.ShipmentModel {
	modelsSlice := make([]models.ShipmentModel, 0, len(shipments))
	for idx, shipment := range shipments {
		id := shipment.ID
		if id == "" {
			id = fmt.Sprintf("%s-shipment-%d", orderID, idx)
		}
		modelsSlice = append(modelsSlice, models.ShipmentModel{
			ID:            id,
			OrderID:       orderID,
			TrackingNo:    shipment.TrackingNo,
			Status:        shipment.Status,
			FromAddress:   shipment.FromAddress,
			ToAddress:     shipment.ToAddress,
			ShipmentItems: shipmentItemEntitiesToModels(id, shipment.Items),
			CreatedAt:     time.Now().UTC().Unix(),
			UpdatedAt:     time.Now().UTC().Unix(),
		})
	}
	return modelsSlice
}

func shipmentItemEntitiesToModels(shipmentID string, items []entity.OrderItem) []models.ShipmentItemModel {
	modelsSlice := make([]models.ShipmentItemModel, 0, len(items))
	for idx, item := range items {
		id := item.ID
		if id == "" {
			id = fmt.Sprintf("%s-shipment-item-%d", shipmentID, idx)
		}
		modelsSlice = append(modelsSlice, models.ShipmentItemModel{
			ID:          id,
			ShipmentID:  shipmentID,
			Sku:         item.Sku,
			Name:        sqlNullString(item.Name),
			Description: sqlNullString(item.Description),
			Price:       item.Price,
			Discount:    item.Discount,
			Quantity:    int(item.Quantity),
			CreatedAt:   time.Now().UTC().Unix(),
			UpdatedAt:   time.Now().UTC().Unix(),
		})
	}
	return modelsSlice
}

func paymentEntitiesToModels(orderID string, payments []entity.Payment) []models.PaymentModel {
	modelsSlice := make([]models.PaymentModel, 0, len(payments))
	for idx, payment := range payments {
		id := payment.ID
		if id == "" {
			id = fmt.Sprintf("%s-payment-%d", orderID, idx)
		}
		modelsSlice = append(modelsSlice, models.PaymentModel{
			ID:              id,
			OrderID:         orderID,
			TransactionCode: payment.TransactionCode,
			Amount:          payment.Amount,
			Method:          payment.Method,
			CreatedAt:       time.Now().UTC().Unix(),
			UpdatedAt:       time.Now().UTC().Unix(),
		})
	}
	return modelsSlice
}

func sqlNullString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: value, Valid: true}
}

func generateOrderItemID(orderID, sku string, idx int) string {
	return fmt.Sprintf("%s-order-item-%s-%d", orderID, sku, idx)
}
