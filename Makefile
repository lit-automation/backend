
phony: all generate lint run docker-build docker-run

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

run:
	@go build ./src/slr-api
	@./slr-api

docker-build:
	@docker build -f ./Dockerfile-builder -t slr-builder .
	@docker build --build-arg base_img=slr-builder -f ./Dockerfile -t slr-api .

docker-run:
	@docker run -e JWT_KEY=${JWT_KEY} -e JWT_KEY_PUB=${JWT_KEY_PUB} -e PG_HOST=${PG_HOST} -e PG_PASSWORD=${PG_PASSWORD} -e PG_USERNAME=${PG_USERNAME} -p 9001:9001 slr-api