package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	storetypes "cosmossdk.io/store/types"
	"github.com/comdex-official/comdex/x/tokenmint/expected"
	"github.com/comdex-official/comdex/x/tokenmint/types"
)

type Keeper struct {
	cdc   codec.BinaryCodec
	key   storetypes.StoreKey
	bank  expected.BankKeeper
	asset expected.AssetKeeper
}

func NewKeeper(cdc codec.BinaryCodec, key storetypes.StoreKey, bank expected.BankKeeper, asset expected.AssetKeeper) Keeper {
	return Keeper{
		cdc:   cdc,
		key:   key,
		bank:  bank,
		asset: asset,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Store(ctx sdk.Context) storetypes.KVStore {
	return ctx.KVStore(k.key)
}
