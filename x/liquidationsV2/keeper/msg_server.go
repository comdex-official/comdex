package keeper

import (
	"context"
	"cosmossdk.io/errors"
	"github.com/comdex-official/comdex/x/liquidationsV2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"strconv"
)

type msgServer struct {
	keeper Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) MsgLiquidateInternalKeeper(c context.Context, req *types.MsgLiquidateInternalKeeperRequest) (*types.MsgLiquidateInternalKeeperResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := m.keeper.MsgLiquidate(ctx, req.From, req.LiqType, req.Id); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeLiquidateInternalKeeper,
			sdk.NewAttribute(types.AttributeKeyLiqType, strconv.FormatUint(req.LiqType, 10)),
			sdk.NewAttribute(types.AttributeKeyID, strconv.FormatUint(req.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyCreator, req.From),
		),
	})
	return &types.MsgLiquidateInternalKeeperResponse{}, nil
}

func (m msgServer) MsgAppReserveFunds(c context.Context, req *types.MsgAppReserveFundsRequest) (*types.MsgAppReserveFundsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := m.keeper.MsgAppReserveFundsFn(ctx, req.From, req.AppId, req.AssetId, req.TokenQuantity); err != nil {
		return nil, err
	}
	return &types.MsgAppReserveFundsResponse{}, nil
}

func (m msgServer) MsgLiquidateExternalKeeper(c context.Context, req *types.MsgLiquidateExternalKeeperRequest) (*types.MsgLiquidateExternalKeeperResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if err := m.keeper.MsgLiquidateExternal(ctx, req.From, req.AppId, req.Owner, req.CollateralToken, req.DebtToken, req.CollateralAssetId, req.DebtAssetId, req.IsDebtCmst); err != nil {
		return nil, err
	}
	return &types.MsgLiquidateExternalKeeperResponse{}, nil
}

func (m msgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if m.keeper.authority != req.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", m.keeper.authority, req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	m.keeper.SetParams(ctx, req.Params)

	return &types.MsgUpdateParamsResponse{}, nil
}
