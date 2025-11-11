package middleware

import (
	"buycar/biz/model/user"
	"buycar/biz/pack"
	"buycar/biz/service"
	"buycar/pkg/constants"
	"buycar/pkg/errno"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hertz-contrib/jwt"
	"github.com/satori/go.uuid"

	"github.com/cloudwego/hertz/pkg/app"
)

var (
	AccessTokenJwtMiddleware  *jwt.HertzJWTMiddleware
	RefreshTokenJwtMiddleware *jwt.HertzJWTMiddleware
)

type JwtCustomClaims struct {
	UserId int64  `json:"userid"`
	UUID   string `json:"uuid"`
}

func AccessTokenJwt() {
	var err error
	AccessTokenJwtMiddleware, err = jwt.New(&jwt.HertzJWTMiddleware{
		Realm:                       "LS",
		Key:                         []byte("AccessToken_key"),
		Timeout:                     12 * time.Hour,
		MaxRefresh:                  12 * time.Hour,
		WithoutDefaultTokenHeadName: true,
		TokenLookup:                 "header: Authorization",
		IdentityKey:                 constants.IdentityKey,

		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*JwtCustomClaims); ok {
				return jwt.MapClaims{
					AccessTokenJwtMiddleware.IdentityKey: v.UserId,
					constants.TokenType:                  "access",
					constants.UUID:                       v.UUID,
				}
			}
			return jwt.MapClaims{}
		},

		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			resp := &JwtCustomClaims{
				UserId: int64(claims[RefreshTokenJwtMiddleware.IdentityKey].(float64)),
				UUID:   claims[constants.UUID].(string),
			}
			return resp
		},

		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			pack.BuildFailResponse(c, errno.AuthNoToken)
		},

		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			c.Set("Access-Token", token)
		},

		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			var loginStruct user.LoginReq
			if err := c.BindAndValidate(&loginStruct); err != nil {
				return nil, err
			}
			users, err := service.NewUserService(ctx, c).LoginIn(&loginStruct)
			if err != nil {
				return nil, err
			}
			c.Set(constants.ContextUid, users.UserID)
			claims := &JwtCustomClaims{
				UserId: users.UserID,
				UUID:   uuid.NewV1().String(),
			}
			return claims, nil
		},
	})
	if err != nil {
		log.Fatal("JWT 错误：" + err.Error())
	}
}

func RefreshTokenJwt() {
	var err error
	RefreshTokenJwtMiddleware, err = jwt.New(&jwt.HertzJWTMiddleware{
		Realm:                       "LS",
		Key:                         []byte("refresh_token_key"),
		Timeout:                     time.Hour * 72,
		WithoutDefaultTokenHeadName: true,
		TokenLookup:                 "header: Refresh-Token",
		IdentityKey:                 constants.IdentityKey,

		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*JwtCustomClaims); ok {
				return jwt.MapClaims{
					AccessTokenJwtMiddleware.IdentityKey: v.UserId,
					constants.TokenType:                  "refresh",
					constants.UUID:                       v.UUID,
				}
			}
			return jwt.MapClaims{}
		},

		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			resp := &JwtCustomClaims{
				UserId: int64(claims[RefreshTokenJwtMiddleware.IdentityKey].(float64)),
				UUID:   claims[constants.UUID].(string),
			}
			return resp
		},

		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			pack.BuildFailResponse(c, errno.AuthNoToken)
		},

		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			c.Set("Refresh-Token", token)
		},

		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			userId := service.GetUidFormContext(c)
			uuidStr := service.GetUuidFormContext(c)

			claims := &JwtCustomClaims{
				UserId: userId,
				UUID:   uuidStr,
			}

			return claims, nil

		},
	})
	if err != nil {
		log.Fatal("JWT 错误：" + err.Error())
	}
}

func GenerateAccessToken(c *app.RequestContext) {

	userId := service.GetUidFormContext(c)
	uuidStr := service.GetUuidFormContext(c)

	data := &JwtCustomClaims{
		UserId: userId,
		UUID:   uuidStr,
	}

	tokenString, _, _ := AccessTokenJwtMiddleware.TokenGenerator(data)
	c.Header("New-Access-Token", tokenString)

}

func IsAccessTokenAvailable(ctx context.Context, c *app.RequestContext) bool {
	claims, err := AccessTokenJwtMiddleware.GetClaimsFromJWT(ctx, c)
	if err != nil {
		return false
	}
	// 验证token类型是否为access
	if tokenType, ok := claims[constants.TokenType].(string); !ok || tokenType != "access" {
		return false
	}

	switch v := claims["exp"].(type) {
	case nil:
		return false
	case float64:
		if int64(v) < AccessTokenJwtMiddleware.TimeFunc().Unix() {
			return false
		}
	case json.Number:
		n, err := v.Int64()
		if err != nil {
			return false
		}
		if n < AccessTokenJwtMiddleware.TimeFunc().Unix() {
			return false
		}
	default:
		return false
	}
	c.Set("JWT_PAYLOAD", claims)
	identity := AccessTokenJwtMiddleware.IdentityHandler(ctx, c)
	if identity != nil {
		c.Set(constants.IdentityKey, identity.(*JwtCustomClaims).UserId)
		c.Set(constants.UUID, identity.(*JwtCustomClaims).UUID)
	}
	if !AccessTokenJwtMiddleware.Authorizator(identity, ctx, c) {
		return false
	}

	return true

}

func IsRefreshTokenAvailable(ctx context.Context, c *app.RequestContext) bool {
	claims, err := RefreshTokenJwtMiddleware.GetClaimsFromJWT(ctx, c)
	if err != nil {
		return false
	}

	// 验证token类型是否为refresh
	if tokenType, ok := claims[constants.TokenType].(string); !ok || tokenType != "refresh" {
		return false
	}

	switch v := claims["exp"].(type) {
	case nil:
		return false
	case float64:
		if int64(v) < RefreshTokenJwtMiddleware.TimeFunc().Unix() {
			return false
		}
	case json.Number:
		n, err := v.Int64()
		if err != nil {
			return false
		}
		if n < RefreshTokenJwtMiddleware.TimeFunc().Unix() {
			return false
		}
	default:
		return false
	}

	c.Set("JWT_PAYLOAD", claims)
	identity := RefreshTokenJwtMiddleware.IdentityHandler(ctx, c)
	if identity != nil {
		c.Set(constants.IdentityKey, identity.(*JwtCustomClaims).UserId)
		c.Set(constants.UUID, identity.(*JwtCustomClaims).UUID)
	}
	if !RefreshTokenJwtMiddleware.Authorizator(identity, ctx, c) {
		return false
	}

	return true
}

func InitJWT() {
	AccessTokenJwt()
	RefreshTokenJwt()
	errInit := AccessTokenJwtMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("AccessTokenJwtMiddleware.MiddlewareInit() 错误：" + errInit.Error())
	}

	errInit = RefreshTokenJwtMiddleware.MiddlewareInit()
	if errInit != nil {
		log.Fatal("RefreshTokenJwtMiddleware.MiddlewareInit() 错误：" + errInit.Error())
	}
}
