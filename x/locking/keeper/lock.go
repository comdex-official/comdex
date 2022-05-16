package keeper

import (
	"time"

	"github.com/comdex-official/comdex/x/locking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetAccountWithSameLockDurationAndDenom(ctx sdk.Context, owner sdk.AccAddress, denom string, duration time.Duration) []types.Lock {

	locks := []types.Lock{}

	locksByOwner, found := k.GetLockByOwner(ctx, owner.String())
	if !found {
		return locks
	}

	for _, lockId := range locksByOwner.LockIds {
		lock, found := k.GetLockById(ctx, lockId)
		if !found {
			continue
		}
		if lock.Coin.Denom == denom && lock.Duration == duration {
			locks = append(locks, lock)
		}
	}

	return locks
}

// AddTokensToLock locks more tokens into a bonding
// This also saves the lock to the store.
func (k Keeper) AddTokensToLockByID(ctx sdk.Context, lockID uint64, coin sdk.Coin) (*types.Lock, error) {
	lock, found := k.GetLockById(ctx, lockID)
	if !found {
		return nil, types.ErrInvalidLockId
	}
	owner, _ := sdk.AccAddressFromBech32(lock.Owner)
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, owner, types.ModuleName, sdk.NewCoins(coin)); err != nil {
		return nil, err
	}
	lock.Coin.Amount = lock.Coin.Amount.Add(coin.Amount)
	k.SetLock(ctx, lock)
	return &lock, nil
}

// Create New Lock
func (k Keeper) NewLockTokens(ctx sdk.Context, owner sdk.AccAddress, duration time.Duration, coin sdk.Coin) (*types.Lock, error) {
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, owner, types.ModuleName, sdk.NewCoins(coin)); err != nil {
		return nil, err
	}
	lastLockId := k.GetLockID(ctx)

	lockId := lastLockId + 1
	lock := types.NewLock(ctx, lockId, owner, duration, coin)

	k.SetLockID(ctx, lockId)
	k.SetLock(ctx, lock)

	locksByOwner, found := k.GetLockByOwner(ctx, owner.String())
	if !found {
		locksByOwner = types.LockByOwner{
			Owner:   owner.String(),
			LockIds: []uint64{},
		}
	}
	locksByOwner.LockIds = append(locksByOwner.LockIds, lock.Id)
	k.SetLockByOwner(ctx, locksByOwner)
	return &lock, nil
}

func (k Keeper) UpdateOrDeleteLock(ctx sdk.Context, isUpdate bool, lockId uint64, unlockCoin sdk.Coin) error {
	lock, found := k.GetLockById(ctx, lockId)
	if !found {
		return types.ErrInvalidLockId
	}

	// Update lock if tokens are being unlocked partially
	if isUpdate {
		if unlockCoin.Amount.GTE(lock.Coin.Amount) {
			return types.ErrInvalidUnlockingAmount
		}
		lock.Coin.Amount = lock.Coin.Amount.Sub(unlockCoin.Amount)
		k.SetLock(ctx, lock)
		return nil
	}

	// Delete lock if 100% tokens are being unlocked
	k.DeleteLock(ctx, lock.Id)
	lockByOwner, _ := k.GetLockByOwner(ctx, lock.Owner)
	updatedLockIds := []uint64{}
	for _, lId := range lockByOwner.LockIds {
		if lId != lock.Id {
			updatedLockIds = append(updatedLockIds, lId)
		}
	}
	lockByOwner.LockIds = updatedLockIds
	k.SetLockByOwner(ctx, lockByOwner)
	return nil
}

func (k Keeper) NewBeginUnlockTokens(ctx sdk.Context, owner sdk.AccAddress, lock types.Lock, unlockCoin sdk.Coin) (uint64, error) {
	lastUnlockId := k.GetUnlockingID(ctx)
	unlockId := lastUnlockId + 1

	newUnlockingTokens := types.NewUnlock(
		unlockId,
		owner,
		lock.Duration,
		ctx.BlockTime().Add(lock.Duration),
		unlockCoin,
	)

	unlockingByOwner, found := k.GetUnlockByOwner(ctx, owner.String())
	if !found {
		unlockingByOwner = types.UnlockingByOwner{
			Owner:        owner.String(),
			UnlockingIds: []uint64{},
		}
	}
	unlockingByOwner.UnlockingIds = append(unlockingByOwner.UnlockingIds, unlockId)

	k.SetUnlockingID(ctx, unlockId)
	k.SetUnlocking(ctx, newUnlockingTokens)
	k.SetUnlockByOwner(ctx, unlockingByOwner)
	return unlockId, nil
}

func (k Keeper) UpdateUnlockingByOwner(ctx sdk.Context, owner string, maturedUnlockId uint64) {
	unlocksByOwner, found := k.GetUnlockByOwner(ctx, owner)
	if found {
		updatedUnlockingIds := []uint64{}
		for _, unlockId := range unlocksByOwner.UnlockingIds {
			if unlockId != maturedUnlockId {
				updatedUnlockingIds = append(updatedUnlockingIds, unlockId)
			}
		}
		unlocksByOwner.UnlockingIds = updatedUnlockingIds
		k.SetUnlockByOwner(ctx, unlocksByOwner)
	}
}

func (k Keeper) DeleteMaturedUnlocks(ctx sdk.Context) {
	allUnlockingPositions := k.GetAllUnlockingPositions(ctx)

	for _, unlocking := range allUnlockingPositions {
		if ctx.BlockTime().After(unlocking.EndTime) {
			owner, _ := sdk.AccAddressFromBech32(unlocking.Owner)
			err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, owner, sdk.NewCoins(unlocking.Coin))
			if err == nil {
				k.DeleteUnlocking(ctx, unlocking.Id)
				k.UpdateUnlockingByOwner(ctx, unlocking.Owner, unlocking.Id)
			}
		}
	}
}
