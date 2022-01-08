package keeper

import (
	"context"
	"github.com/comdex-official/comdex/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ types.MsgServiceServer = (*msgServer)(nil)
)

type msgServer struct {
	Keeper
}

func (k *msgServer) FetchPriceData(ctx context.Context, data *types.MsgFetchPriceData) (*types.MsgFetchPriceDataResponse, error) {
	panic("implement me")
}

func NewMsgServiceServer(keeper Keeper) types.MsgServiceServer {
	return &msgServer{
		Keeper: keeper,
	}
}

func (k *msgServer) MsgRemoveMarketForAsset(c context.Context, msg *types.MsgRemoveMarketForAssetRequest) (*types.MsgRemoveMarketForAssetResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if msg.From != k.assetKeeper.Admin(ctx) {
		return nil, types.ErrorUnauthorized
	}

	if !k.HasMarketForAsset(ctx, msg.Id) {
		return nil, types.ErrorMarketForAssetDoesNotExist
	}

	k.DeleteMarketForAsset(ctx, msg.Id)
	return &types.MsgRemoveMarketForAssetResponse{}, nil
}

func (k *msgServer) MsgAddMarket(c context.Context, msg *types.MsgAddMarketRequest) (*types.MsgAddMarketResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if msg.From != k.assetKeeper.Admin(ctx) {
		return nil, types.ErrorUnauthorized
	}
	if !k.HasAsset(ctx, msg.Id){
		return nil, types.ErrorAssetDoesNotExist
	}
	if k.HasMarket(ctx, msg.Symbol) {
		return nil, types.ErrorDuplicateMarket
	}

	var (
		market = types.Market{
			Symbol:   msg.Symbol,
			ScriptID: msg.ScriptID,
		}
	)

	k.SetMarket(ctx, market)
	k.SetPriceForMarket(ctx,msg.Symbol,0)
	ID := k.assetKeeper.GetAssetID(ctx)
	k.SetMarketForAsset(ctx, ID, msg.Symbol )
	return &types.MsgAddMarketResponse{}, nil
}

func (k *msgServer) MsgUpdateMarket(c context.Context, msg *types.MsgUpdateMarketRequest) (*types.MsgUpdateMarketResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if msg.From != k.assetKeeper.Admin(ctx) {
		return nil, types.ErrorUnauthorized
	}

	market, found := k.GetMarket(ctx, msg.Symbol)
	if !found {
		return nil, types.ErrorMarketDoesNotExist
	}

	if msg.ScriptID != 0 {
		market.ScriptID = msg.ScriptID
	}

	k.SetMarket(ctx, market)
	ID := k.assetKeeper.GetAssetID(ctx)
	k.SetMarketForAsset(ctx, ID, msg.Symbol )
	return &types.MsgUpdateMarketResponse{}, nil
}
