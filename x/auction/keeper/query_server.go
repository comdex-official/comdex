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

func (q *QueryServer) QuerySurplusAuction(c context.Context, req *types.QuerySurplusAuctionRequest) (res *types.QuerySurplusAuctionResponse, err error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx  = sdk.UnwrapSDKContext(c)
		item auctiontypes.SurplusAuction
	)
	if req.History {
		item, err = q.GetHistorySurplusAuction(ctx, req.AppId, req.AuctionMappingId, req.AuctionId)
	} else {
		item, err = q.GetSurplusAuction(ctx, req.AppId, req.AuctionMappingId, req.AuctionId)
	}
	if err != nil {
		return nil, err
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
	if req.History {
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

func (q *QueryServer) QuerySurplusBiddings(c context.Context, req *types.QuerySurplusBiddingsRequest) (res *types.QuerySurplusBiddingsResponse, err error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx  = sdk.UnwrapSDKContext(c)
		item []auctiontypes.SurplusBiddings
	)
	if req.History {
		item = q.GetHistorySurplusUserBiddings(ctx, req.Bidder, req.AppId)
	} else {
		item = q.GetSurplusUserBiddings(ctx, req.Bidder, req.AppId)
	}

	return &types.QuerySurplusBiddingsResponse{
		Bidder:   req.Bidder,
		Biddings: item,
	}, nil
}
func (q *QueryServer) QueryDebtAuction(c context.Context, req *types.QueryDebtAuctionRequest) (res *types.QueryDebtAuctionResponse, err error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx  = sdk.UnwrapSDKContext(c)
		item auctiontypes.DebtAuction
	)
	if req.History {
		item, err = q.GetHistoryDebtAuction(ctx, req.AppId, req.AuctionMappingId, req.AuctionId)
	} else {
		item, err = q.GetDebtAuction(ctx, req.AppId, req.AuctionMappingId, req.AuctionId)
	}
	if err != nil {
		return res, err
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
	if req.History {
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
		ctx  = sdk.UnwrapSDKContext(c)
		item []auctiontypes.DebtBiddings
	)
	if req.History {
		item = q.GetHistoryDebtUserBiddings(ctx, req.Bidder, req.AppId)
	} else {
		item = q.GetDebtUserBiddings(ctx, req.Bidder, req.AppId)
	}

	return &types.QueryDebtBiddingsResponse{
		Bidder:   req.Bidder,
		Biddings: item,
	}, nil
}

func (q *QueryServer) QueryDutchAuction(c context.Context, req *types.QueryDutchAuctionRequest) (res *types.QueryDutchAuctionResponse, err error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx  = sdk.UnwrapSDKContext(c)
		item auctiontypes.DutchAuction
	)
	if req.History {
		item, _ = q.GetHistoryDutchAuction(ctx, req.AppId, req.AuctionMappingId, req.AuctionId)
	} else {
		item, _ = q.GetDutchAuction(ctx, req.AppId, req.AuctionMappingId, req.AuctionId)
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
	if req.History {
		key = types.HistoryAuctionTypeKey(req.AppId, types.DutchString)
	} else {
		key = types.AuctionTypeKey(req.AppId, types.DutchString)
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
		ctx  = sdk.UnwrapSDKContext(c)
		item []auctiontypes.DutchBiddings
	)
	if req.History {
		item = q.GetHistoryDutchUserBiddings(ctx, req.Bidder, req.AppId)
	} else {
		item = q.GetDutchUserBiddings(ctx, req.Bidder, req.AppId)
	}

	return &types.QueryDutchBiddingsResponse{
		Bidder:   req.Bidder,
		Biddings: item,
	}, nil
}

func (q *QueryServer) QueryBiddingsForAuction(c context.Context, req *types.QueryDutchBiddingsRequest) (*types.QueryDutchBiddingsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx  = sdk.UnwrapSDKContext(c)
		item []auctiontypes.DutchBiddings
	)
	if req.History {
		item = q.GetHistoryDutchUserBiddings(ctx, req.Bidder, req.AppId)
	} else {
		item = q.GetDutchUserBiddings(ctx, req.Bidder, req.AppId)
	}

	return &types.QueryDutchBiddingsResponse{
		Bidder:   req.Bidder,
		Biddings: item,
	}, nil
}

func (q *QueryServer) QueryProtocolStatistics(c context.Context, req *types.QueryProtocolStatisticsRequest) (*types.QueryProtocolStatisticsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.ProtocolStatistics
		ctx   = sdk.UnwrapSDKContext(c)
		key   []byte
	)

	key = types.ProtocolStatisticsAppIdKey(req.AppId)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), key),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.ProtocolStatistics
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

	return &types.QueryProtocolStatisticsResponse{
		Stats:      items,
		Pagination: pagination,
	}, nil
}

func (q *QueryServer) QueryAuctionParams(c context.Context, req *auctiontypes.QueryAuctionParamRequest) (*auctiontypes.QueryAuctionParamResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	item, found := q.GetAuctionParams(ctx, req.AppId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "Auction Params not exist for id %d", req.AppId)
	}

	return &types.QueryAuctionParamResponse{
		AuctionParams: item,
	}, nil
}
