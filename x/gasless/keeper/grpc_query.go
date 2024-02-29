package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/gasless/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper.
type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

// Params queries the parameters of the gasless module.
func (k Querier) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	var params types.Params
	k.Keeper.paramSpace.GetParamSet(ctx, &params)
	return &types.QueryParamsResponse{Params: params}, nil
}

func (k Querier) MessagesAndContracts(c context.Context, _ *types.QueryMessagesAndContractsRequest) (*types.QueryMessagesAndContractsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	messages := k.GetAvailableMessages(ctx)
	contractsDetails := k.GetAllAvailableContracts(ctx)
	contracts := []*types.ContractDetails{}
	for _, c := range contractsDetails {
		contract := c
		contracts = append(contracts, &contract)
	}
	return &types.QueryMessagesAndContractsResponse{
		Messages:  messages,
		Contracts: contracts,
	}, nil
}

func (k Querier) GasProvider(c context.Context, req *types.QueryGasProviderRequest) (*types.QueryGasProviderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.GasProviderId == 0 {
		return nil, status.Error(codes.InvalidArgument, "gas provider id cannot be 0")
	}

	ctx := sdk.UnwrapSDKContext(c)

	gp, found := k.GetGasProvider(ctx, req.GasProviderId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "gas provider with id %d doesn't exist", req.GasProviderId)
	}

	gasTankBalances := k.bankKeeper.GetAllBalances(ctx, sdk.MustAccAddressFromBech32(gp.GasTank))
	return &types.QueryGasProviderResponse{
		GasProvider: types.NewGasProviderResponse(gp, gasTankBalances),
	}, nil
}

func (k Querier) GasProviders(c context.Context, req *types.QueryGasProvidersRequest) (*types.QueryGasProvidersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)

	keyPrefix := types.GetAllGasProvidersKey()
	gpGetter := func(_, value []byte) types.GasProvider {
		return types.MustUnmarshalGasProvider(k.cdc, value)
	}
	gpStore := prefix.NewStore(store, keyPrefix)
	var gasProviders []types.GasProviderResponse

	pageRes, err := query.FilteredPaginate(gpStore, req.Pagination, func(key, value []byte, accumulate bool) (bool, error) {
		gp := gpGetter(key, value)
		if accumulate {
			gasTankBalances := k.bankKeeper.GetAllBalances(ctx, sdk.MustAccAddressFromBech32(gp.GasTank))
			gasProviders = append(gasProviders, types.NewGasProviderResponse(gp, gasTankBalances))
		}

		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryGasProvidersResponse{
		GasProviders: gasProviders,
		Pagination:   pageRes,
	}, nil
}

func (k Querier) GasConsumer(c context.Context, req *types.QueryGasConsumerRequest) (*types.QueryGasConsumerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if _, err := sdk.AccAddressFromBech32(req.Consumer); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid consumer address")
	}

	ctx := sdk.UnwrapSDKContext(c)

	gc, found := k.GetGasConsumer(ctx, sdk.MustAccAddressFromBech32(req.Consumer))
	if !found {
		return nil, status.Errorf(codes.NotFound, "gas consumer %s not found", req.Consumer)
	}
	return &types.QueryGasConsumerResponse{
		GasConsumer: gc,
	}, nil
}

func (k Querier) GasConsumers(c context.Context, req *types.QueryGasConsumersRequest) (*types.QueryGasConsumersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)

	keyPrefix := types.GetAllGasConsumersKey()
	gcGetter := func(_, value []byte) types.GasConsumer {
		return types.MustUnmarshalGasConsumer(k.cdc, value)
	}
	gcStore := prefix.NewStore(store, keyPrefix)
	var gasConsumers []types.GasConsumer

	pageRes, err := query.FilteredPaginate(gcStore, req.Pagination, func(key, value []byte, accumulate bool) (bool, error) {
		gc := gcGetter(key, value)
		if accumulate {
			gasConsumers = append(gasConsumers, gc)
		}

		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryGasConsumersResponse{
		GasConsumers: gasConsumers,
		Pagination:   pageRes,
	}, nil
}
