// @Title
// @Author  zls  2023/9/16 10:29
package common

import "github.com/gin-gonic/gin"

// code 0-成功  非0 失败
func ResultMsg(code int, msg string) gin.H {
	return gin.H{
		"code": code,
		"msg":  msg,
		"data": nil,
	}
}

func ResultOk() gin.H {
	return gin.H{
		"code": 0,
		"msg":  "操作成功",
		"data": nil,
	}
}

func ResultOkMsg(msg string) gin.H {
	return gin.H{
		"code": 0,
		"msg":  msg,
		"data": nil,
	}
}

func ResultFail(msg string) gin.H {
	return gin.H{
		"code": 1,
		"msg":  msg,
		"data": nil,
	}
}
