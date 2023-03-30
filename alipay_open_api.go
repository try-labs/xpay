package alipay

import (
	"encoding/json"
	"fmt"
)

/**
 * @author: Sam
 * @since:
 * @date: 2023/3/23 10:01
 * @desc:
 */
var _ IAliPayRequest = &OpenAuthTokenAppReq{}

type OpenAuthTokenAppReq struct {
	GrantType    string `json:"grant_type"`              // 必选	20 授权方式。支持： 1.authorization_code，表示换取使用用户授权码code换取授权令牌access_token。 2.refresh_token，表示使用refresh_token刷新获取新授权令牌。
	Code         string `json:"code,omitempty"`          // 可选 40 授权码，用户对应用授权后得到。本参数在 grant_type 为 authorization_code 时必填；为 refresh_token 时不填。 4b203fe6c11548bcabd8da5bb087a83b
	RefreshToken string `json:"refresh_token,omitempty"` //	可选	40 刷新令牌，上次换取访问令牌时得到。本参数在 grant_type 为 authorization_code 时不填；为 refresh_token 时必填，且该值来源于此接口的返回值 app_refresh_token（即至少需要通过 grant_type=authorization_code 调用此接口一次才能获取）。	201208134b203fe6c11548bcabd8da5bb087a83b
	BaseAliPayRequest
}

func (r *OpenAuthTokenAppReq) DoValidate() error {
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

func (r *OpenAuthTokenAppReq) RequestApi() string {
	return "alipay.system.oauth.token"
}

type OpenAuthTokenAppRes struct {
	*OpenAuthTokenAppResContent `json:"alipay_open_auth_token_app_response,omitempty"`
	*CommonRes                  `json:"error_response,omitempty"`
	SignCertSn
}

func (r *OpenAuthTokenAppRes) String() string {
	buff, _ := json.Marshal(r)
	return string(buff)
}

type OpenAuthTokenAppResContent struct {
	CommonRes
	UserId          string `json:"user_id"`           // 必选	16 授权商户的user_id 2088102150527498
	AuthAppId       string `json:"auth_app_id"`       // 必选	20 授权商户的appid 2013121100055554
	AppAuthToken    string `json:"app_auth_token"`    // 必选	40 应用授权令牌 201509BBeff9351ad1874306903e96b91d248A36
	AppRefreshToken string `json:"app_refresh_token"` // 必选	40 刷新令牌 201509BBdcba1e3347de4e75ba3fed2c9abebE36
	ExpiresIn       string `json:"expires_in"`        // 必选	16 该字段已作废，应用令牌长期有效，接入方不需要消费该字段 123456
	ReExpiresIn     string `json:"re_expires_in"`     // 必选	16 刷新令牌的有效时间（从接口调用时间作为起始时间），单位到秒 123456
}
