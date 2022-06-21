package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = (*MsgAddAssetRequest)(nil)
	_ sdk.Msg = (*MsgUpdateAssetRequest)(nil)
	_ sdk.Msg = (*MsgAddPairRequest)(nil)
)

func NewMsgAddPairRequest(from sdk.AccAddress, assetIn, assetOut uint64, liquidationRatio sdk.Dec) *MsgAddPairRequest {
	return &MsgAddPairRequest{
		From:     from.String(),
		AssetIn:  assetIn,
		AssetOut: assetOut,
	}
}

func (m *MsgAddPairRequest) ValidateBasic() error {
	if m.From == "" {
		return errors.Wrap(ErrorInvalidFrom, "from cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return errors.Wrapf(ErrorInvalidFrom, "%s", err)
	}
	if m.AssetIn == 0 {
		return errors.Wrap(ErrorInvalidID, "asset_in cannot be zero")
	}
	if m.AssetOut == 0 {
		return errors.Wrap(ErrorInvalidID, "asset_out cannot be zero")
	}

	return nil
}

func (m *MsgAddPairRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgAddAssetRequest(from sdk.AccAddress, name, denom string, decimals int64) *MsgAddAssetRequest {
	return &MsgAddAssetRequest{
		From:     from.String(),
		Name:     name,
		Denom:    denom,
		Decimals: decimals,
	}
}

func (m *MsgAddAssetRequest) ValidateBasic() error {
	if m.From == "" {
		return errors.Wrap(ErrorInvalidFrom, "from cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return errors.Wrapf(ErrorInvalidFrom, "%s", err)
	}
	if m.Name == "" {
		return errors.Wrap(ErrorInvalidName, "name cannot be empty")
	}
	if len(m.Name) > MaxAssetNameLength {
		return errors.Wrapf(ErrorInvalidName, "name length cannot be greater than %d", MaxAssetNameLength)
	}
	if m.Denom == "" {
		return errors.Wrapf(ErrorInvalidDenom, "denom cannot be empty")
	}
	if err := sdk.ValidateDenom(m.Denom); err != nil {
		return errors.Wrapf(ErrorInvalidDenom, "%s", err)
	}
	if m.Decimals < 0 {
		return errors.Wrapf(ErrorInvalidDecimals, "decimals cannot be less than zero")
	}

	return nil
}

func (m *MsgAddAssetRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgUpdateAssetRequest(from sdk.AccAddress, id uint64, name, denom string, decimals int64) *MsgUpdateAssetRequest {
	return &MsgUpdateAssetRequest{
		From:     from.String(),
		Id:       id,
		Name:     name,
		Denom:    denom,
		Decimals: decimals,
	}
}

func (m *MsgUpdateAssetRequest) ValidateBasic() error {
	if m.From == "" {
		return errors.Wrap(ErrorInvalidFrom, "from cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return errors.Wrapf(ErrorInvalidFrom, "%s", err)
	}
	if m.Name != "" {
		if len(m.Name) > MaxAssetNameLength {
			return errors.Wrapf(ErrorInvalidName, "name length cannot be greater than %d", MaxAssetNameLength)
		}
	}
	if m.Denom != "" {
		if err := sdk.ValidateDenom(m.Denom); err != nil {
			return errors.Wrapf(ErrorInvalidDenom, "%s", err)
		}
	}

	return nil
}

func (m *MsgUpdateAssetRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}
