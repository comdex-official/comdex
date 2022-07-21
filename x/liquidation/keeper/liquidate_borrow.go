package keeper

import (
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	"github.com/comdex-official/comdex/x/liquidation/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
)

func (k Keeper) LiquidateBorrows(ctx sdk.Context) error {
	borrowIds, found := k.GetBorrows(ctx)
	if !found {
		return nil
	}
	for _, v := range borrowIds.BorrowIDs {
		borrowPos, found := k.GetBorrow(ctx, v)
		if !found {
			continue
		}
		lendPair, _ := k.GetLendPair(ctx, borrowPos.PairID)
		lendPos, _ := k.GetLend(ctx, borrowPos.LendingID)
		pool, _ := k.GetPool(ctx, lendPos.PoolID)
		assetIn, _ := k.GetAsset(ctx, lendPair.AssetIn)
		assetOut, _ := k.GetAsset(ctx, lendPair.AssetOut)
		var currentCollateralizationRatio sdk.Dec

		liqThreshold, _ := k.GetAssetRatesStats(ctx, lendPair.AssetIn)
		if borrowPos.BridgedAssetAmount.Amount.Equal(sdk.ZeroInt()) {
			currentCollateralizationRatio, _ = k.CalculateLendCollaterlizationRatio(ctx, borrowPos.AmountIn.Amount, assetIn, borrowPos.UpdatedAmountOut, assetOut)
		} else {
			firstBridgedAsset, _ := k.GetAsset(ctx, pool.FirstBridgedAssetID)
			secondBridgedAsset, _ := k.GetAsset(ctx, pool.SecondBridgedAssetID)
			if borrowPos.BridgedAssetAmount.Denom == firstBridgedAsset.Denom {
				currentCollateralizationRatio, _ = k.CalculateLendCollaterlizationRatio(ctx, borrowPos.BridgedAssetAmount.Amount, firstBridgedAsset, borrowPos.UpdatedAmountOut, assetOut)
			} else {
				currentCollateralizationRatio, _ = k.CalculateLendCollaterlizationRatio(ctx, borrowPos.BridgedAssetAmount.Amount, secondBridgedAsset, borrowPos.UpdatedAmountOut, assetOut)
			}
		}

		if sdk.Dec.GT(currentCollateralizationRatio, liqThreshold.LiquidationThreshold) {
			err := k.CreateLockedBorrow(ctx, borrowPos, currentCollateralizationRatio, lendPos.AppID)
			if err != nil {
				continue
			}
			k.DeleteBorrow(ctx, v)
			err = k.UpdateUserBorrowIDMapping(ctx, lendPos.Owner, v, false)
			if err != nil {
				continue
			}
			err = k.UpdateBorrowIDByOwnerAndPoolMapping(ctx, lendPos.Owner, v, lendPair.AssetOutPoolID, false)
			if err != nil {
				continue
			}
			err = k.UpdateBorrowIdsMapping(ctx, v, false)
			if err != nil {
				continue
			}

		}
	}

	return nil
}

func (k Keeper) CreateLockedBorrow(ctx sdk.Context, borrow lendtypes.BorrowAsset, collateralizationRatio sdk.Dec, appID uint64) error {
	lockedVaultID := k.GetLockedVaultID(ctx)
	lendPos, _ := k.GetLend(ctx, borrow.LendingID)

	kind := &types.LockedVault_BorrowMetaData{
		BorrowMetaData: &types.BorrowMetaData{
			LendingId:          borrow.LendingID,
			IsStableBorrow:     borrow.IsStableBorrow,
			StableBorrowRate:   borrow.StableBorrowRate,
			BridgedAssetAmount: borrow.BridgedAssetAmount,
		},
	}
	var value = types.LockedVault{
		LockedVaultId:                lockedVaultID + 1,
		AppId:                        appID,
		AppVaultTypeId:               strconv.FormatUint(appID, 10),
		OriginalVaultId:              strconv.FormatUint(borrow.ID, 10),
		ExtendedPairId:               borrow.PairID,
		Owner:                        lendPos.Owner,
		AmountIn:                     borrow.AmountIn.Amount,
		AmountOut:                    borrow.AmountOut.Amount,
		UpdatedAmountOut:             borrow.AmountOut.Amount.Add(borrow.Interest_Accumulated),
		Initiator:                    types.ModuleName,
		IsAuctionComplete:            false,
		IsAuctionInProgress:          false,
		CrAtLiquidation:              collateralizationRatio,
		CurrentCollaterlisationRatio: collateralizationRatio,
		CollateralToBeAuctioned:      sdk.ZeroDec(),
		LiquidationTimestamp:         ctx.BlockTime(),
		SellOffHistory:               nil,
		InterestAccumulated:          borrow.Interest_Accumulated,
		Kind:                         kind,
	}

	k.SetLockedVault(ctx, value)
	k.SetLockedVaultID(ctx, value.LockedVaultId)
	return nil
}

func (k Keeper) UpdateLockedBorrows(ctx sdk.Context) error {
	lockedVaults := k.GetLockedVaults(ctx)
	if len(lockedVaults) == 0 {
		return nil
	}

	for _, lockedVault := range lockedVaults {

		pair, found := k.GetLendPair(ctx, lockedVault.ExtendedPairId)
		if !found {
			continue
		}

		liqThreshold, _ := k.GetAssetRatesStats(ctx, pair.AssetIn)
		unliquidatePointPercentage := liqThreshold.LiquidationThreshold

		assetRatesStats, found := k.GetAssetRatesStats(ctx, pair.AssetIn)
		if !found {
			return lendtypes.ErrorAssetStatsNotFound
		}

		if (!lockedVault.IsAuctionInProgress && !lockedVault.IsAuctionComplete) || (lockedVault.IsAuctionComplete && lockedVault.CurrentCollaterlisationRatio.LTE(unliquidatePointPercentage)) {

			assetIn, found := k.GetAsset(ctx, pair.AssetIn)
			if !found {
				continue
			}
			assetOut, found := k.GetAsset(ctx, pair.AssetOut)
			if !found {
				continue
			}
			collateralizationRatio, err := k.CalculateLendCollaterlizationRatio(ctx, lockedVault.AmountIn, assetIn, lockedVault.UpdatedAmountOut, assetOut)
			if err != nil {
				continue
			}

			assetInPrice, _ := k.GetPriceForAsset(ctx, assetIn.Id)
			assetOutPrice, _ := k.GetPriceForAsset(ctx, assetOut.Id)
			deductionPercentage, _ := sdk.NewDecFromStr("1.0")
			c := assetRatesStats.LiquidationThreshold
			b := deductionPercentage.Add(assetRatesStats.LiquidationPenalty)
			totalIn := lockedVault.AmountIn.Mul(sdk.NewIntFromUint64(assetInPrice)).ToDec()
			totalOut := lockedVault.UpdatedAmountOut.Mul(sdk.NewIntFromUint64(assetOutPrice)).ToDec()

			factor1 := c.Mul(totalIn)
			factor2 := b.Mul(c)
			numerator := totalOut.Sub(factor1)
			denominator := deductionPercentage.Sub(factor2)
			selloffAmount := numerator.Quo(denominator)

			var collateralToBeAuctioned sdk.Dec

			if selloffAmount.GTE(totalIn) {
				collateralToBeAuctioned = totalIn
			} else {

				collateralToBeAuctioned = selloffAmount
			}
			updatedLockedVault := lockedVault
			updatedLockedVault.CurrentCollaterlisationRatio = collateralizationRatio
			updatedLockedVault.CollateralToBeAuctioned = collateralToBeAuctioned
			k.SetLockedVault(ctx, updatedLockedVault)

		}

	}
	return nil
}

func (k Keeper) UnliquidateLockedBorrows(ctx sdk.Context) error {
	lockedVaults := k.GetLockedVaults(ctx)
	if len(lockedVaults) == 0 {
		return nil
	}
	for _, lockedVault := range lockedVaults {

		if lockedVault.IsAuctionComplete {
			//also calculate the current collaterlization ration to ensure there is no sudden changes
			userAddress, err := sdk.AccAddressFromBech32(lockedVault.Owner)
			if err != nil {
				continue
			}

			pair, found := k.GetLendPair(ctx, lockedVault.ExtendedPairId)
			if !found {
				continue
			}

			liqThreshold, _ := k.GetAssetRatesStats(ctx, pair.AssetIn)
			unliquidatePointPercentage := liqThreshold.LiquidationThreshold

			assetIn, found := k.GetAsset(ctx, pair.AssetIn)
			if !found {
				continue
			}
			assetOut, found := k.GetAsset(ctx, pair.AssetOut)
			if !found {
				continue
			}
			if lockedVault.AmountOut.IsZero() {
				err := k.CreateLockedVaultHistory(ctx, lockedVault)
				if err != nil {
					return err
				}
				k.DeleteBorrowForAddressByPair(ctx, userAddress, lockedVault.ExtendedPairId)
				k.DeleteLockedVault(ctx, lockedVault.LockedVaultId)
				if err := k.SendCoinFromModuleToAccount(ctx, vaulttypes.ModuleName, userAddress, sdk.NewCoin(assetIn.Denom, lockedVault.AmountIn)); err != nil {
					continue
				}
				continue
			}
			newCalculatedCollateralizationRatio, err := k.CalculateLendCollaterlizationRatio(ctx, lockedVault.AmountIn, assetIn, lockedVault.UpdatedAmountOut, assetOut)
			if err != nil {
				continue
			}
			if newCalculatedCollateralizationRatio.LT(unliquidatePointPercentage) {
				updatedLockedVault := lockedVault
				updatedLockedVault.CurrentCollaterlisationRatio = newCalculatedCollateralizationRatio
				k.SetLockedVault(ctx, updatedLockedVault)
				continue
			}
			if newCalculatedCollateralizationRatio.GTE(unliquidatePointPercentage) {
				err := k.CreateLockedVaultHistory(ctx, lockedVault)
				if err != nil {
					return err
				}
				k.DeleteBorrowForAddressByPair(ctx, userAddress, lockedVault.ExtendedPairId)
				k.CreteNewBorrow(ctx, lockedVault)
				k.DeleteLockedVault(ctx, lockedVault.LockedVaultId)
			}
		}
	}

	return nil
}
