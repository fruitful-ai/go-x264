# go-x264

Go bindings for x264 — statically linked, zero-latency H.264 encoding.

No system `libx264` required. `make` downloads a pre-built release asset (linux/amd64) or builds from source.

## Prerequisites

- Go 1.26+

### Source build (non-linux/amd64 or first run)

```bash
sudo apt install make gcc curl bzip2   # Debian / Ubuntu
xcode-select --install                 # macOS
```

## Installation

```bash
go get fruitful.ch/x264
make
```

## Quick Start

```go
package main

import (
    "log"
    "fruitful.ch/x264"
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
make          # build x264 + Go package
make test     # run tests
make clean    # remove build artifacts
```

## License

MIT
