package service

import (
	"errors"
	"sgin/model"
	"sgin/pkg/app"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 订单模块的管理 （创建订单、删除订单、查看订单、修改订单，主要是修改状态）

type OrderServer struct{}

func NewOrderServer() *OrderServer {
	return &OrderServer{}
}

func (o *OrderServer) CreateOrder(ctx *app.Context, order *model.ReqOrderCreate) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	_orderNo := uuid.New().String() // 关联编号
	_order := &model.Order{
		OrderNo:          _orderNo,
		UserID:           order.UserID,
		TotalAmount:      order.TotalAmount,
		Status:           order.Status,
		ReceiverName:     order.ReceiverName,
		ReceiverPhone:    order.ReceiverPhone,
		ReceiverEmail:    order.ReceiverEmail,
		ReceiverCountry:  order.ReceiverCountry,
		ReceiverProvince: order.ReceiverProvince,
		ReceiverCity:     order.ReceiverCity,
		ReceiverAddress:  order.ReceiverAddress,
		ReceiverZip:      order.ReceiverZip,
		ReceiverRemark:   order.ReceiverRemark,
		PaidAt:           now,
		DeliveredAt:      now,
		CompletedAt:      now,
		ClosedAt:         now,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	err := ctx.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&_order).Error
		if err != nil {
			ctx.Logger.Error("create order failed", err)
			tx.Rollback()
			return errors.New("create order failed")
		}
		orderItemList := make([]*model.OrderItem, 0)

		for _, item := range order.OrderItems {
			orderItem := model.OrderItem{
				OrderID:        _orderNo,
				ProductItemID:  item.ProductItemID,
				Quantity:       item.Quantity,
				Price:          item.Price,
				TotalAmount:    item.TotalAmount,
				DiscountAmount: item.DiscountAmount,
				Discount:       item.Discount,
				DiscountPrice:  item.DiscountPrice,
				CreatedAt:      now,
				UpdatedAt:      now,
			}
			orderItemList = append(orderItemList, &orderItem)
		}
		err = tx.Create(&orderItemList).Error
		if err != nil {
			ctx.Logger.Error("create order item failed", err)
			tx.Rollback()
			return errors.New("create order item failed")
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (o *OrderServer) DeleteOrder(ctx *app.Context, orderNo string) error {
	err := ctx.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("order_no = ?", orderNo).Updates(&model.Order{
			IsDeleted: 1,
		}).Error

		if err != nil {
			ctx.Logger.Error("Failed to delete order", err)
			tx.Rollback()
			return errors.New("failed to delete order")
		}
		err = tx.Where("order_id = ?", orderNo).Updates(&model.OrderItem{
			IsDeleted: 1,
		}).Error

		if err != nil {
			ctx.Logger.Error("Failed to delete order item", err)
			tx.Rollback()
			return errors.New("failed to delete order item")
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (o *OrderServer) GetOrderDetail(ctx *app.Context, orderNo string) (*model.ReqOrderCreate, error) {
	order := &model.Order{}
	orderItems := &[]model.OrderItem{}

	err := ctx.DB.Where("order_no = ?", orderNo).First(order).Error
	if err != nil {
		ctx.Logger.Error("Failed to get order by orderNo", err)
		return nil, errors.New("failed to order by orderNo")
	}

	err = ctx.DB.Where("order_id = ?", orderNo).Find(orderItems).Error
	if err != nil {
		ctx.Logger.Error("Failed to get order item by orderNo", err)
		return nil, errors.New("failed to order item by orderNo")
	}
	orderResp := &model.ReqOrderCreate{
		Order:      *order,
		OrderItems: *orderItems,
	}
	return orderResp, nil
}

func (o *OrderServer) UpdateOrder(ctx *app.Context, reqOrder *model.ReqOrderUpdate) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	order := &model.Order{
		Status:    reqOrder.Status,
		UpdatedAt: now,
	}
	if reqOrder.Status == "2" {
		order.PaidAt = now
	}
	if reqOrder.Status == "3" {
		order.DeliveredAt = now
	}
	if reqOrder.Status == "4" {
		order.CompletedAt = now
	}
	if reqOrder.Status == "5" {
		order.ClosedAt = now
	}
	err := ctx.DB.Where("order_no = ?", reqOrder.OrderNo).Updates(order).Error
	if err != nil {
		ctx.Logger.Error("Failed to update order", err)
		return errors.New("failed to update order")
	}
	return nil
}
