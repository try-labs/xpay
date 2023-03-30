package alipay

import (
	"encoding/json"
	"fmt"
)

/**
 * @author: Sam
 * @since:
 * @date: 2023/3/20 14:10
 * @desc:
 */

type OauthTokenReq struct {
	GrantType    string `json:"grant_type"`              // 必选	20 授权方式。支持： 1.authorization_code，表示换取使用用户授权码code换取授权令牌access_token。 2.refresh_token，表示使用refresh_token刷新获取新授权令牌。
	Code         string `json:"code,omitempty"`          // 可选 40 授权码，用户对应用授权后得到。本参数在 grant_type 为 authorization_code 时必填；为 refresh_token 时不填。 4b203fe6c11548bcabd8da5bb087a83b
	RefreshToken string `json:"refresh_token,omitempty"` //	可选	40 刷新令牌，上次换取访问令牌时得到。本参数在 grant_type 为 authorization_code 时不填；为 refresh_token 时必填，且该值来源于此接口的返回值 app_refresh_token（即至少需要通过 grant_type=authorization_code 调用此接口一次才能获取）。	201208134b203fe6c11548bcabd8da5bb087a83b
	BaseAliPayRequest
}

var _ IAliPayRequest = &OauthTokenReq{}

func (r *OauthTokenReq) RequestApi() string {
	return "alipay.system.oauth.token"
}

func (r *OauthTokenReq) DoValidate() error {
	if r.GrantType != GrantAuthorizationCode && r.GrantType != GrantRefreshToken {
		return fmt.Errorf("参数grant_type必传")
	}
	if r.GrantType == GrantAuthorizationCode && r.Code == "" {
		return fmt.Errorf("参数code在grant_type 为 authorization_code 时必填；为 refresh_token 时不填。")
	}
	if r.GrantType == GrantRefreshToken && r.RefreshToken == "" {
		return fmt.Errorf("参数code在grant_type 为 refresh_token  时必填；为 authorization_code 时不填。")
	}
	return nil
}

/*
OauthTokenRes 发生错误时返回结构
调用接口返回10000是表示接口调用成功，但是 不等于业务逻辑处理成功。
{
"error_response": {
"code": "20000",
"msg": "Service Currently Unavailable",
"sub_code": "isp.unknow-error",
"sub_msg": "系统繁忙"
},
"sign": "ERITJKEIJKJHKKKKKKKHJEREEEEEEEEEEE"
}
*/
type OauthTokenRes struct {
	*OauthTokenResContent `json:"alipay_system_oauth_token_response,omitempty"`
	*CommonRes            `json:"error_response,omitempty"`
	SignCertSn
}

//	func (r *OauthTokenRes) SetSubCode(subCode string) {
//		r.SubCode = subCode
//	}
func (r *OauthTokenRes) String() string {
	buff, _ := json.Marshal(r)
	return string(buff)
}

type OauthTokenResContent struct {
	CommonRes
	UserId       string `json:"user_id"`              // 必选	16 支付宝用户的唯一标识。以2088开头的16位数字。
	AccessToken  string `json:"access_token"`         // 必选    40 访问令牌。通过该令牌调用需要授权类接口
	ExpiresIn    string `json:"expires_in"`           // 必选	16 访问令牌的有效时间，单位是秒。 3600
	RefreshToken string `json:"refresh_token"`        // 必选	40 刷新令牌。通过该令牌可以刷新access_token
	ReExpiresIn  string `json:"re_expires_in"`        // 必选	16 刷新令牌的有效时间，单位是秒。 3600
	AuthStart    string `json:"auth_start,omitempty"` // 可选	32 授权token开始时间，作为有效期计算的起点 2010-11-11 11:11:11
}

////////////////////////////////////////////////

var _ IAliPayRequest = &UserInfoShareReq{}

type UserInfoShareReq struct {
	AuthToken string `json:"auth_token"` // 必选	40 用户授权令牌，同 access_token（用户访问令牌）。针对用户授权接口，获取用户相关数据时，用于标识用户授权关系。需使用 auth_code（用户授权码）换取此令牌，详情见 用户授权
	BaseAliPayRequest
}

func (r *UserInfoShareReq) DoValidate() error {
	length := len(r.AuthToken)
	if length != 40 {
		return fmt.Errorf("当前用户授权令牌auth_token参数长度为%d，不符合长度40要求", length)
	}
	return nil
}

func (r *UserInfoShareReq) RequestApi() string {
	return "alipay.user.info.share"
}

type UserInfoShareRes struct {
	UserInfoShareResContent `json:"alipay_user_info_share_response"`
	SignCertSn
}
type UserInfoShareResContent struct {
	CommonRes
	UserId   string `json:"user_id"`             // 必选	16 支付宝用户的userId。 2088102104794936
	Avatar   string `json:"avatar,omitempty"`    // 必选	400 用户头像地址。 注意：如果没有数据（用户未设置）时不会返回该信息，请做好容错。 http://tfsimg.alipay.com/images/partner/T1uIxXXbpXXXXXXXX
	City     string `json:"city,omitempty"`      // 必选	20 市名称。 安庆
	NickName string `json:"nick_name,omitempty"` // 必选	50  用户昵称。 注意：如果没有数据（用户未设置）时不会返回该信息，请做好容错。 支付宝小二
	Province string `json:"province,omitempty"`  // 必选	20 省份名称。 安徽省
	Gender   string `json:"gender,omitempty"`    // 可选	10 性别。枚举值如下： F：女性； M：男性。
}

func (r *UserInfoShareRes) String() string {
	buff, _ := json.Marshal(r)
	return string(buff)
}

///////////////////////////////////////////

var _ IAliPayRequest = &UserInfoAuthReq{}

type UserInfoAuthReq struct {
	Scopes []string `json:"scopes"` //	必选	1024	接口权限值，枚举如下：
	/*
	* auth_base：以auth_base为scope发起的网页授权，用于获取进入页面的用户的 userId，并且是静默授权并自动跳转到回调页的。用户感知的就是直接进入了回调页（通常是业务页面）。
	* auth_user：以auth_user为scope发起的网页授权，是用来获取用户的基本信息的（比如头像、昵称等）。但这种授权需要用户手动同意，用户同意后，就可在授权后获取到该用户的基本信息。
	 */
	State string `json:"state"` // 必选	100	商户自定义参数，只允许base64字符（长度小于等于100）。 说明：
	/*
	* 传入时将在用户授权后，重定向到redirect_uri 时会原样回传给商户。 为防止CSRF攻击。
	* 建议开发者请求授权时传入state参数，该参数要做到既不可预测，又可以证明客户端和当前第三方网站的登录认证状态存在关联。	init
	 */
	// 自己添加
	ReturnUrl string `json:"return_url,omitempty" url:"return_url,omitempty"` // 可选	256 HTTP/HTTPS开头字符串
	BaseAliPayRequest
}

func (r *UserInfoAuthReq) DoValidate() error {
	if len(r.Scopes) == 0 {
		return fmt.Errorf("接口权限值 scopes 参数必填")
	}
	return nil
}

func (r *UserInfoAuthReq) RequestApi() string {
	return "alipay.user.info.auth"
}

type UserInfoAuthRes struct {
	UserInfoAuthResContent `json:"alipay_user_info_auth_response"`
	AlipayCertSn           string `json:"alipay_cert_sn,omitempty"`
	Sign                   string `json:"sign"` // 签名
}
type UserInfoAuthResContent struct {
	CommonRes
}

func (r *UserInfoAuthRes) SetSign(sign string) {
	r.Sign = sign
}

func (r *UserInfoAuthRes) SetAlipayCertSn(alipayCertSn string) {
	r.AlipayCertSn = alipayCertSn
}
