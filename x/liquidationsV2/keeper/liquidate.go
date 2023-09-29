package keeper

import (
	"fmt"
	utils "github.com/comdex-official/comdex/types"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	auctionsV2types "github.com/comdex-official/comdex/x/auctionsV2/types"
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	"github.com/comdex-official/comdex/x/liquidationsV2/types"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) Liquidate(ctx sdk.Context) error {

	err := k.LiquidateVaults(ctx, 0)
	if err != nil {
		return err
	}

	err = k.LiquidateBorrows(ctx, 1)
	if err != nil {
		return err
	}

	err = k.LiquidateForSurplusAndDebt(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Liquidate Vaults function can liquidate all vaults created using the vault module.
//All vaults are looped and check if their underlying app has enabled liquidations.

func (k Keeper) LiquidateVaults(ctx sdk.Context, offsetCounterId uint64) error {
	params := k.GetParams(ctx)

	//This allows us to loop over a slice of vaults per block , which doesnt stresses the abci.
	//Eg: if there exists 1,000,000 vaults  and the batch size is 100,000. then at every block 100,000 vaults will be looped and it will take
	//a total of 10 blocks to loop over all vaults.
	liquidationOffsetHolder, found := k.GetLiquidationOffsetHolder(ctx, types.VaultLiquidationsOffsetPrefix, offsetCounterId)
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

			err := k.LiquidateIndividualVault(ctx, vault.Id, "", false)
			if err != nil {

				return fmt.Errorf(err.Error())
				//or maybe continue
			}

			return nil
		})
	}

	liquidationOffsetHolder.CurrentOffset = uint64(end)
	liquidationOffsetHolder.AppId = offsetCounterId
	k.SetLiquidationOffsetHolder(ctx, types.VaultLiquidationsOffsetPrefix, liquidationOffsetHolder)

	return nil

}

func (k Keeper) LiquidateIndividualVault(ctx sdk.Context, vaultID uint64, liquidator string, isInternalkeeper bool) error {

	vault, found := k.vault.GetVault(ctx, vaultID)
	if !found {
		return fmt.Errorf("Vault ID not found  %d", vault.Id)
	}

	//Checking ESM status and / or kill switch status
	esmStatus, found := k.esm.GetESMStatus(ctx, vault.AppId)
	klwsParams, _ := k.esm.GetKillSwitchData(ctx, vault.AppId)
	if (found && esmStatus.Status) || klwsParams.BreakerEnable {
		return fmt.Errorf("kill Switch Or ESM is enabled For Liquidation")
	}

	//Checking if app has enabled liquidations or not
	whitelistingData, found := k.GetLiquidationWhiteListing(ctx, vault.AppId)
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
		vault, _ = k.vault.GetVault(ctx, vault.Id)

		totalOut = vault.AmountOut.Add(vault.InterestAccumulated).Add(vault.ClosingFeeAccumulated)
		collateralizationRatio, err = k.vault.CalculateCollateralizationRatio(ctx, vault.ExtendedPairVaultID, vault.AmountIn, totalOut)
		if err != nil {
			return fmt.Errorf("error Calculating CR in Liquidation, liquidate_vaults.go for vaultID %d", vault.Id)
		}
		//Calculating Liquidation Fees
		feesToBeCollected := sdk.NewDecFromInt(totalOut).Mul(extPair.LiquidationPenalty).TruncateInt()

		//Calculating auction bonus to be given
		auctionBonusToBeGiven := sdk.ZeroInt()

		//Checking if the vault getting liquidated is a cmst based vault or not
		//This is primarily to infer that primary market will consider cmst at $1 at the time of buying it
		isCMST := !extPair.AssetOutOraclePrice

		//Creating locked vault struct , which will trigger auction
		//This function will only trigger Dutch auction
		//before creating locked vault, checking that Dutch auction is already there in the whitelisted liquidation data
		if !whitelistingData.IsDutchActivated {
			return fmt.Errorf("error , dutch auction not activated by the app, this function is only to trigger dutch auctions")

		}

		if vault.AmountIn.GT(sdk.ZeroInt()) {
			err := k.bank.SendCoinsFromModuleToModule(ctx, vaulttypes.ModuleName, auctionsV2types.ModuleName, sdk.NewCoins(k.ReturnCoin(ctx, pair.AssetIn, vault.AmountIn)))
			if err != nil {
				return fmt.Errorf("error , not enough token in vault to transfer %s", vault.AmountIn)
			}
		}
		err = k.CreateLockedVault(ctx, vault.Id, vault.ExtendedPairVaultID, vault.Owner, k.ReturnCoin(ctx, pair.AssetIn, vault.AmountIn), k.ReturnCoin(ctx, pair.AssetOut, totalOut), k.ReturnCoin(ctx, pair.AssetIn, vault.AmountIn), k.ReturnCoin(ctx, pair.AssetOut, totalOut), collateralizationRatio, vault.AppId, isInternalkeeper, liquidator, "", feesToBeCollected, auctionBonusToBeGiven, "vault", whitelistingData.IsDutchActivated, isCMST, pair.AssetIn, pair.AssetOut)
		if err != nil {
			return fmt.Errorf("error Creating Locked Vaults in Liquidation, liquidate_vaults.go for Vault %d", vault.Id)
		}
		length := k.vault.GetLengthOfVault(ctx)
		k.vault.SetLengthOfVault(ctx, length-1)

		//Removing data from existing structs

		var rewards rewardstypes.VaultInterestTracker
		rewards.AppMappingId = vault.AppId
		rewards.VaultId = vault.Id
		k.rewards.DeleteVaultInterestTracker(ctx, rewards)
		k.vault.DeleteAddressFromAppExtendedPairVaultMapping(ctx, vault.ExtendedPairVaultID, vault.Id, vault.AppId)
		k.vault.DeleteUserVaultExtendedPairMapping(ctx, vault.Owner, vault.AppId, vault.ExtendedPairVaultID)
		k.vault.DeleteVault(ctx, vault.Id)

	}

	return nil
}

func (k Keeper) ReturnCoin(ctx sdk.Context, assetID uint64, amount sdk.Int) sdk.Coin {
	asset, _ := k.asset.GetAsset(ctx, assetID)
	return sdk.NewCoin(asset.Denom, amount)
}

func (k Keeper) CreateLockedVault(ctx sdk.Context, OriginalVaultId, ExtendedPairId uint64, Owner string, AmountIn, AmountOut, CollateralToBeAuctioned, TargetDebt sdk.Coin, collateralizationRatio sdk.Dec, appID uint64, isInternalKeeper bool, internalKeeperAddress string, externalKeeperAddress string, feesToBeCollected sdk.Int, bonusToBeGiven sdk.Int, initiatorType string, auctionType bool, isDebtCmst bool, collateralID uint64, DebtId uint64) error {
	lockedVaultID := k.GetLockedVaultID(ctx)

	value := types.LockedVault{
		LockedVaultId:                lockedVaultID + 1,
		AppId:                        appID,
		OriginalVaultId:              OriginalVaultId,
		ExtendedPairId:               ExtendedPairId,
		Owner:                        Owner,
		CollateralToken:              AmountIn,
		DebtToken:                    AmountOut, //just a representation of the total debt the vault had incurred at the time of liquidation. // Target debt is a correct measure of what will get collected in the auction from bidders.
		CurrentCollaterlisationRatio: collateralizationRatio,
		CollateralToBeAuctioned:      AmountIn,
		TargetDebt:                   AmountOut.Add(sdk.NewCoin(AmountOut.Denom, feesToBeCollected)), //to add debt+liquidation+auction bonus here----
		LiquidationTimestamp:         ctx.BlockTime(),
		FeeToBeCollected:             feesToBeCollected, //just for calculation purpose
		BonusToBeGiven:               bonusToBeGiven,    //just for calculation purpose
		IsInternalKeeper:             isInternalKeeper,
		InternalKeeperAddress:        internalKeeperAddress,
		ExternalKeeperAddress:        externalKeeperAddress,
		InitiatorType:                initiatorType,
		AuctionType:                  auctionType,
		IsDebtCmst:                   isDebtCmst,
		CollateralAssetId:            collateralID,
		DebtAssetId:                  DebtId,
	}
	//To understand a condition in which case target debt becomes equal to dollar value of collateral token
	//at some point in the auction
	//1. what happens in that case
	//2. what if the bid on the auction makes the auction lossy,
	//should be used the liquidation penalty ? most probably yes to cover the difference.
	//what if then liquidation penalty still falls short, should we then reduce the auction bonus from the debt , to make things even?
	//will this be enough to make sure auction does not get bid due to collateral not being able to cover the debt?
	//can a case occur in which liquidation penalty and auction bonus are still not enough?

	// if english auction check if it is enabled or not
	if !value.AuctionType {
		liquidationWhitelistingAppData, found := k.GetLiquidationWhiteListing(ctx, appID)
		if !found || !liquidationWhitelistingAppData.IsEnglishActivated {
			return fmt.Errorf("Auction could not be initiated for %s ", types.ErrEnglishAuctionDisabled)
		}
	}

	k.SetLockedVault(ctx, value)
	k.SetLockedVaultID(ctx, value.LockedVaultId)
	//Call auction activator
	err := k.auctionsV2.AuctionActivator(ctx, value)
	if err != nil {
		return fmt.Errorf("Auction could not be initiated for %d ", err)
	}
	//struct for auction will stay same for english and Dutch
	// based on type received from

	return nil
}

func (k Keeper) LiquidateBorrows(ctx sdk.Context, offsetCounterId uint64) error {
	borrows, found := k.lend.GetBorrows(ctx)
	params := k.GetParams(ctx)
	if !found {
		ctx.Logger().Error("Params Not Found in Liquidation")
		return nil
	}
	liquidationOffsetHolder, found := k.GetLiquidationOffsetHolder(ctx, types.VaultLiquidationsOffsetPrefix, offsetCounterId)
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
		err := k.LiquidateIndividualBorrow(ctx, newBorrowIDs[l], "", false)
		if err != nil {
			return err
		}
	}
	liquidationOffsetHolder.CurrentOffset = uint64(end)
	k.SetLiquidationOffsetHolder(ctx, types.VaultLiquidationsOffsetPrefix, liquidationOffsetHolder)

	return nil
}

func (k Keeper) LiquidateIndividualBorrow(ctx sdk.Context, borrowID uint64, liquidator string, isInternalkeeper bool) error {
	borrowPos, found := k.lend.GetBorrow(ctx, borrowID)
	if !found {
		return fmt.Errorf("vault ID not found %d", borrowID)
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
		return fmt.Errorf("error in calculating Borrow Interest before liquidation %d", err)
	}
	if !borrowPos.StableBorrowRate.Equal(sdk.ZeroDec()) {
		borrowPos, err = k.lend.ReBalanceStableRates(ctx, borrowPos)
		if err != nil {
			return fmt.Errorf("error in re-balance stable rate check before liquidation")
		}
	}

	LiquidationThreshold := liqThreshold.LiquidationThreshold
	if lendPair.IsEModeEnabled {
		LiquidationThreshold = liqThreshold.ELiquidationThreshold
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
		if sdk.Dec.GT(currentCollateralizationRatio, LiquidationThreshold) {
			err = k.UpdateLockedBorrows(ctx, borrowPos, lendPos.Owner, lendPos.AppID, currentCollateralizationRatio, liqThreshold, lendPair, pool, assetIn, liquidator, isInternalkeeper)
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
			if sdk.Dec.GT(currentCollateralizationRatio, LiquidationThreshold.Mul(liqThresholdBridgedAssetOne.LiquidationThreshold)) {
				err = k.UpdateLockedBorrows(ctx, borrowPos, lendPos.Owner, lendPos.AppID, currentCollateralizationRatio, liqThreshold, lendPair, pool, assetIn, liquidator, isInternalkeeper)
				if err != nil {
					return fmt.Errorf("error in second condition UpdateLockedBorrows in UpdateLockedBorrows, liquidate_borrow.go for ID ")
				}
			}
		} else {
			currentCollateralizationRatio, err = k.lend.CalculateCollateralizationRatio(ctx, borrowPos.AmountIn.Amount, assetIn, borrowPos.AmountOut.Amount.Add(borrowPos.InterestAccumulated.TruncateInt()), assetOut)
			if err != nil {
				return err
			}

			if sdk.Dec.GT(currentCollateralizationRatio, LiquidationThreshold.Mul(liqThresholdBridgedAssetTwo.LiquidationThreshold)) {
				err = k.UpdateLockedBorrows(ctx, borrowPos, lendPos.Owner, lendPos.AppID, currentCollateralizationRatio, liqThreshold, lendPair, pool, assetIn, liquidator, isInternalkeeper)
				if err != nil {
					return fmt.Errorf("error in third condition UpdateLockedBorrows in UpdateLockedBorrows, liquidate_borrow.go for ID ")
				}
			}
		}
	}
	return nil
}

func (k Keeper) UpdateLockedBorrows(ctx sdk.Context, borrow lendtypes.BorrowAsset, owner string, appID uint64, currentCollateralizationRatio sdk.Dec, assetRatesStats lendtypes.AssetRatesParams, lendPair lendtypes.Extended_Pair, pool lendtypes.Pool, assetIn assettypes.Asset, liquidator string, isInternalkeeper bool) error {
	lendPos, _ := k.lend.GetLend(ctx, borrow.LendingID)
	whitelistingData, found := k.GetLiquidationWhiteListing(ctx, appID)
	if !found {
		return fmt.Errorf("Liquidation not enabled for App ID  %d", appID)
	}
	borrow.IsLiquidated = true
	k.lend.SetBorrow(ctx, borrow)
	pair, _ := k.lend.GetLendPair(ctx, borrow.PairID)
	cAsset, _ := k.asset.GetAsset(ctx, assetRatesStats.CAssetID)
	//Calculating Liquidation Fees
	feesToBeCollected := sdk.NewDecFromInt(borrow.AmountOut.Amount).Mul(assetRatesStats.LiquidationPenalty).TruncateInt()

	//Calculating auction bonus to be given
	auctionBonusToBeGiven := sdk.NewDecFromInt(borrow.AmountOut.Amount).Mul(assetRatesStats.LiquidationBonus).TruncateInt()

	err := k.bank.SendCoinsFromModuleToModule(ctx, pool.ModuleName, auctionsV2types.ModuleName, sdk.NewCoins(sdk.NewCoin(assetIn.Denom, borrow.AmountIn.Amount)))
	if err != nil {
		return err
	}

	err = k.bank.BurnCoins(ctx, pool.ModuleName, sdk.NewCoins(sdk.NewCoin(cAsset.Denom, borrow.AmountIn.Amount)))
	if err != nil {
		return err
	}

	err = k.CreateLockedVault(ctx, borrow.ID, borrow.PairID, owner, sdk.NewCoin(assetIn.Denom, borrow.AmountIn.Amount), borrow.AmountOut, borrow.AmountIn, borrow.AmountOut, currentCollateralizationRatio, appID, isInternalkeeper, liquidator, "", feesToBeCollected, auctionBonusToBeGiven, "lend", whitelistingData.IsDutchActivated, false, pair.AssetIn, pair.AssetOut)
	if err != nil {
		return err
	}

	k.lend.UpdateBorrowStats(ctx, lendPair, borrow.IsStableBorrow, borrow.AmountOut.Amount, false)
	lendPos.AmountIn.Amount = lendPos.AmountIn.Amount.Sub(borrow.AmountIn.Amount)
	k.lend.UpdateLendStats(ctx, lendPos.AssetID, lendPos.PoolID, borrow.AmountIn.Amount, false)
	if !lendPos.AmountIn.Amount.GT(sdk.ZeroInt()) {
		// delete lend position
		k.lend.DeleteLendForAddressByAsset(ctx, lendPos.Owner, lendPos.ID)
		k.lend.DeleteIDFromAssetStatsMapping(ctx, lendPos.PoolID, lendPos.AssetID, borrow.LendingID, true)
		k.lend.DeleteLend(ctx, lendPos.ID)
	} else {
		k.lend.SetLend(ctx, lendPos)
	}

	return nil
}

func (k Keeper) MsgLiquidate(ctx sdk.Context, liquidator string, liqType, id uint64) error {
	if liqType == 0 {
		err := k.LiquidateIndividualVault(ctx, id, liquidator, true)
		if err != nil {
			return err
		}
	} else if liqType == 1 {
		err := k.LiquidateIndividualBorrow(ctx, id, liquidator, true)
		if err != nil {
			return err
		}
	}

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

func (k Keeper) WhitelistLiquidation(ctx sdk.Context, msg types.LiquidationWhiteListing) error {
	k.SetLiquidationWhiteListing(ctx, msg)
	return nil
}

func (k Keeper) LiquidateForSurplusAndDebt(ctx sdk.Context) error {
	auctionMapData, _ := k.collector.GetAllAuctionMappingForApp(ctx)
	for _, data := range auctionMapData {
		killSwitchParams, _ := k.esm.GetKillSwitchData(ctx, data.AppId)
		if !data.IsAuctionActive && !killSwitchParams.BreakerEnable {
			err := k.CheckStatsForSurplusAndDebt(ctx, data.AppId, data.AssetId)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func (k Keeper) CheckStatsForSurplusAndDebt(ctx sdk.Context, appID, assetID uint64) error {
	collector, found := k.collector.GetCollectorLookupTable(ctx, appID, assetID)
	if !found {
		return nil
	}
	auctionLookupTable, _ := k.collector.GetAuctionMappingForApp(ctx, appID, assetID)
	// coin denomination will be of 2 type: Auctioned Asset the asset which is being sold; i.e. Collateral Token // CMST
	// Asset required to bid on Collateral Asset; i.e. Debt Token // HARBOR

	netFeeCollectedData, found := k.collector.GetNetFeeCollectedData(ctx, appID, assetID)
	if !found {
		return nil
	}

	//surplus condition
	collateralAssetID := collector.CollectorAssetId //cmst
	debtAssetID := collector.SecondaryAssetId       //harbor

	// for debt Auction
	if netFeeCollectedData.NetFeesCollected.LTE(collector.DebtThreshold.Sub(collector.LotSize)) && auctionLookupTable.IsDebtAuction {
		// net = 200 debtThreshold = 500 , lotSize = 100
		collateralToken, debtToken := k.DebtTokenAmount(ctx, collateralAssetID, debtAssetID, collector.LotSize, collector.DebtLotSize)
		err := k.CreateLockedVault(ctx, 0, 0, "", collateralToken, debtToken, collateralToken, debtToken, sdk.ZeroDec(), appID, false, "", "", sdk.ZeroInt(), sdk.ZeroInt(), "debt", false, true, collateralAssetID, debtAssetID)
		if err != nil {
			return err
		}
		auctionLookupTable.IsAuctionActive = true
		err1 := k.collector.SetAuctionMappingForApp(ctx, auctionLookupTable)
		if err1 != nil {
			return err1
		}
	}

	// for surplus auction
	if netFeeCollectedData.NetFeesCollected.GTE(collector.SurplusThreshold.Add(collector.LotSize)) && auctionLookupTable.IsSurplusAuction {
		// net = 900 surplusThreshold = 500 , lotSize = 100
		amount := collector.LotSize
		collateralToken, debtToken := k.SurplusTokenAmount(ctx, collateralAssetID, debtAssetID, amount)

		// check to see if we have amount in collector
		_, err := k.collector.GetAmountFromCollector(ctx, appID, assetID, collateralToken.Amount)
		if err != nil {
			return err
		}
		err = k.CreateLockedVault(ctx, 0, 0, "", collateralToken, debtToken, collateralToken, debtToken, sdk.ZeroDec(), appID, false, "", "", sdk.ZeroInt(), sdk.ZeroInt(), "surplus", false, false, collateralAssetID, debtAssetID)
		if err != nil {
			return err
		}
		auctionLookupTable.IsAuctionActive = true
		err1 := k.collector.SetAuctionMappingForApp(ctx, auctionLookupTable)
		if err1 != nil {
			return err1
		}
	}

	return nil
}

//	Surplus  ``` | Debt`
//
// --Collateral 		cmst		harbor
// debt				harbor		cmst
func (k Keeper) DebtTokenAmount(ctx sdk.Context, DebtAssetID, CollateralAssetId uint64, lotSize, debtLotSize sdk.Int) (collateralToken, debtToken sdk.Coin) {
	collateralAsset, found1 := k.asset.GetAsset(ctx, CollateralAssetId)
	debtAsset, found2 := k.asset.GetAsset(ctx, DebtAssetID)
	if !found1 || !found2 {
		return sdk.Coin{}, sdk.Coin{}
	}
	return sdk.NewCoin(collateralAsset.Denom, debtLotSize), sdk.NewCoin(debtAsset.Denom, lotSize)
}

func (k Keeper) SurplusTokenAmount(ctx sdk.Context, CollateralAssetId, DebtAssetID uint64, lotSize sdk.Int) (collateralToken, debtToken sdk.Coin) {
	collateralAsset, found1 := k.asset.GetAsset(ctx, CollateralAssetId)
	debtAsset, found2 := k.asset.GetAsset(ctx, DebtAssetID)
	if !found1 || !found2 {
		return sdk.Coin{}, sdk.Coin{}
	}
	return sdk.NewCoin(collateralAsset.Denom, lotSize), sdk.NewCoin(debtAsset.Denom, sdk.ZeroInt())
}

func (k Keeper) MsgAppReserveFundsFn(ctx sdk.Context, from string, appId, assetId uint64, tokenQuantity sdk.Coin) error {
	asset, found := k.asset.GetAsset(ctx, assetId)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}

	if asset.Denom != tokenQuantity.Denom {
		return assettypes.ErrorInvalidDenom
	}

	_, found = k.asset.GetApp(ctx, appId)
	if !found {
		return assettypes.ErrorUnknownAppType
	}

	appReserveFunds, found := k.GetAppReserveFunds(ctx, appId, assetId)
	if !found {
		appReserveFunds = types.AppReserveFunds{
			AppId:         appId,
			AssetId:       assetId,
			TokenQuantity: tokenQuantity,
		}
	} else {
		appReserveFunds.TokenQuantity.Amount = appReserveFunds.TokenQuantity.Amount.Add(tokenQuantity.Amount)
	}

	addr, err := sdk.AccAddressFromBech32(from)
	if err != nil {
		return err
	}

	err = k.bank.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, sdk.NewCoins(tokenQuantity))
	if err != nil {
		return err
	}
	k.SetAppReserveFunds(ctx, appReserveFunds)

	// trigger AppReserveFundsTxData

	assetTxData := types.AssetTxData{
		AssetId:       assetId,
		TxType:        "credit",
		TokenQuantity: tokenQuantity,
	}

	appReserveFundsTxData, found := k.GetAppReserveFundsTxData(ctx, appId)
	if !found {
		appReserveFundsTxData.AppId = appId
	}
	appReserveFundsTxData.AssetTxData = append(appReserveFundsTxData.AssetTxData, assetTxData)
	k.SetAppReserveFundsTxData(ctx, appReserveFundsTxData)
	return nil
}

func (k Keeper) WithdrawAppReserveFundsFn(ctx sdk.Context, appId, assetId uint64, tokenQuantity sdk.Coin) error {
	appReserveFunds, found := k.GetAppReserveFunds(ctx, appId, assetId)
	if !found {
		return types.ErrorInvalidAppOrAssetData
	}

	if appReserveFunds.TokenQuantity.Amount.Sub(tokenQuantity.Amount).GTE(sdk.ZeroInt()) {
		if tokenQuantity.Amount.GT(sdk.ZeroInt()) {
			err := k.bank.SendCoinsFromModuleToModule(ctx, types.ModuleName, auctionsV2types.ModuleName, sdk.NewCoins(tokenQuantity))
			if err != nil {
				return err
			}
		}
	}
	appReserveFunds.TokenQuantity.Amount = appReserveFunds.TokenQuantity.Amount.Sub(tokenQuantity.Amount)
	k.SetAppReserveFunds(ctx, appReserveFunds)
	// trigger AppReserveFundsTxData
	assetTxData := types.AssetTxData{
		AssetId:       assetId,
		TxType:        "debt",
		TokenQuantity: tokenQuantity,
	}
	appReserveFundsTxData, _ := k.GetAppReserveFundsTxData(ctx, appId)
	appReserveFundsTxData.AssetTxData = append(appReserveFundsTxData.AssetTxData, assetTxData)
	k.SetAppReserveFundsTxData(ctx, appReserveFundsTxData)
	return nil
}

func (k Keeper) SetAppReserveFunds(ctx sdk.Context, appReserveFunds types.AppReserveFunds) {
	var (
		store = k.Store(ctx)
		key   = types.AppReserveFundsKey(appReserveFunds.AppId, appReserveFunds.AssetId)
		value = k.cdc.MustMarshal(&appReserveFunds)
	)

	store.Set(key, value)
}

func (k Keeper) GetAppReserveFunds(ctx sdk.Context, appId, assetId uint64) (appReserveFunds types.AppReserveFunds, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.AppReserveFundsKey(appId, assetId)
		value = store.Get(key)
	)

	if value == nil {
		return appReserveFunds, false
	}

	k.cdc.MustUnmarshal(value, &appReserveFunds)
	return appReserveFunds, true
}

func (k Keeper) SetAppReserveFundsTxData(ctx sdk.Context, appReserveFundsTxData types.AppReserveFundsTxData) {
	var (
		store = k.Store(ctx)
		key   = types.AppReserveFundsTxDataKey(appReserveFundsTxData.AppId)
		value = k.cdc.MustMarshal(&appReserveFundsTxData)
	)

	store.Set(key, value)
}

func (k Keeper) GetAppReserveFundsTxData(ctx sdk.Context, appId uint64) (appReserveFundsTxData types.AppReserveFundsTxData, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.AppReserveFundsTxDataKey(appId)
		value = store.Get(key)
	)

	if value == nil {
		return appReserveFundsTxData, false
	}

	k.cdc.MustUnmarshal(value, &appReserveFundsTxData)
	return appReserveFundsTxData, true
}

func (k Keeper) MsgLiquidateExternal(ctx sdk.Context, from string, appID uint64, owner string, collateralToken, debtToken sdk.Coin, collateralAssetId, debtAssetId uint64, isDebtCmst bool) error {
	// check if the assets exists
	// check if reserve funds are added for the debt or not
	// send tokens from the liquidator's address to the auction module
	auctionParams, found := k.auctionsV2.GetAuctionParams(ctx)
	if !found {
		return auctionsV2types.ErrAuctionParamsNotFound
	}

	_, found = k.asset.GetAsset(ctx, collateralAssetId)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	_, found = k.asset.GetAsset(ctx, debtAssetId)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}

	appReserveFunds, found := k.GetAppReserveFunds(ctx, appID, debtAssetId)
	if !found || appReserveFunds.TokenQuantity.Amount.LTE(sdk.NewInt(0)) {
		return fmt.Errorf("reserve funds not added for debt asset id")
	}

	feeToBeCollected := sdk.NewDecFromInt(debtToken.Amount).Mul(auctionParams.LiquidationPenalty).TruncateInt()
	feetoken := sdk.NewCoin(debtToken.Denom, feeToBeCollected)
	//Calculating auction bonus to be given
	bonusToBeGiven := sdk.NewDecFromInt(debtToken.Amount).Mul(auctionParams.AuctionBonus).TruncateInt()

	addr, err := sdk.AccAddressFromBech32(from)
	err = k.bank.SendCoinsFromAccountToModule(ctx, addr, auctionsV2types.ModuleName, sdk.NewCoins(collateralToken))
	if err != nil {
		return err
	}
	err = k.CreateLockedVault(ctx, 0, 0, owner, collateralToken, debtToken, collateralToken, debtToken.Add(feetoken), sdk.ZeroDec(), appID, false, "", from, feeToBeCollected, bonusToBeGiven, "external", true, isDebtCmst, collateralAssetId, debtAssetId)
	if err != nil {
		return fmt.Errorf("error Creating Locked Vaults in Liquidation, liquidate_vaults.go for External liquidation ")
	}

	return nil
}

func (k Keeper) MsgCloseDutchAuctionForBorrow(ctx sdk.Context, liquidationData types.LockedVault, auctionData auctionsV2types.Auction) error {
	// send money back to the debt pool (assetOut pool)
	// liquidation penalty to the reserve and interest to the pool
	// send token to the bidder
	// if cross pool borrow, settle the transit asset to it's native pool
	// close the borrow and update the stats

	borrowPos, _ := k.lend.GetBorrow(ctx, liquidationData.OriginalVaultId)
	pair, _ := k.lend.GetLendPair(ctx, borrowPos.PairID)
	pool, _ := k.lend.GetPool(ctx, pair.AssetOutPoolID)
	poolAssetLBMappingData, _ := k.lend.GetAssetStatsByPoolIDAndAssetID(ctx, pair.AssetOutPoolID, pair.AssetOut)
	amountToPool := liquidationData.TargetDebt
	assetOutStats, _ := k.lend.GetAssetRatesParams(ctx, pair.AssetOut)
	assetInStats, _ := k.lend.GetAssetRatesParams(ctx, pair.AssetIn)
	cAsset, _ := k.asset.GetAsset(ctx, assetOutStats.CAssetID)
	lend, _ := k.lend.GetLend(ctx, borrowPos.LendingID)

	// sending tokens debt tokens to the pool
	err := k.bank.SendCoinsFromModuleToModule(ctx, auctionsV2types.ModuleName, pool.ModuleName, sdk.NewCoins(amountToPool))
	if err != nil {
		return err
	}

	reservePoolRecords, found := k.lend.GetBorrowInterestTracker(ctx, liquidationData.OriginalVaultId)
	if !found {
		reservePoolRecords = lendtypes.BorrowInterestTracker{
			BorrowingId:         liquidationData.OriginalVaultId,
			ReservePoolInterest: sdk.ZeroDec(),
		}
	}

	// sending liquidation penalty
	liquidationPenalty := assetInStats.LiquidationPenalty
	if pair.IsEModeEnabled {
		liquidationPenalty = assetInStats.ELiquidationPenalty
	}
	liqPenaltyAmount := sdk.NewDecFromInt(borrowPos.AmountOut.Amount).Mul(liquidationPenalty).TruncateInt()
	err = k.lend.UpdateReserveBalances(ctx, pair.AssetOut, pool.ModuleName, sdk.NewCoin(borrowPos.AmountOut.Denom, liqPenaltyAmount), true)
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
	allReserveStats.AmountInFromLiqPenalty = allReserveStats.AmountInFromLiqPenalty.Add(liqPenaltyAmount)

	// calculating amount sent to be reserve pool and the debt pool
	// after recovering interest some part of the interest goes into the reserve pool
	// and for the remaining quantity equivalent number of cToken is minted to be given to the lenders upon calculate interest and rewards
	amtToReservePool := reservePoolRecords.ReservePoolInterest
	if amtToReservePool.TruncateInt().GT(sdk.ZeroInt()) {
		amount := sdk.NewCoin(auctionData.DebtToken.Denom, amtToReservePool.TruncateInt())
		err = k.lend.UpdateReserveBalances(ctx, pair.AssetOut, pool.ModuleName, amount, true)
		if err != nil {
			return err
		}
		allReserveStats.AmountInFromRepayments = allReserveStats.AmountInFromRepayments.Add(amount.Amount)
	}
	k.lend.SetAllReserveStatsByAssetID(ctx, allReserveStats)
	// amount minted in the debt pool
	amtToMint := (borrowPos.InterestAccumulated.Sub(amtToReservePool)).TruncateInt()
	if amtToMint.GT(sdk.ZeroInt()) {
		err = k.bank.MintCoins(ctx, pool.ModuleName, sdk.NewCoins(sdk.NewCoin(cAsset.Denom, amtToMint)))
		if err != nil {
			return err
		}
		poolAssetLBMappingData.TotalInterestAccumulated = poolAssetLBMappingData.TotalInterestAccumulated.Add(amtToMint)
		k.lend.SetAssetStatsByPoolIDAndAssetID(ctx, poolAssetLBMappingData)
	}
	// if borrow position is having bridged asset then return to the initial pool
	if borrowPos.BridgedAssetAmount.Amount.GT(sdk.NewInt(0)) {
		assetInPool, _ := k.lend.GetPool(ctx, lend.PoolID)
		err = k.bank.SendCoinsFromModuleToModule(ctx, pool.ModuleName, assetInPool.ModuleName, sdk.NewCoins(borrowPos.BridgedAssetAmount))
		if err != nil {
			return err
		}
	}

	k.lend.DeleteIDFromAssetStatsMapping(ctx, pair.AssetOutPoolID, pair.AssetOut, liquidationData.OriginalVaultId, false)
	k.lend.DeleteBorrowIDFromUserMapping(ctx, liquidationData.Owner, borrowPos.LendingID, liquidationData.OriginalVaultId)
	k.lend.DeleteBorrow(ctx, liquidationData.OriginalVaultId)
	k.lend.DeleteBorrowInterestTracker(ctx, liquidationData.OriginalVaultId)
	return nil
}
