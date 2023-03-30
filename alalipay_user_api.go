package alipay

import (
	"encoding/json"
	"fmt"
)

/**
 * @author: Sam
 * @since:
 * @date: 2023/3/23 23:01
 * @desc:
 */

var _ IAliPayRequest = &UserCertifyOpenInitializeReq{}

type UserCertifyOpenInitializeReq struct {
	OuterOrderNo string `json:"outer_order_no"` // 必选	32 商户请求的唯一标识，商户要保证其唯一性，值为32位长度的字母数字组合。建议：前面几位字符是商户自定义的简称，中间可以使用一段时间，后段可以使用一个随机或递增序列
	BizCode      string `json:"biz_code"`       // 必选	32 认证场景码。入参支持的认证场景码和商户签约的认证场景相关，取值如下:
	//FACE：多因子人脸认证
	//CERT_PHOTO：多因子证照认证
	//CERT_PHOTO_FACE ：多因子证照和人脸认证
	//SMART_FACE：多因子快捷认证
	IdentityParam  OpenCertifyIdentityParam  `json:"identity_param"`  // 必选	 需要验证的身份信息
	MerchantConfig OpenCertifyMerchantConfig `json:"merchant_config"` // 必选 商户个性化配置
	BaseAliPayRequest
}

func (r *UserCertifyOpenInitializeReq) DoValidate() error {
	outerOrderNoLen := len(r.OuterOrderNo)
	if outerOrderNoLen == 0 || outerOrderNoLen > 32 {
		return fmt.Errorf("商户请求的唯一标识outer_order_no参数长度为%d，不符合长度1-32个字符的要求", outerOrderNoLen)
	}
	bizCodeLen := len(r.BizCode)
	if bizCodeLen == 0 || bizCodeLen > 32 {
		return fmt.Errorf("认证场景码biz_code参数长度为%d，不符合长度1-32个字符的要求", bizCodeLen)
	}
	return nil
}

func (r *UserCertifyOpenInitializeReq) RequestApi() string {
	return "alipay.user.certify.open.initialize"
}

type OpenCertifyIdentityParam struct {
	IdentityType string `json:"identity_type"`       // 必选	30 1.若本人验证，使用CERT_INFO； 2.若代他人验证，使用AGENT_CERT_INFO； 枚举值 证件信息: CERT_INFO 代理人证件信息: AGENT_CERT_INFO
	CertType     string `json:"cert_type,omitempty"` // 可选	100
	/*cert_type：
	1、若为身份证，填IDENTITY_CARD；
	2、若为港澳居民来往内地通行证，填HOME_VISIT_PERMIT_HK_MC；
	3、若为台湾居民来往内地通行证，填HOME_VISIT_PERMIT_TAIWAN；
	4、若为港澳居民居住证，填RESIDENCE_PERMIT_HK_MC；
	5、若为台湾居民居住证，填RESIDENCE_PERMIT_TAIWAN；
	枚举值
	身份证: IDENTITY_CARD
	港澳居民来往内地通行证: HOME_VISIT_PERMIT_HK_MC
	台湾居民来往内地通行证: HOME_VISIT_PERMIT_TAIWAN
	港澳居民居住证: RESIDENCE_PERMIT_HK_MC
	台湾居民居住证: RESIDENCE_PERMIT_TAIWAN
	注意事项: 在identity_type为CERT_INFO或者AGENT_CERT_INFO时，该字段必填
	*/
	CertName string `json:"cert_name,omitempty"` // 可选	50 填入真实姓名 注意事项 在identity_type为CERT_INFO或者AGENT_CERT_INFO时，该字段必填
	CertNo   string `json:"cert_no,omitempty"`   // 可选	30  填入姓名相匹配的证件号码 注意事项 在identity_type为CERT_INFO或者AGENT_CERT_INFO时，该字段必填
}

type OpenCertifyMerchantConfig struct {
	FaceReserveStrategy string `json:"face_reserve_strategy,omitempty"` //	可选	32 不传默认为reserve 枚举值 保存活体人脸: reserve 不保存活体人脸: never
	ReturnUrl           string `json:"return_url,omitempty"`            //	必选	4096 认证成功后需要跳转的地址，一般为商户业务页面；若无跳转地址可填空字符"";

}

type UserCertifyOpenInitializeRes struct {
	UserCertifyOpenInitializeResContent `json:"alipay_user_certify_open_initialize_response"`
	SignCertSn
}

type UserCertifyOpenInitializeResContent struct {
	CommonRes
	CertifyId string // 必选	32 本次申请操作的唯一标识，商户需要记录，后续的操作都需要用到
}

func (r *UserCertifyOpenInitializeRes) String() string {
	buff, _ := json.Marshal(r)
	return string(buff)
}

// ///////////////////////////////////////////
var _ IAliPayRequest = &UserCertifyOpenQueryReq{}

type UserCertifyOpenQueryReq struct {
	CertifyId string `json:"certify_id"` // 必选	32 本次申请操作的唯一标识，通过alipay.user.certify.open.initialize(身份认证初始化服务)接口同步响应获取。
	BaseAliPayRequest
}

func (r *UserCertifyOpenQueryReq) DoValidate() error {
	certifyIdLen := len(r.CertifyId)
	if certifyIdLen == 0 || certifyIdLen > 32 {
		return fmt.Errorf("本次申请操作的唯一标识certify_id参数长度为%d，不符合长度1-32个字符的要求", certifyIdLen)
	}
	return nil
}

func (r *UserCertifyOpenQueryReq) RequestApi() string {
	return "alipay.user.certify.open.query"
}

type UserCertifyOpenQueryRes struct {
	UserCertifyOpenQueryResContent `json:"alipay_user_certify_open_query_response"`
	SignCertSn
}
type UserCertifyOpenQueryResContent struct {
	CommonRes
	Passed       string `json:"passed"`                  // 必选	16 是否通过，通过为T，不通过为F
	IdentityInfo string `json:"identity_info,omitempty"` // 可选	1024 认证的主体信息，一般的认证场景返回为空
	MaterialInfo string `json:"material_info,omitempty"` // 可选	2867200 认证主体附件信息，主要为图片类材料，一般的认证场景都是返回空
	FailReason   string `json:"fail_reason,omitempty"`   // 可选	512 认证错误原因，扩展字段，默认不返回
}

func (r *UserCertifyOpenQueryRes) String() string {
	buff, _ := json.Marshal(r)
	return string(buff)
}

////////////////////////////////////////////////

var _ IAliPayRequest = &UserCertifyOpenCertifyReq{}

type UserCertifyOpenCertifyReq struct {
	ReturnUrl string `json:"return_url,omitempty" url:"return_url,omitempty"` // 可选	256 HTTP/HTTPS开头字符串
	CertifyId string `json:"certify_id"`                                      //	必选	32 本次申请操作的唯一标识，由开放认证初始化接口调用后生成，后续的操作都需要用到
	IAliPayRequest
}

func (r *UserCertifyOpenCertifyReq) DoValidate() error {
	certifyIdLen := len(r.CertifyId)
	if certifyIdLen == 0 || certifyIdLen > 32 {
		return fmt.Errorf("本次申请操作的唯一标识certify_id参数长度为%d，不符合长度1-32个字符的要求", certifyIdLen)
	}
	return nil
}

func (r *UserCertifyOpenCertifyReq) RequestApi() string {
	return "alipay.user.certify.open.certify"
}

////////////////////////////////////////////////
