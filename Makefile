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
DOCKER := $(shell which docker)
NOW=$(shell date +'%Y-%m-%d_%T')
COMMIT:=$(shell git log -1 --format='%H')
SHORT_COMMIT:=$(shell git rev-parse --short=7 HEAD)
CHAIN_VERSION:=$(shell cat chain.version)
LATEST_VERSION=$(shell curl -s "https://api.github.com/repos/arkeonetwork/arkeo/releases/latest" | jq -r '.tag_name')
SENTINEL_VERSION:=$(shell cat sentinel.version)
TAG?=latest
ldflags = -X github.com/arkeonetwork/arkeo/x/arkeo/configs.Version=$(LATEST_VERSION) \
		  -X github.com/arkeonetwork/arkeo/sentinel.Version=$(SENTINEL_VERSION) \
          -X github.com/arkeonetwork/arkeo/x/arkeo/configs.GitCommit=$(COMMIT) \
          -X github.com/arkeonetwork/arkeo/x/arkeo/configs.BuildTime=${NOW} \
		  -X github.com/cosmos/cosmos-sdk/version.Name=Arkeo \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=arkeo \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(LATEST_VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X github.com/cosmos/cosmos-sdk/version.BuildTags=$(TAG)

# golang settings
TEST_DIR?="./..."
BUILD_FLAGS := -ldflags '$(ldflags)' -tags ${TAG}
TEST_BUILD_FLAGS := -parallel=1 -tags=mocknet -test.short=true
GOBIN?=${GOPATH}/bin
BINARIES=./cmd/arkeod ./cmd/sentinel ./cmd/directory/indexer ./cmd/directory/api
TESNET_BUILD_FLAGS :=  -ldflags '$(ldflags)'  -tags=testnet
# pull branch name from CI if unset and available
ifdef CI_COMMIT_BRANCH
	BRANCH?=$(shell echo ${CI_COMMIT_BRANCH})
	BUILDTAG?=$(shell echo ${CI_COMMIT_BRANCH})
endif

# image build settings
BRANCH?=$(shell git rev-parse --abbrev-ref HEAD)
GITREF=$(shell git rev-parse --short HEAD)
BUILDTAG?=$(shell git rev-parse --abbrev-ref HEAD)
GORELEASER_CROSS_VERSION = v1.21.9
GORELEASER_VERSION = v1.21.0
DOCKER := $(shell which docker)

# Release Env Variable
RELEASE ?= false
GORELEASER_SKIP_VALIDATE ?= false
GORELEASER_SKIP_PUBLISH ?= false

########################################################################################
# Targets
########################################################################################

# ------------------------------ Build ------------------------------

build:
	go build ${BUILD_FLAGS} ${BINARIES}

install:
	go install ${BUILD_FLAGS} ${BINARIES}

# ------------------------------ Docker Build ------------------------------

# Detect OS and architecture
OS := $(shell uname -s)
ARCH := $(shell uname -m)

# Determine the Docker build command and image architecture based on OS and architecture
ifeq ($(OS),Darwin)
	BUILD_OS:=darwin
    IMAGE_ARCH := arm64
    DOCKER_BUILD := docker-build-cross
	TESTNET_BUILD:= release-dry-run-cross
else ifeq ($(OS),Linux)
    ifeq ($(ARCH),x86_64)
		BUILD_OS:=linux
        IMAGE_ARCH := amd64
        DOCKER_BUILD := docker-build
		TESTNET_BUILD:= release-dry-run
    else ifeq ($(ARCH),arm64)
		BUILD_OS:=linux
        IMAGE_ARCH := arm64
        DOCKER_BUILD := docker-build-cross
		TESTNET_BUILD:= release-dry-run-cross
    else

    endif
else
    $(error Unsupported OS: $(OS))
endif

# Docker build target
build-docker:
	@$(MAKE) $(DOCKER_BUILD)

docker-build:
	$(DOCKER) run \
		--rm \
		-e BUILD_TAG=$(TAG) \
		-e RELEASE=$(RELEASE) \
		-e GITHUB_TOKEN="$(GITHUB_TOKEN)" \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/github.com/arkeonetwork/arkeo \
		-w /go/src/github.com/arkeonetwork/arkeo \
		ghcr.io/goreleaser/goreleaser:$(GORELEASER_VERSION) \
		--clean \
		--snapshot

docker-build-cross:
	$(DOCKER) run \
		--rm \
		-e CGO_ENABLED=1 \
		-e BUILD_TAG=$(TAG) \
		-e RELEASE=$(RELEASE)\
		-e SENTINEL_VERSION=$(SENTINEL_VERSION)\
		-e GITHUB_TOKEN="$(GITHUB_TOKEN)" \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		-f .goreleaser-cross.yaml \
		--clean \
		--snapshot 

build-docker-localnet:
	@docker build . --file Dockerfile.localnet -t ${IMAGE}:${TAG}

# localnet: build-docker
# 	IMAGE_TAG=$(SHORT_COMMIT)-$(IMAGE_ARCH) docker-compose -f docker-compose-localnet.yaml  up


localnet: build-docker-localnet
	IMAGE_TAG=${IMAGE}:${TAG} docker-compose -f docker-compose-localnet.yaml  up
# ------------------------------    Testnet   ------------------------------

install-testnet-binary:
	@TAG=testnet $(MAKE) $(TESTNET_BUILD)
	@cp dist/arkeod-$(BUILD_OS)-$(IMAGE_ARCH)_$(BUILD_OS)_$(IMAGE_ARCH)$(if $(filter amd64,$(IMAGE_ARCH)),_v1)/arkeod $(HOME)/go/bin

testnet-fullnode:
	@docker run --rm -it -p 1317:1317 -p 26656:26656 -p 26657:26657 -v ./scripts:/scripts --entrypoint /scripts/fullnode.sh ghcr.io/arkeonetwork/arkeo:latest

# ------------------------------ Housekeeping ------------------------------

lint:
	golangci-lint run
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/lcd/statik/statik.go" | xargs goimports -w -local github.com/arkeonetwork/arkeo

# ------------------------------ Unit Tests ------------------------------

test-coverage:
	@go test ${TEST_BUILD_FLAGS} -v -coverprofile=coverage.txt -covermode count ${TEST_DIR}
	sed -i '/\.pb\.go:/d' coverage.txt

coverage-report: test-coverage
	@go tool cover -html=coverage.txt

tools:
	go install ${BUILD_FLAGS} ./tools/signhere ./tools/curleo ./tools/mock-daemon

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
	@DOCKER_BUILDKIT=1 docker-compose -f ./test/regression/docker-compose.yml run -i -p 1317:1317 -p 3636:3636 -p 26657 --rm --build arkeo

test-regression-ci: test-regression

test-regression-coverage:
	@go tool cover -html=test/regression/mnt/coverage/coverage.txt

# internal target used in docker build
_build-test-regression:
	@go install -ldflags '$(ldflags)' -tags=testnet,regtest ${BINARIES}
	@go install ${BUILD_FLAGS} ./tools/mock-daemon
	@go build -ldflags '$(ldflags)' -cover -tags=testnet,regtest -o /regtest/cover-arkeod ./cmd/arkeod
	@go build -ldflags '$(ldflags)' -cover -tags=testnet,regtest -o /regtest/cover-sentinel ./cmd/sentinel
	@go build -ldflags '$(ldflags)' -cover -tags=testnet,regtest -o /regtest/cover-directory-api ./cmd/directory/api
	@go build -ldflags '$(ldflags)' -cover -tags=testnet,regtest -o /regtest/cover-directory-indexer ./cmd/directory/indexer
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

DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf:1.36.0

containerProtoVer=0.15.0
containerProtoImage=ghcr.io/cosmos/proto-builder:$(containerProtoVer)
containerProtoGen=$(PROJECT_NAME)-proto-gen-$(containerProtoVer)
containerProtoFmt=$(PROJECT_NAME)-proto-fmt-$(containerProtoVer)
containerProtoGenSwagger=$(PROJECT_NAME)-proto-gen-swagger-$(containerProtoVer)

proto-all: proto-format proto-lint proto-gen proto-swagger-gen
.PHONY: proto-all proto-gen proto-lint proto-check-breaking proto-format proto-swagger-gen

protob:
	@./scripts/protocgen.sh

proto-gen:
	@echo "Generating Protobuf files"
	$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(containerProtoImage) sh ./scripts/protocgen.sh
	./scripts/dontcover.sh ./x/arkeo
	./scripts/dontcover.sh ./x/claim

proto-swagger-gen:
	@echo "Generating Swagger of Protobuf"
	$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(containerProtoImage) sh  ./scripts/protoc-swagger-gen.sh

proto-format:
	@echo "Formatting Protobuf files"
	$(DOCKER) run --rm -v $(CURDIR):/workspace \
	--workdir /workspace $(containerProtoImage) \
	find ./ -name *.proto -exec clang-format -i {} \;

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

PACKAGE_NAME          := github.com/arkeonetwork/arkeo
GOLANG_CROSS_VERSION  ?= v1.21.9

SYSROOT_DIR     ?= sysroots
SYSROOT_ARCHIVE ?= sysroots.tar.bz2

.PHONY: sysroot-pack
sysroot-pack:
	@tar cf - $(SYSROOT_DIR) -P | pv -s $[$(du -sk $(SYSROOT_DIR) | awk '{print $1}') * 1024] | pbzip2 > $(SYSROOT_ARCHIVE)

.PHONY: sysroot-unpack
sysroot-unpack:
	@pv $(SYSROOT_ARCHIVE) | pbzip2 -cd | tar -xf -

########################################################################################
#         Dry Run			  
########################################################################################

.PHONY: release-dry-run
release-dry-run-cross:
	$(DOCKER) run \
		--rm \
		-e CGO_ENABLED=1 \
		-e BUILD_TAG=$(TAG) \
		-e RELEASE=$(RELEASE)\
		-e GITHUB_TOKEN="$(GITHUB_TOKEN)" \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		-f .goreleaser-cross.yaml \
		--clean --skip=validate --skip=publish

release-dry-run:
	$(DOCKER) run \
		--rm \
		-e CGO_ENABLED=1 \
		-e BUILD_TAG=$(TAG) \
		-e RELEASE=$(RELEASE)\
		-e GITHUB_TOKEN="$(GITHUB_TOKEN)" \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser:${GORELEASER_VERSION} \
		--clean --skip=validate --skip=publish

sentinel-release-dry-run-cross:
	$(DOCKER) run \
		--rm \
		-e CGO_ENABLED=1 \
		-e BUILD_TAG=$(TAG) \
		-e RELEASE=$(RELEASE)\
		-e SENTINEL_VERSION=$(SENTINEL_VERSION)\
		-e GITHUB_TOKEN="$(GITHUB_TOKEN)" \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		-f .goreleaser-sentinel-cross.yaml \
		--clean --skip=validate --skip=publish

sentinel-release-dry-run:
	$(DOCKER) run \
		--rm \
		-e CGO_ENABLED=1 \
		-e BUILD_TAG=$(TAG) \
		-e RELEASE=$(RELEASE)\
		-e SENTINEL_VERSION=$(SENTINEL_VERSION)\
		-e GITHUB_TOKEN="$(GITHUB_TOKEN)" \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser:${GORELEASER_VERSION} \
		-f .goreleaser-sentinel.yaml \
		--clean --skip=validate --skip=publish

directory-release-dry-run-cross:
	$(DOCKER) run \
		--rm \
		-e CGO_ENABLED=1 \
		-e BUILD_TAG=$(TAG) \
		-e RELEASE=$(RELEASE)\
		-e SENTINEL_VERSION=$(SENTINEL_VERSION)\
		-e GITHUB_TOKEN="$(GITHUB_TOKEN)" \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		-f .goreleaser-directory-cross.yaml \
		--clean --skip=validate --skip=publish

directory-release-dry-run:
	$(DOCKER) run \
		--rm \
		-e CGO_ENABLED=1 \
		-e BUILD_TAG=$(TAG) \
		-e RELEASE=$(RELEASE)\
		-e SENTINEL_VERSION=$(SENTINEL_VERSION)\
		-e GITHUB_TOKEN="$(GITHUB_TOKEN)" \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser:${GORELEASER_VERSION} \
		-f .goreleaser-directory.yaml \
		--clean --skip=validate --skip=publish


########################################################################################
#     Mainnet Releases		 
########################################################################################

.PHONY: releases
release:
	$(DOCKER) run \
		--rm \
		-e CGO_ENABLED=1 \
		-e BUILD_TAG=$(TAG) \
		-e RELEASE=$(RELEASE)\
		-e GITHUB_TOKEN="$(GITHUB_TOKEN)" \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser:${GORELEASER_VERSION} \
		--clean \
		--skip-validate=$(GORELEASER_SKIP_VALIDATE) \
		--skip-publish=$(GORELEASER_SKIP_PUBLISH)

release-cross:
	$(DOCKER) run \
		--rm \
		-e CGO_ENABLED=1 \
		-e BUILD_TAG=$(TAG) \
		-e RELEASE=$(RELEASE)\
		-e GITHUB_TOKEN="$(GITHUB_TOKEN)" \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		-f .goreleaser-cross.yaml \
		--clean \
		--skip-validate=$(GORELEASER_SKIP_VALIDATE) \
		--skip-publish=$(GORELEASER_SKIP_PUBLISH)


sentinel-release-cross:
	$(DOCKER) run \
		--rm \
		-e CGO_ENABLED=1 \
		-e BUILD_TAG=$(TAG) \
		-e RELEASE=$(RELEASE)\
		-e SENTINEL_VERSION=$(SENTINEL_VERSION)\
		-e GITHUB_TOKEN="$(GITHUB_TOKEN)" \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		-f .goreleaser-sentinel-cross.yaml \
		--clean \
		--skip-validate=$(GORELEASER_SKIP_VALIDATE) \
		--skip-publish=$(GORELEASER_SKIP_PUBLISH)

sentinel-release:
	$(DOCKER) run \
		--rm \
		-e CGO_ENABLED=1 \
		-e BUILD_TAG=$(TAG) \
		-e RELEASE=$(RELEASE)\
		-e SENTINEL_VERSION=$(SENTINEL_VERSION)\
		-e GITHUB_TOKEN="$(GITHUB_TOKEN)" \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser:${GORELEASER_VERSION} \
		-f .goreleaser-sentinel.yaml \
		--clean \
		--skip-validate=$(GORELEASER_SKIP_VALIDATE) \
		--skip-publish=$(GORELEASER_SKIP_PUBLISH)

directory-release-cross:
	$(DOCKER) run \
		--rm \
		-e CGO_ENABLED=1 \
		-e BUILD_TAG=$(TAG) \
		-e RELEASE=$(RELEASE)\
		-e GITHUB_TOKEN="$(GITHUB_TOKEN)" \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		-f .goreleaser-directory-cross.yaml \
		--clean \
		--skip-validate=$(GORELEASER_SKIP_VALIDATE) \
		--skip-publish=$(GORELEASER_SKIP_PUBLISH)

directory-release:
	$(DOCKER) run \
		--rm \
		-e CGO_ENABLED=1 \
		-e BUILD_TAG=$(TAG) \
		-e RELEASE=$(RELEASE)\
		-e GITHUB_TOKEN="$(GITHUB_TOKEN)" \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser:${GORELEASER_VERSION} \
		-f .goreleaser-directory.yaml \
		--clean \
		--skip-validate=$(GORELEASER_SKIP_VALIDATE) \
		--skip-publish=$(GORELEASER_SKIP_PUBLISH)


########################################################################################
#     Testnet Releases		  
########################################################################################

release-testnet:
	$(DOCKER) run \
		--rm \
		-e CGO_ENABLED=1 \
		-e BUILD_TAG="testnet" \
		-e RELEASE=$(RELEASE)\
		-e GITHUB_TOKEN="$(GITHUB_TOKEN)" \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser:${GORELEASER_VERSION} \
		--clean \
		--skip-validate=$(GORELEASER_SKIP_VALIDATE) \
		--skip-publish=$(GORELEASER_SKIP_PUBLISH)

release-testnet-cross:
	$(DOCKER) run \
		--rm \
		-e CGO_ENABLED=1 \
		-e BUILD_TAG="testnet" \
		-e RELEASE=$(RELEASE)\
		-e GITHUB_TOKEN="$(GITHUB_TOKEN)" \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		-f .goreleaser-cross.yaml \
		--clean \
		--skip-validate=$(GORELEASER_SKIP_VALIDATE) \
		--skip-publish=$(GORELEASER_SKIP_PUBLISH)


sentinel-testnet-cross:
	$(DOCKER) run \
		--rm \
		-e CGO_ENABLED=1 \
		-e BUILD_TAG="testnet" \
		-e RELEASE=$(RELEASE)\
		-e SENTINEL_VERSION=$(SENTINEL_VERSION)\
		-e GITHUB_TOKEN="$(GITHUB_TOKEN)" \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		-f .goreleaser-sentinel-cross.yaml \
		--clean \
		--skip-validate=$(GORELEASER_SKIP_VALIDATE) \
		--skip-publish=$(GORELEASER_SKIP_PUBLISH)

sentinel-testnet-release:
	$(DOCKER) run \
		--rm \
		-e CGO_ENABLED=1 \
		-e BUILD_TAG="testnet" \
		-e RELEASE=$(RELEASE)\
		-e SENTINEL_VERSION=$(SENTINEL_VERSION)\
		-e GITHUB_TOKEN="$(GITHUB_TOKEN)" \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser:${GORELEASER_VERSION} \
		-f .goreleaser-sentinel.yaml \
		--clean \
		--skip-validate=$(GORELEASER_SKIP_VALIDATE) \
		--skip-publish=$(GORELEASER_SKIP_PUBLISH)

directory-testnet-cross:
	$(DOCKER) run \
		--rm \
		-e CGO_ENABLED=1 \
		-e BUILD_TAG="testnet" \
		-e RELEASE=$(RELEASE)\
		-e GITHUB_TOKEN="$(GITHUB_TOKEN)" \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		-f .goreleaser-directory-cross.yaml \
		--clean \
		--skip-validate=$(GORELEASER_SKIP_VALIDATE) \
		--skip-publish=$(GORELEASER_SKIP_PUBLISH)

directory-testnet:
	$(DOCKER) run \
		--rm \
		-e CGO_ENABLED=1 \
		-e BUILD_TAG="testnet" \
		-e RELEASE=$(RELEASE)\
		-e GITHUB_TOKEN="$(GITHUB_TOKEN)" \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser:${GORELEASER_VERSION} \
		-f .goreleaser-directory.yaml \
		--clean \
		--skip-validate=$(GORELEASER_SKIP_VALIDATE) \
		--skip-publish=$(GORELEASER_SKIP_PUBLISH)