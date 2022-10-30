format:
	@go install mvdan.cc/gofumpt@latest
	gofumpt -l -w -extra .

lint:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1
	golangci-lint run ./...

test:
	go test ./...

check: format lint test
	go mod tidy
