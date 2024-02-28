package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/gasless/types"
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
