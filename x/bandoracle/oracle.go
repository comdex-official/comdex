package bandoracle

import (
	"github.com/bandprotocol/bandchain-packet/obi"
	"github.com/bandprotocol/bandchain-packet/packet"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"

	"github.com/comdex-official/comdex/x/bandoracle/types"
)

func (im IBCModule) handleOraclePacket(
	ctx sdk.Context,
	modulePacket channeltypes.Packet,
) (channeltypes.Acknowledgement, error) {
	var ack channeltypes.Acknowledgement
	var modulePacketData packet.OracleResponsePacketData
	fetchPriceMsg := im.keeper.GetFetchPriceMsg(ctx)
	if modulePacket.DestinationChannel != fetchPriceMsg.SourceChannel {
		return channeltypes.Acknowledgement{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,
			"Module packet destination channel and source channel mismatch")
	}
	if err := types.ModuleCdc.UnmarshalJSON(modulePacket.GetData(), &modulePacketData); err != nil {
		return ack, nil
	}

	switch modulePacketData.GetClientID() {
	case types.FetchPriceClientIDKey:
		var fetchPriceResult types.FetchPriceResult
		if err := obi.Decode(modulePacketData.Result, &fetchPriceResult); err != nil {
			ack = channeltypes.NewErrorAcknowledgement(err)
			return ack, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,
				"cannot decode the fetchPrice received packet")
		}
		im.keeper.SetFetchPriceResult(ctx, types.OracleRequestID(modulePacketData.RequestID), fetchPriceResult)
		// TODO: FetchPrice market data reception logic //nolint:godox

	default:
		err := sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal,
			"market received packet not found: %s", modulePacketData.GetClientID())
		ack = channeltypes.NewErrorAcknowledgement(err)
		return ack, err
	}
	ack = channeltypes.NewResultAcknowledgement(
		types.ModuleCdc.MustMarshalJSON(
			packet.NewOracleRequestPacketAcknowledgement(modulePacketData.RequestID),
		),
	)
	return ack, nil
}

func (im IBCModule) handleOracleAcknowledgment(
	ctx sdk.Context,
	ack channeltypes.Acknowledgement,
	modulePacket channeltypes.Packet,
) (*sdk.Result, error) {
	switch resp := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Result:
		var oracleAck packet.OracleRequestPacketAcknowledgement
		err := types.ModuleCdc.UnmarshalJSON(resp.Result, &oracleAck)
		if err != nil {
			return nil, nil
		}

		var data packet.OracleRequestPacketData
		if err = types.ModuleCdc.UnmarshalJSON(modulePacket.GetData(), &data); err != nil {
			return nil, nil
		}
		requestID := types.OracleRequestID(oracleAck.RequestID)

		switch data.GetClientID() {
		case types.FetchPriceClientIDKey:
			var fetchPriceData types.FetchPriceCallData
			if err = obi.Decode(data.GetCalldata(), &fetchPriceData); err != nil {
				return nil, sdkerrors.Wrap(err,
					"cannot decode the fetchPrice market acknowledgment packet")
			}
			im.keeper.SetLastFetchPriceID(ctx, requestID)
			return &sdk.Result{}, nil

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal,
				"market acknowledgment packet not found: %s", data.GetClientID())
		}
	}
	return nil, nil
}
