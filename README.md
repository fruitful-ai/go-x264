# go-x264

Go bindings for x264 — dynamically linked, zero-latency H.264 encoding.

## Prerequisites

- Go 1.26+
- libx264 (shared library + headers)

```bash
# Debian / Ubuntu
sudo apt install libx264-dev gcc

# macOS
brew install x264

# Fedora
sudo dnf install libx264-devel
```

## Installation

```bash
go get github.com/fruitful-ai/go-x264
```

## Quick Start

```go
package main

import (
    "log"
    "github.com/fruitful-ai/go-x264"
)

func main() {
    enc, err := x264.New(640, 480, 30)
    if err != nil {
        log.Fatal(err)
    }
    defer enc.Close()

    yuv := make([]byte, 640*480*3/2)
    data, err := enc.EncodeYUV(yuv)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("encoded %d bytes", len(data))
}
```

## Commands

```bash
go build ./...     # build
go test -v ./...   # test
```

## License

MIT
