package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/asset/types"
)

func (k Keeper) HandleProposalAddAsset(ctx sdk.Context, p *types.AddAssetsProposal) error {
	return k.AddAssetRecords(ctx, p.Assets)
}

func (k Keeper) HandleProposalAddMultipleAsset(ctx sdk.Context, p *types.AddMultipleAssetsProposal) error {
	return k.AddMultipleAssetRecords(ctx, p.Assets...)
}

func (k Keeper) HandleProposalUpdateAsset(ctx sdk.Context, p *types.UpdateAssetProposal) error {
	return k.UpdateAssetRecords(ctx, p.Asset)
}

func (k Keeper) HandleProposalAddPairs(ctx sdk.Context, p *types.AddPairsProposal) error {
	return k.AddPairsRecords(ctx, p.Pairs)
}

func (k Keeper) HandleProposalAddMultiplePairs(ctx sdk.Context, p *types.AddMultiplePairsProposal) error {
	return k.AddMultiplePairsRecords(ctx, p.Pairs...)
}

func (k Keeper) HandleProposalUpdatePair(ctx sdk.Context, p *types.UpdatePairProposal) error {
	return k.UpdatePairRecords(ctx, p.Pairs)
}

func (k Keeper) HandleUpdateGovTimeInApp(ctx sdk.Context, p *types.UpdateGovTimeInAppProposal) error {
	return k.UpdateGovTimeInApp(ctx, p.GovTime)
}

func (k Keeper) HandleAddAppRecords(ctx sdk.Context, p *types.AddAppProposal) error {
	return k.AddAppRecords(ctx, p.App)
}

func (k Keeper) HandleAddAssetInAppRecords(ctx sdk.Context, p *types.AddAssetInAppProposal) error {
	return k.AddAssetInAppRecords(ctx, p.App)
}

func (k Keeper) HandleProposalAddMultipleAssetPair(ctx sdk.Context, p *types.AddMultipleAssetsPairsProposal) error {
	return k.AddMultipleAssetPairRecords(ctx, p.AssetsPair...)
}
