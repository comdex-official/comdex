package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	utils "github.com/comdex-official/comdex/types"
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
		_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
			borrowPos, found := k.lend.GetBorrow(ctx, newBorrowIDs[l])
			if !found {
				return nil
			}
			if borrowPos.IsLiquidated {
				return nil
			}
			lendPair, _ := k.lend.GetLendPair(ctx, borrowPos.PairID)
			lendPos, found := k.lend.GetLend(ctx, borrowPos.LendingID)
			if !found {
				return fmt.Errorf("lend Pos Not Found in Liquidation, liquidate_borrow.go for ID %d", borrowPos.LendingID)
			}
			killSwitchParams, _ := k.esm.GetKillSwitchData(ctx, lendPos.AppID)
			if killSwitchParams.BreakerEnable {
				return fmt.Errorf("kill Switch is enabled in Liquidation, liquidate_borrow.go for ID %d", lendPos.AppID)
			}
			// calculating and updating the interest accumulated before checking for liquidations
			borrowPos, err := k.lend.CalculateBorrowInterestForLiquidation(ctx, borrowPos.ID)
			if err != nil {
				return fmt.Errorf("error in calculating Borrow Interest before liquidation")
			}
			if !borrowPos.StableBorrowRate.Equal(sdk.ZeroDec()) {
				borrowPos, err = k.lend.ReBalanceStableRates(ctx, borrowPos)
				if err != nil {
					return fmt.Errorf("error in re-balance stable rate check before liquidation")
				}
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
			pair, _ := k.lend.GetLendPair(ctx, borrowPos.PairID)
			// there are three possible cases
			// 	a. if borrow is from same pool
			//  b. if borrow is from first transit asset
			//  c. if borrow is from second transit asset
			if borrowPos.BridgedAssetAmount.Amount.Equal(sdk.ZeroInt()) { // first condition
				currentCollateralizationRatio, err = k.lend.CalculateCollateralizationRatio(ctx, borrowPos.AmountIn.Amount, assetIn, borrowPos.AmountOut.Amount.Add(borrowPos.InterestAccumulated.TruncateInt()), assetOut)
				if err != nil {
					return err
				}
				if sdk.Dec.GT(currentCollateralizationRatio, liqThreshold.LiquidationThreshold) {
					// after checking the currentCollateralizationRatio with LiquidationThreshold if borrow is to be liquidated then
					// CreateLockedBorrow function is called
					lockedVault, err := k.CreateLockedBorrow(ctx, borrowPos, currentCollateralizationRatio, lendPos.AppID)
					if err != nil {
						return fmt.Errorf("error in first condition CreateLockedBorrow in Liquidation, liquidate_borrow.go for ID %d", borrowPos.LendingID)
					}
					err = k.UpdateLockedBorrows(ctx, lockedVault)
					if err != nil {
						return fmt.Errorf("error in first condition UpdateLockedBorrows in UpdateLockedBorrows , liquidate_borrow.go for ID %d", lockedVault.LockedVaultId)
					}
					k.lend.UpdateBorrowStats(ctx, pair, borrowPos.IsStableBorrow, borrowPos.AmountOut.Amount, false)

				}
			} else {
				if borrowPos.BridgedAssetAmount.Denom == firstBridgedAsset.Denom {
					currentCollateralizationRatio, err = k.lend.CalculateCollateralizationRatio(ctx, borrowPos.AmountIn.Amount, assetIn, borrowPos.AmountOut.Amount.Add(borrowPos.InterestAccumulated.TruncateInt()), assetOut)
					if err != nil {
						return err
					}
					if sdk.Dec.GT(currentCollateralizationRatio, liqThreshold.LiquidationThreshold.Mul(liqThresholdBridgedAssetOne.LiquidationThreshold)) {
						lockedVault, err := k.CreateLockedBorrow(ctx, borrowPos, currentCollateralizationRatio, lendPos.AppID)
						if err != nil {
							return fmt.Errorf("error in second condition CreateLockedBorrow in Liquidation, liquidate_borrow.go for ID %d", borrowPos.LendingID)
						}
						err = k.UpdateLockedBorrows(ctx, lockedVault)
						if err != nil {
							return fmt.Errorf("error in second condition UpdateLockedBorrows in UpdateLockedBorrows, liquidate_borrow.go for ID %d", lockedVault.LockedVaultId)
						}
						k.lend.UpdateBorrowStats(ctx, pair, borrowPos.IsStableBorrow, borrowPos.AmountOut.Amount, false)

					}
				} else {
					currentCollateralizationRatio, err = k.lend.CalculateCollateralizationRatio(ctx, borrowPos.AmountIn.Amount, assetIn, borrowPos.AmountOut.Amount.Add(borrowPos.InterestAccumulated.TruncateInt()), assetOut)
					if err != nil {
						return err
					}

					if sdk.Dec.GT(currentCollateralizationRatio, liqThreshold.LiquidationThreshold.Mul(liqThresholdBridgedAssetTwo.LiquidationThreshold)) {
						lockedVault, err := k.CreateLockedBorrow(ctx, borrowPos, currentCollateralizationRatio, lendPos.AppID)
						if err != nil {
							return fmt.Errorf("error in third condition CreateLockedBorrow in Liquidation, liquidate_borrow.go for ID %d", borrowPos.LendingID)
						}
						err = k.UpdateLockedBorrows(ctx, lockedVault)
						if err != nil {
							return fmt.Errorf("error in third condition UpdateLockedBorrows in UpdateLockedBorrows, liquidate_borrow.go for ID %d", lockedVault.LockedVaultId)
						}
						k.lend.UpdateBorrowStats(ctx, pair, borrowPos.IsStableBorrow, borrowPos.AmountOut.Amount, false)

					}
				}
			}
			return nil
		})
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
	lockedVault := types.LockedVault{
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
	return lockedVault, nil
}

func (k Keeper) UpdateLockedBorrows(ctx sdk.Context, updatedLockedVault types.LockedVault) error {
	pair, _ := k.lend.GetLendPair(ctx, updatedLockedVault.ExtendedPairId)
	borrowMetaData := updatedLockedVault.GetBorrowMetaData()
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
		// Checking required flags
		if (!updatedLockedVault.IsAuctionInProgress && !updatedLockedVault.IsAuctionComplete) || (updatedLockedVault.IsAuctionComplete && updatedLockedVault.CurrentCollaterlisationRatio.GTE(unliquidatePointPercentage)) {
			assetIn, _ := k.asset.GetAsset(ctx, pair.AssetIn)
			assetOut, _ := k.asset.GetAsset(ctx, pair.AssetOut)
			collateralizationRatio, err := k.lend.CalculateCollateralizationRatio(ctx, updatedLockedVault.AmountIn, assetIn, updatedLockedVault.UpdatedAmountOut, assetOut)
			if err != nil {
				// ctx.Logger().Error("Error Calculating CR in Liquidation, liquidate_borrow.go for locked vault ID %d", lockedVault.LockedVaultId)
				return err
			}

			assetInTotal, _ := k.market.CalcAssetPrice(ctx, assetIn.Id, updatedLockedVault.AmountIn)
			assetOutTotal, _ := k.market.CalcAssetPrice(ctx, assetOut.Id, updatedLockedVault.UpdatedAmountOut)

			deductionPercentage, _ := sdk.NewDecFromStr("1.0")

			var c sdk.Dec
			if !borrowMetaData.BridgedAssetAmount.Amount.Equal(sdk.ZeroInt()) {
				if borrowMetaData.BridgedAssetAmount.Denom == firstBridgeAsset.Denom {
					c = assetRatesStats.Ltv.Mul(firstBridgeAssetStats.Ltv)
				} else {
					c = assetRatesStats.Ltv.Mul(secondBridgeAssetStats.Ltv)
				}
			} else {
				c = assetRatesStats.Ltv
			}
			// calculations for finding selloff amount and liquidationDeductionAmount
			b := deductionPercentage.Add(assetRatesStats.LiquidationPenalty.Add(assetRatesStats.LiquidationBonus))
			totalIn := assetInTotal
			totalOut := assetOutTotal
			factor1 := c.Mul(totalIn)
			factor2 := b.Mul(c)
			numerator := totalOut.Sub(factor1)
			denominator := deductionPercentage.Sub(factor2)
			selloffAmount := numerator.Quo(denominator) // Dollar Value
			aip, _ := k.market.CalcAssetPrice(ctx, assetIn.Id, sdk.OneInt())
			liquidationDeductionAmt := selloffAmount.Mul(assetRatesStats.LiquidationPenalty.Add(assetRatesStats.LiquidationBonus))
			liquidationDeductionAmount := liquidationDeductionAmt.Quo(aip) // To be subtracted from AmountIn along with sellOff amt

			bonusToBidderAmount := (selloffAmount.Mul(assetRatesStats.LiquidationBonus)).Quo(aip)
			penaltyToReserveAmount := (selloffAmount.Mul(assetRatesStats.LiquidationPenalty)).Quo(aip)
			sellOffAmt := selloffAmount.Quo(aip)
			err = k.bank.SendCoinsFromModuleToModule(ctx, pool.ModuleName, auctiontypes.ModuleName, sdk.NewCoins(sdk.NewCoin(assetIn.Denom, bonusToBidderAmount.Add(sellOffAmt).TruncateInt())))
			if err != nil {
				return err
			}
			err = k.lend.UpdateReserveBalances(ctx, pair.AssetIn, pool.ModuleName, sdk.NewCoin(assetIn.Denom, penaltyToReserveAmount.TruncateInt()), true)
			if err != nil {
				return err
			}

			allReserveStats, found := k.lend.GetAllReserveStatsByAssetID(ctx, pair.AssetIn)
			if !found {
				allReserveStats = lendtypes.AllReserveStats{
					AssetID:                        pair.AssetIn,
					AmountOutFromReserveToLenders:  sdk.ZeroInt(),
					AmountOutFromReserveForAuction: sdk.ZeroInt(),
					AmountInFromLiqPenalty:         sdk.ZeroInt(),
					AmountInFromRepayments:         sdk.ZeroInt(),
					TotalAmountOutToLenders:        sdk.ZeroInt(),
				}
			}
			allReserveStats.AmountInFromLiqPenalty = allReserveStats.AmountInFromLiqPenalty.Add(penaltyToReserveAmount.TruncateInt())
			k.lend.SetAllReserveStatsByAssetID(ctx, allReserveStats)

			cAsset, _ := k.asset.GetAsset(ctx, assetRatesStats.CAssetID)
			// totalDeduction is the sum of liquidationDeductionAmount and selloffAmount
			totalDeduction := liquidationDeductionAmount.Add(sellOffAmt).TruncateInt() // Total deduction from amountIn also reduce to lend Position amountIn
			borrowPos, _ := k.lend.GetBorrow(ctx, updatedLockedVault.OriginalVaultId)
			borrowPos.IsLiquidated = true
			if totalDeduction.GTE(updatedLockedVault.AmountIn) { // rare case only
				lendPos.AmountIn.Amount = lendPos.AmountIn.Amount.Sub(updatedLockedVault.AmountIn)
				// also global lend data is subtracted by totalDeduction amount
				assetStats, _ := k.lend.GetAssetStatsByPoolIDAndAssetID(ctx, lendPos.PoolID, lendPos.AssetID)
				assetStats.TotalLend = assetStats.TotalLend.Sub(updatedLockedVault.AmountIn)
				// setting the updated global lend data
				k.lend.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats)
				updatedLockedVault.AmountIn = sdk.ZeroInt()
				borrowPos.AmountIn.Amount = sdk.ZeroInt()
			} else {
				updatedLockedVault.AmountIn = updatedLockedVault.AmountIn.Sub(totalDeduction)
				lendPos.AmountIn.Amount = lendPos.AmountIn.Amount.Sub(totalDeduction)
				// also global lend data is subtracted by totalDeduction amount
				assetStats, _ := k.lend.GetAssetStatsByPoolIDAndAssetID(ctx, lendPos.PoolID, lendPos.AssetID)
				assetStats.TotalLend = assetStats.TotalLend.Sub(totalDeduction)
				// setting the updated global lend data
				k.lend.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats)
				borrowPos.AmountIn.Amount = borrowPos.AmountIn.Amount.Sub(totalDeduction)
			}

			// users cToken present in pool's module will be burnt
			// update borrow position
			// update lend position
			err = k.bank.BurnCoins(ctx, pool.ModuleName, sdk.NewCoins(sdk.NewCoin(cAsset.Denom, totalDeduction)))
			if err != nil {
				return err
			}
			k.lend.SetLend(ctx, lendPos)
			k.lend.SetBorrow(ctx, borrowPos)
			updatedLockedVault.CurrentCollaterlisationRatio = collateralizationRatio
			updatedLockedVault.CollateralToBeAuctioned = selloffAmount
			k.SetLockedVault(ctx, updatedLockedVault)
			k.SetLockedVaultID(ctx, updatedLockedVault.LockedVaultId)
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
		assetOutPool, found := k.lend.GetPool(ctx, pair.AssetOutPoolID)
		if !found {
			return lendtypes.ErrPoolNotFound
		}
		assetStats, _ := k.lend.GetAssetRatesParams(ctx, pair.AssetIn)
		assetIn, _ := k.asset.GetAsset(ctx, pair.AssetIn)
		assetOut, _ := k.asset.GetAsset(ctx, pair.AssetOut)
		cAssetIn, _ := k.asset.GetAsset(ctx, assetStats.CAssetID)

		if lockedVault.IsAuctionComplete {
			// clearing borrow interest from borrow position after auction
			poolAssetLBMappingData, _ := k.lend.GetAssetStatsByPoolIDAndAssetID(ctx, pair.AssetOutPoolID, pair.AssetOut)
			assetOutStats, _ := k.lend.GetAssetRatesParams(ctx, pair.AssetOut)
			cAsset, _ := k.asset.GetAsset(ctx, assetOutStats.CAssetID)
			reservePoolRecords, found := k.lend.GetBorrowInterestTracker(ctx, lockedVault.OriginalVaultId)
			if !found {
				reservePoolRecords = lendtypes.BorrowInterestTracker{
					BorrowingId:         lockedVault.OriginalVaultId,
					ReservePoolInterest: sdk.ZeroDec(),
				}
			}
			borrowPos, _ := k.lend.GetBorrow(ctx, lockedVault.OriginalVaultId)
			amtToReservePool := reservePoolRecords.ReservePoolInterest
			if amtToReservePool.TruncateInt().GT(sdk.ZeroInt()) {
				amount := sdk.NewCoin(assetOut.Denom, amtToReservePool.TruncateInt())
				err := k.lend.UpdateReserveBalances(ctx, pair.AssetOut, assetOutPool.ModuleName, amount, true)
				if err != nil {
					return err
				}
				allReserveStats, found := k.lend.GetAllReserveStatsByAssetID(ctx, pair.AssetOut)
				if !found {
					allReserveStats = lendtypes.AllReserveStats{
						AssetID:                        pair.AssetOut,
						AmountOutFromReserveToLenders:  sdk.ZeroInt(),
						AmountOutFromReserveForAuction: sdk.ZeroInt(),
						AmountInFromLiqPenalty:         sdk.ZeroInt(),
						AmountInFromRepayments:         sdk.ZeroInt(),
						TotalAmountOutToLenders:        sdk.ZeroInt(),
					}
				}
				allReserveStats.AmountInFromRepayments = allReserveStats.AmountInFromRepayments.Add(amount.Amount)
				k.lend.SetAllReserveStatsByAssetID(ctx, allReserveStats)
			}
			amtToMint := (borrowPos.InterestAccumulated.Sub(amtToReservePool)).TruncateInt()
			if amtToMint.GT(sdk.ZeroInt()) {
				err := k.bank.MintCoins(ctx, assetOutPool.ModuleName, sdk.NewCoins(sdk.NewCoin(cAsset.Denom, amtToMint)))
				if err != nil {
					return err
				}
				poolAssetLBMappingData.TotalInterestAccumulated = poolAssetLBMappingData.TotalInterestAccumulated.Add(amtToMint)
				k.lend.SetAssetStatsByPoolIDAndAssetID(ctx, poolAssetLBMappingData)
			}
			borrowPos.InterestAccumulated = borrowPos.InterestAccumulated.Sub(sdk.NewDecFromInt(borrowPos.InterestAccumulated.TruncateInt()))
			reservePoolRecords.ReservePoolInterest = reservePoolRecords.ReservePoolInterest.Sub(sdk.NewDecFromInt(amtToReservePool.TruncateInt())) // the decimal precision is maintained
			k.lend.SetBorrowInterestTracker(ctx, reservePoolRecords)
			k.lend.SetBorrow(ctx, borrowPos)

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
					k.lend.DeleteIDFromAssetStatsMapping(ctx, pair.AssetOutPoolID, pair.AssetOut, lockedVault.OriginalVaultId, false)
					k.lend.DeleteBorrowIDFromUserMapping(ctx, lendPos.Owner, lendPos.ID, lockedVault.OriginalVaultId)
					k.lend.DeleteBorrow(ctx, lockedVault.OriginalVaultId)
					k.lend.DeleteBorrowInterestTracker(ctx, lockedVault.OriginalVaultId)
					return nil
				}
				if lockedVault.AmountIn.IsZero() {
					k.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
					k.lend.DeleteIDFromAssetStatsMapping(ctx, pair.AssetOutPoolID, pair.AssetOut, lockedVault.OriginalVaultId, false)
					k.lend.DeleteBorrowIDFromUserMapping(ctx, lendPos.Owner, lendPos.ID, lockedVault.OriginalVaultId)
					k.lend.DeleteBorrow(ctx, lockedVault.OriginalVaultId)
					k.lend.DeleteBorrowInterestTracker(ctx, lockedVault.OriginalVaultId)
					return nil
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
						k.lend.DeleteIDFromAssetStatsMapping(ctx, pair.AssetOutPoolID, pair.AssetOut, lockedVault.OriginalVaultId, false)
						k.lend.DeleteBorrowIDFromUserMapping(ctx, lendPos.Owner, lendPos.ID, lockedVault.OriginalVaultId)
						k.lend.DeleteBorrow(ctx, lockedVault.OriginalVaultId)
						k.lend.DeleteBorrowInterestTracker(ctx, lockedVault.OriginalVaultId)
						return nil
					}
					if lockedVault.AmountIn.IsZero() {
						k.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
						k.lend.DeleteIDFromAssetStatsMapping(ctx, pair.AssetOutPoolID, pair.AssetOut, lockedVault.OriginalVaultId, false)
						k.lend.DeleteBorrowIDFromUserMapping(ctx, lendPos.Owner, lendPos.ID, lockedVault.OriginalVaultId)
						k.lend.DeleteBorrow(ctx, lockedVault.OriginalVaultId)
						k.lend.DeleteBorrowInterestTracker(ctx, lockedVault.OriginalVaultId)
						return nil
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
						k.lend.DeleteIDFromAssetStatsMapping(ctx, pair.AssetOutPoolID, pair.AssetOut, lockedVault.OriginalVaultId, false)
						k.lend.DeleteBorrowIDFromUserMapping(ctx, lendPos.Owner, lendPos.ID, lockedVault.OriginalVaultId)
						k.lend.DeleteBorrow(ctx, lockedVault.OriginalVaultId)
						k.lend.DeleteBorrowInterestTracker(ctx, lockedVault.OriginalVaultId)
						return nil
					}
					if lockedVault.AmountIn.IsZero() {
						k.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
						k.lend.DeleteIDFromAssetStatsMapping(ctx, pair.AssetOutPoolID, pair.AssetOut, lockedVault.OriginalVaultId, false)
						k.lend.DeleteBorrowIDFromUserMapping(ctx, lendPos.Owner, lendPos.ID, lockedVault.OriginalVaultId)
						k.lend.DeleteBorrow(ctx, lockedVault.OriginalVaultId)
						k.lend.DeleteBorrowInterestTracker(ctx, lockedVault.OriginalVaultId)
						return nil
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
