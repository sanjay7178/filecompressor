// compress/shannon_fano.go
package compress

import (
	"errors"
	"sort"
	// "bytes"
)

type SFNode struct {
	Symbol byte
	Freq   int
	Code   string
}

type ShannonFanoCompressor struct{}

func NewShannonFanoCompressor() *ShannonFanoCompressor {
	return &ShannonFanoCompressor{}
}

func (sf *ShannonFanoCompressor) buildFrequencyTable(data []byte) []*SFNode {
	freqs := make(map[byte]int)
	for _, b := range data {
		freqs[b]++
	}

	nodes := make([]*SFNode, 0)
	for sym, freq := range freqs {
		nodes = append(nodes, &SFNode{Symbol: sym, Freq: freq})
	}

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Freq > nodes[j].Freq
	})

	return nodes
}

func (sf *ShannonFanoCompressor) divide(nodes []*SFNode, start, end int) {
	if start >= end {
		return
	}

	if start+1 == end {
		nodes[start].Code += "0"
		return
	}

	total := 0
	for i := start; i < end; i++ {
		total += nodes[i].Freq
	}

	sum := 0
	diff := total
	mid := start

	for i := start; i < end; i++ {
		sum += nodes[i].Freq
		if abs(2*sum-total) < diff {
			diff = abs(2*sum - total)
			mid = i
		}
	}

	for i := start; i <= mid; i++ {
		nodes[i].Code += "0"
	}
	for i := mid + 1; i < end; i++ {
		nodes[i].Code += "1"
	}

	sf.divide(nodes, start, mid+1)
	sf.divide(nodes, mid+1, end)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (sf *ShannonFanoCompressor) Compress(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}

	// Build frequency table and generate codes
	nodes := sf.buildFrequencyTable(data)
	sf.divide(nodes, 0, len(nodes))

	// Create lookup table
	codeTable := make(map[byte]string)
	for _, node := range nodes {
		codeTable[node.Symbol] = node.Code
	}

	// Convert data to bits
	var bits string
	for _, b := range data {
		bits += codeTable[b]
	}

	// Create header with symbol table
	header := make([]byte, 0)
	header = append(header, byte(len(nodes))) // Number of symbols

	for _, node := range nodes {
		header = append(header, node.Symbol)
		header = append(header, byte(len(node.Code)))

		// Convert code string to bytes
		codeBytes := make([]byte, (len(node.Code)+7)/8)
		for i := 0; i < len(node.Code); i++ {
			if node.Code[i] == '1' {
				codeBytes[i/8] |= 1 << uint(7-(i%8))
			}
		}
		header = append(header, codeBytes...)
	}

	// Convert compressed bits to bytes
	compressed := make([]byte, (len(bits)+7)/8)
	for i := 0; i < len(bits); i++ {
		if bits[i] == '1' {
			compressed[i/8] |= 1 << uint(7-(i%8))
		}
	}

	// Combine header and compressed data
	result := make([]byte, len(header)+len(compressed)+1)
	result[0] = byte(len(header))
	copy(result[1:], header)
	copy(result[1+len(header):], compressed)

	return result, nil
}

func (sf *ShannonFanoCompressor) Decompress(compressed []byte) ([]byte, error) {
	if len(compressed) == 0 {
		return nil, nil
	}

	// Read header
	headerLen := int(compressed[0])
	if headerLen+1 >= len(compressed) {
		return nil, errors.New("invalid compressed data")
	}

	header := compressed[1 : headerLen+1]
	numSymbols := int(header[0])
	pos := 1

	// Rebuild code table
	codeTable := make(map[string]byte)
	for i := 0; i < numSymbols && pos < len(header); i++ {
		symbol := header[pos]
		pos++
		codeLen := int(header[pos])
		pos++

		codeBytes := header[pos:(pos + (codeLen+7)/8)]
		var code string
		for bit := 0; bit < codeLen; bit++ {
			if (codeBytes[bit/8] & (1 << uint(7-(bit%8)))) != 0 {
				code += "1"
			} else {
				code += "0"
			}
		}
		codeTable[code] = symbol
		pos += (codeLen + 7) / 8
	}

	// Decompress data
	data := compressed[headerLen+1:]
	var result []byte
	currentCode := ""

	for i := 0; i < len(data)*8; i++ {
		if (data[i/8] & (1 << uint(7-(i%8)))) != 0 {
			currentCode += "1"
		} else {
			currentCode += "0"
		}

		if symbol, ok := codeTable[currentCode]; ok {
			result = append(result, symbol)
			currentCode = ""
		}
	}

	return result, nil
}
