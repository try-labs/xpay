package alipay

import "encoding/json"

/**
 * @author: Sam
 * @since:
 * @date: 2023/3/17 15:13
 * @desc:
 */

var _ IAliPayRequest = (*DataBillBalanceQueryReq)(nil)

type DataBillBalanceQueryReq struct {
	BillUserId string `json:"bill_user_id,omitempty"` // 可选	16 目标查询账户（仅支持部分场景，查询自身时候不需要传递当前字段）。 2088123456789012
	BaseAliPayRequest
}

func (r *DataBillBalanceQueryReq) RequestApi() string {
	return "alipay.data.bill.balance.query"
}

func (r *DataBillBalanceQueryReq) DoValidate() error {
	return nil
}

type DataBillBalanceQueryRes struct {
	DataBillBalanceQueryResContent `json:"alipay_data_bill_balance_query_response"`
	SignCertSn
}

type DataBillBalanceQueryResContent struct {
	CommonRes
	TotalAmount     string `json:"total_amount"`            // 必选	32 支付宝账户余额 10000.00
	AvailableAmount string `json:"available_amount"`        // 必选	32 账户可用余额 9000.00
	FreezeAmount    string `json:"freeze_amount"`           // 必选	32 冻结金额 1000.00
	SettleAmount    string `json:"settle_amount,omitempty"` // 可选	32 待结算金额 500.00
}

// /////////////////////////////////////////////
var _ IAliPayRequest = &DataBillBailQueryReq{}

type DataBillBailQueryReq struct {
	StartTime  string `json:"start_time"`             // 必选	20 保证金流水创建时间的起始范围 2019-01-01 00:00:00
	EndTime    string `json:"end_time"`               // 必选	20 保证金流水创建时间的结束范围。与起始时间间隔不超过31天。查询结果为起始时间至结束时间的左闭右开区间 2019-01-02 00:00:00
	BailType   string `json:"bail_type"`              // 必选	20 保证金类型，目前支持TMALL_BAIL-天猫保证金，TAOBAO_BAIL-淘宝保证金 TMALL_BAIL
	TransLogId string `json:"trans_log_id,omitempty"` // 可选	255 保证金流水号。如果查询参数中指定流水号，则只查询流水号相关的记录 20190101***
	BizOrigNo  string `json:"biz_orig_no,omitempty"`  // 可选	255 业务基础订单号。如果查询参数中指定订单号，则只查询相关的记录 1***
	BaseAliPayRequest
}

func (r *DataBillBailQueryReq) RequestApi() string {
	return "alipay.data.bill.bail.query"
}

func (r *DataBillBailQueryReq) DoValidate() error {
	return nil
}

type DataBillBailQueryRes struct {
	DataBillBailQueryResContent `json:"alipay_data_bill_bail_query_response"`
	SignCertSn
}

func (r *DataBillBailQueryRes) String() string {
	buff, _ := json.Marshal(r)
	return string(buff)
}

type DataBillBailQueryResContent struct {
	CommonRes
	DetailList []BailDetailResult `json:"detail_list"`
}

type BailDetailResult struct {
	TransDt    string `json:"trans_dt,omitempty"`     // 必选	20 业务发生时间 2019-01-01 00:00:00
	TransLogId string `json:"trans_log_id,omitempty"` // 必选	255 保证金业务流水号 20190101***
	BailType   string `json:"bail_type,omitempty"`    // 必选	255 保证金类型描述，仅供参考 天猫保证金
	Amount     string `json:"amount,omitempty"`       // 必选	32 保证金收支金额 10.00
	Balance    string `json:"balance,omitempty"`      // 必选	32 保证金余额 1000.00
	Memo       string `json:"memo,omitempty"`         // 可选	255 保证金说明 保证金冻结
	BizDesc    string `json:"biz_desc,omitempty"`     // 可选	255 业务描述，资金收支对应的详细业务场景信息 余额账户迁入
	BizOrigNo  string `json:"biz_orig_no,omitempty"`  // 可选	255 业务基础订单号，资金收支对应的原始业务订单唯一识别编号 1***
}
