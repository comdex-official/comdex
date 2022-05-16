package keeper

import (
	"context"

	"github.com/comdex-official/comdex/x/incentives/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

func (m msgServer) CreateGauge(goCtx context.Context, msg *types.MsgCreateGauge) (*types.MsgCreateGaugeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := m.Keeper.ValidateMsgCreateCreateGauge(ctx, msg)
	if err != nil {
		return nil, err
	}

	newGauge, err := m.Keeper.NewGauge(ctx, msg)
	if err != nil {
		return nil, err
	}

	gaugeIdsByTriggerDuration, err := m.Keeper.GetUpdatedGaugeIdsByTriggerDurationObj(ctx, newGauge.TriggerDuration, newGauge.Id)
	if err != nil {
		return nil, err
	}

	from, _ := sdk.AccAddressFromBech32(newGauge.From)
	err = m.Keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, sdk.NewCoins(newGauge.DepositAmount))
	if err != nil {
		return nil, err
	}

	_, found := m.Keeper.GetEpochInfoByDuration(ctx, newGauge.TriggerDuration)
	if !found {
		newEpochInfo := m.Keeper.NewEpochInfo(ctx, newGauge.TriggerDuration)
		m.Keeper.SetEpochInfoByDuration(ctx, newEpochInfo)
	}

	m.Keeper.SetGaugeID(ctx, newGauge.Id)
	m.Keeper.SetGauge(ctx, newGauge)
	m.Keeper.SetGaugeIdsByTriggerDuration(ctx, gaugeIdsByTriggerDuration)

	return &types.MsgCreateGaugeResponse{}, nil
}
