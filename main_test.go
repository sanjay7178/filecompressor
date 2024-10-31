package main

import (
	"bytes"
	// "filecompressor/compress"
	"io/ioutil" 
	"os"
	"testing"
	"strings"
)

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
		"This is a simple test string",
		"AAAAABBBCCCCCCDDEEEEEE", // Good for RLE
		"mississippi", // Good for Huffman/Shannon-Fano
		strings.Repeat("test", 100), // Longer input
		"", // Empty string
	}

	for _, data := range testData {
		for _, algo := range algorithms {
			t.Run(algo+"_"+data[:min(10,len(data))], func(t *testing.T) {
				// Create temp input file
				tmpfile, err := ioutil.TempFile("", "test")
				if err != nil {
					t.Fatal(err)
				}
				defer os.Remove(tmpfile.Name())

				if _, err := tmpfile.Write([]byte(data)); err != nil {
					t.Fatal(err) 
				}
				if err := tmpfile.Close(); err != nil {
					t.Fatal(err)
				}

				// Test compression
				compressFile := tmpfile.Name() + ".comp"
				os.Args = []string{"cmd", "-algo=" + algo, tmpfile.Name()}
				main()

				compressed, err := ioutil.ReadFile(compressFile) 
				if err != nil {
					t.Fatalf("Failed to read compressed file: %v", err)
				}
				if len(compressed) == 0 && len(data) > 0 {
					t.Fatal("Compressed file is empty")
				}

				// Test decompression
				os.Args = []string{"cmd", "-d", compressFile}
				main()

				decompressed, err := ioutil.ReadFile(tmpfile.Name())
				if err != nil {
					t.Fatalf("Failed to read decompressed file: %v", err)
				}

				if !bytes.Equal([]byte(data), decompressed) {
					t.Fatal("Decompressed data doesn't match original")
				}

				// Cleanup
				os.Remove(compressFile)
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