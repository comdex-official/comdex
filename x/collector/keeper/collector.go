package keeper

import (
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	"github.com/comdex-official/comdex/x/collector/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

func (k *Keeper) GetAmountFromCollector(ctx sdk.Context, appId, asset_id uint64, amount sdk.Int) (sdk.Int, error) {
	netfeeData, found := k.GetNetFeeCollectedData(ctx, appId)
	var returnedFee sdk.Int
	if !found {
		return returnedFee, types.ErrorDataDoesNotExists
	}

	for _, data := range netfeeData.AssetIdToFeeCollected {
		if data.AssetId == asset_id {
			if !(data.NetFeesCollected.Sub(amount).GT(sdk.ZeroInt())) {
				return returnedFee, types.ErrorRequestedAmtExceedsCollectedFee
			} else {
				asset, _ := k.GetAsset(ctx, asset_id)
				if err := k.SendCoinFromModuleToModule(ctx, types.ModuleName, auctiontypes.ModuleName, sdk.NewCoins(sdk.NewCoin(asset.Denom, data.NetFeesCollected.Sub(amount)))); err != nil {
					return returnedFee, err
				}
				k.SetNetFeeCollectedData(ctx, appId, asset_id, data.NetFeesCollected.Sub(amount))
			}
		}

	}
	return returnedFee, nil
}

func (k *Keeper) UpdateCollector(ctx sdk.Context, appId, asset_id uint64, CollectedStabilityFee, CollectedClosingFee, CollectedOpeningFee, LiquidationRewardsCollected sdk.Int) error {

	if !k.HasAsset(ctx, asset_id) {
		return types.ErrorAssetDoesNotExist
	}

	collectorData, found := k.GetAppidToAssetCollectorMapping(ctx, appId)

	if !found {
		//create a new instance of AppId To AssetCollectorMapping
		var collectorNewData types.AppIdToAssetCollectorMapping
		collectorNewData.AppId = appId

		var assetIdCollect types.AssetIdCollectorMappping
		assetIdCollect.AssetId = asset_id

		var newCollector types.CollectorData

		newCollector.CollectedClosingFee = CollectedClosingFee
		newCollector.CollectedOpeningFee = CollectedOpeningFee
		newCollector.CollectedStabilityFee = CollectedStabilityFee
		newCollector.LiquidationRewardsCollected = LiquidationRewardsCollected
		assetIdCollect.Collector = newCollector

		collectorNewData.AssetCollector = append(collectorNewData.AssetCollector, &assetIdCollect)

		k.SetAppidToAssetCollectorMapping(ctx, collectorNewData)
	} else {

		var check = 0 // makes it 1 if assetId does not exists for appId

		for _, data := range collectorData.AssetCollector {

			if data.AssetId != asset_id { //if does not exist then create a new instance
				check++
				var collectorNewData types.AppIdToAssetCollectorMapping
				collectorNewData.AppId = appId

				var assetIdCollect types.AssetIdCollectorMappping
				assetIdCollect.AssetId = asset_id

				var newCollector types.CollectorData

				newCollector.CollectedClosingFee = CollectedClosingFee
				newCollector.CollectedOpeningFee = CollectedOpeningFee
				newCollector.CollectedStabilityFee = CollectedStabilityFee
				newCollector.LiquidationRewardsCollected = LiquidationRewardsCollected
				assetIdCollect.Collector = newCollector

				collectorNewData.AssetCollector = append(collectorNewData.AssetCollector, &assetIdCollect)

				k.SetAppidToAssetCollectorMapping(ctx, collectorNewData)

				return nil
			} else {
				continue
			}
		}

		if check == 0 {
			var collectorNewData types.AppIdToAssetCollectorMapping
			collectorNewData.AppId = appId

			var assetIdCollect types.AssetIdCollectorMappping
			assetIdCollect.AssetId = asset_id

			var newCollector types.CollectorData

			newCollector.CollectedClosingFee = CollectedClosingFee
			newCollector.CollectedOpeningFee = CollectedOpeningFee
			newCollector.CollectedStabilityFee = CollectedStabilityFee
			newCollector.LiquidationRewardsCollected = LiquidationRewardsCollected
			assetIdCollect.Collector = newCollector

			collectorNewData.AssetCollector = append(collectorNewData.AssetCollector, &assetIdCollect)

			k.SetAppidToAssetCollectorMapping(ctx, collectorNewData)
		}
	}
	return nil

}

func (k *Keeper) SetAppidToAssetCollectorMapping(ctx sdk.Context, appAssetCollectorData types.AppIdToAssetCollectorMapping) {

	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.AppidToAssetCollectorMappingKey(appAssetCollectorData.AppId)
		value = k.cdc.MustMarshal(&appAssetCollectorData)
	)

	store.Set(key, value)

}

func (k *Keeper) GetAppidToAssetCollectorMapping(ctx sdk.Context, app_id uint64) (appAssetCollectorData types.AppIdToAssetCollectorMapping, found bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.AppidToAssetCollectorMappingKey(app_id)
		value = store.Get(key)
	)

	if value == nil {
		return appAssetCollectorData, false
	}

	k.cdc.MustUnmarshal(value, &appAssetCollectorData)
	return appAssetCollectorData, true
}

//////////////////////////////111111111111111111111

func (k *Keeper) SetCollectorLookupTable(ctx sdk.Context, records ...types.CollectorLookupTable) error {
	for _, msg := range records {
		if !k.HasAsset(ctx, msg.CollectorAssetId) {
			return types.ErrorAssetDoesNotExist
		}
		if !k.HasAsset(ctx, msg.SecondaryAssetId) {
			return types.ErrorAssetDoesNotExist
		}
		if msg.CollectorAssetId == msg.SecondaryAssetId {
			return types.ErrorDuplicateAssetDenoms
		}
		appDenom, found := k.GetAppToDenomsMapping(ctx, msg.AppId)
		if found {
			//check if assetdenom already exists
			var check = 0
			for _, data := range appDenom.AssetIds {
				if data == msg.CollectorAssetId {
					check++
				}
			}
			if check > 0 {
				return types.ErrorDuplicateCollectorDenomForApp
			}
			// if denom is new then append
			appDenom.AssetIds = append(appDenom.AssetIds, msg.CollectorAssetId)
			k.SetAppToDenomsMapping(ctx, msg.AppId, appDenom)

		} else {
			//initialize the mappping
			var appDenomNew types.AppToDenomsMapping
			appDenomNew.AppId = msg.AppId
			appDenomNew.AssetIds = append(appDenomNew.AssetIds, msg.CollectorAssetId)
			k.SetAppToDenomsMapping(ctx, msg.AppId, appDenomNew)
		}

		var Collector = types.CollectorLookupTable{
			AppId:            msg.AppId,
			CollectorAssetId: msg.CollectorAssetId,
			SecondaryAssetId: msg.SecondaryAssetId,
			SurplusThreshold: msg.SurplusThreshold,
			DebtThreshold:    msg.DebtThreshold,
			LockerSavingRate: msg.LockerSavingRate,
			LotSize:          msg.LotSize,
			BidFactor:        msg.BidFactor,
		}
		accmLookup, _ := k.GetCollectorLookupTable(ctx, msg.AppId)
		accmLookup.AppId = msg.AppId
		accmLookup.AssetrateInfo = append(accmLookup.AssetrateInfo, Collector)

		var (
			store = ctx.KVStore(k.storeKey)
			key   = types.CollectorLookupTableMappingKey(msg.AppId)
			value = k.cdc.MustMarshal(&accmLookup)
		)

		store.Set(key, value)
	}
	return nil

}

func (k *Keeper) GetCollectorLookupTable(ctx sdk.Context, app_id uint64) (collectorLookup types.CollectorLookup, found bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.CollectorLookupTableMappingKey(app_id)
		value = store.Get(key)
	)

	if value == nil {
		return collectorLookup, false
	}

	k.cdc.MustUnmarshal(value, &collectorLookup)
	return collectorLookup, true
}

func (k *Keeper) GetCollectorLookupByAsset(ctx sdk.Context, app_id, asset_id uint64) (collectorLookup types.CollectorLookup, found bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.CollectorLookupTableMappingKey(app_id)
		value = store.Get(key)
	)

	if value == nil {
		return collectorLookup, false
	}

	k.cdc.MustUnmarshal(value, &collectorLookup)

	var newCollectorLoopkup types.CollectorLookup
	var assetRateInfo types.CollectorLookupTable
	for _, msg := range collectorLookup.AssetrateInfo {
		if msg.CollectorAssetId == asset_id {
			newCollectorLoopkup.AppId = msg.AppId
			assetRateInfo.AppId = msg.AppId
			assetRateInfo.CollectorAssetId = msg.CollectorAssetId
			assetRateInfo.SecondaryAssetId = msg.SecondaryAssetId
			assetRateInfo.SurplusThreshold = msg.SurplusThreshold
			assetRateInfo.DebtThreshold = msg.DebtThreshold
			assetRateInfo.LockerSavingRate = msg.LockerSavingRate
			assetRateInfo.LotSize = msg.LotSize
			assetRateInfo.BidFactor = msg.BidFactor
			newCollectorLoopkup.AssetrateInfo = append(newCollectorLoopkup.AssetrateInfo, assetRateInfo)
		}
	}
	return newCollectorLoopkup, true
}

// set denoms for appId in Collector LookupTable
func (k *Keeper) SetAppToDenomsMapping(ctx sdk.Context, app_id uint64, appToDenom types.AppToDenomsMapping) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.CollectorForDenomKey(app_id)
		value = k.cdc.MustMarshal(&appToDenom)
	)

	store.Set(key, value)
}

// get denoms for appId in Collector LookupTable
func (k *Keeper) GetAppToDenomsMapping(ctx sdk.Context, AppId uint64) (appToDenom types.AppToDenomsMapping, found bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.CollectorForDenomKey(AppId)
		value = store.Get(key)
	)

	if value == nil {
		return appToDenom, false
	}

	k.cdc.MustUnmarshal(value, &appToDenom)

	return appToDenom, true
}

///////////////////////////////22222222222222222222

func (k *Keeper) SetAppIdToAuctionMappingForAsset(ctx sdk.Context, appAssetAuctionData types.HistoricalAuction) {

	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.AppIdToAuctionMappingForAssetKey(appAssetAuctionData.AppId)
		value = k.cdc.MustMarshal(&appAssetAuctionData)
	)

	store.Set(key, value)

}

func (k *Keeper) GetAppIdToAuctionMappingForAsset(ctx sdk.Context, app_id uint64) (appAssetAuctionData types.HistoricalAuction, found bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.AppIdToAuctionMappingForAssetKey(app_id)
		value = store.Get(key)
	)

	if value == nil {
		return appAssetAuctionData, false
	}

	k.cdc.MustUnmarshal(value, &appAssetAuctionData)
	return appAssetAuctionData, true
}

////////////////////////////////////333333333333333

func (k *Keeper) SetAuctionMappingForApp(ctx sdk.Context, records ...types.CollectorAuctionLookupTable) error {
	for _, msg := range records {
		var collectorAuctionLookup types.CollectorAuctionLookupTable
		collectorAuctionLookup.AppId = msg.AppId
		collectorAuctionLookup.AssetIdToAuctionLookup = msg.AssetIdToAuctionLookup
		var (
			store = ctx.KVStore(k.storeKey)
			key   = types.AppIdToAuctionMappingKey(msg.AppId)
			value = k.cdc.MustMarshal(&collectorAuctionLookup)
		)

		store.Set(key, value)

	}
	return nil
}

func (k *Keeper) GetAuctionMappingForApp(ctx sdk.Context, app_id uint64) (collectorAuctionLookupTable types.CollectorAuctionLookupTable, found bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.AppIdToAuctionMappingKey(app_id)
		value = store.Get(key)
	)

	if value == nil {
		return collectorAuctionLookupTable, false
	}

	k.cdc.MustUnmarshal(value, &collectorAuctionLookupTable)
	return collectorAuctionLookupTable, true
}

func (k *Keeper) SetCollectorAuctionLookupTable(ctx sdk.Context, records ...types.CollectorAuctionLookupTable) error {
	for _, msg := range records {

		var appAuction = types.CollectorAuctionLookupTable{
			AppId:                  msg.AppId,
			AssetIdToAuctionLookup: msg.AssetIdToAuctionLookup,
		}
		var (
			store = ctx.KVStore(k.storeKey)
			key   = types.CollectorAuctionLookupKey(msg.AppId)
			value = k.cdc.MustMarshal(&appAuction)
		)

		store.Set(key, value)

	}
	return nil
}

func (k *Keeper) GetCollectorAuctionLookupTable(ctx sdk.Context, app_id uint64) (appIdToAuctionData types.CollectorAuctionLookupTable, found bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.CollectorAuctionLookupKey(app_id)
		value = store.Get(key)
	)

	if value == nil {
		return appIdToAuctionData, false
	}

	k.cdc.MustUnmarshal(value, &appIdToAuctionData)
	return appIdToAuctionData, true
}

///////////////////

func (k *Keeper) SetNetFeeCollectedData(ctx sdk.Context, app_id, asset_id uint64, fee sdk.Int) error {

	collectorData, found := k.GetAppidToAssetCollectorMapping(ctx, app_id)
	if !found {
		return types.ErrorDataDoesNotExists
	}
	var netcollected types.NetFeeCollectedData
	var assetCollected types.AssetIdToFeeCollected
	netcollected.AppId = app_id

	var netcollectedfee sdk.Int
	for _, data := range collectorData.AssetCollector {
		if data.AssetId == asset_id {
			assetCollected.AssetId = asset_id
			netcollectedfee = data.Collector.CollectedClosingFee.Add(data.Collector.CollectedOpeningFee).Add(data.Collector.CollectedStabilityFee).Add(fee)
		}
	}
	assetCollected.NetFeesCollected = netcollectedfee
	netcollected.AssetIdToFeeCollected = append(netcollected.AssetIdToFeeCollected, assetCollected)
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.NetFeeCollectedDataKey(app_id)
		value = k.cdc.MustMarshal(
			&protobuftypes.Int64Value{
				Value: netcollectedfee.Int64(),
			},
		)
	)

	store.Set(key, value)

	return nil
}

func (k *Keeper) GetNetFeeCollectedData(ctx sdk.Context, app_id uint64) (netFeeData types.NetFeeCollectedData, found bool) {

	collectorData, found := k.GetAppidToAssetCollectorMapping(ctx, app_id)
	if !found {
		return netFeeData, false
	}
	var assetCollector types.AssetIdCollectorMappping
	for _, data := range collectorData.AssetCollector {

		assetCollector.AssetId = data.AssetId
		assetCollector.Collector = data.Collector

	}
	collectorData.AssetCollector = append(collectorData.AssetCollector, &assetCollector)

	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.NetFeeCollectedDataKey(app_id)
		value = store.Get(key)
	)

	if value == nil {
		return netFeeData, false
	}

	k.cdc.MustUnmarshal(value, &netFeeData)
	return netFeeData, true
}

func (k *Keeper) WasmSetCollectorLookupTable(ctx sdk.Context, AppId, CollectorAssetId, SecondaryAssetId, SurplusThreshold, DebtThreshold uint64, LockerSavingRate sdk.Dec, LotSize uint64, BidFactor sdk.Dec) error {

	if !k.HasAsset(ctx, CollectorAssetId) {
		return types.ErrorAssetDoesNotExist
	}
	if !k.HasAsset(ctx, SecondaryAssetId) {
		return types.ErrorAssetDoesNotExist
	}
	if CollectorAssetId == SecondaryAssetId {
		return types.ErrorDuplicateAssetDenoms
	}
	appDenom, found := k.GetAppToDenomsMapping(ctx, AppId)
	if found {
		//check if assetdenom already exists
		var check = 0
		for _, data := range appDenom.AssetIds {
			if data == CollectorAssetId {
				check++
			}
		}
		if check > 0 {
			return types.ErrorDuplicateCollectorDenomForApp
		}
		// if denom is new then append
		appDenom.AssetIds = append(appDenom.AssetIds, CollectorAssetId)
		k.SetAppToDenomsMapping(ctx, AppId, appDenom)

	} else {
		//initialize the mappping
		var appDenomNew types.AppToDenomsMapping
		appDenomNew.AppId = AppId
		appDenomNew.AssetIds = append(appDenomNew.AssetIds, CollectorAssetId)
		k.SetAppToDenomsMapping(ctx, AppId, appDenomNew)
	}

	var Collector = types.CollectorLookupTable{
		AppId:            AppId,
		CollectorAssetId: CollectorAssetId,
		SecondaryAssetId: SecondaryAssetId,
		SurplusThreshold: SurplusThreshold,
		DebtThreshold:    DebtThreshold,
		LockerSavingRate: &LockerSavingRate,
		LotSize:          LotSize,
		BidFactor:        &BidFactor,
	}
	accmLookup, _ := k.GetCollectorLookupTable(ctx, AppId)
	accmLookup.AppId = AppId
	accmLookup.AssetrateInfo = append(accmLookup.AssetrateInfo, Collector)

	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.CollectorLookupTableMappingKey(AppId)
		value = k.cdc.MustMarshal(&accmLookup)
	)

	store.Set(key, value)
	return nil
}

func (k *Keeper) WasmSetCollectorLookupTableQuery(ctx sdk.Context, AppId, CollectorAssetId, SecondaryAssetId uint64) error {

	if !k.HasAsset(ctx, CollectorAssetId) {
		return types.ErrorAssetDoesNotExist
	}
	if !k.HasAsset(ctx, SecondaryAssetId) {
		return types.ErrorAssetDoesNotExist
	}
	if CollectorAssetId == SecondaryAssetId {
		return types.ErrorDuplicateAssetDenoms
	}
	appDenom, found := k.GetAppToDenomsMapping(ctx, AppId)
	if found {
		//check if assetdenom already exists
		var check = 0
		for _, data := range appDenom.AssetIds {
			if data == CollectorAssetId {
				check++
			}
		}
		if check > 0 {
			return types.ErrorDuplicateCollectorDenomForApp
		}
	}
	return nil
}

func (k *Keeper) WasmSetAuctionMappingForApp(ctx sdk.Context, AppId uint64, AssetId []uint64, IsSurplusAuction, IsDebtAuction, IsAuctionActive []bool) error {

	var collectorAuctionLookup types.CollectorAuctionLookupTable
	collectorAuctionLookup.AppId = AppId
	var AssetIdToAuctionLookups []types.AssetIdToAuctionLookupTable
	for i := range AssetId {
		AssetIdToAuctionLookup := types.AssetIdToAuctionLookupTable{
			AssetId:          AssetId[i],
			IsSurplusAuction: IsSurplusAuction[i],
			IsDebtAuction:    IsDebtAuction[i],
			IsAuctionActive:  IsAuctionActive[i],
		}
		AssetIdToAuctionLookups = append(AssetIdToAuctionLookups, AssetIdToAuctionLookup)
	}
	collectorAuctionLookup.AssetIdToAuctionLookup = AssetIdToAuctionLookups
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.AppIdToAuctionMappingKey(AppId)
		value = k.cdc.MustMarshal(&collectorAuctionLookup)
	)

	store.Set(key, value)

	return nil
}
