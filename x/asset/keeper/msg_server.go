package keeper

import (
	"context"
	"fmt"

	bandobi "github.com/bandprotocol/bandchain-packet/obi"
	bandpacket "github.com/bandprotocol/bandchain-packet/packet"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	ibchost "github.com/cosmos/ibc-go/modules/core/24-host"

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

func (k *msgServer) MsgFetchPrice(c context.Context, msg *types.MsgFetchPriceRequest) (*types.MsgFetchPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var (
		calldata = types.Calldata{
			Symbols:    []string{},
			Multiplier: k.OracleMultiplier(ctx),
		}
	)

	for _, symbol := range msg.Symbols {
		market, found := k.GetMarket(ctx, symbol)
		if !found {
			return nil, nil
		}
		if market.ScriptID != msg.ScriptID {
			return nil, nil
		}

		calldata.Symbols = append(calldata.Symbols, market.Symbol)
	}

	channel, found := k.channel.GetChannel(ctx, msg.SourcePort, msg.SourceChannel)
	if !found {
		return nil, ibcchanneltypes.ErrChannelNotFound
	}

	sequence, found := k.channel.GetNextSequenceSend(ctx, msg.SourcePort, msg.SourceChannel)
	if !found {
		return nil, ibcchanneltypes.ErrSequenceSendNotFound
	}

	capability, found := k.scoped.GetCapability(ctx, ibchost.ChannelCapabilityPath(msg.SourcePort, msg.SourceChannel))
	if !found {
		return nil, ibcchanneltypes.ErrChannelCapabilityNotFound
	}

	var (
		id     = k.GetCalldataID(ctx)
		packet = ibcchanneltypes.Packet{
			Sequence:           sequence,
			SourcePort:         msg.SourcePort,
			SourceChannel:      msg.SourceChannel,
			DestinationPort:    channel.GetCounterparty().GetPortID(),
			DestinationChannel: channel.GetCounterparty().GetChannelID(),
			Data: bandpacket.OracleRequestPacketData{
				ClientID:       fmt.Sprintf("%d", id),
				OracleScriptID: msg.ScriptID,
				Calldata:       bandobi.MustEncode(calldata),
				AskCount:       k.OracleAskCount(ctx),
				MinCount:       k.OracleMinCount(ctx),
				FeeLimit:       msg.FeeLimit,
				PrepareGas:     msg.PrepareGas,
				ExecuteGas:     msg.ExecuteGas,
			}.GetBytes(),
			TimeoutHeight:    msg.TimeoutHeight,
			TimeoutTimestamp: msg.TimeoutTimestamp,
		}
	)

	k.SetCalldata(ctx, id, calldata)
	k.SetCalldataID(ctx, id+1)

	if err := k.channel.SendPacket(ctx, capability, packet); err != nil {
		return nil, err
	}

	return &types.MsgFetchPriceResponse{}, nil
}
