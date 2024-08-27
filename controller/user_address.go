package controller

import (
	"net/http"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/service"
)

type UserAddressController struct {
	UserAddressService *service.UserAddressService
}

// @Summary 创建用户地址
// @Description 创建用户地址
// @Tags 用户地址
// @Accept  json
// @Produce  json
// @Param param body model.UserAddress true "用户地址参数"
// @Success 200 {object} model.UserAddress
// @Router /api/v1/user_address/create [post]
func (u *UserAddressController) CreateUserAddress(ctx *app.Context) {
	var param model.UserAddress
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := u.UserAddressService.CreateUserAddress(ctx, &param); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(param)
}

// @Summary 更新用户地址
// @Description 更新用户地址
// @Tags 用户地址
// @Accept  json
// @Produce  json
// @Param param body model.UserAddress true "用户地址参数"
// @Success 200 {object} model.UserAddress
// @Router /api/v1/user_address/update [post]
func (u *UserAddressController) UpdateUserAddress(ctx *app.Context) {
	var param model.UserAddress
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := u.UserAddressService.UpdateUserAddress(ctx, &param); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(param)
}

// @Summary 删除用户地址
// @Description 删除用户地址
// @Tags 用户地址
// @Accept  json
// @Produce  json
// @Param param body model.ReqUuidParam true "用户地址 UUID"
// @Success 200  {object} model.StringDataResponse "ok"
// @Router /api/v1/user_address/delete [post]
func (u *UserAddressController) DeleteUserAddress(ctx *app.Context) {
	var param model.ReqUuidParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := u.UserAddressService.DeleteUserAddress(ctx, param.Uuid); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess("ok")
}

// @Summary 获取用户地址信息
// @Description 获取用户地址信息
// @Tags 用户地址
// @Accept  json
// @Produce  json
// @Param param body model.ReqUuidParam true "用户地址 UUID"
// @Success 200 {object} model.UserAddress
// @Router /api/v1/user_address/info [post]
func (u *UserAddressController) GetUserAddressInfo(ctx *app.Context) {
	var param model.ReqUuidParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	address, err := u.UserAddressService.GetUserAddressByID(ctx, param.Uuid)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(address)
}

// @Summary 获取用户地址列表
// @Description 获取用户地址列表
// @Tags 用户地址
// @Accept  json
// @Produce  json
// @Success 200 {object} model.PagedResponse
// @Router /api/v1/user_address/list [post]
func (u *UserAddressController) GetUserAddressList(ctx *app.Context) {

	userId := ctx.GetString("user_id")
	if userId == "" {
		ctx.JSONError(http.StatusBadRequest, "user_id is required")
		return
	}

	addresses, err := u.UserAddressService.GetUserAddressesByUserID(ctx, userId)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(addresses)
}
