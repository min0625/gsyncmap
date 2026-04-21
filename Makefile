fix:
	golangci-lint run -v --fix ./...

lint:
	golangci-lint run -v ./...

test:
	go test -v -race -failfast ./...

check: lint test
