package helper

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

func HashPassword(password, salt string) string {

	combined := password + salt // Combine password and salt
	// Create MD5 hash
	hasher := md5.New()
	hasher.Write([]byte(combined))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	return hashedPassword // Return the hashed password and the salt
}

// AES 加密和解密函数

// Encrypt 加密函数，使用 AES 加密
func Encrypt(plainText, key string) (string, error) {
	// 转换密钥为字节数组
	keyBytes := []byte(key)
	if len(keyBytes) != 32 {
		return "", fmt.Errorf("密钥必须是 32 字节长度")
	}

	// 创建 AES 块加密器
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// 使用随机生成的 IV（初始化向量）
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// 加密数据
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], []byte(plainText))

	// 返回 Base64 编码的加密结果
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt 解密函数，使用 AES 解密
func Decrypt(encryptedText, key string) (string, error) {
	// 转换密钥为字节数组
	keyBytes := []byte(key)
	if len(keyBytes) != 32 {
		return "", fmt.Errorf("密钥必须是 32 字节长度")
	}

	// 解码 Base64 编码的密文
	cipherText, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	// 检查密文长度
	if len(cipherText) < aes.BlockSize {
		return "", fmt.Errorf("密文太短")
	}

	// 分离 IV 和实际密文
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	// 创建 AES 块加密器
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// 解密数据
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)
	plainText := string(bytes.Trim(cipherText, "\x00"))

	// 返回解密结果
	return string(plainText), nil
}
