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
// @Summary 创建产品
// @Description 创建产品
// @Tags 产品
// @Accept  json
// @Produce  json
// @Param param body model.ReqProductCreate true "产品参数"
// @Success 200 {object} model.StringDataResponse  "Created product"
// @Router /api/v1/product/create [post]
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
// @Summary 获取产品列表
// @Description 获取产品列表
// @Tags 产品
// @Accept  json
// @Produce  json
// @Param  param body model.ReqProductQueryParam true "产品参数"
// @Success 200 {object} model.ProductListPageResponse "产品列表"
// @Router /api/v1/product/list [post]
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
// @Summary 删除产品
// @Description 删除产品
// @Tags 产品
// @Accept  json
// @Produce  json
// @Param param body model.ReqProductDeleteParam true "产品UUID"
// @Success 200 {object} model.StringDataResponse "Deleted product"
// @Router /api/v1/product/delete [post]
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
// @Summary 获取产品信息
// @Description 获取产品信息
// @Tags 产品
// @Accept  json
// @Produce  json
// @Param param body model.ReqUuidParam true "产品UUID"
// @Success 200 {object} model.ProductRes "产品信息"
// @Router /api/v1/product/info [post]
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

// GetShowProductInfo
// @Summary 获取产品前端展示信息
// @Description 获取产品前端展示信息
// @Tags 产品
// @Accept  json
// @Produce  json
// @Param param body model.ReqUuidParam true "产品UUID"
// @Success 200 {object} model.ProductShowItemInfoResponse "产品信息"
// @Router /api/v1/f/product/info [post]
func (p *ProductController) GetShowProductInfo(ctx *app.Context) {
	// 创建参数
	params := &model.ReqUuidParam{}
	// 绑定参数
	if err := ctx.Bind(params); err != nil {
		ctx.Logger.Error("Failed to bind params", err)
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	// 获取产品信息
	info, err := p.ProductService.GetShowProductInfo(ctx, params.Uuid)
	if err != nil {
		ctx.Logger.Error("Failed to get product info", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(info)
}

// GetShowProductList
// @Summary 获取产品前端展示列表
// @Description 获取产品前端展示列表
// @Tags 产品
// @Accept  json
// @Produce  json
// @Param param body model.ReqProductQueryParam true "产品参数"
// @Success 200 {object} model.ProductShowListPageResponse "产品列表"
// @Router /api/v1/f/product/list [post]
func (p *ProductController) GetShowProductList(ctx *app.Context) {

	// 创建参数
	params := &model.ReqProductQueryParam{}
	// 绑定参数
	if err := ctx.Bind(params); err != nil {
		ctx.Logger.Error("Failed to bind params", err)
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	// 获取产品列表
	list, err := p.ProductService.GetShowProductList(ctx, params)
	if err != nil {

		ctx.Logger.Error("Failed to get product list", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(list)
}

// GetProductItemList
// @Summary 获取产品SKU列表
// @Description 获取产品SKU列表
// @Tags 产品
// @Accept  json
// @Produce  json
// @Param param body model.ReqProductQueryParam true "产品参数"
// @Success 200 {object} model.ProductItemListPageResponse "产品SKU列表"
// @Router /api/v1/product/item/list [post]
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
// @Summary 获取产品SKU信息
// @Description 获取产品SKU信息
// @Tags 产品
// @Accept  json
// @Produce  json
// @Param param body model.ReqUuidParam true "产品SKU UUID"
// @Success 200 {object} model.ProductItemInfoResponse "产品SKU信息"
// @Router /api/v1/product/item/info [post]
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
// @Summary 删除产品SKU
// @Description 删除产品SKU
// @Tags 产品
// @Accept  json
// @Produce  json
// @Param param body model.ReqProductDeleteParam true "产品SKU UUID"
// @Success 200 {object} model.StringDataResponse "Deleted product"
// @Router /api/v1/product/item/delete [post]
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
