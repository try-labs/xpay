package alipay

/**
 * @author: Sam
 * @since: 1.0.0
 * @date: 2023/3/15 10:40
 * @desc: async notify
 */

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"
)

type NotifyReq struct {
	NotifyTime       string `json:"notify_time"`                  // 必填 通知的发送时间。格式为 yyyy-MM-dd HH:mm:ss 2018-10-21 15:45:22
	NotifyType       string `json:"notify_type"`                  // 必填 64 通知类型 trade_status_sync
	NotifyId         string `json:"notify_id"`                    // 必填  128 通知校验 ID ac05099524730693a8b330c45cf72da943
	Charset          string `json:"charset"`                      // 必填 10 编码格式。如 utf-8、gbk、gb312等。 utf-8
	Version          string `json:"version"`                      // 必填  3 调用的接口版本。固定为1.0
	SignType         string `json:"sign_type"`                    // 必填 10 签名类型。签名算法类型，目前支持RSA2和RSA，推荐使用 RSA2
	Sign             string `json:"sign"`                         // 必填  344 签名。详情可查看 异步返回结果的验签 601510b7970e52cc63db0f44997cf70e
	AuthAppId        string `json:"auth_app_id"`                  // 必填 32 授权方的APPID。由于本接口暂不开放第三方应用授权，因此 auth_app_id=app_id
	TradeNo          string `json:"trade_no"`                     // 必填 64 支付宝交易号，支付宝交易凭证号。 2013112011001004330000121536
	AppId            string `json:"app_id"`                       // 必填  32 支付宝应用的APPID。支付宝分配给开发者的应用 ID 2019082200007148
	OutTradeNo       string `json:"out_trade_no"`                 // 必填 64 商家订单号。原支付请求的商家订单号 6823789339978248
	OutBizNo         string `json:"out_biz_no,omitempty"`         // 可选 64 商家业务号。商家业务ID，通常是退款通知中返回的退款申请流水号 HZRF001
	BuyerId          string `json:"buyer_id,omitempty"`           // 可选  16 买家支付宝账号 ID。以 2088 开头的纯 16 位数字 20881***524333
	SellerId         string `json:"seller_id,omitempty"`          // 可选  30 卖家支付宝账号 ID。以 2088 开头的纯 16 位数字 20881***2239364
	TradeStatus      string `json:"trade_status,omitempty"`       // 可选 32 交易状态。交易目前所处状态，详情可查看下表 交易状态说明 TRADE_CLOSED
	TotalAmount      string `json:"total_amount,omitempty"`       // 可选 11 订单金额。本次交易支付订单金额，单位为人民币（元），精确到小数点后 2 位 20.00
	ReceiptAmount    string `json:"receipt_amount,omitempty"`     // 可选 11 实收金额。商家在交易中实际收到的款项，单位为人民币（元），精确到小数点后 2 位 15.00
	InvoiceAmount    string `json:"invoice_amount,omitempty"`     // 可选 11 开票金额。用户在交易中支付的可开发票的金额，单位为人民币（元），精确到小数点后 2 位 13.88
	BuyerPayAmount   string `json:"buyer_pay_amount,omitempty"`   // 可选 11 用户在交易中支付的金额，单位为人民币（元），精确到小数点后 2 位 12.00
	PointAmount      string `json:"point_amount,omitempty"`       // 可选 11 使用集分宝支付金额，单位为人民币（元），精确到小数点后 2 位 12.00
	RefundFee        string `json:"refund_fee,omitempty"`         // 可选 11 总退款金额。退款通知中，返回总退款金额，单位为人民币（元），精确到小数点后 2 位 2.58
	Subject          string `json:"subject,omitempty"`            // 可选 256 订单标题/商品标题/交易标题/订单关键字等，是请求时对应参数，会在通知中原样传回 XXXX交易
	Body             string `json:"body,omitempty"`               // 可选 400 商品描述。该订单的备注、描述、明细等。对应请求时的 body 参数，会在通知中原样传回 XXX交易内容
	GmtCreate        string `json:"gmt_create,omitempty"`         // 可选 交易创建时间。格式为 yyyy-MM-dd HH:mm:ss 2018-08-25 15:34:42
	GmtPayment       string `json:"gmt_payment,omitempty"`        // 可选 交易付款时间。格式为 yyyy-MM-dd HH:mm:ss 2018-08-25 15:34:42
	GmtRefund        string `json:"gmt_refund,omitempty"`         // 可选 交易退款时间。格式为 yyyy-MM-dd HH:mm:ss.S 2018-08-26 10:34:44.340
	GmtClose         string `json:"gmt_close,omitempty"`          // 可选 交易结束时间。格式为 yyyy-MM-dd HH:mm:ss 2018-08-26 16:32:30
	FundBillList     string `json:"fund_bill_list,omitempty"`     // 可选 512 支付金额信息。支付成功的各个渠道金额信息。详情可查看下文 资金明细信息说明 [{"amount":"15.00","fundChannel":"ALIPAYACCOUNT"}]
	VocherDetailList string `json:"vocher_detail_list,omitempty"` // 可选 512 优惠券信息。本交易支付时所使用的所有优惠券信息。详情可查看下表 优惠券信息说明 [{"amount":"0.20","merchantContribute":"0.00","name":"一键创建券模板名称","otherContribute":"0.20","type":"ALIAPY_BIZ_VOUCHER","memo":"学生卡8折优惠"}]
	PassbackParams   string `json:"passback_params,omitempty"`    // 可选 512 回传参数，公共回传参数，如果请求时传递了该参数，则返回的异步通知会原样传回。本参数必须进行 UrlEncode 之后才可传入。 merchantBizType%3d3C%26merchantBizNo%3d201601001111
	// 证书签名特有
	AlipayCertSn string `json:"alipay_cert_sn,omitempty"`
}

func (r NotifyReq) String() string {
	buff, _ := json.Marshal(r)
	return string(buff)
}

type VocherDetail struct {
	VoucherId  string `json:"voucher_id"`            // 必填 32 券 ID 2015102600073002039000002D5O
	TemplateId string `json:"template_id,omitempty"` // 可选 64 券模板 ID 20171030000730015359000EMZP0
	Name       string `json:"name"`                  // 必填 64 券名称 5元代金券
	Type       string `json:"type"`                  // 必填 32 优惠类型。当前支持以下几种主要类型：
	//ALIPAY_BIZ_VOUCHER：商家全场券。
	//ALIPAY_COMMON_ITEM_VOUCHER：商家单品券。
	//ALIPAY_CASH_VOUCHER：平台优惠券，支付宝或第三方出资。
	//ALICREDIT_INTFREE_VOUCHER：花呗分期券，该券仅做订单外的工作呗分期费用减免，并不抵扣订单内支付金额。
	//注意：不排除未来新增其它类型的可能，商家接入时请注意兼容性，避免硬编码。
	Amount                string              `json:"amount"`                            // 必填 11 优惠金额。优惠金额中，由商家出资的金额。 10.00
	MerchantContribute    string              `json:"merchant_contribute,omitempty"`     // 可选 11 商家出资金额。优惠金额中，由商家出资的金额。 9.00
	OtherContribute       string              `json:"other_contribute,omitempty"`        // 可选 11 其他出资方出资金额。可能是支付宝，可能是品牌商，或者其他方，也可能是他们的共同出资。 1.00
	OtherContributeDetail []*ContributeDetail `json:"other_contribute_detail,omitempty"` // 可选 优惠券的其他出资方明细
	Memo                  string              `json:"memo,omitempty"`                    // 可选  256 优惠券备注信息。 学生专用优惠
}

type ContributeDetail struct {
	ContributeType   string `json:"contribute_type,omitempty"`   // 可选 32 出资方类型，如品牌商出资、支付宝平台出资等。 PLATFORM
	ContributeAmount string `json:"contribute_amount,omitempty"` // 可选 8 出资方金额
}

// 通知逻辑
func (r *Client) doNotify(request *http.Request) (*NotifyReq, error) {
	var err error
	if err = request.ParseForm(); err != nil {
		return nil, err
	}
	urlValues := request.Form
	notifyParamMap := make(map[string]string, len(urlValues))
	for k, v := range urlValues {
		if len(v) == 1 {
			notifyParamMap[k] = v[0]
		}
	}
	log.Println("待验签参数", notifyParamMap)
	var buff []byte
	if buff, err = json.Marshal(notifyParamMap); err != nil {
		return nil, err
	}
	var notifyParam = new(NotifyReq)
	if err = json.Unmarshal(buff, notifyParam); err != nil {
		return nil, err
	}
	var keyValueList = make([]string, 0, len(notifyParamMap))
	for key, value := range notifyParamMap {
		// 删除sign、sign_type
		if key == ExcludeKeySign || key == ExcludeKeySignType {
			continue
		}
		keyValueList = append(keyValueList, key+"="+value)
	}
	sort.Strings(keyValueList)
	if err = r.VerifySign(AsyncVerificationScene, notifyParam.Sign, []byte(strings.Join(keyValueList, "&")), notifyParam.AlipayCertSn); err != nil {
		log.Println("校验参数err", err)
		return nil, err
	}
	log.Println("verification success")
	return notifyParam, err
}

// AsyncNotify 异步通知
func (r *Client) AsyncNotify(request *http.Request) (*NotifyReq, error) {
	log.Println("async notify verification ")
	return r.doNotify(request)
}

// SyncNotify 同步通知
func (r *Client) SyncNotify(request *http.Request) (*NotifyReq, error) {
	log.Println("sync notify verification ")
	return r.doNotify(request)
}
