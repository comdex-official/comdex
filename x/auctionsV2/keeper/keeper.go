package keeper

import (
	"context"
	"fmt"

	"github.com/comdex-official/comdex/x/liquidationsV2/expected"
	"github.com/comdex-official/comdex/x/liquidationsV2/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"
)

type (
	Keeper struct {
		cdc         codec.BinaryCodec
		storeKey    sdk.StoreKey
		memKey      sdk.StoreKey
		paramstore  paramtypes.Subspace
		liquidation expected.LiquidationV2Keeper
		bankKeeper  types.BankKeeper
	}
)

func (k Keeper) Params(ctx context.Context, request *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,
	liquidation expected.LiquidationV2Keeper,
	bankKeeper types.BankKeeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{

		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
		liquidation: liquidation,
		bankKeeper: bankKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
