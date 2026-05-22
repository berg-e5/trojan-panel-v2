package model

import "time"

// AuditLog 审计日志
type AuditLog struct {
	Id         *uint      `ddb:"id" json:"id"`
	UserId     *uint      `ddb:"user_id" json:"userId"`
	Username   *string    `ddb:"username" json:"username"`
	Action     *string    `ddb:"action" json:"action"`
	TargetType *string    `ddb:"target_type" json:"targetType"`
	TargetId   *uint      `ddb:"target_id" json:"targetId"`
	Detail     *string    `ddb:"detail" json:"detail"`
	Ip         *string    `ddb:"ip" json:"ip"`
	UserAgent  *string    `ddb:"user_agent" json:"userAgent"`
	CreatedAt  *time.Time `ddb:"created_at" json:"createdAt"`
}

// AuditAction 操作类型常量
const (
	AuditActionCreate  = "create"
	AuditActionUpdate  = "update"
	AuditActionDelete  = "delete"
	AuditActionEnable  = "enable"
	AuditActionDisable = "disable"
	AuditActionReset   = "reset"
	AuditActionLogin   = "login"
	AuditActionLogout  = "logout"
)

// AuditTargetType 目标类型常量
const (
	AuditTargetAccount = "account"
	AuditTargetNode    = "node"
	AuditTargetSystem  = "system"
)
