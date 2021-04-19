package tcp

import (
	"bytes"
	"crypto/aes"
	"encoding/hex"
	"strings"
)

/**
* 输入输出都是string类型，简化使用
 */
func EncryptString(plantText, key string) (string, error) {
	ret, err := Encrypt([]byte(plantText), []byte(key))
	if err != nil {
		return plantText, err
	}

	return strings.ToUpper(hex.EncodeToString(ret)), nil
}

func Encrypt(plantText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key) //选择加密算法
	if err != nil {
		return nil, err
	}

	plantText = ECBZeroPadding(plantText)

	blockModel := NewECBEncrypter(block)

	ciphertext := make([]byte, len(plantText))

	blockModel.CryptBlocks(ciphertext, plantText)
	return ciphertext, nil
}

/**
*  ECB/NoPadding
*
 */
func ECBZeroPadding(plantText []byte) []byte {
	bit16 := len(plantText) % 16

	if bit16 == 0 {
		return plantText
	} else {
		// 计算补足的位数，填补16的位数，例如 10 = 16, 17 = 32, 33 = 48
		return append(plantText, bytes.Repeat([]byte{byte(0)}, 16-bit16)...)
	}
}

func DecryptString(ciphertext, key string) (string, error) {
	ciphertext1, err := hex.DecodeString(ciphertext)
	if err != nil {
		return ciphertext, err
	}

	ret, err1 := Decrypt(ciphertext1, []byte(key))

	if err1 != nil {
		return ciphertext, err1
	}

	return string(ret), nil
}

func Decrypt(ciphertext, key []byte) ([]byte, error) {
	keyBytes := []byte(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}
	blockModel := NewECBDecrypter(block)
	plantText := make([]byte, len(ciphertext))
	blockModel.CryptBlocks(plantText, ciphertext)
	plantText = ECBUnZeroPadding(plantText)
	return plantText, nil
}

func ECBUnZeroPadding(plantText []byte) []byte {
	l := len(plantText)

	//大于16，需要把结尾是0的截取掉
	if l >= 16 {
		for i := l; i > l-16; i-- {
			if plantText[i-1] != byte(0) {
				return plantText[0:i]
			}
		}
	}

	return plantText
}
