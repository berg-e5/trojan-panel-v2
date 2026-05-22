package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	redisgo "github.com/gomodule/redigo/redis"
	"trojan-panel/dao"
	redispkg "trojan-panel/dao/redis"
	"trojan-panel/model/constant"
	"trojan-panel/model/vo"
)

// RefreshTokenExpireMs Refresh Token 过期时间 7 天（毫秒）
const RefreshTokenExpireMs = int64(7 * 24 * 60 * 60 * 1000)

// generateRefreshToken 生成一个随机 Refresh Token
func generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// SaveRefreshToken 将 Refresh Token 存入 Redis，key = trojan-panel:refresh:{username}
func SaveRefreshToken(username string, token string) error {
	key := "trojan-panel:refresh:" + username
	_, err := redispkg.Client.String.Set(key, token, RefreshTokenExpireMs).Result()
	return err
}

// GenerateRefreshTokenOnly 生成并存储 Refresh Token，供 api 层在登录时调用
func GenerateRefreshTokenOnly(username string) (string, error) {
	token, err := generateRefreshToken()
	if err != nil {
		return "", err
	}
	if err := SaveRefreshToken(username, token); err != nil {
		return "", err
	}
	return token, nil
}

// GetRefreshToken 获取用户当前的 Refresh Token
func GetRefreshToken(username string) (string, error) {
	key := "trojan-panel:refresh:" + username
	result, err := redispkg.Client.String.Get(key).Result()
	if err != nil && err != redisgo.ErrNil {
		return "", err
	}
	if result == nil {
		return "", nil
	}
	s, _ := redisgo.String(result, nil)
	return s, nil
}

// DeleteRefreshToken 删除 Refresh Token（注销/顶出）
func DeleteRefreshToken(username string) error {
	key := "trojan-panel:refresh:" + username
	_, err := redispkg.Client.Key.Del(key).Result()
	return err
}

// RefreshAccessToken 用 Refresh Token 换新的 Access Token
// 返回新的 accessToken 和 refreshToken（滚动更新）
func RefreshAccessToken(refreshToken string) (string, string, error) {
	// 遍历所有用户查找对应的 refresh token
	// 生产环境建议用 Set NX 存 username:token 做反向索引
	scanResult, err := redispkg.Client.Key.Keys("trojan-panel:refresh:*").Strings()
	if err != nil {
		return "", "", errors.New(constant.SysError)
	}

	var matchedUsername string
	for _, key := range scanResult {
		storedResult, err := redispkg.Client.String.Get(key).Result()
		if err != nil && err != redisgo.ErrNil {
			continue
		}
		if storedResult != nil {
			stored, _ := redisgo.String(storedResult, nil)
			if stored == refreshToken {
				matchedUsername = key[len("trojan-panel:refresh:"):]
				break
			}
		}
	}

	if matchedUsername == "" {
		return "", "", errors.New(constant.TokenExpiredError)
	}

	// 查用户信息，重新生成 Access Token
	account, err := SelectAccountByUsername(&matchedUsername)
	if err != nil || account == nil {
		return "", "", errors.New(constant.TokenExpiredError)
	}
	if account.Deleted != nil && *account.Deleted != 0 {
		return "", "", errors.New(constant.AccountDisabled)
	}

	roles, err := dao.SelectRoleNameByParentId(account.RoleId, true)
	if err != nil {
		return "", "", errors.New(constant.SysError)
	}
	accountVo := vo.AccountVo{
		Id:       *account.Id,
		Username: *account.Username,
		RoleId:   *account.RoleId,
		Deleted:  *account.Deleted,
		Roles:    roles,
	}

	newAccessToken, err := GenToken(accountVo)
	if err != nil {
		return "", "", errors.New(constant.SysError)
	}

	// 滚动更新 Refresh Token（防泄露重用）
	newRefreshToken, err := generateRefreshToken()
	if err != nil {
		return "", "", errors.New(constant.SysError)
	}
	if err := SaveRefreshToken(matchedUsername, newRefreshToken); err != nil {
		return "", "", errors.New(constant.SysError)
	}

	return newAccessToken, newRefreshToken, nil
}
