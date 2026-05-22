package api

import (
	"github.com/gin-gonic/gin"
	"trojan-panel/model/constant"
	"trojan-panel/model/vo"
	"trojan-panel/service"
)

// RefreshToken 刷新 Access Token
func RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}

	newAccessToken, newRefreshToken, err := service.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}

	vo.Success(vo.AccountLoginVo{
		Token:        newAccessToken,
		RefreshToken: newRefreshToken,
	}, c)
}
