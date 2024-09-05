package paymentmethod

import (
	"context"
	"encoding/json"
	"fmt"
	"sgin/pkg/app"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
)

// Alipay 支付宝支付
type Alipay struct {
	// AppID：支付宝分配给开发者的应用ID
	// privateKey: 应用私钥
	// isProd：是否是正式环境，沙箱环境请选择新版沙箱应用。
	AppID            string `json:"app_id"`
	PrivateKey       string `json:"private_key"` // 应用私钥
	IsProd           bool   `json:"is_prod"`
	AppPublicCert    string `json:"app_public_cert"`    // 应用公钥证书
	AlipayRootCert   string `json:"alipay_root_cert"`   // 支付宝根证书
	AlipayPublicCert string `json:"alipay_public_cert"` // 支付宝公钥证书
}

// CreateOrder 创建订单
func (a *Alipay) CreateOrder(ctx *app.Context, orderId string, currencyCode string, amount float64, description string) (r *alipay.TradePayResponse, err error) {

	ctx.Logger.Info("Create order id:", orderId)
	ctx.Logger.Info("Create order amount:", amount)
	ctx.Logger.Info("Create order description:", description)

	// 打印支付宝配置信息
	b, _ := json.Marshal(a)
	ctx.Logger.Info("Alipay config:", string(b))

	client, err := alipay.NewClient(a.AppID, a.PrivateKey, false)
	if err != nil {
		ctx.Logger.Error("Failed to create alipay client:", err)
		return
	}

	client.SetLocation(alipay.LocationShanghai). // 设置时区，不设置或出错均为默认服务器时间
							SetCharset(alipay.UTF8).                                               // 设置字符编码，不设置默认 utf-8
							SetSignType(alipay.RSA2).                                              // 设置签名类型，不设置默认 RSA2
							SetReturnUrl("http://sgin-shop.biggerforum.org/api/v1/alipay/return"). // 设置返回URL
							SetNotifyUrl("http://sgin-shop.biggerforum.org/api/v1/alipay/notify"). // 设置异步通知URL
							SetAppAuthToken("")                                                    // 设置第三方应用授权

	// 创建订单

	client.AutoVerifySign([]byte(a.AlipayPublicCert))

	err = client.SetCertSnByContent([]byte(a.AppPublicCert), []byte(a.AlipayRootCert), []byte(a.AlipayPublicCert))
	if err != nil {
		ctx.Logger.Error("Failed to set cert sn by content:", err)
		return
	}

	// 初始化 BodyMap
	bm := make(gopay.BodyMap)
	bm.Set("subject", "条码支付").
		Set("scene", "bar_code").
		Set("product_code", "FAST_INSTANT_TRADE_PAY").
		Set("out_trade_no", orderId).
		Set("total_amount", fmt.Sprintf("%.2f", amount))

	aliRsp, err := client.TradePay(context.Background(), bm)
	if err != nil {
		if bizErr, ok := alipay.IsBizError(err); ok {
			//xlog.Errorf("%+v", bizErr)
			// do something
			ctx.Logger.Errorf("Failed to trade pay: %+v", bizErr)
			return
		}
		//	xlog.Errorf("client.TradePay(%+v),err:%+v", bm, err)
		ctx.Logger.Errorf("Failed to trade pay: %+v， bm: %+v", err, bm)
		return
	}

	r = aliRsp
	return r, nil

}
