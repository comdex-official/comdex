package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/lend/types"
)

func (k Keeper) HandleAddWhitelistedPairsRecords(ctx sdk.Context, p *types.LendPairsProposal) error {
	return k.AddLendPairsRecords(ctx, p.Pairs)
}

func (k Keeper) HandleAddPoolRecords(ctx sdk.Context, p *types.AddPoolsProposal) error {
	return k.AddPoolRecords(ctx, p.Pool)
}

func (k Keeper) HandleAddAssetToPairRecords(ctx sdk.Context, p *types.AddAssetToPairProposal) error {
	return k.AddAssetToPair(ctx, p.AssetToPairMapping)
}

func (k Keeper) HandleAddMultipleAssetToPairRecords(ctx sdk.Context, p *types.AddMultipleAssetToPairProposal) error {
	return k.AddMultipleAssetToPair(ctx, p.AssetToPairSingleMapping)
}

func (k Keeper) HandleAddAssetRatesParamsRecords(ctx sdk.Context, p *types.AddAssetRatesParams) error {
	return k.AddAssetRatesParams(ctx, p.AssetRatesParams)
}

func (k Keeper) HandleAddAuctionParamsRecords(ctx sdk.Context, p *types.AddAuctionParamsProposal) error {
	return k.AddAuctionParamsData(ctx, p.AuctionParams)
}

func (k Keeper) HandleMultipleAddWhitelistedPairsRecords(ctx sdk.Context, p *types.MultipleLendPairsProposal) error {
	return k.AddLendPairsRecords(ctx, p.Pairs...)
}

func (k Keeper) HandleAddPoolPairsRecords(ctx sdk.Context, p *types.AddPoolPairsProposal) error {
	return k.AddPoolsPairsRecords(ctx, p.PoolPairs)
}

func (k Keeper) HandleAddAssetRatesPoolPairsRecords(ctx sdk.Context, p *types.AddAssetRatesPoolPairsProposal) error {
	return k.AddAssetRatesPoolPairs(ctx, p.AssetRatesPoolPairs)
}
