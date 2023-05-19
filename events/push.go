package events

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
)

// EventPush 消息推送
type EventPush struct {
	pushBarkURL     string
	pushDingtalkURL string
	pushName        string
	myIP            string
	msgChan         chan string
	client          *http.Client
}

var push *EventPush

// DefaultEventPush Default
func DefaultEventPush() *EventPush {
	if push == nil {
		push = &EventPush{
			myIP:    getmyIP(),
			msgChan: make(chan string, 100),
			client: &http.Client{
				Timeout: time.Second * 10,
			},
		}
		go push.loop()
	}
	return push
}

// SetConfig SetConfig
func (p *EventPush) SetConfig(name, barkUrl, dingtalkUrl string) {
	p.pushName = name
	p.pushBarkURL = barkUrl
	p.pushDingtalkURL = dingtalkUrl
}

func getmyIP() string {
	res, _ := http.Get("http://vip3.duodianzhekou.com/ip.php")
	if res != nil && res.Body != nil {
		defer res.Body.Close()
		by, _ := ioutil.ReadAll(res.Body)
		return string(by)
	}
	return ""
}

// PushMsg 发送消息
func (p *EventPush) PushMsg(msg string) {
	p.msgChan <- msg
}
func (p *EventPush) loop() {
	for content := range p.msgChan {
		if p.pushBarkURL != "" {
			// 发送
			// log.WithFields(log.Fields{"msg": msg}).Info("pushBarkURL is null")
			// return
			msg := fmt.Sprintf("%s IP:%s,%s", p.pushName, p.myIP, content)
			// pushBarkURL := "https://api.day.app/UHsDZHcVgcWbtaAkfCDsUT/"
			url := fmt.Sprintf("%s%s?sound=birdsong&group=监控提醒test", p.pushBarkURL, url.QueryEscape(msg))
			log.WithFields(log.Fields{"msg": msg}).Info("push msg")
			req, _ := http.NewRequest("GET", url, nil)
			resp, err := p.client.Do(req)
			if err != nil {
				log.WithFields(log.Fields{"msg": msg}).Error("push msg", err)
			}
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			log.WithFields(log.Fields{"content": string(body)}).Info("push msg result")
		}
		if p.pushDingtalkURL != "" {
			markTxt := fmt.Sprintf("#### %s事件\\n\\n > IP:%s\\n\\n > 内容:%s\\n\\n", p.pushName, p.myIP, content)
			msg := fmt.Sprintf(`{"msgtype": "markdown", "markdown": {"title": "%s", "text": "%s"}, "at": {"atMobiles": [], "isAtAll": false}}`, p.pushName, markTxt)
			req, _ := http.NewRequest("POST", p.pushDingtalkURL, bytes.NewBuffer([]byte(msg)))
			// req.Header.Set("X-Custom-Header", "myvalue")
			req.Header.Set("Content-Type", "application/json")
			resp, err := p.client.Do(req)
			if err != nil {
				log.WithFields(log.Fields{"msg": msg}).Error("push msg", err)
				continue
			}
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			log.WithFields(log.Fields{"content": string(body)}).Info("push msg result")
		}
	}
}
