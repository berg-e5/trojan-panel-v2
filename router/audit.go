package router

import (
	"github.com/gin-gonic/gin"
	"trojan-panel/api"
)

func initAuditRouter(trojanApi *gin.RouterGroup) {
	audit := trojanApi.Group("/audit")
	{
		audit.GET("/page", api.SelectAuditLogPage)
	}
}
