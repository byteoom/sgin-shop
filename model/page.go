package model

const (
	// PageStatusDraft 页面状态为草稿
	PageStatusDraft = "draft"
	// PageStatusPublished 页面状态为已发布
	PageStatusPublished = "published"
)

type Page struct {
	ID        uint   `gorm:"primary_key" json:"id"`                         // ID 是页面的主键
	UUID      string `gorm:"type:char(36);index" json:"uuid"`               // UUID 是页面的唯一标识符
	Title     string `gorm:"type:varchar(100)" json:"title"`                // Title 是页面的标题
	Slug      string `gorm:"type:varchar(255);unique" json:"slug"`          // Slug 是页面的路径 (URL)
	Status    string `json:"status" gorm:"column:status;type:varchar(100)"` // Status 是页面的状态 (如 "draft", "published")
	Data      string `gorm:"type:longtext" json:"data"`                     // Data 存储页面的详细内容，较大字符串
	Ext       string `gorm:"type:longtext" json:"ext"`                      // Ext 存储页面的扩展内容，较大字符串
	CreatedAt string `gorm:"autoCreateTime" json:"created_at"`              // CreatedAt 页面创建时间
	UpdatedAt string `gorm:"autoUpdateTime" json:"updated_at"`              // UpdatedAt 页面最后更新时间
}

type ReqPageCreate struct {
	Title  string `json:"title" binding:"required"` // Title 是页面的标题
	Slug   string `json:"slug" binding:"-"`         // Slug 是页面的路径 (URL)
	Status string `json:"status" binding:"-"`       // Status 是页面的状态 (如 "draft", "published")
	Data   string `json:"data" binding:"-"`         // Data 存储页面的详细内容，较大字符串
	Ext    string `json:"ext" binding:"-"`          // Ext 存储页面的扩展内容，较大字符串
}

type ReqPageQueryParam struct {
	Title string `json:"title"` // Title 是页面的标题
	Pagination
}
