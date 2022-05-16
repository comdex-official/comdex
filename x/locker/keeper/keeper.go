package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/comdex-official/comdex/x/locker/expected"
	"github.com/comdex-official/comdex/x/locker/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type Keeper struct {
	cdc    codec.BinaryCodec
	key    sdk.StoreKey
	paramstore paramtypes.Subspace
	bank   expected.BankKeeper
	asset  expected.AssetKeeper
	oracle expected.OracleKeeper
}

func NewKeeper(cdc codec.BinaryCodec, key sdk.StoreKey, ps paramtypes.Subspace, bank expected.BankKeeper, asset expected.AssetKeeper, oracle expected.OracleKeeper) Keeper {
	return Keeper{
		cdc:    cdc,
		key:    key,
		paramstore: ps,
		bank:   bank,
		asset:  asset,
		oracle: oracle,
	}
}

func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.key)
}
