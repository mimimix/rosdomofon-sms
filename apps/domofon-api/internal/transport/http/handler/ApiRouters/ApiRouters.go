package ApiRouters

import (
	"github.com/gin-gonic/gin"
)

type ApiRouters struct {
	Public  *gin.RouterGroup
	Private *gin.RouterGroup
}

func CreateApiRoutes(gin *gin.Engine) *ApiRouters {
	gin.MaxMultipartMemory = 1 << 20
	publicRoute := gin.Group("/api")

	return &ApiRouters{
		Public: publicRoute,
	}
}
