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

func (qp QueryPlugin) GetAppInfo(ctx sdk.Context, appMappingId uint64) (int64, int64, uint64, error) {
	MinGovDeposit, GovTimeInSeconds, AssetId, err := qp.assetKeeper.GetAppWasmQuery(ctx, appMappingId)
	if err != nil {
		return MinGovDeposit, GovTimeInSeconds, AssetId, nil
	}
	return MinGovDeposit, GovTimeInSeconds, AssetId, nil
}

func (qp QueryPlugin) GetAssetInfo(ctx sdk.Context, Id uint64) (string, error) {
	assetDenom := qp.assetKeeper.GetAssetDenom(ctx, Id)
	return assetDenom, nil
}

func (qp QueryPlugin) GetTokenMint(ctx sdk.Context, appMappingId, assetId uint64) (int64, error) {
	tokenData, err := qp.tokenMintKeeper.GetAssetDataInTokenMintByAppSupply(ctx, appMappingId, assetId)
	if err != true {
		return tokenData, nil
	}
	return tokenData, nil
}

func (qp QueryPlugin) GetRemoveWhitelistAppIdLockerRewardsCheck(ctx sdk.Context, appMappingId uint64, assetId []uint64) (found bool, err string) {
	found, err = qp.rewardsKeeper.GetRemoveWhitelistAppIdLockerRewardsCheck(ctx, appMappingId, assetId)
	return found, err
}

func (qp QueryPlugin) GetWhitelistAppIdVaultInterestCheck(ctx sdk.Context, appMappingId uint64) (found bool, err string) {
	found, err = qp.rewardsKeeper.GetWhitelistAppIdVaultInterestCheck(ctx, appMappingId)
	return found, err
}
func (qp QueryPlugin) GetWhitelistAppIdLockerRewardsCheck(ctx sdk.Context, appMappingId uint64, assetId []uint64) (found bool, err string) {
	found, err = qp.rewardsKeeper.GetWhitelistAppIdLockerRewardsCheck(ctx, appMappingId, assetId)

	return found, err
}

func (qp QueryPlugin) GetExternalLockerRewardsCheck(ctx sdk.Context, appMappingId, assetId uint64) (found bool, err string) {
	found, err = qp.rewardsKeeper.GetExternalLockerRewardsCheck(ctx, appMappingId, assetId)
	return found, err
}

func (qp QueryPlugin) GetExternalVaultRewardsCheck(ctx sdk.Context, appMappingId, assetId uint64) (found bool, err string) {
	found, err = qp.rewardsKeeper.GetExternalVaultRewardsCheck(ctx, appMappingId, assetId)
	return found, err
}

func (qp QueryPlugin) CollectorLookupTableQueryCheck(ctx sdk.Context, AppMappingId, CollectorAssetId, SecondaryAssetId uint64) (found bool, err string) {
	found, err = qp.collectorKeeper.WasmSetCollectorLookupTableQuery(ctx, AppMappingId, CollectorAssetId, SecondaryAssetId)
	return found, err
}

func (qp QueryPlugin) ExtendedPairsVaultRecordsQueryCheck(ctx sdk.Context, AppMappingId, PairId uint64, StabilityFee, ClosingFee, DrawDownFee sdk.Dec, DebtCeiling, DebtFloor uint64, PairName string) (found bool, err string) {
	found, err = qp.assetKeeper.WasmAddExtendedPairsVaultRecordsQuery(ctx, AppMappingId, PairId, StabilityFee, ClosingFee, DrawDownFee, DebtCeiling, DebtFloor, PairName)
	return found, err
}

func (qp QueryPlugin) AuctionMappingForAppQueryCheck(ctx sdk.Context, AppMappingId uint64) (found bool, err string) {
	found, err = qp.collectorKeeper.WasmSetAuctionMappingForAppQuery(ctx, AppMappingId)
	return found, err
}

func (qp QueryPlugin) WhiteListedAssetQueryCheck(ctx sdk.Context, AppMappingId, AssetId uint64) (found bool, err string) {
	found, err = qp.lockerKeeper.WasmAddWhiteListedAssetQuery(ctx, AppMappingId, AssetId)
	return found, err
}

func (qp QueryPlugin) UpdateLsrInPairsVaultQueryCheck(ctx sdk.Context, AppMappingId, ExtPairId uint64) (found bool, err string) {
	found, err = qp.assetKeeper.WasmUpdateLsrInPairsVaultQuery(ctx, AppMappingId, ExtPairId)
	return found, err
}

func (qp QueryPlugin) UpdateLsrInCollectorLookupTableQueryCheck(ctx sdk.Context, AppMappingId, AssetId uint64) (found bool, err string) {
	found, err = qp.collectorKeeper.WasmUpdateLsrInCollectorLookupTableQuery(ctx, AppMappingId, AssetId)
	return found, err
}

func (qp QueryPlugin) WasmRemoveWhitelistAppIdVaultInterestQueryCheck(ctx sdk.Context, AppMappingId uint64) (found bool, err string) {
	found, err = qp.rewardsKeeper.WasmRemoveWhitelistAppIdVaultInterestQuery(ctx, AppMappingId)
	return found, err
}

func (qp QueryPlugin) WasmRemoveWhitelistAssetLockerQueryCheck(ctx sdk.Context, AppMappingId, AssetId uint64) (found bool, err string) {
	found, err = qp.rewardsKeeper.WasmRemoveWhitelistAssetLockerQuery(ctx, AppMappingId, AssetId)
	return found, err
}

func (qp QueryPlugin) WasmWhitelistAppIdLiquidationQueryCheck(ctx sdk.Context, AppMappingId uint64) (found bool, err string) {
	found, err = qp.liquidationKeeper.WasmWhitelistAppIdLiquidationQuery(ctx, AppMappingId)
	return found, err
}

func (qp QueryPlugin) WasmRemoveWhitelistAppIdLiquidationQueryCheck(ctx sdk.Context, AppMappingId uint64) (found bool, err string) {
	found, err = qp.liquidationKeeper.WasmRemoveWhitelistAppIdLiquidationQuery(ctx, AppMappingId)
	return found, err
}
