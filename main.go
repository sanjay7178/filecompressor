// main.go
package main

import (
	"encoding/hex"
	"filecompressor/compress"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	var algorithms string
	var decompress bool
	var verbose bool

	flag.StringVar(&algorithms, "algo", "lzw", "Compression algorithms (comma-separated: lzw,huffman,rle,sf,bwt)")
	flag.BoolVar(&decompress, "d", false, "Decompress mode")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: compress [-d] [-v] [-algo=<algorithm>] <filename>")
		os.Exit(1)
	}

	algoList := strings.Split(algorithms, ",")
	chain := make([]compress.Compressor, 0)

	for _, algo := range algoList {
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
		default:
			fmt.Printf("Unknown algorithm: %s\n", algo)
			os.Exit(1)
		}
	}

	compressor := compress.NewCompressionChain(chain...)
	filename := flag.Arg(0)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("Original size: %d bytes\n", len(data))
	}

	var result []byte
	if decompress {
		if !strings.HasSuffix(filename, ".comp") {
			fmt.Printf("Error: Compressed file must have .comp extension\n")
			os.Exit(1)
		}

		result, err = compressor.Decompress(data)
		if err != nil {
			fmt.Printf("Error during decompression: %v\n", err)
			os.Exit(1)
		}

		outfile := strings.TrimSuffix(filename, ".comp")
		err = ioutil.WriteFile(outfile, result, 0644)
		if err != nil {
			fmt.Printf("Error writing decompressed file: %v\n", err)
			os.Exit(1)
		}

		if verbose {
			fmt.Printf("Decompressed size: %d bytes\n", len(result))
		}
		fmt.Printf("Successfully decompressed to: %s\n", outfile)

	} else {
		result, err = compressor.Compress(data)
		if err != nil {
			fmt.Printf("Error during compression: %v\n", err)
			os.Exit(1)
		}

		if verbose {
			fmt.Printf("Compressed size: %d bytes\n", len(result))
			fmt.Printf("Compression ratio: %.2f%%\n", float64(len(result))/float64(len(data))*100)
			fmt.Printf("First 32 bytes: %s\n", hex.EncodeToString(result[:min(32, len(result))]))
		}

		// Verify compression
		if len(result) == 0 {
			fmt.Printf("Error: Compression produced empty result\n")
			os.Exit(1)
		}

		outfile := filename + ".comp"
		err = ioutil.WriteFile(outfile, result, 0644)
		if err != nil {
			fmt.Printf("Error writing compressed file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully compressed to: %s using algorithms: %s\n", outfile, algorithms)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
