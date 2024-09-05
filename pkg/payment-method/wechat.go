package paymentmethod

import (
	"context"
	"sgin/pkg/app"
	"time"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
)

type Wechat struct {
	// mchid：商户ID 或者服务商模式的 sp_mchid
	// serialNo：商户证书的证书序列号
	// apiV3Key：apiV3Key，商户平台获取
	// privateKey：私钥 apiclient_key.pem 读取后的内容

	Mchid      string `json:"mchid"`       // 商户ID 或者服务商模式的 sp_mchid
	SerialNo   string `json:"serial_no"`   // 商户证书的证书序列号
	ApiV3Key   string `json:"api_v3_key"`  // apiV3Key，商户平台获取
	PrivateKey string `json:"private_key"` // 私钥 apiclient_key.pem 读取后的内容
}

// CreateOrder 创建订单
func (w *Wechat) CreateOrder(ctx *app.Context, orderId string, currencyCode string, amount float64, description string) (r *wechat.PrepayRsp, err error) {

	client, err := wechat.NewClientV3(w.Mchid, w.SerialNo, w.ApiV3Key, w.PrivateKey)
	if err != nil {
		ctx.Logger.Error("Failed to create wechat client:", err)
		return
	}

	// 设置微信平台API证书和序列号（推荐开启自动验签，无需手动设置证书公钥等信息）
	//client.SetPlatformCert([]byte(""), "")

	// 启用自动同步返回验签，并定时更新微信平台API证书（开启自动验签时，无需单独设置微信平台API证书和序列号）
	err = client.AutoVerifySign()
	if err != nil {
		ctx.Logger.Error("Failed to auto verify sign:", err)
		return
	}

	expire := time.Now().Add(10 * time.Minute).Format(time.RFC3339)
	// 初始化 BodyMap
	bm := make(gopay.BodyMap)
	bm.Set("sp_appid", "sp_appid").
		Set("sp_mchid", "sp_mchid").
		Set("sub_mchid", "sub_mchid").
		Set("description", "测试Jsapi支付商品").
		Set("out_trade_no", orderId).
		Set("time_expire", expire).
		Set("notify_url", "http://sgin-shop.biggerforum.org/api/v1/wechat_pay/return").
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", amount).
				Set("currency", currencyCode)
		}).
		SetBodyMap("payer", func(bm gopay.BodyMap) {
			bm.Set("sp_openid", "asdas")
		})

	wxRsp, err := client.V3TransactionJsapi(context.Background(), bm)
	if err != nil {
		ctx.Logger.Error("Failed to create wechat order:", err)
		return
	}
	if wxRsp.Code == wechat.Success {
		ctx.Logger.Info("Create wechat order successfully:", wxRsp.Response)
		r = wxRsp
		return
	}

	ctx.Logger.Error("Failed to create wechat order:", wxRsp.Error, wxRsp.Code)

	return
}
