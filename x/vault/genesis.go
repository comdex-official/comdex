package vault

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/vault/keeper"
	"github.com/comdex-official/comdex/x/vault/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, state *types.GenesisState) {

	for _, item := range state.Vaults {
		k.SetVault(ctx, item)
	}

	for _, item := range state.StableMintVault {
		k.SetStableMintVault(ctx, item)
	}

	// for _, item := range state.AppExtendedPairVaultMapping {
	// 	err := k.SetAppExtendedPairVaultMappingData(ctx, item)
	// 	if err != nil {
	// 		return
	// 	}
	// }

	// for _, item := range state.UserVaultAssetMapping {
	// 	k.SetUserVaultExtendedPairMapping(ctx, item)
	// }
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.NewGenesisState(
		k.GetVaults(ctx),
		k.GetStableMintVaults(ctx),
		// k.GetAllAppExtendedPairVaultMapping(ctx),
		// k.GetAllUserVaultExtendedPairMapping(ctx),
	)
}
