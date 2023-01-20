package keeper

import (
	"fmt"
	"math"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

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
					amountRewardedTracker := sdk.NewCoin(v.TotalRewards.Denom, sdk.ZeroInt())
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
						userShare := (locker.NetBalance.ToDec()).Quo(totalShare.ToDec()) // getting share percentage
						availableRewards := v.AvailableRewards                           // Available Rewards
						Duration := v.DurationDays - int64(epoch.Count)                  // duration left (total duration - current count)

						epochRewards := availableRewards.Amount.ToDec().Quo(sdk.NewDec(Duration))
						dailyRewards := userShare.Mul(epochRewards)
						user, _ := sdk.AccAddressFromBech32(locker.Depositor)
						finalDailyRewards := dailyRewards.TruncateInt()
						// after calculating final daily rewards, the amount is sent to the user
						if finalDailyRewards.GT(sdk.ZeroInt()) {
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
					amountRewardedTracker := sdk.NewCoin(v.TotalRewards.Denom, sdk.ZeroInt())

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
						individualUserShare := userVault.AmountOut.ToDec().Quo(sdk.NewDecFromInt(appExtPairVaultData.TokenMintedAmount)) // getting share percentage
						Duration := v.DurationDays - int64(epoch.Count)                                                                  // duration left (total duration - current count)
						epochRewards := (totalRewards.Amount.ToDec()).Quo(sdk.NewDec(Duration))
						dailyRewards := individualUserShare.Mul(epochRewards)
						finalDailyRewards := dailyRewards.TruncateInt()

						user, _ := sdk.AccAddressFromBech32(userVault.Owner)
						if finalDailyRewards.GT(sdk.ZeroInt()) {
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
	amount sdk.Int, lsr sdk.Dec, bTime int64,
) (sdk.Dec, error) {
	currentTime := ctx.BlockTime().Unix()

	secondsElapsed := currentTime - bTime
	if secondsElapsed < types.Int64Zero {
		return sdk.ZeroDec(), sdkerrors.Wrap(types.ErrNegativeTimeElapsed, fmt.Sprintf("%d seconds", secondsElapsed))
	}
	//{(1+ Annual Interest Rate)^(No of seconds per block/No. of seconds in a year)}-1

	yearsElapsed := sdk.NewDec(secondsElapsed).QuoInt64(types.SecondsPerYear)
	perc := lsr.String()
	a, _ := sdk.NewDecFromStr("1")
	b, _ := sdk.NewDecFromStr(perc)
	factor1 := a.Add(b)
	intPerBlockFactor := math.Pow(factor1.MustFloat64(), yearsElapsed.MustFloat64())
	intAccPerBlock := intPerBlockFactor - types.Float64One
	amtFloat := amount.ToDec().MustFloat64()
	newAmount := intAccPerBlock * amtFloat

	// s := fmt.Sprint(newAmount)
	s := strconv.FormatFloat(newAmount, 'f', 18, 64)
	newAm, err := sdk.NewDecFromStr(s)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	return newAm, nil
}

func (k Keeper) OraclePriceForRewards(ctx sdk.Context, id uint64, amt sdk.Int) (sdk.Dec, bool) {
	asset, found := k.asset.GetAsset(ctx, id)
	if !found {
		return sdk.ZeroDec(), false
	}

	price, found := k.marketKeeper.GetTwa(ctx, asset.Id)
	if !found {
		return sdk.ZeroDec(), false
	}

	// if price is not active and twa is 0 return false
	if !price.IsPriceActive && price.Twa == 0 {
		return sdk.ZeroDec(), false
	}
	// if price is not active and DiscardedHeightDiff is not -1
	if price.DiscardedHeightDiff != -1 {
		priceInactiveBlockCount := ctx.BlockHeight() - price.DiscardedHeightDiff
		// if price is inactive since 600 block and also twa is 0 return error else continue with the old price
		if priceInactiveBlockCount >= types.DefaultAllowedBlocksForPriceInactive {
			return sdk.ZeroDec(), false
		}
	}

	numerator := sdk.NewDecFromInt(amt).Mul(sdk.NewDecFromInt(sdk.NewIntFromUint64(price.Twa)))
	denominator := sdk.NewDecFromInt(asset.Decimals)
	return numerator.Quo(denominator), true
}

func (k Keeper) DistributeExtRewardLend(ctx sdk.Context) error {
	// Give external rewards to borrowers for opening a vault with specific assetID
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
					// we will only consider the borrows of the pool and assetID defined
					// initializing totalBorrowedAmt $ value to store total borrowed across different assetIDs for given cPool
					totalBorrowedAmt := sdk.ZeroInt()
					rewardsAssetPoolData := v.RewardsAssetPoolData
					for _, assetID := range rewardsAssetPoolData.AssetId {
						amt, _ := k.CalculateTotalBorrowedAmtByFarmers(ctx, assetID, rewardsAssetPoolData.CPoolId, rewardsAssetPoolData.CSwapAppId, v.MasterPoolId)
						totalBorrowedAmt = totalBorrowedAmt.Add(amt.TruncateInt()) // in $USD
					}
					// calculating totalAPR
					rewardAsset, found := k.asset.GetAssetForDenom(ctx, v.TotalRewards.Denom)
					if !found {
						continue
					}
					totalRewardAmt, found := k.OraclePriceForRewards(ctx, rewardAsset.Id, v.AvailableRewards.Amount)
					if !found {
						continue
					}
					if totalBorrowedAmt.LTE(sdk.ZeroInt()) {
						continue
					}
					dailyRewardAmt := totalRewardAmt.Quo(sdk.NewDec(v.DurationDays - int64(epoch.Count)))
					totalAPR := dailyRewardAmt.Quo(sdk.NewDecFromInt(totalBorrowedAmt))

					// initializing amountRewardedTracker to keep a track of daily rewards given to locker owners
					amountRewardedTracker := sdk.NewCoin(v.TotalRewards.Denom, sdk.ZeroInt())

					for _, assetID := range rewardsAssetPoolData.AssetId { // iterating over assetIDs
						borrowIDs, _ := k.lend.GetAssetStatsByPoolIDAndAssetID(ctx, rewardsAssetPoolData.CPoolId, assetID)
						for _, borrowID := range borrowIDs.BorrowIds { // iterating over borrowIDs
							borrow, found := k.lend.GetBorrow(ctx, borrowID)
							if !found {
								continue
							}
							lend, found := k.lend.GetLend(ctx, borrow.LendingID)
							if !found {
								continue
							}
							user, _ := sdk.AccAddressFromBech32(lend.Owner)
							pair, found := k.lend.GetLendPair(ctx, borrow.PairID)
							if !found {
								continue
							}
							borrowAmt, found := k.OraclePriceForRewards(ctx, pair.AssetOut, borrow.AmountOut.Amount)
							if !found {
								continue
							}
							liqFound, found := k.CheckMinOfBorrowersLiquidityAndBorrow(ctx, user, v.MasterPoolId, rewardsAssetPoolData.CSwapAppId, borrowAmt)
							if !found {
								continue
							}

							finalDailyRewardsPerUser := liqFound.Mul(totalAPR)
							if finalDailyRewardsPerUser.TruncateInt().GT(sdk.ZeroInt()) {
								amountRewardedTracker = amountRewardedTracker.Add(sdk.NewCoin(v.TotalRewards.Denom, finalDailyRewardsPerUser.TruncateInt()))
								err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, user, sdk.NewCoins(sdk.NewCoin(v.TotalRewards.Denom, finalDailyRewardsPerUser.TruncateInt())))
								if err != nil {
									continue
								}
							}
						}
					}
					// after all the borrowers are rewarded
					// setting the starting time to next day
					epoch.Count = epoch.Count + types.UInt64One
					epoch.StartingTime = timeNow + types.SecondsPerDay
					k.SetEpochTime(ctx, epoch)

					// setting the available rewards by subtracting the amount sent per epoch for the ext rewards
					v.AvailableRewards.Amount = v.AvailableRewards.Amount.Sub(amountRewardedTracker.Amount)
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

func (k Keeper) CalculateTotalBorrowedAmtByFarmers(ctx sdk.Context, assetID, poolID, appID uint64, masterPoolID int64) (sdk.Dec, bool) {
	borrowByPoolIDAssetID, found := k.lend.GetAssetStatsByPoolIDAndAssetID(ctx, poolID, assetID)
	if !found {
		return sdk.ZeroDec(), false
	}

	amt := sdk.ZeroDec()
	for _, id := range borrowByPoolIDAssetID.BorrowIds {
		borrowPos, found := k.lend.GetBorrow(ctx, id)
		if !found {
			return sdk.ZeroDec(), false
		}
		lendPos, found := k.lend.GetLend(ctx, borrowPos.LendingID)
		if !found {
			return sdk.ZeroDec(), false
		}
		if borrowPos.IsLiquidated {
			continue
		}
		pair, found := k.lend.GetLendPair(ctx, borrowPos.PairID)
		if !found {
			return sdk.ZeroDec(), false
		}
		borrowAmt, found := k.OraclePriceForRewards(ctx, pair.AssetOut, borrowPos.AmountOut.Amount)
		if !found {
			return sdk.ZeroDec(), false
		}
		addr, _ := sdk.AccAddressFromBech32(lendPos.Owner)
		minAmt, found := k.CheckMinOfBorrowersLiquidityAndBorrow(ctx, addr, masterPoolID, appID, borrowAmt)
		if !found {
			continue
		}
		amt = amt.Add(minAmt)
	}

	return amt, true
}

func (k Keeper) CheckMinOfBorrowersLiquidityAndBorrow(ctx sdk.Context, addr sdk.AccAddress, masterPoolID int64, appID uint64, borrowAmount sdk.Dec) (sdk.Dec, bool) {
	farmedCoin, found := k.liquidityKeeper.GetActiveFarmer(ctx, appID, uint64(masterPoolID), addr)
	if !found {
		return sdk.ZeroDec(), false
	}
	pool, pair, ammPool, err := k.liquidityKeeper.GetAMMPoolInterfaceObject(ctx, appID, uint64(masterPoolID))
	if err != nil {
		return sdk.ZeroDec(), false
	}
	poolCoin := sdk.NewCoin(pool.PoolCoinDenom, farmedCoin.FarmedPoolCoin.Amount)
	x, y, err := k.liquidityKeeper.CalculateXYFromPoolCoin(ctx, ammPool, poolCoin)
	if err != nil {
		return sdk.ZeroDec(), false
	}

	quoteCoinAsset, _ := k.asset.GetAssetForDenom(ctx, pair.QuoteCoinDenom)
	baseCoinAsset, _ := k.asset.GetAssetForDenom(ctx, pair.BaseCoinDenom)
	priceQuoteCoin, found := k.OraclePriceForRewards(ctx, quoteCoinAsset.Id, x)
	if !found {
		return sdk.ZeroDec(), false
	}
	priceBaseCoin, found := k.OraclePriceForRewards(ctx, baseCoinAsset.Id, y)
	if !found {
		return sdk.ZeroDec(), false
	}

	return sdk.MinDec(priceQuoteCoin.Add(priceBaseCoin), borrowAmount), true
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
						amountRewardedTracker := sdk.NewCoin(extRew.TotalRewards.Denom, sdk.ZeroInt())
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
							farmedAmount = sdk.ZeroInt()
						}
						lendAsset, found := k.lend.UserAssetLends(ctx, user.String(), asset.Id) // commodo lend pos
						if !found {
							lendAsset = sdk.ZeroInt()
						}
						lockerAmt := sdk.ZeroInt()
						lockerLookupData, found := k.locker.GetUserLockerAssetMapping(ctx, user.String(), extRew.AppId, asset.Id) // locker data
						if found {
							lockerData, _ := k.locker.GetLocker(ctx, lockerLookupData.LockerId)
							lockerAmt = lockerData.NetBalance
						}

						eligibleRewardAmt := sdk.ZeroInt()
						if (userBalance.Amount.Add(farmedAmount).Add(lendAsset).Add(lockerAmt)).GTE(userReward.Amount) {
							eligibleRewardAmt = userReward.Amount
						} else {
							eligibleRewardAmt = userBalance.Amount.Add(farmedAmount).Add(lendAsset).Add(lockerAmt)
						}

						totalMintedData := sdk.ZeroInt()
						getAllExtpairData, _ := k.asset.GetPairsVaults(ctx)
						for _, stableExtPairData := range getAllExtpairData {
							if stableExtPairData.AppId == extRew.AppId && stableExtPairData.IsStableMintVault {
								appExtPairVaultData, _ := k.vault.GetAppExtendedPairVaultMappingData(ctx, stableExtPairData.AppId, stableExtPairData.Id)
								totalMintedData = totalMintedData.Add(appExtPairVaultData.TokenMintedAmount)
							}
						}

						individualUserShare := eligibleRewardAmt.ToDec().Quo(sdk.NewDecFromInt(totalMintedData)) // getting share percentage
						Duration := extRew.DurationDays - int64(epoch.Count)                                     // duration left (total duration - current count)
						epochRewards := (totalRewards.Amount.ToDec()).Quo(sdk.NewDec(Duration))
						dailyRewards := individualUserShare.Mul(epochRewards)
						finalDailyRewards := dailyRewards.TruncateInt()

						if finalDailyRewards.GT(sdk.ZeroInt()) {
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
