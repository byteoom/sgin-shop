package controller

import (
	"net/http"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/service"
)

type PaymentController struct {
	PaymentService *service.PaymentService
}

// NewPaymentController 创建一个新的PaymentController实例
func NewPaymentController(paymentService *service.PaymentService) *PaymentController {
	return &PaymentController{
		PaymentService: paymentService,
	}
}

func (c *PaymentController) CreatePayment(ctx *app.Context) {
	var payment model.Payment
	if err := ctx.ShouldBindJSON(&payment); err != nil {
		ctx.JSONError(http.StatusBadRequest, "无效的请求数据")
		return
	}
	newPayment, err := c.PaymentService.CreatePayment(ctx, &payment)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(newPayment)
}

func (c *PaymentController) GetPaymentByUUID(ctx *app.Context) {
	param := &model.Payment{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	payment, err := c.PaymentService.GetPaymentByUUID(ctx, param.Uuid)
	if err != nil {
		if err.Error() == "payment not found" {
			ctx.JSONError(http.StatusNotFound, "付款记录未找到")
			return
		}
		ctx.JSONError(http.StatusInternalServerError, "内部服务器错误")
		return
	}
	ctx.JSONSuccess(payment)
}

func (c *PaymentController) UpdatePayment(ctx *app.Context) {
	param := &model.Payment{}
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, "无效的请求数据")
		return
	}

	err := c.PaymentService.UpdatePayment(ctx, param)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(param)
}

func (c *PaymentController) DeletePayment(ctx *app.Context) {
	param := &model.Payment{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	err := c.PaymentService.DeletePayment(ctx, param.Uuid)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, "内部服务器错误")
		return
	}

	ctx.JSONSuccess(param)
}

func (c *PaymentController) GetPaymentList(ctx *app.Context) {
	var params model.ReqPaymentQueryParam
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSONError(http.StatusBadRequest, "无效的查询参数")
		return
	}
	pagedResponse, err := c.PaymentService.GetPaymentList(ctx, &params)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, "内部服务器错误")
		return
	}
	ctx.JSONSuccess(pagedResponse)
}
