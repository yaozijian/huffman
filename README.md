
哈夫曼编解码库

```
package main

import (
	"fmt"

	"github.com/yaozijian/huffman"
)

func main() {

	coder := huffman.NewHuffmanCoder()

	coder.Append('A', 5)
	coder.Append('B', 4)
	coder.Append('C', 3)
	coder.Append('D', 2)
	coder.Append('E', 1)

	fmt.Println(coder.Tree())

	in := "ABCDE"
	out, err := coder.Encode(in)
	fmt.Printf("编码输入: %s 输出: %s 错误: %v\n", in, out, err)

	in = out
	out, err = coder.Decode(in)
	fmt.Printf("解码输入: %s 输出: %s 错误: %v\n", in, out, err)
}
```