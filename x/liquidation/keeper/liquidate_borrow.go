package keeper

import (
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	"github.com/comdex-official/comdex/x/liquidation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		liqThresholdBridgedAssetOne, _ := k.GetAssetRatesStats(ctx, pool.FirstBridgedAssetID)
		liqThresholdBridgedAssetTwo, _ := k.GetAssetRatesStats(ctx, pool.SecondBridgedAssetID)
		firstBridgedAsset, _ := k.GetAsset(ctx, pool.FirstBridgedAssetID)
		if borrowPos.BridgedAssetAmount.Amount.Equal(sdk.ZeroInt()) {
			currentCollateralizationRatio, _ = k.CalculateLendCollaterlizationRatio(ctx, borrowPos.AmountIn.Amount, assetIn, borrowPos.UpdatedAmountOut, assetOut)
			if sdk.Dec.GT(currentCollateralizationRatio, liqThreshold.LiquidationThreshold) {
				err := k.CreateLockedBorrow(ctx, borrowPos, currentCollateralizationRatio, lendPos.AppID)
				if err != nil {
					continue
				}
				k.UpdateBorrowStats(ctx, lendPair, borrowPos, borrowPos.AmountOut.Amount, false)
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

		} else {
			if borrowPos.BridgedAssetAmount.Denom == firstBridgedAsset.Denom {
				currentCollateralizationRatio, _ = k.CalculateLendCollaterlizationRatio(ctx, borrowPos.AmountIn.Amount, assetIn, borrowPos.UpdatedAmountOut, assetOut)
				if sdk.Dec.GT(currentCollateralizationRatio, liqThreshold.LiquidationThreshold.Mul(liqThresholdBridgedAssetOne.LiquidationThreshold)) {
					err := k.CreateLockedBorrow(ctx, borrowPos, currentCollateralizationRatio, lendPos.AppID)
					if err != nil {
						continue
					}
					k.UpdateBorrowStats(ctx, lendPair, borrowPos, borrowPos.AmountOut.Amount, false)
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
			} else {
				currentCollateralizationRatio, _ = k.CalculateLendCollaterlizationRatio(ctx, borrowPos.AmountIn.Amount, assetIn, borrowPos.UpdatedAmountOut, assetOut)

				if sdk.Dec.GT(currentCollateralizationRatio, liqThreshold.LiquidationThreshold.Mul(liqThresholdBridgedAssetTwo.LiquidationThreshold)) {
					err := k.CreateLockedBorrow(ctx, borrowPos, currentCollateralizationRatio, lendPos.AppID)
					if err != nil {
						continue
					}
					k.UpdateBorrowStats(ctx, lendPair, borrowPos, borrowPos.AmountOut.Amount, false)
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
		}
		lockedVaultID := k.GetLockedVaultID(ctx)
		err := k.UpdateLockedBorrows(ctx, lendPos.AppID, lockedVaultID)
		if err != nil {
			return nil
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
		OriginalVaultId:              borrow.ID,
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

func (k Keeper) UpdateLockedBorrows(ctx sdk.Context, appID, id uint64) error {
	lockedVault, _ := k.GetLockedVault(ctx, appID, id)
	pair, _ := k.GetLendPair(ctx, lockedVault.ExtendedPairId)
	borrowMetaData := lockedVault.GetBorrowMetaData()
	if borrowMetaData != nil {
		lendPos, _ := k.GetLend(ctx, borrowMetaData.LendingId)
		pool, _ := k.GetPool(ctx, lendPos.PoolID)
		var unliquidatePointPercentage sdk.Dec
		firstBridgeAsset, _ := k.GetAsset(ctx, pool.FirstBridgedAssetID)
		firstBridgeAssetStats, _ := k.GetAssetRatesStats(ctx, pool.FirstBridgedAssetID)
		secondBridgeAssetStats, _ := k.GetAssetRatesStats(ctx, pool.SecondBridgedAssetID)
		liqThreshold, _ := k.GetAssetRatesStats(ctx, pair.AssetIn)

		if !borrowMetaData.BridgedAssetAmount.Amount.Equal(sdk.ZeroInt()) {
			if borrowMetaData.BridgedAssetAmount.Denom == firstBridgeAsset.Denom {
				unliquidatePointPercentage = liqThreshold.LiquidationThreshold.Mul(firstBridgeAssetStats.LiquidationThreshold)
			} else {
				unliquidatePointPercentage = liqThreshold.LiquidationThreshold.Mul(secondBridgeAssetStats.LiquidationThreshold)
			}
		} else {
			unliquidatePointPercentage = liqThreshold.LiquidationThreshold
		}

		assetRatesStats, _ := k.GetAssetRatesStats(ctx, pair.AssetIn)

		if (!lockedVault.IsAuctionInProgress && !lockedVault.IsAuctionComplete) || (lockedVault.IsAuctionComplete && lockedVault.CurrentCollaterlisationRatio.GTE(unliquidatePointPercentage)) {

			assetIn, _ := k.GetAsset(ctx, pair.AssetIn)

			assetOut, _ := k.GetAsset(ctx, pair.AssetOut)

			collateralizationRatio, err := k.CalculateLendCollaterlizationRatio(ctx, lockedVault.AmountIn, assetIn, lockedVault.UpdatedAmountOut, assetOut)
			if err != nil {
				return nil
			}

			assetInPrice, _ := k.GetPriceForAsset(ctx, assetIn.Id)
			assetOutPrice, _ := k.GetPriceForAsset(ctx, assetOut.Id)
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
			totalIn := lockedVault.AmountIn.Mul(sdk.NewIntFromUint64(assetInPrice)).ToDec()
			totalOut := lockedVault.UpdatedAmountOut.Mul(sdk.NewIntFromUint64(assetOutPrice)).ToDec()
			factor1 := c.Mul(totalIn)
			factor2 := b.Mul(c)
			numerator := totalOut.Sub(factor1)
			denominator := deductionPercentage.Sub(factor2)
			selloffAmount := numerator.Quo(denominator)
			updatedLockedVault := lockedVault
			if lockedVault.SellOffHistory == nil {
				aip := sdk.NewDec(int64(assetInPrice))
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
				lendPos.UpdatedAmountIn = lendPos.UpdatedAmountIn.Sub(sdk.NewInt(liquidationDeductionAmount.TruncateInt64()))
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
		firstBridgedAsset, _ := k.GetAsset(ctx, assetInPool.FirstBridgedAssetID)
		userAddress, _ := sdk.AccAddressFromBech32(lockedVault.Owner)
		pair, _ := k.GetLendPair(ctx, lockedVault.ExtendedPairId)
		assetStats, _ := k.GetAssetRatesStats(ctx, pair.AssetIn)
		assetIn, _ := k.GetAsset(ctx, pair.AssetIn)
		assetOut, _ := k.GetAsset(ctx, pair.AssetOut)
		cAssetIn, _ := k.GetAsset(ctx, assetStats.CAssetID)

		if lockedVault.IsAuctionComplete {
			if borrowMetadata.BridgedAssetAmount.IsZero() {
				//also calculate the current collaterlization ratio to ensure there is no sudden changes
				liqThreshold, _ := k.GetAssetRatesStats(ctx, pair.AssetIn)
				unliquidatePointPercentage := liqThreshold.LiquidationThreshold

				if lockedVault.AmountOut.IsZero() {
					err := k.CreateLockedVaultHistory(ctx, lockedVault)
					if err != nil {
						return err
					}
					k.DeleteBorrowForAddressByPair(ctx, userAddress, lockedVault.ExtendedPairId)
					k.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
					if err = k.SendCoinFromModuleToAccount(ctx, assetInPool.ModuleName, userAddress, sdk.NewCoin(cAssetIn.Denom, lockedVault.AmountIn)); err != nil {
						return err
					}
					lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Add(lockedVault.AmountIn)
					k.SetLend(ctx, lendPos)
				}
				newCalculatedCollateralizationRatio, _ := k.CalculateLendCollaterlizationRatio(ctx, lockedVault.AmountIn, assetIn, lockedVault.UpdatedAmountOut, assetOut)
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
					k.DeleteBorrowForAddressByPair(ctx, userAddress, lockedVault.ExtendedPairId)
					k.CreteNewBorrow(ctx, lockedVault)
					k.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
				}
			} else {
				if borrowMetadata.BridgedAssetAmount.Denom == firstBridgedAsset.Denom {
					liqThresholdAssetIn, _ := k.GetAssetRatesStats(ctx, pair.AssetIn)
					liqThresholdFirstBridgedAsset, _ := k.GetAssetRatesStats(ctx, assetInPool.FirstBridgedAssetID)
					liqThreshold := liqThresholdAssetIn.LiquidationThreshold.Mul(liqThresholdFirstBridgedAsset.LiquidationThreshold)
					unliquidatePointPercentage := liqThreshold

					if lockedVault.AmountOut.IsZero() {
						err := k.CreateLockedVaultHistory(ctx, lockedVault)
						if err != nil {
							return err
						}
						k.DeleteBorrowForAddressByPair(ctx, userAddress, lockedVault.ExtendedPairId)
						k.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
						if err = k.SendCoinFromModuleToAccount(ctx, assetInPool.ModuleName, userAddress, sdk.NewCoin(cAssetIn.Denom, lockedVault.AmountIn)); err != nil {
							return err
						}
						lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Add(lockedVault.AmountIn)
						k.SetLend(ctx, lendPos)
					}
					newCalculatedCollateralizationRatio, _ := k.CalculateLendCollaterlizationRatio(ctx, lockedVault.AmountIn, assetIn, lockedVault.UpdatedAmountOut, assetOut)
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
						k.DeleteBorrowForAddressByPair(ctx, userAddress, lockedVault.ExtendedPairId)
						k.CreteNewBorrow(ctx, lockedVault)
						k.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
					}
				} else {
					liqThresholdAssetIn, _ := k.GetAssetRatesStats(ctx, pair.AssetIn)
					liqThresholdSecondBridgedAsset, _ := k.GetAssetRatesStats(ctx, assetInPool.SecondBridgedAssetID)
					liqThreshold := liqThresholdAssetIn.LiquidationThreshold.Mul(liqThresholdSecondBridgedAsset.LiquidationThreshold)
					unliquidatePointPercentage := liqThreshold

					if lockedVault.AmountOut.IsZero() {
						err := k.CreateLockedVaultHistory(ctx, lockedVault)
						if err != nil {
							return err
						}
						k.DeleteBorrowForAddressByPair(ctx, userAddress, lockedVault.ExtendedPairId)
						k.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
						if err = k.SendCoinFromModuleToAccount(ctx, assetInPool.ModuleName, userAddress, sdk.NewCoin(cAssetIn.Denom, lockedVault.AmountIn)); err != nil {
							return err
						}
						lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Add(lockedVault.AmountIn)
						k.SetLend(ctx, lendPos)
					}
					newCalculatedCollateralizationRatio, _ := k.CalculateLendCollaterlizationRatio(ctx, lockedVault.AmountIn, assetIn, lockedVault.UpdatedAmountOut, assetOut)
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
						k.DeleteBorrowForAddressByPair(ctx, userAddress, lockedVault.ExtendedPairId)
						k.CreteNewBorrow(ctx, lockedVault)
						k.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
					}
				}
			}
		}
	}
	return nil
}
