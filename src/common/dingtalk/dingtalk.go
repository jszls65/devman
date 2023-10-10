// @Title
// @Author  zls  2023/9/21 23:48
package dingtalk

import (
	"bytes"
	"dev-utils/config"
	"encoding/json"
	"log"
	"net/http"
)

type Message struct {
	Type string `json:"msgtype"`
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
	At struct {
		IsAtAll   string   `json:"isAtAll"`
		AtMobiles []string `json:"atMobiles"`
	} `json:"at"`
}

func SendText(content string, phone string) {
	log.Println("发送钉钉消息: ", content)
	if !config.Conf.DingTalk.Enable {
		log.Println("钉钉消息未打开, 消息无法发送")
		return
	}
	msg := Message{
		Type: "text",
		Text: struct {
			Content string `json:"content"`
		}{Content: content},
		At: struct {
			IsAtAll   string   `json:"isAtAll"`
			AtMobiles []string `json:"atMobiles"`
		}{IsAtAll: "false", AtMobiles: []string{phone}},
	}
	body, _ := json.Marshal(msg)
	log.Println("钉钉消息内容: ", string(body))
	req, _ := http.NewRequest("POST",
		config.Conf.DingTalk.Url,
		bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	_, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("发送钉钉消息失败:", err.Error())
	}

}
