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
	//nolint
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

func NewMsgWhitelistAsset(appMappingId uint64, from sdk.AccAddress, assetId []uint64) *WhitelistAsset {
	return &WhitelistAsset{
		AppMappingId: appMappingId,
		From:         from.String(),
		AssetId:      assetId,
	}
}

func (m *WhitelistAsset) Route() string {
	return RouterKey
}

func (m *WhitelistAsset) Type() string {
	return ModuleName
}

func (m *WhitelistAsset) ValidateBasic() error {

	return nil
}

func (m *WhitelistAsset) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *WhitelistAsset) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgRemoveWhitelistAsset(appMappingId uint64, from sdk.AccAddress, assetId uint64) *RemoveWhitelistAsset {
	return &RemoveWhitelistAsset{
		AppMappingId: appMappingId,
		From:         from.String(),
		AssetId:      assetId,
	}
}

func (m *RemoveWhitelistAsset) Route() string {
	return RouterKey
}

func (m *RemoveWhitelistAsset) Type() string {
	return ModuleName
}

func (m *RemoveWhitelistAsset) ValidateBasic() error {

	return nil
}

func (m *RemoveWhitelistAsset) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *RemoveWhitelistAsset) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgWhitelistAppIdVault(appMappingId uint64, from sdk.AccAddress) *WhitelistAppIdVault {
	return &WhitelistAppIdVault{
		AppMappingId: appMappingId,
		From:         from.String(),
	}
}

func (m *WhitelistAppIdVault) Route() string {
	return RouterKey
}

func (m *WhitelistAppIdVault) Type() string {
	return ModuleName
}

func (m *WhitelistAppIdVault) ValidateBasic() error {

	return nil
}

func (m *WhitelistAppIdVault) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *WhitelistAppIdVault) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgRemoveWhitelistAppIdVault(appMappingId uint64, from sdk.AccAddress) *RemoveWhitelistAppIdVault {
	return &RemoveWhitelistAppIdVault{
		AppMappingId: appMappingId,
		From:         from.String(),
	}
}

func (m *RemoveWhitelistAppIdVault) Route() string {
	return RouterKey
}

func (m *RemoveWhitelistAppIdVault) Type() string {
	return ModuleName
}

func (m *RemoveWhitelistAppIdVault) ValidateBasic() error {

	return nil
}

func (m *RemoveWhitelistAppIdVault) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *RemoveWhitelistAppIdVault) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgActivateExternalRewardsLockers(appMappingId uint64, AssetId uint64, TotalRewards sdk.Coin, DurationDays, MinLockupTimeSeconds int64, from sdk.AccAddress) *ActivateExternalRewardsLockers {
	return &ActivateExternalRewardsLockers{
		AppMappingId:         appMappingId,
		AssetId:              AssetId,
		TotalRewards:         TotalRewards,
		DurationDays:         DurationDays,
		MinLockupTimeSeconds: MinLockupTimeSeconds,
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

	return nil
}

func (m *ActivateExternalRewardsLockers) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *ActivateExternalRewardsLockers) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.Depositor)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}

func NewMsgActivateExternalVaultLockers(appMappingId uint64, extendedPairId uint64, TotalRewards sdk.Coin, DurationDays, MinLockupTimeSeconds int64, from sdk.AccAddress) *ActivateExternalRewardsVault {
	return &ActivateExternalRewardsVault{
		AppMappingId:         appMappingId,
		Extended_Pair_Id:     extendedPairId,
		TotalRewards:         TotalRewards,
		DurationDays:         DurationDays,
		MinLockupTimeSeconds: MinLockupTimeSeconds,
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

	return nil
}

func (m *ActivateExternalRewardsVault) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *ActivateExternalRewardsVault) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.Depositor)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{from}
}
