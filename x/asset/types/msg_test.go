package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	addrs = []sdk.AccAddress{
		sdk.AccAddress("test1"),
		sdk.AccAddress("test2"),
	}
)

func TestNewMsgAddMarketRequest(t *testing.T) {
	tests := []struct {
		description string
		from        sdk.AccAddress
		symbol      string
		scriptID    uint64
		expectPass  bool
	}{
		{"add market", addrs[0], "symStr", 1, true},
	}

	for _, tc := range tests {
		m := NewMsgAddMarketRequest(
			tc.from,
			tc.symbol,
			tc.scriptID,
		)
		if tc.expectPass {
			require.NoError(t, m.ValidateBasic(), "test: %v", tc.description)
		} else {
			require.NoError(t, m.ValidateBasic(), "test: %v", tc.description)
		}
	}

}
