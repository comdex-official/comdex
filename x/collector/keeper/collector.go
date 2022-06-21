package keeper

import (
	"github.com/comdex-official/comdex/app/wasm/bindings"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	"github.com/comdex-official/comdex/x/collector/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetAmountFromCollector returns amount from the collector
func (k *Keeper) GetAmountFromCollector(ctx sdk.Context, appID, assetID uint64, amount sdk.Int) (sdk.Int, error) {
	netFeeData, found := k.GetNetFeeCollectedData(ctx, appID)
	var returnedFee sdk.Int
	if !found {
		return returnedFee, types.ErrorDataDoesNotExists
	}

	for _, data := range netFeeData.AssetIdToFeeCollected {
		if data.AssetId == assetID {
			if !(data.NetFeesCollected.Sub(amount).GT(sdk.ZeroInt())) {
				return returnedFee, types.ErrorRequestedAmtExceedsCollectedFee
			}
			asset, _ := k.GetAsset(ctx, assetID)
			err := k.SendCoinFromModuleToModule(ctx, types.ModuleName, auctiontypes.ModuleName, sdk.NewCoins(sdk.NewCoin(asset.Denom, amount)))
			if err != nil {
				return returnedFee, err
			}
			err = k.DecreaseNetFeeCollectedData(ctx, appID, assetID, amount)
			if err != nil {
				return sdk.Int{}, err
			}
		}
	}
	return returnedFee, nil
}

func (k *Keeper) DecreaseNetFeeCollectedData(ctx sdk.Context, appID, assetID uint64, amount sdk.Int) error {
	collectorData, found := k.GetNetFeeCollectedData(ctx, appID)
	if !found {
		return types.ErrorDataDoesNotExists
	}
	var netCollected types.NetFeeCollectedData
	var assetCollected types.AssetIdToFeeCollected
	netCollected.AppId = appID

	var netCollectedFee sdk.Int
	for _, data := range collectorData.AssetIdToFeeCollected {
		if data.AssetId == assetID {
			assetCollected.AssetId = assetID
			netCollectedFee = data.NetFeesCollected.Sub(amount)
		}
	}
	assetCollected.NetFeesCollected = netCollectedFee
	netCollected.AssetIdToFeeCollected = append(netCollected.AssetIdToFeeCollected, assetCollected)
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.NetFeeCollectedDataKey(appID)
		value = k.cdc.MustMarshal(&netCollected)
	)

	store.Set(key, value)

	return nil
}

// UpdateCollector update collector store
func (k *Keeper) UpdateCollector(ctx sdk.Context, appID, assetID uint64, collectedStabilityFee, collectedClosingFee, collectedOpeningFee, liquidationRewardsCollected sdk.Int) error {
	if !k.HasAsset(ctx, assetID) {
		return types.ErrorAssetDoesNotExist
	}

	collectorData, found := k.GetAppidToAssetCollectorMapping(ctx, appID)
	if !found {
		//create a new instance of AppId To AssetCollectorMapping
		var collectorNewData types.AppIdToAssetCollectorMapping
		collectorNewData.AppId = appID

		var assetIDCollect types.AssetIdCollectorMapping
		assetIDCollect.AssetId = assetID

		var newCollector types.CollectorData
		newCollector.CollectedClosingFee = collectedClosingFee
		newCollector.CollectedOpeningFee = collectedOpeningFee
		newCollector.CollectedStabilityFee = collectedStabilityFee
		newCollector.LiquidationRewardsCollected = liquidationRewardsCollected
		assetIDCollect.Collector = newCollector
		collectorNewData.AssetCollector = append(collectorNewData.AssetCollector, assetIDCollect)

		k.SetAppidToAssetCollectorMapping(ctx, collectorNewData)
		err := k.SetNetFeeCollectedData(ctx, appID, assetID,
			newCollector.CollectedClosingFee.
				Add(newCollector.CollectedOpeningFee).
				Add(newCollector.CollectedStabilityFee).
				Add(newCollector.LiquidationRewardsCollected))
		if err != nil {
			return err
		}
	} else {
		var check = 0 // makes it 1 if assetID exists for appId
		for _, data := range collectorData.AssetCollector {
			if data.AssetId == assetID {
				check++
				var collectorNewData types.AppIdToAssetCollectorMapping
				collectorNewData.AppId = appID

				var assetIDCollect types.AssetIdCollectorMapping
				assetIDCollect.AssetId = assetID

				var newCollector types.CollectorData
				newCollector.CollectedClosingFee = data.Collector.CollectedClosingFee.Add(collectedClosingFee)
				newCollector.CollectedOpeningFee = data.Collector.CollectedOpeningFee.Add(collectedOpeningFee)
				newCollector.CollectedStabilityFee = data.Collector.CollectedStabilityFee.Add(collectedStabilityFee)
				newCollector.LiquidationRewardsCollected = sdk.ZeroInt()
				newCollector.LiquidationRewardsCollected = data.Collector.LiquidationRewardsCollected.Add(newCollector.LiquidationRewardsCollected)
				assetIDCollect.Collector = newCollector

				collectorNewData.AssetCollector = append(collectorNewData.AssetCollector, assetIDCollect)
				k.SetAppidToAssetCollectorMapping(ctx, collectorNewData)
				err := k.SetNetFeeCollectedData(ctx, appID, assetID,
					collectedClosingFee.
						Add(collectedOpeningFee).
						Add(collectedStabilityFee).
						Add(liquidationRewardsCollected))
				if err != nil {
					return err
				}

				return nil
			}
		}

		if check == 0 {
			var collectorNewData types.AppIdToAssetCollectorMapping
			collectorNewData.AppId = appID

			var assetIDCollect types.AssetIdCollectorMapping
			assetIDCollect.AssetId = assetID
			var newCollector types.CollectorData

			newCollector.CollectedClosingFee = collectedClosingFee
			newCollector.CollectedOpeningFee = collectedOpeningFee
			newCollector.CollectedStabilityFee = collectedStabilityFee
			newCollector.LiquidationRewardsCollected = liquidationRewardsCollected
			assetIDCollect.Collector = newCollector

			collectorNewData.AssetCollector = append(collectorNewData.AssetCollector, assetIDCollect)

			k.SetAppidToAssetCollectorMapping(ctx, collectorNewData)
			err := k.SetNetFeeCollectedData(ctx, appID, assetID,
				newCollector.CollectedClosingFee.
					Add(newCollector.CollectedOpeningFee).
					Add(newCollector.CollectedStabilityFee).
					Add(newCollector.LiquidationRewardsCollected))
			if err != nil {
				return err
			}
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
func (k *Keeper) GetAppidToAssetCollectorMapping(ctx sdk.Context, appID uint64) (appAssetCollectorData types.AppIdToAssetCollectorMapping, found bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.AppidToAssetCollectorMappingKey(appID)
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
		_, found := k.GetMintGenesisTokenData(ctx, msg.AppId, msg.SecondaryAssetId)
		if !found {
			return types.ErrorAssetNotAddedForGenesisMinting
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
			DebtLotSize:      msg.DebtLotSize,
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

func (k *Keeper) SetCollectorLookupTableForWasm(ctx sdk.Context, records ...types.CollectorLookupTable) error {
	for _, msg := range records {
		accmLookup, _ := k.GetCollectorLookupTable(ctx, msg.AppId)
		accmLookup.AppId = msg.AppId
		aa := accmLookup.AssetRateInfo
		for j, v := range aa {
			if v.CollectorAssetId == msg.CollectorAssetId {
				v.LockerSavingRate = msg.LockerSavingRate
				accmLookup.AssetRateInfo[j] = v
				var (
					store = ctx.KVStore(k.storeKey)
					key   = types.CollectorLookupTableMappingKey(msg.AppId)
					value = k.cdc.MustMarshal(&accmLookup)
				)
				store.Set(key, value)
			}
		}
	}
	return nil

}

// GetCollectorLookupTable returns collector lookup table
func (k *Keeper) GetCollectorLookupTable(ctx sdk.Context, appID uint64) (collectorLookup types.CollectorLookup, found bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.CollectorLookupTableMappingKey(appID)
		value = store.Get(key)
	)

	if value == nil {
		return collectorLookup, false
	}

	k.cdc.MustUnmarshal(value, &collectorLookup)
	return collectorLookup, true
}

// GetCollectorLookupByAsset return collector lookup data queried on asset
func (k *Keeper) GetCollectorLookupByAsset(ctx sdk.Context, appID, assetID uint64) (collectorLookupTable types.CollectorLookupTable, found bool) {
	collectorLookup, found := k.GetCollectorLookupTable(ctx, appID)
	if !found {
		return collectorLookupTable, false
	}

	var assetRateInfo types.CollectorLookupTable
	for _, data := range collectorLookup.AssetRateInfo {
		if data.CollectorAssetId == assetID {
			assetRateInfo = data
		}
	}
	return assetRateInfo, true
}

// SetAppToDenomsMapping set denoms for appId in Collector LookupTable
func (k *Keeper) SetAppToDenomsMapping(ctx sdk.Context, appID uint64, appToDenom types.AppToDenomsMapping) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.CollectorForDenomKey(appID)
		value = k.cdc.MustMarshal(&appToDenom)
	)

	store.Set(key, value)
}

// GetAppToDenomsMapping get denoms for appId in Collector LookupTable
func (k *Keeper) GetAppToDenomsMapping(ctx sdk.Context, appID uint64) (appToDenom types.AppToDenomsMapping, found bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.CollectorForDenomKey(appID)
		value = store.Get(key)
	)

	if value == nil {
		return appToDenom, false
	}

	k.cdc.MustUnmarshal(value, &appToDenom)

	return appToDenom, true
}

// SetAuctionMappingForApp sets auction map data for app/product
func (k *Keeper) SetAuctionMappingForApp(ctx sdk.Context, records ...types.CollectorAuctionLookupTable) error {
	for _, msg := range records {
		_, found := k.GetApp(ctx, msg.AppId)
		if !found {
			return types.ErrorAppDoesNotExist
		}
		_, found1 := k.GetAuctionParams(ctx, msg.AppId)
		if !found1 {
			return types.ErrorAuctionParmsNotSet
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
			key   = types.AppIDToAuctionMappingKey(msg.AppId)
			value = k.cdc.MustMarshal(&collectorAuctionLookup)
		)

		store.Set(key, value)
	}
	return nil
}

// GetAuctionMappingForApp gets auction map data for app/product
func (k *Keeper) GetAuctionMappingForApp(ctx sdk.Context, appID uint64) (collectorAuctionLookupTable types.CollectorAuctionLookupTable, found bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.AppIDToAuctionMappingKey(appID)
		value = store.Get(key)
	)

	if value == nil {
		return collectorAuctionLookupTable, false
	}

	k.cdc.MustUnmarshal(value, &collectorAuctionLookupTable)
	return collectorAuctionLookupTable, true
}

func (k *Keeper) GetAllAuctionMappingForApp(ctx sdk.Context) (collectorAuctionLookupTable []types.CollectorAuctionLookupTable, found bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		iter  = sdk.KVStorePrefixIterator(store, types.AppIDToAuctionMappingPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var table types.CollectorAuctionLookupTable
		k.cdc.MustUnmarshal(iter.Value(), &table)
		collectorAuctionLookupTable = append(collectorAuctionLookupTable, table)
	}
	if collectorAuctionLookupTable == nil {
		return nil, false
	}

	return collectorAuctionLookupTable, true
}

func (k *Keeper) SetNetFeeCollectedData(ctx sdk.Context, appID, assetID uint64, fee sdk.Int) error {
	collectorData, found := k.GetNetFeeCollectedData(ctx, appID)
	if !found {
		var netCollected types.NetFeeCollectedData
		var assetCollected types.AssetIdToFeeCollected
		netCollected.AppId = appID

		var netCollectedFee sdk.Int

		assetCollected.AssetId = assetID
		netCollectedFee = fee

		assetCollected.NetFeesCollected = netCollectedFee
		netCollected.AssetIdToFeeCollected = append(netCollected.AssetIdToFeeCollected, assetCollected)
		var (
			store = ctx.KVStore(k.storeKey)
			key   = types.NetFeeCollectedDataKey(appID)
			value = k.cdc.MustMarshal(&netCollected)
		)

		store.Set(key, value)
	} else {
		var netCollected types.NetFeeCollectedData
		var assetCollected types.AssetIdToFeeCollected
		netCollected.AppId = appID

		var netCollectedFee sdk.Int
		for _, data := range collectorData.AssetIdToFeeCollected {
			if data.AssetId == assetID {
				assetCollected.AssetId = assetID
				netCollectedFee = data.NetFeesCollected.Add(fee)
			}
		}
		assetCollected.NetFeesCollected = netCollectedFee
		netCollected.AssetIdToFeeCollected = append(netCollected.AssetIdToFeeCollected, assetCollected)
		var (
			store = ctx.KVStore(k.storeKey)
			key   = types.NetFeeCollectedDataKey(appID)
			value = k.cdc.MustMarshal(&netCollected)
		)

		store.Set(key, value)
	}

	return nil
}

// GetNetFeeCollectedData sets net fees collected
func (k *Keeper) GetNetFeeCollectedData(ctx sdk.Context, appID uint64) (netFeeData types.NetFeeCollectedData, found bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.NetFeeCollectedDataKey(appID)
		value = store.Get(key)
	)

	if value == nil {
		return netFeeData, false
	}

	k.cdc.MustUnmarshal(value, &netFeeData)
	return netFeeData, true
}

func (k *Keeper) WasmSetCollectorLookupTable(ctx sdk.Context, collectorBindings *bindings.MsgSetCollectorLookupTable) error {
	if !k.HasAsset(ctx, collectorBindings.CollectorAssetId) {
		return types.ErrorAssetDoesNotExist
	}
	if !k.HasAsset(ctx, collectorBindings.SecondaryAssetId) {
		return types.ErrorAssetDoesNotExist
	}
	if collectorBindings.CollectorAssetId == collectorBindings.SecondaryAssetId {
		return types.ErrorDuplicateAssetDenoms
	}
	appDenom, found := k.GetAppToDenomsMapping(ctx, collectorBindings.AppMappingId)
	if found {
		//check if asset denom already exists
		var check = 0
		for _, data := range appDenom.AssetIds {
			if data == collectorBindings.CollectorAssetId {
				check++
			}
		}
		if check > 0 {
			return types.ErrorDuplicateCollectorDenomForApp
		}
		// if denom is new then append
		appDenom.AssetIds = append(appDenom.AssetIds, collectorBindings.CollectorAssetId)
		k.SetAppToDenomsMapping(ctx, collectorBindings.AppMappingId, appDenom)
	} else {
		//initialize the mapping
		var appDenomNew types.AppToDenomsMapping
		appDenomNew.AppId = collectorBindings.AppMappingId
		appDenomNew.AssetIds = append(appDenomNew.AssetIds, collectorBindings.CollectorAssetId)
		k.SetAppToDenomsMapping(ctx, collectorBindings.AppMappingId, appDenomNew)
	}

	var Collector = types.CollectorLookupTable{
		AppId:            collectorBindings.AppMappingId,
		CollectorAssetId: collectorBindings.CollectorAssetId,
		SecondaryAssetId: collectorBindings.SecondaryAssetId,
		SurplusThreshold: collectorBindings.SurplusThreshold,
		DebtThreshold:    collectorBindings.DebtThreshold,
		LockerSavingRate: collectorBindings.LockerSavingRate,
		LotSize:          collectorBindings.LotSize,
		BidFactor:        collectorBindings.BidFactor,
		DebtLotSize:      collectorBindings.DebtLotSize,
	}
	accmLookup, _ := k.GetCollectorLookupTable(ctx, collectorBindings.AppMappingId)
	accmLookup.AppId = collectorBindings.AppMappingId
	accmLookup.AssetRateInfo = append(accmLookup.AssetRateInfo, Collector)

	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.CollectorLookupTableMappingKey(collectorBindings.AppMappingId)
		value = k.cdc.MustMarshal(&accmLookup)
	)

	store.Set(key, value)
	return nil
}

func (k *Keeper) WasmSetCollectorLookupTableQuery(ctx sdk.Context, appID, collectorAssetID, secondaryAssetID uint64) (bool, string) {
	if !k.HasAsset(ctx, collectorAssetID) {
		return false, types.ErrorAssetDoesNotExist.Error()
	}
	if !k.HasAsset(ctx, secondaryAssetID) {
		return false, types.ErrorAssetDoesNotExist.Error()
	}
	if collectorAssetID == secondaryAssetID {
		return false, types.ErrorDuplicateAssetDenoms.Error()
	}
	appDenom, found := k.GetAppToDenomsMapping(ctx, appID)
	if found {
		//check if asset denom already exists
		var check = 0
		for _, data := range appDenom.AssetIds {
			if data == collectorAssetID {
				check++
			}
		}
		if check > 0 {
			return false, types.ErrorDuplicateCollectorDenomForApp.Error()
		}
	}
	return true, ""
}

func (k *Keeper) WasmSetAuctionMappingForApp(ctx sdk.Context, auctionMappingBinding *bindings.MsgSetAuctionMappingForApp) error {
	var collectorAuctionLookup types.CollectorAuctionLookupTable
	collectorAuctionLookup.AppId = auctionMappingBinding.AppMappingId
	var assetIDToAuctionLookups []types.AssetIdToAuctionLookupTable
	for i := range auctionMappingBinding.AssetId {
		assetIDToAuctionLookup := types.AssetIdToAuctionLookupTable{
			AssetId:             auctionMappingBinding.AssetId[i],
			IsSurplusAuction:    auctionMappingBinding.IsSurplusAuction[i],
			IsDebtAuction:       auctionMappingBinding.IsDebtAuction[i],
			AssetOutOraclePrice: auctionMappingBinding.AssetOutOraclePrice[i],
			AssetOutPrice:       auctionMappingBinding.AssetOutPrice[i],
		}
		assetIDToAuctionLookups = append(assetIDToAuctionLookups, assetIDToAuctionLookup)
	}
	collectorAuctionLookup.AssetIdToAuctionLookup = assetIDToAuctionLookups
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.AppIDToAuctionMappingKey(auctionMappingBinding.AppMappingId)
		value = k.cdc.MustMarshal(&collectorAuctionLookup)
	)

	store.Set(key, value)

	return nil
}

func (k *Keeper) WasmSetAuctionMappingForAppQuery(ctx sdk.Context, appID uint64) (bool, string) {
	_, _ = k.GetAppidToAssetCollectorMapping(ctx, appID)

	return true, ""
}

func (k *Keeper) WasmUpdateCollectorLookupTable(ctx sdk.Context, updateLsrInColBinding *bindings.MsgUpdateLsrInCollectorLookupTable) error {
	var Collector types.CollectorLookupTable
	accmLookup, _ := k.GetCollectorLookupTable(ctx, updateLsrInColBinding.AppMappingId)

	for _, data := range accmLookup.AssetRateInfo {
		if data.CollectorAssetId == updateLsrInColBinding.AssetId {
			Collector.CollectorAssetId = updateLsrInColBinding.AssetId
			Collector.AppId = data.AppId
			Collector.BidFactor = updateLsrInColBinding.BidFactor
			Collector.DebtThreshold = updateLsrInColBinding.DebtThreshold
			Collector.SurplusThreshold = updateLsrInColBinding.SurplusThreshold
			Collector.LockerSavingRate = updateLsrInColBinding.LSR
			Collector.LotSize = updateLsrInColBinding.LotSize
			Collector.SecondaryAssetId = data.SecondaryAssetId
			Collector.DebtLotSize = updateLsrInColBinding.DebtLotSize
		}
	}
	err := k.SetCollectorLookupTableForWasm(ctx, Collector)
	if err != nil {
		return err
	}
	return nil
}

func (k *Keeper) WasmUpdateCollectorLookupTableQuery(ctx sdk.Context, appID, assetID uint64) (bool, string) {
	_, found := k.GetCollectorLookupByAsset(ctx, appID, assetID)
	if !found {
		return false, types.ErrorDataDoesNotExists.Error()
	}
	return true, ""
}
