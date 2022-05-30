package keeper

import (
	"context"

	"github.com/comdex-official/comdex/x/auction/types"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = (*QueryServer)(nil)

type QueryServer struct {
	Keeper
}

func NewQueryServiceServer(k Keeper) types.QueryServer {
	return &QueryServer{
		Keeper: k,
	}
}

func (q *QueryServer) QueryParams(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	var (
		ctx    = sdk.UnwrapSDKContext(c)
		params = q.GetParams(ctx)
	)

	return &types.QueryParamsResponse{
		Params: params,
	}, nil
}

func (q *QueryServer) QuerySurplusAuction(c context.Context, req *types.QuerySurplusAuctionRequest) (*types.QuerySurplusAuctionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx   = sdk.UnwrapSDKContext(c)
		found bool
		item  auctiontypes.SurplusAuction
	)
	if req.History == true {
		item, found = q.GetHistorySurplusAuction(ctx, req.AppId, req.AuctionMappingId, req.AuctionId)
	} else {
		item, found = q.GetSurplusAuction(ctx, req.AppId, req.AuctionMappingId, req.AuctionId)
	}
	if !found {
		return nil, status.Errorf(codes.NotFound, "auction does not exist for id %d", req.AuctionId)
	}

	return &types.QuerySurplusAuctionResponse{
		Auction: item,
	}, nil
}

func (q *QueryServer) QuerySurplusAuctions(c context.Context, req *types.QuerySurplusAuctionsRequest) (*types.QuerySurplusAuctionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.SurplusAuction
		ctx   = sdk.UnwrapSDKContext(c)
		key   []byte
	)
	if req.History == true {
		key = types.HistoryAuctionTypeKey(req.AppId, types.SurplusString)
	} else {
		key = types.AuctionTypeKey(req.AppId, types.SurplusString)
	}
	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), key),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.SurplusAuction
			if err := q.cdc.Unmarshal(value, &item); err != nil {
				return false, err
			}

			if accumulate {
				items = append(items, item)
			}

			return true, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QuerySurplusAuctionsResponse{
		Auctions:   items,
		Pagination: pagination,
	}, nil
}

func (q *QueryServer) QuerySurplusBiddings(c context.Context, req *types.QuerySurplusBiddingsRequest) (*types.QuerySurplusBiddingsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx   = sdk.UnwrapSDKContext(c)
		found bool
		item  []auctiontypes.SurplusBiddings
	)
	if req.History == true {
		item, found = q.GetHistorySurplusUserBiddings(ctx, req.Bidder, req.AppId)
	} else {
		item, found = q.GetSurplusUserBiddings(ctx, req.Bidder, req.AppId)
	}
	if !found {
		return &types.QuerySurplusBiddingsResponse{
			Bidder:   req.Bidder,
			Biddings: []types.SurplusBiddings{},
		}, nil
	}

	return &types.QuerySurplusBiddingsResponse{
		Bidder:   req.Bidder,
		Biddings: item,
	}, nil
}

func (q *QueryServer) QueryDebtAuction(c context.Context, req *types.QueryDebtAuctionRequest) (*types.QueryDebtAuctionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx   = sdk.UnwrapSDKContext(c)
		found bool
		item  auctiontypes.DebtAuction
	)
	if req.History == true {
		item, found = q.GetHistoryDebtAuction(ctx, req.AppId, req.AuctionMappingId, req.AuctionId)
	} else {
		item, found = q.GetDebtAuction(ctx, req.AppId, req.AuctionMappingId, req.AuctionId)
	}
	if !found {
		return nil, status.Errorf(codes.NotFound, "auction does not exist for id %d", req.AuctionId)
	}

	return &types.QueryDebtAuctionResponse{
		Auction: item,
	}, nil
}

func (q *QueryServer) QueryDebtAuctions(c context.Context, req *types.QueryDebtAuctionsRequest) (*types.QueryDebtAuctionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.DebtAuction
		ctx   = sdk.UnwrapSDKContext(c)
		key   []byte
	)
	if req.History == true {
		key = types.HistoryAuctionTypeKey(req.AppId, types.DebtString)
	} else {
		key = types.AuctionTypeKey(req.AppId, types.DebtString)
	}

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), key),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.DebtAuction
			if err := q.cdc.Unmarshal(value, &item); err != nil {
				return false, err
			}

			if accumulate {
				items = append(items, item)
			}

			return true, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDebtAuctionsResponse{
		Auctions:   items,
		Pagination: pagination,
	}, nil
}

func (q *QueryServer) QueryDebtBiddings(c context.Context, req *types.QueryDebtBiddingsRequest) (*types.QueryDebtBiddingsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx   = sdk.UnwrapSDKContext(c)
		found bool
		item  []auctiontypes.DebtBiddings
	)
	if req.History == true {
		item, found = q.GetHistoryDebtUserBiddings(ctx, req.Bidder, req.AppId)
	} else {
		item, found = q.GetDebtUserBiddings(ctx, req.Bidder, req.AppId)
	}
	if !found {
		return &types.QueryDebtBiddingsResponse{
			Bidder:   req.Bidder,
			Biddings: []types.DebtBiddings{},
		}, nil
	}

	return &types.QueryDebtBiddingsResponse{
		Bidder:   req.Bidder,
		Biddings: item,
	}, nil
}

func (q *QueryServer) QueryDutchAuction(c context.Context, req *types.QueryDutchAuctionRequest) (*types.QueryDutchAuctionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx   = sdk.UnwrapSDKContext(c)
		found bool
		item  auctiontypes.DutchAuction
	)
	if req.History == true {
		item, found = q.GetHistoryDutchAuction(ctx, req.AppId, req.AuctionMappingId, req.AuctionId)
	} else {
		item, found = q.GetDutchAuction(ctx, req.AppId, req.AuctionMappingId, req.AuctionId)
	}
	if !found {
		return nil, status.Errorf(codes.NotFound, "auction does not exist for id %d", req.AuctionId)
	}

	return &types.QueryDutchAuctionResponse{
		Auction: item,
	}, nil
}

func (q *QueryServer) QueryDutchAuctions(c context.Context, req *types.QueryDutchAuctionsRequest) (*types.QueryDutchAuctionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.DutchAuction
		ctx   = sdk.UnwrapSDKContext(c)
		key   []byte
	)
	if req.History == true {
		key = types.HistoryAuctionTypeKey(req.AppId, types.DutchString)
	} else {
		key = types.HistoryAuctionTypeKey(req.AppId, types.DutchString)
	}
	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), key),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.DutchAuction
			if err := q.cdc.Unmarshal(value, &item); err != nil {
				return false, err
			}

			if accumulate {
				items = append(items, item)
			}

			return true, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryDutchAuctionsResponse{
		Auctions:   items,
		Pagination: pagination,
	}, nil
}

func (q *QueryServer) QueryDutchBiddings(c context.Context, req *types.QueryDutchBiddingsRequest) (*types.QueryDutchBiddingsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx   = sdk.UnwrapSDKContext(c)
		found bool
		item  []auctiontypes.DutchBiddings
	)
	if req.History == true {
		item, found = q.GetHistoryDutchUserBiddings(ctx, req.Bidder, req.AppId)
	} else {
		item, found = q.GetDutchUserBiddings(ctx, req.Bidder, req.AppId)
	}
	if !found {
		return &types.QueryDutchBiddingsResponse{
			Bidder:   req.Bidder,
			Biddings: []types.DutchBiddings{},
		}, nil
	}

	return &types.QueryDutchBiddingsResponse{
		Bidder:   req.Bidder,
		Biddings: item,
	}, nil
}
