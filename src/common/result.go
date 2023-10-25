// @Title
// @Author  zls  2023/9/16 10:29
package common

import "github.com/gin-gonic/gin"

// code 0-成功  非0 失败
func ResultMsg(code int, msg string) gin.H {
	return gin.H{
		"code": code,
		"msg":  msg,
	}
}

func ResultOk() gin.H {
	return gin.H{
		"code": 0,
		"msg":  "操作成功",
	}
}

func ResultOkMsg(msg string) gin.H {
	return gin.H{
		"code": 0,
		"msg":  msg,
	}
}

func ResultOkMsgData(msg string, data interface{}) gin.H {
	return gin.H{
		"code": 0,
		"msg":  msg,
		"data": data,
	}
}
func ResultFail(msg string) gin.H {
	return gin.H{
		"code": 1,
		"msg":  msg,
	}
}
