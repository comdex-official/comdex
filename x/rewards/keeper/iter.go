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
		klwsParams, _ := k.GetKillSwitchData(ctx, v.AppMappingId)
		if klwsParams.BreakerEnable {
			return esmtypes.ErrCircuitBreakerEnabled
		}
		esmStatus, found := k.GetESMStatus(ctx, v.AppMappingId)
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
					lockerLookup, _ := k.GetLockerLookupTable(ctx, v.AppMappingId, v.AssetId)
					totalShare := lockerLookup.DepositedAmount

					// initializing amountRewardedTracker to keep a track of daily rewards given to locker owners
					amountRewardedTracker := sdk.NewCoin(v.TotalRewards.Denom, sdk.ZeroInt())
					for _, lockerID := range lockerLookup.LockerIds {
						locker, found := k.GetLocker(ctx, lockerID)
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
							amountRewardedTracker = amountRewardedTracker.Add(sdk.NewCoin(availableRewards.Denom, finalDailyRewards))
							err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, user, sdk.NewCoin(availableRewards.Denom, finalDailyRewards))
							if err != nil {
								continue
							}
						}
					}
					// after all the locker owners are rewarded
					// setting the starting time to next day
					epoch.Count = epoch.Count + types.UInt64One
					epoch.StartingTime = timeNow + types.Int64SecondsInADay
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
		klwsParams, _ := k.GetKillSwitchData(ctx, v.AppMappingId)
		if klwsParams.BreakerEnable {
			return esmtypes.ErrCircuitBreakerEnabled
		}
		esmStatus, found := k.GetESMStatus(ctx, v.AppMappingId)
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
					appExtPairVaultData, _ := k.GetAppExtendedPairVaultMappingData(ctx, v.AppMappingId, v.Extended_Pair_Id)

					// initializing amountRewardedTracker to keep a track of daily rewards given to locker owners
					amountRewardedTracker := sdk.NewCoin(v.TotalRewards.Denom, sdk.ZeroInt())

					for _, vaultID := range appExtPairVaultData.VaultIds {
						totalRewards := v.AvailableRewards
						userVault, found := k.GetVault(ctx, vaultID)
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
							err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, user, sdk.NewCoin(totalRewards.Denom, finalDailyRewards))
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

	yearsElapsed := sdk.NewDec(secondsElapsed).QuoInt64(types.SecondsPerYear).MustFloat64()
	perc := lsr.String()
	a, _ := sdk.NewDecFromStr("1")
	b, _ := sdk.NewDecFromStr(perc)
	factor1 := a.Add(b).MustFloat64()
	intPerBlockFactor := math.Pow(factor1, yearsElapsed)
	intAccPerBlock := intPerBlockFactor - types.Float64One
	amtFloat, _ := strconv.ParseFloat(amount.String(), 64)
	newAmount := intAccPerBlock * amtFloat

	s := fmt.Sprint(newAmount)
	newAm, err := sdk.NewDecFromStr(s)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	return newAm, nil
}

func (k Keeper) DistributeExtRewardLend(ctx sdk.Context) error {
	// Give external rewards to borrowers for opening a vault with specific assetID
	extRewards := k.GetExternalRewardLends(ctx)
	for _, v := range extRewards {
		klwsParams, _ := k.GetKillSwitchData(ctx, v.AppMappingId)
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
						borrowByPoolIDAssetID, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, rewardsAssetPoolData.CPoolId, assetID)
						price, err := k.CalcAssetPrice(ctx, assetID, borrowByPoolIDAssetID.TotalBorrowed.Add(borrowByPoolIDAssetID.TotalStableBorrowed))
						if err != nil {
							return err
						}
						totalBorrowedAmt = totalBorrowedAmt.Add(price)
					}
					// calculating totalAPR
					rewardAsset, found := k.GetAssetForDenom(ctx, v.TotalRewards.Denom)
					if !found {
						continue
					}
					totalRewardAmt, _ := k.CalcAssetPrice(ctx, rewardAsset.Id, v.TotalRewards.Amount)
					totalAPR := sdk.NewDecFromInt(totalRewardAmt).Quo(sdk.NewDecFromInt(totalBorrowedAmt))
					var inverseRatesSum sdk.Dec
					// inverting the rate to enable low apr for assets which are more borrowed
					for _, assetID := range rewardsAssetPoolData.AssetId {
						inverseRate := k.InvertingRates(ctx, assetID, rewardsAssetPoolData.CPoolId, totalRewardAmt)
						inverseRatesSum = inverseRatesSum.Add(inverseRate)
					}

					// initializing amountRewardedTracker to keep a track of daily rewards given to locker owners
					amountRewardedTracker := sdk.NewCoin(v.TotalRewards.Denom, sdk.ZeroInt())

					for _, assetID := range rewardsAssetPoolData.AssetId { // iterating over assetIDs
						borrowIDs, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, rewardsAssetPoolData.CPoolId, assetID)
						for _, borrowID := range borrowIDs.BorrowIds { // iterating over borrowIDs
							borrow, found := k.GetBorrow(ctx, borrowID)
							if !found {
								continue
							}
							lend, found := k.GetLend(ctx, borrow.LendingID)
							if !found {
								continue
							}
							user, _ := sdk.AccAddressFromBech32(lend.Owner)
							liqFound := k.CheckBorrowersLiquidity(ctx, user, v.MasterPoolId, rewardsAssetPoolData.CSwapAppId, sdk.NewIntFromUint64(rewardsAssetPoolData.CSwapMinLockAmount))
							if !liqFound {
								continue
							}
							inverseRate := k.InvertingRates(ctx, assetID, rewardsAssetPoolData.CPoolId, totalRewardAmt)
							numerator := totalAPR.Mul(inverseRate)
							finalAPR := numerator.Quo(inverseRatesSum)
							finalDailyRewardsNumerator := sdk.NewDecFromInt(borrow.AmountOut.Amount).Mul(finalAPR)
							daysInYear, _ := sdk.NewDecFromStr(types.DaysInYear)
							finalDailyRewardsPerUser := finalDailyRewardsNumerator.Quo(daysInYear)

							if finalDailyRewardsPerUser.TruncateInt().GT(sdk.ZeroInt()) {
								amountRewardedTracker = amountRewardedTracker.Sub(sdk.NewCoin(v.TotalRewards.Denom, finalDailyRewardsPerUser.TruncateInt()))
								err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, user, sdk.NewCoin(v.TotalRewards.Denom, finalDailyRewardsPerUser.TruncateInt()))
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

func (k Keeper) InvertingRates(ctx sdk.Context, assetID, poolID uint64, totalRewardAmt sdk.Int) sdk.Dec {
	assetBorrowedByPoolIDAndAssetID, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, poolID, assetID)
	assetBorrowedByPoolIDAndAssetIDAmt, _ := k.CalcAssetPrice(ctx, assetID, assetBorrowedByPoolIDAndAssetID.TotalBorrowed.Add(assetBorrowedByPoolIDAndAssetID.TotalStableBorrowed))
	tempRate := sdk.NewDecFromInt(assetBorrowedByPoolIDAndAssetIDAmt).Quo(sdk.NewDecFromInt(totalRewardAmt))
	inverseRate := sdk.OneDec().Sub(tempRate)
	return inverseRate
}

func (k Keeper) CheckBorrowersLiquidity(ctx sdk.Context, addr sdk.AccAddress, masterPoolID, appID uint64, amount sdk.Int) bool {
	farmedCoin, found := k.GetActiveFarmer(ctx, appID, masterPoolID, addr)
	if !found {
		return false
	}

	pool, pair, ammPool, err := k.GetAMMPoolInterfaceObject(ctx, appID, masterPoolID)
	if err != nil {
		return false
	}
	poolCoin := sdk.NewCoin(pool.PoolCoinDenom, farmedCoin.FarmedPoolCoin.Amount)
	x, y, err := k.CalculateXYFromPoolCoin(ctx, ammPool, poolCoin)
	if err != nil {
		return false
	}

	quoteCoinAsset, _ := k.GetAssetForDenom(ctx, pair.QuoteCoinDenom)
	baseCoinAsset, _ := k.GetAssetForDenom(ctx, pair.BaseCoinDenom)
	priceQuoteCoin, err := k.CalcAssetPrice(ctx, quoteCoinAsset.Id, x)
	if err != nil {
		return false
	}
	priceBaseCoin, err := k.CalcAssetPrice(ctx, baseCoinAsset.Id, y)
	if err != nil {
		return false
	}
	if priceQuoteCoin.Add(priceBaseCoin).GTE(amount) {
		return true
	}
	return false
}
