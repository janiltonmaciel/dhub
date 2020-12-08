.SILENT: help
SHELL = /bin/bash
.DEFAULT_GOAL := help

PROJECT := dhub
GITHUB_TOKEN := $(shell git config --get github.version-gen-token || echo $$GITHUB_TOKEN)

NOW = $(shell date +"%d-%m-%YT%H:%M:%S")
GIT_COMMIT = $(shell git rev-parse --short HEAD)
GIT_TAG = $(shell git describe --tags)
LDFLAGS := -X 'main.version=${GIT_TAG}'
LDFLAGS += -X 'main.commit=${GIT_COMMIT}'
LDFLAGS += -X 'main.date=${NOW}'
LDFLAGS += -X 'main.token=${GITHUB_TOKEN}'


#############  SETUP  #############

## Setup of the project
setup:
	@go mod download
	@brew install goreleaser/tap/goreleaser


#############  RUN  #############

## Run project
run:
	@go run main.go


#############  TESTS  #############

## Runs the project unit tests
test:
	@go test -timeout 10s  -v -covermode atomic -cover -coverprofile coverage.out
	@go vet . 2>&1 | grep -v '^vendor\/' | grep -v '^exit\ status\ 1' || true

## Run all the tests and opens the coverage report
test-cover: test
	go tool cover -html=coverage.out
	@rm coverage.out 2>/dev/null || true

## Run all the tests and code checks
test-ci: lint test

lint:
	@echo "*** Start lint ***"
	@script -q /dev/null golangci-lint run --enable-all -D gochecknoglobals,lll,wsl,maligned,nlreturn --print-issued-lines=false --out-format=colored-line-number | awk '{print $0; count++} END {print "\nCount: " count}'
	@echo "*** Finish ***"


#############  DEPLOY  #############

## Build project
build:
	@echo "Building $(PROJECT)"
	export GITHUB_TOKEN=$(GITHUB_TOKEN); \
	go build -ldflags "$(LDFLAGS) -s -w" -o $(PROJECT) main.go

## Build Docker Image
docker:
	echo "Building Docker Image of $(PROJECT)"
	docker build --target release -t $(PROJECT) .

git-tag:
	@printf "\n"; \
	read -p "Tag ($(TAG)): "; \
	if [ ! "$$REPLY" ]; then \
		printf "\n${COLOR_RED}"; \
		echo "Invalid tag."; \
		exit 1; \
	fi; \
	TAG=$$REPLY; \
	if git rev-parse $$TAG >/dev/null 2>&1; then \
		echo "TAG: ${TAG}"; \
	else \
		sed -i.bak "s/download\/[^/]*/download\/$$TAG/g" README.md && \
		sed -i.bak "s/statiks_[^_]*/statiks_$$TAG/g" README.md  && \
		rm README.md.bak 2>/dev/null; \
		git commit README.md -m "Update README.md with release $$TAG"; \
		git tag -s $$TAG -m "$$TAG"; \
	fi;

## Release of the project
release: git-tag
	@if [ ! "$(GITHUB_TOKEN)" ]; then \
		echo "github token should be configurated."; \
		exit 1; \
	fi; \
	export GITHUB_TOKEN=$(GITHUB_TOKEN); \
	goreleaser release --rm-dist; \
	git push origin master; \
	echo "Release - OK"

push-release:
	export GITHUB_TOKEN=$(GITHUB_TOKEN); \
	goreleaser release --rm-dist --skip-validate


#############  OTHERS  #############

COLOR_RESET = \033[0m
COLOR_COMMAND = \033[36m
COLOR_YELLOW = \033[33m

## Prints this help
help:
	printf "${COLOR_YELLOW}${PROJECT}\n------\n${COLOR_RESET}"
	awk '/^[a-zA-Z\-\_0-9\.%]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "${COLOR_COMMAND}$$ make %s${COLOR_RESET} %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST) | sort
	printf "\n"
