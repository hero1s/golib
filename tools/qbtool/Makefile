# Go parameters
GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOBIN=$(GOPATH)/bin
GOINSTALL=$(GOCMD) install

all: test build
build:
	rm -rf ./qbtool
	$(GOBUILD) -o qbtool main.go

install: build
	mv qbtool $(GOBIN)
test:
	$(GOTEST) -v ./...

clean:
	rm -rf ./qbtool
