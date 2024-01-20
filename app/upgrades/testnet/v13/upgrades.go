package v13

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	bandoraclemodulekeeper "github.com/comdex-official/comdex/x/bandoracle/keeper"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	exported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	ibctmmigrations "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint/migrations"
)

func CreateUpgradeHandlerV13(
	mm *module.Manager,
	configurator module.Configurator,
	cdc codec.Codec,
	capabilityStoreKey *storetypes.KVStoreKey,
	capabilityKeeper *capabilitykeeper.Keeper,
	wasmKeeper wasmkeeper.Keeper,
	paramsKeeper paramskeeper.Keeper,
	consensusParamsKeeper consensusparamkeeper.Keeper,
	IBCKeeper ibckeeper.Keeper,
	GovKeeper govkeeper.Keeper,
	StakingKeeper stakingkeeper.Keeper,
	MintKeeper mintkeeper.Keeper,
	SlashingKeeper slashingkeeper.Keeper,
	bandoracleKeeper bandoraclemodulekeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		sdk.UnwrapSDKContext(ctx).Logger().Info("Applying test net upgrade - v.13.1.0")
		logger := sdk.UnwrapSDKContext(ctx).Logger().With("upgrade", UpgradeName)

		// Migrate Tendermint consensus parameters from x/params module to a deprecated x/consensus module.
		// The old params module is required to still be imported in your app.go in order to handle this migration.
		sdk.UnwrapSDKContext(ctx).Logger().Info("Migrating tendermint consensus params from x/params to x/consensus...")
		legacyParamSubspace := paramsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable())
		baseapp.MigrateParams(sdk.UnwrapSDKContext(ctx), legacyParamSubspace, &consensusParamsKeeper.ParamsStore)

		// ibc v4-to-v5
		// https://github.com/cosmos/ibc-go/blob/v7.1.0/docs/migrations/v4-to-v5.md
		// -- nothing --

		// TODO: check if v5-v6 is required ??
		// https://github.com/cosmos/ibc-go/blob/v7.1.0/docs/migrations/v5-to-v6.md

		// ibc v6-to-v7
		// https://github.com/cosmos/ibc-go/blob/v7.1.0/docs/migrations/v6-to-v7.md#chains
		// (optional) prune expired tendermint consensus states to save storage space
		sdk.UnwrapSDKContext(ctx).Logger().Info("Pruning expired tendermint consensus states...")
		if _, err := ibctmmigrations.PruneExpiredConsensusStates(sdk.UnwrapSDKContext(ctx), cdc, IBCKeeper.ClientKeeper); err != nil {
			return nil, err
		}

		// ibc v7-to-v7.1
		// https://github.com/cosmos/ibc-go/blob/v7.1.0/docs/migrations/v7-to-v7_1.md#09-localhost-migration
		// explicitly update the IBC 02-client params, adding the localhost client type
		params := IBCKeeper.ClientKeeper.GetParams(sdk.UnwrapSDKContext(ctx))
		params.AllowedClients = append(params.AllowedClients, exported.Localhost)
		IBCKeeper.ClientKeeper.SetParams(sdk.UnwrapSDKContext(ctx), params)
		logger.Info(fmt.Sprintf("updated ibc client params %v", params))

		// Run migrations
		logger.Info(fmt.Sprintf("pre migrate version map: %v", fromVM))
		vm, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return nil, err
		}
		logger.Info(fmt.Sprintf("post migrate version map: %v", vm))

		//TODO: confirm the initial deposit
		// update gov params to use a 20% initial deposit ratio, allowing us to remote the ante handler
		govParams, err := GovKeeper.Params.Get(ctx)
		if err != nil {
			return nil, err
		}
		govParams.MinInitialDepositRatio = sdkmath.LegacyNewDec(20).Quo(sdkmath.LegacyNewDec(100)).String()
		if err := GovKeeper.Params.Set(ctx, govParams); err != nil {
			return nil, err
		}
		logger.Info(fmt.Sprintf("updated gov params to %v", govParams))

		// x/Mint
		// Double blocks per year (from 6 seconds to 3 = 2x blocks per year)
		mintParams, err := MintKeeper.Params.Get(ctx)
		if err != nil {
			return nil, err
		}
		mintParams.BlocksPerYear *= 2
		if err = MintKeeper.Params.Set(ctx, mintParams); err != nil {
			return nil, err
		}
		logger.Info(fmt.Sprintf("updated minted blocks per year logic to %v", mintParams))

		// x/Slashing
		// Double slashing window due to double blocks per year
		slashingParams, err := SlashingKeeper.GetParams(ctx)
		if err != nil {
			return nil, err
		}
		slashingParams.SignedBlocksWindow *= 2
		if err := SlashingKeeper.SetParams(ctx, slashingParams); err != nil {
			return nil, err
		}
		logger.Info(fmt.Sprintf("updated slashing params to %v", slashingParams))

		// update wasm to permissionless
		wasmParams := wasmKeeper.GetParams(ctx)
		wasmParams.CodeUploadAccess = wasmtypes.AllowEverybody
		wasmKeeper.SetParams(ctx, wasmParams)
		logger.Info(fmt.Sprintf("updated wasm params to %v", wasmParams))

		// update discard BH of oracle
		bandData := bandoracleKeeper.GetFetchPriceMsg(sdk.UnwrapSDKContext(ctx))
		if bandData.Size() > 0 {
			bandData.AcceptedHeightDiff = 6000
			bandoracleKeeper.SetFetchPriceMsg(sdk.UnwrapSDKContext(ctx), bandData)
			logger.Info(fmt.Sprintf("updated bandData to %v", bandData))
		}
		return vm, err
	}
}
