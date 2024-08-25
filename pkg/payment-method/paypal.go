package paymentmethod

type PayPal struct {
	Email      string `json:"email"`       // 收款人邮箱
	MerchantId string `json:"merchant_id"` // 商户ID
	Clientid   string `json:"clientid"`    // 客户端ID
	Secret     string `json:"secret"`      // 客户端密钥
}
