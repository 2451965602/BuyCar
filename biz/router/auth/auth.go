package auth

import (
	"buycar/biz/middleware"
	"buycar/biz/pack"
	"buycar/pkg/errno"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

func Auth() []app.HandlerFunc {
	return append(make([]app.HandlerFunc, 0),
		DoubleTokenAuthFunc(),
	)
}

func DoubleTokenAuthFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if !middleware.IsAccessTokenAvailable(ctx, c) {
			if !middleware.IsRefreshTokenAvailable(ctx, c) {
				pack.BuildFailResponse(c, errno.AuthInvalid)
				c.Abort()

				return
			}
			middleware.GenerateAccessToken(c)
		}

		c.Next(ctx)
	}
}
