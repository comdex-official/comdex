package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = (*MsgKillRequest)(nil)
)

const (
	DefaultCircuitBreakerEnabled = false
)

func NewMsgKillRequest(from sdk.AccAddress, condition bool) *MsgKillRequest {
	return &MsgKillRequest{
		From:          from.String(),
		BreakerEnable: condition,
	}

}

func (m *MsgKillRequest) ValidateBasic() error {

	if m.From == "" {
		return errors.Wrap(errors.ErrInvalidAddress, "from cannot be empty")
	}

	return nil
}

func (m *MsgKillRequest) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}
