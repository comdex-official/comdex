package keeper

import (
	"github.com/comdex-official/comdex/x/bandoracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
