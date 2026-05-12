.PHONY: all build test clean

all: build

build:
	go build ./...

test:
	go test -v ./...

clean:
	rm -rf internal/x264
	go clean -cache
