package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	protobuftypes "github.com/gogo/protobuf/types"

	"github.com/comdex-official/comdex/x/lend/types"
)

func (k Keeper) AddLendPairsRecords(ctx sdk.Context, records ...types.Extended_Pair) error {
	for _, msg := range records {
		_, found := k.GetLendPair(ctx, msg.Id)
		if found {
			return types.ErrorDuplicateLendPair
		}
		if msg.AssetIn == msg.AssetOut {
			return types.ErrorAssetsCanNotBeSame
		}

		var (
			id   = k.GetLendPairID(ctx)
			pair = types.Extended_Pair{
				Id:              id + 1,
				AssetIn:         msg.AssetIn,
				AssetOut:        msg.AssetOut,
				IsInterPool:     msg.IsInterPool,
				AssetOutPoolID:  msg.AssetOutPoolID,
				MinUsdValueLeft: msg.MinUsdValueLeft,
			}
		)

		k.SetLendPairID(ctx, pair.Id)
		k.SetLendPair(ctx, pair)
	}
	return nil
}

func (k Keeper) AddPoolRecords(ctx sdk.Context, pool types.Pool) error {
	for _, v := range pool.AssetData {
		_, found := k.Asset.GetAsset(ctx, v.AssetID)
		if !found {
			return types.ErrorAssetDoesNotExist
		}
	}

	poolID := k.GetPoolID(ctx)
	newPool := types.Pool{
		PoolID:     poolID + 1,
		ModuleName: pool.ModuleName,
		CPoolName:  pool.CPoolName,
		AssetData:  pool.AssetData,
	}
	for _, v := range pool.AssetData {
		var assetStats types.PoolAssetLBMapping
		assetStats.PoolID = newPool.PoolID
		assetStats.AssetID = v.AssetID
		assetStats.TotalBorrowed = sdk.ZeroInt()
		assetStats.TotalStableBorrowed = sdk.ZeroInt()
		assetStats.TotalLend = sdk.ZeroInt()
		assetStats.TotalInterestAccumulated = sdk.ZeroInt()
		k.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats)
		k.UpdateAPR(ctx, newPool.PoolID, v.AssetID)
		reserveBuybackStats, found := k.GetReserveBuybackAssetData(ctx, v.AssetID)
		if !found {
			reserveBuybackStats.AssetID = v.AssetID
			reserveBuybackStats.ReserveAmount = sdk.ZeroInt()
			reserveBuybackStats.BuybackAmount = sdk.ZeroInt()
			k.SetReserveBuybackAssetData(ctx, reserveBuybackStats)
			reserveStat := types.AllReserveStats{
				AssetID:                        v.AssetID,
				AmountOutFromReserveToLenders:  sdk.ZeroInt(),
				AmountOutFromReserveForAuction: sdk.ZeroInt(),
				AmountInFromLiqPenalty:         sdk.ZeroInt(),
				AmountInFromRepayments:         reserveBuybackStats.BuybackAmount.Add(reserveBuybackStats.ReserveAmount),
				TotalAmountOutToLenders:        sdk.ZeroInt(),
			}
			k.SetAllReserveStatsByAssetID(ctx, reserveStat)
		}
	}

	k.SetPool(ctx, newPool)
	k.SetPoolID(ctx, newPool.PoolID)
	return nil
}

func (k Keeper) AddPoolsPairsRecords(ctx sdk.Context, pool types.PoolPairs) error {
	for _, v := range pool.AssetData {
		_, found := k.Asset.GetAsset(ctx, v.AssetID)
		if !found {
			return types.ErrorAssetDoesNotExist
		}
	}

	poolID := k.GetPoolID(ctx)
	newPool := types.Pool{
		PoolID:     poolID + 1,
		ModuleName: pool.ModuleName,
		CPoolName:  pool.CPoolName,
		AssetData:  pool.AssetData,
	}
	var assetIDPair []uint64
	var mainAssetCurrentPool uint64
	for _, v := range pool.AssetData {
		var assetStats types.PoolAssetLBMapping
		assetStats.PoolID = newPool.PoolID
		assetStats.AssetID = v.AssetID
		assetStats.TotalBorrowed = sdk.ZeroInt()
		assetStats.TotalStableBorrowed = sdk.ZeroInt()
		assetStats.TotalLend = sdk.ZeroInt()
		assetStats.TotalInterestAccumulated = sdk.ZeroInt()
		k.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats)
		k.UpdateAPR(ctx, newPool.PoolID, v.AssetID)
		reserveBuybackStats, found := k.GetReserveBuybackAssetData(ctx, v.AssetID)
		if !found {
			reserveBuybackStats.AssetID = v.AssetID
			reserveBuybackStats.ReserveAmount = sdk.ZeroInt()
			reserveBuybackStats.BuybackAmount = sdk.ZeroInt()
			k.SetReserveBuybackAssetData(ctx, reserveBuybackStats)
			reserveStat := types.AllReserveStats{
				AssetID:                        v.AssetID,
				AmountOutFromReserveToLenders:  sdk.ZeroInt(),
				AmountOutFromReserveForAuction: sdk.ZeroInt(),
				AmountInFromLiqPenalty:         sdk.ZeroInt(),
				AmountInFromRepayments:         reserveBuybackStats.BuybackAmount.Add(reserveBuybackStats.ReserveAmount),
				TotalAmountOutToLenders:        sdk.ZeroInt(),
			}
			k.SetAllReserveStatsByAssetID(ctx, reserveStat)
		}
		assetIDPair = append(assetIDPair, v.AssetID)
		if v.AssetTransitType == 1 {
			mainAssetCurrentPool = v.AssetID
		}
	}

	k.SetPool(ctx, newPool)
	k.SetPoolID(ctx, newPool.PoolID)
	// Add same pool pairs for the specific pool
	// Add inter-pool pairs by fetching all the previous pools
	// create pair mapping for individual pairs

	intraPoolPairs := CreatePairs(assetIDPair) // returns {{1,2},{2,1}}
	for _, pairs := range intraPoolPairs {
		id := k.GetLendPairID(ctx)
		pair := types.Extended_Pair{
			Id:              id + 1,
			AssetIn:         pairs[0],
			AssetOut:        pairs[1],
			IsInterPool:     false,
			AssetOutPoolID:  poolID + 1,
			MinUsdValueLeft: pool.MinUsdValueLeft,
		}
		k.SetLendPairID(ctx, pair.Id)
		k.SetLendPair(ctx, pair)

		// var assetToPair types.AssetToPairMapping
		assetToPair, found := k.GetAssetToPair(ctx, pair.AssetIn, poolID+1)
		if !found {
			assetToPair.AssetID = pair.AssetIn
			assetToPair.PoolID = poolID + 1 // here asset IN and Out Pool are same
			assetToPair.PairID = append(assetToPair.PairID, pair.Id)
		} else {
			assetToPair.PairID = append(assetToPair.PairID, pair.Id)
		}
		k.SetAssetToPair(ctx, assetToPair)
	}

	// Plan:
	// first call all previous pools and get their transit asset type == 1
	// store data and, first create pair as current cPool main asset as assetIn and other {{poolID, mainAssetID}} as assetOut
	// In second step reverse the pair

	allPools := k.GetPools(ctx) // implement changes to get only active pools
	var pools []types.Pool
	for _, v := range allPools {
		if k.IsPoolDepreciated(ctx, v.PoolID) {
			continue
		}
		pools = append(pools, v)
	}
	var diffPoolID, mainAssetID []uint64
	if len(pools) > 1 {
		for _, cPool := range pools {
			if cPool.PoolID == poolID+1 {
				continue
			}
			for _, ad := range cPool.AssetData {
				if ad.AssetTransitType == 1 {
					diffPoolID = append(diffPoolID, cPool.PoolID)
					mainAssetID = append(mainAssetID, ad.AssetID)
					break
				}
			}
		}

		// here we have different cPoolID and their main assetID
		for i := range diffPoolID {
			// first creating pair with different cPool as asset In
			id := k.GetLendPairID(ctx)
			pair := types.Extended_Pair{
				Id:              id + 1,
				AssetIn:         mainAssetID[i],
				AssetOut:        mainAssetCurrentPool,
				IsInterPool:     true,
				AssetOutPoolID:  poolID + 1,
				MinUsdValueLeft: pool.MinUsdValueLeft,
			}
			k.SetLendPairID(ctx, pair.Id)
			k.SetLendPair(ctx, pair)

			assetToPair, found := k.GetAssetToPair(ctx, pair.AssetIn, diffPoolID[i])
			if !found {
				assetToPair.AssetID = pair.AssetIn
				assetToPair.PoolID = diffPoolID[i] // here asset IN and Out Pool are same
				assetToPair.PairID = append(assetToPair.PairID, pair.Id)
			} else {
				assetToPair.PairID = append(assetToPair.PairID, pair.Id)
			}
			k.SetAssetToPair(ctx, assetToPair)
		}

		for i := range diffPoolID {
			// Second creating pair with current cPool as asset In
			id := k.GetLendPairID(ctx)
			pair := types.Extended_Pair{
				Id:              id + 1,
				AssetIn:         mainAssetCurrentPool,
				AssetOut:        mainAssetID[i],
				IsInterPool:     true,
				AssetOutPoolID:  diffPoolID[i],
				MinUsdValueLeft: pool.MinUsdValueLeft,
			}
			k.SetLendPairID(ctx, pair.Id)
			k.SetLendPair(ctx, pair)

			assetToPair, found := k.GetAssetToPair(ctx, pair.AssetIn, poolID+1)
			if !found {
				assetToPair.AssetID = pair.AssetIn
				assetToPair.PoolID = poolID + 1 // here asset IN and Out Pool are same
				assetToPair.PairID = append(assetToPair.PairID, pair.Id)
			} else {
				assetToPair.PairID = append(assetToPair.PairID, pair.Id)
			}
			k.SetAssetToPair(ctx, assetToPair)
		}
	}
	return nil
}

func CreatePairs(numbers []uint64) [][2]uint64 {
	var result [][2]uint64
	for i := 0; i < len(numbers); i++ {
		for j := i + 1; j < len(numbers); j++ {
			result = append(result, [2]uint64{numbers[i], numbers[j]})
			result = append(result, [2]uint64{numbers[j], numbers[i]})
		}
	}
	return result
}

func (k Keeper) AddAssetToPair(ctx sdk.Context, assetToPair types.AssetToPairMapping) error {
	_, found := k.Asset.GetAsset(ctx, assetToPair.AssetID)
	if !found {
		return types.ErrorAssetDoesNotExist
	}
	_, found = k.GetPool(ctx, assetToPair.PoolID)
	if !found {
		return types.ErrPoolNotFound
	}
	for _, v := range assetToPair.PairID {
		_, found = k.GetLendPair(ctx, v)
		if !found {
			return types.ErrorPairDoesNotExist
		}
	}
	k.SetAssetToPair(ctx, assetToPair)

	return nil
}

func (k Keeper) AddMultipleAssetToPair(ctx sdk.Context, assetToPair []types.AssetToPairSingleMapping) error {
	for _, assetPair := range assetToPair {
		assetToPairMap, found := k.GetAssetToPair(ctx, assetPair.AssetID, assetPair.PoolID)
		if found {
			_, found = k.GetLendPair(ctx, assetPair.PairID)
			if !found {
				return types.ErrorPairDoesNotExist
			}
			assetToPairMap.PairID = append(assetToPairMap.PairID, assetPair.PairID)
			k.SetAssetToPair(ctx, assetToPairMap)
		} else {
			var assetToPairMap types.AssetToPairMapping
			_, found := k.Asset.GetAsset(ctx, assetPair.AssetID)
			if !found {
				return types.ErrorAssetDoesNotExist
			}
			_, found = k.GetPool(ctx, assetPair.PoolID)
			if !found {
				return types.ErrPoolNotFound
			}
			_, found = k.GetLendPair(ctx, assetPair.PairID)
			if !found {
				return types.ErrorPairDoesNotExist
			}
			assetToPairMap.AssetID = assetPair.AssetID
			assetToPairMap.PoolID = assetPair.PoolID
			assetToPairMap.PairID = append(assetToPairMap.PairID, assetPair.PairID)
			k.SetAssetToPair(ctx, assetToPairMap)
		}
	}

	return nil
}

func (k Keeper) SetLendPairID(ctx sdk.Context, id uint64) {
	var (
		store = k.Store(ctx)
		key   = types.LendPairIDKey
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: id,
			},
		)
	)

	store.Set(key, value)
}

func (k Keeper) SetLendPair(ctx sdk.Context, pair types.Extended_Pair) {
	var (
		store = k.Store(ctx)
		key   = types.LendPairKey(pair.Id)
		value = k.cdc.MustMarshal(&pair)
	)

	store.Set(key, value)
}

func (k Keeper) GetLendPair(ctx sdk.Context, id uint64) (pair types.Extended_Pair, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.LendPairKey(id)
		value = store.Get(key)
	)

	if value == nil {
		return pair, false
	}

	k.cdc.MustUnmarshal(value, &pair)
	return pair, true
}

func (k Keeper) GetLendPairs(ctx sdk.Context) (pairs []types.Extended_Pair) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.LendPairKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var pair types.Extended_Pair
		k.cdc.MustUnmarshal(iter.Value(), &pair)
		pairs = append(pairs, pair)
	}

	return pairs
}

func (k Keeper) GetLendPairID(ctx sdk.Context) uint64 {
	var (
		store = k.Store(ctx)
		key   = types.LendPairIDKey
		value = store.Get(key)
	)

	if value == nil {
		return 0
	}

	var count protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &count)

	return count.GetValue()
}

func (k Keeper) AddAssetRatesParams(ctx sdk.Context, records ...types.AssetRatesParams) error {
	for _, msg := range records {
		assetRatesParams := types.AssetRatesParams{
			AssetID:              msg.AssetID,
			UOptimal:             msg.UOptimal,
			Base:                 msg.Base,
			Slope1:               msg.Slope1,
			Slope2:               msg.Slope2,
			EnableStableBorrow:   msg.EnableStableBorrow,
			StableBase:           msg.StableBase,
			StableSlope1:         msg.StableSlope1,
			StableSlope2:         msg.StableSlope2,
			Ltv:                  msg.Ltv,
			LiquidationThreshold: msg.LiquidationThreshold,
			LiquidationPenalty:   msg.LiquidationPenalty,
			LiquidationBonus:     msg.LiquidationBonus,
			ReserveFactor:        msg.ReserveFactor,
			CAssetID:             msg.CAssetID,
		}

		k.SetAssetRatesParams(ctx, assetRatesParams)
	}
	return nil
}

func (k Keeper) AddAuctionParamsData(ctx sdk.Context, param types.AuctionParams) error {
	var (
		store = k.Store(ctx)
		key   = types.AuctionParamKey(param.AppId)
		value = k.cdc.MustMarshal(&param)
	)

	store.Set(key, value)

	return nil
}

func (k Keeper) GetAddAuctionParamsData(ctx sdk.Context, appID uint64) (auctionParams types.AuctionParams, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.AuctionParamKey(appID)
		value = store.Get(key)
	)

	if value == nil {
		return auctionParams, false
	}

	k.cdc.MustUnmarshal(value, &auctionParams)
	return auctionParams, true
}

func (k Keeper) GetAllAddAuctionParamsData(ctx sdk.Context) (auctionParams []types.AuctionParams) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.AuctionParamPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var asset types.AuctionParams
		k.cdc.MustUnmarshal(iter.Value(), &asset)
		auctionParams = append(auctionParams, asset)
	}
	return auctionParams
}

func (k Keeper) SetAssetRatesParams(ctx sdk.Context, assetRatesParams types.AssetRatesParams) {
	var (
		store = k.Store(ctx)
		key   = types.AssetRatesParamsKey(assetRatesParams.AssetID)
		value = k.cdc.MustMarshal(&assetRatesParams)
	)

	store.Set(key, value)
}

func (k Keeper) GetAssetRatesParams(ctx sdk.Context, assetID uint64) (assetRatesParams types.AssetRatesParams, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.AssetRatesParamsKey(assetID)
		value = store.Get(key)
	)

	if value == nil {
		return assetRatesParams, false
	}

	k.cdc.MustUnmarshal(value, &assetRatesParams)
	return assetRatesParams, true
}

func (k Keeper) GetAllAssetRatesParams(ctx sdk.Context) (assetRatesParams []types.AssetRatesParams) {
	var (
		store = k.Store(ctx)
		iter  = sdk.KVStorePrefixIterator(store, types.AssetRatesParamsKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var asset types.AssetRatesParams
		k.cdc.MustUnmarshal(iter.Value(), &asset)
		assetRatesParams = append(assetRatesParams, asset)
	}
	return assetRatesParams
}

func (k Keeper) AddAssetRatesPoolPairs(ctx sdk.Context, msg types.AssetRatesPoolPairs) error {
	_, found := k.GetAssetRatesParams(ctx, msg.AssetID)
	if found {
		return types.ErrorAssetRatesParamsAlreadyExists
	}

	assetRatesParams := types.AssetRatesParams{
		AssetID:              msg.AssetID,
		UOptimal:             msg.UOptimal,
		Base:                 msg.Base,
		Slope1:               msg.Slope1,
		Slope2:               msg.Slope2,
		EnableStableBorrow:   msg.EnableStableBorrow,
		StableBase:           msg.StableBase,
		StableSlope1:         msg.StableSlope1,
		StableSlope2:         msg.StableSlope2,
		Ltv:                  msg.Ltv,
		LiquidationThreshold: msg.LiquidationThreshold,
		LiquidationPenalty:   msg.LiquidationPenalty,
		LiquidationBonus:     msg.LiquidationBonus,
		ReserveFactor:        msg.ReserveFactor,
		CAssetID:             msg.CAssetID,
	}

	k.SetAssetRatesParams(ctx, assetRatesParams)

	poolPairs := types.PoolPairs{
		ModuleName:      msg.ModuleName,
		CPoolName:       msg.CPoolName,
		AssetData:       msg.AssetData,
		MinUsdValueLeft: msg.MinUsdValueLeft,
	}
	err := k.AddPoolsPairsRecords(ctx, poolPairs)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) AddPoolDepreciate(ctx sdk.Context, msg types.PoolDepreciate) error {
	depreciatedPoolRecords, found := k.GetPoolDepreciateRecords(ctx)
	if !found {
		k.SetPoolDepreciateRecords(ctx, msg)
	}
	depreciatedPoolRecords.IndividualPoolDepreciate = append(depreciatedPoolRecords.IndividualPoolDepreciate, msg.IndividualPoolDepreciate...)
	k.SetPoolDepreciateRecords(ctx, depreciatedPoolRecords)
	return nil
}

func (k Keeper) SetPoolDepreciateRecords(ctx sdk.Context, msg types.PoolDepreciate) {
	var (
		store = k.Store(ctx)
		key   = types.DepreciatedPoolPrefix
		value = k.cdc.MustMarshal(&msg)
	)

	store.Set(key, value)
}

func (k Keeper) GetPoolDepreciateRecords(ctx sdk.Context) (msg types.PoolDepreciate, found bool) {
	var (
		store = k.Store(ctx)
		key   = types.DepreciatedPoolPrefix
		value = store.Get(key)
	)

	if value == nil {
		return msg, false
	}

	k.cdc.MustUnmarshal(value, &msg)
	return msg, true
}

func (k Keeper) IsPoolDepreciated(ctx sdk.Context, poolID uint64) bool {
	depreciatedPoolRecords, found := k.GetPoolDepreciateRecords(ctx)
	if !found {
		return false
	}
	for _, v := range depreciatedPoolRecords.IndividualPoolDepreciate {
		if v.PoolID == poolID {
			return true
		}
	}
	return false
}

func (k Keeper) DeletePoolAndTransferFunds(ctx sdk.Context) error {
	poolDepRecords, found := k.GetPoolDepreciateRecords(ctx)
	if !found {
		return nil
	}
	for _, poolDepRecord := range poolDepRecords.IndividualPoolDepreciate {
		// condition when the proposal is passed and pool isn't deprecated yet
		if !poolDepRecord.IsPoolDepreciated {
			pool, _ := k.GetPool(ctx, poolDepRecord.PoolID)
			var firstAssetID, secondAssetID, thirdAssetID uint64
			// for getting transit assets details
			for _, data := range pool.AssetData {
				if data.AssetTransitType == 1 {
					firstAssetID = data.AssetID
				}
				if data.AssetTransitType == 2 {
					secondAssetID = data.AssetID
				}
				if data.AssetTransitType == 3 {
					thirdAssetID = data.AssetID
				}
			}
			firstAsset, _ := k.Asset.GetAsset(ctx, firstAssetID)
			secondAsset, _ := k.Asset.GetAsset(ctx, secondAssetID)
			thirdAsset, _ := k.Asset.GetAsset(ctx, thirdAssetID)
			LBMappingFirstAsset, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, poolDepRecord.PoolID, firstAssetID)
			LBMappingSecondAsset, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, poolDepRecord.PoolID, secondAssetID)
			LBMappingThirdAsset, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, poolDepRecord.PoolID, thirdAssetID)
			// condition when all lend and borrow positions are deleted
			if LBMappingFirstAsset.LendIds == nil && LBMappingFirstAsset.BorrowIds == nil && LBMappingSecondAsset.LendIds == nil && LBMappingSecondAsset.BorrowIds == nil && LBMappingThirdAsset.LendIds == nil && LBMappingThirdAsset.BorrowIds == nil {
				// transfer excess funds to the reserve
				// make IsPoolDepreciated true
				// delete Pool
				modBalFirstAsset := k.bank.GetBalance(ctx, authtypes.NewModuleAddress(pool.ModuleName), firstAsset.Denom)
				modBalSecondAsset := k.bank.GetBalance(ctx, authtypes.NewModuleAddress(pool.ModuleName), secondAsset.Denom)
				modBalThirdAsset := k.bank.GetBalance(ctx, authtypes.NewModuleAddress(pool.ModuleName), thirdAsset.Denom)
				err := k.UpdateReserveBalances(ctx, firstAssetID, pool.ModuleName, modBalFirstAsset, true)
				if err != nil {
					return err
				}
				err = k.UpdateReserveBalances(ctx, secondAssetID, pool.ModuleName, modBalSecondAsset, true)
				if err != nil {
					return err
				}
				err = k.UpdateReserveBalances(ctx, thirdAssetID, pool.ModuleName, modBalThirdAsset, true)
				if err != nil {
					return err
				}
				poolDepRecord.IsPoolDepreciated = true
				k.SetPoolDepreciateRecords(ctx, poolDepRecords)
				k.DeletePool(ctx, poolDepRecord.PoolID)
			}
		}
	}
	return nil
}
