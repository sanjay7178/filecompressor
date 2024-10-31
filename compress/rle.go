// compress/rle.go
package compress

type RLECompressor struct{}

func NewRLECompressor() *RLECompressor {
	return &RLECompressor{}
}

func (rc *RLECompressor) Compress(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}

	var result []byte
	count := 1
	current := data[0]

	for i := 1; i < len(data); i++ {
		if data[i] == current && count < 255 {
			count++
		} else {
			result = append(result, byte(count), current)
			current = data[i]
			count = 1
		}
	}
	result = append(result, byte(count), current)

	return result, nil
}

func (rc *RLECompressor) Decompress(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}

	var result []byte
	for i := 0; i < len(data); i += 2 {
		count := int(data[i])
		char := data[i+1]
		for j := 0; j < count; j++ {
			result = append(result, char)
		}
	}

	return result, nil
}
