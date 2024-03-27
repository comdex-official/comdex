package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/lend/types"
)

func (k Keeper) GetUtilisationRatioByPoolIDAndAssetID(ctx sdk.Context, poolID, assetID uint64) (sdk.Dec, error) {
	pool, _ := k.GetPool(ctx, poolID)
	asset, _ := k.Asset.GetAsset(ctx, assetID)
	moduleBalance := k.ModuleBalance(ctx, pool.ModuleName, asset.Denom)
	assetStats, found := k.GetAssetStatsByPoolIDAndAssetID(ctx, poolID, assetID)
	if !found {
		return sdk.ZeroDec(), types.ErrAssetStatsNotFound
	}
	if sdk.NewDecFromInt(moduleBalance).Add(sdk.NewDecFromInt(assetStats.TotalBorrowed.Add(assetStats.TotalStableBorrowed))).IsZero() {
		return sdk.ZeroDec(), nil
	}
	utilizationRatio := (sdk.NewDecFromInt(assetStats.TotalBorrowed.Add(assetStats.TotalStableBorrowed))).Quo(sdk.NewDecFromInt(moduleBalance).Add(sdk.NewDecFromInt(assetStats.TotalBorrowed.Add(assetStats.TotalStableBorrowed))))
	return utilizationRatio, nil
}

func (k Keeper) GetBorrowAPRByAssetID(ctx sdk.Context, poolID, assetID uint64, IsStableBorrow bool) (borrowAPY sdk.Dec, err error) {
	assetRatesStats, found := k.GetAssetRatesParams(ctx, assetID)
	if !found {
		return sdk.ZeroDec(), types.ErrorAssetStatsNotFound
	}
	currentUtilisationRatio, err := k.GetUtilisationRatioByPoolIDAndAssetID(ctx, poolID, assetID)
	if err != nil {
		return sdk.ZeroDec(), err
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
		utilisationDenominator := sdk.OneDec().Sub(assetRatesStats.UOptimal)
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
	utilisationDenominator := sdk.OneDec().Sub(assetRatesStats.UOptimal)
	utilisationRatio := utilisationNumerator.Quo(utilisationDenominator)
	multiplicationFactor := utilisationRatio.Mul(assetRatesStats.StableSlope2)
	borrowAPY = assetRatesStats.StableBase.Add(assetRatesStats.StableSlope1).Add(multiplicationFactor)
	return borrowAPY, nil
}

func (k Keeper) GetLendAPRByAssetIDAndPoolID(ctx sdk.Context, poolID, assetID uint64) (lendAPY sdk.Dec, err error) {
	assetRatesStats, found := k.GetAssetRatesParams(ctx, assetID)
	if !found {
		return sdk.ZeroDec(), types.ErrorAssetStatsNotFound
	}
	borrowAPY, err := k.GetBorrowAPRByAssetID(ctx, poolID, assetID, false)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	currentUtilisationRatio, err := k.GetUtilisationRatioByPoolIDAndAssetID(ctx, poolID, assetID)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	mulFactor := sdk.OneDec().Sub(assetRatesStats.ReserveFactor)
	lendAPY = borrowAPY.Mul(currentUtilisationRatio).Mul(mulFactor)

	return lendAPY, nil
}

func (k Keeper) GetAverageBorrowRate(ctx sdk.Context, poolID, assetID uint64) (sdk.Dec, error) {
	assetStats, _ := k.UpdateAPR(ctx, poolID, assetID)
	factor1 := assetStats.BorrowApr.Mul(sdk.NewDecFromInt(assetStats.TotalBorrowed))
	factor2 := assetStats.StableBorrowApr.Mul(sdk.NewDecFromInt(assetStats.TotalStableBorrowed))
	numerator := factor1.Add(factor2)
	denominator := sdk.NewDecFromInt(assetStats.TotalStableBorrowed.Add(assetStats.TotalBorrowed))

	if denominator.LTE(sdk.ZeroDec()) {
		return sdk.ZeroDec(), types.ErrAverageBorrowRate
	}
	averageBorrowRate := numerator.Quo(denominator)
	return averageBorrowRate, nil
}

func (k Keeper) GetSavingRate(ctx sdk.Context, poolID, assetID uint64) (savingRate sdk.Dec, err error) {
	assetRatesStats, found := k.GetAssetRatesParams(ctx, assetID)
	if !found {
		return sdk.ZeroDec(), types.ErrorAssetRatesParamsNotFound
	}
	averageBorrowRate, err := k.GetAverageBorrowRate(ctx, poolID, assetID)
	if err != nil {
		return sdk.Dec{}, err
	}
	utilizationRatio, err := k.GetUtilisationRatioByPoolIDAndAssetID(ctx, poolID, assetID)
	if err != nil {
		return sdk.Dec{}, err
	}
	factor1 := sdk.OneDec().Sub(assetRatesStats.ReserveFactor)
	savingRate = averageBorrowRate.Mul(utilizationRatio).Mul(factor1)
	return savingRate, nil
}

func (k Keeper) GetReserveRate(ctx sdk.Context, poolID, assetID uint64) (reserveRate sdk.Dec, err error) {
	averageBorrowRate, err := k.GetAverageBorrowRate(ctx, poolID, assetID)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	savingRate, err := k.GetSavingRate(ctx, poolID, assetID)
	if err != nil {
		return sdk.Dec{}, err
	}
	if averageBorrowRate != sdk.ZeroDec() {
		reserveRate = averageBorrowRate.Sub(savingRate)
		return reserveRate, nil
	}
	return sdk.ZeroDec(), nil
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
