// compress/lzw.go
package compress

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type Dictionary struct {
	entries  map[string]int
	nextCode int
}

func NewDictionary() *Dictionary {
	dict := &Dictionary{
		entries:  make(map[string]int),
		nextCode: 256, // Reserve first 256 codes for single bytes
	}

	// Initialize with single byte sequences
	for i := 0; i < 256; i++ {
		dict.entries[string([]byte{byte(i)})] = i
	}

	return dict
}

type LZWCompressor struct {
	dict *Dictionary
}

func NewLZWCompressor() *LZWCompressor {
	return &LZWCompressor{
		dict: NewDictionary(),
	}
}

func (lzw *LZWCompressor) Compress(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}

	dict := NewDictionary()
	current := []byte{}
	result := []int{}

	for _, b := range data {
		probe := append(current, b)

		if _, exists := dict.entries[string(probe)]; exists {
			current = probe
		} else {
			result = append(result, dict.entries[string(current)])
			dict.entries[string(probe)] = dict.nextCode
			dict.nextCode++
			current = []byte{b}
		}
	}

	if len(current) > 0 {
		result = append(result, dict.entries[string(current)])
	}

	var buf bytes.Buffer
	for _, code := range result {
		err := binary.Write(&buf, binary.LittleEndian, uint16(code))
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (lzw *LZWCompressor) Decompress(compressed []byte) ([]byte, error) {
	if len(compressed) == 0 {
		return nil, nil
	}

	codes := make([]int, 0)
	buf := bytes.NewReader(compressed)
	for {
		var code uint16
		err := binary.Read(buf, binary.LittleEndian, &code)
		if err != nil {
			break
		}
		codes = append(codes, int(code))
	}

	dict := make(map[int][]byte)
	for i := 0; i < 256; i++ {
		dict[i] = []byte{byte(i)}
	}
	nextCode := 256

	result := []byte{}

	if len(codes) == 0 {
		return result, nil
	}

	current := dict[codes[0]]
	result = append(result, current...)

	for i := 1; i < len(codes); i++ {
		var word []byte

		if entry, ok := dict[codes[i]]; ok {
			word = entry
		} else if codes[i] == nextCode {
			word = append(current, current[0])
		} else {
			return nil, errors.New("invalid compressed data")
		}

		result = append(result, word...)
		dict[nextCode] = append(current, word[0])
		nextCode++
		current = word
	}

	return result, nil
}
