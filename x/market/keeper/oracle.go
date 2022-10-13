package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	assetTypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/market/types"
)

func (k Keeper) SetTwa(ctx sdk.Context, twa types.TimeWeightedAverage) {
	var (
		store = k.Store(ctx)
		key   = types.TwaKey(twa.AssetID)
		value = k.cdc.MustMarshal(&twa)
	)

	store.Set(key, value)
}

func (k Keeper) GetTwa(ctx sdk.Context, id uint64) (twa types.TimeWeightedAverage, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.TwaKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return twa, false
	}

	k.cdc.MustUnmarshal(value, &twa)
	return twa, true
}

func (k Keeper) GetAllTwa(ctx sdk.Context) (twa []types.TimeWeightedAverage) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.TwaKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var data types.TimeWeightedAverage
		k.cdc.MustUnmarshal(iter.Value(), &data)
		twa = append(twa, data)
	}

	return twa
}

func (k Keeper) UpdatePriceList(ctx sdk.Context, id, scriptID, rate uint64) {
	twa, found := k.GetTwa(ctx, id)
	if !found {
		twa.AssetID = id
		twa.ScriptID = scriptID
		twa.Twa = 0
		twa.IsPriceActive = false
		twa.PriceValue = append(twa.PriceValue, rate)
		twa.CurrentIndex = 1
		k.SetTwa(ctx, twa)
	} else {
		if twa.IsPriceActive {
			twa.PriceValue[twa.CurrentIndex] = rate
			twa.CurrentIndex = twa.CurrentIndex + 1
			twa.Twa = k.CalculateTwa(ctx, twa)
			if twa.CurrentIndex == 30 {
				twa.CurrentIndex = 0
			}
			k.SetTwa(ctx, twa)
		} else {
			twa.PriceValue = append(twa.PriceValue, rate)
			twa.CurrentIndex = twa.CurrentIndex + 1
			if twa.CurrentIndex == 30 {
				twa.IsPriceActive = true
				twa.CurrentIndex = 0
				twa.Twa = k.CalculateTwa(ctx, twa)
			}
			k.SetTwa(ctx, twa)
		}
	}
}

func (k Keeper) CalculateTwa(ctx sdk.Context, twa types.TimeWeightedAverage) uint64 {
	var sum uint64
	for _, price := range twa.PriceValue {
		sum += price
	}
	twa.Twa = sum / 30
	return twa.Twa
}

func (k Keeper) GetLatestPrice(ctx sdk.Context, id uint64) (price uint64, err error) {
	twa, found := k.GetTwa(ctx, id)
	if found && twa.IsPriceActive {
		return twa.PriceValue[twa.CurrentIndex], nil
	}
	return 0, types.ErrorPriceNotActive
}

func (k Keeper) CalcAssetPrice(ctx sdk.Context, id uint64, amt sdk.Int) (price sdk.Int, err error) {
	asset, found := k.GetAsset(ctx, id)
	if !found {
		return sdk.ZeroInt(), assetTypes.ErrorAssetDoesNotExist
	}
	twa, found := k.GetTwa(ctx, id)
	if found && twa.IsPriceActive {
		numerator := sdk.NewDecFromInt(amt).Mul(sdk.NewDecFromInt(sdk.NewIntFromUint64(twa.Twa)))
		denominator := sdk.NewDecFromInt(sdk.NewIntFromUint64(uint64(asset.Decimals)))
		result := numerator.Quo(denominator)
		return result.TruncateInt(), nil
	}
	return sdk.ZeroInt(), types.ErrorPriceNotActive
}
