package keeper

import (
	"context"
	"fmt"

<<<<<<< HEAD
	"github.com/comdex-official/comdex/x/liquidationsV2/expected"
	"github.com/comdex-official/comdex/x/liquidationsV2/types"
=======
	"github.com/tendermint/tendermint/libs/log"

	"github.com/comdex-official/comdex/x/auctionsV2/expected"
	"github.com/comdex-official/comdex/x/auctionsV2/types"
>>>>>>> 6cc24557 (updating expected keepers)
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"
)

type (
	Keeper struct {
<<<<<<< HEAD
		cdc         codec.BinaryCodec
		storeKey    sdk.StoreKey
		memKey      sdk.StoreKey
		paramstore  paramtypes.Subspace
		liquidation expected.LiquidationV2Keeper
		bankKeeper  types.BankKeeper
=======
		cdc            codec.BinaryCodec
		storeKey       sdk.StoreKey
		memKey         sdk.StoreKey
		paramstore     paramtypes.Subspace
		bankKeeper     types.BankKeeper
		liquidationsV2 expected.LiquidationsV2Keeper
>>>>>>> 6cc24557 (updating expected keepers)
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
<<<<<<< HEAD
	liquidation expected.LiquidationV2Keeper,
=======
>>>>>>> 6cc24557 (updating expected keepers)
	bankKeeper types.BankKeeper,
	liquidationsV2Keeper expected.LiquidationsV2Keeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{

<<<<<<< HEAD
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
		liquidation: liquidation,
		bankKeeper: bankKeeper,
=======
		cdc:            cdc,
		storeKey:       storeKey,
		memKey:         memKey,
		paramstore:     ps,
		bankKeeper:     bankKeeper,
		liquidationsV2: liquidationsV2Keeper,
>>>>>>> 6cc24557 (updating expected keepers)
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
