suggestion
==========

Simple search query suggest service support by python/golang

### Purpose

> I just want to build a simple and powerful search query suggest service.

> Python and Golang version

用Python/Golang实现最简单的搜索框下拉提示服务!

### Live Demo

Coming soon!

### quick run

    git clone https://github.com/wklken/suggestion.git
    cd suggestion
    python suggest.py
    go run suggest.go


### Tutorial

#### build tree in simple code and search

build:

    n = Node("")
    #default weight=0, 其他的参数可以任意加,搜索返回结果再从node中将放入对应的值取出,这里放入一个othervalue值
    add(n, u'he', othervalue="v-he")
    add(n, u'her', weight=0, othervalue="v-her")
    add(n, u'hero', weight=10, othervalue="v-hero")
    add(n, u'hera', weight=3, othervalue="v-hera")

search:

    for key, node in search(n, u'h'):
        #print key, node, node.othervalue, id(node)
        print key, node, node.weight

result:

    search h:
    hero 10
    hera 3
    he 0
    her 0

#### build with data file

file format:

    format:    words\tweight
    coding:    utf-8
    require:   weight type(int)

    eg:  植物大战僵尸\t1000

build tree:

    tree = build("./test_data", is_case_sensitive=False)

search

    print u'search 植物'
    for key, node in search(tree, u'植物', limit=10):
        print key, node.weight

result:

    search 植物
    植物大战僵尸 154717704
    植物大战僵尸年度中文版 44592048
    植物大战僵尸ol 43566752
    植物联盟 4244331
    植物大战僵尸2 630955
    植物大战外星人 530403
    植物小鸟大战僵尸 128907
    植物大作战 52909
    植物精灵 50475
    植物秘境：深入未知 43468

### Golang version

    The same logical as Python version.
    Just run:
        go run suggest.go

### Others

1. It's suggested that build your own cache layer through Memcached/Redis.

   That will help a lot!(much better,reduce the cpu load)

2. Kill the goose that how to use the wheel

   If your data word count less than 100,000 , Use the simple version.(Suggest)

   Simple version is easy enough to modify code your self.(It's quick enough,just take more Mem)

   Otherwise, Use the Double-Array Trie Tree version.(Less Mem, quick seach but the logical is a little complicated)

3. Encoding

   Make sure that your data file encoding is utf-8

### Performance

Coming soon!

### TODO

1. rebuild it with double-array trie tree

2. get the performance tree

3. change code to support weight type float/double

### Connect

If you have any suggestions or questions, Open an issue!

Also, pull requests!

I will check that weekly.

---------------

wklken(Pythonista/Vimer)

Email: wklken@yeah.net

Blog: http://www.wklken.me

Github: https://github.com/wklken


