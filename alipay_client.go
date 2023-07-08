package alipay

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

/**
 * @author: Sam
 * @since: 1.0.0
 * @date: 2023/3/15 17:46
 * @desc: alipay client
 */

type ClientOptsFunc []ClientOptFunc

func (opts ClientOptsFunc) apply(client *Client) {
	for _, optFunc := range opts {
		optFunc(client)
	}
}

type ClientOptFunc func(*Client)

// SetClientOptLocation  setup Location
func SetClientOptLocation(location *time.Location) ClientOptFunc {
	return func(client *Client) {
		client.location = location
	}
}

// SetClientOptHttpClient  setup HttpClient
func SetClientOptHttpClient(httpClient *http.Client) ClientOptFunc {
	return func(client *Client) {
		client.httpClient = httpClient
	}
}

// SetClientOptIsProd setup isProd
func SetClientOptIsProd(isProd bool) ClientOptFunc {
	return func(client *Client) {
		client.isProd = isProd
	}
}

func SetServerUrl(serverUrl string) ClientOptFunc {
	return func(client *Client) {
		client.serverUrl = serverUrl
	}
}

type Client struct {
	mu        sync.Mutex
	serverUrl string
	// 是否时生产环境
	isProd bool
	// 时区
	location *time.Location
	// http 请求客户端
	httpClient *http.Client
	SignVerifier
	RequestObjectBuilder
}

func NewClient(signVerifier SignVerifier, optsFunc ...ClientOptFunc) (*Client, error) {
	client := &Client{
		serverUrl:    SandboxGatewayURL,
		httpClient:   http.DefaultClient,
		location:     time.Local,
		SignVerifier: signVerifier,
		RequestObjectBuilder: &RequestAliPayObjectBuilder{
			location: time.Local,
		},
	}
	ClientOptsFunc(optsFunc).apply(client)
	if client.isProd {
		client.serverUrl = ProductionGatewayURL
	}
	return client, nil
}

// TradePagePay alipay.trade.page.pay(统一收单下单并支付页面接口) https://opendocs.alipay.com/open/028r8t
func (r *Client) TradePagePay(req TradePagePayReq) (*url.URL, error) {
	var err error
	if err = req.DoValidate(); err != nil {
		return nil, err
	}
	var commonReqParam *CommonReqParam
	if commonReqParam, err = r.buildRequestObject(&req, WithNotifyUrl(req.NotifyUrl), WithReturnUrl(req.ReturnUrl)); err != nil {
		return nil, err
	}
	r.SetSignContent(commonReqParam)
	var encode string
	if encode, err = r.Encode(); err != nil {
		return nil, err
	}
	return url.Parse(r.serverUrl + "?" + encode)
}

// TradeQuery alipay.trade.query(统一收单交易查询) https://opendocs.alipay.com/open/028woa?scene=common
func (r *Client) TradeQuery(ctx context.Context, req TradeQueryReq) (*TradeQueryRes, error) {
	res := new(TradeQueryRes)
	err := r.DoRequest(ctx, &req, res)
	return res, err
}

// TradeClose alipay.trade.close(统一收单交易关闭接口) https://opendocs.alipay.com/open/028wob
func (r *Client) TradeClose(ctx context.Context, req TradeCloseReq) (*TradeCloseRes, error) {
	res := new(TradeCloseRes)
	err := r.DoRequest(ctx, &req, res, WithNotifyUrl(req.NotifyUrl))
	return res, err
}

// TradeRefund alipay.trade.refund(统一收单交易退款接口) https://opendocs.alipay.com/open/028sm9
func (r *Client) TradeRefund(ctx context.Context, req TradeRefundReq) (*TradeRefundRes, error) {
	res := new(TradeRefundRes)
	err := r.DoRequest(ctx, &req, res)
	return res, err
}

// TradeFastPayRefundQuery alipay.trade.fastpay.refund.query(统一收单交易退款查询) https://opendocs.alipay.com/open/028sma
func (r *Client) TradeFastPayRefundQuery(ctx context.Context, req TradeFastPayRefundQueryReq) (*TradeFastPayRefundQueryRes, error) {
	res := new(TradeFastPayRefundQueryRes)
	err := r.DoRequest(ctx, &req, res)
	return res, err
}

// DataServiceBillDownloadUrlQuery alipay.data.dataservice.bill.downloadurl.query(查询对账单下载地址) https://opendocs.alipay.com/open/028woc
func (r *Client) DataServiceBillDownloadUrlQuery(ctx context.Context, req DataServiceBillDownloadUrlQueryReq) (*DataServiceBillDownloadUrlQueryRes, error) {
	res := new(DataServiceBillDownloadUrlQueryRes)
	err := r.DoRequest(ctx, &req, res)
	return res, err
}

// TradeWapPay alipay.trade.wap.pay(手机网站支付接口2.0) https://opendocs.alipay.com/open/02ivbs?scene=21&ref=api
func (r *Client) TradeWapPay(req TradeWapPayReq) (*url.URL, error) {
	var commonReqParam *CommonReqParam
	var err error
	if err = req.DoValidate(); err != nil {
		return nil, err
	}
	if commonReqParam, err = r.buildRequestObject(&req, WithNotifyUrl(req.NotifyUrl), WithReturnUrl(req.ReturnUrl)); err != nil {
		return nil, err
	}
	r.SetSignContent(commonReqParam)
	var encode string
	if encode, err = r.Encode(); err != nil {
		return nil, err
	}
	return url.Parse(r.serverUrl + "?" + encode)
}

// TradeAppPay alipay.trade.app.pay(app支付接口2.0) https://opendocs.alipay.com/open/02e7gq?ref=api&scene=20
func (r *Client) TradeAppPay(req TradeAppPayReq) (*url.URL, error) {
	var commonReqParam *CommonReqParam
	var err error
	if err = req.DoValidate(); err != nil {
		return nil, err
	}
	if commonReqParam, err = r.buildRequestObject(&req, WithNotifyUrl(req.NotifyUrl), WithReturnUrl(req.ReturnUrl)); err != nil {
		return nil, err
	}
	r.SetSignContent(commonReqParam)
	var encode string
	if encode, err = r.Encode(); err != nil {
		return nil, err
	}
	return url.Parse(r.serverUrl + "?" + encode)
}

// TradePreCreate https://opendocs.alipay.com/open/02ekfg?scene=19 alipay.trade.precreate(统一收单线下交易预创建)
func (r *Client) TradePreCreate(ctx context.Context, req TradePreCreateReq) (*TradePreCreateRes, error) {
	res := new(TradePreCreateRes)
	err := r.DoRequest(ctx, &req, res)
	return res, err
}

// TradeCancel alipay.trade.cancel(统一收单交易撤销接口) https://opendocs.alipay.com/open/02ekfi
func (r *Client) TradeCancel(ctx context.Context, req TradeCancelReq) (*TradeCancelRes, error) {
	res := new(TradeCancelRes)
	err := r.DoRequest(ctx, &req, res)
	return res, err
}

// SystemOauthToken alipay.system.oauth.token(换取授权访问令牌) https://opendocs.alipay.com/open/02ahjv
func (r *Client) SystemOauthToken(ctx context.Context, req OauthTokenReq) (*OauthTokenRes, error) {
	res := new(OauthTokenRes)
	err := r.DoRequest(ctx, &req, res)
	return res, err
}

// UserInfoShare alipay.user.info.share(支付宝会员授权信息查询接口) https://opendocs.alipay.com/open/02ailg
func (r *Client) UserInfoShare(ctx context.Context, req UserInfoShareReq) (*UserInfoShareRes, error) {
	res := new(UserInfoShareRes)
	err := r.DoRequest(ctx, &req, res)
	return res, err
}

// DataBillBalanceQuery alipay.data.bill.balance.query(支付宝商家账户当前余额查询) https://opendocs.alipay.com/open/02awe3
func (r *Client) DataBillBalanceQuery(ctx context.Context, req DataBillBalanceQueryReq) (*DataBillBalanceQueryRes, error) {
	res := new(DataBillBalanceQueryRes)
	err := r.DoRequest(ctx, &req, res)
	return res, err
}

// DataBillBailQuery alipay.data.bill.bail.query(支付宝商家账户保证金查询) https://opendocs.alipay.com/open/02awe2
func (r *Client) DataBillBailQuery(ctx context.Context, req DataBillBailQueryReq) (*DataBillBailQueryRes, error) {
	res := new(DataBillBailQueryRes)
	err := r.DoRequest(ctx, &req, res)
	return res, err
}

// TradeCreate alipay.trade.create(统一收单交易创建接口) https://opendocs.alipay.com/mini/03l5wn
func (r *Client) TradeCreate(ctx context.Context, req TradeCreateReq) (*TradeCreateRes, error) {
	res := new(TradeCreateRes)
	err := r.DoRequest(ctx, &req, res, WithNotifyUrl(req.NotifyUrl))
	return res, err
}

// TradePay alipay.trade.pay(统一收单交易支付接口) https://opendocs.alipay.com/open/02ekfp?scene=32
func (r *Client) TradePay(ctx context.Context, req TradePayReq) (*TradePayRes, error) {
	res := new(TradePayRes)
	err := r.DoRequest(ctx, &req, res, WithNotifyUrl(req.NotifyUrl))
	return res, err
}

// TradeOrderInfoSync alipay.trade.orderinfo.sync(支付宝订单信息同步接口) https://opendocs.alipay.com/open/02cnou
func (r *Client) TradeOrderInfoSync(ctx context.Context, req TradeOrderInfoSyncReq) (*TradeOrderInfoSyncRes, error) {
	res := new(TradeOrderInfoSyncRes)
	err := r.DoRequest(ctx, &req, res)
	return res, err
}

// FundAccountQuery alipay.fund.account.query(支付宝资金账户资产查询接口) https://opendocs.alipay.com/open/02byuq?scene=c76aa8f1c54e4b8b8ffecfafc4d3c31c
func (r *Client) FundAccountQuery(ctx context.Context, req FundAccountQueryReq) (*FundAccountQueryRes, error) {
	res := new(FundAccountQueryRes)
	err := r.DoRequest(ctx, &req, res)
	return res, err
}

// UserInfoAuth alipay.user.info.auth(用户登录授权) https://opendocs.alipay.com/open/02aile
func (r *Client) UserInfoAuth(req UserInfoAuthReq) (*url.URL, error) {
	var err error
	if err = req.DoValidate(); err != nil {
		return nil, err
	}
	var commonReqParam *CommonReqParam
	if commonReqParam, err = r.buildRequestObject(&req, WithReturnUrl(req.ReturnUrl)); err != nil {
		return nil, err
	}
	r.SetSignContent(commonReqParam)
	var encode string
	if encode, err = r.Encode(); err != nil {
		return nil, err
	}
	return url.Parse(r.serverUrl + "?" + encode)
}

// OpenAuthTokenApp alipay.open.auth.token.app(换取应用授权令牌) https://opendocs.alipay.com/isv/04h3uf
func (r *Client) OpenAuthTokenApp(ctx context.Context, req OpenAuthTokenAppReq) (*OpenAuthTokenAppRes, error) {
	res := new(OpenAuthTokenAppRes)
	err := r.DoRequest(ctx, &req, res)
	return res, err
}

// UserCertifyOpenInitialize alipay.user.certify.open.initialize(身份认证初始化服务) https://opendocs.alipay.com/open/02ahjy
func (r *Client) UserCertifyOpenInitialize(ctx context.Context, req UserCertifyOpenInitializeReq) (*UserCertifyOpenInitializeRes, error) {
	res := new(UserCertifyOpenInitializeRes)
	err := r.DoRequest(ctx, &req, res)
	return res, err
}

// UserCertifyOpenQuery alipay.user.certify.open.query(身份认证记录查询) https://opendocs.alipay.com/open/02ahjw
func (r *Client) UserCertifyOpenQuery(ctx context.Context, req UserCertifyOpenQueryReq) (*UserCertifyOpenQueryRes, error) {
	res := new(UserCertifyOpenQueryRes)
	err := r.DoRequest(ctx, &req, res)
	return res, err
}

// UserCertifyOpenCertify alipay.user.certify.open.certify(身份认证开始认证) https://opendocs.alipay.com/open/02ahk0
func (r *Client) UserCertifyOpenCertify(ctx context.Context, req UserCertifyOpenCertifyReq) (*url.URL, error) {
	var err error
	if err = req.DoValidate(); err != nil {
		return nil, err
	}
	var commonReqParam *CommonReqParam
	if commonReqParam, err = r.buildRequestObject(&req, WithReturnUrl(req.ReturnUrl)); err != nil {
		return nil, err
	}
	r.SetSignContent(commonReqParam)
	var encode string
	if encode, err = r.Encode(); err != nil {
		return nil, err
	}
	return url.Parse(r.serverUrl + "?" + encode)
}

// Deprecated: 此接口已过时，推荐使用接口 FundTransUniTransfer()，相关升级指南 https://opendocs.alipay.com/open/00ou7f
// FundTransToAccountTransfer alipay.fund.trans.toaccount.transfer(单笔转账到支付宝账户接口) https://opendocs.alipay.com/apis/00fka9
func (r *Client) FundTransToAccountTransfer(ctx context.Context, req FundTransToAccountTransferReq) (*FundTransToAccountTransferRes, error) {
	res := new(FundTransToAccountTransferRes)
	err := r.DoRequest(ctx, &req, res)
	return res, err
}

// FundTransOrderQuery alipay.fund.trans.order.query(查询转账订单接口) https://opendocs.alipay.com/apis/api_28/alipay.fund.trans.order.query
func (r *Client) FundTransOrderQuery(ctx context.Context, req FundTransOrderQueryReq) (*FundTransOrderQueryRes, error) {
	res := new(FundTransOrderQueryRes)
	err := r.DoRequest(ctx, &req, res)
	return res, err
}

// FundTransUniTransfer alipay.fund.trans.uni.transfer(单笔转账接口) https://opendocs.alipay.com/open/02byuo
func (r *Client) FundTransUniTransfer(ctx context.Context, req FundTransUniTransferReq) (*FundTransUniTransferRes, error) {
	res := new(FundTransUniTransferRes)
	err := r.DoRequest(ctx, &req, res)
	return res, err
}

// FundTransCommonQuery alipay.fund.trans.common.query(转账业务单据查询接口) https://opendocs.alipay.com/open/02byup?scene=f9fece54d41f49cbbd00dc73655a01a4
func (r *Client) FundTransCommonQuery(ctx context.Context, req FundTransCommonQueryReq) (*FundTransCommonQueryRes, error) {
	res := new(FundTransCommonQueryRes)
	err := r.DoRequest(ctx, &req, res)
	return res, err
}

// CommerceCityFacilitatorVoucherGenerate alipay.commerce.cityfacilitator.voucher.generate(地铁购票核销码发码) https://opendocs.alipay.com/open/02ars7
func (r *Client) CommerceCityFacilitatorVoucherGenerate(ctx context.Context, req CommerceCityFacilitatorVoucherGenerateReq) (*CommerceCityFacilitatorVoucherGenerateRes, error) {
	res := new(CommerceCityFacilitatorVoucherGenerateRes)
	err := r.DoRequest(ctx, &req, res)
	return res, err
}

// CommerceCityFacilitatorVoucherRefund alipay.commerce.cityfacilitator.voucher.refund(地铁购票发码退款) https://opendocs.alipay.com/open/02ars8
func (r *Client) CommerceCityFacilitatorVoucherRefund(ctx context.Context, req CommerceCityFacilitatorVoucherRefundReq) (*CommerceCityFacilitatorVoucherRefundRes, error) {
	res := new(CommerceCityFacilitatorVoucherRefundRes)
	err := r.DoRequest(ctx, &req, res)
	return res, err
}

// CommerceCityFacilitatorStationQuery alipay.commerce.cityfacilitator.station.query(地铁购票站点数据查询) https://opendocs.alipay.com/open/02ars9
func (r *Client) CommerceCityFacilitatorStationQuery(ctx context.Context, req CommerceCityFacilitatorStationQueryReq) (*CommerceCityFacilitatorStationQueryRes, error) {
	res := new(CommerceCityFacilitatorStationQueryRes)
	err := r.DoRequest(ctx, &req, res)
	return res, err
}

// CommerceCityFacilitatorVoucherBatchQuery alipay.commerce.cityfacilitator.voucher.batchquery(地铁购票订单批量查询) https://opendocs.alipay.com/open/02aqvy
func (r *Client) CommerceCityFacilitatorVoucherBatchQuery(ctx context.Context, req CommerceCityFacilitatorVoucherBatchQueryReq) (*CommerceCityFacilitatorVoucherBatchQueryRes, error) {
	res := new(CommerceCityFacilitatorVoucherBatchQueryRes)
	err := r.DoRequest(ctx, &req, res)
	return res, err
}

// DoRequest 发送请求
func (r *Client) DoRequest(ctx context.Context, req IAliPayRequest, responseParam ResponseSigner, opts ...commonParamOpt) error {
	if err := req.DoValidate(); err != nil {
		return err
	}
	buff, err := r.doRequest(ctx, req, opts...)
	if err != nil {
		return err
	}
	resBuff := responseBuff(buff)
	if resBuff.IsHtmlError() {
		return ErrRequest
	}
	if err = json.Unmarshal(buff, responseParam); err != nil {
		return err
	}
	return r.VerifySign(SyncVerificationScene, responseParam.GetSign(), buff, responseParam.GetAlipayCertSn())
}

// DoRequest 发送请求
//func (r *Client) DoRequest(ctx context.Context, req IAliPayRequest, opts ...commonParamOpt) ([]byte, interface{}, error) {
//	buff, err := r.doRequest(ctx, req, opts...)
//	if err != nil {
//		return nil, nil, err
//	}
//	responseParam := req.ResponseParam()
//	resBuff := responseBuff(buff)
//	if resBuff.IsHtmlError() {
//		return nil, nil, ErrRequest
//	}
//	if resBuff.IsErrorResponse() && !req.IsNormalErrResponse() {
//		errRes := new(ErrorResponse)
//		if err = json.Unmarshal(buff, errRes); err != nil {
//			return []byte{}, nil, err
//		}
//		responseParam.SetCode(errRes.Code)
//		responseParam.SetMsg(errRes.Msg)
//		responseParam.SetSubCode(errRes.SubCode)
//		responseParam.SetSubMsg(errRes.SubMsg)
//		responseParam.SetAlipayCertSn(errRes.AlipayCertSn)
//		responseParam.SetSign(errRes.Sign)
//		return buff, responseParam, err
//	}
//	if err = json.Unmarshal(buff, responseParam); err != nil {
//		return nil, nil, err
//	}
//	return buff, responseParam, err
//}

// 具体请求
func (r *Client) doRequest(ctx context.Context, req IAliPayRequest, opts ...commonParamOpt) ([]byte, error) {
	deadline, ok := ctx.Deadline()
	if ok && time.Now().After(deadline) {
		return nil, ErrRequestTimeout
	}
	var err error
	var newRequest *http.Request
	var commonReqParam *CommonReqParam
	if commonReqParam, err = r.buildRequestObject(req, opts...); err != nil {
		return nil, err
	}
	r.SetSignContent(commonReqParam)
	var encode string
	if encode, err = r.Encode(); err != nil {
		return nil, err
	}
	if newRequest, err = http.NewRequestWithContext(ctx, req.RequestHttpMethod(), r.serverUrl, strings.NewReader(encode)); err != nil {
		return nil, err
	}
	newRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	var response *http.Response
	if response, err = r.httpClient.Do(newRequest); err != nil {
		return nil, err
	}
	defer func() {
		if response == nil {
			return
		}
		response.Body.Close()
	}()
	return io.ReadAll(response.Body)
}
