package keeper

import (
	"fmt"
	"time"

	"github.com/comdex-official/comdex/x/bonding/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
)

func (k Keeper) setSyntheticBondingObject(ctx sdk.Context, synthLock *types.SyntheticLock) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := proto.Marshal(synthLock)
	if err != nil {
		return err
	}
	store.Set(syntheticLockStoreKey(synthLock.UnderlyingLockId, synthLock.SynthDenom), bz)
	if !synthLock.EndTime.Equal(time.Time{}) {
		store.Set(syntheticLockTimeStoreKey(synthLock.UnderlyingLockId, synthLock.SynthDenom, synthLock.EndTime), bz)
	}
	return nil
}

func (k Keeper) deleteSyntheticBondingObject(ctx sdk.Context, lockID uint64, synthdenom string) {
	store := ctx.KVStore(k.storeKey)
	synthLock, _ := k.GetSyntheticBonding(ctx, lockID, synthdenom)
	if synthLock != nil && !synthLock.EndTime.Equal(time.Time{}) {
		store.Delete(syntheticLockTimeStoreKey(lockID, synthdenom, synthLock.EndTime))
	}
	store.Delete(syntheticLockStoreKey(lockID, synthdenom))
}

func (k Keeper) GetUnderlyingLock(ctx sdk.Context, synthlock types.SyntheticLock) types.PeriodLock {
	lock, err := k.GetLockByID(ctx, synthlock.UnderlyingLockId)
	if err != nil {
		panic(err) // Synthetic lock MUST have underlying lock
	}
	return *lock
}

func (k Keeper) GetSyntheticBonding(ctx sdk.Context, lockID uint64, synthdenom string) (*types.SyntheticLock, error) {
	synthLock := types.SyntheticLock{}
	store := ctx.KVStore(k.storeKey)
	synthLockKey := syntheticLockStoreKey(lockID, synthdenom)
	if !store.Has(synthLockKey) {
		return nil, fmt.Errorf("synthetic lock with ID %d and synth denom %s does not exist", lockID, synthdenom)
	}
	bz := store.Get(synthLockKey)
	err := proto.Unmarshal(bz, &synthLock)
	return &synthLock, err
}

func (k Keeper) GetAllSyntheticBondingsByBonding(ctx sdk.Context, lockID uint64) []types.SyntheticLock {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, combineKeys(types.KeyPrefixSyntheticBonding, sdk.Uint64ToBigEndian(lockID)))
	defer iterator.Close()

	synthLocks := []types.SyntheticLock{}
	for ; iterator.Valid(); iterator.Next() {
		synthLock := types.SyntheticLock{}
		err := proto.Unmarshal(iterator.Value(), &synthLock)
		if err != nil {
			panic(err)
		}
		synthLocks = append(synthLocks, synthLock)
	}
	return synthLocks
}

func (k Keeper) GetAllSyntheticBondingsByAddr(ctx sdk.Context, owner sdk.AccAddress) []types.SyntheticLock {
	synthLocks := []types.SyntheticLock{}
	locks := k.GetAccountPeriodLocks(ctx, owner)
	for _, lock := range locks {
		synthLocks = append(synthLocks, k.GetAllSyntheticBondingsByBonding(ctx, lock.ID)...)
	}
	return synthLocks
}

func (k Keeper) HasAnySyntheticBondings(ctx sdk.Context, lockID uint64) bool {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, combineKeys(types.KeyPrefixSyntheticBonding, sdk.Uint64ToBigEndian(lockID)))
	defer iterator.Close()
	return iterator.Valid()
}

func (k Keeper) GetAllSyntheticBondings(ctx sdk.Context) []types.SyntheticLock {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixSyntheticBonding)
	defer iterator.Close()

	synthLocks := []types.SyntheticLock{}
	for ; iterator.Valid(); iterator.Next() {
		synthLock := types.SyntheticLock{}
		err := proto.Unmarshal(iterator.Value(), &synthLock)
		if err != nil {
			panic(err)
		}
		synthLocks = append(synthLocks, synthLock)
	}
	return synthLocks
}

// CreateSyntheticBonding create synthetic bonding with lock id and synthdenom.
func (k Keeper) CreateSyntheticBonding(ctx sdk.Context, lockID uint64, synthDenom string, unlockDuration time.Duration, isUnlocking bool) error {
	// Note: synthetic bonding is doing everything same as bonding except coin movement
	// There is no relationship between unbonding and bonding synthetic bonding, it's managed separately
	// Accumulation store works without caring about unlocking synthetic or not

	_, err := k.GetSyntheticBonding(ctx, lockID, synthDenom)
	if err == nil {
		return types.ErrSyntheticBondingAlreadyExists
	}

	lock, err := k.GetLockByID(ctx, lockID)
	if err != nil {
		return err
	}

	endTime := time.Time{}
	if isUnlocking { // end time is set automatically if it's unlocking bonding
		if unlockDuration > lock.Duration {
			return types.ErrSyntheticDurationLongerThanNative
		}
		endTime = ctx.BlockTime().Add(unlockDuration)
	}

	// set synthetic bonding object
	synthLock := types.SyntheticLock{
		UnderlyingLockId: lockID,
		SynthDenom:       synthDenom,
		EndTime:          endTime,
		Duration:         unlockDuration,
	}
	err = k.setSyntheticBondingObject(ctx, &synthLock)
	if err != nil {
		return err
	}

	// add lock refs into not unlocking queue
	err = k.addSyntheticLockRefs(ctx, *lock, synthLock)
	if err != nil {
		return err
	}

	coin, err := lock.SingleCoin()
	if err != nil {
		return err
	}

	k.accumulationStore(ctx, synthLock.SynthDenom).Increase(accumulationKey(unlockDuration), coin.Amount)
	return nil
}

// DeleteSyntheticBonding delete synthetic bonding with lock id and synthdenom.
func (k Keeper) DeleteSyntheticBonding(ctx sdk.Context, lockID uint64, synthdenom string) error {
	synthLock, err := k.GetSyntheticBonding(ctx, lockID, synthdenom)
	if err != nil {
		return err
	}

	lock, err := k.GetLockByID(ctx, lockID)
	if err != nil {
		return err
	}

	// update lock for synthetic lock
	lock.EndTime = synthLock.EndTime

	k.deleteSyntheticBondingObject(ctx, lockID, synthdenom)

	// delete lock refs from the unlocking queue
	err = k.deleteSyntheticLockRefs(ctx, *lock, *synthLock)
	if err != nil {
		return err
	}

	// remove from accumulation store
	coin, err := lock.SingleCoin()
	if err != nil {
		return err
	}
	k.accumulationStore(ctx, synthLock.SynthDenom).Decrease(accumulationKey(lock.Duration), coin.Amount)
	return nil
}

func (k Keeper) DeleteAllMaturedSyntheticLocks(ctx sdk.Context) {
	iterator := k.iteratorBeforeTime(ctx, combineKeys(types.KeyPrefixSyntheticLockTimestamp), ctx.BlockTime())
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		synthLock := types.SyntheticLock{}
		err := proto.Unmarshal(iterator.Value(), &synthLock)
		if err != nil {
			panic(err)
		}
		err = k.DeleteSyntheticBonding(ctx, synthLock.UnderlyingLockId, synthLock.SynthDenom)
		if err != nil {
			// TODO: When underlying lock is deleted for a reason while synthetic bonding exists, panic could happen
			panic(err)
		}
	}
}
