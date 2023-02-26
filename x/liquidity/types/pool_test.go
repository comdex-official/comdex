package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/comdex-official/comdex/x/liquidity/types"
)

func TestPoolReserveAddress(t *testing.T) {
	for _, tc := range []struct {
		appID    uint64
		poolID   uint64
		expected string
	}{
		{1, 1, "cosmos1s83g2v83mtmc4y4wcf3mw8204flrt0wdlp2mmjnykn2csxcrss5ql5hh0m"},
		{1, 2, "cosmos12g0duka5c5zgjc2yskztfqg8wytnxc34f58evzrf7nrfnnvyjqgq5lx7zk"},
		{2, 1, "cosmos1khz4nd0duzvk4cm3glz3czncnq5ecp77gdh58k558k3wh460rn6qx4e3m0"},
		{2, 2, "cosmos1fl2dgfs62uv9srk8tp7w4xazlxxzwx9fzqhq5c260v5fusj88r2qwcw23l"},
	} {
		t.Run("", func(t *testing.T) {
			require.Equal(t, tc.expected, types.PoolReserveAddress(tc.appID, tc.poolID).String())
		})
	}
}

func TestPoolCoinDenom(t *testing.T) {
	for _, tc := range []struct {
		appID    uint64
		poolId   uint64
		expected string
	}{
		{1, 1, "pool1-1"},
		{1, 10, "pool1-10"},
		{2, 1, "pool2-1"},
		{2, 10, "pool2-10"},
		{6, 18446744073709551615, "pool6-18446744073709551615"},
	} {
		t.Run("", func(t *testing.T) {
			poolCoinDenom := types.PoolCoinDenom(tc.appID, tc.poolId)
			require.Equal(t, tc.expected, poolCoinDenom)
		})
	}
}

func TestParsePoolCoinDenomFailure(t *testing.T) {
	for _, tc := range []struct {
		appID      uint64
		denom      string
		expectsErr bool
	}{
		{1, "pool1-1", false},
		{2, "pool2-10", false},
		{3, "pool3-18446744073709551615", false},
		{4, "pool3-18446744073709551616", true},
		{5, "pool5-01", true},
		{6, "pool6--10", true},
		{7, "pool7-+10", true},
		{8, "stake8-1", true},
		{9, "denom9-1", true},
		{10, "pool10-1", true},
	} {
		t.Run("", func(t *testing.T) {
			_, poolID, err := types.ParsePoolCoinDenom(tc.denom)
			if tc.expectsErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.denom, types.PoolCoinDenom(tc.appID, poolID))
			}
		})
	}
}

func TestFarmCoinDenom(t *testing.T) {
	for _, tc := range []struct {
		appID    uint64
		poolId   uint64
		expected string
	}{
		{1, 1, "farm1-1"},
		{1, 10, "farm1-10"},
		{2, 1, "farm2-1"},
		{2, 10, "farm2-10"},
		{6, 18446744073709551615, "farm6-18446744073709551615"},
	} {
		t.Run("", func(t *testing.T) {
			poolCoinDenom := types.FarmCoinDenom(tc.appID, tc.poolId)
			require.Equal(t, tc.expected, poolCoinDenom)
		})
	}
}

func TestParseFarmCoinDenomFailure(t *testing.T) {
	for _, tc := range []struct {
		appID      uint64
		denom      string
		expectsErr bool
	}{
		{1, "farm1-1", false},
		{2, "farm2-10", false},
		{3, "farm3-18446744073709551615", false},
		{4, "farm3-18446744073709551616", true},
		{5, "farm5-01", true},
		{6, "farm6--10", true},
		{7, "farm7-+10", true},
		{8, "stake8-1", true},
		{9, "denom9-1", true},
		{10, "farm10-1", true},
	} {
		t.Run("", func(t *testing.T) {
			_, poolID, err := types.ParseFarmCoinDenom(tc.denom)
			if tc.expectsErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.denom, types.FarmCoinDenom(tc.appID, poolID))
			}
		})
	}
}

func TestIsFarmCoinDenom(t *testing.T) {
	for _, tc := range []struct {
		denom string
		valid bool
	}{
		{"farm1-1", true},
		{"farm11-1", false},
		{"farm1-1123456", true},
		{"pool1-1123456", false},
		{"stake", false},
		{"cmdx", false},
		{"pool1-1", false},
		{"farmer1-1", false},
		{"farm1-1-1", false},
		{"1farm1-1", false},
		{"farm1farm1-1", false},
		{"farm9-11234", true},
		{"farm2-10", true},
	} {
		t.Run("", func(t *testing.T) {
			isvalid := types.IsFarmCoinDenom(tc.denom)
			require.Equal(t, tc.valid, isvalid)
		})
	}
}
