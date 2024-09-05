package controller

import (
	"encoding/json"
	"net/http"
	"sgin/pkg/app"
	paymentmethod "sgin/pkg/payment-method"
	"sgin/service"

	"github.com/go-pay/gopay/alipay"
)

type AlipayController struct {
	PaymentMethodService *service.PaymentMethodService
}

// return 回调
func (a *AlipayController) Return(ctx *app.Context) {
	ctx.Logger.Info("Alipay return")
	// 解析异步通知的参数
	// req：*http.Request
	notifyReq, err := alipay.ParseNotifyToBodyMap(ctx.Request) // c.Request 是 gin 框架的写法
	if err != nil {
		ctx.Logger.Info("ParseNotifyToBodyMap err:", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	payment, err := a.PaymentMethodService.GetPaymentMethodInfo(ctx, "", "alipay")
	if err != nil {
		ctx.Logger.Error("Failed to get payment method info", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	alipayClient := &paymentmethod.Alipay{}
	err = json.Unmarshal([]byte(payment.Config), alipayClient)
	if err != nil {
		ctx.Logger.Error("Failed to unmarshal paypal config", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	reqbody, _ := json.Marshal(notifyReq)
	ctx.Logger.Info("notifyReq:", string(reqbody))

	// 支付宝异步通知验签（公钥证书模式）
	ok, err := alipay.VerifySignWithCert([]byte(alipayClient.AlipayPublicCert), notifyReq)
	if err != nil {
		ctx.Logger.Error("VerifySignWithCert err:", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	if !ok {
		ctx.Logger.Error("VerifySignWithCert failed")
		ctx.JSONError(http.StatusInternalServerError, "VerifySignWithCert failed")
		return
	}

	ctx.Logger.Info("VerifySignWithCert success")

	// 如果需要，可将 BodyMap 内数据，Unmarshal 到指定结构体指针 ptr
	//err = notifyReq.Unmarshal(ptr)

	// ====异步通知，返回支付宝平台的信息====
	// 文档：https://opendocs.alipay.com/open/203/105286
	// 程序执行完后必须打印输出“success”（不包含引号）。如果商户反馈给支付宝的字符不是success这7个字符，支付宝服务器会不断重发通知，直到超过24小时22分钟。一般情况下，25小时以内完成8次通知（通知的间隔频率一般是：4m,10m,10m,1h,2h,6h,15h）

	// 此写法是 gin 框架返回支付宝的写法
	ctx.String(http.StatusOK, "%s", "success")
}

// return 回调
func (a *AlipayController) Notify(ctx *app.Context) {
	ctx.Logger.Info("Alipay Notify")
	// 解析异步通知的参数
	// req：*http.Request
	notifyReq, err := alipay.ParseNotifyToBodyMap(ctx.Request) // c.Request 是 gin 框架的写法
	if err != nil {
		ctx.Logger.Info("ParseNotifyToBodyMap err:", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	payment, err := a.PaymentMethodService.GetPaymentMethodInfo(ctx, "", "alipay")
	if err != nil {
		ctx.Logger.Error("Failed to get payment method info", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	alipayClient := &paymentmethod.Alipay{}
	err = json.Unmarshal([]byte(payment.Config), alipayClient)
	if err != nil {
		ctx.Logger.Error("Failed to unmarshal paypal config", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	reqbody, _ := json.Marshal(notifyReq)
	ctx.Logger.Info("notifyReq:", string(reqbody))

	// 支付宝异步通知验签（公钥证书模式）
	ok, err := alipay.VerifySignWithCert([]byte(alipayClient.AlipayPublicCert), notifyReq)
	if err != nil {
		ctx.Logger.Error("VerifySignWithCert err:", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	if !ok {
		ctx.Logger.Error("VerifySignWithCert failed")
		ctx.JSONError(http.StatusInternalServerError, "VerifySignWithCert failed")
		return
	}

	ctx.Logger.Info("VerifySignWithCert success")

	// 如果需要，可将 BodyMap 内数据，Unmarshal 到指定结构体指针 ptr
	//err = notifyReq.Unmarshal(ptr)

	// ====异步通知，返回支付宝平台的信息====
	// 文档：https://opendocs.alipay.com/open/203/105286
	// 程序执行完后必须打印输出“success”（不包含引号）。如果商户反馈给支付宝的字符不是success这7个字符，支付宝服务器会不断重发通知，直到超过24小时22分钟。一般情况下，25小时以内完成8次通知（通知的间隔频率一般是：4m,10m,10m,1h,2h,6h,15h）

	// 此写法是 gin 框架返回支付宝的写法
	ctx.String(http.StatusOK, "%s", "success")
}
