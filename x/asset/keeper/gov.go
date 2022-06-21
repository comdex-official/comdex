package keeper

import (
	"github.com/comdex-official/comdex/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) HandleProposalAddAsset(ctx sdk.Context, p *types.AddAssetsProposal) error {
	return k.AddAssetRecords(ctx, p.Assets...)
}

func (k Keeper) HandleProposalUpdateAsset(ctx sdk.Context, p *types.UpdateAssetProposal) error {
	return k.UpdateAssetRecords(ctx, p.Asset)
}

func (k Keeper) HandleProposalAddPairs(ctx sdk.Context, p *types.AddPairsProposal) error {
	return k.AddPairsRecords(ctx, p.Pairs...)
}

func (k Keeper) HandleAddWhitelistedAssetRecords(ctx sdk.Context, p *types.AddWhitelistedAssetsProposal) error {
	return k.AddWhitelistedAssetRecords(ctx, p.Assets...)
}

func (k Keeper) HandleUpdateWhitelistedAssetRecords(ctx sdk.Context, p *types.UpdateWhitelistedAssetProposal) error {
	return k.UpdateWhitelistedAssetRecords(ctx, p.Asset)

}

func (k Keeper) HandleAddWhitelistedPairsRecords(ctx sdk.Context, p *types.AddWhitelistedPairsProposal) error {
	return k.AddWhitelistedPairsRecords(ctx, p.Pairs...)
}

func (k Keeper) HandleUpdateWhitelistedPairRecords(ctx sdk.Context, p *types.UpdateWhitelistedPairProposal) error {
	return k.UpdateWhitelistedPairRecords(ctx, p.Pair)
}

func (k Keeper) HandleAddAppMappingRecords(ctx sdk.Context, p *types.AddAppMappingProposal) error {
	return k.AddAppMappingRecords(ctx, p.App...)
}
func (k Keeper) HandleAddAssetMappingRecords(ctx sdk.Context, p *types.AddAssetMappingProposal) error {
	return k.AddAssetMappingRecords(ctx, p.App...)
}

func (k Keeper) HandleAddExtendedPairsVaultRecords(ctx sdk.Context, p *types.AddExtendedPairsVaultProposal) error {
	return k.AddExtendedPairsVaultRecords(ctx, p.Pairs...)
}
