package keeper

import (
	"context"
	"strconv"

	"github.com/comdex-official/comdex/x/locking/types"
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

func (m msgServer) LockTokens(goCtx context.Context, msg *types.MsgLockTokens) (*types.MsgLockTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	locks := m.Keeper.GetAccountWithSameLockDurationAndDenom(ctx, owner, msg.Coin.Denom, msg.Duration)

	// if existing lock with same duration and denom exists, just add there
	if len(locks) == 1 {
		lock := locks[0]

		_, err = m.Keeper.AddTokensToLockByID(ctx, lock.Id, msg.Coin)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
		}

		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.TypeEvtAddTokensToLock,
				sdk.NewAttribute(types.AttributeLockID, strconv.Itoa(int(lock.Id))),
				sdk.NewAttribute(types.AttributeLockOwner, msg.Owner),
				sdk.NewAttribute(types.AttributeLockAmount, msg.Coin.String()),
			),
		})
		return &types.MsgLockTokensResponse{}, nil
	} else if len(locks) > 1 {
		return &types.MsgLockTokensResponse{}, types.ErrorDuplicateLockExists
	}

	lock, err := m.Keeper.NewLockTokens(ctx, owner, msg.Duration, msg.Coin)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeEvtLockTokens,
			sdk.NewAttribute(types.AttributeLockID, strconv.Itoa(int(lock.Id))),
			sdk.NewAttribute(types.AttributeLockOwner, lock.Owner),
			sdk.NewAttribute(types.AttributeLockAmount, lock.Coin.String()),
			sdk.NewAttribute(types.AttributeLockDuration, lock.Duration.String()),
		),
	})

	return &types.MsgLockTokensResponse{}, nil
}

func (m msgServer) BeginUnlockTokens(goCtx context.Context, msg *types.MsgBeginUnlockingTokens) (*types.MsgBeginUnlockingTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	lock, found := m.Keeper.GetLockByID(ctx, msg.LockId)
	if !found {
		return nil, types.ErrInvalidLockID
	}

	if lock.Owner != owner.String() {
		return nil, types.ErrInvalidLockOwner
	}

	if lock.Coin.Denom != msg.Coin.Denom {
		return nil, types.ErrInvalidUnlockingCoinDenom
	}

	if msg.Coin.Amount.GT(lock.Coin.Amount) {
		return nil, types.ErrInvalidUnlockingAmount
	}

	// update lock on partial unlock else delete it
	isLockUpdate := false
	if msg.Coin.Amount.LT(lock.Coin.Amount) {
		isLockUpdate = true
	}
	//nolint
	m.Keeper.UpdateOrDeleteLock(ctx, isLockUpdate, lock.Id, msg.Coin)
	unlockID, _ := m.Keeper.NewBeginUnlockTokens(ctx, owner, lock, msg.Coin)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeEvtBeginUnlock,
			sdk.NewAttribute(types.AttributeUnlockID, strconv.Itoa(int(unlockID))),
			sdk.NewAttribute(types.AttributeLockOwner, lock.Owner),
			sdk.NewAttribute(types.AttributeUnLockAmount, msg.Coin.String()),
			sdk.NewAttribute(types.AttributeLockDuration, lock.Duration.String()),
		),
	})

	return &types.MsgBeginUnlockingTokensResponse{}, nil
}
