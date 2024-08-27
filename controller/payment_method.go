package controller

import (
	"encoding/json"
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
func (p *PaymentMethodController) GetPaymentMethodAll(ctx *app.Context) {
	r, err := p.PaymentMethodService.GetAvailablePaymentMethodList(ctx)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(r)
}

// CreatePaypalPayment
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
