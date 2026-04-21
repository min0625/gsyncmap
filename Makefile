NEW_FROM_REV ?= HEAD

fix:
	go mod tidy
	golangci-lint run -v --new-from-rev=$(NEW_FROM_REV) --fix ./...

lint:
	golangci-lint run -v --new-from-rev=$(NEW_FROM_REV) ./...

test:
	go test -v -race -failfast ./...

check-tidy:
	go mod tidy -diff

check: check-tidy lint test
