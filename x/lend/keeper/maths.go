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
			assetStats, _ := k.GetAssetStatsByPoolIdAndAssetId(ctx, assetId, poolId)
			assetStats.BorrowApr = borrowAPY
			k.SetAssetStatsByPoolIdAndAssetId(ctx, assetStats)
			return borrowAPY, nil
		} else {
			utilisationNumerator := currentUtilisationRatio.Sub(assetRatesStats.UOptimal)
			utilisationDenominator := sdk.OneDec().Sub(assetRatesStats.UOptimal)
			utilisationRatio := utilisationNumerator.Quo(utilisationDenominator)
			multiplicationFactor := utilisationRatio.Mul(assetRatesStats.Slope2)
			borrowAPY = assetRatesStats.Base.Add(assetRatesStats.Slope1).Add(multiplicationFactor)
			assetStats, _ := k.GetAssetStatsByPoolIdAndAssetId(ctx, assetId, poolId)
			assetStats.BorrowApr = borrowAPY
			k.SetAssetStatsByPoolIdAndAssetId(ctx, assetStats)
			return borrowAPY, nil
		}
	} else {
		if currentUtilisationRatio.LT(assetRatesStats.UOptimal) {
			utilisationRatio := currentUtilisationRatio.Quo(assetRatesStats.UOptimal)
			multiplicationFactor := utilisationRatio.Mul(assetRatesStats.StableSlope1)
			borrowAPY = assetRatesStats.StableBase.Add(multiplicationFactor)
			assetStats, _ := k.GetAssetStatsByPoolIdAndAssetId(ctx, assetId, poolId)
			assetStats.StableBorrowApr = borrowAPY
			k.SetAssetStatsByPoolIdAndAssetId(ctx, assetStats)
			return borrowAPY, nil
		} else {
			utilisationNumerator := currentUtilisationRatio.Sub(assetRatesStats.UOptimal)
			utilisationDenominator := sdk.OneDec().Sub(assetRatesStats.UOptimal)
			utilisationRatio := utilisationNumerator.Quo(utilisationDenominator)
			multiplicationFactor := utilisationRatio.Mul(assetRatesStats.StableSlope2)
			borrowAPY = assetRatesStats.StableBase.Add(assetRatesStats.StableSlope1).Add(multiplicationFactor)
			assetStats, _ := k.GetAssetStatsByPoolIdAndAssetId(ctx, assetId, poolId)
			assetStats.StableBorrowApr = borrowAPY
			k.SetAssetStatsByPoolIdAndAssetId(ctx, assetStats)
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
	assetStats, _ := k.GetAssetStatsByPoolIdAndAssetId(ctx, assetId, poolId)
	assetStats.LendApr = lendAPY
	k.SetAssetStatsByPoolIdAndAssetId(ctx, assetStats)

	return lendAPY, nil
}
