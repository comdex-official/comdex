package keeper

import (
	"context"
	"github.com/comdex-official/comdex/x/auctionsV2/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = QueryServer{}

type QueryServer struct {
	Keeper
}

func NewQueryServer(k Keeper) types.QueryServer {
	return &QueryServer{
		Keeper: k,
	}
}

func (q QueryServer) Auction(c context.Context, req *types.QueryAuctionRequest) (res *types.QueryAuctionResponse, err error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx  = sdk.UnwrapSDKContext(c)
		item types.Auction
	)
	if req.History {
		auctionHistorical, _ := q.GetAuctionHistorical(ctx, req.AuctionId)
		item = *auctionHistorical.AuctionHistorical
	} else {
		item, err = q.GetAuction(ctx, req.AuctionId)
	}
	if err != nil {
		return nil, err
	}

	return &types.QueryAuctionResponse{
		Auction: item,
	}, nil
}

func (q QueryServer) Auctions(c context.Context, req *types.QueryAuctionsRequest) (*types.QueryAuctionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.Auction
		ctx   = sdk.UnwrapSDKContext(c)
		key   []byte
	)
	if req.History {
		key = types.AuctionHistoricalKeyPrefix
	} else {
		key = types.AuctionKeyPrefix
	}
	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), key),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.Auction
			if err := q.cdc.Unmarshal(value, &item); err != nil {
				return false, err
			}

			if req.AuctionType == 1 {
				if accumulate {
					if item.AuctionType {
						items = append(items, item)
					}
				}
			} else if req.AuctionType == 2 {
				lockedVault, _ := q.LiquidationsV2.GetLockedVault(ctx, item.AppId, item.LockedVaultId)
				if !item.AuctionType && lockedVault.InitiatorType == "surplus" {
					items = append(items, item)
				}
			} else if req.AuctionType == 3 {
				lockedVault, _ := q.LiquidationsV2.GetLockedVault(ctx, item.AppId, item.LockedVaultId)
				if !item.AuctionType && lockedVault.InitiatorType == "debt" {
					items = append(items, item)
				}
			}
			return true, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAuctionsResponse{
		Auctions:   items,
		Pagination: pagination,
	}, nil
}

func (q QueryServer) Bids(c context.Context, req *types.QueryBidsRequest) (*types.QueryBidsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx   = sdk.UnwrapSDKContext(c)
		key   []byte
		items []types.Bid
	)
	if req.History {
		key = types.GetBidHistoricalKey(req.Bidder)
	} else {
		key = types.GetUserBidHistoricalKey(req.Bidder)
	}

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), key),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.Bid
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

	return &types.QueryBidsResponse{
		Bidder:     req.Bidder,
		Bids:       items,
		Pagination: pagination,
	}, nil
}

func (q QueryServer) AuctionParams(c context.Context, req *types.QueryAuctionParamsRequest) (*types.QueryAuctionParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		ctx  = sdk.UnwrapSDKContext(c)
		item types.AuctionParams
	)

	item, found := q.GetAuctionParams(ctx)
	if !found {
		return nil, types.ErrAuctionParamsNotFound
	}

	return &types.QueryAuctionParamsResponse{
		AuctionParams: item,
	}, nil
}

func (q QueryServer) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{Params: q.GetParams(ctx)}, nil
}

func (q QueryServer) UserLimitBids(c context.Context, req *types.QueryUserLimitBidsRequest) (*types.QueryUserLimitBidsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		item types.LimitOrderBidsForUser
		ctx  = sdk.UnwrapSDKContext(c)
	)

	item, found := q.GetUserLimitBidDataByAddress(ctx, req.Bidder)
	if !found {
		return nil, nil
	}

	return &types.QueryUserLimitBidsResponse{
		LimitOrderBids: item,
	}, nil
}

func (q QueryServer) UserLimitBidsByAssetID(c context.Context, req *types.QueryUserLimitBidsByAssetIDRequest) (*types.QueryUserLimitBidsByAssetIDResponse, error) {

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.LimitOrderBid
		ctx   = sdk.UnwrapSDKContext(c)
		key   []byte
	)
	key = types.LimitBidKeyForAssetID(req.DebtTokenId, req.CollateralTokenId)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), key),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.LimitOrderBid
			if err := q.cdc.Unmarshal(value, &item); err != nil {
				return false, err
			}

			if accumulate {
				if item.BidderAddress == req.Bidder {
					items = append(items, item)
				}
			}

			return true, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryUserLimitBidsByAssetIDResponse{
		Bidder:         req.Bidder,
		LimitOrderBids: items,
		Pagination:     pagination,
	}, nil
}

func (q QueryServer) LimitBids(c context.Context, req *types.QueryLimitBidsRequest) (*types.QueryLimitBidsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.LimitOrderBid
		ctx   = sdk.UnwrapSDKContext(c)
		key   []byte
	)
	key = types.LimitBidKeyForAssetID(req.DebtTokenId, req.CollateralTokenId)

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), key),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.LimitOrderBid
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

	return &types.QueryLimitBidsResponse{
		LimitOrderBids: items,
		Pagination:     pagination,
	}, nil
}

func (q QueryServer) LimitBidProtocolData(c context.Context, req *types.QueryLimitBidProtocolDataRequest) (*types.QueryLimitBidProtocolDataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	var (
		items []types.LimitBidProtocolData
		ctx   = sdk.UnwrapSDKContext(c)
		key   []byte
	)
	key = types.MarketBidProtocolKeyPrefix

	pagination, err := query.FilteredPaginate(
		prefix.NewStore(q.Store(ctx), key),
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			var item types.LimitBidProtocolData
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

	return &types.QueryLimitBidProtocolDataResponse{
		LimitBidProtocolData: items,
		Pagination:           pagination,
	}, nil
}
