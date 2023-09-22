// @Title
// @Author  zls  2023/9/21 23:48
package dingtalk

import (
	"bytes"
	"dev-utils/config"
	"encoding/json"
	"net/http"
)

type Message struct {
	Type string `json:"msgtype"`
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
}

func SendText(content string) {
	msg := Message{
		Type: "text",
		Text: struct {
			Content string `json:"content"`
		}{Content: content},
	}
	body, _ := json.Marshal(msg)
	req, _ := http.NewRequest("POST",
		config.Conf.DingTalk.Url,
		bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	_, _ = http.DefaultClient.Do(req)
}
