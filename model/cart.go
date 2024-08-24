package model

// Cart 购物车
type Cart struct {
	ID   int64  `json:"id" gorm:"primary_key"`
	Uuid string `json:"uuid" gorm:"type:varchar(36);unique_index"`
	// 用户ID
	UserID string `json:"user_id" gorm:"index"`
	// 产品ID
	ProductItemUuid string `json:"product_item_uuid" gorm:"type:varchar(36);index"`
	// 数量
	Quantity  int    `json:"quantity"`
	CreatedAt string `gorm:"autoCreateTime" json:"created_at"` // CreatedAt 记录了创建的时间
	UpdatedAt string `gorm:"autoUpdateTime" json:"updated_at"` // UpdatedAt 记录了最后更新的时间
}

type CartProductItemRes struct {
	Cart
	ProductItem *ProductItemRes `json:"product_item"`
}

type ReqCartQueryParam struct {
	// 用户ID
	UserID string `json:"user_id"`
	Pagination
}
