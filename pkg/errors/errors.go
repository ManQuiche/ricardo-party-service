package errors

import (
	"github.com/gin-gonic/gin"
)

func GinErrorHandler(gtx *gin.Context, err error, code int) bool {
	if err != nil {
		_ = gtx.Error(err)
		gtx.AbortWithStatusJSON(code, gin.H{"status": false, "message": err.Error()})
		return true
	}
	return false
}
