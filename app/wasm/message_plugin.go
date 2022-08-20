package wasm

import (
	"encoding/json"
	esmkeeper "github.com/comdex-official/comdex/x/esm/keeper"
	vaultkeeper "github.com/comdex-official/comdex/x/vault/keeper"

	auctionkeeper "github.com/comdex-official/comdex/x/auction/keeper"
	liquidationkeeper "github.com/comdex-official/comdex/x/liquidation/keeper"
	tokenmintkeeper "github.com/comdex-official/comdex/x/tokenmint/keeper"

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

func CustomMessageDecorator(lockerKeeper lockerkeeper.Keeper, rewardsKeeper rewardskeeper.Keeper,
	assetKeeper assetkeeper.Keeper, collectorKeeper collectorkeeper.Keeper, liquidationKeeper liquidationkeeper.Keeper,
	auctionKeeper auctionkeeper.Keeper, tokenMintKeeper tokenmintkeeper.Keeper, esmKeeper esmkeeper.Keeper, vaultKeeper vaultkeeper.Keeper) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			wrapped:           old,
			lockerKeeper:      lockerKeeper,
			rewardsKeeper:     rewardsKeeper,
			assetKeeper:       assetKeeper,
			collectorKeeper:   collectorKeeper,
			liquidationKeeper: liquidationKeeper,
			auctionKeeper:     auctionKeeper,
			tokenMintKeeper:   tokenMintKeeper,
			esmKeeper:         esmKeeper,
			vaultKeeper:       vaultKeeper,
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
	auctionKeeper     auctionkeeper.Keeper
	tokenMintKeeper   tokenmintkeeper.Keeper
	esmKeeper         esmkeeper.Keeper
	vaultKeeper       vaultkeeper.Keeper
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
		if comdexMsg.MsgWhitelistAppIDLockerRewards != nil {
			return m.whitelistAppIDLockerRewards(ctx, contractAddr, comdexMsg.MsgWhitelistAppIDLockerRewards)
		}
		if comdexMsg.MsgWhitelistAppIDVaultInterest != nil {
			return m.whitelistAppIDVaultInterest(ctx, contractAddr, comdexMsg.MsgWhitelistAppIDVaultInterest)
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
		if comdexMsg.MsgUpdatePairsVault != nil {
			return m.UpdatePairsVault(ctx, contractAddr, comdexMsg.MsgUpdatePairsVault)
		}
		if comdexMsg.MsgUpdateCollectorLookupTable != nil {
			return m.UpdateCollectorLookupTable(ctx, contractAddr, comdexMsg.MsgUpdateCollectorLookupTable)
		}
		if comdexMsg.MsgRemoveWhitelistAssetLocker != nil {
			return m.RemoveWhitelistAssetLocker(ctx, contractAddr, comdexMsg.MsgRemoveWhitelistAssetLocker)
		}
		if comdexMsg.MsgRemoveWhitelistAppIDVaultInterest != nil {
			return m.RemoveWhitelistAppIDVaultInterest(ctx, contractAddr, comdexMsg.MsgRemoveWhitelistAppIDVaultInterest)
		}
		if comdexMsg.MsgWhitelistAppIDLiquidation != nil {
			return m.WhitelistAppIDLiquidation(ctx, contractAddr, comdexMsg.MsgWhitelistAppIDLiquidation)
		}
		if comdexMsg.MsgRemoveWhitelistAppIDLiquidation != nil {
			return m.RemoveWhitelistAppIDLiquidation(ctx, contractAddr, comdexMsg.MsgRemoveWhitelistAppIDLiquidation)
		}
		if comdexMsg.MsgAddAuctionParams != nil {
			return m.AddAuctionParams(ctx, contractAddr, comdexMsg.MsgAddAuctionParams)
		}
		if comdexMsg.MsgBurnGovTokensForApp != nil {
			return m.BurnGovTokensForApp(ctx, contractAddr, comdexMsg.MsgBurnGovTokensForApp)
		}
		if comdexMsg.MsgAddESMTriggerParams != nil {
			return m.AddESMTriggerParams(ctx, contractAddr, comdexMsg.MsgAddESMTriggerParams)
		}
		if comdexMsg.MsgEmissionRewards != nil {
			return m.ExecuteAddEmissionRewards(ctx, contractAddr, comdexMsg.MsgEmissionRewards)
		}
		if comdexMsg.MsgFoundationEmission != nil {
			return m.ExecuteFoundationEmission(ctx, contractAddr, comdexMsg.MsgFoundationEmission)
		}
		if comdexMsg.MsgRebaseMint != nil {
			return m.ExecuteMsgRebaseMint(ctx, contractAddr, comdexMsg.MsgRebaseMint)
		}
		if comdexMsg.MsgGetSurplusFund != nil {
			return m.ExecuteMsgGetSurplusFund(ctx, contractAddr, comdexMsg.MsgGetSurplusFund)
		}
	}
	return m.wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
}

func (m *CustomMessenger) whitelistAssetLocker(ctx sdk.Context, contractAddr sdk.AccAddress, whiteListAsset *bindings.MsgWhiteListAssetLocker) ([]sdk.Event, [][]byte, error) {
	err := WhiteListAsset(m.lockerKeeper, ctx, contractAddr.String(), whiteListAsset)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "white list asset")
	}
	return nil, nil, nil
}

func WhiteListAsset(lockerKeeper lockerkeeper.Keeper, ctx sdk.Context, contractAddr string,
	whiteListAsset *bindings.MsgWhiteListAssetLocker) error {
	msg := lockertypes.MsgAddWhiteListedAssetRequest{
		From:    contractAddr,
		AppId:   whiteListAsset.AppID,
		AssetId: whiteListAsset.AssetID,
	}
	_, err := lockerKeeper.AddWhiteListedAsset(sdk.WrapSDKContext(ctx), &msg)

	if err != nil {
		return err
	}

	return nil
}

func (m *CustomMessenger) whitelistAppIDLockerRewards(ctx sdk.Context, contractAddr sdk.AccAddress, whiteListAsset *bindings.MsgWhitelistAppIDLockerRewards) ([]sdk.Event, [][]byte, error) {
	err := WhitelistAppIDLockerRewards(m.rewardsKeeper, ctx, contractAddr.String(), whiteListAsset)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "white list appId locker rewards")
	}
	return nil, nil, nil
}

func WhitelistAppIDLockerRewards(rewardsKeeper rewardskeeper.Keeper, ctx sdk.Context, contractAddr string,
	whiteListAsset *bindings.MsgWhitelistAppIDLockerRewards) error {
	msg := rewardstypes.WhitelistAsset{
		From:         contractAddr,
		AppMappingId: whiteListAsset.AppID,
		AssetId:      whiteListAsset.AssetIDs,
	}
	_, err := rewardsKeeper.Whitelist(sdk.WrapSDKContext(ctx), &msg)

	if err != nil {
		return err
	}

	return nil
}

func (m *CustomMessenger) whitelistAppIDVaultInterest(ctx sdk.Context, contractAddr sdk.AccAddress, whiteListAsset *bindings.MsgWhitelistAppIDVaultInterest) ([]sdk.Event, [][]byte, error) {
	err := WhitelistAppIDVaultInterest(m.rewardsKeeper, ctx, contractAddr.String(), whiteListAsset)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "white list appId vault Interest")
	}
	return nil, nil, nil
}

func WhitelistAppIDVaultInterest(rewardsKeeper rewardskeeper.Keeper, ctx sdk.Context, contractAddr string,
	whiteListAsset *bindings.MsgWhitelistAppIDVaultInterest) error {
	msg := rewardstypes.WhitelistAppIdVault{

		From:         contractAddr,
		AppMappingId: whiteListAsset.AppID,
	}
	_, err := rewardsKeeper.WhitelistAppVault(sdk.WrapSDKContext(ctx), &msg)

	if err != nil {
		return err
	}

	return nil
}

func GetState(addr, denom, blockHeight, target string) (sdk.Coin, error) {
	state, _ := lockerkeeper.QueryState(addr, denom, blockHeight, target)
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
	pairVaultBinding *bindings.MsgAddExtendedPairsVault) error {
	err := assetKeeper.WasmAddExtendedPairsVaultRecords(ctx, pairVaultBinding)
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
	collectorBindings *bindings.MsgSetCollectorLookupTable) error {
	err := collectorKeeper.WasmSetCollectorLookupTable(ctx, collectorBindings)
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
	auctionMappingBinding *bindings.MsgSetAuctionMappingForApp) error {
	err := collectorKeeper.WasmSetAuctionMappingForApp(ctx, auctionMappingBinding)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) UpdatePairsVault(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgUpdatePairsVault) ([]sdk.Event, [][]byte, error) {
	err := MsgUpdatePairsVault(m.assetKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "UpdatePairsVault error")
	}
	return nil, nil, nil
}

func MsgUpdatePairsVault(assetKeeper assetkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	updatePairVault *bindings.MsgUpdatePairsVault) error {
	err := assetKeeper.WasmUpdatePairsVault(ctx, updatePairVault)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) UpdateCollectorLookupTable(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgUpdateCollectorLookupTable) ([]sdk.Event, [][]byte, error) {
	err := MsgUpdateCollectorLookupTable(m.collectorKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "UpdateCollectorLookupTable error")
	}
	return nil, nil, nil
}

func MsgUpdateCollectorLookupTable(collectorKeeper collectorkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	updateColBinding *bindings.MsgUpdateCollectorLookupTable) error {
	err := collectorKeeper.WasmUpdateCollectorLookupTable(ctx, updateColBinding)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) RemoveWhitelistAssetLocker(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgRemoveWhitelistAssetLocker) ([]sdk.Event, [][]byte, error) {
	err := MsgRemoveWhitelistAssetLocker(m.rewardsKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "RemoveWhitelistAssetRewards error")
	}
	return nil, nil, nil
}

func MsgRemoveWhitelistAssetLocker(rewardsKeeper rewardskeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	a *bindings.MsgRemoveWhitelistAssetLocker) error {
	err := rewardsKeeper.WasmRemoveWhitelistAssetLocker(ctx, a.AppID, a.AssetID)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) RemoveWhitelistAppIDVaultInterest(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgRemoveWhitelistAppIDVaultInterest) ([]sdk.Event, [][]byte, error) {
	err := MsgRemoveWhitelistAppIDVaultInterest(m.rewardsKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "RemoveWhitelistAppIdVaultInterest error")
	}
	return nil, nil, nil
}

func MsgRemoveWhitelistAppIDVaultInterest(rewardsKeeper rewardskeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	a *bindings.MsgRemoveWhitelistAppIDVaultInterest) error {
	err := rewardsKeeper.WasmRemoveWhitelistAppIDVaultInterest(ctx, a.AppMappingID)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) WhitelistAppIDLiquidation(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgWhitelistAppIDLiquidation) ([]sdk.Event, [][]byte, error) {
	err := MsgWhitelistAppIDLiquidation(m.liquidationKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "WhitelistAppIdLiquidation error")
	}
	return nil, nil, nil
}

func MsgWhitelistAppIDLiquidation(liquidationKeeper liquidationkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	a *bindings.MsgWhitelistAppIDLiquidation) error {
	err := liquidationKeeper.WasmWhitelistAppIDLiquidation(ctx, a.AppID)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) RemoveWhitelistAppIDLiquidation(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgRemoveWhitelistAppIDLiquidation) ([]sdk.Event, [][]byte, error) {
	err := MsgRemoveWhitelistAppIDLiquidation(m.liquidationKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "RemoveWhitelistAppIdLiquidation error")
	}
	return nil, nil, nil
}

func MsgRemoveWhitelistAppIDLiquidation(liquidationKeeper liquidationkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	a *bindings.MsgRemoveWhitelistAppIDLiquidation) error {
	err := liquidationKeeper.WasmRemoveWhitelistAppIDLiquidation(ctx, a.AppID)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) AddAuctionParams(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgAddAuctionParams) ([]sdk.Event, [][]byte, error) {
	err := MsgAddAuctionParams(m.auctionKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "AddAuctionParams error")
	}
	return nil, nil, nil
}

func MsgAddAuctionParams(auctionKeeper auctionkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	auctionParamsBinding *bindings.MsgAddAuctionParams) error {
	err := auctionKeeper.AddAuctionParams(ctx, auctionParamsBinding)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) BurnGovTokensForApp(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgBurnGovTokensForApp) ([]sdk.Event, [][]byte, error) {
	err := MsgBurnGovTokensForApp(m.tokenMintKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "BurnGovTokensForApp error")
	}
	return nil, nil, nil
}

func MsgBurnGovTokensForApp(tokenMintKeeper tokenmintkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	a *bindings.MsgBurnGovTokensForApp) error {
	err := tokenMintKeeper.BurnGovTokensForApp(ctx, a.AppID, a.From, a.Amount)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) AddESMTriggerParams(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgAddESMTriggerParams) ([]sdk.Event, [][]byte, error) {
	err := MsgAddESMTriggerParams(m.esmKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "BurnGovTokensForApp error")
	}
	return nil, nil, nil
}

func MsgAddESMTriggerParams(esmKeeper esmkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	a *bindings.MsgAddESMTriggerParams) error {
	err := esmKeeper.AddESMTriggerParamsForApp(ctx, a)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) ExecuteAddEmissionRewards(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgEmissionRewards) ([]sdk.Event, [][]byte, error) {
	err := MsgAddEmissionRewards(m.vaultKeeper, ctx, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "Emission rewards error")
	}
	return nil, nil, nil
}

func MsgAddEmissionRewards(vaultKeeper vaultkeeper.Keeper, ctx sdk.Context,
	a *bindings.MsgEmissionRewards) error {
	err := vaultKeeper.WasmMsgAddEmissionRewards(ctx, a.AppID, a.Amount, a.ExtendedPair, a.VotingRatio)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) ExecuteFoundationEmission(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgFoundationEmission) ([]sdk.Event, [][]byte, error) {
	err := MsgFoundationEmission(m.tokenMintKeeper, ctx, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "Foundation Emission rewards error")
	}
	return nil, nil, nil
}

func MsgFoundationEmission(tokenmintKeeper tokenmintkeeper.Keeper, ctx sdk.Context,
	a *bindings.MsgFoundationEmission) error {
	err := tokenmintKeeper.WasmMsgFoundationEmission(ctx, a.AppID, a.Amount, a.FoundationAddress)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) ExecuteMsgRebaseMint(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgRebaseMint) ([]sdk.Event, [][]byte, error) {
	err := MsgRebaseMint(m.tokenMintKeeper, ctx, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "Foundation Emission rewards error")
	}
	return nil, nil, nil
}

func MsgRebaseMint(tokenmintKeeper tokenmintkeeper.Keeper, ctx sdk.Context,
	a *bindings.MsgRebaseMint) error {
	err := tokenmintKeeper.WasmMsgRebaseMint(ctx, a.AppID, a.Amount, a.ContractAddr)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) ExecuteMsgGetSurplusFund(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgGetSurplusFund) ([]sdk.Event, [][]byte, error) {
	err := MsgGetSurplusFund(m.collectorKeeper, ctx, a, contractAddr)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "Execute surplus fund rewards error")
	}
	return nil, nil, nil
}

func MsgGetSurplusFund(collectorKeeper collectorkeeper.Keeper, ctx sdk.Context,
	a *bindings.MsgGetSurplusFund, contractAddr sdk.AccAddress) error {
	err := collectorKeeper.WasmMsgGetSurplusFund(ctx, a.AppID, a.AssetID, contractAddr, a.Amount)
	if err != nil {
		return err
	}
	return nil
}
