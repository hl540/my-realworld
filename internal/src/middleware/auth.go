package middleware

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/hl540/my-realworld/internal/src/errors"
	"github.com/hl540/my-realworld/internal/src/util"
)

func NewJwt(secretKey string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if _, ok := transport.FromServerContext(ctx); ok {
				// 解析header中的token
				tokenStr := util.ParseTokenStr(ctx)
				if tokenStr == "" {
					return nil, errors.NewHTTPError(401, "body", "Authorization token is required")
				}
				// 解析jwt内容，设置上下文
				jwt, err := util.NewJwtByToken(secretKey, tokenStr)
				if err != nil {
					return nil, errors.NewHTTPError(401, "body", err.Error())
				}
				ctx = util.SetContext(ctx, util.AuthKey{}, &util.UserInfo{
					UserID:    int64(jwt.GetInt(util.UserID)),
					UserName:  jwt.GetString(util.UserName),
					UserEmail: jwt.GetString(util.UserEmail),
				})
				return handler(ctx, req)
			}
			return handler(ctx, req)
		}
	}
}
