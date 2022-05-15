package keeper

import (
	"github.com/comdex-official/comdex/x/collector/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) AddLookupTableRecords(ctx sdk.Context, records ...types.AccumulatorLookupTable) error {
	for _, msg := range records {
		var (
			store = ctx.KVStore(k.storeKey)
			key   = types.MsgLookupTableRecordsDataKey
			v     = types.NewMsgLookupTableRecords(
				msg.AccumulatorTokenDenom,
				msg.SecondaryTokenDenom,
				msg.SurplusThreshold,
				msg.DebtThreshold,
				msg.LockerSavingRate,
			)
			value = k.cdc.MustMarshal(v)
		)

		store.Set(key, value)
	}
	return nil

}
