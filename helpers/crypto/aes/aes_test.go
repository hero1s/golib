package aes

import (
	"testing"

	"github.com/hero1s/golib/log"
)

var (
	secretKey = "GYBh3Rmey7nNzR/NpV0vAw=="
)

func TestAesCBCEncryptToString(t *testing.T) {
	originData := "www.gopay.ink"
	log.Debugf("originData:", originData)
	encryptData, err := AesCBCEncryptToString([]byte(originData), secretKey)
	if err != nil {
		log.Errorf("AesCBCEncryptToString:", err)
		return
	}
	log.Debugf("encryptData:", encryptData)
	origin, err := AesDecryptToBytes(encryptData, secretKey)
	if err != nil {
		log.Errorf("AesDecryptToBytes:", err)
		return
	}
	log.Debugf("origin:", string(origin))
}
