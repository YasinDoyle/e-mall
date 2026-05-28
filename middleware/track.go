package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/utils/track"
)

func Jaeger() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, requestCtx, err := track.StartRequestSpan(c.Request.Context(), c.FullPath(), c.Request.Header)
		if err != nil {
			c.Next()
			return
		}
		defer span.Finish()

		c.Request = c.Request.WithContext(requestCtx)
		c.Set(consts.SpanCTX, span.Context())
		c.Next()
	}
}
