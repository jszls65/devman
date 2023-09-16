// @Title
// @Author  zls  2023/9/16 10:29
package common

import "github.com/gin-gonic/gin"

func Result(msg string, data any) gin.H {
	return gin.H{
		"msg":  msg,
		"data": data,
	}
}

func ResultMsg(msg string) gin.H {
	return gin.H{
		"msg":  msg,
		"data": nil,
	}
}

func ResultOk() gin.H {
	return gin.H{
		"msg":  "操作成功",
		"data": nil,
	}
}

func ResultFail() gin.H {
	return gin.H{
		"msg":  "操作失败",
		"data": nil,
	}
}
