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
	startTime time.Time,
	gaugeTypeID uint64,
	triggerDuration time.Duration,
	depositAmount sdk.Coin,
	totalTriggers uint64,
) *MsgCreateGauge {
	return &MsgCreateGauge{
		From:            from.String(),
		StartTime:       startTime,
		GaugeTypeId:     gaugeTypeID,
		TriggerDuration: triggerDuration,
		DepositAmount:   depositAmount,
		TotalTriggers:   totalTriggers,
		Kind:            nil,
	}
}

// Route Implements MsgCreateGauge.
func (m MsgCreateGauge) Route() string { return RouterKey }

// Type Implements MsgCreateGauge.
func (m MsgCreateGauge) Type() string { return TypeMsgCreateGauge }

// ValidateBasic Implements baic validations for MsgCreateGauge.
func (m MsgCreateGauge) ValidateBasic() error {
	isValidGaugeTypeID := false
	for _, gaugeTypeID := range ValidGaugeTypeIds {
		if gaugeTypeID == m.GaugeTypeId {
			isValidGaugeTypeID = true
			break
		}
	}
	if !isValidGaugeTypeID {
		err := fmt.Sprintf("gauge-type-id %d is invalid, available gauge type ids are %v", m.GaugeTypeId, ValidGaugeTypeIds)
		return fmt.Errorf(err)
	}

	if m.TriggerDuration <= 0 {
		return fmt.Errorf("duration should be positive: %d < 0", m.TriggerDuration)
	}
	if m.DepositAmount.Amount.IsNegative() || m.DepositAmount.Amount.IsZero() {
		return fmt.Errorf("invalid coin amount: %d < 0", m.DepositAmount.Amount)
	}

	if m.DepositAmount.Amount.LT(sdk.NewIntFromUint64(m.TotalTriggers)) {
		return fmt.Errorf("deposit amount : %d smaller than total triggers %d", m.DepositAmount.Amount, m.TotalTriggers)
	}

	return nil
}

// GetSignBytes Implements MsgCreateGauge.
func (m MsgCreateGauge) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// GetSigners Implements MsgCreateGauge.
func (m MsgCreateGauge) GetSigners() []sdk.AccAddress {
	owner, _ := sdk.AccAddressFromBech32(m.From)
	return []sdk.AccAddress{owner}
}
