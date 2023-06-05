package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = (*MsgAddAsset)(nil)

// Message types for the liquidity module.
const (
	TypeMsgAddAsset = "add_asset"
)

// NewMsgAddAsset returns a new MsgAddAsset.
func NewMsgAddAsset(
	creator sdk.AccAddress,
	name string,
	denom string,
	decimals uint64,
	isOnChain bool,
	isOraclePriceRequired bool,
	isCdpMintable bool,
) *MsgAddAsset {
	return &MsgAddAsset{
		Creator: creator.String(),
		Asset: Asset{
			Name:                  name,
			Denom:                 denom,
			Decimals:              sdk.NewIntFromUint64(decimals),
			IsOnChain:             isOnChain,
			IsOraclePriceRequired: isOraclePriceRequired,
			IsCdpMintable:         isCdpMintable,
		},
	}
}

func (msg MsgAddAsset) Route() string { return RouterKey }

func (msg MsgAddAsset) Type() string { return TypeMsgAddAsset }

func (msg MsgAddAsset) ValidateBasic() error {
	if err := msg.Asset.Validate(); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address: %v", err)
	}
	return nil
}

func (msg MsgAddAsset) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgAddAsset) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgAddAsset) GetCreator() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return addr
}
