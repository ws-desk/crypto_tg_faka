package tg_bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/net/proxy"
	"gopay/internal/exts/config"
	my_log "gopay/internal/exts/log"
	"net/http"
	"net/url"
)

//type SilentLogger struct{}
//
//func (l SilentLogger) Println(v ...interface{})               {}
//func (l SilentLogger) Printf(format string, v ...interface{}) {}
//
//type BotLogger interface {
//	Println(v ...interface{})
//	Printf(format string, v ...interface{})
//}

var Bot *tgbotapi.BotAPI

func InitTGBot() {
	client := &http.Client{}

	if config.GetSiteConfig().Proxy.EnableProxy == true {
		tgProxyURL, err := url.Parse(fmt.Sprintf("%s://%s:%d", config.GetSiteConfig().Proxy.Protocol, config.SiteConfig.Proxy.Host, config.SiteConfig.Proxy.Port))
		if err != nil {
			panic(fmt.Sprintf("Failed to parse proxy: %s\n", err))
		}
		tgDialer, err := proxy.FromURL(tgProxyURL, proxy.Direct)
		if err != nil {
			panic(fmt.Sprintf("Failed to obtain proxy dialer: %s\n", err))
		}
		tgTransport := &http.Transport{
			Dial: tgDialer.Dial,
		}
		client.Transport = tgTransport
	}

	my_log.LogInfo("正在连接TG bot")
	var err error
	Bot, err = tgbotapi.NewBotAPIWithClient(config.GetSiteConfig().TgBotToken, "https://api.telegram.org/bot%s/%s", client)
	if err != nil {
		panic(err)
	}
	Bot.Debug = config.SiteConfig.EnableTGBotDebug
	my_log.LogInfo("成功连接TG bot")
}

func SendAdmin(msgText string) {
	<-config.SendAdminLimit.C

	msg := tgbotapi.NewMessage(config.SiteConfig.AdminTGID, msgText)
	msg.DisableWebPagePreview = true
	msg.ParseMode = "HTML"
	Bot.Send(msg)
}
func DeleteMsg(chatID int64, msgID int) error {
	deleteConfig := tgbotapi.DeleteMessageConfig{
		ChatID:    chatID,
		MessageID: msgID,
	}
	if _, err := Bot.Request(deleteConfig); err != nil {
		return err
	}
	return nil
}
