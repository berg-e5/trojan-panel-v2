package api

import (
	"github.com/gin-gonic/gin"
	"trojan-panel/model/constant"
	"trojan-panel/model/vo"
	"trojan-panel/service"
)

// SelectAuditLogPage 分页查询审计日志
func SelectAuditLogPage(c *gin.Context) {
	var req struct {
		Username   *string `form:"username"`
		Action     *string `form:"action"`
		TargetType *string `form:"targetType"`
		StartTime  *uint   `form:"startTime"`
		EndTime    *uint   `form:"endTime"`
		PageNum    *uint   `form:"pageNum" binding:"required"`
		PageSize   *uint   `form:"pageSize" binding:"required"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		vo.Fail(constant.ValidateFailed, c)
		return
	}

	logs, total, err := service.SelectAuditLogPage(
		req.Username, req.Action, req.TargetType,
		req.StartTime, req.EndTime,
		req.PageNum, req.PageSize,
	)
	if err != nil {
		vo.Fail(err.Error(), c)
		return
	}

	vo.Success(gin.H{
		"logs":  logs,
		"total": total,
	}, c)
}
