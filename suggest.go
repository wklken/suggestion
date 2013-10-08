package main

import (
    "fmt"
)

type Node struct {
    Data map[string]*Node
    //= make(map[string]*Node)
    Is_leaf bool
    Weight int
}

func (n *Node) Init() {
    n.Data = make(map[string]*Node)
    n.Is_leaf = false
    n.Weight = 0

}

func (n *Node) Has_next() bool {
    return len(n.Data) > 0
}

func (n *Node) Add(keyword string, subnode *Node) {
    n.Data[keyword] = subnode
}

func (n *Node) Get(keyword string) *Node {
    return n.Data[keyword]
}

func (n *Node) Get_the_top_node(prefix string) *Node {
    top := n
    for _, c := range prefix {
        top = top.Get(string(c))
        if top != nil {
            continue
        } else {
            return nil
        }
    }
    return top
}

func Depth_walk(node *Node) map[string]int {
    result := make(map[string]int)
    if node.Is_leaf {
        result[""] = node.Weight
    }

    if node.Has_next() {
        for k, _ := range node.Data {
            s := Depth_walk(node.Get(k))
            for sk, sv := range s {
                result[k+sk] = sv
            }
        }
        return result
    } else {
        result[""] = node.Weight
        return result
    }
}

func Search(node *Node, prefix string, limit int) map[string]int {
    node = node.Get_the_top_node(prefix)

    result := make(map[string]int)

    if node.Is_leaf {
        result[prefix] = node.Weight
    }

    d := Depth_walk(node)
    for suffix, weight := range d {
        result[prefix+suffix] = weight
    }

    return result
}

func (n *Node) Str() string {
    return "<Node>"
}



func main() {
    fmt.Println("print line")
    n := new(Node)
    n.Init()

    hn := new(Node)
    hn.Init()

    n.Add("h", hn)

    en := new(Node)
    en.Init()
    en.Is_leaf = true

    hn.Add("e", en)

    ln := new(Node)
    ln.Init()
    rn := new(Node)
    rn.Init()
    rn.Is_leaf = true

    en.Add("r", rn)
    en.Add("l", ln)

    on := new(Node)
    on.Init()

    rn.Add("o", on)

    ln2 := new(Node)
    ln2.Init()
    ln.Add("l", ln2)

    on2 := new(Node)
    on2.Init()
    ln2.Add("o", on2)


    fmt.Println(Search(n, "he", 10))
    //p := new(Node)
    //n.Data["p"] = p
    //fmt.Println(n.Has_next())
    //fmt.Println(n.Get("p"))
    //fmt.Println(n.Get("g"))

    //fmt.Println(n.Get_the_top_node("he"))
}
