BINARY:=devbox

all: build

build:
	@go build -o devbox cmd/devbox/main.go

test:
	@go test ./...

clean:
	@go clean
	@$(RM) -f $(BINARY)
