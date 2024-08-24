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

// CreatePayment 创建一个新的付款记录
// @Summary 创建付款记录
// @Tags 付款
// @Accept json
// @Produce json
// @Param payment body model.Payment true "创建付款记录"
// @Success 200 {object} model.Payment
// @Failure 400 {string} string "错误信息"
// @Failure 500 {string} string "内部服务器错误"
// @Router /payments [post]
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

// GetPaymentByUUID 根据UUID获取付款记录
// @Summary 根据UUID获取付款记录
// @Tags 付款
// @Accept json
// @Produce json
// @Param uuid path string true "付款记录的UUID"
// @Success 200 {object} model.Payment
// @Failure 404 {string} string "未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /payments/{uuid} [get]
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

// UpdatePayment 更新付款记录
// @Summary 更新付款记录
// @Tags 付款
// @Accept json
// @Produce json
// @Param uuid path string true "付款记录的UUID"
// @Param payment body model.Payment true "更新的付款记录"
// @Success 200 {object} model.Payment
// @Failure 404 {string} string "未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /payments/{uuid} [put]
// UpdatePayment 更新付款记录
// ... 其他注释保持不变 ...
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

// DeletePayment 删除付款记录
// @Summary 删除付款记录
// @Tags 付款
// @Accept json
// @Produce json
// @Param uuid path string true "付款记录的UUID"
// @Success 200 {string} string "ok"
// @Failure 404 {string} string "未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /payments/{uuid} [delete]
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

// GetPaymentList 获取付款记录列表
// @Summary 获取付款记录列表
// @Tags 付款
// @Accept json
// @Produce json
// @Param page query int 0 "页码"
// @Param page_size query int 10 "每页记录数"
// @Param user_id query int 0 "用户ID过滤"
// @Param order_id query int 0 "订单ID过滤"
// @Success 200 {object} model.PagedResponse
// @Failure 500 {string} string "内部服务器错误"
// @Router /payments/list [get]
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
