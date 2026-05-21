package service

// 通知系统集成指南
// ====================
// 本文件展示通知系统在各业务逻辑中的集成方式。
// 由于 Account model 暂时没有 TelegramID 字段，邮件通知可直接接入，
// Telegram 通知需等 model 扩展后启用。
//
// 集成示例如下：

// --------------------------------------------------
// 1. 账号创建时通知（service/account.go CreateAccount）
// --------------------------------------------------
// 在 CreateAccount 成功创建账号后加入：
//
//  event := NotificationEvent{
//      Type:     NotifyAccountCreated,
//      Username: account.Username,
//      Email:    account.Email != nil ? *account.Email : "",
//      TelegramID: account.TelegramID != nil ? *account.TelegramID : "",
//  }
//  _ = SendNotification(event) // 异步发送，不阻塞主流程
//
// --------------------------------------------------
// 2. 账号到期扫描定时任务（service/account.go CronScanAccountExpireWarn）
// --------------------------------------------------
// 在现有到期邮件逻辑中替换为统一通知：
//
//  event := NotificationEvent{
//      Type:     NotifyAccountExpireWarn,
//      Username: account.Username,
//      Email:    *account.Email,
//      Data:     map[string]interface{}{"days": expireWarnDay},
//  }
//  _ = SendNotification(event)
//
// --------------------------------------------------
// 3. 账号被禁用时（service/account.go DisableAccount）
// --------------------------------------------------
//  event := NotificationEvent{
//      Type:       NotifyAccountDisabled,
//      Username:   username,
//      Email:      email,
//      TelegramID: telegramID,
//  }
//  _ = SendNotification(event)
//
// --------------------------------------------------
// 4. 月度流量重置（service/account.go CronResetDownloadAndUploadMonth）
// --------------------------------------------------
//  event := NotificationEvent{
//      Type:     NotifyMonthlyReset,
//      Username: username,
//      Email:    email,
//  }
//  _ = SendNotification(event)
//
// --------------------------------------------------
// 5. 节点上下线（service/node.go 相关函数）
// --------------------------------------------------
//  // 上线
//  event := NotificationEvent{
//      Type: NotifyNodeOnline,
//      Data: map[string]interface{}{"nodeName": nodeName},
//  }
//  _ = SendNotification(event)
//
//  // 离线
//  event := NotificationEvent{
//      Type: NotifyNodeOffline,
//      Data: map[string]interface{}{"nodeName": nodeName},
//  }
//  _ = SendNotification(event)
//
// --------------------------------------------------
// 6. 系统告警（管理员主动触发）
// --------------------------------------------------
//  event := NotificationEvent{
//      Type: NotifySystemAlert,
//      Data: map[string]interface{}{"message": "磁盘使用率超过 90%"},
//  }
//  _ = SendNotification(event)
//
// --------------------------------------------------
// 注意事项：
// - SendNotification 为异步调用，不影响主业务逻辑响应速度
// - Telegram 通知需数据库增加 telegram_enable / telegram_bot_token / telegram_admin_id 字段
// - 邮件通知 EmailEnable 字段已存在于 SystemVo，可直接使用
// --------------------------------------------------
