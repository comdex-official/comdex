package keeper

import (
	"fmt"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetFeeSource(ctx sdk.Context, sdkTx sdk.Tx, feePayer sdk.AccAddress, fee sdk.Coins) sdk.AccAddress {
	if len(sdkTx.GetMsgs()) > 1 {
		return feePayer
	}

	msg := sdkTx.GetMsgs()[0]
	msgTypeURL := sdk.MsgTypeURL(msg)

	isContract := false
	var contractAddress string

	executeContractMessage, ok := msg.(*wasmtypes.MsgExecuteContract)
	if ok {
		isContract = true
		contractAddress = executeContractMessage.GetContract()
	}

	fmt.Println(msgTypeURL)
	if isContract {
		fmt.Println(contractAddress)
	}

	return feePayer
}
