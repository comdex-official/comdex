package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDeRegisterContract = "de_register_contract"

var _ sdk.Msg = &MsgDeRegisterContract{}

func NewMsgDeRegisterContract(
	securityAddress string,
	gameID uint64,
) *MsgDeRegisterContract {
	return &MsgDeRegisterContract{
		SecurityAddress: securityAddress,
		GameId: gameID,
	}
}

func (msg *MsgDeRegisterContract) Route() string {
	return RouterKey
}

func (msg *MsgDeRegisterContract) Type() string {
	return TypeMsgDeRegisterContract
}

func (msg *MsgDeRegisterContract) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.SecurityAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeRegisterContract) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeRegisterContract) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.SecurityAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}
