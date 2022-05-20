package keeper

import (
	"github.com/comdex-official/comdex/x/collector/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) UpdateCollector(ctx sdk.Context, appId, asset_id uint64, collector types.CollectorData) error {

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

		newCollector.CollectedClosingFee = collector.CollectedClosingFee
		newCollector.CollectedOpeningFee = collector.CollectedOpeningFee
		newCollector.CollectedStabilityFee = collector.CollectedStabilityFee
		newCollector.LiquidationRewardsCollected = collector.LiquidationRewardsCollected
		newCollector.NetFeesCollected = newCollector.CollectedClosingFee.Add(newCollector.CollectedOpeningFee)
		newCollector.NetFeesCollected = newCollector.NetFeesCollected.Add(newCollector.CollectedStabilityFee)
		assetIdCollect.Collector= &newCollector
		
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

				newCollector.CollectedClosingFee = collector.CollectedClosingFee
				newCollector.CollectedOpeningFee = collector.CollectedOpeningFee
				newCollector.CollectedStabilityFee = collector.CollectedStabilityFee
				newCollector.LiquidationRewardsCollected = collector.LiquidationRewardsCollected
				newCollector.NetFeesCollected = newCollector.CollectedClosingFee.Add(newCollector.CollectedOpeningFee)
				newCollector.NetFeesCollected = newCollector.NetFeesCollected.Add(newCollector.CollectedStabilityFee)
				assetIdCollect.Collector= &newCollector
				
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

			newCollector.CollectedClosingFee = collector.CollectedClosingFee
			newCollector.CollectedOpeningFee = collector.CollectedOpeningFee
			newCollector.CollectedStabilityFee = collector.CollectedStabilityFee
			newCollector.LiquidationRewardsCollected = collector.LiquidationRewardsCollected
			newCollector.NetFeesCollected = newCollector.CollectedClosingFee.Add(newCollector.CollectedOpeningFee)
			newCollector.NetFeesCollected = newCollector.NetFeesCollected.Add(newCollector.CollectedStabilityFee)
			assetIdCollect.Collector= &newCollector
			
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
		if !k.HasAssetForDenom(ctx, msg.CollectorDenom) {
			return types.ErrorAssetDoesNotExist
		}
		if !k.HasAssetForDenom(ctx, msg.SecondaryDenom) {
			return types.ErrorAssetDoesNotExist
		}
		if msg.CollectorDenom == msg.SecondaryDenom {
			return types.ErrorDuplicateAssetDenoms
		}
		appDenom, found := k.GetAppToDenomsMapping(ctx, msg.AppId)
		if found {
			//check if assetdenom already exists
			var check = 0
			for _, data := range appDenom.AssetDenoms {
				if data == msg.CollectorDenom {
					check++
				}
			}
			if check > 0 {
				return types.ErrorDuplicateCollectorDenomForApp
			}
			// if denom is new then append
			appDenom.AssetDenoms = append(appDenom.AssetDenoms, msg.CollectorDenom)
			k.SetAppToDenomsMapping(ctx, msg.AppId, appDenom)

		} else {
			//initialize the mappping
			var appDenomNew types.AppToDenomsMapping
			appDenomNew.AppId = msg.AppId
			appDenomNew.AssetDenoms = append(appDenomNew.AssetDenoms, msg.CollectorDenom)
			k.SetAppToDenomsMapping(ctx, msg.AppId, appDenomNew)
		}
		
			var Collector = types.CollectorLookupTable{
				AppId: msg.AppId,
				CollectorDenom: msg.CollectorDenom,
				SecondaryDenom: msg.SecondaryDenom,
				SurplusThreshold: msg.SurplusThreshold,
				DebtThreshold: msg.DebtThreshold,
				LockerSavingRate: msg.LockerSavingRate,
				LotSize: msg.LotSize,
				BidFactor: msg.BidFactor,
			}
			assetInfo, found := k.GetAssetForDenom(ctx, msg.CollectorDenom)
			accmLookup, found := k.GetCollectorLookupTable(ctx, msg.AppId)
			accmLookup.AssetId = assetInfo.Id
			accmLookup.AssetrateInfo = append(accmLookup.AssetrateInfo, &Collector)
			
		var(
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
