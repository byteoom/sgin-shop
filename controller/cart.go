package controller

import (
	"net/http"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/service"
)

type CartController struct {
	CartService *service.CartService
}

// 查询购物车列表
// @Summary 查询购物车列表
// @Tags 购物车
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.ReqCartQueryParam false "查询参数"
// @Success 200 {object} model.Cart
// @Router /api/v1/cart/list [post]
func (c *CartController) GetCartList(ctx *app.Context) {
	param := &model.ReqCartQueryParam{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	userId := ctx.GetString("user_id")

	if userId == "" {
		ctx.JSONError(http.StatusBadRequest, "user_id is required")
		return
	}

	param.UserID = userId

	carts, err := c.CartService.GetCartList(ctx, param)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(carts)
}

// 创建购物车
// @Summary 创建购物车
// @Tags 购物车
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.Cart true "购物车信息"
// @Success 200 {object} model.Cart
// @Router /api/v1/cart/create [post]
func (c *CartController) CreateCart(ctx *app.Context) {
	param := &model.Cart{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	userId := ctx.GetString("user_id")
	if userId == "" {
		ctx.JSONError(http.StatusBadRequest, "user_id is required")
		return
	}
	param.UserID = userId

	err := c.CartService.CreateCart(ctx, param)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(param)
}

// 更新购物车
// @Summary 更新购物车
// @Tags 购物车
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.Cart true "购物车信息"
// @Success 200 {object} model.Cart
// @Router /api/v1/cart/update [post]
func (c *CartController) UpdateCart(ctx *app.Context) {
	param := &model.Cart{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := c.CartService.UpdateCart(ctx, param)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(param)
}

// UpdateCartItemCount
// @Summary 更新购物车商品数量
// @Tags 购物车
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.ReqCartItemCountParam true "更新参数"
// @Success 200 {string} string "ok"
// @Router /api/v1/cart/update/count [post]
func (c *CartController) UpdateCartItemCount(ctx *app.Context) {
	param := &model.ReqCartItemCountParam{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := c.CartService.UpdateCartItemCount(ctx, param.Uuid, param.Quantity)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess("ok")
}

// 删除购物车
// @Summary 删除购物车
// @Tags 购物车
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.ReqUuidParam true "删除参数"
// @Success 200 {string} string "ok"
// @Router /api/v1/cart/delete [post]
func (c *CartController) DeleteCart(ctx *app.Context) {
	param := &model.ReqUuidParam{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := c.CartService.DeleteCart(ctx, param.Uuid)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess("ok")
}

// 查询购物车详情
// @Summary 查询购物车详情
// @Tags 购物车
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.ReqUuidParam true "查询参数"
// @Success 200 {object} model.Cart
// @Router /api/v1/cart/info [post]
func (c *CartController) GetCartInfo(ctx *app.Context) {
	param := &model.ReqUuidParam{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	cart, err := c.CartService.GetCartByUUID(ctx, param.Uuid)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(cart)
}
