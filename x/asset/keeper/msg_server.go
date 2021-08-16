package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	ibcchanneltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	ibchost "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"

	"github.com/comdex-official/comdex/x/asset/types"
)

var (
	_ types.MsgServiceServer = (*msgServer)(nil)
)

type msgServer struct {
	Keeper
}

func NewMsgServiceServer(keeper Keeper) types.MsgServiceServer {
	return &msgServer{
		Keeper: keeper,
	}
}

func (k *msgServer) MsgFetchPriceFromBand(c context.Context, msg *types.MsgFetchPriceFromBandRequest) (*types.MsgFetchPriceFromBandResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	source, found := k.channel.GetChannel(ctx, types.PortId, msg.SourceChannel)
	if !found {
		return nil, errors.Wrapf(ibcchanneltypes.ErrChannelNotFound, "channel %s port %s", msg.SourceChannel, types.PortId)
	}

	sequence, found := k.channel.GetNextSequenceSend(ctx, types.PortId, msg.SourceChannel)
	if !found {
		return nil, errors.Wrapf(ibcchanneltypes.ErrSequenceSendNotFound, "channel %s port %s", msg.SourceChannel, types.PortId)
	}

	capability, found := k.scoped.GetCapability(ctx, ibchost.ChannelCapabilityPath(types.PortId, msg.SourceChannel))
	if !found {
		return nil, ibcchanneltypes.ErrChannelCapabilityNotFound
	}

	var (
		packet = ibcchanneltypes.Packet{
			Sequence:           sequence,
			SourcePort:         types.PortId,
			SourceChannel:      msg.SourceChannel,
			DestinationPort:    source.GetCounterparty().GetPortID(),
			DestinationChannel: source.GetCounterparty().GetChannelID(),
			Data:               msg.PacketData.GetBytes(),
			TimeoutHeight:      msg.TimeoutHeight,
			TimeoutTimestamp:   msg.TimeoutTimestamp,
		}
	)

	if err := k.channel.SendPacket(ctx, capability, packet); err != nil {
		return nil, err
	}

	return &types.MsgFetchPriceFromBandResponse{}, nil
}
