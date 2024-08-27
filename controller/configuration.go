package controller

import (
	"net/http"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/service"
)

type ConfigurationController struct {
	ConfigurationService *service.ConfigurationService
}

// @Summary 创建配置
// @Description 创建配置
// @Tags 配置
// @Accept  json
// @Produce  json
// @Param param body model.ReqConfigCreate true "配置参数"
// @Success 200 {object} model.ConfigurationInfoResponse
// @Router /api/v1/configuration/create [post]
func (c *ConfigurationController) CreateConfiguration(ctx *app.Context) {
	var param model.Configuration
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := c.ConfigurationService.CreateOrUpdateConfiguration(ctx, &param); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(param)
}

// @Summary 更新配置
// @Description 更新配置
// @Tags 配置
// @Accept  json
// @Produce  json
// @Param param body model.Configuration true "配置参数"
// @Success 200 {object} model.ConfigurationInfoResponse
// @Router /api/v1/configuration/update [post]
func (c *ConfigurationController) UpdateConfiguration(ctx *app.Context) {
	var param model.Configuration
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := c.ConfigurationService.UpdateConfiguration(ctx, &param); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(param)
}

// @Summary 获取配置信息
// @Description 获取配置信息
// @Tags 配置
// @Accept  json
// @Produce  json
// @Param param body model.ReqConfigQueryParam true "配置ID"
// @Success 200 {object} model.ConfigurationInfoResponse
// @Router /api/v1/configuration/info [post]
func (c *ConfigurationController) GetConfigurationInfo(ctx *app.Context) {
	var param model.ReqConfigQueryParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	config, err := c.ConfigurationService.GetConfigurationByCategoryAndName(ctx, param.Category, param.Name)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(config)
}

// @Summary 获取配置列表
// @Description 获取配置列表
// @Tags 配置
// @Accept  json
// @Produce  json
// @Param param body model.ReqConfigQueryParam true "查询参数"
// @Success 200 {object} model.ConfigurationPageResponse
// @Router /api/v1/configuration/list [post]
func (c *ConfigurationController) GetConfigurationList(ctx *app.Context) {
	param := &model.ReqConfigQueryParam{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	configs, err := c.ConfigurationService.GetConfigurationList(ctx, param)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(configs)
}

// GetConfigurationMapByCategory
// @Summary 获取配置Map
// @Description 获取配置Map
// @Tags 配置
// @Accept  json
// @Produce  json
// @Param param body model.ReqConfigCategoryParam true "查询参数"
// @Success 200 {object} model.ConfigurationMapResponse
// @Router /api/v1/configuration/category_map [post]
func (c *ConfigurationController) GetConfigurationMapByCategory(ctx *app.Context) {
	var param model.ReqConfigQueryParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	configs, err := c.ConfigurationService.GetConfigurationMapByCategory(ctx, param.Category)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(configs)
}

// CreateConfigurationMapByCategory
// @Summary 创建配置Map
// @Description 创建配置Map
// @Tags 配置
// @Accept  json
// @Produce  json
// @Param param body map[string]string true "配置Map参数, key为配置名称, value为配置值, category为配置分类(必须)"
// @Success 200 {object} model.StringDataResponse "ok"
// @Router /api/v1/configuration/category_create_map [post]
func (c *ConfigurationController) CreateConfigurationMapByCategory(ctx *app.Context) {
	mparam := make(map[string]string, 0)
	if err := ctx.ShouldBindJSON(&mparam); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	category, ok := mparam["category"]
	if !ok {
		ctx.JSONError(http.StatusBadRequest, "category is required")
		return
	}
	delete(mparam, "category")

	if err := c.ConfigurationService.CreateOrUpdateConfigurationMap(ctx, mparam, category); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess("ok")
}
