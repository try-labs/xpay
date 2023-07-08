package alipay

import (
	"errors"
)

/**
 * @author: Sam
 * @since: 1.0.0
 * @date: 2023/3/12 17:46
 * @desc: alipay const
 */

const (
	AppId           = ""
	PrivateKey      = ""
	PublicKey       = ""
	AlipayPublicKey = ""
)

const (
	OtherAppId      = ""
	OtherPrivateKey = ""
)
const LineBreak = "\n"

const (
	SandboxGatewayURL    = "https://openapi.alipaydev.com/gateway.do"
	NewSandboxServerUrl  = "https://openapi-sandbox.dl.alipaydev.com/gateway.do"
	ProductionGatewayURL = "https://openapi.alipay.com/gateway.do"
	ProductionMAPIURL    = "https://mapi.alipay.com/gateway.do"

	FormatJson     = "JSON"
	CharsetUTF8    = "utf-8"
	SignTypeRSA2   = "RSA2"
	ApiVersion     = "1.0"
	CertificateEnd = "-----END CERTIFICATE-----"
)

type TradeStatus string

const (
	TradeWaitBuyerPay TradeStatus = "WAIT_BUYER_PAY" //（交易创建，等待买家付款）
	TradeClosed       TradeStatus = "TRADE_CLOSED"   //（未付款交易超时关闭，或支付完成后全额退款）
	TradeSuccess      TradeStatus = "TRADE_SUCCESS"  //（交易支付成功）
	TradeFinished     TradeStatus = "TRADE_FINISHED" //（交易结束，不可退款）
)

// 销售产品码 product_code
const (
	// FastInstantTradePay PC网站
	FastInstantTradePay string = "FAST_INSTANT_TRADE_PAY"
	// QuickWapWay 手机网站
	QuickWapWay string = "QUICK_WAP_WAY"
	// QuickMsecurityPay app支付
	QuickMsecurityPay string = "QUICK_MSECURITY_PAY"
	// FaceToFacePayment 当面付产品
	FaceToFacePayment string = "FACE_TO_FACE_PAYMENT"
	TransAccountNoPwd string = "TRANS_ACCOUNT_NO_PWD"
)

// VerificationScene 校验场景
type VerificationScene int

const (
	AsyncVerificationScene VerificationScene = 1
	SyncVerificationScene  VerificationScene = 2
)

var ErrTrans = errors.New("xpay:transform error")
var ErrNotContainsSignData = errors.New("xpay:not contains sign data error")
var ErrRequestTimeout = errors.New("xpay: request timeout error")
var ErrRequest = errors.New("xpay: request  error")

// exclude key
const (
	ExcludeKeySign     string = "sign"
	ExcludeKeySignType string = "sign_type"
)

const (
	GrantAuthorizationCode string = "authorization_code"
	GrantRefreshToken      string = "refresh_token"
)

const (
	CreditAuth   string = "CREDIT_AUTH"   // 信用授权场景
	CreditDeduct string = "CREDIT_DEDUCT" // 信用代扣场景
)

const (
	BizCodeFace          string = "FACE"
	BizCodeCertPhoto     string = "CERT_PHOTO"
	BizCodeCertPhotoFace string = "CERT_PHOTO_FACE"
	BizCodeSmartFace     string = "SMART_FACE"
)
