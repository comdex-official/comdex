package keeper

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/comdex-official/comdex/x/liquidationsV2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetGenLiquidationWhiteListing(ctx sdk.Context) (liquidationWhiteListing []types.LiquidationWhiteListing) {
	var (
		store = k.Store(ctx)
		iter  = storetypes.KVStorePrefixIterator(store, types.LiquidationWhiteListingKeyPrefix)
	)

	defer func(iter storetypes.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var data types.LiquidationWhiteListing
		k.cdc.MustUnmarshal(iter.Value(), &data)
		liquidationWhiteListing = append(liquidationWhiteListing, data)
	}

	return liquidationWhiteListing
}

func (k Keeper) GetGenAppReserveFunds(ctx sdk.Context) (appReserveFunds []types.AppReserveFunds) {
	var (
		store = k.Store(ctx)
		iter  = storetypes.KVStorePrefixIterator(store, types.AppReserveFundsKeyPrefix)
	)

	defer func(iter storetypes.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var data types.AppReserveFunds
		k.cdc.MustUnmarshal(iter.Value(), &data)
		appReserveFunds = append(appReserveFunds, data)
	}

	return appReserveFunds
}
