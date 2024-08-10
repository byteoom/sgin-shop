package model

// 产品
type Product struct {
	ID   int64  `json:"id" gorm:"primary_key"`
	Uuid string `json:"uuid" gorm:"type:varchar(36);unique_index"`
	// 产品分类
	ProductCategoryUuid string `json:"product_category_uuid" gorm:"type:varchar(36);index"`
	// 产品名称
	Name string `json:"name" gorm:"type:varchar(100)"`
	// 产品描述
	Description string `json:"description" gorm:"type:varchar(255)"`

	Images string `json:"images" gorm:"comment:产品图片"`
	// 产品视频
	Videos string `json:"videos" gorm:"comment:产品视频"`

	CreatedAt string `gorm:"autoCreateTime" json:"created_at"` // CreatedAt 记录了创建的时间
	UpdatedAt string `gorm:"autoUpdateTime" json:"updated_at"` // UpdatedAt 记录了最后更新的时间

}

// 产品变体
type ProductVariants struct {
	ID        int64  `json:"id" gorm:"primary_key"`
	Uuid      string `json:"uuid" gorm:"type:varchar(36);unique_index"`
	ProductID int64  `json:"product_id" gorm:"index"`
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
	Variants []ReqProductVariantsCreate `json:"variants" binding:"-"`
}

type ReqProductVariantsCreate struct {
	ReqProdcutItemCommonCreate
	// 产品变体名称
	Name string `json:"name" binding:"required"`
	// 产品变体描述
	Description string `json:"description" binding:"required"`
	// 产品变体Option
	Options []ReqProductVariantsOptionCreate `json:"options" binding:"-"`
}

type ReqProductVariantsOptionCreate struct {
	ReqProdcutItemCommonCreate
	// 产品变体Option名称
	Name string `json:"name" binding:"required"`
	// 单位 例如: 个、件、套、箱
	Unit string `json:"unit" binding:"required"`
	// 产品变体Option描述
	Description string `json:"description" binding:"required"`

	// 产品变体Option值
	Values []ReqProductVariantsOptionValueCreate `json:"values" binding:"-"`
}

type ReqProductVariantsOptionValueCreate struct {
	ReqProdcutItemCommonCreate
	// 产品变体Option值名称
	Name string `json:"name" binding:"required"`
	// 单位 例如: 个、件、套、箱
	Unit string `json:"unit" binding:"required"`
	// 产品变体Option值描述
	Description string `json:"description" binding:"-"`
}
