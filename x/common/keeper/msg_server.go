package keeper

import (
	"context"
	"fmt"

	"github.com/comdex-official/comdex/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		return nil, err
	}

	// Validation such that only the user who instantiated the contract can register contract
	contractAddr, _ := sdk.AccAddressFromBech32(msg.ContractAddr)
	contractInfo := k.conOps.GetContractInfo(ctx, contractAddr)

	// TODO: Add wasm fixture to write unit tests to verify this behavior
	if contractInfo.Creator != msg.Creator {
		return nil, sdkerrors.ErrUnauthorized
	}

	allContracts := k.GetAllContract(ctx)

	for _, data := range allContracts {
		if data.GameId == msg.GameId || data.ContractAddr == msg.ContractAddr{
			return &types.MsgRegisterContractResponse{}, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "contract already registered for this game ID")
		}
	}

	err := k.SetContract(ctx, *msg)
	if err != nil {
		ctx.Logger().Error("failed to set new contract")
		return &types.MsgRegisterContractResponse{}, err
	}
	
	return &types.MsgRegisterContractResponse{}, nil
}