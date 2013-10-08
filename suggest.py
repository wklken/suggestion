#!/usr/bin/env python
# encoding: utf-8
"""
@date:    20131001
@version: 0.2
@author:  wklken@yeah.net
@desc:    搜索下拉提示，基于后台提供数据，建立数据结构(trie)，用户输入query前缀时，可以提示对应query前缀补全

@update:
    20131001 基本结构，新增，搜索等基本功能
    20131005 增加缓存功能，当缓存打开，用户搜索某个前缀超过一定次数时，进行缓存，减少搜索时间
    20131006 增加PuppetNode,可以自定义返回节点中内容，用户存放下拉的其他属性，例如图片，分类等

@TODO:
    test case

    返回leaf node放在result里面，并且里面有cache，这样合适？
"""

#这是实现cache的一种方式，也可以使用redis/memcached在外部做缓存

#一旦打开，search时会对每个节点做cache，当增加删除节点时，其路径上的cache会被清除,搜索时间降低了一个数量级
#代价：内存消耗, 不需要时可以关闭,或者通过CACHED_THREHOLD调整缓存数量
CACHED = True
#CACHED = False
#注意，CACHED_SIZE >= search中的limit，保证search从缓存能获取到足够多的结果
CACHED_SIZE = 10
#被搜索超过多少次后才加入缓存
CACHED_THREHOLD = 1

############### start ######################
class PuppetNode(object):
    pass


class Node(dict):
    def __init__(self, key, is_leaf=False, weight=0):
        #节点字符
        self.key = key
        #是否叶子节点
        self.is_leaf = is_leaf
        #节点权重, 某个词最后一个字节点代表其权重，其余中间节点权重为0，无意义
        self.weight = weight

        #缓存
        self.cache = []
        #节点前缀搜索次数，可以用于搜索query数据分析
        self.search_count = 0

    @staticmethod
    def get_puppet_fields():
        """
        傀儡节点拷贝的项
        NOTICE:可以加入其他属性，例如图标，分类等在展示时需要用到的，根据需要修改数据文件格式和build方法
        """
        return ['weight']

    def make_puppet_node(self):
        """
        最终返回值,防止节点cache字段嵌套导致的问题,仅配置需要用到的属性值
        """
        n = PuppetNode()
        for field in Node.get_puppet_fields():
            setattr(n, field, getattr(self, field))
        return n


    def __str__(self):
        return '<Node %s %s> %s' % (self.key, self.is_leaf, self.items())


    def add_subnode(self, node):
        """
        添加子节点
        """
        self.update({node.key: node})


    def get_subnode(self, key):
        """
        获取子节点
        """
        return self.get(key)


    def has_subnode(self):
        """
        判断是否存在子节点
        """
        return len(self) > 0


    def get_top_node(self, prefix):
        """
        获取一个前缀的最后一个节点，相当于补全的顶部节点
        """
        top = self

        for k in prefix:
            top = top.get_subnode(k)
            if top is None:
                return None

            if top.has_subnode():
                continue
            else:
                break
        return top


def depth_walk(node):
    """
    递归，深度优先遍历,返回每个节点所代表的key以及叶子节点的权重
    """
    result = []
    if node.is_leaf:
        result.append(('', node.make_puppet_node()))

    if node.has_subnode():
        for k in node.iterkeys():
            s = depth_walk(node.get(k))
            result.extend([(k+name, n) for name,n in s])
        return result
    else:
        return [('', node.make_puppet_node())]


def search(node, prefix, limit=None):
    """
    搜索一个前缀下的所有单词列表 递归式
    """
    node = node.get_top_node(prefix)
    if node is None:
        return []

    node.search_count += 1

    if CACHED and node.cache:
        print '...................read from cache'
        return node.cache[:limit] if limit is not None else node.cache

    result = []

    result.extend([(prefix + i, n) for i, n in depth_walk(node)])

    result.sort(key=lambda x: -x[1].weight)

    if CACHED and node.search_count >= CACHED_THREHOLD:
        print '...................write to cache'
        node.cache = result[:CACHED_SIZE]

    return result[:limit] if limit is not None else result

def add(node, keyword, weight=0):
    """
    加入一个单词
    """
    one_node = node

    index = 0
    last_index = len(keyword) - 1
    for c in keyword:
        if c not in one_node:
            if index != last_index:
                one_node.add_subnode(Node(c, weight=weight))
            else:
                one_node.add_subnode(Node(c, is_leaf=True, weight=weight))
            one_node = one_node.get_subnode(c)
        else:
            one_node = one_node.get_subnode(c)

            if CACHED:
                one_node.cache = []

            if index == last_index:
                one_node.is_leaf = True
                one_node.weight = weight
        index += 1

def delete(node, keyword, judge_leaf=False):
    """
    删除一个单词
    """
    # 空关键词，传入参数有问题，或者递归调用到了根节点,直接返回
    if not keyword:
        return

    top_node = node.get_top_node(keyword)

    #清理缓存
    if CACHED:
        top_node.cache = []

    #递归往上，遇到节点是某个关键词节点时，要退出
    if judge_leaf:
        if top_node.is_leaf:
            return
    #非递归，调用delete
    else:
        if not top_node.is_leaf:
            return

    #存在子节点，去除其标志 done
    if top_node.has_subnode():
        top_node.is_leaf = False
        return
    #不存在子节点，逐层检查删除节点
    else:
        this_node = top_node

        prefix = keyword[:-1]
        top_node = node.get_top_node(prefix)
        del top_node[this_node.key]
        delete(node, prefix, judge_leaf=True)


##############################
#  增补功能 读数据文件建立树 #
##############################

def build(file_path, is_case_sensitive=False):
    """
    从文件构建数据结构, 文件必须utf-8编码
    """
    node = Node("")
    f = open(file_path)
    for line in f:
        line = line.strip()
        if not isinstance(line,unicode):
            line = line.decode('utf-8')
        parts = line.split('\t')
        name = parts[0]
        if not is_case_sensitive:
            name = name.lower()
        add(node, name, int(parts[1]))
    f.close()
    return node


if __name__ == '__main__':
    print "TRACK ================= test 1"
    node = Node("")

    hn = Node("h")
    node.add_subnode(hn)

    an = Node("a")
    en = Node("e", is_leaf=True)
    hn.add_subnode(an)
    hn.add_subnode(en)

    ln = Node("l")
    an.add_subnode(ln)

    on = Node("o")
    ln.add_subnode(on)

    ln1 = Node("l")
    rn = Node("r", is_leaf=True)
    en.add_subnode(ln1)
    en.add_subnode(rn)

    ln2 = Node("l")
    an2 = Node("a")
    ln1.add_subnode(ln2)
    ln1.add_subnode(an2)

    on2 = Node("o", is_leaf=True)
    ln2.add_subnode(on2)
    bn = Node("b", is_leaf=True)
    an2.add_subnode(bn)

    on1 = Node("o", is_leaf=True)
    rn.add_subnode(on1)

    add(node, 'hec')
    add(node, 'helac')
    add(node, 'hec')

    #print get_top_node(node, 'he')
    print search(node, 'he')

    print "TRACK ================= test11"

    n = Node("")
    add(n, u'h')
    add(n, u'he')
    add(n, u'her')

    print 'test_get top node'
    print n.get_top_node(u'her')
    print n.get_top_node(u'he')
    print n.get_top_node(u'h')

    print "TRACK ================= test2"

    n = build('./game_info')
    for v,w in search(n, u'植物'):
        print v, w.weight
    print 'len', len(search(n, u'植物'))

    import sys
    print 'total size:', sys.getsizeof(n)

    for v,w in search(n, u'植物'):
        print v, w.weight
    search(n, u'植物')
    print 'len', len(search(n, u'植物'))

    import time
    b = time.clock()
    search(n, u'植物')
    print time.clock() - b

    print "TRACK ================= test3"

    n = Node("")
    add(n, u'h')
    add(n, u'he')
    add(n, u'her')

    print 'before delete', n
    for v, w in search(n, u'h'):
        print v, w

    print "TRACK ================= delete begin"
    delete(n, u'he')

    print 'after delete', n, n is None
    for v, w in search(n, u'h'):
        print v, w.weight

    #print 'delete h-done he her heer'
    #print "TRACK ================= after delete"

    #for v, w in search(n, u'he'):
        #print v, w


