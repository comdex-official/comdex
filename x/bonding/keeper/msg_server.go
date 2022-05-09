package keeper

import (
	"context"
	"fmt"

	"github.com/comdex-official/comdex/x/bonding/types"
	"github.com/comdex-official/comdex/x/bonding/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type msgServer struct {
	keeper *Keeper
}

// NewMsgServerImpl returns an instance of MsgServer.
func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{
		keeper: keeper,
	}
}

var _ types.MsgServer = msgServer{}

func (server msgServer) LockTokens(goCtx context.Context, msg *types.MsgLockTokens) (*types.MsgLockTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	fmt.Println("on line 28....")

	// we only allow locks with one denom for now
	if msg.Coins.Len() != 1 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("Bondings can only have one denom per lockID, got %v", msg.Coins))
	}
	fmt.Println("on line 35....")

	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}
	fmt.Println("on line 41....")

	if len(msg.Coins) == 1 {
		locks := server.keeper.GetAccountLockedDurationNotUnlockingOnly(ctx, owner, msg.Coins[0].Denom, msg.Duration)
		// if existing lock with same duration and denom exists, just add there
		if len(locks) > 0 {
			fmt.Println("on line 44....")
			lock := locks[0]
			if lock.Coins.Len() != 1 {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
			}
			fmt.Println("on line 49....")
			if lock.Owner != owner.String() {
				return nil, types.ErrNotLockOwner
			}
			fmt.Println("on line 53....")
			_, err = server.keeper.AddTokensToLockByID(ctx, lock.ID, msg.Coins)
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
			}
			fmt.Println("on line 58....")

			ctx.EventManager().EmitEvents(sdk.Events{
				sdk.NewEvent(
					types.TypeEvtAddTokensToLock,
					sdk.NewAttribute(types.AttributePeriodLockID, utils.Uint64ToString(locks[0].ID)),
					sdk.NewAttribute(types.AttributePeriodLockOwner, msg.Owner),
					sdk.NewAttribute(types.AttributePeriodLockAmount, msg.Coins.String()),
				),
			})
			return &types.MsgLockTokensResponse{ID: locks[0].ID}, nil
		}
	}

	lock, err := server.keeper.LockTokens(ctx, owner, msg.Coins, msg.Duration)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	fmt.Println("on line 76....")
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeEvtLockTokens,
			sdk.NewAttribute(types.AttributePeriodLockID, utils.Uint64ToString(lock.ID)),
			sdk.NewAttribute(types.AttributePeriodLockOwner, lock.Owner),
			sdk.NewAttribute(types.AttributePeriodLockAmount, lock.Coins.String()),
			sdk.NewAttribute(types.AttributePeriodLockDuration, lock.Duration.String()),
			sdk.NewAttribute(types.AttributePeriodLockUnlockTime, lock.EndTime.String()),
		),
	})

	return &types.MsgLockTokensResponse{ID: lock.ID}, nil
}

func (server msgServer) BeginUnlocking(goCtx context.Context, msg *types.MsgBeginUnlocking) (*types.MsgBeginUnlockingResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	lock, err := server.keeper.GetLockByID(ctx, msg.ID)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	err = server.keeper.BeginUnlock(ctx, lock.ID, msg.Coins)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	if msg.Owner != lock.Owner {
		return nil, sdkerrors.Wrap(types.ErrNotLockOwner, fmt.Sprintf("msg sender(%s) and lock owner(%s) does not match", msg.Owner, lock.Owner))
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		createBeginUnlockEvent(lock),
	})

	return &types.MsgBeginUnlockingResponse{}, nil
}

func (server msgServer) BeginUnlockingAll(goCtx context.Context, msg *types.MsgBeginUnlockingAll) (*types.MsgBeginUnlockingAllResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return nil, err
	}

	unlocks, err := server.keeper.BeginUnlockAllNotUnlockings(ctx, owner)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	// Create the events for this message
	unlockedCoins := server.keeper.getCoinsFromLocks(unlocks)
	events := sdk.Events{
		sdk.NewEvent(
			types.TypeEvtBeginUnlockAll,
			sdk.NewAttribute(types.AttributePeriodLockOwner, msg.Owner),
			sdk.NewAttribute(types.AttributeUnlockedCoins, unlockedCoins.String()),
		),
	}
	for _, lock := range unlocks {
		events = events.AppendEvent(createBeginUnlockEvent(&lock))
	}
	ctx.EventManager().EmitEvents(events)

	return &types.MsgBeginUnlockingAllResponse{}, nil
}

func createBeginUnlockEvent(lock *types.PeriodLock) sdk.Event {
	return sdk.NewEvent(
		types.TypeEvtBeginUnlock,
		sdk.NewAttribute(types.AttributePeriodLockID, utils.Uint64ToString(lock.ID)),
		sdk.NewAttribute(types.AttributePeriodLockOwner, lock.Owner),
		sdk.NewAttribute(types.AttributePeriodLockDuration, lock.Duration.String()),
		sdk.NewAttribute(types.AttributePeriodLockUnlockTime, lock.EndTime.String()),
	)
}
