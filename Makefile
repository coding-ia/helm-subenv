GOFMT_FILES?=$$(find . -name '*.go')
export GO111MODULE=on

default: build

build:
	go install

fmt:
	gofmt -w $(GOFMT_FILES)