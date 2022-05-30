package keeper

import (
	"time"

	"github.com/comdex-official/comdex/x/locking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetAccountWithSameLockDurationAndDenom(
	ctx sdk.Context,
	//nolint
	owner sdk.AccAddress,
	denom string,
	duration time.Duration,
) []types.Lock {
	locks := []types.Lock{}

	locksByOwner, found := k.GetLockByOwner(ctx, owner.String())
	if !found {
		return locks
	}

	for _, lockID := range locksByOwner.LockIds {
		lock, found := k.GetLockByID(ctx, lockID)
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
	lock, found := k.GetLockByID(ctx, lockID)
	if !found {
		return nil, types.ErrInvalidLockID
	}
	owner, _ := sdk.AccAddressFromBech32(lock.Owner)
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, owner, types.ModuleName, sdk.NewCoins(coin)); err != nil {
		return nil, err
	}
	lock.Coin.Amount = lock.Coin.Amount.Add(coin.Amount)
	k.SetLock(ctx, lock)
	return &lock, nil
}

// Create New Lock.
func (k Keeper) NewLockTokens(
	ctx sdk.Context,
	//nolint
	owner sdk.AccAddress,
	duration time.Duration,
	coin sdk.Coin,
) (*types.Lock, error) {
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, owner, types.ModuleName, sdk.NewCoins(coin)); err != nil {
		return nil, err
	}
	lastLockID := k.GetLockID(ctx)

	lockID := lastLockID + 1
	lock := types.NewLock(ctx, lockID, owner, duration, coin)

	k.SetLockID(ctx, lockID)
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

func (k Keeper) UpdateOrDeleteLock(ctx sdk.Context, isUpdate bool, lockID uint64, unlockCoin sdk.Coin) error {
	lock, found := k.GetLockByID(ctx, lockID)
	if !found {
		return types.ErrInvalidLockID
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
	for _, lID := range lockByOwner.LockIds {
		if lID != lock.Id {
			updatedLockIds = append(updatedLockIds, lID)
		}
	}
	lockByOwner.LockIds = updatedLockIds
	k.SetLockByOwner(ctx, lockByOwner)
	return nil
}

func (k Keeper) NewBeginUnlockTokens(
	ctx sdk.Context,
	//nolint
	owner sdk.AccAddress,
	lock types.Lock,
	unlockCoin sdk.Coin,
) (uint64, error) {
	lastUnlockID := k.GetUnlockingID(ctx)
	unlockID := lastUnlockID + 1

	newUnlockingTokens := types.NewUnlock(
		unlockID,
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
	unlockingByOwner.UnlockingIds = append(unlockingByOwner.UnlockingIds, unlockID)

	k.SetUnlockingID(ctx, unlockID)
	k.SetUnlocking(ctx, newUnlockingTokens)
	k.SetUnlockByOwner(ctx, unlockingByOwner)
	return unlockID, nil
}

func (k Keeper) UpdateUnlockingByOwner(ctx sdk.Context, owner string, maturedUnlockID uint64) {
	unlocksByOwner, found := k.GetUnlockByOwner(ctx, owner)
	if found {
		updatedUnlockingIds := []uint64{}
		for _, unlockID := range unlocksByOwner.UnlockingIds {
			if unlockID != maturedUnlockID {
				updatedUnlockingIds = append(updatedUnlockingIds, unlockID)
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
