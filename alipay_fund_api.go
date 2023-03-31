package alipay

import (
	"encoding/json"
	"fmt"
)

/**
 * @author: Sam
 * @since:
 * @date: 2023/3/22 17:35
 * @desc: 转账到
 */

var _ IAliPayRequest = &FundAccountQueryReq{}

type FundAccountQueryReq struct {
	AlipayUserId string `json:"alipay_user_id"`         //	必选	28 支付宝会员 id。 2088301409188095
	AccountType  string `json:"account_type,omitempty"` // 特殊可选	30 查询的账号类型，查询余额账户值为ACCTRANS_ACCOUNT。必填。
	baseAliPayRequest
}

func (r *FundAccountQueryReq) DoValidate() error {
	if len(r.AlipayUserId) == 0 {
		return fmt.Errorf("支付宝会员alipay_user_id 参数必传")
	}
	return nil
}

func (r *FundAccountQueryReq) RequestApi() string {
	return "alipay.fund.account.query"
}

type FundAccountQueryRes struct {
	FundAccountQueryResContent `json:"alipay_fund_account_query_response"`
	SignCertSn
}

type FundAccountQueryResContent struct {
	CommonRes
	AvailableAmount string `json:"available_amount"`        // 必选	15 账户可用余额，单位元，精确到小数点后两位。 26.45
	FreezeAmount    string `json:"freeze_amount,omitempty"` // 可选	15 当前支付宝账户的实时冻结余额 11.11
}

func (r *FundAccountQueryRes) String() string {
	buff, _ := json.Marshal(r)
	return string(buff)
}

//////////////////////////////////////////////

var _ IAliPayRequest = &FundTransToAccountTransferReq{}

type FundTransToAccountTransferReq struct {
	OutBizNo  string `json:"out_biz_no"` // 必选	64 商户转账唯一订单号。发起转账来源方定义的转账单据ID，用于将转账回执通知给来源方。 不同来源方给出的ID可以重复，同一个来源方必须保证其ID的唯一性。 只支持半角英文、数字，及“-”、“_”。
	PayeeType string `json:"payee_type"` // 必选	20 收款方账户类型。可取值：
	/*
		1、ALIPAY_USERID：支付宝账号对应的支付宝唯一用户号。以2088开头的16位纯数字组成。
		2、ALIPAY_LOGONID：支付宝登录号，支持邮箱和手机号格式。
		2、ALIPAY_OPENID：支付宝openid
	*/
	PayeeAccount  string `json:"payee_account"`             // 必选	100 收款方账户。与payee_type配合使用。付款方和收款方不能是同一个账户。
	Amount        string `json:"amount"`                    // 必选	16 转账金额，单位：元。 只支持2位小数，小数点前最大支持13位，金额必须大于等于0.1元。 最大转账金额以实际签约的限额为准。
	PayerShowName string `json:"payer_show_name,omitempty"` // 可选	100 付款方姓名（最长支持100个英文/50个汉字）。显示在收款方的账单详情页。如果该字段不传，则默认显示付款方的支付宝认证姓名或单位名称。
	PayeeRealName string `json:"payee_real_name,omitempty"` // 可选	100 收款方真实姓名（最长支持100个英文/50个汉字）。 如果本参数不为空，则会校验该账户在支付宝登记的实名是否与收款方真实姓名一致。
	Remark        string `json:"remark,omitempty"`          // 可选	200 转账备注（支持200个英文/100个汉字）。 当付款方为企业账户，且转账金额达到（大于等于）50000元，remark不能为空。收款方可见，会展示在收款用户的收支详情中。
	baseAliPayRequest
}

func (r *FundTransToAccountTransferReq) DoValidate() error {
	if len(r.OutBizNo) == 0 {
		return fmt.Errorf("商户转账唯一订单号out_biz_no参数必传")
	}
	return nil
}

func (r *FundTransToAccountTransferReq) RequestApi() string {
	return "alipay.fund.trans.toaccount.transfer"
}

type FundTransToAccountTransferRes struct {
	FundTransToAccountTransferResContent `json:"alipay_fund_trans_toaccount_transfer_response"`
	SignCertSn
}

type FundTransToAccountTransferResContent struct {
	CommonRes
	OutBizNo string `json:"out_biz_no"`         // 必选	64 商户转账唯一订单号：发起转账来源方定义的转账单据号。请求时对应的参数，原样返回。
	OrderId  string `json:"order_id,omitempty"` // 可选	64 支付宝转账单据号，成功一定返回，失败可能不返回也可能返回。
	PayDate  string `json:"pay_date,omitempty"` // 可选	20 支付时间：格式为yyyy-MM-dd HH:mm:ss，仅转账成功返回。 2013-01-01 08:08:08
}

func (r *FundTransToAccountTransferRes) String() string {
	buff, _ := json.Marshal(r)
	return string(buff)
}

/////////////////////////////////////////////

type FundTransOrderQueryReq struct {
	OutBizNo string `json:"out_biz_no"`         // 必选	64 商户转账唯一订单号：发起转账来源方定义的转账单据号。请求时对应的参数，原样返回。
	OrderId  string `json:"order_id,omitempty"` // 可选	64 支付宝转账单据号，成功一定返回，失败可能不返回也可能返回。
	baseAliPayRequest
}

func (r *FundTransOrderQueryReq) DoValidate() error {
	if len(r.OutBizNo) == 0 {
		return fmt.Errorf("商户转账唯一订单号out_biz_no参数必传")
	}
	return nil
}

func (r *FundTransOrderQueryReq) RequestApi() string {
	return "alipay.fund.trans.order.query"
}

type FundTransOrderQueryRes struct {
	FundTransOrderQueryResContent `json:"alipay_fund_trans_order_query_response"`
	SignCertSn
}

type FundTransOrderQueryResContent struct {
	CommonRes
	OrderId string `json:"order_id,omitempty"` // 可选	64 支付宝转账单据号，成功一定返回，失败可能不返回也可能返回。
	Status  string `json:"status,omitempty"`   // 可选	10 转账单据状态。
	/*
		SUCCESS：成功（配合"单笔转账到银行账户接口"产品使用时, 同一笔单据多次查询有可能从成功变成退票状态）；
		FAIL：失败（具体失败原因请参见error_code以及fail_reason返回值）；
		INIT：等待处理；
		DEALING：处理中；
		REFUND：退票（仅配合"单笔转账到银行账户接口"产品使用时会涉及, 具体退票原因请参见fail_reason返回值）；
		UNKNOWN：状态未知。
	*/
	PayDate        string `json:"pay_date,omitempty"`         // 可选	20 支付时间，格式为yyyy-MM-dd HH:mm:ss，转账失败不返回。 2013-01-01 08:08:08
	ArrivalTimeEnd string `json:"arrival_time_end,omitempty"` // 可选	20 预计到账时间，转账到银行卡专用，格式为yyyy-MM-dd HH:mm:ss，转账受理失败不返回。 注意： 此参数为预计时间，可能与实际到账时间有较大误差，不能作为实际到账时间使用，仅供参考用途。
	OrderFee       string `json:"order_fee,omitempty"`        // 可选	20 预计收费金额（元），转账到银行卡专用，数字格式，精确到小数点后2位，转账失败或转账受理失败不返回。 0.02
	FailReason     string `json:"fail_reason,omitempty"`      // 可选	100 查询到的订单状态为FAIL失败或REFUND退票时，返回具体的原因。 单笔额度超限
	OutBizNo       string `json:"out_biz_no,omitempty"`       // 可选	64 发起转账来源方定义的转账单据号。 该参数的赋值均以查询结果中 的 out_biz_no 为准。 如果查询失败，不返回该参数。
	ErrorCode      string `json:"error_code,omitempty"`       // 可选	100 查询失败时，本参数为错误代 码。 查询成功不返回。 对于退票订单，不返回该参数。

}

func (r *FundTransOrderQueryRes) String() string {
	buff, _ := json.Marshal(r)
	return string(buff)
}

/////////////////////////////////////////////

var _ IAliPayRequest = &FundTransUniTransferReq{}

type FundTransUniTransferReq struct {
	OutBizNo       string      `json:"out_biz_no"`                // 必选	64 商家侧唯一订单号，由商家自定义。对于不同转账请求，商家需保证该订单号在自身系统唯一。
	TransAmount    string      `json:"trans_amount"`              // 必选	20 订单总金额，单位为元，不支持千位分隔符，精确到小数点后两位，取值范围[0.1,100000000]。
	ProductCode    string      `json:"product_code"`              // 必选	64 销售产品码。单笔无密转账固定为 TRANS_ACCOUNT_NO_PWD。
	BizScene       string      `json:"biz_scene"`                 // 必选	64 业务场景。单笔无密转账固定为 DIRECT_TRANSFER。
	OrderTitle     string      `json:"order_title"`               // 必选	128 转账业务的标题，用于在支付宝用户的账单里显示。
	PayeeInfo      Participant `json:"payee_info"`                // 必选        收款方信息
	Remark         string      `json:"remark,omitempty"`          // 可选	200 业务备注。 201905代发
	BusinessParams string      `json:"business_params,omitempty"` // 可选	2048 转账业务请求的扩展参数，支持传入的扩展参数如下：
	/*
		payer_show_name_use_alias：是否展示付款方别名，可选，收款方在支付宝账单中可见。枚举支持：
		* true：展示别名，将展示商家支付宝在商家中心 商户信息 > 商户基本信息 页面配置的 商户别名。
		* false：不展示别名。默认为 false。
		{"payer_show_name_use_alias":"true"}
	*/
	baseAliPayRequest
}

func (r *FundTransUniTransferReq) DoValidate() error {
	if len(r.OutBizNo) == 0 {
		return fmt.Errorf("商户转账唯一订单号out_biz_no参数必传")
	}
	return nil
}

func (r *FundTransUniTransferReq) RequestApi() string {
	return "alipay.fund.trans.uni.transfer"
}

type Participant struct {
	Identity string `json:"identity"` // 	必选	64 参与方的标识 ID。
	/*
		当 identity_type=ALIPAY_USER_ID 时，填写支付宝用户 UID。示例值：2088123412341234。
		当 identity_type=ALIPAY_LOGON_ID 时，填写支付宝登录号。示例值：186xxxxxxxx。
	*/
	IdentityType string `json:"identity_type"` //	必选	64 参与方的标识类型，目前支持
	/* 如下枚举：
	ALIPAY_USER_ID：支付宝会员的用户 ID，可通过 获取会员信息 能力获取。
	ALIPAY_LOGON_ID：支付宝登录号，支持邮箱和手机号格式。
	*/
	Name string `json:"name,omitempty"` // 可选	128 参与方真实姓名。如果非空，将校验收款支付宝账号姓名一致性。 当 identity_type=ALIPAY_LOGON_ID 时，本字段必填。若传入该属性，则在支付宝回单中将会显示这个属性。
}

type FundTransUniTransferRes struct {
	FundTransUniTransferResContent `json:"alipay_fund_trans_uni_transfer_response"`
	SignCertSn
}
type FundTransUniTransferResContent struct {
	CommonRes
	OutBizNo       string `json:"out_biz_no"`        // 必选	64 商户订单号
	OrderId        string `json:"order_id"`          // 必选	32 支付宝转账订单号
	PayFundOrderId string `json:"pay_fund_order_id"` // 必选	32 支付宝支付资金流水号
	Status         string `json:"status,omitempty"`  // 可选	32 转账单据状态。 SUCCESS（该笔转账交易成功）：成功； FAIL：失败（具体失败原因请参见error_code以及fail_reason返回值）；
	TransDate      string `json:"trans_date"`        // 必选	30 订单支付时间，格式为yyyy-MM-dd HH:mm:ss 2019-08-21 00:00:00
}

func (r *FundTransUniTransferRes) String() string {
	buff, _ := json.Marshal(r)
	return string(buff)
}

///////////////////////////////////

var _ IAliPayRequest = &FundTransCommonQueryReq{}

type FundTransCommonQueryReq struct {
	ProductCode string `json:"product_code,omitempty"` //	可选	64
	/*
		销售产品码，商家和支付宝签约的产品码，如果传递了out_biz_no则该字段为必传。可传值如下：
		STD_RED_PACKET：现金红包
		TRANS_ACCOUNT_NO_PWD：单笔无密转账到支付宝账户
		TRANS_BANKCARD_NO_PWD：单笔无密转账到银行卡
	*/
	BizScene string `json:"biz_scene,omitempty"` //  可选 64 描述特定的业务场景，如果传递了out_biz_no则该字段为必传。可取的业务场景如下：
	/*
		PERSONAL_PAY：C2C现金红包-发红包；
		PERSONAL_COLLECTION：C2C现金红包-领红包；
		REFUND：C2C现金红包-红包退回；
		DIRECT_TRANSFER：B2C现金红包、单笔无密转账
	*/
	OutBizNo string `json:"out_biz_no,omitempty"` // 可选	64 商户转账唯一订单号，发起转账来源方定义的转账单据ID。
	/*
		本参数和order_id（支付宝转账单据号）、pay_fund_order_id（支付宝支付资金流水号）三者不能同时为空。
		当三者同时传入时，将用pay_fund_order_id（支付宝支付资金流水号）进行查询，忽略其余两者；
		当本参数和支付宝转账单据号同时提供时，将用支付宝转账单据号进行查询，忽略本参数。
	*/
	OrderId string `json:"order_id,omitempty"` // 可选	32
	/*
		支付宝转账单据号。 本参数和out_biz_no（商户转账唯一订单号）、pay_fund_order_id（支付宝支付资金流水号）三者不能同时为空。
		当三者同时传入时，将用pay_fund_order_id（支付宝支付资金流水号）进行查询，忽略其余两者；
		当本参数和pay_fund_order_id（支付宝支付资金流水号）同时提供时，将用支付宝支付资金流水号进行查询，忽略本参数；
		当本参数和out_biz_no（商户转账唯一订单号）同时提供时，将用本参数进行查询，忽略商户转账唯一订单号。
	*/
	PayFundOrderId string `json:"pay_fund_order_id,omitempty"` // 可选	32
	/*
		支付宝支付资金流水号。本参数和支付宝转账单据号、商户转账唯一订单号三者不能同时为空。
		当本参数和out_biz_no（商户转账唯一订单号）、order_id（支付宝转账单据号）同时提供时，将用本参数进行查询，忽略其余两者；
		当本参数和order_id（支付宝转账单据号）同时提供时，将用本参数进行查询，忽略支付宝转账单据号；
		当本参数和out_biz_no（商户转账唯一订单号）同时提供时，将用本参数进行查询，忽略商户转账唯一订单号。
	*/
	baseAliPayRequest
}

func (r *FundTransCommonQueryReq) DoValidate() error {
	if len(r.OutBizNo) == 0 {
		return fmt.Errorf("商户转账唯一订单号out_biz_no参数必传")
	}
	return nil
}

func (r *FundTransCommonQueryReq) RequestApi() string {
	return "alipay.fund.trans.common.query"
}

type FundTransCommonQueryRes struct {
	FundTransCommonQueryResContent `json:"alipay_fund_trans_common_query_response"`
	SignCertSn
}
type FundTransCommonQueryResContent struct {
	CommonRes
	OrderId            string `json:"order_id"`                        // 必选	32 支付宝转账单据号，查询失败不返回。
	PayFundOrderId     string `json:"pay_fund_order_id,omitempty"`     // 可选	32 支付宝支付资金流水号，转账失败不返回。
	OutBizNo           string `json:"out_biz_no"`                      // 必选	64 商户订单号
	TransAmount        string `json:"trans_amount,omitempty"`          // 可选	16 付款金额，收银台场景下付款成功后的支付金额，订单状态为SUCCESS才返回，其他状态不返回。 付款金额，单位为元，精确到小数点后两位：32.00
	Status             string `json:"status"`                          // 必选	64 转账单据状态。可能出现的状态如下： SUCCESS：转账成功； WAIT_PAY：等待支付； CLOSED：订单超时关闭； FAIL：失败（适用于"单笔转账到银行卡"）； DEALING：处理中（适用于"单笔转账到银行卡"）； REFUND：退票（适用于"单笔转账到银行卡"）； alipay.fund.trans.app.pay涉及的状态： WAIT_PAY、SUCCESS、CLOSED alipay.fund.trans.refund涉及的状态：SUCCESS alipay.fund.trans.uni.transfer涉及的状态：SUCCESS、FAIL、DEALING、REFUND
	PayDate            string `json:"pay_date,omitempty"`              // 可选	20 支付时间，格式为yyyy-MM-dd HH:mm:ss，转账失败不返回。2013-01-01 08:08:08
	ArrivalTimeEnd     string `json:"arrival_time_end,omitempty"`      // 可选	20 预计到账时间，转账到银行卡专用，格式为yyyy-MM-dd HH:mm:ss，转账受理失败不返回。 注意： 此参数为预计时间，可能与实际到账时间有较大误差，不能作为实际到账时间使用，仅供参考用途。
	OrderFee           string `json:"order_fee,omitempty"`             // 可选	20 预计收费金额（元），转账到银行卡专用，数字格式，精确到小数点后2位，转账失败或转账受理失败不返回。 0.02
	ErrorCode          string `json:"error_code,omitempty"`            // 可选	64 查询到的订单状态为FAIL失败或REFUND退票时，返回错误代码
	FailReason         string `json:"fail_reason,omitempty"`           // 可选	128 查询到的订单状态为FAIL失败或REFUND退票时，返回具体的原因。
	SubOrderErrorCode  string `json:"sub_order_error_code,omitempty"`  // 可选	64 特殊场景提供，当子单出现异常导致主单失败或者退款时，会提供此字段，用于透出子单具体的错误场景
	SubOrderFailReason string `json:"sub_order_fail_reason,omitempty"` // 可选	128 特殊场景提供，当子单出现异常导致主单失败或者退款时，会提供此字段，用于透出子单具体的错误场景
	SubOrderStatus     string `json:"sub_order_status,omitempty"`      // 可选	64 退款子单失败状态
}

func (r *FundTransCommonQueryRes) String() string {
	buff, _ := json.Marshal(r)
	return string(buff)
}
