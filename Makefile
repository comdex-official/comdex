#!/usr/bin/make -f

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%H')

# don't override user values
ifeq (,$(VERSION))
  VERSION := $(shell git describe --tags)
  # if VERSION is empty, then populate it with branch's name and raw commit hash
  ifeq (,$(VERSION))
    VERSION := $(BRANCH)-$(COMMIT)
  endif
endif

PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
LEDGER_ENABLED ?= true
SDK_PACK := $(shell go list -m github.com/cosmos/cosmos-sdk | sed  's/ /\@/g')
TM_VERSION := $(shell go list -m github.com/tendermint/tendermint | sed 's:.* ::') # grab everything after the space in "github.com/tendermint/tendermint v0.34.7"
BUILDDIR ?= $(CURDIR)/build

export GO111MODULE = on

# process build tags

build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq (cleveldb,$(findstring cleveldb,$(COMDEX_BUILD_OPTIONS)))
  build_tags += gcc cleveldb
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=comdex \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=comdex \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)" \
			-X github.com/tendermint/tendermint/version.TMCoreSemVer=$(TM_VERSION)

ifeq (cleveldb,$(findstring cleveldb,$(COMDEX_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ifeq ($(LINK_STATICALLY),true)
  ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif
ifeq (,$(findstring nostrip,$(COMDEX_BUILD_OPTIONS)))
  ldflags += -w -s
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'
# check for nostrip option
ifeq (,$(findstring nostrip,$(COMDEX_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
endif

#$(info $$BUILD_FLAGS is [$(BUILD_FLAGS)])


all: install

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/comdex

build:
	go build $(BUILD_FLAGS) -o bin/comdex ./cmd/comdex



###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

test: test-unit
test-all: test-unit test-ledger-mock test-race test-cover

TEST_PACKAGES=./...
TEST_TARGETS := test-unit test-unit-amino test-unit-proto test-ledger-mock test-race test-ledger test-race

# Test runs-specific rules. To add a new test target, just add
# a new rule, customise ARGS or TEST_PACKAGES ad libitum, and
# append the new rule to the TEST_TARGETS list.
test-unit: ARGS=-tags='cgo ledger test_ledger_mock norace'
test-unit-amino: ARGS=-tags='ledger test_ledger_mock test_amino norace'
test-ledger: ARGS=-tags='cgo ledger norace'
test-ledger-mock: ARGS=-tags='ledger test_ledger_mock norace'
test-race: ARGS=-race -tags='cgo ledger test_ledger_mock'
test-race: TEST_PACKAGES=$(PACKAGES_NOSIMULATION)
$(TEST_TARGETS): run-tests

# check-* compiles and collects tests without running them
# note: go test -c doesn't support multiple packages yet (https://github.com/golang/go/issues/15513)
CHECK_TEST_TARGETS := check-test-unit check-test-unit-amino
check-test-unit: ARGS=-tags='cgo ledger test_ledger_mock norace'
check-test-unit-amino: ARGS=-tags='ledger test_ledger_mock test_amino norace'
$(CHECK_TEST_TARGETS): EXTRA_ARGS=-run=none
$(CHECK_TEST_TARGETS): run-tests

run-tests:
ifneq (,$(shell which tparse 2>/dev/null))
	go test -mod=readonly -json $(ARGS) $(EXTRA_ARGS) $(TEST_PACKAGES) | tparse
else
	go test -mod=readonly $(ARGS)  $(EXTRA_ARGS) $(TEST_PACKAGES)
endif

.PHONY: run-tests test test-all $(TEST_TARGETS)

test-sim-nondeterminism:
	@echo "Running non-determinism test..."
	@go test -mod=readonly $(SIMAPP) -run TestAppStateDeterminism -Enabled=true \
		-NumBlocks=100 -BlockSize=200 -Commit=true -Period=0 -v -timeout 24h

test-sim-custom-genesis-fast:
	@echo "Running custom genesis simulation..."
	@echo "By default, ${HOME}/.comdex/config/genesis.json will be used."
	@go test -mod=readonly $(SIMAPP) -run TestFullAppSimulation -Genesis=${HOME}/.comdex/config/genesis.json \
		-Enabled=true -NumBlocks=100 -BlockSize=200 -Commit=true -Seed=99 -Period=5 -v -timeout 24h

test-sim-import-export: runsim
	@echo "Running application import/export simulation. This may take several minutes..."
	@$(BINDIR)/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) -ExitOnFail 50 5 TestAppImportExport

test-sim-after-import: runsim
	@echo "Running application simulation-after-import. This may take several minutes..."
	@$(BINDIR)/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) -ExitOnFail 50 5 TestAppSimulationAfterImport

test-sim-custom-genesis-multi-seed: runsim
	@echo "Running multi-seed custom genesis simulation..."
	@echo "By default, ${HOME}/.comdex/config/genesis.json will be used."
	@$(BINDIR)/runsim -Genesis=${HOME}/.comdex/config/genesis.json -SimAppPkg=$(SIMAPP) -ExitOnFail 400 5 TestFullAppSimulation

test-sim-multi-seed-long: runsim
	@echo "Running long multi-seed application simulation. This may take awhile!"
	@$(BINDIR)/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) -ExitOnFail 500 50 TestFullAppSimulation

test-sim-multi-seed-short: runsim
	@echo "Running short multi-seed application simulation. This may take awhile!"
	@$(BINDIR)/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) -ExitOnFail 50 10 TestFullAppSimulation

test-sim-benchmark-invariants:
	@echo "Running simulation invariant benchmarks..."
	@go test -mod=readonly $(SIMAPP) -benchmem -bench=BenchmarkInvariants -run=^$ \
	-Enabled=true -NumBlocks=1000 -BlockSize=200 \
	-Period=1 -Commit=true -Seed=57 -v -timeout 24h

.PHONY: \
test-sim-nondeterminism \
test-sim-custom-genesis-fast \
test-sim-import-export \
test-sim-after-import \
test-sim-custom-genesis-multi-seed \
test-sim-multi-seed-short \
test-sim-multi-seed-long \
test-sim-benchmark-invariants

SIM_NUM_BLOCKS ?= 500
SIM_BLOCK_SIZE ?= 200
SIM_COMMIT ?= true

test-sim-benchmark:
	@echo "Running application benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -mod=readonly -benchmem -run=^$$ $(SIMAPP) -bench ^BenchmarkFullAppSimulation$$  \
		-Enabled=true -NumBlocks=$(SIM_NUM_BLOCKS) -BlockSize=$(SIM_BLOCK_SIZE) -Commit=$(SIM_COMMIT) -timeout 24h

test-sim-profile:
	@echo "Running application benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -mod=readonly -benchmem -run=^$$ $(SIMAPP) -bench ^BenchmarkFullAppSimulation$$ \
		-Enabled=true -NumBlocks=$(SIM_NUM_BLOCKS) -BlockSize=$(SIM_BLOCK_SIZE) -Commit=$(SIM_COMMIT) -timeout 24h -cpuprofile cpu.out -memprofile mem.out

.PHONY: test-sim-profile test-sim-benchmark

test-cover:
	@export VERSION=$(VERSION); bash -x scripts/test_cover.sh
.PHONY: test-cover

benchmark:
	@go test -mod=readonly -bench=. $(PACKAGES_NOSIMULATION)
.PHONY: benchmark