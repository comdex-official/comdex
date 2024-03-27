package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"

	esmkeeper "github.com/comdex-official/comdex/x/esm/keeper"
	vaultkeeper "github.com/comdex-official/comdex/x/vault/keeper"

	auctionkeeper "github.com/comdex-official/comdex/x/auction/keeper"
	liquidationkeeper "github.com/comdex-official/comdex/x/liquidation/keeper"
	tokenmintkeeper "github.com/comdex-official/comdex/x/tokenmint/keeper"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/comdex-official/comdex/app/wasm/bindings"
	assetkeeper "github.com/comdex-official/comdex/x/asset/keeper"
	collectorkeeper "github.com/comdex-official/comdex/x/collector/keeper"
	gaslessKeeper "github.com/comdex-official/comdex/x/gasless/keeper"
	liquidityKeeper "github.com/comdex-official/comdex/x/liquidity/keeper"
	lockerkeeper "github.com/comdex-official/comdex/x/locker/keeper"
	lockertypes "github.com/comdex-official/comdex/x/locker/types"
	rewardskeeper "github.com/comdex-official/comdex/x/rewards/keeper"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	tokenfactorykeeper "github.com/comdex-official/comdex/x/tokenfactory/keeper"
	tokenfactorytypes "github.com/comdex-official/comdex/x/tokenfactory/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func CustomMessageDecorator(lockerKeeper lockerkeeper.Keeper, rewardsKeeper rewardskeeper.Keeper,
	assetKeeper assetkeeper.Keeper, collectorKeeper collectorkeeper.Keeper, liquidationKeeper liquidationkeeper.Keeper,
	auctionKeeper auctionkeeper.Keeper, tokenMintKeeper tokenmintkeeper.Keeper, esmKeeper esmkeeper.Keeper, vaultKeeper vaultkeeper.Keeper, liquiditykeeper liquidityKeeper.Keeper,
	bankkeeper bankkeeper.Keeper, tokenfactorykeeper tokenfactorykeeper.Keeper, gaslesskeeper gaslessKeeper.Keeper,
) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			wrapped:            old,
			lockerKeeper:       lockerKeeper,
			rewardsKeeper:      rewardsKeeper,
			assetKeeper:        assetKeeper,
			collectorKeeper:    collectorKeeper,
			liquidationKeeper:  liquidationKeeper,
			auctionKeeper:      auctionKeeper,
			tokenMintKeeper:    tokenMintKeeper,
			esmKeeper:          esmKeeper,
			vaultKeeper:        vaultKeeper,
			liquiditykeeper:    liquiditykeeper,
			bankkeeper:         bankkeeper,
			tokenfactorykeeper: tokenfactorykeeper,
			gaslesskeeper:      gaslesskeeper,
		}
	}
}

type CustomMessenger struct {
	wrapped            wasmkeeper.Messenger
	lockerKeeper       lockerkeeper.Keeper
	rewardsKeeper      rewardskeeper.Keeper
	assetKeeper        assetkeeper.Keeper
	collectorKeeper    collectorkeeper.Keeper
	liquidationKeeper  liquidationkeeper.Keeper
	auctionKeeper      auctionkeeper.Keeper
	tokenMintKeeper    tokenmintkeeper.Keeper
	esmKeeper          esmkeeper.Keeper
	vaultKeeper        vaultkeeper.Keeper
	liquiditykeeper    liquidityKeeper.Keeper
	bankkeeper         bankkeeper.Keeper
	tokenfactorykeeper tokenfactorykeeper.Keeper
	gaslesskeeper      gaslessKeeper.Keeper
}

var _ wasmkeeper.Messenger = (*CustomMessenger)(nil)

var comdex1 = []string{"comdex17p9rzwnnfxcjp32un9ug7yhhzgtkhvl9jfksztgw5uh69wac2pgs4jg6dx", "comdex1nc5tatafv6eyq7llkr2gv50ff9e22mnf70qgjlv737ktmt4eswrqdfklyz"}
var testnet3 = []string{"comdex1qwlgtx52gsdu7dtp0cekka5zehdl0uj3fhp9acg325fvgs8jdzksjvgq6q", "comdex1ghd753shjuwexxywmgs4xz7x2q732vcnkm6h2pyv9s6ah3hylvrqfy9rd8"}

func (m *CustomMessenger) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) ([]sdk.Event, [][]byte, error) {
	if msg.Custom != nil {
		// only handle the happy path where this is really minting / swapping ...
		// leave everything else for the wrapped version
		var comdexMsg bindings.ComdexMessages
		if err := json.Unmarshal(msg.Custom, &comdexMsg); err != nil {
			return nil, nil, errorsmod.Wrap(err, "comdex msg error")
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
		if comdexMsg.MsgEmissionPoolRewards != nil {
			return m.ExecuteAddEmissionPoolRewards(ctx, contractAddr, comdexMsg.MsgEmissionPoolRewards)
		}

		// only handle the happy path where this is really creating / minting / swapping ...
		// leave everything else for the wrapped version
		var contractMsg bindings.TokenFactoryMsg
		if err := json.Unmarshal(msg.Custom, &contractMsg); err != nil {
			return nil, nil, errorsmod.Wrap(err, "token factory msg")
		}

		if contractMsg.CreateDenom != nil {
			return m.createDenom(ctx, contractAddr, contractMsg.CreateDenom)
		}
		if contractMsg.MintTokens != nil {
			return m.mintTokens(ctx, contractAddr, contractMsg.MintTokens)
		}
		if contractMsg.ChangeAdmin != nil {
			return m.changeAdmin(ctx, contractAddr, contractMsg.ChangeAdmin)
		}
		if contractMsg.BurnTokens != nil {
			return m.burnTokens(ctx, contractAddr, contractMsg.BurnTokens)
		}
		if contractMsg.SetMetadata != nil {
			return m.setMetadata(ctx, contractAddr, contractMsg.SetMetadata)
		}
		if contractMsg.ForceTransfer != nil {
			return m.forceTransfer(ctx, contractAddr, contractMsg.ForceTransfer)
		}
	}
	return m.wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
}

func (m *CustomMessenger) whitelistAssetLocker(ctx sdk.Context, contractAddr sdk.AccAddress, whiteListAsset *bindings.MsgWhiteListAssetLocker) ([]sdk.Event, [][]byte, error) {
	if ctx.ChainID() == "comdex-1" {
		if contractAddr.String() != comdex1[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	} else if ctx.ChainID() == "comdex-test3" {
		if contractAddr.String() != testnet3[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	}
	err := WhiteListAsset(m.lockerKeeper, ctx, contractAddr.String(), whiteListAsset)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "white list asset")
	}
	return nil, nil, nil
}

func WhiteListAsset(lockerKeeper lockerkeeper.Keeper, ctx sdk.Context, contractAddr string,
	whiteListAsset *bindings.MsgWhiteListAssetLocker,
) error {
	msg := lockertypes.MsgAddWhiteListedAssetRequest{
		From:    contractAddr,
		AppId:   whiteListAsset.AppID,
		AssetId: whiteListAsset.AssetID,
	}
	_, err := lockerKeeper.AddWhiteListedAsset(ctx, &msg)
	if err != nil {
		return err
	}

	return nil
}

func (m *CustomMessenger) whitelistAppIDLockerRewards(ctx sdk.Context, contractAddr sdk.AccAddress, whiteListAsset *bindings.MsgWhitelistAppIDLockerRewards) ([]sdk.Event, [][]byte, error) {
	if ctx.ChainID() == "comdex-1" {
		if contractAddr.String() != comdex1[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	} else if ctx.ChainID() == "comdex-test3" {
		if contractAddr.String() != testnet3[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	}
	err := WhitelistAppIDLockerRewards(m.rewardsKeeper, ctx, contractAddr.String(), whiteListAsset)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "white list appId locker rewards")
	}
	return nil, nil, nil
}

func WhitelistAppIDLockerRewards(rewardsKeeper rewardskeeper.Keeper, ctx sdk.Context, contractAddr string,
	whiteListAsset *bindings.MsgWhitelistAppIDLockerRewards,
) error {
	msg := rewardstypes.WhitelistAsset{
		From:         contractAddr,
		AppMappingId: whiteListAsset.AppID,
		AssetId:      whiteListAsset.AssetID,
	}
	_, err := rewardsKeeper.Whitelist(ctx, &msg)
	if err != nil {
		return err
	}

	return nil
}

func (m *CustomMessenger) whitelistAppIDVaultInterest(ctx sdk.Context, contractAddr sdk.AccAddress, whiteListAsset *bindings.MsgWhitelistAppIDVaultInterest) ([]sdk.Event, [][]byte, error) {
	if ctx.ChainID() == "comdex-1" {
		if contractAddr.String() != comdex1[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	} else if ctx.ChainID() == "comdex-test3" {
		if contractAddr.String() != testnet3[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	}
	err := WhitelistAppIDVaultInterest(m.rewardsKeeper, ctx, contractAddr.String(), whiteListAsset)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "white list appId vault Interest")
	}
	return nil, nil, nil
}

func WhitelistAppIDVaultInterest(rewardsKeeper rewardskeeper.Keeper, ctx sdk.Context, contractAddr string,
	whiteListAsset *bindings.MsgWhitelistAppIDVaultInterest,
) error {
	msg := rewardstypes.WhitelistAppIdVault{
		From:         contractAddr,
		AppMappingId: whiteListAsset.AppID,
	}
	_, err := rewardsKeeper.WhitelistAppVault(ctx, &msg)
	if err != nil {
		return err
	}

	return nil
}

func (m *CustomMessenger) AddExtendedPairsVault(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgAddExtendedPairsVault) ([]sdk.Event, [][]byte, error) {
	if ctx.ChainID() == "comdex-1" {
		if contractAddr.String() != comdex1[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	} else if ctx.ChainID() == "comdex-test3" {
		if contractAddr.String() != testnet3[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	}
	err := MsgAddExtendedPairsVault(m.assetKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "AddExtendedPairsVault error")
	}
	return nil, nil, nil
}

func MsgAddExtendedPairsVault(assetKeeper assetkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	pairVaultBinding *bindings.MsgAddExtendedPairsVault,
) error {
	err := assetKeeper.WasmAddExtendedPairsVaultRecords(ctx, pairVaultBinding)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) SetCollectorLookupTable(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgSetCollectorLookupTable) ([]sdk.Event, [][]byte, error) {
	if ctx.ChainID() == "comdex-1" {
		if contractAddr.String() != comdex1[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	} else if ctx.ChainID() == "comdex-test3" {
		if contractAddr.String() != testnet3[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	}
	err := MsgSetCollectorLookupTable(m.collectorKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "SetCollectorLookupTable error")
	}
	return nil, nil, nil
}

func MsgSetCollectorLookupTable(collectorKeeper collectorkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	collectorBindings *bindings.MsgSetCollectorLookupTable,
) error {
	err := collectorKeeper.WasmSetCollectorLookupTable(ctx, collectorBindings)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) SetAuctionMappingForApp(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgSetAuctionMappingForApp) ([]sdk.Event, [][]byte, error) {
	if ctx.ChainID() == "comdex-1" {
		if contractAddr.String() != comdex1[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	} else if ctx.ChainID() == "comdex-test3" {
		if contractAddr.String() != testnet3[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	}
	err := MsgSetAuctionMappingForApp(m.collectorKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "SetAuctionMappingForApp error")
	}
	return nil, nil, nil
}

func MsgSetAuctionMappingForApp(collectorKeeper collectorkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	auctionMappingBinding *bindings.MsgSetAuctionMappingForApp,
) error {
	err := collectorKeeper.WasmSetAuctionMappingForApp(ctx, auctionMappingBinding)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) UpdatePairsVault(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgUpdatePairsVault) ([]sdk.Event, [][]byte, error) {
	if ctx.ChainID() == "comdex-1" {
		if contractAddr.String() != comdex1[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	} else if ctx.ChainID() == "comdex-test3" {
		if contractAddr.String() != testnet3[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	}
	err := MsgUpdatePairsVault(m.assetKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "UpdatePairsVault error")
	}
	return nil, nil, nil
}

func MsgUpdatePairsVault(assetKeeper assetkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	updatePairVault *bindings.MsgUpdatePairsVault,
) error {
	err := assetKeeper.WasmUpdatePairsVault(ctx, updatePairVault)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) UpdateCollectorLookupTable(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgUpdateCollectorLookupTable) ([]sdk.Event, [][]byte, error) {
	if ctx.ChainID() == "comdex-1" {
		if contractAddr.String() != comdex1[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	} else if ctx.ChainID() == "comdex-test3" {
		if contractAddr.String() != testnet3[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	}
	err := MsgUpdateCollectorLookupTable(m.collectorKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "UpdateCollectorLookupTable error")
	}
	return nil, nil, nil
}

func MsgUpdateCollectorLookupTable(collectorKeeper collectorkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	updateColBinding *bindings.MsgUpdateCollectorLookupTable,
) error {
	err := collectorKeeper.WasmUpdateCollectorLookupTable(ctx, updateColBinding)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) RemoveWhitelistAssetLocker(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgRemoveWhitelistAssetLocker) ([]sdk.Event, [][]byte, error) {
	if ctx.ChainID() == "comdex-1" {
		if contractAddr.String() != comdex1[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	} else if ctx.ChainID() == "comdex-test3" {
		if contractAddr.String() != testnet3[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	}
	err := MsgRemoveWhitelistAssetLocker(m.rewardsKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "RemoveWhitelistAssetRewards error")
	}
	return nil, nil, nil
}

func MsgRemoveWhitelistAssetLocker(rewardsKeeper rewardskeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	a *bindings.MsgRemoveWhitelistAssetLocker,
) error {
	err := rewardsKeeper.WasmRemoveWhitelistAssetLocker(ctx, a.AppID, a.AssetID)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) RemoveWhitelistAppIDVaultInterest(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgRemoveWhitelistAppIDVaultInterest) ([]sdk.Event, [][]byte, error) {
	if ctx.ChainID() == "comdex-1" {
		if contractAddr.String() != comdex1[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	} else if ctx.ChainID() == "comdex-test3" {
		if contractAddr.String() != testnet3[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	}
	err := MsgRemoveWhitelistAppIDVaultInterest(m.rewardsKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "RemoveWhitelistAppIdVaultInterest error")
	}
	return nil, nil, nil
}

func MsgRemoveWhitelistAppIDVaultInterest(rewardsKeeper rewardskeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	a *bindings.MsgRemoveWhitelistAppIDVaultInterest,
) error {
	err := rewardsKeeper.WasmRemoveWhitelistAppIDVaultInterest(ctx, a.AppMappingID)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) WhitelistAppIDLiquidation(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgWhitelistAppIDLiquidation) ([]sdk.Event, [][]byte, error) {
	if ctx.ChainID() == "comdex-1" {
		if contractAddr.String() != comdex1[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	} else if ctx.ChainID() == "comdex-test3" {
		if contractAddr.String() != testnet3[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	}
	err := MsgWhitelistAppIDLiquidation(m.liquidationKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "WhitelistAppIdLiquidation error")
	}
	return nil, nil, nil
}

func MsgWhitelistAppIDLiquidation(liquidationKeeper liquidationkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	a *bindings.MsgWhitelistAppIDLiquidation,
) error {
	err := liquidationKeeper.WasmWhitelistAppIDLiquidation(ctx, a.AppID)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) RemoveWhitelistAppIDLiquidation(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgRemoveWhitelistAppIDLiquidation) ([]sdk.Event, [][]byte, error) {
	if ctx.ChainID() == "comdex-1" {
		if contractAddr.String() != comdex1[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	} else if ctx.ChainID() == "comdex-test3" {
		if contractAddr.String() != testnet3[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	}
	err := MsgRemoveWhitelistAppIDLiquidation(m.liquidationKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "RemoveWhitelistAppIdLiquidation error")
	}
	return nil, nil, nil
}

func MsgRemoveWhitelistAppIDLiquidation(liquidationKeeper liquidationkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	a *bindings.MsgRemoveWhitelistAppIDLiquidation,
) error {
	err := liquidationKeeper.WasmRemoveWhitelistAppIDLiquidation(ctx, a.AppID)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) AddAuctionParams(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgAddAuctionParams) ([]sdk.Event, [][]byte, error) {
	if ctx.ChainID() == "comdex-1" {
		if contractAddr.String() != comdex1[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	} else if ctx.ChainID() == "comdex-test3" {
		if contractAddr.String() != testnet3[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	}
	err := MsgAddAuctionParams(m.auctionKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "AddAuctionParams error")
	}
	return nil, nil, nil
}

func MsgAddAuctionParams(auctionKeeper auctionkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	auctionParamsBinding *bindings.MsgAddAuctionParams,
) error {
	err := auctionKeeper.AddAuctionParams(ctx, auctionParamsBinding)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) BurnGovTokensForApp(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgBurnGovTokensForApp) ([]sdk.Event, [][]byte, error) {
	if ctx.ChainID() == "comdex-1" {
		if contractAddr.String() != comdex1[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	} else if ctx.ChainID() == "comdex-test3" {
		if contractAddr.String() != testnet3[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	}
	err := MsgBurnGovTokensForApp(m.tokenMintKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "BurnGovTokensForApp error")
	}
	return nil, nil, nil
}

func MsgBurnGovTokensForApp(tokenMintKeeper tokenmintkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	a *bindings.MsgBurnGovTokensForApp,
) error {
	err := tokenMintKeeper.BurnGovTokensForApp(ctx, a.AppID, a.From, a.Amount)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) AddESMTriggerParams(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgAddESMTriggerParams) ([]sdk.Event, [][]byte, error) {
	if ctx.ChainID() == "comdex-1" {
		if contractAddr.String() != comdex1[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	} else if ctx.ChainID() == "comdex-test3" {
		if contractAddr.String() != testnet3[0] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	}
	err := MsgAddESMTriggerParams(m.esmKeeper, ctx, contractAddr, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "BurnGovTokensForApp error")
	}
	return nil, nil, nil
}

func MsgAddESMTriggerParams(esmKeeper esmkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress,
	a *bindings.MsgAddESMTriggerParams,
) error {
	err := esmKeeper.AddESMTriggerParamsForApp(ctx, a)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) ExecuteAddEmissionRewards(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgEmissionRewards) ([]sdk.Event, [][]byte, error) {
	if ctx.ChainID() == "comdex-1" {
		if contractAddr.String() != comdex1[1] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	} else if ctx.ChainID() == "comdex-test3" {
		if contractAddr.String() != testnet3[1] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	}
	err := MsgAddEmissionRewards(m.vaultKeeper, ctx, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "Emission rewards error")
	}
	return nil, nil, nil
}

func (m *CustomMessenger) ExecuteAddEmissionPoolRewards(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgEmissionPoolRewards) ([]sdk.Event, [][]byte, error) {
	if ctx.ChainID() == "comdex-1" {
		if contractAddr.String() != comdex1[1] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	} else if ctx.ChainID() == "comdex-test3" {
		if contractAddr.String() != testnet3[1] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	}
	err := MsgAddEmissionPoolRewards(m.liquiditykeeper, ctx, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "Emission pool rewards error")
	}
	return nil, nil, nil
}

func MsgAddEmissionRewards(vaultKeeper vaultkeeper.Keeper, ctx sdk.Context,
	a *bindings.MsgEmissionRewards,
) error {
	err := vaultKeeper.WasmMsgAddEmissionRewards(ctx, a.AppID, a.Amount, a.ExtendedPair, a.VotingRatio)
	if err != nil {
		return err
	}
	return nil
}

func MsgAddEmissionPoolRewards(liquiditykeeper liquidityKeeper.Keeper, ctx sdk.Context,
	a *bindings.MsgEmissionPoolRewards,
) error {
	err := liquiditykeeper.WasmMsgAddEmissionPoolRewards(ctx, a.AppID, a.CswapAppID, a.Amount, a.Pools, a.VotingRatio)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) ExecuteFoundationEmission(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgFoundationEmission) ([]sdk.Event, [][]byte, error) {
	if ctx.ChainID() == "comdex-1" {
		if contractAddr.String() != comdex1[1] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	} else if ctx.ChainID() == "comdex-test3" {
		if contractAddr.String() != testnet3[1] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	}
	err := MsgFoundationEmission(m.tokenMintKeeper, ctx, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "Foundation Emission rewards error")
	}
	return nil, nil, nil
}

func MsgFoundationEmission(tokenmintKeeper tokenmintkeeper.Keeper, ctx sdk.Context,
	a *bindings.MsgFoundationEmission,
) error {
	err := tokenmintKeeper.WasmMsgFoundationEmission(ctx, a.AppID, a.Amount, a.FoundationAddress)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) ExecuteMsgRebaseMint(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgRebaseMint) ([]sdk.Event, [][]byte, error) {
	if ctx.ChainID() == "comdex-1" {
		if contractAddr.String() != comdex1[1] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	} else if ctx.ChainID() == "comdex-test3" {
		if contractAddr.String() != testnet3[1] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	}
	err := MsgRebaseMint(m.tokenMintKeeper, ctx, a)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "Foundation Emission rewards error")
	}
	return nil, nil, nil
}

func MsgRebaseMint(tokenmintKeeper tokenmintkeeper.Keeper, ctx sdk.Context,
	a *bindings.MsgRebaseMint,
) error {
	err := tokenmintKeeper.WasmMsgRebaseMint(ctx, a.AppID, a.Amount, a.ContractAddr)
	if err != nil {
		return err
	}
	return nil
}

func (m *CustomMessenger) ExecuteMsgGetSurplusFund(ctx sdk.Context, contractAddr sdk.AccAddress, a *bindings.MsgGetSurplusFund) ([]sdk.Event, [][]byte, error) {
	if ctx.ChainID() == "comdex-1" {
		if contractAddr.String() != comdex1[1] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	} else if ctx.ChainID() == "comdex-test3" {
		if contractAddr.String() != testnet3[1] {
			return nil, nil, sdkerrors.ErrInvalidAddress
		}
	}
	err := MsgGetSurplusFund(m.collectorKeeper, ctx, a, contractAddr)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "Execute surplus fund rewards error")
	}
	return nil, nil, nil
}

func MsgGetSurplusFund(collectorKeeper collectorkeeper.Keeper, ctx sdk.Context,
	a *bindings.MsgGetSurplusFund, contractAddr sdk.AccAddress,
) error {
	err := collectorKeeper.WasmMsgGetSurplusFund(ctx, a.AppID, a.AssetID, contractAddr, a.Amount)
	if err != nil {
		return err
	}
	return nil
}

// createDenom creates a new token denom
func (m *CustomMessenger) createDenom(ctx sdk.Context, contractAddr sdk.AccAddress, createDenom *bindings.CreateDenom) ([]sdk.Event, [][]byte, error) {
	bz, err := PerformCreateDenom(m.tokenfactorykeeper, m.bankkeeper, ctx, contractAddr, createDenom)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform create denom")
	}
	// TODO: double check how this is all encoded to the contract
	return nil, [][]byte{bz}, nil
}

// PerformCreateDenom is used with createDenom to create a token denom; validates the msgCreateDenom.
func PerformCreateDenom(f tokenfactorykeeper.Keeper, b bankkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, createDenom *bindings.CreateDenom) ([]byte, error) {
	if createDenom == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "create denom null create denom"}
	}

	msgServer := tokenfactorykeeper.NewMsgServerImpl(f)

	msgCreateDenom := tokenfactorytypes.NewMsgCreateDenom(contractAddr.String(), createDenom.Subdenom)

	if err := msgCreateDenom.ValidateBasic(); err != nil {
		return nil, errorsmod.Wrap(err, "failed validating MsgCreateDenom")
	}

	// Create denom
	resp, err := msgServer.CreateDenom(
		sdk.WrapSDKContext(ctx),
		msgCreateDenom,
	)
	if err != nil {
		return nil, errorsmod.Wrap(err, "creating denom")
	}

	if createDenom.Metadata != nil {
		newDenom := resp.NewTokenDenom
		err := PerformSetMetadata(f, b, ctx, contractAddr, newDenom, *createDenom.Metadata)
		if err != nil {
			return nil, errorsmod.Wrap(err, "setting metadata")
		}
	}

	return resp.Marshal()
}

// mintTokens mints tokens of a specified denom to an address.
func (m *CustomMessenger) mintTokens(ctx sdk.Context, contractAddr sdk.AccAddress, mint *bindings.MintTokens) ([]sdk.Event, [][]byte, error) {
	err := PerformMint(m.tokenfactorykeeper, m.bankkeeper, ctx, contractAddr, mint)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform mint")
	}
	return nil, nil, nil
}

// PerformMint used with mintTokens to validate the mint message and mint through token factory.
func PerformMint(f tokenfactorykeeper.Keeper, b bankkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, mint *bindings.MintTokens) error {
	if mint == nil {
		return wasmvmtypes.InvalidRequest{Err: "mint token null mint"}
	}
	rcpt, err := parseAddress(mint.MintToAddress)
	if err != nil {
		return err
	}

	coin := sdk.Coin{Denom: mint.Denom, Amount: mint.Amount}
	sdkMsg := tokenfactorytypes.NewMsgMint(contractAddr.String(), coin)

	if err = sdkMsg.ValidateBasic(); err != nil {
		return err
	}

	// Mint through token factory / message server
	msgServer := tokenfactorykeeper.NewMsgServerImpl(f)
	_, err = msgServer.Mint(sdk.WrapSDKContext(ctx), sdkMsg)
	if err != nil {
		return errorsmod.Wrap(err, "minting coins from message")
	}

	if b.BlockedAddr(rcpt) {
		return errorsmod.Wrapf(err, "minting coins to blocked address %s", rcpt.String())
	}

	err = b.SendCoins(ctx, contractAddr, rcpt, sdk.NewCoins(coin))
	if err != nil {
		return errorsmod.Wrap(err, "sending newly minted coins from message")
	}
	return nil
}

// changeAdmin changes the admin.
func (m *CustomMessenger) changeAdmin(ctx sdk.Context, contractAddr sdk.AccAddress, changeAdmin *bindings.ChangeAdmin) ([]sdk.Event, [][]byte, error) {
	err := ChangeAdmin(m.tokenfactorykeeper, ctx, contractAddr, changeAdmin)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "failed to change admin")
	}
	return nil, nil, nil
}

// ChangeAdmin is used with changeAdmin to validate changeAdmin messages and to dispatch.
func ChangeAdmin(f tokenfactorykeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, changeAdmin *bindings.ChangeAdmin) error {
	if changeAdmin == nil {
		return wasmvmtypes.InvalidRequest{Err: "changeAdmin is nil"}
	}
	newAdminAddr, err := parseAddress(changeAdmin.NewAdminAddress)
	if err != nil {
		return err
	}

	changeAdminMsg := tokenfactorytypes.NewMsgChangeAdmin(contractAddr.String(), changeAdmin.Denom, newAdminAddr.String())
	if err := changeAdminMsg.ValidateBasic(); err != nil {
		return err
	}

	msgServer := tokenfactorykeeper.NewMsgServerImpl(f)
	_, err = msgServer.ChangeAdmin(sdk.WrapSDKContext(ctx), changeAdminMsg)
	if err != nil {
		return errorsmod.Wrap(err, "failed changing admin from message")
	}
	return nil
}

// burnTokens burns tokens.
func (m *CustomMessenger) burnTokens(ctx sdk.Context, contractAddr sdk.AccAddress, burn *bindings.BurnTokens) ([]sdk.Event, [][]byte, error) {
	err := PerformBurn(m.tokenfactorykeeper, ctx, contractAddr, burn)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform burn")
	}
	return nil, nil, nil
}

// PerformBurn performs token burning after validating tokenBurn message.
func PerformBurn(f tokenfactorykeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, burn *bindings.BurnTokens) error {
	if burn == nil {
		return wasmvmtypes.InvalidRequest{Err: "burn token null mint"}
	}

	coin := sdk.Coin{Denom: burn.Denom, Amount: burn.Amount}
	sdkMsg := tokenfactorytypes.NewMsgBurn(contractAddr.String(), coin)
	if burn.BurnFromAddress != "" {
		sdkMsg = tokenfactorytypes.NewMsgBurnFrom(contractAddr.String(), coin, burn.BurnFromAddress)
	}

	if err := sdkMsg.ValidateBasic(); err != nil {
		return err
	}

	// Burn through token factory / message server
	msgServer := tokenfactorykeeper.NewMsgServerImpl(f)
	_, err := msgServer.Burn(sdk.WrapSDKContext(ctx), sdkMsg)
	if err != nil {
		return errorsmod.Wrap(err, "burning coins from message")
	}
	return nil
}

// forceTransfer moves tokens.
func (m *CustomMessenger) forceTransfer(ctx sdk.Context, contractAddr sdk.AccAddress, forcetransfer *bindings.ForceTransfer) ([]sdk.Event, [][]byte, error) {
	err := PerformForceTransfer(m.tokenfactorykeeper, ctx, contractAddr, forcetransfer)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform force transfer")
	}
	return nil, nil, nil
}

// PerformForceTransfer performs token moving after validating tokenForceTransfer message.
func PerformForceTransfer(f tokenfactorykeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, forcetransfer *bindings.ForceTransfer) error {
	if forcetransfer == nil {
		return wasmvmtypes.InvalidRequest{Err: "force transfer null"}
	}

	_, err := parseAddress(forcetransfer.FromAddress)
	if err != nil {
		return err
	}

	_, err = parseAddress(forcetransfer.ToAddress)
	if err != nil {
		return err
	}

	coin := sdk.Coin{Denom: forcetransfer.Denom, Amount: forcetransfer.Amount}
	sdkMsg := tokenfactorytypes.NewMsgForceTransfer(contractAddr.String(), coin, forcetransfer.FromAddress, forcetransfer.ToAddress)

	if err := sdkMsg.ValidateBasic(); err != nil {
		return err
	}

	// Transfer through token factory / message server
	msgServer := tokenfactorykeeper.NewMsgServerImpl(f)
	_, err = msgServer.ForceTransfer(sdk.WrapSDKContext(ctx), sdkMsg)
	if err != nil {
		return errorsmod.Wrap(err, "force transferring from message")
	}
	return nil
}

// createDenom creates a new token denom
func (m *CustomMessenger) setMetadata(ctx sdk.Context, contractAddr sdk.AccAddress, setMetadata *bindings.SetMetadata) ([]sdk.Event, [][]byte, error) {
	err := PerformSetMetadata(m.tokenfactorykeeper, m.bankkeeper, ctx, contractAddr, setMetadata.Denom, setMetadata.Metadata)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform create denom")
	}
	return nil, nil, nil
}

// PerformSetMetadata is used with setMetadata to add new metadata
// It also is called inside CreateDenom if optional metadata field is set
func PerformSetMetadata(f tokenfactorykeeper.Keeper, b bankkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, denom string, metadata bindings.Metadata) error {
	// ensure contract address is admin of denom
	auth, err := f.GetAuthorityMetadata(ctx, denom)
	if err != nil {
		return err
	}
	if auth.Admin != contractAddr.String() {
		return wasmvmtypes.InvalidRequest{Err: "only admin can set metadata"}
	}

	// ensure we are setting proper denom metadata (bank uses Base field, fill it if missing)
	if metadata.Base == "" {
		metadata.Base = denom
	} else if metadata.Base != denom {
		// this is the key that we set
		return wasmvmtypes.InvalidRequest{Err: "Base must be the same as denom"}
	}

	// Create and validate the metadata
	bankMetadata := WasmMetadataToSdk(metadata)
	if err := bankMetadata.Validate(); err != nil {
		return err
	}

	b.SetDenomMetaData(ctx, bankMetadata)
	return nil
}

// GetFullDenom is a function, not method, so the message_plugin can use it
func GetFullDenom(contract string, subDenom string) (string, error) {
	// Address validation
	if _, err := parseAddress(contract); err != nil {
		return "", err
	}
	fullDenom, err := tokenfactorytypes.GetTokenDenom(contract, subDenom)
	if err != nil {
		return "", errorsmod.Wrap(err, "validate sub-denom")
	}

	return fullDenom, nil
}

// parseAddress parses address from bech32 string and verifies its format.
func parseAddress(addr string) (sdk.AccAddress, error) {
	parsed, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return nil, errorsmod.Wrap(err, "address from bech32")
	}
	err = sdk.VerifyAddressFormat(parsed)
	if err != nil {
		return nil, errorsmod.Wrap(err, "verify address format")
	}
	return parsed, nil
}

func WasmMetadataToSdk(metadata bindings.Metadata) banktypes.Metadata {
	denoms := []*banktypes.DenomUnit{}
	for _, unit := range metadata.DenomUnits {
		denoms = append(denoms, &banktypes.DenomUnit{
			Denom:    unit.Denom,
			Exponent: unit.Exponent,
			Aliases:  unit.Aliases,
		})
	}
	return banktypes.Metadata{
		Description: metadata.Description,
		Display:     metadata.Display,
		Base:        metadata.Base,
		Name:        metadata.Name,
		Symbol:      metadata.Symbol,
		DenomUnits:  denoms,
	}
}

func SdkMetadataToWasm(metadata banktypes.Metadata) *bindings.Metadata {
	denoms := []bindings.DenomUnit{}
	for _, unit := range metadata.DenomUnits {
		denoms = append(denoms, bindings.DenomUnit{
			Denom:    unit.Denom,
			Exponent: unit.Exponent,
			Aliases:  unit.Aliases,
		})
	}
	return &bindings.Metadata{
		Description: metadata.Description,
		Display:     metadata.Display,
		Base:        metadata.Base,
		Name:        metadata.Name,
		Symbol:      metadata.Symbol,
		DenomUnits:  denoms,
	}
}
