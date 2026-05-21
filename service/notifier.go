package service

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"trojan-panel/model/constant"
	"trojan-panel/model/dto"
	"trojan-panel/model/vo"
)

// NotificationType 通知类型
type NotificationType string

const (
	NotifyAccountCreated     NotificationType = "account_created"
	NotifyAccountDisabled    NotificationType = "account_disabled"
	NotifyAccountEnabled     NotificationType = "account_enabled"
	NotifyAccountExpireWarn NotificationType = "account_expire_warn"
	NotifyAccountExpired     NotificationType = "account_expired"
	NotifyTrafficWarning80  NotificationType = "traffic_warning_80"
	NotifyTrafficWarning100 NotificationType = "traffic_warning_100"
	NotifyNodeOnline         NotificationType = "node_online"
	NotifyNodeOffline        NotificationType = "node_offline"
	NotifySystemAlert        NotificationType = "system_alert"
	NotifyMonthlyReset       NotificationType = "monthly_reset"
)

// NotificationEvent 通知事件
type NotificationEvent struct {
	Type       NotificationType
	Username   string
	Email      string
	TelegramID string
	Data       map[string]interface{}
}

// SendNotification 统一发送通知（邮件 + Telegram）
func SendNotification(event NotificationEvent) error {
	name := constant.SystemName
	systemVo, err := SelectSystemByName(&name)
	if err != nil {
		return err
	}

	title, content := buildNotificationContent(event)
	var lastErr error

	// 邮件通知
	if systemVo.EmailEnable == 1 && event.Email != "" {
		if err := sendEmailNotification(event.Email, title, content); err != nil {
			logrus.Errorf("邮件通知发送失败: %v", err)
			lastErr = err
		} else {
			logrus.Infof("邮件通知发送成功: %s -> %s", event.Type, event.Email)
		}
	}

	// Telegram 通知（仅当系统配置了 TelegramBotToken 时启用）
	if hasTelegramEnabled(&systemVo) && event.TelegramID != "" {
		if err := sendTelegramNotification(event.TelegramID, title+"\n\n"+content); err != nil {
			logrus.Errorf("Telegram通知发送失败: %v", err)
			lastErr = err
		} else {
			logrus.Infof("Telegram通知发送成功: %s -> %s", event.Type, event.TelegramID)
		}
	}

	return lastErr
}

// buildNotificationContent 根据事件类型构造消息
func buildNotificationContent(event NotificationEvent) (title, content string) {
	switch event.Type {
	case NotifyAccountCreated:
		title = "✅ 账号创建成功"
		content = fmt.Sprintf("账号: %s\n欢迎使用！请及时查看您的节点信息。", event.Username)

	case NotifyAccountDisabled:
		title = "⚠️ 账号已被禁用"
		content = fmt.Sprintf("账号: %s\n原因: 到期或违规，请续期后联系管理员。", event.Username)

	case NotifyAccountEnabled:
		title = "✅ 账号已恢复"
		content = fmt.Sprintf("账号: %s\n您的账号已恢复正常使用。", event.Username)

	case NotifyAccountExpireWarn:
		days := 0
		if v, ok := event.Data["days"].(int); ok {
			days = v
		}
		title = "📅 账号到期提醒"
		content = fmt.Sprintf("账号: %s\n还有 %d 天到期，请及时续期！", event.Username, days)

	case NotifyAccountExpired:
		title = "❌ 账号已到期"
		content = fmt.Sprintf("账号: %s\n已到期，账号已被禁用，请联系管理员。", event.Username)

	case NotifyTrafficWarning80:
		title = "📊 流量使用提醒"
		content = fmt.Sprintf("账号: %s\n流量已使用超过 80%%，请注意流量使用情况。", event.Username)

	case NotifyTrafficWarning100:
		title = "🚫 流量已用尽"
		content = fmt.Sprintf("账号: %s\n流量已用尽，账号已被暂停，请联系管理员续期。", event.Username)

	case NotifyNodeOnline:
		name := ""
		if v, ok := event.Data["nodeName"].(string); ok {
			name = v
		}
		title = "🟢 节点上线"
		content = fmt.Sprintf("节点: %s 已恢复在线", name)

	case NotifyNodeOffline:
		name := ""
		if v, ok := event.Data["nodeName"].(string); ok {
			name = v
		}
		title = "🔴 节点离线"
		content = fmt.Sprintf("节点: %s 已离线，请检查服务状态", name)

	case NotifySystemAlert:
		msg := ""
		if v, ok := event.Data["message"].(string); ok {
			msg = v
		}
		title = "🚨 系统告警"
		content = fmt.Sprintf("告警内容: %s", msg)

	case NotifyMonthlyReset:
		title = "📅 月度流量重置"
		content = fmt.Sprintf("账号: %s\n本月流量已重置，上月使用量已清零。", event.Username)

	default:
		title = "通知"
		content = fmt.Sprintf("账号: %s\n类型: %s", event.Username, event.Type)
	}
	return
}

// hasTelegramEnabled 检测系统是否配置了 Telegram
// TODO: 正式上线时替换为实际的字段检查
// 目前返回 false，正式对接时需要数据库增加 telegram_enable / telegram_bot_token 字段
func hasTelegramEnabled(systemVo *vo.SystemVo) bool {
	return false
}

// sendEmailNotification 发送邮件
func sendEmailNotification(toEmail, subject, content string) error {
	name := constant.SystemName
	systemVo, err := SelectSystemByName(&name)
	if err != nil {
		return err
	}
	if systemVo.EmailEnable == 0 {
		return fmt.Errorf("邮件功能未启用")
	}

	emailDto := dto.SendEmailDto{
		FromEmailName: systemVo.SystemName,
		ToEmails:      []string{toEmail},
		Subject:       subject,
		Content:       content,
	}
	return SendEmail(&emailDto)
}

// sendTelegramNotification 发送 Telegram 消息
func sendTelegramNotification(chatID, text string) error {
	return SendTelegramMessage(chatID, text)
}
