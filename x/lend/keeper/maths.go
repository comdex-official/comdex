package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/lend/types"
)

func (k Keeper) GetUtilisationRatioByPoolIDAndAssetID(ctx sdk.Context, poolID, assetID uint64) (sdkmath.LegacyDec, error) {
	pool, _ := k.GetPool(ctx, poolID)
	asset, _ := k.Asset.GetAsset(ctx, assetID)
	moduleBalance := k.ModuleBalance(ctx, pool.ModuleName, asset.Denom)
	assetStats, found := k.GetAssetStatsByPoolIDAndAssetID(ctx, poolID, assetID)
	if !found {
		return sdkmath.LegacyZeroDec(), types.ErrAssetStatsNotFound
	}
	if sdkmath.LegacyNewDecFromInt(moduleBalance).Add(sdkmath.LegacyNewDecFromInt(assetStats.TotalBorrowed.Add(assetStats.TotalStableBorrowed))).IsZero() {
		return sdkmath.LegacyZeroDec(), nil
	}
	utilizationRatio := (sdkmath.LegacyNewDecFromInt(assetStats.TotalBorrowed.Add(assetStats.TotalStableBorrowed))).Quo(sdkmath.LegacyNewDecFromInt(moduleBalance).Add(sdkmath.LegacyNewDecFromInt(assetStats.TotalBorrowed.Add(assetStats.TotalStableBorrowed))))
	return utilizationRatio, nil
}

func (k Keeper) GetBorrowAPRByAssetID(ctx sdk.Context, poolID, assetID uint64, IsStableBorrow bool) (borrowAPY sdkmath.LegacyDec, err error) {
	assetRatesStats, found := k.GetAssetRatesParams(ctx, assetID)
	if !found {
		return sdkmath.LegacyZeroDec(), types.ErrorAssetStatsNotFound
	}
	currentUtilisationRatio, err := k.GetUtilisationRatioByPoolIDAndAssetID(ctx, poolID, assetID)
	if err != nil {
		return sdkmath.LegacyZeroDec(), err
	}
	// for normal borrow
	if !IsStableBorrow {
		if currentUtilisationRatio.LT(assetRatesStats.UOptimal) {
			utilisationRatio := currentUtilisationRatio.Quo(assetRatesStats.UOptimal)
			multiplicationFactor := utilisationRatio.Mul(assetRatesStats.Slope1)
			borrowAPY = assetRatesStats.Base.Add(multiplicationFactor)
			return borrowAPY, nil
		}
		utilisationNumerator := currentUtilisationRatio.Sub(assetRatesStats.UOptimal)
		utilisationDenominator := sdkmath.LegacyOneDec().Sub(assetRatesStats.UOptimal)
		utilisationRatio := utilisationNumerator.Quo(utilisationDenominator)
		multiplicationFactor := utilisationRatio.Mul(assetRatesStats.Slope2)
		borrowAPY = assetRatesStats.Base.Add(assetRatesStats.Slope1).Add(multiplicationFactor)
		return borrowAPY, nil
	} // for stable borrow
	if currentUtilisationRatio.LT(assetRatesStats.UOptimal) {
		utilisationRatio := currentUtilisationRatio.Quo(assetRatesStats.UOptimal)
		multiplicationFactor := utilisationRatio.Mul(assetRatesStats.StableSlope1)
		borrowAPY = assetRatesStats.StableBase.Add(multiplicationFactor)
		return borrowAPY, nil
	}
	utilisationNumerator := currentUtilisationRatio.Sub(assetRatesStats.UOptimal)
	utilisationDenominator := sdkmath.LegacyOneDec().Sub(assetRatesStats.UOptimal)
	utilisationRatio := utilisationNumerator.Quo(utilisationDenominator)
	multiplicationFactor := utilisationRatio.Mul(assetRatesStats.StableSlope2)
	borrowAPY = assetRatesStats.StableBase.Add(assetRatesStats.StableSlope1).Add(multiplicationFactor)
	return borrowAPY, nil
}

func (k Keeper) GetLendAPRByAssetIDAndPoolID(ctx sdk.Context, poolID, assetID uint64) (lendAPY sdkmath.LegacyDec, err error) {
	assetRatesStats, found := k.GetAssetRatesParams(ctx, assetID)
	if !found {
		return sdkmath.LegacyZeroDec(), types.ErrorAssetStatsNotFound
	}
	borrowAPY, err := k.GetBorrowAPRByAssetID(ctx, poolID, assetID, false)
	if err != nil {
		return sdkmath.LegacyZeroDec(), err
	}
	currentUtilisationRatio, err := k.GetUtilisationRatioByPoolIDAndAssetID(ctx, poolID, assetID)
	if err != nil {
		return sdkmath.LegacyZeroDec(), err
	}
	mulFactor := sdkmath.LegacyOneDec().Sub(assetRatesStats.ReserveFactor)
	lendAPY = borrowAPY.Mul(currentUtilisationRatio).Mul(mulFactor)

	return lendAPY, nil
}

func (k Keeper) GetAverageBorrowRate(ctx sdk.Context, poolID, assetID uint64) (sdkmath.LegacyDec, error) {
	assetStats, _ := k.UpdateAPR(ctx, poolID, assetID)
	factor1 := assetStats.BorrowApr.Mul(sdkmath.LegacyNewDecFromInt(assetStats.TotalBorrowed))
	factor2 := assetStats.StableBorrowApr.Mul(sdkmath.LegacyNewDecFromInt(assetStats.TotalStableBorrowed))
	numerator := factor1.Add(factor2)
	denominator := sdkmath.LegacyNewDecFromInt(assetStats.TotalStableBorrowed.Add(assetStats.TotalBorrowed))

	if denominator.LTE(sdkmath.LegacyZeroDec()) {
		return sdkmath.LegacyZeroDec(), types.ErrAverageBorrowRate
	}
	averageBorrowRate := numerator.Quo(denominator)
	return averageBorrowRate, nil
}

func (k Keeper) GetSavingRate(ctx sdk.Context, poolID, assetID uint64) (savingRate sdkmath.LegacyDec, err error) {
	assetRatesStats, found := k.GetAssetRatesParams(ctx, assetID)
	if !found {
		return sdkmath.LegacyZeroDec(), types.ErrorAssetRatesParamsNotFound
	}
	averageBorrowRate, err := k.GetAverageBorrowRate(ctx, poolID, assetID)
	if err != nil {
		return sdkmath.LegacyDec{}, err
	}
	utilizationRatio, err := k.GetUtilisationRatioByPoolIDAndAssetID(ctx, poolID, assetID)
	if err != nil {
		return sdkmath.LegacyDec{}, err
	}
	factor1 := sdkmath.LegacyOneDec().Sub(assetRatesStats.ReserveFactor)
	savingRate = averageBorrowRate.Mul(utilizationRatio).Mul(factor1)
	return savingRate, nil
}

func (k Keeper) GetReserveRate(ctx sdk.Context, poolID, assetID uint64) (reserveRate sdkmath.LegacyDec, err error) {
	averageBorrowRate, err := k.GetAverageBorrowRate(ctx, poolID, assetID)
	if err != nil {
		return sdkmath.LegacyZeroDec(), err
	}
	savingRate, err := k.GetSavingRate(ctx, poolID, assetID)
	if err != nil {
		return sdkmath.LegacyDec{}, err
	}
	if averageBorrowRate != sdkmath.LegacyZeroDec() {
		reserveRate = averageBorrowRate.Sub(savingRate)
		return reserveRate, nil
	}
	return sdkmath.LegacyZeroDec(), nil
}

func (k Keeper) UpdateAPR(ctx sdk.Context, poolID, assetID uint64) (PoolAssetLBData types.PoolAssetLBMapping, found bool) {
	poolAssetLBData, found := k.GetAssetStatsByPoolIDAndAssetID(ctx, poolID, assetID)
	if !found {
		return poolAssetLBData, false
	}
	lendAPR, _ := k.GetLendAPRByAssetIDAndPoolID(ctx, poolID, assetID)
	borrowAPR, _ := k.GetBorrowAPRByAssetID(ctx, poolID, assetID, false)
	stableBorrowAPR, _ := k.GetBorrowAPRByAssetID(ctx, poolID, assetID, true)
	currentUtilisationRatio, _ := k.GetUtilisationRatioByPoolIDAndAssetID(ctx, poolID, assetID)
	PoolAssetLBData = types.PoolAssetLBMapping{
		PoolID:                   poolAssetLBData.PoolID,
		AssetID:                  poolAssetLBData.AssetID,
		LendIds:                  poolAssetLBData.LendIds,
		BorrowIds:                poolAssetLBData.BorrowIds,
		TotalBorrowed:            poolAssetLBData.TotalBorrowed,
		TotalStableBorrowed:      poolAssetLBData.TotalStableBorrowed,
		TotalLend:                poolAssetLBData.TotalLend,
		TotalInterestAccumulated: poolAssetLBData.TotalInterestAccumulated,
		LendApr:                  lendAPR,
		BorrowApr:                borrowAPR,
		StableBorrowApr:          stableBorrowAPR,
		UtilisationRatio:         currentUtilisationRatio,
	}
	return PoolAssetLBData, true
}
