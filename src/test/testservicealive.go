package main

import "dev-utils/src/job"

// 测试服务存活

func main() {
	c := job.ServiceAliveCheck{}
	c.Run()
}
