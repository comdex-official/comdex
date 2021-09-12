package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateP(t *testing.T) {
	invalidPairs := []Pair{
		{Id: 0,AssetIn: 10,AssetOut: 20,LiquidationRatio: sdk.NewDec(1) },
		{Id: 1,AssetIn: 0,AssetOut: 20,LiquidationRatio: sdk.NewDec(1) },
		{Id: 1,AssetIn: 10,AssetOut: 0,LiquidationRatio: sdk.NewDec(1) },
		//{Id: 1,AssetIn: 10,AssetOut: 20,LiquidationRatio: sdk.NewDec(1.0) },
		{Id: 1,AssetIn: 10,AssetOut: 20,LiquidationRatio: sdk.NewDec(-1) },
	}

	validPairs := Pair{
		Id: 1,AssetIn: 1,AssetOut: 20,LiquidationRatio: sdk.NewDec(1),
	}

	for _, pair := range invalidPairs{
			err := pair.Validate()
		require.Error(t, err)
	}

		err := validPairs.Validate()
		require.NoError(t, err)
}
