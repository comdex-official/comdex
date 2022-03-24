package types

import (
	"errors"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = (*MsgDepositMintingRewardAmountRequest)(nil)
	_ sdk.Msg = (*MsgUpdateMintRewardStartTimeRequest)(nil)
)

func NewMsgDepositMintingRewardAmount(mintingrewardId uint64, from sdk.AccAddress, startTimeStamp time.Time) *MsgDepositMintingRewardAmountRequest {
	return &MsgDepositMintingRewardAmountRequest{
		MintingRewardId: mintingrewardId,
		StartTimestamp:  startTimeStamp,
		From:            from.String(),
	}
}

func NewMsgUpdateMintRewardStartTime(mintingrewardId uint64, from sdk.AccAddress, newStartTimeStamp time.Time) *MsgUpdateMintRewardStartTimeRequest {
	return &MsgUpdateMintRewardStartTimeRequest{
		MintingRewardId:   mintingrewardId,
		NewStartTimestamp: newStartTimeStamp,
		From:              from.String(),
	}
}

func (m *MsgDepositMintingRewardAmountRequest) ValidateBasic() error {
	if m.MintingRewardId == 0 {
		return errors.New("invalid minting rewards id")
	}
	_, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "--from address cannot be empty or invalid")
	}
	return nil
}

func (m *MsgDepositMintingRewardAmountRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func (m *MsgUpdateMintRewardStartTimeRequest) ValidateBasic() error {
	if m.MintingRewardId == 0 {
		return errors.New("invalid minting rewards id")
	}
	_, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "--from address cannot be empty or invalid")
	}
	return nil
}

func (m *MsgUpdateMintRewardStartTimeRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}
