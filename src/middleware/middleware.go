package middleware

import (
	"github.com/gin-gonic/gin"
)

func UserMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetHeader("X-User-Id")
		if userId == "" {
			ctx.Next()
			return
		}

		ctx.Set("userId", userId)
		ctx.Next()
	}
}
