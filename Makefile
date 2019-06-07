BINARY = abscraper

GOOS =  $(shell go env GOOS)
GOARCH =  $(shell go env GOARCH)
GOBUILD =  go build -o bin/$(BINARY)-$(GOOS)-$(GOARCH)

all: clean fmt build

build:
	$(GOBUILD) main.go
	ln -sf $(BINARY)-$(GOOS)-$(GOARCH) bin/$(BINARY)

clean:
	rm -rf bin

fmt:
	go fmt ./...
