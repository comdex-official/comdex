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

func TestNewMsgAddMarketRequest_GetSigners(t *testing.T) {
	msgsWithInvalidFrom := []MsgAddMarketRequest{
		{
			From:     "",
			Symbol:   "",
			ScriptID: 0,
		},
		{
			From:     "randomString",
			Symbol:   "",
			ScriptID: 0,
		},
		{
			From:     "5t3y445wiu4",
			Symbol:   "",
			ScriptID: 0,
		},
	}

	msgWithValidFrom := MsgAddMarketRequest{
		From:     "cosmos1cs644d07zvrmcray3uflmn3lwz7gyecyle8vn7",
		Symbol:   "cmdx",
		ScriptID: 1,
	}

	for _, msg := range msgsWithInvalidFrom {
		funcThatPanics := func() {
			msg.GetSigners()
		}
		require.Panics(t, funcThatPanics)
	}

	require.NotPanics(t, func() {
		msgWithValidFrom.GetSigners()
	})
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

func TestMsgUpdateMarketRequest_GetSigners(t *testing.T) {
	msgsWithInvalidFrom := []MsgUpdateMarketRequest{
		{
			From:     "",
			Symbol:   "",
			ScriptID: 0,
		},
		{
			From:     "randomString",
			Symbol:   "",
			ScriptID: 0,
		},
		{
			From:     "5t3y445wiu4",
			Symbol:   "",
			ScriptID: 0,
		},
	}

	msgWithValidFrom := MsgUpdateMarketRequest{
		From:     "cosmos1cs644d07zvrmcray3uflmn3lwz7gyecyle8vn7",
		Symbol:   "cmdx",
		ScriptID: 1,
	}

	for _, msg := range msgsWithInvalidFrom {
		funcThatPanics := func() {
			msg.GetSigners()
		}
		require.Panics(t, funcThatPanics)
	}

	require.NotPanics(t, func() {
		msgWithValidFrom.GetSigners()
	})
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

func TestMsgAddAssetRequest_GetSigners(t *testing.T) {
	msgsWithInvalidFrom := []MsgAddAssetRequest{
		{
			From:     "",
			Name:     "",
			Denom:    "",
			Decimals: 0,
		},
		{
			From:     "randomString",
			Name:     "",
			Denom:    "",
			Decimals: 0,
		},
		{
			From:     "5t3y445wiu4",
			Name:     "",
			Denom:    "",
			Decimals: 0,
		},
	}

	msgWithValidFrom := MsgAddAssetRequest{
		From:     "cosmos1cs644d07zvrmcray3uflmn3lwz7gyecyle8vn7",
		Name:     "cmdx",
		Denom:    "system",
		Decimals: 1,
	}

	for _, msg := range msgsWithInvalidFrom {
		funcThatPanics := func() {
			msg.GetSigners()
		}
		require.Panics(t, funcThatPanics)
	}

	require.NotPanics(t, func() {
		msgWithValidFrom.GetSigners()
	})
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

func TestMsgUpdateAssetRequest_GetSigners(t *testing.T) {
	msgsWithInvalidFrom := []MsgUpdateAssetRequest{
		{
			From:     "",
			Id:       0,
			Name:     "",
			Denom:    "",
			Decimals: 0,
		},
		{
			From:     "randomString",
			Id:       0,
			Name:     "",
			Denom:    "",
			Decimals: 0,
		},
		{
			From:     "5t3y445wiu4",
			Id:       0,
			Name:     "",
			Denom:    "",
			Decimals: 0,
		},
	}

	msgWithValidFrom := MsgUpdateAssetRequest{
		From:     "cosmos1cs644d07zvrmcray3uflmn3lwz7gyecyle8vn7",
		Id:       1,
		Name:     "cmdx",
		Denom:    "system",
		Decimals: 1,
	}

	for _, msg := range msgsWithInvalidFrom {
		funcThatPanics := func() {
			msg.GetSigners()
		}
		require.Panics(t, funcThatPanics)
	}

	require.NotPanics(t, func() {
		msgWithValidFrom.GetSigners()
	})
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

func TestMsgAddMarketForAssetRequest_GetSigners(t *testing.T) {
	msgsWithInvalidFrom := []MsgAddMarketForAssetRequest{
		{
			From:   "",
			Id:     0,
			Symbol: "",
		},
		{
			From:   "randomString",
			Id:     0,
			Symbol: "",
		},
		{
			From:   "5t3y445wiu4",
			Id:     0,
			Symbol: "",
		},
	}

	msgWithValidFrom := MsgAddMarketForAssetRequest{
		From:   "cosmos1cs644d07zvrmcray3uflmn3lwz7gyecyle8vn7",
		Id:     1,
		Symbol: "cmdx",
	}

	for _, msg := range msgsWithInvalidFrom {
		funcThatPanics := func() {
			msg.GetSigners()
		}
		require.Panics(t, funcThatPanics)
	}

	require.NotPanics(t, func() {
		msgWithValidFrom.GetSigners()
	})
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

func TestMsgRemoveMarketForAssetRequest_GetSigners(t *testing.T) {
	msgsWithInvalidFrom := []MsgRemoveMarketForAssetRequest{
		{
			From:   "",
			Id:     0,
			Symbol: "",
		},
		{
			From:   "randomString",
			Id:     0,
			Symbol: "",
		},
		{
			From:   "5t3y445wiu4",
			Id:     0,
			Symbol: "",
		},
	}

	msgWithValidFrom := MsgRemoveMarketForAssetRequest{
		From:   "cosmos1cs644d07zvrmcray3uflmn3lwz7gyecyle8vn7",
		Id:     1,
		Symbol: "cmdx",
	}

	for _, msg := range msgsWithInvalidFrom {
		funcThatPanics := func() {
			msg.GetSigners()
		}
		require.Panics(t, funcThatPanics)
	}

	require.NotPanics(t, func() {
		msgWithValidFrom.GetSigners()
	})
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

func TestMsgAddPairRequest_GetSigners(t *testing.T) {
	msgsWithInvalidFrom := []MsgAddPairRequest{
		{
			From:             "",
			AssetIn:          0,
			AssetOut:         0,
			LiquidationRatio: sdk.NewDec(0),
		},
		{
			From:             "randomString",
			AssetIn:          0,
			AssetOut:         0,
			LiquidationRatio: sdk.NewDec(0),
		},
		{
			From:             "5t3y445wiu4",
			AssetIn:          0,
			AssetOut:         0,
			LiquidationRatio: sdk.NewDec(0),
		},
	}

	msgWithValidFrom := MsgAddPairRequest{
		From:             "cosmos1cs644d07zvrmcray3uflmn3lwz7gyecyle8vn7",
		AssetIn:          1,
		AssetOut:         0,
		LiquidationRatio: sdk.NewDec(0),
	}

	for _, msg := range msgsWithInvalidFrom {
		funcThatPanics := func() {
			msg.GetSigners()
		}
		require.Panics(t, funcThatPanics)
	}

	require.NotPanics(t, func() {
		msgWithValidFrom.GetSigners()
	})
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

func TestMsgUpdatePairRequest_GetSigners(t *testing.T) {
	msgsWithInvalidFrom := []MsgUpdatePairRequest{
		{
			From:             "",
			Id:               0,
			LiquidationRatio: sdk.NewDec(0),
		},
		{
			From:             "randomString",
			Id:               0,
			LiquidationRatio: sdk.NewDec(0),
		},
		{
			From:             "5t3y445wiu4",
			Id:               0,
			LiquidationRatio: sdk.NewDec(0),
		},
	}

	msgWithValidFrom := MsgUpdatePairRequest{
		From:             "cosmos1cs644d07zvrmcray3uflmn3lwz7gyecyle8vn7",
		Id:               0,
		LiquidationRatio: sdk.NewDec(0),
	}

	for _, msg := range msgsWithInvalidFrom {
		funcThatPanics := func() {
			msg.GetSigners()
		}
		require.Panics(t, funcThatPanics)
	}

	require.NotPanics(t, func() {
		msgWithValidFrom.GetSigners()
	})
}

func TestNewMsgFetchPriceRequest(t *testing.T) {
	timeOutHeight := ibcclienttypes.Height{
		1, 1,
	}
	timeOutStamp := []string{"abc"}
	feeLimit := sdk.Coins{sdk.NewCoin("cmdx", sdk.NewInt(1000))}

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
		{"valid_Form", addrs[0], "srsStr", "channelStr", timeOutHeight, 1, timeOutStamp, 1, feeLimit, 1, 1, true},
		{"invalid_form", sdk.AccAddress{}, "srsStr", "channelStr", timeOutHeight, 1, timeOutStamp, 1, feeLimit, 1, 1, false},
		{"invalid_err", addrs[2], "srsStr", "channelStr", timeOutHeight, 1, timeOutStamp, 1, feeLimit, 1, 1, false},
		{"invalid_srsPort", addrs[0], "a&", "channelStr", timeOutHeight, 1, timeOutStamp, 1, feeLimit, 1, 1, false},
		{"invalid_srsChannel", addrs[0], "a1", "5s", timeOutHeight, 1, timeOutStamp, 1, feeLimit, 1, 1, false},
		{"invalid_symbolNil", addrs[0], "a1", "channelStr", timeOutHeight, 1, nil, 1, feeLimit, 1, 1, false},
		{"invalid_symbolEmpty", addrs[0], "a1", "channelStr", timeOutHeight, 1, []string{}, 1, feeLimit, 1, 1, false},
		{"invalid_scriptId", addrs[0], "a1", "channelStr", timeOutHeight, 1, timeOutStamp, 0, feeLimit, 1, 1, false},
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

func TestMsgFetchPriceRequest_GetSigners(t *testing.T) {
	msgsWithInvalidFrom := []MsgFetchPriceRequest{
		{
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
		},
		{
			From:             "randomString",
			SourcePort:       "",
			SourceChannel:    "",
			TimeoutHeight:    ibcclienttypes.Height{},
			TimeoutTimestamp: 0,
			Symbols:          nil,
			ScriptID:         0,
			FeeLimit:         nil,
			PrepareGas:       0,
			ExecuteGas:       0,
		},
		{
			From:             "5t3y445wiu4",
			SourcePort:       "",
			SourceChannel:    "",
			TimeoutHeight:    ibcclienttypes.Height{},
			TimeoutTimestamp: 0,
			Symbols:          nil,
			ScriptID:         0,
			FeeLimit:         nil,
			PrepareGas:       0,
			ExecuteGas:       0,
		},
	}

	msgWithValidFrom := MsgFetchPriceRequest{
		From:             "cosmos1cs644d07zvrmcray3uflmn3lwz7gyecyle8vn7",
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

	for _, msg := range msgsWithInvalidFrom {
		funcThatPanics := func() {
			msg.GetSigners()
		}
		require.Panics(t, funcThatPanics)
	}

	require.NotPanics(t, func() {
		msgWithValidFrom.GetSigners()
	})
}
