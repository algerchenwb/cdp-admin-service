package helper

import (
	"cdp-admin-service/internal/helper/cache"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var secretKey = []byte("0004c9ed-d901-475e-a955-dbde3609995f") // 密钥

// 自定义 Claims
type CustomClaims struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	RoleId   int64  `json:"role_id"`
	Platform int64  `json:"platform"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// 生成 Token
func GenerateToken(userId int64, userName string, roleId int64, platform int64, isAdmin bool, c *cache.Cache, timeouts int) (string, error) {
	claims := CustomClaims{
		UserID:   userId,
		UserName: userName,
		RoleId:   roleId,
		Platform: platform,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(48 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// 解析 Token
func ParseToken(tokenString string) (*CustomClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	// 验证 Token 有效性并提取 Claims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
