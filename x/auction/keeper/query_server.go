package keeper

import (
	"context"

	"github.com/comdex-official/comdex/x/auction/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServiceServer = (*queryServer)(nil)

type queryServer struct {
	Keeper
}

func NewQueryServiceServer(k Keeper) types.QueryServiceServer {
	return &queryServer{
		Keeper: k,
	}
}

func (q *queryServer) QueryParams(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	var (
		ctx    = sdk.UnwrapSDKContext(c)
		params = q.GetParams(ctx)
	)

	return &types.QueryParamsResponse{
		Params: params,
	}, nil
}

func (q *queryServer) QuerySurplusAuction(c context.Context, req *types.QuerySurplusAuctionRequest) (*types.QuerySurplusAuctionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	item, found := q.GetSurplusAuction(ctx, req.Id)
	if !found {
		return nil, status.Errorf(codes.NotFound, "auction does not exist for id %d", req.Id)
	}

	return &types.QuerySurplusAuctionResponse{
		Auction: item,
	}, nil
}

func (q *queryServer) QuerySurplusAuctions(c context.Context, req *types.QuerySurplusAuctionsRequest) (*types.QuerySurplusAuctionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.SurplusAuction
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.CollateralAuctionKeyPrefix),
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

func (q *queryServer) QuerySurplusBiddings(c context.Context, req *types.QuerySurplusBiddingsRequest) (*types.QuerySurplusBiddingsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	item, found := q.GetSurplusUserBiddings(ctx, req.Bidder)
	if !found {
		return &types.QuerySurplusBiddingsResponse{
			Bidder:   req.Bidder,
			Biddings: []types.Biddings{},
		}, nil
	}

	userBiddings := []types.Biddings{}
	for _, biddingId := range item.BiddingIds {
		bidding, found := q.GetSurplusBidding(ctx, biddingId)
		if !found {
			continue
		}
		userBiddings = append(userBiddings, bidding)
	}

	return &types.QuerySurplusBiddingsResponse{
		Bidder:   req.Bidder,
		Biddings: userBiddings,
	}, nil
}

func (q *queryServer) QueryDebtAuction(c context.Context, req *types.QueryDebtAuctionRequest) (*types.QueryDebtAuctionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	item, found := q.GetDebtAuction(ctx, req.Id)
	if !found {
		return nil, status.Errorf(codes.NotFound, "auction does not exist for id %d", req.Id)
	}

	return &types.QueryDebtAuctionResponse{
		Auction: item,
	}, nil
}

func (q *queryServer) QueryDebtAuctions(c context.Context, req *types.QueryDebtAuctionsRequest) (*types.QueryDebtAuctionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.DebtAuction
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.CollateralAuctionKeyPrefix),
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

func (q *queryServer) QueryDebtBiddings(c context.Context, req *types.QueryDebtBiddingsRequest) (*types.QueryDebtBiddingsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	item, found := q.GetDebtUserBiddings(ctx, req.Bidder)
	if !found {
		return &types.QueryDebtBiddingsResponse{
			Bidder:   req.Bidder,
			Biddings: []types.Biddings{},
		}, nil
	}

	userBiddings := []types.Biddings{}
	for _, biddingId := range item.BiddingIds {
		bidding, found := q.GetDebtBidding(ctx, biddingId)
		if !found {
			continue
		}
		userBiddings = append(userBiddings, bidding)
	}

	return &types.QueryDebtBiddingsResponse{
		Bidder:   req.Bidder,
		Biddings: userBiddings,
	}, nil
}

func (q *queryServer) QueryDutchAuction(c context.Context, req *types.QueryDutchAuctionRequest) (*types.QueryDutchAuctionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	item, found := q.GetDutchAuction(ctx, req.Id)
	if !found {
		return nil, status.Errorf(codes.NotFound, "auction does not exist for id %d", req.Id)
	}

	return &types.QueryDutchAuctionResponse{
		Auction: item,
	}, nil
}

func (q *queryServer) QueryDutchAuctions(c context.Context, req *types.QueryDutchAuctionsRequest) (*types.QueryDutchAuctionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.DutchAuction
		ctx   = sdk.UnwrapSDKContext(c)
	)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), types.CollateralAuctionKeyPrefix),
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

func (q *queryServer) QueryDutchBiddings(c context.Context, req *types.QueryDutchBiddingsRequest) (*types.QueryDutchBiddingsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx = sdk.UnwrapSDKContext(c)
	)

	item, found := q.GetDutchUserBiddings(ctx, req.Bidder)
	if !found {
		return &types.QueryDutchBiddingsResponse{
			Bidder:   req.Bidder,
			Biddings: []types.DutchBiddings{},
		}, nil
	}

	userBiddings := []types.DutchBiddings{}
	for _, biddingId := range item.BiddingIds {
		bidding, found := q.GetDutchBidding(ctx, biddingId)
		if !found {
			continue
		}
		userBiddings = append(userBiddings, bidding)
	}

	return &types.QueryDutchBiddingsResponse{
		Bidder:   req.Bidder,
		Biddings: userBiddings,
	}, nil
}
