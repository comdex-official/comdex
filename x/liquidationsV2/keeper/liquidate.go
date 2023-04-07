package keeper

import (
	"fmt"

	utils "github.com/comdex-official/comdex/types"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	"github.com/comdex-official/comdex/x/liquidationsV2/types"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) Liquidate(ctx sdk.Context) error {

	err := k.LiquidateVaults(ctx)
	if err != nil {
		return err
	}

	err = k.LiquidateBorrows(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Liquidate Vaults function can liquidate all vaults created using the vault module.
//All vauts are looped and check if their underlying app has enabled liquidations.

func (k Keeper) LiquidateVaults(ctx sdk.Context) error {
	params := k.GetParams(ctx)

	//This allows us to loop over a slice of vaults per block , which doesnt stresses the abci.
	//Eg: if there exists 1,000,000 vaults  and the batch size is 100,000. then at every block 100,000 vaults will be looped and it will take
	//a total of 10 blocks to loop over all vaults.
	liquidationOffsetHolder, found := k.GetLiquidationOffsetHolder(ctx, types.VaultLiquidationsOffsetPrefix)
	if !found {
		liquidationOffsetHolder = types.NewLiquidationOffsetHolder(0)
	}
	// Fetching all  vaults
	totalVaults := k.vault.GetVaults(ctx)
	// Getting length of all vaults
	lengthOfVaults := int(k.vault.GetLengthOfVault(ctx))
	// Creating start and end slice
	start, end := types.GetSliceStartEndForLiquidations(lengthOfVaults, int(liquidationOffsetHolder.CurrentOffset), int(params.LiquidationBatchSize))
	if start == end {
		liquidationOffsetHolder.CurrentOffset = 0
		start, end = types.GetSliceStartEndForLiquidations(lengthOfVaults, int(liquidationOffsetHolder.CurrentOffset), int(params.LiquidationBatchSize))
	}
	newVaults := totalVaults[start:end]
	for _, vault := range newVaults {
		_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {

			//Checking ESM status and / or kill switch status
			esmStatus, found := k.esm.GetESMStatus(ctx, vault.AppId)
			klwsParams, _ := k.esm.GetKillSwitchData(ctx, vault.AppId)
			if (found && esmStatus.Status) || klwsParams.BreakerEnable {
				ctx.Logger().Error("Kill Switch Or ESM is enabled For Liquidation, liquidate_vaults.go for AppID %d", vault.AppId)
				continue
			}

			//Checking if app has enabled liquidations or not
			_, found = k.GetAppIDByAppForLiquidation(ctx, vault.AppId)

			if !found {
				return fmt.Errorf("Liquidation not enabled for App ID  %d", vault.AppId)
			}

			// Checking extended pair vault data for Minimum collateralisation ratio
			extPair, _ := k.asset.GetPairsVault(ctx, vault.ExtendedPairVaultID)
			liqRatio := extPair.MinCr
			totalOut := vault.AmountOut.Add(vault.InterestAccumulated).Add(vault.ClosingFeeAccumulated)
			collateralizationRatio, err := k.vault.CalculateCollateralizationRatio(ctx, vault.ExtendedPairVaultID, vault.AmountIn, totalOut)
			if err != nil {
				return fmt.Errorf("error Calculating CR in Liquidation, liquidate_vaults.go for vault ID %d", vault.Id)
			}
			if collateralizationRatio.LT(liqRatio) {
				totalDebt := vault.AmountOut.Add(vault.InterestAccumulated)
				err1 := k.rewards.CalculateVaultInterest(ctx, vault.AppId, vault.ExtendedPairVaultID, vault.Id, totalDebt, vault.BlockHeight, vault.BlockTime.Unix())
				if err1 != nil {
					return fmt.Errorf("error Calculating vault interest in Liquidation, liquidate_vaults.go for vaultID %d", vault.Id)
				}
				//Callling vault to use the updated values of the vault
				vault, _ := k.vault.GetVault(ctx, vault.Id)

				totalOut := vault.AmountOut.Add(vault.InterestAccumulated).Add(vault.ClosingFeeAccumulated)
				collateralizationRatio, err := k.vault.CalculateCollateralizationRatio(ctx, vault.ExtendedPairVaultID, vault.AmountIn, totalOut)
				if err != nil {
					return fmt.Errorf("error Calculating CR in Liquidation, liquidate_vaults.go for vaultID %d", vault.Id)
				}

				//Creating locked vault struct , which will trigger auction
				err = k.CreateLockedVault(ctx, vault.Id, vault.ExtendedPairVaultID, vault.Owner, vault.AmountIn, totalOut, collateralizationRatio, vault.AppId, false, false, "", "")
				if err != nil {
					return fmt.Errorf("error Creating Locked Vaults in Liquidation, liquidate_vaults.go for Vault %d", vault.Id)
				}
				length := k.vault.GetLengthOfVault(ctx)
				k.vault.SetLengthOfVault(ctx, length-1)

				//Removing data from existing structs
				k.vault.DeleteVault(ctx, vault.Id)
				var rewards rewardstypes.VaultInterestTracker
				rewards.AppMappingId = vault.AppId
				rewards.VaultId = vault.Id
				k.rewards.DeleteVaultInterestTracker(ctx, rewards)
				k.vault.DeleteAddressFromAppExtendedPairVaultMapping(ctx, vault.ExtendedPairVaultID, vault.Id, vault.AppId)
			}
			return nil
		})
	}

	liquidationOffsetHolder.CurrentOffset = uint64(end)
	k.SetLiquidationOffsetHolder(ctx, types.VaultLiquidationsOffsetPrefix, liquidationOffsetHolder)

	return nil

}

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
			pool, _ := k.lend.GetPool(ctx, lendPos.PoolID)
			assetIn, _ := k.asset.GetAsset(ctx, lendPair.AssetIn)
			assetOut, _ := k.asset.GetAsset(ctx, lendPair.AssetOut)
			liqThreshold, _ := k.lend.GetAssetRatesParams(ctx, lendPair.AssetIn)
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

			liqThresholdBridgedAssetOne, _ := k.lend.GetAssetRatesParams(ctx, firstTransitAssetID)
			liqThresholdBridgedAssetTwo, _ := k.lend.GetAssetRatesParams(ctx, secondTransitAssetID)
			firstBridgedAsset, _ := k.asset.GetAsset(ctx, firstTransitAssetID)
			secondBridgedAsset, _ := k.asset.GetAsset(ctx, secondTransitAssetID)

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
					err = k.UpdateLockedBorrows(ctx, lockedVault)
					if err != nil {
						return fmt.Errorf("error in first condition UpdateLockedBorrows in UpdateLockedBorrows , liquidate_borrow.go for ID %d", lockedVault.LockedVaultId)
					}
					k.lend.UpdateBorrowStats(ctx, lendPair, borrowPos.IsStableBorrow, borrowPos.AmountOut.Amount, false)
				}
			} else {
				if borrowPos.BridgedAssetAmount.Denom == firstBridgedAsset.Denom {
					currentCollateralizationRatio, err = k.lend.CalculateCollateralizationRatio(ctx, borrowPos.AmountIn.Amount, assetIn, borrowPos.AmountOut.Amount.Add(borrowPos.InterestAccumulated.TruncateInt()), assetOut)
					if err != nil {
						return err
					}
					if sdk.Dec.GT(currentCollateralizationRatio, liqThreshold.LiquidationThreshold.Mul(liqThresholdBridgedAssetOne.LiquidationThreshold)) {
						err = k.UpdateLockedBorrows(ctx, lockedVault)
						if err != nil {
							return fmt.Errorf("error in second condition UpdateLockedBorrows in UpdateLockedBorrows, liquidate_borrow.go for ID %d", lockedVault.LockedVaultId)
						}
						k.lend.UpdateBorrowStats(ctx, lendPair, borrowPos.IsStableBorrow, borrowPos.AmountOut.Amount, false)
					}
				} else {
					currentCollateralizationRatio, err = k.lend.CalculateCollateralizationRatio(ctx, borrowPos.AmountIn.Amount, assetIn, borrowPos.AmountOut.Amount.Add(borrowPos.InterestAccumulated.TruncateInt()), assetOut)
					if err != nil {
						return err
					}

					if sdk.Dec.GT(currentCollateralizationRatio, liqThreshold.LiquidationThreshold.Mul(liqThresholdBridgedAssetTwo.LiquidationThreshold)) {
						err = k.UpdateLockedBorrows(ctx, lockedVault, lendPair)
						if err != nil {
							return fmt.Errorf("error in third condition UpdateLockedBorrows in UpdateLockedBorrows, liquidate_borrow.go for ID %d", lockedVault.LockedVaultId)
						}
						k.lend.UpdateBorrowStats(ctx, lendPair, borrowPos.IsStableBorrow, borrowPos.AmountOut.Amount, false)
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

func (k Keeper) UpdateLockedBorrows(ctx sdk.Context, updatedLockedVault types.LockedVault, lendPos lendtypes.LendAsset, pool lendtypes.Pool, borrow lendtypes.BorrowAsset, assetRatesStats lendtypes.AssetRatesParams, assetIn, assetOut, firstBridgeAsset assettypes.Asset) error {
	firstBridgeAssetStats, _ := k.lend.GetAssetRatesParams(ctx, firstBridgeAsset.Id)
	secondBridgeAssetStats, _ := k.lend.GetAssetRatesParams(ctx, firstBridgeAsset.Id)

	assetInTotal, _ := k.market.CalcAssetPrice(ctx, assetIn.Id, updatedLockedVault.AmountIn)
	assetOutTotal, _ := k.market.CalcAssetPrice(ctx, assetOut.Id, updatedLockedVault.AmountOut)

	deductionPercentage, _ := sdk.NewDecFromStr("1.0")

	var c sdk.Dec
	if !borrow.BridgedAssetAmount.Amount.Equal(sdk.ZeroInt()) {
		if borrow.BridgedAssetAmount.Denom == firstBridgeAsset.Denom {
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
	err := k.bank.SendCoinsFromModuleToModule(ctx, pool.ModuleName, auctiontypes.ModuleName, sdk.NewCoins(sdk.NewCoin(assetIn.Denom, bonusToBidderAmount.Add(sellOffAmt).TruncateInt())))
	if err != nil {
		return err
	}
	err = k.lend.UpdateReserveBalances(ctx, assetIn.Id, pool.ModuleName, sdk.NewCoin(assetIn.Denom, penaltyToReserveAmount.TruncateInt()), true)
	if err != nil {
		return err
	}

	cAsset, _ := k.asset.GetAsset(ctx, assetRatesStats.CAssetID)
	// totalDeduction is the sum of liquidationDeductionAmount and selloffAmount
	totalDeduction := liquidationDeductionAmount.Add(sellOffAmt).TruncateInt() // Total deduction from amountIn also reduce to lend Position amountIn
	borrow.IsLiquidated = true
	if totalDeduction.GTE(updatedLockedVault.AmountIn) { // rare case only
		lendPos.AmountIn.Amount = lendPos.AmountIn.Amount.Sub(updatedLockedVault.AmountIn)
		// also global lend data is subtracted by totalDeduction amount
		assetStats, _ := k.lend.GetAssetStatsByPoolIDAndAssetID(ctx, lendPos.PoolID, lendPos.AssetID)
		assetStats.TotalLend = assetStats.TotalLend.Sub(updatedLockedVault.AmountIn)
		// setting the updated global lend data
		k.lend.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats)
		updatedLockedVault.AmountIn = sdk.ZeroInt()
		borrow.AmountIn.Amount = sdk.ZeroInt()
	} else {
		updatedLockedVault.AmountIn = updatedLockedVault.AmountIn.Sub(totalDeduction)
		lendPos.AmountIn.Amount = lendPos.AmountIn.Amount.Sub(totalDeduction)
		// also global lend data is subtracted by totalDeduction amount
		assetStats, _ := k.lend.GetAssetStatsByPoolIDAndAssetID(ctx, lendPos.PoolID, lendPos.AssetID)
		assetStats.TotalLend = assetStats.TotalLend.Sub(totalDeduction)
		// setting the updated global lend data
		k.lend.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats)
		borrow.AmountIn.Amount = borrow.AmountIn.Amount.Sub(totalDeduction)
	}

	// users cToken present in pool's module will be burnt
	// update borrow position
	// update lend position
	err = k.bank.BurnCoins(ctx, pool.ModuleName, sdk.NewCoins(sdk.NewCoin(cAsset.Denom, totalDeduction)))
	if err != nil {
		return err
	}
	k.lend.SetLend(ctx, lendPos)
	k.lend.SetBorrow(ctx, borrow)
	updatedLockedVault.CollateralToBeAuctioned = selloffAmount
	k.SetLockedVault(ctx, updatedLockedVault)
	k.SetLockedVaultID(ctx, updatedLockedVault.LockedVaultId)

	return nil
}

func (k Keeper) CreateLockedVault(ctx sdk.Context, OriginalVaultId, ExtendedPairId uint64, Owner string, AmountIn sdk.Int, AmountOut sdk.Int, collateralizationRatio sdk.Dec, appID uint64, isInternalKeeper bool, isExternalKeeper bool, internalKeeperAddress string, externalKeeperAddress string) error {
	lockedVaultID := k.GetLockedVaultID(ctx)

	value := types.LockedVault{
		LockedVaultId:                lockedVaultID + 1,
		AppId:                        appID,
		OriginalVaultId:              OriginalVaultId,
		ExtendedPairId:               ExtendedPairId,
		Owner:                        Owner,
		AmountIn:                     AmountIn,
		AmountOut:                    AmountOut,
		CurrentCollaterlisationRatio: collateralizationRatio,
		CollateralToBeAuctioned:      AmountIn,
		TargetDebt:                   AmountOut,
		LiquidationTimestamp:         ctx.BlockTime(),
		IsInternalKeeper:             false,
		InternalKeeperAddress:        "",
		IsExternalKeeper:             "",
		ExternalKeeperAddress:        "",
	
	}

	k.SetLockedVault(ctx, value)
	k.SetLockedVaultID(ctx, value.LockedVaultId)

	return nil
}
