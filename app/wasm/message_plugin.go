package wasm

import (
	"encoding/json"
	liquidationkeeper "github.com/comdex-official/comdex/x/liquidation/keeper"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/comdex-official/comdex/app/wasm/bindings"
	assetkeeper "github.com/comdex-official/comdex/x/asset/keeper"
	collectorkeeper "github.com/comdex-official/comdex/x/collector/keeper"
	lockerkeeper "github.com/comdex-official/comdex/x/locker/keeper"
	lockertypes "github.com/comdex-official/comdex/x/locker/types"
	rewardskeeper "github.com/comdex-official/comdex/x/rewards/keeper"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func CustomMessageDecorator(lockerKeeper lockerkeeper.Keeper, rewardsKeeper rewardskeeper.Keeper, assetKeeper assetkeeper.Keeper, collectorKeeper collectorkeeper.Keeper, liquidationKeeper liquidationkeeper.Keeper) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			wrapped:           old,
			lockerKeeper:      lockerKeeper,
			rewardsKeeper:     rewardsKeeper,
			assetKeeper:       assetKeeper,
			collectorKeeper:   collectorKeeper,
			liquidationKeeper: liquidationKeeper,
		}
	}
}

type CustomMessenger struct {
	wrapped           wasmkeeper.Messenger
	lockerKeeper      lockerkeeper.Keeper
	rewardsKeeper     rewardskeeper.Keeper
	assetKeeper       assetkeeper.Keeper
	collectorKeeper   collectorkeeper.Keeper
	liquidationKeeper liquidationkeeper.Keeper
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
		if comdexMsg.MsgWhiteListAssetLocker != nil {
			return m.whitelistAssetLocker(ctx, contractAddr, comdexMsg.MsgWhiteListAssetLocker)
		}
		if comdexMsg.MsgWhitelistAppIdLockerRewards != nil {
			return m.whitelistAppIdLockerRewards(ctx, contractAddr, comdexMsg.MsgWhitelistAppIdLockerRewards)
		}
		if comdexMsg.MsgWhitelistAppIdVaultInterest != nil {
			return m.whitelistAppIdVaultInterest(ctx, contractAddr, comdexMsg.MsgWhitelistAppIdVaultInterest)
		}
		if comdexMsg.MsgAddExtendedPairsVault != nil {
			return m.AddExtendedPairsVault(ctx, contractAddr, comdexMsg.MsgAddExtendedPairsVault)
		}
		if comdexMsg.MsgSetCollectorLookupTable != nil {
			return m.SetCollectorLookupTable(ctx, contractAddr, comdexMsg.MsgSetCollectorLookupTable)
		}
		if comdexMsg.MsgSetAuctionMappingForApp != nil {
			return m.SetAuctionMappingForApp(ctx, contractAddr, comdexMsg.MsgSetAuctionMappingForApp)
		}
		if comdexMsg.MsgUpdateLsrInPairsVault != nil {
			return m.UpdateLsrInPairsVault(ctx, contractAddr, comdexMsg.MsgUpdateLsrInPairsVault)
		}
		if comdexMsg.MsgUpdateLsrInCollectorLookupTable != nil {
			return m.UpdateLsrInCollectorLookupTable(ctx, contractAddr, comdexMsg.MsgUpdateLsrInCollectorLookupTable)
		}
		if comdexMsg.MsgRemoveWhitelistAssetLocker != nil {
			return m.RemoveWhitelistAssetLocker(ctx, contractAddr, comdexMsg.MsgUpdateLsrInCollectorLookupTable)
		}
		if comdexMsg.MsgRemoveWhitelistAppIdVaultInterest != nil {
			return m.RemoveWhitelistAppIdVaultInterest(ctx, contractAddr, comdexMsg.MsgUpdateLsrInCollectorLookupTable)
		}
		if comdexMsg.MsgWhitelistAppIdLiquidation != nil {
			return m.WhitelistAppIdLiquidation(ctx, contractAddr, comdexMsg.MsgUpdateLsrInCollectorLookupTable)
		}
		if comdexMsg.MsgRemoveWhitelistAppIdLiquidation != nil {
			return m.RemoveWhitelistAppIdLiquidation(ctx, contractAddr, comdexMsg.MsgUpdateLsrInCollectorLookupTable)
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
		AssetId:      whiteListAsset.AssetId,
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

func (m *CustomMessenger) AddExtendedPairsVault(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgAddExtendedPairsVault) ([]sdk.Event, [][]byte, error) {
	err := MsgAddExtendedPairsVault(m.assetKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "AddExtendedPairsVault error")
	}
	return nil, nil, nil
}

func MsgAddExtendedPairsVault(assetKeeper assetkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	a *bindings.MsgAddExtendedPairsVault) error {
	err := assetKeeper.WasmAddExtendedPairsVaultRecords(ctx, a.AppMappingId, a.PairId, a.LiquidationRatio, a.StabilityFee, a.ClosingFee, a.LiquidationPenalty, a.DrawDownFee, a.IsVaultActive, a.DebtCeiling, a.DebtFloor, a.IsPsmPair, a.MinCr, a.PairName, a.AssetOutOraclePrice, a.AssetOutPrice)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) SetCollectorLookupTable(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgSetCollectorLookupTable) ([]sdk.Event, [][]byte, error) {
	err := MsgSetCollectorLookupTable(m.collectorKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "SetCollectorLookupTable error")
	}
	return nil, nil, nil
}

func MsgSetCollectorLookupTable(collectorKeeper collectorkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	a *bindings.MsgSetCollectorLookupTable) error {
	err := collectorKeeper.WasmSetCollectorLookupTable(ctx, a.AppMappingId, a.CollectorAssetId, a.SecondaryAssetId, a.SurplusThreshold, a.DebtThreshold, a.LockerSavingRate, a.LotSize, a.BidFactor)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) SetAuctionMappingForApp(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgSetAuctionMappingForApp) ([]sdk.Event, [][]byte, error) {
	err := MsgSetAuctionMappingForApp(m.collectorKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "SetAuctionMappingForApp error")
	}
	return nil, nil, nil
}

func MsgSetAuctionMappingForApp(collectorKeeper collectorkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	a *bindings.MsgSetAuctionMappingForApp) error {
	err := collectorKeeper.WasmSetAuctionMappingForApp(ctx, a.AppMappingId, a.AssetId, a.IsSurplusAuction, a.IsDebtAuction, a.AssetOutOraclePrice, a.AssetOutPrice)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) UpdateLsrInPairsVault(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgUpdateLsrInPairsVault) ([]sdk.Event, [][]byte, error) {
	err := MsgUpdateLsrInPairsVault(m.assetKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "UpdateLsrInPairsVault error")
	}
	return nil, nil, nil
}

func MsgUpdateLsrInPairsVault(assetKeeper assetkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	a *bindings.MsgUpdateLsrInPairsVault) error {
	err := assetKeeper.WasmUpdateLsrInPairsVault(ctx, a.AppMappingId, a.ExtPairId, a.LiquidationRatio, a.StabilityFee, a.ClosingFee,
		a.LiquidationPenalty, a.DrawDownFee, a.MinCr, a.DebtCeiling, a.DebtFloor)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) UpdateLsrInCollectorLookupTable(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgUpdateLsrInCollectorLookupTable) ([]sdk.Event, [][]byte, error) {
	err := MsgUpdateLsrInCollectorLookupTable(m.collectorKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "UpdateLsrInCollectorLookupTable error")
	}
	return nil, nil, nil
}

func MsgUpdateLsrInCollectorLookupTable(collectorKeeper collectorkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	a *bindings.MsgUpdateLsrInCollectorLookupTable) error {
	err := collectorKeeper.WasmUpdateLsrInCollectorLookupTable(ctx, a.AppMappingId, a.AssetId, a.LSR)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) RemoveWhitelistAssetLocker(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgUpdateLsrInCollectorLookupTable) ([]sdk.Event, [][]byte, error) {
	err := MsgRemoveWhitelistAssetLocker(m.rewardsKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "RemoveWhitelistAssetRewards error")
	}
	return nil, nil, nil
}

func MsgRemoveWhitelistAssetLocker(rewardsKeeper rewardskeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	a *bindings.MsgUpdateLsrInCollectorLookupTable) error {
	err := rewardsKeeper.WasmRemoveWhitelistAssetLocker(ctx, a.AppMappingId, a.AssetId)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) RemoveWhitelistAppIdVaultInterest(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgUpdateLsrInCollectorLookupTable) ([]sdk.Event, [][]byte, error) {
	err := MsgRemoveWhitelistAppIdVaultInterest(m.rewardsKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "RemoveWhitelistAppIdVaultInterest error")
	}
	return nil, nil, nil
}

func MsgRemoveWhitelistAppIdVaultInterest(rewardsKeeper rewardskeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	a *bindings.MsgUpdateLsrInCollectorLookupTable) error {
	err := rewardsKeeper.WasmRemoveWhitelistAppIdVaultInterest(ctx, a.AppMappingId)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) WhitelistAppIdLiquidation(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgUpdateLsrInCollectorLookupTable) ([]sdk.Event, [][]byte, error) {
	err := MsgWhitelistAppIdLiquidation(m.liquidationKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "WhitelistAppIdLiquidation error")
	}
	return nil, nil, nil
}

func MsgWhitelistAppIdLiquidation(liquidationKeeper liquidationkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	a *bindings.MsgUpdateLsrInCollectorLookupTable) error {
	err := liquidationKeeper.WasmWhitelistAppIdLiquidation(ctx, a.AppMappingId)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) RemoveWhitelistAppIdLiquidation(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgUpdateLsrInCollectorLookupTable) ([]sdk.Event, [][]byte, error) {
	err := MsgRemoveWhitelistAppIdLiquidation(m.liquidationKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "RemoveWhitelistAppIdLiquidation error")
	}
	return nil, nil, nil
}

func MsgRemoveWhitelistAppIdLiquidation(liquidationKeeper liquidationkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	a *bindings.MsgUpdateLsrInCollectorLookupTable) error {
	err := liquidationKeeper.WasmRemoveWhitelistAppIdLiquidation(ctx, a.AppMappingId)
	if err != nil {
		return err
	}
	return nil
}
