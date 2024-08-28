package controller

import (
	"net/http"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/service"
)

type PermissionController struct {
	PermissionService *service.PermissionService
}

// CreatePermission 创建新的权限
// @Summary 创建新的权限
// @Tags 权限
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.Permission true "权限信息"
// @Success 200 {object} model.PermissionInfoResponse "ok"
// @Router /api/v1/permission/create [post]
func (p *PermissionController) CreatePermission(ctx *app.Context) {
	var param model.Permission
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := p.PermissionService.CreatePermission(ctx, &param); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(param)
}

// UpdatePermission 更新权限信息
// @Summary 更新权限信息
// @Tags 权限
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.Permission true "权限信息"
// @Success 200 {object} model.PermissionInfoResponse "ok"
// @Router /api/v1/permission/update [post]
func (p *PermissionController) UpdatePermission(ctx *app.Context) {
	var param model.Permission
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := p.PermissionService.UpdatePermission(ctx, &param); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(param)
}

// DeletePermission 删除权限
// @Summary 删除权限
// @Tags 权限
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.ReqUuidParam true "权限信息"
// @Success 200 {object} model.StringDataResponse "ok"
// @Router /api/v1/permission/delete [post]
func (p *PermissionController) DeletePermission(ctx *app.Context) {
	var param model.ReqUuidParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := p.PermissionService.DeletePermission(ctx, param.Uuid); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess("ok")
}

// GetPermissionInfo 获取权限信息
// @Summary 获取权限信息
// @Tags 权限
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.ReqUuidParam true "权限信息"
// @Success 200 {object} model.PermissionInfoResponse "ok"
// @Router /api/v1/permission/info [post]
func (p *PermissionController) GetPermissionInfo(ctx *app.Context) {
	var param model.ReqUuidParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	permission, err := p.PermissionService.GetPermissionByUUID(ctx, param.Uuid)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(permission)
}

// GetPermissionList 获取权限列表
// @Summary 获取权限列表
// @Tags 权限
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param params body model.ReqPermissionQueryParam true "权限查询参数"
// @Success 200 {object} model.PermissionListResponse "ok"
// @Router /api/v1/permission/list [post]
func (p *PermissionController) GetPermissionList(ctx *app.Context) {
	var param model.ReqPermissionQueryParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	permissions, err := p.PermissionService.GetPermissionList(ctx, &param)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(permissions)
}
