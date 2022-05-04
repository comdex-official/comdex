package keeper

import (
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
	for _, msg := range records {
		if k.HasAssetForDenom(ctx, msg.Denom) {
			return types.ErrorDuplicateAsset
		}

		var (
			id    = k.GetAssetID(ctx)
			asset = types.Asset{
				Id:                   id + 1,
				Name:                 msg.Name,
				Denom:                msg.Denom,
				Decimals:             msg.Decimals,
				CollateralWeight:     msg.CollateralWeight,
				LiquidationThreshold: msg.LiquidationThreshold,
				IsBridgedAsset:       msg.IsBridgedAsset,
			}
		)

		k.SetAssetID(ctx, asset.Id)
		k.SetAsset(ctx, asset)
		k.SetAssetForDenom(ctx, asset.Denom, asset.Id)

	}

	return nil
}

func (k *Keeper) UpdateWhitelistedAssetRecords(ctx sdk.Context, msg types.Asset) error {
	asset, found := k.GetAsset(ctx, msg.Id)
	if !found {
		return types.ErrorAssetDoesNotExist
	}
	if len(msg.Name) > 0 {
		asset.Name = msg.Name
	}
	if len(msg.Denom) > 0 {
		if k.HasAssetForDenom(ctx, msg.Denom) {
			return types.ErrorDuplicateAsset
		}
		asset.Denom = msg.Denom

		k.DeleteAssetForDenom(ctx, asset.Denom)
		k.SetAssetForDenom(ctx, asset.Denom, asset.Id)

	}
	if msg.Decimals >= 0 {
		asset.Decimals = msg.Decimals
	}
	if !msg.CollateralWeight.IsZero() {
		asset.CollateralWeight = msg.CollateralWeight
	}
	if !msg.LiquidationThreshold.IsZero() {
		asset.LiquidationThreshold = msg.LiquidationThreshold
	}

	k.SetAsset(ctx, asset)
	return nil
}

func (k *Keeper) AddPairsRecords(ctx sdk.Context, records ...types.Pair) error {
	for _, msg := range records {
		if !k.HasAsset(ctx, msg.Asset_1) {
			return types.ErrorAssetDoesNotExist
		}
		if !k.HasAsset(ctx, msg.Asset_2) {
			return types.ErrorAssetDoesNotExist
		}

		var (
			id   = k.GetPairID(ctx)
			pair = types.Pair{
				Id:                    id + 1,
				Asset_1:               msg.Asset_1,
				Asset_2:               msg.Asset_2,
				ModuleAcc:             msg.ModuleAcc,
				BaseBorrowRateAsset_1: msg.BaseBorrowRateAsset_1,
				BaseLendRateAsset_1:   msg.BaseLendRateAsset_1,
				BaseBorrowRateAsset_2: msg.BaseBorrowRateAsset_2,
				BaseLendRateAsset_2:   msg.BaseLendRateAsset_2,
			}
		)

		k.SetPairID(ctx, pair.Id)
		k.SetPair(ctx, pair)
	}
	return nil
}

func (k *Keeper) UpdatePairRecords(ctx sdk.Context, msg types.Pair) error {

	pair, found := k.GetPair(ctx, msg.Id)
	if !found {
		return types.ErrorPairDoesNotExist
	}

	if msg.Asset_1 > 0 {
		pair.Asset_1 = msg.Asset_1
	}

	if msg.Asset_2 > 0 {
		pair.Asset_2 = msg.Asset_2
	}

	if len(msg.ModuleAcc) > 0 {
		pair.ModuleAcc = msg.ModuleAcc
	}
	k.SetPair(ctx, pair)
	return nil
}
