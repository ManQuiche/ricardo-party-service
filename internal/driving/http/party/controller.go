package party

import "github.com/gin-gonic/gin"

type Controller interface {
	Create(gtx *gin.Context)
}
