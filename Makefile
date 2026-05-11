.PHONY: all build test clean

all: build

build:
	./build_x264.sh
	go build ./...

test:
	./build_x264.sh
	go test -v ./...

clean:
	rm -rf internal/x264
	go clean -cache
