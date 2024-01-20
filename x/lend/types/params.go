package types

import (
	storetypes "cosmossdk.io/store/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

const (
	LendGas                       = storetypes.Gas(32763)
	WithdrawGas                   = storetypes.Gas(22763)
	DepositGas                    = storetypes.Gas(22763)
	CloseLendGas                  = storetypes.Gas(22763)
	BorrowAssetGas                = storetypes.Gas(22763)
	DrawAssetGas                  = storetypes.Gas(22763)
	RepayAssetGas                 = storetypes.Gas(22763)
	DepositBorrowAssetGas         = storetypes.Gas(22763)
	CloseBorrowAssetGas           = storetypes.Gas(22763)
	BorrowAssetAlternateGas       = storetypes.Gas(22763)
	CalculateInterestAndRewardGas = storetypes.Gas(22763)
	RepayWithdrawGas              = storetypes.Gas(22763)
)

const (
	AppName        = "commodo"
	AppID          = uint64(3)
	Uint64Zero     = uint64(0)
	Uint64Two      = uint64(2)
	Perc1          = string("0.2")
	Perc2          = string("0.9")
	DollarOneValue = string("1000000")
)

// ParamKeyTable the param key table for launch module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance.
func NewParams() Params {
	return Params{}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// Validate validates the set of params.
func (p Params) Validate() error {
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
