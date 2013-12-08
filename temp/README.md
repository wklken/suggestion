suggestion
==========

### 简介

下拉提示：在搜索框输入一个词，根据输入词匹配前缀，下拉框提示搜索系统中有的搜索词

具体方法和说明可以查看源代码，200行左右,可以根据自己需求任意修改

python 版本已经完成，实际环境线上30w key，不加缓存情况下占用内存800M，每天百万级别请求毫无压力(8核/16G/双进程起服务)

开启缓存(用memcached或redis或自带)，可以有效提升响应速度，降低cpu load

golang 版本书写中(golang盲，开始翻书中)

### 使用

####测试1

在代码中建立树



结果


###测试2

读取数据文件建立树


数据文件格式

    关键词\t权重   且存成utf-8编码

    植物大战僵尸\t11111

搜索


结果


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

TODO:


--------

The End!

wklken (凌岳/Pythonista/vim党党员)

Email: wklken@yeah.net

Github: https://github.com/wklken

Blog: http://wklken.me

2013-10-13 于深圳


