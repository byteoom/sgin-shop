package model

// 支付方式
type PaymentMethod struct {
	ID   int64  `json:"id" gorm:"primary_key"`
	Uuid string `json:"uuid" gorm:"type:varchar(36);unique_index"`

	// 支付方式名称
	Name string `json:"name" gorm:"type:varchar(100)"`

	// 支付方式code
	Code string `json:"code" gorm:"type:varchar(100);unique_index"`
	// 支付方式描述
	Description string `json:"description" gorm:"type:varchar(255)"`

	// 支付方式图标
	Icon string `json:"icon" gorm:"type:varchar(255)"`

	// 支付方式url
	Url string `json:"url" gorm:"type:varchar(255)"`

	// 支付方式状态
	Status int `json:"status"` // 1: 启用 2: 禁用

	// Config 支付方式配置
	Config string `json:"config" gorm:"type:text"`
}

type ReqPaymentMethodQueryParam struct {
	Name   string `json:"name"`   // 支付方式名称，用于过滤
	Code   string `json:"code"`   // 支付方式code，用于过滤
	Status int    `json:"status"` // 支付方式状态，用于过滤
	Pagination
}

type ReqPaymentMethodInfoParam struct {
	Uuid string `json:"uuid" binding:"-"` // 支付方式uuid
	Code string `json:"code" binding:"-"` // 支付方式code
}

type PaymentMethodRes struct {
	PaymentMethod
	// 是否已经配置
	IsConfig bool `json:"is_config"`
}

type ReqPaypalClientIdParam struct {
	Env string `json:"env" binding:"required"` // 环境
}

type ReqPaymentOrderCreateParam struct {
	OrderID string `json:"order_id" binding:"required"` // 订单ID
}

type ReqPaypalOrderCreateParam struct {
	Name     string  `json:"productName" binding:"required"` // 名称
	Amount   float64 `json:"amount" binding:"required"`      // 金额
	Currency string  `json:"currency" binding:"-"`           // 货币
}
