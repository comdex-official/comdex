package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	"github.com/comdex-official/comdex/x/liquidation/types"
)

func (k Keeper) LiquidateBorrows(ctx sdk.Context) error {
	borrows, found := k.lend.GetBorrows(ctx)
	params := k.GetParams(ctx)
	if !found {
		ctx.Logger().Error("Params Not Found in Liquidation, liquidate_borrow.go")
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
		borrowPos, found := k.lend.GetBorrow(ctx, newBorrowIDs[l])
		if !found {
			ctx.Logger().Error("Borrow Pos Not Found in Liquidation, liquidate_borrow.go for ID %d", newBorrowIDs[l])
			continue
		}
		if borrowPos.IsLiquidated {
			continue
		}
		lendPair, _ := k.lend.GetLendPair(ctx, borrowPos.PairID)
		lendPos, found := k.lend.GetLend(ctx, borrowPos.LendingID)
		if !found {
			ctx.Logger().Error("Lend Pos Not Found in Liquidation, liquidate_borrow.go for ID %d", borrowPos.LendingID)
			continue
		}

		// calculating and updating the interest accumulated before checking for liquidations
		err := k.lend.MsgCalculateBorrowInterest(ctx, lendPos.Owner, borrowPos.ID)
		if err != nil {
			ctx.Logger().Error("error in calculating Borrow Interest before liquidation")
			continue
		}
		borrowPos, _ = k.lend.GetBorrow(ctx, newBorrowIDs[l])
		if !borrowPos.StableBorrowRate.Equal(sdk.ZeroDec()) {
			borrowPos, err = k.lend.ReBalanceStableRates(ctx, borrowPos)
			if err != nil {
				ctx.Logger().Error("error in re-balance stable rate check before liquidation")
				continue
			}
		}

		killSwitchParams, _ := k.esm.GetKillSwitchData(ctx, lendPos.AppID)
		if killSwitchParams.BreakerEnable {
			ctx.Logger().Error("Kill Switch is enabled in Liquidation, liquidate_borrow.go for ID %d", lendPos.AppID)
			continue
		}
		// calculating and updating the interest accumulated before checking for liquidations
		err1 := k.lend.MsgCalculateBorrowInterest(ctx, lendPos.Owner, borrowPos.ID)
		if err1 != nil {
			ctx.Logger().Error("Error in calculating Borrow Interest before liquidation, liquidate_borrow.go for ID %d", borrowPos.ID)
		}
		pool, _ := k.lend.GetPool(ctx, lendPos.PoolID)
		assetIn, _ := k.asset.GetAsset(ctx, lendPair.AssetIn)
		assetOut, _ := k.asset.GetAsset(ctx, lendPair.AssetOut)

		var currentCollateralizationRatio sdk.Dec
		var firstTransitAssetID, secondTransitAssetID uint64
		// for getting transit assets details
		for _, data := range pool.AssetData {
			if data.AssetTransitType == 2 {
				firstTransitAssetID = data.AssetID
			}
			if data.AssetTransitType == 3 {
				secondTransitAssetID = data.AssetID
			}
		}

		liqThreshold, _ := k.lend.GetAssetRatesParams(ctx, lendPair.AssetIn)
		liqThresholdBridgedAssetOne, _ := k.lend.GetAssetRatesParams(ctx, firstTransitAssetID)
		liqThresholdBridgedAssetTwo, _ := k.lend.GetAssetRatesParams(ctx, secondTransitAssetID)
		firstBridgedAsset, _ := k.asset.GetAsset(ctx, firstTransitAssetID)
		// there are three possible cases
		// 	a. if borrow is from same pool
		//  b. if borrow is from first transit asset
		//  c. if borrow is from second transit asset
		if borrowPos.BridgedAssetAmount.Amount.Equal(sdk.ZeroInt()) { // first condition
			currentCollateralizationRatio, _ = k.lend.CalculateCollateralizationRatio(ctx, borrowPos.AmountIn.Amount, assetIn, borrowPos.AmountOut.Amount.Add(borrowPos.InterestAccumulated.TruncateInt()), assetOut)
			if sdk.Dec.GT(currentCollateralizationRatio, liqThreshold.LiquidationThreshold) {
				// after checking the currentCollateralizationRatio with LiquidationThreshold if borrow is to be liquidated then
				// CreateLockedBorrow function is called
				lockedVault, err := k.CreateLockedBorrow(ctx, borrowPos, currentCollateralizationRatio, lendPos.AppID)
				if err != nil {
					ctx.Logger().Error("Error in first condition CreateLockedBorrow in Liquidation, liquidate_borrow.go for ID %d", borrowPos.LendingID)
					continue
				}
				borrowPos.IsLiquidated = true // isLiquidated flag is set to true
				k.lend.SetBorrow(ctx, borrowPos)
				err = k.UpdateLockedBorrows(ctx, lockedVault)
				if err != nil {
					ctx.Logger().Error("Error in first condition UpdateLockedBorrows in Liquidation, liquidate_borrow.go for ID %d", lockedVault.LockedVaultId)
					continue
				}
			}
		} else {
			if borrowPos.BridgedAssetAmount.Denom == firstBridgedAsset.Denom {
				currentCollateralizationRatio, _ = k.lend.CalculateCollateralizationRatio(ctx, borrowPos.AmountIn.Amount, assetIn, borrowPos.AmountOut.Amount.Add(borrowPos.InterestAccumulated.TruncateInt()), assetOut)
				if sdk.Dec.GT(currentCollateralizationRatio, liqThreshold.LiquidationThreshold.Mul(liqThresholdBridgedAssetOne.LiquidationThreshold)) {
					lockedVault, err := k.CreateLockedBorrow(ctx, borrowPos, currentCollateralizationRatio, lendPos.AppID)
					if err != nil {
						ctx.Logger().Error("Error in second condition CreateLockedBorrow in Liquidation, liquidate_borrow.go for ID %d", borrowPos.LendingID)
						continue
					}
					borrowPos.IsLiquidated = true
					k.lend.SetBorrow(ctx, borrowPos)
					err = k.UpdateLockedBorrows(ctx, lockedVault)
					if err != nil {
						ctx.Logger().Error("Error in second condition UpdateLockedBorrows in Liquidation, liquidate_borrow.go for ID %d", lockedVault.LockedVaultId)
						continue
					}
				}
			} else {
				currentCollateralizationRatio, _ = k.lend.CalculateCollateralizationRatio(ctx, borrowPos.AmountIn.Amount, assetIn, borrowPos.AmountOut.Amount.Add(borrowPos.InterestAccumulated.TruncateInt()), assetOut)

				if sdk.Dec.GT(currentCollateralizationRatio, liqThreshold.LiquidationThreshold.Mul(liqThresholdBridgedAssetTwo.LiquidationThreshold)) {
					lockedVault, err := k.CreateLockedBorrow(ctx, borrowPos, currentCollateralizationRatio, lendPos.AppID)
					if err != nil {
						ctx.Logger().Error("Error in third condition CreateLockedBorrow in Liquidation, liquidate_borrow.go for ID %d", borrowPos.LendingID)
						continue
					}
					borrowPos.IsLiquidated = true
					k.lend.SetBorrow(ctx, borrowPos)
					err = k.UpdateLockedBorrows(ctx, lockedVault)
					if err != nil {
						ctx.Logger().Error("Error in third condition UpdateLockedBorrows in Liquidation, liquidate_borrow.go for ID %d", lockedVault.LockedVaultId)
						continue
					}
				}
			}
		}
	}
	liquidationOffsetHolder.CurrentOffset = uint64(end)
	k.SetLiquidationOffsetHolder(ctx, types.VaultLiquidationsOffsetPrefix, liquidationOffsetHolder)

	return nil
}

func (k Keeper) CreateLockedBorrow(ctx sdk.Context, borrow lendtypes.BorrowAsset, collateralizationRatio sdk.Dec, appID uint64) (types.LockedVault, error) {
	lockedVaultID := k.GetLockedVaultID(ctx)
	lendPos, _ := k.lend.GetLend(ctx, borrow.LendingID)
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
	return value, nil
}

func (k Keeper) UpdateLockedBorrows(ctx sdk.Context, lockedVault types.LockedVault) error {
	pair, _ := k.lend.GetLendPair(ctx, lockedVault.ExtendedPairId)
	borrowMetaData := lockedVault.GetBorrowMetaData()
	if borrowMetaData != nil {
		lendPos, _ := k.lend.GetLend(ctx, borrowMetaData.LendingId)
		pool, _ := k.lend.GetPool(ctx, lendPos.PoolID)
		var unliquidatePointPercentage sdk.Dec
		// retrieving transit asset details from cPool
		var firstTransitAssetID, secondTransitAssetID uint64
		for _, data := range pool.AssetData {
			if data.AssetTransitType == 2 {
				firstTransitAssetID = data.AssetID
			}
			if data.AssetTransitType == 3 {
				secondTransitAssetID = data.AssetID
			}
		}

		firstBridgeAsset, _ := k.asset.GetAsset(ctx, firstTransitAssetID)
		firstBridgeAssetStats, _ := k.lend.GetAssetRatesParams(ctx, firstTransitAssetID)
		secondBridgeAssetStats, _ := k.lend.GetAssetRatesParams(ctx, secondTransitAssetID)
		liqThreshold, _ := k.lend.GetAssetRatesParams(ctx, pair.AssetIn)

		// finding unLiquidate Point percentage
		if !borrowMetaData.BridgedAssetAmount.Amount.Equal(sdk.ZeroInt()) { // if bridged asset is being used for borrow (inter-pool borrow)
			if borrowMetaData.BridgedAssetAmount.Denom == firstBridgeAsset.Denom {
				unliquidatePointPercentage = liqThreshold.LiquidationThreshold.Mul(firstBridgeAssetStats.LiquidationThreshold)
			} else {
				unliquidatePointPercentage = liqThreshold.LiquidationThreshold.Mul(secondBridgeAssetStats.LiquidationThreshold)
			}
		} else { // same pool borrow
			unliquidatePointPercentage = liqThreshold.LiquidationThreshold
		}

		assetRatesStats, _ := k.lend.GetAssetRatesParams(ctx, pair.AssetIn)
		updatedLockedVault := lockedVault
		// Checking required flags
		if (!lockedVault.IsAuctionInProgress && !lockedVault.IsAuctionComplete) || (lockedVault.IsAuctionComplete && lockedVault.CurrentCollaterlisationRatio.GTE(unliquidatePointPercentage)) {
			assetIn, _ := k.asset.GetAsset(ctx, pair.AssetIn)
			assetOut, _ := k.asset.GetAsset(ctx, pair.AssetOut)
			collateralizationRatio, err := k.lend.CalculateCollateralizationRatio(ctx, lockedVault.AmountIn, assetIn, lockedVault.UpdatedAmountOut, assetOut)
			if err != nil {
				ctx.Logger().Error("Error Calculating CR in Liquidation, liquidate_borrow.go for locked vault ID %d", lockedVault.LockedVaultId)
				return nil
			}
			assetInprice, _ := k.market.GetTwa(ctx, assetIn.Id)
			assetOutprice, _ := k.market.GetTwa(ctx, assetOut.Id)

			//assetInTotal, _ := k.market.CalcAssetPrice(ctx, assetIn.Id, lockedVault.AmountIn)
			//assetOutTotal, _ := k.market.CalcAssetPrice(ctx, assetOut.Id, lockedVault.AmountOut)
			assetInTotalUint := lockedVault.AmountIn.Mul(sdk.NewIntFromUint64(assetInprice.Twa))
			assetOutTotalUint := lockedVault.AmountOut.Mul(sdk.NewIntFromUint64(assetOutprice.Twa))
			assetInTotal := sdk.NewDecFromInt(assetInTotalUint)
			assetOutTotal := sdk.NewDecFromInt(assetOutTotalUint)
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
			// calculations for finding selloff amount and liquidationDeductionAmount
			penalty := assetRatesStats.LiquidationPenalty.Add(assetRatesStats.LiquidationBonus)
			b := deductionPercentage.Add(penalty)
			totalIn := assetInTotal
			totalOut := assetOutTotal
			factor1 := c.Mul(totalIn)
			factor2 := b.Mul(c)
			numerator := totalOut.Sub(factor1)
			denominator := deductionPercentage.Sub(factor2)
			selloffAmount := numerator.Quo(denominator)

			// using this check as we want to deduct the liquidation penalty and auction bonus from the borrower position only once

			// TODO: revisit : DONE
			aip := sdk.NewDecFromInt(sdk.NewIntFromUint64(assetInprice.Twa))
			liquidationDeductionAmt := selloffAmount.Mul(penalty)
			liquidationDeductionAmount := liquidationDeductionAmt.Quo(aip)
			bonusToBidderAmt := selloffAmount.Mul(assetRatesStats.LiquidationBonus)

			bonusToBidderAmount := bonusToBidderAmt.Quo(aip)
			penaltyToReserveAmt := selloffAmount.Mul(assetRatesStats.LiquidationPenalty)
			penaltyToReserveAmount := penaltyToReserveAmt.Quo(aip)
			err = k.bank.SendCoinsFromModuleToModule(ctx, pool.ModuleName, auctiontypes.ModuleName, sdk.NewCoins(sdk.NewCoin(assetIn.Denom, sdk.NewInt(bonusToBidderAmount.TruncateInt64()))))
			if err != nil {
				return err
			}
			err = k.lend.UpdateReserveBalances(ctx, pair.AssetIn, pool.ModuleName, sdk.NewCoin(assetIn.Denom, sdk.NewInt(penaltyToReserveAmount.TruncateInt64())), true)
			if err != nil {
				return err
			}
			cAsset, _ := k.asset.GetAsset(ctx, assetRatesStats.CAssetID)
			// totalDeduction is the sum of liquidationDeductionAmount and selloffAmount
			sellOffAmt := selloffAmount.Quo(aip)
			totalDeduction := liquidationDeductionAmount.TruncateInt().Add(sellOffAmt.TruncateInt())
			updatedLockedVault.AmountIn = updatedLockedVault.AmountIn.Sub(totalDeduction)
			lendPos.AmountIn.Amount = lendPos.AmountIn.Amount.Sub(totalDeduction)
			if totalDeduction.GTE(updatedLockedVault.AmountIn) { // rare case only
				updatedLockedVault.AmountIn = sdk.ZeroInt()
				lendPos.AmountIn.Amount = sdk.ZeroInt()
			}
			// also global lend data is subtracted by totalDeduction amount
			assetStats, _ := k.lend.GetAssetStatsByPoolIDAndAssetID(ctx, lendPos.PoolID, lendPos.AssetID)
			assetStats.TotalLend = assetStats.TotalLend.Sub(totalDeduction)
			// setting the updated global lend data
			k.lend.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats)

			// users cToken present in pool's module will be burnt
			err = k.bank.BurnCoins(ctx, pool.ModuleName, sdk.NewCoins(sdk.NewCoin(cAsset.Denom, sdk.NewInt(penaltyToReserveAmount.TruncateInt64()))))
			if err != nil {
				return err
			}
			k.lend.SetLend(ctx, lendPos)

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
		// now the auction will be started from the auction module for the lockedVault

		err := k.auction.LendDutchActivator(ctx, updatedLockedVault)
		if err != nil {
			ctx.Logger().Error("error in dutch lend activator")
			return err
		}
	}
	return nil
}

func (k Keeper) UnLiquidateLockedBorrows(ctx sdk.Context, appID, id uint64, dutchAuction auctiontypes.DutchAuction) error {
	lockedVault, _ := k.GetLockedVault(ctx, appID, id)
	borrowMetadata := lockedVault.GetBorrowMetaData()
	if borrowMetadata != nil {
		lendPos, found := k.lend.GetLend(ctx, borrowMetadata.LendingId)
		if !found {
			return lendtypes.ErrLendNotFound
		}
		assetInPool, found := k.lend.GetPool(ctx, lendPos.PoolID)
		if !found {
			return lendtypes.ErrPoolNotFound
		}
		var firstTransitAssetID, secondTransitAssetID uint64
		for _, data := range assetInPool.AssetData {
			if data.AssetTransitType == 2 {
				firstTransitAssetID = data.AssetID
			}
			if data.AssetTransitType == 3 {
				secondTransitAssetID = data.AssetID
			}
		}
		firstBridgedAsset, _ := k.asset.GetAsset(ctx, firstTransitAssetID)
		userAddress, _ := sdk.AccAddressFromBech32(lockedVault.Owner)
		pair, _ := k.lend.GetLendPair(ctx, lockedVault.ExtendedPairId)
		assetStats, _ := k.lend.GetAssetRatesParams(ctx, pair.AssetIn)
		assetIn, _ := k.asset.GetAsset(ctx, pair.AssetIn)
		assetOut, _ := k.asset.GetAsset(ctx, pair.AssetOut)
		cAssetIn, _ := k.asset.GetAsset(ctx, assetStats.CAssetID)

		if lockedVault.IsAuctionComplete {
			if borrowMetadata.BridgedAssetAmount.IsZero() {
				// also calculate the current collaterlization ratio to ensure there is no sudden changes
				liqThreshold, _ := k.lend.GetAssetRatesParams(ctx, pair.AssetIn)
				unliquidatePointPercentage := liqThreshold.LiquidationThreshold

				if lockedVault.AmountOut.IsZero() {
					err := k.CreateLockedVaultHistory(ctx, lockedVault)
					if err != nil {
						return err
					}
					k.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
					if err = k.bank.SendCoinsFromModuleToAccount(ctx, assetInPool.ModuleName, userAddress, sdk.NewCoins(sdk.NewCoin(cAssetIn.Denom, lockedVault.AmountIn))); err != nil {
						return err
					}
				}
				newCalculatedCollateralizationRatio, _ := k.lend.CalculateCollateralizationRatio(ctx, lockedVault.AmountIn, assetIn, lockedVault.UpdatedAmountOut, assetOut)
				if newCalculatedCollateralizationRatio.GT(unliquidatePointPercentage) {
					updatedLockedVault := lockedVault
					updatedLockedVault.CurrentCollaterlisationRatio = newCalculatedCollateralizationRatio
					updatedLockedVault.SellOffHistory = append(updatedLockedVault.SellOffHistory, dutchAuction.String())
					k.SetLockedVault(ctx, updatedLockedVault)
					err := k.UpdateLockedBorrows(ctx, updatedLockedVault)
					if err != nil {
						ctx.Logger().Error("Error in UnLiquidateLockedBorrows first condition UpdateLockedBorrows in Liquidation, liquidate_borrow.go for ID %d", updatedLockedVault.LockedVaultId)
						return nil
					}
				}
				if newCalculatedCollateralizationRatio.LTE(unliquidatePointPercentage) {
					err := k.CreateLockedVaultHistory(ctx, lockedVault)
					if err != nil {
						return err
					}
					k.lend.CreteNewBorrow(ctx, lockedVault)
					k.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
				}
			} else {
				if borrowMetadata.BridgedAssetAmount.Denom == firstBridgedAsset.Denom {
					liqThresholdAssetIn, _ := k.lend.GetAssetRatesParams(ctx, pair.AssetIn)
					liqThresholdFirstBridgedAsset, _ := k.lend.GetAssetRatesParams(ctx, firstTransitAssetID)
					liqThreshold := liqThresholdAssetIn.LiquidationThreshold.Mul(liqThresholdFirstBridgedAsset.LiquidationThreshold)
					unliquidatePointPercentage := liqThreshold

					if lockedVault.AmountOut.IsZero() {
						err := k.CreateLockedVaultHistory(ctx, lockedVault)
						if err != nil {
							return err
						}
						k.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
						if err = k.bank.SendCoinsFromModuleToAccount(ctx, assetInPool.ModuleName, userAddress, sdk.NewCoins(sdk.NewCoin(cAssetIn.Denom, lockedVault.AmountIn))); err != nil {
							return err
						}
					}
					newCalculatedCollateralizationRatio, _ := k.lend.CalculateCollateralizationRatio(ctx, lockedVault.AmountIn, assetIn, lockedVault.UpdatedAmountOut, assetOut)
					if newCalculatedCollateralizationRatio.GT(unliquidatePointPercentage) {
						updatedLockedVault := lockedVault
						updatedLockedVault.CurrentCollaterlisationRatio = newCalculatedCollateralizationRatio
						updatedLockedVault.SellOffHistory = append(updatedLockedVault.SellOffHistory, dutchAuction.String())
						k.SetLockedVault(ctx, updatedLockedVault)
						err := k.UpdateLockedBorrows(ctx, updatedLockedVault)
						if err != nil {
							ctx.Logger().Error("Error in UnLiquidateLockedBorrows second condition UpdateLockedBorrows in Liquidation, liquidate_borrow.go for ID %d", updatedLockedVault.LockedVaultId)
							return nil
						}
					}
					if newCalculatedCollateralizationRatio.LTE(unliquidatePointPercentage) {
						err := k.CreateLockedVaultHistory(ctx, lockedVault)
						if err != nil {
							return err
						}
						k.lend.CreteNewBorrow(ctx, lockedVault)
						k.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
					}
				} else {
					liqThresholdAssetIn, _ := k.lend.GetAssetRatesParams(ctx, pair.AssetIn)
					liqThresholdSecondBridgedAsset, _ := k.lend.GetAssetRatesParams(ctx, secondTransitAssetID)
					liqThreshold := liqThresholdAssetIn.LiquidationThreshold.Mul(liqThresholdSecondBridgedAsset.LiquidationThreshold)
					unliquidatePointPercentage := liqThreshold

					if lockedVault.AmountOut.IsZero() {
						err := k.CreateLockedVaultHistory(ctx, lockedVault)
						if err != nil {
							return err
						}
						k.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
						if err = k.bank.SendCoinsFromModuleToAccount(ctx, assetInPool.ModuleName, userAddress, sdk.NewCoins(sdk.NewCoin(cAssetIn.Denom, lockedVault.AmountIn))); err != nil {
							return err
						}
					}
					newCalculatedCollateralizationRatio, _ := k.lend.CalculateCollateralizationRatio(ctx, lockedVault.AmountIn, assetIn, lockedVault.UpdatedAmountOut, assetOut)
					if newCalculatedCollateralizationRatio.GT(unliquidatePointPercentage) {
						updatedLockedVault := lockedVault
						updatedLockedVault.CurrentCollaterlisationRatio = newCalculatedCollateralizationRatio
						updatedLockedVault.SellOffHistory = append(updatedLockedVault.SellOffHistory, dutchAuction.String())
						k.SetLockedVault(ctx, updatedLockedVault)
						err := k.UpdateLockedBorrows(ctx, updatedLockedVault)
						if err != nil {
							ctx.Logger().Error("Error in UnLiquidateLockedBorrows third condition UpdateLockedBorrows in Liquidation, liquidate_borrow.go for ID %d", updatedLockedVault.LockedVaultId)
							return nil
						}
					}
					if newCalculatedCollateralizationRatio.LTE(unliquidatePointPercentage) {
						err := k.CreateLockedVaultHistory(ctx, lockedVault)
						if err != nil {
							return err
						}
						k.lend.CreteNewBorrow(ctx, lockedVault)
						k.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
					}
				}
			}
		}
	}
	return nil
}
