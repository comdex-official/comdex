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
