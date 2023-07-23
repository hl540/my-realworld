package util

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"strings"
)

// SetContext 设置上下文value
func SetContext(ctx context.Context, key interface{}, value interface{}) context.Context {
	return context.WithValue(ctx, key, value)
}

// GetContext 获取上下文value
func GetContext(ctx context.Context, key interface{}) interface{} {
	return ctx.Value(key)
}

const UserID = "user_id"

type AuthKey struct{}

// GetUserID 通过上下文获取username
func GetUserID(ctx context.Context) int {
	data, ok := GetContext(ctx, AuthKey{}).(jwt.MapClaims)
	if !ok {
		return 0
	}
	if _, ok := data[UserID]; !ok {
		return 0
	}
	return int(data[UserID].(float64))
}

// ParseTokenStr 从header中解析token
func ParseTokenStr(ctx context.Context) string {
	if tr, ok := transport.FromServerContext(ctx); ok {
		// 解析header
		authStr := tr.RequestHeader().Get("Authorization")
		auths := strings.SplitN(authStr, " ", 2)
		if len(auths) != 2 || auths[0] != "Token" {
			return ""
		}
		return auths[1]
	}
	return ""
}

// ParseJwtStr 解析jwt
func ParseJwtStr(tokenStr string, secretKey string) (jwt.MapClaims, error) {
	// 解析jwt内容
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// 验证jwt算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	// 解析内容
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, jwt.ErrTokenInvalidClaims
	}
}

// MakeJwtString 生成jwt
func MakeJwtString(data jwt.MapClaims, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	return token.SignedString([]byte(secretKey))
}

// MakePassword 生成加密后的密码
func MakePassword(oldPassword string, secretKey string) string {
	h := md5.New()
	io.WriteString(h, oldPassword+secretKey)
	return string(h.Sum(nil))
}
