package decmath

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMedian(t *testing.T) {
	require := require.New(t)
	prices := []sdk.Dec{
		sdk.MustNewDecFromStr("1.12"),
		sdk.MustNewDecFromStr("1.07"),
		sdk.MustNewDecFromStr("1.11"),
		sdk.MustNewDecFromStr("1.2"),
	}

	median, err := Median(prices)
	require.NoError(err)
	require.Equal(sdk.MustNewDecFromStr("1.115"), median)

	// test empty prices list
	median, err = Median([]sdk.Dec{})
	require.ErrorIs(err, ErrEmptyList)
}

func TestMedianDeviation(t *testing.T) {
	require := require.New(t)
	prices := []sdk.Dec{
		sdk.MustNewDecFromStr("1.12"),
		sdk.MustNewDecFromStr("1.07"),
		sdk.MustNewDecFromStr("1.11"),
		sdk.MustNewDecFromStr("1.2"),
	}
	median := sdk.MustNewDecFromStr("1.115")

	medianDeviation, err := MedianDeviation(median, prices)
	require.NoError(err)
	require.Equal(sdk.MustNewDecFromStr("0.048218253804964775"), medianDeviation)

	// test empty prices list
	medianDeviation, err = MedianDeviation(median, []sdk.Dec{})
	require.ErrorIs(err, ErrEmptyList)
}

func TestAverage(t *testing.T) {
	require := require.New(t)
	prices := []sdk.Dec{
		sdk.MustNewDecFromStr("1.12"),
		sdk.MustNewDecFromStr("1.07"),
		sdk.MustNewDecFromStr("1.11"),
		sdk.MustNewDecFromStr("1.2"),
	}

	average, err := Average(prices)
	require.NoError(err)
	require.Equal(sdk.MustNewDecFromStr("1.125"), average)

	// test empty prices list
	average, err = Average([]sdk.Dec{})
	require.ErrorIs(err, ErrEmptyList)
}

func TestMin(t *testing.T) {
	require := require.New(t)
	prices := []sdk.Dec{
		sdk.MustNewDecFromStr("1.12"),
		sdk.MustNewDecFromStr("1.07"),
		sdk.MustNewDecFromStr("1.11"),
		sdk.MustNewDecFromStr("1.2"),
	}

	min, err := Min(prices)
	require.NoError(err)
	require.Equal(sdk.MustNewDecFromStr("1.07"), min)

	// test empty prices list
	min, err = Min([]sdk.Dec{})
	require.ErrorIs(err, ErrEmptyList)
}

func TestMax(t *testing.T) {
	require := require.New(t)
	prices := []sdk.Dec{
		sdk.MustNewDecFromStr("1.12"),
		sdk.MustNewDecFromStr("1.07"),
		sdk.MustNewDecFromStr("1.11"),
		sdk.MustNewDecFromStr("1.2"),
	}

	max, err := Max(prices)
	require.NoError(err)
	require.Equal(sdk.MustNewDecFromStr("1.2"), max)

	// test empty prices list
	max, err = Max([]sdk.Dec{})
	require.ErrorIs(err, ErrEmptyList)
}

func TestNewDecFromFloat(t *testing.T) {
	testCases := []struct {
		name       string
		float      float64
		dec        sdk.Dec
		expectPass bool
	}{
		{
			name:       "max float64 precision",
			float:      1.000_000_000_000_001,
			dec:        sdk.MustNewDecFromStr("1.000000000000001"),
			expectPass: true,
		},
		{
			name:       "over max float64 precision",
			float:      1.000_000_000_000_000_1,
			dec:        sdk.MustNewDecFromStr("1"),
			expectPass: true,
		},
		{
			name:       "simple float",
			float:      2999999.9,
			dec:        sdk.MustNewDecFromStr("2999999.9"),
			expectPass: true,
		},
		{
			name:       "negative float",
			float:      -10.598,
			dec:        sdk.MustNewDecFromStr("-10.598"),
			expectPass: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dec, err := NewDecFromFloat(tc.float)
			if tc.expectPass {
				require.NoError(t, err)
				require.Equal(t, tc.dec, dec)
			} else {
				require.Error(t, err)
			}
		})
	}
}
