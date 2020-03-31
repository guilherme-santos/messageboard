help: ## Display this help
	@ echo "Please use \`make <target>' where <target> is one of:"
	@ echo
	@ grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36m%-16s\033[0m - %s\n", $$1, $$2}'
	@ echo

export version?=latest

build: param-version ## Build the docker image
	docker-compose build

push: param-version ## Push the docker image to our private docker registry
	docker-compose push

up: ## Run the service on docker-compose locally
	@docker-compose up -d

down: ## Stop the service on docker-compose locally
	@docker-compose down

param-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Param \"$*\" is missing, use: make $(MAKECMDGOALS) $*=<value>"; \
		exit 1; \
	fi

# To be used internally or during the development

export GIT_TAG ?= $(shell git describe --tags --exact-match 2> /dev/null || git symbolic-ref -q --short HEAD)
export GIT_COMMIT ?= $(shell git rev-parse --short HEAD)

MAKE := $(MAKE) -f $(MAKEFILE_LIST)
LDFLAGS = -s -w -X "main.tag=$(GIT_TAG)" -X "main.gitCommit=$(GIT_COMMIT)" -X "main.buildTime=$(shell date -u +%FT%T%z)"

go-generate:
	@go generate ./...

testflag?=-race $(flag)
test: ## Run unit tests, set testcase=<testcase> and flag=-v if you need them
	go test -failfast ./... $(testflag) $(if $(testcase),-run "$(testcase)")

test-coverage: override flag+=-coverprofile coverage.out ## Show test coverage
test-coverage: coverage.out
	go tool cover -html=coverage.out

coverage.out:
	testflag="$(flag)" $(MAKE) test

go-build: ## Build the binaries
	go build -v -ldflags "$(LDFLAGS)" ./cmd/messageboard

go-install: ## Build the binaries statically and install it
	CGO_ENABLED=0 go install -v -ldflags "$(LDFLAGS)" -a -installsuffix cgo ./cmd/messageboard

run: go-build ## Build and run the app locally
	@./messageboard
