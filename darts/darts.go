package darts

/*
TODO:
关注： 1.性能测试   2.消耗内存
代码行数变成 300 行，并且优化
1. 大小写的开关 构建字典,key全部转成小写,匹配的时候,也转成小写

DONE:
1. 当前，只有整体词出现在词典中时，才会有结果，含部分词的不行，需要处理 [DONE]
2. 读词的时候，需要去除左右空白  [DONE]
3. 对空行，要懂得跳过，而不是报错 [DONE]

*/

/*
1. 加注释(核心算法和原算法）
2. 保留核心算法，删除无用函数
3. 元数据和排序方法，变成可配置的
4. 重构部分代码的写法, 思考更好的方式，更节约内存的方式，以及golang的常用方法
5. 增加web端
6. 测试10w词，处理
7. 测试拼音+拼音首字母，处理
8. 加入logging模块
9. 增加缓存支持，memcached or redis
10. 发布版本
11. 看下源代码，确认实现方式是否一致
http://www.chasen.org/~taku/software/darts/
https://code.google.com/p/darts-clone/w/list
*/

// import
import (
	"bufio"
	"encoding/gob"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ==========================================================
// 元数据
type Term struct {
	Id    int //唯一的id，插入时生成
	Key   []rune
	Value int
	//可以包含其他需要的自定义信息
}

/*
元数据排序方法
*/

//  - 底层支持方法: 按权重排序 - 在对结果排序时使用
type Terms []Term

func (termSlice Terms) Len() int { return len(termSlice) }

func (termSlice Terms) Less(i, j int) bool { return termSlice[i].Value > termSlice[j].Value }

func (termSlice Terms) Swap(i, j int) { termSlice[i], termSlice[j] = termSlice[j], termSlice[i] }

//   - 底层支持方法: 按字典序排序 - 在构建字典时使用
type literalTerms []Term

func (termSlice literalTerms) Len() int { return len(termSlice) }
func (termSlice literalTerms) Less(i, j int) bool {
	var l int
	if len(termSlice[i].Key) < len(termSlice[j].Key) {
		l = len(termSlice[i].Key)
	} else {
		l = len(termSlice[j].Key)
	}

	for m := 0; m < l; m++ {
		if termSlice[i].Key[m] < termSlice[j].Key[m] {
			return true
		} else if termSlice[i].Key[m] == termSlice[j].Key[m] {
			continue
		} else {
			return false
		}
	}
	if len(termSlice[i].Key) < len(termSlice[j].Key) {
		return true
	}
	return false
}
func (termSlice literalTerms) Swap(i, j int) { termSlice[i], termSlice[j] = termSlice[j], termSlice[i] }

// ==========================================================
type node struct {
	code               rune /*Key_type*/
	depth, left, right int
	key                string
}

type Darts struct {
	Base             []int
	Check            []int
	KeyCount         int
	TermList         []map[int]bool // 存储所有包含到这个前缀的后缀id
	KeyString2IntMap map[string]int // 存储字面值 - ID(int) 的映射 多个字面值-同一个ID
	KeyInt2InfoMap   map[int]Term   // 存储 id - detailInfo 的映射
	Used             []int
}

// 构建过程中需要存储一些中间状态数据
type dartsBuild struct {
	darts        Darts
	size         int
	key          [][]rune /*Key_type*/
	nextCheckPos int
	err          int
}

// ==========================================================
// for Darts Build
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// variable key should be sorted ascendingly
func Build(key [][]rune /*Key_type*/, keyString2IntMap map[string]int, keyInt2InfoMap map[int]Term) Darts {
	var d = new(dartsBuild)

	d.key = key
	d.resize(512)

	d.darts.Base[0] = 1
	d.darts.KeyInt2InfoMap = keyInt2InfoMap
	d.darts.KeyString2IntMap = keyString2IntMap
	d.darts.KeyCount = 0

	d.nextCheckPos = 0

	var rootNode node
	rootNode.depth = 0
	rootNode.left = 0
	rootNode.right = len(key)

	siblings := d.fetch(rootNode)
	// 从root的第一层子节点开始，执行递归插入
	d.insert(siblings)

	if d.err < 0 {
		panic("Build error")
	}
	return d.darts
}

/*
重新分配空间
*/
func (d *dartsBuild) resize(newSize int) {
	// fmt.Println("DEBUG - resize", newSize, len(d.darts.Base), cap(d.darts.Base))
	// fmt.Println("DEBUG - resize", newSize, len(d.darts.Check), cap(d.darts.Check))
	// fmt.Println("DEBUG - resize", newSize, len(d.darts.TermList), cap(d.darts.TermList))
	// fmt.Println("DEBUG - resize", newSize, len(d.darts.Used), cap(d.darts.Used))

	if newSize > cap(d.darts.Base) {
		d.darts.Base = append(d.darts.Base, make([]int, (newSize-len(d.darts.Base)))...)
		d.darts.Check = append(d.darts.Check, make([]int, (newSize-len(d.darts.Check)))...)
		d.darts.TermList = append(d.darts.TermList, make([]map[int]bool, (newSize-len(d.darts.TermList)))...)
		d.darts.Used = append(d.darts.Used, make([]int, (newSize-len(d.darts.Used)))...)
	} else {
		d.darts.Base = d.darts.Base[:newSize]
		d.darts.Check = d.darts.Check[:newSize]
		d.darts.TermList = d.darts.TermList[:newSize]
		d.darts.Used = d.darts.Used[:newSize]
	}
}

/*
获取第一层子节点
*/
func (d *dartsBuild) fetch(parent node) []node {
	var siblings = make([]node, 0, 2)
	if d.err < 0 {
		return siblings[0:0]
	}

	var prev rune = /*Key_type*/ 0

	for i := parent.left; i < parent.right; i++ {
		if len(d.key[i]) < parent.depth {
			continue
		}

		tmp := d.key[i]

		var cur rune = /*Key_type*/ 0
		if len(d.key[i]) != parent.depth {
			cur = tmp[parent.depth] + 1
		}

		if prev > cur {
			fmt.Println(prev, cur, i, parent.depth, d.key[i])
			fmt.Println(d.key[i])
			panic("fetch error 1")
			d.err = -3
			return siblings[0:0]
		}

		if cur != prev || len(siblings) == 0 {
			var tmpNode node
			tmpNode.depth = parent.depth + 1
			tmpNode.code = cur
			tmpNode.key = string(d.key[i])
			tmpNode.left = i
			if len(siblings) != 0 {
				siblings[len(siblings)-1].right = i
			}

			siblings = append(siblings, tmpNode)
		}

		prev = cur
	}

	if len(siblings) != 0 {
		// ? 做什么用的
		siblings[len(siblings)-1].right = parent.right
	}

	return siblings
}

/**
递归插入方法 - 深度优先
*/
func (d *dartsBuild) insert(siblings []node) (int, []int) {
	var begin int = 0
	var keys []int = make([]int, 0)

	if d.err < 0 {
		panic("insert error")
		return begin, keys
	}

	var pos int = max(int(siblings[0].code)+1, d.nextCheckPos) - 1
	var nonZeroNum int = 0
	first := false

	// 如果空间不足，补足空间
	if len(d.darts.Base) <= pos {
		d.resize(pos + 1)
	}

	for {
	next:
		pos++

		// 如果空间不足，补足空间
		if len(d.darts.Base) <= pos {
			d.resize(pos + 1)
		}

		if d.darts.Check[pos] > 0 {
			nonZeroNum++
			continue
		} else if !first {
			d.nextCheckPos = pos
			first = true
		}

		begin = pos - int(siblings[0].code)
		if len(d.darts.Base) <= (begin + int(siblings[len(siblings)-1].code)) {
			d.resize(begin + int(siblings[len(siblings)-1].code) + 400)
		}

		if d.darts.Used[begin] == 1 {
			continue
		}

		for i := 1; i < len(siblings); i++ {
			if begin+int(siblings[i].code) >= len(d.darts.Base) {
				fmt.Println(len(d.darts.Base), begin+int(siblings[i].code), begin+int(siblings[len(siblings)-1].code))
			}
			if 0 != d.darts.Check[begin+int(siblings[i].code)] {
				goto next
			}
		}
		break
	}

	if float32(nonZeroNum)/float32(pos-d.nextCheckPos+1) >= 0.95 {
		d.nextCheckPos = pos
	}
	// d.darts.Used[begin] = true
	d.darts.Used[begin] = 1
	d.size = max(d.size, begin+int(siblings[len(siblings)-1].code)+1)

	for i := 0; i < len(siblings); i++ {
		d.darts.Check[begin+int(siblings[i].code)] = begin
		keys = append(keys, d.darts.KeyString2IntMap[siblings[i].key])
	}

	// 递归处理下一层
	for i := 0; i < len(siblings); i++ {
		newSiblings := d.fetch(siblings[i])

		// 如果是叶节点，多一个节点end，其base为负
		if len(newSiblings) == 0 {
			d.darts.Base[begin+int(siblings[i].code)] = -d.darts.KeyCount - 1
			d.darts.KeyCount++
		} else {
			// 非叶节点，递归调用
			h, subkeys := d.insert(newSiblings)
			d.darts.Base[begin+int(siblings[i].code)] = h

			// 将这个节点的所有后缀节点的key加入到termlist
			for j := 0; j < len(subkeys); j++ {
				keys = append(keys, subkeys[j])
			}
		}
	}

	for i := 0; i < len(keys); i++ {
		termMap := d.darts.TermList[begin]
		if termMap == nil {
			d.darts.TermList[begin] = make(map[int]bool)
		}

		// TODO: sort, 保留最低个数, 每次插入要比较下用哪个

		d.darts.TermList[begin][keys[i]] = true
		// 在这里需要控制个数，即，保存子节点的个数，只存top
		// 如果能存order，则可以去除掉查询结果的sort方法
	}

	return begin, keys
}

// ==================================================
// for Darts

/*
精确查找某个字符串是否存在
*/
func (d Darts) ExactMatch(key []rune /*Key_type*/, nodePos int) bool {
	b := d.Base[nodePos]
	var p int

	for i := 0; i < len(key); i++ {
		p = b + int(key[i]) + 1
		if b == d.Check[p] {
			b = d.Base[p]
		} else {
			return false
		}
	}

	p = b
	n := d.Base[p]
	if b == d.Check[p] && n < 0 {
		return true
	}

	return false
}

/*
搜索，获取以某个字符串开始的所有子节点
*/
func (d Darts) Search(key []rune /*Key_type*/, nodePos int) (results Terms) {

	b := d.Base[nodePos]
	var p int

	// 定位到最后一个key节点
	for i := 0; i < len(key); i++ {
		//p = b

		p = b + int(key[i]) + 1
		// p could be bigger than the index
		if p >= len(d.Check) {
			return results
		}

		if b == d.Check[p] {
			b = d.Base[p]
		} else {
			return results
		}
	}

	// 获取节点的所有key列表
	for k, _ := range d.TermList[b] {
		// TODO: trigger of if match self
		// if string(d.KeyInt2InfoMap[k].Key) == string(key) {
		// continue
		// }
		results = append(results, d.KeyInt2InfoMap[k])
	}

	// 排序,返回
	sort.Sort(results)
	return results
}

// ======================================================

// ======================================================

/*
读入文件，构建double-array-trie
*/
func Import(inFile, outFile string) (Darts, error) {
	// 输入文件
	unifile, erri := os.Open(inFile)
	if erri != nil {
		return Darts{}, erri
	}
	defer unifile.Close()

	// 存储到本地文件一份，后续支持直接load
	ofile, erro := os.Create(outFile)
	if erro != nil {
		return Darts{}, erro
	}
	defer ofile.Close()

	// 读取文件，对所有key进行排序
	terms := make(literalTerms, 0, 130000)
	uniLineReader := bufio.NewReaderSize(unifile, 400)
	line, _, bufErr := uniLineReader.ReadLine()
	// TODO: 加容错
	for nil == bufErr {
		rst := strings.Split(string(line), "\t")

		// 跳过异常的行
		if len(rst) < 2 {
			line, _, bufErr = uniLineReader.ReadLine()
			continue
		}

		// 去除左右空白
		key := []rune(strings.TrimSpace(rst[0]))
		value, _ := strconv.Atoi(rst[1])

		terms = append(terms, Term{0, key, value})

		line, _, bufErr = uniLineReader.ReadLine()
	}
	sort.Sort(terms)

	// 获取排序后的keys数组和values数组
	keys := make([][]rune, len(terms))
	values := make([]int, len(terms))

	//建造映射 提供排序时的返回
	keyString2IntMap := make(map[string]int)
	keyInt2InfoMap := make(map[int]Term)

	for i := 0; i < len(terms); i++ {
		keys[i] = terms[i].Key
		values[i] = terms[i].Value

		terms[i].Id = i
		keyString2IntMap[string(terms[i].Key)] = i
		keyInt2InfoMap[i] = terms[i]
	}

	// 开始构建darts
	fmt.Printf("input dict length: %v\n", len(keys))
	round := len(keys)
	var d Darts
	d = Build(keys[:round], keyString2IntMap, keyInt2InfoMap)

	// 确保所有key都被加入
	fmt.Printf("build out length %v\n", len(d.Base))
	t := time.Now()
	for i := 0; i < round; i++ {
		if true != d.ExactMatch(keys[i], 0) {
			err := fmt.Errorf("missing key %s, %v, %d, %v, %v", string(keys[i]), keys[i], i, keys[i-1], keys[i+1])
			return d, err
		}
	}
	fmt.Println(time.Since(t))

	// 写入本地文件
	enc := gob.NewEncoder(ofile)
	enc.Encode(d)

	// 返回
	return d, nil
}

/*
从已有文件直接加载到内存
*/
func Load(filename string) (Darts, error) {
	var dict Darts
	file, err := os.Open(filename)
	if err != nil {
		return Darts{}, err
	}
	defer file.Close()

	dec := gob.NewDecoder(file)
	dec.Decode(&dict)
	return dict, nil
}
