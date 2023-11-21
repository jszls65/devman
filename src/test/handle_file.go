package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	// fileName := "/Users/zhangliansheng/Downloads/超时的.log"
	fileName := "/Users/zhangliansheng/Downloads/2023-11-12.0.error.log"
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		log.Fatalln("--->", err)
		return
	}
	reader := bufio.NewReader(file)
	// 按行读取  ReaderString
	index := 0
	lineMap := make(map[int]string)
	dupMap := make(map[string]int)
	for {
		index++
		lineStr, err := reader.ReadString('\n')
		// fmt.Println("第",index,"行: ", lineStr)
		if err == io.EOF {
			break
		}

		// fmt.Println("第", index, "行: ", "newUrlStr:", lineStr)
		lineMap[index] = lineStr
		// 排除的行 !strings.Contains(lineStr, ",") || !strings.Contains(lineStr, "/") || strings.Contains(lineStr, "调用亚马逊接口失败") ||
		if strings.Contains(lineStr, "amazon-advertising-api-snapshots-prod-feamazon.s3.amazonaws.com") ||
			strings.Contains(lineStr, "amazon-advertising-api-snapshots-prod-euamazon.s3.amazonaws.com") ||
			strings.Contains(lineStr, "amazon-advertising-api-snapshots-prod-usamazon.s3.amazonaws.com") ||
			strings.Contains(lineStr, "sd-snapshotfiles-na-prod.s3.amazonaws.com") ||
			!strings.Contains(lineStr, "snapshot") {
			continue
		}
		// 处理行
		lineSplit := strings.Split(lineStr, ",")
		if len(lineSplit) < 2 {
			continue
		}

		urlStr := strings.Split(lineSplit[1], "/")
		newUrlStr := strings.Trim(urlStr[0], " ")

		dupMap[newUrlStr] = index

	}

	fmt.Println("------------------------------")
	for _, v := range dupMap {
		fmt.Println("行号:", v, "   ", lineMap[v])
	}
}
