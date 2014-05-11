package main

import (
    "./darts"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    //"strings"
)

type ReturnJson struct {
    Ok     bool     `json:"ok,string"`
    Data   []string `json:"data,string"`
    Reason string   `json:"reason,string"`
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

    var result ReturnJson

    r.ParseForm() //解析参数，默认是不会解析的
    keyword := r.Form["keyword"]
    if len(keyword) == 0 {
        result.Ok = false
        result.Reason = "Keyword required!"
    } else {
        results := dart.Search([]rune(keyword[0]), 0)

        result.Ok = true
        for i := 0; i < len(results); i++ {
            //fmt.Println(string(results[i].Key), results[i].Value)
            result.Data = append(result.Data, string(results[i].Key))
        }
        result.Data = result.Data[:10]
    }

    b, err := json.Marshal(result)
    if err != nil {
        fmt.Println("json err:", err)
    }
    fmt.Fprintf(w, string(b)) //这个写入到w的是输出到客户端的
    return
}

func main() {
    http.HandleFunc("/suggest/", simpleSuggest) //设置访问的路由
    err := http.ListenAndServe(":9090", nil)    //设置监听的端口
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
