package vault

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/vault/keeper"
	"github.com/comdex-official/comdex/x/vault/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, state *types.GenesisState) {
	var (
		vaultID       uint64
		stableVaultID uint64
	)

	for _, item := range state.Vaults {
		if item.Id > vaultID {
			vaultID = item.Id
		}
		k.SetVault(ctx, item)
	}

	for _, item := range state.StableMintVault {
		if item.Id > stableVaultID {
			stableVaultID = item.Id
		}
		k.SetStableMintVault(ctx, item)
	}

	for _, item := range state.AppExtendedPairVaultMapping {
		k.SetAppExtendedPairVaultMappingData(ctx, item)
	}

	for _, item := range state.UserVaultAssetMapping {
		k.SetUserAppExtendedPairMappingData(ctx, item)
	}
	k.SetIDForVault(ctx, vaultID)
	k.SetIDForStableVault(ctx, stableVaultID)
	k.SetLengthOfVault(ctx, state.LengthOfVaults)
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.NewGenesisState(
		k.GetVaults(ctx),
		k.GetStableMintVaults(ctx),
		k.GetAllAppExtendedPairVaultMapping(ctx),
		k.GetAllUserVaultExtendedPairMapping(ctx),
		k.GetLengthOfVault(ctx),
	)
}
