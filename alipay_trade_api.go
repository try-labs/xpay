package alipay

import (
	"encoding/json"
	"fmt"
)

/**
 * @author: Sam
 * @since:
 * @date: 2023/3/15 17:52
 * @desc:
 */

// 强制TradePagePayReq实现接口IAliPayRequest，另外的一种方式见github.com/sirupsen/logrus@v1.9.0/logrus.go:116 => _ StdLogger = &Entry{}
// github.com/gin-gonic/gin@v1.8.1/routergroup.go:61 =>var _ IRouter = &RouterGroup{}
// go/go1.20.1/src/encoding/json/stream.go:281 => var _ Marshaler = (*RawMessage)(nil)
var _ IAliPayRequest = (*TradePagePayReq)(nil)

type TradePagePayReq struct {
	OutTradeNo  string `json:"out_trade_no"`          // 	必选	64 商户订单号。 由商家自定义，64个字符以内，仅支持字母、数字、下划线且需保证在商户端不重复。
	TotalAmount string `json:"total_amount"`          // 	必选	11 订单总金额，单位为元，精确到小数点后两位，取值范围为 [0.01,100000000]。金额不能为0。
	Subject     string `json:"subject"`               // 	必选	256 订单标题。注意：不可使用特殊字符，如 /，=，& 等。
	ProductCode string `json:"product_code"`          // 	必选	64 销售产品码，与支付宝签约的产品码名称。注：目前电脑支付场景下仅支持FAST_INSTANT_TRADE_PAY
	QrPayMode   string `json:"qr_pay_mode,omitempty"` // 	可选	2 销售产品码，PC扫码支付的方式。
	//支持前置模式和跳转模式。前置模式是将二维码前置到商户的订单确认页的模式。需要商户在自己的页面中以 iframe 方式请求支付宝页面。具体支持的枚举值有以下几种：
	//0：订单码-简约前置模式，对应 iframe 宽度不能小于600px，高度不能小于300px；
	//1：订单码-前置模式，对应iframe 宽度不能小于 300px，高度不能小于600px；
	//3：订单码-迷你前置模式，对应 iframe 宽度不能小于 75px，高度不能小于75px；
	//4：订单码-可定义宽度的嵌入式二维码，商户可根据需要设定二维码的大小。
	//
	//跳转模式下，用户的扫码界面是由支付宝生成的，不在商户的域名下。支持传入的枚举值有：
	//2：订单码-跳转模式
	QrcodeWidth     string         `json:"qrcode_width,omitempty"`      // 可选 4 商户自定义二维码宽度。注：qr_pay_mode=4时该参数有效
	GoodsDetail     []*GoodsDetail `json:"goods_detail,omitempty"`      // 可选 订单包含的商品列表信息，json格式。
	TimeExpire      string         `json:"time_expire,omitempty"`       // 可选 订单绝对超时时间。 格式为yyyy-MM-dd HH:mm:ss。超时时间范围：1m~15d。 注：time_expire和timeout_express两者只需传入一个或者都不传，两者均传入时，优先使用time_expire。
	SubMerchant     *SubMerchant   `json:"sub_merchant,omitempty"`      // 可选  二级商户信息。 直付通模式和机构间连模式下必传，其它场景下不需要传入。
	ExtendParams    *ExtendParams  `json:"extend_params,omitempty"`     // 可选 业务扩展参数
	BusinessParams  string         `json:"business_params,omitempty"`   // 可选	512 商户传入业务信息，具体值要和支付宝约定，应用于安全，营销等参数直传场景，格式为json格式 {"data":"123"}
	PromoParams     string         `json:"promo_params,omitempty"`      // 可选	512 优惠参数。为 JSON 格式。注：仅与支付宝协商后可用 {"storeIdType":"1"}
	IntegrationType string         `json:"integration_type,omitempty"`  // 可选	16 请求后页面的集成方式。枚举值： ALIAPP：支付宝钱包内 PCWEB：PC端访问 默认值为PCWEB。
	RequestFromUrl  string         `json:"request_from_url,omitempty"`  // 可选	256 请求来源地址。如果使用ALIAPP的集成方式，用户中途取消支付会返回该地址。 https://
	StoreId         string         `json:"store_id,omitempty"`          // 可选	32 商户门店编号。指商户创建门店时输入的门店编号。
	MerchantOrderNo string         `json:"merchant_order_no,omitempty"` // 可选	32 商户原始订单号，最大长度限制 32 位
	ExtUserInfo     *ExtUserInfo   `json:"ext_user_info,omitempty"`     // 可选    外部指定买家
	InvoiceInfo     *InvoiceInfo   `json:"invoice_info,omitempty"`      // 可选    开票信息
	// 自己添加
	ReturnUrl string `json:"-" url:"-"` // 可选	256 HTTP/HTTPS开头字符串
	NotifyUrl string `json:"-" url:"-"` // 可选	256	 支付宝服务器主动通知商户服务器里指定的页面http/https路径。
	BaseAliPayRequest
}

func (r *TradePagePayReq) DoValidate() error {
	outTradeNoLen := len(r.OutTradeNo)
	if outTradeNoLen == 0 || outTradeNoLen > 64 {
		return fmt.Errorf("当前商户订单号长度为%d，不符合长度1-64要求", outTradeNoLen)
	}
	subjectLen := len(r.Subject)
	if subjectLen == 0 || subjectLen > 256 {
		return fmt.Errorf("当前订单标题长度为%d，不符合长度1-256要求", subjectLen)
	}
	if r.ProductCode != FastInstantTradePay {
		return fmt.Errorf("销售产品码必须是%s", FastInstantTradePay)
	}
	if len(r.GoodsDetail) > 0 {
		for _, detail := range r.GoodsDetail {
			if err := detail.DoValidate(); err != nil {
				return err
			}
		}
	}
	return nil
}

//func (r *TradePagePayReq) RequestHttpMethod() string {
//	return ""
//}

func (r *TradePagePayReq) RequestApi() string {
	return "alipay.trade.page.pay"
}

func NewTradePagePayReq(outTradeNo, totalAmount, subject string, opts ...TradePagePayOpt) *TradePagePayReq {
	req := &TradePagePayReq{
		OutTradeNo:  outTradeNo,
		TotalAmount: totalAmount,
		Subject:     subject,
		ProductCode: FastInstantTradePay,
	}
	for _, opt := range opts {
		opt(req)
	}
	return req

}

type InvoiceInfo struct {
	KeyInfo InvoiceKeyInfo `json:"key_info"` // 必选	200	 开票关键信息
	Details string         `json:"details"`  // 必选	400	 开票内容 注：json数组格式 [{"code":"100294400","name":"服饰","num":"2","sumPrice":"200.00","taxRate":"6%"}]
}
type InvoiceKeyInfo struct {
	IsSupportInvoice    bool   `json:"is_support_invoice"`    // 必选	5 该交易是否支持开票 true
	InvoiceMerchantName string `json:"invoice_merchant_name"` // 必选	80 开票商户名称：商户品牌简称|商户门店简称 ABC|003
	TaxNum              string `json:"tax_num"`               // 必选	30 税号 1464888883494
}
type SubMerchant struct {
	MerchantId   string `json:"merchant_id"`             // 必选	16 间连受理商户的支付宝商户编号，通过间连商户入驻后得到。间连业务下必传，并且需要按规范传递受理商户编号。
	MerchantType string `json:"merchant_type,omitempty"` // 可选	32 二级商户编号类型。 枚举值： alipay:支付宝分配的间联商户编号； 目前仅支持alipay，默认可以不传
}

type GoodsDetail struct {
	GoodsId        string `json:"goods_id"`                  // 必选	32	 商品的编号
	AliPayGoodsId  string `json:"alipay_goods_id,omitempty"` // 可选	32	 支付宝定义的统一商品编号
	GoodsName      string `json:"goods_name"`                // 必选	256  商品名称
	Quantity       int    `json:"quantity"`                  // 必选	10	 商品数量
	Price          int    `json:"price"`                     // 必选	9	 商品单价，单位为元
	GoodsCategory  string `json:"goods_category,omitempty"`  // 可选	24	 商品类目
	CategoriesTree string `json:"categories_tree,omitempty"` // 可选	128	 商品类目树，从商品类目根节点到叶子节点的类目id组成，类目id值使用|分割
	Body           string `json:"body,omitempty"`            // alipay.trade.wap.pay(手机网站支付接口2.0)专属,可选	1000 商品描述信息 特价手机
	ShowURL        string `json:"show_url,omitempty"`        // 可选 400 商品的展示地址
}

func (r *GoodsDetail) DoValidate() error {
	goodsIdLen := len(r.GoodsId)
	if goodsIdLen == 0 || goodsIdLen > 32 {
		return fmt.Errorf("当前商品的编号长度为%d，不符合长度1-64要求", goodsIdLen)
	}
	showUrlLen := len(r.ShowURL)
	if showUrlLen > 0 && showUrlLen > 400 {
		return fmt.Errorf("当前商品的展示地址长度为%d，不符合长度1-128要求", showUrlLen)
	}
	if r.Quantity <= 0 || r.Quantity > 10000000000 {
		return fmt.Errorf("当前商品数量为%d，不符合最大1-10000000000的数量要求", r.Quantity)
	}
	return nil
}

func (r *GoodsDetail) SetAliPayGoodsId(alipayGoodsId string) *GoodsDetail {
	r.AliPayGoodsId = alipayGoodsId
	return r
}

func (r *GoodsDetail) SetCategoriesTree(categoriesTree string) *GoodsDetail {
	r.CategoriesTree = categoriesTree
	return r
}

func (r *GoodsDetail) SetShowURL(showUrl string) *GoodsDetail {
	r.ShowURL = showUrl
	return r
}

func NewGoodsDetail(goodsId, goodsName string, quantity, price int) *GoodsDetail {
	goods := &GoodsDetail{
		GoodsId:   goodsId,
		GoodsName: goodsName,
		Quantity:  quantity,
		Price:     price,
	}
	return goods
}

type TradePagePayOpt func(req *TradePagePayReq)

func WithGoodsDetail(goodsDetail []*GoodsDetail) TradePagePayOpt {
	return func(r *TradePagePayReq) {
		r.GoodsDetail = goodsDetail
	}
}

func WithProductCode(productCode string) TradePagePayOpt {
	return func(r *TradePagePayReq) {
		r.ProductCode = productCode
	}
}

func WithQrcodeWidth(qrcodeWidth string) TradePagePayOpt {
	return func(r *TradePagePayReq) {
		r.QrcodeWidth = qrcodeWidth
	}
}

// ///////////////////////////////////////////

type TradeQueryReq struct {
	OutTradeNo   string   `json:"out_trade_no,omitempty"`  // 特殊可选 64 商户订单号。订单支付时传入的商户订单号,和支付宝交易号不能同时为空。trade_no,out_trade_no如果同时存在优先取trade_no
	TradeNo      string   `json:"trade_no,omitempty"`      // 特殊可选 64 支付宝交易号，和商户订单号不能同时为空
	OrgPid       string   `json:"org_pid,omitempty"`       // 可选 16 银行间联模式下有用，其它场景请不要使用；双联通过该参数指定需要查询的交易所属收单机构的pid;
	QueryOptions []string `json:"query_options,omitempty"` // 可选 1024
	//查询选项，商户传入该参数可定制本接口同步响应额外返回的信息字段，数组格式。支持枚举如下：trade_settle_info：返回的交易结算信息，包含分账、补差等信息；
	//fund_bill_list：交易支付使用的资金渠道；
	//voucher_detail_list：交易支付时使用的所有优惠券信息；
	//discount_goods_detail：交易支付所使用的单品券优惠的商品优惠信息；
	//mdiscount_amount：商家优惠金额；
	BaseAliPayRequest
}

func (r *TradeQueryReq) DoValidate() error {
	if len(r.OutTradeNo) == 0 && len(r.TradeNo) == 0 {
		return fmt.Errorf("商户订单号out_trade_no和支付宝交易号trade_no不能同时为空")
	}
	return nil
}

var _ IAliPayRequest = (*TradeQueryReq)(nil)

//func (r *TradeQueryReq) RequestHttpMethod() string {
//	return "POST"
//}

func (r *TradeQueryReq) RequestApi() string {
	return "alipay.trade.query"
}

type TradeQueryRes struct {
	TradeQueryResContent `json:"alipay_trade_query_response"`
	SignCertSn
}

func (r *TradeQueryRes) String() string {
	buff, _ := json.Marshal(r)
	return string(buff)
}

type TradeQueryResContent struct {
	CommonRes
	OutTradeNo            string           `json:"out_trade_no"`                       // 必选 64 商户订单号。订单支付时传入的商户订单号,和支付宝交易号不能同时为空。trade_no,out_trade_no如果同时存在优先取trade_no
	TradeNo               string           `json:"trade_no"`                           // 必选 64 支付宝交易号，和商户订单号不能同时为空
	TradeStatus           TradeStatus      `json:"trade_status"`                       // 必选 32 交易状态：WAIT_BUYER_PAY（交易创建，等待买家付款）、TRADE_CLOSED（未付款交易超时关闭，或支付完成后全额退款）、TRADE_SUCCESS（交易支付成功）、TRADE_FINISHED（交易结束，不可退款）
	TotalAmount           string           `json:"total_amount"`                       // 必选 11 交易的订单金额，单位为元，两位小数。该参数的值为支付时传入的total_amount
	TransCurrency         string           `json:"trans_currency,omitempty"`           // 可选 标价币种 8，该参数的值为支付时传入的trans_currency
	SettleCurrency        string           `json:"settle_currency,omitempty"`          // 可选 8 订单结算币种 对应支付接口传入的settle_currency，支持英镑：GBP、港币：HKD、美元：USD、新加坡元：SGD、日元：JPY、加拿大元：CAD、澳元：AUD、欧元：EUR、新西兰元：NZD、韩元：KRW、泰铢：THB、瑞士法郎：CHF、瑞典克朗：SEK、丹麦克朗：DKK、挪威克朗：NOK、马来西亚林吉特：MYR、印尼卢比：IDR、菲律宾比索：PHP、毛里求斯卢比：MUR、以色列新谢克尔：ILS、斯里兰卡卢比：LKR、俄罗斯卢布：RUB、阿联酋迪拉姆：AED、捷克克朗：CZK、南非兰特：ZAR
	SettleAmount          string           `json:"settle_amount,omitempty"`            // 可选 11 结算币种订单金额
	PayCurrency           string           `json:"pay_currency,omitempty"`             // 可选 8 订单支付币种
	PayAmount             string           `json:"pay_amount,omitempty"`               // 可选 11 支付币种订单金额
	SettleTransRate       string           `json:"settle_trans_rate,omitempty"`        // 可选 11 结算币种兑换标价币种汇率
	TransPayRate          string           `json:"trans_pay_rate,omitempty"`           // 可选 11 标价币种兑换支付币种汇率
	BuyerPayAmount        string           `json:"buyer_pay_amount,omitempty"`         // 可选 11 买家实付金额，单位为元，两位小数。该金额代表该笔交易买家实际支付的金额，不包含商户折扣等金额
	PointAmount           string           `json:"point_amount,omitempty"`             // 可选 11 积分支付的金额，单位为元，两位小数。该金额代表该笔交易中用户使用积分支付的金额，比如集分宝或者支付宝实时优惠等
	InvoiceAmount         string           `json:"invoice_amount,omitempty"`           // 可选 11 交易中用户支付的可开具发票的金额，单位为元，两位小数。该金额代表该笔交易中可以给用户开具发票的金额
	SendPayDate           string           `json:"send_pay_date,omitempty"`            // 可选 32 本次交易打款给卖家的时间
	ReceiptAmount         string           `json:"receipt_amount,omitempty"`           // 可选 11 实收金额，单位为元，两位小数。该金额为本笔交易，商户账户能够实际收到的金额
	StoreId               string           `json:"store_id,omitempty"`                 // 可选 32 商户门店编号
	TerminalId            string           `json:"terminal_id,omitempty"`              // 可选 32 商户机具终端编号
	FundBillList          []*FundBill      `json:"fund_bill_list,omitempty"`           // 可选 交易支付使用的资金渠道。只有在签约中指定需要返回资金明细，或者入参的query_options中指定时才返回该字段信息。
	StoreName             string           `json:"store_name,omitempty"`               // 可选 512 请求交易支付中的商户店铺的名称
	BuyerUserId           string           `json:"buyer_user_id,omitempty"`            // 可选 16 买家在支付宝的用户id
	IndustrySepcDetailGov string           `json:"industry_sepc_detail_gov,omitempty"` // 可选 4096 行业特殊信息-统筹相关
	IndustrySepcDetailAcc string           `json:"industry_sepc_detail_acc,omitempty"` // 可选 4096 行业特殊信息-个账相关
	ChargeAmount          string           `json:"charge_amount,omitempty"`            // 可选 11 该笔交易针对收款方的收费金额；
	ChargeFlags           string           `json:"charge_flags,omitempty"`             // 可选 64 费率活动标识，当交易享受活动优惠费率时，返回该活动的标识；
	SettlementId          string           `json:"settlement_id,omitempty"`            // 可选 64 支付清算编号，用于清算对账使用；
	TradeSettleInfo       *TradeSettleInfo `json:"trade_settle_info,omitempty"`        // 可选 返回的交易结算信息，包含分账、补差等信息
	AuthTradePayMode      string           `json:"auth_trade_pay_mode,omitempty"`      // 可选 64 预授权支付模式，该参数仅在信用预授权支付场景下返回。信用预授权支付：CREDIT_PREAUTH_PAY
	BuyerUserType         string           `json:"buyer_user_type,omitempty"`          // 可选 18 买家用户类型。CORPORATE:企业用户；PRIVATE:个人用户。
	MdiscountAmount       string           `json:"mdiscount_amount,omitempty"`         // 可选 11 商家优惠金额
	DiscountAmount        string           `json:"discount_amount,omitempty"`          // 可选 11 平台优惠金额
	BuyerUserName         string           `json:"buyer_user_name,omitempty"`          // 可选 买家名称；
	Subject               string           `json:"subject,omitempty"`                  // 可选 256 订单标题；
	Body                  string           `json:"body,omitempty"`                     // 可选 1000 订单描述;只在银行间联交易场景下返回该信息
	AlipaySubMerchantId   string           `json:"alipay_sub_merchant_id,omitempty"`   // 可选 32 间连商户在支付宝端的商户编号；
	ExtInfos              string           `json:"ext_infos,omitempty"`                // 可选 1024 交易额外信息，特殊场景下与支付宝约定返回。
	PassbackParams        string           `json:"passback_params,omitempty"`          // 可选 512 公用回传参数。返回支付时传入的passback_params参数信息
	DiscountGoodsDetail   string           `json:"discount_goods_detail"`              // 本次交易支付所使用的单品券优惠的商品优惠信息
	VoucherDetailList     []*VoucherDetail `json:"voucher_detail_list,omitempty"`      // 本交易支付时使用的所有优惠券信息
}
type FundBill struct {
	FundChannel string `json:"fund_channel"`          // 交易使用的资金渠道，详见 支付渠道列表 https://opendocs.alipay.com/open/common/103259
	BankCode    string `json:"bank_code"`             // 银行卡支付时的银行代码
	Amount      string `json:"amount"`                // 该支付工具类型所使用的金额
	RealAmount  string `json:"real_amount,omitempty"` // 可选 11	渠道实际付款金额
}

type VoucherDetail struct {
	Id                         string `json:"id"`                                     // 券id
	Name                       string `json:"name"`                                   // 券名称
	Type                       string `json:"type"`                                   // 当前有三种类型： ALIPAY_FIX_VOUCHER - 全场代金券, ALIPAY_DISCOUNT_VOUCHER - 折扣券, ALIPAY_ITEM_VOUCHER - 单品优惠
	Amount                     string `json:"amount"`                                 // 优惠券面额，它应该会等于商家出资加上其他出资方出资
	MerchantContribute         string `json:"merchant_contribute"`                    // 商家出资（特指发起交易的商家出资金额）
	OtherContribute            string `json:"other_contribute"`                       // 其他出资方出资金额，可能是支付宝，可能是品牌商，或者其他方，也可能是他们的一起出资
	Memo                       string `json:"memo"`                                   // 优惠券备注信息
	TemplateId                 string `json:"template_id,omitempty"`                  // 可选	64 券模板id
	PurchaseBuyerContribute    string `json:"purchase_buyer_contribute,omitempty"`    // 可选	8 如果使用的这张券是用户购买的，则该字段代表用户在购买这张券时用户实际付款的金额 2.01
	PurchaseMerchantContribute string `json:"purchase_merchant_contribute,omitempty"` // 可选	8 如果使用的这张券是用户购买的，则该字段代表用户在购买这张券时商户优惠的金额 1.03
	PurchaseAntContribute      string `json:"purchase_ant_contribute,omitempty"`      // 可选	8 如果使用的这张券是用户购买的，则该字段代表用户在购买这张券时平台优惠的金额 0.82
}

type TradeSettleInfo struct {
	TradeSettleDetailList []*TradeSettleDetail `json:"trade_settle_detail_list"`
}

type TradeSettleDetail struct {
	OperationType     string `json:"operation_type"`
	OperationSerialNo string `json:"operation_serial_no"`
	OperationDate     string `json:"operation_dt"`
	TransOut          string `json:"trans_out"`
	TransIn           string `json:"trans_in"`
	Amount            string `json:"amount"`
}

var _ IAliPayRequest = (*TradeCloseReq)(nil)

type TradeCloseReq struct {
	OutTradeNo   string `json:"out_trade_no,omitempty"` // 特殊可选 64 商户订单号。订单支付时传入的商户订单号,和支付宝交易号不能同时为空。trade_no,out_trade_no如果同时存在优先取trade_no
	TradeNo      string `json:"trade_no,omitempty"`     // 特殊可选 64 支付宝交易号，和商户订单号不能同时为空
	QueryOptions string `json:"operator_id,omitempty"`  // 可选 28 商家操作员编号 id，由商家自定义。

	// 自行添加
	NotifyUrl string `json:"notify_url,omitempty" url:"notify_url,omitempty"` // 可选	256	 支付宝服务器主动通知商户服务器里指定的页面http/https路径。
	BaseAliPayRequest
}

func (r *TradeCloseReq) DoValidate() error {
	if len(r.OutTradeNo) == 0 && len(r.TradeNo) == 0 {
		return fmt.Errorf("商户订单号out_trade_no和支付宝交易号trade_no不能同时为空")
	}
	return nil
}

func (r *TradeCloseReq) RequestApi() string {
	return "alipay.trade.close"
}

type TradeCloseRes struct {
	TradeCloseResContent `json:"alipay_trade_close_response"`
	SignCertSn
}

type TradeCloseResContent struct {
	CommonRes
	OutTradeNo string `json:"out_trade_no,omitempty"` // 特殊可选 64 商户订单号。订单支付时传入的商户订单号,和支付宝交易号不能同时为空。trade_no,out_trade_no如果同时存在优先取trade_no
	TradeNo    string `json:"trade_no,omitempty"`     // 特殊可选 64 支付宝交易号，和商户订单号不能同时为空
}

func (r *TradeCloseRes) String() string {
	buff, _ := json.Marshal(r)
	return string(buff)
}

//////////////////////////////////////////////////////////////////////

type TradeRefundReq struct {
	OutTradeNo              string                          `json:"out_trade_no,omitempty"`              // 特殊可选	64 商户订单号。 订单支付时传入的商户订单号，商家自定义且保证商家系统中唯一。与支付宝交易号 trade_no 不能同时为空。
	TradeNo                 string                          `json:"trade_no,omitempty"`                  // 特殊可选	64 支付宝交易号。 和商户订单号 out_trade_no 不能同时为空。
	RefundAmount            string                          `json:"refund_amount"`                       // 必选	11  退款金额。需要退款的金额，该金额不能大于订单金额，单位为元，支持两位小数。
	RefundReason            string                          `json:"refund_reason,omitempty"`             // 可选	256 退款原因说明。 商家自定义，将在会在商户和用户的pc退款账单详情中展示
	OutRequestNo            string                          `json:"out_request_no,omitempty"`            // 可选	64 退款请求号。 标识一次退款请求，需要保证在交易号下唯一，如需部分退款，则此参数必传。 注：针对同一次退款请求，如果调用接口失败或异常了，重试时需要保证退款请求号不能变更，防止该笔交易重复退款。支付宝会保证同样的退款请求号多次请求只会退一次。
	RefundRoyaltyParameters []*OpenApiRoyaltyDetailInfoPojo `json:"refund_royalty_parameters,omitempty"` // 退分账明细信息。
	//注： 1.当面付且非直付通模式无需传入退分账明细，系统自动按退款金额与订单金额的比率，从收款方和分账收入方退款，不支持指定退款金额与退款方。
	//2.直付通模式，电脑网站支付，手机 APP 支付，手机网站支付产品，须在退款请求中明确是否退分账，从哪个分账收入方退，退多少分账金额；如不明确，默认从收款方退款，收款方余额不足退款失败。不支持系统按比率退款。
	QueryOptions []string `json:"query_options,omitempty"` // 可选	1024 查询选项。 商户通过上送该参数来定制同步需要额外返回的信息字段，数组格式。支持：refund_detail_item_list：退款使用的资金渠道；deposit_back_info：触发银行卡冲退信息通知；
	BaseAliPayRequest
}

type OpenApiRoyaltyDetailInfoPojo struct {
	RoyaltyType  string `json:"royalty_type,omitempty"`   // 可选	32	 分账类型. 普通分账为：transfer; 补差为：replenish; 为空默认为分账transfer; transfer
	TransOut     string `json:"trans_out,omitempty"`      // 可选	16 支出方账户。如果支出方账户类型为userId，本参数为支出方的支付宝账号对应的支付宝唯一用户号，以2088开头的纯16位数字；如果支出方类型为loginName，本参数为支出方的支付宝登录号。 泛金融类商户分账时，该字段不要上送。
	TransOutType string `json:"trans_out_type,omitempty"` // 可选	64	 支出方账户类型。userId表示是支付宝账号对应的支付宝唯一用户号;loginName表示是支付宝登录号； 泛金融类商户分账时，该字段不要上送。
	TransInType  string `json:"trans_in_type,omitempty"`  // 可选	64	 收入方账户类型。userId表示是支付宝账号对应的支付宝唯一用户号;cardAliasNo表示是卡编号;loginName表示是支付宝登录号； userId
	TransIn      string `json:"trans_in"`                 // 必选	16	 收入方账户。如果收入方账户类型为userId，本参数为收入方的支付宝账号对应的支付宝唯一用户号，以2088开头的纯16位数字；如果收入方类型为cardAliasNo，本参数为收入方在支付宝绑定的卡编号；如果收入方类型为loginName，本参数为收入方的支付宝登录号；
	Amount       string `json:"amount,omitempty"`         // 可选	9	 分账的金额，单位为元 0.1
	Desc         string `json:"desc,omitempty"`           // 可选	1000 分账描述 分账给2088101126708402
	RoyaltyScene string `json:"royalty_scene,omitempty"`  // 可选	256	 可选值：达人佣金、平台服务费、技术服务费、其他 达人佣金
	TransInName  string `json:"trans_in_name,omitempty"`  // 可选	64	 分账收款方姓名，上送则进行姓名与支付宝账号的一致性校验，校验不一致则分账失败。不上送则不进行姓名校验 张三
}

var _ IAliPayRequest = (*TradeRefundReq)(nil)

func (r TradeRefundReq) RequestApi() string {
	return "alipay.trade.refund"
}

func (r TradeRefundReq) DoValidate() error {
	if len(r.OutTradeNo) == 0 && len(r.TradeNo) == 0 {
		return fmt.Errorf("商户订单号out_trade_no和支付宝交易号trade_no不能同时为空")
	}
	return nil
}

type TradeRefundRes struct {
	TradeRefundResContent `json:"alipay_trade_refund_response"`
	SignCertSn
}

func (r *TradeRefundRes) String() string {
	buff, _ := json.Marshal(r)
	return string(buff)
}

type TradeRefundResContent struct {
	CommonRes
	TradeNo              string              `json:"trade_no"`                          // 必选	64	支付宝交易号
	OutTradeNo           string              `json:"out_trade_no"`                      // 必选	64	商户订单号
	BuyerLogonId         string              `json:"buyer_logon_id,omitempty"`          // 必选	100 用户的登录id
	FundChange           string              `json:"fund_change,omitempty"`             // 必选  1  本次退款是否发生了资金变化 示例值:Y
	RefundFee            string              `json:"refund_fee,omitempty"`              // 必选  11 退款总金额。 指该笔交易累计已经退款成功的金额 示例值:88.88
	RefundDetailItemList []*TradeFundBill    `json:"refund_detail_item_list,omitempty"` // 可选   退款使用的资金渠道。 只有在签约中指定需要返回资金明细，或者入参的query_options中指定时才返回该字段信息。
	StoreName            string              `json:"store_name,omitempty"`              // 可选 512 交易在支付时候的门店名称
	BuyerUserId          string              `json:"buyer_user_id,omitempty"`           // 可选 28 买家在支付宝的用户id
	SendBackFee          string              `json:"send_back_fee,omitempty"`           // 可选 11 本次商户实际退回金额。 说明：如需获取该值，需在入参query_options中传入 refund_detail_item_list。
	RefundHybAmount      string              `json:"refund_hyb_amount,omitempty"`       // 可选 	11  本次请求退惠营宝金额 示例值:88.88
	RefundChargeInfoList []*RefundChargeInfo `json:"refund_charge_info_list,omitempty"` // 可选   退费信息
}

type TradeFundBill struct {
	FundChannel string `json:"fund_channel,omitempty"` // 必选	32 交易使用的资金渠道，详见 支付渠道列表
	Amount      string `json:"amount,omitempty"`       // 必选	32  该支付工具类型所使用的金额
	RealAmount  string `json:"real_amount,omitempty"`  // 可选	11 渠道实际付款金额
	FundType    string `json:"fund_type,omitempty"`    // 可选	32 渠道所使用的资金类型,目前只在资金渠道(fund_channel)是银行卡渠道(BANKCARD)的情况下才返回该信息(DEBIT_CARD:借记卡,CREDIT_CARD:信用卡,MIXED_CARD:借贷合一卡)
}

type RefundChargeInfo struct {
	RefundChargeFee        string          `json:"refund_charge_fee,omitempty"`          // 可选 11	实退费用
	SwitchFeeRate          string          `json:"switch_fee_rate,omitempty"`            // 可选 64	签约费率
	ChargeType             string          `json:"charge_type,omitempty"`                // 可选 64	收单手续费trade，花呗分期手续hbfq，其他手续费charge
	RefundSubFeeDetailList []*RefundSubFee `json:"refund_sub_fee_detail_list,omitempty"` // 可选   组合支付退费明细
}

type RefundSubFee struct {
	RefundChargeFee string `json:"refund_charge_fee,omitempty"` // 可选 11	实退费用
	SwitchFeeRate   string `json:"switch_fee_rate,omitempty"`   // 可选 64	签约费率
}

//////////////////////////////////////////////////////////////////////

// TradeFastPayRefundQueryReq 商户可使用该接口查询自已通过alipay.trade.refund提交的退款请求是否执行成功。
type TradeFastPayRefundQueryReq struct {
	OutTradeNo   string   `json:"out_trade_no,omitempty"`  // 与 TradeNo 二选一
	TradeNo      string   `json:"trade_no,omitempty"`      // 与 OutTradeNo 二选一
	OutRequestNo string   `json:"out_request_no"`          // 必须 64 退款请求号。 请求退款接口时，传入的退款请求号，如果在退款请求时未传入，则该值为创建交易时的商户订单号。
	QueryOptions []string `json:"query_options,omitempty"` // 可选 1024 查询选项，商户通过上送该参数来定制同步需要额外返回的信息字段，数组格式。枚举支持：
	//refund_detail_item_list：本次退款使用的资金渠道；
	//gmt_refund_pay：退款执行成功的时间；
	//deposit_back_info：银行卡冲退信息；
	BaseAliPayRequest
}

var _ IAliPayRequest = (*TradeFastPayRefundQueryReq)(nil)

func (r TradeFastPayRefundQueryReq) DoValidate() error {
	if len(r.OutTradeNo) == 0 && len(r.TradeNo) == 0 {
		return fmt.Errorf("商户订单号out_trade_no和支付宝交易号trade_no不能同时为空")
	}
	if len(r.OutRequestNo) == 0 {
		return fmt.Errorf("退款请求号out_request_no不能为空")
	}
	return nil
}

func (r TradeFastPayRefundQueryReq) RequestApi() string {
	return "alipay.trade.fastpay.refund.query"
}

type TradeFastPayRefundQueryRes struct {
	TradeFastPayRefundQueryResContent `json:"alipay_trade_fastpay_refund_query_response"`
	SignCertSn
}

func (r *TradeFastPayRefundQueryRes) String() string {
	buff, _ := json.Marshal(r)
	return string(buff)
}

type TradeFastPayRefundQueryResContent struct {
	CommonRes
	TradeNo      string `json:"trade_no"`       // 支付宝交易号
	OutTradeNo   string `json:"out_trade_no"`   // 创建交易传入的商户订单号
	OutRequestNo string `json:"out_request_no"` // 本笔退款对应的退款请求号
	RefundReason string `json:"refund_reason"`  // 发起退款时，传入的退款原因
	TotalAmount  string `json:"total_amount"`   // 发该笔退款所对应的交易的订单金额
	RefundAmount string `json:"refund_amount"`  // 本次退款请求，对应的退款金额
	RefundStatus string `json:"refund_status"`  // 退款状态。枚举值：
	//REFUND_SUCCESS 退款处理成功；
	//未返回该字段表示退款请求未收到或者退款失败；
	//注：如果退款查询发起时间早于退款时间，或者间隔退款发起时间太短，可能出现退款查询时还没处理成功，后面又处理成功的情况，建议商户在退款发起后间隔10秒以上再发起退款查询请求。
	RefundRoyaltys       []*RefundRoyalty    `json:"refund_royaltys"`                   // 可选 退分账明细信息
	GMTRefundPay         string              `json:"gmt_refund_pay"`                    // 可选 退款时间。
	RefundDetailItemList []*RefundDetailItem `json:"refund_detail_item_list,omitempty"` // 可选 本次退款使用的资金渠道；
	SendBackFee          string              `json:"send_back_fee"`                     // 可选 本次商户实际退回金额；
	DepositBackInfo      []*DepositBackInfo  `json:"deposit_back_info,omitempty"`       // 可选 银行卡冲退信息； 默认不返回该信息，需要在入参的query_options中指定"deposit_back_info"值时才返回该字段信息。
	RefundHybAmount      string              `json:"refund_hyb_amount,omitempty"`       // 可选 本次请求退惠营宝金额 示例值:88.88
	RefundChargeInfoList []*RefundChargeInfo `json:"refund_charge_info_list,omitempty"` // 可选  组合支付退费明细
}

type RefundDetailItem struct {
	FundChannel string `json:"fund_channel"`          // 必选 交易使用的资金渠道，详见 支付渠道列表
	Amount      string `json:"amount"`                // 必选 该支付工具类型所使用的金额
	RealAmount  string `json:"real_amount,omitempty"` // 可选 渠道实际付款金额
	FundType    string `json:"fund_type,omitempty"`   // 可选 渠道所使用的资金类型,目前只在资金渠道(fund_channel)是银行卡渠道(BANKCARD)的情况下才返回该信息(DEBIT_CARD:借记卡,CREDIT_CARD:信用卡,MIXED_CARD:借贷合一卡)
}

type RefundRoyalty struct {
	RefundAmount  string `json:"refund_amount"`             // 必选	9	 退分账金额
	RoyaltyType   string `json:"royalty_type,omitempty"`    // 可选	32 分账类型. 普通分账为：transfer; 补差为：replenish; 为空默认为分账transfer;
	ResultCode    string `json:"result_code"`               // 必选	32 退分账结果码 SUCCESS
	TransOut      string `json:"trans_out,omitempty"`       // 可选	28 转出人支付宝账号对应用户ID
	TransOutEmail string `json:"trans_out_email,omitempty"` // 可选	64	 转出人支付宝账号
	TransIn       string `json:"trans_in,omitempty"`        // 可选	28	 转入人支付宝账号对应用户ID
	TransInEmail  string `json:"trans_in_email,omitempty"`  // 可选  64    转入人支付宝账号
}

type DepositBackInfo struct {
	HasDepositBack     string `json:"has_deposit_back"`                // 可选 是否存在银行卡冲退信息
	DBackStatus        string `json:"dback_status"`                    // 可选 银行卡冲退状态。S-成功，F-失败，P-处理中。银行卡冲退失败，资金自动转入用户支付宝余额。
	DBackAmount        string `json:"dback_amount"`                    // 可选 银行卡冲退金额
	BankAckTime        string `json:"bank_ack_time"`                   // 可选 银行响应时间，格式为yyyy-MM-dd HH:mm:ss
	ESTBankReceiptTime string `json:"est_bank_receipt_time"`           // 可选 预估银行到账时间，格式为yyyy-MM-dd HH:mm:ss
	IsUseEnterprisePay bool   `json:"is_use_enterprise_pay,omitempty"` // 可选 是否包含因公付资产
}

// TradeOrderSettle 统一收单交易结算接口请求参数 https://docs.open.alipay.com/api_1/alipay.trade.order.settle
type TradeOrderSettle struct {
	AppAuthToken      string              `json:"-"`                  // 可选
	OutRequestNo      string              `json:"out_request_no"`     // 必须 结算请求流水号 开发者自行生成并保证唯一性
	TradeNo           string              `json:"trade_no"`           // 必须 支付宝订单号
	RoyaltyParameters []*RoyaltyParameter `json:"royalty_parameters"` // 必须 分账明细信息
	OperatorId        string              `json:"operator_id"`        //可选 操作员id
}

type RoyaltyParameter struct {
	TransOut         string  `json:"trans_out"`                   // 可选 分账支出方账户，类型为userId，本参数为要分账的支付宝账号对应的支付宝唯一用户号。以2088开头的纯16位数字。
	TransIn          string  `json:"trans_in"`                    // 可选 分账收入方账户，类型为userId，本参数为要分账的支付宝账号对应的支付宝唯一用户号。以2088开头的纯16位数字。
	Amount           float64 `json:"amount"`                      // 可选 分账的金额，单位为元
	AmountPercentage float64 `json:"amount_percentage,omitempty"` // 可选 分账信息中分账百分比。取值范围为大于0，少于或等于100的整数。
	Desc             string  `json:"desc"`                        // 可选 分账描述
}

//////////////////////////////////////////////////////////////////////

// DataServiceBillDownloadUrlQueryReq alipay.data.dataservice.bill.downloadurl.query(查询对账单下载地址) https://opendocs.alipay.com/open/028woc
type DataServiceBillDownloadUrlQueryReq struct {
	BillType string `json:"bill_type"` // 必选	20
	//账单类型，商户通过接口或商户经开放平台授权后其所属服务商通过接口可以获取以下账单类型，支持：
	//trade：商户基于支付宝交易收单的业务账单；
	//signcustomer：基于商户支付宝余额收入及支出等资金变动的账务账单；
	//merchant_act：营销活动账单，包含营销活动的发放，核销记录
	//trade_zft_merchant：直付通二级商户查询交易的业务账单；
	//zft_acc：直付通平台商查询二级商户流水使用，返回所有二级商户流水。
	BillDate string `json:"bill_date"` // 必选	15 账单时间：
	//* 日账单格式为yyyy-MM-dd，最早可下载2016年1月1日开始的日账单。不支持下载当日账单，只能下载前一日24点前的账单数据（T+1），当日数据一般于次日 9 点前生成，特殊情况可能延迟。
	//* 月账单格式为yyyy-MM，最早可下载2016年1月开始的月账单。不支持下载当月账单，只能下载上一月账单数据，当月账单一般在次月 3 日生成，特殊情况可能延迟。
	SMid string `json:"smid,omitempty"` // 可选	20 二级商户smid，这个参数只在bill_type是trade_zft_merchant时才能使用 2088123412341234
	BaseAliPayRequest
}

var _ IAliPayRequest = (*DataServiceBillDownloadUrlQueryReq)(nil)

func (r DataServiceBillDownloadUrlQueryReq) DoValidate() error {
	if len(r.BillType) == 0 {
		return fmt.Errorf("账单类型bill_type不能为空")
	}
	if len(r.BillDate) == 0 {
		return fmt.Errorf("账单时间bill_date不能为空")
	}
	return nil
}

func (r DataServiceBillDownloadUrlQueryReq) RequestApi() string {
	return "alipay.data.dataservice.bill.downloadurl.query"
}

type DataServiceBillDownloadUrlQueryRes struct {
	DataServiceBillDownloadUrlQueryResContent `json:"alipay_data_dataservice_bill_downloadurl_query_response"`
	SignCertSn
}

type DataServiceBillDownloadUrlQueryResContent struct {
	CommonRes
	BillDownloadUrl string `json:"bill_download_url"` // 必选 2048 账单下载地址链接，获取连接后30秒后未下载，链接地址失效。
}

//////////////////////////////////////////////////////////////////////

// TradeWapPayReq https://opendocs.alipay.com/open/02ivbs?ref=api&scene=21
type TradeWapPayReq struct {
	OutTradeNo      string         `json:"out_trade_no"`           // 必选	64 商户网站唯一订单号 70501111111S001111119
	TotalAmount     string         `json:"total_amount"`           // 必选	9 订单总金额。 单位为元，精确到小数点后两位，取值范围：[0.01,100000000] 。
	Subject         string         `json:"subject"`                // 必选	256 订单标题。 注意：不可使用特殊字符，如 /，=，& 等。
	ProductCode     string         `json:"product_code"`           // 必选	64 销售产品码。 QUICK_WAP_WAY 。
	AuthToken       string         `json:"auth_token,omitempty"`   // 可选	40 针对用户授权接口，获取用户相关数据时，用于标识用户授权关系 appopenBb64d181d0146481ab6a762c00714cC27
	QuitUrl         string         `json:"quit_url,omitempty"`     // 可选	400 用户付款中途退出返回商户网站的地址 http://www.taobao.com/product/113714.html
	GoodsDetail     []*GoodsDetail `json:"goods_detail,omitempty"` // 可选 	订单包含的商品列表信息，json格式，其它说明详见商品明细说明
	ExtendParams    *ExtendParams  `json:"extend_params,omitempty"`
	BusinessParams  string         `json:"business_params,omitempty"`   // 可选	512 商户传入业务信息，具体值要和支付宝约定，应用于安全，营销等参数直传场景，格式为json格式 {"data":"123"}
	PassbackParams  string         `json:"passback_params,omitempty"`   // 可选	512 公用回传参数，如果请求时传递了该参数，则返回给商户时会回传该参数。支付宝只会在同步返回（包括跳转回商户网站）和异步通知时将该参数原样返回。本参数必须进行UrlEncode之后才可以发送给支付宝。 merchantBizType%3d3C%26merchantBizNo%3d2016010101111
	MerchantOrderNo string         `json:"merchant_order_no,omitempty"` // 可选	32 商户原始订单号，最大长度限制32位 	20161008001
	ExtUserInfo     *ExtUserInfo   `json:"ext_user_info,omitempty"`
	ReturnUrl       string         `json:"-" url:"-"` // 可选	256 HTTP/HTTPS开头字符串
	NotifyUrl       string         `json:"-" url:"-"` // 可选	256	 支付宝服务器主动通知商户服务器里指定的页面http/https路径。
	BaseAliPayRequest
}

func NewTradeWapPayReq(outTradeNo, totalAmount, subject string, opts ...TradeWapPayOpt) *TradeWapPayReq {
	req := &TradeWapPayReq{
		OutTradeNo:  outTradeNo,
		TotalAmount: totalAmount,
		Subject:     subject,
		ProductCode: "QUICK_WAP_WAY",
	}
	for _, opt := range opts {
		opt(req)
	}
	return req
}

type TradeWapPayOpt func(req *TradeWapPayReq)

func WithGoodsDetailWapPayOpt(goodsDetail []*GoodsDetail) TradeWapPayOpt {
	return func(r *TradeWapPayReq) {
		r.GoodsDetail = goodsDetail
	}
}

func WithNotifyUrlWapPayOpt(notifyUrl string) TradeWapPayOpt {
	return func(param *TradeWapPayReq) {
		if len(notifyUrl) == 0 {
			return
		}
		param.NotifyUrl = notifyUrl
	}
}

func WithReturnUrlWapPayOpt(returnUrl string) TradeWapPayOpt {
	return func(param *TradeWapPayReq) {
		if len(returnUrl) == 0 {
			return
		}
		param.ReturnUrl = returnUrl
	}
}

type ExtendParams struct {
	SysServiceProviderId string `json:"sys_service_provider_id,omitempty"` // 可选	64 系统商编号 该参数作为系统商返佣数据提取的依据，请填写系统商签约协议的PID
	HbFqNum              string `json:"hb_fq_num,omitempty"`               // 可选	5 使用花呗分期要进行的分期数
	HbFqSellerPercent    string `json:"hb_fq_seller_percent,omitempty"`    // 可选	3 使用花呗分期需要卖家承担的手续费比例的百分值，传入100代表100%
	IndustryRefluxInfo   string `json:"industry_reflux_info,omitempty"`    // 可选	512 行业数据回流信息, 详见：地铁支付接口参数补充说明 {\"scene_code\":\"metro_tradeorder\",\"channel\":\"xxxx\",\"scene_data\":{\"asset_name\":\"ALIPAY\"}}
	CardType             string `json:"card_type,omitempty"`               // 可选	32 卡类型 S0JP0000
	SpecifiedSellerName  string `json:"specified_seller_name,omitempty"`   // 可选	32 特殊场景下，允许商户指定交易展示的卖家名称 XXX的跨境小铺
}
type ExtUserInfo struct {
	Name     string `json:"name,omitempty"`      // 可选	16 指定买家姓名。 注： need_check_info=T或fix_buyer=T时该参数才有效
	Mobile   string `json:"mobile,omitempty"`    // 可选	20 指定买家手机号。 注：该参数暂不校验 16587658765
	CertType string `json:"cert_type,omitempty"` //可选	32 指定买家证件类型。 枚举值：
	//IDENTITY_CARD：身份证；
	//PASSPORT：护照；
	//OFFICER_CARD：军官证；
	//SOLDIER_CARD：士兵证；
	//HOKOU：户口本。如有其它类型需要支持，请与蚂蚁金服工作人员联系。
	//注： need_check_info=T或fix_buyer=T时该参数才有效，支付宝会比较买家在支付宝留存的证件类型与该参数传入的值是否匹配。

	CertNo string `json:"cert_no,omitempty"` // 可选	64 买家证件号。 注：need_check_info=T或fix_buyer=T时该参数才有效，支付宝会比较买家在支付宝留存的证件号码与该参数传入的值是否匹配。
	MinAge string `json:"min_age,omitempty"` // 可选	3 允许的最小买家年龄。 买家年龄必须大于等于所传数值 注：
	//1. need_check_info=T时该参数才有效
	//2. min_age为整数，必须大于等于0

	FixBuyer string `json:"fix_buyer,omitempty"` // 可选	8
	//是否强制校验买家身份。
	//需要强制校验传：T;
	//不需要强制校验传：F或者不传；
	//当传T时，接口上必须指定cert_type、cert_no和name信息且支付宝会校验传入的信息跟支付买家的信息都匹配，否则报错。
	//默认为不校验。

	NeedCheckInfo string `json:"need_check_info,omitempty"` // 可选	1
	//是否强制校验买家信息；
	//需要强制校验传：T;
	//不需要强制校验传：F或者不传；
	//当传T时，支付宝会校验支付买家的信息与接口上传递的cert_type、cert_no、name或age是否匹配，只有接口传递了信息才会进行对应项的校验；只要有任何一项信息校验不匹配交易都会失败。如果传递了need_check_info，但是没有传任何校验项，则不进行任何校验。
	//默认为不校验。
	IdentityHash string `json:"identity_hash,omitempty"` // 可选	128 买家加密身份信息。当指定了此参数且指定need_check_info=T时，支付宝会对买家身份进行校验，校验逻辑为买家姓名、买家证件号拼接后的字符串，以sha256算法utf-8编码计算hash，若与传入的值不匹配则会拦截本次支付。注意：如果同时指定了用户明文身份信息（name，cert_type，cert_no中任意一个），则忽略identity_hash以明文参数校验。 27bfcd1dee4f22c8fe8a2374af9b660419d1361b1c207e9b41a754a113f38fcc
}

var _ IAliPayRequest = (*TradeWapPayReq)(nil)

func (r TradeWapPayReq) DoValidate() error {
	outTradeNoLen := len(r.OutTradeNo)
	if outTradeNoLen == 0 || outTradeNoLen > 64 {
		return fmt.Errorf("当前商户订单号长度为%d，不符合长度1-64要求", outTradeNoLen)
	}
	subjectLen := len(r.Subject)
	if subjectLen == 0 || subjectLen > 256 {
		return fmt.Errorf("当前订单标题长度为%d，不符合长度1-256要求", subjectLen)
	}
	if r.ProductCode != QuickWapWay {
		return fmt.Errorf("销售产品码必须是%s", QuickWapWay)
	}
	if len(r.GoodsDetail) > 0 {
		for _, detail := range r.GoodsDetail {
			if err := detail.DoValidate(); err != nil {
				return err
			}
		}
	}
	return nil
}

//func (r TradeWapPayReq) RequestHttpMethod() string {
//	return ""
//}

func (r TradeWapPayReq) RequestApi() string {
	return "alipay.trade.wap.pay"
}

func (r TradeWapPayReq) RequestApiVersion() string {
	return "2.0"
}

/////////////////////////////////////////////////

type TradeAppPayReq struct {
	OutTradeNo      string         `json:"out_trade_no"`                // 必选	64 商户订单号。 由商家自定义，64个字符以内，仅支持字母、数字、下划线且需保证在商户端不重复。
	TotalAmount     string         `json:"total_amount"`                // 必选	11 订单总金额，单位为元，精确到小数点后两位，取值范围为 [0.01,100000000]。金额不能为0。
	Subject         string         `json:"subject"`                     // 必选	256 订单标题。注意：不可使用特殊字符，如 /，=，& 等。
	GoodsDetail     []*GoodsDetail `json:"goods_detail,omitempty"`      // 可选 订单包含的商品列表信息，json格式。
	TimeExpire      string         `json:"time_expire,omitempty"`       // 可选 订单绝对超时时间。 格式为yyyy-MM-dd HH:mm:ss。超时时间范围：1m~15d。 注：time_expire和timeout_express两者只需传入一个或者都不传，两者均传入时，优先使用time_expire。
	ExtendParams    *ExtendParams  `json:"extend_params,omitempty"`     // 可选 业务扩展参数
	BusinessParams  string         `json:"business_params,omitempty"`   // 可选	512 商户传入业务信息，具体值要和支付宝约定，应用于安全，营销等参数直传场景，格式为json格式 {"data":"123"}
	PassbackParams  string         `json:"passback_params,omitempty"`   // 可选 512 回传参数，公共回传参数，如果请求时传递了该参数，则返回的异步通知会原样传回。本参数必须进行 UrlEncode 之后才可传入。 merchantBizType%3d3C%26merchantBizNo%3d201601001111
	MerchantOrderNo string         `json:"merchant_order_no,omitempty"` // 可选	32 商户原始订单号，最大长度限制 32 位
	ExtUserInfo     *ExtUserInfo   `json:"ext_user_info,omitempty"`     // 可选    外部指定买家
	QueryOptions    []string       `json:"query_options,omitempty"`     // 可选 1024 返回参数选项。 商户通过传递该参数来定制同步需要额外返回的信息字段，数组格式。包括但不限于：["hyb_amount","enterprise_pay_info"]
	ReturnUrl       string         `json:"-" url:"-"`                   // 可选	256 HTTP/HTTPS开头字符串
	NotifyUrl       string         `json:"-" url:"-"`                   // 可选	256	 支付宝服务器主动通知商户服务器里指定的页面http/https路径。
	BaseAliPayRequest
}

func (r *TradeAppPayReq) DoValidate() error {
	outTradeNoLen := len(r.OutTradeNo)
	if outTradeNoLen == 0 || outTradeNoLen > 64 {
		return fmt.Errorf("当前商户订单号长度为%d，不符合长度1-64要求", outTradeNoLen)
	}
	subjectLen := len(r.Subject)
	if subjectLen == 0 || subjectLen > 256 {
		return fmt.Errorf("当前订单标题长度为%d，不符合长度1-256要求", subjectLen)
	}
	if len(r.GoodsDetail) > 0 {
		for _, detail := range r.GoodsDetail {
			if err := detail.DoValidate(); err != nil {
				return err
			}
		}
	}
	return nil
}

type TradeAppPayRes struct {
	CommonRes
	OutTradeNo      string `json:"out_trade_no"`      // 必选	64	商户网站唯一订单号 70501111111S001111119
	TradeNo         string `json:"trade_no"`          // 必选	64	 该交易在支付宝系统中 2014112400001000340011111118
	TotalAmount     string `json:"total_amount"`      // 必选	9 该笔订单的资金总额，单位为人民币（元），取值范围为 0.01~100000000.00，精确到小数点后两位。 9.00
	SellerId        string `json:"seller_id"`         // 必选	16 收款支付宝账号对应的支付宝唯一用户号。 以2088开头的纯16位数字 2088111111116894
	MerchantOrderNo string `json:"merchant_order_no"` // 必选	32	商户原始订单号，最大长度限制32位 20161008001
}

func (r *TradeAppPayReq) RequestApi() string {
	return "alipay.trade.app.pay"
}

func (r *TradeAppPayReq) RequestApiVersion() string {
	return "2.0"
}

/////////////////////////////////////////////////

// TradePreCreateReq alipay.trade.precreate(统一收单线下交易预创建) https://opendocs.alipay.com/open/02ekfg?scene=19
type TradePreCreateReq struct {
	OutTradeNo         string         `json:"out_trade_no"`                  // 必选	64 商户网站唯一订单号 70501111111S001111119
	TotalAmount        string         `json:"total_amount"`                  // 必选	9 订单总金额。 单位为元，精确到小数点后两位，取值范围：[0.01,100000000] 。
	Subject            string         `json:"subject"`                       // 必选	256 订单标题。 注意：不可使用特殊字符，如 /，=，& 等。
	ProductCode        string         `json:"product_code"`                  // 必选	64 销售产品码。 销售产品码。如果签约的是当面付快捷版，则传 OFFLINE_PAYMENT；其它支付宝当面付产品传 FACE_TO_FACE_PAYMENT；不传则默认使用 FACE_TO_FACE_PAYMENT。
	SellerId           string         `json:"seller_id,omitempty"`           // 可选  30 卖家支付宝用户 ID。 如果该值为空，则默认为商户签约账号对应的支付宝用户 ID。不允许收款账号与付款方账号相同
	Body               string         `json:"body,omitempty"`                // 可选 128 订单附加信息。 如果请求时传递了该参数，将在异步通知、对账单中原样返回，同时会在商户和用户的pc账单详情中作为交易描述展示
	GoodsDetail        []*GoodsDetail `json:"goods_detail,omitempty"`        // 可选 订单包含的商品列表信息，json格式。
	ExtendParams       *ExtendParams  `json:"extend_params,omitempty"`       // 可选 业务扩展参数
	BusinessParams     string         `json:"business_params,omitempty"`     // 可选	512 商户传入业务信息，具体值要和支付宝约定，应用于安全，营销等参数直传场景，格式为json格式 {"data":"123"}
	DiscountableAmount string         `json:"discountable_amount,omitempty"` // 可选	11 可打折金额。参与优惠计算的金额，单位为元，精确到小数点后两位，取值范围为 [0.01,100000000]。如果该值未传入，但传入了【订单总金额】和【不可打折金额】，则该值默认为【订单总金额】-【不可打折金额】 80.00
	StoreId            string         `json:"store_id,omitempty"`            // 可选	32 商户门店编号。指商户创建门店时输入的门店编号。
	OperatorId         string         `json:"operator_id,omitempty"`         // 可选 28 操作员id
	TerminalId         string         `json:"terminal_id,omitempty"`         // 可选	32 商户机具终端编号
	MerchantOrderNo    string         `json:"merchant_order_no,omitempty"`   // 可选	32 商户原始订单号，最大长度限制 32 位
	BaseAliPayRequest
}

func (r *TradePreCreateReq) DoValidate() error {
	outTradeNoLen := len(r.OutTradeNo)
	if outTradeNoLen == 0 || outTradeNoLen > 64 {
		return fmt.Errorf("当前商户订单号长度为%d，不符合长度1-64要求", outTradeNoLen)
	}
	subjectLen := len(r.Subject)
	if subjectLen == 0 || subjectLen > 256 {
		return fmt.Errorf("当前订单标题长度为%d，不符合长度1-256要求", subjectLen)
	}
	if len(r.GoodsDetail) > 0 {
		for _, detail := range r.GoodsDetail {
			if err := detail.DoValidate(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *TradePreCreateReq) RequestApi() string {
	return "alipay.trade.precreate"
}

type TradePreCreateRes struct {
	TradePreCreateResContent `json:"alipay_trade_precreate_response"`
	SignCertSn
}

type TradePreCreateResContent struct {
	CommonRes
	OutTradeNo string `json:"out_trade_no"` // 必选	64 商户的订单号 6823789339978248
	QrCode     string `json:"qr_code"`      // 必选	1024 当前预下单请求生成的二维码码串，有效时间2小时，可以用二维码生成工具根据该码串值生成对应的二维码 https://qr.alipay.com/bavh4wjlxf12tper3a
}

/////////////////////////////////////////////////

var _ IAliPayRequest = &TradeCancelReq{}

type TradeCancelReq struct {
	OutTradeNo string `json:"out_trade_no,omitempty"` // 特殊可选	64 商户订单号。 订单支付时传入的商户订单号，商家自定义且保证商家系统中唯一。与支付宝交易号 trade_no 不能同时为空。
	TradeNo    string `json:"trade_no,omitempty"`     // 特殊可选	64 支付宝交易号。 和商户订单号 out_trade_no 不能同时为空。
	BaseAliPayRequest
}

func (r *TradeCancelReq) DoValidate() error {
	outTradeNoLen := len(r.OutTradeNo)
	tradeNo := len(r.TradeNo)
	if outTradeNoLen == 0 && tradeNo == 0 {
		return fmt.Errorf("商户订单号out_trade_no和支付宝交易号trade_no不能同时为空")
	}
	if outTradeNoLen > 64 {
		return fmt.Errorf("当前商户订单号长度为%d，不符合长度1-64要求", outTradeNoLen)
	}
	if tradeNo > 64 {
		return fmt.Errorf("当前订单标题长度为%d，不符合长度1-256要求", tradeNo)
	}
	return nil
}
func (r *TradeCancelReq) RequestApi() string {
	return "alipay.trade.cancel"
}

type TradeCancelRes struct {
	TradeCancelResContent `json:"alipay_trade_cancel_response"`
	SignCertSn
}

func (r *TradeCancelRes) String() string {
	buff, _ := json.Marshal(r)
	return string(buff)
}

type TradeCancelResContent struct {
	CommonRes
	TradeNo    string `json:"trade_no"`     // 必选	64 支付宝交易号; 当发生交易关闭或交易退款时返回； 2013112011001004330000121536
	OutTradeNo string `json:"out_trade_no"` // 必选	64 商户订单号 6823789339978248
	RetryFlag  string `json:"retry_flag"`   // 必选	1 是否需要重试 N
	Action     string `json:"action"`       // 必选	10 本次撤销触发的交易动作,接口调用成功且交易存在时返回。可能的返回值： close：交易未支付，触发关闭交易动作，无退款； refund：交易已支付,触发交易退款动作； 未返回：未查询到交易，或接口调用失败；
}

/////////////////////////////////////

type TradeCreateReq struct {
	OutTradeNo     string         `json:"out_trade_no"`              // 	必选	64 商户订单号。 由商家自定义，64个字符以内，仅支持字母、数字、下划线且需保证在商户端不重复。
	TotalAmount    string         `json:"total_amount"`              // 	必选	9 订单总金额，单位为元，精确到小数点后两位，取值范围为 [0.01,100000000]。金额不能为0。
	Subject        string         `json:"subject"`                   // 	必选	256 订单标题。注意：不可使用特殊字符，如 /，=，& 等。
	ProductCode    string         `json:"product_code"`              // 	可选	64 产品码。 商家和支付宝签约的产品码。 枚举值（点击查看签约情况）： FACE_TO_FACE_PAYMENT：当面付产品； 默认值为FACE_TO_FACE_PAYMENT。
	BuyerId        string         `json:"buyer_id"`                  // 必选  28 买家支付宝用户ID。 2088开头的16位纯数字，小程序场景下获取用户ID请参考：用户授权; 其它场景下获取用户ID请参考：网页授权获取用户信息; 注：交易的买家与卖家不能相同。
	SellerId       string         `json:"seller_id,omitempty"`       // 可选  28 	 卖家支付宝用户ID。 当需要指定收款账号时，通过该参数传入，如果该值为空，则默认为商户签约账号对应的支付宝用户ID。 收款账号优先级规则：门店绑定的收款账户>请求传入的seller_id>商户签约账号对应的支付宝用户ID； 注：直付通和机构间联场景下seller_id无需传入或者保持跟pid一致； 如果传入的seller_id与pid不一致，需要联系支付宝小二配置收款关系；
	Body           string         `json:"body,omitempty"`            // 可选 	128 订单附加信息。 如果请求时传递了该参数，将在异步通知、对账单中原样返回，同时会在商户和用户的pc账单详情中作为交易描述展示
	GoodsDetail    []*GoodsDetail `json:"goods_detail,omitempty"`    // 可选 订单包含的商品列表信息，json格式。
	TimeExpire     string         `json:"time_expire,omitempty"`     // 可选 32 订单绝对超时时间。 格式为yyyy-MM-dd HH:mm:ss。超时时间范围：1m~15d。 注：time_expire和timeout_express两者只需传入一个或者都不传，两者均传入时，优先使用time_expire。
	TimeoutExpress string         `json:"timeout_express,omitempty"` // 可选 6 订单相对超时时间。从交易创建时间开始计算。
	//该笔订单允许的最晚付款时间，逾期将关闭交易。取值范围：1m～15d。m-分钟，h-小时，d-天，1c-当天（1c-当天的情况下，无论交易何时创建，都在0点关闭）。 该参数数值不接受小数点， 如 1.5h，可转换为 90m。
	//当面付场景默认值为3h。 注：time_expire和timeout_express两者只需传入一个或者都不传，如果两者都传，优先使用time_expire。
	SettleInfo           *SettleInfo          `json:"settle_info,omitempty"`           // 可选  描述结算信息，json格式。
	ExtendParams         *ExtendParams        `json:"extend_params,omitempty"`         // 可选 业务扩展参数，具体传参数见官方接口文档
	BusinessParams       string               `json:"business_params,omitempty"`       // 可选	512 商户传入业务信息，具体值要和支付宝约定，应用于安全，营销等参数直传场景，格式为json格式 {"data":"123"}
	DiscountableAmount   string               `json:"discountable_amount,omitempty"`   // 可选	9 参与优惠计算的金额，单位为元，精确到小数点后两位，取值范围[0.01,100000000]。 如果同时传入了【可打折金额】、【不可打折金额】和【订单总金额】，则必须满足如下条件：【订单总金额】=【可打折金额】+【不可打折金额】。 如果订单金额全部参与优惠计算，则【可打折金额】和【不可打折金额】都无需传入。
	UndiscountableAmount string               `json:"undiscountable_amount,omitempty"` // 可选	9 不可打折金额。 不参与优惠计算的金额，单位为元，精确到小数点后两位，取值范围[0.01,100000000]。 如果同时传入了【可打折金额】、【不可打折金额】和【订单总金额】，则必须满足如下条件：【订单总金额】=【可打折金额】+【不可打折金额】。 如果订单金额全部参与优惠计算，则【可打折金额】和【不可打折金额】都无需传入
	StoreId              string               `json:"store_id,omitempty"`              // 可选	32 商户门店编号。指商户创建门店时输入的门店编号。
	OperatorId           string               `json:"operator_id,omitempty"`           // 可选 28 操作员id
	TerminalId           string               `json:"terminal_id,omitempty"`           // 可选	32 商户机具终端编号
	LogisticsDetail      *LogisticsDetail     `json:"logistics_detail,omitempty"`      // 可选 物流信息
	ReceiverAddressInfo  *ReceiverAddressInfo `json:"receiver_address_info,omitempty"` // 可选 收货人及地址信息
	QueryOptions         []string             `json:"query_options,omitempty"`         // 可选 1024 返回参数选项。 商户通过传递该参数来定制需要额外返回的信息字段，数组格式。包括但不限于：["enterprise_pay_info","hyb_amount"]
	BkAgentReqInfo       *BkAgentReqInfo      `json:"bkagent_req_info,omitempty"`      // 可选 间联交易下，由收单机构上送的信息
	// 自己添加
	NotifyUrl string `json:"-" url:"-"` // 可选	256	 支付宝服务器主动通知商户服务器里指定的页面http/https路径。
	BaseAliPayRequest
}

func (r *TradeCreateReq) DoValidate() error {
	outTradeNoLen := len(r.OutTradeNo)
	if outTradeNoLen == 0 {
		return fmt.Errorf("商户订单号out_trade_no不能为空")
	}
	if outTradeNoLen > 64 {
		return fmt.Errorf("当前商户订单号长度为%d，不符合长度1-64要求", outTradeNoLen)
	}

	return nil
}
func (r *TradeCreateReq) RequestApi() string {
	return "alipay.trade.create"
}

type BkAgentReqInfo struct {
	AcqCode    string // 可选	14 收单机构在清算组织登记或分配的机构代码
	MerchCode  string // 	可选	32 收单机构在清算组织登记的商户编码 123412341234
	DeviceType string // 可选	2 终端设备类型，受理方可参考终端 注册时的设备类型填写，取值如下:
	//01:自动柜员机(含 ATM 和 CDM) 和多媒体自助终端
	//02:传统 POS
	//03:mPOS
	//04:智能 POS
	//05:II 型固定电话
	//06:云闪付终端;
	//07:保留使用;
	//08:手机 POS;
	//09:刷脸付终端;
	//10:条码支付受理终端;
	//11:条码支付辅助受理终端;
	//12:行业终端(公交、地铁用于指
	//定行业的终端);
	//13:MIS 终端;
	SerialNum string // 可选	50 终端类型填写为 02、 03、04、05、06、08、09 或 10 时，必 须填写终端序列号。 123123123123
	Location  string // 可选	32 终端设备实时经纬度信息，格式为纬度/经度，+表示北纬、东经，-表示南纬、西经。 +37.28/-121.268
}
type ReceiverAddressInfo struct {
	Name         string `json:"name,omitempty"`          // 可选	512 收货人的姓名 张三
	Address      string `json:"address,omitempty"`       // 可选	512 收货地址 上海市浦东新区陆家嘴银城中路501号
	Mobile       string `json:"mobile,omitempty"`        // 可选	60 收货人手机号 13120180615
	Zip          string `json:"zip,omitempty"`           // 可选	40 收货地址邮编 200120
	DivisionCode string `json:"division_code,omitempty"` // 可选	16 中国标准城市区域码

}

type LogisticsDetail struct {
	LogisticsType string `json:"logistics_type,omitempty"` // 可选	32 物流类型, POST 平邮, EXPRESS 其他快递, VIRTUAL 虚拟物品, EMS EMS, DIRECT 无需物流。
}
type SettleInfo struct {
	SettleDetailInfoList []SettleDetailInfo `json:"settle_detail_infos"`          // 必选	10	 结算详细信息，json数组，目前只支持一条。
	SettlePeriodTime     string             `json:"settle_period_time,omitempty"` // 可选	10 该笔订单的超期自动确认结算时间，到达期限后，将自动确认结算。此字段只在签约账期结算模式时有效。取值范围：1d～365d。d-天。 该参数数值不接受小数点。
}
type SettleDetailInfo struct {
	TransInType string `json:"trans_in_type"` // 必选	64 结算收款方的账户类型。
	//cardAliasNo：结算收款方的银行卡编号;
	//userId：表示是支付宝账号对应的支付宝唯一用户号;
	//loginName：表示是支付宝登录号；
	//defaultSettle：表示结算到商户进件时设置的默认结算账号，结算主体为门店时不支持传defaultSettle；
	TransIn          string `json:"trans_in"`                     // 必选	64 结算收款方。当结算收款方类型是cardAliasNo时，本参数为用户在支付宝绑定的卡编号；结算收款方类型是userId时，本参数为用户的支付宝账号对应的支付宝唯一用户号，以2088开头的纯16位数字；当结算收款方类型是loginName时，本参数为用户的支付宝登录号；当结算收款方类型是defaultSettle时，本参数不能传值，保持为空。
	SummaryDimension string `json:"summary_dimension,omitempty"`  // 可选	64 结算汇总维度，按照这个维度汇总成批次结算，由商户指定。 目前需要和结算收款方账户类型为cardAliasNo配合使用
	SettleEntityId   string `json:"settle_entity_id,omitempty"`   // 可选	64 结算主体标识。当结算主体类型为SecondMerchant时，为二级商户的SecondMerchantID；当结算主体类型为Store时，为门店的外标。
	SettleEntityType string `json:"settle_entity_type,omitempty"` // 可选	32 结算主体类型。 二级商户:SecondMerchant;商户或者直连商户门店:Store SecondMerchant、Store
	Amount           string `json:"amount"`                       // 必选	9 结算的金额，单位为元。在创建订单和支付接口时必须和交易金额相同。在结算确认接口时必须等于交易金额减去已退款金额。直付通账期模式下，如使用部分结算能力、传递了actual_amount字段，则忽略本字段的校验、可不传。
	ActualAmount     string `json:"actual_amount,omitempty"`      // 可选	9 仅在直付通账期模式下，当一笔交易需要分多次发起部分确认结算时使用，表示本次确认结算的实际结算金额。传递本字段后，原amount字段不再生效，结算金额以本字段为准。如已经发生过部分确认结算、不传递本字段则默认按剩余待结算金额一次性结算。
}

type TradeCreateRes struct {
	TradeCreateResContent `json:"alipay_trade_create_response"`
	SignCertSn
}

func (r *TradeCreateRes) String() string {
	buff, _ := json.Marshal(r)
	return string(buff)
}

type TradeCreateResContent struct {
	CommonRes
	OutTradeNo string `json:"out_trade_no"` // 必选	64 商户订单号 20150423001001
	TradeNo    string `json:"trade_no"`     // 必选	64 支付宝交易号 2015042321001004720200028594
}

///////////////////////////////////////////////

type TradePayReq struct {
	OutTradeNo     string         `json:"out_trade_no"`              // 必选	64 商户订单号。 由商家自定义，64个字符以内，仅支持字母、数字、下划线且需保证在商户端不重复。
	TotalAmount    string         `json:"total_amount"`              // 必选	9 订单总金额，单位为元，精确到小数点后两位，取值范围为 [0.01,100000000]。金额不能为0。
	Subject        string         `json:"subject"`                   // 必选	256 订单标题。注意：不可使用特殊字符，如 /，=，& 等。
	AuthCode       string         `json:"auth_code"`                 // 必选	64 支付授权码。 当面付场景传买家的付款码（25~30开头的长度为16~24位的数字，实际字符串长度以开发者获取的付款码长度为准）或者刷脸标识串（fp开头的35位字符串）。
	Scene          string         `json:"scene"`                     // 必选	32 支付场景。 枚举值： bar_code：当面付条码支付场景； security_code：当面付刷脸支付场景，对应的auth_code为fp开头的刷脸标识串； 默认值为bar_code。
	ProductCode    string         `json:"product_code,omitempty"`    // 可选	64 产品码。 商家和支付宝签约的产品码。 当面付场景下，如果签约的是当面付快捷版，则传 OFFLINE_PAYMENT; 其它支付宝当面付产品传 FACE_TO_FACE_PAYMENT； 不传则默认使用FACE_TO_FACE_PAYMENT。
	SellerId       string         `json:"seller_id,omitempty"`       // 可选	28 卖家支付宝用户ID。 当需要指定收款账号时，通过该参数传入，如果该值为空，则默认为商户签约账号对应的支付宝用户ID。 收款账号优先级规则：门店绑定的收款账户>请求传入的seller_id>商户签约账号对应的支付宝用户ID； 注：直付通和机构间联场景下seller_id无需传入或者保持跟pid一致；如果传入的seller_id与pid不一致，需要联系支付宝小二配置收款关系；
	GoodsDetail    []*GoodsDetail `json:"goods_detail,omitempty"`    // 可选 订单包含的商品列表信息，json格式。
	ExtendParams   *ExtendParams  `json:"extend_params,omitempty"`   // 可选 业务扩展参数
	BusinessParams string         `json:"business_params,omitempty"` // 可选	 商户传入业务信息，具体值要和支付宝约定，应用于安全，营销等参数直传场景，格式为json格式
	PromoParam     string         `json:"promo_params,omitempty"`    // 可选	 优惠明细参数，通过此属性补充营销参数。 注：仅与支付宝协商后可用。
	StoreId        string         `json:"store_id,omitempty"`        // 可选	32 商户门店编号。指商户创建门店时输入的门店编号。
	OperatorId     string         `json:"operator_id,omitempty"`     // 可选 28 操作员id
	TerminalId     string         `json:"terminal_id,omitempty"`     // 可选	32 商户机具终端编号
	QueryOptions   []string       `json:"query_options,omitempty"`   // 可选 1024 返回参数选项。 商户通过传递该参数来定制需要额外返回的信息字段，数组格式。包括但不限于：["enterprise_pay_info","hyb_amount"]
	// 自己添加
	NotifyUrl string `json:"-" url:"-"` // 可选	256	 支付宝服务器主动通知商户服务器里指定的页面http/https路径。
	BaseAliPayRequest
}

func (r *TradePayReq) DoValidate() error {
	outTradeNoLen := len(r.OutTradeNo)
	if outTradeNoLen == 0 {
		return fmt.Errorf("商户订单号out_trade_no不能为空")
	}
	if outTradeNoLen > 64 {
		return fmt.Errorf("当前商户订单号长度为%d，不符合长度1-64要求", outTradeNoLen)
	}
	authCodeLen := len(r.AuthCode)
	if authCodeLen < 16 || authCodeLen > 24 {
		return fmt.Errorf("支付授权码auth_code长度为%d，不符合长度64个字符的要求", outTradeNoLen)
	}
	subjectLen := len(r.Subject)
	if subjectLen == 0 || subjectLen > 256 {
		return fmt.Errorf("订单标题subject长度为%d，不符合长度1-256个字符的要求", outTradeNoLen)
	}

	return nil
}
func (r *TradePayReq) RequestApi() string {
	return "alipay.trade.pay"
}

type TradePayRes struct {
	TradePayResContent `json:"alipay_trade_pay_response"`
	SignCertSn
}
type TradePayResContent struct {
	CommonRes
	TradeNo             string           `json:"trade_no"`                      // 必须 64 支付宝交易号
	OutTradeNo          string           `json:"out_trade_no"`                  // 必选	64	商户订单号
	BuyerLogonId        string           `json:"buyer_logon_id,omitempty"`      // 必选	100 用户的登录id
	TotalAmount         string           `json:"total_amount,omitempty"`        // 可选 11 订单金额。本次交易支付订单金额，单位为人民币（元），精确到小数点后 2 位 20.00
	ReceiptAmount       string           `json:"receipt_amount,omitempty"`      // 可选 11 实收金额。商家在交易中实际收到的款项，单位为人民币（元），精确到小数点后 2 位 15.00
	BuyerPayAmount      string           `json:"buyer_pay_amount,omitempty"`    // 可选 11 用户在交易中支付的金额，单位为人民币（元），精确到小数点后 2 位 12.00
	PointAmount         string           `json:"point_amount,omitempty"`        // 可选 11 积分支付的金额，单位为元，两位小数。该金额代表该笔交易中用户使用积分支付的金额，比如集分宝或者支付宝实时优惠等
	InvoiceAmount       string           `json:"invoice_amount,omitempty"`      // 可选 11 开票金额。用户在交易中支付的可开发票的金额，单位为人民币（元），精确到小数点后 2 位 13.88
	GmtPayment          string           `json:"gmt_payment,omitempty"`         // 可选 交易付款时间。格式为 yyyy-MM-dd HH:mm:ss 2018-08-25 15:34:42
	FundBillList        []*FundBill      `json:"fund_bill_list"`                // 必选  交易支付使用的资金渠道。只有在签约中指定需要返回资金明细，或者入参的query_options中指定时才返回该字段信息。
	StoreName           string           `json:"store_name,omitempty"`          // 可选	512 发生支付交易的商户门店名称 证大五道口店
	DiscountGoodsDetail string           `json:"discount_goods_detail"`         // 可选 5120 本次交易支付所使用的单品券优惠的商品优惠信息
	BuyerUserId         string           `json:"buyer_user_id,omitempty"`       // 可选 28 买家在支付宝的用户id
	VoucherDetailList   []*VoucherDetail `json:"voucher_detail_list,omitempty"` // 可选 本交易支付时使用的所有优惠券信息
	MdiscountAmount     string           `json:"mdiscount_amount,omitempty"`    // 特殊可选 11 商家优惠金额
	DiscountAmount      string           `json:"discount_amount,omitempty"`     // 特殊可选 11 平台优惠金额
}

/////////////////////////////////////////////////////////////

var _ IAliPayRequest = &TradeOrderInfoSyncReq{}

type TradeOrderInfoSyncReq struct {
	TradeNo       string `json:"trade_no"`                  // 必选	64 支付宝交易号 2018061021001004680073956707
	OrigRequestNo string `json:"orig_request_no,omitempty"` // 可选	64 原始业务请求单号。如对某一次退款进行履约时，该字段传退款时的退款请求号 HZ01RF001
	OutRequestNo  string `json:"out_request_no"`            // 必选	64 外部请求号，商家自定义。标识一笔交易多次请求，同一笔交易多次信息同步时需要保证唯一。 HZ01RF001
	BizType       string `json:"biz_type"`                  // 必选	64 交易信息同步对应的业务类型，具体值与支付宝约定； 信用授权场景下传CREDIT_AUTH 信用代扣场景下传CREDIT_DEDUCT
	OrderBizInfo  string `json:"order_biz_info,omitempty"`  // 可选	2018
	/*商户传入同步信息，具体值要和支付宝约定；用于芝麻信用租车、单次授权等信息同步场景，格式为json格式。 状态枚举如下：
	COMPLETE：同步用户已履约
	适用场景：发起扣款后，芝麻生成待履约记录，如果用户通过其他方式完成订单支付，请反馈该状态，芝麻将完结待履约记录对用户形成一条良好履约记录；
	CLOSED： 同步履约已取消
	适用场景：发起扣款后，芝麻生成待履约记录，如果发现该笔扣款无效需要取消，请反馈该状态，芝麻将取消用户待履约记录；
	VIOLATED： 用户已违约
	适用场景：如果用户在约定时间（具体根据行业约定，有一定宽限期）内未完成订单支付，反馈该状态，芝麻将对用户记录一条负面记录，请谨慎使用； {"status":"COMPLETE"}
	*/
	BaseAliPayRequest
}

func (r *TradeOrderInfoSyncReq) RequestApi() string {
	return "alipay.trade.orderinfo.sync"
}

func (r *TradeOrderInfoSyncReq) DoValidate() error {
	tradeNoLen := len(r.TradeNo)
	if tradeNoLen == 0 {
		return fmt.Errorf("商户订单号out_trade_no不能为空")
	}
	if tradeNoLen > 64 {
		return fmt.Errorf("当前商户订单号长度为%d，不符合长度1-64要求", tradeNoLen)
	}

	outRequestNo := len(r.OutRequestNo)
	if outRequestNo == 0 {
		return fmt.Errorf("外部请求号out_trade_no不能为空")
	}
	if outRequestNo > 64 {
		return fmt.Errorf("外部请求号out_trade_no参数长度为%d，不符合长度1-64要求", tradeNoLen)
	}
	bizTypeLen := len(r.BizType)

	if bizTypeLen > 64 || bizTypeLen < 0 {
		return fmt.Errorf("参数biz_type参数长度为%d，不符合长度1-64要求", tradeNoLen)
	}
	return nil
}

type TradeOrderInfoSyncRes struct {
	TradeOrderInfoSyncResContent `json:"alipay_trade_orderinfo_sync_response"`
	SignCertSn
}

type TradeOrderInfoSyncResContent struct {
	CommonRes
	TradeNo     string `json:"trade_no,omitempty"`      // 必选 64 支付宝交易号，和商户订单号不能同时为空
	OutTradeNo  string `json:"out_trade_no,omitempty"`  // 必选 64 商户订单号。订单支付时传入的商户订单号,和支付宝交易号不能同时为空。trade_no,out_trade_no如果同时存在优先取trade_no
	BuyerUserId string `json:"buyer_user_id,omitempty"` // 可选 28 买家在支付宝的用户id
}
