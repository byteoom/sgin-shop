package controller

import (
	"io/ioutil"
	"sgin/pkg/app"
)

type PaypalController struct {
}

// return 回调
func (p *PaypalController) Return(ctx *app.Context) {

	// 获取body

	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.Logger.Error("Failed to read body:", err)
		ctx.JSONError(500, "Failed to read body")
		return
	}

	ctx.Logger.Info("Paypal Return:", string(body))

	// 获取参数
	urlPath := ctx.Request.URL.Path
	ctx.Logger.Info("Paypal Return URL Path:", urlPath)

	// 获取method
	method := ctx.Request.Method
	ctx.Logger.Info("Paypal Return Method:", method)

	ctx.JSONSuccess("Paypal Return")
}

// cancel 回调
func (p *PaypalController) Cancel(ctx *app.Context) {

	// 获取body

	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.Logger.Error("Failed to read body:", err)
		ctx.JSONError(500, "Failed to read body")
		return
	}

	ctx.Logger.Info("Paypal Cancel:", string(body))

	// 获取参数
	urlPath := ctx.Request.URL.Path
	ctx.Logger.Info("Paypal Cancel URL Path:", urlPath)

	// 获取method
	method := ctx.Request.Method
	ctx.Logger.Info("Paypal Cancel Method:", method)

	ctx.JSONSuccess("Paypal Cancel")
}
