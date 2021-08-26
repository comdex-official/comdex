package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	ibcclienttypes "github.com/cosmos/ibc-go/modules/core/02-client/types"
	ibchost "github.com/cosmos/ibc-go/modules/core/24-host"
)

var (
	_ sdk.Msg = (*MsgAddMarketRequest)(nil)
	_ sdk.Msg = (*MsgUpdateMarketRequest)(nil)
	_ sdk.Msg = (*MsgAddAssetRequest)(nil)
	_ sdk.Msg = (*MsgUpdateAssetRequest)(nil)
	_ sdk.Msg = (*MsgAddMarketForAssetRequest)(nil)
	_ sdk.Msg = (*MsgRemoveMarketForAssetRequest)(nil)
	_ sdk.Msg = (*MsgAddPairRequest)(nil)
	_ sdk.Msg = (*MsgUpdatePairRequest)(nil)
	_ sdk.Msg = (*MsgFetchPriceRequest)(nil)
)

func NewMsgAddMarketRequest(from sdk.AccAddress, symbol string, scriptID uint64) *MsgAddMarketRequest {
	return &MsgAddMarketRequest{
		From:     from.String(),
		Symbol:   symbol,
		ScriptID: scriptID,
	}
}

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

func NewMsgUpdateMarketRequest(from sdk.AccAddress, symbol string, scriptID uint64) *MsgUpdateMarketRequest {
	return &MsgUpdateMarketRequest{
		From:     from.String(),
		Symbol:   symbol,
		ScriptID: scriptID,
	}
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

func NewMsgAddAssetRequest(from sdk.AccAddress, name, denom string) *MsgAddAssetRequest {
	return &MsgAddAssetRequest{
		From:  from.String(),
		Name:  name,
		Denom: denom,
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

	return nil
}

func (m *MsgAddAssetRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgUpdateAssetRequest(from sdk.AccAddress, id uint64, name, denom string) *MsgUpdateAssetRequest {
	return &MsgUpdateAssetRequest{
		From:  "",
		ID:    0,
		Name:  "",
		Denom: "",
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

func NewMsgAddMarketForAssetRequest(from sdk.AccAddress, id uint64, symbol string) *MsgAddMarketForAssetRequest {
	return &MsgAddMarketForAssetRequest{
		From:   from.String(),
		ID:     id,
		Symbol: symbol,
	}
}

func (m *MsgAddMarketForAssetRequest) ValidateBasic() error {
	if m.From == "" {
		return errors.Wrap(ErrorInvalidFrom, "from cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return errors.Wrapf(ErrorInvalidFrom, "%s", err)
	}
	if m.ID == 0 {
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

func (m *MsgAddMarketForAssetRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgRemoveMarketForAssetRequest(from sdk.AccAddress, id uint64, symbol string) *MsgRemoveMarketForAssetRequest {
	return &MsgRemoveMarketForAssetRequest{
		From:   from.String(),
		ID:     id,
		Symbol: symbol,
	}
}

func (m *MsgRemoveMarketForAssetRequest) ValidateBasic() error {
	if m.From == "" {
		return errors.Wrap(ErrorInvalidFrom, "from cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return errors.Wrapf(ErrorInvalidFrom, "%s", err)
	}
	if m.ID == 0 {
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

func NewMsgAddPairRequest(from sdk.AccAddress, assetIn, assetOut uint64, liquidationRatio sdk.Dec) *MsgAddPairRequest {
	return &MsgAddPairRequest{
		From:             from.String(),
		AssetIn:          assetIn,
		AssetOut:         assetOut,
		LiquidationRatio: liquidationRatio,
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
	if m.LiquidationRatio.IsNil() {
		return errors.Wrap(ErrorInvalidLiquidationRatio, "liquidation_ratio cannot be nil")
	}
	if m.LiquidationRatio.IsNegative() {
		return errors.Wrap(ErrorInvalidLiquidationRatio, "liquidation_ratio cannot be negative")
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

func NewMsgUpdatePairRequest(from sdk.AccAddress, id uint64, assetIn, assetOut uint64, liquidationRatio sdk.Dec) *MsgUpdatePairRequest {
	return &MsgUpdatePairRequest{
		From:             from.String(),
		ID:               id,
		AssetIn:          assetIn,
		AssetOut:         assetOut,
		LiquidationRatio: liquidationRatio,
	}
}

func (m *MsgUpdatePairRequest) ValidateBasic() error {
	if m.From == "" {
		return errors.Wrap(ErrorInvalidFrom, "from cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return errors.Wrapf(ErrorInvalidFrom, "%s", err)
	}
	if !m.LiquidationRatio.IsNil() {
		if m.LiquidationRatio.IsNegative() {
			return errors.Wrap(ErrorInvalidLiquidationRatio, "liquidation_ratio cannot be negative")
		}
	}

	return nil
}

func (m *MsgUpdatePairRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgFetchPriceRequest(
	from sdk.AccAddress,
	sourcePort, sourceChannel string,
	timeoutHeight ibcclienttypes.Height,
	timeoutTimestamp uint64,
	symbols []string,
	scriptID uint64,
	feeLimit sdk.Coins,
	prepareGas, executeGas uint64,
) *MsgFetchPriceRequest {
	return &MsgFetchPriceRequest{
		From:             from.String(),
		SourcePort:       sourcePort,
		SourceChannel:    sourceChannel,
		TimeoutHeight:    timeoutHeight,
		TimeoutTimestamp: timeoutTimestamp,
		Symbols:          symbols,
		ScriptID:         scriptID,
		FeeLimit:         feeLimit,
		PrepareGas:       prepareGas,
		ExecuteGas:       executeGas,
	}
}

func (m *MsgFetchPriceRequest) ValidateBasic() error {
	if m.From == "" {
		return errors.Wrap(ErrorInvalidFrom, "from cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.From); err != nil {
		return errors.Wrapf(ErrorInvalidFrom, "%s", err)
	}
	if err := ibchost.PortIdentifierValidator(m.SourcePort); err != nil {
		return errors.Wrapf(ErrorInvalidSourcePort, "%s", err)
	}
	if err := ibchost.ChannelIdentifierValidator(m.SourceChannel); err != nil {
		return errors.Wrapf(ErrorInvalidSourceChannel, "%s", err)
	}
	if m.Symbols == nil {
		return errors.Wrapf(ErrorInvalidSymbols, "symbols cannot be nil")
	}
	if len(m.Symbols) == 0 {
		return errors.Wrapf(ErrorInvalidSymbols, "symbols cannot be empty")
	}
	if m.ScriptID == 0 {
		return errors.Wrapf(ErrorInvalidScriptID, "script_id cannot be zero")
	}

	return nil
}

func (m *MsgFetchPriceRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}
