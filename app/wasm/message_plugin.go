package wasm

import (
	"encoding/json"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/comdex-official/comdex/app/wasm/bindings"
	lockerkeeper "github.com/comdex-official/comdex/x/locker/keeper"
	lockertypes "github.com/comdex-official/comdex/x/locker/types"
	rewardskeeper "github.com/comdex-official/comdex/x/rewards/keeper"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func CustomMessageDecorator(lockerKeeper lockerkeeper.Keeper, rewardsKeeper rewardskeeper.Keeper) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			wrapped:       old,
			lockerKeeper:  lockerKeeper,
			rewardsKeeper: rewardsKeeper,
		}
	}
}

type CustomMessenger struct {
	wrapped       wasmkeeper.Messenger
	lockerKeeper  lockerkeeper.Keeper
	rewardsKeeper rewardskeeper.Keeper
}

var _ wasmkeeper.Messenger = (*CustomMessenger)(nil)

func (m *CustomMessenger) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) ([]sdk.Event, [][]byte, error) {
	if msg.Custom != nil {
		// only handle the happy path where this is really minting / swapping ...
		// leave everything else for the wrapped version
		var comdexMsg bindings.ComdexMessages
		if err := json.Unmarshal(msg.Custom, &comdexMsg); err != nil {
			return nil, nil, sdkerrors.Wrap(err, "comdex msg error")
		}
		if &comdexMsg.MsgWhiteListAssetLocker != nil {
			return m.whitelistAssetLocker(ctx, contractAddr, &comdexMsg.MsgWhiteListAssetLocker)
		}
		if &comdexMsg.MsgWhitelistAppIdLockerRewards != nil {
			return m.whitelistAppIdLockerRewards(ctx, contractAddr, &comdexMsg.MsgWhitelistAppIdLockerRewards)
		}
		if &comdexMsg.MsgWhitelistAppIdVaultInterest != nil {
			return m.whitelistAppIdVaultInterest(ctx, contractAddr, &comdexMsg.MsgWhitelistAppIdVaultInterest)
		}
	}
	return m.wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
}

func (m *CustomMessenger) whitelistAssetLocker(ctx sdk.Context, contractAddr sdk.AccAddress, whiteListAsset *bindings.MsgWhiteListAssetLocker) ([]sdk.Event, [][]byte, error) {
	err := WhiteListAsset(m.lockerKeeper, ctx, contractAddr, whiteListAsset)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "white list asset")
	}
	return nil, nil, nil
}

func WhiteListAsset(lockerKeeper lockerkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	whiteListAsset *bindings.MsgWhiteListAssetLocker) error {

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

func (m *CustomMessenger) whitelistAppIdLockerRewards(ctx sdk.Context, contractAddr sdk.AccAddress, whiteListAsset *bindings.MsgWhitelistAppIdLockerRewards) ([]sdk.Event, [][]byte, error) {
	err := WhitelistAppIdLockerRewards(m.rewardsKeeper, ctx, contractAddr, whiteListAsset)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "white list appId locker rewards")
	}
	return nil, nil, nil
}

func WhitelistAppIdLockerRewards(rewardsKeeper rewardskeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	whiteListAsset *bindings.MsgWhitelistAppIdLockerRewards) error {

	msg := rewardstypes.WhitelistAsset{
		From:         contractAddr.String(),
		AppMappingId: whiteListAsset.AppMappingId,
		AssetId:      whiteListAsset.AssetId,
	}
	msgServer := rewardskeeper.NewMsgServerImpl(rewardsKeeper)
	_, err := msgServer.Whitelist(sdk.WrapSDKContext(ctx), &msg)

	if err != nil {
		return err
	}

	return nil

}

func (m *CustomMessenger) whitelistAppIdVaultInterest(ctx sdk.Context, contractAddr sdk.AccAddress, whiteListAsset *bindings.MsgWhitelistAppIdVaultInterest) ([]sdk.Event, [][]byte, error) {
	err := WhitelistAppIdVaultInterest(m.rewardsKeeper, ctx, contractAddr, whiteListAsset)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "white list appId vault Interest")
	}
	return nil, nil, nil
}

func WhitelistAppIdVaultInterest(rewardsKeeper rewardskeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	whiteListAsset *bindings.MsgWhitelistAppIdVaultInterest) error {

	msg := rewardstypes.WhitelistAppIdVault{

		From:         contractAddr.String(),
		AppMappingId: whiteListAsset.AppMappingId,
	}
	msgServer := rewardskeeper.NewMsgServerImpl(rewardsKeeper)
	_, err := msgServer.WhitelistAppVault(sdk.WrapSDKContext(ctx), &msg)

	if err != nil {
		return err
	}

	return nil

}

func GetState(addr, denom, blockheight, target string) (sdk.Coin, error) {
	state, _ := lockerkeeper.QueryState(addr, denom, blockheight, target)
	return *state, nil
}
