package keeper

import (
	"fmt"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	auctiontypes "github.com/comdex-official/comdex/x/auctionsV2/types"
	lendtypes "github.com/comdex-official/comdex/x/lend/types"

	utils "github.com/comdex-official/comdex/types"
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

			err := k.LiquidateIndividualVault(ctx, vault.Id)
			if err != nil {

				return fmt.Errorf(err.Error())
				//or maybe continue
			}

			return nil
		})
	}

	liquidationOffsetHolder.CurrentOffset = uint64(end)
	k.SetLiquidationOffsetHolder(ctx, types.VaultLiquidationsOffsetPrefix, liquidationOffsetHolder)

	return nil

}

func (k Keeper) LiquidateIndividualVault(ctx sdk.Context, vaultID uint64) error {

	vault, found := k.vault.GetVault(ctx, vaultID)
	if !found {
		return fmt.Errorf("Vault ID not found  %d", vault.Id)
	}

	//Checking ESM status and / or kill switch status
	esmStatus, found := k.esm.GetESMStatus(ctx, vault.AppId)
	klwsParams, _ := k.esm.GetKillSwitchData(ctx, vault.AppId)
	if (found && esmStatus.Status) || klwsParams.BreakerEnable {
		return fmt.Errorf("Kill Switch Or ESM is enabled For Liquidation,")
	}

	//Checking if app has enabled liquidations or not
	whitelistingData, found := k.GetAppIDByAppForLiquidation(ctx, vault.AppId)
	if !found {
		return fmt.Errorf("Liquidation not enabled for App ID  %d", vault.AppId)
	}

	// Checking extended pair vault data for Minimum collateralisation ratio
	extPair, _ := k.asset.GetPairsVault(ctx, vault.ExtendedPairVaultID)
	liqRatio := extPair.MinCr
	pair, _ := k.asset.GetPair(ctx, extPair.PairId)
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
		//Calling vault to use the updated values of the vault
		vault, _ := k.vault.GetVault(ctx, vault.Id)

		totalOut := vault.AmountOut.Add(vault.InterestAccumulated).Add(vault.ClosingFeeAccumulated)
		collateralizationRatio, err := k.vault.CalculateCollateralizationRatio(ctx, vault.ExtendedPairVaultID, vault.AmountIn, totalOut)
		if err != nil {
			return fmt.Errorf("error Calculating CR in Liquidation, liquidate_vaults.go for vaultID %d", vault.Id)
		}
		//Calculating Liquidation Fees
		feesToBeCollected := sdk.NewDecFromInt(totalOut).Mul(extPair.LiquidationPenalty).TruncateInt()

		//Calculating auction bonus to be given
		auctionBonusToBeGiven := sdk.NewDecFromInt(totalOut).Mul(extPair.LiquidationPenalty).TruncateInt()

		//Checking if the vault getting liquidated is a cmst based vault or not
		//This is primarily to infer that primary market will consider cmst at $1 at the time of buying it
		isCMST := !extPair.AssetOutOraclePrice

		//Creating locked vault struct , which will trigger auction
		err = k.CreateLockedVault(ctx, vault.Id, vault.ExtendedPairVaultID, vault.Owner, k.ReturnCoin(ctx, pair.AssetIn, vault.AmountIn), k.ReturnCoin(ctx, pair.AssetOut, totalOut), k.ReturnCoin(ctx, pair.AssetIn, vault.AmountIn), k.ReturnCoin(ctx, pair.AssetOut, totalOut), collateralizationRatio, vault.AppId, false, false, "", "", feesToBeCollected, auctionBonusToBeGiven, "vault", whitelistingData.AuctionType, isCMST, extPair.PairId)
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
}
func (k Keeper) ReturnCoin(ctx sdk.Context, assetID uint64, amount sdk.Int) sdk.Coin {
	asset, _ := k.asset.GetAsset(ctx, assetID)
	return sdk.NewCoin(asset.Denom, amount)
}

func (k Keeper) CreateLockedVault(ctx sdk.Context, OriginalVaultId, ExtendedPairId uint64, Owner string, AmountIn, AmountOut, CollateralToBeAuctioned, TargetDebt sdk.Coin, collateralizationRatio sdk.Dec, appID uint64, isInternalKeeper bool, isExternalKeeper bool, internalKeeperAddress string, externalKeeperAddress string, feesToBeCollected sdk.Int, bonusToBeGiven sdk.Int, initiatorType string, auctionType bool, isDebtCmst bool, pairId unit64) error {
	lockedVaultID := k.GetLockedVaultID(ctx)

	value := types.LockedVault{
		LockedVaultId:                lockedVaultID + 1,
		AppId:                        appID,
		OriginalVaultId:              OriginalVaultId,
		ExtendedPairId:               ExtendedPairId,
		Owner:                        Owner,
		CollateralToken:              AmountIn,
		DebtToken:                    AmountOut,//just a representation of the total debt the vault had incurred at the time of liquidation. // Target debt is a correct measure of what will get collected in the auction from bidders.
		CurrentCollaterlisationRatio: collateralizationRatio,
		CollateralToBeAuctioned:      AmountIn,
		TargetDebt:                   AmountOut.Add(sdk.NewCoin(AmountOut.Denom,feesToBeCollected)).Add(sdk.NewCoin(AmountOut.Denom,bonusToBeGiven)),//to add debt+liquidation+auction bonus here---- 
		LiquidationTimestamp:         ctx.BlockTime(),
		FeeToBeCollected:             feesToBeCollected,//just for calculation purpose
		BonusToBeGiven:               bonusToBeGiven,//just for calculation purpose
		IsInternalKeeper:             isInternalKeeper,
		InternalKeeperAddress:        internalKeeperAddress,
		IsExternalKeeper:             isExternalKeeper,
		ExternalKeeperAddress:        externalKeeperAddress,
		InitiatorType:                initiatorType,
		AuctionType:                  auctionType,
		IsDebtCmst:                   isDebtCmst,
		PairId:                       pairId,
	}
	//To understand a condition in which case target debt becomes equal to dollar value of collateral token 
	//at some point in the auction
	//1. what happens in that case
	//2. what if the bid on the auction makes the auction lossy, 
	//should be use the liquidation penalty ? most probably yes to cover the difference.
	//what if then liquidation penalty still falls short, should we then reduce the auction bonus from the debt , to make things even?
	//will this be enough to make sure auction does not not gets bid due to collateral not being able to cover the debt?
	//can a case occur in which liquidation penalty and auction bonus are still not enough?


	k.SetLockedVault(ctx, value)
	k.SetLockedVaultID(ctx, value.LockedVaultId)
	//Call auction activator
	err := k.auctionsV2.AuctionActivator(ctx, value)
	if err != nil {
		return fmt.Errorf("Auction could not be initiated for %d %d", value, err)
	}
	//struct for auction will stay same for english and dutch
	// based on type recieved from

	return nil
}

func (k Keeper) LiquidateBorrows(ctx sdk.Context) error {
	borrows, found := k.lend.GetBorrows(ctx)
	params := k.GetParams(ctx)
	if !found {
		ctx.Logger().Error("Params Not Found in Liquidation, liquidate_borrow.go")
		return nil
	}
	liquidationOffsetHolder, found := k.GetLiquidationOffsetHolder(ctx, types.VaultLiquidationsOffsetPrefix)
	if !found {
		liquidationOffsetHolder = types.NewLiquidationOffsetHolder(0)
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
					err = k.UpdateLockedBorrows(ctx, pool, borrowPos, lendPos.Owner, liqThreshold, assetIn, assetOut, firstBridgedAsset, lendPos.AppID, currentCollateralizationRatio)
					if err != nil {
						return fmt.Errorf("error in first condition UpdateLockedBorrows in UpdateLockedBorrows , liquidate_borrow.go for ID ")
					}
				}
			} else {
				if borrowPos.BridgedAssetAmount.Denom == firstBridgedAsset.Denom {
					currentCollateralizationRatio, err = k.lend.CalculateCollateralizationRatio(ctx, borrowPos.AmountIn.Amount, assetIn, borrowPos.AmountOut.Amount.Add(borrowPos.InterestAccumulated.TruncateInt()), assetOut)
					if err != nil {
						return err
					}
					if sdk.Dec.GT(currentCollateralizationRatio, liqThreshold.LiquidationThreshold.Mul(liqThresholdBridgedAssetOne.LiquidationThreshold)) {
						err = k.UpdateLockedBorrows(ctx, pool, borrowPos, lendPos.Owner, liqThreshold, assetIn, assetOut, firstBridgedAsset, lendPos.AppID, currentCollateralizationRatio)
						if err != nil {
							return fmt.Errorf("error in second condition UpdateLockedBorrows in UpdateLockedBorrows, liquidate_borrow.go for ID ")
						}
					}
				} else {
					currentCollateralizationRatio, err = k.lend.CalculateCollateralizationRatio(ctx, borrowPos.AmountIn.Amount, assetIn, borrowPos.AmountOut.Amount.Add(borrowPos.InterestAccumulated.TruncateInt()), assetOut)
					if err != nil {
						return err
					}

					if sdk.Dec.GT(currentCollateralizationRatio, liqThreshold.LiquidationThreshold.Mul(liqThresholdBridgedAssetTwo.LiquidationThreshold)) {
						err = k.UpdateLockedBorrows(ctx, pool, borrowPos, lendPos.Owner, liqThreshold, assetIn, assetOut, firstBridgedAsset, lendPos.AppID, currentCollateralizationRatio)
						if err != nil {
							return fmt.Errorf("error in third condition UpdateLockedBorrows in UpdateLockedBorrows, liquidate_borrow.go for ID ")
						}
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

func (k Keeper) UpdateLockedBorrows(ctx sdk.Context, pool lendtypes.Pool, borrow lendtypes.BorrowAsset, owner string, assetRatesStats lendtypes.AssetRatesParams, assetIn, assetOut, firstBridgeAsset assettypes.Asset, appID uint64, currentCollateralizationRatio sdk.Dec) error {
	firstBridgeAssetStats, _ := k.lend.GetAssetRatesParams(ctx, firstBridgeAsset.Id)
	secondBridgeAssetStats, _ := k.lend.GetAssetRatesParams(ctx, firstBridgeAsset.Id)

	assetInTotal, _ := k.market.CalcAssetPrice(ctx, assetIn.Id, borrow.AmountIn.Amount)
	assetOutTotal, _ := k.market.CalcAssetPrice(ctx, assetOut.Id, borrow.AmountOut.Amount)

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
	aop, _ := k.market.CalcAssetPrice(ctx, assetOut.Id, sdk.OneInt())
	bonusToBidderAmount := (selloffAmount.Mul(assetRatesStats.LiquidationBonus)).Quo(aop)
	penaltyToReserveAmount := (selloffAmount.Mul(assetRatesStats.LiquidationPenalty)).Quo(aop)
	sellOffAmt := selloffAmount.Quo(aip)
	//TODO: sellOffAmt At oracle price currently
	err := k.bank.SendCoinsFromModuleToModule(ctx, pool.ModuleName, auctiontypes.ModuleName, sdk.NewCoins(sdk.NewCoin(assetIn.Denom, sellOffAmt.TruncateInt())))
	if err != nil {
		return err
	}
	borrow.IsLiquidated = true
	k.lend.SetBorrow(ctx, borrow)
	//updatedLockedVault.CollateralToBeAuctioned = selloffAmount.TruncateInt()
	err = k.CreateLockedVault(ctx, borrow.ID, borrow.PairID, owner, borrow.AmountIn, borrow.AmountOut, borrow.AmountIn, sdk.NewCoin(assetIn.Denom, sellOffAmt.TruncateInt()), currentCollateralizationRatio, appID, false, false, "", "", bonusToBidderAmount.Add(penaltyToReserveAmount).TruncateInt(), "", true)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) LiquidateIndividualBorrow(ctx sdk.Context, borrowID uint64) error {
	borrowPos, found := k.lend.GetBorrow(ctx, borrowID)
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
			err = k.UpdateLockedBorrows(ctx, pool, borrowPos, lendPos.Owner, liqThreshold, assetIn, assetOut, firstBridgedAsset, lendPos.AppID, currentCollateralizationRatio)
			if err != nil {
				return fmt.Errorf("error in first condition UpdateLockedBorrows in UpdateLockedBorrows , liquidate_borrow.go for ID ")
			}
		}
	} else {
		if borrowPos.BridgedAssetAmount.Denom == firstBridgedAsset.Denom {
			currentCollateralizationRatio, err = k.lend.CalculateCollateralizationRatio(ctx, borrowPos.AmountIn.Amount, assetIn, borrowPos.AmountOut.Amount.Add(borrowPos.InterestAccumulated.TruncateInt()), assetOut)
			if err != nil {
				return err
			}
			if sdk.Dec.GT(currentCollateralizationRatio, liqThreshold.LiquidationThreshold.Mul(liqThresholdBridgedAssetOne.LiquidationThreshold)) {
				err = k.UpdateLockedBorrows(ctx, pool, borrowPos, lendPos.Owner, liqThreshold, assetIn, assetOut, firstBridgedAsset, lendPos.AppID, currentCollateralizationRatio)
				if err != nil {
					return fmt.Errorf("error in second condition UpdateLockedBorrows in UpdateLockedBorrows, liquidate_borrow.go for ID ")
				}
			}
		} else {
			currentCollateralizationRatio, err = k.lend.CalculateCollateralizationRatio(ctx, borrowPos.AmountIn.Amount, assetIn, borrowPos.AmountOut.Amount.Add(borrowPos.InterestAccumulated.TruncateInt()), assetOut)
			if err != nil {
				return err
			}

			if sdk.Dec.GT(currentCollateralizationRatio, liqThreshold.LiquidationThreshold.Mul(liqThresholdBridgedAssetTwo.LiquidationThreshold)) {
				err = k.UpdateLockedBorrows(ctx, pool, borrowPos, lendPos.Owner, liqThreshold, assetIn, assetOut, firstBridgedAsset, lendPos.AppID, currentCollateralizationRatio)
				if err != nil {
					return fmt.Errorf("error in third condition UpdateLockedBorrows in UpdateLockedBorrows, liquidate_borrow.go for ID ")
				}
			}
		}
	}
	return nil
}

func (k Keeper) MsgLiquidate(ctx sdk.Context, liquidator string, liqType, id uint64) error {
	if liqType == 0 {
		err := k.LiquidateIndividualVault(ctx, id)
		if err != nil {
			return err
		}
	} else if liqType == 1 {
		err := k.LiquidateIndividualBorrow(ctx, id)
		if err != nil {
			return err
		}
	} else {
		// TODO: for other apps
	}
	// TODO: send liquidation bonus to liquidator address logic
	return nil
}

func (k Keeper) SetLiquidationWhiteListing(ctx sdk.Context, liquidationWhiteListing types.LiquidationWhiteListing) {
	var (
		store = k.Store(ctx)
		key   = types.LiquidationWhiteListingKey(liquidationWhiteListing.AppId)
		value = k.cdc.MustMarshal(&liquidationWhiteListing)
	)

	store.Set(key, value)
}

func (k Keeper) GetLiquidationWhiteListing(ctx sdk.Context, appId uint64) (liquidationWhiteListing types.LiquidationWhiteListing, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LiquidationWhiteListingKey(appId)
		value = store.Get(key)
	)

	if value == nil {
		return liquidationWhiteListing, false
	}

	k.cdc.MustUnmarshal(value, &liquidationWhiteListing)
	return liquidationWhiteListing, true
}
