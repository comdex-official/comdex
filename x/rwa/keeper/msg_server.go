package keeper

import (
	"context"
	"github.com/comdex-official/comdex/x/rwa/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ types.MsgServer = (*msgServer)(nil)
)

type msgServer struct {
	keeper Keeper
}

func NewMsgServer(keeper Keeper) types.MsgServer {
	return &msgServer{
		keeper: keeper,
	}
}

func (m msgServer) MsgCreateRwaUser(c context.Context, request *types.CreateRwaUserRequest) (*types.CreateRwaUserResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var rwaUser types.RwaUser
	rwaUser.AccountAddress = request.From
	rwaUser.Name = request.Name
	rwaUser.Address = request.Address
	rwaUser.Owner = request.Owner
	rwaUser.EmailId = request.EmailId
	rwaUser.Jurisdiction = request.Jurisdiction
	rwaUser.KycStatus = types.KycUnverified

	err := m.keeper.SetRwaUser(ctx, rwaUser)
	if err != nil {
		return nil, err
	}

	return &types.CreateRwaUserResponse{}, nil
}

func (m msgServer) MsgCreateCounterParty(c context.Context, request *types.CreateCounterPartyRequest) (*types.CreateCounterPartyResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var counterParty types.Counterparty
	counterParty.Sender = request.From
	counterParty.Receiver = request.SendTo

	err := m.keeper.SetCounterParty(ctx, counterParty)
	if err != nil {
		return nil, err
	}

	counterParty, found := m.keeper.GetCounterparty(ctx, counterParty.Id)

	if found {
		rwaUser1, found1 := m.keeper.GetRwaUSer(ctx, request.From)
		if found1 {
			rwaUser1.CounterpartyList = append(rwaUser1.CounterpartyList, counterParty.Id)
			err = m.keeper.SetRwaUser(ctx, rwaUser1)
			if err != nil {
				return nil, err
			}
		}

		rwaUser2, found2 := m.keeper.GetRwaUSer(ctx, request.SendTo)
		if found2 {
			rwaUser2.CounterpartyList = append(rwaUser1.CounterpartyList, counterParty.Id)
			err = m.keeper.SetRwaUser(ctx, rwaUser2)
			if err != nil {
				return nil, err
			}
		}
	}
	return &types.CreateCounterPartyResponse{}, nil
}

func (m msgServer) MsgUpdateCounterParty(c context.Context, request *types.UpdateCounterPartyRequest) (*types.UpdateCounterPartyResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	counterParty, found := m.keeper.GetCounterparty(ctx, request.Id)

	if !found {
		return &types.UpdateCounterPartyResponse{}, types.ErrCounterPartyNotFound
	}
	counterParty.Accepted = request.Accepted

	err := m.keeper.SetCounterParty(ctx, counterParty)
	if err != nil {
		return nil, err
	}

	return &types.UpdateCounterPartyResponse{}, nil
}

func (m msgServer) MsgUpdateRwaUser(c context.Context, request *types.UpdateRwaUserRequest) (*types.UpdateRwaUserResponse, error) {
	panic("implementation pending")
}

func (m msgServer) MsgUpdateKYC(c context.Context, request *types.UpdateKYCRequest) (*types.UpdateKYCResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	rwaUser, found := m.keeper.GetRwaUSer(ctx, request.From)

	if !found {
		return &types.UpdateKYCResponse{}, types.ErrUserNotFound
	}

	rwaUser.KycStatus = request.Kyc
	rwaUser.KycType = request.KycType

	err := m.keeper.SetRwaUser(ctx, rwaUser)
	if err != nil {
		return nil, err
	}

	return &types.UpdateKYCResponse{}, nil
}

func (m msgServer) MsgCreateInvoice(c context.Context, request *types.CreateInvoiceRequest) (*types.CreateInvoiceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var invoice types.Invoice
	invoice.From = request.From
	invoice.Receiver = request.Receiver
	invoice.Amount = request.Amount
	invoice.AmountPaid = request.AmountPaid
	invoice.CounterpartyId = request.CounterpartyId
	invoice.NftId = request.NftId
	invoice.Receivable = request.Receivable
	invoice.ServiceType = request.ServiceType
	invoice.Status = request.Status

	err := m.keeper.SetInvoice(ctx, invoice)
	if err != nil {
		return nil, err
	}

	return &types.CreateInvoiceResponse{}, nil
}
