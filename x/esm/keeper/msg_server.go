package keeper

import (
	"context"
	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/esm/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
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
	ctx.GasMeter().ConsumeGas(types.DepositESMGas, "DepositESMGas")

	return &types.MsgDepositESMResponse{}, nil
}

func (m msgServer) ExecuteESM(goCtx context.Context, execute *types.MsgExecuteESM) (*types.MsgExecuteESMResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	appID := execute.AppId

	if err := m.keeper.ExecuteESM(ctx, execute.Depositor, appID); err != nil {
		return nil, err
	}
	ctx.GasMeter().ConsumeGas(types.ExecuteESMGas, "ExecuteESMGas")

	return &types.MsgExecuteESMResponse{}, nil
}

func (m msgServer) MsgKillSwitch(c context.Context, msg *types.MsgKillRequest) (*types.MsgKillResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if !m.keeper.Admin(ctx, msg.From) {
		return nil, types.ErrorUnauthorized
	}

	if err := m.keeper.SetKillSwitchData(ctx, *msg.KillSwitchParams); err != nil {
		return nil, err
	}
	ctx.GasMeter().ConsumeGas(types.MsgKillSwitchGas, "MsgKillSwitchGas")

	return &types.MsgKillResponse{}, nil
}

func (m msgServer) MsgCollateralRedemption(c context.Context, req *types.MsgCollateralRedemptionRequest) (*types.MsgCollateralRedemptionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	esmStatus, found := m.keeper.GetESMStatus(ctx, req.AppId)
	status := false
	if found {
		status = esmStatus.Status
	}

	if ctx.BlockTime().Before(esmStatus.EndTime) && status {
		return nil, types.ErrCoolOffPeriodRemains
	}
	if err := m.keeper.CalculateCollateral(ctx, req.AppId, req.Amount, req.From); err != nil {
		return nil, err
	}

	ctx.GasMeter().ConsumeGas(types.MsgCollateralRedemptionGas, "MsgCollateralRedemptionGas")

	return nil, nil
}

func (m msgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if m.keeper.authority != req.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", m.keeper.authority, req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	m.keeper.SetParams(ctx, req.Params)

	return &types.MsgUpdateParamsResponse{}, nil
}