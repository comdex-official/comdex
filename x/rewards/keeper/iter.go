package keeper

import (
	"fmt"
	"math"
	"strconv"

	errorsmod "cosmossdk.io/errors"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	"github.com/comdex-official/comdex/x/rewards/types"
)

func (k Keeper) DistributeExtRewardLocker(ctx sdk.Context) error {
	// Give external rewards to locker owners for creating locker with specific assetID
	extRewards := k.GetExternalRewardsLockers(ctx)
	for _, v := range extRewards {
		klwsParams, _ := k.esm.GetKillSwitchData(ctx, v.AppMappingId)
		if klwsParams.BreakerEnable {
			return esmtypes.ErrCircuitBreakerEnabled
		}
		esmStatus, found := k.esm.GetESMStatus(ctx, v.AppMappingId)
		status := false
		if found {
			status = esmStatus.Status
		}
		if status {
			return esmtypes.ErrESMAlreadyExecuted
		}
		// checking if rewards are active

		if v.IsActive {
			epoch, _ := k.GetEpochTime(ctx, v.EpochId)
			et := epoch.StartingTime
			timeNow := ctx.BlockTime().Unix()

			// here the epoch starting time is set to the next day whenever any external locker reward is distributed
			// so when the epoch starting time is less than current time then the condition becomes true and flow passes through the function

			if et < timeNow {
				if epoch.Count < uint64(v.DurationDays) { // rewards will be given till the duration defined in the ext rewards
					// getting the total share of Deposited amount of Lockers of specific assetID and AppID
					lockerLookup, _ := k.locker.GetLockerLookupTable(ctx, v.AppMappingId, v.AssetId)
					totalShare := lockerLookup.DepositedAmount

					// initializing amountRewardedTracker to keep a track of daily rewards given to locker owners
					amountRewardedTracker := sdk.NewCoin(v.TotalRewards.Denom, sdkmath.ZeroInt())
					for _, lockerID := range lockerLookup.LockerIds {
						locker, found := k.locker.GetLocker(ctx, lockerID)
						if !found {
							continue
						}
						// checking if the locker was not created just to claim the external rewards, so we apply a basic check here.
						// last day don't check min lockup time, so we should have no remaining amount left
						if int64(epoch.Count) != v.DurationDays-1 {
							if timeNow-locker.CreatedAt.Unix() < v.MinLockupTimeSeconds {
								continue
							}
						}
						userShare := (sdkmath.LegacyNewDecFromInt(locker.NetBalance)).Quo(sdkmath.LegacyNewDecFromInt(totalShare)) // getting share percentage
						availableRewards := v.AvailableRewards                                                                     // Available Rewards
						Duration := v.DurationDays - int64(epoch.Count)                                                            // duration left (total duration - current count)

						epochRewards := sdkmath.LegacyNewDecFromInt(availableRewards.Amount).Quo(sdkmath.LegacyNewDec(Duration))
						dailyRewards := userShare.Mul(epochRewards)
						user, _ := sdk.AccAddressFromBech32(locker.Depositor)
						finalDailyRewards := dailyRewards.TruncateInt()
						// after calculating final daily rewards, the amount is sent to the user
						if finalDailyRewards.GT(sdkmath.ZeroInt()) {
							amountRewardedTracker.Amount = amountRewardedTracker.Amount.Add(finalDailyRewards)
							err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, user, sdk.NewCoins(sdk.NewCoin(availableRewards.Denom, finalDailyRewards)))
							if err != nil {
								continue
							}
						}
					}
					// after all the locker owners are rewarded
					// setting the starting time to next day
					epoch.Count = epoch.Count + types.UInt64One
					epoch.StartingTime = timeNow + types.SecondsPerDay
					k.SetEpochTime(ctx, epoch)

					// setting the available rewards by subtracting the amount sent per epoch for the ext rewards
					v.AvailableRewards.Amount = v.AvailableRewards.Amount.Sub(amountRewardedTracker.Amount)
					k.SetExternalRewardsLockers(ctx, v)
				} else {
					v.IsActive = false
					k.SetExternalRewardsLockers(ctx, v)
				}
			}
		}
	}
	return nil
}

func (k Keeper) DistributeExtRewardVault(ctx sdk.Context) error {
	// Give external rewards to vault owners for opening a vault with specific assetID
	extRewards := k.GetExternalRewardVaults(ctx)
	for _, v := range extRewards {
		klwsParams, _ := k.esm.GetKillSwitchData(ctx, v.AppMappingId)
		if klwsParams.BreakerEnable {
			return esmtypes.ErrCircuitBreakerEnabled
		}
		esmStatus, found := k.esm.GetESMStatus(ctx, v.AppMappingId)
		status := false
		if found {
			status = esmStatus.Status
		}
		if status {
			return esmtypes.ErrESMAlreadyExecuted
		}
		// checking if rewards are active
		if v.IsActive {
			epoch, _ := k.GetEpochTime(ctx, v.EpochId)
			et := epoch.StartingTime
			timeNow := ctx.BlockTime().Unix()

			// here the epoch starting time is set to the next day whenever any external vault reward is distributed
			// so when the epoch starting time is less than current time then the condition becomes true and flow passes through the function

			if et < timeNow {
				if epoch.Count < uint64(v.DurationDays) { // rewards will be given till the duration defined in the ext rewards
					appExtPairVaultData, _ := k.vault.GetAppExtendedPairVaultMappingData(ctx, v.AppMappingId, v.ExtendedPairId)

					// initializing amountRewardedTracker to keep a track of daily rewards given to locker owners
					amountRewardedTracker := sdk.NewCoin(v.TotalRewards.Denom, sdkmath.ZeroInt())

					for _, vaultID := range appExtPairVaultData.VaultIds {
						totalRewards := v.AvailableRewards
						userVault, found := k.vault.GetVault(ctx, vaultID)
						if !found {
							continue
						}
						// checking if the locker was not created just to claim the external rewards, so we apply a basic check here.
						// last day don't check min lockup time, so we should have no remaining amount left
						if int64(epoch.Count) != v.DurationDays-1 {
							if timeNow-userVault.CreatedAt.Unix() < v.MinLockupTimeSeconds {
								continue
							}
						}
						individualUserShare := sdkmath.LegacyNewDecFromInt(userVault.AmountOut).Quo(sdkmath.LegacyNewDecFromInt(appExtPairVaultData.TokenMintedAmount)) // getting share percentage
						Duration := v.DurationDays - int64(epoch.Count)                                                                                                 // duration left (total duration - current count)
						epochRewards := (sdkmath.LegacyNewDecFromInt(totalRewards.Amount)).Quo(sdkmath.LegacyNewDec(Duration))
						dailyRewards := individualUserShare.Mul(epochRewards)
						finalDailyRewards := dailyRewards.TruncateInt()

						user, _ := sdk.AccAddressFromBech32(userVault.Owner)
						if finalDailyRewards.GT(sdkmath.ZeroInt()) {
							amountRewardedTracker = amountRewardedTracker.Add(sdk.NewCoin(totalRewards.Denom, finalDailyRewards))
							err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, user, sdk.NewCoins(sdk.NewCoin(totalRewards.Denom, finalDailyRewards)))
							if err != nil {
								continue
							}
						}
					}
					// after all the vault owners are rewarded
					// setting the starting time to next day
					epoch.Count = epoch.Count + types.UInt64One
					epoch.StartingTime = timeNow + types.SecondsPerDay
					k.SetEpochTime(ctx, epoch)

					// setting the available rewards by subtracting the amount sent per epoch for the ext rewards
					v.AvailableRewards.Amount = v.AvailableRewards.Amount.Sub(amountRewardedTracker.Amount)

					k.SetExternalRewardVault(ctx, v)
				} else {
					v.IsActive = false
					k.SetExternalRewardVault(ctx, v)
				}
			}
		}
	}
	return nil
}

// calculate new locker rewards
func (k Keeper) CalculationOfRewards(
	ctx sdk.Context,
	amount sdkmath.Int, lsr sdkmath.LegacyDec, bTime int64,
) (sdkmath.LegacyDec, error) {
	currentTime := ctx.BlockTime().Unix()

	secondsElapsed := currentTime - bTime
	if secondsElapsed < types.Int64Zero {
		return sdkmath.LegacyZeroDec(), errorsmod.Wrap(types.ErrNegativeTimeElapsed, fmt.Sprintf("%d seconds", secondsElapsed))
	}
	//{(1+ Annual Interest Rate)^(No of seconds per block/No. of seconds in a year)}-1

	yearsElapsed := sdkmath.LegacyNewDec(secondsElapsed).QuoInt64(types.SecondsPerYear)
	perc := lsr.String()
	a, _ := sdkmath.LegacyNewDecFromStr("1")
	b, _ := sdkmath.LegacyNewDecFromStr(perc)
	factor1 := a.Add(b)
	intPerBlockFactor := math.Pow(factor1.MustFloat64(), yearsElapsed.MustFloat64())
	intAccPerBlock := intPerBlockFactor - types.Float64One
	amtFloat := sdkmath.LegacyNewDecFromInt(amount).MustFloat64()
	newAmount := intAccPerBlock * amtFloat

	// s := fmt.Sprint(newAmount)
	s := strconv.FormatFloat(newAmount, 'f', 18, 64)
	newAm, err := sdkmath.LegacyNewDecFromStr(s)
	if err != nil {
		return sdkmath.LegacyZeroDec(), err
	}
	return newAm, nil
}

func (k Keeper) OraclePriceForRewards(ctx sdk.Context, id uint64, amt sdkmath.Int) (sdkmath.LegacyDec, bool) {
	asset, found := k.asset.GetAsset(ctx, id)
	if !found {
		return sdkmath.LegacyZeroDec(), false
	}

	price, found := k.marketKeeper.GetTwa(ctx, asset.Id)
	if !found {
		return sdkmath.LegacyZeroDec(), false
	}

	// if price is not active and twa is 0 return false
	if !price.IsPriceActive && price.Twa <= 0 {
		return sdkmath.LegacyZeroDec(), false
	}

	numerator := sdkmath.LegacyNewDecFromInt(amt).Mul(sdkmath.LegacyNewDecFromInt(sdkmath.NewIntFromUint64(price.Twa)))
	denominator := sdkmath.LegacyNewDecFromInt(asset.Decimals)
	return numerator.Quo(denominator), true
}

func (k Keeper) DistributeExtRewardLend(ctx sdk.Context) error {
	// Give external rewards to borrowers for opening a vault with specific assetID
	var addrArr []string
	var amountArr []sdkmath.LegacyDec
	totalAmount := sdkmath.NewInt(0)
	extRewards := k.GetExternalRewardLends(ctx)
	for _, v := range extRewards {
		klwsParams, _ := k.esm.GetKillSwitchData(ctx, v.AppMappingId)
		if klwsParams.BreakerEnable {
			return esmtypes.ErrCircuitBreakerEnabled
		}
		if v.IsActive {
			epoch, _ := k.GetEpochTime(ctx, v.EpochId)
			et := epoch.StartingTime
			timeNow := ctx.BlockTime().Unix()
			// checking if rewards are active
			if et < timeNow {
				if epoch.Count < uint64(v.DurationDays) {
					rewardsAssetPoolData := v.RewardsAssetPoolData
					assetID := rewardsAssetPoolData.AssetId[0]
					borrowByPoolIDAssetID, found := k.lend.GetAssetStatsByPoolIDAndAssetID(ctx, rewardsAssetPoolData.CPoolId, assetID)
					if !found {
						return nil
					}
					for _, borrowId := range borrowByPoolIDAssetID.BorrowIds {
						borrowPos, found := k.lend.GetBorrow(ctx, borrowId)
						if !found {
							continue
						}
						if borrowPos.IsLiquidated {
							continue
						}
						borrowAmt, found := k.OraclePriceForRewards(ctx, assetID, borrowPos.AmountOut.Amount)
						if !found {
							continue
						}
						lendPos, _ := k.lend.GetLend(ctx, borrowPos.LendingID)
						addr, _ := sdk.AccAddressFromBech32(lendPos.Owner)
						minAmt, found := k.CheckMinOfBorrowersLiquidityAndBorrow(ctx, addr, v.MasterPoolId, rewardsAssetPoolData.CSwapAppId, borrowAmt)
						if !found {
							continue
						}
						addrArr = append(addrArr, lendPos.Owner)
						amountArr = append(amountArr, minAmt)
						totalAmount = totalAmount.Add(minAmt.TruncateInt())
					}
					rewardAsset, found := k.asset.GetAssetForDenom(ctx, v.TotalRewards.Denom)
					if !found {
						continue
					}
					totalRewardAmt, found := k.OraclePriceForRewards(ctx, rewardAsset.Id, v.AvailableRewards.Amount)
					if !found {
						continue
					}
					if totalAmount.LTE(sdkmath.ZeroInt()) {
						continue
					}
					dailyRewardAmt := totalRewardAmt.Quo(sdkmath.LegacyNewDec(v.DurationDays - int64(epoch.Count)))
					totalAPR := dailyRewardAmt.Quo(sdkmath.LegacyNewDecFromInt(totalAmount))
					amountRewardedTracker := sdkmath.NewInt(0)
					for i, borrower := range addrArr {
						user, _ := sdk.AccAddressFromBech32(borrower)
						finalDailyRewardsPerUser := amountArr[i].Mul(totalAPR)
						if finalDailyRewardsPerUser.TruncateInt().GT(sdkmath.ZeroInt()) {
							amountRewardedTracker = amountRewardedTracker.Add(finalDailyRewardsPerUser.TruncateInt())
							err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, user, sdk.NewCoins(sdk.NewCoin(v.TotalRewards.Denom, finalDailyRewardsPerUser.TruncateInt())))
							if err != nil {
								continue
							}
						}
					}
					epoch.Count = epoch.Count + types.UInt64One
					epoch.StartingTime = timeNow + types.SecondsPerDay
					k.SetEpochTime(ctx, epoch)
					v.AvailableRewards.Amount = v.AvailableRewards.Amount.Sub(amountRewardedTracker)
					k.SetExternalRewardLend(ctx, v)
				} else {
					v.IsActive = false
					k.SetExternalRewardLend(ctx, v)
				}
			}
		}
	}
	return nil
}

func (k Keeper) CheckMinOfBorrowersLiquidityAndBorrow(ctx sdk.Context, addr sdk.AccAddress, masterPoolID int64, appID uint64, borrowAmount sdkmath.LegacyDec) (sdkmath.LegacyDec, bool) {
	farmedCoin, found := k.liquidityKeeper.GetActiveFarmer(ctx, appID, uint64(masterPoolID), addr)
	if !found {
		return sdkmath.LegacyZeroDec(), false
	}
	deserializerKit, err := k.liquidityKeeper.GetPoolTokenDesrializerKit(ctx, appID, uint64(masterPoolID))
	if err != nil {
		return sdkmath.LegacyZeroDec(), false
	}
	x, y, err := k.liquidityKeeper.CalculateXYFromPoolCoin(ctx, deserializerKit, farmedCoin.FarmedPoolCoin)
	if err != nil {
		return sdkmath.LegacyZeroDec(), false
	}

	quoteCoinAsset, _ := k.asset.GetAssetForDenom(ctx, deserializerKit.Pair.QuoteCoinDenom)
	baseCoinAsset, _ := k.asset.GetAssetForDenom(ctx, deserializerKit.Pair.BaseCoinDenom)
	priceQuoteCoin, found := k.OraclePriceForRewards(ctx, quoteCoinAsset.Id, x)
	if !found {
		return sdkmath.LegacyZeroDec(), false
	}
	priceBaseCoin, found := k.OraclePriceForRewards(ctx, baseCoinAsset.Id, y)
	if !found {
		return sdkmath.LegacyZeroDec(), false
	}

	return sdkmath.LegacyMinDec(priceQuoteCoin.Add(priceBaseCoin), borrowAmount), true
}

func (k Keeper) CombinePSMUserPositions(ctx sdk.Context) error {
	// Step 3 Elaborated
	// call app function
	// call all adddresses app wise
	// Join user positions for psm rewards that have completed the 1 day epoch
	// after combining them delete there one, ignore ones that have not completed an epoch
	// do this for all positions.

	extRewardAllAppData := k.GetAllExternalRewardStableVault(ctx)
	for _, extRewardAppData := range extRewardAllAppData {
		appStableVaultsData, found := k.vault.GetStableMintVaultRewardsByApp(ctx, extRewardAppData.AppId)
		if found {
			for _, appStableVaultData := range appStableVaultsData {
				if (uint64(ctx.BlockHeight()) - appStableVaultData.BlockHeight) > uint64(extRewardAppData.AcceptedBlockHeight) {
					// First checking if that exists or deleted
					_, found := k.vault.GetStableMintVaultRewards(ctx, appStableVaultData)
					if !found {
						continue
					}
					// using address from one user value to get all, then checking the epoch duration limit, and for those who have crossed it, joining it together.
					userStableVaultsData, found := k.vault.GetStableMintVaultUserRewards(ctx, appStableVaultData.AppId, appStableVaultData.User)
					if !found {
						continue
					}
					//****looping over the different data, but keeping in mind to ignore the one being used as initial data (appStableVaultData)****
					for _, individualVault := range userStableVaultsData {
						if ((uint64(ctx.BlockHeight()) - individualVault.BlockHeight) > uint64(extRewardAppData.AcceptedBlockHeight)) && (individualVault.BlockHeight != appStableVaultData.BlockHeight) {
							appStableVaultData.Amount = appStableVaultData.Amount.Add(individualVault.Amount)
							k.vault.DeleteStableMintVaultRewards(ctx, individualVault)
						}
					}
					k.vault.SetStableMintVaultRewards(ctx, appStableVaultData)
				}
			}
		}
	}
	return nil
}

// Stable Mint Rewards Rewards
// 1. Make a DS that take app ID for activating rewards, along with other necessary params (eg. cswap id , commodo id, else they could be 0) along with rewards quantity and epoch
// 2. Create, Deposit, Withdraw functions only save data if DS in 1. is active.
// 3. Using that 1. DS , the CombinePSMUserPositions runs for those apps and combine the rewards for addresses that have completeed  min1 epoch (app specific)
// 4. Reward function will run and check epoch deadline, (balance + lockerbal + lpFarming+ commodo)>=mint balance , then give rewards on whichever is less.

func (k Keeper) DistributeExtRewardStableVault(ctx sdk.Context) error {
	// Give external rewards to users who mint via stable vault with specific assetID
	// extRewards, _ := k.vault.GetStableMintVaultRewardsByApp(ctx, appID)
	// extRewardsProp, _ := k.GetExternalRewardStableVaultByApp(ctx, appID)
	extRewardsProp := k.GetAllExternalRewardStableVault(ctx)
	for _, extRew := range extRewardsProp {
		extRewards, _ := k.vault.GetStableMintVaultRewardsByApp(ctx, extRew.AppId)
		epoch, _ := k.GetEpochTime(ctx, extRew.EpochId)
		et := epoch.StartingTime
		timeNow := ctx.BlockTime().Unix()
		for _, userReward := range extRewards {
			extPair, _ := k.asset.GetPairsVault(ctx, userReward.StableExtendedPairId)
			pair, _ := k.asset.GetPair(ctx, extPair.PairId)
			asset, _ := k.asset.GetAsset(ctx, pair.AssetOut)
			if extRew.IsActive {
				// here the epoch starting time is set to the next day whenever any external vault reward is distributed
				// so when the epoch starting time is less than current time then the condition becomes true and flow passes through the function
				// checking if rewards are active
				if et < timeNow {
					if epoch.Count < uint64(extRew.DurationDays) { // rewards will be given till the duration defined in the ext rewards
						// initializing amountRewardedTracker to keep a track of daily rewards given to stableVault users
						amountRewardedTracker := sdk.NewCoin(extRew.TotalRewards.Denom, sdkmath.ZeroInt())
						totalRewards := extRew.AvailableRewards

						// checking if the locker was not created just to claim the external rewards, so we apply a basic check here.
						// last day don't check min lockup time, so we should have no remaining amount left
						if int64(epoch.Count) != extRew.DurationDays-1 {
							if ctx.BlockHeight()-int64(userReward.BlockHeight) < extRew.AcceptedBlockHeight {
								continue
							}
						}
						user, err := sdk.AccAddressFromBech32(userReward.User)
						if err != nil {
							return err
						}
						userBalance := k.bank.GetBalance(ctx, user, asset.Denom)                                                 // userbal
						farmedAmount, err := k.liquidityKeeper.GetAmountFarmedForAssetID(ctx, extRew.CswapAppId, asset.Id, user) // cswap farm
						if err != nil {
							farmedAmount = sdkmath.ZeroInt()
						}
						lendAsset, found := k.lend.UserAssetLends(ctx, user.String(), asset.Id) // commodo lend pos
						if !found {
							lendAsset = sdkmath.ZeroInt()
						}
						lockerAmt := sdkmath.ZeroInt()
						lockerLookupData, found := k.locker.GetUserLockerAssetMapping(ctx, user.String(), extRew.AppId, asset.Id) // locker data
						if found {
							lockerData, _ := k.locker.GetLocker(ctx, lockerLookupData.LockerId)
							lockerAmt = lockerData.NetBalance
						}

						eligibleRewardAmt := sdkmath.ZeroInt()
						if (userBalance.Amount.Add(farmedAmount).Add(lendAsset).Add(lockerAmt)).GTE(userReward.Amount) {
							eligibleRewardAmt = userReward.Amount
						} else {
							eligibleRewardAmt = userBalance.Amount.Add(farmedAmount).Add(lendAsset).Add(lockerAmt)
						}

						totalMintedData := sdkmath.ZeroInt()
						getAllExtpairData, _ := k.asset.GetPairsVaults(ctx)
						for _, stableExtPairData := range getAllExtpairData {
							if stableExtPairData.AppId == extRew.AppId && stableExtPairData.IsStableMintVault {
								appExtPairVaultData, _ := k.vault.GetAppExtendedPairVaultMappingData(ctx, stableExtPairData.AppId, stableExtPairData.Id)
								totalMintedData = totalMintedData.Add(appExtPairVaultData.TokenMintedAmount)
							}
						}

						individualUserShare := sdkmath.LegacyNewDecFromInt(eligibleRewardAmt).Quo(sdkmath.LegacyNewDecFromInt(totalMintedData)) // getting share percentage
						Duration := extRew.DurationDays - int64(epoch.Count)                                                                    // duration left (total duration - current count)
						epochRewards := (sdkmath.LegacyNewDecFromInt(totalRewards.Amount)).Quo(sdkmath.LegacyNewDec(Duration))
						dailyRewards := individualUserShare.Mul(epochRewards)
						finalDailyRewards := dailyRewards.TruncateInt()

						if finalDailyRewards.GT(sdkmath.ZeroInt()) {
							amountRewardedTracker = amountRewardedTracker.Add(sdk.NewCoin(totalRewards.Denom, finalDailyRewards))
							err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, user, sdk.NewCoins(sdk.NewCoin(totalRewards.Denom, finalDailyRewards)))
							if err != nil {
								continue
							}
						}

						// setting the available rewards by subtracting the amount sent per epoch for the ext rewards
						extRew.AvailableRewards.Amount = extRew.AvailableRewards.Amount.Sub(amountRewardedTracker.Amount)
						k.SetExternalRewardStableVault(ctx, extRew)
					} else {
						extRew.IsActive = false
						k.SetExternalRewardStableVault(ctx, extRew)
					}
				}
			}
		}
		// after all the vault owners are rewarded
		// setting the starting time to next day
		epoch.Count = epoch.Count + types.UInt64One
		epoch.StartingTime = timeNow + types.SecondsPerDay
		k.SetEpochTime(ctx, epoch)
	}
	return nil
}
