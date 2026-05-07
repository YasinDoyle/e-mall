package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/utils/ctl"
	"github.com/YasinDoyle/e-mall/utils/e"
	jwt "github.com/YasinDoyle/e-mall/utils/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		code = e.SUCCESS
		accessToken := c.GetHeader("access_token")
		refreshToken := c.GetHeader("refresh_token")
		if accessToken == "" {
			code = e.InvalidParams
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": e.GetMsg(code),
				"data":    "Token不能为空",
			})
			c.Abort()
			return
		}
		newAccessToken, newRefreshToken, err := jwt.ParseRefreshToken(accessToken, refreshToken)
		if err != nil {
			code = e.ErrorAuthCheckTokenFail
		}
		if code != e.SUCCESS {
			c.JSON(http.StatusOK, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
				"data":   "鉴权失败",
				"error":  err.Error(),
			})
			c.Abort()
			return
		}
		claims, err := jwt.ParseToken(newAccessToken)
		if err != nil {
			code = e.ErrorAuthCheckTokenFail
			c.JSON(http.StatusOK, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
				"data":   err.Error(),
			})
			c.Abort()
			return
		}
		SetToken(c, newAccessToken, newRefreshToken)
		c.Request = c.Request.WithContext(ctl.NewContext(c.Request.Context(), &ctl.UserInfo{Id: claims.ID}))
		ctl.InitUserInfo(c.Request.Context())
		c.Next()
	}
}

func SetToken(c *gin.Context, accessToken, refreshToken string) {
	secure := IsHttps(c)
	c.Header(consts.AccessTokenHeader, accessToken)
	c.Header(consts.RefreshTokenHeader, refreshToken)
	c.SetCookie(consts.AccessTokenHeader, accessToken, consts.MaxAge, "/", "", secure, true)
	c.SetCookie(consts.RefreshTokenHeader, refreshToken, consts.MaxAge, "/", "", secure, true)
}

func IsHttps(c *gin.Context) bool {
	if c.GetHeader(consts.HeaderForwardedProto) == "https" || c.Request.TLS != nil {
		return true
	}
	return false
}
