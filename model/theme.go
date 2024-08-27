package model

type Theme struct {
	ID          int64  `json:"id" gorm:"primary_key"`                // 主题ID
	Uuid        string `json:"uuid" gorm:"type:varchar(36);unique"`  // 主题UUID
	Name        string `json:"name" gorm:"type:varchar(100)"`        // 主题名称
	Slug        string `json:"slug" gorm:"type:varchar(100)"`        // 主题别名/slug
	Version     string `json:"version" gorm:"type:varchar(10)"`      // 主题版本
	Author      string `json:"author" gorm:"type:varchar(100)"`      // 作者名称
	AuthorURI   string `json:"author_uri" gorm:"type:varchar(255)"`  // 作者链接
	Description string `json:"description" gorm:"type:text"`         // 主题描述
	Tags        string `json:"tags" gorm:"type:varchar(255)"`        // 主题标签，逗号分隔
	TextDomain  string `json:"text_domain" gorm:"type:varchar(100)"` // 主题语言域
	License     string `json:"license" gorm:"type:varchar(100)"`     // 主题许可协议
	LicenseURI  string `json:"license_uri" gorm:"type:varchar(255)"` // 许可协议链接
	Screenshot  string `json:"screenshot" gorm:"type:varchar(255)"`  // 主题截图链接
	Status      int    `json:"status" gorm:"type:int"`               // 主题状态，1为启用，0为禁用

	Options string `json:"options" gorm:"text"` // 主题选项（嵌入结构体）
}

// 主题功能选项结构体
type ThemeOptions struct {
	CustomLogo       bool `json:"custom_logo"`       // 是否支持自定义Logo
	CustomHeader     bool `json:"custom_header"`     // 是否支持自定义页眉
	CustomBackground bool `json:"custom_background"` // 是否支持自定义背景
	Menus            bool `json:"menus"`             // 是否支持菜单
	Widgets          bool `json:"widgets"`           // 是否支持小工具
	PostFormats      bool `json:"post_formats"`      // 是否支持文章格式
	ThemeCustomizer  bool `json:"theme_customizer"`  // 是否支持主题定制器
}

// 主题元信息
type ThemeMeta struct {
	ThemeUuid string `json:"theme_id" gorm:"index"`             // 关联的主题ID
	MetaKey   string `json:"meta_key" gorm:"type:varchar(100)"` // 元数据键
	MetaValue string `json:"meta_value" gorm:"type:text"`       // 元数据值
}
