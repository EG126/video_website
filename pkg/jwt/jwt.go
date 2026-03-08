package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"time"
	"video_website/pkg/errno"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func Init(secret string) {
	jwtSecret = []byte(secret)
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	// hex.EncodeToString():随机字节切片编码为十六进制字符串
	return hex.EncodeToString(b), nil
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		log.Printf("ParseToken error: %v", err) // 添加此行
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errno.TokenExpired
		}
		return nil, errno.TokenInvalid
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errno.TokenInvalid
}
