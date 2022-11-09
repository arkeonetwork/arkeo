
########################################################################################
# Environment Checks
########################################################################################

CHECK_ENV:=$(shell ./scripts/check-env.sh)
ifneq ($(CHECK_ENV),)
$(error Check environment dependencies.)
endif

########################################################################################
# Config
########################################################################################

.PHONY: build test tools export healthcheck run-mocknet build-mocknet stop-mocknet ps-mocknet reset-mocknet logs-mocknet openapi

# compiler flags
NOW=$(shell date +'%Y-%m-%d_%T')
COMMIT:=$(shell git log -1 --format='%H')
VERSION:=$(shell cat version)
TAG?=testnet
ldflags = -X gitlab.com/cbarraford/mercury/config.Version=$(VERSION) \
          -X gitlab.com/cbarraford/mercury/config.GitCommit=$(COMMIT) \
          -X gitlab.com/cbarraford/mercury/config.BuildTime=${NOW} \
		  -X github.com/cosmos/cosmos-sdk/version.Name=Mercury \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=mercury \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X github.com/cosmos/cosmos-sdk/version.BuildTags=$(TAG)

# golang settings
TEST_DIR?="./..."
BUILD_FLAGS := -ldflags '$(ldflags)' -tags ${TAG}
TEST_BUILD_FLAGS := -parallel=1 -tags=mocknet
GOBIN?=${GOPATH}/bin
BINARIES=./cmd/mercuryd ./cmd/sentinel

# pull branch name from CI if unset and available
ifdef CI_COMMIT_BRANCH
	BRANCH?=$(shell echo ${CI_COMMIT_BRANCH})
	BUILDTAG?=$(shell echo ${CI_COMMIT_BRANCH})
endif

# image build settings
BRANCH?=$(shell git rev-parse --abbrev-ref HEAD)
GITREF=$(shell git rev-parse --short HEAD)
BUILDTAG?=$(shell git rev-parse --abbrev-ref HEAD)

########################################################################################
# Targets
########################################################################################

# ------------------------------ Build ------------------------------

build:
	go build ${BUILD_FLAGS} ${BINARIES}

install:
	go install ${BUILD_FLAGS} ${BINARIES}

# ------------------------------ Housekeeping ------------------------------

format:
	@git ls-files '*.go' | grep -v -e '^docs/' | xargs gofumpt -w

lint:
	@./scripts/lint.sh
	@go build ${BINARIES}
	@./scripts/trunk check --no-fix --upstream origin/master

lint-ci:
	@./scripts/lint.sh
	@go build ${BINARIES}
	@./scripts/trunk check --all --no-progress --monitor=false

# ------------------------------ Testing ------------------------------

test-coverage:
	@go test ${TEST_BUILD_FLAGS} -v -coverprofile=coverage.txt -covermode count ${TEST_DIR}
	sed -i '/\.pb\.go:/d' coverage.txt

coverage-report: test-coverage
	@go tool cover -html=coverage.txt

tools:
	go install ./tools/signhere

test-coverage-sum:
	@go run gotest.tools/gotestsum --junitfile report.xml --format testname -- ${TEST_BUILD_FLAGS} -v -coverprofile=coverage.txt -covermode count ${TEST_DIR}
	sed -i '/\.pb\.go:/d' coverage.txt
	@GOFLAGS='${TEST_BUILD_FLAGS}' go run github.com/boumenot/gocover-cobertura < coverage.txt > coverage.xml
	@go tool cover -func=coverage.txt
	@go tool cover -html=coverage.txt -o coverage.html

test:
	@CGO_ENABLED=0 go test ${TEST_BUILD_FLAGS} ${TEST_DIR}

test-race:
	@go test -race ${TEST_BUILD_FLAGS} ${TEST_DIR}

test-watch:
	@gow -c test ${TEST_BUILD_FLAGS} ${TEST_DIR}

