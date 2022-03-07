package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	KeyMintRewards = []byte("MintRewards")
)

var (
	DefaultMintRewards = []*MintingRewardsV1{
		{
			Id:                0,
			AllowedCollateral: "ucmdx",
			AllowedCassets:    []string{"ucgold", "ucsilver", "ucoil"},
			TotalRewards:      sdk.NewCoin("ucmdx", sdk.NewInt(0)),
			CassetMaxCap:      0,
			DurationDays:      0,
			IsActive:          false,
		},
	}
)

var _ paramtypes.ParamSet = (*Params)(nil)

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams() Params {
	return Params{
		MintRewards: DefaultMintRewards,
	}
}

func DefaultParams() Params {
	return NewParams()
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMintRewards, &p.MintRewards, validateMintRewards),
	}
}

func (p Params) Validate() error {
	return nil
}

func validateMintRewards(i interface{}) error {
	_, ok := i.([]*MintingRewardsV1)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}
