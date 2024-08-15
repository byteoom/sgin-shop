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

// GetProductList 获取产品列表
func (p *ProductController) GetProductList(ctx *app.Context) {
	// 创建参数
	params := &model.ReqProductQueryParam{}
	// 绑定参数
	if err := ctx.Bind(params); err != nil {
		ctx.Logger.Error("Failed to bind params", err)
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	// 获取产品列表
	list, err := p.ProductService.ProductList(ctx, params)
	if err != nil {
		ctx.Logger.Error("Failed to get product list", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(list)
}

// DeleteProduct
func (p *ProductController) DeleteProduct(ctx *app.Context) {
	// 创建参数
	params := &model.ReqProductDeleteParam{}
	// 绑定参数
	if err := ctx.Bind(params); err != nil {
		ctx.Logger.Error("Failed to bind params", err)
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	// 删除产品
	err := p.ProductService.DeleteProductByUUIDList(ctx, params.UUids)
	if err != nil {
		ctx.Logger.Error("Failed to delete product", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess("Deleted product")
}

// GetProductInfo
func (p *ProductController) GetProductInfo(ctx *app.Context) {
	// 创建参数
	params := &model.ReqUuidParam{}
	// 绑定参数
	if err := ctx.Bind(params); err != nil {
		ctx.Logger.Error("Failed to bind params", err)
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	// 获取产品信息
	info, err := p.ProductService.GetProductInfo(ctx, params.Uuid)
	if err != nil {
		ctx.Logger.Error("Failed to get product info", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(info)
}

// GetProductItemList
func (p *ProductController) GetProductItemList(ctx *app.Context) {
	// 创建参数
	params := &model.ReqProductQueryParam{}
	// 绑定参数
	if err := ctx.Bind(params); err != nil {
		ctx.Logger.Error("Failed to bind params", err)
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	// 获取产品列表
	list, err := p.ProductService.GetProductSkuList(ctx, params)
	if err != nil {
		ctx.Logger.Error("Failed to get product list", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(list)
}

// GetProductItemInfo
func (p *ProductController) GetProductItemInfo(ctx *app.Context) {
	// 创建参数
	params := &model.ReqUuidParam{}
	// 绑定参数
	if err := ctx.Bind(params); err != nil {
		ctx.Logger.Error("Failed to bind params", err)
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	// 获取产品信息
	info, err := p.ProductService.GetProductSkuInfo(ctx, params.Uuid)
	if err != nil {
		ctx.Logger.Error("Failed to get product info", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(info)
}

// DeleteProductItem
func (p *ProductController) DeleteProductItem(ctx *app.Context) {
	// 创建参数
	params := &model.ReqProductDeleteParam{}
	// 绑定参数
	if err := ctx.Bind(params); err != nil {
		ctx.Logger.Error("Failed to bind params", err)
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	// 删除产品
	err := p.ProductService.DeleteProductSkuByUUIDList(ctx, params.UUids)
	if err != nil {
		ctx.Logger.Error("Failed to delete product", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess("Deleted product")
}
