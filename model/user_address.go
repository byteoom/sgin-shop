package model

type UserAddress struct {
	ID     uint   `json:"id" gorm:"primary_key"`
	Uuid   string `json:"uuid" gorm:"type:char(36);unique" ` // 用户地址唯一标识
	UserID uint   `json:"user_id" gorm:"index"`              // 用户ID
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

	// 备注
	Remark string `json:"remark" gorm:"type:varchar(255)"`
	// 地址类型
	AddressType string `json:"address_type" gorm:"type:varchar(100)"` // 公司、家庭、学校 等
	// 是否默认地址
	IsDefault bool   `json:"is_default" gorm:"default:false"`  // 是否默认地址
	CreatedAt string `gorm:"autoCreateTime" json:"created_at"` // CreatedAt 记录了创建的时间
	UpdatedAt string `gorm:"autoUpdateTime" json:"updated_at"` // UpdatedAt 记录了最后更新的时间
}
