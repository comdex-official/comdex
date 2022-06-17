package types

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Liquidity params default values.
const (
	DefaultBatchSize        uint64        = 1
	DefaultTickPrecision    uint64        = 6
	DefaultMaxOrderLifespan time.Duration = 24 * time.Hour
	DefaultFeeDenom         string        = "ucmdx"
)

// Liquidity params default values.
var (
	DefaultMinInitialPoolCoinSupply = sdk.NewInt(1_000_000_000_000)
	DefaultPairCreationFee          = sdk.NewCoins(sdk.NewInt64Coin(DefaultFeeDenom, 200_000_000))
	DefaultPoolCreationFee          = sdk.NewCoins(sdk.NewInt64Coin(DefaultFeeDenom, 200_000_000))
	DefaultMinInitialDepositAmount  = sdk.NewInt(1000000)
	DefaultMaxPriceLimitRatio       = sdk.NewDecWithPrec(1, 1) // 10%
	DefaultSwapFeeRate              = sdk.NewDecWithPrec(3, 3) // 0.3%
	DefaultWithdrawFeeRate          = sdk.ZeroDec()
	DefaultDepositExtraGas          = sdk.Gas(60000)
	DefaultWithdrawExtraGas         = sdk.Gas(64000)
	DefaultOrderExtraGas            = sdk.Gas(37000)
	DefaultSwapFeeDistrDenom        = DefaultFeeDenom
	DefaultSwapFeeBurnRate          = sdk.NewDecWithPrec(5, 1) //50%
)

var (
	AppId                    = "AppId"
	BatchSize                = "BatchSize"
	TickPrecision            = "TickPrecision"
	FeeCollectorAddress      = "FeeCollectorAddress"
	DustCollectorAddress     = "DustCollectorAddress"
	MinInitialPoolCoinSupply = "MinInitialPoolCoinSupply"
	PairCreationFee          = "PairCreationFee"
	PoolCreationFee          = "PoolCreationFee"
	MinInitialDepositAmount  = "MinInitialDepositAmount"
	MaxPriceLimitRatio       = "MaxPriceLimitRatio"
	MaxOrderLifespan         = "MaxOrderLifespan"
	SwapFeeRate              = "SwapFeeRate"
	WithdrawFeeRate          = "WithdrawFeeRate"
	DepositExtraGas          = "DepositExtraGas"
	WithdrawExtraGas         = "WithdrawExtraGas"
	OrderExtraGas            = "OrderExtraGas"
	SwapFeeDistrDenom        = "SwapFeeDistrDenom"
	SwapFeeBurnRate          = "SwapFeeBurnRate"
)

var UpdatableKeys = []string{
	BatchSize,
	TickPrecision,
	MinInitialPoolCoinSupply,
	PairCreationFee,
	PoolCreationFee,
	MinInitialDepositAmount,
	MaxPriceLimitRatio,
	MaxOrderLifespan,
	SwapFeeRate,
	WithdrawFeeRate,
	DepositExtraGas,
	WithdrawExtraGas,
	OrderExtraGas,
	SwapFeeDistrDenom,
	SwapFeeBurnRate,
}

// DeriveFeeCollectorAddress returns a unique address of the fee collector.
func DeriveFeeCollectorAddress(appID uint64) sdk.AccAddress {
	return DeriveAddress(
		AddressType32Bytes,
		ModuleName,
		strings.Join([]string{FeeCollectorAddressPrefix, strconv.FormatUint(appID, 10)}, ModuleAddressNameSplitter))
}

// DeriveFeeCollectorAddress returns a unique address of the fee collector.
func DeriveDustCollectorAddress(appID uint64) sdk.AccAddress {
	return DeriveAddress(
		AddressType32Bytes,
		ModuleName,
		strings.Join([]string{DustCollectorAddress, strconv.FormatUint(appID, 10)}, ModuleAddressNameSplitter))
}

// DefaultParams returns a default params for the liquidity module.
func DefaultGenericParams(appID uint64) GenericParams {
	return GenericParams{
		AppId:                    appID,
		BatchSize:                DefaultBatchSize,
		TickPrecision:            DefaultTickPrecision,
		FeeCollectorAddress:      DeriveFeeCollectorAddress(appID).String(),
		DustCollectorAddress:     DeriveDustCollectorAddress(appID).String(),
		MinInitialPoolCoinSupply: DefaultMinInitialPoolCoinSupply,
		PairCreationFee:          DefaultPairCreationFee,
		PoolCreationFee:          DefaultPoolCreationFee,
		MinInitialDepositAmount:  DefaultMinInitialDepositAmount,
		MaxPriceLimitRatio:       DefaultMaxPriceLimitRatio,
		MaxOrderLifespan:         DefaultMaxOrderLifespan,
		SwapFeeRate:              DefaultSwapFeeRate,
		WithdrawFeeRate:          DefaultWithdrawFeeRate,
		DepositExtraGas:          DefaultDepositExtraGas,
		WithdrawExtraGas:         DefaultWithdrawExtraGas,
		OrderExtraGas:            DefaultOrderExtraGas,
		SwapFeeDistrDenom:        DefaultSwapFeeDistrDenom,
		SwapFeeBurnRate:          DefaultSwapFeeBurnRate,
	}
}

func KeyParseValidateFuncMap() map[string][]interface{} {
	return map[string][]interface{}{
		AppId:                    {ParseStringToUint, validateAppId},
		BatchSize:                {ParseStringToUint, validateBatchSize},
		TickPrecision:            {ParseStringToUint, validateTickPrecision},
		FeeCollectorAddress:      {ParseString, validateFeeCollectorAddress},
		DustCollectorAddress:     {ParseString, validateDustCollectorAddress},
		MinInitialPoolCoinSupply: {ParseStringToInt, validateMinInitialPoolCoinSupply},
		PairCreationFee:          {ParseStringToCoins, validatePairCreationFee},
		PoolCreationFee:          {ParseStringToCoins, validatePoolCreationFee},
		MinInitialDepositAmount:  {ParseStringToInt, validateMinInitialDepositAmount},
		MaxPriceLimitRatio:       {ParseStringToDec, validateMaxPriceLimitRatio},
		MaxOrderLifespan:         {ParseStringToDuration, validateMaxOrderLifespan},
		SwapFeeRate:              {ParseStringToDec, validateSwapFeeRate},
		WithdrawFeeRate:          {ParseStringToDec, validateWithdrawFeeRate},
		DepositExtraGas:          {ParseStringToGas, validateExtraGas},
		WithdrawExtraGas:         {ParseStringToGas, validateExtraGas},
		OrderExtraGas:            {ParseStringToGas, validateExtraGas},
		SwapFeeDistrDenom:        {ParseString, validateSwapFeeDistrDenom},
		SwapFeeBurnRate:          {ParseStringToDec, validateSwapFeeBurnRate},
	}
}

func ParseString(value string) (interface{}, error) {
	return value, nil
}

func ParseStringToUint(value string) (interface{}, error) {
	return strconv.ParseUint(value, 10, 64)
}

func ParseStringToInt(value string) (interface{}, error) {
	parsedValue, ok := sdk.NewIntFromString(value)
	if !ok {
		return sdk.Int{}, fmt.Errorf("invalid parameter type: %T", value)
	}
	return parsedValue, nil
}

func ParseStringToCoins(value string) (interface{}, error) {
	return sdk.ParseCoinsNormalized(value)
}

func ParseStringToDec(value string) (interface{}, error) {
	return sdk.NewDecFromStr(value)
}

func ParseStringToDuration(value string) (interface{}, error) {
	return time.ParseDuration(value)
}

func ParseStringToGas(value string) (interface{}, error) {
	gas, err := ParseStringToUint(value)
	if err != nil {
		return sdk.Gas(0), nil
	}
	g, _ := gas.(uint64)
	return sdk.Gas(g), nil
}

func validateAppId(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("app id must be positive: %d", v)
	}

	return nil
}

func validateBatchSize(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %d", i)
	}

	if v == 0 {
		return fmt.Errorf("batch size must be positive: %d", v)
	}

	return nil
}

func validateTickPrecision(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateFeeCollectorAddress(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if _, err := sdk.AccAddressFromBech32(v); err != nil {
		return fmt.Errorf("invalid fee collector address: %w", err)
	}

	return nil
}

func validateDustCollectorAddress(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if _, err := sdk.AccAddressFromBech32(v); err != nil {
		return fmt.Errorf("invalid dust collector address: %w", err)
	}

	return nil
}

func validateMinInitialPoolCoinSupply(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("min initial pool coin supply must not be nil")
	}

	if !v.IsPositive() {
		return fmt.Errorf("min initial pool coin supply must be positive: %s", v)
	}

	return nil
}

func validatePairCreationFee(i interface{}) error {
	v, ok := i.(sdk.Coins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if err := v.Validate(); err != nil {
		return fmt.Errorf("invalid pair creation fee: %w", err)
	}

	return nil
}

func validatePoolCreationFee(i interface{}) error {
	v, ok := i.(sdk.Coins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if err := v.Validate(); err != nil {
		return fmt.Errorf("invalid pool creation fee: %w", err)
	}

	return nil
}

func validateMinInitialDepositAmount(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("minimum initial deposit amount must not be negative: %s", v)
	}

	return nil
}

func validateMaxPriceLimitRatio(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("max price limit ratio must not be negative: %s", v)
	}

	return nil
}

func validateMaxOrderLifespan(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0 {
		return fmt.Errorf("max order lifespan must not be negative: %s", v)
	}

	return nil
}

func validateSwapFeeRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("swap fee rate must not be negative: %s", v)
	}

	if v.GTE(sdk.OneDec()) {
		return fmt.Errorf("swap fee rate cannot exceed 1 i.e 100 perc. : %s", v)
	}

	return nil
}

func validateWithdrawFeeRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("withdraw fee rate must not be negative: %s", v)
	}
	if v.GTE(sdk.OneDec()) {
		return fmt.Errorf("withdraw fee rate cannot exceed 1 i.e 100 perc. : %s", v)
	}

	return nil
}

func validateExtraGas(i interface{}) error {
	_, ok := i.(sdk.Gas)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateSwapFeeDistrDenom(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateSwapFeeBurnRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("swap fee burn rate must not be negative: %s", v)
	}

	if v.GTE(sdk.OneDec()) {
		return fmt.Errorf("swap fee burn rate cannot exceed 1 i.e 100 perc. : %s", v)
	}

	return nil
}
