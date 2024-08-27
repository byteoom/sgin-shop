package model

const (
	ResourceTypeFolder = "folder" // 文件夹
	ResourceTypeFile   = "file"   // 文件
)

// 资源
type Resource struct {
	ID   int64  `json:"id" gorm:"primary_key"`
	Uuid string `json:"uuid" gorm:"type:varchar(36);unique_index"`
	// 资源名称
	Name string `json:"name" gorm:"type:varchar(100)"`
	// 资源描述
	Description string `json:"description" gorm:"type:varchar(255)"`
	// 父级资源
	ParentUuid string `json:"parent_uuid" gorm:"type:varchar(36);index"`

	Md5 string `json:"md5" gorm:"type:varchar(70);index"` // 文件MD5

	// 资源类型
	Type string `json:"type" gorm:"type:varchar(100)"` // 资源类型 文件夹、文件 等

	// 文件类型
	MimeType string `json:"mime_type" gorm:"type:varchar(100)"` // 文件类型

	// 文件大小
	Size int64 `json:"size" gorm:"default:0"` // 文件大小

	// 文件路径
	Path string `json:"path" gorm:"type:varchar(255)"`

	// 文件路径
	Address string `json:"address" gorm:"type:varchar(255)"` // 文件路径

	CreatedAt string `gorm:"autoCreateTime" json:"created_at"` // CreatedAt 记录了创建的时间
	UpdatedAt string `gorm:"autoUpdateTime" json:"updated_at"` // UpdatedAt 记录了最后更新的时间
}

// 创建文件夹
type ReqResourceCreateFolder struct {
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
	Path        string `json:"path" form:"path"`
}

// 移动资源
type ReqResourceMove struct {
	UuidList   []string `json:"uuid_list" form:"uuid_list"`     // 资源uuid
	ParentUuid string   `json:"parent_uuid" form:"parent_uuid"` // 父级资源
}

type ReqResourceQueryParam struct {
	// 资源名称
	Path       string `json:"path" form:"path"` // 资源路径
	Name       string `json:"name" form:"name"`
	ParentUuid string `json:"parent_uuid" form:"parent_uuid"` // 父级资源
	MimeType   string `json:"mime_type" form:"mime_type"`     // 文件类型
	Pagination
}

type ResourceRes struct {
	Resource
	Children []*ResourceRes `json:"children"`
}

type ReqResourceDeleteParam struct {
	Uuid string `json:"uuid" form:"uuid"`
}
