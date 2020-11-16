package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"time"
)

func queryLog(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("server:", r.Form["server"])
	cmd := exec.Command("/bin/bash", "-c", "grep g main.go")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Execute Command failed:" + err.Error())
	} else {
		fmt.Fprintf(w, string(output))
		fmt.Printf("finished with output:\n%s", string(output))
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
