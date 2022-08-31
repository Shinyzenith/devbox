BINARY:=devbox

all: build

build:
	@go build ./cmd/devbox/

test:
	@go test ./...

clean:
	@go clean
	@$(RM) -f $(BINARY)
