package keeper

import (
	assetTypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/market/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

func (k Keeper) DeleteTwaData(ctx sdk.Context, assetID uint64) {
	var (
		store = k.Store(ctx)
		key   = types.TwaKey(assetID)
	)

	store.Delete(key)
}

func (k Keeper) UpdatePriceList(ctx sdk.Context, id, scriptID, rate, twaBatch uint64) {
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
			twa.Twa = k.CalculateTwa(ctx, twa, twaBatch)
			if twa.CurrentIndex == twaBatch {
				twa.CurrentIndex = 0
			}
			k.SetTwa(ctx, twa)
		} else {
			twa.PriceValue = append(twa.PriceValue, rate)
			twa.CurrentIndex = twa.CurrentIndex + 1
			if twa.CurrentIndex == twaBatch {
				twa.IsPriceActive = true
				twa.CurrentIndex = 0
				twa.Twa = k.CalculateTwa(ctx, twa, twaBatch)
			}
			k.SetTwa(ctx, twa)
		}
	}
}

func (k Keeper) CalculateTwa(ctx sdk.Context, twa types.TimeWeightedAverage, twaBatch uint64) uint64 {
	var sum uint64
	for i := 0; i < int(twaBatch); i++ {
		sum = sum + twa.PriceValue[i]
	}
	twa.Twa = sum / twaBatch
	return twa.Twa
}

func (k Keeper) GetLatestPrice(ctx sdk.Context, id uint64) (price uint64, err error) {
	twa, found := k.GetTwa(ctx, id)
	if found && twa.IsPriceActive {
		return twa.PriceValue[twa.CurrentIndex], nil
	}
	return 0, types.ErrorPriceNotActive
}

func (k Keeper) CalcAssetPrice(ctx sdk.Context, id uint64, amt sdk.Int) (price sdk.Dec, err error) {
	asset, found := k.assetKeeper.GetAsset(ctx, id)
	if !found {
		return sdk.ZeroDec(), assetTypes.ErrorAssetDoesNotExist
	}
	twa, found := k.GetTwa(ctx, id)
	if found && twa.IsPriceActive {
		numerator := sdk.NewDecFromInt(amt).Mul(sdk.NewDecFromInt(sdk.NewIntFromUint64(twa.Twa)))
		denominator := sdk.NewDecFromInt(sdk.NewIntFromUint64(uint64(asset.Decimals)))
		return numerator.Quo(denominator), nil
	}
	return sdk.ZeroDec(), types.ErrorPriceNotActive
}

func (k Keeper) CalcAssetPrice(ctx sdk.Context, id uint64, amt sdk.Int) (price sdk.Dec, err error) {
    asset, found := k.assetKeeper.GetAsset(ctx, id)
    if !found {
        return sdk.ZeroDec(), assetTypes.ErrorAssetDoesNotExist
    }
    // twa, found := k.GetTwa(ctx, id)
    var rate uint64
    if id == 1 {
        rate = 11632845
    }
    if id == 2 {
        rate = 14053
    }
    if id == 3 {
        rate = 1000000
    }
    if id == 4 {
        rate = 2200000
    }
    if id == 10 {
        rate = 1297119384
    }
    if found {
        numerator := sdk.NewDecFromInt(amt).Mul(sdk.NewDecFromInt(sdk.NewIntFromUint64(rate)))
        denominator := sdk.NewDecFromInt(sdk.NewIntFromUint64(uint64(asset.Decimals)))
        return numerator.Quo(denominator), nil
    }
    return sdk.ZeroDec(), types.ErrorPriceNotActive
}

func (k Keeper) GetTwa(ctx sdk.Context, id uint64) (twa types.TimeWeightedAverage, found bool) {
    twa.AssetID = id
    twa.IsPriceActive = true
    if id == 1 {
        twa.Twa = 11632845
    }
    if id == 2 {
        twa.Twa = 140530
    }
    if id == 3 {
        twa.Twa = 1000000
    }
    if id == 4 {
        twa.Twa = 2200000
    }
    if id == 10 {
        twa.Twa = 1297119384
    }

    return twa, true
}

