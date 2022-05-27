package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	KeyLiquidationPenaltyPercent       = []byte("LiquidationPenaltyPercent")
	KeyAuctionDiscountPercent          = []byte("AuctionDiscountPercent")
	KeyAuctionDurationSeconds          = []byte("AuctionDurationSeconds")
	KeyDebtMintTokenDecreasePercentage = []byte("DebtMintTokenDecreasePercentage")
	KeyBuffer                          = []byte("Buffer")
	KeyCusp                            = []byte("Cusp")
	KeyTau                             = []byte("Tau")
	KeyDutchDecreasePercentage         = []byte("DutchDecreasePercentage")
	KeyChost                           = []byte("Chost")
	KeyStep                            = []byte("Step")
	KeyPriceFunctionType               = []byte("PriceFunctionType")
	KeySurplusId                       = []byte("SurplusId")
	KeyDebtId                          = []byte("DebtId")
	KeyDutchId                         = []byte("DutchId")
)

var (
	DefaultLiquidationPenaltyPercent       = "0.15"
	DefaultAuctionDiscountPercent          = "0.05"
	DefaultAuctionDurationSeconds          = uint64(60)
	DefaultDebtMintTokenDecreasePercentage = sdk.MustNewDecFromStr("0.03")
	DefaultBuffer                          = sdk.MustNewDecFromStr("1.2")
	DefaultCusp                            = sdk.MustNewDecFromStr("0.6")
	DefaultTau                             = sdk.NewInt(3600)
	DefaultDutchDecreasePercentage         = sdk.MustNewDecFromStr("0.01")
	DefaultChost                           = sdk.MustNewDecFromStr("10")
	DefaultStep                            = sdk.NewInt(360)
	DefaultPriceFunctionType               = uint64(1)
	DefaultSurplusId                       = uint64(1)
	DefaultDebtId                          = uint64(2)
	DefaultDutchId                         = uint64(3)
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(liquidationPenaltyPercent string,
	auctionDiscountPercent string,
	auctionDurationSeconds uint64,
	debtMintTokenDecreasePercentage sdk.Dec,
	buffer sdk.Dec,
	cusp sdk.Dec,
	tau sdk.Int,
	dutchDecreasePercentage sdk.Dec,
	chost sdk.Dec,
	step sdk.Int, priceFunctionType uint64, surplusId uint64, debtId uint64, dutchId uint64) Params {
	return Params{
		LiquidationPenaltyPercent:       liquidationPenaltyPercent,
		AuctionDiscountPercent:          auctionDiscountPercent,
		AuctionDurationSeconds:          auctionDurationSeconds,
		DebtMintTokenDecreasePercentage: debtMintTokenDecreasePercentage,
		Buffer:                          buffer,
		Cusp:                            cusp,
		Tau:                             tau,
		DutchDecreasePercentage:         dutchDecreasePercentage,
		Chost:                           chost,
		Step:                            step,
		PriceFunctionType:               priceFunctionType,
		SurplusId:                       surplusId,
		DebtId:                          debtId,
		DutchId:                         dutchId,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultLiquidationPenaltyPercent,
		DefaultAuctionDiscountPercent,
		DefaultAuctionDurationSeconds,
		DefaultDebtMintTokenDecreasePercentage,
		DefaultBuffer,
		DefaultCusp,
		DefaultTau,
		DefaultDutchDecreasePercentage,
		DefaultChost,
		DefaultStep,
		DefaultPriceFunctionType,
		DefaultSurplusId,
		DefaultDebtId,
		DefaultDutchId,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyLiquidationPenaltyPercent, &p.LiquidationPenaltyPercent, validateLiquidationPenalty),
		paramtypes.NewParamSetPair(KeyAuctionDiscountPercent, &p.AuctionDiscountPercent, validateAuctionDiscount),
		paramtypes.NewParamSetPair(KeyAuctionDurationSeconds, &p.AuctionDurationSeconds, validateAuctionDuration),
		paramtypes.NewParamSetPair(KeyDebtMintTokenDecreasePercentage, &p.DebtMintTokenDecreasePercentage, validatePercentage),
		paramtypes.NewParamSetPair(KeyBuffer, &p.Buffer, validateBuffer),
		paramtypes.NewParamSetPair(KeyCusp, &p.Cusp, validateCusp),
		paramtypes.NewParamSetPair(KeyTau, &p.Tau, validateTau),
		paramtypes.NewParamSetPair(KeyDutchDecreasePercentage, &p.DutchDecreasePercentage, validateDutchDecreasePercentage),
		paramtypes.NewParamSetPair(KeyChost, &p.Chost, validateChost),
		paramtypes.NewParamSetPair(KeyStep, &p.Step, validateStep),
		paramtypes.NewParamSetPair(KeyPriceFunctionType, &p.PriceFunctionType, validatePriceFunctionType),
		paramtypes.NewParamSetPair(KeySurplusId, &p.SurplusId, validateAuctionId),
		paramtypes.NewParamSetPair(KeyDebtId, &p.DebtId, validateAuctionId),
		paramtypes.NewParamSetPair(KeyDutchId, &p.DutchId, validateAuctionId),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	for _, v := range []struct {
		value     interface{}
		validator func(interface{}) error
	}{
		{p.LiquidationPenaltyPercent, validateLiquidationPenalty},
		{p.AuctionDiscountPercent, validateAuctionDiscount},
		{p.AuctionDurationSeconds, validateAuctionDuration},
		{p.DebtMintTokenDecreasePercentage, validatePercentage},
		{p.Buffer, validateBuffer},
		{p.Cusp, validateCusp},
		{p.Tau, validateTau},
		{p.DutchDecreasePercentage, validateDutchDecreasePercentage},
		{p.Chost, validateChost},
		{p.Step, validateStep},
		{p.PriceFunctionType, validatePriceFunctionType},
		{p.SurplusId, validateAuctionId},
		{p.DebtId, validateAuctionId},
		{p.DutchId, validateAuctionId},
	} {
		if err := v.validator(v.value); err != nil {
			return err
		}
	}
	return nil
}

// String implements the Stringer interface.
// func (p Params) String() string {
// 	out, _ := yaml.Marshal(p)
// 	return string(out)
// }

func validateLiquidationPenalty(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	q, _ := sdk.NewDecFromStr(v)
	u, _ := sdk.NewDecFromStr("0.01")
	if q.LT(u) {
		return fmt.Errorf("liquidation penalty cannot be less than 1 percent")
	}
	return nil
}

func validateAuctionDiscount(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	q, _ := sdk.NewDecFromStr(v)
	u, _ := sdk.NewDecFromStr("0.01")
	if q.LT(u) {
		return fmt.Errorf("auction discount cannot be less than 1 percent")
	}
	return nil
}

func validateAuctionDuration(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v < 1 {
		return fmt.Errorf("auction duration cannot be less than 1 hour")
	}
	return nil
}

func validatePercentage(i interface{}) error {
	q, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	u, _ := sdk.NewDecFromStr("0.01")
	if q.LT(u) {
		return fmt.Errorf("decrease percentage cannot be less than 1 percent")
	}
	return nil
}

func validateBuffer(i interface{}) error {
	q, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	u, _ := sdk.NewDecFromStr("1")
	if q.LTE(u) {
		return fmt.Errorf("buffer cannot be less than 1")
	}
	return nil
}

func validateCusp(i interface{}) error {
	q, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	u, _ := sdk.NewDecFromStr("0.01")
	if q.LT(u) {
		return fmt.Errorf("cusp cannot be less than 0.01")
	}
	return nil
}

func validateDutchDecreasePercentage(i interface{}) error {
	q, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	u, _ := sdk.NewDecFromStr("0.5")
	if q.GT(u) {
		return fmt.Errorf("dutch decrease percentage cannot be less than 10.5")
	}
	return nil
}

func validateChost(i interface{}) error {
	q, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	u, _ := sdk.NewDecFromStr("1")
	if q.LT(u) {
		return fmt.Errorf("chost cannot be less than 1 ")
	}
	return nil
}

func validatePriceFunctionType(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v < 0 {
		return fmt.Errorf("price function type cannot be less than 0")
	}
	return nil
}

func validateTau(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.LT(sdk.NewInt(1)) {
		return fmt.Errorf("tau cannot be less than 1 second")
	}
	return nil
}

func validateStep(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.LT(sdk.NewInt(1)) {
		return fmt.Errorf("step cannot be less than 1")
	}
	return nil
}

func validateAuctionId(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v < 0 {
		return fmt.Errorf("auction id cannot be less than 0")
	}
	return nil
}
