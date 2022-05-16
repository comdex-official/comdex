package types

import (
	time "time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewLock returns a new instance of lock.
func NewLock(ctx sdk.Context, id uint64, owner sdk.AccAddress, duration time.Duration, coin sdk.Coin) Lock {
	return Lock{
		Id:        id,
		Owner:     owner.String(),
		CreatedAt: ctx.BlockTime(),
		Duration:  duration,
		Coin:      coin,
	}
}

// NewUnLock returns a new instance of unlock.
func NewUnlock(id uint64, owner sdk.AccAddress, duration time.Duration, endTime time.Time, coin sdk.Coin) Unlocking {
	return Unlocking{
		Id:       id,
		Owner:    owner.String(),
		Duration: duration,
		EndTime:  endTime,
		Coin:     coin,
	}
}
