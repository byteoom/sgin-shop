package service

import (
	"errors"
	"time"

	"sgin/model"
	"sgin/pkg/app"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService struct {
}

func NewOrderService() *OrderService {
	return &OrderService{}
}

// CreateOrder creates a new order along with its items and receiver details
func (s *OrderService) CreateOrder(ctx *app.Context, req *model.ReqOrderCreate) (*model.Order, error) {
	// Create the order
	order := &model.Order{
		OrderNo:          uuid.New().String(),
		UserID:           req.UserId,
		Status:           model.OrderStatusPending,
		ReceiverName:     req.Receiver.ReceiverName,
		ReceiverPhone:    req.Receiver.ReceiverPhone,
		ReceiverEmail:    req.Receiver.ReceiverEmail,
		ReceiverCountry:  req.Receiver.ReceiverCountry,
		ReceiverProvince: req.Receiver.ReceiverProvince,
		ReceiverCity:     req.Receiver.ReceiverCity,
		ReceiverAddress:  req.Receiver.ReceiverAddress,
		ReceiverZip:      req.Receiver.ReceiverZip,
		ReceiverRemark:   req.Receiver.ReceiverRemark,
		CreatedAt:        time.Now().Format(time.DateTime),
		UpdatedAt:        time.Now().Format(time.DateTime),
	}

	productItemUuids := make([]string, 0)
	for _, item := range req.Items {
		productItemUuids = append(productItemUuids, item.ProductItemID)
	}

	productItemMap, err := NewProductService().GetProductItemByUUIDList(ctx, productItemUuids)
	if err != nil {
		return nil, err
	}

	err = ctx.DB.Transaction(func(tx *gorm.DB) error {

		orderItems := make([]*model.OrderItem, 0)

		// Create the order items
		for _, item := range req.Items {
			productItem, ok := productItemMap[item.ProductItemID]
			if !ok {
				return errors.New("product item not found")
			}

			orderItem := &model.OrderItem{
				OrderID:       order.OrderNo,
				ProductItemID: item.ProductItemID,
				Quantity:      item.Quantity,
				Price:         productItem.Price,
				TotalAmount:   productItem.Price * float64(item.Quantity),
				// Additional calculations for price, discount, etc., can be added here
				CreatedAt: time.Now().Format(time.DateTime),
				UpdatedAt: time.Now().Format(time.DateTime),
			}

			order.TotalAmount += orderItem.TotalAmount
			orderItems = append(orderItems, orderItem)

		}

		// Create the order in the database
		err := tx.Create(order).Error
		if err != nil {
			ctx.Logger.Error("Failed to create order", err)
			tx.Rollback()
			return errors.New("failed to create order")
		}

		// Create the order items
		err = tx.Create(orderItems).Error
		if err != nil {
			ctx.Logger.Error("Failed to create order items", err)
			tx.Rollback()
			return errors.New("failed to create order items")
		}

		return nil

	})

	if err != nil {
		return nil, err
	}

	return order, nil
}

// 创建订单，根据购物车
func (s *OrderService) CreateOrderByCart(ctx *app.Context, req *model.ReqOrderCreate) (*model.Order, error) {
	// Create the order
	order := &model.Order{
		OrderNo:          uuid.New().String(),
		UserID:           req.UserId,
		Status:           model.OrderStatusPending,
		ReceiverName:     req.Receiver.ReceiverName,
		ReceiverPhone:    req.Receiver.ReceiverPhone,
		ReceiverEmail:    req.Receiver.ReceiverEmail,
		ReceiverCountry:  req.Receiver.ReceiverCountry,
		ReceiverProvince: req.Receiver.ReceiverProvince,
		ReceiverCity:     req.Receiver.ReceiverCity,
		ReceiverAddress:  req.Receiver.ReceiverAddress,
		ReceiverZip:      req.Receiver.ReceiverZip,
		ReceiverRemark:   req.Receiver.ReceiverRemark,
		CreatedAt:        time.Now().Format(time.DateTime),
		UpdatedAt:        time.Now().Format(time.DateTime),
	}

	// 获取购物车商品

	cartProductMap, err := NewCartService().GetCartByUUIDList(ctx, req.CartUuids)
	if err != nil {
		return nil, err
	}

	err = ctx.DB.Transaction(func(tx *gorm.DB) error {

		orderItems := make([]*model.OrderItem, 0)

		for _, cartUuid := range req.CartUuids {
			if cartItem, ok := cartProductMap[cartUuid]; ok {

				orderItem := &model.OrderItem{
					OrderID:       order.OrderNo,
					ProductItemID: cartItem.ProductItemUuid,
					Quantity:      cartItem.Quantity,
					Price:         cartItem.ProductItem.Price,
					TotalAmount:   cartItem.ProductItem.Price * float64(cartItem.Quantity),
					// Additional calculations for price, discount, etc., can be added here
					CreatedAt: time.Now().Format(time.DateTime),
					UpdatedAt: time.Now().Format(time.DateTime),
				}

				order.TotalAmount += orderItem.TotalAmount
				orderItems = append(orderItems, orderItem)
			}
		}

		// Create the order in the database
		err := tx.Create(order).Error
		if err != nil {
			ctx.Logger.Error("Failed to create order", err)
			tx.Rollback()
			return errors.New("failed to create order")
		}

		// Create the order items
		err = tx.Create(orderItems).Error
		if err != nil {
			ctx.Logger.Error("Failed to create order items", err)
			tx.Rollback()
			return errors.New("failed to create order items")
		}

		return nil

	})

	if err != nil {
		return nil, err
	}

	return order, nil
}

// GetOrderByID retrieves an order by its ID
func (s *OrderService) GetOrderByID(ctx *app.Context, uuidStr string) (*model.Order, error) {
	order := &model.Order{}
	err := ctx.DB.Where("uuid = ?", uuidStr).First(order).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("order not found")
		}
		ctx.Logger.Error("Failed to get order by ID", err)
		return nil, errors.New("failed to get order by ID")
	}

	return order, nil
}

// 根据订单号获取订单商品列表信息
func (s *OrderService) GetOrderItemsByOrderNo(ctx *app.Context, orderNo string) ([]*model.OrderItemRes, error) {
	orderItems := make([]*model.OrderItem, 0)
	err := ctx.DB.Where("order_no = ?", orderNo).Find(&orderItems).Error
	if err != nil {
		ctx.Logger.Error("Failed to get order items by order no", err)
		return nil, errors.New("failed to get order items by order no")
	}

	productItemUuids := make([]string, 0)
	for _, item := range orderItems {
		productItemUuids = append(productItemUuids, item.ProductItemID)
	}

	productItemMap, err := NewProductService().GetProductItemByUUIDList(ctx, productItemUuids)
	if err != nil {
		return nil, err
	}

	orderItemRes := make([]*model.OrderItemRes, 0)
	for _, item := range orderItems {
		itemRes := &model.OrderItemRes{
			OrderItem: *item,
		}
		if productItem, ok := productItemMap[item.ProductItemID]; ok {
			itemRes.ProductItem = productItem
		}
		orderItemRes = append(orderItemRes, itemRes)

	}

	return orderItemRes, nil
}

// UpdateOrderStatus updates the status of an existing order
func (s *OrderService) UpdateOrderStatus(ctx *app.Context, id string, status string) error {
	err := ctx.DB.Model(&model.Order{}).Where("uuid = ?", id).Updates(map[string]interface{}{
		"status":     status,
		"updated_at": time.Now().Format(time.DateTime),
	}).Error
	if err != nil {
		ctx.Logger.Error("Failed to update order status", err)
		return errors.New("failed to update order status")
	}
	return nil
}

// DeleteOrder deletes an order by its ID
func (s *OrderService) DeleteOrder(ctx *app.Context, uuid string) error {

	err := ctx.DB.Transaction(func(tx *gorm.DB) error {
		// Delete the order items
		err := tx.Where("order_id = ?", uuid).Delete(&model.OrderItem{}).Error
		if err != nil {
			ctx.Logger.Error("Failed to delete order items", err)
			tx.Rollback()
			return errors.New("failed to delete order items")
		}

		// Delete the order
		err = tx.Where("uuid = ?", uuid).Delete(&model.Order{}).Error
		if err != nil {
			ctx.Logger.Error("Failed to delete order", err)
			tx.Rollback()
			return errors.New("failed to delete order")
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// GetOrderList retrieves a list of orders based on query parameters
func (s *OrderService) GetOrderList(ctx *app.Context, params *model.ReqOrderQueryParam) (r *model.PagedResponse, err error) {
	var (
		orders []*model.Order
		total  int64
	)

	db := ctx.DB.Model(&model.Order{})

	if params.UserID != "" {
		db = db.Where("user_id = ?", params.UserID)
	}

	if params.Status != "" {
		db = db.Where("status = ?", params.Status)
	}

	if params.OrderNo != "" {
		db = db.Where("order_no = ?", params.OrderNo)
	}
	err = db.Count(&total).Error
	if err != nil {
		ctx.Logger.Error("Failed to get order count", err)
		return nil, errors.New("failed to get order count")
	}

	err = db.Order("id DESC").Offset(params.GetOffset()).Limit(params.PageSize).Error
	if err != nil {
		ctx.Logger.Error("Failed to get order list", err)
		return nil, errors.New("failed to get order list")
	}

	return &model.PagedResponse{
		Total: total,
		Data:  orders,
	}, nil
}
