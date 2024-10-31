// compress/chain.go
package compress

type CompressionChain struct {
	compressors []Compressor
}

func NewCompressionChain(compressors ...Compressor) *CompressionChain {
	return &CompressionChain{compressors: compressors}
}

func (cc *CompressionChain) Compress(data []byte) ([]byte, error) {
	result := data
	var err error

	for _, c := range cc.compressors {
		result, err = c.Compress(result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (cc *CompressionChain) Decompress(data []byte) ([]byte, error) {
	result := data
	var err error

	// Decompress in reverse order
	for i := len(cc.compressors) - 1; i >= 0; i-- {
		result, err = cc.compressors[i].Decompress(result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
