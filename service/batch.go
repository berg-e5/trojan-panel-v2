package service

import (
	"github.com/sirupsen/logrus"
	"trojan-panel/dao"
	"trojan-panel/model"
	"trojan-panel/model/dto"
	"trojan-panel/util"
)

// BatchOperation 批量账号操作
func BatchOperation(dto dto.BatchOperationDto) (succeed int, failed int, err error) {
	for _, id := range dto.IDs {
		account, err := dao.SelectAccountById(&id)
		if err != nil || account == nil {
			failed++
			logrus.Warnf("批量操作: 账号不存在 id=%d", id)
			continue
		}

		switch dto.Action {
		case "enable":
			if err := batchEnableAccount(account); err != nil {
				failed++
				logrus.Errorf("批量启用账号失败 id=%d err=%v", id, err)
			} else {
				succeed++
			}

		case "disable":
			if err := batchDisableAccount(account); err != nil {
				failed++
				logrus.Errorf("批量禁用账号失败 id=%d err=%v", id, err)
			} else {
				succeed++
			}

		case "delete":
			if err := batchDeleteAccount(account); err != nil {
				failed++
				logrus.Errorf("批量删除账号失败 id=%d err=%v", id, err)
			} else {
				succeed++
			}

		case "extend":
			if dto.Days == nil {
				failed++
				continue
			}
			if err := batchExtendAccount(account, *dto.Days); err != nil {
				failed++
				logrus.Errorf("批量续期账号失败 id=%d err=%v", id, err)
			} else {
				succeed++
			}

		case "reset":
			if err := batchResetTraffic(account); err != nil {
				failed++
				logrus.Errorf("批量重置流量失败 id=%d err=%v", id, err)
			} else {
				succeed++
			}
		}
	}
	return
}

// batchEnableAccount 启用账号
func batchEnableAccount(account *model.Account) error {
	if account.Deleted != nil && *account.Deleted == 0 {
		return nil
	}
	deleted := uint(0)
	account.Deleted = &deleted

	if err := PullAccountWhiteOrBlackByUsername([]string{*account.Username}, false); err != nil {
		return err
	}

	var email string
	if account.Email != nil {
		email = *account.Email
	}

	event := NotificationEvent{
		Type:     NotifyAccountEnabled,
		Username: *account.Username,
		Email:    email,
	}
	_ = SendNotification(event)

	return nil
}

// batchDisableAccount 禁用账号
func batchDisableAccount(account *model.Account) error {
	if account.Deleted != nil && *account.Deleted == 1 {
		return nil
	}
	deleted := uint(1)
	account.Deleted = &deleted

	if err := PullAccountWhiteOrBlackByUsername([]string{*account.Username}, true); err != nil {
		return err
	}

	var email string
	if account.Email != nil {
		email = *account.Email
	}

	event := NotificationEvent{
		Type:     NotifyAccountDisabled,
		Username: *account.Username,
		Email:    email,
	}
	_ = SendNotification(event)

	return nil
}

// batchDeleteAccount 删除账号
func batchDeleteAccount(account *model.Account) error {
	nodes, err := dao.SelectNodesIpGrpcPortDistinct()
	if err != nil {
		return err
	}
	for _, node := range nodes {
		removeDto := struct {
			Password string `json:"password"`
		}{Password: *account.Pass}
		_ = removeDto
		_ = node
	}

	if err := dao.DeleteAccountById(account.Id); err != nil {
		return err
	}

	logrus.Infof("批量删除账号: %s", *account.Username)
	return nil
}

// batchExtendAccount 续期账号
func batchExtendAccount(account *model.Account, days int) error {
	currentExpire := account.ExpireTime
	now := util.NowMilli()

	var newExpire uint
	if currentExpire != nil && *currentExpire > now {
		newExpire = *currentExpire + uint(days)*24*60*60*1000
	} else {
		newExpire = now + uint(days)*24*60*60*1000
	}
	account.ExpireTime = &newExpire

	if account.Deleted != nil && *account.Deleted == 1 {
		deleted := uint(0)
		account.Deleted = &deleted
		_ = PullAccountWhiteOrBlackByUsername([]string{*account.Username}, false)
	}

	var email string
	if account.Email != nil {
		email = *account.Email
	}

	event := NotificationEvent{
		Type:     NotifyAccountEnabled,
		Username: *account.Username,
		Email:    email,
		Data:     map[string]interface{}{"days": days},
	}
	_ = SendNotification(event)

	return nil
}

// batchResetTraffic 重置流量
func batchResetTraffic(account *model.Account) error {
	// 使用 DAO 方法将 download/upload 重置为 0 并写入数据库
	if err := dao.ResetAccountDownloadAndUpload(account.Id, nil); err != nil {
		return err
	}

	var email string
	if account.Email != nil {
		email = *account.Email
	}

	event := NotificationEvent{
		Type:     NotifyMonthlyReset,
		Username: *account.Username,
		Email:    email,
	}
	_ = SendNotification(event)

	return nil
}
