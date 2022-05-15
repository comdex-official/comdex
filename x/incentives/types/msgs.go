package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// constants.
const (
	TypeMsgCreateGauge = "create_gauge"
)

var _ sdk.Msg = &MsgCreateGauge{}

// NewMsgCreateGauge creates a message to add a new gauge.
func NewMsgCreateGauge(
	from sdk.AccAddress,
	gaugeTypeId uint64,
	triggerDuration time.Duration,
	depositAmount sdk.Coin,
	totalTriggers uint64,
) *MsgCreateGauge {
	return &MsgCreateGauge{
		From:            from.String(),
		GaugeTypeId:     gaugeTypeId,
		TriggerDuration: triggerDuration,
		DepositAmount:   depositAmount,
		TotalTriggers:   totalTriggers,
		Kind:            nil,
	}
}

func (m MsgCreateGauge) Route() string { return RouterKey }
func (m MsgCreateGauge) Type() string  { return TypeMsgCreateGauge }
func (m MsgCreateGauge) ValidateBasic() error {

	isValidGaugeTypeId := false
	for _, gaugeTypeId := range ValidGaugeTypeIds {
		if gaugeTypeId == m.GaugeTypeId {
			isValidGaugeTypeId = true
			break
		}
	}
	if !isValidGaugeTypeId {
		err := fmt.Sprintf("gauge-type-id %d is invalid, available gauge type ids are %v", m.GaugeTypeId, ValidGaugeTypeIds)
		return fmt.Errorf(err)
	}

	if m.TriggerDuration <= 0 {
		return fmt.Errorf("duration should be positive: %d < 0", m.TriggerDuration)
	}
	if m.DepositAmount.Amount.IsNegative() || m.DepositAmount.Amount.IsZero() {
		return fmt.Errorf("invalid coin amount: %d < 0", m.DepositAmount.Amount)
	}

	return nil
}

func (m MsgCreateGauge) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgCreateGauge) GetSigners() []sdk.AccAddress {
	owner, _ := sdk.AccAddressFromBech32(m.From)
	return []sdk.AccAddress{owner}
}
