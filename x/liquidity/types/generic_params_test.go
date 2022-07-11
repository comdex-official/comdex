package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/liquidity/types"
)

func TestGenericParamsValidate(t *testing.T) {
	for _, tc := range []struct {
		name     string
		malleate func(*types.GenericParams)
		errStr   string
	}{
		{
			"default params",
			func(params *types.GenericParams) {},
			"",
		},
		{
			"zero AppId",
			func(params *types.GenericParams) {
				params.AppId = 0
			},
			"app id must be positive: 0",
		},
		{
			"zero BatchSize",
			func(params *types.GenericParams) {
				params.BatchSize = 0
			},
			"batch size must be positive: 0",
		},
		{
			"invalid FeeCollectorAddress",
			func(params *types.GenericParams) {
				params.FeeCollectorAddress = "invalidaddr"
			},
			"invalid fee collector address: decoding bech32 failed: invalid separator index -1",
		},
		{
			"invalid DustCollectorAddress",
			func(params *types.GenericParams) {
				params.DustCollectorAddress = "invalidaddr"
			},
			"invalid dust collector address: decoding bech32 failed: invalid separator index -1",
		},
		{
			"negative MinInitialPoolCoinSupply",
			func(params *types.GenericParams) {
				params.MinInitialPoolCoinSupply = sdk.NewInt(-1)
			},
			"min initial pool coin supply must be positive: -1",
		},
		{
			"zero MinInitialPoolCoinSupply",
			func(params *types.GenericParams) {
				params.MinInitialPoolCoinSupply = sdk.ZeroInt()
			},
			"min initial pool coin supply must be positive: 0",
		},
		{
			"invalid PairCreationFee",
			func(params *types.GenericParams) {
				params.PairCreationFee = sdk.Coins{sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: sdk.ZeroInt()}}
			},
			"invalid pair creation fee: coin 0stake amount is not positive",
		},
		{
			"invalid PoolCreationFee",
			func(params *types.GenericParams) {
				params.PoolCreationFee = sdk.Coins{sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: sdk.ZeroInt()}}
			},
			"invalid pool creation fee: coin 0stake amount is not positive",
		},
		{
			"negative MinInitialDepositAmount",
			func(params *types.GenericParams) {
				params.MinInitialDepositAmount = sdk.NewInt(-1)
			},
			"minimum initial deposit amount must not be negative: -1",
		},
		{
			"negative MaxPriceLimitRatio",
			func(params *types.GenericParams) {
				params.MaxPriceLimitRatio = sdk.NewDec(-1)
			},
			"max price limit ratio must not be negative: -1.000000000000000000",
		},
		{
			"negative MaxOrderLifespan",
			func(params *types.GenericParams) {
				params.MaxOrderLifespan = -1
			},
			"max order lifespan must not be negative: -1ns",
		},
		{
			"negative SwapFeeRate",
			func(params *types.GenericParams) {
				params.SwapFeeRate = sdk.NewDec(-1)
			},
			"swap fee rate must not be negative: -1.000000000000000000",
		},
		{
			"overflow SwapFeeRate",
			func(params *types.GenericParams) {
				params.SwapFeeRate = sdk.NewDec(2)
			},
			"swap fee rate cannot exceed 1 i.e 100 perc. : 2.000000000000000000",
		},
		{
			"negative WithdrawFeeRate",
			func(params *types.GenericParams) {
				params.WithdrawFeeRate = sdk.NewDec(-1)
			},
			"withdraw fee rate must not be negative: -1.000000000000000000",
		},
		{
			"overflow WithdrawFeeRate",
			func(params *types.GenericParams) {
				params.WithdrawFeeRate = sdk.NewDec(2)
			},
			"withdraw fee rate cannot exceed 1 i.e 100 perc. : 2.000000000000000000",
		},
		{
			"negative SwapFeeBurnRate",
			func(params *types.GenericParams) {
				params.SwapFeeBurnRate = sdk.NewDec(-1)
			},
			"swap fee burn rate must not be negative: -1.000000000000000000",
		},
		{
			"overflow SwapFeeBurnRate",
			func(params *types.GenericParams) {
				params.SwapFeeBurnRate = sdk.NewDec(2)
			},
			"swap fee burn rate cannot exceed 1 i.e 100 perc. : 2.000000000000000000",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			params := types.DefaultGenericParams(1)
			tc.malleate(&params)
			err := params.Validate()
			if tc.errStr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.errStr)
			}
		})
	}
}
