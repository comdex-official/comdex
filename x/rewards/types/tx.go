package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgCreateGauge = "create_gauge"
)

var _ sdk.Msg = &MsgCreateGauge{}

// NewMsgCreateGauge creates a message to add a new gauge.
func NewMsgCreateGauge(
	appID uint64,

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
		AppId:           appID,
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
		return fmt.Errorf("invalid coin amount: %s < 0", m.DepositAmount.Amount)
	}

	if m.DepositAmount.Amount.LT(sdk.NewIntFromUint64(m.TotalTriggers)) {
		return fmt.Errorf("deposit amount : %s smaller than total triggers %d", m.DepositAmount.Amount, m.TotalTriggers)
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

func NewMsgActivateExternalRewardsLockers(
	appMappingID uint64,
	assetID uint64,
	totalRewards sdk.Coin,
	durationDays, minLockupTimeSeconds int64,
	// nolint
	from sdk.AccAddress,
) *ActivateExternalRewardsLockers {
	return &ActivateExternalRewardsLockers{
		AppMappingId:         appMappingID,
		AssetId:              assetID,
		TotalRewards:         totalRewards,
		DurationDays:         durationDays,
		MinLockupTimeSeconds: minLockupTimeSeconds,
		Depositor:            from.String(),
	}
}

func (m *ActivateExternalRewardsLockers) Route() string {
	return RouterKey
}

func (m *ActivateExternalRewardsLockers) Type() string {
	return ModuleName
}

func (m *ActivateExternalRewardsLockers) ValidateBasic() error {
	if m.AppMappingId <= 0 {
		return fmt.Errorf("app id should be positive: %d > 0", m.AppMappingId)
	}
	if m.AssetId <= 0 {
		return fmt.Errorf("asset id should be positive: %d > 0", m.AssetId)
	}
	if m.TotalRewards.IsZero() {
		return fmt.Errorf("TotalRewards should be positive: > 0")
	}
	if m.DurationDays <= 0 {
		return fmt.Errorf("DurationDays should be positive: %d > 0", m.DurationDays)
	}
	if m.MinLockupTimeSeconds <= 0 {
		return fmt.Errorf("MinLockupTimeSeconds should be positive: %d > 0", m.MinLockupTimeSeconds)
	}
	return nil
}

func (m *ActivateExternalRewardsLockers) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *ActivateExternalRewardsLockers) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.GetDepositor())
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgActivateExternalVaultLockers(
	appMappingID uint64,
	extendedPairID uint64,
	totalRewards sdk.Coin,
	durationDays, minLockupTimeSeconds int64,
	// nolint
	from sdk.AccAddress,
) *ActivateExternalRewardsVault {
	return &ActivateExternalRewardsVault{
		AppMappingId:         appMappingID,
		Extended_Pair_Id:     extendedPairID,
		TotalRewards:         totalRewards,
		DurationDays:         durationDays,
		MinLockupTimeSeconds: minLockupTimeSeconds,
		Depositor:            from.String(),
	}
}

func (m *ActivateExternalRewardsVault) Route() string {
	return RouterKey
}

func (m *ActivateExternalRewardsVault) Type() string {
	return ModuleName
}

func (m *ActivateExternalRewardsVault) ValidateBasic() error {
	if m.AppMappingId <= 0 {
		return fmt.Errorf("app id should be positive: %d > 0", m.AppMappingId)
	}
	if m.Extended_Pair_Id <= 0 {
		return fmt.Errorf("asset id should be positive: %d > 0", m.Extended_Pair_Id)
	}
	if m.TotalRewards.IsZero() {
		return fmt.Errorf("TotalRewards should be positive: > 0")
	}
	if m.DurationDays <= 0 {
		return fmt.Errorf("DurationDays should be positive: %d > 0", m.DurationDays)
	}
	if m.MinLockupTimeSeconds <= 0 {
		return fmt.Errorf("MinLockupTimeSeconds should be positive: %d > 0", m.MinLockupTimeSeconds)
	}
	return nil
}

func (m *ActivateExternalRewardsVault) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *ActivateExternalRewardsVault) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.GetDepositor())
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}
