package wasm

import (
	"encoding/json"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	wasmbindings "github.com/comdex-official/comdex/app/wasm/bindings"
	lockerkeeper "github.com/comdex-official/comdex/x/locker/keeper"
	lockertypes "github.com/comdex-official/comdex/x/locker/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func CustomMessageDecorator(lockerKeeper lockerkeeper.Keeper) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			wrapped:      old,
			lockerKeeper: lockerKeeper,
		}
	}
}

type CustomMessenger struct {
	wrapped      wasmkeeper.Messenger
	lockerKeeper lockerkeeper.Keeper
}

var _ wasmkeeper.Messenger = (*CustomMessenger)(nil)

func (m *CustomMessenger) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) ([]sdk.Event, [][]byte, error) {
	if msg.Custom != nil {
		// only handle the happy path where this is really minting / swapping ...
		// leave everything else for the wrapped version
		var contractMsg wasmbindings.ComdexMessages
		if err := json.Unmarshal(msg.Custom, &contractMsg); err != nil {
			return nil, nil, sdkerrors.Wrap(err, "comdex msg error")
		}
		if contractMsg.WhiteListAssetLocker != nil {
			return m.whitelistAssetLocker(ctx, contractAddr, contractMsg.WhiteListAssetLocker)
		}
	}
	return m.wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
}

func (m *CustomMessenger) whitelistAssetLocker(ctx sdk.Context, contractAddr sdk.AccAddress, whiteListAsset *wasmbindings.WhiteListAssetLocker) ([]sdk.Event, [][]byte, error) {
	err := WhiteListAsset(m.lockerKeeper, ctx, contractAddr, whiteListAsset)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "white list asset")
	}
	return nil, nil, nil
}

func WhiteListAsset(lockerKeeper lockerkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	whiteListAsset *wasmbindings.WhiteListAssetLocker) error {

	msg := lockertypes.MsgAddWhiteListedAssetRequest{
		From:         contractAddr.String(),
		AppMappingId: whiteListAsset.AppMappingId,
		AssetID:      whiteListAsset.AssetId,
	}
	msgServer := lockerkeeper.NewMsgServiceServer(lockerKeeper)
	_, err := msgServer.MsgAddWhiteListedAsset(sdk.WrapSDKContext(ctx), &msg)

	if err != nil {
		return err
	}

	return nil

}
