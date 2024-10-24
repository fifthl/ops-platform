package utilModel

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"github.com/astaxie/beego"
	"log"
)

/*
AK加密/解密
*/

func pKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// AES解密
func aesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = pKCS7UnPadding(origData)
	return origData, nil
}

func Decrypt(id, secret string) (string, string) {
	decodeId, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		log.Println("读取ID失败: ", err.Error())
	}

	decodeSecret, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		log.Println("读取Secret失败: ", err.Error())
	}

	ID, err := aesDecrypt(decodeId, []byte(beego.AppConfig.String("decodePasswd")))
	if err != nil {
		panic(err)
	}

	Secret, err := aesDecrypt(decodeSecret, []byte(beego.AppConfig.String("decodePasswd")))
	if err != nil {
		panic(err)
	}

	return string(ID), string(Secret)
}
