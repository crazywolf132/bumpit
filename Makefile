BINARY_NAME=bumpit
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date +%FT%T%z)
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"
MAIN_GO=./cmd/bumpit/main.go

.PHONY: build clean test lint install uninstall

default: build

build:
	go build ${LDFLAGS} -o ${BINARY_NAME} -v $(MAIN_GO)

install:
	go install ${LDFLAGS}

clean:
	go clean
	rm -f ${BINARY_NAME}

test:
	go test ./... -v

lint:
	@if command -v golangci-lint >/dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint is not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi

# Create a new release
release:
	@if [ -z "$(VERSION)" ]; then \
		echo "Usage: make release VERSION=v1.0.0"; \
		exit 1; \
	fi
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)
