package dao

import (
	"trojan-panel/model"
)

// SelectAllAccounts 获取所有账号（用于统计）
func SelectAllAccounts() ([]model.Account, error) {
	var accounts []model.Account
	rows, err := db.Query("SELECT id, username, pass, email, quota, upload, download, last_login_time, expire_time, deleted, role_id, preset_expire, preset_quota FROM account")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a model.Account
		var username string
		var email, pass *string
		var quota, upload, download *int
		var lastLoginTime *uint
		var expireTime *uint
		var deleted, roleId, presetExpire *uint
		var presetQuota *int

		err := rows.Scan(&a.Id, &username, &pass, &email, &quota, &upload, &download, &lastLoginTime, &expireTime, &deleted, &roleId, &presetExpire, &presetQuota)
		if err != nil {
			continue
		}
		a.Username = &username
		a.Pass = pass
		a.Email = email
		a.Quota = quota
		a.Upload = upload
		a.Download = download
		a.LastLoginTime = lastLoginTime
		a.ExpireTime = expireTime
		a.Deleted = deleted
		a.RoleId = roleId
		a.PresetExpire = presetExpire
		a.PresetQuota = presetQuota
		accounts = append(accounts, a)
	}
	return accounts, nil
}

// SelectAllNodes 获取所有节点（用于统计）
func SelectAllNodes() ([]model.Node, error) {
	var nodes []model.Node
	rows, err := db.Query("SELECT id, node_server_id, name, node_type_id FROM node")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var n model.Node
		var nodeTypeId, nodeServerId *uint
		var name *string

		err := rows.Scan(&n.Id, &nodeServerId, &name, &nodeTypeId)
		if err != nil {
			continue
		}
		n.NodeServerId = nodeServerId
		n.Name = name
		n.NodeTypeId = nodeTypeId
		nodes = append(nodes, n)
	}
	return nodes, nil
}
