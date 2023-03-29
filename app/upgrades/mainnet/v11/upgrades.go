package v11

import (
	"embed"
	"fmt"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	icatypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/types"
	transfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	ibcratelimittypes "github.com/osmosis-labs/osmosis/v15/x/ibc-rate-limit/types"
)

//go:embed rate_limiter.wasm
var embedFs embed.FS

func SetupIBCRateLimitingContract(
	ctx sdk.Context,
	wasmKeeper wasmkeeper.Keeper,
	accountKeeper authkeeper.AccountKeeper,
	paramKeeper paramskeeper.Keeper,
) error {
	govModule := accountKeeper.GetModuleAddress(govtypes.ModuleName)
	code, err := embedFs.ReadFile("rate_limiter.wasm")
	if err != nil {
		return err
	}
	contractKeeper := wasmkeeper.NewGovPermissionKeeper(wasmKeeper)
	instantiateConfig := wasmtypes.AccessConfig{Permission: wasmtypes.AccessTypeOnlyAddress, Address: govModule.String()}
	codeID, _, err := contractKeeper.Create(ctx, govModule, code, &instantiateConfig)
	if err != nil {
		return err
	}

	transferModule := accountKeeper.GetModuleAddress(transfertypes.ModuleName)

	initMsgBz := []byte(fmt.Sprintf(`{
           "gov_module":  "%s",
           "ibc_module":"%s",
           "paths": []
        }`,
		govModule, transferModule))

	addr, _, err := contractKeeper.Instantiate(ctx, codeID, govModule, govModule, initMsgBz, "rate limiting contract", nil)
	if err != nil {
		return err
	}
	addrStr, err := sdk.Bech32ifyAddressBytes("ucmdx", addr)
	if err != nil {
		return err
	}
	params, err := ibcratelimittypes.NewParams(addrStr)
	if err != nil {
		return err
	}
	paramSpace, ok := paramKeeper.GetSubspace(ibcratelimittypes.ModuleName)
	if !ok {
		return sdkerrors.New("rate-limiting-upgrades", 2, "can't create paramspace")
	}
	paramSpace.SetParamSet(ctx, &params)
	return nil
}

func CreateUpgradeHandlerV11(
	mm *module.Manager,
	configurator module.Configurator,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		fromVM[icatypes.ModuleName] = mm.Modules[icatypes.ModuleName].ConsensusVersion()

		vm, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return nil, err
		}
		return vm, err
	}
}
