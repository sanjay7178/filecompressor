// compress/huffman.go
package compress

import (
	"container/heap"
	"errors"
	// "encoding/binary"
)

type HuffmanNode struct {
	Char  byte
	Freq  int
	Left  *HuffmanNode
	Right *HuffmanNode
}

type HuffmanHeap []*HuffmanNode

func (h HuffmanHeap) Len() int            { return len(h) }
func (h HuffmanHeap) Less(i, j int) bool  { return h[i].Freq < h[j].Freq }
func (h HuffmanHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *HuffmanHeap) Push(x interface{}) { *h = append(*h, x.(*HuffmanNode)) }
func (h *HuffmanHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type HuffmanCompressor struct{}

func NewHuffmanCompressor() *HuffmanCompressor {
	return &HuffmanCompressor{}
}

func (hc *HuffmanCompressor) buildTree(data []byte) *HuffmanNode {
	freqs := make(map[byte]int)
	for _, b := range data {
		freqs[b]++
	}

	h := &HuffmanHeap{}
	heap.Init(h)

	for char, freq := range freqs {
		heap.Push(h, &HuffmanNode{Char: char, Freq: freq})
	}

	for h.Len() > 1 {
		left := heap.Pop(h).(*HuffmanNode)
		right := heap.Pop(h).(*HuffmanNode)
		heap.Push(h, &HuffmanNode{
			Freq:  left.Freq + right.Freq,
			Left:  left,
			Right: right,
		})
	}

	return heap.Pop(h).(*HuffmanNode)
}

func (hc *HuffmanCompressor) buildCodes(node *HuffmanNode, code string, codes map[byte]string) {
	if node == nil {
		return
	}
	if node.Left == nil && node.Right == nil {
		codes[node.Char] = code
		return
	}
	hc.buildCodes(node.Left, code+"0", codes)
	hc.buildCodes(node.Right, code+"1", codes)
}

func (hc *HuffmanCompressor) serializeTree(node *HuffmanNode) []byte {
	if node == nil {
		return []byte{}
	}

	if node.Left == nil && node.Right == nil {
		return []byte{1, node.Char} // Leaf node: 1 + character
	}

	// Internal node: 0 + left subtree + right subtree
	left := hc.serializeTree(node.Left)
	right := hc.serializeTree(node.Right)
	result := append([]byte{0}, left...)
	return append(result, right...)
}

func (hc *HuffmanCompressor) deserializeTree(data []byte) (*HuffmanNode, int) {
	if len(data) == 0 {
		return nil, 0
	}

	if data[0] == 1 {
		return &HuffmanNode{Char: data[1]}, 2
	}

	left, leftSize := hc.deserializeTree(data[1:])
	right, rightSize := hc.deserializeTree(data[1+leftSize:])

	return &HuffmanNode{Left: left, Right: right}, 1 + leftSize + rightSize
}

func (hc *HuffmanCompressor) Compress(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}

	tree := hc.buildTree(data)
	codes := make(map[byte]string)
	hc.buildCodes(tree, "", codes)

	// Serialize tree structure
	treeBytes := hc.serializeTree(tree)

	// Build compressed data
	var bits string
	for _, b := range data {
		bits += codes[b]
	}

	// Convert bits to bytes
	compressed := make([]byte, len(treeBytes)+1+(len(bits)+7)/8)
	compressed[0] = byte(len(treeBytes))
	copy(compressed[1:], treeBytes)

	offset := len(treeBytes) + 1
	for i := 0; i < len(bits); i += 8 {
		end := i + 8
		if end > len(bits) {
			end = len(bits)
		}
		byteStr := bits[i:end]
		val := 0
		for j := 0; j < len(byteStr); j++ {
			if byteStr[j] == '1' {
				val |= 1 << uint(7-j)
			}
		}
		compressed[offset+i/8] = byte(val)
	}

	return compressed, nil
}

func (hc *HuffmanCompressor) Decompress(compressed []byte) ([]byte, error) {
	if len(compressed) == 0 {
		return nil, nil
	}

	// Read tree size and deserialize tree
	treeSize := int(compressed[0])
	tree, _ := hc.deserializeTree(compressed[1 : treeSize+1])

	// Read compressed data
	var result []byte
	node := tree

	for i := treeSize + 1; i < len(compressed); i++ {
		byte := compressed[i]
		for bit := 7; bit >= 0; bit-- {
			if node.Left == nil && node.Right == nil {
				result = append(result, node.Char)
				node = tree
			}

			if (byte & (1 << uint(bit))) != 0 {
				node = node.Right
			} else {
				node = node.Left
			}

			if node == nil {
				return nil, errors.New("invalid compressed data")
			}
		}
	}

	// Handle last character
	if node != nil && node.Left == nil && node.Right == nil {
		result = append(result, node.Char)
	}

	return result, nil
}
