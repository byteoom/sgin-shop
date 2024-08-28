package controller

import (
	"net/http"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/service"
)

type UserPermissionController struct {
	UserPermissionService *service.UserPermissionService
}

// CreateUserPermission 创建新的用户权限关联
// @Summary 创建新的用户权限关联
// @Tags 用户权限关联
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.ReqPermissionUserCreate true "用户权限关联信息"
// @Success 200 {object} model.StringDataResponse	"ok"
// @Router /api/v1/permission_user/create [post]
func (u *UserPermissionController) CreateUserPermission(ctx *app.Context) {
	var param model.ReqPermissionUserCreate
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := u.UserPermissionService.CreateUserPermission(ctx, &param); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess("ok")
}

// UpdateUserPermission 更新用户权限关联信息
// @Summary 更新用户权限关联信息
// @Tags 用户权限关联
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.UserPermission true "用户权限关联信息"
// @Success 200 {object} model.UserPermissionInfoResponse	"ok"
// @Router /api/v1/permission_user/update [post]
func (u *UserPermissionController) UpdateUserPermission(ctx *app.Context) {
	var param model.UserPermission
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := u.UserPermissionService.UpdateUserPermission(ctx, &param); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(param)
}

// DeleteUserPermission 删除用户权限关联
// @Summary 删除用户权限关联
// @Tags 用户权限关联
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.ReqUuidParam true "用户权限关联信息"
// @Success 200 {object} model.StringDataResponse	"ok"
// @Router /api/v1/permission_user/delete [post]
func (u *UserPermissionController) DeleteUserPermission(ctx *app.Context) {
	var param model.ReqUuidParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := u.UserPermissionService.DeleteUserPermission(ctx, param.Uuid); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess("ok")
}

// GetUserPermissionInfo 获取用户权限关联信息
// @Summary 获取用户权限关联信息
// @Tags 用户权限关联
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.ReqUuidParam true "用户权限关联信息"
// @Success 200 {object} model.UserPermissionListResponse	"ok"
// @Router /api/v1/permission_user/info [post]
func (u *UserPermissionController) GetUserPermissionInfo(ctx *app.Context) {
	var param model.ReqUuidParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	userPermission, err := u.UserPermissionService.GetUserPermissionByUserUUID(ctx, param.Uuid)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(userPermission)
}

// GetUserPermissionList 获取用户权限关联列表
// @Summary 获取用户权限关联列表
// @Tags 用户权限关联
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.ReqUserPermissionQueryParam true "用户权限关联查询参数"
// @Success 200 {object} model.UserPermissionPageResponse	"ok"
// @Router /api/v1/permission_user/list [post]
func (u *UserPermissionController) GetUserPermissionList(ctx *app.Context) {
	var param model.ReqUserPermissionQueryParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	userPermissions, err := u.UserPermissionService.GetUserPermissionList(ctx, &param)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(userPermissions)
}
