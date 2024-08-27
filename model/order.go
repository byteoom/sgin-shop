package model

const (
	// 订单状态
	OrderStatusPending   = "pending"   // 待支付
	OrderStatusPaid      = "paid"      // 已支付
	OrderStatusDelivered = "delivered" // 已发货
	OrderStatusCompleted = "completed" // 已完成
	OrderStatusClosed    = "closed"    // 已关闭
)

// 订单
type Order struct {
	ID int64 `json:"id" gorm:"primary_key"`
	// 订单编号
	OrderNo string `json:"order_no" gorm:"type:varchar(100);unique_index"`
	// 用户ID
	UserID string `json:"user_id" gorm:"index"`
	// 订单总金额
	TotalAmount float64 `json:"total_amount"`
	// 订单状态 1:待支付 2:已支付 3:已发货 4:已完成 5:已关闭
	Status string `json:"status" gorm:"default:1"`

	// 收货人姓名
	ReceiverName string `json:"receiver_name" gorm:"type:varchar(100)"`
	// 收货人电话
	ReceiverPhone string `json:"receiver_phone" gorm:"type:varchar(20)"`

	// 收货人邮箱
	ReceiverEmail string `json:"receiver_email" gorm:"type:varchar(100)"`

	// 收货人国家
	ReceiverCountry string `json:"receiver_country" gorm:"type:varchar(100)"`
	// 收货人省份
	ReceiverProvince string `json:"receiver_province" gorm:"type:varchar(100)"`
	// 收货人城市
	ReceiverCity string `json:"receiver_city" gorm:"type:varchar(100)"`

	// 收货人地址
	ReceiverAddress string `json:"receiver_address" gorm:"type:varchar(255)"`
	// 收货人邮编
	ReceiverZip string `json:"receiver_zip" gorm:"type:varchar(10)"`
	// 收货人备注
	ReceiverRemark string `json:"receiver_remark" gorm:"type:varchar(255)"`

	// 支付时间
	PaidAt string `json:"paid_at"`
	// 发货时间
	DeliveredAt string `json:"delivered_at"`
	// 完成时间
	CompletedAt string `json:"completed_at"`
	// 关闭时间
	ClosedAt  string `json:"closed_at"`
	CreatedAt string `gorm:"autoCreateTime" json:"created_at"` // CreatedAt 记录了创建的时间
	UpdatedAt string `gorm:"autoUpdateTime" json:"updated_at"` // UpdatedAt 记录了最后更新的时间
}

type OrderRes struct {
	Order
	Items []*OrderItemRes `json:"items"` // 订单商品
}

// 收货人信息
type OrderReceiver struct {
	// 收货人姓名
	ReceiverName string `json:"receiver_name" gorm:"type:varchar(100)"`
	// 收货人电话
	ReceiverPhone string `json:"receiver_phone" gorm:"type:varchar(20)"`

	// 收货人邮箱
	ReceiverEmail string `json:"receiver_email" gorm:"type:varchar(100)"`

	// 收货人国家
	ReceiverCountry string `json:"receiver_country" gorm:"type:varchar(100)"`
	// 收货人省份
	ReceiverProvince string `json:"receiver_province" gorm:"type:varchar(100)"`
	// 收货人城市
	ReceiverCity string `json:"receiver_city" gorm:"type:varchar(100)"`

	// 收货人地址
	ReceiverAddress string `json:"receiver_address" gorm:"type:varchar(255)"`
	// 收货人邮编
	ReceiverZip string `json:"receiver_zip" gorm:"type:varchar(10)"`
	// 收货人备注
	ReceiverRemark string `json:"receiver_remark" gorm:"type:varchar(255)"`
}

// 订单商品
type OrderItem struct {
	ID int64 `json:"id" gorm:"primary_key"`
	// 订单ID
	OrderID string `json:"order_id" gorm:"index"`
	// 商品ID
	ProductItemID string `json:"product_item_id" gorm:"index"`
	// 商品数量
	Quantity int `json:"quantity"`
	// 商品单价
	Price float64 `json:"price"`
	// 商品总价
	TotalAmount float64 `json:"total_amount"`

	// 折扣金额
	DiscountAmount float64 `json:"discount_amount"`
	// 折扣
	Discount float64 `json:"discount"`
	// 折扣价
	DiscountPrice float64 `json:"discount_price"`
	CreatedAt     string  `gorm:"autoCreateTime" json:"created_at"` // CreatedAt 记录了创建的时间
	UpdatedAt     string  `gorm:"autoUpdateTime" json:"updated_at"` // UpdatedAt 记录了最后更新的时间
}

type OrderItemRes struct {
	OrderItem
	ProductItem *ProductItemRes `json:"product_item"`
}

type ReqOrderCreate struct {
	UserId   string        `json:"user_id"`
	Receiver OrderReceiver `json:"receiver"`

	Items []ReqOrderItemCreate `json:"items"`

	CartUuids []string `json:"cart_uuids"` // 购物车ID列表
}

type ReqOrderItemCreate struct {
	ProductItemID string `json:"product_item_id"`
	Quantity      int    `json:"quantity"`
}

type ReqOrderQueryParam struct {
	UserID  string `json:"user_id"`  // 用户ID，用于过滤
	Status  string `json:"status"`   // 订单状态，用于过滤
	OrderNo string `json:"order_no"` // 订单编号，用于过滤
	Pagination
}
