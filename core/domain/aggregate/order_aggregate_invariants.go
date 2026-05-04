package aggregate

import (
	"boilerplate-ddd/core/domain/entity"
	"fmt"
	"time"
)

func (agg *OrderAggregate) CreateShipment(data *EventCreateShipment) error {
	if data == nil {
		return ErrNilData
	}

	if err := agg.checkShipmentInvariants(data.SKUs); err != nil {
		return err
	}

	newShipment := entity.Shipment{
		AbstractEntity: entity.AbstractEntity{
			ID:        data.ID,
			CreatedAt: data.Time,
			UpdatedAt: data.Time,
		},
		TrackingNo:  data.TrackingNo,
		Status:      data.Status,
		FromAddress: data.FromAddress,
		ToAddress:   data.ToAddress,
		Items:       agg.buildShipmentItemsFromSKUs(data.SKUs, data.Time),
	}

	agg.shipments = append(agg.shipments, newShipment)

	return nil
}

func (agg *OrderAggregate) checkShipmentInvariants(skus []string) error {
	if len(skus) == 0 {
		return ErrShipmentEmpty
	}

	orderedQtyBySKU := make(map[string]int64, len(agg.orderItems))
	for _, item := range agg.orderItems {
		if item.Sku == "" {
			return ErrSkuEmpty
		}

		if item.Quantity <= 0 {
			return fmt.Errorf("%w: sku=%s", ErrOrderItemInvalidQty, item.Sku)
		}

		orderedQtyBySKU[item.Sku] += item.Quantity
	}

	shippedQtyBySKU := make(map[string]int64)
	for _, shipment := range agg.shipments {
		for _, item := range shipment.Items {
			if item.Sku == "" {
				return ErrSkuEmpty
			}

			if item.Quantity <= 0 {
				return fmt.Errorf("%w: sku=%s", ErrShipmentItemInvalidQty, item.Sku)
			}

			shippedQtyBySKU[item.Sku] += item.Quantity
		}
	}

	requestQtyBySKU := make(map[string]int64, len(skus))
	for _, sku := range skus {
		if sku == "" {
			return ErrSkuEmpty
		}

		orderedQty, exists := orderedQtyBySKU[sku]
		if !exists {
			return fmt.Errorf("%w: sku=%s", ErrSkuNotBelongToOrder, sku)
		}

		requestQtyBySKU[sku]++

		if requestQtyBySKU[sku] > orderedQty {
			return fmt.Errorf(
				"%w: sku=%s requested=%d ordered=%d",
				ErrShipmentQtyExceeded,
				sku,
				requestQtyBySKU[sku],
				orderedQty,
			)
		}
	}

	for sku, requestQty := range requestQtyBySKU {
		orderedQty := orderedQtyBySKU[sku]
		alreadyShippedQty := shippedQtyBySKU[sku]

		if alreadyShippedQty+requestQty > orderedQty {
			return fmt.Errorf(
				"%w: sku=%s already_shipped=%d requested=%d ordered=%d",
				ErrShipmentQtyExceeded,
				sku,
				alreadyShippedQty,
				requestQty,
				orderedQty,
			)
		}
	}

	return nil
}

func (agg *OrderAggregate) buildShipmentItemsFromSKUs(skus []string, now time.Time) []entity.OrderItem {
	orderItemBySKU := make(map[string]entity.OrderItem, len(agg.orderItems))
	for _, item := range agg.orderItems {
		orderItemBySKU[item.Sku] = item
	}

	qtyBySKU := make(map[string]int64, len(skus))
	for _, sku := range skus {
		qtyBySKU[sku]++
	}

	items := make([]entity.OrderItem, 0, len(qtyBySKU))
	for sku, qty := range qtyBySKU {
		orderItem := orderItemBySKU[sku]

		items = append(items, entity.OrderItem{
			AbstractEntity: entity.AbstractEntity{
				ID:        orderItem.ID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			Sku:         orderItem.Sku,
			Name:        orderItem.Name,
			Description: orderItem.Description,
			Price:       orderItem.Price,
			Discount:    orderItem.Discount,
			Quantity:    qty,
		})
	}

	return items
}
