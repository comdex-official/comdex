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
	extRewards := k.GetExternalRewardsLockers(ctx)
	for i, v := range extRewards {
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
		if v.IsActive {
			epochTime, _ := k.GetEpochTime(ctx, v.EpochId)
			et := epochTime.StartingTime
			timeNow := ctx.BlockTime().Unix()

			if et < timeNow {
				if extRewards[i].IsActive {
					epoch, _ := k.GetEpochTime(ctx, v.EpochId)

					if epoch.Count < uint64(extRewards[i].DurationDays) {
						lockerLookup, _ := k.GetLockerLookupTable(ctx, v.AppMappingId, v.AssetId)
						totalShare := lockerLookup.DepositedAmount
						for _, lockerID := range lockerLookup.LockerIds {
							locker, found := k.GetLocker(ctx, lockerID)
							if !found {
								continue
							}
							userShare := (locker.NetBalance.ToDec()).Quo(totalShare.ToDec())
							totalRewards := k.GetExternalRewardsLocker(ctx, v.Id).TotalRewards
							Duration := k.GetExternalRewardsLocker(ctx, v.Id).DurationDays
							rewardsPerEpoch := (totalRewards.Amount.ToDec()).Quo(sdk.NewInt(Duration).ToDec())
							dailyRewards := userShare.Mul(rewardsPerEpoch)
							user, _ := sdk.AccAddressFromBech32(locker.Depositor)
							finalDailyRewards := sdk.NewInt(dailyRewards.TruncateInt64())
							if finalDailyRewards.GT(sdk.ZeroInt()) {
								err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, user, sdk.NewCoin(totalRewards.Denom, finalDailyRewards))
								if err != nil {
									continue
								}
							}
							epoch.Count = epoch.Count + types.UInt64One
							epoch.StartingTime = timeNow + types.Int64SecondsInADay
							k.SetEpochTime(ctx, epoch)
						}
					} else {
						extRewards[i].IsActive = false
						k.SetExternalRewardsLockers(ctx, extRewards[i])
					}
				}
			}
		}
	}
	return nil
}

func (k Keeper) DistributeExtRewardVault(ctx sdk.Context) error {
	extRewards := k.GetExternalRewardVaults(ctx)
	for i, v := range extRewards {
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
		if v.IsActive {
			epochTime, _ := k.GetEpochTime(ctx, v.EpochId)
			et := epochTime.StartingTime

			timeNow := ctx.BlockTime().Unix()

			if et < timeNow {
				if extRewards[i].IsActive {
					epoch, _ := k.GetEpochTime(ctx, v.EpochId)
					if epoch.Count < uint64(extRewards[i].DurationDays) {
						appExtPairVaultData, _ := k.GetAppExtendedPairVaultMappingData(ctx, v.AppMappingId, v.Extended_Pair_Id)
						for _, vaultID := range appExtPairVaultData.VaultIds {
							totalRewards := v.TotalRewards
							userVault, found := k.GetVault(ctx, vaultID)
							if !found {
								continue
							}
							individualUserShare := sdk.NewDec(userVault.AmountOut.Int64()).Quo(sdk.NewDec(appExtPairVaultData.CollateralLockedAmount.Int64()))
							Duration := v.DurationDays
							rewardsPerEpoch := sdk.NewDec((totalRewards.Amount).Quo(sdk.NewInt(Duration)).Int64())
							dailyRewards := individualUserShare.Mul(rewardsPerEpoch)
							finalDailyRewards := sdk.NewInt(dailyRewards.TruncateInt64())

							user, _ := sdk.AccAddressFromBech32(userVault.Owner)
							if finalDailyRewards.GT(sdk.ZeroInt()) {
								err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, user, sdk.NewCoin(totalRewards.Denom, finalDailyRewards))
								if err != nil {
									continue
								}
							}
							epoch.Count = epoch.Count + types.UInt64One
							epoch.StartingTime = timeNow + types.SecondsPerDay
							k.SetEpochTime(ctx, epoch)
						}
					} else {
						extRewards[i].IsActive = false
						k.SetExternalRewardVault(ctx, extRewards[i])
					}
				}
			}
		}
	}
	return nil
}

// calculate new locker rewards
func (k Keeper) CalculationOfRewards(
	ctx sdk.Context,
	// nolint
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
	extRewards := k.GetExternalRewardLends(ctx)
	for i, v := range extRewards {
		klwsParams, _ := k.GetKillSwitchData(ctx, v.AppMappingId)
		if klwsParams.BreakerEnable {
			return esmtypes.ErrCircuitBreakerEnabled
		}
		if v.IsActive {
			epochTime, _ := k.GetEpochTime(ctx, v.EpochId)
			et := epochTime.StartingTime

			timeNow := ctx.BlockTime().Unix()

			if et < timeNow {
				if extRewards[i].IsActive {
					epoch, _ := k.GetEpochTime(ctx, v.EpochId)
					if epoch.Count < uint64(extRewards[i].DurationDays) {
						// we will only consider the borrows of the pool and assetID defined
						totalBorrowedAmt := uint64(0)
						for _, rewardsAssetPoolData := range v.RewardsAssetPoolData {
							for _, assetID := range rewardsAssetPoolData.AssetId {
								borrowByPoolIDAssetID := k.getlbdatabypoolandassetid()
								price, err := k.CalcAssetPrice(ctx, assetID, borrowByPoolIDAssetID)
								if err != nil {
									return err
								}
								totalBorrowedAmt = totalBorrowedAmt + price
							}
						}
						totalRewardAmt, _ := k.CalcAssetPrice(ctx, v.RewardAssetId, v.TotalRewards.Amount)
						totalAPR := sdk.NewDec(int64(totalRewardAmt)).Quo(sdk.NewDec(int64(totalBorrowedAmt)))
						var inverseRatesSum sdk.Dec
						for _, rewardsAssetPoolData := range v.RewardsAssetPoolData {
							for _, assetID := range rewardsAssetPoolData.AssetId {
								inverseRate := k.InversingRates(ctx, assetID, rewardsAssetPoolData.CPoolId, totalRewardAmt)
								inverseRatesSum = inverseRatesSum.Add(inverseRate)
							}
						}

						for _, rewardsAssetPoolData := range v.RewardsAssetPoolData {
							for _, assetID := range rewardsAssetPoolData.AssetId {
								borrowIDs := k.getlbdatabypoolandassetid()
								for _, borrowID := range borrowIDs {
									borrow, _ := k.GetBorrow(ctx, borrowID)
									lend, _ := k.GetLend(ctx, borrow.LendingID)
									inverseRate := k.InversingRates(ctx, assetID, rewardsAssetPoolData.CPoolId, totalRewardAmt)
									numerator := totalAPR.Mul(inverseRate)
									finalAPR := numerator.Quo(inverseRatesSum)
									finalDailyRewardsNumerator := sdk.NewDecFromInt(borrow.AmountOut.Amount).Mul(finalAPR)
									daysInYear, _ := sdk.NewDecFromStr(types.DaysInYear)
									finalDailyRewardsPerUser := finalDailyRewardsNumerator.Quo(daysInYear)
									user, _ := sdk.AccAddressFromBech32(lend.Owner)
									if finalDailyRewardsPerUser.TruncateInt().GT(sdk.ZeroInt()) {
										err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, user, sdk.NewCoin(v.TotalRewards.Denom, finalDailyRewardsPerUser.TruncateInt()))
										if err != nil {
											continue
										}
									}
								}
							}
						}
						epoch.Count = epoch.Count + types.UInt64One
						epoch.StartingTime = timeNow + types.SecondsPerDay
						k.SetEpochTime(ctx, epoch)

					} else {
						extRewards[i].IsActive = false
						k.SetExternalRewardLend(ctx, extRewards[i])
					}
				}
			}
		}
	}
	return nil
}

func (k Keeper) InversingRates(ctx sdk.Context, assetID, poolID, totalRewardAmt uint64) sdk.Dec {
	assetBorrowedByPoolIDandAssetID := k.getlbdatabypoolandassetid(assetID, poolID)
	assetBorrowedByPoolIDandAssetIDAmt, _ := k.CalcAssetPrice(ctx, assetID, assetBorrowedByPoolIDandAssetID)
	tempRate := sdk.NewDec(int64(assetBorrowedByPoolIDandAssetIDAmt)).Quo(sdk.NewDec(int64(totalRewardAmt)))
	inverseRate := sdk.OneDec().Sub(tempRate)
	return inverseRate
}
