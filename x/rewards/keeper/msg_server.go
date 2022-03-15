package keeper

import (
	"context"

	"github.com/comdex-official/comdex/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

func NewMsgServiceServer(keeper Keeper) types.MsgServiceServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServiceServer = msgServer{}

func (k msgServer) MsgDepositMintingRewards(goCtx context.Context, msg *types.MsgDepositMintingRewardAmountRequest) (*types.MsgDepositMintingRewardAmountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}
	err = k.TransferDeposits(ctx, msg.MintingRewardId, from, msg.StartTimestamp)
	if err != nil {
		return &types.MsgDepositMintingRewardAmountResponse{}, err
	}
	return &types.MsgDepositMintingRewardAmountResponse{}, nil
}

func (k msgServer) MsgUpdateMintRewardStartTime(goCtx context.Context, msg *types.MsgUpdateMintRewardStartTimeRequest) (*types.MsgUpdateMintRewardStartTimeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}
	err = k.UpdateMintRewardStartTime(ctx, msg.MintingRewardId, from, msg.NewStartTimestamp)
	if err != nil {
		return &types.MsgUpdateMintRewardStartTimeResponse{}, err
	}
	return &types.MsgUpdateMintRewardStartTimeResponse{}, nil
}
