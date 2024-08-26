package controller

import (
	"net/http"
	"sgin/model"
	"sgin/pkg/app"
	"sgin/service"
)

type CurrencyController struct {
	CurrencyService *service.CurrencyService
}

// @Summary 创建币种
// @Description 创建币种
// @Tags 币种
// @Accept  json
// @Produce  json
// @Param param body model.Currency true "币种参数"
// @Success 200 {object} model.Currency
// @Router /api/v1/currency/create [post]
func (c *CurrencyController) CreateCurrency(ctx *app.Context) {
	var param model.Currency
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := c.CurrencyService.CreateCurrency(ctx, &param); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(param)
}

// @Summary 更新币种
// @Description 更新币种
// @Tags 币种
// @Accept  json
// @Produce  json
// @Param param body model.Currency true "币种参数"
// @Success 200 {object} model.Currency
// @Router /api/v1/currency/update [post]
func (c *CurrencyController) UpdateCurrency(ctx *app.Context) {
	var param model.Currency
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := c.CurrencyService.UpdateCurrency(ctx, &param); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(param)
}

// @Summary 删除币种
// @Description 删除币种
// @Tags 币种
// @Accept  json
// @Produce  json
// @Param param body model.ReqUuidParam true "币种UUID"
// @Success 200 {string} string	"ok"
// @Router /api/v1/currency/delete [post]
func (c *CurrencyController) DeleteCurrency(ctx *app.Context) {
	var param model.ReqUuidParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	if err := c.CurrencyService.DeleteCurrency(ctx, param.Uuid); err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess("ok")
}

// @Summary 获取币种信息
// @Description 获取币种信息
// @Tags 币种
// @Accept  json
// @Produce  json
// @Param param body model.ReqUuidParam true "币种UUID"
// @Success 200 {object} model.Currency
// @Router /api/v1/currency/info [post]
func (c *CurrencyController) GetCurrencyInfo(ctx *app.Context) {
	var param model.ReqUuidParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}
	currency, err := c.CurrencyService.GetCurrencyByUUID(ctx, param.Uuid)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(currency)
}

// @Summary 获取币种列表
// @Description 获取币种列表
// @Tags 币种
// @Accept  json
// @Produce  json
// @Param param body model.ReqCurrencyQueryParam true "查询参数"
// @Success 200 {object} model.CurrencyQueryResponse
// @Router /api/v1/currency/list [post]
func (c *CurrencyController) GetCurrencyList(ctx *app.Context) {
	param := &model.ReqCurrencyQueryParam{}
	if err := ctx.ShouldBindJSON(param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	currencies, err := c.CurrencyService.GetCurrencyList(ctx, param)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess(currencies)
}

// @Summary 获取全部可用的币种
// @Description 获取全部可用的币种
// @Tags 币种
// @Accept  json
// @Produce  json
// @Success 200 {array} model.Currency
// @Router /api/v1/currency/all [get]
func (c *CurrencyController) GetAllCurrency(ctx *app.Context) {
	currencies, err := c.CurrencyService.GetAvailableCurrencyList(ctx)
	if err != nil {
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSONSuccess(currencies)
}
