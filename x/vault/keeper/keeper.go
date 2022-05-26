package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/comdex-official/comdex/x/vault/expected"
	"github.com/comdex-official/comdex/x/vault/types"
)

type Keeper struct {
	cdc       codec.BinaryCodec
	key       sdk.StoreKey
	bank      expected.BankKeeper
	asset     expected.AssetKeeper
	oracle    expected.Marketkeeper
	collector expected.CollectorKeeper
	esm       expected.EsmKeeper
}

//
//func (k *Keeper) CalculateCollaterlizationRatio(ctx sdk.Context, amountIn sdk.Int, assetIn assettypes.Asset, amountOut sdk.Int, assetOut assettypes.Asset) (sdk.Dec, error) {
//	//TODO implement me
//	panic("implement me")
//}

//func (k *Keeper) BurnCAssets(ctx sdk.Context, moduleName string, collateralDenom string, denom string, amount sdk.Int) error {
//	//TODO implement me
//	panic("implement me")
//}

func NewKeeper(cdc codec.BinaryCodec, key sdk.StoreKey, bank expected.BankKeeper, asset expected.AssetKeeper, oracle expected.Marketkeeper, collector expected.CollectorKeeper, esm expected.EsmKeeper) Keeper {
	return Keeper{
		cdc:       cdc,
		key:       key,
		bank:      bank,
		asset:     asset,
		oracle:    oracle,
		collector: collector,
		esm:       esm,
	}
}

func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.key)
}
