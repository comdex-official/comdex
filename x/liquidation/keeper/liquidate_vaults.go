package keeper

import (
	"strconv"
	"time"

	"github.com/comdex-official/comdex/x/liquidation/types"
	vaulttypes "github.com/comdex-official/comdex/x/vault/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (k Keeper) LiquidateVaults(ctx sdk.Context) error {
	appIds := k.GetAppIds(ctx).WhitelistedAppMappingIds

	for i := range appIds {

		vaultsMap, _ := k.GetAppExtendedPairVaultMapping(ctx, appIds[i])

		vaults := vaultsMap.ExtendedPairVaults
		for j := range vaults {

			vaultIds := vaults[j].VaultIds
			for l := range vaultIds {
				vault, found := k.GetVault(ctx, vaultIds[l])
				if !found {
					continue
				}

				extPair, _ := k.GetPairsVault(ctx, vault.ExtendedPairVaultID)

				liqRatio := extPair.LiquidationRatio
				totalOut := vault.AmountOut.Add(vault.InterestAccumulated).Add(vault.ClosingFeeAccumulated)
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

	lockedVaultId := k.GetLockedVaultID(ctx)

	var value = types.LockedVault{
		LockedVaultId:                lockedVaultId + 1,
		AppMappingId:                 appId,
		AppVaultTypeId:               strconv.FormatUint(appId, 10),
		OriginalVaultId:              vault.Id,
		ExtendedPairId:               vault.ExtendedPairVaultID,
		Owner:                        vault.Owner,
		AmountIn:                     vault.AmountIn,
		AmountOut:                    vault.AmountOut,
		UpdatedAmountOut:             vault.AmountOut.Add(vault.InterestAccumulated).Add(vault.ClosingFeeAccumulated),
		Initiator:                    types.ModuleName,
		IsAuctionComplete:            false,
		IsAuctionInProgress:          false,
		CrAtLiquidation:              collateralizationRatio,
		CurrentCollaterlisationRatio: collateralizationRatio,
		CollateralToBeAuctioned:      sdk.ZeroDec(),
		LiquidationTimestamp:         time.Now(),
		SellOffHistory:               nil,
	}

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
	LockedVaultToApp, _ := k.GetLockedVaultByAppId(ctx, lockedVault.AppMappingId)
	LockedVaultToApp.LockedVault = append(LockedVaultToApp.LockedVault, &lockedVault)

	newLockedVaultToApp := types.LockedVaultToAppMapping{
		AppMappingId: lockedVault.AppMappingId,
		LockedVault:  LockedVaultToApp.LockedVault,
	}
	k.SetLockedVaultByAppId(ctx, newLockedVaultToApp)
}

func (k Keeper) SetLockedVaultByAppId(ctx sdk.Context, msg types.LockedVaultToAppMapping) {
	var (
		store = k.Store(ctx)
		key   = types.AppIDLockedVaultMappingKey(msg.AppMappingId)
		value = k.cdc.MustMarshal(&msg)
	)

	store.Set(key, value)
}

func (k *Keeper) GetLockedVaultByAppId(ctx sdk.Context, appMappingId uint64) (msg types.LockedVaultToAppMapping, found bool) {
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

func (k Keeper) CreateLockedVaultHistory(ctx sdk.Context, lockedVault types.LockedVault) error {

	lockedVaultId := k.GetLockedVaultIDHistory(ctx)
	k.SetLockedVaultHistory(ctx, lockedVault, lockedVaultId)
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
				if (!lockedVault.IsAuctionInProgress && !lockedVault.IsAuctionComplete) || (lockedVault.IsAuctionComplete && lockedVault.CurrentCollaterlisationRatio.LTE(ExtPair.MinCr)) {
					pair, _ := k.GetPair(ctx, ExtPair.PairId)
					assetIn, found := k.GetAsset(ctx, pair.AssetIn)
					if !found {
						continue
					}

					assetOut, found := k.GetAsset(ctx, pair.AssetOut)
					if !found {
						continue
					}

					collateralizationRatio, err := k.CalculateCollaterlizationRatio(ctx, ExtPair.PairId, lockedVault.AmountIn, lockedVault.UpdatedAmountOut)
					if err != nil {
						continue
					}

					//Asset Price in Dollar Terms to find how how much is to be auctioned

					extendedPairVault, found := k.GetPairsVault(ctx, lockedVault.ExtendedPairId)
					if !found {
						return types.ErrorExtendedPairVaultDoesNotExists
					}
					var assetOutPrice uint64

					if extendedPairVault.AssetOutOraclePrice {

						//If oracle Price required for the assetOut
						assetOutPrice, found = k.GetPriceForAsset(ctx, assetOut.Id)

						if !found {
							return types.ErrorPriceDoesNotExist
						}
					} else {
						//If oracle Price is not required for the assetOut
						assetOutPrice = extendedPairVault.AssetOutPrice

					}

					assetInPrice, _ := k.GetPriceForAsset(ctx, assetIn.Id)
					// assetOutPrice, _ := k.GetPriceForAsset(ctx, assetOut.Id)

					totalIn := lockedVault.AmountIn.Mul(sdk.NewIntFromUint64(assetInPrice)).ToDec()
					totalOut := lockedVault.AmountOut.Mul(sdk.NewIntFromUint64(assetOutPrice)).ToDec()
					//Selloff Collateral Calculation
					//Assuming that the collateral to be sold is 1 unit, so finding out how much is going to be deducted from the
					//collateral which will account as repaying the user's debt

					numpart1 := ExtPair.MinCr.Mul(totalOut)
					numpart2 := ExtPair.LiquidationPenalty.Add(sdk.MustNewDecFromStr("1"))
					numerator1 := numpart1.Sub(totalIn)
					numerator := numerator1.Mul(numpart2)
					denominator1 := ExtPair.MinCr.Sub(numpart2)
					denominator := denominator1.Mul((sdk.NewIntFromUint64(assetInPrice).ToDec()))
					collateralAuctioned := numerator.Quo(denominator)
					collateralToBeAuctioned := collateralAuctioned.Mul(sdk.NewDec(1000000))

					// safeLiquidationFactor, _ := sdk.NewDecFromStr(types.SafeLiquidationFactor)
					// deductionPercentage, _ := sdk.NewDecFromStr("1.0")
					// auctionDeduction := (deductionPercentage).Sub(ExtPair.LiquidationPenalty)
					// multiplicationFactor := auctionDeduction.Mul(ExtPair.MinCr.Add(safeLiquidationFactor))
					// asssetOutMultiplicationFactor := totalOut.Mul(ExtPair.MinCr.Add(safeLiquidationFactor))
					// assetsDifference := totalIn.Sub(asssetOutMultiplicationFactor)
					// //Substracting again from 1 unit to find the selloff multiplication factor
					// selloffMultiplicationFactor := deductionPercentage.Sub(multiplicationFactor)
					// selloffAmount := assetsDifference.Quo(selloffMultiplicationFactor)
					// var collateralToBeAuctioned sdk.Dec

					if collateralToBeAuctioned.GTE(totalIn) || collateralToBeAuctioned.IsNegative(){
						collateralToBeAuctioned = totalIn
					}
					updatedLockedVault := lockedVault
					updatedLockedVault.CurrentCollaterlisationRatio = collateralizationRatio
					updatedLockedVault.CollateralToBeAuctioned = collateralToBeAuctioned
					//updatedLockedVault.IsAuctionComplete = false
					k.SetLockedVault(ctx, updatedLockedVault)
					k.UpdateLockedVaultsAppMapping(ctx, updatedLockedVault)
				}
			}
		}
	}
	return nil
}

func (k Keeper) UnliquidateLockedVaults(ctx sdk.Context) error {
	//not unliquidating
	//appIds := k.GetAppIds(ctx).WhitelistedAppMappingIds

	lockedVaults := k.GetLockedVaults(ctx)

	for _, lockedVault := range lockedVaults {

		if lockedVault.IsAuctionComplete {
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

			if lockedVault.AmountIn.IsZero() && lockedVault.AmountOut.IsZero() {
				err := k.CreateLockedVaultHistory(ctx, lockedVault)
				if err != nil {
					return err
				}
				k.UpdateUserVaultExtendedPairMapping(ctx, lockedVault.ExtendedPairId, lockedVault.Owner, lockedVault.AppMappingId)
				k.DeleteLockedVault(ctx, lockedVault.LockedVaultId)
				continue
			}

			if lockedVault.AmountOut.IsZero() {

				err := k.CreateLockedVaultHistory(ctx, lockedVault)
				if err != nil {
					return err
				}

				k.UpdateUserVaultExtendedPairMapping(ctx, lockedVault.ExtendedPairId, lockedVault.Owner, lockedVault.AppMappingId)

				k.DeleteLockedVault(ctx, lockedVault.LockedVaultId)
				if err := k.SendCoinFromModuleToAccount(ctx, vaulttypes.ModuleName, userAddress, sdk.NewCoin(assetIn.Denom, lockedVault.AmountIn)); err != nil {
					continue
				}
				continue
			}
			newCalculatedCollateralizationRatio, err := k.CalculateCollaterlizationRatio(ctx, extPair.PairId, lockedVault.AmountIn, lockedVault.UpdatedAmountOut)
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
				k.UpdateUserVaultExtendedPairMapping(ctx, lockedVault.ExtendedPairId, lockedVault.Owner, lockedVault.AppMappingId)

				err = k.CreteNewVault(ctx, lockedVault.Owner, lockedVault.AppMappingId, lockedVault.ExtendedPairId, lockedVault.AmountIn, lockedVault.AmountOut)
				if err != nil {
					return err
				}

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

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k *Keeper) SetLockedVaultID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultIdKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)
	store.Set(key, value)
}

func (k *Keeper) GetLockedVaultID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultIdKey
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
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

func (k *Keeper) SetLockedVault(ctx sdk.Context, lockedVault types.LockedVault) {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultKey(lockedVault.LockedVaultId)
		value = k.cdc.MustMarshal(&lockedVault)
	)
	store.Set(key, value)
}

func (k *Keeper) SetLockedVaultHistory(ctx sdk.Context, lockedVault types.LockedVault, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultHistoryKey(id)
		value = k.cdc.MustMarshal(&lockedVault)
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

func (k *Keeper) GetLockedVault(ctx sdk.Context, id uint64) (lockedVault types.LockedVault, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LockedVaultKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return lockedVault, false
	}

	k.cdc.MustUnmarshal(value, &lockedVault)
	return lockedVault, true
}

func (k *Keeper) GetLockedVaults(ctx sdk.Context) (lockedVaults []types.LockedVault) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.LockedVaultKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {

		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var lockedVault types.LockedVault
		k.cdc.MustUnmarshal(iter.Value(), &lockedVault)
		lockedVaults = append(lockedVaults, lockedVault)
	}

	return lockedVaults
}

func (k *Keeper) SetFlagIsAuctionInProgress(ctx sdk.Context, id uint64, flag bool) error {

	lockedVault, found := k.GetLockedVault(ctx, id)
	if !found {
		return types.LockedVaultDoesNotExist
	}
	lockedVault.IsAuctionInProgress = flag
	k.SetLockedVault(ctx, lockedVault)
	return nil
}

func (k *Keeper) SetFlagIsAuctionComplete(ctx sdk.Context, id uint64, flag bool) error {

	lockedVault, found := k.GetLockedVault(ctx, id)
	if !found {
		return types.LockedVaultDoesNotExist
	}
	lockedVault.IsAuctionComplete = flag
	k.SetLockedVault(ctx, lockedVault)
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
