FORMATTER=mvdan.cc/gofumpt@latest
LINTER=github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: default
default: vet fix fmt lint

.PHONY: vet
vet:
	@echo "go vet"
	@go vet ./...

.PHONY: fix
fix:
	@echo "go fix"
	@go fix ./...

.PHONY: fmt
fmt:
	@echo "go fmt"
	@go run $(FORMATTER) -l -w .

.PHONY: lint
lint:
	@echo "go lint"
	@go run $(LINTER) run

.PHONY: run
run:
	go run main.go

.PHONY: build
build:
	go build -o delete main.go

.PHONY: clean
clean:
	rm -f delete

.PHONY: install
install:
	mv delete ~/bin
