package keeper

import (
	"fmt"
	collectortypes "github.com/comdex-official/comdex/x/collector/types"
	lockertypes "github.com/comdex-official/comdex/x/locker/types"
	"github.com/comdex-official/comdex/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"math"
	"strconv"
)

//IterateLocker does reward calculation for locker
func (k Keeper) IterateLocker(ctx sdk.Context) error {
	rewards := k.GetRewards(ctx)
	for _, v := range rewards {
		appMappingId := v.App_mapping_ID
		assetIds := v.Asset_ID

		for i := range assetIds {

			CollectorLookup, found := k.GetCollectorLookupByAsset(ctx, appMappingId, assetIds[i])
			if !found {
				continue
			}

			LockerProductAssetMapping, _ := k.GetLockerLookupTable(ctx, appMappingId)
			lockers := LockerProductAssetMapping.Lockers
			for _, v := range lockers {
				if v.AssetId == assetIds[i] {
					lockerIds := v.LockerIds
					for w := range lockerIds {
						locker, _ := k.GetLocker(ctx, lockerIds[w])
						balance := locker.NetBalance
						rewards, err := k.CalculateRewards(ctx, balance, CollectorLookup.LockerSavingRate)
						if err != nil {
							return nil
						}
			
						// update the lock position
						returnsAcc := locker.ReturnsAccumulated
				
						updatedReturnsAcc := rewards.Add(returnsAcc)
						netBalance := locker.NetBalance.Add(rewards)
						updatedLocker := lockertypes.Locker{
							LockerId:           locker.LockerId,
							Depositor:          locker.Depositor,
							ReturnsAccumulated: updatedReturnsAcc,
							NetBalance:         netBalance,
							CreatedAt:          locker.CreatedAt,
							AssetDepositId:     locker.AssetDepositId,
							IsLocked:           locker.IsLocked,
							AppMappingId:       locker.AppMappingId,
						}
						netfeecollectedData, _ := k.GetNetFeeCollectedData(ctx, locker.AppMappingId)
						for _, p := range netfeecollectedData.AssetIdToFeeCollected {
							if p.AssetId == locker.AssetDepositId {
								asset, _ := k.GetAsset(ctx, p.AssetId)
								//updatedNetFee := p.NetFeesCollected.Sub(rewards)
								//err := k.SetNetFeeCollectedData(ctx, locker.AppMappingId, locker.AssetDepositId, updatedNetFee)
								err := k.DecreaseNetFeeCollectedData(ctx, locker.AppMappingId, locker.AssetDepositId, rewards)
								err = k.SendCoinFromModuleToModule(ctx, collectortypes.ModuleName, lockertypes.ModuleName, sdk.NewCoins(sdk.NewCoin(asset.Denom, rewards)))
								if err != nil {
									return err
								}
								if err != nil {
									return err
								}
							}
						}
						k.UpdateLocker(ctx, updatedLocker)
					}
				}
			}
		}
	}
	return nil
}

//CalculateRewards does per block rewards/interest calculation
func (k Keeper) CalculateRewards(ctx sdk.Context, amount sdk.Int, lsr sdk.Dec) (sdk.Int, error) {

	currentTime := ctx.BlockTime().Unix()

	prevInterestTime := k.GetLastInterestTime(ctx)
	if prevInterestTime == 0 {
		prevInterestTime = currentTime
	}

	secondsElapsed := currentTime - prevInterestTime
	if secondsElapsed < 0 {
		return sdk.ZeroInt(), sdkerrors.Wrap(types.ErrNegativeTimeElapsed, fmt.Sprintf("%d seconds", secondsElapsed))
	}

	//{(1+ Annual Interest Rate)^(No of seconds per block/No. of seconds in a year)}-1

	yearsElapsed := sdk.NewDec(secondsElapsed).QuoInt64(types.SecondsPerYear).MustFloat64()
	perc := lsr.String()
	a, _ := sdk.NewDecFromStr("1")
	b, _ := sdk.NewDecFromStr(perc)
	factor1 := a.Add(b).MustFloat64()
	intPerBlockfactor := math.Pow(factor1, yearsElapsed)
	intAccPerBlock := intPerBlockfactor - 1
	amtFloat, _ := strconv.ParseFloat(amount.String(), 64)
	newAmount := intAccPerBlock * amtFloat

	err := k.SetLastInterestTime(ctx, currentTime)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	return sdk.NewInt(int64(newAmount)), nil
}

//IterateVaults does interest calculation for vaults
func (k Keeper) IterateVaults(ctx sdk.Context, appMappingId uint64) error {
	extVaultMapping, _ := k.GetAppExtendedPairVaultMapping(ctx, appMappingId)
	for _, v := range extVaultMapping.ExtendedPairVaults {
		vaultIds := v.VaultIds
		for j, _ := range vaultIds {
			vault, _ := k.GetVault(ctx, vaultIds[j])
			ExtPairVault, _ := k.GetPairsVault(ctx, vault.ExtendedPairVaultID)
			StabilityFee := ExtPairVault.StabilityFee

			if StabilityFee != sdk.ZeroDec() {

				interest, _ := k.CalculateRewards(ctx, vault.AmountOut, StabilityFee)
				intAcc := vault.InterestAccumulated
				updatedIntAcc := (intAcc).Add(interest)
				vault.InterestAccumulated = updatedIntAcc
				vault.AmountOut = vault.AmountOut.Add(interest)
				k.SetVault(ctx, vault)
			}
		}
	}
	return nil
}

//DistributeExtRewardLocker does distribution of external locker rewards
func (k Keeper) DistributeExtRewardLocker(ctx sdk.Context) error {
	extRewards := k.GetExternalRewardsLockers(ctx)
	for i, v := range extRewards {
		epochTime, _ := k.GetEpochTime(ctx, v.EpochId)
		et := epochTime.StartingTime
		timeNow := ctx.BlockTime().Unix()

		if et < timeNow {

			if extRewards[i].IsActive == true {
				//count, _ := k.GetExternalRewardsLockersCounter(ctx, extRewards[i].Id)
				epoch, _ := k.GetEpochTime(ctx, v.EpochId)

				if epoch.Count < uint64(extRewards[i].DurationDays) {
					lockerLookup, _ := k.GetLockerLookupTable(ctx, v.AppMappingId)
					for _, u := range lockerLookup.Lockers {
						if u.AssetId == v.AssetId {
							lockerIds := u.LockerIds
							totalShare := u.DepositedAmount
							for w, _ := range lockerIds {
								locker, _ := k.GetLocker(ctx, lockerIds[w])
								userShare := locker.NetBalance.Quo(totalShare)
								totalRewards := k.GetExternalRewardsLocker(ctx, v.Id).TotalRewards
								Duration := k.GetExternalRewardsLocker(ctx, v.Id).DurationDays
								rewardsPerEpoch := (totalRewards.Amount).Quo(sdk.NewInt(Duration))
								dailyRewards := userShare.Mul(rewardsPerEpoch)
								user, _ := sdk.AccAddressFromBech32(locker.Depositor)
								err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, user, sdk.NewCoin(totalRewards.Denom, dailyRewards))
								if err != nil {
									return err
								}
								epoch.Count = epoch.Count + 1
								k.SetEpochTime(ctx, epoch)
							}
						}
					}
				} else {
					extRewards[i].IsActive = false
					k.SetExternalRewardsLockers(ctx, extRewards[i])
				}
			}
			epoch := types.EpochTime{
				StartingTime: et + 84600,
			}
			k.SetEpochTime(ctx, epoch)
		}

	}
	return nil
}

func (k Keeper) DistributeExtRewardVault(ctx sdk.Context) error {
	extRewards := k.GetExternalRewardVaults(ctx)
	for i, v := range extRewards {
		epochTime, _ := k.GetEpochTime(ctx, v.EpochId)
		et := epochTime.StartingTime
		timeNow := ctx.BlockTime().Unix()

		if et < timeNow {
			if extRewards[i].IsActive == true {
				epoch, _ := k.GetEpochTime(ctx, v.EpochId)
				if epoch.Count < uint64(extRewards[i].DurationDays) {
					appExtPairVaultData, _ := k.GetAppExtendedPairVaultMapping(ctx, v.AppMappingId)
					for _, u := range appExtPairVaultData.ExtendedPairVaults {
						for _, w := range u.VaultIds {
							totalRewards := v.TotalRewards
							userVault, _ := k.GetVault(ctx, w)
							userShare := userVault.AmountOut.Quo(u.CollateralLockedAmount)
							Duration := v.DurationDays
							rewardsPerEpoch := (totalRewards.Amount).Quo(sdk.NewInt(Duration))
							dailyRewards := userShare.Mul(rewardsPerEpoch)
							user, _ := sdk.AccAddressFromBech32(userVault.Owner)
							err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, user, sdk.NewCoin(totalRewards.Denom, dailyRewards))
							if err != nil {
								return err
							}
							epoch.Count = epoch.Count + 1
							k.SetEpochTime(ctx, epoch)
						}
					}
				} else {
					extRewards[i].IsActive = false
					k.SetExternalRewardVault(ctx, extRewards[i])
				}
			}
			epoch := types.EpochTime{
				StartingTime: et + 84600,
			}
			k.SetEpochTime(ctx, epoch)
		}
	}

	return nil
}
