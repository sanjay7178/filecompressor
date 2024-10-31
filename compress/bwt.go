// compress/bwt.go
package compress

import (
	"bytes"
	"errors"
	"sort"
)

type BWTCompressor struct {
	blockSize int
}

func NewBWTCompressor(blockSize int) *BWTCompressor {
	return &BWTCompressor{blockSize: blockSize}
}

func (bwt *BWTCompressor) transform(data []byte) ([]byte, int) {
	n := len(data)
	rotations := make([][]byte, n)

	// Create all rotations
	for i := 0; i < n; i++ {
		rotation := make([]byte, n)
		for j := 0; j < n; j++ {
			rotation[j] = data[(i+j)%n]
		}
		rotations[i] = rotation
	}

	// Sort rotations
	sort.Slice(rotations, func(i, j int) bool {
		return bytes.Compare(rotations[i], rotations[j]) < 0
	})

	// Find original string index and get last column
	result := make([]byte, n)
	originalIndex := 0
	for i := 0; i < n; i++ {
		result[i] = rotations[i][n-1]
		if bytes.Equal(rotations[i], data) {
			originalIndex = i
		}
	}

	return result, originalIndex
}

func (bwt *BWTCompressor) inverseTransform(data []byte, originalIndex int) []byte {
	n := len(data)
	table := make([][]byte, n)

	// Initialize with single bytes
	for i := 0; i < n; i++ {
		table[i] = []byte{data[i]}
	}

	// Sort table repeatedly
	for i := 0; i < n-1; i++ {
		sort.Slice(table, func(i, j int) bool {
			return bytes.Compare(table[i], table[j]) < 0
		})

		// Prepend transformed data
		for j := 0; j < n; j++ {
			table[j] = append([]byte{data[j]}, table[j]...)
		}
	}

	return table[originalIndex]
}

func (bwt *BWTCompressor) Compress(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}

	var result []byte
	blockCount := (len(data) + bwt.blockSize - 1) / bwt.blockSize

	// Write number of blocks
	result = append(result, byte(blockCount))

	// Process each block
	for i := 0; i < blockCount; i++ {
		start := i * bwt.blockSize
		end := start + bwt.blockSize
		if end > len(data) {
			end = len(data)
		}

		block := data[start:end]
		transformed, index := bwt.transform(block)

		// Write block metadata
		result = append(result, byte(len(block)))
		result = append(result, byte(index))

		// Write transformed block
		result = append(result, transformed...)
	}

	return result, nil
}

func (bwt *BWTCompressor) Decompress(compressed []byte) ([]byte, error) {
	if len(compressed) == 0 {
		return nil, nil
	}

	var result []byte
	pos := 0

	// Read number of blocks
	blockCount := int(compressed[pos])
	pos++

	// Process each block
	for i := 0; i < blockCount && pos < len(compressed); i++ {
		// Read block metadata
		blockSize := int(compressed[pos])
		pos++
		originalIndex := int(compressed[pos])
		pos++

		if pos+blockSize > len(compressed) {
			return nil, errors.New("invalid compressed data")
		}

		// Read and inverse transform block
		block := compressed[pos : pos+blockSize]
		decompressed := bwt.inverseTransform(block, originalIndex)
		result = append(result, decompressed...)

		pos += blockSize
	}

	return result, nil
}
