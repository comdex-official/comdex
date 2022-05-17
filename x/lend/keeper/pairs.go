package keeper

import (
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (k *Keeper) AddLendPairsRecords(ctx sdk.Context, records ...types.Extended_Pair) error {

	for _, msg := range records {

		_, found := k.GetLendPair(ctx, msg.Id)
		if found {
			return types.ErrorDuplicateLendPair
		}

		var (
			id   = k.GetLendPairID(ctx)
			pair = types.Extended_Pair{
				Id:                     id + 1,
				AssetIn:                msg.AssetIn,
				AssetOut:               msg.AssetOut,
				ModuleAcc:              msg.ModuleAcc,
				BaseBorrowRateAssetIn:  msg.BaseBorrowRateAssetIn,
				BaseLendRateAssetIn:    msg.BaseLendRateAssetIn,
				BaseBorrowRateAssetOut: msg.BaseBorrowRateAssetOut,
				BaseLendRateAssetOut:   msg.BaseLendRateAssetOut,
			}
		)

		k.SetLendPairID(ctx, pair.Id)
		k.SetLendPair(ctx, pair)
	}
	return nil
}

func (k *Keeper) UpdateLendPairRecords(ctx sdk.Context, msg types.Extended_Pair) error {
	pair, found := k.GetLendPair(ctx, msg.Id)
	if !found {
		return types.ErrorPairDoesNotExist
	}

	if len(msg.ModuleAcc) > 0 {
		pair.ModuleAcc = msg.ModuleAcc
	}
	if !msg.BaseBorrowRateAssetIn.IsZero() {
		pair.BaseBorrowRateAssetIn = msg.BaseBorrowRateAssetIn
	}
	if !msg.BaseBorrowRateAssetOut.IsZero() {
		pair.BaseBorrowRateAssetOut = msg.BaseBorrowRateAssetOut
	}
	if !msg.BaseLendRateAssetIn.IsZero() {
		pair.BaseLendRateAssetIn = msg.BaseLendRateAssetIn
	}
	if !msg.BaseLendRateAssetOut.IsZero() {
		pair.BaseLendRateAssetOut = msg.BaseLendRateAssetOut
	}

	k.SetLendPair(ctx, pair)
	return nil
}

func (k *Keeper) SetLendPairID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LendPairIDKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) SetLendPair(ctx sdk.Context, pair types.Extended_Pair) {
	var (
		store = k.Store(ctx)
		key   = types.LendPairKey(pair.Id)
		value = k.cdc.MustMarshal(&pair)
	)

	store.Set(key, value)
}

func (k *Keeper) GetLendPair(ctx sdk.Context, id uint64) (pair types.Extended_Pair, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LendPairKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return pair, false
	}

	k.cdc.MustUnmarshal(value, &pair)
	return pair, true
}

func (k *Keeper) GetLendPairID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.LendPairIDKey
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var count protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &count)

	return count.GetValue()
}
