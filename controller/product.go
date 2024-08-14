package controller

import (
	"net/http"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/service"
)

type ProductController struct {
	ProductService *service.ProductService
}

// ProductCreate 创建产品
func (p *ProductController) ProductCreate(ctx *app.Context) {
	// 创建参数
	params := &model.ReqProductCreate{}
	// 绑定参数
	if err := ctx.Bind(params); err != nil {
		ctx.Logger.Error("Failed to bind params", err)
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	// 创建产品
	err := p.ProductService.ProductCreate(ctx, params)
	if err != nil {
		ctx.Logger.Error("Failed to create product", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess("Created product")
}
