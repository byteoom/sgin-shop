package service

import (
	"errors"
	"gorm.io/gorm"
	"time"

	"github.com/google/uuid"
	"sgin/model"
	"sgin/pkg/app"
)

type PaymentService struct {
}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

// CreatePayment 创建一个新的付款记录
func (s *PaymentService) CreatePayment(ctx *app.Context, payment *model.Payment) (*model.Payment, error) {
	payment.Uuid = uuid.New().String()
	payment.PaidAt = time.Now().Format(time.RFC3339)
	payment.CreatedAt = time.Now().Format(time.RFC3339)
	payment.UpdatedAt = payment.CreatedAt
	err := ctx.DB.Create(&payment).Error
	if err != nil {
		ctx.Logger.Error("Failed to create payment", err)
		return nil, errors.New("failed to create payment")
	}

	return payment, nil
}

// GetPaymentByUUID 根据UUID获取付款记录
func (s *PaymentService) GetPaymentByUUID(ctx *app.Context, uuid string) (*model.Payment, error) {
	var payment model.Payment
	err := ctx.DB.Where("uuid = ?", uuid).First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		ctx.Logger.Error("Failed to get payment by UUID", err)
		return nil, errors.New("failed to get payment by UUID")
	}
	return &payment, nil
}

// UpdatePayment 更新付款记录
func (s *PaymentService) UpdatePayment(ctx *app.Context, payment **model.Payment) error {
	now := time.Now()
	(*payment).UpdatedAt = now.Format(time.RFC3339)
	err := ctx.DB.Updates(payment).Error
	if err != nil {
		ctx.Logger.Error("Failed to update payment", err)
		return errors.New("failed to update payment")
	}

	return nil
}

// DeletePayment 删除付款记录
func (s *PaymentService) DeletePayment(ctx *app.Context, uuid string) error {
	err := ctx.DB.Where("uuid = ?", uuid).Delete(&model.Payment{}).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete payment", err)
		return errors.New("failed to delete payment")
	}

	return nil
}

// GetPaymentList 获取付款记录列表
func (s *PaymentService) GetPaymentList(ctx *app.Context, params *model.ReqPaymentQueryParam) (*model.PagedResponse, error) {
	var (
		payments []*model.Payment
		total    int64
	)
	db := ctx.DB.Model(&model.Payment{})
	// 应用UserID和OrderID过滤条件
	if params != nil {
		if params.UserID != 0 {
			db = db.Where("user_id = ?", params.UserID)
		}
		if params.OrderID != 0 {
			db = db.Where("order_id = ?", params.OrderID)
		}
	}
	// 计算偏移量
	offset := (params.Page - 1) * params.PageSize
	err := db.Count(&total).Error
	if err != nil {
		ctx.Logger.Error("Failed to get payment count", err)
		return nil, errors.New("failed to get payment count")
	}
	err = db.Offset(offset).Limit(params.PageSize).Find(&payments).Error
	if err != nil {
		ctx.Logger.Error("Failed to get payment list", err)
		return nil, errors.New("failed to get payment list")
	}

	return &model.PagedResponse{
		Total: total,
		Data:  payments,
	}, nil
}
