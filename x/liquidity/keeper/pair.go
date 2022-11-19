package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/petrichormoney/petri/x/liquidity/types"
)

// getNextPairIdWithUpdate increments pair id by one and set it.
func (k Keeper) getNextPairIDWithUpdate(ctx sdk.Context, appID uint64) uint64 {
	id := k.GetLastPairID(ctx, appID) + 1
	k.SetLastPairID(ctx, appID, id)
	return id
}

// getNextOrderIdWithUpdate increments the pair's last order id and returns it.
func (k Keeper) getNextOrderIDWithUpdate(ctx sdk.Context, pair types.Pair) uint64 {
	id := pair.LastOrderId + 1
	pair.LastOrderId = id
	k.SetPair(ctx, pair)
	return id
}

// ValidateMsgCreatePair validates types.MsgCreatePair.
func (k Keeper) ValidateMsgCreatePair(ctx sdk.Context, msg *types.MsgCreatePair) error {
	_, found := k.assetKeeper.GetApp(ctx, msg.AppId)
	if !found {
		return sdkerrors.Wrapf(types.ErrInvalidAppID, "app id %d not found", msg.AppId)
	}

	if !k.assetKeeper.HasAssetForDenom(ctx, msg.BaseCoinDenom) {
		return sdkerrors.Wrapf(types.ErrAssetNotWhiteListed, "asset with denom  %s is not white listed", msg.BaseCoinDenom)
	}

	if !k.assetKeeper.HasAssetForDenom(ctx, msg.QuoteCoinDenom) {
		return sdkerrors.Wrapf(types.ErrAssetNotWhiteListed, "asset with denom  %s is not white listed", msg.QuoteCoinDenom)
	}

	if _, found := k.GetPairByDenoms(ctx, msg.AppId, msg.BaseCoinDenom, msg.QuoteCoinDenom); found {
		return types.ErrPairAlreadyExists
	}
	return nil
}

// CreatePair handles types.MsgCreatePair and creates a pair.
func (k Keeper) CreatePair(ctx sdk.Context, msg *types.MsgCreatePair, isViaProp bool) (types.Pair, error) {
	if err := k.ValidateMsgCreatePair(ctx, msg); err != nil {
		return types.Pair{}, err
	}

	params, err := k.GetGenericParams(ctx, msg.AppId)
	if err != nil {
		return types.Pair{}, sdkerrors.Wrap(err, "params retreval failed")
	}

	// ignore fee collection if the request is from proposal
	if !isViaProp {
		// Send the pair creation fee to the fee collector.
		feeCollectorAddr, _ := sdk.AccAddressFromBech32(params.FeeCollectorAddress)
		if err := k.bankKeeper.SendCoins(ctx, msg.GetCreator(), feeCollectorAddr, params.PairCreationFee); err != nil {
			return types.Pair{}, sdkerrors.Wrap(err, "insufficient pair creation fee")
		}
	}

	id := k.getNextPairIDWithUpdate(ctx, msg.AppId)
	pair := types.NewPair(id, msg.BaseCoinDenom, msg.QuoteCoinDenom, msg.AppId)
	k.SetPair(ctx, pair)
	k.SetPairIndex(ctx, msg.AppId, pair.BaseCoinDenom, pair.QuoteCoinDenom, pair.Id)
	k.SetPairLookupIndex(ctx, msg.AppId, pair.BaseCoinDenom, pair.QuoteCoinDenom, pair.Id)
	k.SetPairLookupIndex(ctx, msg.AppId, pair.QuoteCoinDenom, pair.BaseCoinDenom, pair.Id)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreatePair,
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyBaseCoinDenom, msg.BaseCoinDenom),
			sdk.NewAttribute(types.AttributeKeyQuoteCoinDenom, msg.QuoteCoinDenom),
			sdk.NewAttribute(types.AttributeKeyPairID, strconv.FormatUint(pair.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyEscrowAddress, pair.EscrowAddress),
		),
	})

	return pair, nil
}
