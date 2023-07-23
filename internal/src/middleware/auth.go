package middleware

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/hl540/my-realworld/internal/src/errors"
	"github.com/hl540/my-realworld/internal/src/util"
)

func NewJwt(jwtKey string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if _, ok := transport.FromServerContext(ctx); ok {
				// 解析header中的token
				jwtStr := util.ParseTokenStr(ctx)
				// 解析jwt
				tokenInfo, err := util.ParseJwtStr(jwtStr, jwtKey)
				if err != nil {
					return nil, errors.NewHTTPError(401, "body", err.Error())
				}
				ctx = util.SetContext(ctx, util.AuthKey{}, tokenInfo)
				return handler(ctx, req)
			}
			return handler(ctx, req)
		}
	}
}
