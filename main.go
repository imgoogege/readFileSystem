package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

func Gin() {
	r := gin.Default()
	// 此文件 读取html
	r.GET("/", func(ctx *gin.Context) {
		tem, err := template.ParseFiles("./index.html") // 将文件 导出为完整的HTML
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"data":    "首页载入出错",
				"success": "error",
			})
			return
		}
		tem.Execute(ctx.Writer, nil)
	})

	// 此文件下载文件。
	r.POST("/result", func(ctx *gin.Context) {
		start := ctx.PostForm("startNum")
		end := ctx.PostForm("endNum")
		// 将从前端传过来的string转换为int
		startNum, err := strconv.Atoi(start)
		if err != nil {
			ctx.JSON(502, gin.H{
				"message": "你应该输入数字，你输入的非可用字段",
			})
		}
		endNum, err := strconv.Atoi(end)
		if err != nil {
			ctx.JSON(502, gin.H{
				"message": "你应该输入数字，你输入的非可用字段",
			})
		}
		// 读取文件名称，并且排序
		fileNames := ReadFileName(ctx)
		// 读取文件。并且写入文件。
		if err = ReadFile(ctx, nil, fileNames, startNum, endNum); err != nil {
			ctx.JSON(502, gin.H{
				"message": "读取文件错误",
				"error":   err,
			})
		}

		ctx.JSON(200, gin.H{
			"message": "已经打包好数据，请在本地寻找",
			"error":   "nil",
		})
	})
	r.Run()
}

// 读取数据，并且给定writer的接口，并且指定开始行数和结束的行数。
func ReadFile(ctx *gin.Context, writer io.Writer, fileNames []string, startNum, endNum int) error {
	result := make([][]string, 0)
	for _, fileName := range fileNames {
		j := 1
		f, err := os.Open("./file/" + fileName + ".dat")
		if err != nil {
			ctx.JSON(503, gin.H{
				"message": "无法读取文件",
				"error":   err,
			})
		}
		rd := bufio.NewReader(f)
		for {
			str, err := rd.ReadString('\n')
			if io.EOF == err || err != nil {
				break
			}
			// An artificial input source.
			if j >= startNum && j <= endNum {
				re := make([]string, 0)
				scanner := bufio.NewScanner(strings.NewReader(str))
				// Set the split function for the scanning operation.
				scanner.Split(bufio.ScanWords)
				// Count the words.
				scanTT := 0
				for scanner.Scan() {
					if scanTT < 5 {
						re = append(re, scanner.Text())
					}
					scanTT++
				}
				if err := scanner.Err(); err != nil {
					ctx.JSON(503, gin.H{
						"message": "无法scannner",
					})
				}
				result = append(result, re)
			}
			j++

		}

		f.Close()
	}
	writeFileToCsv(ctx, result)
	DeleteAllFile(ctx)
	return nil
}

// 读取文件
func ReadFileName(ctx *gin.Context) []string {
	result := make([]string, 0)
	files, err := ioutil.ReadDir("./file")
	if err != nil {
		ctx.JSON(502, gin.H{
			"message": "无法读取数据名称",
			"error":   err,
		})
	}
	for _, f := range files {
		result = append(result, f.Name())
	}
	for i, v := range result {
		if len(v) == 5 {
			result[i] = v[:1]
		} else if len(v) == 6 {
			result[i] = v[:2]
		} else if len(v) == 7 {
			result[i] = v[:3]
		} else {
			ctx.JSON(502, gin.H{
				"message": "设置问题，数据太多了，已经超过1000了无法操作",
			})
		}
	}
	resut := []int{}
	for i := range result {
		v, _ := strconv.Atoi(result[i])
		resut = append(resut, v)
	}
	sort.Ints(resut)
	for i := range resut {
		t := strconv.FormatInt(int64(resut[i]), 10)
		result[i] = t
	}
	fmt.Println("测试数据")
	return result
}

// csv
func writeFileToCsv(ctx *gin.Context, re [][]string) {
	file, err := os.Create("./aoligei.csv")
	if err != nil {
		ctx.JSON(502, gin.H{
			"message": "无法创建文件",
			"error":   err,
		})
	}

	w := csv.NewWriter(file)
	w.Write([]string{"时间", "电流", "极化", "电流", "极化"})
	for _, r := range re {
		if err := w.Write(r); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
		// Write any buffered data to the underlying writer (standard output).
	}
	w.Flush()
	if err := w.Error(); err != nil {
		ctx.JSON(502, gin.H{
			"message": "写入数据错误",
			"error":   err,
		})
	}

}

func DeleteAllFile(ctx *gin.Context) {
	err := os.RemoveAll("./file")
	if err != nil {
		ctx.JSON(502, gin.H{
			"message": "无法删除文件",
			"error":   err,
		})
	}
}
func main() {
	Gin()
}
