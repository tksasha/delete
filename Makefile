MAIN=cmd/delete/main.go
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
	go run $(MAIN)

.PHONY: build
build:
	go build -o delete $(MAIN)

.PHONY: clean
clean:
	rm -f delete

.PHONY: install
install:
	mv delete ~/bin

.PHONY: test
test: build
	@echo "preparing test fs..."
	@go run test/prepare.go
	@./delete test/__fs
	@make clean
