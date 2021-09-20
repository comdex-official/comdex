package types

import (
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/cosmos-sdk/x/gov/types"
	ibcclienttypes "github.com/cosmos/ibc-go/modules/core/02-client/types"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	addrs = []sdk.AccAddress{
		sdk.AccAddress("test1"),
		sdk.AccAddress(""),
		sdk.AccAddress("zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"),
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
		{"valid_Conditions", addrs[0], "cmdx", 1, true},
		{"invalid_Form", sdk.AccAddress{}, "cmdx", 1, false},
		{"invalid_err", addrs[2], "cmdx", 1, false},
		{"invalid_Symbol", addrs[0], "", 1, false},
		{"invalid_symbol_length", addrs[0], "Symbol_Length", 1, false},
		{"invalid_scriptId", addrs[0], "cmdx", 0, false},
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

func TestNewMsgUpdateMarketRequest(t *testing.T) {
	tests := []struct {
		description string
		from        sdk.AccAddress
		symbol      string
		scriptID    uint64
		expectPass  bool
	}{
		{"validForm", addrs[0], "cmdx", 1, true},
		{"invalid_form", sdk.AccAddress{}, "cmdx", 1, false},
		{"invalid_err", addrs[2], "cmdx", 1, false},
		{"invalidSymbol", addrs[0], "symbol_length", 1, false},
	}

	for _, tc := range tests {
		m := NewMsgUpdateMarketRequest(
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
		{"invalid_err", addrs[2], "name", sdk.DefaultBondDenom, 1, false},
		{"invalid_Name", addrs[0], "", sdk.DefaultBondDenom, 1, false},
		{"invalidNameLength", addrs[0], "name_length_greater", sdk.DefaultBondDenom, 1, false},
		{"invalidDenom", addrs[0], "name", "", 1, false},
		{"invalidDenomErr", addrs[0], "name", sdk.DefaultCoinDenomRegex(), 1, false},
		{"invalidDecimal", addrs[0], "name", sdk.DefaultBondDenom, -1, false},
	}

	for _, tc := range tests {
		m := NewMsgAddAssetRequest(
			tc.from,
			tc.name,
			tc.denom,
			tc.decimal,
		)

		if tc.expectPass {
			require.NoError(t, m.ValidateBasic(), "test: %v", tc.description)
		} else {
			require.Error(t, m.ValidateBasic(), "test: %v", tc.description)
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
		{"invalid_err", addrs[2], 1, "cmdx", "system", 1, false},
		{"invalid_name", addrs[0], 1, "name_length_greater", "system", 1, false},
		{"invalid_denom", addrs[0], 1, "cmdx", "denom_length_greater", 1, false},
	}

	for _, tc := range tests {
		m := NewMsgUpdateAssetRequest(
			tc.from,
			tc.id,
			tc.name,
			tc.denom,
			tc.decimal,
		)

		if tc.expectPass {
			require.NoError(t, m.ValidateBasic(), "test: %v", tc.description)
		} else {
			require.Error(t, m.ValidateBasic(), "test: %v", tc.description)
		}
	}

}

func TestNewMsgAddMarketForAssetRequest(t *testing.T) {
	tests := []struct {
		description string
		from        sdk.AccAddress
		id          uint64
		symbol      string
		expectPass  bool
	}{
		{"valid_form", addrs[0], 1, "cmdx", true},
		{"invalid_Form", sdk.AccAddress{}, 1, "cmdx", false},
		{"invalid_err", addrs[2], 1, "cmdx", false},
		{"invalid_id", addrs[0], 0, "cmdx", false},
		{"invalid_Symbol", addrs[0], 1, "", false},
		{"invalid_symbol+length", addrs[0], 1, "Symbol_length", false},
	}

	for _, tc := range tests {
		m := NewMsgAddMarketForAssetRequest(
			tc.from,
			tc.id,
			tc.symbol,
		)

		if tc.expectPass {
			require.NoError(t, m.ValidateBasic(), "test: %v", tc.description)
		} else {
			require.Error(t, m.ValidateBasic(), "test: %v", tc.description)
		}
	}

}

func TestNewMsgRemoveMarketForAssetRequest(t *testing.T) {
	tests := []struct {
		description string
		from        sdk.AccAddress
		id          uint64
		symbol      string
		expectPass  bool
	}{
		{"valid_form", addrs[0], 1, "cmdx", true},
		{"invalid_Form", sdk.AccAddress{}, 1, "cmdx", false},
		{"invalid_err", addrs[2], 1, "cmdx", false},
		{"invalid_id", addrs[0], 0, "cmdx", false},
		{"invalid_Symbol", addrs[0], 1, "", false},
		{"invalid_symbol+length", addrs[0], 1, "Symbol_length", false},
	}

	for _, tc := range tests {
		m := NewMsgRemoveMarketForAssetRequest(
			tc.from,
			tc.id,
			tc.symbol,
		)

		if tc.expectPass {
			require.NoError(t, m.ValidateBasic(), "test: %v", tc.description)
		} else {
			require.Error(t, m.ValidateBasic(), "test: %v", tc.description)
		}
	}

}

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
		{"invalid_err", addrs[2], 1, 20, sdk.NewDec(1), false},
		{"invalid_assetIn", addrs[0], 0, 20, sdk.NewDec(1), false},
		{"invalid_assetOut", addrs[0], 1, 0, sdk.NewDec(1), false},
		{"invalid_liquidationRatio", addrs[0], 1, 20, sdk.Dec{}, false},
		{"invalid_liquidationNegative", addrs[0], 1, 20, sdk.NewDec(-1), false},
	}

	for _, tc := range tests {
		m := NewMsgAddPairRequest(
			tc.from,
			tc.assetIn,
			tc.assetOut,
			tc.liquidationRation,
		)

		if tc.expectPass {
			require.NoError(t, m.ValidateBasic(), "test: %v", tc.description)
		} else {
			require.Error(t, m.ValidateBasic(), "test: %v", tc.description)
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
		{"invalid_err", addrs[2], 1, sdk.NewDec(1), false},
		{"invalid_liquidationRation", addrs[0], 1, sdk.NewDec(-1), false},
	}

	for _, tc := range tests {
		m := NewMsgUpdatePairRequest(
			tc.from,
			tc.id,
			tc.liquidationRation,
		)

		if tc.expectPass {
			require.NoError(t, m.ValidateBasic(), "test: %v", tc.description)
		} else {
			require.Error(t, m.ValidateBasic(), "test: %v", tc.description)
		}
	}

}

func TestNewMsgFetchPriceRequest(t *testing.T) {
	a := ibcclienttypes.Height{
		1, 1,
	}
	b := []string{"abc"}
	c := sdk.Coins{sdk.NewCoin("cmdx", sdk.NewInt(1000))}

	tests := []struct {
		description      string
		from             sdk.AccAddress
		sourcePort       string
		sourceChannel    string
		timeoutHeight    ibcclienttypes.Height
		timeoutTimestamp uint64
		symbols          []string
		scriptID         uint64
		feeLimit         github_com_cosmos_cosmos_sdk_types.Coins
		prepareGas       uint64
		executeGas       uint64
		expectPass       bool
	}{
		{"valid_Form", addrs[0], "srsStr", "channelStr", a, 1, b, 1, c, 1, 1, true},
		{"invalid_form", sdk.AccAddress{}, "srsStr", "channelStr", a, 1, b, 1, c, 1, 1, false},
		{"invalid_err", addrs[2], "srsStr", "channelStr", a, 1, b, 1, c, 1, 1, false},
		{"invalid_srsPort", addrs[0], "a&", "channelStr", a, 1, b, 1, c, 1, 1, false},
		{"invalid_srsChannel", addrs[0], "a1", "5s", a, 1, b, 1, c, 1, 1, false},
		{"invalid_symbolNil", addrs[0], "a1", "channelStr", a, 1, nil, 1, c, 1, 1, false},
		{"invalid_symbolEmpty", addrs[0], "a1", "channelStr", a, 1, []string{}, 1, c, 1, 1, false},
		{"invalid_scriptId", addrs[0], "a1", "channelStr", a, 1, b, 0, c, 1, 1, false},
	}

	for _, tc := range tests {
		m := NewMsgFetchPriceRequest(
			tc.from,
			tc.sourcePort,
			tc.sourceChannel,
			tc.timeoutHeight,
			tc.timeoutTimestamp,
			tc.symbols,
			tc.scriptID,
			tc.feeLimit,
			tc.prepareGas,
			tc.executeGas,
		)

		if tc.expectPass {
			require.NoError(t, m.ValidateBasic(), "test: %v", tc.description)
		} else {
			require.Error(t, m.ValidateBasic(), "test: %v", tc.description)
		}

	}
}

func TestGet_Signers(t *testing.T) {
	 err := MsgFetchPriceRequest{
		From:             "",
		SourcePort:       "",
		SourceChannel:    "",
		TimeoutHeight:    ibcclienttypes.Height{},
		TimeoutTimestamp: 0,
		Symbols:          nil,
		ScriptID:         0,
		FeeLimit:         nil,
		PrepareGas:       0,
		ExecuteGas:       0,
	}
	err.GetSigners()
	require.Error(t, &err)

}
