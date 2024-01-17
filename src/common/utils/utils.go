// @Title
// @Author  zls  2023/9/27 16:14
package utils

import (
	structsm "devman/src/structs/datamap"
	"net/http"
	"strings"
	"sync"
)

// 发送http请求
func HttpClient(method string, url string, body string) (string, bool) {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return err.Error(), false
	}
	req.Header.Add("Content-Type", "application/json")
	//req.Header.Add("Authorization", "Bearer ")
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return err.Error(), false
	}
	defer resp.Body.Close()
	return "", true
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
