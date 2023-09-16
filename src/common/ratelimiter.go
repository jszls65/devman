// @Title
// @Author  zls  2023/9/15 17:11
package common

import (
	"golang.org/x/time/rate"
	"time"
)

var limiterMap = make(map[string]*rate.Limiter, 0)

// Limiter 创建限速器
// 参数 scene 场景, 作为map的key
// every : 每秒往桶中放令牌的数量
// total: 令牌桶容量
func Limiter(scene string, every int, total int) *rate.Limiter {
	if limiterMap[scene] != nil {
		return limiterMap[scene]
	}
	limiter := rate.NewLimiter(rate.Every(time.Second*time.Duration(every)), total)
	limiterMap[scene] = limiter
	return limiter
}
