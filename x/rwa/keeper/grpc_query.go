package keeper

import (
	"context"
	"github.com/comdex-official/comdex/x/rwa/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

func (q QueryServer) QueryRwaUser(c context.Context, request *types.RwaUserRequest) (*types.RwaUserResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	rwaUser, found := q.GetRwaUSer(ctx, request.AccountAddress)
	if !found {
		return &types.RwaUserResponse{}, nil
	}

	return &types.RwaUserResponse{
		RwaUser: &rwaUser,
	}, nil
}

func (q QueryServer) QueryCounterParty(c context.Context, request *types.CounterPartyRequest) (*types.CounterPartyResponse, error) {

	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	counterParty, found := q.GetCounterparty(ctx, request.Id)
	if !found {
		return &types.CounterPartyResponse{}, nil
	}

	return &types.CounterPartyResponse{
		CounterParty: &counterParty,
	}, nil
}

func (q QueryServer) QueryInvoice(ctx context.Context, request *types.InvoiceRequest) (*types.InvoiceResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q QueryServer) QueryInvoiceAsSender(ctx context.Context, request *types.InvoiceSenderRequest) (*types.InvoiceResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (q QueryServer) QueryInvoiceAsReceiver(ctx context.Context, request *types.InvoiceReceiverRequest) (*types.InvoiceResponse, error) {
	//TODO implement me
	panic("implement me")
}
