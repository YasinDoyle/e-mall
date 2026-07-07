package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/YasinDoyle/e-mall/repository/db/dao"
	"github.com/YasinDoyle/e-mall/utils/ctl"
	"github.com/YasinDoyle/e-mall/utils/e"
)

// AdminAuthMiddleware 管理员鉴权，须在 AuthMiddleware 之后挂载
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		u, err := ctl.GetUserInfo(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": e.ErrorAuthCheckTokenFail,
				"msg":    "请先登录",
			})
			c.Abort()
			return
		}
		user, err := dao.NewUserDao(c.Request.Context()).GetUserById(u.Id)
		if err != nil || !user.IsAdmin {
			c.JSON(http.StatusOK, gin.H{
				"status": e.ErrorAuthCheckTokenFail,
				"msg":    "无管理员权限",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
