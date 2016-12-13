BINARY		= travis-artifacts
BIN_DIR		= bin

SHELL		:= /bin/bash
GOOS		:= $(shell go env GOOS)
GOARCH		:= $(shell go env GOARCH)

default: clean install test build

install:
	go get -t -v ./...

test:
	if [ "${TRAVIS}" == "true" ]; then		\
		go get golang.org/x/tools/cmd/cover;\
		go get github.com/mattn/goveralls;	\
		goveralls;							\
	else									\
		go test -v ./...;					\
	fi

build:
	test -d $(BIN_DIR) || mkdir -p $(BIN_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BIN_DIR)/$(BINARY) cmd/travis-artifacts/main.go

clean:
	if [ -d $(BIN_DIR) ]; then rm -r $(BIN_DIR); fi

.PHONY: default test build clean
