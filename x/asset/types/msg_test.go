package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	addrs = []sdk.AccAddress{
		sdk.AccAddress("test1"),
		sdk.AccAddress(""),
	}
)

func TestNewMsgAddPairRequest(t *testing.T) {
	tests := []struct {
		description       string
		from              sdk.AccAddress
		assetIn           uint64
		assetOut          uint64
		liquidationRation sdk.Dec
		expectPass        bool
	}{
		{"valid_Form", addrs[0], 1, 20, sdk.NewDec(1), true},
		{"invalid_form", sdk.AccAddress{}, 1, 20, sdk.NewDec(1), false},
		{"invalid_address", addrs[1], 1, 20, sdk.NewDec(1), false},
		{"invalid_assetIn", addrs[0], 0, 20, sdk.NewDec(1), false},
		{"invalid_assetOut", addrs[0], 1, 0, sdk.NewDec(1), false},
		{"invalid_liquidationRatio", addrs[0], 1, 20, sdk.Dec{}, false},
		{"invalid_liquidationNegative", addrs[0], 1, 20, sdk.NewDec(-1), false},
	}

	for _, tc := range tests {
		msg := NewMsgAddPairRequest(
			tc.from,
			tc.assetIn,
			tc.assetOut,
			tc.liquidationRation,
		)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(),msg.GetSigners(), "test: %v", tc.description)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", tc.description)
		}
	}
}

func TestNewMsgUpdatePairRequest(t *testing.T) {
	tests := []struct {
		description       string
		from              sdk.AccAddress
		id                uint64
		liquidationRation sdk.Dec
		expectPass        bool
	}{
		{"valid_from", addrs[0], 1, sdk.NewDec(1), true},
		{"invalid_form", sdk.AccAddress{}, 1, sdk.NewDec(1), false},
		{"invalid_address", addrs[1], 1, sdk.NewDec(1), false},
		{"invalid_liquidationRation", addrs[0], 1, sdk.NewDec(-1), false},
	}

	for _, tc := range tests {
		msg := NewMsgUpdatePairRequest(
			tc.from,
			tc.id,
			tc.liquidationRation,
		)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), msg.GetSigners(),"test: %v", tc.description)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", tc.description)
		}
	}

}

func TestNewMsgAddAssetRequest(t *testing.T) {
	tests := []struct {
		description string
		from        sdk.AccAddress
		name        string
		denom       string
		decimal     int64
		expectPass  bool
	}{
		{"validForm", addrs[0], "name", sdk.DefaultBondDenom, 1, true},
		{"invalid_form", sdk.AccAddress{}, "name", sdk.DefaultBondDenom, 1, false},
		{"invalid_err", addrs[1], "name", sdk.DefaultBondDenom, 1, false},
		{"invalid_Name", addrs[0], "", sdk.DefaultBondDenom, 1, false},
		{"invalidNameLength", addrs[0], "name_length_greater", sdk.DefaultBondDenom, 1, false},
		{"invalidDenom", addrs[0], "name", "", 1, false},
		{"invalidDenomErr", addrs[0], "name", sdk.DefaultCoinDenomRegex(), 1, false},
		{"invalidDecimal", addrs[0], "name", sdk.DefaultBondDenom, -1, false},
	}

	for _, tc := range tests {
		msg := NewMsgAddAssetRequest(
			tc.from,
			tc.name,
			tc.denom,
			tc.decimal,
		)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), msg.GetSigners(),"test: %v", tc.description)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", tc.description)
		}
	}
}

func TestNewMsgUpdateAssetRequest(t *testing.T) {
	tests := []struct {
		description string
		from        sdk.AccAddress
		id          uint64
		name        string
		denom       string
		decimal     int64
		expectPass  bool
	}{
		{"valid_form", addrs[0], 1, "cmdx", "system", 1, true},
		{"invalid_form", sdk.AccAddress{}, 1, "cmdx", "system", 1, false},
		{"invalid_err", addrs[1], 1, "cmdx", "system", 1, false},
		{"invalid_name", addrs[0], 1, "name_length_greater", "system", 1, false},
		{"invalid_denom", addrs[0], 1, "cmdx", "denom_length_greater", 1, false},
	}

	for _, tc := range tests {
		msg := NewMsgUpdateAssetRequest(
			tc.from,
			tc.id,
			tc.name,
			tc.denom,
			tc.decimal,
		)

		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), msg.GetSigners(),"test: %v", tc.description)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", tc.description)
		}
	}

}



