package controller

import (
	"net/http"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/service"
)

type OrderController struct {
	OrderService *service.OrderServer
}

func (o *OrderController) CreateOrder(ctx *app.Context) {
	params := &model.ReqOrderCreate{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := o.OrderService.CreateOrder(ctx, params)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(nil)
}

func (o *OrderController) GetOrderDetail(ctx *app.Context) {
	params := &model.ReqOrderGet{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	order, err := o.OrderService.GetOrderDetail(ctx, params.OrderNo)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(order)
}

func (o *OrderController) UpdateOrder(ctx *app.Context) {
	params := &model.ReqOrderUpdate{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	err := o.OrderService.UpdateOrder(ctx, params)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(nil)
}

func (o *OrderController) DeleteOrder(ctx *app.Context) {
	params := &model.ReqOrderGet{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	err := o.OrderService.DeleteOrder(ctx, params.OrderNo)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(nil)
}
