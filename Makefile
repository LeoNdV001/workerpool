.DEFAULT_GOAL := help

# Colors used
NOCOLOR=\033[0m
INPUTCOLOR=\033[0;36m
OUTPUTCOLOR=\033[0;32m
SOMECOLOR=\033[01;31m
ANOTHERCOLOR=\033[01;34m

PKG := "github.com/LeoNdV001/workerpool"
PKG_LIST := $(shell go list ${PKG}/... | grep -v .mocks)
START=@echo "$(ANOTHERCOLOR) > Starting $@ ...$(NOCOLOR)"
END=@echo "$(OUTPUTCOLOR) > Finished $@ ...$(NOCOLOR)"

makefile: help

.PHONY: help
## help: Lists all the make-commands
help: Makefile
	@echo "$(ANOTHERCOLOR) > Choose a make command from the following:"
	@(echo " -----------: -----------\n COMMAND: DESCRIPTION\n -----------: -----------$(NOCOLOR)\n" &&\
 	(grep -h "## " $(MAKEFILE_LIST) | grep -v '%' |  sed -e 's/##//' | sed -e '1d')) | sed -e 's/: /:/' | column -t -s \:;

.PHONY: upgrade
## upgrade: Upgrade all go dependencies
upgrade: | ;
	$(START)
	@go clean -modcache
	@go get -u -d ./...
	@go mod tidy;
	$(END);

.PHONY: tidy
## tidy: Runs go mod tidy command
tidy: | ;
	$(START)
	@go clean -modcache
	@go mod tidy;
	$(END);

.PHONY: test
## test: Runs the unit tests
test: | ;
	$(START)
	@go clean -testcache
	@go test ./... -v -count=1 -failfast;
	$(END);

.PHONY: test-cov
## test-cov: Runs test coverage locally
test-cov: | ;
	$(START)
	@go clean -testcache
	@go test ${PKG_LIST} -coverprofile coverage.out
	@go tool cover -func coverage.out | grep total
	@go tool cover -html=coverage.out -o coverage.html;
	$(END);

.PHONY: lint
## lint: Runs linting utility
lint: | ;
	$(START)
	@golangci-lint run -c ./.golangci.yml --timeout=3m;
	$(END);

.PHONY: lint-fix
## lint-fix: Runs linting utility and automatically fixes issues (if possible)
lint-fix: | ;
	$(START)
	@golangci-lint run -c ./.golangci.yml --timeout=3m --fix;
	$(END);

.PHONY: fmt
## fmt: Runs gofmt
fmt: | ;
	$(START)
	@gofmt -s -w .;
	$(END);

.PHONY: fumpt
## fumpt: Runs gofumpt
fumpt: | ;
	$(START)
	@gofumpt -w .;
	$(END);

.PHONY: vuln
## vuln: Runs the Go Vulnerability Check
vuln: | ;
	$(START)
	@govulncheck ./...;
	$(END);