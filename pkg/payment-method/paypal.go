package paymentmethod

import (
	"encoding/json"
	"fmt"
	"sgin/pkg/app"
	"sgin/pkg/utils"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/paypal"
	"github.com/go-pay/xlog"
)

type PayPal struct {
	Email      string `json:"email"`       // 收款人邮箱
	MerchantId string `json:"merchant_id"` // 商户ID
	Clientid   string `json:"clientid"`    // 客户端ID
	Secret     string `json:"secret"`      // 客户端密钥
	Env        string `json:"env"`         // 环境 sandbox: 沙盒 production: 正式环境
}

// Create Orders
// 创建订单
func (p *PayPal) CreateOrder(ctx *app.Context, orderId string, currencyCode string, amount float64, description string) (r *paypal.OrderDetail, err error) {
	// 创建订单

	// 打印paypal配置信息
	b, _ := json.Marshal(p)
	ctx.Logger.Info("Paypal config:", string(b))

	// 初始化PayPal支付客户端
	client, err := paypal.NewClient(p.Clientid, p.Secret, false)
	if err != nil {
		ctx.Logger.Error("Failed to create paypal client:", err)
		xlog.Error(err)
		return
	}

	token, err := client.GetAccessToken() // 获取 access_token

	if err != nil {
		ctx.Logger.Error("Failed to get access token:", err)
		xlog.Error(err)
		return
	}
	btoken, _ := json.Marshal(token)
	ctx.Logger.Info("Get access token successfully:", string(btoken))

	// 自定义配置http请求接收返回结果body大小，默认 10MB
	// client.SetBodySize() // 没有特殊需求，可忽略此配置

	// 打开Debug开关，输出日志，默认关闭
	client.DebugSwitch = gopay.DebugOn

	ctx.Logger.Info("Create order id:", orderId)

	// Create Orders example
	var pus []*paypal.PurchaseUnit
	var item = &paypal.PurchaseUnit{
		ReferenceId: orderId,
		Amount: &paypal.Amount{
			CurrencyCode: currencyCode,
			Value:        fmt.Sprintf("%.2f", amount),
		},
	}
	pus = append(pus, item)

	bm := make(gopay.BodyMap)
	bm.Set("intent", "CAPTURE").
		Set("purchase_units", pus).
		SetBodyMap("payment_source", func(b gopay.BodyMap) {
			b.SetBodyMap("paypal", func(bb gopay.BodyMap) {
				bb.SetBodyMap("experience_context", func(bbb gopay.BodyMap) {
					bbb.Set("brand_name", "gopay").
						Set("locale", "en-US").
						Set("shipping_preference", "NO_SHIPPING").
						Set("user_action", "PAY_NOW").
						Set("return_url", "http://sgin-shop.biggerforum.org/api/v1/paypal/return").
						Set("cancel_url", "http://sgin-shop.biggerforum.org/api/v1/paypal/cancel")
				})
			})
		})

	reqbody := bm.JsonBody()

	ctx.Logger.Info("Create order request body:", reqbody)

	ppRsp, err := client.CreateOrder(ctx, bm)
	if err != nil {
		xlog.Error(err)
		ctx.Logger.Error("Failed to create order", err)
		return
	}
	b, _ = json.Marshal(ppRsp)
	if ppRsp.Code != 200 {
		// do something

		ctx.Logger.Error("Failed to create order : ", ppRsp.Error, ppRsp.Code)
		ctx.Logger.Error("Failed to create order : ", string(b))
		return
	}

	r = ppRsp.Response
	ctx.Logger.Info("Create order successfully:", string(b))
	return r, nil

}

// Create Orders
// 创建订单
func (p *PayPal) CreateSandBoxOrder(ctx *app.Context, currencyCode string, amount float64, description string) (r *paypal.OrderDetail, err error) {
	// 创建订单

	// 打印paypal配置信息
	b, _ := json.Marshal(p)
	ctx.Logger.Info("Paypal config:", string(b))

	// 初始化PayPal支付客户端
	client, err := paypal.NewClient(p.Clientid, p.Secret, false)
	if err != nil {
		ctx.Logger.Error("Failed to create paypal client:", err)
		xlog.Error(err)
		return
	}

	token, err := client.GetAccessToken() // 获取 access_token

	if err != nil {
		ctx.Logger.Error("Failed to get access token:", err)
		xlog.Error(err)
		return
	}
	btoken, _ := json.Marshal(token)
	ctx.Logger.Info("Get access token successfully:", string(btoken))

	// 自定义配置http请求接收返回结果body大小，默认 10MB
	// client.SetBodySize() // 没有特殊需求，可忽略此配置

	// 打开Debug开关，输出日志，默认关闭
	client.DebugSwitch = gopay.DebugOn

	orderId := utils.GenerateOrderID()

	ctx.Logger.Info("Create order id:", orderId)

	// Create Orders example
	var pus []*paypal.PurchaseUnit
	var item = &paypal.PurchaseUnit{
		ReferenceId: orderId,
		Amount: &paypal.Amount{
			CurrencyCode: currencyCode,
			Value:        fmt.Sprintf("%.2f", amount),
		},
	}
	pus = append(pus, item)

	bm := make(gopay.BodyMap)
	bm.Set("intent", "CAPTURE").
		Set("purchase_units", pus).
		SetBodyMap("payment_source", func(b gopay.BodyMap) {
			b.SetBodyMap("paypal", func(bb gopay.BodyMap) {
				bb.SetBodyMap("experience_context", func(bbb gopay.BodyMap) {
					bbb.Set("brand_name", "gopay").
						Set("locale", "en-US").
						Set("shipping_preference", "NO_SHIPPING").
						Set("user_action", "PAY_NOW").
						Set("return_url", "http://sgin-shop.biggerforum.org/api/v1/paypal/return").
						Set("cancel_url", "http://sgin-shop.biggerforum.org/api/v1/paypal/cancel")
				})
			})
		})

	reqbody := bm.JsonBody()

	ctx.Logger.Info("Create order request body:", reqbody)

	ppRsp, err := client.CreateOrder(ctx, bm)
	if err != nil {
		xlog.Error(err)
		ctx.Logger.Error("Failed to create order", err)
		return
	}
	b, _ = json.Marshal(ppRsp)
	if ppRsp.Code != 200 {
		// do something

		ctx.Logger.Error("Failed to create order : ", ppRsp.Error, ppRsp.Code)
		ctx.Logger.Error("Failed to create order : ", string(b))
		return
	}

	r = ppRsp.Response
	ctx.Logger.Info("Create order successfully:", string(b))
	return r, nil

}
