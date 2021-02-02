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

//Java 的RSA方式 myhash := crypto.SHA1
func RsaSign(data []byte, privateKey []byte,myhash crypto.Hash) (string, error) {
	hashInstance := myhash.New()
	hashInstance.Write(data)
	hashed := hashInstance.Sum(nil)

	block, _ := pem.Decode(privateKey)
	if block == nil {
		return "", errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Parse private key  error:%v", err))
	}
	p := priv.(*rsa.PrivateKey)
	signature, err := rsa.SignPKCS1v15(rand.Reader, p, myhash, hashed)
	return base64.StdEncoding.EncodeToString(signature), err
}

// Java RSA 公钥验证 myhash := crypto.SHA1
func RsaVerify(data []byte, publicKeyStr []byte,myhash crypto.Hash) error {
	// 2、选择hash算法，对需要签名的数据进行hash运算
	hashInstance := myhash.New()
	hashInstance.Write(data)
	hashed := hashInstance.Sum(nil)

	block, _ := pem.Decode(publicKeyStr)
	if block == nil {
		return errors.New("公钥信息错误！")
	}
	// 3、解析DER编码的公钥，生成公钥接口
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	publicKey := publicKeyInterface.(*rsa.PublicKey)
	// 4、RSA验证数字签名（参数是公钥对象、哈希类型、签名文件的哈希串、签名后的字节）
	return rsa.VerifyPKCS1v15(publicKey, myhash, hashed, nil)
}