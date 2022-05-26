package wasm

import (
	"encoding/json"
	lockertypes "github.com/comdex-official/comdex/x/locker/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// WasmMsgParser - wasm msg parser for staking msgs
type WasmMsgParser struct{}

// NewWasmMsgParser returns bank wasm msg parser
func NewWasmMsgParser() WasmMsgParser {
	return WasmMsgParser{}
}

func (WasmMsgParser) Parse(contractAddr sdk.AccAddress, wasmMsg WhiteListAssetLocker) (sdk.Msg, error) {

	msg := wasmMsg
	cosmosMsg := &lockertypes.MsgAddWhiteListedAssetRequest{
		From:         contractAddr.String(),
		AppMappingId: msg.appMapId,
		AssetID:      msg.assetId,
	}

	return cosmosMsg, cosmosMsg.ValidateBasic()
}

// ParseCustom implements custom parser
func (WasmMsgParser) ParseCustom(_ sdk.AccAddress, _ json.RawMessage) (sdk.Msg, error) {
	return nil, nil
}
