package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	addrs       = []sdk.AccAddress{
		sdk.AccAddress("test1"),
		sdk.AccAddress(""),
	}
	k sdk.Int
)

func TestNewMsgCreateRequest(t *testing.T) {
	tests := []struct {
		description    string
		From	       sdk.AccAddress
		PairID         uint64
		AmountIn       sdk.Int
		AmountOut      sdk.Int
		expectPass     bool
	}{
		{"create vault", addrs[0], 1, sdk.NewInt(100), sdk.NewInt(100), true},
		{"Invalid Address", addrs[1], 1, sdk.NewInt(100), sdk.NewInt(100), false},
		{"AmountIn IsZero", addrs[0], 1, sdk.NewInt(0), sdk.NewInt(100), false},
		{"AmountIn IsNegative", addrs[0], 1, sdk.NewInt(-1), sdk.NewInt(100), false},
		{"AmountIn IsNil", addrs[0], 1, k, sdk.NewInt(100), false},
		{"AmountOut IsZero", addrs[0], 1, sdk.NewInt(100), sdk.NewInt(0), false},
		{"AmountOut IsNegative", addrs[0], 1, sdk.NewInt(100), sdk.NewInt(-1), false},
		{"AmountOut IsNil", addrs[0], 1, sdk.NewInt(100), k, false},

		}

	for _, tc := range tests {
		msg := NewMsgCreateRequest(
			tc.From,
			tc.PairID,
			tc.AmountIn,
			tc.AmountOut,
		)
		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(),msg.GetSignBytes(), msg.GetSigners(),msg.Route(), msg.Type(), "test: %v", tc.description)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", tc.description)
		}
	}
}

func TestNewMsgDepositRequest(t *testing.T) {
	tests := []struct {
		description    string
		From           sdk.AccAddress
		ID             uint64
		Amount         sdk.Int
		expectPass     bool
	}{
		{"Msg Deposit Request", addrs[0], 1, sdk.NewInt(100), true},
		{"Invalid From", addrs[1], 1, sdk.NewInt(100), false},
		{"Invalid ID", addrs[0], 0, sdk.NewInt(100), false},
		{"Amount IsNil", addrs[0], 1, k, false},
		{"Amount IsNegative", addrs[0], 1, sdk.NewInt(-1), false},
		{"Amount IsZero", addrs[0], 1, sdk.NewInt(0), false},

	}

	for _, tc := range tests {
		msg := NewMsgDepositRequest(
			tc.From,
			tc.ID,
			tc.Amount,
		)
		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(),msg.GetSignBytes(), msg.GetSigners(),msg.Route(), msg.Type(), "test: %v", tc.description)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", tc.description)
		}
	}
}

func TestNewMsgWithdrawRequest(t *testing.T) {
	tests := []struct {
		description    string
		From           sdk.AccAddress
		ID             uint64
		Amount         sdk.Int
		expectPass     bool
	}{
		{"", addrs[0], 1, sdk.NewInt(100), true},
		{"Invalid From", addrs[1], 1, sdk.NewInt(100), false},
		{"Invalid ID", addrs[0], 0, sdk.NewInt(100), false},
		{"Amount IsNil", addrs[0], 1, k, false},
		{"Amount IsNegative", addrs[0], 1, sdk.NewInt(-1), false},
		{"Amount IsZero", addrs[0], 1, sdk.NewInt(0), false},

	}

	for _, tc := range tests {
		msg := NewMsgWithdrawRequest(
			tc.From,
			tc.ID,
			tc.Amount,
		)
		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(),msg.GetSignBytes(), msg.GetSigners(),msg.Route(), msg.Type(), "test: %v", tc.description)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", tc.description)
		}
	}
}

func TestNewMsgDrawRequest(t *testing.T) {
	tests := []struct {
		description    string
		From           sdk.AccAddress
		ID             uint64
		Amount         sdk.Int
		expectPass     bool
	}{
		{"draw debt", addrs[0], 1, sdk.NewInt(100), true},
		{"Invalid From", addrs[1], 1, sdk.NewInt(100), false},
		{"Invalid ID", addrs[0], 0, sdk.NewInt(100), false},
		{"Amount IsNil", addrs[0], 1, k, false},
		{"Amount IsNegative", addrs[0], 1, sdk.NewInt(-1), false},
		{"Amount IsZero", addrs[0], 1, sdk.NewInt(0), false},

	}

	for _, tc := range tests {
		msg := NewMsgDrawRequest(
			tc.From,
			tc.ID,
			tc.Amount,
		)
		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(),msg.GetSignBytes(), msg.GetSigners(),msg.Route(), msg.Type(), "test: %v", tc.description)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", tc.description)
		}
	}
}

func TestNewMsgRepayDebtRequest(t *testing.T) {
	tests := []struct {
		description string
		From           sdk.AccAddress
		ID             uint64
		Amount         sdk.Int
		expectPass  bool
	}{
		{"repay debt", addrs[0], 1, sdk.NewInt(100), true},
		{"Invalid From", addrs[1], 1, sdk.NewInt(100), false},
		{"Invalid ID", addrs[0], 0, sdk.NewInt(100), false},
		{"Amount IsNil", addrs[0], 1, k, false},
		{"Amount IsNegative", addrs[0], 1, sdk.NewInt(-1), false},
		{"Amount IsZero", addrs[0], 1, sdk.NewInt(0), false},

	}

	for _, tc := range tests {
		msg := NewMsgRepayRequest(
			tc.From,
			tc.ID,
			tc.Amount,
		)
		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(),msg.GetSignBytes(), msg.GetSigners(),msg.Route(), msg.Type(), "test: %v", tc.description)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", tc.description)
		}
	}
}

func TestNewMsgLiquidateRequest(t *testing.T) {
	tests := []struct {
		description    string
		From           sdk.AccAddress
		ID             uint64
		expectPass     bool
	}{
		{"Msg Liquidate Request", addrs[0], 1, true},
		{"Msg Liquidate Request", addrs[1], 1, false},
		{"Msg Liquidate Request", addrs[0], 0, false},


	}

	for _, tc := range tests {
		msg := NewMsgLiquidateRequest(
			tc.From,
			tc.ID,
		)
		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(),msg.GetSignBytes(), msg.GetSigners(),msg.Route(), msg.Type(), "test: %v", tc.description)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", tc.description)
		}
	}
}