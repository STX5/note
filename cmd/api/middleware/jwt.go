package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"note/cmd/api/rpc"
	"note/hertz_gen/api"
	"note/kitex_gen/user"
	"note/pkg/consts"
	"note/pkg/errno"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	"github.com/hertz-contrib/jwt"
)

var JwtMiddleware *jwt.HertzJWTMiddleware

func InitJWT() {
	JwtMiddleware, _ = jwt.New(&jwt.HertzJWTMiddleware{
		Key:           []byte(consts.SecretKey),
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour,
		// 检索身份的键，默认为 identity (identity是key，value就是用户ID)
		IdentityKey: consts.IdentityKey,
		// 获取身份信息的函数，与 IdentityKey 配合使用
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			userid, _ := claims[consts.IdentityKey].(json.Number).Int64()
			return &api.User{
				UserID: userid,
			}
		},
		// 登陆成功后为向 token 中添加自定义负载信息
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(int64); ok {
				return jwt.MapClaims{
					consts.IdentityKey: v,
				}
			}
			return jwt.MapClaims{}
		},
		// 配合 HertzJWTMiddleware.LoginHandler 使用，
		// 登录时触发，用于认证用户的登录信息
		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			var err error
			var req api.CheckUserRequest
			if err = c.BindAndValidate(&req); err != nil {
				return "", jwt.ErrInvalidAuthHeader
			}
			if len(req.Username) == 0 || len(req.Password) == 0 {
				return "", jwt.ErrMissingLoginValues
			}
			return rpc.CheckUser(context.Background(), &user.CheckUserRequest{
				Username: req.Username,
				Password: req.Password,
			})
		},
		// 登录的响应函数，作为 LoginHandler 的响应结果
		// // 在 LoginHandler 内调用
		// // h.POST("/login", authMiddleware.LoginHandler)
		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, utils.H{
				"code":   errno.Success.ErrCode,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		},
		// jwt 验证流程失败的响应
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(http.StatusOK, utils.H{
				"code":    errno.AuthorizationFailedErr.ErrCode,
				"message": message,
			})
		},
		// jwt 校验流程发生错误时响应所包含的错误信息
		HTTPStatusMessageFunc: func(e error, ctx context.Context, c *app.RequestContext) string {
			switch t := e.(type) {
			case errno.ErrNo:
				return t.ErrMsg
			default:
				return t.Error()
			}
		},
		ParseOptions: []jwtv4.ParserOption{jwtv4.WithJSONNumber()},
	})
}
