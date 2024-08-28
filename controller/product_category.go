package controller

import (
	"net/http"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/service"
)

type ProductCategoryController struct {
	CategoryService *service.ProductCategoryService
}

// 获取分类列表
// @Summary 获取分类列表
// @Description 获取分类列表
// @Tags 产品分类
// @Accept  json
// @Produce  json
// @Success 200 {object} model.ProductCategoryListResponse "分类列表"
// @Router /api/v1/product_category/list [post]
func (c *ProductCategoryController) GetCategoryList(ctx *app.Context) {
	categories, err := c.CategoryService.GetCategoryList(ctx)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(categories)
}

// 创建分类
// @Summary 创建分类
// @Description 创建分类
// @Tags 产品分类
// @Accept  json
// @Produce  json
// @Param param body model.ProductCategory true "分类参数"
// @Success 200 {object} model.ProductCategoryInfoResponse "Created category"
// @Router /api/v1/product_category/create [post]
func (c *ProductCategoryController) CreateCategory(ctx *app.Context) {
	category := &model.ProductCategory{}
	if err := ctx.ShouldBindJSON(category); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := c.CategoryService.CreateCategory(ctx, category)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(category)
}

// 更新分类
// @Summary 更新分类
// @Description 更新分类
// @Tags 产品分类
// @Accept  json
// @Produce  json
// @Param param body model.ProductCategory true "分类参数"
// @Success 200 {object} model.ProductCategoryInfoResponse "Updated category"
// @Router /api/v1/product_category/update [post]
func (c *ProductCategoryController) UpdateCategory(ctx *app.Context) {
	category := &model.ProductCategory{}
	if err := ctx.ShouldBindJSON(category); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := c.CategoryService.UpdateCategory(ctx, category)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(category)
}

// 删除分类
// @Summary 删除分类
// @Description 删除分类
// @Tags 产品分类
// @Accept  json
// @Produce  json
// @Param param body model.ReqUuidParam true "分类参数"
// @Success 200 {object} model.StringDataResponse  "Deleted successfully"
// @Router /api/v1/product_category/delete [post]
func (c *ProductCategoryController) DeleteCategory(ctx *app.Context) {
	param := &model.ReqUuidParam{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := c.CategoryService.DeleteCategory(ctx, param.Uuid)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess("Deleted successfully")
}

// 获取分类详情
// @Summary 获取分类详情
// @Description 获取分类详情
// @Tags 产品分类
// @Accept  json
// @Produce  json
// @Param param body model.ReqUuidParam true "分类UUID"
// @Success 200 {object} model.ProductCategoryInfoResponse "分类详情"
// @Router /api/v1/product_category/info [post]
func (c *ProductCategoryController) GetCategory(ctx *app.Context) {
	param := &model.ReqUuidParam{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	category, err := c.CategoryService.GetCategoryByUUID(ctx, param.Uuid)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(category)
}

// 获取所有分类
// @Summary 用户前端获取所有分类
// @Description 用户前端获取所有分类
// @Tags 产品分类
// @Accept  json
// @Produce  json
// @Success 200 {object} model.ProductCategoryListResponse "分类列表"
// @Router /api/v1/f/product_category/all [post]
func (c *ProductCategoryController) GetAllCategory(ctx *app.Context) {
	categories, err := c.CategoryService.GetAllCategory(ctx)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(categories)
}
