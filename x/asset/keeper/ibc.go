package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	ibchost "github.com/cosmos/ibc-go/modules/core/24-host"
)

func (k Keeper) IsBound(ctx sdk.Context, id string) bool {
	_, found := k.scoped.GetCapability(ctx, ibchost.PortPath(id))
	return found
}

func (k *Keeper) BindPort(ctx sdk.Context, id string) error {
	capability := k.port.BindPort(ctx, id)
	return k.ClaimCapability(ctx, capability, ibchost.PortPath(id))
}

func (k *Keeper) ChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	capability, found := k.scoped.GetCapability(ctx, ibchost.ChannelCapabilityPath(portID, channelID))
	if !found {
		return ibcchanneltypes.ErrChannelCapabilityNotFound
	}

	return k.channel.ChanCloseInit(ctx, portID, channelID, capability)
}
