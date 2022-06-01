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
		if !gotit {
			return types.ErrorPairDoesNotExist
		}

		var id = k.GetPairsVaultID(ctx)

		extendedPairVault, _ := k.GetPairsVaults(ctx)

		if len(extendedPairVault) > 0 {
			for _, data := range extendedPairVault {
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
			Id:                  id + 1,
			AppMappingId:        msg.AppMappingId,
			PairId:              msg.PairId,
			LiquidationRatio:    msg.LiquidationRatio,
			StabilityFee:        msg.StabilityFee,
			ClosingFee:          msg.ClosingFee,
			LiquidationPenalty:  msg.LiquidationPenalty,
			DrawDownFee:         msg.DrawDownFee,
			IsVaultActive:       msg.IsVaultActive,
			DebtCeiling:         msg.DebtCeiling,
			DebtFloor:           msg.DebtFloor,
			IsPsmPair:           msg.IsPsmPair,
			MinCr:               msg.MinCr,
			PairName:            msg.PairName,
			AssetOutOraclePrice: msg.AssetOutOraclePrice,
			AsssetOutPrice:      msg.AsssetOutPrice,
		}

		k.SetPairsVaultID(ctx, app.Id)
		k.SetPairsVault(ctx, app)
	}
	return nil
}

func (k *Keeper) WasmAddExtendedPairsVaultRecords(ctx sdk.Context, AppMappingId, PairId uint64, LiquidationRatio, StabilityFee, ClosingFee, LiquidationPenalty, DrawDownFee sdk.Dec, IsVaultActive bool, debtCeiling, debtFloor uint64, IsPsmPair bool, MinCr sdk.Dec, PairName string, AssetOutOraclePrice bool, AssetOutPrice uint64) error {

	DebtCeiling := sdk.NewInt(int64(debtCeiling))
	DebtFloor := sdk.NewInt(int64(debtFloor))

	_, found := k.GetApp(ctx, AppMappingId)
	if !found {
		return types.ErrorUnknownAppType
	}
	_, gotit := k.GetPair(ctx, PairId)
	if !gotit {
		return types.ErrorPairDoesNotExist
	}

	var id = k.GetPairsVaultID(ctx)

	extendedPairVault, _ := k.GetPairsVaults(ctx)

	if len(extendedPairVault) > 0 {
		for _, data := range extendedPairVault {
			if (data.PairName == PairName) && (data.AppMappingId == AppMappingId) {
				return types.ErrorPairNameForID
			}
		}
	}
	if DebtFloor.GTE(DebtCeiling) {
		return types.ErrorDebtFloorIsGreaterThanDebtCeiling
	}
	if !(StabilityFee.GTE(sdk.ZeroDec()) && StabilityFee.LT(sdk.OneDec())) {
		return types.ErrorFeeShouldNotBeGTOne
	}
	if !(ClosingFee.GTE(sdk.ZeroDec()) && ClosingFee.LT(sdk.OneDec())) {
		return types.ErrorFeeShouldNotBeGTOne
	}
	if !(DrawDownFee.GTE(sdk.ZeroDec()) && DrawDownFee.LT(sdk.OneDec())) {
		return types.ErrorFeeShouldNotBeGTOne
	}
	var app = types.ExtendedPairVault{
		Id:                  id + 1,
		AppMappingId:        AppMappingId,
		PairId:              PairId,
		LiquidationRatio:    LiquidationRatio,
		StabilityFee:        StabilityFee,
		ClosingFee:          ClosingFee,
		LiquidationPenalty:  LiquidationPenalty,
		DrawDownFee:         DrawDownFee,
		IsVaultActive:       IsVaultActive,
		DebtCeiling:         DebtCeiling,
		DebtFloor:           DebtFloor,
		IsPsmPair:           IsPsmPair,
		MinCr:               MinCr,
		PairName:            PairName,
		AssetOutOraclePrice: AssetOutOraclePrice,
		AsssetOutPrice:      AssetOutPrice,
	}

	k.SetPairsVaultID(ctx, app.Id)
	k.SetPairsVault(ctx, app)

	return nil
}

func (k *Keeper) WasmAddExtendedPairsVaultRecordsQuery(ctx sdk.Context, AppMappingId, PairId uint64, StabilityFee, ClosingFee, DrawDownFee sdk.Dec, debtCeiling, debtFloor uint64, PairName string) (bool, string) {

	DebtCeiling := sdk.NewInt(int64(debtCeiling))
	DebtFloor := sdk.NewInt(int64(debtFloor))

	_, found := k.GetApp(ctx, AppMappingId)
	if !found {
		return false, types.ErrorUnknownAppType.Error()
	}
	_, gotit := k.GetPair(ctx, PairId)
	if !gotit {
		return false, types.ErrorPairDoesNotExist.Error()
	}
	extendedPairVault, _ := k.GetPairsVaults(ctx)

	if len(extendedPairVault) > 0 {
		for _, data := range extendedPairVault {
			if (data.PairName == PairName) && (data.AppMappingId == AppMappingId) {
				return false, types.ErrorPairNameForID.Error()
			}
		}
	}
	if DebtFloor.GTE(DebtCeiling) {
		return false, types.ErrorDebtFloorIsGreaterThanDebtCeiling.Error()
	}
	if !(StabilityFee.GTE(sdk.ZeroDec()) && StabilityFee.LT(sdk.OneDec())) {
		return false, types.ErrorFeeShouldNotBeGTOne.Error()
	}
	if !(ClosingFee.GTE(sdk.ZeroDec()) && ClosingFee.LT(sdk.OneDec())) {
		return false, types.ErrorFeeShouldNotBeGTOne.Error()
	}
	if !(DrawDownFee.GTE(sdk.ZeroDec()) && DrawDownFee.LT(sdk.OneDec())) {
		return false, types.ErrorFeeShouldNotBeGTOne.Error()
	}

	return true, ""
}

func (k *Keeper) WasmUpdateLsrInPairsVault(ctx sdk.Context, app_id, ex_pair_id uint64, liq_ratio, stab_fee, close_fee, penalty, 
	draw_down_fee, min_cr sdk.Dec, debtCeiling, debtFloor sdk.Int) error {

	var ExtPairVaultData types.ExtendedPairVault
	pairVaults, found := k.GetPairsVaults(ctx)
	if !found {
		return types.ErrorPairDoesNotExist
	}
	var count = 0
	for _, data := range pairVaults {
		if data.AppMappingId == app_id && data.Id == ex_pair_id {
			count++
			ExtPairVaultData.Id = data.Id
			ExtPairVaultData.PairId = data.PairId
			ExtPairVaultData.AppMappingId = data.AppMappingId
			ExtPairVaultData.LiquidationRatio = liq_ratio
			ExtPairVaultData.StabilityFee = stab_fee
			ExtPairVaultData.ClosingFee = close_fee
			ExtPairVaultData.LiquidationPenalty = penalty
			ExtPairVaultData.DrawDownFee = draw_down_fee
			ExtPairVaultData.IsVaultActive = data.IsVaultActive
			ExtPairVaultData.DebtCeiling = debtCeiling
			ExtPairVaultData.DebtFloor = debtFloor
			ExtPairVaultData.IsPsmPair = data.IsPsmPair
			ExtPairVaultData.MinCr = min_cr
			ExtPairVaultData.PairName = data.PairName
			ExtPairVaultData.AssetOutOraclePrice = data.AssetOutOraclePrice
			ExtPairVaultData.AsssetOutPrice = data.AsssetOutPrice
		}
	}
	if count == 0 {
		return types.ErrorExtendedPairDoesNotExistForTheApp
	}

	var (
		store = k.Store(ctx)
		key   = types.PairsKey(app_id)
		value = k.cdc.MustMarshal(&ExtPairVaultData)
	)

	store.Set(key, value)
	return nil
}

func (k *Keeper) WasmUpdateLsrInPairsVaultQuery(ctx sdk.Context, appId, exPairId uint64) (bool, string) {
	pairVaults, found := k.GetPairsVaults(ctx)
	if !found {
		return false, types.ErrorPairDoesNotExist.Error()
	}
	var count = 0
	for _, data := range pairVaults {
		if data.AppMappingId == appId && data.Id == exPairId {
			count++
		}
	}
	if count == 0 {
		return false, types.ErrorExtendedPairDoesNotExistForTheApp.Error()
	}
	return true, ""
}
