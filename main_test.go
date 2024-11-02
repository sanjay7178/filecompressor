// main_test.go
package main

import (
    "bytes"
    "filecompressor/compress"
    "io/ioutil"
    "os"
    "testing"
    "path/filepath"
    "strings"
)

// Helper function to handle compression
func compressData(data []byte, algorithms string) ([]byte, error) {
    chain := make([]compress.Compressor, 0)
    
    for _, algo := range strings.Split(algorithms, ",") {
        switch algo {
        case "lzw":
            chain = append(chain, compress.NewLZWCompressor())
        case "huffman":
            chain = append(chain, compress.NewHuffmanCompressor())
        case "rle":
            chain = append(chain, compress.NewRLECompressor())
        case "sf":
            chain = append(chain, compress.NewShannonFanoCompressor())
        case "bwt":
            chain = append(chain, compress.NewBWTCompressor(1024))
        }
    }
    
    compressor := compress.NewCompressionChain(chain...)
    return compressor.Compress(data)
}

func TestAllAlgorithms(t *testing.T) {
    algorithms := []string{
        "lzw",
        "huffman",
        "rle",
        "sf",
        "bwt",
        "lzw,huffman",
        "huffman,rle",
        "sf,bwt",
        "lzw,huffman,rle",
        "sf,bwt,huffman",
        "lzw,huffman,rle,sf,bwt",
    }

    testData := []string{
        "This is a test string",
        "Hello world!",
        "Repeated repeated repeated repeated data",
        "Binary\x00\x01\x02\x03data\xff\xfe\xfd\xfc",
    }

    for _, algo := range algorithms {
        for _, data := range testData {
            t.Run(algo+"_"+data[:10], func(t *testing.T) {
                input := []byte(data)
                
                // Compress
                compressed, err := compressData(input, algo)
                if err != nil {
                    t.Fatalf("Compression failed: %v", err)
                }
                
                if len(compressed) == 0 {
                    t.Error("Compressed data is empty")
                }

                // Write to temporary file
                tmpFile := filepath.Join(os.TempDir(), "test.comp")
                if err := ioutil.WriteFile(tmpFile, compressed, 0644); err != nil {
                    t.Fatalf("Failed to write compressed file: %v", err)
                }
                defer os.Remove(tmpFile)

                // Read and decompress
                compressedData, err := ioutil.ReadFile(tmpFile)
                if err != nil {
                    t.Fatalf("Failed to read compressed file: %v", err)
                }

                chain := make([]compress.Compressor, 0)
                for _, a := range strings.Split(algo, ",") {
                    switch a {
                    case "lzw":
                        chain = append(chain, compress.NewLZWCompressor())
                    case "huffman":
                        chain = append(chain, compress.NewHuffmanCompressor())
                    case "rle":
                        chain = append(chain, compress.NewRLECompressor())
                    case "sf":
                        chain = append(chain, compress.NewShannonFanoCompressor())
                    case "bwt":
                        chain = append(chain, compress.NewBWTCompressor(1024))
                    }
                }

                compressor := compress.NewCompressionChain(chain...)
                decompressed, err := compressor.Decompress(compressedData)
                if err != nil {
                    t.Fatalf("Decompression failed: %v", err)
                }

                if !bytes.Equal(input, decompressed) {
                    t.Errorf("Data mismatch after compression/decompression\nInput: %q\nOutput: %q", input, decompressed)
                }
            })
        }
    }
}

func TestLargeFile(t *testing.T) {
	// Create 1MB test file
	data := make([]byte, 1024*1024)
	for i := range data {
		data[i] = byte(i % 256)
	}

	tmpfile, err := ioutil.TempFile("", "largefile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test with different block sizes for BWT
	blockSizes := []string{"bwt", "lzw,bwt", "bwt,huffman"}
	
	for _, algo := range blockSizes {
		t.Run(algo, func(t *testing.T) {
			compressFile := tmpfile.Name() + ".comp"
			os.Args = []string{"cmd", "-algo=" + algo, tmpfile.Name()}
			main()

			os.Args = []string{"cmd", "-d", compressFile}
			main()

			decompressed, err := ioutil.ReadFile(tmpfile.Name())
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(data, decompressed) {
				t.Fatal("Large file decompression failed")
			}

			os.Remove(compressFile)
		})
	}
}