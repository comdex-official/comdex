package keeper

import (
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetUtilisationRatioByPoolIdAndAssetId(ctx sdk.Context, poolId, assetId uint64) (sdk.Dec, error) {
	pool, _ := k.GetPool(ctx, poolId)
	asset, _ := k.GetAsset(ctx, assetId)
	moduleBalance := k.ModuleBalance(ctx, pool.ModuleName, asset.Denom)
	assetStats, found := k.GetAssetStatsByPoolIdAndAssetId(ctx, assetId, poolId)
	if !found {
		return sdk.ZeroDec(), types.ErrAssetStatsNotFound
	}
	utilizationRatio := assetStats.TotalBorrowed.ToDec().Quo(moduleBalance.ToDec())
	return utilizationRatio, nil
}

func (k Keeper) GetBorrowAPRByAssetId(ctx sdk.Context, poolId, assetId uint64, IsStableBorrow bool) (borrowAPY sdk.Dec, err error) {
	assetRatesStats, found := k.GetAssetRatesStats(ctx, assetId)
	if !found {
		return sdk.ZeroDec(), types.ErrorAssetStatsNotFound
	}
	currentUtilisationRatio, err := k.GetUtilisationRatioByPoolIdAndAssetId(ctx, poolId, assetId)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	if !IsStableBorrow {
		if currentUtilisationRatio.LT(assetRatesStats.UOptimal) {
			utilisationRatio := currentUtilisationRatio.Quo(assetRatesStats.UOptimal)
			multiplicationFactor := utilisationRatio.Mul(assetRatesStats.Slope1)
			borrowAPY = assetRatesStats.Base.Add(multiplicationFactor)
			return borrowAPY, nil
		} else {
			utilisationNumerator := currentUtilisationRatio.Sub(assetRatesStats.UOptimal)
			utilisationDenominator := sdk.OneDec().Sub(assetRatesStats.UOptimal)
			utilisationRatio := utilisationNumerator.Quo(utilisationDenominator)
			multiplicationFactor := utilisationRatio.Mul(assetRatesStats.Slope2)
			borrowAPY = assetRatesStats.Base.Add(assetRatesStats.Slope1).Add(multiplicationFactor)
			return borrowAPY, nil
		}
	} else {
		if currentUtilisationRatio.LT(assetRatesStats.UOptimal) {
			utilisationRatio := currentUtilisationRatio.Quo(assetRatesStats.UOptimal)
			multiplicationFactor := utilisationRatio.Mul(assetRatesStats.StableSlope1)
			borrowAPY = assetRatesStats.StableBase.Add(multiplicationFactor)
			return borrowAPY, nil
		} else {
			utilisationNumerator := currentUtilisationRatio.Sub(assetRatesStats.UOptimal)
			utilisationDenominator := sdk.OneDec().Sub(assetRatesStats.UOptimal)
			utilisationRatio := utilisationNumerator.Quo(utilisationDenominator)
			multiplicationFactor := utilisationRatio.Mul(assetRatesStats.StableSlope2)
			borrowAPY = assetRatesStats.StableBase.Add(assetRatesStats.StableSlope1).Add(multiplicationFactor)
			return borrowAPY, nil
		}
	}
}

func (k Keeper) GetLendAPRByAssetIdAndPoolId(ctx sdk.Context, poolId, assetId uint64) (lendAPY sdk.Dec, err error) {
	assetRatesStats, found := k.GetAssetRatesStats(ctx, assetId)
	if !found {
		return sdk.ZeroDec(), types.ErrorAssetStatsNotFound
	}
	borrowAPY, err := k.GetBorrowAPRByAssetId(ctx, poolId, assetId, false)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	currentUtilisationRatio, err := k.GetUtilisationRatioByPoolIdAndAssetId(ctx, poolId, assetId)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	mulFactor := sdk.OneDec().Sub(assetRatesStats.ReserveFactor)
	lendAPY = borrowAPY.Mul(currentUtilisationRatio).Mul(mulFactor)

	return lendAPY, nil
}

func (k Keeper) GetAverageBorrowRate(ctx sdk.Context, poolId, assetId uint64) (averageBorrowRate sdk.Dec, err error) {

	assetStats, _ := k.UpdateAPR(ctx, assetId, poolId)
	factor1 := assetStats.BorrowApr.Mul(sdk.Dec(assetStats.TotalBorrowed))
	factor2 := assetStats.StableBorrowApr.Mul(sdk.Dec(assetStats.TotalStableBorrowed))
	numerator := factor1.Add(factor2)
	denominator := sdk.Dec(assetStats.TotalStableBorrowed).Add(sdk.Dec(assetStats.TotalBorrowed))
	averageBorrowRate = numerator.Quo(denominator)
	return averageBorrowRate, nil
}

func (k Keeper) GetSavingRate(ctx sdk.Context, poolId, assetId uint64) (savingRate sdk.Dec, err error) {
	assetRatesStats, _ := k.GetAssetRatesStats(ctx, assetId)
	averageBorrowRate, err := k.GetAverageBorrowRate(ctx, poolId, assetId)
	if err != nil {
		return sdk.Dec{}, err
	}
	utilizationRatio, err := k.GetUtilisationRatioByPoolIdAndAssetId(ctx, poolId, assetId)
	if err != nil {
		return sdk.Dec{}, err
	}
	factor1 := sdk.OneDec().Sub(assetRatesStats.ReserveFactor)
	savingRate = averageBorrowRate.Mul(utilizationRatio).Mul(factor1)
	return savingRate, nil
}

func (k Keeper) GetReserveRate(ctx sdk.Context, poolId, assetId uint64) (reserveRate sdk.Dec, err error) {
	averageBorrowRate, err := k.GetAverageBorrowRate(ctx, poolId, assetId)
	if err != nil {
		return sdk.Dec{}, err
	}
	savingRate, err := k.GetSavingRate(ctx, poolId, assetId)
	if err != nil {
		return sdk.Dec{}, err
	}
	reserveRate = averageBorrowRate.Sub(savingRate)
	return reserveRate, nil
}

func (k Keeper) UpdateAPR(ctx sdk.Context, poolId, assetId uint64) (AssetStats types.AssetStats, found bool) {
	lendAPR, _ := k.GetLendAPRByAssetIdAndPoolId(ctx, poolId, assetId)
	borrowAPR, _ := k.GetBorrowAPRByAssetId(ctx, poolId, assetId, false)
	stableBorrowAPR, _ := k.GetBorrowAPRByAssetId(ctx, poolId, assetId, true)
	currentUtilisationRatio, _ := k.GetUtilisationRatioByPoolIdAndAssetId(ctx, poolId, assetId)
	AssetStats = types.AssetStats{
		LendApr:          lendAPR,
		BorrowApr:        borrowAPR,
		StableBorrowApr:  stableBorrowAPR,
		UtilisationRatio: currentUtilisationRatio,
	}
	return AssetStats, true
}
