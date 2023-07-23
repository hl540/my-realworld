package util

import (
	"context"
	"crypto/md5"
	"io"
)

// SetContext 设置上下文value
func SetContext(ctx context.Context, key interface{}, value interface{}) context.Context {
	return context.WithValue(ctx, key, value)
}

// GetContext 获取上下文value
func GetContext(ctx context.Context, key interface{}) interface{} {
	return ctx.Value(key)
}

// MakePassword 生成加密后的密码
func MakePassword(oldPassword string, secretKey string) string {
	h := md5.New()
	io.WriteString(h, oldPassword+secretKey)
	return string(h.Sum(nil))
}
