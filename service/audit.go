package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"trojan-panel/dao"
	"trojan-panel/model"
)

// WriteAudit 记录审计日志
func WriteAudit(c *gin.Context, action, targetType string, targetId *uint, detail interface{}) {
	account := GetCurrentAccount(c)

	userId := account.Id
	username := account.Username
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	var detailStr string
	if detail != nil {
		bytes, _ := json.Marshal(detail)
		detailStr = string(bytes)
	}

	log := &model.AuditLog{
		UserId:     &userId,
		Username:   &username,
		Action:     &action,
		TargetType: &targetType,
		TargetId:   targetId,
		Detail:     &detailStr,
		Ip:         &ip,
		UserAgent:  &userAgent,
	}

	go func() {
		_ = dao.InsertAuditLog(log)
	}()
}

// SelectAuditLogPage 查询审计日志（分页）
func SelectAuditLogPage(
	username *string,
	action *string,
	targetType *string,
	startTime *uint,
	endTime *uint,
	pageNum *uint,
	pageSize *uint,
) ([]model.AuditLog, int, error) {
	return dao.SelectAuditLogPage(username, action, targetType, startTime, endTime, pageNum, pageSize)
}
