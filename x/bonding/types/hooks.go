package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BondingHooks interface {
	AfterAddTokensToLock(ctx sdk.Context, address sdk.AccAddress, lockID uint64, amount sdk.Coins)
	OnTokenLocked(ctx sdk.Context, address sdk.AccAddress, lockID uint64, amount sdk.Coins, lockDuration time.Duration, unlockTime time.Time)
	OnStartUnlock(ctx sdk.Context, address sdk.AccAddress, lockID uint64, amount sdk.Coins, lockDuration time.Duration, unlockTime time.Time)
	OnTokenUnlocked(ctx sdk.Context, address sdk.AccAddress, lockID uint64, amount sdk.Coins, lockDuration time.Duration, unlockTime time.Time)
	OnTokenSlashed(ctx sdk.Context, lockID uint64, amount sdk.Coins)
}

var _ BondingHooks = MultiBondingHooks{}

// combine multiple gamm hooks, all hook functions are run in array sequence.
type MultiBondingHooks []BondingHooks

func NewMultiBondingHooks(hooks ...BondingHooks) MultiBondingHooks {
	return hooks
}

func (h MultiBondingHooks) AfterAddTokensToLock(ctx sdk.Context, address sdk.AccAddress, lockID uint64, amount sdk.Coins) {
	for i := range h {
		h[i].AfterAddTokensToLock(ctx, address, lockID, amount)
	}
}

func (h MultiBondingHooks) OnTokenLocked(ctx sdk.Context, address sdk.AccAddress, lockID uint64, amount sdk.Coins, lockDuration time.Duration, unlockTime time.Time) {
	for i := range h {
		h[i].OnTokenLocked(ctx, address, lockID, amount, lockDuration, unlockTime)
	}
}

func (h MultiBondingHooks) OnStartUnlock(ctx sdk.Context, address sdk.AccAddress, lockID uint64, amount sdk.Coins, lockDuration time.Duration, unlockTime time.Time) {
	for i := range h {
		h[i].OnStartUnlock(ctx, address, lockID, amount, lockDuration, unlockTime)
	}
}

func (h MultiBondingHooks) OnTokenUnlocked(ctx sdk.Context, address sdk.AccAddress, lockID uint64, amount sdk.Coins, lockDuration time.Duration, unlockTime time.Time) {
	for i := range h {
		h[i].OnTokenUnlocked(ctx, address, lockID, amount, lockDuration, unlockTime)
	}
}

func (h MultiBondingHooks) OnTokenSlashed(ctx sdk.Context, lockID uint64, amount sdk.Coins) {
	for i := range h {
		h[i].OnTokenSlashed(ctx, lockID, amount)
	}
}
