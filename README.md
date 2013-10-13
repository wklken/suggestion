suggestion
==========

### 简介

下拉提示：在搜索框输入一个词，根据输入词匹配前缀，下拉框提示搜索系统中有的搜索词

具体方法和说明可以查看源代码，200行左右

python 版本已经完成，实际环境线上30w key，不加缓存情况下占用内存800M，开启缓存(用memcached或redis或自带)，可以有效提升响应速度，降低cpu load

golang 版本书写中(golang盲，开始翻书中)

### 使用

####测试1

在代码中建立树

    n = Node("")
    #default weight=0, 后面的参数可以任意加,搜索返回结果再从node中将放入对应的值取出,这里放入一个othervalue值
    add(n, u'he', othervalue="v-he")
    add(n, u'her', weight=0, othervalue="v-her")
    add(n, u'hero', weight=10, othervalue="v-hero")
    add(n, u'hera', weight=3, othervalue="v-hera")

进行搜索

    for key, node in search(n, u'h'):
        print key, node, node.othervalue, id(node)

结果

    search h:
    hero <Node key:o is_leaf:True weight:10 Subnodes: []> v-hero 140563304390448
    hera <Node key:a is_leaf:True weight:3 Subnodes: []> v-hera 140563304390752
    he <Node key:e is_leaf:True weight:0 Subnodes: [(u'r', {u'a': {}, u'o': {}})]> v-he 140563304376768
    her <Node key:r is_leaf:True weight:0 Subnodes: [(u'a', {}), (u'o', {})]> v-her 140563304377808

###测试2

读取数据文件建立树

    tree = build("./test_data", is_case_sensitive=False)

搜索

    print u'search 植物'
    for key, node in search(tree, u'植物', limit=10):
        print key, node.weight

结果

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

### 参数和函数的说明

1. 可以查看源代码中对应参数和方法的注解

2. 查看docstring

        Python 2.7.5 (default, Aug 25 2013, 00:04:04)
        [GCC 4.2.1 Compatible Apple LLVM 5.0 (clang-500.0.68)] on darwin
        Type "help", "copyright", "credits" or "license" for more information.
        >>> import suggest
        >>> help(suggest)
        Help on module suggest:
        ......

3. 简要说明

        CACHED = True #是否开启默认缓存
        CACHED_SIZE = 10 #缓存大小
        CACHED_THREHOLD = 10 #节点被搜索超过一定次数时，才会加缓存

        class Node(dict)  #节点对象
        def depth_walk(node) #深度遍历节点的方法

        def add(node, keyword, weight=0, **kwargs) #给树加一个前缀关键词
        def delete(node, keyword, judge_leaf=False) #删除树里面的某个前缀
        def search(node, prefix, limit=None, is_case_sensitive=False) #搜索方法，其中is_case_sensitive=False时，将会把prefx转为小写进行遍历
                                                                    #所以，建立索引的时候，需要确认放到树中的，大小写
                                                                    #search和build的is_case_sensitive保持一致
        def build(file_path, is_case_sensitive=False) #从数据文件建立树的方法，可以根据自己数据文件格式重定义


--------

The End!

wklken (凌岳/Pythonista/vim党党员)

Email: wklken@yeah.net

Github: https://github.com/wklken

Blog: http://wklken.me

2013-10-13 于深圳


