package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestVault(t *testing.T) {
	var amount sdk.Int
	tests := []struct {
		description    string
		ID             uint64
		PairID         uint64
		Owner          string
		AmountIn       sdk.Int
		AmountOut      sdk.Int
		expectPass     bool
	}{
		{"create Vault", 1, 1, "cosmos1yples84d8avjlmegn90663mmjs4tardwjmltrm", sdk.NewInt(100),sdk.NewInt(100), true },
		{"Empty ID", 0, 1, "cosmos1yples84d8avjlmegn90663mmjs4tardwjmltrm", sdk.NewInt(100),sdk.NewInt(100), false },
		{"Empty PairID", 1, 0, "cosmos1yples84d8avjlmegn90663mmjs4tardwjmltrm", sdk.NewInt(100),sdk.NewInt(100), false },
		{"Empty Owner", 1, 1, "", sdk.NewInt(100),sdk.NewInt(100), false },
		{"Invalid Owner", 1, 1, "sentinelles84d8avjlmegn90663mmjs4tardwjmltrm", sdk.NewInt(100),sdk.NewInt(100), false },
		{"AmountIn IsNil", 1, 1, "cosmos1yples84d8avjlmegn90663mmjs4tardwjmltrm", amount,sdk.NewInt(100), false },
		{"AmountIn IsNegative", 1, 1, "cosmos1yples84d8avjlmegn90663mmjs4tardwjmltrm", sdk.NewInt(-1),sdk.NewInt(100), false },
		{"AmountOut IsNil", 1, 1, "cosmos1yples84d8avjlmegn90663mmjs4tardwjmltrm", sdk.NewInt(100),amount, false },
		{"AmountOut IsNegative", 1, 1, "cosmos1yples84d8avjlmegn90663mmjs4tardwjmltrm", sdk.NewInt(100),sdk.NewInt(-1), false },

	}

	for _, tc := range tests {
		msg := Vault{
			tc.ID,
			tc.PairID,
			tc.Owner,
			tc.AmountIn,
			tc.AmountOut,
		}
		if tc.expectPass {
			require.NoError(t, msg.Validate() , "test: %v", tc.description)
		} else {
			require.Error(t, msg.Validate(), "test: %v", tc.description)
		}
	}
}
