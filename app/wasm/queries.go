package wasm

import (
	assetKeeper "github.com/comdex-official/comdex/x/asset/keeper"
	collectorkeeper "github.com/comdex-official/comdex/x/collector/keeper"
	esmKeeper "github.com/comdex-official/comdex/x/esm/keeper"
	liquidationKeeper "github.com/comdex-official/comdex/x/liquidation/keeper"
	lockerkeeper "github.com/comdex-official/comdex/x/locker/keeper"
	rewardsKeeper "github.com/comdex-official/comdex/x/rewards/keeper"
	tokenMintKeeper "github.com/comdex-official/comdex/x/tokenmint/keeper"
	vaultKeeper "github.com/comdex-official/comdex/x/vault/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type QueryPlugin struct {
	assetKeeper       *assetKeeper.Keeper
	lockerKeeper      *lockerkeeper.Keeper
	tokenMintKeeper   *tokenMintKeeper.Keeper
	rewardsKeeper     *rewardsKeeper.Keeper
	collectorKeeper   *collectorkeeper.Keeper
	liquidationKeeper *liquidationKeeper.Keeper
	esmKeeper         *esmKeeper.Keeper
	vaultKeeper       *vaultKeeper.Keeper
}

func NewQueryPlugin(
	assetKeeper *assetKeeper.Keeper,
	lockerKeeper *lockerkeeper.Keeper,
	tokenMintKeeper *tokenMintKeeper.Keeper,
	rewardsKeeper *rewardsKeeper.Keeper,
	collectorKeeper *collectorkeeper.Keeper,
	liquidation *liquidationKeeper.Keeper,
	esmKeeper *esmKeeper.Keeper,
	vaultKeeper *vaultKeeper.Keeper,

) *QueryPlugin {
	return &QueryPlugin{
		assetKeeper:       assetKeeper,
		lockerKeeper:      lockerKeeper,
		tokenMintKeeper:   tokenMintKeeper,
		rewardsKeeper:     rewardsKeeper,
		collectorKeeper:   collectorKeeper,
		liquidationKeeper: liquidation,
		esmKeeper:         esmKeeper,
		vaultKeeper:       vaultKeeper,
	}
}

func (qp QueryPlugin) GetAppInfo(ctx sdk.Context, appID uint64) (int64, int64, uint64, error) {
	MinGovDeposit, GovTimeInSeconds, AssetID, err := qp.assetKeeper.GetAppWasmQuery(ctx, appID)
	if err != nil {
		return MinGovDeposit, GovTimeInSeconds, AssetID, nil
	}
	return MinGovDeposit, GovTimeInSeconds, AssetID, nil
}

func (qp QueryPlugin) GetAssetInfo(ctx sdk.Context, ID uint64) (string, error) {
	assetDenom := qp.assetKeeper.GetAssetDenom(ctx, ID)
	return assetDenom, nil
}

func (qp QueryPlugin) GetTokenMint(ctx sdk.Context, appID, assetID uint64) (int64, error) {
	tokenData, found := qp.tokenMintKeeper.GetAssetDataInTokenMintByAppSupply(ctx, appID, assetID)
	if !found {
		return tokenData, nil
	}
	return tokenData, nil
}

func (qp QueryPlugin) GetRemoveWhitelistAppIDLockerRewardsCheck(ctx sdk.Context, appID uint64, assetIDs uint64) (found bool, err string) {
	found, err = qp.rewardsKeeper.GetRemoveWhitelistAppIDLockerRewardsCheck(ctx, appID, assetIDs)
	return found, err
}

func (qp QueryPlugin) GetWhitelistAppIDVaultInterestCheck(ctx sdk.Context, appID uint64) (found bool, err string) {
	found, err = qp.rewardsKeeper.GetWhitelistAppIDVaultInterestCheck(ctx, appID)
	return found, err
}
func (qp QueryPlugin) GetWhitelistAppIDLockerRewardsCheck(ctx sdk.Context, appID uint64, assetID uint64) (found bool, err string) {
	found, err = qp.rewardsKeeper.GetWhitelistAppIDLockerRewardsCheck(ctx, appID, assetID)

	return found, err
}

func (qp QueryPlugin) GetExternalLockerRewardsCheck(ctx sdk.Context, appID, assetID uint64) (found bool, err string) {
	found, err = qp.rewardsKeeper.GetExternalLockerRewardsCheck(ctx, appID, assetID)
	return found, err
}

func (qp QueryPlugin) GetExternalVaultRewardsCheck(ctx sdk.Context, appID, assetID uint64) (found bool, err string) {
	found, err = qp.rewardsKeeper.GetExternalVaultRewardsCheck(ctx, appID, assetID)
	return found, err
}

func (qp QueryPlugin) CollectorLookupTableQueryCheck(ctx sdk.Context, appID, collectorAssetID, secondaryAssetID uint64) (found bool, err string) {
	found, err = qp.collectorKeeper.WasmSetCollectorLookupTableQuery(ctx, appID, collectorAssetID, secondaryAssetID)
	return found, err
}

func (qp QueryPlugin) ExtendedPairsVaultRecordsQueryCheck(ctx sdk.Context, appID, pairID uint64, StabilityFee, ClosingFee, DrawDownFee sdk.Dec, DebtCeiling, DebtFloor uint64, PairName string) (found bool, err string) {
	found, err = qp.assetKeeper.WasmAddExtendedPairsVaultRecordsQuery(ctx, appID, pairID, StabilityFee, ClosingFee, DrawDownFee, DebtCeiling, DebtFloor, PairName)
	return found, err
}

func (qp QueryPlugin) AuctionMappingForAppQueryCheck(ctx sdk.Context, appID uint64) (found bool, err string) {
	found, err = qp.collectorKeeper.WasmSetAuctionMappingForAppQuery(ctx, appID)
	return found, err
}

func (qp QueryPlugin) WhiteListedAssetQueryCheck(ctx sdk.Context, appID, assetID uint64) (found bool, err string) {
	found, err = qp.lockerKeeper.WasmAddWhiteListedAssetQuery(ctx, appID, assetID)
	return found, err
}

func (qp QueryPlugin) UpdatePairsVaultQueryCheck(ctx sdk.Context, appID, extPairID uint64) (found bool, err string) {
	found, err = qp.assetKeeper.WasmUpdatePairsVaultQuery(ctx, appID, extPairID)
	return found, err
}

func (qp QueryPlugin) UpdateCollectorLookupTableQueryCheck(ctx sdk.Context, appID, assetID uint64) (found bool, err string) {
	found, err = qp.collectorKeeper.WasmUpdateCollectorLookupTableQuery(ctx, appID, assetID)
	return found, err
}

func (qp QueryPlugin) WasmRemoveWhitelistAppIDVaultInterestQueryCheck(ctx sdk.Context, appID uint64) (found bool, err string) {
	found, err = qp.rewardsKeeper.WasmRemoveWhitelistAppIDVaultInterestQuery(ctx, appID)
	return found, err
}

func (qp QueryPlugin) WasmRemoveWhitelistAssetLockerQueryCheck(ctx sdk.Context, appID, assetID uint64) (found bool, err string) {
	found, err = qp.rewardsKeeper.WasmRemoveWhitelistAssetLockerQuery(ctx, appID, assetID)
	return found, err
}

func (qp QueryPlugin) WasmWhitelistAppIDLiquidationQueryCheck(ctx sdk.Context, appID uint64) (found bool, err string) {
	found, err = qp.liquidationKeeper.WasmWhitelistAppIDLiquidationQuery(ctx, appID)
	return found, err
}

func (qp QueryPlugin) WasmRemoveWhitelistAppIDLiquidationQueryCheck(ctx sdk.Context, appID uint64) (found bool, err string) {
	found, err = qp.liquidationKeeper.WasmRemoveWhitelistAppIDLiquidationQuery(ctx, appID)
	return found, err
}

func (qp QueryPlugin) WasmAddESMTriggerParamsQueryCheck(ctx sdk.Context, appID uint64) (found bool, err string) {
	found, err = qp.esmKeeper.WasmAddESMTriggerParamsQuery(ctx, appID)
	return found, err
}

func (qp QueryPlugin) WasmExtendedPairByApp(ctx sdk.Context, appID uint64) (extendedPairIDs []uint64, found bool) {
	extendedPairIDs, found = qp.assetKeeper.WasmExtendedPairByAppQuery(ctx, appID)
	return extendedPairIDs, found
}

func (qp QueryPlugin) WasmCheckSurplusReward(ctx sdk.Context, appID, assetID uint64) (amount sdk.Coin) {
	//TO DO : add extended pair app query
	amount = qp.collectorKeeper.WasmCheckSurplusRewardQuery(ctx, appID, assetID)
	return amount
}

func (qp QueryPlugin) WasmCheckWhitelistedAsset(ctx sdk.Context, denom string) (found bool) {
	//TO DO : add extended pair app query
	found = qp.assetKeeper.WasmCheckWhitelistedAssetQuery(ctx, denom)
	return found
}
