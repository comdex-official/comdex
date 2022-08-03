package keeper

import (
	"context"
	"github.com/comdex-official/comdex/x/esm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
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

func (m msgServer) DepositESM(goCtx context.Context, deposit *types.MsgDepositESM) (*types.MsgDepositESMResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	appID := deposit.AppId

	if err := m.keeper.DepositESM(ctx, deposit.Depositor, appID, deposit.Amount); err != nil {
		return nil, err
	}

	return &types.MsgDepositESMResponse{}, nil
}

func (m msgServer) ExecuteESM(goCtx context.Context, execute *types.MsgExecuteESM) (*types.MsgExecuteESMResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	appID := execute.AppId

	if err := m.keeper.ExecuteESM(ctx, execute.Depositor, appID); err != nil {
		return nil, err
	}

	return &types.MsgExecuteESMResponse{}, nil
}

func (k msgServer) MsgKillSwitch(c context.Context, msg *types.MsgKillRequest) (*types.MsgKillResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if !k.keeper.Admin(ctx, msg.From) {
		return nil, types.ErrorUnauthorized
	}

	if err := k.keeper.SetKillSwitchData(ctx, *msg.KillSwitchParams); err != nil {
		return nil, err
	}

	return &types.MsgKillResponse{}, nil
}

func (k msgServer) MsgCollateralRedemption(c context.Context, req *types.MsgCollateralRedemptionRequest) (*types.MsgCollateralRedemptionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	esmStatus, found := k.keeper.GetESMStatus(ctx, req.AppId)
	status := false
	if found {
		status = esmStatus.Status
	}

	if ctx.BlockTime().Before(esmStatus.EndTime) && status {
		return nil, types.ErrCoolOffPeriodRemains
	}
	if ctx.BlockTime().After(esmStatus.EndTime) && status {
		esmDataAfterCoolOff, _ := k.keeper.GetDataAfterCoolOff(ctx, req.AppId)
		for _, v := range esmDataAfterCoolOff.DebtAsset {
			debtAsset, _ := k.keeper.GetAsset(ctx, v.AssetID)
			if req.Amount.Denom == debtAsset.Denom {
				if !req.Amount.Amount.LTE(v.Amount) {
					return nil, types.ErrorInvalidAmount
				}
				if err := k.keeper.CalculateCollateral(ctx, req.AppId, req.Amount, esmDataAfterCoolOff, req.From); err != nil {
					return nil, err
				}

			} else {
				return nil, types.ErrInvalidAsset
			}
		}
	}

	return nil, nil
}
