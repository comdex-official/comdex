package wasm

import (
	assetKeeper "github.com/comdex-official/comdex/x/asset/keeper"
	lockerkeeper "github.com/comdex-official/comdex/x/locker/keeper"
	rewardsKeeper "github.com/comdex-official/comdex/x/rewards/keeper"
	tokenMintKeeper "github.com/comdex-official/comdex/x/tokenmint/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type QueryPlugin struct {
	assetKeeper     *assetKeeper.Keeper
	lockerKeeper    *lockerkeeper.Keeper
	tokenMintKeeper *tokenMintKeeper.Keeper
	rewardsKeeper   *rewardsKeeper.Keeper
}

func NewQueryPlugin(
	assetKeeper *assetKeeper.Keeper,
	lockerKeeper *lockerkeeper.Keeper,
	tokenMintKeeper *tokenMintKeeper.Keeper,
	rewardsKeeper *rewardsKeeper.Keeper,

) *QueryPlugin {
	return &QueryPlugin{
		assetKeeper:     assetKeeper,
		lockerKeeper:    lockerKeeper,
		tokenMintKeeper: tokenMintKeeper,
		rewardsKeeper:   rewardsKeeper,
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
