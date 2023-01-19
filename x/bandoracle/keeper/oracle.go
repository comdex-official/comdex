package keeper

import (
	"time"

	"github.com/bandprotocol/bandchain-packet/obi"
	"github.com/bandprotocol/bandchain-packet/packet"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v4/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v4/modules/core/24-host"
	protobuftypes "github.com/gogo/protobuf/types"

	"github.com/comdex-official/comdex/x/bandoracle/types"
)

func (k Keeper) SetFetchPriceResult(ctx sdk.Context, requestID types.OracleRequestID, result types.FetchPriceResult) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.FetchPriceResultStoreKey(requestID), k.cdc.MustMarshal(&result))
}

// GetFetchPriceResult returns the FetchPrice by requestId.
func (k Keeper) GetFetchPriceResult(ctx sdk.Context, id types.OracleRequestID) (types.FetchPriceResult, error) {
	bz := ctx.KVStore(k.storeKey).Get(types.FetchPriceResultStoreKey(id))
	if bz == nil {
		return types.FetchPriceResult{}, sdkerrors.Wrapf(types.ErrRequestIDNotAvailable,
			"GetResult: Result for request ID %d is not available.", id,
		)
	}
	var result types.FetchPriceResult
	k.cdc.MustUnmarshal(bz, &result)
	return result, nil
}

// GetLastFetchPriceID return the id from the last FetchPrice request.
func (k Keeper) GetLastFetchPriceID(ctx sdk.Context) int64 {
	bz := ctx.KVStore(k.storeKey).Get(types.KeyPrefix(types.LastFetchPriceIDKey))
	intV := protobuftypes.Int64Value{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &intV)
	return intV.GetValue()
}

// SetLastFetchPriceID saves the id from the last FetchPrice request.
func (k Keeper) SetLastFetchPriceID(ctx sdk.Context, id types.OracleRequestID) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyPrefix(types.LastFetchPriceIDKey),
		k.cdc.MustMarshalLengthPrefixed(&protobuftypes.Int64Value{Value: int64(id)}))
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
	assets := k.assetKeeper.GetAssets(ctx)
	for _, asset := range assets {
		if asset.IsOraclePriceRequired {
			symbol = append(symbol, asset.Name)
		}
	}

	encodedCallData := obi.MustEncode(types.FetchPriceCallData{Symbols: symbol, Multiplier: 1000000})

	packetData := packet.NewOracleRequestPacketData(
		msg.ClientID,
		msg.OracleScriptID,
		encodedCallData,
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

func (k Keeper) SetFetchPriceMsg(ctx sdk.Context, msg types.MsgFetchPriceData) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.MsgDataKey
		v     = types.NewMsgFetchPriceData(
			types.ModuleName,
			types.OracleScriptID(msg.OracleScriptID),
			msg.SourceChannel,
			nil,
			msg.AskCount,
			msg.MinCount,
			msg.FeeLimit,
			msg.PrepareGas,
			msg.ExecuteGas,
			msg.TwaBatchSize,
			msg.AcceptedHeightDiff,
		)
		value = k.cdc.MustMarshal(v)
	)

	store.Set(key, value)
}

func (k Keeper) GetFetchPriceMsg(ctx sdk.Context) types.MsgFetchPriceData {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.MsgDataKey
		value = store.Get(key)
	)
	var msg types.MsgFetchPriceData
	if value != nil {
		k.cdc.MustUnmarshal(value, &msg)
	}

	return msg
}

func (k Keeper) GetLastBlockHeight(ctx sdk.Context) int64 {
	bz := ctx.KVStore(k.storeKey).Get(types.KeyPrefix(types.LastBlockHeightKey))
	intV := protobuftypes.Int64Value{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &intV)
	return intV.GetValue()
}

func (k Keeper) SetLastBlockHeight(ctx sdk.Context, height int64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyPrefix(types.LastBlockHeightKey),
		k.cdc.MustMarshalLengthPrefixed(&protobuftypes.Int64Value{Value: height}))
}

func (k Keeper) AddFetchPriceRecords(ctx sdk.Context, price types.MsgFetchPriceData) error {
	k.SetFetchPriceMsg(ctx, price)
	k.SetLastBlockHeight(ctx, ctx.BlockHeight())
	k.SetCheckFlag(ctx, false)
	k.SetDiscardData(ctx, types.DiscardData{BlockHeight: -1, DiscardBool: false})
	allTwa := k.market.GetAllTwa(ctx)
	for _, data := range allTwa {
		k.market.DeleteTwaData(ctx, data.AssetID)
	}
	return nil
}

func (k Keeper) OraclePriceValidationByRequestID(ctx sdk.Context, req int64) bool {
	currentReqID := k.GetLastFetchPriceID(ctx)

	return currentReqID != req
}

func (k Keeper) SetOracleValidationResult(ctx sdk.Context, res bool) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyPrefix(types.OracleValidationResultKey),
		k.cdc.MustMarshalLengthPrefixed(&protobuftypes.BoolValue{Value: res}))
}

func (k Keeper) GetOracleValidationResult(ctx sdk.Context) bool {
	bz := ctx.KVStore(k.storeKey).Get(types.KeyPrefix(types.OracleValidationResultKey))
	boolV := protobuftypes.BoolValue{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &boolV)
	return boolV.GetValue()
}

func (k Keeper) SetTempFetchPriceID(ctx sdk.Context, id int64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyPrefix(types.TempFetchPriceIDKey),
		k.cdc.MustMarshalLengthPrefixed(&protobuftypes.Int64Value{Value: id}))
}

func (k Keeper) GetTempFetchPriceID(ctx sdk.Context) int64 {
	bz := ctx.KVStore(k.storeKey).Get(types.KeyPrefix(types.TempFetchPriceIDKey))
	intV := protobuftypes.Int64Value{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &intV)
	return intV.GetValue()
}

func (k Keeper) SetCheckFlag(ctx sdk.Context, flag bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.CheckFlagKey
		value = k.cdc.MustMarshal(
			&protobuftypes.BoolValue{
				Value: flag,
			},
		)
	)

	store.Set(key, value)
}

func (k Keeper) GetCheckFlag(ctx sdk.Context) bool {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.CheckFlagKey
		value = store.Get(key)
	)

	var id protobuftypes.BoolValue
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}

func (k Keeper) SetDiscardData(ctx sdk.Context, disData types.DiscardData) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.DiscardFlagKey
		value = k.cdc.MustMarshal(&disData)
	)

	store.Set(key, value)
}

func (k Keeper) GetDiscardData(ctx sdk.Context) types.DiscardData {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.DiscardFlagKey
		value = store.Get(key)
	)

	var disData types.DiscardData
	if value != nil {
		k.cdc.MustUnmarshal(value, &disData)
	}
	return disData
}
