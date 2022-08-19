package keeper

import (
	"fmt"
	"math"
	"strconv"

	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	lockertypes "github.com/comdex-official/comdex/x/locker/types"
	"github.com/comdex-official/comdex/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

//IterateLocker does reward calculation for locker.
func (k Keeper) IterateLocker(ctx sdk.Context) error {
	rewards := k.GetRewards(ctx)
	for _, r := range rewards {
		klwsParams, _ := k.GetKillSwitchData(ctx, r.App_mapping_ID)
		if klwsParams.BreakerEnable {
			return esmtypes.ErrCircuitBreakerEnabled
		}
		esmStatus, found := k.GetESMStatus(ctx, r.App_mapping_ID)
		status := false
		if found {
			status = esmStatus.Status
		}
		if status {
			return esmtypes.ErrESMAlreadyExecuted
		}
		appMappingID := r.App_mapping_ID
		assetIds := r.Asset_ID

		for i := range assetIds {
			CollectorLookup, found := k.GetCollectorLookupTable(ctx, appMappingID, assetIds[i])
			if !found {
				continue
			}
			lockers := k.GetLockers(ctx)
			for _, v := range lockers {
				if v.AppId == appMappingID {
					locker := v
					balance := locker.NetBalance
					reward, err := k.CalculateRewards(ctx, balance, CollectorLookup.LockerSavingRate)
					if err != nil {
						continue
					}

					// update the lock position
					returnsAcc := locker.ReturnsAccumulated

					lockerRewardsTracker, found := k.GetLockerRewardTracker(ctx, locker.LockerId, r.App_mapping_ID)
					if !found {
						lockerRewardsTracker = types.LockerRewardsTracker{
							LockerId:           locker.LockerId,
							AppMappingId:       r.App_mapping_ID,
							RewardsAccumulated: sdk.ZeroDec(),
						}
					}
					lockerRewardsTracker.RewardsAccumulated = lockerRewardsTracker.RewardsAccumulated.Add(reward)
					newReward := sdk.ZeroInt()
					if lockerRewardsTracker.RewardsAccumulated.GTE(sdk.OneDec()) {
						newReward = lockerRewardsTracker.RewardsAccumulated.TruncateInt()
						newRewardDec := sdk.NewDec(newReward.Int64())
						lockerRewardsTracker.RewardsAccumulated = lockerRewardsTracker.RewardsAccumulated.Sub(newRewardDec)
					}
					k.SetLockerRewardTracker(ctx, lockerRewardsTracker)
					locker.NetBalance = locker.NetBalance.Add(newReward)
					locker.ReturnsAccumulated = returnsAcc.Add(newReward)
					netFeeCollectedData, found := k.GetNetFeeCollectedData(ctx, locker.AppId, locker.AssetDepositId)
					if !found {
						continue
					}
					asset, _ := k.GetAsset(ctx, locker.AssetDepositId)
					err = k.DecreaseNetFeeCollectedData(ctx, locker.AppId, locker.AssetDepositId, newReward, netFeeCollectedData)
					if err != nil {
						continue
					}
					if newReward.GT(sdk.ZeroInt()) {
						err = k.SendCoinFromModuleToModule(ctx, collectortypes.ModuleName, lockertypes.ModuleName, sdk.NewCoins(sdk.NewCoin(asset.Denom, newReward)))
						if err != nil {
							continue
						}
					}
					lockerRewardsMapping, found := k.GetLockerTotalRewardsByAssetAppWise(ctx, appMappingID, locker.AssetDepositId)
					if !found {
						var lockerReward lockertypes.LockerTotalRewardsByAssetAppWise
						lockerReward.AppId = locker.AppId
						lockerReward.AssetId = locker.AssetDepositId
						lockerReward.TotalRewards = sdk.ZeroInt().Add(newReward)
						err = k.SetLockerTotalRewardsByAssetAppWise(ctx, lockerReward)
						if err != nil {
							continue
						}
					} else {
						lockerRewardsMapping.TotalRewards = lockerRewardsMapping.TotalRewards.Add(newReward)

						err = k.SetLockerTotalRewardsByAssetAppWise(ctx, lockerRewardsMapping)
						if err != nil {
							continue
						}
					}
					k.UpdateLocker(ctx, locker)
				}
			}
		}
	}
	return nil
}

//CalculateRewards does per block rewards/interest calculation .
func (k Keeper) CalculateRewards(
	ctx sdk.Context,
	// nolint
	amount sdk.Int, lsr sdk.Dec,
) (sdk.Dec, error) {
	currentTime := ctx.BlockTime().Unix()

	prevInterestTime := k.GetLastInterestTime(ctx)
	if prevInterestTime == types.Int64Zero {
		prevInterestTime = currentTime
	}
	secondsElapsed := currentTime - prevInterestTime
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

//IterateVaults does interest calculation for vaults .
func (k Keeper) IterateVaults(ctx sdk.Context, appMappingID uint64) error {
	extVaultMapping, found := k.GetAppMappingData(ctx, appMappingID)
	if !found {
		return types.ErrVaultNotFound
	}
	klwsParams, _ := k.GetKillSwitchData(ctx, appMappingID)
	if klwsParams.BreakerEnable {
		return esmtypes.ErrCircuitBreakerEnabled
	}
	esmStatus, found := k.GetESMStatus(ctx, appMappingID)
	status := false
	if found {
		status = esmStatus.Status
	}
	if status {
		return esmtypes.ErrESMAlreadyExecuted
	}
	for _, v := range extVaultMapping {
		vaultIds := v.VaultIds
		for j := range vaultIds {
			vault, found := k.GetVault(ctx, vaultIds[j])
			if !found {
				continue
			}
			ExtPairVault, _ := k.GetPairsVault(ctx, vault.ExtendedPairVaultID)
			StabilityFee := ExtPairVault.StabilityFee
			if vault.ExtendedPairVaultID != types.UInt64Zero {
				if StabilityFee.GT(sdk.ZeroDec()) {
					interest, err := k.CalculateRewards(ctx, vault.AmountOut, StabilityFee)
					if err != nil {
						continue
					}
					vaultInterestTracker, found := k.GetVaultInterestTracker(ctx, vault.ExtendedPairVaultID, appMappingID)
					if !found {
						vaultInterestTracker = types.VaultInterestTracker{
							VaultId:             vault.ExtendedPairVaultID,
							AppMappingId:        appMappingID,
							InterestAccumulated: sdk.ZeroDec(),
						}
					}
					vaultInterestTracker.InterestAccumulated = vaultInterestTracker.InterestAccumulated.Add(interest)
					newInterest := sdk.ZeroInt()
					if vaultInterestTracker.InterestAccumulated.GTE(sdk.OneDec()) {
						newInterest = vaultInterestTracker.InterestAccumulated.TruncateInt()
						newInterestDec := sdk.NewDec(newInterest.Int64())
						vaultInterestTracker.InterestAccumulated = vaultInterestTracker.InterestAccumulated.Sub(newInterestDec)
					}
					k.SetVaultInterestTracker(ctx, vaultInterestTracker)
					intAcc := vault.InterestAccumulated
					updatedIntAcc := (intAcc).Add(newInterest)
					vault.InterestAccumulated = updatedIntAcc
					k.SetVault(ctx, vault)
				}
			}
		}
	}
	return nil
}

//DistributeExtRewardLocker does distribution of external locker rewards .
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
					//count, _ := k.GetExternalRewardsLockersCounter(ctx, extRewards[i].Id)
					epoch, _ := k.GetEpochTime(ctx, v.EpochId)

					if epoch.Count < uint64(extRewards[i].DurationDays) {
						lockerLookup, _ := k.GetLockerLookupTableByApp(ctx, v.AppMappingId)
						for _, u := range lockerLookup {
							if u.AssetId == v.AssetId {
								lockerIds := u.LockerIds
								totalShare := u.DepositedAmount
								for w := range lockerIds {
									locker, _ := k.GetLocker(ctx, lockerIds[w])

									userShare := (locker.NetBalance.ToDec()).Quo(totalShare.ToDec())
									totalRewards := k.GetExternalRewardsLocker(ctx, v.Id).TotalRewards
									Duration := k.GetExternalRewardsLocker(ctx, v.Id).DurationDays
									rewardsPerEpoch := (totalRewards.Amount.ToDec()).Quo(sdk.NewInt(Duration).ToDec())
									dailyRewards := userShare.Mul(rewardsPerEpoch)
									user, _ := sdk.AccAddressFromBech32(locker.Depositor)
									finalDailyRewards := sdk.NewInt(dailyRewards.TruncateInt64())

									err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, user, sdk.NewCoin(totalRewards.Denom, finalDailyRewards))
									if err != nil {
										continue
									}
									epoch.Count = epoch.Count + types.UInt64One
									epoch.StartingTime = timeNow + types.Int64SecondsInADay
									k.SetEpochTime(ctx, epoch)
								}
							}
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
						appExtPairVaultData, _ := k.GetAppMappingData(ctx, v.AppMappingId)
						for _, u := range appExtPairVaultData {
							for _, w := range u.VaultIds {
								totalRewards := v.TotalRewards
								userVault, _ := k.GetVault(ctx, w)

								individualUserShare := sdk.NewDec(userVault.AmountOut.Int64()).Quo(sdk.NewDec(u.CollateralLockedAmount.Int64()))
								Duration := v.DurationDays
								rewardsPerEpoch := sdk.NewDec((totalRewards.Amount).Quo(sdk.NewInt(Duration)).Int64())
								dailyRewards := individualUserShare.Mul(rewardsPerEpoch)
								finalDailyRewards := sdk.NewInt(dailyRewards.TruncateInt64())

								user, _ := sdk.AccAddressFromBech32(userVault.Owner)
								err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, user, sdk.NewCoin(totalRewards.Denom, finalDailyRewards))
								if err != nil {
									continue
								}
								epoch.Count = epoch.Count + types.UInt64One
								epoch.StartingTime = timeNow + types.SecondsPerDay
								k.SetEpochTime(ctx, epoch)
							}
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
