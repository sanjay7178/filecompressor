// main.go
package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "os"
    "strings"
    "filecompressor/compress"
)

func main() {
    var algorithms string
    var decompress bool

    flag.StringVar(&algorithms, "algo", "lzw", "Compression algorithms (comma-separated: lzw,huffman,rle,sf,bwt)")
    flag.BoolVar(&decompress, "d", false, "Decompress mode")
    flag.Parse()

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

    if flag.NArg() < 1 {
        fmt.Println("Usage: compress [-d] [-algo=<algorithm>] <filename>")
        os.Exit(1)
    }

    filename := flag.Arg(0)
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        fmt.Printf("Error reading file: %v\n", err)
        os.Exit(1)
    }

    var result []byte
    if decompress {
        // Validate file extension
        if !strings.HasSuffix(filename, ".comp") {
            fmt.Printf("Error: Compressed file must have .comp extension\n")
            os.Exit(1)
        }

        result, err = compressor.Decompress(data)
        if err != nil {
            fmt.Printf("Error during decompression: %v\n", err)
            os.Exit(1)
        }

        // Remove .comp extension
        outfile := strings.TrimSuffix(filename, ".comp")
        err = ioutil.WriteFile(outfile, result, 0644)
        if err != nil {
            fmt.Printf("Error writing decompressed file: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Successfully decompressed to: %s\n", outfile)

    } else {
        result, err = compressor.Compress(data)
        if err != nil {
            fmt.Printf("Error during compression: %v\n", err)
            os.Exit(1)
        }

        // Use consistent .comp extension
        outfile := filename + ".comp"
        err = ioutil.WriteFile(outfile, result, 0644)
        if err != nil {
            fmt.Printf("Error writing compressed file: %v\n", err)
            os.Exit(1)
        }
        fmt.Printf("Successfully compressed to: %s using algorithms: %s\n", outfile, algorithms)
    }

    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
}



