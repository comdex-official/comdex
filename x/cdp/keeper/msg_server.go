package keeper

import (
	"context"

	"github.com/comdex-official/comdex/x/cdp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ types.MsgServiceServer = (*msgServer)(nil)
)

type msgServer struct {
	Keeper
}

func NewMsgServiceServer(keeper Keeper) types.MsgServiceServer {
	return &msgServer{
		Keeper: keeper,
	}
}

func (k *msgServer) MsgCreate(c context.Context, msg *types.MsgCreateRequest) (*types.MsgCreateResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	if k.HasCDPForAddressByPair(ctx, from, msg.PairID) {
		return nil, types.ErrorDuplicateCDP
	}

	pair, found := k.GetPair(ctx, msg.PairID)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}

	assetIn, found := k.GetAsset(ctx, pair.AssetIn)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	assetOut, found := k.GetAsset(ctx, pair.AssetOut)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	assetInPrice, found := k.GetPriceForAsset(ctx, assetIn.ID)
	if !found {
		return nil, types.ErrorPriceDoesNotExist
	}

	assetOutPrice, found := k.GetPriceForAsset(ctx, assetOut.ID)
	if !found {
		return nil, types.ErrorPriceDoesNotExist
	}

	totalIn := msg.AmountIn.Mul(sdk.NewIntFromUint64(assetInPrice)).QuoRaw(assetIn.Decimals).ToDec()
	if totalIn.IsZero() {
		return nil, types.ErrorInvalidAmount
	}

	totalOut := msg.AmountOut.Mul(sdk.NewIntFromUint64(assetOutPrice)).QuoRaw(assetOut.Decimals).ToDec()
	if totalOut.IsZero() {
		return nil, types.ErrorInvalidAmount
	}

	if totalIn.Quo(totalOut).LT(pair.LiquidationRatio) {
		return nil, types.ErrorInvalidAmountRatio
	}

	if err := k.SendCoinFromAccountToModule(ctx, from, types.ModuleName, sdk.NewCoin(assetIn.Denom, msg.AmountIn)); err != nil {
		return nil, err
	}
	if err := k.MintCoin(ctx, types.ModuleName, sdk.NewCoin(assetOut.Denom, msg.AmountOut)); err != nil {
		return nil, err
	}
	if err := k.SendCoinFromModuleToAccount(ctx, types.ModuleName, from, sdk.NewCoin(assetOut.Denom, msg.AmountOut)); err != nil {
		return nil, err
	}

	var (
		id  = k.GetID(ctx)
		cdp = types.CDP{
			ID:        id + 1,
			PairID:    msg.PairID,
			Owner:     msg.From,
			AmountIn:  msg.AmountIn,
			AmountOut: msg.AmountOut,
		}
	)

	k.SetID(ctx, id+1)
	k.SetCDP(ctx, cdp)
	k.SetCDPForAddressByPair(ctx, from, cdp.PairID, cdp.ID)

	return &types.MsgCreateResponse{}, nil
}

func (k *msgServer) MsgDeposit(c context.Context, msg *types.MsgDepositRequest) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	cdp, found := k.GetCDP(ctx, msg.ID)
	if !found {
		return nil, types.ErrorCDPDoesNotExist
	}
	if msg.From != cdp.Owner {
		return nil, types.ErrorUnauthorized
	}

	pair, found := k.GetPair(ctx, cdp.PairID)
	if !found {
		return nil, types.ErrorPairDoesNotExist
	}

	asset, found := k.GetAsset(ctx, pair.AssetIn)
	if !found {
		return nil, types.ErrorAssetDoesNotExist
	}

	if err := k.SendCoinFromAccountToModule(ctx, from, types.ModuleName, sdk.NewCoin(asset.Denom, msg.Amount)); err != nil {
		return nil, err
	}

	cdp.AmountIn = cdp.AmountIn.Add(msg.Amount)
	k.SetCDP(ctx, cdp)

	return &types.MsgDepositResponse{}, nil
}

func (k *msgServer) MsgWithdraw(c context.Context, msg *types.MsgWithdrawRequest) (*types.MsgWithdrawResponse, error) {
	panic("implement me")
}

func (k *msgServer) MsgDraw(c context.Context, msg *types.MsgDrawRequest) (*types.MsgDrawResponse, error) {
	panic("implement me")
}

func (k *msgServer) MsgRepay(c context.Context, msg *types.MsgRepayRequest) (*types.MsgRepayResponse, error) {
	panic("implement me")
}

func (k *msgServer) MsgLiquidate(c context.Context, msg *types.MsgLiquidateRequest) (*types.MsgLiquidateResponse, error) {
	panic("implement me")
}
