package model

const (
	// 站点配置
	ConfigCategorySite = "site"
	// 邮件配置
	ConfigCategoryEmail = "email"
)

const (
	// 站点url
	ConfigNameSiteUrl = "site_url"
	// 站点标题
	ConfigNameSiteTitle = "site_title"
	// 站点副标题
	ConfigNameSiteSubTitle = "site_sub_title"
	// 站点描述
	ConfigNameSiteDescription = "site_description"
	// 站点关键字
	ConfigNameSiteKeyword = "site_keyword"
	// 站点版权信息
	ConfigNameSiteCopyRight = "site_copy_right"
	// 站点备案号
	ConfigNameSiteRecordNo = "site_record_no"
	// 站点备案链接
	ConfigNameSiteRecordUrl = "site_record_url"
	// 站点logo
	ConfigNameSiteLogo = "site_logo"
	// 站点favicon
	ConfigNameSiteFavicon = "site_favicon"
	// 站点语言
	ConfigNameSiteLanguage = "site_language"
)

const (
	// SMTP 服务器
	ConfigNameEmailSmtpHost = "email_smtp_host"
	// SMTP 端口
	ConfigNameEmailSmtpPort = "email_smtp_port"
	// SMTP 发件人
	ConfigNameEmailSmtpUser = "email_smtp_user"
	// SMTP 密码
	ConfigNameEmailSmtpPass = "email_smtp_pass"
)

type Configuration struct {
	Id        int    `json:"id"`
	Category  string `json:"category"`
	Name      string `json:"name"`
	Value     string `json:"value"`
	CreatedAt string `json:"created_at" gorm:"autoCreateTime;comment:'创建时间'"` // 创建时间
	UpdatedAt string `json:"updated_at" gorm:"autoUpdateTime;comment:'更新时间'"` // 更新时间
}

type ReqConfigQueryParam struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Pagination
}
