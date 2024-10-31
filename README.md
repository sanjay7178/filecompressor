# Go File Compression Tool

A versatile file compression tool implemented in Go that supports multiple compression algorithms and chaining them together.

## Features

- Multiple compression algorithms:
    - Huffman Coding
    - Shannon-Fano Coding
    - Burrows-Wheeler Transform (BWT)
    - Run-Length Encoding (RLE)
    - LZW Compression

- Algorithm chaining capability
- Command-line interface
- Supports both compression and decompression

## Installation

```bash
git clone https://github.com/yourusername/filecompressor
cd filecompressor
go build
```

## Usage

Basic usage:
```bash
./filecompressor [-d] [-algo=<algorithms>] <filename>
```

Options:
- `-d`: Decompress mode
- `-algo`: Comma-separated list of compression algorithms (default: "lzw")
    - Available algorithms: lzw, huffman, rle, sf, bwt

Examples:
```bash
# Compress using default LZW
./filecompressor myfile.txt

# Compress using multiple algorithms
./filecompressor -algo=huffman,bwt,rle myfile.txt

# Decompress a file
./filecompressor -d myfile.txt.comp
```

## Implementation Details

- Uses a chain of compression algorithms that can be applied sequentially
- Each algorithm implements the Compressor interface
- Compressed files use the `.comp` extension
- Supports various compression techniques including:
    - Huffman coding with tree serialization
    - Shannon-Fano coding with frequency-based division
    - Burrows-Wheeler Transform with configurable block size

## Contributing

Feel free to submit issues, fork the repository, and create pull requests for any improvements.

## License

MIT

