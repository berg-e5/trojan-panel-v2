package service

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"strconv"
	"trojan-panel/model/constant"
)

var bot = new(tgbotapi.BotAPI)

func NewTelegramBotApi() (*tgbotapi.BotAPI, error) {
	var err error
	// 从数据库中查询 api token
	apiToken := ""
	bot, err = tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		logrus.Errorf("new bot api err: %v", err)
		return nil, errors.New(constant.TelegramBotApiError)
	}
	logrus.Infof("Authorized on account %s", bot.Self.UserName)
	return bot, nil
}

func GetUpdatesChan() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return bot.GetUpdatesChan(u)
}

// SendTelegramMessage 发送 Telegram 消息
func SendTelegramMessage(chatID, text string) error {
	if bot == nil {
		if _, err := NewTelegramBotApi(); err != nil {
			return err
		}
	}

	cid, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		return errors.New("invalid telegram chat ID")
	}

	msg := tgbotapi.NewMessage(cid, text)
	msg.ParseMode = "HTML"

	_, err = bot.Send(msg)
	if err != nil {
		logrus.Errorf("send telegram message err: %v", err)
		return err
	}
	return nil
}

// SendTelegramMessageToAdmin 向管理员发送 Telegram 消息
// TODO: 需在 SystemVo 中增加 telegram_admin_id 字段
func SendTelegramMessageToAdmin(text string) error {
	return errors.New("telegram admin id not configured, please add telegram_admin_id to system config")
}
