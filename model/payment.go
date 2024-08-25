package model

const (
	PaymentStatusPending = "pending" // 待支付
	PaymentStatusPaid    = "paid"    // 已支付
	// 取消支付
	PaymentStatusCanceled = "canceled"
)

// 付款信息
type Payment struct {
	ID   int64  `json:"id" gorm:"primary_key"`
	Uuid string `json:"uuid" gorm:"type:varchar(36);unique_index"`

	// 用户ID
	UserID string `json:"user_id" gorm:"index"`
	// 订单ID
	OrderID string `json:"order_id" gorm:"index"`
	// 付款金额
	Amount float64 `json:"amount"`
	// 付款状态  待支付 已经支付
	Status string `json:"status"`

	// 付款方式
	Method string `json:"method"`
	// 付款渠道
	Channel string `json:"channel"`
	// 付款渠道订单号
	ChannelOrderNo string `json:"channel_order_no"`
	// 付款渠道交易号
	ChannelTransactionNo string `json:"channel_transaction_no"`

	ChannelStatus string `json:"channel_status"` // 付款渠道状态

	ChannelData string `json:"channel_data"` // 付款渠道数据

	// 付款时间
	PaidAt    string `json:"paid_at"`
	CreatedAt string `gorm:"autoCreateTime" json:"created_at"` // CreatedAt 记录了创建的时间
	UpdatedAt string `gorm:"autoUpdateTime" json:"updated_at"` // UpdatedAt 记录了最后更新的时间
}

type ReqPaymentQueryParam struct {
	UserID  string `json:"user_id"`  // 用户ID，用于过滤
	OrderID string `json:"order_id"` // 订单ID，用于过滤
	Pagination
}
