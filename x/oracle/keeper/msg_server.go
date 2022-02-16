package keeper

import (
	"context"
	"fmt"
	bandobi "github.com/bandprotocol/bandchain-packet/obi"
	bandpacket "github.com/bandprotocol/bandchain-packet/packet"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v2/modules/core/04-channel/types"
	ibchost "github.com/cosmos/ibc-go/v2/modules/core/24-host"

	"github.com/comdex-official/comdex/x/oracle/types"
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

func (k *msgServer) MsgRemoveMarketForAsset(c context.Context, msg *types.MsgRemoveMarketForAssetRequest) (*types.MsgRemoveMarketForAssetResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if msg.From != k.assetKeeper.Admin(ctx) {
		return nil, types.ErrorUnauthorized
	}

	if !k.HasMarketForAsset(ctx, msg.Id) {
		return nil, types.ErrorMarketForAssetDoesNotExist
	}

	k.DeleteMarketForAsset(ctx, msg.Id)
	return &types.MsgRemoveMarketForAssetResponse{}, nil
}

func (k *msgServer) MsgAddMarket(c context.Context, msg *types.MsgAddMarketRequest) (*types.MsgAddMarketResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if msg.From != k.assetKeeper.Admin(ctx) {
		return nil, types.ErrorUnauthorized
	}

	if k.HasMarket(ctx, msg.Symbol) {
		return nil, types.ErrorDuplicateMarket
	}

	var (
		market = types.Market{
			Symbol:   msg.Symbol,
			ScriptID: msg.ScriptID,
		}
	)

	k.SetMarket(ctx, market)
	return &types.MsgAddMarketResponse{}, nil
}

func (k *msgServer) MsgUpdateMarket(c context.Context, msg *types.MsgUpdateMarketRequest) (*types.MsgUpdateMarketResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if msg.From != k.assetKeeper.Admin(ctx) {
		return nil, types.ErrorUnauthorized
	}

	market, found := k.GetMarket(ctx, msg.Symbol)
	if !found {
		return nil, types.ErrorMarketDoesNotExist
	}

	if msg.ScriptID != 0 {
		market.ScriptID = msg.ScriptID
	}

	k.SetMarket(ctx, market)
	return &types.MsgUpdateMarketResponse{}, nil
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
			return nil, types.ErrorMarketDoesNotExist
		}
		if market.ScriptID != msg.ScriptID {
			return nil, types.ErrorScriptIDMismatch
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

	k.SetCalldataID(ctx, id+1)
	k.SetCalldata(ctx, id, calldata)

	if err := k.channel.SendPacket(ctx, capability, packet); err != nil {
		return nil, err
	}

	return &types.MsgFetchPriceResponse{}, nil
}
