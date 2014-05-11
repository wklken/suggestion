package main

import (
    "./darts"
    "fmt"
)

func main() {
    d, err := darts.Import("data.txt", "data.lib")
    fmt.Println(err)

    //fmt.Println(d.KeyString2IntMap)
    //fmt.Println(d.KeyInt2InfoMap)

    if err == nil {
        fmt.Println("搜索: 植物大战")
        results := d.Search([]rune("植物大战"), 0)

        fmt.Println("Result Len:", len(results))

        for i := 0; i < len(results); i++ {
            fmt.Println(string(results[i].Key), results[i].Value)
        }

    }
}
