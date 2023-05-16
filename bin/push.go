package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const pushBarkURL = "https://api.day.app/UHsDZHcVgcWbtaAkfCDsUT/"

func main() {
	// args := os.Args
	// str := fmt.Sprintf("len=%d, %v", len(args), args)
	// ioutil.WriteFile("./log.txt", []byte(str), 0666)
	file, err := os.OpenFile("./logs/push.log", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	//及时关闭file句柄
	defer file.Close()
	// var data string
	file.WriteString(fmt.Sprintf("%s: 启动\n", time.Now().Format("2006-01-02 15:04:05")))
	// fmt.Println(strings.Fields("a b"))
	msg := ""
	for {
		fmt.Println("READY")
		// fmt.Scanln(&data)
		inputReader := bufio.NewReader(os.Stdin)
		input, err := inputReader.ReadString('\n')
		if err != nil {
			break
		}
		if input == "" {
			continue
		}
		msg = input
		// ioutil.WriteFile("./log.txt", []byte(data), 0666)
		file.WriteString(fmt.Sprintf("%s: %s\n", time.Now().Format("2006-01-02 15:04:05"), input))
		if strings.HasPrefix(input, "ver:3.0") {
			// 再读一行数据
			keyV := strings.Fields(input)
			for _, v := range keyV {
				arr := strings.Split(v, ":")
				if arr[0] == "len" {
					n, err := strconv.Atoi(arr[1])
					if err == nil && n > 0 {
						//读
						// read n bytes
						b := make([]byte, n)
						for i := 0; i < n; i++ {
							b[i], err = inputReader.ReadByte()
							if err != nil {
								break
							}
						}
						// log.Warn("events readResult ok:", string(b))
						file.WriteString(string(b))
						file.WriteString("\n")
						msg = string(b)
					}
				}
			}
		}

		fmt.Print("RESULT 2\n")
		fmt.Print("OK")

		go pushBarkMsg(msg)
	}
	fmt.Println("end")
}

func pushBarkMsg(msg string) {
	//发送消息
	url := fmt.Sprintf("%s%s?sound=birdsong&group=监控提醒", pushBarkURL, url.QueryEscape(msg))
	// log.Info("push msg", url)
	_, err := http.Get(url)
	if err != nil {
		// log.Error(err)
	}
}
