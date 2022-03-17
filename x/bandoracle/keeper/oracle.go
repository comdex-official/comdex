package keeper

import (
	"fmt"
	"time"

	"github.com/bandprotocol/bandchain-packet/obi"
	"github.com/bandprotocol/bandchain-packet/packet"
	"github.com/comdex-official/comdex/x/bandoracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v2/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v2/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v2/modules/core/24-host"
	gogotypes "github.com/gogo/protobuf/types"
)

func (k Keeper) SetFetchPriceResult(ctx sdk.Context, requestID types.OracleRequestID, result types.FetchPriceResult) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.FetchPriceResultStoreKey(requestID), k.cdc.MustMarshal(&result))
}

// GetFetchPriceResult returns the FetchPrice by requestId
func (k Keeper) GetFetchPriceResult(ctx sdk.Context, id types.OracleRequestID) (types.FetchPriceResult, error) {
	bz := ctx.KVStore(k.storeKey).Get(types.FetchPriceResultStoreKey(id))
	if bz == nil {
		return types.FetchPriceResult{}, sdkerrors.Wrapf(types.ErrSample,
			"GetResult: Result for request ID %d is not available.", id,
		)
	}
	var result types.FetchPriceResult
	k.cdc.MustUnmarshal(bz, &result)
	return result, nil
}

// GetLastFetchPriceID return the id from the last FetchPrice request
func (k Keeper) GetLastFetchPriceID(ctx sdk.Context) int64 {
	bz := ctx.KVStore(k.storeKey).Get(types.KeyPrefix(types.LastFetchPriceIDKey))
	intV := gogotypes.Int64Value{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &intV)
	return intV.GetValue()
}

// SetLastFetchPriceID saves the id from the last FetchPrice request
func (k Keeper) SetLastFetchPriceID(ctx sdk.Context, id types.OracleRequestID) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyPrefix(types.LastFetchPriceIDKey),
		k.cdc.MustMarshalLengthPrefixed(&gogotypes.Int64Value{Value: int64(id)}))
}

func (k Keeper) FetchPrice(ctx sdk.Context, msg types.MsgFetchPriceData) (*types.MsgFetchPriceDataResponse, error) {

	sourcePort := types.PortID
	sourceChannelEnd, found := k.channelKeeper.GetChannel(ctx, sourcePort, msg.SourceChannel)
	if !found {
		return nil, nil
	}
	destinationPort := sourceChannelEnd.GetCounterparty().GetPortID()
	destinationChannel := sourceChannelEnd.GetCounterparty().GetChannelID()

	// get the next sequence
	sequence, found := k.channelKeeper.GetNextSequenceSend(ctx, sourcePort, msg.SourceChannel)
	if !found {
		return nil, nil
	}

	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, msg.SourceChannel))
	if !ok {
		return nil, nil
	}

	var symbol []string
	assets := k.GetAssets(ctx)
	for _, asset := range assets {
		symbol = append(symbol, asset.Name)
	}

	encodedCalldata := obi.MustEncode(types.FetchPriceCallData{symbol, 1000000})

	packetData := packet.NewOracleRequestPacketData(
		msg.ClientID,
		msg.OracleScriptID,
		encodedCalldata,
		msg.AskCount,
		msg.MinCount,
		msg.FeeLimit,
		msg.PrepareGas,
		msg.ExecuteGas,
	)
	err := k.channelKeeper.SendPacket(ctx, channelCap, channeltypes.NewPacket(
		packetData.GetBytes(),
		sequence,
		sourcePort,
		msg.SourceChannel,
		destinationPort,
		destinationChannel,
		clienttypes.NewHeight(0, 0),
		uint64(ctx.BlockTime().UnixNano()+int64(10*time.Minute)), // Arbitrary timestamp timeout for now
	))
	if err != nil {
		return nil, nil
	}

	return &types.MsgFetchPriceDataResponse{}, nil
}

func (k *Keeper) SetFetchPriceMsg(ctx sdk.Context, msg types.MsgFetchPriceData) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.MsgdataKey
		v = types.NewMsgFetchPriceData(
			types.ModuleName,
			types.OracleScriptID(msg.OracleScriptID),
			msg.SourceChannel,
			nil,
			msg.AskCount,
			msg.MinCount,
			msg.FeeLimit,
			msg.PrepareGas,
			msg.ExecuteGas,
		)
		value = k.cdc.MustMarshal(v)
	)

	store.Set(key, value)
	fmt.Println("price is being set in store")
	fmt.Println(v)
}

func (k *Keeper) GetFetchPriceMsg(ctx sdk.Context) types.MsgFetchPriceData {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.MsgdataKey
		value = store.Get(key)
	)

	if value == nil {
		fmt.Println("msg value nil")
	}

	var msg types.MsgFetchPriceData
	k.cdc.MustUnmarshal(value, &msg)

	return msg
}

func (k Keeper) GetLastBlockheight(ctx sdk.Context) int64 {
	bz := ctx.KVStore(k.storeKey).Get(types.KeyPrefix(types.LastBlockheightKey))
	intV := gogotypes.Int64Value{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &intV)
	return intV.GetValue()
}

func (k Keeper) SetLastBlockheight(ctx sdk.Context, height int64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyPrefix(types.LastBlockheightKey),
		k.cdc.MustMarshalLengthPrefixed(&gogotypes.Int64Value{Value: int64(height)}))
}

func (k Keeper) AddFetchPriceRecords(ctx sdk.Context, price types.MsgFetchPriceData) error {
	k.SetFetchPriceMsg(ctx, price)
	k.SetLastBlockheight(ctx, ctx.BlockHeight())
	return nil
}
