package util

import (
	"context"
	"github.com/go-kratos/kratos/v2/transport"
	"strings"
)

const UserID = "user_id"
const UserName = "user_name"
const UserEmail = "user_email"

type AuthKey struct{}

type UserInfo struct {
	UserID    uint
	UserName  string
	UserEmail string
}

// GetUserInfo 通过上下文获取username
func GetUserInfo(ctx context.Context) *UserInfo {
	user, ok := GetContext(ctx, AuthKey{}).(*UserInfo)
	if !ok {
		return nil
	}
	return user
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
