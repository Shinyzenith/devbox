BINARY:=devbox
PREFIX:=/usr
ARCH := $(shell uname -p)
FMT_REQUIRED:=$(shell gofmt -l $(shell find . -type f -iname *.go))
LDFLAGS=-s -w -X github.com/shinyzenith/devbox/version.Version=$(shell git rev-parse HEAD) -linkmode external -extldflags "-static"

all: zig_static

build: tidy
	go build ./cmd/$(BINARY)
	strip $(BINARY)

zig_static: tidy
	CGO_ENALBED=1 CC="zig cc -target ${ARCH}-linux-musl" go build -v \
				-ldflags '${LDFLAGS}' \
				./cmd/$(BINARY)/
	strip $(BINARY)

musl_static: tidy
	CGO_ENALBED=1 CC="musl-gcc" go build -v \
				-ldflags '${LDFLAGS}' \
				./cmd/$(BINARY)/
	strip $(BINARY)

dependencies:
	go get -u ./...
	go mod tidy

install:
	mv $(BINARY) $(PREFIX)/bin/$(BINARY)

check:
	@echo $(FMT_REQUIRED)
	@test -z $(FMT_REQUIRED)
	go vet ./...

test:
	go test ./...

tidy:
	go mod tidy

clean:
	go clean
	$(RM) -f $(BINARY)

.PHONY:all build test install clean dependencies musl_static zig_static tidy check
