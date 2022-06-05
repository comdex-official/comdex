package keeper

import (
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	"github.com/comdex-official/comdex/x/collector/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protobuftypes "github.com/gogo/protobuf/types"
)

// GetAmountFromCollector returns amount from the collector
func (k *Keeper) GetAmountFromCollector(ctx sdk.Context, appId, assetId uint64, amount sdk.Int) (sdk.Int, error) {
	netFeeData, found := k.GetNetFeeCollectedData(ctx, appId)
	var returnedFee sdk.Int
	if !found {
		return returnedFee, types.ErrorDataDoesNotExists
	}

	for _, data := range netFeeData.AssetIdToFeeCollected {
		if data.AssetId == assetId {
			if !(data.NetFeesCollected.Sub(amount).GT(sdk.ZeroInt())) {
				return returnedFee, types.ErrorRequestedAmtExceedsCollectedFee
			} else {
				asset, _ := k.GetAsset(ctx, assetId)
				if err := k.SendCoinFromModuleToModule(ctx, types.ModuleName, auctiontypes.ModuleName, sdk.NewCoins(sdk.NewCoin(asset.Denom, data.NetFeesCollected.Sub(amount)))); err != nil {
					return returnedFee, err
				}
				err := k.SetNetFeeCollectedData(ctx, appId, assetId, data.NetFeesCollected.Sub(amount))
				if err != nil {
					return sdk.Int{}, err
				}
			}
		}

	}
	return returnedFee, nil
}

// UpdateCollector update collector store
func (k *Keeper) UpdateCollector(ctx sdk.Context, appId, assetId uint64, CollectedStabilityFee, CollectedClosingFee, CollectedOpeningFee, LiquidationRewardsCollected sdk.Int) error {

	if !k.HasAsset(ctx, assetId) {
		return types.ErrorAssetDoesNotExist
	}

	collectorData, found := k.GetAppidToAssetCollectorMapping(ctx, appId)

	if !found {
		//create a new instance of AppId To AssetCollectorMapping
		var collectorNewData types.AppIdToAssetCollectorMapping
		collectorNewData.AppId = appId

		var assetIdCollect types.AssetIdCollectorMapping
		assetIdCollect.AssetId = assetId

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

			if data.AssetId != assetId { //if does not exist then create a new instance
				check++
				var collectorNewData types.AppIdToAssetCollectorMapping
				collectorNewData.AppId = appId

				var assetIdCollect types.AssetIdCollectorMapping
				assetIdCollect.AssetId = assetId

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

			var assetIdCollect types.AssetIdCollectorMapping
			assetIdCollect.AssetId = assetId

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

// SetAppidToAssetCollectorMapping update collector with app_id and asset
func (k *Keeper) SetAppidToAssetCollectorMapping(ctx sdk.Context, appAssetCollectorData types.AppIdToAssetCollectorMapping) {

	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.AppidToAssetCollectorMappingKey(appAssetCollectorData.AppId)
		value = k.cdc.MustMarshal(&appAssetCollectorData)
	)

	store.Set(key, value)

}

// GetAppidToAssetCollectorMapping returns app_id to asset mapping for collector
func (k *Keeper) GetAppidToAssetCollectorMapping(ctx sdk.Context, appId uint64) (appAssetCollectorData types.AppIdToAssetCollectorMapping, found bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.AppidToAssetCollectorMappingKey(appId)
		value = store.Get(key)
	)

	if value == nil {
		return appAssetCollectorData, false
	}

	k.cdc.MustUnmarshal(value, &appAssetCollectorData)
	return appAssetCollectorData, true
}

// SetCollectorLookupTable updates the collector lookup store
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
		accmLookup.AssetRateInfo = append(accmLookup.AssetRateInfo, Collector)

		var (
			store = ctx.KVStore(k.storeKey)
			key   = types.CollectorLookupTableMappingKey(msg.AppId)
			value = k.cdc.MustMarshal(&accmLookup)
		)

		store.Set(key, value)
	}
	return nil

}

// GetCollectorLookupTable returns collector lookup table
func (k *Keeper) GetCollectorLookupTable(ctx sdk.Context, appId uint64) (collectorLookup types.CollectorLookup, found bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.CollectorLookupTableMappingKey(appId)
		value = store.Get(key)
	)

	if value == nil {
		return collectorLookup, false
	}

	k.cdc.MustUnmarshal(value, &collectorLookup)
	return collectorLookup, true
}

// GetCollectorLookupByAsset return collector lookup data queried on asset
func (k *Keeper) GetCollectorLookupByAsset(ctx sdk.Context, appId, assetId uint64) (collectorLookupTable types.CollectorLookupTable, found bool) {
	collectorLookup, found := k.GetCollectorLookupTable(ctx, appId)
	if !found {
		return collectorLookupTable, false
	}

	var assetRateInfo types.CollectorLookupTable
	for _, data := range collectorLookup.AssetRateInfo {
		if data.CollectorAssetId == assetId {
			assetRateInfo = data
		}
	}
	return assetRateInfo, true
}

// SetAppToDenomsMapping set denoms for appId in Collector LookupTable
func (k *Keeper) SetAppToDenomsMapping(ctx sdk.Context, appId uint64, appToDenom types.AppToDenomsMapping) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.CollectorForDenomKey(appId)
		value = k.cdc.MustMarshal(&appToDenom)
	)

	store.Set(key, value)
}

// GetAppToDenomsMapping get denoms for appId in Collector LookupTable
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

// SetAppIdToAuctionMappingForAsset updates app_id for asset in auction data map
func (k *Keeper) SetAppIdToAuctionMappingForAsset(ctx sdk.Context, appAssetAuctionData types.HistoricalAuction) {

	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.AppIdToAuctionMappingForAssetKey(appAssetAuctionData.AppId)
		value = k.cdc.MustMarshal(&appAssetAuctionData)
	)

	store.Set(key, value)

}

// GetAppIdToAuctionMappingForAsset gets app_id for asset in auction data map
func (k *Keeper) GetAppIdToAuctionMappingForAsset(ctx sdk.Context, appId uint64) (appAssetAuctionData types.HistoricalAuction, found bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.AppIdToAuctionMappingForAssetKey(appId)
		value = store.Get(key)
	)

	if value == nil {
		return appAssetAuctionData, false
	}

	k.cdc.MustUnmarshal(value, &appAssetAuctionData)
	return appAssetAuctionData, true
}

// SetAuctionMappingForApp sets auction map data for app/product
func (k *Keeper) SetAuctionMappingForApp(ctx sdk.Context, records ...types.CollectorAuctionLookupTable) error {
	for _, msg := range records {
		_, found := k.GetApp(ctx, msg.AppId)
		if !found {
			return types.ErrorAppDoesNotExist
		}
		var collectorAuctionLookup types.CollectorAuctionLookupTable
		collectorAuctionLookup.AppId = msg.AppId
		collectorAuctionLookup.AssetIdToAuctionLookup = msg.AssetIdToAuctionLookup
		for _, data := range collectorAuctionLookup.AssetIdToAuctionLookup {
			_, found := k.GetAsset(ctx, data.AssetId)
			if !found {
				return types.ErrorAssetDoesNotExist
			}
		}
		var (
			store = ctx.KVStore(k.storeKey)
			key   = types.AppIdToAuctionMappingKey(msg.AppId)
			value = k.cdc.MustMarshal(&collectorAuctionLookup)
		)

		store.Set(key, value)

	}
	return nil
}

// GetAuctionMappingForApp gets auction map data for app/product
func (k *Keeper) GetAuctionMappingForApp(ctx sdk.Context, appId uint64) (collectorAuctionLookupTable types.CollectorAuctionLookupTable, found bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.AppIdToAuctionMappingKey(appId)
		value = store.Get(key)
	)

	if value == nil {
		return collectorAuctionLookupTable, false
	}

	k.cdc.MustUnmarshal(value, &collectorAuctionLookupTable)
	return collectorAuctionLookupTable, true
}

// SetCollectorAuctionLookupTable sets collector lookup auction table
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

// GetCollectorAuctionLookupTable gets collector lookup auction table
func (k *Keeper) GetCollectorAuctionLookupTable(ctx sdk.Context, appId uint64) (appIdToAuctionData types.CollectorAuctionLookupTable, found bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.CollectorAuctionLookupKey(appId)
		value = store.Get(key)
	)

	if value == nil {
		return appIdToAuctionData, false
	}

	k.cdc.MustUnmarshal(value, &appIdToAuctionData)
	return appIdToAuctionData, true
}

// SetNetFeeCollectedData sets net fees collected
func (k *Keeper) SetNetFeeCollectedData(ctx sdk.Context, appId, assetId uint64, fee sdk.Int) error {

	collectorData, found := k.GetAppidToAssetCollectorMapping(ctx, appId)
	if !found {
		return types.ErrorDataDoesNotExists
	}
	var netCollected types.NetFeeCollectedData
	var assetCollected types.AssetIdToFeeCollected
	netCollected.AppId = appId

	var netCollectedFee sdk.Int
	for _, data := range collectorData.AssetCollector {
		if data.AssetId == assetId {
			assetCollected.AssetId = assetId
			netCollectedFee = data.Collector.CollectedClosingFee.Add(data.Collector.CollectedOpeningFee).Add(data.Collector.CollectedStabilityFee).Add(fee)
		}
	}
	assetCollected.NetFeesCollected = netCollectedFee
	netCollected.AssetIdToFeeCollected = append(netCollected.AssetIdToFeeCollected, assetCollected)
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.NetFeeCollectedDataKey(appId)
		value = k.cdc.MustMarshal(
			&protobuftypes.Int64Value{
				Value: netCollectedFee.Int64(),
			},
		)
	)

	store.Set(key, value)

	return nil
}

// GetNetFeeCollectedData sets net fees collected
func (k *Keeper) GetNetFeeCollectedData(ctx sdk.Context, appId uint64) (netFeeData types.NetFeeCollectedData, found bool) {

	collectorData, found := k.GetAppidToAssetCollectorMapping(ctx, appId)
	if !found {
		return netFeeData, false
	}
	var assetCollector types.AssetIdCollectorMapping
	for _, data := range collectorData.AssetCollector {

		assetCollector.AssetId = data.AssetId
		assetCollector.Collector = data.Collector

	}
	collectorData.AssetCollector = append(collectorData.AssetCollector, &assetCollector)

	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.NetFeeCollectedDataKey(appId)
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
		LockerSavingRate: LockerSavingRate,
		LotSize:          LotSize,
		BidFactor:        BidFactor,
	}
	accmLookup, _ := k.GetCollectorLookupTable(ctx, AppId)
	accmLookup.AppId = AppId
	accmLookup.AssetRateInfo = append(accmLookup.AssetRateInfo, Collector)

	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.CollectorLookupTableMappingKey(AppId)
		value = k.cdc.MustMarshal(&accmLookup)
	)

	store.Set(key, value)
	return nil
}

func (k *Keeper) WasmSetCollectorLookupTableQuery(ctx sdk.Context, AppId, CollectorAssetId, SecondaryAssetId uint64) (bool, string) {

	if !k.HasAsset(ctx, CollectorAssetId) {
		return false, types.ErrorAssetDoesNotExist.Error()
	}
	if !k.HasAsset(ctx, SecondaryAssetId) {
		return false, types.ErrorAssetDoesNotExist.Error()
	}
	if CollectorAssetId == SecondaryAssetId {
		return false, types.ErrorDuplicateAssetDenoms.Error()
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
			return false, types.ErrorDuplicateCollectorDenomForApp.Error()
		}
	}
	return true, ""
}

func (k *Keeper) WasmSetAuctionMappingForApp(ctx sdk.Context, AppId uint64, AssetId []uint64, IsSurplusAuction, IsDebtAuction, AssetOutOraclePrice []bool, AssetOutPrice []uint64) error {

	var collectorAuctionLookup types.CollectorAuctionLookupTable
	collectorAuctionLookup.AppId = AppId
	var AssetIdToAuctionLookups []types.AssetIdToAuctionLookupTable
	for i := range AssetId {
		AssetIdToAuctionLookup := types.AssetIdToAuctionLookupTable{
			AssetId:             AssetId[i],
			IsSurplusAuction:    IsSurplusAuction[i],
			IsDebtAuction:       IsDebtAuction[i],
			AssetOutOraclePrice: AssetOutOraclePrice[i],
			AssetOutPrice:       AssetOutPrice[i],
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

func (k *Keeper) WasmSetAuctionMappingForAppQuery(ctx sdk.Context, AppId uint64) (bool, string) {

	_, _ = k.GetAppidToAssetCollectorMapping(ctx, AppId)

	return true, ""
}

func (k *Keeper) WasmUpdateLsrInCollectorLookupTable(ctx sdk.Context, appId, assetId uint64, lsr sdk.Dec) error {

	var Collector types.CollectorLookupTable
	accmLookup, _ := k.GetCollectorLookupTable(ctx, appId)

	for _, data := range accmLookup.AssetRateInfo {
		if data.CollectorAssetId == assetId {
			Collector.CollectorAssetId = assetId
			Collector.AppId = data.AppId
			Collector.BidFactor = data.BidFactor
			Collector.DebtThreshold = data.DebtThreshold
			Collector.LockerSavingRate = lsr
			Collector.LotSize = data.LotSize
			Collector.SecondaryAssetId = data.SecondaryAssetId
			Collector.DebtThreshold = data.DebtThreshold
		}
	}
	k.SetCollectorLookupTable(ctx, Collector)
	// var (
	// 	store = ctx.KVStore(k.storeKey)
	// 	key   = types.CollectorLookupTableMappingKey(appId)
	// 	value = k.cdc.MustMarshal(&accmLookup)
	// )

	// store.Set(key, value)
	return nil
}

func (k *Keeper) WasmUpdateLsrInCollectorLookupTableQuery(ctx sdk.Context, appId, assetId uint64) (bool, string) {
	_, found := k.GetCollectorLookupByAsset(ctx, appId, assetId)
	if !found {
		return false, types.ErrorDataDoesNotExists.Error()
	}
	return true, ""
}
