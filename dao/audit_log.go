package dao

import (
	"fmt"
	"time"

	"github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
	"github.com/sirupsen/logrus"
	"trojan-panel/model"
)

// InsertAuditLog 写入审计日志
func InsertAuditLog(log *model.AuditLog) error {
	data := map[string]interface{}{
		"user_id":     *log.UserId,
		"username":    *log.Username,
		"action":      *log.Action,
		"target_type": *log.TargetType,
		"ip":          *log.Ip,
		"created_at":  time.Now(),
	}
	if log.TargetId != nil {
		data["target_id"] = *log.TargetId
	}
	if log.Detail != nil {
		data["detail"] = *log.Detail
	}
	if log.UserAgent != nil {
		data["user_agent"] = *log.UserAgent
	}

	query, values, err := builder.BuildInsert("audit_log", []map[string]interface{}{data})
	if err != nil {
		logrus.Errorf("build audit log insert err: %v", err)
		return err
	}
	if _, err = db.Exec(query, values...); err != nil {
		logrus.Errorf("insert audit log err: %v", err)
		return err
	}
	return nil
}

// SelectAuditLogPage 分页查询审计日志
func SelectAuditLogPage(
	username *string,
	action *string,
	targetType *string,
	startTime *uint,
	endTime *uint,
	pageNum *uint,
	pageSize *uint,
) ([]model.AuditLog, int, error) {
	where := map[string]interface{}{}
	var args []interface{}

	if username != nil && *username != "" {
		where["username"] = *username
	}
	if action != nil && *action != "" {
		where["action"] = *action
	}
	if targetType != nil && *targetType != "" {
		where["target_type"] = *targetType
	}

	// 拼接 WHERE 子句
	whereClause := ""
	for k, v := range where {
		if whereClause != "" {
			whereClause += " AND "
		}
		whereClause += k + " = ?"
		args = append(args, v)
	}

	// 时间范围
	if startTime != nil && endTime != nil {
		if whereClause != "" {
			whereClause += " AND "
		}
		whereClause += "created_at >= ? AND created_at <= ?"
		args = append(args,
			time.Unix(int64(*startTime)/1000, 0),
			time.Unix(int64(*endTime)/1000, 0),
		)
	}

	if whereClause != "" {
		whereClause = "WHERE " + whereClause
	}

	// 统计总数
	countQuery := "SELECT COUNT(1) FROM audit_log " + whereClause
	var total int
	if err := db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		logrus.Errorf("count audit log err: %v", err)
		return nil, 0, err
	}

	// 分页查询
	offset := (*pageNum - 1) * *pageSize
	query := fmt.Sprintf(
		"SELECT id, user_id, username, action, target_type, target_id, detail, ip, user_agent, created_at FROM audit_log %s ORDER BY id DESC LIMIT ? OFFSET ?",
		whereClause,
	)
	args = append(args, *pageSize, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		logrus.Errorln(err.Error())
		return nil, 0, err
	}
	defer rows.Close()

	var logs []model.AuditLog
	if err = scanner.Scan(rows, &logs); err != nil {
		logrus.Errorln(err.Error())
		return nil, 0, err
	}
	return logs, total, nil
}
