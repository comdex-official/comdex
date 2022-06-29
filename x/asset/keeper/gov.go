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

func (k Keeper) HandleUpdateGovTimeInAppMapping(ctx sdk.Context, p *types.UpdateGovTimeInAppMappingProposal) error {
	return k.UpdateGovTimeInAppMapping(ctx, p.GovTime)
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
