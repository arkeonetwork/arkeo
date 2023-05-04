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

.PHONY: build test tools

# compiler flags
IMAGE="arkeo"
PROJECT_NAME= arkeo
DOCKER         := $(shell which docker)
NOW=$(shell date +'%Y-%m-%d_%T')
COMMIT:=$(shell git log -1 --format='%H')
VERSION:=$(shell cat version)
TAG?=latest
ldflags = -X github.com/arkeonetwork/arkeo/config.Version=$(VERSION) \
          -X github.com/arkeonetwork/arkeo/config.GitCommit=$(COMMIT) \
          -X github.com/arkeonetwork/arkeo/config.BuildTime=${NOW} \
		  -X github.com/cosmos/cosmos-sdk/version.Name=Arkeo \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=arkeo \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X github.com/cosmos/cosmos-sdk/version.BuildTags=$(TAG)

# golang settings
TEST_DIR?="./..."
BUILD_FLAGS := -ldflags '$(ldflags)' -tags ${TAG}
TEST_BUILD_FLAGS := -parallel=1 -tags=mocknet -test.short=true
GOBIN?=${GOPATH}/bin
BINARIES=./cmd/arkeod ./cmd/sentinel ./cmd/directory/directory-indexer ./cmd/directory/directory-api

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

# ------------------------------ Docker Build ------------------------------

docker-build: proto-gen
	@docker build . --file Dockerfile -t ${IMAGE}:${TAG}

docker-run:
	@docker run --rm -it -p 1317:1317 -p 26656:26656 -p 26657:26657 ${IMAGE}:${TAG}

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

# ------------------------------ Unit Tests ------------------------------

test-coverage:
	@go test ${TEST_BUILD_FLAGS} -v -coverprofile=coverage.txt -covermode count ${TEST_DIR}
	sed -i '/\.pb\.go:/d' coverage.txt

coverage-report: test-coverage
	@go tool cover -html=coverage.txt

tools:
	go install ${BUILD_FLAGS} ./tools/signhere ./tools/curleo

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

# ------------------------------ Regression Tests ------------------------------

test-regression:
	@DOCKER_BUILDKIT=1 docker-compose -f ./test/regression/docker-compose.yml run --build arkeo

test-regression-ci:
	docker-compose --version
	DOCKER_BUILDKIT=1 docker-compose -f ./test/regression/docker-compose.yml run --build -- arkeo make _test-regression

test-regression-coverage:
	@go tool cover -html=test/regression/mnt/coverage/coverage.txt

# internal target used in docker build
_build-test-regression:
	@go install -ldflags '$(ldflags)' -tags=testnet,regtest ${BINARIES}
	@go build -ldflags '$(ldflags)' -cover -tags=testnet,regtest -o /regtest/cover-arkeod ./cmd/arkeod
	@go build -ldflags '$(ldflags)' -cover -tags=testnet,regtest -o /regtest/cover-sentinel ./cmd/sentinel
	@go build -ldflags '$(ldflags)' -cover -tags=testnet,regtest -o /regtest/cover-directory-api ./cmd/directory/directory-api
	@go build -ldflags '$(ldflags)' -cover -tags=testnet,regtest -o /regtest/cover-directory-indexer ./cmd/directory/directory-indexer
	@go build -ldflags '$(ldflags)' -tags testnet -o /regtest/regtest ./test/regression/cmd

# internal target used in test run
_test-regression:
	@rm -rf /mnt/coverage && mkdir -p /mnt/coverage
	@cd test/regression && /regtest/regtest
	@go tool covdata textfmt -i /mnt/coverage -o /mnt/coverage/coverage.txt
	@go tool cover -func /mnt/coverage/coverage.txt > /mnt/coverage/func-coverage.txt
	@awk '/^total:/ {print "Regression Coverage: " $$3}' /mnt/coverage/func-coverage.txt
	@chown -R ${UID}:${GID} /mnt

########################################################################################
# Protobuf
########################################################################################

DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf:1.9.0

containerProtoVer=v0.7
containerProtoImage=tendermintdev/sdk-proto-gen:$(containerProtoVer)
containerProtoGen=$(PROJECT_NAME)-proto-gen-$(containerProtoVer)
containerProtoFmt=$(PROJECT_NAME)-proto-fmt-$(containerProtoVer)
containerProtoGenSwagger=$(PROJECT_NAME)-proto-gen-swagger-$(containerProtoVer)

proto-all: proto-format proto-lint proto-gen proto-swagger-gen
.PHONY: proto-all proto-gen proto-lint proto-check-breaking proto-format proto-swagger-gen

proto-gen:
	@echo "Generating Protobuf files"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoGen}$$"; then docker start -a $(containerProtoGen); else docker run --name $(containerProtoGen) -v $(CURDIR):/workspace --workdir /workspace $(containerProtoImage) \
		sh ./scripts/protocgen.sh; fi

proto-swagger-gen:
	@echo "Generating Swagger of Protobuf"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoGenSwagger}$$"; then docker start -a $(containerProtoGenSwagger); else docker run --name $(containerProtoGenSwagger) -v $(CURDIR):/workspace --workdir /workspace $(containerProtoImage) \
		sh ./scripts/protoc-swagger-gen.sh; fi

proto-format:
	@echo "Formatting Protobuf files"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoFmt}$$"; then docker start -a $(containerProtoFmt); else docker run --name $(containerProtoFmt) -v $(CURDIR):/workspace --workdir /workspace tendermintdev/docker-build-proto \
		find ./ -name "*.proto" -exec sh -c 'clang-format -style=file -i {}' \; ; fi

proto-lint:
	@echo "Linting Protobuf files"
	@$(DOCKER_BUF) lint --error-format=json

proto-check-breaking:
	@echo "Checking for breaking changes"
	@$(DOCKER_BUF) breaking --against $(HTTPS_GIT)#branch=main

# arkeod binaries
dist:
	rm -rf bin && mkdir -p bin/linux_amd64 bin/linux_arm64 bin/darwin_amd64 bin/darwin_arm64
	env GOOS=linux GOARCH=amd64 go build -o bin/linux_amd64 ${BUILD_FLAGS} ./cmd/arkeod
	env GOOS=linux GOARCH=arm64 go build -o bin/linux_arm64 ${BUILD_FLAGS} ./cmd/arkeod
	env GOOS=darwin GOARCH=amd64 go build -o bin/darwin_amd64 ${BUILD_FLAGS} ./cmd/arkeod
	env GOOS=darwin GOARCH=arm64 go build -o bin/darwin_arm64 ${BUILD_FLAGS} ./cmd/arkeod

	env GOOS=linux GOARCH=amd64 go build -o bin/linux_amd64 ${BUILD_FLAGS} ./tools/curleo
	env GOOS=linux GOARCH=arm64 go build -o bin/linux_arm64 ${BUILD_FLAGS} ./tools/curleo
	env GOOS=darwin GOARCH=amd64 go build -o bin/darwin_amd64 ${BUILD_FLAGS} ./tools/curleo
	env GOOS=darwin GOARCH=arm64 go build -o bin/darwin_arm64 ${BUILD_FLAGS} ./tools/curleo

	env GOOS=linux GOARCH=amd64 go build -o bin/linux_amd64 ${BUILD_FLAGS} ./tools/signhere
	env GOOS=linux GOARCH=arm64 go build -o bin/linux_arm64 ${BUILD_FLAGS} ./tools/signhere
	env GOOS=darwin GOARCH=amd64 go build -o bin/darwin_amd64 ${BUILD_FLAGS} ./tools/signhere
	env GOOS=darwin GOARCH=arm64 go build -o bin/darwin_arm64 ${BUILD_FLAGS} ./tools/signhere
	
	/usr/bin/tar -C bin/linux_amd64 --uid 0  --no-fflags --no-mac-metadata --strip-components 1 -czvf bin/arkeo_linux_amd64.tar.gz .
	/usr/bin/tar -C bin/linux_arm64 --uid 0  --no-fflags --no-mac-metadata --strip-components 1 -czvf bin/arkeo_linux_arm64.tar.gz .
	/usr/bin/tar -C bin/darwin_amd64 --uid 0  --no-fflags --no-mac-metadata --strip-components 1 -czvf bin/arkeo_darwin_amd64.tar.gz .
	/usr/bin/tar -C bin/darwin_arm64 --uid 0  --no-fflags --no-mac-metadata --strip-components 1 -czvf bin/arkeo_darwin_arm64.tar.gz .

	pushd bin && \
	sha256sum arkeo_linux_amd64.tar.gz > ../docs/testing/sums/arkeo_linux_amd64.sha256 && \
	sha256sum arkeo_linux_arm64.tar.gz > ../docs/testing/sums/arkeo_linux_arm64.sha256 && \
	sha256sum arkeo_darwin_amd64.tar.gz > ../docs/testing/sums/arkeo_darwin_amd64.sha256 && \
	sha256sum arkeo_darwin_arm64.tar.gz > ../docs/testing/sums/arkeo_darwin_arm64.sha256 && \
	popd

	rm -rf bin/linux_amd64 bin/linux_arm64 bin/darwin_amd64 bin/darwin_arm64
