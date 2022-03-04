package types

import "github.com/cosmos/cosmos-sdk/types"

func NewGenesisState(assets []Asset, pairs []Pair, params Params) *GenesisState {
	return &GenesisState{
		Assets:  assets,
		Pairs:   pairs,
		Params:  params,
	}
}

func DefaultGenesisState() *GenesisState {
	liquidationRatio := DefaultParams().LiquidationRatio
	newliqRatio,_ := types.NewDecFromStr(liquidationRatio)
	return NewGenesisState(
		[]Asset{
			{1 ,"ATOM", "ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9", 1000000},
			{2, "XAU", "ucgold", 1000000},
			{3, "XAG", "ucsilver", 1000000},
			{4, "OIL", "ucoil", 1000000},
			{5, "UST", "ibc/4294C3DB67564CF4A0B2BFACC8415A59B38243F6FF9E288FBA34F9B4823BA16E", 1000000},
			{6, "CMDX", "ucmdx", 1000000}},
		[]Pair{
			{1,6,2, newliqRatio},
			{2,6,3, newliqRatio},
			{3,6,4, newliqRatio},
			{4,1,2, newliqRatio},
			{5,1,3, newliqRatio},
			{6,1,4, newliqRatio},
			{7,5,2, newliqRatio},
			{8,5,3, newliqRatio},
			{9,5,4, newliqRatio},
		},
		DefaultParams(),
	)
}

func ValidateGenesis(_ *GenesisState) error {
	return nil
}
