package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

const (
	XRealIP       = "X-Real-IP"
	XForwardedFor = "X-Forwarded-For"
)

const (
	ServerFile     = "static/server.json"
	LogPathFile    = "static/logpath.json"
	CustomPathFile = "static/customlogpath.json"
)

var serverMap = make(map[string]interface{})
var customPathMap = make(map[string]string)

/*
type IP struct {
	Address string `json:"address"`
	Port    string `json:"port"`
}
*/

func loadConfig() bool {
	data, err := ioutil.ReadFile("static/server.json")
	if err != nil {
		fmt.Println("Load Config Error")
		return false
	}
	json.Unmarshal(data, &serverMap)

	data, err = ioutil.ReadFile(CustomPathFile)
	if err != nil {
		fmt.Println("Load Path Config Error")
		return false
	}
	json.Unmarshal(data, &customPathMap)
	return true
}

func getIp(r *http.Request) string {
	addr := r.RemoteAddr
	if ip := r.Header.Get(XRealIP); ip != "" {
		addr = ip
	} else if ip = r.Header.Get(XForwardedFor); ip != "" {
		addr = ip
	} else {
		addr, _, _ = net.SplitHostPort(addr)
	}
	if addr == "::1" {
		addr = "127.0.0.1"
	}
	return addr
}

func getIpAjax(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, getIp(r))
}

func saveServer() {
	file, _ := os.OpenFile("static/server.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	encoder := json.NewEncoder(file)
	err := encoder.Encode(serverMap)
	if err == nil {
		fmt.Println("save server success")
	}
	file, _ = os.OpenFile("static/customlogpath.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	encoder = json.NewEncoder(file)
	err = encoder.Encode(customPathMap)
	if err == nil {
		fmt.Println("save path success")
	}
}

func showPage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("main.html")
	t.Execute(w, nil)
}

func addSrv(w http.ResponseWriter, r *http.Request) {
	srvName := r.FormValue("logsrvname")
	srvIp := r.FormValue("logsrvip")
	srvPath := r.FormValue("logsrvpath")
	clientIp := getIp(r)
	serverMap[clientIp+srvName] = srvIp
	customPathMap[clientIp+srvName] = srvPath
	result := "<option value='" + clientIp + srvName + "'>" + srvName + "</option>"
	fmt.Fprintf(w, result)
	saveServer()
}

func forward(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	r.ParseForm()

	server := r.FormValue("server")
	if server == "" {
		fmt.Println("len err")
		return
	}
	if strings.Index(r.RequestURI, "query") != -1 && r.FormValue("key0") == "" {
		fmt.Fprintf(w, "请输入查询关键字!")
		return
	}

	ip, ok := serverMap[server]
	if !ok {
		fmt.Fprintf(w, "无此服务器")
		fmt.Fprintf(w, server)
		return
	}

	if path, ok := customPathMap[server]; ok {
		_, file := filepath.Split(r.FormValue("log"))
		fmt.Println(file)
		r.Form["log"] = []string{path + "/" + file}
	}
	r.Body = ioutil.NopCloser(strings.NewReader(r.Form.Encode()))
	r.ContentLength = int64(len(r.Form.Encode()))
	//r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	remote, _ := url.Parse("http://" + ip.(string) + ":9001")
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
}

func syncWatch() {
	url := "http://127.0.0.1:9001/syncwatch"
	contentType := "application/json;charset=utf-8"
	//b := "hello"
	//b := []byte("Hello, Server")
	//body := bytes.NewBuffer(b)
	str := strings.NewReader("hello")
	resp, err := http.Post(url, contentType, str)
	if err != nil {
		log.Println("Post failed:", err)
		return
	}
	defer resp.Body.Close()
}

func main() {
	if !loadConfig() {
		return
	}
	fmt.Println(serverMap)

	ip, _ := serverMap["qa68"]
	fmt.Println(reflect.TypeOf(ip).Name() == "string")
	ip, _ = serverMap["qa70"]
	fmt.Println(ip.(map[string]interface{})["address"].(string) + ip.(map[string]interface{})["port"].(string))

	http.HandleFunc("/", showPage)
	http.HandleFunc("/addsrv", addSrv)
	http.HandleFunc("/getip", getIpAjax)
	http.HandleFunc("/query", forward)
	http.HandleFunc("/gettime", forward)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// 同步监控信息 与日志无关
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			syncWatch()
		}
	}()

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
