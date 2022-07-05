package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = (*MsgAddMarketRequest)(nil)
	_ sdk.Msg = (*MsgUpdateMarketRequest)(nil)
	_ sdk.Msg = (*MsgRemoveMarketForAssetRequest)(nil)
)

func (m *MsgAddMarketRequest) ValidateBasic() error {
	if m.From == "" {
		return errors.Wrap(ErrorInvalidFrom, "from cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return errors.Wrapf(ErrorInvalidFrom, "%s", err)
	}
	if m.Symbol == "" {
		return errors.Wrap(ErrorInvalidSymbol, "symbol cannot be empty")
	}
	if len(m.Symbol) > MaxMarketSymbolLength {
		return errors.Wrapf(ErrorInvalidSymbol, "symbol length cannot be greater than %d", MaxMarketSymbolLength)
	}
	if m.ScriptID == 0 {
		return errors.Wrapf(ErrorInvalidScriptID, "script_id cannot be zero")
	}

	return nil
}

func (m *MsgAddMarketRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func (m *MsgUpdateMarketRequest) ValidateBasic() error {
	if m.From == "" {
		return errors.Wrap(ErrorInvalidFrom, "from cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return errors.Wrapf(ErrorInvalidFrom, "%s", err)
	}
	if m.Symbol != "" {
		if len(m.Symbol) > MaxMarketSymbolLength {
			return errors.Wrapf(ErrorInvalidSymbol, "symbol length cannot be greater than %d", MaxMarketSymbolLength)
		}
	}

	return nil
}

func (m *MsgUpdateMarketRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func (m *MsgRemoveMarketForAssetRequest) ValidateBasic() error {
	if m.From == "" {
		return errors.Wrap(ErrorInvalidFrom, "from cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return errors.Wrapf(ErrorInvalidFrom, "%s", err)
	}
	if m.Id == 0 {
		return errors.Wrap(ErrorInvalidID, "id cannot be zero")
	}
	if m.Symbol == "" {
		return errors.Wrap(ErrorInvalidSymbol, "symbol cannot be empty")
	}
	if len(m.Symbol) > MaxMarketSymbolLength {
		return errors.Wrapf(ErrorInvalidSymbol, "symbol length cannot be greater than %d", MaxMarketSymbolLength)
	}

	return nil
}

func (m *MsgRemoveMarketForAssetRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}
