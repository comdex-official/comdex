package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRegisterContract = "register_contract"

var _ sdk.Msg = &MsgRegisterContract{}

func NewMsgRegisterContract(
	creator string,
	gameName string,
	gameId uint64,
	contractAddress string,
) *MsgRegisterContract {
	return &MsgRegisterContract{
		Creator: creator,
		GameName: gameName,
		GameId: gameId,
		ContractAddr: contractAddress,
	}
}

func (msg *MsgRegisterContract) Route() string {
	return RouterKey
}

func (msg *MsgRegisterContract) Type() string {
	return TypeMsgRegisterContract
}

func (msg *MsgRegisterContract) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
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
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.ContractAddr)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid contract address (%s)", err)
	}

	return nil
}
