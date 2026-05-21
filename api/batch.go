package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"trojan-panel/model/dto"
	"trojan-panel/model/vo"
	"trojan-panel/service"
)

type BatchAccountReq struct {
	IDs    []uint `json:"ids" form:"ids" binding:"required,min=1"`
	Action string `json:"action" form:"action" binding:"required,oneof=enable disable delete extend reset"`
	Days   *int   `json:"days" form:"days"`
}

// BatchAccount 批量账号操作
func BatchAccount(c *gin.Context) {
	var req BatchAccountReq
	if err := c.ShouldBindJSON(&req); err != nil {
		vo.Fail(err.Error(), c)
		return
	}

	batchDto := dto.BatchOperationDto{
		IDs:    req.IDs,
		Action: req.Action,
		Days:   req.Days,
	}

	succeed, failed, err := service.BatchOperation(batchDto)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}

	vo.Success(gin.H{
		"succeed": succeed,
		"failed":  failed,
		"total":   len(req.IDs),
		"message": fmt.Sprintf("成功 %d 个，失败 %d 个", succeed, failed),
	}, c)
}
