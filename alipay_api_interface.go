package alipay

import (
	"encoding/json"
	"time"
)

/**
 * @author: Sam
 * @since: 1.0.0
 * @date: 2023/3/05 17:53
 * @desc:
 */

type Validator interface {
	DoValidate() error
}

// GenericResponse 返回通用数据类型 generic try
type GenericResponse interface {
	//*url.URL | *TradeQueryRes
}

type Requester interface {
	// RequestHttpMethod http请求接口的方法
	RequestHttpMethod() string
	// RequestApiVersion 请求接口的对应的版本号
	RequestApiVersion() string
	// RequestApi 针对alipay请求的接口
	RequestApi() string
}

type IAliPayRequest interface {
	Requester
	Validator
}

type BaseAliPayRequest struct {
}

func (r *BaseAliPayRequest) RequestHttpMethod() string {
	return "POST"
}

func (r *BaseAliPayRequest) RequestApi() string {
	panic("please implement yourself")
}

func (r *BaseAliPayRequest) RequestApiVersion() string {
	return "1.0"
}

type ResponseSigner interface {
	// SetSign 设置返回的签名
	SetSign(sign string)
	GetSign() string

	// SetAlipayCertSn 设置返回的密钥序列号
	SetAlipayCertSn(alipayCertSn string)
	GetAlipayCertSn() string
}

type CommonReqParam struct {
	AppId            string `json:"app_id" url:"app_id"`                                     // 必选	最大长度32 支付宝分配给开发者的应用ID 2014072300007148
	Method           string `json:"method" url:"method"`                                     // 必选	128 接口名称 alipay.trade.page.pay
	Format           string `json:"format,omitempty" url:"format"`                           // 可选	40 仅支持JSON JSON
	ReturnUrl        string `json:"return_url,omitempty" url:"return_url,omitempty"`         // 可选	256 HTTP/HTTPS开头字符串
	Charset          string `json:"charset" url:"charset"`                                   // 必选	10 请求使用的编码格式，如utf-8,gbk,gb2312等
	SignType         string `json:"sign_type" url:"sign_type"`                               // 必选	10 商户生成签名字符串所使用的签名算法类型，目前支持RSA2和RSA，推荐使用RSA2
	Sign             string `json:"sign" url:"sign"`                                         // 必选	344 商户请求参数的签名串
	Timestamp        string `json:"timestamp" url:"timestamp"`                               // 必选	19 发送请求的时间，格式"yyyy-MM-dd HH:mm:ss" 2014-07-24 03:07:50
	Version          string `json:"version" url:"version"`                                   // 必选 3 调用的接口版本，固定为：1.0
	NotifyUrl        string `json:"notify_url,omitempty" url:"notify_url,omitempty"`         // 可选	256	 支付宝服务器主动通知商户服务器里指定的页面http/https路径。
	AppAuthToken     string `json:"app_auth_token,omitempty" url:"app_auth_token,omitempty"` // 可选	40 详见应用授权概述
	BizContent       string `json:"biz_content" url:"biz_content"`                           // 必选	无长度限制 请求参数的集合，最大长度不限，除公共参数外所有请求参数都必须放在这个参数中传递，具体参照各产品快速接入文档
	AppCertSn        string `json:"app_cert_sn" url:"app_cert_sn,omitempty"`                 // 可选	具体参照各产品快速接入文档
	AlipayRootCertSn string `json:"alipay_root_cert_sn" url:"alipay_root_cert_sn,omitempty"` // 可选	具体参照各产品快速接入文档
}

func newCommonParam(method, timestamp string) *CommonReqParam {
	res := CommonReqParam{
		Method:    method,
		Format:    FormatJson,
		Charset:   CharsetUTF8,
		SignType:  SignTypeRSA2,
		Timestamp: timestamp,
		Version:   ApiVersion,
	}
	return &res
}

type commonParamOpt func(*CommonReqParam)

func WithNotifyUrl(notifyUrl string) func(*CommonReqParam) {
	return func(param *CommonReqParam) {
		if len(notifyUrl) == 0 {
			return
		}
		param.NotifyUrl = notifyUrl
	}
}

func WithReturnUrl(returnUrl string) func(*CommonReqParam) {
	return func(param *CommonReqParam) {
		if len(returnUrl) == 0 {
			return
		}
		param.ReturnUrl = returnUrl
	}
}

func WithAppCertSn(appCertSn string) func(*CommonReqParam) {
	return func(param *CommonReqParam) {
		if len(appCertSn) == 0 {
			return
		}
		param.AppCertSn = appCertSn
	}
}

func WithAlipayRootCertSn(alipayRootCertSn string) func(*CommonReqParam) {
	return func(param *CommonReqParam) {
		if len(alipayRootCertSn) == 0 {
			return
		}
		param.AlipayRootCertSn = alipayRootCertSn
	}
}

type CommonRes struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	SubCode string `json:"sub_code,omitempty"`
	SubMsg  string `json:"sub_msg,omitempty"`
}

func (r *CommonRes) SetCode(code string) {
	r.Code = code
}

func (r *CommonRes) SetMsg(msg string) {
	r.Msg = msg
}

func (r *CommonRes) SetSubCode(subCode string) {
	r.SubCode = subCode
}

func (r *CommonRes) SetSubMsg(subMsg string) {
	r.SubMsg = subMsg
}

// Success 文档： https://opendoc.alipay.com/common/02km9f
func (r *CommonRes) Success() bool {
	return r.Code == "10000"
}

func (r *CommonRes) Fail() bool {
	return !r.Success()
}

type SignCertSn struct {
	AlipayCertSn string `json:"alipay_cert_sn,omitempty"`
	Sign         string `json:"sign"` // 签名
}

func (r *SignCertSn) SetSign(sign string) {
	r.Sign = sign
}

func (r *SignCertSn) SetAlipayCertSn(alipayCertSn string) {
	r.AlipayCertSn = alipayCertSn
}

func (r *SignCertSn) GetSign() string {
	return r.Sign
}

func (r *SignCertSn) GetAlipayCertSn() string {
	return r.AlipayCertSn
}

type ErrorResponse struct {
	CommonRes    `json:"error_response"`
	AlipayCertSn string `json:"alipay_cert_sn,omitempty"`
	Sign         string `json:"sign"` // 签名
}

type RequestObjectBuilder interface {
	// RequestApiObject 请求api的对象
	buildRequestObject(req IAliPayRequest, opts ...commonParamOpt) (*CommonReqParam, error)
}
type RequestAliPayObjectBuilder struct {
	// 时区
	location *time.Location
}

func (r *RequestAliPayObjectBuilder) buildRequestObject(req IAliPayRequest, opts ...commonParamOpt) (*CommonReqParam, error) {
	var buff []byte
	var err error
	if buff, err = json.Marshal(req); err != nil {
		return nil, err
	}
	commonParam := newCommonParam(req.RequestApi(), time.Now().In(r.location).Format(time.DateTime))
	for _, paramOpt := range opts {
		paramOpt(commonParam)
	}
	commonParam.BizContent = string(buff)
	return commonParam, nil
}
