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
func (c *ProductCategoryController) GetCategoryList(ctx *app.Context) {
	categories, err := c.CategoryService.GetCategoryList(ctx)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(categories)
}

// 创建分类
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
func (c *ProductCategoryController) GetAllCategory(ctx *app.Context) {
	categories, err := c.CategoryService.GetAllCategory(ctx)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(categories)
}
