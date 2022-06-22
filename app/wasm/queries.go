package wasm

import (
	assetKeeper "github.com/comdex-official/comdex/x/asset/keeper"
	collectorkeeper "github.com/comdex-official/comdex/x/collector/keeper"
	liquidationKeeper "github.com/comdex-official/comdex/x/liquidation/keeper"
	lockerkeeper "github.com/comdex-official/comdex/x/locker/keeper"
	rewardsKeeper "github.com/comdex-official/comdex/x/rewards/keeper"
	tokenMintKeeper "github.com/comdex-official/comdex/x/tokenmint/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type QueryPlugin struct {
	assetKeeper       *assetKeeper.Keeper
	lockerKeeper      *lockerkeeper.Keeper
	tokenMintKeeper   *tokenMintKeeper.Keeper
	rewardsKeeper     *rewardsKeeper.Keeper
	collectorKeeper   *collectorkeeper.Keeper
	liquidationKeeper *liquidationKeeper.Keeper
}

func NewQueryPlugin(
	assetKeeper *assetKeeper.Keeper,
	lockerKeeper *lockerkeeper.Keeper,
	tokenMintKeeper *tokenMintKeeper.Keeper,
	rewardsKeeper *rewardsKeeper.Keeper,
	collectorKeeper *collectorkeeper.Keeper,
	liquidation *liquidationKeeper.Keeper,

) *QueryPlugin {
	return &QueryPlugin{
		assetKeeper:       assetKeeper,
		lockerKeeper:      lockerKeeper,
		tokenMintKeeper:   tokenMintKeeper,
		rewardsKeeper:     rewardsKeeper,
		collectorKeeper:   collectorKeeper,
		liquidationKeeper: liquidation,
	}
}

func (qp QueryPlugin) GetAppInfo(ctx sdk.Context, appMappingID uint64) (int64, int64, uint64, error) {
	MinGovDeposit, GovTimeInSeconds, AssetID, err := qp.assetKeeper.GetAppWasmQuery(ctx, appMappingID)
	if err != nil {
		return MinGovDeposit, GovTimeInSeconds, AssetID, nil
	}
	return MinGovDeposit, GovTimeInSeconds, AssetID, nil
}

func (qp QueryPlugin) GetAssetInfo(ctx sdk.Context, ID uint64) (string, error) {
	assetDenom := qp.assetKeeper.GetAssetDenom(ctx, ID)
	return assetDenom, nil
}

func (qp QueryPlugin) GetTokenMint(ctx sdk.Context, appMappingID, assetID uint64) (int64, error) {
	tokenData, found := qp.tokenMintKeeper.GetAssetDataInTokenMintByAppSupply(ctx, appMappingID, assetID)
	if !found {
		return tokenData, nil
	}
	return tokenData, nil
}

func (qp QueryPlugin) GetRemoveWhitelistAppIDLockerRewardsCheck(ctx sdk.Context, appMappingID uint64, assetIDs []uint64) (found bool, err string) {
	found, err = qp.rewardsKeeper.GetRemoveWhitelistAppIDLockerRewardsCheck(ctx, appMappingID, assetIDs)
	return found, err
}

func (qp QueryPlugin) GetWhitelistAppIDVaultInterestCheck(ctx sdk.Context, appMappingID uint64) (found bool, err string) {
	found, err = qp.rewardsKeeper.GetWhitelistAppIDVaultInterestCheck(ctx, appMappingID)
	return found, err
}
func (qp QueryPlugin) GetWhitelistAppIDLockerRewardsCheck(ctx sdk.Context, appMappingID uint64, assetIDs []uint64) (found bool, err string) {
	found, err = qp.rewardsKeeper.GetWhitelistAppIDLockerRewardsCheck(ctx, appMappingID, assetIDs)

	return found, err
}

func (qp QueryPlugin) GetExternalLockerRewardsCheck(ctx sdk.Context, appMappingID, assetID uint64) (found bool, err string) {
	found, err = qp.rewardsKeeper.GetExternalLockerRewardsCheck(ctx, appMappingID, assetID)
	return found, err
}

func (qp QueryPlugin) GetExternalVaultRewardsCheck(ctx sdk.Context, appMappingID, assetID uint64) (found bool, err string) {
	found, err = qp.rewardsKeeper.GetExternalVaultRewardsCheck(ctx, appMappingID, assetID)
	return found, err
}

func (qp QueryPlugin) CollectorLookupTableQueryCheck(ctx sdk.Context, appMappingID, collectorAssetID, secondaryAssetID uint64) (found bool, err string) {
	found, err = qp.collectorKeeper.WasmSetCollectorLookupTableQuery(ctx, appMappingID, collectorAssetID, secondaryAssetID)
	return found, err
}

func (qp QueryPlugin) ExtendedPairsVaultRecordsQueryCheck(ctx sdk.Context, appMappingID, pairID uint64, StabilityFee, ClosingFee, DrawDownFee sdk.Dec, DebtCeiling, DebtFloor uint64, PairName string) (found bool, err string) {
	found, err = qp.assetKeeper.WasmAddExtendedPairsVaultRecordsQuery(ctx, appMappingID, pairID, StabilityFee, ClosingFee, DrawDownFee, DebtCeiling, DebtFloor, PairName)
	return found, err
}

func (qp QueryPlugin) AuctionMappingForAppQueryCheck(ctx sdk.Context, appMappingID uint64) (found bool, err string) {
	found, err = qp.collectorKeeper.WasmSetAuctionMappingForAppQuery(ctx, appMappingID)
	return found, err
}

func (qp QueryPlugin) WhiteListedAssetQueryCheck(ctx sdk.Context, appMappingID, assetID uint64) (found bool, err string) {
	found, err = qp.lockerKeeper.WasmAddWhiteListedAssetQuery(ctx, appMappingID, assetID)
	return found, err
}

func (qp QueryPlugin) UpdatePairsVaultQueryCheck(ctx sdk.Context, appMappingID, extPairID uint64) (found bool, err string) {
	found, err = qp.assetKeeper.WasmUpdatePairsVaultQuery(ctx, appMappingID, extPairID)
	return found, err
}

func (qp QueryPlugin) UpdateCollectorLookupTableQueryCheck(ctx sdk.Context, appMappingID, assetID uint64) (found bool, err string) {
	found, err = qp.collectorKeeper.WasmUpdateCollectorLookupTableQuery(ctx, appMappingID, assetID)
	return found, err
}

func (qp QueryPlugin) WasmRemoveWhitelistAppIDVaultInterestQueryCheck(ctx sdk.Context, appMappingID uint64) (found bool, err string) {
	found, err = qp.rewardsKeeper.WasmRemoveWhitelistAppIDVaultInterestQuery(ctx, appMappingID)
	return found, err
}

func (qp QueryPlugin) WasmRemoveWhitelistAssetLockerQueryCheck(ctx sdk.Context, appMappingID, assetID uint64) (found bool, err string) {
	found, err = qp.rewardsKeeper.WasmRemoveWhitelistAssetLockerQuery(ctx, appMappingID, assetID)
	return found, err
}

func (qp QueryPlugin) WasmWhitelistAppIDLiquidationQueryCheck(ctx sdk.Context, appMappingID uint64) (found bool, err string) {
	found, err = qp.liquidationKeeper.WasmWhitelistAppIdLiquidationQuery(ctx, appMappingID)
	return found, err
}

func (qp QueryPlugin) WasmRemoveWhitelistAppIDLiquidationQueryCheck(ctx sdk.Context, appMappingID uint64) (found bool, err string) {
	found, err = qp.liquidationKeeper.WasmRemoveWhitelistAppIdLiquidationQuery(ctx, appMappingID)
	return found, err
}
