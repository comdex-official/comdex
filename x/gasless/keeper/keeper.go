package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/comdex-official/comdex/x/gasless/expected"
	"github.com/comdex-official/comdex/x/gasless/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
)

// Keeper of the gasless store.
type Keeper struct {
	cdc               codec.BinaryCodec
	storeKey          storetypes.StoreKey
	paramSpace        paramstypes.Subspace
	interfaceRegistry codectypes.InterfaceRegistry

	accountKeeper expected.AccountKeeper
	bankKeeper    expected.BankKeeper
	wasmKeeper    *wasmkeeper.Keeper
}

// NewKeeper creates a new gasless Keeper instance.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	paramSpace paramstypes.Subspace,
	interfaceRegistry codectypes.InterfaceRegistry,
	accountKeeper expected.AccountKeeper,
	bankKeeper expected.BankKeeper,
	wasmKeeper *wasmkeeper.Keeper,
) Keeper {
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:               cdc,
		storeKey:          storeKey,
		paramSpace:        paramSpace,
		interfaceRegistry: interfaceRegistry,
		accountKeeper:     accountKeeper,
		bankKeeper:        bankKeeper,
		wasmKeeper:        wasmKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
