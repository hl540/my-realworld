package util

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/cast"
)

type JWT struct {
	token     *jwt.Token
	data      jwt.MapClaims
	secretKey string
}

func NewJwtByData(secretKey string, data map[string]interface{}) *JWT {
	return &JWT{
		data:      data,
		secretKey: secretKey,
	}
}

func NewJwtByToken(secretKey, tokenStr string) (*JWT, error) {
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
	data, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return &JWT{
		token:     token,
		data:      data,
		secretKey: secretKey,
	}, nil
}

func (j *JWT) GetInt(key string) int {
	value, ok := j.data[key]
	if !ok {
		return 0
	}
	return cast.ToInt(value)
}

func (j *JWT) GetString(key string) string {
	value, ok := j.data[key]
	if !ok {
		return ""
	}
	return cast.ToString(value)
}

func (j *JWT) Token() (string, error) {
	j.token = jwt.NewWithClaims(jwt.SigningMethodHS256, j.data)
	return j.token.SignedString([]byte(j.secretKey))
}
