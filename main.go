package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"github.com/gin-gonic/gin"
	"os"
	"sort"
)

func Gin(){
	r := gin.Default()
}
// 读取数据，并且给定writer的接口，并且指定开始行数和结束的行数。
func ReadFile(writer io.Writer,fileName string,startNum,endNum int)error{
	i := 1
	f, err := os.Open(fileName)
	i ++
	if err != nil {
		panic(err)
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	j := 1
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil || io.EOF == err { // 达到这个数据该跳转的时候就break。
			return nil
		}
		if j >=startNum && j <= endNum {
			fmt.Println(line)
		}
		j++
	}
	return nil
}
func ReadFileName()[]string  {
	result := make([]string,0)
	files, _ := ioutil.ReadDir("./file")
	for _, f := range files {
		result = append(result,f.Name())
	}
	sort.Strings(result)
	return result
}
func main() {

	ReadFile(nil,"./1.dat",10,11)
}
