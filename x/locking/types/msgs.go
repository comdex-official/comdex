package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// constants.
const (
	TypeMsgLockTokens  = "lock_tokens"
	TypeMsgBeginUnlock = "begin_unlock"
)

var _ sdk.Msg = &MsgLockTokens{}

// NewMsgLockTokens creates a message to lock tokens.
func NewMsgLockTokens(
	//nolint
	owner sdk.AccAddress,
	duration time.Duration,
	coin sdk.Coin,
) *MsgLockTokens {
	return &MsgLockTokens{
		Owner:    owner.String(),
		Duration: duration,
		Coin:     coin,
	}
}

func (m MsgLockTokens) Route() string { return RouterKey }
func (m MsgLockTokens) Type() string  { return TypeMsgLockTokens }
func (m MsgLockTokens) ValidateBasic() error {
	if m.Duration <= 0 {
		return fmt.Errorf("duration should be positive: %d < 0", m.Duration)
	}
	if m.Coin.Amount.IsNegative() || m.Coin.Amount.IsZero() {
		return fmt.Errorf("invalid coin amount: %d < 0", m.Coin.Amount)
	}
	return nil
}

func (m MsgLockTokens) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgLockTokens) GetSigners() []sdk.AccAddress {
	owner, _ := sdk.AccAddressFromBech32(m.Owner)
	return []sdk.AccAddress{owner}
}

var _ sdk.Msg = &MsgBeginUnlockingTokens{}

// NewMsgLockTokens creates a message to lock tokens.
func NewMsgBeginUnlockingTokens(
	//nolint
	owner sdk.AccAddress,
	lockID uint64,
	coin sdk.Coin,
) *MsgBeginUnlockingTokens {
	return &MsgBeginUnlockingTokens{
		Owner:  owner.String(),
		LockId: lockID,
		Coin:   coin,
	}
}

func (m MsgBeginUnlockingTokens) Route() string { return RouterKey }
func (m MsgBeginUnlockingTokens) Type() string  { return TypeMsgLockTokens }
func (m MsgBeginUnlockingTokens) ValidateBasic() error {
	if m.LockId <= 0 {
		return fmt.Errorf("invalid lock_id: %d < 0", m.LockId)
	}
	if m.Coin.Amount.IsNegative() || m.Coin.Amount.IsZero() {
		return fmt.Errorf("invalid coin amount: %d < 0", m.Coin.Amount)
	}
	return nil
}

func (m MsgBeginUnlockingTokens) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgBeginUnlockingTokens) GetSigners() []sdk.AccAddress {
	owner, _ := sdk.AccAddressFromBech32(m.Owner)
	return []sdk.AccAddress{owner}
}
