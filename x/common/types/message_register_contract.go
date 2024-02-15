package types

import (
	"slices"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRegisterContract = "register_contract"

var _ sdk.Msg = &MsgRegisterContract{}

func NewMsgRegisterContract(
	securityAddress string,
	gameName string,
	contractAddress string,
	gameType uint64,
) *MsgRegisterContract {
	return &MsgRegisterContract{
		SecurityAddress: securityAddress,
		GameName: gameName,
		ContractAddress: contractAddress,
		GameType: gameType,
	}
}

func (msg *MsgRegisterContract) Route() string {
	return RouterKey
}

func (msg *MsgRegisterContract) Type() string {
	return TypeMsgRegisterContract
}

func (msg *MsgRegisterContract) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.SecurityAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRegisterContract) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRegisterContract) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.SecurityAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.ContractAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid contract address (%s)", err)
	}
	gameType := []uint64{1, 2, 3}
	if !slices.Contains(gameType, msg.GameType) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidType, "invalid game type, should be 1,2 or 3")
	}

	return nil
}
