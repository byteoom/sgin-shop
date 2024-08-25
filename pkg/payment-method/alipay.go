package paymentmethod

// Alipay 支付宝支付
type Alipay struct {
	// AppID：支付宝分配给开发者的应用ID
	// privateKey: 应用私钥
	// isProd：是否是正式环境，沙箱环境请选择新版沙箱应用。
	AppID      string `json:"app_id"`
	PrivateKey string `json:"private_key"`
	IsProd     bool   `json:"is_prod"`
}
