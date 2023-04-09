BINARY:=devbox
PREFIX:=/usr

LDFLAGS=-s -w -X github.com/shinyzenith/devbox/version.Version=$(shell git rev-parse HEAD) -linkmode external -extldflags "-static"
ARCH := $(shell uname -p)

all: zig_static

build: tidy
	go build ./cmd/$(BINARY)
	$(MAKE) -s strip

zig_static: tidy
	CGO_ENALBED=1 CC="zig cc -target ${ARCH}-linux-musl" go build -v \
				-ldflags '${LDFLAGS}' \
				./cmd/$(BINARY)/
	$(MAKE) -s strip

musl_static: tidy
	CGO_ENALBED=1 CC="musl-gcc" go build -v \
				-ldflags '${LDFLAGS}' \
				./cmd/$(BINARY)/
	@$(MAKE) -s strip

dependencies:
	go get -u ./...
	go mod tidy

install:
	mv $(BINARY) $(PREFIX)/bin/$(BINARY)

test:
	go test ./...

tidy:
	go mod tidy

strip:
	strip $(BINARY)

clean:
	go clean
	$(RM) -f $(BINARY)

.PHONY:all build test install clean dependencies musl_static zig_static tidy strip
