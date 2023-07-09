package server

import (
	"context"

	v1 "github.com/hl540/my-realworld/api/my_realworld/v1"
	"github.com/hl540/my-realworld/internal/conf"
	"github.com/hl540/my-realworld/internal/service"
	"github.com/hl540/my-realworld/internal/src/middleware/auth"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, myRealworld *service.MyRealworldService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			authMiddleware(c),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterMyRealworldHTTPServer(srv, myRealworld)
	return srv
}

// 带白名单的鉴权组件
func authMiddleware(c *conf.Server) middleware.Middleware {
	// auth中间件
	authMiddleware := auth.NewJwt(c.Jwt.GetSecretKey())
	// 白名单
	return selector.Server(authMiddleware).Match(func(ctx context.Context, operation string) bool {
		for _, path := range c.GetJwt().GetWhitePath() {
			if path == operation {
				return false
			}
		}
		return true
	}).Build()
}
