package model

// 产品分类
type ProductCategory struct {
	ID   int64  `json:"id" gorm:"primary_key"`
	Uuid string `json:"uuid" gorm:"type:varchar(36);unique_index"`
	// 分类名称
	Name string `json:"name" gorm:"type:varchar(100)"`
	Icon string `json:"icon" gorm:"type:varchar(100)"` // 分类图标
	// 分类描述
	Description string `json:"description" gorm:"type:varchar(255)"`
	// 父级分类
	ParentUuid string `json:"parent_uuid" gorm:"type:varchar(36);index"`

	Sort int `json:"sort" gorm:"default:0"` // 排序
	// 状态 1:启用 2:禁用
	Status    int    `json:"status" gorm:"default:1"`          // 状态 1:启用 2:禁用
	CreatedAt string `gorm:"autoCreateTime" json:"created_at"` // CreatedAt 记录了创建的时间
	UpdatedAt string `gorm:"autoUpdateTime" json:"updated_at"` // UpdatedAt 记录了最后更新的时间
}
