package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/comdex-official/comdex/x/liquidity/types"
)

func TestPair_Validate(t *testing.T) {
	for _, tc := range []struct {
		name        string
		malleate    func(pair *types.Pair)
		expectedErr string
	}{
		{
			"happy case",
			func(pair *types.Pair) {},
			"",
		},
		{
			"zero id",
			func(pair *types.Pair) {
				pair.Id = 0
			},
			"pair id must not be 0",
		},
		{
			"invalid base coin denom",
			func(pair *types.Pair) {
				pair.BaseCoinDenom = "invalliddenom!"
			},
			"invalid base coin denom: invalid denom: invalliddenom!",
		},
		{
			"invalid quote coin denom",
			func(pair *types.Pair) {
				pair.QuoteCoinDenom = "invaliddenom!"
			},
			"invalid quote coin denom: invalid denom: invaliddenom!",
		},
		{
			"invalid escrow address",
			func(pair *types.Pair) {
				pair.EscrowAddress = "invalidaddr"
			},
			"invalid escrow address invalidaddr: decoding bech32 failed: invalid separator index -1",
		},
		{
			"",
			func(pair *types.Pair) {
				p := sdk.NewDec(-1)
				pair.LastPrice = &p
			},
			"last price must be positive: -1.000000000000000000",
		},
		{
			"",
			func(pair *types.Pair) {
				pair.CurrentBatchId = 0
			},
			"current batch id must not be 0",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			pair := types.NewPair(1, "denom1", "denom2", 1)
			tc.malleate(&pair)
			err := pair.Validate()
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}

func TestPairEscrowAddress(t *testing.T) {
	for _, tc := range []struct {
		appID    uint64
		pairId   uint64
		expected string
	}{
		{1, 1, "cosmos1url34vfv5a5a7esm2aapklqelh2mzuwe34vvgc94eh752t8mcjeqla0n0v"},
		{1, 2, "cosmos174fjyc3w25m7pku6sww9w03e9s8phphe4h88euxypsp5ekjngdkqp596l4"},
		{2, 1, "cosmos1mjjgv53lgef35ywsfd4hwduyyptqreev2k9ngh2p84plhrc69srqdmdydt"},
		{2, 2, "cosmos1ljpkf5rve73e77vthnc83m6gxjqznag5st4czyur89sjd2gtclts4rz2wu"},
	} {
		t.Run("", func(t *testing.T) {
			require.Equal(t, tc.expected, types.PairEscrowAddress(tc.appID, tc.pairId).String())
		})
	}
}

func TestSwapFeeCollectorAddress(t *testing.T) {
	for _, tc := range []struct {
		appID    uint64
		pairId   uint64
		expected string
	}{
		{1, 1, "cosmos19a7w3ferywxjst035636dzktx94xyh22u64pwee3fl62sennw5hsw8erx3"},
		{1, 2, "cosmos1sfhtyclc7mvz356z578yc44jqr2pe3tpmcswjngqc3mvh7r2kmcs9pd44e"},
		{2, 1, "cosmos1v9h2ymqyhu34py8shdl59j62gw3h4qykxns76s0epedyq890ca9qvvjw87"},
		{2, 2, "cosmos1cs9rytu8mpsxnpaw05rzh8km03fl4ujjlrwmq9eu8z8wzrqap3zsd0rx2w"},
	} {
		t.Run("", func(t *testing.T) {
			require.Equal(t, tc.expected, types.PairSwapFeeCollectorAddress(tc.appID, tc.pairId).String())
		})
	}
}
