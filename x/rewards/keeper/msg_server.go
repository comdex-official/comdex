package keeper

import (
	"context"

	"github.com/comdex-official/comdex/x/rewards/types"
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
	err = m.Keeper.bank.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, sdk.NewCoins(newGauge.DepositAmount))
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

func (m msgServer) Whitelist(goCtx context.Context, msg *types.WhitelistAsset) (*types.MsgWhitelistAssetResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := m.Keeper.WhitelistAsset(ctx, msg.AppMappingId, msg.AssetId); err != nil {
		return nil, err
	}
	return &types.MsgWhitelistAssetResponse{}, nil
}

func (m msgServer) RemoveWhitelist(goCtx context.Context, msg *types.RemoveWhitelistAsset) (*types.MsgRemoveWhitelistAssetResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := m.Keeper.RemoveWhitelistAsset(ctx, msg.AppMappingId, msg.AssetId); err != nil {
		return nil, err
	}
	return &types.MsgRemoveWhitelistAssetResponse{}, nil
}

func (m msgServer) WhitelistAppVault(goCtx context.Context, msg *types.WhitelistAppIdVault) (*types.MsgWhitelistAppIdVaultResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := m.Keeper.WhitelistAppIdVault(ctx, msg.AppMappingId); err != nil {
		return nil, err
	}
	return &types.MsgWhitelistAppIdVaultResponse{}, nil
}

func (m msgServer) RemoveWhitelistAppVault(goCtx context.Context, msg *types.RemoveWhitelistAppIdVault) (*types.MsgRemoveWhitelistAppIdVaultResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := m.Keeper.RemoveWhitelistAppIdVault(ctx, msg.AppMappingId); err != nil {
		return nil, err
	}
	return &types.MsgRemoveWhitelistAppIdVaultResponse{}, nil
}

func (m msgServer) ExternalRewardsLockers(goCtx context.Context, msg *types.ActivateExternalRewardsLockers) (*types.ActivateExternalRewardsLockersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	Depositor, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return nil, err
	}
	if err := m.Keeper.ActExternalRewardsLockers(ctx, msg.AppMappingId, msg.AssetId, msg.TotalRewards, msg.DurationDays, Depositor, msg.MinLockupTimeSeconds); err != nil {
		return nil, err
	}
	return &types.ActivateExternalRewardsLockersResponse{}, nil
}

func (m msgServer) ExternalRewardsVault(goCtx context.Context, msg *types.ActivateExternalRewardsVault) (*types.ActivateExternalRewardsVaultResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	Depositor, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return nil, err
	}
	if err := m.Keeper.ActExternalRewardsVaults(ctx, msg.AppMappingId, msg.Extended_Pair_Id, msg.TotalRewards, msg.DurationDays, Depositor, msg.MinLockupTimeSeconds); err != nil {
		return nil, err
	}
	return &types.ActivateExternalRewardsVaultResponse{}, nil
}
