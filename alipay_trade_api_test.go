package alipay

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"testing"
	"time"
)

/**
 * @author: Sam
 * @since:
 * @date: 2023/3/10 17:52
 * @desc:
 */

var client *Client

func init() {
	var err error
	signStrategy := NewCertSignStrategy(OtherAppId, OtherPrivateKey, "appPublicCert.crt", "alipayRootCert.crt", "alipayPublicCert.crt")
	client, err = NewClient(signStrategy, SetClientOptIsProd(true))

	//signStrategy := NewNormalRSA2SignStrategy(AppId,PrivateKey, PublicKey, AlipayPublicKey)
	//client, err = NewAlipayClient(signStrategy)
	if err != nil {
		fmt.Println("初始化失败, 错误信息为", err, client)
		os.Exit(-1)
	}
}

func TestClient_AsyncNotify(t *testing.T) {
	notifyParamMap := make(map[string]string, 0)
	notifyParamMap["app_id"] = "2016091800539057"
	notifyParamMap["auth_app_id"] = "2016091800539057"
	notifyParamMap["buyer_id"] = "2088102176305546"
	notifyParamMap["buyer_pay_amount"] = "1.69"
	notifyParamMap["charset"] = "utf-8"
	notifyParamMap["fund_bill_list"] = "[{\"amount\":\"1.69\",\"fundChannel\":\"ALIPAYACCOUNT\"}]"
	notifyParamMap["gmt_create"] = "2023-03-15 11:24:05"
	notifyParamMap["gmt_payment"] = "2023-03-15 11:27:57"
	notifyParamMap["invoice_amount"] = "1.69"
	notifyParamMap["notify_id"] = "2023031500222112759005540527043634"
	notifyParamMap["notify_time"] = "2023-03-15 11:28:00"
	notifyParamMap["notify_type"] = "trade_status_sync"
	notifyParamMap["out_trade_no"] = "123456789022"
	notifyParamMap["point_amount"] = "0.00"
	notifyParamMap["receipt_amount"] = "1.69"
	notifyParamMap["seller_id"] = "2088102176166347"
	notifyParamMap["sign"] = "F3nlxkcKoCtginCmxscHhDRxTEpxP8jh+TeQCEP87UCC/in+g9JQ9zRVlxRq/819OYykVG+wY2D0vwNpiiiujZaOuGiH2Yl/VKgnG3CSyCZfbsoa02nydc4pV87ZbU9xwYF8CX3RODbTRNWC589uY1OXi1L4YzkPT1AoZVNu+LtZhrVlWEYGOW0Isjd95ZEG9rkdDk4LfcAAW9A6lY9IExFF0S7nRm6UoGRzxI/ViHw3g/gyMBTBb/MAzST3J4PmtB3o6YVkJDbOuY5b9wFi9e2r22xccYqdLbCNYN24Qdhlu+MCkSH5mDiDGgc0jMgOEs5BVzULIIYE+biAs05QXw=="
	notifyParamMap["sign_type"] = "RSA2"
	notifyParamMap["subject"] = "测试title"
	notifyParamMap["total_amount"] = "1.69"
	notifyParamMap["trade_no"] = "2023031522001405540505760265"
	notifyParamMap["trade_status"] = "TRADE_SUCCESS"
	notifyParamMap["version"] = "1.0"
	var err error
	var buff []byte
	if buff, err = json.Marshal(notifyParamMap); err != nil {
		panic(err)
	}
	var notifyParam = new(NotifyReq)
	if err = json.Unmarshal(buff, notifyParam); err != nil {
		panic(err)
	}
	var keyValueList = make([]string, 0, len(notifyParamMap))
	for key, value := range notifyParamMap {
		// 删除sign、sign_type
		if key == "sign" || key == "sign_type" {
			continue
		}
		keyValueList = append(keyValueList, key+"="+value)
	}
	sort.Strings(keyValueList)
	src := strings.Join(keyValueList, "&")
	err = client.VerifySign(AsyncVerificationScene, notifyParam.Sign, []byte(src), notifyParam.AlipayCertSn)
	log.Println("校验参数err", err)
}

func TestClient_TradePagePay(t *testing.T) {
	goods := NewGoodsDetail("1111", "iphone3", 1, 100).SetShowURL("https://ms.bdimg.com/pacific/0/pic/-1225338224_-1800436947.jpg")
	goodsDetail := make([]*GoodsDetail, 0, 1)
	goodsDetail = append(goodsDetail, goods)
	req := NewTradePagePayReq("210122262212", "100.20", "测试title", WithGoodsDetail(goodsDetail))
	req.NotifyUrl = "http://106.14.196.12:8081/order/asyncCallBack"
	req.ReturnUrl = "http://106.14.196.12:8081/syncCallBack"
	result, err := client.TradePagePay(*req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func TestClient_TradeWapPay(t *testing.T) {
	goods := NewGoodsDetail("1111", "iphone3", 1, 100).SetShowURL("https://ms.bdimg.com/pacific/0/pic/-1225338224_-1800436947.jpg")
	goodsDetail := make([]*GoodsDetail, 0, 1)
	goodsDetail = append(goodsDetail, goods)
	req := TradeWapPayReq{}
	req.OutTradeNo = "1211010101221122"
	req.TotalAmount = "1001.00"
	req.Subject = "测试产品"
	req.GoodsDetail = goodsDetail
	req.NotifyUrl = "http://106.14.196.12:8081/asyncCallBack"
	req.ReturnUrl = "http://106.14.196.12:8081/syncCallBack"
	req.ProductCode = QuickWapWay
	result, err := client.TradeWapPay(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func TestClient_TradeAppPay(t *testing.T) {
	goods := NewGoodsDetail("1111", "iphone3", 1, 100).SetShowURL("https://ms.bdimg.com/pacific/0/pic/-1225338224_-1800436947.jpg")
	goodsDetail := make([]*GoodsDetail, 0, 1)
	goodsDetail = append(goodsDetail, goods)
	req := TradeAppPayReq{}
	req.OutTradeNo = "221010101221122"
	req.TotalAmount = "1001.00"
	req.Subject = "测试产品"
	req.GoodsDetail = goodsDetail
	req.NotifyUrl = "http://106.14.196.12:8081/asyncCallBack"
	req.ReturnUrl = "http://106.14.196.12:8081/syncCallBack"
	result, err := client.TradeAppPay(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}
func TestClient_TradeQuery(t *testing.T) {
	req := TradeQueryReq{}
	req.OutTradeNo = "20230321171318"
	req.QueryOptions = []string{"TRADE_SETTLE_INFO"}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Duration(time.Second))
	defer cancelFunc()
	result, err := client.TradeQuery(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)

}
func TestClient_TradeClose(t *testing.T) {
	req := TradeCloseReq{}
	req.OutTradeNo = "21110122262212"

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Duration(time.Second))
	defer cancelFunc()
	result, err := client.TradeClose(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)

}

func TestClient_TradeRefund(t *testing.T) {
	req := TradeRefundReq{}
	req.OutTradeNo = "20230321171318"
	req.OutRequestNo = "202303211713181"
	req.RefundAmount = "0.10"
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Duration(time.Second))
	defer cancelFunc()
	result, err := client.TradeRefund(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)

}

func TestClient_TradeFastPayRefundQuery(t *testing.T) {
	req := TradeFastPayRefundQueryReq{}
	req.OutTradeNo = "trade_no_20170623021124"
	req.OutRequestNo = "111"
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Duration(time.Second))
	defer cancelFunc()
	result, err := client.TradeFastPayRefundQuery(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func TestClient_DataServiceBillDownloadUrlQuery(t *testing.T) {
	req := DataServiceBillDownloadUrlQueryReq{}
	req.BillType = "trade"
	req.BillDate = "2019-01-01"
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Duration(time.Second))
	defer cancelFunc()
	result, err := client.DataServiceBillDownloadUrlQuery(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func TestClient_TradePreCreate(t *testing.T) {
	goods := NewGoodsDetail("1111", "iphone3", 1, 100).SetShowURL("https://ms.bdimg.com/pacific/0/pic/-1225338224_-1800436947.jpg")
	goodsDetail := make([]*GoodsDetail, 0, 1)
	goodsDetail = append(goodsDetail, goods)
	req := TradePreCreateReq{}
	req.OutTradeNo = "2310101012211222"
	req.TotalAmount = "1001.00"
	req.Subject = "测试产品"
	req.GoodsDetail = goodsDetail
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Duration(time.Second))
	defer cancelFunc()
	result, err := client.TradePreCreate(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func TestClient_TradeCancel(t *testing.T) {
	req := TradeCancelReq{}
	req.OutTradeNo = "trade_no_20170623021124"
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Duration(time.Second))
	defer cancelFunc()
	result, err := client.TradeCancel(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}
func TestClient_DataBillBalanceQuery(t *testing.T) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Duration(time.Second))
	defer cancelFunc()
	req := DataBillBalanceQueryReq{}
	result, err := client.DataBillBalanceQuery(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func TestClient_DataBillBailQuery(t *testing.T) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Duration(time.Second))
	defer cancelFunc()
	req := DataBillBailQueryReq{
		StartTime: "2023-03-01 00:00:00",
		EndTime:   "2023-03-010 00:00:00",
		BailType:  "TMALL_BAIL",
	}
	result, err := client.DataBillBailQuery(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func TestClient_TradeCreate(t *testing.T) {
	req := TradeCreateReq{}
	req.OutTradeNo = "231010101221122"
	req.TotalAmount = "101.00"
	req.Subject = "测试产品"
	req.ProductCode = FaceToFacePayment
	req.BuyerId = "2088102146225135"
	req.NotifyUrl = "http://106.14.196.12:8081/asyncCallBack"
	ctx, cancelFunc := context.WithTimeout(context.Background(), 500*time.Duration(time.Second))
	defer cancelFunc()
	result, err := client.TradeCreate(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func TestClient_TradePay(t *testing.T) {
	goods := NewGoodsDetail("1111", "iphone3", 1, 100).SetShowURL("https://ms.bdimg.com/pacific/0/pic/-1225338224_-1800436947.jpg")
	goodsDetail := make([]*GoodsDetail, 0, 1)
	goodsDetail = append(goodsDetail, goods)
	req := &TradePayReq{}
	req.OutTradeNo = "20230322111111123"
	req.TotalAmount = "56058987.00"
	req.Subject = "100个iPhone14 pro max "
	req.Scene = "bar_code"
	req.AuthCode = "285516572327851289"
	req.GoodsDetail = goodsDetail
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Duration(time.Second))
	defer cancelFunc()
	result, err := client.TradePay(ctx, *req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func TestClient_TradeOrderInfoSync(t *testing.T) {
	req := &TradeOrderInfoSyncReq{}
	req.TradeNo = "2023031522001405540505760265"
	req.OutRequestNo = "20230322111111123"
	req.BizType = CreditAuth
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Duration(time.Second))
	defer cancelFunc()
	result, err := client.TradeOrderInfoSync(ctx, *req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func TestClient_FundAccountQuery(t *testing.T) {
	req := &FundAccountQueryReq{}
	req.AlipayUserId = "2088102176166347"
	req.AccountType = "ACCTRANS_ACCOUNT"
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Duration(time.Second))
	defer cancelFunc()
	result, err := client.FundAccountQuery(ctx, *req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func TestClient_SystemOauthToken(t *testing.T) {
	req := OauthTokenReq{GrantType: GrantAuthorizationCode, Code: "111"}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 500*time.Duration(time.Second))
	defer cancelFunc()
	result, err := client.SystemOauthToken(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func TestClient_UserInfoAuth(t *testing.T) {
	req := &UserInfoAuthReq{}
	req.Scopes = []string{"auth_user"}
	req.State = "init"
	req.ReturnUrl = "http://106.14.196.12:8081/"
	result, err := client.UserInfoAuth(*req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func TestClient_UserInfoShare(t *testing.T) {
	req := UserInfoShareReq{AuthToken: "ERITJKEIJKJHKKKKKKKHJEREEEEEEEEEEE111111"}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Duration(time.Second))
	defer cancelFunc()
	result, err := client.UserInfoShare(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func TestClient_OpenAuthTokenApp(t *testing.T) {
	req := OpenAuthTokenAppReq{GrantType: GrantAuthorizationCode, Code: "111"}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Duration(time.Second))
	defer cancelFunc()
	result, err := client.OpenAuthTokenApp(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func TestClient_UserCertifyOpenInitialize(t *testing.T) {
	req := UserCertifyOpenInitializeReq{}
	req.OuterOrderNo = "ZGYD201809132323000001234"
	req.BizCode = BizCodeFace
	merchantConfig := OpenCertifyMerchantConfig{
		ReturnUrl: "http://www.baidu.com",
	}
	req.MerchantConfig = merchantConfig
	identityParam := OpenCertifyIdentityParam{
		IdentityType: "CERT_INFO",
		CertType:     "IDENTITY_CARD",
		CertName:     "张三",
		CertNo:       "342423199011233456",
	}
	req.IdentityParam = identityParam
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Duration(time.Second))
	defer cancelFunc()
	result, err := client.UserCertifyOpenInitialize(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func TestClient_UserCertifyOpenQuery(t *testing.T) {
	req := UserCertifyOpenQueryReq{}
	req.CertifyId = "OC201809253000000393900404029253"
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Duration(time.Second))
	defer cancelFunc()
	result, err := client.UserCertifyOpenQuery(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func TestClient_UserCertifyOpenCertify(t *testing.T) {
	req := UserCertifyOpenCertifyReq{}
	req.CertifyId = "OC201809253000000393900404029253"
	req.ReturnUrl = "http://www.baidu.com"
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Duration(time.Second))
	defer cancelFunc()
	result, err := client.UserCertifyOpenCertify(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func TestClient_FundTransToAccountTransfer(t *testing.T) {
	req := FundTransToAccountTransferReq{}
	req.OutBizNo = "OC201809253000000393900404029253"
	req.PayeeType = "ALIPAY_LOGONID"
	req.PayeeAccount = "13951604344"
	req.Amount = "100.00"
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Duration(time.Second))
	defer cancelFunc()
	result, err := client.FundTransToAccountTransfer(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func TestClient_FundTransOrderQuery(t *testing.T) {
	req := FundTransOrderQueryReq{}
	req.OutBizNo = "OC201809253000000393900404029253"
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Duration(time.Second))
	defer cancelFunc()
	result, err := client.FundTransOrderQuery(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func TestClient_FundTransUniTransfer(t *testing.T) {
	req := FundTransUniTransferReq{}
	req.TransAmount = "10.20"
	req.ProductCode = TransAccountNoPwd
	req.OutBizNo = "OC201809253000000393900404029253"
	req.BizScene = "DIRECT_TRANSFER"
	req.OrderTitle = "测试"
	participant := Participant{
		Identity:     "13951604344",
		IdentityType: "ALIPAY_LOGON_ID",
		Name:         "sam",
	}
	req.PayeeInfo = participant
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Duration(time.Second))
	defer cancelFunc()
	result, err := client.FundTransUniTransfer(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func TestClient_FundTransCommonQuery(t *testing.T) {
	req := FundTransCommonQueryReq{}
	req.ProductCode = TransAccountNoPwd
	req.OutBizNo = "OC201809253000000393900404029253"
	req.BizScene = "DIRECT_TRANSFER"
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Duration(time.Second))
	defer cancelFunc()
	result, err := client.FundTransCommonQuery(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func TestClient_CommerceCityFacilitatorVoucherGenerate(t *testing.T) {
	req := CommerceCityFacilitatorVoucherGenerateReq{}
	req.CityCode = "440300"
	req.TradeNo = "OC201809253000000393900404029253"
	req.TradeFee = "10.00"
	req.TicketNum = "2"
	req.TicketType = "oneway"
	req.TicketPrice = "5.00"
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Duration(time.Second))
	defer cancelFunc()
	result, err := client.CommerceCityFacilitatorVoucherGenerate(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}
func TestClient_CommerceCityFacilitatorVoucherRefund(t *testing.T) {
	req := CommerceCityFacilitatorVoucherRefundReq{}
	req.CityCode = "440300"
	req.TradeNo = "OC201809253000000393900404029253"

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Duration(time.Second))
	defer cancelFunc()
	result, err := client.CommerceCityFacilitatorVoucherRefund(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func TestClient_CommerceCityFacilitatorStationQuery(t *testing.T) {
	req := CommerceCityFacilitatorStationQueryReq{}
	req.CityCode = "440300"

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Duration(time.Second))
	defer cancelFunc()
	result, err := client.CommerceCityFacilitatorStationQuery(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func TestClient_CommerceCityFacilitatorVoucherBatchQuery(t *testing.T) {
	req := CommerceCityFacilitatorVoucherBatchQueryReq{}
	req.CityCode = "440300"
	req.TradeNos = []string{"0123456789"}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Duration(time.Second))
	defer cancelFunc()
	result, err := client.CommerceCityFacilitatorVoucherBatchQuery(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func TestName(t *testing.T) {
	//type Options struct {
	//	Query   string `url:"q"`
	//	ShowAll bool   `url:"all"`
	//	Page    int    `url:"-"`
	//}
	//opt := Options{Query: "foo", ShowAll: true, Page: 10}
	//v, _ := query.Values(opt)
	//fmt.Print(v.Encode())
	//amountStr := fmt.Sprintf("%d", 100)
	//length := len(amountStr)
	//start := amountStr[:length-2]
	//end := amountStr[length-2:]
	//fmt.Println(fmt.Sprintf("%s.%s", start, end))
	//fmt.Println(strings.Trim("---a--Hello, Gophers-----", "-a"))
	//fmt.Println(strings.TrimSpace("  77d8d8d  "))
}
