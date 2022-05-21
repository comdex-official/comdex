package keeper

import (
	"fmt"
	lockertypes "github.com/comdex-official/comdex/x/locker/types"
	"github.com/comdex-official/comdex/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) Iterate(ctx sdk.Context, appMappingId uint64, assetIds []uint64) error {
	CollectorAppAsset, _ := k.GetAppToDenomsMapping(ctx, appMappingId)
	for i := range assetIds {
		found := uint64InSlice(assetIds[i], CollectorAppAsset.AssetIds)
		if !found {
			return types.ErrAssetIdDoesNotExist
		}
		CollectorLookup, _ := k.GetCollectorLookupByAsset(ctx, appMappingId, assetIds[i])
		for _, j := range CollectorLookup.AssetrateInfo {
			LockerProductAssetMapping, _ := k.GetLockerLookupTable(ctx, appMappingId)
			lockers := LockerProductAssetMapping.Lockers
			for _, v := range lockers {
				if v.AssetId == assetIds[i] {
					lockerIds := v.LockerIds
					for w := range lockerIds {
						locker, _ := k.GetLocker(ctx, lockerIds[w])
						balance := locker.NetBalance
						rewards, err := k.CalculateRewards(ctx, balance, *j.LockerSavingRate)
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
						k.UpdateLocker(ctx, updatedLocker)
					}
				}
			}
		}
	}
	return nil
}

func (k Keeper) CalculateRewards(ctx sdk.Context, amount sdk.Int, LockerSavingsRate sdk.Dec) (sdk.Int, error) {

	currentTime := ctx.BlockTime().Unix()

	prevInterestTime := k.GetLastInterestTime(ctx)
	if prevInterestTime == 0 {
		prevInterestTime = currentTime
	}

	secondsElapsed := currentTime - prevInterestTime
	if secondsElapsed < 0 {
		return sdk.ZeroInt(), sdkerrors.Wrap(types.ErrNegativeTimeElapsed, fmt.Sprintf("%d seconds", secondsElapsed))
	}
	yearsElapsed := sdk.NewDec(secondsElapsed).QuoInt64(types.SecondsPerYear)

	newAmount := sdk.NewDecFromInt(amount.Mul(sdk.Int(LockerSavingsRate))).Mul(yearsElapsed).QuoInt64(100)

	err := k.SetLastInterestTime(ctx, currentTime)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	return sdk.Int(newAmount), nil
}
