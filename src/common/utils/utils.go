// @Title
// @Author  zls  2023/9/27 16:14
package utils

import (
	structsm "devman/src/structs/datamap"
	"io"
	"net/http"
	"strings"
	"sync"
)

// 发送http请求-get
func SendHttpRequstGet(url string, head map[string]string) (string, error) {
	return sendHttpRequst("GET", url, "", head)
}


// 发送http请求
func sendHttpRequst(method string, url string, requestBody string, head map[string]string) (string, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	if len(head) != 0 {
		for key, val := range head {
			req.Header.Add(key, val)
		}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return "", err
	}
	defer resp.Body.Close()
	respBodyB, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	bodyStr := string(respBodyB)
	return bodyStr, nil
}

var lock sync.RWMutex

// 加了读写锁的Map
func PutMap(mapVal map[string][]structsm.TableInfo, key string, val []structsm.TableInfo) {
	lock.Lock()
	defer lock.Unlock()
	mapVal[key] = val
}

// 加了读写锁的Map
func GetMap(mapVal map[string][]structsm.TableInfo, key string) ([]structsm.TableInfo, bool) {
	lock.Lock()
	defer lock.Unlock()
	infos, ok := mapVal[key]
	return infos, ok
}
