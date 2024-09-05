package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sgin/model"
	"sgin/pkg/app"
	paymentmethod "sgin/pkg/payment-method"
	"sgin/service"

	"github.com/google/uuid"
)

type PaymentMethodController struct {
	PaymentMethodService *service.PaymentMethodService
	OrderService         *service.OrderService
	PaymentService       *service.PaymentService
}

// @Summary 创建支付方式
// @Description 创建支付方式
// @Tags 支付方式
// @Accept  json
// @Produce  json
// @Param param body model.PaymentMethod true "支付方式参数"
// @Success 200 {object} model.PaymentMethod
// @Router /api/v1/payment-method/create [post]
func (p *PaymentMethodController) CreatePaymentMethod(ctx *app.Context) {
	var param model.PaymentMethod
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := p.PaymentMethodService.CreatePaymentMethod(ctx, &param); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(param)
}

// @Summary 更新支付方式
// @Description 更新支付方式
// @Tags 支付方式
// @Accept  json
// @Produce  json
// @Param param body model.PaymentMethod true "支付方式参数"
// @Success 200 {object} model.PaymentMethod
// @Router /api/v1/payment-method/update [post]
func (p *PaymentMethodController) UpdatePaymentMethodConfig(ctx *app.Context) {
	var param model.PaymentMethod
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := p.PaymentMethodService.UpdatePaymentMethodConfig(ctx, param.Uuid, param.Config); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(param)
}

// UpdatePaymentMethod

func (p *PaymentMethodController) UpdatePaymentMethodStatus(ctx *app.Context) {
	var param model.PaymentMethod
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := p.PaymentMethodService.UpdatePaymentMethodStatus(ctx, param.Uuid, param.Status); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(param)
}

// @Summary 删除支付方式
// @Description 删除支付方式
// @Tags 支付方式
// @Accept  json
// @Produce  json
// @Param param body model.ReqUuidParam true "支付方式uuid"
// @Success 200 {string} string "ok"
// @Router /api/v1/payment-method/delete [post]
func (p *PaymentMethodController) DeletePaymentMethod(ctx *app.Context) {
	var param model.ReqUuidParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := p.PaymentMethodService.DeletePaymentMethod(ctx, param.Uuid); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess("ok")
}

// @Summary 获取支付方式列表
// @Description 获取支付方式列表
// @Tags 支付方式
// @Accept  json
// @Produce  json
// @Param param body model.ReqPaymentMethodQueryParam false "支付方式查询参数"
// @Success 200 {object} app.Response
// @Router /api/v1/payment-method/list [post]
func (p *PaymentMethodController) GetPaymentMethodList(ctx *app.Context) {
	var param model.ReqPaymentMethodQueryParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	r, err := p.PaymentMethodService.GetPaymentMethodList(ctx, &param)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(r)
}

// @Summary 获取支付方式信息
// @Description 获取支付方式信息
// @Tags 支付方式
// @Accept  json
// @Produce  json
// @Param param body model.ReqUuidParam true "支付方式uuid"
// @Success 200 {object} model.PaymentMethod
// @Router /api/v1/payment-method/info [post]
func (p *PaymentMethodController) GetPaymentMethodInfo(ctx *app.Context) {
	var param model.ReqPaymentMethodInfoParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	r, err := p.PaymentMethodService.GetPaymentMethodInfo(ctx, param.Uuid, param.Code)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(r)
}

// GetPaymentMethodAll
// @Summary 获取所有可用支付方式
// @Description 获取所有可用支付方式
// @Tags 支付方式
// @Accept  json
// @Produce  json
// @Success 200 {object} model.PaymentMethod
// @Router /api/v1/f/payment_method/all [post]
func (p *PaymentMethodController) GetPaymentMethodAll(ctx *app.Context) {
	r, err := p.PaymentMethodService.GetAvailablePaymentMethodList(ctx)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(r)
}

// CreatePaypalPayment
// @Summary 创建Paypal支付
// @Description 创建Paypal支付
// @Tags 支付方式
// @Accept  json
// @Produce  json
// @Param param body model.ReqPaymentOrderCreateParam true "支付订单参数"
// @Success 200 {object} model.PaypalOrderDetailResponse
// @Router /api/v1/payment_method/paypal/create [post]
func (p *PaymentMethodController) CreatePaypalPayment(ctx *app.Context) {

	var param model.ReqPaymentOrderCreateParam

	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	userId := ctx.GetString("user_id")
	if userId == "" {
		ctx.Logger.Error("Unauthorized")
		ctx.JSONError(http.StatusUnauthorized, "Unauthorized")
		return
	}

	// 获取订单信息
	order, err := p.OrderService.GetOrderByID(ctx, param.OrderID)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	// 获取支付方式信息
	payment, err := p.PaymentMethodService.GetPaymentMethodInfo(ctx, "", "paypal")
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	mdata := make(map[string]interface{})
	err = json.Unmarshal([]byte(payment.Config), &mdata)
	if err != nil {
		ctx.Logger.Error("Failed to unmarshal payment config", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	v, ok := mdata["production"]
	if !ok {
		ctx.Logger.Error("Failed to get production config")
		ctx.JSONError(http.StatusInternalServerError, "production not found")
		return
	}

	b, err := json.Marshal(v)
	if err != nil {
		ctx.Logger.Error("Failed to marshal production config", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	paypal := &paymentmethod.PayPal{}
	err = json.Unmarshal(b, paypal)
	if err != nil {
		ctx.Logger.Error("Failed to unmarshal paypal config", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	r, err := paypal.CreateOrder(ctx, order.OrderNo, "USD", order.TotalAmount, order.OrderNo)
	if err != nil {
		ctx.Logger.Error("Failed to create order", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	b, _ = json.Marshal(r)

	paymentInfo := &model.Payment{
		Uuid:           uuid.New().String(),
		UserID:         userId,
		OrderID:        order.OrderNo,
		Amount:         order.TotalAmount,
		Status:         model.PaymentStatusPending,
		Method:         "paypal",
		Channel:        "web",
		ChannelOrderNo: r.Id,
		ChannelStatus:  r.Status,
		ChannelData:    string(b),
	}

	_, err = p.PaymentService.CreatePayment(ctx, paymentInfo)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(r)
}

// CreateAlipayPayment
// @Summary 创建支付宝支付
// @Description 创建支付宝支付
// @Tags 支付方式
// @Accept  json
// @Produce  json
// @Param param body model.ReqPaymentOrderCreateParam true "支付订单参数"
// @Success 200 {object} model.AlipayOrderDetailResponse
// @Router /api/v1/payment_method/alipay/create [post]
func (p *PaymentMethodController) CreateAlipayPayment(ctx *app.Context) {

	var param model.ReqPaymentOrderCreateParam

	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	userId := ctx.GetString("user_id")
	if userId == "" {
		ctx.Logger.Error("Unauthorized")
		ctx.JSONError(http.StatusUnauthorized, "Unauthorized")
		return
	}

	// 获取订单信息
	order, err := p.OrderService.GetOrderByID(ctx, param.OrderID)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	// 获取支付方式信息
	payment, err := p.PaymentMethodService.GetPaymentMethodInfo(ctx, "", "alipay")
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	alipayClient := &paymentmethod.Alipay{}
	err = json.Unmarshal([]byte(payment.Config), alipayClient)
	if err != nil {
		ctx.Logger.Error("Failed to unmarshal alipay config", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	r, err := alipayClient.CreateOrder(ctx, order.OrderNo, "USD", order.TotalAmount, order.OrderNo)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	b, _ := json.Marshal(r)

	ctx.Logger.Info("Create alipay order successfully:", string(b))

	paymentInfo := &model.Payment{
		Uuid:    uuid.New().String(),
		UserID:  userId,
		OrderID: order.OrderNo,
		Amount:  order.TotalAmount,
		Status:  model.PaymentStatusPending,
		Method:  "alipay",

		ChannelOrderNo: r.Response.TradeNo,
		ChannelStatus:  r.Response.Code,
		ChannelData:    string(b),
	}

	_, err = p.PaymentService.CreatePayment(ctx, paymentInfo)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(r)
}

// CreateWechatPayment
// @Summary 创建微信支付
// @Description 创建微信支付
// @Tags 支付方式
// @Accept  json
// @Produce  json
// @Param param body model.ReqPaymentOrderCreateParam true "支付订单参数"
// @Success 200 {object} model.WechatOrderDetailResponse
// @Router /api/v1/payment_method/wechat/create [post]
func (p *PaymentMethodController) CreateWechatPayment(ctx *app.Context) {

	var param model.ReqPaymentOrderCreateParam

	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	userId := ctx.GetString("user_id")
	if userId == "" {
		ctx.Logger.Error("Unauthorized")
		ctx.JSONError(http.StatusUnauthorized, "Unauthorized")
		return
	}

	// 获取订单信息
	order, err := p.OrderService.GetOrderByID(ctx, param.OrderID)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	// 获取支付方式信息
	payment, err := p.PaymentMethodService.GetPaymentMethodInfo(ctx, "", "wechat")
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	wechatClient := &paymentmethod.Wechat{}
	err = json.Unmarshal([]byte(payment.Config), wechatClient)
	if err != nil {
		ctx.Logger.Error("Failed to unmarshal wechat config", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	r, err := wechatClient.CreateOrder(ctx, order.OrderNo, "CNY", order.TotalAmount, order.OrderNo)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	b, _ := json.Marshal(r)

	ctx.Logger.Info("Create wechat order successfully:", string(b))

	paymentInfo := &model.Payment{
		Uuid:    uuid.New().String(),
		UserID:  userId,
		OrderID: order.OrderNo,
		Amount:  order.TotalAmount,
		Status:  model.PaymentStatusPending,
		Method:  "wechat",

		ChannelOrderNo: r.Response.PrepayId,
		ChannelStatus:  "",
		ChannelData:    string(b),
	}

	_, err = p.PaymentService.CreatePayment(ctx, paymentInfo)

	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(r)
}

// CreatePaypalPaymentSandboxTest
func (p *PaymentMethodController) CreatePaypalPaymentSandboxTest(ctx *app.Context) {

	var param model.ReqPaypalOrderCreateParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	// 获取沙盒测试支付信息

	payment, err := p.PaymentMethodService.GetPaymentMethodInfo(ctx, "", "paypal")
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	mdata := make(map[string]interface{})
	err = json.Unmarshal([]byte(payment.Config), &mdata)
	if err != nil {
		ctx.Logger.Error("Failed to unmarshal payment config", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	v, ok := mdata["sandbox"]
	if !ok {
		ctx.Logger.Error("Failed to get sandbox config")
		ctx.JSONError(http.StatusInternalServerError, "sandbox not found")
		return
	}

	b, err := json.Marshal(v)
	if err != nil {
		ctx.Logger.Error("Failed to marshal sandbox config", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	paypal := &paymentmethod.PayPal{}
	err = json.Unmarshal(b, paypal)
	if err != nil {
		ctx.Logger.Error("Failed to unmarshal paypal config", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	r, err := paypal.CreateSandBoxOrder(ctx, "USD", param.Amount, param.Name)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(r)
}

// GetPaypalClientID
func (p *PaymentMethodController) GetPaypalClientID(ctx *app.Context) {

	params := &model.ReqPaypalClientIdParam{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	payment, err := p.PaymentMethodService.GetPaymentMethodInfo(ctx, "", "paypal")
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	mdata := make(map[string]interface{})
	err = json.Unmarshal([]byte(payment.Config), &mdata)
	if err != nil {
		ctx.Logger.Error("Failed to unmarshal payment config", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	v, ok := mdata[params.Env]
	if !ok {
		ctx.Logger.Error("Failed to get sandbox config")
		ctx.JSONError(http.StatusInternalServerError, "sandbox not found")
		return
	}

	b, err := json.Marshal(v)
	if err != nil {
		ctx.Logger.Error("Failed to marshal sandbox config", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	paypal := &paymentmethod.PayPal{}
	err = json.Unmarshal(b, paypal)
	if err != nil {
		ctx.Logger.Error("Failed to unmarshal paypal config", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(paypal.Clientid)
}

// SetAlipayConfig
// @Summary 设置支付宝配置
// @Description 设置支付宝配置
// @Tags 支付方式
// @Accept  json
// @Produce  json
// @Param private_key formData file true "支付宝私钥"
// @Param app_public_key formData file true "支付宝公钥"
// @Param alipay_public_cert formData file true "支付宝公钥证书"
// @Param alipay_root_cert formData file true "支付宝根证书"
// @Param app_id formData string true "app_id"
// @Success 200 {object} model.StringDataResponse "ok"
// @Router /api/v1/payment-method/alipay/config [post]
func (p *PaymentMethodController) SetAlipayConfig(ctx *app.Context) {

	// 从form 获取私钥文件
	file, err := ctx.FormFile("private_key")
	if err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	// 读取私钥文件内容
	f, err := file.Open()
	if err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	defer f.Close()

	privateKey, err := ioutil.ReadAll(f)
	if err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	// 从form 获取支付宝公钥文件
	file2, err := ctx.FormFile("app_public_key")
	if err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	// 读取支付宝公钥文件内容
	f2, err := file2.Open()
	if err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	defer f2.Close()

	appPublicKey, err := ioutil.ReadAll(f2)

	// 从form 获取支付宝公钥证书文件
	file3, err := ctx.FormFile("alipay_public_cert")
	if err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	// 读取支付宝公钥证书文件内容
	f3, err := file3.Open()
	if err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	defer f3.Close()

	alipayPublicCert, err := ioutil.ReadAll(f3)
	if err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	// 从form 获取支付宝根证书文件
	file4, err := ctx.FormFile("alipay_root_cert")
	if err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	// 读取支付宝根证书文件内容
	f4, err := file4.Open()
	if err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	defer f4.Close()

	alipayRootCert, err := ioutil.ReadAll(f4)
	if err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	appid := ctx.PostForm("app_id")

	if appid == "" {
		ctx.JSONError(http.StatusBadRequest, "app_id is required")
		return
	}

	// 保存支付宝配置
	alipay := &paymentmethod.Alipay{
		AppID:         appid,
		PrivateKey:    string(privateKey),
		IsProd:        true,
		AppPublicCert: string(appPublicKey),

		AlipayRootCert:   string(alipayRootCert),
		AlipayPublicCert: string(alipayPublicCert),
	}

	b, err := json.Marshal(alipay)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	// 获取支付宝的payment
	payment, err := p.PaymentMethodService.GetPaymentMethodInfo(ctx, "", "alipay")
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	// 更新支付宝配置
	err = p.PaymentMethodService.UpdatePaymentMethodConfig(ctx, payment.Uuid, string(b))
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess("ok")
}

// SetWechatConfig
// @Summary 设置微信支付配置
// @Description 设置微信支付配置
// @Tags 支付方式
// @Accept  json
// @Produce  json
// @Param serial_no formData string true "serial_no"
// @Param mch_id formData string true "mch_id"
// @Param api_key formData string true
// @Param key formData file true "key"
// @Success 200 {object} model.StringDataResponse "ok"
// @Router /api/v1/payment-method/wechat/config [post]
func (p *PaymentMethodController) SetWechatConfig(ctx *app.Context) {

	serialNo := ctx.PostForm("serial_no")
	mchID := ctx.PostForm("mch_id")
	apiKey := ctx.PostForm("api_key")

	// 从form 获取微信支付证书文件
	file, err := ctx.FormFile("key")
	if err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	// 读取微信支付证书文件内容
	f, err := file.Open()
	if err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	defer f.Close()

	key, err := ioutil.ReadAll(f)
	if err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	// 保存微信支付配置
	wechat := &paymentmethod.Wechat{
		SerialNo:   serialNo,
		Mchid:      mchID,
		ApiV3Key:   apiKey,
		PrivateKey: string(key),
	}

	b, err := json.Marshal(wechat)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	// 获取微信支付的payment
	payment, err := p.PaymentMethodService.GetPaymentMethodInfo(ctx, "", "wechat")
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	// 更新微信支付配置
	err = p.PaymentMethodService.UpdatePaymentMethodConfig(ctx, payment.Uuid, string(b))
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess("ok")
}
