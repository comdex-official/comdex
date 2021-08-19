package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	ibcclienttypes "github.com/cosmos/ibc-go/modules/core/02-client/types"
	ibchost "github.com/cosmos/ibc-go/modules/core/24-host"
)

var (
	_ sdk.Msg = (*MsgFetchPriceRequest)(nil)
)

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
