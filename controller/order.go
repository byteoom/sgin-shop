package controller

import (
	"net/http"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/service"
)

type OrderController struct {
	OrderService *service.OrderService
}

// 查询订单列表
// @Summary 查询订单列表
// @Tags 订单
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.ReqOrderQueryParam false "查询参数"
// @Success 200 {object} model.OrderQueryResponse
// @Router /api/v1/order/list [post]
func (c *OrderController) GetOrderList(ctx *app.Context) {
	param := &model.ReqOrderQueryParam{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	orders, err := c.OrderService.GetOrderList(ctx, param)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(orders)
}

// 创建订单
// @Summary 创建订单
// @Tags 订单
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.ReqOrderCreate true "订单信息"
// @Success 200 {object} model.Order
// @Router /api/v1/order/create [post]
func (c *OrderController) CreateOrder(ctx *app.Context) {
	param := &model.ReqOrderCreate{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	userId := ctx.GetString("user_id")
	if userId == "" {
		ctx.JSONError(http.StatusUnauthorized, "user_id is empty")
		return
	}

	param.UserId = userId

	if len(param.Items) == 0 && len(param.CartUuids) == 0 {
		ctx.JSONError(http.StatusBadRequest, "items and cartUuids can't be empty at the same time")
		return
	}

	if len(param.Items) > 0 && len(param.CartUuids) > 0 {
		ctx.JSONError(http.StatusBadRequest, "items and cartUuids can't be set at the same time")
		return
	}

	if len(param.Items) > 0 {
		order, err := c.OrderService.CreateOrder(ctx, param)
		if err != nil {
			ctx.JSONError(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSONSuccess(order)
	} else {
		order, err := c.OrderService.CreateOrderByCart(ctx, param)
		if err != nil {
			ctx.JSONError(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSONSuccess(order)
	}
}

// 删除订单
// @Summary 删除订单
// @Tags 订单
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.ReqOrderDeleteParam true "删除参数"
// @Success 200 {object} app.Response
// @Router /api/v1/order/delete [post]
func (c *OrderController) DeleteOrder(ctx *app.Context) {
	param := &model.ReqUuidParam{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := c.OrderService.DeleteOrder(ctx, param.Uuid)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess("ok")
}

// 查询订单详情
// @Summary 查询订单详情
// @Tags 订单
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.ReqOrderInfoParam true "查询参数"
// @Success 200 {object} model.Order
// @Router /api/v1/order/info [post]
func (c *OrderController) GetOrderInfo(ctx *app.Context) {
	param := &model.ReqUuidParam{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	order, err := c.OrderService.GetOrderByID(ctx, param.Uuid)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(order)
}

// GetOrderItemList

func (c *OrderController) GetOrderItemList(ctx *app.Context) {
	param := &model.ReqUuidParam{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	items, err := c.OrderService.GetOrderItemsByOrderNo(ctx, param.Uuid)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(items)
}
