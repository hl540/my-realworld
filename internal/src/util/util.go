package util

import (
	"context"
	"crypto/md5"
	"encoding/hex"
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

func MD5(str string) string {
	w := md5.New()
	io.WriteString(w, str)
	bw := w.Sum(nil)
	return hex.EncodeToString(bw)
}

func IntDefault(i int, d int) int {
	if i == 0 {
		return d
	}
	return i
}
