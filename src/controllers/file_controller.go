// @Title
// @Author  zls  2023/11/8 21:35
package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"time"
)

type FileController struct {
}

// DownloadFile 下载文件, 备份db文件
func (con FileController) DownloadFile(context *gin.Context) {
	file, err := os.Open("../dev-utils.db")
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "../dev-utils.db文件不存在！"})
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("关闭文件异常:", err.Error())
		}
	}(file)

	stat, _ := file.Stat()
	newFileName := stat.Name()
	newFileName = newFileName[:len(newFileName)-3]
	newFileName = newFileName + "-" + time.Now().Format(time.DateOnly) + ".db"
	fmt.Println("newFileName:", newFileName)
	context.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", newFileName))
	context.Header("Content-Type", "application/octet-stream")
	_, _ = io.Copy(context.Writer, file)
}
