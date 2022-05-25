package keeper

import (
	"github.com/comdex-official/comdex/x/locking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

// LOCKED TOKENS

func (k *Keeper) GetLockID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.LockIDKey
		value = store.Get(key)
	)
	if value == nil {
		return 0
	}
	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)
	return id.GetValue()
}

func (k *Keeper) SetLockID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LockIDKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)
	store.Set(key, value)
}

func (k *Keeper) SetLock(ctx sdk.Context, lock types.Lock) {
	var (
		store = k.Store(ctx)
		key   = types.GetLockKey(lock.Id)
		value = k.cdc.MustMarshal(&lock)
	)
	store.Set(key, value)
}

func (k *Keeper) DeleteLock(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.GetLockKey(id)
	)
	store.Delete(key)
}

func (k *Keeper) GetLockByID(ctx sdk.Context, id uint64) (lock types.Lock, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.GetLockKey(id)
		value = store.Get(key)
	)
	if value == nil {
		return lock, false
	}
	k.cdc.MustUnmarshal(value, &lock)
	return lock, true
}

func (k *Keeper) SetLockByOwner(ctx sdk.Context, lockByOwner types.LockByOwner) {
	var (
		store = k.Store(ctx)
		key   = types.GetLockByOwnerKey(lockByOwner.Owner)
		value = k.cdc.MustMarshal(&lockByOwner)
	)
	store.Set(key, value)
}

func (k *Keeper) GetLockByOwner(ctx sdk.Context, owner string) (lockByOwner types.LockByOwner, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.GetLockByOwnerKey(owner)
		value = store.Get(key)
	)
	if value == nil {
		return lockByOwner, false
	}
	k.cdc.MustUnmarshal(value, &lockByOwner)
	return lockByOwner, true
}

// TOKENS IN UNLOCKING PERIOD

func (k *Keeper) GetUnlockingID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.UnlockingIDKey
		value = store.Get(key)
	)
	if value == nil {
		return 0
	}
	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)
	return id.GetValue()
}

func (k *Keeper) SetUnlockingID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.UnlockingIDKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)
	store.Set(key, value)
}

func (k *Keeper) SetUnlocking(ctx sdk.Context, unlock types.Unlocking) {
	var (
		store = k.Store(ctx)
		key   = types.GetUnlockKey(unlock.Id)
		value = k.cdc.MustMarshal(&unlock)
	)
	store.Set(key, value)
}

func (k *Keeper) DeleteUnlocking(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.GetUnlockKey(id)
	)
	store.Delete(key)
}

func (k *Keeper) GetUnlockingByID(ctx sdk.Context, id uint64) (unlock types.Unlocking, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.GetUnlockKey(id)
		value = store.Get(key)
	)
	if value == nil {
		return unlock, false
	}
	k.cdc.MustUnmarshal(value, &unlock)
	return unlock, true
}

func (k *Keeper) SetUnlockByOwner(ctx sdk.Context, unlockingByOwner types.UnlockingByOwner) {
	var (
		store = k.Store(ctx)
		key   = types.GetUnlockByOwnerKey(unlockingByOwner.Owner)
		value = k.cdc.MustMarshal(&unlockingByOwner)
	)
	store.Set(key, value)
}

func (k *Keeper) GetUnlockByOwner(ctx sdk.Context, owner string) (unlockingByOwner types.UnlockingByOwner, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.GetUnlockByOwnerKey(owner)
		value = store.Get(key)
	)
	if value == nil {
		return unlockingByOwner, false
	}
	k.cdc.MustUnmarshal(value, &unlockingByOwner)
	return unlockingByOwner, true
}

func (k *Keeper) GetAllUnlockingPositions(ctx sdk.Context) (unlockings []types.Unlocking) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.UnlockKeyPrefix)
	)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var unlocking types.Unlocking
		k.cdc.MustUnmarshal(iter.Value(), &unlocking)
		unlockings = append(unlockings, unlocking)
	}
	return unlockings
}
