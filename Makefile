.PHONY: build test lint fmt vet clean

build:
	go build ./...

test:
	go test -race ./...

lint:
	golangci-lint run ./...

fmt:
	gofmt -s -w .
	goimports -w -local github.com/latticehq/latticesdk .

vet:
	go vet ./...

clean:
	rm -rf bin/ dist/

check: fmt vet lint test
