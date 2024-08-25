package service

import (
	"errors"

	"sgin/model"
	"sgin/pkg/app"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentMethodService struct {
}

func NewPaymentMethodService() *PaymentMethodService {
	return &PaymentMethodService{}
}

// Create a new payment method
func (s *PaymentMethodService) CreatePaymentMethod(ctx *app.Context, paymentMethod *model.PaymentMethod) error {

	paymentMethod.Uuid = uuid.New().String()

	err := ctx.DB.Create(paymentMethod).Error
	if err != nil {
		ctx.Logger.Error("Failed to create payment method", err)
		return errors.New("failed to create payment method")
	}
	return nil
}

// Get a payment method by UUID
func (s *PaymentMethodService) GetPaymentMethodInfo(ctx *app.Context, uuid string, code string) (*model.PaymentMethod, error) {
	paymentMethod := &model.PaymentMethod{}
	err := ctx.DB.Where("uuid = ? OR code  = ?", uuid, code).First(paymentMethod).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("payment method not found")
		}
		ctx.Logger.Error("Failed to get payment method by UUID", err)
		return nil, errors.New("failed to get payment method by UUID")
	}
	return paymentMethod, nil
}

// 根据code 获取支付方式
func (s *PaymentMethodService) GetPaymentMethodByCode(ctx *app.Context, code string) (*model.PaymentMethod, error) {
	paymentMethod := &model.PaymentMethod{}
	err := ctx.DB.Where("code = ?", code).First(paymentMethod).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("payment method not found")
		}
		ctx.Logger.Error("Failed to get payment method by code", err)
		return nil, errors.New("failed to get payment method by code")
	}
	return paymentMethod, nil
}

// Update an existing payment method
func (s *PaymentMethodService) UpdatePaymentMethod(ctx *app.Context, paymentMethod *model.PaymentMethod) error {
	err := ctx.DB.Save(paymentMethod).Error
	if err != nil {
		ctx.Logger.Error("Failed to update payment method", err)
		return errors.New("failed to update payment method")
	}
	return nil
}

// 更新状态
func (s *PaymentMethodService) UpdatePaymentMethodStatus(ctx *app.Context, uuid string, status int) error {
	err := ctx.DB.Model(&model.PaymentMethod{}).Where("uuid = ?", uuid).Update("status", status).Error
	if err != nil {
		ctx.Logger.Error("Failed to update payment method status", err)
		return errors.New("failed to update payment method status")
	}
	return nil
}

// 更新配置
func (s *PaymentMethodService) UpdatePaymentMethodConfig(ctx *app.Context, uuid string, config string) error {
	err := ctx.DB.Model(&model.PaymentMethod{}).Where("uuid = ?", uuid).Update("config", config).Error
	if err != nil {
		ctx.Logger.Error("Failed to update payment method config", err)
		return errors.New("failed to update payment method config")
	}
	return nil
}

// Delete a payment method by UUID
func (s *PaymentMethodService) DeletePaymentMethod(ctx *app.Context, uuid string) error {
	err := ctx.DB.Where("uuid = ?", uuid).Delete(&model.PaymentMethod{}).Error
	if err != nil {
		ctx.Logger.Error("Failed to delete payment method", err)
		return errors.New("failed to delete payment method")
	}
	return nil
}

// Get a list of payment methods with optional filters
func (s *PaymentMethodService) GetPaymentMethodList(ctx *app.Context, params *model.ReqPaymentMethodQueryParam) (r *model.PagedResponse, err error) {
	var (
		paymentMethods []*model.PaymentMethod
		total          int64
	)

	query := ctx.DB.Model(&model.PaymentMethod{})

	if params.Name != "" {
		query = query.Where("name LIKE ?", "%"+params.Name+"%")
	}

	if params.Code != "" {
		query = query.Where("code LIKE ?", "%"+params.Code+"%")
	}

	if params.Status != 0 {
		query = query.Where("status = ?", params.Status)
	}

	err = query.Count(&total).Error
	if err != nil {
		ctx.Logger.Error("Failed to count payment methods", err)
		return nil, errors.New("failed to count payment methods")
	}

	err = query.Offset(params.GetOffset()).Limit(params.PageSize).Find(&paymentMethods).Error
	if err != nil {
		ctx.Logger.Error("Failed to get payment method list", err)
		return nil, errors.New("failed to get payment method list")
	}

	res := make([]*model.PaymentMethodRes, 0)
	for _, v := range paymentMethods {

		isConfig := false
		if v.Config != "" {
			isConfig = true
		}
		v.Config = ""

		item := &model.PaymentMethodRes{
			PaymentMethod: *v,
			IsConfig:      isConfig,
		}
		res = append(res, item)
	}

	return &model.PagedResponse{
		Total:    total,
		Data:     res,
		Current:  params.Current,
		PageSize: params.PageSize,
	}, nil
}

// 获取所有可用的支付方式
func (s *PaymentMethodService) GetAvailablePaymentMethodList(ctx *app.Context) ([]*model.PaymentMethod, error) {
	var paymentMethods []*model.PaymentMethod
	err := ctx.DB.Where("status = ?", 1).Find(&paymentMethods).Error
	if err != nil {
		ctx.Logger.Error("Failed to get available payment methods", err)
		return nil, errors.New("failed to get available payment methods")
	}

	for _, v := range paymentMethods {
		v.Config = ""
	}

	return paymentMethods, nil
}
