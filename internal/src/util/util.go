package util

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt/v4"
	"strings"
)

type authKey struct{}

// NewContext put auth info into context
func NewContext(ctx context.Context, info jwt.Claims) context.Context {
	return context.WithValue(ctx, authKey{}, info)
}

// FromContext extract auth info from context
func FromContext(ctx context.Context) (token jwt.MapClaims, ok bool) {
	token, ok = ctx.Value(authKey{}).(jwt.MapClaims)
	return
}

const UserName = "username"

func GetUserNameFromContext(ctx context.Context) string {
	auth, ok := FromContext(ctx)
	if !ok {
		return ""
	}
	userName, ok := auth[UserName]
	if !ok {
		return ""
	}
	return userName.(string)
}

func ParseJwtToken(ctx context.Context) string {
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

// ParseJwt 解析jwt
func ParseJwt(tokenStr string, secretKey string) (jwt.MapClaims, error) {
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
