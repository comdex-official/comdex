package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/errors"
	"github.com/comdex-official/comdex/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) RegisterContract(goCtx context.Context, msg *types.MsgRegisterContract) (*types.MsgRegisterContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		ctx.Logger().Error(fmt.Sprintf("request invalid: %s", err))
		return &types.MsgRegisterContractResponse{}, err
	}

	// Validation such that only the user who instantiated the contract can register contract
	contractAddr, err := sdk.AccAddressFromBech32(msg.ContractAddress)
	if err != nil {
		return &types.MsgRegisterContractResponse{}, sdkerrors.ErrInvalidAddress
	}
	contractInfo := k.conOps.GetContractInfo(ctx, contractAddr)

	// check if sender is authorized
	exists := k.CheckSecurityAddress(ctx, msg.SecurityAddress)
	if !exists {
		return &types.MsgRegisterContractResponse{}, sdkerrors.ErrUnauthorized
	}

	allContracts := k.GetAllContract(ctx)

	for _, data := range allContracts {
		if data.ContractAddress == msg.ContractAddress{
			return &types.MsgRegisterContractResponse{}, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "contract already registered")
		}
	}
	gameID := k.GetGameID(ctx)
	contract := types.WhitelistedContract {
		GameId: gameID+1,
		SecurityAddress: msg.SecurityAddress,
		ContractAdmin: contractInfo.Admin,
		GameName: msg.GameName,
		ContractAddress: msg.ContractAddress,
	}

	err = k.SetContract(ctx, contract)
	if err != nil {
		ctx.Logger().Error("failed to set new contract")
		return &types.MsgRegisterContractResponse{}, err
	}
	k.SetGameID(ctx, gameID+1)
	
	return &types.MsgRegisterContractResponse{}, nil
}

func (k msgServer) DeRegisterContract(goCtx context.Context, msg *types.MsgDeRegisterContract) (*types.MsgDeRegisterContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		ctx.Logger().Error(fmt.Sprintf("request invalid: %s", err))
		return &types.MsgDeRegisterContractResponse{}, err
	}

	// Get Game info from Game Id
	gameInfo, found := k.GetContract(ctx, msg.GameId)
	if !found {
		return &types.MsgDeRegisterContractResponse{}, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "no contract found for this game ID")
	}

	// Validation such that only the user who instantiated the contract can register contract
	contractAddr, err := sdk.AccAddressFromBech32(gameInfo.ContractAddress)
	if err != nil {
		return &types.MsgDeRegisterContractResponse{}, sdkerrors.ErrInvalidAddress
	}
	contractInfo := k.conOps.GetContractInfo(ctx, contractAddr)

	// check if sender is authorized
	exists := k.CheckSecurityAddress(ctx, msg.SecurityAddress)
	if !exists && contractInfo.Admin != msg.SecurityAddress{
		return &types.MsgDeRegisterContractResponse{}, sdkerrors.ErrUnauthorized
	}

	k.DeleteContract(ctx, msg.GameId)
	
	return &types.MsgDeRegisterContractResponse{}, nil
}

func (k msgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if k.authority != req.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.SetParams(ctx, req.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}