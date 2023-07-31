package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestRewardBandString(t *testing.T) {
	rb := RewardBand{
		SymbolDenom: "ojo",
		RewardBand:  sdk.OneDec(),
	}
	require.Equal(t, rb.String(), "symbol_denom: ojo\nreward_band: \"1.000000000000000000\"\n")

	rbl := RewardBandList{rb}
	require.Equal(t, rbl.String(), "symbol_denom: ojo\nreward_band: \"1.000000000000000000\"")
}

func TestRewardBandEqual(t *testing.T) {
	rb := RewardBand{
		SymbolDenom: "ojo",
		RewardBand:  sdk.OneDec(),
	}
	rb2 := RewardBand{
		SymbolDenom: "ojo",
		RewardBand:  sdk.OneDec(),
	}
	rb3 := RewardBand{
		SymbolDenom: "inequal",
		RewardBand:  sdk.OneDec(),
	}

	require.True(t, rb.Equal(&rb2))
	require.False(t, rb.Equal(&rb3))
	require.False(t, rb2.Equal(&rb3))
}

func TestRewardBandDenomFinder(t *testing.T) {
	rbl := RewardBandList{
		{
			SymbolDenom: "foo",
			RewardBand:  sdk.OneDec(),
		},
		{
			SymbolDenom: "bar",
			RewardBand:  sdk.ZeroDec(),
		},
	}

	band, err := rbl.GetBandFromDenom("foo")
	require.NoError(t, err)
	require.Equal(t, band, sdk.OneDec())

	band, err = rbl.GetBandFromDenom("bar")
	require.NoError(t, err)
	require.Equal(t, band, sdk.ZeroDec())

	band, err = rbl.GetBandFromDenom("baz")
	require.Error(t, err)
	require.Equal(t, band, sdk.ZeroDec())
}
