package service

import (
	"trojan-panel/dao"
	"trojan-panel/model/constant"
	"trojan-panel/util"
)

// UserStatsVo 用户统计
type UserStatsVo struct {
	TotalUsers   int     `json:"totalUsers"`
	ActiveUsers  int     `json:"activeUsers"`   // 7天内有登录的
	InactiveUsers int    `json:"inactiveUsers"` // 7天内无登录的
	DisabledUsers int    `json:"disabledUsers"` // 已禁用
	NewUsers7Days int    `json:"newUsers7Days"` // 7天内新增
}

// TrafficStatsVo 流量统计
type TrafficStatsVo struct {
	TotalBandwidth uint64 `json:"totalBandwidth"` // 总流量配额 Bytes
	UsedBandwidth  uint64 `json:"usedBandwidth"`  // 已用流量 Bytes
	UnusedBandwidth uint64 `json:"unusedBandwidth"` // 剩余流量 Bytes
	UsagePercent   float64 `json:"usagePercent"`   // 使用率 %
}

// NodeStatsVo 节点统计
type NodeStatsVo struct {
	TotalNodes   int `json:"totalNodes"`
	OnlineNodes  int `json:"onlineNodes"`
	OfflineNodes int `json:"offlineNodes"`
}

// ProtocolStatsVo 协议分布
type ProtocolStatsVo struct {
	Protocols []ProtocolCount `json:"protocols"`
}

// ProtocolCount 单个协议的节点数量
type ProtocolCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// GetUserStats 获取用户统计
func GetUserStats() (*UserStatsVo, error) {
	users, err := dao.SelectAllAccounts()
	if err != nil {
		return nil, err
	}

	vo := &UserStatsVo{
		TotalUsers: len(users),
	}

	now := util.NowMilli()
	sevenDaysAgo := now - 7*24*60*60*1000

	for _, user := range users {
		if user.Deleted != nil && *user.Deleted == 1 {
			vo.DisabledUsers++
		}
		if user.LastLoginTime != nil && *user.LastLoginTime >= sevenDaysAgo {
			vo.ActiveUsers++
		} else {
			vo.InactiveUsers++
		}
		// 粗略估算7天内新增（createTime 字段缺失，用 presetExpire 配合估算）
		if user.PresetExpire != nil && *user.PresetExpire > 0 && user.LastLoginTime != nil && *user.LastLoginTime == 0 {
			// 无登录记录视为新注册账号（简化判断）
			vo.NewUsers7Days++
		}
	}

	return vo, nil
}

// GetTrafficStats 获取流量统计
func GetTrafficStats() (*TrafficStatsVo, error) {
	users, err := dao.SelectAllAccounts()
	if err != nil {
		return nil, err
	}

	var total, used uint64
	for _, user := range users {
		if user.Quota != nil {
			total += uint64(*user.Quota)
		}
		if user.Upload != nil && user.Download != nil {
			used += uint64(*user.Upload) + uint64(*user.Download)
		}
	}

	var percent float64
	if total > 0 {
		percent = float64(used) / float64(total) * 100
	}

	return &TrafficStatsVo{
		TotalBandwidth:  total,
		UsedBandwidth:   used,
		UnusedBandwidth: total - used,
		UsagePercent:    percent,
	}, nil
}

// GetNodeStats 获取节点统计
func GetNodeStats() (*NodeStatsVo, error) {
	nodes, err := dao.SelectAllNodes()
	if err != nil {
		return nil, err
	}

	vo := &NodeStatsVo{
		TotalNodes: len(nodes),
	}

	// 节点状态需要通过节点服务器健康检查获取
	// 目前暂时按全部在线，后续可接入健康检查数据
	vo.OnlineNodes = len(nodes)
	vo.OfflineNodes = 0

	return vo, nil
}

// GetProtocolStats 获取各协议的节点分布
func GetProtocolStats() (*ProtocolStatsVo, error) {
	nodes, err := dao.SelectAllNodes()
	if err != nil {
		return nil, err
	}

	countMap := make(map[uint]int)
	for _, node := range nodes {
		if node.NodeTypeId != nil {
			countMap[*node.NodeTypeId]++
		}
	}

	var protocols []ProtocolCount
	for nodeTypeId, count := range countMap {
		var name string
		switch nodeTypeId {
		case constant.Xray:
			name = "Xray (VLESS/VMess/SS)"
		case constant.TrojanGo:
			name = "Trojan-Go"
		case constant.Hysteria:
			name = "Hysteria"
		case constant.NaiveProxy:
			name = "NaiveProxy"
		case constant.Hysteria2:
			name = "Hysteria2"
		default:
			name = "Unknown"
		}
		protocols = append(protocols, ProtocolCount{Name: name, Count: count})
	}

	return &ProtocolStatsVo{Protocols: protocols}, nil
}
