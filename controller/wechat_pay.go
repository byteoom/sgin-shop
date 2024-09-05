package controller

import (
	"encoding/json"
	"net/http"
	"sgin/pkg/app"
	paymentmethod "sgin/pkg/payment-method"
	"sgin/service"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
)

type WechatPayController struct {
	PaymentMethodService *service.PaymentMethodService
}

// return 回调
func (w *WechatPayController) Return(ctx *app.Context) {

	notifyReq, err := wechat.V3ParseNotify(ctx.Request)
	if err != nil {
		ctx.Logger.Info("ParseNotifyToBodyMap err:", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	payment, err := w.PaymentMethodService.GetPaymentMethodInfo(ctx, "", "wechat")
	if err != nil {
		ctx.Logger.Error("Failed to get payment method info", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	wechatClient := &paymentmethod.Wechat{}
	err = json.Unmarshal([]byte(payment.Config), wechatClient)
	if err != nil {
		ctx.Logger.Error("Failed to unmarshal wechat config", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	client, err := wechat.NewClientV3(wechatClient.Mchid, wechatClient.SerialNo, wechatClient.ApiV3Key, wechatClient.PrivateKey)
	if err != nil {
		ctx.Logger.Error("Failed to create wechat client:", err)
		return
	}

	reqbody, _ := json.Marshal(notifyReq)
	ctx.Logger.Info("notifyReq:", string(reqbody))

	// 验签
	// 获取微信平台证书
	certMap := client.WxPublicKeyMap()
	// 验证异步通知的签名
	err = notifyReq.VerifySignByPKMap(certMap)
	if err != nil {
		ctx.Logger.Error("VerifySignByPKMap err:", err)
		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	// ====↓↓↓====异步通知应答====↓↓↓====
	// 退款通知http应答码为200且返回状态码为SUCCESS才会当做商户接收成功，否则会重试。
	// 注意：重试过多会导致微信支付端积压过多通知而堵塞，影响其他正常通知。

	// 此写法是 gin 框架返回微信的写法
	ctx.JSON(http.StatusOK, &wechat.V3NotifyRsp{Code: gopay.SUCCESS, Message: "成功"})

}
