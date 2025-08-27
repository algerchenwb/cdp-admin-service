package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/vlog"
)

var _secret []byte = []byte("#HvL%$o0oNNoOZnk#o2qbqCeQB1iXeIR")

func Compare(encryptPassword string, password string, salt string) (bool, error) {
	originPassword, err := Decode(encryptPassword)
	if err != nil {
		return false, err
	}
	md5pass := MD5(originPassword + salt)
	return password == md5pass, nil
}
func Decode(encryptPassword string) (string, error) {
	encrypt, err := base64.StdEncoding.DecodeString(encryptPassword)
	if err != nil {
		//vlog.Errorf("base64.StdEncoding.DecodeString(%v) failure, err[%v]", encryptPassword, err)
		return "", err
	}
	originByte, err := aesDecrypt(encrypt, _secret)
	if err != nil {
		//vlog.Errorf("utils.aesDecrypt failure, err[%v]", err)
		return "", err
	}
	return string(originByte), nil
}

func aesDecrypt(encrypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	vlog.Infof("blockSize:%v", blockSize)
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	originData := make([]byte, len(encrypted))
	blockMode.CryptBlocks(originData, encrypted)
	originData = pKCS7UnPadding(originData)
	return originData, nil
}
func pKCS7UnPadding(originData []byte) []byte {
	length := len(originData)
	unpadding := int(originData[length-1])
	return originData[:(length - unpadding)]
}

func MD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}
