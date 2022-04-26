package keeper

import (
	"fmt"
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) HandleProposalAddAsset(ctx sdk.Context, p *types.AddWhitelistedAssetsProposal) error {
	return k.AddWhitelistedAssetRecords(ctx, p.Assets...)
}

func (k Keeper) HandleProposalUpdateAsset(ctx sdk.Context, p *types.UpdateWhitelistedAssetProposal) error {
	return k.UpdateWhitelistedAssetRecords(ctx, p.Asset)
}

func (k Keeper) HandleProposalAddPairs(ctx sdk.Context, p *types.AddWhitelistedPairsProposal) error {
	return k.AddPairsRecords(ctx, p.Pairs...)
}

func (k Keeper) HandleProposalUpdatePair(ctx sdk.Context, p *types.UpdateWhitelistedPairProposal) error {
	return k.UpdatePairRecords(ctx, p.Pair)
}

func (k *Keeper) AddWhitelistedAssetRecords(ctx sdk.Context, records ...types.Asset) error {
	/*	for _, msg := range records {
		if k.HasAssetForDenom(ctx, msg.Denom) {
			return types.ErrorDuplicateAsset
		}

		var (
			id    = k.GetAssetID(ctx)
			asset = types.Asset{
				Id:       id + 1,
				Name:     msg.Name,
				Denom:    msg.Denom,
				Decimals: msg.Decimals,
			}
		)

		k.SetAssetID(ctx, asset.Id)
		k.SetAsset(ctx, asset)
		k.SetAssetForDenom(ctx, asset.Denom, asset.Id)

	}*/
	fmt.Println("this works")
	fmt.Println("this works")
	fmt.Println("this works")
	fmt.Println("this works")

	return nil
}

func (k *Keeper) UpdateWhitelistedAssetRecords(ctx sdk.Context, msg types.Asset) error {
	return nil
}

func (k *Keeper) AddPairsRecords(ctx sdk.Context, records ...types.Pair) error {
	return nil
}

func (k *Keeper) UpdatePairRecords(ctx sdk.Context, msg types.Pair) error {
	return nil
}
