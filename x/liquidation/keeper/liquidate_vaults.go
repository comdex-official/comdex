package keeper

import (
	"fmt"
	"github.com/comdex-official/comdex/x/liquidation/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
	"strconv"
	"time"
)

func (k Keeper) LiquidateVaults(ctx sdk.Context) error {
	fmt.Println("Inside Liquidate vaults _____________")
	appIds := k.GetAppIds(ctx).WhitelistedAppMappingIds
	fmt.Println(appIds)
	for i := range appIds {
		fmt.Println("Inside Liquidate vaults 1_____________")
		vaultsMap, _ := k.GetAppExtendedPairVaultMapping(ctx, appIds[i])

		vaults := vaultsMap.ExtendedPairVaults
		for j := range vaults {
			fmt.Println("Inside Liquidate vaults 2______________")
			vaultIds := vaults[j].VaultIds
			for l := range vaultIds {
				vault, _ := k.GetVault(ctx, vaultIds[l])

				extPair, _ := k.GetPairsVault(ctx, vault.ExtendedPairVaultID)

				liqRatio := extPair.LiquidationRatio
				totalOut := vault.AmountOut.Add(*vault.InterestAccumulated).Add(*vault.ClosingFeeAccumulated)
				collateralitzationRatio, err := k.CalculateCollaterlizationRatio(ctx, vault.ExtendedPairVaultID, vault.AmountIn, totalOut)
				if err != nil {
					continue
				}

				if sdk.Dec.LT(collateralitzationRatio, liqRatio) {
					err := k.CreateLockedVault(ctx, vault, collateralitzationRatio, appIds[i])
					if err != nil {
						continue
					}
					k.DeleteVault(ctx, vault.Id)
					k.DeleteAddressFromAppExtendedPairVaultMapping(ctx, vault.ExtendedPairVaultID, vault.Id, appIds[i])
				}
			}
		}
	}
	return nil
}

func (k Keeper) CreateLockedVault(ctx sdk.Context, vault vaulttypes.Vault, collateralizationRatio sdk.Dec, appId uint64) error {

	lockedVaultId := k.GetLockedVaultIDbyApp(ctx, appId)

	var value = types.LockedVault{
		LockedVaultId:                lockedVaultId + 1,
		AppMappingId:                 appId,
		AppVaultTypeId:               strconv.FormatUint(appId, 10),
		OriginalVaultId:              vault.Id,
		ExtendedPairId:               vault.ExtendedPairVaultID,
		Owner:                        vault.Owner,
		AmountIn:                     vault.AmountIn,
		AmountOut:                    vault.AmountOut,
		UpdatedAmountOut:             vault.AmountOut.Add(*vault.InterestAccumulated).Add(*vault.ClosingFeeAccumulated),
		Initiator:                    types.ModuleName,
		IsAuctionComplete:            false,
		IsAuctionInProgress:          false,
		CrAtLiquidation:              collateralizationRatio,
		CurrentCollaterlisationRatio: collateralizationRatio,
		CollateralToBeAuctioned:      nil,
		LiquidationTimestamp:         time.Now(),
		SellOffHistory:               nil,
	}
	fmt.Println("CreateLockedVault")
	k.SetLockedVault(ctx, value)
	k.SetLockedVaultID(ctx, value.LockedVaultId)
	//k.UpdateLockedVaultsAppMapping(ctx, value)

	//Create a new Data Structure with the current Params
	//Set nil for all the values not available right now
	//New function will loop over locked vaults to set all values so that they can be auctioned, seperately
	//Auction will then use the selloff amount set by lockedvault function to update params .
	//Unliquidate will take place after all the events trigger.
	return nil

}

func (k Keeper) UpdateLockedVaultsAppMapping(ctx sdk.Context, lockedVault types.LockedVault) {
	LockedVaultToApp, _ := k.GetLockedVaultbyAppId(ctx, lockedVault.LockedVaultId)
	LockedVaultToApp.LockedVault = append(LockedVaultToApp.LockedVault, lockedVault)

	newLockedVaultToApp := types.LockedVaultToAppMapping{
		AppMappingId: lockedVault.AppMappingId,
		LockedVault:  LockedVaultToApp.LockedVault,
	}
	k.SetLockedVaultbyAppId(ctx, newLockedVaultToApp)
}

func (k Keeper) SetLockedVaultbyAppId(ctx sdk.Context, msg types.LockedVaultToAppMapping) {
	var (
		store = k.Store(ctx)
		key   = types.AppIDLockedVaultMappingKey(msg.AppMappingId)
		value = k.cdc.MustMarshal(&msg)
	)

	store.Set(key, value)
}

func (k *Keeper) GetLockedVaultbyAppId(ctx sdk.Context, appMappingId uint64) (msg types.LockedVaultToAppMapping, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.AppIDLockedVaultMappingKey(appMappingId)
		value = store.Get(key)
	)

	if value == nil {
		return msg, false
	}

	k.cdc.MustUnmarshal(value, &msg)
	return msg, true
}

func (k Keeper) CreateLockedVaultHistoy(ctx sdk.Context, lockedvault types.LockedVault) error {

	lockedVaultId := k.GetLockedVaultIDHistory(ctx)
	k.SetLockedVaultHistory(ctx, lockedvault, lockedVaultId)
	k.SetLockedVaultIDHistory(ctx, lockedVaultId+1)

	return nil

}

//for first time to update the collateralization value & sell off amount
//and if auction is complete and cr is less than 1.6
func (k Keeper) UpdateLockedVaults(ctx sdk.Context) error {
	appIds := k.GetAppIds(ctx).WhitelistedAppMappingIds
	for _, v := range appIds {
		newUpdatedVaults := k.GetLockedVaults(ctx)
		for _, lockedVault := range newUpdatedVaults {
			if lockedVault.AppMappingId == v {
				ExtPair, _ := k.GetPairsVault(ctx, lockedVault.ExtendedPairId)
				fmt.Println("locked vault", lockedVault)
				fmt.Println("!lockedVault.IsAuctionInProgress && !lockedVault.IsAuctionComplete", !lockedVault.IsAuctionInProgress && !lockedVault.IsAuctionComplete)
				fmt.Println("lockedVault.IsAuctionComplete && lockedVault.CurrentCollaterlisationRatio.LTE(ExtPair.MinCr)", lockedVault.IsAuctionComplete && lockedVault.CurrentCollaterlisationRatio.LTE(ExtPair.MinCr))
				if (!lockedVault.IsAuctionInProgress && !lockedVault.IsAuctionComplete) || (lockedVault.IsAuctionComplete && lockedVault.CurrentCollaterlisationRatio.LTE(ExtPair.MinCr)) {
					fmt.Println("(!lockedVault.IsAuctionInProgress && !lockedVault.IsAuctionComplete) || (lockedVault.IsAuctionComplete && lockedVault.CurrentCollaterlisationRatio.LTE(ExtPair.MinCr))", (!lockedVault.IsAuctionInProgress && !lockedVault.IsAuctionComplete) || (lockedVault.IsAuctionComplete && lockedVault.CurrentCollaterlisationRatio.LTE(ExtPair.MinCr)))
					pair, _ := k.GetPair(ctx, ExtPair.PairId)
					assetIn, found := k.GetAsset(ctx, pair.AssetIn)
					if !found {
						continue
					}
					fmt.Println("red1______________")
					assetOut, found := k.GetAsset(ctx, pair.AssetOut)
					if !found {
						continue
					}
					fmt.Println("red2______________")
					collateralizationRatio, err := k.CalculateCollaterlizationRatio(ctx, ExtPair.PairId, lockedVault.AmountIn, lockedVault.UpdatedAmountOut)
					if err != nil {
						continue
					}
					fmt.Println("red3______________")
					//Asset Price in Dollar Terms to find how how much is to be auctioned
					assetInPrice, _ := k.GetPriceForAsset(ctx, assetIn.Id)
					assetOutPrice, _ := k.GetPriceForAsset(ctx, assetOut.Id)

					totalIn := lockedVault.AmountIn.Mul(sdk.NewIntFromUint64(assetInPrice)).ToDec()
					totalOut := lockedVault.AmountOut.Mul(sdk.NewIntFromUint64(assetOutPrice)).ToDec()
					fmt.Println("red4______________")
					//Selloff Collateral Calculation
					//Assuming that the collateral to be sold is 1 unit, so finding out how much is going to be deducted from the
					//collateral which will account as repaying the user's debt
					perc, _ := sdk.NewDecFromStr("0.1")
					deductionPercentage, _ := sdk.NewDecFromStr("1.0")
					auctionDeduction := (deductionPercentage).Sub(ExtPair.LiquidationPenalty)
					multiplicationFactor := auctionDeduction.Mul(ExtPair.MinCr.Add(perc))
					asssetOutMultiplicationFactor := totalOut.Mul(ExtPair.MinCr.Add(perc))
					assetsDifference := totalIn.Sub(asssetOutMultiplicationFactor)
					//Substracting again from 1 unit to find the selloff multiplication factor
					selloffMultiplicationFactor := deductionPercentage.Sub(multiplicationFactor)
					selloffAmount := assetsDifference.Quo(selloffMultiplicationFactor)
					fmt.Println("red5______________")
					var collateralToBeAuctioned sdk.Dec

					if selloffAmount.GTE(totalIn) || selloffAmount.IsNegative() {
						collateralToBeAuctioned = totalIn
					} else {

						collateralToBeAuctioned = selloffAmount
					}
					fmt.Println("red6______________")
					updatedLockedVault := lockedVault
					updatedLockedVault.CurrentCollaterlisationRatio = collateralizationRatio
					updatedLockedVault.CollateralToBeAuctioned = &collateralToBeAuctioned
					//updatedLockedVault.IsAuctionComplete = false
					k.SetLockedVault(ctx, updatedLockedVault)
					fmt.Println("red7______________")
					k.UpdateLockedVaultsAppMapping(ctx, updatedLockedVault)
					fmt.Println("red8______________")
				}
			}
		}
	}
	return nil
}

func (k Keeper) UnliquidateLockedVaults(ctx sdk.Context) error {
	//not unliquidating
	fmt.Println("UnliquidateLockedVaults______")
	//appIds := k.GetAppIds(ctx).WhitelistedAppMappingIds

	fmt.Println("UnliquidateLockedVaults______1")
	lockedVaults := k.GetLockedVaults(ctx)
	fmt.Println("UnliquidateLockedVaults______1111111333", lockedVaults)
	for _, lockedVault := range lockedVaults {
		fmt.Println("UnliquidateLockedVaults______8429")
		ExtPair, _ := k.GetPairsVault(ctx, lockedVault.ExtendedPairId)
		if lockedVault.IsAuctionComplete && lockedVault.CurrentCollaterlisationRatio.GTE(ExtPair.MinCr) {
			//also calculate the current collaterlization ration to ensure there is no sudden changes
			userAddress, err := sdk.AccAddressFromBech32(lockedVault.Owner)
			if err != nil {
				continue
			}

			extPair, _ := k.GetPairsVault(ctx, lockedVault.ExtendedPairId)

			pair, found := k.GetPair(ctx, extPair.PairId)
			if !found {
				continue
			}

			unliquidatePointPercentage := extPair.MinCr

			assetIn, found := k.GetAsset(ctx, pair.AssetIn)
			if !found {
				continue
			}
			fmt.Println("UnliquidateLockedVaults______2")
			if lockedVault.AmountOut.IsZero() {
				fmt.Println("UnliquidateLockedVaults______6")
				err := k.CreateLockedVaultHistoy(ctx, lockedVault)
				if err != nil {
					return err
				}
				fmt.Println("UnliquidateLockedVaults______4")
				//k.DeleteAddressFromAppExtendedPairVaultMapping(ctx, lockedVault.ExtendedPairId, lockedVault.OriginalVaultId, lockedVault.AppMappingId)
				fmt.Println("deleting")
				fmt.Println("deleting")
				k.DeleteLockedVault(ctx, lockedVault.LockedVaultId)
				if err := k.SendCoinFromModuleToAccount(ctx, vaulttypes.ModuleName, userAddress, sdk.NewCoin(assetIn.Denom, lockedVault.AmountIn)); err != nil {
					continue
				}
				continue
			}
			fmt.Println("UnliquidateLockedVaults______8")
			newCalculatedCollateralizationRatio, err := k.CalculateCollaterlizationRatio(ctx, extPair.PairId, lockedVault.AmountIn, lockedVault.UpdatedAmountOut)
			fmt.Println("err", err)
			if err != nil {
				continue
			}
			fmt.Println("newCalculatedCollateralizationRatio", newCalculatedCollateralizationRatio)
			if newCalculatedCollateralizationRatio.LT(unliquidatePointPercentage) {
				updatedLockedVault := lockedVault
				updatedLockedVault.CurrentCollaterlisationRatio = newCalculatedCollateralizationRatio
				fmt.Println("UnliquidateLockedVaults")
				k.SetLockedVault(ctx, updatedLockedVault)
				continue
			}
			if newCalculatedCollateralizationRatio.GTE(unliquidatePointPercentage) {
				fmt.Println("wad________1")
				err := k.CreateLockedVaultHistoy(ctx, lockedVault)
				fmt.Println("wad________2", err)
				if err != nil {
					return err
				}
				//k.DeleteAddressFromAppExtendedPairVaultMapping(ctx, lockedVault.ExtendedPairId, lockedVault.OriginalVaultId, lockedVault.AppMappingId)
				fmt.Println("wad________3")
				err = k.CreteNewVault(ctx, lockedVault.Owner, lockedVault.AppMappingId, lockedVault.ExtendedPairId, lockedVault.AmountIn, lockedVault.AmountOut)
				if err != nil {
					return err
				}
				fmt.Println("deleting")
				fmt.Println("deleting")
				k.DeleteLockedVault(ctx, lockedVault.LockedVaultId)

				//======================================NOTE TO BE CHANGED================================================
				//One important thing that we missed is that we need to pop and append the current vault as per the user -> This has bee handled -Vishnu
				//IF all the borrowed amount is repayed , then we need to ensure the unliquidate vault is not called for that particular lockedvault- his vault is automatically closed.
			}
		}
	}

	return nil
}

func (k Keeper) GetModAccountBalances(ctx sdk.Context, accountName string, denom string) sdk.Int {
	macc := k.GetModuleAccount(ctx, accountName)
	return k.GetBalance(ctx, macc.GetAddress(), denom).Amount
}

func (k *Keeper) GetLockedVaultIDbyApp(ctx sdk.Context, appId uint64) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.AppLockedVaultMappingKey(appId)
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) GetLockedVaultIDHistory(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultKeyHistory
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}
	fmt.Println("value", value)

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetLockedVaultID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.AppLockedVaultMappingKey(id)
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)
	store.Set(key, value)
}

func (k *Keeper) SetLockedVaultIDHistory(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultKeyHistory
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)
	store.Set(key, value)
}

func (k *Keeper) SetLockedVault(ctx sdk.Context, locked_vault types.LockedVault) {
	fmt.Println("11111111111")
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultKey(locked_vault.LockedVaultId)
		value = k.cdc.MustMarshal(&locked_vault)
	)
	store.Set(key, value)
}

func (k *Keeper) SetLockedVaultHistory(ctx sdk.Context, locked_vault types.LockedVault, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultHistoryKey(id)
		value = k.cdc.MustMarshal(&locked_vault)
	)
	store.Set(key, value)
}

func (k *Keeper) DeleteLockedVault(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultKey(id)
	)
	store.Delete(key)
}

func (k *Keeper) GetLockedVault(ctx sdk.Context, id uint64) (locked_vault types.LockedVault, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return locked_vault, false
	}

	k.cdc.MustUnmarshal(value, &locked_vault)
	return locked_vault, true
}

func (k *Keeper) GetLockedVaults(ctx sdk.Context) (locked_vaults []types.LockedVault) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.LockedVaultKeyPrefix)
	)

	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var locked_vault types.LockedVault
		k.cdc.MustUnmarshal(iter.Value(), &locked_vault)
		locked_vaults = append(locked_vaults, locked_vault)
	}

	return locked_vaults
}

func (k *Keeper) SetFlagIsAuctionInProgress(ctx sdk.Context, id uint64, flag bool) error {

	locked_vault, found := k.GetLockedVault(ctx, id)
	if !found {
		return types.LockedVaultDoesNotExist
	}
	locked_vault.IsAuctionInProgress = flag
	fmt.Println("SetFlagIsAuctionInProgress")
	k.SetLockedVault(ctx, locked_vault)
	return nil
}

func (k *Keeper) SetFlagIsAuctionComplete(ctx sdk.Context, id uint64, flag bool) error {

	locked_vault, found := k.GetLockedVault(ctx, id)
	if !found {
		return types.LockedVaultDoesNotExist
	}
	locked_vault.IsAuctionComplete = flag
	fmt.Println("SetFlagIsAuctionComplete")
	k.SetLockedVault(ctx, locked_vault)
	return nil
}

/*func (k *Keeper) UpdateAssetQuantitiesInLockedVault(
	ctx sdk.Context,
	collateral_auction auctiontypes.CollateralAuction,
	amountIn sdk.Int,
	assetIn assettypes.Asset,
	amountOut sdk.Int,
	assetOut assettypes.Asset,
) error {

	locked_vault, found := k.GetLockedVault(ctx, collateral_auction.LockedVaultId)
	if !found {
		return types.LockedVaultDoesNotExist
	}
	updatedAmountIn := locked_vault.AmountIn.Sub(amountIn)
	updatedAmountOut := locked_vault.AmountOut.Sub(amountOut)
	updatedCollateralizationRatio, _ := k.CalculateCollaterlizationRatio(ctx, updatedAmountIn, assetIn, updatedAmountOut, assetOut)

	locked_vault.AmountIn = updatedAmountIn
	locked_vault.AmountOut = updatedAmountOut
	locked_vault.CurrentCollaterlisationRatio = updatedCollateralizationRatio
	locked_vault.SellOffHistory = append(locked_vault.SellOffHistory, collateral_auction.String())
	k.SetLockedVault(ctx, locked_vault)
	return nil
}*/

func (k *Keeper) SetAppId(ctx sdk.Context, AppIds types.WhitelistedAppIds) {
	var (
		store = k.Store(ctx)
		key   = types.AppIdsKeyPrefix
		value = k.cdc.MustMarshal(&AppIds)
	)

	store.Set(key, value)
}

func (k *Keeper) GetAppIds(ctx sdk.Context) (appIds types.WhitelistedAppIds) {
	var (
		store = k.Store(ctx)
		key   = types.AppIdsKeyPrefix
		value = store.Get(key)
	)

	if value == nil {
		return appIds
	}

	k.cdc.MustUnmarshal(value, &appIds)
	return appIds
}
