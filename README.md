
#### XPay
支付宝Golang版本的SDK


[![Github](https://img.shields.io/github/followers/try-labs?label=Follow&style=social)](https://github.com/try-labs)
[![Golang](https://img.shields.io/badge/golang-1.16-brightgreen.svg)](https://golang.google.cn)


#### Golang版本要求
由于Golang对于泛型方法暂未支持，遂SDK中的实现并没有使用泛型。当前SDK对Golang版本没有特殊要求，
但防止一些意外问题发生，推荐最低使用Golang1.16版本。支持设置

#### 安装
```shell
go get -u github.com/try-labs/xpay
```

#### SDK和支付宝接口对应规则
例如支付宝的手机网站支付接口名称为``alipay.trade.wap.pay`` ，则SDK中的实现方法名为 ``TradeWapPay()``，对应的请求结构体为``TradeWapPayReq``，支付宝响应结构体为``TradeWapPayRes``


#### 使用
支持支付宝的公钥模式和证书模式,两种模式能同时存在，只能选择其中一种。所有请求返回结果时已经验证过签名，无需再次实现验签。

 - 使用公钥初始化 ``alipay.Client`` 客户端
```Golang

var err error
var client *alipay.Client
signStrategy := NewNormalRSA2SignStrategy(AppId,PrivateKey, PublicKey, AlipayPublicKey)
client, err = NewAlipayClient(signStrategy)
if err != nil {
  fmt.Println("初始化失败, 错误信息为", err, client)
}
```

 - 使用证书初始化 ``alipay.Client`` 客户端
```Golang

var err error
var client *alipay.Client
signStrategy := NewCertSignStrategy(OtherAppId, OtherPrivateKey, "appPublicCert.crt", "alipayRootCert.crt", "alipayPublicCert.crt")
client, err = NewClient(signStrategy)
  if err != nil {
      fmt.Println("初始化失败, 错误信息为", err, client)
  }

```
- 调用具体的接口
 例如调用alipay.trade.page.pay(统一收单下单并支付页面接口)，按照规则SDK对应的方法为 ``TradePagePay()``,
 ```Golang
goods := NewGoodsDetail("1111", "iphone3", 1, 100).SetShowURL("https://ms.bdimg.com/pacific/0/pic/-1225338224_-1800436947.jpg")
goodsDetail := make([]*GoodsDetail, 0, 1)
goodsDetail = append(goodsDetail, goods)
req := NewTradePagePayReq("210122262212", "100.20", "测试title", WithGoodsDetail(goodsDetail))
req.NotifyUrl = "http://xxxx/order/asyncCallBack"
req.ReturnUrl = "http://xxxx/order/syncCallBack"
result, err := client.TradePagePay(*req)
if err != nil {
    fmt.Println(err)
}
fmt.Println(result)
```

#### 接口列表
- [x]  接口前有此标志代表接口已被实现

具体接口情况如下：

- [x] 手机网站支付接口

  alipay.trade.wap.pay - TradeWapPay()

- [x] 电脑网站支付

  alipay.trade.page.pay - TradePagePay()

- [x] 统一收单线下交易查询

  alipay.trade.query - TradeQuery()

- [x] 统一收单交易支付接口

  alipay.trade.pay - TradePay()

- [x] 统一收单交易创建接口

  alipay.trade.create - TradeCreate()

- [x] 统一收单线下交易预创建

  alipay.trade.precreate - TradePreCreate()

- [x] 统一收单交易撤销接口

  alipay.trade.cancel - TradeCancel()

- [x] 统一收单交易关闭接口

  alipay.trade.close - TradeClose()

- [x] 统一收单交易退款接口

  alipay.trade.refund - TradeRefund()

- [x] App支付接口

  alipay.trade.app.pay - TradeAppPay()

- [x] 统一收单交易退款查询

  alipay.trade.fastpay.refund.query - TradeFastPayRefundQuery()

- [x] 支付宝订单信息同步接口

  alipay.trade.orderinfo.sync - TradeOrderInfoSync()

- [x] 查询对账单下载地址

  alipay.data.dataservice.bill.downloadurl.query - DataServiceBillDownloadUrlQuery()

- [x] 支付宝商家账户当前余额查询

  alipay.data.bill.balance.query - DataBillBalanceQuery()

- [x] 支付宝商家账户保证金查询

  alipay.data.bill.bail.query - DataBillBailQuery()

- 申请电子回单(incubating)
  alipay.data.bill.ereceipt.apply - DataBillEreceiptApply()

- 查询电子回单状态(incubating)

alipay.data.bill.ereceipt.query - DataBillEreceiptQuery()

##### 支付宝身份验证
- [x] 身份认证初始化服务

  alipay.user.certify.open.initialize - UserCertifyOpenInitialize()

- [x] 身份认证开始认证

  alipay.user.certify.open.certify - UserCertifyOpenCertify()

- [x] 身份认证记录查询

  alipay.user.certify.open.query - UserCertifyOpenQuery()

##### 会员产品
- [x] 换取授权访问令牌

  alipay.system.oauth.token - SystemOauthToken()

- [x] 支付宝会员授权信息查询

  alipay.user.info.share - UserInfoShare()

- [x]  用户登录授权
  alipay.user.info.auth - UserInfoAuth()
  https://opendocs.alipay.com/open/02aile

##### 周期扣款
- 支付宝个人协议页面签约

  alipay.user.agreement.page.sign - UserAgreementPageSign()

- 支付宝个人代扣协议查询

  alipay.user.agreement.query - UserAgreementQuery()

- 支付宝个人代扣协议解约

  alipay.user.agreement.unsign - UserAgreementUnsign()

##### 转账到支付宝账户

- [x] 单笔转账接口

  alipay.fund.trans.uni.transfer - FundTransUniTransfer()

>  alipay.fund.trans.uni.transfer为新转账接口

- [x] 转账业务单据查询接口

  alipay.fund.trans.common.query - FundTransCommonQuery()

- [x] 支付宝资金账户资产查询接口

  alipay.fund.account.query - FundAccountQuery()

- [x] 单笔转账到支付宝账户接口

  alipay.fund.trans.toaccount.transfer - FundTransToAccountTransfer()

> alipay.fund.trans.toaccount.transfer为老转账接口，新查询接口 alipay.fund.trans.common.query 可以查询老转账接口生成的转账单据， 
> 老查询接口 alipay.fund.trans.order.query 无法查询到新转账接口生成的转账单据。

- [x] 查询转账订单接口

  alipay.fund.trans.order.query - FundTransOrderQuery()
> 此接口为老查询接口

##### 地铁线上购票

- [x] 地铁购票核销码发码 
  alipay.commerce.cityfacilitator.voucher.generate - CommerceCityFacilitatorVoucherGenerate()
  https://opendocs.alipay.com/open/02ars7

- [x] 地铁购票发码退款
  alipay.commerce.cityfacilitator.voucher.refund -  CommerceCityFacilitatorVoucherRefund()
  https://opendocs.alipay.com/open/02ars8

- [x] 地铁购票站点数据查询
  alipay.commerce.cityfacilitator.station.query  - CommerceCityFacilitatorStationQuery()
  https://opendocs.alipay.com/open/02ars9

- [x] 地铁购票订单批量查询
  alipay.commerce.cityfacilitator.voucher.batchquery -  CommerceCityFacilitatorVoucherBatchQuery()
  https://opendocs.alipay.com/open/02aqvy

#### 沙箱环境

支付宝开放平台为每一个应用提供了沙箱环境，供开发人员开发测试使用。
[沙箱地址](https://open.alipay.com/develop/sandbox/app)
沙箱环境是独立的，每一个应用都会有一个商家账号和买家账号。

#### 应用信息配置

参考[官网文档](https://docs.open.alipay.com/200/105894) 进行应用的配置。

本 SDK 中的签名方法默认为 RSA2，采用支付宝提供的 [RSA签名&验签工具](https://docs.open.alipay.com/291/105971) 生成秘钥时，秘钥的格式必须为 PKCS1，秘钥长度推荐 2048。所以在支付宝管理后台请注意配置 RSA2(SHA256)密钥。

生成秘钥对之后，将公钥提供给支付宝（通过支付宝后台上传）对我们请求的数据进行签名验证，我们的代码中将使用私钥对请求数据签名。

请参考 [如何生成 RSA 密钥](https://docs.open.alipay.com/291/105971)。


#### 同步通知
请求接口时需要设置``return_url`` 参数，支付成功之后支付宝将跳转到该地址并通过Get方式携带通知参数，相关文档如下：
[https://opendocs.alipay.com/support/01raw3?ant_source=opendoc_recommend](https://opendocs.alipay.com/support/01raw3?ant_source=opendoc_recommend)

对应的方法为
 ```Golang
 client.SyncNotify()
 ```
#### 异步通知
对于支付产生的交易，支付宝会根据原始支付 API 中传入的异步通知地址 notify_url，通过 POST 请求的形式将支付结果作为参数通知到商家系统。
[https://opendocs.alipay.com/support/01raw4](https://opendocs.alipay.com/support/01raw4)

对应的方法为
 ```Golang
 client.AsyncNotify()
 ```
#### 两种通知的区别
同步通知与异步通知的区别
return_url用于接收同步通知，notify_url用于接收异步通知。
主要区别如下：
- 接收方法不同
同步通知是在同步返回的页面使用get方式来获取的，异步通知是使用post方式发送到异步地址上的，只能使用post方式来接收。

- 返回的参数不同
同步通知返回的参数比异步通知返回的参数是要少一些的，最明显的就是异步返回参数会返回trade_status（交易状态）参数 但是同步返回的参数是没有的，所以建议使用异步返回的参数来确认交易是否成功 如果trade_status=TRADE_SUCCESS 表示就是交易成功了
注：老版接口（如即时到账）会返回trade_status，具体返回参数建议参考文档。

- 发送的时间不同
一般情况下，异步通知比同步通知稍慢（几乎无感知）。同步和异步通知的发送和网络环境、服务器处理速度等因素关系较大，具体是先接收到同步通知还是异步通知，建议以实际接收到的为准。

- 参数返回的逻辑不一样
1. 同步通知是在支付成功后的页面跳转，一般是用于展示给用户的。
付款成功页面会同步跳转到商户请求参数中传入的return_url地址页面。
2. 异步通知参数是直接由支付宝服务器发送到商户的异步地址上的，不存在页面跳转是直接的服务器交互。

>注：同步验签和异步验签方式相同。
#### 关于应用私钥 (privateKey)

应用私钥是我们通过工具生成的私钥，调用支付宝接口的时候，我们需要使用该私钥对参数进行签名。

#### 支持 RSA 签名及验证

采用的是 RSA2 签名

- 加签文档:  [https://opendocs.alipay.com/common/057k53](https://opendocs.alipay.com/common/057k53)、[https://opendocs.alipay.com/support/01rave](https://opendocs.alipay.com/support/01rave)
- 验签文档： [https://opendocs.alipay.com/common/02mse7](https://opendocs.alipay.com/common/02mse7)
- 常见问题： [https://opensupport.alipay.com/support/knowledgeInfo/9483?ant_source=antsupport](https://opensupport.alipay.com/support/knowledgeInfo/9483?ant_source=antsupport)
  [https://opensupport.alipay.com/support/helpcenter/192/201602487632?ant_source=antsupport](https://opensupport.alipay.com/support/helpcenter/192/201602487632?ant_source=antsupport)
