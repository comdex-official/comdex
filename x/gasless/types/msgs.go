package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = (*MsgCreateGasProvider)(nil)
	_ sdk.Msg = (*MsgAuthorizeActors)(nil)
	_ sdk.Msg = (*MsgUpdateGasProviderStatus)(nil)
	_ sdk.Msg = (*MsgUpdateGasProviderConfig)(nil)
)

// Message types for the gasless module.
const (
	TypeMsgCreateGasProvider       = "create_gas_provider"
	TypeMsgAuthorizeActors         = "authorize_actors"
	TypeMsgUpdateGasProviderStatus = "update_gas_provider_status"
	TypeMsgUpdateGasProviderConfig = "update_gas_provider_config"
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
	if msg.MaxTxsCountPerConsumer == 0 {
		return sdkerrors.Wrap(errors.ErrInvalidRequest, "max tx count per consumer must not be 0")
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

// NewMsgAuthorizeActors returns a new MsgAuthorizeActors.
func NewMsgAuthorizeActors(
	gasProviderID uint64,
	provider sdk.AccAddress,
	actors []sdk.AccAddress,
) *MsgAuthorizeActors {
	authorizedActors := []string{}
	for _, actor := range actors {
		authorizedActors = append(authorizedActors, actor.String())
	}
	return &MsgAuthorizeActors{
		GasProviderId: gasProviderID,
		Provider:      provider.String(),
		Actors:        authorizedActors,
	}
}

func (msg MsgAuthorizeActors) Route() string { return RouterKey }

func (msg MsgAuthorizeActors) Type() string { return TypeMsgAuthorizeActors }

func (msg MsgAuthorizeActors) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Provider); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid provider address: %v", err)
	}
	if msg.GasProviderId == 0 {
		return sdkerrors.Wrap(errors.ErrInvalidRequest, "gas provider id must not be 0")
	}

	for _, actor := range msg.Actors {
		if _, err := sdk.AccAddressFromBech32(actor); err != nil {
			return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid actor address - %s : %v", actor, err)
		}
	}
	return nil
}

func (msg MsgAuthorizeActors) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgAuthorizeActors) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

// NewMsgUpdateGasProviderStatus returns a new MsgUpdateGasProviderStatus.
func NewMsgUpdateGasProviderStatus(
	gasProviderID uint64,
	provider sdk.AccAddress,
) *MsgUpdateGasProviderStatus {
	return &MsgUpdateGasProviderStatus{
		GasProviderId: gasProviderID,
		Provider:      provider.String(),
	}
}

func (msg MsgUpdateGasProviderStatus) Route() string { return RouterKey }

func (msg MsgUpdateGasProviderStatus) Type() string { return TypeMsgUpdateGasProviderStatus }

func (msg MsgUpdateGasProviderStatus) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Provider); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid provider address: %v", err)
	}
	if msg.GasProviderId == 0 {
		return sdkerrors.Wrap(errors.ErrInvalidRequest, "gas provider id must not be 0")
	}
	return nil
}

func (msg MsgUpdateGasProviderStatus) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgUpdateGasProviderStatus) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

// NewMsgUpdateGasProviderConfig returns a new MsgUpdateGasProviderConfig.
func NewMsgUpdateGasProviderConfig(
	gasProviderID uint64,
	provider sdk.AccAddress,
	maxFeeUsagePerTx sdkmath.Int,
	maxTxsCountPerConsumer uint64,
	maxFeeUsagePerConsumer sdkmath.Int,
	txsAllowed []string,
	contractsAllowed []string,
) *MsgUpdateGasProviderConfig {
	return &MsgUpdateGasProviderConfig{
		GasProviderId:          gasProviderID,
		Provider:               provider.String(),
		MaxFeeUsagePerTx:       maxFeeUsagePerTx,
		MaxTxsCountPerConsumer: maxTxsCountPerConsumer,
		MaxFeeUsagePerConsumer: maxFeeUsagePerConsumer,
		TxsAllowed:             txsAllowed,
		ContractsAllowed:       contractsAllowed,
	}
}

func (msg MsgUpdateGasProviderConfig) Route() string { return RouterKey }

func (msg MsgUpdateGasProviderConfig) Type() string { return TypeMsgUpdateGasProviderConfig }

func (msg MsgUpdateGasProviderConfig) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Provider); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "invalid provider address: %v", err)
	}
	if msg.GasProviderId == 0 {
		return sdkerrors.Wrap(errors.ErrInvalidRequest, "gas provider id must not be 0")
	}
	if msg.MaxTxsCountPerConsumer == 0 {
		return sdkerrors.Wrap(errors.ErrInvalidRequest, "max tx count per consumer must not be 0")
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

func (msg MsgUpdateGasProviderConfig) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgUpdateGasProviderConfig) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}
