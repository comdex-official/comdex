export GO111MODULE=on

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git rev-parse --short HEAD)

build_tags = netgo
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))
ldflags = -X github.com/cosmos/cosmos-sdk/version.AppName=comdex \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep) \

BUILD_FLAGS += -ldflags "${ldflags}"

GOBIN = $(shell go env GOPATH)/bin

all: verify build

install:
ifeq (${OS},Windows_NT)
	go build -mod=readonly ${BUILD_FLAGS} -o ${GOBIN}/comdex.exe ./node
else
	go build -mod=readonly ${BUILD_FLAGS} -o ${GOBIN}/comdex ./node
endif

build:
ifeq (${OS},Windows_NT)
	go build  ${BUILD_FLAGS} -o ${GOBIN}/comdex.exe ./node
else
	go build  ${BUILD_FLAGS} -o ${GOBIN}/comdex ./node
endif

verify:
	@echo "verifying modules"
	@go mod verify


DOCKER := $(shell which docker)

DOCKER_IMAGE_NAME = comdex-official/comdex
DOCKER_TAG_NAME = latest
DOCKER_CONTAINER_NAME = comdex-container
DOCKER_CMD ?= "/bin/sh"
DOCKER_VOLUME = -v $(CURDIR):/usr/local/app

.PHONY: all install build verify docker-run

proto-gen:
	@echo "Generating Protobuf files"
	$(DOCKER) run --rm -v $(shell go list -f "{{ .Dir }}" \
	-m github.com/cosmos/cosmos-sdk):/workspace/cosmos_sdk_dir\
	 --env COSMOS_SDK_DIR=/workspace/cosmos_sdk_dir \
	 -v $(CURDIR):/workspace --workdir /workspace \
	 tendermintdev/sdk-proto-gen sh ./.script/protocgen.sh