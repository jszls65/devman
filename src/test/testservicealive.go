package main

import "devman/src/job"

// 测试服务存活

func main() {
	c := job.ServiceAliveCheck{}
	c.Run()
}
