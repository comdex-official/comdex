package keeper

import (
	"github.com/comdex-official/comdex/x/asset/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)



func (k Keeper) HandleProposalAddAsset(ctx sdk.Context, p *types.AddAssetsProposal) error {
	return k.AddAssetRecords(ctx, p.Assets...)
}
