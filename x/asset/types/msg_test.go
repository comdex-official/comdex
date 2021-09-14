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
		{"validForm", addrs[0], "cmdx", 1, true},
		{"invalid formNotEmpty",sdk.AccAddress{},"asdfghjkl",1,false},
		{"invalid formNotNil", sdk.AccAddress{},"asdfghjkl",1,false},
		{"invalid symbol", addrs[0], "", 1, false},
		{"invalid symbol_Length", addrs[0], "asdfghjkl", 1, false},
		{"invalid script_ID", addrs[0], "cmdx", 0, false},
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
			require.Error(t, m.ValidateBasic(), "test: %v", tc.description)
		}
	}
}
