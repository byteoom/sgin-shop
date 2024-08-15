package model

type ProductBase struct {
	// 重量
	Weight float64 `json:"weight" gorm:"type:decimal(10,2)"` // 重量
	// 长度
	Length float64 `json:"length" gorm:"type:decimal(10,2)"` // 长度
	// 宽度
	Width float64 `json:"width" gorm:"type:decimal(10,2)"` // 宽度
	// 高度
	Height float64 `json:"height" gorm:"type:decimal(10,2)"` // 高度

	Unit string `json:"unit" gorm:"type:varchar(100)"` // 单位 例如: 个、件、套、箱
}

// 产品
type Product struct {
	ProductBase
	ID   int64  `json:"id" gorm:"primary_key"`
	Uuid string `json:"uuid" gorm:"type:varchar(36);unique_index"`
	// 产品分类
	ProductCategoryUuid string `json:"product_category_uuid" gorm:"type:varchar(36);index"`
	// 产品名称
	Name string `json:"name" gorm:"type:varchar(100)"`
	// 产品描述
	Description string `json:"description" gorm:"type:varchar(255)"`

	// 产品类型
	Type string `json:"type" gorm:"type:varchar(100)"` // 产品类型 全新、二手、虚拟产品

	// 产品状态
	Status string `json:"status" gorm:"type:varchar(100)"` // 产品状态 上架、下架、售罄

	// 产品警戒库存
	StockWarning int64 `json:"stock_warning" gorm:"type:int"` // 产品警戒库存

	// 低于警戒库存是否可售
	StockWarningSell bool `json:"stock_warning_sell" gorm:"type:bool"` // 低于警戒库存是否可售

	Images string `json:"images" gorm:"comment:产品图片"`
	// 产品视频
	Videos string `json:"videos" gorm:"comment:产品视频"`

	CreatedAt string `gorm:"autoCreateTime" json:"created_at"` // CreatedAt 记录了创建的时间
	UpdatedAt string `gorm:"autoUpdateTime" json:"updated_at"` // UpdatedAt 记录了最后更新的时间

}

type ProductRes struct {
	Product
	ImageList       []string        `json:"image_list"`
	ProductCategory ProductCategory `json:"product_category"`
}

// 产品变体
type ProductVariants struct {
	ID          int64  `json:"id" gorm:"primary_key"`
	Uuid        string `json:"uuid" gorm:"type:varchar(36);unique_index"`
	ProductUuid string `json:"product_uuid" gorm:"index"`
	// 产品变体名称
	Name string `json:"name" gorm:"type:varchar(100)"`
	// 产品变体描述
	Description string `json:"description" gorm:"type:varchar(255)"`
	CreatedAt   string `gorm:"autoCreateTime" json:"created_at"` // CreatedAt 记录了创建的时间
	UpdatedAt   string `gorm:"autoUpdateTime" json:"updated_at"` // UpdatedAt 记录了最后更新的时间
}

// 产品变体Option
type ProductVariantsOption struct {
	ID                  int64  `json:"id" gorm:"primary_key"`
	Uuid                string `json:"uuid" gorm:"type:varchar(36);unique_index"`
	ProductUuid         string `json:"product_uuid" gorm:"index"`
	ProductVariantsUuid string `json:"product_variants_uuid" gorm:"type:varchar(36);index"`
	// 产品变体Option名称
	Name string `json:"name" gorm:"type:varchar(100)"`
	// 单位 例如: 个、件、套、箱
	Unit string `json:"unit" gorm:"type:varchar(100)"`
	// 产品变体Option描述
	Description string `json:"description" gorm:"type:varchar(255)"`
	CreatedAt   string `gorm:"autoCreateTime" json:"created_at"` // CreatedAt 记录了创建的时间
	UpdatedAt   string `gorm:"autoUpdateTime" json:"updated_at"` // UpdatedAt 记录了最后更新的时间
}

// 产品变体Option值
type ProductVariantsOptionValue struct {
	ID                        int64  `json:"id" gorm:"primary_key"`
	Uuid                      string `json:"uuid" gorm:"type:varchar(36);unique_index"`
	ProductVariantsOptionUuid string `json:"product_variants_option_uuid" gorm:"type:varchar(36);index"`
	// 产品变体Option值名称
	Name string `json:"name" gorm:"type:varchar(100)"`
	// 单位 例如: 个、件、套、箱
	Unit string `json:"unit" gorm:"type:varchar(100)"`
	// 产品变体Option值描述
	Description string `json:"description" gorm:"type:varchar(255)"`
	CreatedAt   string `gorm:"autoCreateTime" json:"created_at"` // CreatedAt 记录了创建的时间
	UpdatedAt   string `gorm:"autoUpdateTime" json:"updated_at"` // UpdatedAt 记录了最后更新的时间
}

// 产品具体信息
type ProductItem struct {
	ProductBase
	ID   int64  `json:"id" gorm:"primary_key"`
	Uuid string `json:"uuid" gorm:"type:varchar(36);unique_index"`
	// 产品名称
	Name string `json:"name" gorm:"type:varchar(100)"`
	// 产品uuid
	ProductUuid string `json:"product_uuid" gorm:"type:varchar(36);index"`
	// 产品变体uuid
	ProductVariantsUuid string `json:"product_variants_uuid" gorm:"type:varchar(36);index"`
	// 变体optionuuid
	ProductVariantsOptionUuid string `json:"product_variants_option_uuid" gorm:"type:varchar(36);index"`

	//变体option值uuid
	ProductVariantsOptionValueUuid string `json:"product_variants_option_value_uuid" gorm:"type:varchar(36);index"`

	//  产品变体
	Variants string `json:"variants" gorm:"type:text"`

	Images string `json:"images" gorm:"comment:产品图片"`
	// 产品视频
	Videos string `json:"videos" gorm:"comment:产品视频"`
	// 产品描述
	Description string `json:"description" gorm:"type:varchar(255)"`
	// 产品价格
	Price float64 `json:"price" gorm:"type:decimal(10,2)"`
	// 产品折扣
	Discount float64 `json:"discount" gorm:"type:decimal(10,2)"`
	// 产品折扣价
	DiscountPrice float64 `json:"discount_price" gorm:"type:decimal(10,2)"`
	// 产品库存
	Stock int64 `json:"stock" gorm:"type:int"`

	CreatedAt string `gorm:"autoCreateTime" json:"created_at"` // CreatedAt 记录了创建的时间
	UpdatedAt string `gorm:"autoUpdateTime" json:"updated_at"` // UpdatedAt 记录了最后更新的时间
}

type ProductItemRes struct {
	ProductItem
	ImageList   []string    `json:"image_list"`
	ProductInfo *ProductRes `json:"product_info"`
}

type ReqProdcutItemCommonCreate struct {
	// 产品价格
	Price float64 `json:"price" binding:"-"`
	// 产品折扣
	Discount float64 `json:"discount" binding:"-"`
	// 产品折扣价
	DiscountPrice float64 `json:"discount_price" binding:"-"`
	// 产品库存
	Stock int64 `json:"stock" binding:"-"`
}

// 产品创建
type ReqProductCreate struct {
	ReqProdcutItemCommonCreate
	// 产品名称
	Name string `json:"name" binding:"required"`
	// 产品分类
	ProductCategoryUuid string `json:"product_category_uuid" binding:"-"`

	// 产品描述
	Description string `json:"description" binding:"-"`

	Images []string `json:"images" binding:"-"`

	// 产品视频
	Videos []string `json:"videos" binding:"-"`

	// 产品变体
	Variants []*ReqProductVariantsCreate `json:"variants" binding:"-"`

	// 产品变体值
	VariantsVals []map[string]interface{} `json:"variants_vals" binding:"-"`

	Unit string `json:"unit" binding:"-"` // 单位 例如: 个、件、套、箱

	// 产品上架状态
	// 上架、下架、售罄
	ProductStatus string `json:"product_status" binding:"-"` // 产品状态 上架、下架、售罄

	// 预警库存
	StockWarning int64 `json:"stock_warning" binding:"-"`

	// 低于警戒库存是否可售
	StockWarningSell bool `json:"stock_warning_sell" binding:"-"`

	// 长度
	Length float64 `json:"length" binding:"-"` // 长度
	// 宽度
	Width float64 `json:"width" binding:"-"` // 宽度

	// 高度
	Height float64 `json:"height" binding:"-"` // 高度

	// 重量
	Weight float64 `json:"weight" binding:"-"` // 重量
}

type ReqProductVariantsCreate struct {
	ReqProdcutItemCommonCreate
	// 产品变体名称
	Name string `json:"name" binding:"required"`
	// 产品变体描述
	Description string `json:"description" binding:"required"`
	// 产品变体Option
	Options []string `json:"options" binding:"-"`
}

type ReqProductVariantsOptionCreate struct {
	ReqProdcutItemCommonCreate
	// 产品变体Option名称
	Name string `json:"name" binding:"required"`
	// 单位 例如: 个、件、套、箱
	Unit string `json:"unit" binding:"required"`
	// 产品变体Option描述
	Description string `json:"description" binding:"required"`
}
