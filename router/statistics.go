package router

import (
	"github.com/gin-gonic/gin"
	"trojan-panel/api"
)

func initStatisticsRouter(trojanApi *gin.RouterGroup) {
	statistics := trojanApi.Group("/statistics")
	{
		statistics.GET("/userStats", api.UserStats)
		statistics.GET("/trafficStats", api.TrafficStats)
		statistics.GET("/nodeStats", api.NodeStats)
		statistics.GET("/protocolStats", api.ProtocolStats)
	}
}
