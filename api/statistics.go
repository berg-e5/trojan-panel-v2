package api

import (
	"github.com/gin-gonic/gin"
	"trojan-panel/model/vo"
	"trojan-panel/service"
)

func UserStats(c *gin.Context) {
	stats, err := service.GetUserStats()
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(stats, c)
}

func TrafficStats(c *gin.Context) {
	stats, err := service.GetTrafficStats()
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(stats, c)
}

func NodeStats(c *gin.Context) {
	stats, err := service.GetNodeStats()
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(stats, c)
}

func ProtocolStats(c *gin.Context) {
	stats, err := service.GetProtocolStats()
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}
	vo.Success(stats, c)
}
