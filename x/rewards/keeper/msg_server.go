package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	assettypes "github.com/comdex-official/comdex/x/asset/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	"github.com/comdex-official/comdex/x/rewards/types"
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

	err := m.Keeper.ValidateMsgCreateGauge(ctx, msg)
	if err != nil {
		return nil, err
	}

	err = m.Keeper.CreateNewGauge(ctx, msg, false)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateGaugeResponse{}, nil
}

func (m msgServer) ExternalRewardsLockers(goCtx context.Context, msg *types.ActivateExternalRewardsLockers) (*types.ActivateExternalRewardsLockersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	klwsParams, _ := m.esm.GetKillSwitchData(ctx, msg.AppMappingId)
	if klwsParams.BreakerEnable {
		return nil, esmtypes.ErrCircuitBreakerEnabled
	}
	esmStatus, found := m.esm.GetESMStatus(ctx, msg.AppMappingId)
	status := false
	if found {
		status = esmStatus.Status
	}
	if status {
		return nil, esmtypes.ErrESMAlreadyExecuted
	}
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
	klwsParams, _ := m.esm.GetKillSwitchData(ctx, msg.AppMappingId)
	if klwsParams.BreakerEnable {
		return nil, esmtypes.ErrCircuitBreakerEnabled
	}
	esmStatus, found := m.esm.GetESMStatus(ctx, msg.AppMappingId)
	status := false
	if found {
		status = esmStatus.Status
	}
	if status {
		return nil, esmtypes.ErrESMAlreadyExecuted
	}
	Depositor, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return nil, err
	}
	pairVault, found := m.asset.GetPairsVault(ctx, msg.ExtendedPairId)
	if !found {
		return nil, assettypes.ErrorExtendedPairDoesNotExistForTheApp
	}
	if pairVault.IsStableMintVault {
		return nil, types.ErrStablemintVaultFound
	}
	if err := m.Keeper.ActExternalRewardsVaults(ctx, msg.AppMappingId, msg.ExtendedPairId, msg.DurationDays, msg.MinLockupTimeSeconds, msg.TotalRewards, Depositor); err != nil {
		return nil, err
	}
	return &types.ActivateExternalRewardsVaultResponse{}, nil
}

func (m msgServer) ExternalRewardsLend(goCtx context.Context, msg *types.ActivateExternalRewardsLend) (*types.ActivateExternalRewardsLendResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	klwsParams, _ := m.esm.GetKillSwitchData(ctx, msg.AppMappingId)
	if klwsParams.BreakerEnable {
		return nil, esmtypes.ErrCircuitBreakerEnabled
	}
	esmStatus, found := m.esm.GetESMStatus(ctx, msg.AppMappingId)
	status := false
	if found {
		status = esmStatus.Status
	}
	if status {
		return nil, esmtypes.ErrESMAlreadyExecuted
	}
	_, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return nil, err
	}
	if err = m.Keeper.AddLendExternalRewards(ctx, *msg); err != nil {
		return nil, err
	}
	return &types.ActivateExternalRewardsLendResponse{}, nil
}

func (m msgServer) ExternalRewardsStableMint(goCtx context.Context, msg *types.ActivateExternalRewardsStableMint) (*types.ActivateExternalRewardsStableMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	klwsParams, _ := m.esm.GetKillSwitchData(ctx, msg.AppId)
	if klwsParams.BreakerEnable {
		return nil, esmtypes.ErrCircuitBreakerEnabled
	}
	esmStatus, found := m.esm.GetESMStatus(ctx, msg.AppId)
	status := false
	if found {
		status = esmStatus.Status
	}
	if status {
		return nil, esmtypes.ErrESMAlreadyExecuted
	}
	Depositor, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return nil, err
	}
	if err := m.Keeper.ActExternalRewardsStableVaults(ctx, msg.AppId, msg.CswapAppId, msg.CommodoAppId, msg.DurationDays, msg.MinLockupTimeSeconds, msg.TotalRewards, Depositor); err != nil {
		return nil, err
	}
	return &types.ActivateExternalRewardsStableMintResponse{}, nil
}
