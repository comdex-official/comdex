package types

import (
	sdkmath "cosmossdk.io/math"
	"fmt"
)

func (m *Extended_Pair) Validate() error {
	if m.AssetIn == 0 {
		return fmt.Errorf("asset_in cannot be zero")
	}
	if m.AssetOut == 0 {
		return fmt.Errorf("asset_out cannot be zero")
	}
	if m.AssetOutPoolID == 0 {
		return fmt.Errorf("AssetOutPoolID cannot be zero")
	}
	if m.MinUsdValueLeft == 0 {
		return fmt.Errorf("MinUsdValueLeft cannot be zero")
	}

	return nil
}

func (m *Pool) Validate() error {
	if len(m.CPoolName) >= 20 {
		return ErrInvalidLengthCPoolName
	}
	if m.AssetData == nil {
		return fmt.Errorf("AssetData cannot be nil")
	}
	return nil
}

func (m *AssetToPairMapping) Validate() error {
	if m.AssetID == 0 {
		return fmt.Errorf("assetID cannot be zero")
	}
	if m.PairID == nil {
		return fmt.Errorf("PairIDs cannot be nil")
	}
	return nil
}

func (m *AssetToPairSingleMapping) Validate() error {
	if m.AssetID == 0 {
		return fmt.Errorf("assetID cannot be zero")
	}
	if m.PairID == 0 {
		return fmt.Errorf("PairIDs cannot be zero")
	}
	if m.PoolID == 0 {
		return fmt.Errorf("PoolIDs cannot be zero")
	}
	return nil
}

func (m *AssetRatesParams) Validate() error {
	if m.AssetID == 0 {
		return fmt.Errorf("assetID cannot be zero")
	}
	if m.UOptimal.LTE(sdkmath.LegacyZeroDec()) {
		return fmt.Errorf("UOptimal cannot be zero")
	}
	if m.Base.LTE(sdkmath.LegacyZeroDec()) {
		return fmt.Errorf("base cannot be zero")
	}
	if m.Slope1.LTE(sdkmath.LegacyZeroDec()) {
		return fmt.Errorf("slope1 cannot be zero")
	}
	if m.Slope2.LTE(sdkmath.LegacyZeroDec()) {
		return fmt.Errorf("slope2 cannot be zero")
	}
	if m.StableBase.LT(sdkmath.LegacyZeroDec()) {
		return fmt.Errorf("StableBase cannot be less than zero")
	}
	if m.StableSlope1.LT(sdkmath.LegacyZeroDec()) {
		return fmt.Errorf("StableSlope1 cannot be less than zero")
	}
	if m.StableSlope2.LT(sdkmath.LegacyZeroDec()) {
		return fmt.Errorf("StableSlope2 cannot be less than zero")
	}
	if m.LiquidationThreshold.LTE(sdkmath.LegacyZeroDec()) {
		return fmt.Errorf("LiquidationThreshold cannot be zero")
	}
	if m.LiquidationBonus.LTE(sdkmath.LegacyZeroDec()) {
		return fmt.Errorf("LiquidationBonus cannot be zero")
	}
	if m.LiquidationPenalty.LTE(sdkmath.LegacyZeroDec()) {
		return fmt.Errorf("LiquidationPenalty cannot be zero")
	}
	if m.Ltv.LTE(sdkmath.LegacyZeroDec()) {
		return fmt.Errorf("ltv cannot be zero")
	}
	if m.ReserveFactor.LTE(sdkmath.LegacyZeroDec()) {
		return fmt.Errorf("ReserveFactor cannot be zero")
	}
	if m.CAssetID == 0 {
		return fmt.Errorf("cAssetID cannot be zero")
	}
	return nil
}

func (m *PoolPairs) Validate() error {
	if len(m.CPoolName) >= 20 {
		return ErrInvalidLengthCPoolName
	}
	if m.AssetData == nil {
		return fmt.Errorf("AssetData cannot be nil")
	}
	return nil
}

func (m *AssetRatesPoolPairs) Validate() error {
	if m.AssetID == 0 {
		return fmt.Errorf("assetID cannot be zero")
	}
	if m.UOptimal.LTE(sdkmath.LegacyZeroDec()) {
		return fmt.Errorf("UOptimal cannot be zero")
	}
	if m.Base.LTE(sdkmath.LegacyZeroDec()) {
		return fmt.Errorf("base cannot be zero")
	}
	if m.Slope1.LTE(sdkmath.LegacyZeroDec()) {
		return fmt.Errorf("slope1 cannot be zero")
	}
	if m.Slope2.LTE(sdkmath.LegacyZeroDec()) {
		return fmt.Errorf("slope2 cannot be zero")
	}
	if m.StableBase.LT(sdkmath.LegacyZeroDec()) {
		return fmt.Errorf("StableBase cannot be less than zero")
	}
	if m.StableSlope1.LT(sdkmath.LegacyZeroDec()) {
		return fmt.Errorf("StableSlope1 cannot be less than zero")
	}
	if m.StableSlope2.LT(sdkmath.LegacyZeroDec()) {
		return fmt.Errorf("StableSlope2 cannot be less than zero")
	}
	if m.LiquidationThreshold.LTE(sdkmath.LegacyZeroDec()) {
		return fmt.Errorf("LiquidationThreshold cannot be zero")
	}
	if m.LiquidationBonus.LTE(sdkmath.LegacyZeroDec()) {
		return fmt.Errorf("LiquidationBonus cannot be zero")
	}
	if m.LiquidationPenalty.LTE(sdkmath.LegacyZeroDec()) {
		return fmt.Errorf("LiquidationPenalty cannot be zero")
	}
	if m.Ltv.LTE(sdkmath.LegacyZeroDec()) {
		return fmt.Errorf("ltv cannot be zero")
	}
	if m.ReserveFactor.LTE(sdkmath.LegacyZeroDec()) {
		return fmt.Errorf("ReserveFactor cannot be zero")
	}
	if m.CAssetID == 0 {
		return fmt.Errorf("cAssetID cannot be zero")
	}
	if len(m.CPoolName) >= 20 {
		return ErrInvalidLengthCPoolName
	}
	if m.AssetData == nil {
		return fmt.Errorf("AssetData cannot be nil")
	}
	return nil
}

func (m *PoolDepreciate) Validate() error {
	if m.IndividualPoolDepreciate == nil {
		return fmt.Errorf("PoolDepreciate struct cannot be empty")
	}
	return nil
}

func (m *EModePairsForProposal) Validate() error {
	if m.EModePairs == nil {
		return fmt.Errorf("EmodePairs struct cannot be empty")
	}
	return nil
}
