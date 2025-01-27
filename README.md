# go-ascv: Official ASCIIVideo (.ascv) Library for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/Vinetwigs/go-ascv.svg)](https://pkg.go.dev/github.com/Vinetwigs/go-ascv)

[![License](https://img.shields.io/badge/License-BSD_3--Clause-orange.svg)](https://opensource.org/licenses/BSD-3-Clause)

The `go-ascv` library is the **official implementation** of the ASCIIVideo (.ascv) file format in Go. ASCIIVideo is a compact, efficient format for storing and handling video sequences of ASCII art. This library provides a simple, powerful API to work with .ascv files, including creation, parsing, compression, and decompression.

- - -

## ✨ Features

*   **Parse and Create \`.ascv\` Files:** Read and generate ASCIIVideo files with ease.
*   **Built-In Compression:** Supports Run-Length Encoding (RLE) for space efficiency.
*   **Flexible Frame Handling:** Efficiently manage ASCII art video frames with VLQ (Variable-Length Quantity) encoding.
*   **Versioning Support:** Fully adheres to the \`.ascv\` v1.0 specification, ensuring compatibility with future updates.
*   **High Performance:** Optimized for speed and low memory usage.

- - -

## 📖 Installation

To install the library, run:

```bash
go get github.com/Vinetwigs/go-ascv@v1.0
```

or get the latest version with:

```bash
go get github.com/Vinetwigs/go-ascv@latest
```

## 🚀 Quick Start

Here’s how you can get started with `go-ascv`:

```go

package main

import (
	"fmt"
	"github.com/Vinetwigs/go-ascv"
)

func main() {
	// Create a new ASCIIVideo file
	video := ascv.NewASCIIVideo(8, 4, 30, true)

	// Add frames
	video.AddFrame("AAAABBBBCCCCDDDD") // Frame 1
	video.AddFrame("XXXXXXXXYYYYZZZZ") // Frame 2

	// Save to file
	err := video.Save("example.ascv")
	if err != nil {
		fmt.Println("Error saving file:", err)
		return
	}

	fmt.Println("ASCIIVideo file saved successfully!")
}
```

- - -

## 📂 File Format Overview

The .ascv file format is designed for space-efficient storage of ASCII art videos. It uses techniques like Run-Length Encoding (RLE) and Variable-Length Quantity (VLQ) for compact frame data representation.

### Header Structure (32 bytes)

| Field | Type | Size | Description |
| --- | --- | --- | --- |
| MAGIC | String | 4 bytes | Magic number identifying the file |
| VERSION | Integer | 1 byte | File format version |
| WIDTH | Integer | 2 bytes | Video width in characters |
| HEIGHT | Integer | 2 bytes | Video height in characters |
| FPS | Integer | 1 byte | Frames per second |
| FRAMES | Integer | 4 bytes | Total number of frames |
| COMP | Integer | 1 byte | Compression flag (0 = off, 1 = on) |
| CHARSET | Integer | 1 byte | Character set identifier |
| RESERVED | Reserved | 16 bytes | Reserved for future extensions |

For a full breakdown of the format, refer to the [specifications](.\docs\ASCV_File_Format_Specifications.pdf).

## 🛠️ Supported Compression

*   **Run-Length Encoding (RLE):** Compresses repeated characters in frames.

Example:

```
    Original: AAAABBBCCDDDDDD
    Compressed: 04A03B02C06D
    
```

*   **VLQ Encoding:** Optimizes storage of variable-length data like frame sizes.

## 🖥️ Compatibility

*   **Go Version:** Requires Go 1.18 or newer.
*   **Platforms:** Cross-platform support (Linux, macOS, Windows).

## 📜 License

This project is licensed under the [BSD 3-Clause License](https://opensource.org/license/bsd-3-clause).

## 🌟 Contributing

We welcome contributions to improve the library! To get started:

*   Fork the repository.
*   Create a new branch for your feature or bugfix.
*   Submit a pull request with a clear description of your changes.

## 📧 Contact

For questions or support, feel free to open an issue!

Made with ❤️ for ASCII art enthusiasts!