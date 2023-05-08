package v11

import (
	// "embed"
	// "fmt"

	// wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	// wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	// sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module"
	// authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	// govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	// paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	// transfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	// ibcratelimit "github.com/osmosis-labs/osmosis/v15/x/ibc-rate-limit"
	// ibcratelimittypes "github.com/osmosis-labs/osmosis/v15/x/ibc-rate-limit/types"
	// packetforwardkeeper "github.com/strangelove-ventures/packet-forward-middleware/v4/router/keeper"
	// packetforwardtypes "github.com/strangelove-ventures/packet-forward-middleware/v4/router/types"
)

// var embedFs embed.FS

// func SetupIBCRateLimitingContract(
// 	ctx sdk.Context,
// 	wasmKeeper wasmkeeper.Keeper,
// 	accountKeeper authkeeper.AccountKeeper,
// 	paramKeeper paramskeeper.Keeper,
// ) error {
// 	govModule := accountKeeper.GetModuleAddress(govtypes.ModuleName)
// 	code, err := embedFs.ReadFile("rate_limiter.wasm")
// 	if err != nil {
// 		return err
// 	}
// 	contractKeeper := wasmkeeper.NewGovPermissionKeeper(wasmKeeper)
// 	instantiateConfig := wasmtypes.AccessConfig{Permission: wasmtypes.AccessTypeOnlyAddress, Address: govModule.String()}
// 	codeID, _, err := contractKeeper.Create(ctx, govModule, code, &instantiateConfig)
// 	if err != nil {
// 		return err
// 	}

// 	transferModule := accountKeeper.GetModuleAddress(transfertypes.ModuleName)

// 	initMsgBz := []byte(fmt.Sprintf(`{
//            "gov_module":  "%s",
//            "ibc_module":"%s",
//            "paths": []
//         }`,
// 		govModule, transferModule))

// 	addr, _, err := contractKeeper.Instantiate(ctx, codeID, govModule, govModule, initMsgBz, "rate limiting contract", nil)
// 	if err != nil {
// 		return err
// 	}
// 	addrStr, err := sdk.Bech32ifyAddressBytes("comdex", addr)
// 	if err != nil {
// 		return err
// 	}
// 	params, err := ibcratelimittypes.NewParams(addrStr)
// 	if err != nil {
// 		return err
// 	}
// 	paramSpace, ok := paramKeeper.GetSubspace(ibcratelimittypes.ModuleName)
// 	if !ok {
// 		return sdkerrors.New("rate-limiting-upgrades", 2, "can't create paramspace")
// 	}
// 	paramSpace.SetParamSet(ctx, &params)
// 	return nil
// }

// func setRateLimits(ctx sdk.Context, accountKeeper authkeeper.AccountKeeper, rateLimitingICS4Wrapper ibcratelimit.ICS4Wrapper, wasmKeeper wasmkeeper.Keeper) {
// 	govModule := accountKeeper.GetModuleAddress(govtypes.ModuleName)
// 	contractKeeper := wasmkeeper.NewGovPermissionKeeper(wasmKeeper)

// 	paths := []string{
// 		`{"add_path": {"channel_id": "any", "denom": "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2",
//           "quotas":
//             [
//               {"name":"ATOM-DAY-1","duration":86400,"send_recv":[30,30]},
//               {"name":"ATOM-DAY-2","duration":129600,"send_recv":[30,30]},
//               {"name":"ATOM-WEEK-1","duration":604800,"send_recv":[60,60]},
//               {"name":"ATOM-WEEK-2","duration":907200,"send_recv":[60,60]}
//             ]
//           }}`,
// 	}

// 	contract := rateLimitingICS4Wrapper.GetContractAddress(ctx)
// 	if contract == "" {
// 		panic("rate limiting contract not set")
// 	}
// 	rateLimitingContract, err := sdk.AccAddressFromBech32(contract)
// 	if err != nil {
// 		panic("contract address improperly formatted")
// 	}
// 	for _, denom := range paths {
// 		_, err := contractKeeper.Execute(ctx, rateLimitingContract, govModule, []byte(denom), nil)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// }

func CreateUpgradeHandlerV11(
	mm *module.Manager,
	configurator module.Configurator,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// if err := SetupIBCRateLimitingContract(ctx, wasmKeeper, accountKeeper, paramKeeper); err != nil {
		// 	return nil, err
		// }
		// packetforwardKeeper.SetParams(ctx, packetforwardtypes.DefaultParams())
		//  N.B.: this is done to avoid initializing genesis for ibcratelimit module.
		// Otherwise, it would overwrite migrations with InitGenesis().
		// See RunMigrations() for details.
		// fromVM[ibcratelimittypes.ModuleName] = 0
		// setRateLimits(ctx, accountKeeper, rateLimitingICS4Wrapper, wasmKeeper)

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
