package crypto

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"github.com/pkg/errors"
)

//md5方法
func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// 获取SHA1
func SHA1(str string) string {
	_sha1 := sha1.New()
	_sha1.Write([]byte(str))
	return hex.EncodeToString(_sha1.Sum([]byte(nil)))
}

// 获取SHA256
func SHA256(str string) string {
	_sha256 := sha256.New()
	_sha256.Write([]byte(str))
	return hex.EncodeToString(_sha256.Sum([]byte(nil)))
}

//MD5WithRsa
func Md5WithRsa(params string, privateKey []byte) (string, error) {
	data := []byte(params)
	hashMd5 := md5.Sum(data)
	hashed := hashMd5[:]

	block, _ := pem.Decode(privateKey)
	if block == nil {
		return "", errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("Parse private key  error:%v", err))
	}
	p := priv.(*rsa.PrivateKey)
	signature, err := rsa.SignPKCS1v15(rand.Reader, p, crypto.MD5, hashed)

	return base64.StdEncoding.EncodeToString(signature), err
}
