package huffman

import (
	"bytes"
	"fmt"
	"sort"
)

type (
	node struct {
		char    rune   //字符
		weight  int    //权重
		left    *node  //左孩子
		right   *node  //右孩子
		parent  *node  //父节点
		visited bool   //访问过?
		huffman string //哈夫曼编码
	}

	nodelist []*node

	traverseFunc func(*node)

	HuffmanCoder struct {
		nodelist
		coder map[rune]string // 字符 --> 哈夫曼编码，用于编码
	}
)

const (
	vert_line         = "┃"
	start_left_child  = "┣━━"
	start_right_child = "┗━━"
	left_child        = '0'
	right_child       = '1'
)

func NewHuffmanCoder() *HuffmanCoder {
	coder := &HuffmanCoder{
		nodelist: []*node{},
		coder:    make(map[rune]string),
	}
	return coder
}

func (coder *HuffmanCoder) Append(char rune, weight int) {
	coder.nodelist.Append(char, weight)
	coder.coder = make(map[rune]string)
	coder.getCoder(coder.nodelist[0])
}

func (coder *HuffmanCoder) Encode(input string) (output string, err error) {
	for _, char := range input {
		if code := coder.coder[char]; len(code) > 0 {
			output += code
		} else {
			err = fmt.Errorf("不可识别的字符: %c", char)
			break
		}
	}
	return
}

func (coder *HuffmanCoder) getCoder(item *node) {
	if item != nil {
		if item.char > 0 {
			coder.coder[item.char] = item.huffman
		}
		coder.getCoder(item.left)
		coder.getCoder(item.right)
	}
}

//---------------------------------------------------------

func (list nodelist) Len() int {
	return len(list)
}

func (list nodelist) Swap(x, y int) {
	list[x], list[y] = list[y], list[x]
}

func (list nodelist) Less(x, y int) bool {
	return list[x].weight < list[y].weight
}

// 插入一个节点
func (list *nodelist) Append(char rune, weight int) {
	// 恢复原始节点
	list.reset()

	// 插入新节点
	item := &node{
		char:   char,
		weight: weight,
	}
	*list = append(*list, item)
	sort.Sort(list)

	// 生成编码
	list.huffmanCode()
}

// 输出编码树
func (list *nodelist) Tree() string {
	out := bytes.NewBuffer(nil)
	list.resetVisted()
	obj := *list
	for _, item := range obj {
		if !item.visited {
			item.tree(0, out)
		}
		out.WriteString("\n\n")
	}
	return out.String()
}

func (list *nodelist) Clear() {
	*list = []*node{}
}

func (list *nodelist) Decode(input string) (output string, err error) {

	if list.Len() == 0 {
		err = fmt.Errorf("空编码树")
		return
	} else if len(input) == 0 {
		err = fmt.Errorf("空输入")
		return
	}

	cur := (*list)[0]
	for _, char := range input {

		if char == left_child {
			cur = cur.left
		} else if char == right_child {
			cur = cur.right
		} else {
			err = fmt.Errorf("非法字符: %c", char)
			break
		}

		if cur == nil {
			err = fmt.Errorf("错误的输入，解码失败")
			break
		} else if cur.char > 0 {
			output += string(cur.char)
			cur = (*list)[0]
		}
	}

	return
}

//-------------------------------------------

// 生成编码
func (list *nodelist) huffmanCode() {
	for list.nextStep() {
	}

	if list.Len() > 0 {
		(*list)[0].huffmanCode("")
	}
}

func (list *nodelist) reset() {

	var newlist nodelist

	obj := *list
	for _, item := range obj {
		item.reset(&newlist)
	}

	*list = newlist
}

func (list *nodelist) nextStep() bool {

	if list.Len() < 2 {
		return false
	}

	obj := *list
	a := obj[0]
	b := obj[1]

	item := &node{
		weight: a.weight + b.weight,
		left:   a,
		right:  b,
	}
	a.parent, b.parent = item, item

	*list = append(obj[2:], item)
	sort.Sort(list)

	return true
}

func (list *nodelist) resetVisted() {
	obj := *list
	for _, item := range obj {
		item.traverse((*node).resetVisited)
	}
}

//----------------------------------

// 遍历
func (item *node) traverse(fn traverseFunc) {
	fn(item)
	if item.left != nil {
		item.left.traverse(fn)
	}
	if item.right != nil {
		item.right.traverse(fn)
	}
}

// 恢复叶子节点到节点列表中
func (item *node) reset(list *nodelist) {

	if item.char > 0 {
		*list = append(*list, &node{char: item.char, weight: item.weight})
	}

	if item.left != nil {
		item.left.reset(list)
	}

	if item.right != nil {
		item.right.reset(list)
	}
}

func (item *node) resetVisited() {
	item.visited = false
}

// 生成哈弗曼编码
func (item *node) huffmanCode(path string) {

	if item.char > 0 {
		item.huffman = path
	}

	if item.left != nil {
		item.left.huffmanCode(path + string(left_child))
	}

	if item.right != nil {
		item.right.huffmanCode(path + string(right_child))
	}
}

/*
	根定义为0层,根的子节点为1层...
*/
func (item *node) tree(depth int, out *bytes.Buffer) {

	var str string

	if item.parent != nil {
		cur := item.parent
		parent := cur.parent
		for parent != nil {
			if parent.left == cur {
				// 这一级的祖先是左孩子
				str = fmt.Sprintf("%-8s", vert_line) + str
			} else {
				// 这一级的祖先是右孩子
				str = fmt.Sprintf("%-8s", " ") + str
			}
			cur = parent
			parent = cur.parent
		}
	}

	if depth > 0 {
		if item.parent != nil && item.parent.left == item && item.parent.right != nil {
			str += start_left_child
		} else {
			str += start_right_child
		}
	}

	if item.char > 0 {
		str += fmt.Sprintf("权重:%-4d 字符:%c 编码:%s", item.weight, item.char, item.huffman)
	} else {
		str += fmt.Sprintf("权重:%-4d", item.weight)
	}

	item.visited = true

	out.WriteString(str)
	out.WriteString("\n")

	// 递归处理左右孩子
	if item.left != nil {
		item.left.tree(depth+1, out)
	}

	if item.right != nil {
		item.right.tree(depth+1, out)
	}
}
