NEW_FROM_REV ?= HEAD

.PHONY: fix
fix:
	go mod tidy
	golangci-lint run --new-from-rev=$(NEW_FROM_REV) --fix ./...

.PHONY: lint
lint:
	golangci-lint config verify
	golangci-lint run --new-from-rev=$(NEW_FROM_REV) ./...

.PHONY: test
test:
	go test -race -failfast ./...

.PHONY: cover
cover:
	go test -race -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

.PHONY: cover-html
cover-html: cover
	go tool cover -html=coverage.out -o coverage.html

.PHONY: check-tidy
check-tidy:
	go mod tidy -diff

.PHONY: check
check: check-tidy lint test
