package controller

import (
	"net/http"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/service"
)

type PaymentMethodController struct {
	PaymentMethodService *service.PaymentMethodService
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
// @Success 200 {object} model.PaymentMethodQueryResponse
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
