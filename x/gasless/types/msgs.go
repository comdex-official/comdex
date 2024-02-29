package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = (*MsgCreateGasProvider)(nil)
)

// Message types for the gasless module.
const (
	TypeMsgCreateGasProvider = "create_gas_provider"
)

// NewMsgCreateGasProvider returns a new MsgCreateGasProvider.
func NewMsgCreateGasProvider(
	creator sdk.AccAddress,
	feeDenom string,
	maxFeeUsagePerTx sdkmath.Int,
	maxTxsCountPerConsumer uint64,
	maxFeeUsagePerConsumer sdkmath.Int,
	txsAllowed []string,
	contractsAllowed []string,
	gasDeposit sdk.Coin,
) *MsgCreateGasProvider {
	return &MsgCreateGasProvider{
		Creator:                creator.String(),
		FeeDenom:               feeDenom,
		MaxFeeUsagePerTx:       maxFeeUsagePerTx,
		MaxTxsCountPerConsumer: maxTxsCountPerConsumer,
		MaxFeeUsagePerConsumer: maxFeeUsagePerConsumer,
		TxsAllowed:             txsAllowed,
		ContractsAllowed:       contractsAllowed,
		GasDeposit:             gasDeposit,
	}
}

func (msg MsgCreateGasProvider) Route() string { return RouterKey }

func (msg MsgCreateGasProvider) Type() string { return TypeMsgCreateGasProvider }

func (msg MsgCreateGasProvider) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid creator address: %v", err)
	}
	if err := sdk.ValidateDenom(msg.FeeDenom); err != nil {
		return sdkerrors.Wrap(errors.ErrInvalidRequest, err.Error())
	}
	if msg.FeeDenom != msg.GasDeposit.Denom {
		return sdkerrors.Wrap(errors.ErrInvalidRequest, "denom mismatch, fee denom and gas_deposit")
	}
	if !msg.MaxFeeUsagePerTx.IsPositive() {
		return sdkerrors.Wrap(errors.ErrInvalidRequest, "max_fee_usage_per_tx should be positive")
	}
	if !msg.MaxFeeUsagePerConsumer.IsPositive() {
		return sdkerrors.Wrap(errors.ErrInvalidRequest, "max_fee_usage_per_consumer should be positive")
	}
	if len(msg.TxsAllowed) == 0 && len(msg.ContractsAllowed) == 0 {
		return sdkerrors.Wrap(errors.ErrInvalidRequest, "atleast one tx or contract is required to initialize")
	}
	return nil
}

func (msg MsgCreateGasProvider) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgCreateGasProvider) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}
