FORMATTER=mvdan.cc/gofumpt@latest
LINTER=github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: default
default: vet fix fmt lint

.PHONY: vet
vet:
	go vet

.PHONY: fix
fix:
	go fix

.PHONY: fmt
fmt:
	@echo "go fmt"
	@go run $(FORMATTER) -l -w .

.PHONY: lint
lint:
	@echo "go lint"
	@go run $(LINTER) run

run:
	@go run main.go

build:
	@go build -o delete main.go

clean:
	@rm -f delete
