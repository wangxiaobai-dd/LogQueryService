package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func queryLog(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// 查询字符串
	execStr := "grep '" + r.Form["key0"][0] + "' " + r.Form["log"][0] + "*.log"
	// 日期
	if _, ok := r.Form["realTime"]; !ok {
		date := r.Form["logdate"][0]
		dateArr := strings.Split(date, "-")
		year := dateArr[0][2:]
		execStr = execStr + "." + year + dateArr[1] + dateArr[2] + "*"
		fmt.Println(dateArr)
		fmt.Println(year)
	}

	// 关键字
	for i := 1; ; i++ {
		if key, ok := r.Form["key"+strconv.Itoa(i)]; ok && len(key[0]) > 0 {
			execStr += " | grep '" + key[0] + "' "
		} else {
			break
		}
	}
	// 排除关键字
	for i := 0; ; i++ {
		if exkey, ok := r.Form["exkey"+strconv.Itoa(i)]; ok && len(exkey[0]) > 0 {
			execStr += " | grep -v '" + exkey[0] + "' "
		} else {
			break
		}
	}
	fmt.Println(execStr)

	cmd := exec.Command("/bin/bash", "-c", execStr)
	output, err := cmd.Output()
	if err != nil {
		if strings.Index(err.Error(), "1") != -1 {
			fmt.Fprintf(w, "无查询结果")
		} else if strings.Index(err.Error(), "2") != -1 {
			fmt.Fprintf(w, "没有此日志文件")
		}
		fmt.Println(err)
	} else {
		fmt.Fprintf(w, string(output))
	}
}

func getTime(w http.ResponseWriter, r *http.Request) {
	timestamp := time.Now().Unix()
	fmt.Println(timestamp)
	fmt.Fprintf(w, strconv.FormatInt(timestamp, 10))
}

func main() {
	http.HandleFunc("/query", queryLog)
	http.HandleFunc("/gettime", getTime)
	err := http.ListenAndServe(":9001", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
