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



.PHONY: all install build verify

.PHONY: clean
clean:
	rm -rf ./bin ./vendor

.PHONY: install
install: mod-vendor
	go install -mod=readonly -tags="${BUILD_TAGS}" -ldflags="${LD_FLAGS}" ./node/cmd

.PHONY: mod-vendor
mod-vendor: tools
	@go mod vendor
	@modvendor -copy="**/*.proto" -include=github.com/cosmos/cosmos-sdk/proto,github.com/cosmos/cosmos-sdk/third_party/proto

.PHONY: proto-gen
proto-gen:
	@.script/protocgen.sh

.PHONY: proto-lint
proto-lint:
	@find proto -name *.proto -exec clang-format-12 -i {} \;

.PHONY: tools
tools:
	@go install github.com/bufbuild/buf/cmd/buf@v0.37.0
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.27.0
	@go install github.com/goware/modvendor@v0.3.0
	@go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.16.0