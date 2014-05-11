suggestion
==========

> 简单的输入框下拉提示服务

### 功能描述

在搜索输入框等位置,用户输入关键词,系统提示可以使用的关键字,提升用户体验

search

    植物

result:

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

具体代码直接查看源文件

### 实现方式1: easymap

使用map方式实现树结构,有python和golang两个版本

quick run:

```shell
git clone https://github.com/wklken/suggestion.git
cd suggestion/easymap
python suggest.py
go run suggest.go
```

适用:

    小型系统, 关键词在 10W 左右(中文+拼音+拼音首字母共30W左右)

优点:

    逻辑简单结构清晰, 代码足够少, 方便扩展
      e.g. 可自行修改存储结构,在返回中加入图片等额外信息

缺点:

    内存占用,30W关键词,平均词长3,占用800M内存
    另外对cpu也有一定消耗

处理和实践:

    加一层redis/memcached, python版本, 单机8进程, 16核, 占用1G内存, 每天总请求量在300-500w左右, qps峰值在 300 左右, 没什么压力[没做过压测....]
    golang版本完全没在生产上试过

### 实现方式2: double-array-trie

使用实现了double-array-trie的darts实现,golang代码

darts实现参考项目: [awsong/go-darts](https://github.com/awsong/go-darts)

double-array-trie文章: [What is Trie](http://en.wikipedia.org/wiki/Trie) | [An Implementation of Double-Array Trie](http://linux.thai.net/~thep/datrie/datrie.html)

quick run

      go run test_run.go

     input dict length: 29
     build out length 65708
     3.55us
     <nil>
     搜索: 植物大战
     Result Len: 10
     植物大战僵尸 154717704
     植物大战僵尸年度中文版 44592048
     植物大战僵尸OL 43566752
     植物大战僵尸2 630955
     植物大战外星人 530403
     植物大战怪兽 29727
     植物大战异形变态版 14773
     植物大战臭虫 5999
     植物大战异形 4456
     植物大战昆虫2无敌版 3419

适用:

    关键词在10w 以上的系统

优点:

    内存占用小, 性能保证

缺点:

    底层依赖double-array-trie,逻辑有点绕,自定义不是很方便

处理和实践:

    加一层redis/memcached

### 数据文件格式

默认文件格式:

    format:    word\tweight
    coding:    utf-8 [must]
    require:   weight type(int)

    eg:  植物大战僵尸\t1000


#### TODO

    集中于darts版本
    1.性能测试
    2.数据结构可自定义
    3.容错处理
    4.大小写,拼音,首字母等处理

#### Change Log


    2013-10-13 created, python版本
    2013-12-14 增加golang版本
    2014-05-11 增加double-array-trie实现的golang 版本

---------------


wklken(Pythonista/Vimer)

Email: wklken@yeah.net

Blog: http://www.wklken.me

Github: https://github.com/wklken

2013-10-13 于深圳
