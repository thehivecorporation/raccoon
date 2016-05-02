## simple makefile to log workflow
.PHONY: all test clean build install integration

all: test install

build:
	go build ./...

install:
	go get ./...

test:
	go test ./...

cover:
	go test -cover ./...

integration-test:
	tests/integration-test

bench:
	go test -run=NONE -bench=. ./...

clean:
	go clean -i ./...
