package alipay

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/google/go-querystring/query"
	"os"
	"sort"
	"strings"
)

/**
 * @author: Sam
 * @since:
 * @date: 2023/3/14 17:50
 * @desc: 签名和验签
 */

type Signer interface {
	// SetSignContent 设置待签名的数据
	SetSignContent(*CommonReqParam)
	// Sign 生成签名
	Sign() (string, error)
	// Encode url.encode
	Encode() (string, error)
}

type Verifier interface {
	// VerifySign 校验签名
	VerifySign(scene VerificationScene, sign string, buff []byte, otherParam ...string) error
}

type SignVerifier interface {
	Signer
	Verifier
}

// 待验证buff
type responseBuff []byte

// GetWaitSignData 获取待验证的字符串
func (r responseBuff) GetWaitSignData() (string, error) {
	str := string(r)
	// {"alipay_trade_query_response":
	indexStart := strings.Index(str, `_response":`) + 11
	indexEnd := strings.Index(str, `,"alipay_cert_sn":`)
	if indexEnd > 0 && indexEnd > indexStart {
		return str[indexStart:indexEnd], nil
	}
	indexEnd = strings.Index(str, `,"sign":`)
	if indexEnd < 0 {
		return "", ErrNotContainsSignData
	}
	return str[indexStart:indexEnd], nil
}

/*
IsBizSuccess 判断业务请求是否成功
{
	"error_response": {
		"code": "40002",
		"msg": "Invalid Arguments",
		"sub_code": "isv.isv.grant-type-invalid",
		"sub_msg": "grant_type参数不正确"
	},
	"alipay_cert_sn": "e97e0b247e2637dac32fe3989c73e1eb",
	"sign": "AxoPq1O"
}
*/

func (r responseBuff) IsErrorResponse() bool {
	str := string(r)
	indexStart := strings.Index(str, `error_response":`)
	return indexStart > -1
}

// IsHtmlError 是否报错,有时直接返回html界面
func (r responseBuff) IsHtmlError() bool {
	str := string(r)
	return strings.Contains(str, `<html>`)
}

type Signature struct {
	// 应用id
	appId string
	// 应用私钥
	appPrivateKey   *rsa.PrivateKey
	waitSignContent *CommonReqParam
}

func (r *Signature) SetSignContent(param *CommonReqParam) {
	param.AppId = r.appId
	r.waitSignContent = param
}

func (r *Signature) Encode() (string, error) {
	tempParam := r.waitSignContent
	sign, err := r.Sign()
	if err != nil {
		return "", err
	}
	tempParam.Sign = sign
	param, err := query.Values(tempParam)
	if err != nil {
		return "", err
	}
	return param.Encode(), nil
}

func (r *Signature) Sign() (string, error) {
	values, err := query.Values(r.waitSignContent)
	if err != nil {
		return "", err
	}
	valueList := make([]string, 0, len(values))
	for key := range values {
		var value = strings.TrimSpace(values.Get(key))
		if len(value) > 0 {
			valueList = append(valueList, key+"="+value)
		}
	}
	sort.Strings(valueList)
	var src = strings.Join(valueList, "&")
	sign, err := RSASignWithKey([]byte(src), r.appPrivateKey, crypto.SHA256)
	if err != nil {
		return "", err
	}
	signStr := base64.StdEncoding.EncodeToString(sign)
	return signStr, nil
}

type NormalRSA2SignStrategy struct {
	Signature
	appPublicKey *rsa.PublicKey
	// 公钥
	aliPayPublicKey *rsa.PublicKey
}

func NewNormalRSA2SignStrategy(appId, privateKey, publicKey, aliPayPublicKey string) SignVerifier {
	var appPrivateKey *rsa.PrivateKey
	var appPublicKey *rsa.PublicKey
	var err error
	if appPrivateKey, err = ParsePKCS1PrivateKey(FormatPKCS1PrivateKey(privateKey)); err != nil {
		if appPrivateKey, err = ParsePKCS8PrivateKey(FormatPKCS8PrivateKey(privateKey)); err != nil {
			panic(err)
		}
	}
	if appPublicKey, err = ParsePublicKey(FormatPublicKey(publicKey)); err != nil {
		panic(err)
	}
	var aliPublicKey *rsa.PublicKey
	if aliPublicKey, err = ParsePublicKey(FormatPublicKey(aliPayPublicKey)); err != nil {
		panic(err)
	}
	strategy := &NormalRSA2SignStrategy{
		Signature:       Signature{appId: appId, appPrivateKey: appPrivateKey},
		appPublicKey:    appPublicKey,
		aliPayPublicKey: aliPublicKey,
	}
	return strategy
}

func (r *NormalRSA2SignStrategy) Encode() (string, error) {
	tempParam := r.waitSignContent
	sign, err := r.Sign()
	if err != nil {
		return "", err
	}
	tempParam.Sign = sign
	param, err := query.Values(tempParam)
	if err != nil {
		return "", err
	}
	return param.Encode(), nil
}

func (r *NormalRSA2SignStrategy) VerifySign(scene VerificationScene, sign string, buff []byte, otherParam ...string) error {
	var err error
	var signBytes []byte
	if signBytes, err = base64.StdEncoding.DecodeString(sign); err != nil {
		return err
	}
	waitSignBuff := buff
	if scene == SyncVerificationScene {
		var waitSignData string
		if waitSignData, err = responseBuff(buff).GetWaitSignData(); err != nil {
			return err
		}
		waitSignBuff = []byte(waitSignData)
	}
	return RSAVerifyWithKey(waitSignBuff, signBytes, r.aliPayPublicKey, crypto.SHA256)
}

type CertSignStrategy struct {
	Signature
	// 应用公钥证书 SN
	appCertSN string
	// 支付宝根证书 SN
	alipayRootCertSn string
	// 支付宝公钥证书序列号
	alipayPublicCertSN string
	// 支付宝公钥证书序列号=>支付宝公钥证书
	alipayPublicKeyList map[string]*rsa.PublicKey
	alipayPublicKey     *rsa.PublicKey
}

func NewCertSignStrategy(appId, privateKey, appPublicCert, alipayRootCert, alipayPublicCert string) SignVerifier {
	var err error
	object := &CertSignStrategy{alipayPublicKeyList: make(map[string]*rsa.PublicKey, 0)}
	if err = object.LoadAppPublicCertFile(appPublicCert); err != nil {
		panic(err)
	}
	if err = object.LoadRootCertFile(alipayRootCert); err != nil {
		panic(err)
	}
	if err = object.LoadPublicCertFile(alipayPublicCert); err != nil {
		panic(err)
	}
	var appPrivateKey *rsa.PrivateKey
	if appPrivateKey, err = ParsePKCS1PrivateKey(FormatPKCS1PrivateKey(privateKey)); err != nil {
		if appPrivateKey, err = ParsePKCS8PrivateKey(FormatPKCS8PrivateKey(privateKey)); err != nil {
			panic(err)
		}
	}
	object.Signature = Signature{appId: appId, appPrivateKey: appPrivateKey}
	return object
}

// LoadAppPublicCertFile 加载应用公钥证书
func (r *CertSignStrategy) LoadAppPublicCertFile(filename string) error {
	buff, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	cert, err := ParseCertificate(buff)
	if err != nil {
		return err
	}
	r.appCertSN = GetCertSN(cert)
	return nil
}

// LoadRootCertFile 加载支付宝根证书
func (r *CertSignStrategy) LoadRootCertFile(filename string) error {
	buff, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	var certStrList = strings.Split(string(buff), CertificateEnd)
	certSNSlice := make([]string, 0, len(certStrList))
	for _, certStr := range certStrList {
		certStr = certStr + CertificateEnd
		cert, _ := ParseCertificate([]byte(certStr))
		if cert != nil && (cert.SignatureAlgorithm == x509.SHA256WithRSA || cert.SignatureAlgorithm == x509.SHA1WithRSA) {
			certSNSlice = append(certSNSlice, GetCertSN(cert))
		}
	}
	r.alipayRootCertSn = strings.Join(certSNSlice, "_")
	return nil
}

// LoadPublicCertFile 加载支付宝公钥证书
func (r *CertSignStrategy) LoadPublicCertFile(filename string) error {
	buff, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	cert, err := ParseCertificate(buff)
	if err != nil {
		return err
	}
	key, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return ErrTrans
	}
	r.alipayPublicCertSN = GetCertSN(cert)
	r.alipayPublicKeyList[r.alipayPublicCertSN] = key
	r.alipayPublicKey = key

	return nil
}

func (r *CertSignStrategy) SetSignContent(param *CommonReqParam) {
	param.AlipayRootCertSn = r.alipayRootCertSn
	param.AppCertSn = r.appCertSN
	param.AppId = r.appId
	r.waitSignContent = param
}

// VerifySign 异步通知验签 公钥、证书两种模式下，异步通知验签方式相同
func (r *CertSignStrategy) VerifySign(scene VerificationScene, sign string, buff []byte, otherParam ...string) error {
	var err error
	var signBytes []byte
	if signBytes, err = base64.StdEncoding.DecodeString(sign); err != nil {
		return err
	}
	if scene == SyncVerificationScene {
		if len(otherParam) == 0 {
			return errors.New("缺少app_cert_sn参数")
		}
		waitSignData, err := responseBuff(buff).GetWaitSignData()
		if err != nil {
			return err
		}
		publicKey := r.alipayPublicKeyList[otherParam[0]]
		if publicKey == nil {
			return fmt.Errorf("证书序列号：%s，对应的公钥证书不存在", otherParam[0])
		}
		return RSAVerifyWithKey([]byte(waitSignData), signBytes, publicKey, crypto.SHA256)
	}
	return RSAVerifyWithKey(buff, signBytes, r.alipayPublicKey, crypto.SHA256)
}
