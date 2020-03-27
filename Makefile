
phony: all generate lint

GOBIN = $(GOPATH)/bin
GOAGEN = $(GOBIN)/goagen
GOIMPORTS = $(GOBIN)/goimports
GOLANGCI_LINT := $(GOPATH)/bin/golangci-lint

generate: src/slr-api/design/*.go | $(GOAGEN) $(GOIMPORTS)
	@echo Generating database models and controllers
	@goagen app -d github.com/wimspaargaren/slr-automation/src/slr-api/design -o src/slr-api
	@goagen main -d github.com/wimspaargaren/slr-automation/src/slr-api/design --regen -o src/slr-api
	@goagen -d github.com/wimspaargaren/slr-automation/src/slr-api/design gen --pkg-path=lab.weave.nl/forks/gorma -o src/slr-api
	@goimports -w src/slr-api/models
	@perl -i -pe's#$(PWD)#\$$(GOPATH)/src/github.com/wimspaargaren/slr-automation/src/slr-api#g' src/slr-api/models/*

lint: | $(GOLANGCI_LINT)
	@echo Linting Go files
	@ $(GOLANGCI_LINT) run --deadline=30m --exclude-use-default=false -v  --disable dupl