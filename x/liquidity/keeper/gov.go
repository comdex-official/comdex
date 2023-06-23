package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/liquidity/types"
)

func (k Keeper) HandelUpdateGenericParamsProposal(ctx sdk.Context, p *types.UpdateGenericParamsProposal) error {
	return k.UpdateGenericParams(ctx, p.AppId, p.Keys, p.Values)
}

func (k Keeper) HandelCreateNewLiquidityPairProposal(ctx sdk.Context, p *types.CreateNewLiquidityPairProposal) error {
	msg := types.NewMsgCreatePair(p.AppId, sdk.MustAccAddressFromBech32(p.From), p.BaseCoinDenom, p.QuoteCoinDenom)
	_, err := k.CreatePair(ctx, msg, true)
	return err
}
