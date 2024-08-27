package model

const (
	CurrencyStatusEnabled = 1 // 启用
)

type Currency struct {
	ID     int64  `json:"id" gorm:"primary_key"`
	Uuid   string `json:"uuid" gorm:"type:varchar(36);unique_index"` // 货币uuid
	Name   string `json:"name" gorm:"type:varchar(100)"`             // 货币名称
	Code   string `json:"code" gorm:"type:varchar(10)"`              // 货币代码
	Symbol string `json:"symbol" gorm:"type:varchar(10)"`            // 货币符号
	Status int    `json:"status"`                                    // 状态 1:启用 2:禁用
}

type ReqCurrencyCreate struct {
	Name   string `json:"name" binding:"required"`   // 货币名称
	Code   string `json:"code" binding:"required"`   // 货币代码
	Symbol string `json:"symbol" binding:"required"` // 货币符号
}

type ReqCurrencyQueryParam struct {
	Name string `json:"name"` // 货币名称
	Pagination
}
