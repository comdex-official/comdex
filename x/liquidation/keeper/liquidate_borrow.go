package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	"github.com/comdex-official/comdex/x/liquidation/types"
)

func (k Keeper) LiquidateBorrows(ctx sdk.Context) error {
	borrows, found := k.GetBorrows(ctx)
	params := k.GetParams(ctx)
	if !found {
		return nil
	}
	liquidationOffsetHolder, found := k.GetLiquidationOffsetHolder(ctx, lendtypes.AppID, types.VaultLiquidationsOffsetPrefix)
	if !found {
		liquidationOffsetHolder = types.NewLiquidationOffsetHolder(lendtypes.AppID, 0)
	}
	borrowIDs := borrows
	start, end := types.GetSliceStartEndForLiquidations(len(borrowIDs), int(liquidationOffsetHolder.CurrentOffset), int(params.LiquidationBatchSize))

	if start == end {
		liquidationOffsetHolder.CurrentOffset = 0
		start, end = types.GetSliceStartEndForLiquidations(len(borrowIDs), int(liquidationOffsetHolder.CurrentOffset), int(params.LiquidationBatchSize))
	}
	newBorrowIDs := borrowIDs[start:end]
	for l := range newBorrowIDs {
		borrowPos, found := k.GetBorrow(ctx, newBorrowIDs[l])
		if !found {
			continue
		}
		lendPair, _ := k.GetLendPair(ctx, borrowPos.PairID)
		lendPos, _ := k.GetLend(ctx, borrowPos.LendingID)
		killSwitchParams, _ := k.GetKillSwitchData(ctx, lendPos.AppID)
		if killSwitchParams.BreakerEnable {
			continue
		}
		pool, _ := k.GetPool(ctx, lendPos.PoolID)
		assetIn, _ := k.GetAsset(ctx, lendPair.AssetIn)
		assetOut, _ := k.GetAsset(ctx, lendPair.AssetOut)
		var currentCollateralizationRatio sdk.Dec

		var firstTransitAssetID, secondTransitAssetID uint64
		for _, data := range pool.AssetData {
			if data.AssetTransitType == 2 {
				firstTransitAssetID = data.AssetID
			}
			if data.AssetTransitType == 3 {
				secondTransitAssetID = data.AssetID
			}
		}

		liqThreshold, _ := k.GetAssetRatesParams(ctx, lendPair.AssetIn)
		liqThresholdBridgedAssetOne, _ := k.GetAssetRatesParams(ctx, firstTransitAssetID)
		liqThresholdBridgedAssetTwo, _ := k.GetAssetRatesParams(ctx, secondTransitAssetID)
		firstBridgedAsset, _ := k.GetAsset(ctx, firstTransitAssetID)
		if borrowPos.BridgedAssetAmount.Amount.Equal(sdk.ZeroInt()) {
			currentCollateralizationRatio, _ = k.CalculateLendCollateralizationRatio(ctx, borrowPos.AmountIn.Amount, assetIn, borrowPos.AmountOut.Amount.Add(borrowPos.InterestAccumulated.TruncateInt()), assetOut)
			if sdk.Dec.GT(currentCollateralizationRatio, liqThreshold.LiquidationThreshold) {
				err := k.CreateLockedBorrow(ctx, borrowPos, currentCollateralizationRatio, lendPos.AppID)
				if err != nil {
					continue
				}
				borrowPos.IsLiquidated = true
				k.SetBorrow(ctx, borrowPos)
				lockedVaultID := k.GetLockedVaultID(ctx)
				err = k.UpdateLockedBorrows(ctx, lendPos.AppID, lockedVaultID+1)
				if err != nil {
					return nil
				}
			}
		} else {
			if borrowPos.BridgedAssetAmount.Denom == firstBridgedAsset.Denom {
				currentCollateralizationRatio, _ = k.CalculateLendCollateralizationRatio(ctx, borrowPos.AmountIn.Amount, assetIn, borrowPos.AmountOut.Amount.Add(borrowPos.InterestAccumulated.TruncateInt()), assetOut)
				if sdk.Dec.GT(currentCollateralizationRatio, liqThreshold.LiquidationThreshold.Mul(liqThresholdBridgedAssetOne.LiquidationThreshold)) {
					err := k.CreateLockedBorrow(ctx, borrowPos, currentCollateralizationRatio, lendPos.AppID)
					if err != nil {
						continue
					}
					borrowPos.IsLiquidated = true
					k.SetBorrow(ctx, borrowPos)
					lockedVaultID := k.GetLockedVaultID(ctx)
					err = k.UpdateLockedBorrows(ctx, lendPos.AppID, lockedVaultID+1)
					if err != nil {
						return nil
					}

				}
			} else {
				currentCollateralizationRatio, _ = k.CalculateLendCollateralizationRatio(ctx, borrowPos.AmountIn.Amount, assetIn, borrowPos.AmountOut.Amount.Add(borrowPos.InterestAccumulated.TruncateInt()), assetOut)

				if sdk.Dec.GT(currentCollateralizationRatio, liqThreshold.LiquidationThreshold.Mul(liqThresholdBridgedAssetTwo.LiquidationThreshold)) {
					err := k.CreateLockedBorrow(ctx, borrowPos, currentCollateralizationRatio, lendPos.AppID)
					if err != nil {
						continue
					}
					borrowPos.IsLiquidated = true
					k.SetBorrow(ctx, borrowPos)
					lockedVaultID := k.GetLockedVaultID(ctx)
					err = k.UpdateLockedBorrows(ctx, lendPos.AppID, lockedVaultID+1)
					if err != nil {
						return nil
					}
				}
			}
		}
	}
	liquidationOffsetHolder.CurrentOffset = uint64(end)
	k.SetLiquidationOffsetHolder(ctx, types.VaultLiquidationsOffsetPrefix, liquidationOffsetHolder)

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
	value := types.LockedVault{
		LockedVaultId:                lockedVaultID + 1,
		AppId:                        appID,
		OriginalVaultId:              borrow.ID,
		ExtendedPairId:               borrow.PairID,
		Owner:                        lendPos.Owner,
		AmountIn:                     borrow.AmountIn.Amount,
		AmountOut:                    borrow.AmountOut.Amount,
		UpdatedAmountOut:             borrow.AmountOut.Amount.Add(borrow.InterestAccumulated.TruncateInt()),
		Initiator:                    types.ModuleName,
		IsAuctionComplete:            false,
		IsAuctionInProgress:          false,
		CrAtLiquidation:              collateralizationRatio,
		CurrentCollaterlisationRatio: collateralizationRatio,
		CollateralToBeAuctioned:      sdk.ZeroDec(),
		LiquidationTimestamp:         ctx.BlockTime(),
		SellOffHistory:               nil,
		InterestAccumulated:          sdk.ZeroInt(),
		Kind:                         kind,
	}
	k.SetLockedVault(ctx, value)
	k.SetLockedVaultID(ctx, value.LockedVaultId)
	return nil
}

func (k Keeper) UpdateLockedBorrows(ctx sdk.Context, appID, id uint64) error {
	lockedVault, _ := k.GetLockedVault(ctx, appID, id)
	pair, _ := k.GetLendPair(ctx, lockedVault.ExtendedPairId)
	borrowMetaData := lockedVault.GetBorrowMetaData()
	if borrowMetaData != nil {
		lendPos, _ := k.GetLend(ctx, borrowMetaData.LendingId)
		pool, _ := k.GetPool(ctx, lendPos.PoolID)
		var unliquidatePointPercentage sdk.Dec
		var firstTransitAssetID, secondTransitAssetID uint64
		for _, data := range pool.AssetData {
			if data.AssetTransitType == 2 {
				firstTransitAssetID = data.AssetID
			}
			if data.AssetTransitType == 3 {
				secondTransitAssetID = data.AssetID
			}
		}

		firstBridgeAsset, _ := k.GetAsset(ctx, firstTransitAssetID)
		firstBridgeAssetStats, _ := k.GetAssetRatesParams(ctx, firstTransitAssetID)
		secondBridgeAssetStats, _ := k.GetAssetRatesParams(ctx, secondTransitAssetID)
		liqThreshold, _ := k.GetAssetRatesParams(ctx, pair.AssetIn)

		if !borrowMetaData.BridgedAssetAmount.Amount.Equal(sdk.ZeroInt()) {
			if borrowMetaData.BridgedAssetAmount.Denom == firstBridgeAsset.Denom {
				unliquidatePointPercentage = liqThreshold.LiquidationThreshold.Mul(firstBridgeAssetStats.LiquidationThreshold)
			} else {
				unliquidatePointPercentage = liqThreshold.LiquidationThreshold.Mul(secondBridgeAssetStats.LiquidationThreshold)
			}
		} else {
			unliquidatePointPercentage = liqThreshold.LiquidationThreshold
		}

		assetRatesStats, _ := k.GetAssetRatesParams(ctx, pair.AssetIn)

		if (!lockedVault.IsAuctionInProgress && !lockedVault.IsAuctionComplete) || (lockedVault.IsAuctionComplete && lockedVault.CurrentCollaterlisationRatio.GTE(unliquidatePointPercentage)) {
			assetIn, _ := k.GetAsset(ctx, pair.AssetIn)

			assetOut, _ := k.GetAsset(ctx, pair.AssetOut)

			collateralizationRatio, err := k.CalculateLendCollateralizationRatio(ctx, lockedVault.AmountIn, assetIn, lockedVault.UpdatedAmountOut, assetOut)
			if err != nil {
				return nil
			}

			assetInTotal, _ := k.CalcAssetPrice(ctx, assetIn.Id, lockedVault.AmountIn)
			assetOutTotal, _ := k.CalcAssetPrice(ctx, assetOut.Id, lockedVault.AmountOut)
			deductionPercentage, _ := sdk.NewDecFromStr("1.0")

			var c sdk.Dec
			if !borrowMetaData.BridgedAssetAmount.Amount.Equal(sdk.ZeroInt()) {
				if borrowMetaData.BridgedAssetAmount.Denom == firstBridgeAsset.Denom {
					c = assetRatesStats.LiquidationThreshold.Mul(firstBridgeAssetStats.Ltv)
				} else {
					c = assetRatesStats.LiquidationThreshold.Mul(secondBridgeAssetStats.Ltv)
				}
			} else {
				c = assetRatesStats.LiquidationThreshold
			}
			penalty := assetRatesStats.LiquidationPenalty.Add(assetRatesStats.LiquidationBonus)
			b := deductionPercentage.Add(penalty)
			totalIn := assetInTotal.ToDec()
			totalOut := assetOutTotal.ToDec()
			factor1 := c.Mul(totalIn)
			factor2 := b.Mul(c)
			numerator := totalOut.Sub(factor1)
			denominator := deductionPercentage.Sub(factor2)
			selloffAmount := numerator.Quo(denominator)
			updatedLockedVault := lockedVault
			if lockedVault.SellOffHistory == nil {
				//TODO: revisit
				aip := sdk.NewDec(int64(1))
				liquidationDeductionAmt := selloffAmount.Mul(penalty)
				liquidationDeductionAmount := liquidationDeductionAmt.Quo(aip)
				bonusToBidderAmt := selloffAmount.Mul(assetRatesStats.LiquidationBonus)

				bonusToBidderAmount := bonusToBidderAmt.Quo(aip)
				penaltyToReserveAmt := selloffAmount.Mul(assetRatesStats.LiquidationPenalty)
				penaltyToReserveAmount := penaltyToReserveAmt.Quo(aip)
				err = k.SendCoinsFromModuleToModule(ctx, pool.ModuleName, auctiontypes.ModuleName, sdk.NewCoins(sdk.NewCoin(assetIn.Denom, sdk.NewInt(bonusToBidderAmount.TruncateInt64()))))
				if err != nil {
					return err
				}
				err = k.UpdateReserveBalances(ctx, pair.AssetIn, pool.ModuleName, sdk.NewCoin(assetIn.Denom, sdk.NewInt(penaltyToReserveAmount.TruncateInt64())), true)
				if err != nil {
					return err
				}
				cAsset, _ := k.GetAsset(ctx, assetRatesStats.CAssetID)
				updatedLockedVault.AmountIn = updatedLockedVault.AmountIn.Sub(sdk.NewInt(liquidationDeductionAmount.TruncateInt64()))
				lendPos.AmountIn.Amount = lendPos.AmountIn.Amount.Sub(sdk.NewInt(liquidationDeductionAmount.TruncateInt64()))
				lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Sub(sdk.NewInt(liquidationDeductionAmount.TruncateInt64()))
				err = k.BurnCoin(ctx, pool.ModuleName, sdk.NewCoin(cAsset.Denom, sdk.NewInt(penaltyToReserveAmount.TruncateInt64())))
				if err != nil {
					return err
				}
				k.SetLend(ctx, lendPos)
			}
			var collateralToBeAuctioned sdk.Dec

			if selloffAmount.GTE(totalIn) {
				collateralToBeAuctioned = totalIn
			} else {
				collateralToBeAuctioned = selloffAmount
			}
			updatedLockedVault.CurrentCollaterlisationRatio = collateralizationRatio
			updatedLockedVault.CollateralToBeAuctioned = collateralToBeAuctioned
			k.SetLockedVault(ctx, updatedLockedVault)
		}
		newUpdatedLockedVault, _ := k.GetLockedVault(ctx, appID, id)
		err := k.LendDutchActivator(ctx, newUpdatedLockedVault)
		if err != nil {
			ctx.Logger().Error("error in dutch activator")
		}
	}
	return nil
}

func (k Keeper) UnLiquidateLockedBorrows(ctx sdk.Context, appID, id uint64, dutchAuction auctiontypes.DutchAuction) error {
	lockedVault, _ := k.GetLockedVault(ctx, appID, id)
	borrowMetadata := lockedVault.GetBorrowMetaData()
	if borrowMetadata != nil {
		lendPos, _ := k.GetLend(ctx, borrowMetadata.LendingId)
		assetInPool, _ := k.GetPool(ctx, lendPos.PoolID)
		var firstTransitAssetID, secondTransitAssetID uint64
		for _, data := range assetInPool.AssetData {
			if data.AssetTransitType == 2 {
				firstTransitAssetID = data.AssetID
			}
			if data.AssetTransitType == 3 {
				secondTransitAssetID = data.AssetID
			}
		}
		firstBridgedAsset, _ := k.GetAsset(ctx, firstTransitAssetID)
		userAddress, _ := sdk.AccAddressFromBech32(lockedVault.Owner)
		pair, _ := k.GetLendPair(ctx, lockedVault.ExtendedPairId)
		assetStats, _ := k.GetAssetRatesParams(ctx, pair.AssetIn)
		assetIn, _ := k.GetAsset(ctx, pair.AssetIn)
		assetOut, _ := k.GetAsset(ctx, pair.AssetOut)
		cAssetIn, _ := k.GetAsset(ctx, assetStats.CAssetID)

		if lockedVault.IsAuctionComplete {
			if borrowMetadata.BridgedAssetAmount.IsZero() {
				// also calculate the current collaterlization ratio to ensure there is no sudden changes
				liqThreshold, _ := k.GetAssetRatesParams(ctx, pair.AssetIn)
				unliquidatePointPercentage := liqThreshold.LiquidationThreshold

				if lockedVault.AmountOut.IsZero() {
					err := k.CreateLockedVaultHistory(ctx, lockedVault)
					if err != nil {
						return err
					}
					k.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
					if err = k.SendCoinFromModuleToAccount(ctx, assetInPool.ModuleName, userAddress, sdk.NewCoin(cAssetIn.Denom, lockedVault.AmountIn)); err != nil {
						return err
					}
					lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Add(lockedVault.AmountIn)
					k.SetLend(ctx, lendPos)
				}
				newCalculatedCollateralizationRatio, _ := k.CalculateLendCollateralizationRatio(ctx, lockedVault.AmountIn, assetIn, lockedVault.UpdatedAmountOut, assetOut)
				if newCalculatedCollateralizationRatio.GT(unliquidatePointPercentage) {
					updatedLockedVault := lockedVault
					updatedLockedVault.CurrentCollaterlisationRatio = newCalculatedCollateralizationRatio
					updatedLockedVault.SellOffHistory = append(updatedLockedVault.SellOffHistory, dutchAuction.String())
					k.SetLockedVault(ctx, updatedLockedVault)
					err := k.UpdateLockedBorrows(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
					if err != nil {
						return nil
					}
				}
				if newCalculatedCollateralizationRatio.LTE(unliquidatePointPercentage) {
					err := k.CreateLockedVaultHistory(ctx, lockedVault)
					if err != nil {
						return err
					}
					k.CreteNewBorrow(ctx, lockedVault)
					k.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
				}
			} else {
				if borrowMetadata.BridgedAssetAmount.Denom == firstBridgedAsset.Denom {
					liqThresholdAssetIn, _ := k.GetAssetRatesParams(ctx, pair.AssetIn)
					liqThresholdFirstBridgedAsset, _ := k.GetAssetRatesParams(ctx, firstTransitAssetID)
					liqThreshold := liqThresholdAssetIn.LiquidationThreshold.Mul(liqThresholdFirstBridgedAsset.LiquidationThreshold)
					unliquidatePointPercentage := liqThreshold

					if lockedVault.AmountOut.IsZero() {
						err := k.CreateLockedVaultHistory(ctx, lockedVault)
						if err != nil {
							return err
						}
						k.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
						if err = k.SendCoinFromModuleToAccount(ctx, assetInPool.ModuleName, userAddress, sdk.NewCoin(cAssetIn.Denom, lockedVault.AmountIn)); err != nil {
							return err
						}
						lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Add(lockedVault.AmountIn)
						k.SetLend(ctx, lendPos)
					}
					newCalculatedCollateralizationRatio, _ := k.CalculateLendCollateralizationRatio(ctx, lockedVault.AmountIn, assetIn, lockedVault.UpdatedAmountOut, assetOut)
					if newCalculatedCollateralizationRatio.GT(unliquidatePointPercentage) {
						updatedLockedVault := lockedVault
						updatedLockedVault.CurrentCollaterlisationRatio = newCalculatedCollateralizationRatio
						updatedLockedVault.SellOffHistory = append(updatedLockedVault.SellOffHistory, dutchAuction.String())
						k.SetLockedVault(ctx, updatedLockedVault)
						err := k.UpdateLockedBorrows(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
						if err != nil {
							return nil
						}
					}
					if newCalculatedCollateralizationRatio.LTE(unliquidatePointPercentage) {
						err := k.CreateLockedVaultHistory(ctx, lockedVault)
						if err != nil {
							return err
						}
						k.CreteNewBorrow(ctx, lockedVault)
						k.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
					}
				} else {
					liqThresholdAssetIn, _ := k.GetAssetRatesParams(ctx, pair.AssetIn)
					liqThresholdSecondBridgedAsset, _ := k.GetAssetRatesParams(ctx, secondTransitAssetID)
					liqThreshold := liqThresholdAssetIn.LiquidationThreshold.Mul(liqThresholdSecondBridgedAsset.LiquidationThreshold)
					unliquidatePointPercentage := liqThreshold

					if lockedVault.AmountOut.IsZero() {
						err := k.CreateLockedVaultHistory(ctx, lockedVault)
						if err != nil {
							return err
						}
						k.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
						if err = k.SendCoinFromModuleToAccount(ctx, assetInPool.ModuleName, userAddress, sdk.NewCoin(cAssetIn.Denom, lockedVault.AmountIn)); err != nil {
							return err
						}
						lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Add(lockedVault.AmountIn)
						k.SetLend(ctx, lendPos)
					}
					newCalculatedCollateralizationRatio, _ := k.CalculateLendCollateralizationRatio(ctx, lockedVault.AmountIn, assetIn, lockedVault.UpdatedAmountOut, assetOut)
					if newCalculatedCollateralizationRatio.GT(unliquidatePointPercentage) {
						updatedLockedVault := lockedVault
						updatedLockedVault.CurrentCollaterlisationRatio = newCalculatedCollateralizationRatio
						updatedLockedVault.SellOffHistory = append(updatedLockedVault.SellOffHistory, dutchAuction.String())
						k.SetLockedVault(ctx, updatedLockedVault)
						err := k.UpdateLockedBorrows(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
						if err != nil {
							return nil
						}
					}
					if newCalculatedCollateralizationRatio.LTE(unliquidatePointPercentage) {
						err := k.CreateLockedVaultHistory(ctx, lockedVault)
						if err != nil {
							return err
						}
						k.CreteNewBorrow(ctx, lockedVault)
						k.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
					}
				}
			}
		}
	}
	return nil
}
