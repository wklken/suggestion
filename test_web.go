package main

import (
	"./darts"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	//"strings"
	// "path"
)

type ValueJson struct {
	Value string `json:"value"`
}

func dartsInit() (darts.Darts, error) {
	d, err := darts.Import("data.txt", "data.lib")
	if err != nil {
		fmt.Println("ERROR: darts initial failed!")
	} else {
		fmt.Println("INFO: darts initial success!")
	}

	return d, err
}

var dart, err = dartsInit()

func simpleSuggest(w http.ResponseWriter, r *http.Request) {

	// valueList := make([]ValueJson, 10)
	var valueList []ValueJson

	r.ParseForm() //解析参数，默认是不会解析的
	keyword := r.Form["keyword"]

	fmt.Println("keyword:", keyword)
	if len(keyword) == 0 {

	} else {
		results := dart.Search([]rune(keyword[0]), 0)
		for i := 0; i < len(results); i++ {
			var value ValueJson
			value.Value = string(results[i].Key)
			valueList = append(valueList, value)
		}

	}

	if len(valueList) > 10 {
		valueList = valueList[:10]
	}

	fmt.Println("return", valueList)

	if len(valueList) > 0 {
		b, err := json.Marshal(valueList)
		if err != nil {
			fmt.Println("json err:", err)
		}
		fmt.Fprintf(w, string(b)) //这个写入到w的是输出到客户端的
	} else {
		fmt.Fprintf(w, "[]") //这个写入到w的是输出到客户端的
	}
	return
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func main() {

	// index
	http.HandleFunc("/", index) //设置访问的路由

	// suggest
	http.HandleFunc("/suggest/", simpleSuggest) //设置访问的路由

	// static
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	err := http.ListenAndServe("0.0.0.0:9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
