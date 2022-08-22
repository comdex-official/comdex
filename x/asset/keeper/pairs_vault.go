package keeper

import (
	"github.com/comdex-official/comdex/app/wasm/bindings"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"

	"github.com/comdex-official/comdex/x/asset/types"
)

func (k Keeper) GetPairsVaultID(ctx sdk.Context) uint64 {
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

func (k Keeper) SetPairsVaultID(ctx sdk.Context, id uint64) {
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

func (k Keeper) SetPairsVault(ctx sdk.Context, app types.ExtendedPairVault) {
	var (
		store = k.Store(ctx)
		key   = types.PairsKey(app.Id)
		value = k.cdc.MustMarshal(&app)
	)

	store.Set(key, value)
}

func (k Keeper) GetPairsVault(ctx sdk.Context, id uint64) (pairs types.ExtendedPairVault, found bool) {
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

func (k Keeper) GetPairsVaults(ctx sdk.Context) (apps []types.ExtendedPairVault, found bool) {
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

func (k Keeper) WasmExtendedPairByAppQuery(ctx sdk.Context, appID uint64) (extId []uint64, found bool) {
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
		var extPair types.ExtendedPairVault
		k.cdc.MustUnmarshal(iter.Value(), &extPair)
		if extPair.AppId == appID {
			extId = append(extId, extPair.Id)
		}
	}
	if extId == nil {
		return nil, false
	}

	return extId, true
}

func (k Keeper) WasmAddExtendedPairsVaultRecords(ctx sdk.Context, pairVaultBinding *bindings.MsgAddExtendedPairsVault) error {
	DebtCeiling := sdk.NewInt(int64(pairVaultBinding.DebtCeiling))
	DebtFloor := sdk.NewInt(int64(pairVaultBinding.DebtFloor))

	_, found := k.GetApp(ctx, pairVaultBinding.AppID)
	if !found {
		return types.ErrorUnknownAppType
	}
	_, pairExists := k.GetPair(ctx, pairVaultBinding.PairID)
	if !pairExists {
		return types.ErrorPairDoesNotExist
	}

	var id = k.GetPairsVaultID(ctx)

	extendedPairVault, _ := k.GetPairsVaults(ctx)

	if len(extendedPairVault) > 0 {
		for _, data := range extendedPairVault {
			if (data.PairName == pairVaultBinding.PairName) && (data.AppId == pairVaultBinding.AppID) {
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
	blockHeight := ctx.BlockHeight()

	if pairVaultBinding.StabilityFee.IsZero() {
		blockHeight = 0
	}
	var app = types.ExtendedPairVault{
		Id:                  id + 1,
		AppId:               pairVaultBinding.AppID,
		PairId:              pairVaultBinding.PairID,
		StabilityFee:        pairVaultBinding.StabilityFee,
		ClosingFee:          pairVaultBinding.ClosingFee,
		LiquidationPenalty:  pairVaultBinding.LiquidationPenalty,
		DrawDownFee:         pairVaultBinding.DrawDownFee,
		IsVaultActive:       pairVaultBinding.IsVaultActive,
		DebtCeiling:         DebtCeiling,
		DebtFloor:           DebtFloor,
		IsStableMintVault:   pairVaultBinding.IsStableMintVault,
		MinCr:               pairVaultBinding.MinCr,
		PairName:            pairVaultBinding.PairName,
		AssetOutOraclePrice: pairVaultBinding.AssetOutOraclePrice,
		AssetOutPrice:       pairVaultBinding.AssetOutPrice,
		MinUsdValueLeft:     pairVaultBinding.MinUsdValueLeft,
		BlockHeight:         blockHeight,
		BlockTime:           ctx.BlockTime(),
	}

	k.SetPairsVaultID(ctx, app.Id)
	k.SetPairsVault(ctx, app)

	return nil
}

func (k Keeper) WasmAddExtendedPairsVaultRecordsQuery(ctx sdk.Context, appID, pairID uint64, StabilityFee, ClosingFee, DrawDownFee sdk.Dec, debtCeiling, debtFloor uint64, PairName string) (bool, string) {
	DebtCeiling := sdk.NewInt(int64(debtCeiling))
	DebtFloor := sdk.NewInt(int64(debtFloor))

	_, found := k.GetApp(ctx, appID)
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
			if (data.PairName == PairName) && (data.AppId == appID) {
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

func (k Keeper) WasmUpdatePairsVault(ctx sdk.Context, updatePairVault *bindings.MsgUpdatePairsVault) error {
	ExtPairVaultData, found := k.GetPairsVault(ctx, updatePairVault.ExtPairID)
	if !found {
		return types.ErrorPairDoesNotExist
	}
	_, found1 := k.rewards.GetAppIDByApp(ctx, updatePairVault.AppID)
	if found1 {
		if ExtPairVaultData.StabilityFee != updatePairVault.StabilityFee {
			if updatePairVault.StabilityFee.IsZero() {
				// run script to distrubyte reward
				k.VaultIterateRewards(ctx, ExtPairVaultData.StabilityFee, ExtPairVaultData.BlockHeight, ExtPairVaultData.BlockTime.Unix(), updatePairVault.AppID, ExtPairVaultData.Id, false)
				ExtPairVaultData.BlockTime = ctx.BlockTime()
				ExtPairVaultData.BlockHeight = 0

			} else if ExtPairVaultData.StabilityFee.IsZero() {
				// do nothing
				ExtPairVaultData.BlockHeight = ctx.BlockHeight()
				ExtPairVaultData.BlockTime = ctx.BlockTime()
			} else if ExtPairVaultData.StabilityFee.GT(sdk.ZeroDec()) && updatePairVault.StabilityFee.GT(sdk.ZeroDec()) {
				// run script to distribute
				k.VaultIterateRewards(ctx, ExtPairVaultData.StabilityFee, ExtPairVaultData.BlockHeight, ExtPairVaultData.BlockTime.Unix(), updatePairVault.AppID, ExtPairVaultData.Id, true)
				ExtPairVaultData.BlockHeight = ctx.BlockHeight()
				ExtPairVaultData.BlockTime = ctx.BlockTime()

			}
		}
	}

	ExtPairVaultData.StabilityFee = updatePairVault.StabilityFee
	ExtPairVaultData.ClosingFee = updatePairVault.ClosingFee
	ExtPairVaultData.LiquidationPenalty = updatePairVault.LiquidationPenalty
	ExtPairVaultData.DrawDownFee = updatePairVault.DrawDownFee
	ExtPairVaultData.DebtCeiling = sdk.NewInt(int64(updatePairVault.DebtCeiling))
	ExtPairVaultData.DebtFloor = sdk.NewInt(int64(updatePairVault.DebtFloor))
	ExtPairVaultData.MinCr = updatePairVault.MinCr
	ExtPairVaultData.MinUsdValueLeft = updatePairVault.MinUsdValueLeft

	k.SetPairsVault(ctx, ExtPairVaultData)

	return nil
}

func (k Keeper) WasmUpdatePairsVaultQuery(ctx sdk.Context, appID, exPairID uint64) (bool, string) {
	pairVaults, found := k.GetPairsVaults(ctx)
	if !found {
		return false, types.ErrorPairDoesNotExist.Error()
	}
	var count = 0
	for _, data := range pairVaults {
		if data.AppId == appID && data.Id == exPairID {
			count++
		}
	}
	if count == 0 {
		return false, types.ErrorExtendedPairDoesNotExistForTheApp.Error()
	}
	return true, ""
}

func (k Keeper) WasmCheckWhitelistedAssetQuery(ctx sdk.Context, denom string) (found bool) {
	found = k.HasAssetForDenom(ctx, denom)
	return found
}

func (k Keeper) VaultIterateRewards(ctx sdk.Context, collectorLsr sdk.Dec, collectorBh, collectorBt int64, appID, pairVaultID uint64, changeTypes bool) {

	extPairVault, found := k.vault.GetAppExtendedPairVaultMappingData(ctx, appID, pairVaultID)
	if found {
		for _, valID := range extPairVault.VaultIds {
			vaultData, found := k.vault.GetVault(ctx, valID)
			if !found {
				continue
			}
			interest := sdk.ZeroDec()
			var err error
			if vaultData.BlockHeight == 0 {
				interest, err = k.rewards.CalculationOfRewards(ctx, vaultData.AmountOut, collectorLsr, collectorBt)
				if err != nil {
					return
				}
			} else {
				interest, err = k.rewards.CalculationOfRewards(ctx, vaultData.AmountOut, collectorLsr, vaultData.BlockTime.Unix())
				if err != nil {
					return
				}
			}

			vaultInterestTracker, found := k.rewards.GetVaultInterestTracker(ctx, extPairVault.ExtendedPairId, appID)
			if !found {
				vaultInterestTracker = rewardstypes.VaultInterestTracker{
					VaultId:             extPairVault.ExtendedPairId,
					AppMappingId:        appID,
					InterestAccumulated: sdk.ZeroDec(),
				}
			}
			vaultInterestTracker.InterestAccumulated = vaultInterestTracker.InterestAccumulated.Add(interest)
			newInterest := sdk.ZeroInt()
			if vaultInterestTracker.InterestAccumulated.GTE(sdk.OneDec()) {
				newInterest = vaultInterestTracker.InterestAccumulated.TruncateInt()
				newInterestDec := sdk.NewDec(newInterest.Int64())
				vaultInterestTracker.InterestAccumulated = vaultInterestTracker.InterestAccumulated.Sub(newInterestDec)
			}
			k.rewards.SetVaultInterestTracker(ctx, vaultInterestTracker)

			// updating user rewards data
			vaultData.BlockTime = ctx.BlockTime()
			if changeTypes {
				vaultData.BlockHeight = ctx.BlockHeight()
			} else {
				vaultData.BlockHeight = 0
			}

			intAcc := vaultData.InterestAccumulated
			updatedIntAcc := (intAcc).Add(newInterest)
			vaultData.InterestAccumulated = updatedIntAcc
			k.vault.SetVault(ctx, vaultData)

		}
	}

}
