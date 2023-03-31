package alipay

import (
	"encoding/json"
	"fmt"
)

/**
 * @author: Sam
 * @since:
 * @date: 2023/3/25 13:13
 * @desc:
 */

var _ IAliPayRequest = &CommerceCityFacilitatorVoucherGenerateReq{}

type CommerceCityFacilitatorVoucherGenerateReq struct {
	CityCode    string `json:"city_code"`            // 必选	30 城市编码请参考查询 中华人民共和国行政区划代码。 已支持城市：广州 440100，深圳 440300，杭州330100。
	TradeNo     string `json:"trade_no"`             // 必选	100 支付宝交易号（交易支付时，必须通过指定sellerId：2088121612215201，将钱支付到指定的中间户中）
	TradeFee    string `json:"trade_fee"`            // 必选	20 订单总金额，元为单位 10.00
	TicketNum   string `json:"ticket_num"`           // 必选	20 地铁票购票数量 5
	TicketType  string `json:"ticket_type"`          // 必选	60 地铁票种类，枚举支持： *oneway。
	SiteBegin   string `json:"site_begin,omitempty"` // 可选	30 起点站站点编码 02490301
	SiteEnd     string `json:"site_end,omitempty"`   // 可选	30 终点站站点编码 02490305
	TicketPrice string `json:"ticket_price"`         // 必选	20 单张票价，元为单价 5.00
	baseAliPayRequest
}

func (r *CommerceCityFacilitatorVoucherGenerateReq) DoValidate() error {
	cityCodeLen := len(r.CityCode)
	if cityCodeLen == 0 || cityCodeLen > 30 {
		return fmt.Errorf("城市编码city_code参数长度为%d，不符合长度1-30个字符的要求", cityCodeLen)
	}
	tradeNoLen := len(r.TradeNo)
	if tradeNoLen == 0 || tradeNoLen > 100 {
		return fmt.Errorf("支付宝交易号trade_no参数长度为%d，不符合长度1-100个字符的要求", tradeNoLen)
	}
	return nil
}

func (r *CommerceCityFacilitatorVoucherGenerateReq) RequestApi() string {
	return "alipay.commerce.cityfacilitator.voucher.generate"
}

type CommerceCityFacilitatorVoucherGenerateRes struct {
	CommerceCityFacilitatorVoucherGenerateResContent `json:"alipay_commerce_cityfacilitator_voucher_generate_response"`
	SignCertSn
}
type CommerceCityFacilitatorVoucherGenerateResContent struct {
	CommonRes
	QrCodeNo    string `json:"qr_code_no"`   //	必选	200 地铁购票二维码编码，可自定义
	TicketNo    string `json:"ticket_no"`    //	必选	20 地铁购票的核销码 0243567832
	ExpiredDate string `json:"expired_date"` //	必选	30 核销码过期时间 2016-03-18 00:00:00
}

func (r *CommerceCityFacilitatorVoucherGenerateRes) String() string {
	buff, _ := json.Marshal(r)
	return string(buff)
}

//////////////////////////////////////////////////////

var _ IAliPayRequest = &CommerceCityFacilitatorVoucherRefundReq{}

type CommerceCityFacilitatorVoucherRefundReq struct {
	CityCode string `json:"city_code"` // 必选	30 城市编码请参考查询 中华人民共和国行政区划代码。 已支持城市：广州 440100，深圳 440300，杭州330100。
	TradeNo  string `json:"trade_no"`  // 必选	100 支付宝交易号（交易支付时，必须通过指定sellerId：2088121612215201，将钱支付到指定的中间户中）
	baseAliPayRequest
}

func (r *CommerceCityFacilitatorVoucherRefundReq) DoValidate() error {
	cityCodeLen := len(r.CityCode)
	if cityCodeLen == 0 || cityCodeLen > 30 {
		return fmt.Errorf("城市编码city_code参数长度为%d，不符合长度1-30个字符的要求", cityCodeLen)
	}
	tradeNoLen := len(r.TradeNo)
	if tradeNoLen == 0 || tradeNoLen > 100 {
		return fmt.Errorf("支付宝交易号trade_no参数长度为%d，不符合长度1-100个字符的要求", tradeNoLen)
	}
	return nil
}

func (r *CommerceCityFacilitatorVoucherRefundReq) RequestApi() string {
	return "alipay.commerce.cityfacilitator.voucher.refund"
}

type CommerceCityFacilitatorVoucherRefundRes struct {
	CommerceCityFacilitatorVoucherRefundResContent `json:"alipay_commerce_cityfacilitator_voucher_refund_response"`
	SignCertSn
}

type CommerceCityFacilitatorVoucherRefundResContent struct {
	CommonRes
}

func (r *CommerceCityFacilitatorVoucherRefundRes) String() string {
	buff, _ := json.Marshal(r)
	return string(buff)
}

/////////////////////////////////////////////////

var _ IAliPayRequest = &CommerceCityFacilitatorStationQueryReq{}

type CommerceCityFacilitatorStationQueryReq struct {
	CityCode string `json:"city_code"` // 必选	30 城市编码请参考查询 中华人民共和国行政区划代码。 已支持城市：广州 440100，深圳 440300，杭州330100。
	baseAliPayRequest
}

func (r *CommerceCityFacilitatorStationQueryReq) DoValidate() error {
	cityCodeLen := len(r.CityCode)
	if cityCodeLen == 0 || cityCodeLen > 30 {
		return fmt.Errorf("城市编码city_code参数长度为%d，不符合长度1-30个字符的要求", cityCodeLen)
	}
	return nil
}

func (r *CommerceCityFacilitatorStationQueryReq) RequestApi() string {
	return "alipay.commerce.cityfacilitator.station.query"
}

type CommerceCityFacilitatorStationQueryRes struct {
	CommerceCityFacilitatorStationQueryResContent `json:"alipay_commerce_cityfacilitator_station_query_response"`
	SignCertSn
}
type CommerceCityFacilitatorStationQueryResContent struct {
	CommonRes
	SupportStarts []StationDetailInfo `json:"support_starts"` // 必选 支持设为起点的站点列表
}

func (r *CommerceCityFacilitatorStationQueryRes) String() string {
	buff, _ := json.Marshal(r)
	return string(buff)
}

type StationDetailInfo struct {
	Name    string `json:"name"`     // 必选	30 站点中文名称 蓝村路
	Code    string `json:"code"`     // 必选	30 站点编码 12000020
	ExtCode string `json:"ext_code"` // 必选	30 站点外部编码 260000210002
}

/////////////////////////////////////

var _ IAliPayRequest = &CommerceCityFacilitatorVoucherBatchQueryReq{}

type CommerceCityFacilitatorVoucherBatchQueryReq struct {
	CityCode string   `json:"city_code"` // 必选	30 城市编码请参考查询 中华人民共和国行政区划代码。 已支持城市：广州 440100，深圳 440300，杭州330100。
	TradeNos []string `json:"trade_nos"` // 必选	800 支付宝交易号列表
	baseAliPayRequest
}

func (r *CommerceCityFacilitatorVoucherBatchQueryReq) DoValidate() error {
	cityCodeLen := len(r.CityCode)
	if cityCodeLen == 0 || cityCodeLen > 30 {
		return fmt.Errorf("城市编码city_code参数长度为%d，不符合长度1-30个字符的要求", cityCodeLen)
	}
	return nil
}
func (r *CommerceCityFacilitatorVoucherBatchQueryReq) RequestApi() string {
	return "alipay.commerce.cityfacilitator.voucher.batchquery"
}

type CommerceCityFacilitatorVoucherBatchQueryRes struct {
	CommerceCityFacilitatorVoucherBatchQueryResContent `json:"alipay_commerce_cityfacilitator_voucher_batchquery_response"`
	SignCertSn
}

type CommerceCityFacilitatorVoucherBatchQueryResContent struct {
	CommonRes
	Tickets []TicketDetailInfo `json:"tickets"`
}

func (r *CommerceCityFacilitatorVoucherBatchQueryRes) String() string {
	buff, _ := json.Marshal(r)
	return string(buff)
}

type TicketDetailInfo struct {
	TradeNo          string `json:"trade_no"`                     // 必选	100 支付宝交易号 0123456789
	Amount           string `json:"amount"`                       // 必选	20 总金额，元为单位 10.00
	StartStation     string `json:"start_station,omitempty"`      // 可选	80 起点站编码 12300002
	EndStation       string `json:"end_station,omitempty"`        // 可选	80 终点站编码 21003002
	Quantity         string `json:"quantity"`                     // 必选	20 票数量 8
	Status           string `json:"status"`                       // 必选	40 订单状态 SUCCESS
	TicketPrice      string `json:"ticket_price"`                 // 必选	20 单价，元为单位 5.00
	StartStationName string `json:"start_station_name,omitempty"` // 可选	80 起点站中文名称 蓝村路
	EndStationName   string `json:"end_station_name,omitempty"`   // 可选	80 终点站中文名称 浦电路
	TicketType       string `json:"ticket_type"`                  // 必选	60 票类型 oneway
}
