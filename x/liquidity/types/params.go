package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Liquidity params default values
const (
	DefaultBatchSize        uint32 = 1
	DefaultTickPrecision    uint32 = 3
	DefaultMaxOrderLifespan        = 24 * time.Hour
)

// Liquidity params default values
var (
	DefaultFeeCollectorAddress      = DeriveAddress(AddressType32Bytes, ModuleName, "FeeCollector")
	DefaultDustCollectorAddress     = DeriveAddress(AddressType32Bytes, ModuleName, "DustCollector")
	DefaultMinInitialPoolCoinSupply = sdk.NewInt(1_000_000_000_000)
	DefaultPairCreationFee          = sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000000))
	DefaultPoolCreationFee          = sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000000))
	DefaultMinInitialDepositAmount  = sdk.NewInt(1000000)
	DefaultMaxPriceLimitRatio       = sdk.NewDecWithPrec(1, 1) // 10%
	DefaultSwapFeeRate              = sdk.ZeroDec()
	DefaultWithdrawFeeRate          = sdk.ZeroDec()
	DefaultDepositExtraGas          = sdk.Gas(60000)
	DefaultWithdrawExtraGas         = sdk.Gas(64000)
	DefaultOrderExtraGas            = sdk.Gas(37000)
	DefaultPoolUnbondingDuration    = sdk.NewInt(21)
)

// General constants
const (
	PoolReserveAddressPrefix  = "PoolReserveAddress"
	PairEscrowAddressPrefix   = "PairEscrowAddress"
	ModuleAddressNameSplitter = "|"
)

var (
	// GlobalEscrowAddress is an escrow for deposit/withdraw requests.
	GlobalEscrowAddress = DeriveAddress(AddressType32Bytes, ModuleName, "GlobalEscrow")
)

var (
	KeyBatchSize                = []byte("BatchSize")
	KeyTickPrecision            = []byte("TickPrecision")
	KeyFeeCollectorAddress      = []byte("FeeCollectorAddress")
	KeyDustCollectorAddress     = []byte("DustCollectorAddress")
	KeyMinInitialPoolCoinSupply = []byte("MinInitialPoolCoinSupply")
	KeyPairCreationFee          = []byte("PairCreationFee")
	KeyPoolCreationFee          = []byte("PoolCreationFee")
	KeyMinInitialDepositAmount  = []byte("MinInitialDepositAmount")
	KeyMaxPriceLimitRatio       = []byte("MaxPriceLimitRatio")
	KeyMaxOrderLifespan         = []byte("MaxOrderLifespan")
	KeySwapFeeRate              = []byte("SwapFeeRate")
	KeyWithdrawFeeRate          = []byte("WithdrawFeeRate")
	KeyDepositExtraGas          = []byte("DepositExtraGas")
	KeyWithdrawExtraGas         = []byte("WithdrawExtraGas")
	KeyOrderExtraGas            = []byte("OrderExtraGas")
	KeyPoolUnbondingDuration    = []byte("PoolUnbondingDuration")
)

var _ paramstypes.ParamSet = (*Params)(nil)

func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default params for the liquidity module.
func DefaultParams() Params {
	return Params{
		BatchSize:                DefaultBatchSize,
		TickPrecision:            DefaultTickPrecision,
		FeeCollectorAddress:      DefaultFeeCollectorAddress.String(),
		DustCollectorAddress:     DefaultDustCollectorAddress.String(),
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
		PoolUnbondingDuration:    DefaultPoolUnbondingDuration,
	}
}

// ParamSetPairs implements ParamSet.
func (params *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyBatchSize, &params.BatchSize, validateBatchSize),
		paramstypes.NewParamSetPair(KeyTickPrecision, &params.TickPrecision, validateTickPrecision),
		paramstypes.NewParamSetPair(KeyFeeCollectorAddress, &params.FeeCollectorAddress, validateFeeCollectorAddress),
		paramstypes.NewParamSetPair(KeyDustCollectorAddress, &params.DustCollectorAddress, validateDustCollectorAddress),
		paramstypes.NewParamSetPair(KeyMinInitialPoolCoinSupply, &params.MinInitialPoolCoinSupply, validateMinInitialPoolCoinSupply),
		paramstypes.NewParamSetPair(KeyPairCreationFee, &params.PairCreationFee, validatePairCreationFee),
		paramstypes.NewParamSetPair(KeyPoolCreationFee, &params.PoolCreationFee, validatePoolCreationFee),
		paramstypes.NewParamSetPair(KeyMinInitialDepositAmount, &params.MinInitialDepositAmount, validateMinInitialDepositAmount),
		paramstypes.NewParamSetPair(KeyMaxPriceLimitRatio, &params.MaxPriceLimitRatio, validateMaxPriceLimitRatio),
		paramstypes.NewParamSetPair(KeyMaxOrderLifespan, &params.MaxOrderLifespan, validateMaxOrderLifespan),
		paramstypes.NewParamSetPair(KeySwapFeeRate, &params.SwapFeeRate, validateSwapFeeRate),
		paramstypes.NewParamSetPair(KeyWithdrawFeeRate, &params.WithdrawFeeRate, validateWithdrawFeeRate),
		paramstypes.NewParamSetPair(KeyDepositExtraGas, &params.DepositExtraGas, validateExtraGas),
		paramstypes.NewParamSetPair(KeyWithdrawExtraGas, &params.WithdrawExtraGas, validateExtraGas),
		paramstypes.NewParamSetPair(KeyOrderExtraGas, &params.OrderExtraGas, validateExtraGas),
		paramstypes.NewParamSetPair(KeyPoolUnbondingDuration, &params.PoolUnbondingDuration, validatePoolUnbondingDuration),
	}
}

// Validate validates Params.
func (params Params) Validate() error {
	for _, field := range []struct {
		val          interface{}
		validateFunc func(i interface{}) error
	}{
		{params.BatchSize, validateBatchSize},
		{params.TickPrecision, validateTickPrecision},
		{params.FeeCollectorAddress, validateFeeCollectorAddress},
		{params.DustCollectorAddress, validateDustCollectorAddress},
		{params.MinInitialPoolCoinSupply, validateMinInitialPoolCoinSupply},
		{params.PairCreationFee, validatePairCreationFee},
		{params.PoolCreationFee, validatePoolCreationFee},
		{params.MinInitialDepositAmount, validateMinInitialDepositAmount},
		{params.MaxPriceLimitRatio, validateMaxPriceLimitRatio},
		{params.MaxOrderLifespan, validateMaxOrderLifespan},
		{params.SwapFeeRate, validateSwapFeeRate},
		{params.WithdrawFeeRate, validateWithdrawFeeRate},
		{params.DepositExtraGas, validateExtraGas},
		{params.WithdrawExtraGas, validateExtraGas},
		{params.OrderExtraGas, validateExtraGas},
		{params.PoolUnbondingDuration, validatePoolUnbondingDuration},
	} {
		if err := field.validateFunc(field.val); err != nil {
			return err
		}
	}
	return nil
}

func validateBatchSize(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("batch size must be positive: %d", v)
	}

	return nil
}

func validateTickPrecision(i interface{}) error {
	_, ok := i.(uint32)
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

	return nil
}

func validateExtraGas(i interface{}) error {
	_, ok := i.(sdk.Gas)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validatePoolUnbondingDuration(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok{
		return fmt.Errorf("invalid parameter type: %T", i)

	}
	if v.IsZero() {
		return fmt.Errorf("Unbonding duration cannot be 0 day: %d", v)
	}
	return nil
}