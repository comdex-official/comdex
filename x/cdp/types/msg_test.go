package types

import (
"testing"

"github.com/stretchr/testify/require"

sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	coinsSingle = sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)
	coinsZero   = sdk.NewCoin(sdk.DefaultBondDenom, sdk.ZeroInt())
	addrs       = []sdk.AccAddress{
		sdk.AccAddress("test1"),
		sdk.AccAddress("test2"),
	}
)

func TestNewMsgCreateCDPRequest(t *testing.T) {
	tests := []struct {
		description    string
		sender         sdk.AccAddress
		collateral     sdk.Coin
		principal      sdk.Coin
		collateralType string
		expectPass     bool
	}{
		{"create cdp", addrs[0], coinsSingle, coinsSingle, "type-a", true},
		{"create cdp no collateral", addrs[0], coinsZero, coinsSingle, "type-a", false},
		{"create cdp no debt", addrs[0], coinsSingle, coinsZero, "type-a", false},
		{"create cdp empty owner", sdk.AccAddress{}, coinsSingle, coinsSingle, "type-a", false},
		{"create cdp empty type", addrs[0], coinsSingle, coinsSingle, "", false},
	}

	for _, tc := range tests {
		msg := NewMsgCreateCDPRequest(
			tc.sender,
			tc.collateral,
			tc.principal,
			tc.collateralType,
		)
		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", tc.description)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", tc.description)
		}
	}
}


func TestNewMsgDrawDebtRequest(t *testing.T) {
	tests := []struct {
		description    string
		sender         sdk.AccAddress
		collateralType string
		principal      sdk.Coin
		expectPass     bool
	}{
		{"draw debt", addrs[0], sdk.DefaultBondDenom, coinsSingle, true},
		{"draw debt no debt", addrs[0], sdk.DefaultBondDenom, coinsZero, false},
		{"draw debt empty owner", sdk.AccAddress{}, sdk.DefaultBondDenom, coinsSingle, false},
		{"draw debt empty denom", sdk.AccAddress{}, "", coinsSingle, false},
	}

	for _, tc := range tests {
		msg := NewMsgDrawDebtRequest(
			tc.sender,
			tc.collateralType,
			tc.principal,
		)
		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", tc.description)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", tc.description)
		}
	}
}

func TestNewMsgRepayDebtRequest(t *testing.T) {
	tests := []struct {
		description string
		sender      sdk.AccAddress
		denom       string
		payment     sdk.Coin
		expectPass  bool
	}{
		{"repay debt", addrs[0], sdk.DefaultBondDenom, coinsSingle, true},
		{"repay debt no payment", addrs[0], sdk.DefaultBondDenom, coinsZero, false},
		{"repay debt empty owner", sdk.AccAddress{}, sdk.DefaultBondDenom, coinsSingle, false},
		{"repay debt empty denom", sdk.AccAddress{}, "", coinsSingle, false},
	}

	for _, tc := range tests {
		msg := NewMsgRepayDebtRequest(
			tc.sender,
			tc.denom,
			tc.payment,
		)
		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", tc.description)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", tc.description)
		}
	}
}
