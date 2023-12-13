package v13

import (
	"fmt"
	assetkeeper "github.com/comdex-official/comdex/x/asset/keeper"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	auctionV2keeper "github.com/comdex-official/comdex/x/auctionsV2/keeper"
	auctionsV2types "github.com/comdex-official/comdex/x/auctionsV2/types"
	lendkeeper "github.com/comdex-official/comdex/x/lend/keeper"
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	liquidationV2keeper "github.com/comdex-official/comdex/x/liquidationsV2/keeper"
	liquidationV2types "github.com/comdex-official/comdex/x/liquidationsV2/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	icqkeeper "github.com/cosmos/ibc-apps/modules/async-icq/v7/keeper"
	icqtypes "github.com/cosmos/ibc-apps/modules/async-icq/v7/types"
	exported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	ibctmmigrations "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint/migrations"
)

func CreateUpgradeHandlerV13(
	mm *module.Manager,
	configurator module.Configurator,
	cdc codec.Codec,
	paramsKeeper paramskeeper.Keeper,
	consensusParamsKeeper consensusparamkeeper.Keeper,
	IBCKeeper ibckeeper.Keeper,
	icqkeeper *icqkeeper.Keeper,
	GovKeeper govkeeper.Keeper,
	assetKeeper assetkeeper.Keeper,
	lendKeeper lendkeeper.Keeper,
	liquidationV2Keeper liquidationV2keeper.Keeper,
	auctionV2Keeper auctionV2keeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Applying main net upgrade - v.13.3.0")
		logger := ctx.Logger().With("upgrade", UpgradeName)

		// Migrate Tendermint consensus parameters from x/params module to a deprecated x/consensus module.
		// The old params module is required to still be imported in your app.go in order to handle this migration.
		ctx.Logger().Info("Migrating tendermint consensus params from x/params to x/consensus...")
		legacyParamSubspace := paramsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable())
		baseapp.MigrateParams(ctx, legacyParamSubspace, &consensusParamsKeeper)

		// ibc v4-to-v5
		// https://github.com/cosmos/ibc-go/blob/v7.1.0/docs/migrations/v4-to-v5.md
		// -- nothing --

		// https://github.com/cosmos/ibc-go/blob/v7.1.0/docs/migrations/v5-to-v6.md

		// ibc v6-to-v7
		// https://github.com/cosmos/ibc-go/blob/v7.1.0/docs/migrations/v6-to-v7.md#chains
		// (optional) prune expired tendermint consensus states to save storage space
		ctx.Logger().Info("Pruning expired tendermint consensus states...")
		if _, err := ibctmmigrations.PruneExpiredConsensusStates(ctx, cdc, IBCKeeper.ClientKeeper); err != nil {
			return nil, err
		}

		// ibc v7-to-v7.1
		// https://github.com/cosmos/ibc-go/blob/v7.1.0/docs/migrations/v7-to-v7_1.md#09-localhost-migration
		// explicitly update the IBC 02-client params, adding the localhost client type
		params := IBCKeeper.ClientKeeper.GetParams(ctx)
		params.AllowedClients = append(params.AllowedClients, exported.Localhost)
		IBCKeeper.ClientKeeper.SetParams(ctx, params)
		logger.Info(fmt.Sprintf("updated ibc client params %v", params))

		// icq params set
		icqparams := icqtypes.DefaultParams()
		icqparams.AllowQueries = append(icqparams.AllowQueries, "/cosmwasm.wasm.v1.Query/SmartContractState")
		icqkeeper.SetParams(ctx, icqparams)

		// Run migrations
		logger.Info(fmt.Sprintf("pre migrate version map: %v", fromVM))
		vm, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return nil, err
		}
		logger.Info(fmt.Sprintf("post migrate version map: %v", vm))

		// update gov params to use a 20% initial deposit ratio, allowing us to remote the ante handler
		govParams := GovKeeper.GetParams(ctx)
		govParams.MinInitialDepositRatio = sdk.NewDec(20).Quo(sdk.NewDec(100)).String()
		if err := GovKeeper.SetParams(ctx, govParams); err != nil {
			return nil, err
		}
		logger.Info(fmt.Sprintf("updated gov params to %v", govParams))

		UpdateLendParams(ctx, lendKeeper, assetKeeper)
		InitializeStates(ctx, liquidationV2Keeper, auctionV2Keeper)

		return vm, err
	}
}

func UpdateLendParams(
	ctx sdk.Context,
	lendKeeper lendkeeper.Keeper,
	assetKeeper assetkeeper.Keeper,
) {

	cSTATOM := assettypes.Asset{
		Name:                  "CSTATOM",
		Denom:                 "ucstatom",
		Decimals:              sdk.NewInt(1000000),
		IsOnChain:             true,
		IsOraclePriceRequired: false,
		IsCdpMintable:         true,
	}
	err := assetKeeper.AddAssetRecords(ctx, cSTATOM)
	if err != nil {
		fmt.Println(err)
	}
	assetID := assetKeeper.GetAssetID(ctx)

	assetRatesParamsSTAtom := lendtypes.AssetRatesParams{
		AssetID:              14,
		UOptimal:             newDec("0.75"),
		Base:                 newDec("0.002"),
		Slope1:               newDec("0.07"),
		Slope2:               newDec("1.25"),
		EnableStableBorrow:   false,
		Ltv:                  newDec("0.7"),
		LiquidationThreshold: newDec("0.75"),
		LiquidationPenalty:   newDec("0.05"),
		LiquidationBonus:     newDec("0.05"),
		ReserveFactor:        newDec("0.2"),
		CAssetID:             assetID,
		IsIsolated:           false,
	}
	lendKeeper.SetAssetRatesParams(ctx, assetRatesParamsSTAtom)

	assetRatesParamsCmdx, _ := lendKeeper.GetAssetRatesParams(ctx, 2)
	assetRatesParamsCmdx.LiquidationPenalty = newDec("0.075")
	assetRatesParamsCmdx.LiquidationBonus = newDec("0.075")
	lendKeeper.SetAssetRatesParams(ctx, assetRatesParamsCmdx)

	assetRatesParamsCmst, _ := lendKeeper.GetAssetRatesParams(ctx, 3)
	assetRatesParamsCmst.LiquidationPenalty = newDec("0.05")
	assetRatesParamsCmst.LiquidationBonus = newDec("0.05")
	lendKeeper.SetAssetRatesParams(ctx, assetRatesParamsCmst)

	cAXLUSDC := assettypes.Asset{
		Name:                  "CAXLUSDC",
		Denom:                 "ucaxlusdc",
		Decimals:              sdk.NewInt(1000000),
		IsOnChain:             true,
		IsOraclePriceRequired: false,
		IsCdpMintable:         true,
	}
	err = assetKeeper.AddAssetRecords(ctx, cAXLUSDC)
	if err != nil {
		fmt.Println(err)
	}
}

func InitializeStates(
	ctx sdk.Context,
	liquidationKeeper liquidationV2keeper.Keeper,
	auctionKeeper auctionV2keeper.Keeper,
) {
	dutchAuctionParams := liquidationV2types.DutchAuctionParam{
		Premium:         newDec("1.15"),
		Discount:        newDec("0.7"),
		DecrementFactor: sdk.NewInt(1),
	}
	englishAuctionParams := liquidationV2types.EnglishAuctionParam{DecrementFactor: sdk.NewInt(1)}

	harborParams := liquidationV2types.LiquidationWhiteListing{
		AppId:               2,
		Initiator:           true,
		IsDutchActivated:    true,
		DutchAuctionParam:   &dutchAuctionParams,
		IsEnglishActivated:  true,
		EnglishAuctionParam: &englishAuctionParams,
		KeeeperIncentive:    sdk.ZeroDec(),
	}

	commodoParams := liquidationV2types.LiquidationWhiteListing{
		AppId:               3,
		Initiator:           true,
		IsDutchActivated:    true,
		DutchAuctionParam:   &dutchAuctionParams,
		IsEnglishActivated:  false,
		EnglishAuctionParam: nil,
		KeeeperIncentive:    sdk.ZeroDec(),
	}

	liquidationKeeper.SetLiquidationWhiteListing(ctx, harborParams)
	liquidationKeeper.SetLiquidationWhiteListing(ctx, commodoParams)

	appReserveFundsTxDataHbr, found := liquidationKeeper.GetAppReserveFundsTxData(ctx, 2)
	if !found {
		appReserveFundsTxDataHbr.AppId = 2
	}
	appReserveFundsTxDataHbr.AssetTxData = append(appReserveFundsTxDataHbr.AssetTxData, liquidationV2types.AssetTxData{})
	liquidationKeeper.SetAppReserveFundsTxData(ctx, appReserveFundsTxDataHbr)

	appReserveFundsTxDataCmdo, found := liquidationKeeper.GetAppReserveFundsTxData(ctx, 3)
	if !found {
		appReserveFundsTxDataCmdo.AppId = 3
	}
	appReserveFundsTxDataCmdo.AssetTxData = append(appReserveFundsTxDataCmdo.AssetTxData, liquidationV2types.AssetTxData{})
	liquidationKeeper.SetAppReserveFundsTxData(ctx, appReserveFundsTxDataCmdo)

	auctionParams := auctionsV2types.AuctionParams{
		AuctionDurationSeconds: 18000,
		Step:                   newDec("0.1"),
		WithdrawalFee:          newDec("0.0005"),
		ClosingFee:             newDec("0.0005"),
		MinUsdValueLeft:        100000,
		BidFactor:              newDec("0.01"),
		LiquidationPenalty:     newDec("0.1"),
		AuctionBonus:           newDec("0.0"),
	}
	auctionKeeper.SetAuctionParams(ctx, auctionParams)
	auctionKeeper.SetParams(ctx, auctionsV2types.Params{})
	auctionKeeper.SetAuctionID(ctx, 0)
	auctionKeeper.SetUserBidID(ctx, 0)

}

func newDec(i string) sdk.Dec {
	dec, _ := sdk.NewDecFromStr(i)
	return dec
}
