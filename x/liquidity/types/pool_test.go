package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/petrichormoney/petri/x/liquidity/types"
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
		{8, "ucre8-1", true},
		{9, "denom9-1", true},
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
