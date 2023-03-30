package alipay

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"strings"
)

/**
 * @author: Sam
 * @since:
 * @date: 2023/3/20 17:30
 * @desc:
 */

const (
	PublicKeyPrefix = "-----BEGIN PUBLIC KEY-----"
	PublicKeySuffix = "-----END PUBLIC KEY-----"

	PKCS1Prefix = "-----BEGIN RSA PRIVATE KEY-----"
	PKCS1Suffix = "-----END RSA PRIVATE KEY-----"

	PKCS8Prefix = "-----BEGIN PRIVATE KEY-----"
	PKCS8Suffix = "-----END PRIVATE KEY-----"

	PublicKeyType     = "PUBLIC KEY"
	PrivateKeyType    = "PRIVATE KEY"
	RSAPrivateKeyType = "RSA PRIVATE KEY"
)

var (
	ErrLoadPrivateKey  = errors.New("xpay: private key failed to load")
	ErrLoadPublicKey   = errors.New("xpay: public key failed to load")
	ErrLoadCertificate = errors.New("xpay: certificate  failed to load")
)

func ParsePKCS1PrivateKey(data []byte) (key *rsa.PrivateKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode(data)
	if block == nil {
		return nil, ErrLoadPrivateKey
	}
	key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, err
}

func ParseCertificate(b []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(b)
	if block == nil {
		return nil, ErrLoadCertificate
	}
	csr, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return csr, nil
}

func GetCertSN(cert *x509.Certificate) string {
	var value = md5.Sum([]byte(cert.Issuer.String() + cert.SerialNumber.String()))
	return hex.EncodeToString(value[:])
}

func RSAVerifyWithKey(ciphertext, sign []byte, key *rsa.PublicKey, hash crypto.Hash) error {
	var h = hash.New()
	h.Write(ciphertext)
	var hashed = h.Sum(nil)
	return rsa.VerifyPKCS1v15(key, hash, hashed, sign)
}

func FormatPublicKey(raw string) []byte {
	return formatKey(raw, PublicKeyPrefix, PublicKeySuffix, 64)
}

func FormatPKCS1PrivateKey(raw string) []byte {
	raw = strings.Replace(raw, PKCS8Prefix, "", 1)
	raw = strings.Replace(raw, PKCS8Suffix, "", 1)
	return formatKey(raw, PKCS1Prefix, PKCS1Suffix, 64)
}

func FormatPKCS8PrivateKey(raw string) []byte {
	raw = strings.Replace(raw, PKCS1Prefix, "", 1)
	raw = strings.Replace(raw, PKCS1Suffix, "", 1)
	return formatKey(raw, PKCS8Prefix, PKCS8Suffix, 64)
}

func ParsePKCS8PrivateKey(data []byte) (key *rsa.PrivateKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode(data)
	if block == nil {
		return nil, ErrLoadPrivateKey
	}
	rawKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := rawKey.(*rsa.PrivateKey)
	if !ok {
		return nil, ErrLoadPrivateKey
	}
	return key, err
}

func ParsePublicKey(data []byte) (key *rsa.PublicKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode(data)
	if block == nil {
		return nil, ErrLoadPublicKey
	}
	var pubInterface interface{}
	pubInterface, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, ErrLoadPublicKey
	}

	return key, err
}

func RSASignWithKey(plaintext []byte, key *rsa.PrivateKey, hash crypto.Hash) ([]byte, error) {
	var h = hash.New()
	h.Write(plaintext)
	var hashed = h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, key, hash, hashed)
}

func formatKey(raw, prefix, suffix string, lineCount int) []byte {
	raw = strings.Replace(raw, prefix, "", 1)
	raw = strings.Replace(raw, suffix, "", 1)
	raw = strings.Replace(raw, " ", "", -1)
	raw = strings.Replace(raw, LineBreak, "", -1)
	raw = strings.Replace(raw, "\r", "", -1)
	raw = strings.Replace(raw, "\t", "", -1)

	var sl = len(raw)
	var c = sl / lineCount
	if sl%lineCount > 0 {
		c = c + 1
	}
	var buf bytes.Buffer
	buf.WriteString(prefix + LineBreak)
	for i := 0; i < c; i++ {
		var b = i * lineCount
		var e = b + lineCount
		if e > sl {
			buf.WriteString(raw[b:])
		} else {
			buf.WriteString(raw[b:e])
		}
		buf.WriteString(LineBreak)
	}
	buf.WriteString(suffix)
	return buf.Bytes()
}
