MODULE_DIRS = . ./tools

gowork:
	go work init . ./tools

tidy:
	$(foreach dir,$(MODULE_DIRS), \
		(cd $(dir) && go mod tidy) &&) true

install: tidy
	cd tools && go install \
		mvdan.cc/gofumpt

fmt: install
	gofumpt -l -w -extra .

lint: install
	golangci-lint run ./...

test:
	go test ./...

check: fmt lint test
