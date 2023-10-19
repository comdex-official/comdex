package v13

import (
	"fmt"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	exported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	ibctmmigrations "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint/migrations"
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
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Applying test net upgrade - v.13.0.0")
		logger := ctx.Logger().With("upgrade", UpgradeName)

		// Migrate Tendermint consensus parameters from x/params module to a deprecated x/consensus module.
		// The old params module is required to still be imported in your app.go in order to handle this migration.
		ctx.Logger().Info("Migrating tendermint consensus params from x/params to x/consensus...")
		legacyParamSubspace := paramsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable())
		baseapp.MigrateParams(ctx, legacyParamSubspace, &consensusParamsKeeper)

		// ibc v4-to-v5
		// https://github.com/cosmos/ibc-go/blob/v7.1.0/docs/migrations/v4-to-v5.md
		// -- nothing --

		// TODO: check if v5-v6 is required ??
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

		// Run migrations
		logger.Info(fmt.Sprintf("pre migrate version map: %v", fromVM))
		vm, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return nil, err
		}
		logger.Info(fmt.Sprintf("post migrate version map: %v", vm))

		//TODO: confirm the initial deposit
		// update gov params to use a 20% initial deposit ratio, allowing us to remote the ante handler
		govParams := GovKeeper.GetParams(ctx)
		govParams.MinInitialDepositRatio = sdk.NewDec(20).Quo(sdk.NewDec(100)).String()
		if err := GovKeeper.SetParams(ctx, govParams); err != nil {
			return nil, err
		}
		logger.Info(fmt.Sprintf("updated gov params to %v", govParams))

		// x/Mint
		// Double blocks per year (from 6 seconds to 3 = 2x blocks per year)
		mintParams := MintKeeper.GetParams(ctx)
		mintParams.BlocksPerYear *= 2
		if err = MintKeeper.SetParams(ctx, mintParams); err != nil {
			return nil, err
		}
		logger.Info(fmt.Sprintf("updated minted blocks per year logic to %v", mintParams))

		// x/Slashing
		// Double slashing window due to double blocks per year
		slashingParams := SlashingKeeper.GetParams(ctx)
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

		return vm, err
	}
}
