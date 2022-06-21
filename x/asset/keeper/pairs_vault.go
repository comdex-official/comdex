package keeper

import (
	"github.com/comdex-official/comdex/app/wasm/bindings"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"

	"github.com/comdex-official/comdex/x/asset/types"
)

func (k *Keeper) GetPairsVaultID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.PairsVaultIDKey
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
		key   = types.PairsVaultIDKey
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

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

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

func (k *Keeper) AddExtendedPairsVaultRecords(ctx sdk.Context, records ...types.ExtendedPairVault) error {
	for _, msg := range records {
		_, found := k.GetApp(ctx, msg.AppMappingId)
		if !found {
			return types.ErrorUnknownAppType
		}
		_, pairExists := k.GetPair(ctx, msg.PairId)
		if !pairExists {
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
			AssetOutPrice:       msg.AssetOutPrice,
			MinUsdValueLeft:     msg.MinUsdValueLeft,
		}

		k.SetPairsVaultID(ctx, app.Id)
		k.SetPairsVault(ctx, app)
	}
	return nil
}

func (k *Keeper) WasmAddExtendedPairsVaultRecords(ctx sdk.Context, pairVaultBinding *bindings.MsgAddExtendedPairsVault) error {
	DebtCeiling := sdk.NewInt(int64(pairVaultBinding.DebtCeiling))
	DebtFloor := sdk.NewInt(int64(pairVaultBinding.DebtFloor))

	_, found := k.GetApp(ctx, pairVaultBinding.AppMappingId)
	if !found {
		return types.ErrorUnknownAppType
	}
	_, pairExists := k.GetPair(ctx, pairVaultBinding.PairId)
	if !pairExists {
		return types.ErrorPairDoesNotExist
	}

	var id = k.GetPairsVaultID(ctx)

	extendedPairVault, _ := k.GetPairsVaults(ctx)

	if len(extendedPairVault) > 0 {
		for _, data := range extendedPairVault {
			if (data.PairName == pairVaultBinding.PairName) && (data.AppMappingId == pairVaultBinding.AppMappingId) {
				return types.ErrorPairNameForID
			}
		}
	}
	if DebtFloor.GTE(DebtCeiling) {
		return types.ErrorDebtFloorIsGreaterThanDebtCeiling
	}
	if !(pairVaultBinding.StabilityFee.GTE(sdk.ZeroDec()) && pairVaultBinding.StabilityFee.LT(sdk.OneDec())) {
		return types.ErrorFeeShouldNotBeGTOne
	}
	if !(pairVaultBinding.ClosingFee.GTE(sdk.ZeroDec()) && pairVaultBinding.ClosingFee.LT(sdk.OneDec())) {
		return types.ErrorFeeShouldNotBeGTOne
	}
	if !(pairVaultBinding.DrawDownFee.GTE(sdk.ZeroDec()) && pairVaultBinding.DrawDownFee.LT(sdk.OneDec())) {
		return types.ErrorFeeShouldNotBeGTOne
	}
	var app = types.ExtendedPairVault{
		Id:                  id + 1,
		AppMappingId:        pairVaultBinding.AppMappingId,
		PairId:              pairVaultBinding.PairId,
		LiquidationRatio:    pairVaultBinding.LiquidationRatio,
		StabilityFee:        pairVaultBinding.StabilityFee,
		ClosingFee:          pairVaultBinding.ClosingFee,
		LiquidationPenalty:  pairVaultBinding.LiquidationPenalty,
		DrawDownFee:         pairVaultBinding.DrawDownFee,
		IsVaultActive:       pairVaultBinding.IsVaultActive,
		DebtCeiling:         DebtCeiling,
		DebtFloor:           DebtFloor,
		IsPsmPair:           pairVaultBinding.IsPsmPair,
		MinCr:               pairVaultBinding.MinCr,
		PairName:            pairVaultBinding.PairName,
		AssetOutOraclePrice: pairVaultBinding.AssetOutOraclePrice,
		AssetOutPrice:       pairVaultBinding.AssetOutPrice,
		MinUsdValueLeft:     pairVaultBinding.MinUsdValueLeft,
	}

	k.SetPairsVaultID(ctx, app.Id)
	k.SetPairsVault(ctx, app)

	return nil
}

func (k *Keeper) WasmAddExtendedPairsVaultRecordsQuery(ctx sdk.Context, appMappingID, pairID uint64, StabilityFee, ClosingFee, DrawDownFee sdk.Dec, debtCeiling, debtFloor uint64, PairName string) (bool, string) {
	DebtCeiling := sdk.NewInt(int64(debtCeiling))
	DebtFloor := sdk.NewInt(int64(debtFloor))

	_, found := k.GetApp(ctx, appMappingID)
	if !found {
		return false, types.ErrorUnknownAppType.Error()
	}
	_, pairExists := k.GetPair(ctx, pairID)
	if !pairExists {
		return false, types.ErrorPairDoesNotExist.Error()
	}
	extendedPairVault, _ := k.GetPairsVaults(ctx)

	if len(extendedPairVault) > 0 {
		for _, data := range extendedPairVault {
			if (data.PairName == PairName) && (data.AppMappingId == appMappingID) {
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

func (k *Keeper) WasmUpdatePairsVault(ctx sdk.Context, updatePairVault *bindings.MsgUpdatePairsVault) error {
	var ExtPairVaultData types.ExtendedPairVault
	pairVaults, found := k.GetPairsVaults(ctx)
	if !found {
		return types.ErrorPairDoesNotExist
	}
	var count = 0
	for _, data := range pairVaults {
		if data.AppMappingId == updatePairVault.AppMappingId && data.Id == updatePairVault.ExtPairId {
			count++
			ExtPairVaultData.Id = data.Id
			ExtPairVaultData.PairId = data.PairId
			ExtPairVaultData.AppMappingId = data.AppMappingId
			ExtPairVaultData.LiquidationRatio = updatePairVault.LiquidationRatio
			ExtPairVaultData.StabilityFee = updatePairVault.StabilityFee
			ExtPairVaultData.ClosingFee = updatePairVault.ClosingFee
			ExtPairVaultData.LiquidationPenalty = updatePairVault.LiquidationPenalty
			ExtPairVaultData.DrawDownFee = updatePairVault.DrawDownFee
			ExtPairVaultData.IsVaultActive = data.IsVaultActive
			ExtPairVaultData.DebtCeiling = sdk.NewInt(int64(updatePairVault.DebtCeiling))
			ExtPairVaultData.DebtFloor = sdk.NewInt(int64(updatePairVault.DebtFloor))
			ExtPairVaultData.IsPsmPair = data.IsPsmPair
			ExtPairVaultData.MinCr = updatePairVault.MinCr
			ExtPairVaultData.PairName = data.PairName
			ExtPairVaultData.AssetOutOraclePrice = data.AssetOutOraclePrice
			ExtPairVaultData.AssetOutPrice = data.AssetOutPrice
			ExtPairVaultData.MinUsdValueLeft = updatePairVault.MinUsdValueLeft
		}
	}
	if count == 0 {
		return types.ErrorExtendedPairDoesNotExistForTheApp
	}

	k.SetPairsVault(ctx, ExtPairVaultData)

	return nil
}

func (k *Keeper) WasmUpdatePairsVaultQuery(ctx sdk.Context, appID, exPairID uint64) (bool, string) {
	pairVaults, found := k.GetPairsVaults(ctx)
	if !found {
		return false, types.ErrorPairDoesNotExist.Error()
	}
	var count = 0
	for _, data := range pairVaults {
		if data.AppMappingId == appID && data.Id == exPairID {
			count++
		}
	}
	if count == 0 {
		return false, types.ErrorExtendedPairDoesNotExistForTheApp.Error()
	}
	return true, ""
}
