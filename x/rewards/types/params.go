package types

import (
	fmt "fmt"
	"strings"
	"time"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	KeyMintRewardTimeStamp = []byte("MintRewardTimestamp")
)

var (
	DefaultMintRewardTimeStamp = "18:30:00"
)

var _ paramtypes.ParamSet = (*Params)(nil)

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(mintRewardTimeStamp string) Params {
	return Params{
		MintRewardTimestamp: mintRewardTimeStamp,
	}
}

func DefaultParams() Params {
	return NewParams(DefaultMintRewardTimeStamp)
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMintRewardTimeStamp, &p.MintRewardTimestamp, validateMintRewardTimeStamp),
	}
}

func (p Params) Validate() error {
	for _, v := range []struct {
		value     interface{}
		validator func(interface{}) error
	}{
		{p.MintRewardTimestamp, validateMintRewardTimeStamp},
	} {
		if err := v.validator(v.value); err != nil {
			return err
		}
	}
	return nil
}

func validateMintRewardTimeStamp(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	const layoutTime = "15:04:05"
	_, err := time.Parse(layoutTime, strings.TrimSpace(v))
	if err != nil {
		return err
	}
	return nil
}
