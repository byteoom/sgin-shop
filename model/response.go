package model

type PagedResponse struct {
	Data     interface{} `json:"data"`
	Current  int         `json:"current"`
	PageSize int         `json:"pageSize"`
	Total    int64       `json:"total"`
}

type ResUserLogin struct {
	Token string `json:"token"`
}

type BaseResponse struct {
	TraceID string `json:"trace_id"` // 请求唯一标识
	Code    int    `json:"code"`     // 状态码
	Message string `json:"message"`  // 提示信息
}

type StringDataResponse struct {
	BaseResponse
	Data string `json:"data"`
}

type UserInfoResponse struct {
	BaseResponse
	Data User `json:"data"`
}

type BasePageResponse struct {
	Current  int   `json:"current"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
}

type UserPageResponse struct {
	BasePageResponse
	Data []User `json:"data"`
}

type UserQueryResponse struct {
	BaseResponse
	Data UserPageResponse `json:"data"`
}

type AppInfoResponse struct {
	BaseResponse
	Data App `json:"data"`
}

type AppPageResponse struct {
	BasePageResponse
	Data []App `json:"data"`
}

type AppQueryResponse struct {
	BaseResponse
	Data AppPageResponse `json:"data"`
}

type MenuInfoResponse struct {
	BaseResponse
	Data Menu `json:"data"`
}

type MenuPageResponse struct {
	BasePageResponse
	Data []Menu `json:"data"`
}

type MenuQueryResponse struct {
	BaseResponse
	Data MenuPageResponse `json:"data"`
}

type RoleInfoResponse struct {
	BaseResponse
	Data Role `json:"data"`
}

type RolePageResponse struct {
	BasePageResponse
	Data []Role `json:"data"`
}

type RoleQueryResponse struct {
	BaseResponse
	Data RolePageResponse `json:"data"`
}

type ServerInfoResponse struct {
	BaseResponse
	Data Server `json:"data"`
}

type ServerPageResponse struct {
	BasePageResponse
	Data []Server `json:"data"`
}

type ServerQueryResponse struct {
	BaseResponse
	Data ServerPageResponse `json:"data"`
}

type TeamInfoResponse struct {
	BaseResponse
	Data Team `json:"data"`
}

type TeamPageResponse struct {
	BasePageResponse
	Data []Team `json:"data"`
}

type TeamQueryResponse struct {
	BaseResponse
	Data TeamPageResponse `json:"data"`
}

type TeamMemberInfoResponse struct {
	BaseResponse
	Data TeamMember `json:"data"`
}

type ProductListPageResponse struct {
	BasePageResponse
	Data []ProductRes `json:"data"`
}

// ProductItemRes
type ProductItemListPageResponse struct {
	BasePageResponse
	Data []ProductItemRes `json:"data"`
}

// ProductItemRes
type ProductItemInfoResponse struct {
	BaseResponse
	Data ProductItemRes `json:"data"`
}

// ProductShow list
type ProductShowListPageResponse struct {
	BasePageResponse
	Data []ProductShow `json:"data"`
}

// ProductShowItem
type ProductShowItemInfoResponse struct {
	BaseResponse
	Data ProductShowItem `json:"data"`
}

// Configuration
type ConfigurationInfoResponse struct {
	BaseResponse
	Data Configuration `json:"data"`
}

// Configuration page list
type ConfigurationPageResponse struct {
	BasePageResponse
	Data []Configuration `json:"data"`
}

// map[string]string
type ConfigurationMapResponse struct {
	BaseResponse
	Data map[string]string `json:"data"`
}

// Currency
type CurrencyInfoResponse struct {
	BaseResponse
	Data Currency `json:"data"`
}

// Currency page list
type CurrencyPageResponse struct {
	BasePageResponse
	Data []Currency `json:"data"`
}

// Currency list
type CurrencyListResponse struct {
	BaseResponse
	Data []Currency `json:"data"`
}

// OrderRes Page
type OrderListPageResponse struct {
	BasePageResponse
	Data []OrderRes `json:"data"`
}

// OrderItemRes List
type OrderItemListResponse struct {
	BaseResponse
	Data []OrderItemRes `json:"data"`
}

// Order
type OrderInfoResponse struct {
	BaseResponse
	Data Order `json:"data"`
}

// PermissionMenu
type PermissionMenuInfoResponse struct {
	BaseResponse
	Data PermissionMenu `json:"data"`
}

// PermissionMenu page list
type PermissionMenuPageResponse struct {
	BasePageResponse
	Data []PermissionMenu `json:"data"`
}

// PermissionMenu list
type PermissionMenuListResponse struct {
	BaseResponse
	Data []PermissionMenu `json:"data"`
}

// PermissionMenuCreate
type PermissionMenuCreateResponse struct {
	BaseResponse
	Data ReqPermissionMenuCreate `json:"data"`
}

// UserPermission
type UserPermissionInfoResponse struct {
	BaseResponse
	Data UserPermission `json:"data"`
}

// UserPermission list
type UserPermissionListResponse struct {
	BaseResponse
	Data []UserPermission `json:"data"`
}

// UserPermission page list
type UserPermissionPageResponse struct {
	BasePageResponse
	Data []UserPermission `json:"data"`
}

// Permission
type PermissionInfoResponse struct {
	BaseResponse
	Data Permission `json:"data"`
}

// PermissionListResponse
type PermissionListResponse struct {
	BasePageResponse
	Data []Permission `json:"data"`
}

// ProductCategory list
type ProductCategoryListResponse struct {
	BaseResponse
	Data []ProductCategory `json:"data"`
}

// ProductCategory
type ProductCategoryInfoResponse struct {
	BaseResponse
	Data ProductCategory `json:"data"`
}

// TeamMember page list
type TeamMemberPageResponse struct {
	BasePageResponse
	Data []TeamUserRes `json:"data"`
}

// Team List Response
type TeamListResponse struct {
	BaseResponse
	Data []UserTeamRes `json:"data"`
}

// MenuListResponse
type MenuListResponse struct {
	BaseResponse
	Data []Menu `json:"data"`
}

// PaypalOrderDetail
type PaypalOrderDetailResponse struct {
	BaseResponse
	Data PaypalOrderDetail `json:"data"`
}

// PaymentMethod list
type PaymentMethodListResponse struct {
	BaseResponse
	Data []PaymentMethod `json:"data"`
}
