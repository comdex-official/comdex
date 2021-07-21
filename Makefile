export GO111MODULE=on

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git rev-parse --short HEAD)

build_tags = netgo
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))
ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=comdexCore \
		  -X github.com/cosmos/cosmos-sdk/version.ServerName=comdexNode \
		  -X github.com/cosmos/cosmos-sdk/version.ClientName=comdexClient \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep) \

BUILD_FLAGS += -ldflags "${ldflags}"

GOBIN = $(shell go env GOPATH)/bin

all: verify build

install:
ifeq (${OS},Windows_NT)
	go build -mod=readonly ${BUILD_FLAGS} -o ${GOBIN}/comdexClient.exe ./client
	go build -mod=readonly ${BUILD_FLAGS} -o ${GOBIN}/comdexNode.exe ./node
else
	go build -mod=readonly ${BUILD_FLAGS} -o ${GOBIN}/comdexClient ./client
	go build -mod=readonly ${BUILD_FLAGS} -o ${GOBIN}/comdexNode ./node
endif

build:
ifeq (${OS},Windows_NT)
	go build  ${BUILD_FLAGS} -o ${GOBIN}/comdexClient.exe ./client
	go build  ${BUILD_FLAGS} -o ${GOBIN}/comdexNode.exe ./node
else
	go build  ${BUILD_FLAGS} -o ${GOBIN}/comdexClient ./client
	go build  ${BUILD_FLAGS} -o ${GOBIN}/comdexNode ./node
endif

verify:
	@echo "verifying modules"
	@go mod verify

.PHONY: all install build verify