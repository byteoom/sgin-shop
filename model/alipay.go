package model

type ErrorResponse struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	SubCode string `json:"sub_code,omitempty"`
	SubMsg  string `json:"sub_msg,omitempty"`
}

type TradePay struct {
	ErrorResponse
	TradeNo             string           `json:"trade_no,omitempty"`
	OutTradeNo          string           `json:"out_trade_no,omitempty"`
	BuyerLogonId        string           `json:"buyer_logon_id,omitempty"`
	TotalAmount         string           `json:"total_amount,omitempty"`
	ReceiptAmount       string           `json:"receipt_amount,omitempty"`
	BuyerPayAmount      string           `json:"buyer_pay_amount,omitempty"`
	PointAmount         string           `json:"point_amount,omitempty"`
	InvoiceAmount       string           `json:"invoice_amount,omitempty"`
	FundBillList        []*TradeFundBill `json:"fund_bill_list"`
	StoreName           string           `json:"store_name,omitempty"`
	BuyerUserId         string           `json:"buyer_user_id,omitempty"`
	BuyerOpenId         string           `json:"buyer_open_id,omitempty"`
	DiscountGoodsDetail string           `json:"discount_goods_detail,omitempty"`
	AsyncPaymentMode    string           `json:"async_payment_mode,omitempty"`
	VoucherDetailList   []*VoucherDetail `json:"voucher_detail_list"`
	AdvanceAmount       string           `json:"advance_amount,omitempty"`
	AuthTradePayMode    string           `json:"auth_trade_pay_mode,omitempty"`
	MdiscountAmount     string           `json:"mdiscount_amount,omitempty"`
	DiscountAmount      string           `json:"discount_amount,omitempty"`
	CreditPayMode       string           `json:"credit_pay_mode"`
	CreditBizOrderId    string           `json:"credit_biz_order_id"`
}

type TradeFundBill struct {
	FundChannel string `json:"fund_channel,omitempty"` // 同步通知里是 fund_channel
	Amount      string `json:"amount,omitempty"`
	RealAmount  string `json:"real_amount,omitempty"`
	FundType    string `json:"fund_type,omitempty"`
}

type VoucherDetail struct {
	Id                         string `json:"id,omitempty"`
	Name                       string `json:"name,omitempty"`
	Type                       string `json:"type,omitempty"`
	Amount                     string `json:"amount,omitempty"`
	MerchantContribute         string `json:"merchant_contribute,omitempty"`
	OtherContribute            string `json:"other_contribute,omitempty"`
	Memo                       string `json:"memo,omitempty"`
	TemplateId                 string `json:"template_id,omitempty"`
	PurchaseBuyerContribute    string `json:"purchase_buyer_contribute,omitempty"`
	PurchaseMerchantContribute string `json:"purchase_merchant_contribute,omitempty"`
	PurchaseAntContribute      string `json:"purchase_ant_contribute,omitempty"`
}
