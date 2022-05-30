package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"

	"github.com/comdex-official/comdex/x/asset/types"
)

func (k *Keeper) GetPairsVaultID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.PairsVaultIDkey
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetPairsVaultID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.PairsVaultIDkey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k *Keeper) SetPairsVault(ctx sdk.Context, app types.ExtendedPairVault) {
	var (
		store = k.Store(ctx)
		key   = types.PairsKey(app.Id)
		value = k.cdc.MustMarshal(&app)
	)

	store.Set(key, value)
}

func (k *Keeper) GetPairsVault(ctx sdk.Context, id uint64) (pairs types.ExtendedPairVault, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.PairsKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return pairs, false
	}

	k.cdc.MustUnmarshal(value, &pairs)
	return pairs, true
}

func (k *Keeper) GetPairsVaults(ctx sdk.Context) (apps []types.ExtendedPairVault, found bool) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.PairsVaultKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var app types.ExtendedPairVault
		k.cdc.MustUnmarshal(iter.Value(), &app)
		apps = append(apps, app)
	}
	if apps == nil {
		return nil, false
	}

	return apps, true
}

// func (k *Keeper) SetPairsVaultForPairId(ctx sdk.Context, pairId uint64, id uint64) {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.PairsForPairIdKey(pairId)
// 		value = k.cdc.MustMarshal(
// 			&protobuftypes.UInt64Value{
// 				Value: id,
// 			},
// 		)
// 	)

// 	store.Set(key, value)
// }

// // checks if extended pair exists for a given asset pair ID
// func (k *Keeper) HasPairsVaultForPairId(ctx sdk.Context, PairId uint64) bool {
// 	var (
// 		store = k.Store(ctx)
// 		key   = types.PairsForPairIdKey(PairId)
// 	)

// 	return store.Has(key)
// }

func (k *Keeper) AddExtendedPairsVaultRecords(ctx sdk.Context, records ...types.ExtendedPairVault) error {

	for _, msg := range records {

		_, found := k.GetApp(ctx, msg.AppMappingId)
		if !found {
			return types.ErrorUnknownAppType
		}
		_, gotit := k.GetPair(ctx, msg.PairId)
		if !gotit{
			return types.ErrorPairDoesNotExist
		}

			var id    = k.GetPairsVaultID(ctx)

			extendedPairVault, _ := k.GetPairsVaults(ctx)

			if len(extendedPairVault) > 0 {
			for _, data := range extendedPairVault{
				if (data.PairName == msg.PairName) && (data.AppMappingId == msg.AppMappingId) {
						return types.ErrorPairNameForID
					}
				}
			}
			if msg.DebtFloor.GTE(msg.DebtCeiling) {
				return types.ErrorDebtFloorIsGreaterThanDebtCeiling
			}
			if !(msg.StabilityFee.GTE(sdk.ZeroDec()) && msg.StabilityFee.LT(sdk.OneDec())) {
				return types.ErrorFeeShouldNotBeGTOne
			}
			if !(msg.ClosingFee.GTE(sdk.ZeroDec()) && msg.ClosingFee.LT(sdk.OneDec())) {
				return types.ErrorFeeShouldNotBeGTOne
			}
			if !(msg.DrawDownFee.GTE(sdk.ZeroDec()) && msg.DrawDownFee.LT(sdk.OneDec())) {
				return types.ErrorFeeShouldNotBeGTOne
			}
			var app = types.ExtendedPairVault{
				Id:       id + 1,
				AppMappingId: msg.AppMappingId,
				PairId: msg.PairId,
				LiquidationRatio: msg.LiquidationRatio,
				StabilityFee: msg.StabilityFee,
				ClosingFee: msg.ClosingFee,
				LiquidationPenalty: msg.LiquidationPenalty,
				DrawDownFee: msg.DrawDownFee,
				IsVaultActive: msg.IsVaultActive,
				DebtCeiling: msg.DebtCeiling,
				DebtFloor: msg.DebtFloor,
				IsPsmPair: msg.IsPsmPair,
				MinCr: msg.MinCr,
				PairName: msg.PairName,
				AssetOutOraclePrice: msg.AssetOutOraclePrice,
				AsssetOutPrice: msg.AsssetOutPrice,
			}

		k.SetPairsVaultID(ctx, app.Id)
		k.SetPairsVault(ctx, app)
	}
	return nil
}