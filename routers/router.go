package routers

import (
	"sgin/controller"
	"sgin/middleware"
	"sgin/pkg/app"

	"io/ioutil"
	"net/http"
	"sgin/service"

	"github.com/gin-gonic/gin"
)

func InitRouter(ctx *app.App) {
	InitSwaggerRouter(ctx)
	InitUserRouter(ctx)
	InitMenuRouter(ctx)
	InitAppRouter(ctx)
	InitVerificationCodeRouter(ctx)
	InitRegisterRouter(ctx)
	InitLoginRouter(ctx)
	InitServerRouter(ctx)
	InitTeamRouter(ctx)
	InitSysLoginLogRouter(ctx)
	InitSysOpLogRouter(ctx)
	InitSysApiRouter(ctx)
	InitPermissionRouter(ctx)
	InitPermissionMenuRouter(ctx)
	InitPermissionUserRouter(ctx)
	InitMenuAPIRouter(ctx)
	InitTeamMemberRouter(ctx)
	InitProductCategoryRouter(ctx)
	InitResourceRouter(ctx)
	InitProductRouter(ctx)
	InitPaymentRouter(ctx)
	InitCartRouter(ctx)
	InitOrderRouter(ctx)
	InitProductFrontRouter(ctx)
	InitPaymentMethodRouter(ctx)
	InitPaypalRouter(ctx)
	InitConfigurationRouter(ctx)
	InitUserAddressRouter(ctx)
	InitCurrencyRouter(ctx)
}

func InitUserRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	v1.Use(middleware.SysOpLogMiddleware(&service.SysOpLogService{}))
	{
		userController := &controller.UserController{
			Service: &service.UserService{},
		}

		v1.POST("/user/create", userController.CreateUser)
		v1.POST("/user/info", userController.GetUserByUUID)
		v1.POST("/user/list", userController.GetUserList)
		v1.POST("/user/update", userController.UpdateUser)
		v1.POST("/user/delete", userController.DeleteUser)
		v1.GET("/user/myinfo", userController.GetMyInfo)
		v1.POST("/user/avatar", userController.UpdateAvatar)
		v1.POST("/user/all", userController.GetAllUsers)

	}

	{
		roleController := &controller.RoleController{
			RoleService: &service.RoleService{},
		}

		v1.POST("/role/create", roleController.CreateRole)
		v1.POST("/role/list", roleController.GetRoleList)
		v1.POST("/role/update", roleController.UpdateRole)
		v1.POST("/role/delete", roleController.DeleteRole)

	}
}

func InitMenuRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	v1.Use(middleware.SysOpLogMiddleware(&service.SysOpLogService{}))
	{
		menuController := &controller.MenuController{
			MenuService: &service.MenuService{},
		}
		v1.POST("/menu/create", menuController.CreateMenu)
		v1.POST("/menu/list", menuController.GetMenuList)
		v1.POST("/menu/update", menuController.UpdateMenu)
		v1.POST("/menu/delete", menuController.DeleteMenu)
		v1.POST("/menu/info", menuController.GetMenuInfo)
	}
}

func InitAppRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	v1.Use(middleware.SysOpLogMiddleware(&service.SysOpLogService{}))
	{
		appController := &controller.AppController{
			AppService: &service.AppService{},
		}
		v1.POST("/app/list", appController.GetAppList)
		v1.POST("/app/create", appController.CreateApp)
		v1.POST("/app/update", appController.UpdateApp)
		v1.POST("/app/delete", appController.DeleteApp)

	}
}

func InitVerificationCodeRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	{
		verificationCodeController := &controller.VerificationCodeController{
			VerificationCodeService: &service.VerificationCodeService{},
		}
		v1.POST("/verification_code/create", verificationCodeController.CreateVerificationCode)
	}
}

// 注册的路由
func InitRegisterRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	{
		registerController := &controller.RegisterController{
			UserService:             &service.UserService{},
			VerificationCodeService: &service.VerificationCodeService{},
		}
		v1.POST("/register", registerController.Register)
	}
}

func InitLoginRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	{
		loginController := &controller.LoginController{
			UserService: &service.UserService{},
		}
		v1.POST("/login", loginController.Login)
	}
}

// 服务的路由
func InitServerRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	v1.Use(middleware.SysOpLogMiddleware(&service.SysOpLogService{}))
	{
		serverController := &controller.ServerController{
			ServerService: &service.ServerService{},
		}
		v1.POST("/server/create", serverController.CreateServer)
		v1.POST("/server/update", serverController.UpdateServer)
		v1.POST("/server/delete", serverController.DeleteServer)
		v1.POST("/server/info", serverController.GetServerInfo)
		v1.POST("/server/list", serverController.GetServerList)
	}
}

// 团队的路由
func InitTeamRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	v1.Use(middleware.SysOpLogMiddleware(&service.SysOpLogService{}))
	{
		teamController := &controller.TeamController{
			TeamService: &service.TeamService{},
		}
		v1.POST("/team/create", teamController.CreateTeam)
		v1.POST("/team/update", teamController.UpdateTeam)
		v1.POST("/team/delete", teamController.DeleteTeam)
		v1.POST("/team/info", teamController.GetTeamInfo)
		v1.POST("/team/list", teamController.GetTeamList)
	}
}

func InitTeamMemberRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	v1.Use(middleware.SysOpLogMiddleware(&service.SysOpLogService{}))
	{
		teamMemberController := &controller.TeamMemberController{
			TeamMemberService: &service.TeamMemberService{},
		}
		v1.POST("/team_member/create", teamMemberController.CreateTeamMember)
		v1.POST("/team_member/delete", teamMemberController.DeleteTeamMember)
		v1.POST("/team_member/list", teamMemberController.GetTeamMemberList)
	}
}

func InitSysApiRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	v1.Use(middleware.SysOpLogMiddleware(&service.SysOpLogService{}))
	{
		apiController := &controller.APIController{
			APIService: &service.APIService{},
		}
		v1.POST("/sys_api/create", apiController.CreateAPI)
		v1.POST("/sys_api/update", apiController.UpdateAPI)
		v1.POST("/sys_api/delete", apiController.DeleteAPI)
		v1.POST("/sys_api/list", apiController.GetAPIList)
		v1.POST("/sys_api/info", apiController.GetAPIInfo)

	}
}

func InitSysOpLogRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	v1.Use(middleware.SysOpLogMiddleware(&service.SysOpLogService{}))
	{
		sysOpLogController := &controller.SysOpLogController{
			SysOpLogService: &service.SysOpLogService{},
		}

		v1.POST("/sysoplog/delete", sysOpLogController.DeleteSysOpLog)
		v1.POST("/sysoplog/info", sysOpLogController.GetSysOpLogInfo)
		v1.POST("/sysoplog/list", sysOpLogController.GetSysOpLogList)
	}
}

func InitSysLoginLogRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	v1.Use(middleware.SysOpLogMiddleware(&service.SysOpLogService{}))
	{
		sysLoginLogController := &controller.SysLoginLogController{
			LoginLogService: &service.SysLoginLogService{},
		}

		v1.POST("/sys_login_log/info", sysLoginLogController.GetLoginLog)
		v1.POST("/sys_login_log/list", sysLoginLogController.GetLoginLogList)
	}
}

func InitPermissionRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	v1.Use(middleware.SysOpLogMiddleware(&service.SysOpLogService{}))
	{
		permissionController := &controller.PermissionController{
			PermissionService: &service.PermissionService{},
		}
		v1.POST("/permission/create", permissionController.CreatePermission)
		v1.POST("/permission/update", permissionController.UpdatePermission)
		v1.POST("/permission/delete", permissionController.DeletePermission)
		v1.POST("/permission/info", permissionController.GetPermissionInfo)
		v1.POST("/permission/list", permissionController.GetPermissionList)
	}
}

func InitPermissionMenuRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	v1.Use(middleware.SysOpLogMiddleware(&service.SysOpLogService{}))
	{
		permissionMenuController := &controller.PermissionMenuController{
			PermissionMenuService: &service.PermissionMenuService{},
		}
		v1.POST("/permission_menu/create", permissionMenuController.CreatePermissionMenu)
		v1.POST("/permission_menu/update", permissionMenuController.UpdatePermissionMenu)
		v1.POST("/permission_menu/delete", permissionMenuController.DeletePermissionMenu)
		v1.POST("/permission_menu/info", permissionMenuController.GetPermissionMenuInfo)
		v1.POST("/permission_menu/info_menu", permissionMenuController.GetPermissionMenuListByPermissionUUID)
		v1.POST("/permission_menu/list", permissionMenuController.GetPermissionMenuList)
	}
}

func InitPermissionUserRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	v1.Use(middleware.SysOpLogMiddleware(&service.SysOpLogService{}))
	{
		permissionUserController := &controller.UserPermissionController{
			UserPermissionService: &service.UserPermissionService{},
		}
		v1.POST("/permission_user/create", permissionUserController.CreateUserPermission)
		v1.POST("/permission_user/update", permissionUserController.UpdateUserPermission)
		v1.POST("/permission_user/delete", permissionUserController.DeleteUserPermission)
		v1.POST("/permission_user/info", permissionUserController.GetUserPermissionInfo)
		v1.POST("/permission_user/list", permissionUserController.GetUserPermissionList)
	}
}

func InitMenuAPIRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	v1.Use(middleware.SysOpLogMiddleware(&service.SysOpLogService{}))
	{
		menuAPIController := &controller.MenuAPIController{
			MenuAPIService: &service.MenuAPIService{},
		}
		v1.POST("/menu_api/create", menuAPIController.CreateMenuAPI)
		v1.POST("/menu_api/update", menuAPIController.UpdateMenuAPI)
		v1.POST("/menu_api/delete", menuAPIController.DeleteMenuAPI)
		v1.POST("/menu_api/info", menuAPIController.GetMenuAPIInfo)
		v1.POST("/menu_api/info_menu", menuAPIController.GetMenuAPIListByMenuUUID)
		v1.POST("/menu_api/info_api", menuAPIController.GetMenuAPIListByAPIUUID)
		v1.POST("/menu_api/list", menuAPIController.GetMenuAPIList)
	}
}

func InitProductCategoryRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	{
		productCategoryController := &controller.ProductCategoryController{
			CategoryService: &service.ProductCategoryService{},
		}
		v1.POST("/product_category/create", productCategoryController.CreateCategory)
		v1.POST("/product_category/list", productCategoryController.GetCategoryList)
		v1.POST("/product_category/update", productCategoryController.UpdateCategory)
		v1.POST("/product_category/delete", productCategoryController.DeleteCategory)
		v1.POST("/product_category/all", productCategoryController.GetAllCategory)
	}
}

func InitResourceRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	{
		resourceController := &controller.ResourceController{
			ResourceService: &service.ResourceService{},
		}
		v1.POST("/resource/create", resourceController.CreateResource)
		v1.POST("/resource/list", resourceController.GetResourceList)
		v1.POST("/resource/update", resourceController.UpdateResource)
		v1.POST("/resource/delete", resourceController.DeleteResource)
		v1.POST("/resource/create_folder", resourceController.CreateFolder)
		v1.POST("/resource/folder_list", resourceController.GetFolderList)
		v1.POST("/resource/move", resourceController.MoveResource)
	}
}

func InitProductRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	{
		productController := &controller.ProductController{
			ProductService: &service.ProductService{},
		}
		v1.POST("/product/create", productController.ProductCreate)
		v1.POST("/product/list", productController.GetProductList)
		// v1.POST("/product/update", productController.UpdateProduct)
		v1.POST("/product/delete", productController.DeleteProduct)
		v1.POST("/product/info", productController.GetProductInfo)

		v1.POST("/product/item/list", productController.GetProductItemList)
		v1.POST("/product/item/delete", productController.DeleteProductItem)
		v1.POST("/product/item/info", productController.GetProductItemInfo)
	}
}

func InitProductFrontRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	{
		productController := &controller.ProductController{
			ProductService: &service.ProductService{},
		}
		v1.POST("/f/product/list", productController.GetShowProductList)
		v1.POST("/f/product/info", productController.GetShowProductInfo)
		v1.POST("/f/product/item/list", productController.GetProductItemList)
	}

	{
		productCategory := &controller.ProductCategoryController{
			CategoryService: &service.ProductCategoryService{},
		}

		v1.POST("/f/product_category/all", productCategory.GetAllCategory)
	}

	paymentMethodController := &controller.PaymentMethodController{
		PaymentMethodService: &service.PaymentMethodService{},
	}

	{
		v1.POST("/f/payment_method/all", paymentMethodController.GetPaymentMethodAll)
	}

}

func InitPaymentMethodRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	{
		paymentMethodController := &controller.PaymentMethodController{
			PaymentMethodService: &service.PaymentMethodService{},
		}
		v1.POST("/payment_method/create", paymentMethodController.CreatePaymentMethod)
		v1.POST("/payment_method/update_status", paymentMethodController.UpdatePaymentMethodStatus)

		// 更新配置
		v1.POST("/payment_method/update_config", paymentMethodController.UpdatePaymentMethodConfig)
		// 获取列表
		v1.POST("/payment_method/list", paymentMethodController.GetPaymentMethodList)

		// 获取详细信息
		v1.POST("/payment_method/info", paymentMethodController.GetPaymentMethodInfo)

		// 创建paypal 支付订单
		v1.POST("/payment_method/paypal/create", paymentMethodController.CreatePaypalPayment)

		// 创建paypal 沙盒支付订单
		//v1.POST("/payment_method/paypal/sandbox/create_test", paymentMethodController.CreatePaypalPaymentSandboxTest)
	}
}

func InitPaypalRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")

	{
		paypalController := &controller.PaypalController{}
		v1.Any("/paypal/return", paypalController.Return)
		v1.Any("/paypal/cancel", paypalController.Cancel)
	}

	{
		paymentMethodController := &controller.PaymentMethodController{
			PaymentMethodService: &service.PaymentMethodService{},
		}
		// 创建paypal 沙盒支付订单
		v1.POST("/payment_method/paypal/sandbox/create_test", paymentMethodController.CreatePaypalPaymentSandboxTest)

		// 获取paypal client id
		v1.POST("/payment_method/paypal/client_id", paymentMethodController.GetPaypalClientID)
	}
}

func InitCartRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	{
		cartController := &controller.CartController{
			CartService: &service.CartService{},
		}
		v1.POST("/cart/add", cartController.CreateCart)
		v1.POST("/cart/list", cartController.GetCartList)
		v1.POST("/cart/delete", cartController.DeleteCart)
	}
}

func InitOrderRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	{
		orderController := &controller.OrderController{
			OrderService: &service.OrderService{},
		}
		v1.POST("/order/create", orderController.CreateOrder)
		v1.POST("/order/list", orderController.GetOrderList)
		v1.POST("/order/delete", orderController.DeleteOrder)
		v1.POST("/order/info", orderController.GetOrderInfo)
		// 获取订单详情
		v1.POST("/order/item/list", orderController.GetOrderItemList)
	}
}

func InitConfigurationRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	{
		configurationController := &controller.ConfigurationController{
			ConfigurationService: &service.ConfigurationService{},
		}
		v1.POST("/configuration/create", configurationController.CreateConfiguration)
		v1.POST("/configuration/list", configurationController.GetConfigurationList)
		v1.POST("/configuration/update", configurationController.UpdateConfiguration)
		v1.POST("/configuration/info", configurationController.GetConfigurationInfo)
		// 根据category获取配置列表
		v1.POST("/configuration/category_map", configurationController.GetConfigurationMapByCategory)

		// 根据category 创建配置
		v1.POST("/configuration/category_create_map", configurationController.CreateConfigurationMapByCategory)
	}
}

func InitUserAddressRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	{
		userAddressController := &controller.UserAddressController{
			UserAddressService: &service.UserAddressService{},
		}
		v1.POST("/user_address/create", userAddressController.CreateUserAddress)
		v1.POST("/user_address/list", userAddressController.GetUserAddressList)
		v1.POST("/user_address/update", userAddressController.UpdateUserAddress)
		v1.POST("/user_address/delete", userAddressController.DeleteUserAddress)
		v1.POST("/user_address/info", userAddressController.GetUserAddressInfo)
	}
}

func InitCurrencyRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	{
		currencyController := &controller.CurrencyController{
			CurrencyService: &service.CurrencyService{},
		}
		v1.POST("/currency/create", currencyController.CreateCurrency)
		v1.POST("/currency/update", currencyController.UpdateCurrency)
		v1.POST("/currency/delete", currencyController.DeleteCurrency)
		v1.POST("/currency/list", currencyController.GetCurrencyList)
		v1.POST("/currency/all", currencyController.GetAllCurrency)
	}
}

func InitSwaggerRouter(ctx *app.App) {
	ctx.GET("/swagger/doc.json", func(c *app.Context) {
		jsonFile, err := ioutil.ReadFile("./docs/swagger.json") // Replace with your actual json file path
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Data(http.StatusOK, "application/json", jsonFile)
	})

	ctx.GET("/swagger/redoc.standalone.js", func(c *app.Context) {
		b, err := ioutil.ReadFile("./swagger/redoc.standalone.js") // Replace with your actual json file path
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", b)
	})

	ctx.GET("/swagger/index.html", func(c *app.Context) {
		b, err := ioutil.ReadFile("./swagger/swagger.html") // Replace with your actual json file path
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", b)
	})
}

// InitPaymentRouter 初始化支付相关的路由
func InitPaymentRouter(ctx *app.App) {
	v1 := ctx.Group(ctx.Config.ApiPrefix + "/v1")
	v1.Use(middleware.LoginCheck())
	{
		paymentController := &controller.PaymentController{
			PaymentService: &service.PaymentService{},
		}
		v1.POST("/payments/info", paymentController.GetPaymentByUUID)
		v1.POST("/payments/update", paymentController.UpdatePayment)
		v1.POST("/payments/delete", paymentController.DeletePayment)
		v1.POST("/payments/list", paymentController.GetPaymentList)
	}
}
