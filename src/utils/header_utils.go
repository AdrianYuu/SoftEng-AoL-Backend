package utils

import "github.com/gin-gonic/gin"

func GetCurrentUserID(ctx *gin.Context) string {
	if ctx == nil {
		return ""
	}

	userId, _ := ctx.Get("userId")
	return userId.(string)
}
