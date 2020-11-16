package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var serverMap = make(map[string]interface{})

func loadConfig() bool {
	data, err := ioutil.ReadFile("static/server.json")
	if err != nil {
		fmt.Println("Load Config Error")
		return false
	}
	json.Unmarshal(data, &serverMap)
	return true
}

func showPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("main.html")
	t.Execute(w, nil)
}

func forward(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	r.ParseForm()

	serverArr := r.Form["server"]
	if len(serverArr) < 1 {
		fmt.Println("len err")
		return
	}
	server := r.Form["server"][0]
	fmt.Println("server:", server)
	fmt.Println("key:", r.Form["key2"])
	ip, ok := serverMap[server]
	if !ok {
		fmt.Fprintf(w, "无此服务器")
		fmt.Fprintf(w, server)
		return
	}

	defer r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	remote, _ := url.Parse("http://" + ip.(string) + ":9001")
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
}

func main() {
	if !loadConfig() {
		return
	}
	fmt.Println(serverMap)

	http.HandleFunc("/", showPage)
	http.HandleFunc("/query", forward)
	http.HandleFunc("/gettime", forward)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
