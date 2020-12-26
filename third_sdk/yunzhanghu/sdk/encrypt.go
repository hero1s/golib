package sdk

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// Sign 生成签名
func Sign(data, mess, timestamp, appKey string) string {
	signPars := fmt.Sprintf("data=%s&mess=%s&timestamp=%s&key=%s", data, mess, timestamp, appKey)
	sign := hmac.New(sha256.New, []byte(appKey))
	sign.Write([]byte(signPars))
	calcMac := hex.EncodeToString(sign.Sum(nil))
	return calcMac
}

// Encrypt data加密
func Encrypt(originData []byte, des3key string) (string, error) {
	crypt, err := TripleDesEncrypt(originData, []byte(des3key))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(crypt), nil
}

// Decrypt data解密
func Decrypt(data string, des3key string) ([]byte, error) {
	crypt, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return TripleDesDecrypt([]byte(crypt), []byte(des3key))
}

// TripleDesEncrypt 3DES加密
func TripleDesEncrypt(originData, des3key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(des3key)
	if err != nil {
		return nil, err
	}
	originData = PKCS5Padding(originData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, des3key[:8])
	crypt := make([]byte, len(originData))
	blockMode.CryptBlocks(crypt, originData)
	return crypt, nil
}

// TripleDesDecrypt 3DES解密
func TripleDesDecrypt(crypt, des3key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(des3key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, des3key[:8])
	originData := make([]byte, len(crypt))
	blockMode.CryptBlocks(originData, crypt)
	originData = PKCS5UnPadding(originData)
	return originData, nil
}

// PKCS5Padding 填充
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS5UnPadding 取消填充
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后⼀一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
