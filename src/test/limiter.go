// @Title
// @Author  zls  2023/9/15 17:00
package main

import (
	"fmt"
	"golang.org/x/time/rate"
	"time"
)

func main() {
	l := rate.NewLimiter(rate.Every(time.Second/10), 10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			for {
				// 判断是否有令牌，如果有就输出
				if l.Allow() {
					fmt.Printf("allow %d\n", i)
				}
				// 每0.5秒请求一次
				time.Sleep(time.Second / 2)
			}
		}(i)
	}
	time.Sleep(time.Second * 10)

}
