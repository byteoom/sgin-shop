package paymentmethod

type Wechat struct {
	// mchid：商户ID 或者服务商模式的 sp_mchid
	// serialNo：商户证书的证书序列号
	// apiV3Key：apiV3Key，商户平台获取
	// privateKey：私钥 apiclient_key.pem 读取后的内容

	Mchid      string `json:"mchid"`
	SerialNo   string `json:"serial_no"`
	ApiV3Key   string `json:"api_v3_key"`
	PrivateKey string `json:"private_key"`
}
