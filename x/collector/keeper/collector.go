package keeper

import (
	"github.com/comdex-official/comdex/app/wasm/bindings"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	"github.com/comdex-official/comdex/x/collector/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetAmountFromCollector returns amount from the collector.
func (k Keeper) GetAmountFromCollector(ctx sdk.Context, appID, assetID uint64, amount sdk.Int) (sdk.Int, error) {
	netFeeData, found := k.GetNetFeeCollectedData(ctx, appID)
	var returnedFee sdk.Int
	if !found {
		return returnedFee, types.ErrorDataDoesNotExists
	}
	if amount.IsNegative() {
		return returnedFee, types.ErrorAmountCanNotBeNegative
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
			err = k.DecreaseNetFeeCollectedData(ctx, appID, assetID, amount, netFeeData)
			if err != nil {
				return sdk.Int{}, err
			}
		}
	}
	returnedFee = amount
	return returnedFee, nil
}

func (k Keeper) DecreaseNetFeeCollectedData(ctx sdk.Context, appID, assetID uint64, amount sdk.Int, collectorData types.NetFeeCollectedData) error {
	var netCollected types.NetFeeCollectedData
	var assetCollected types.AssetIdToFeeCollected
	netCollected.AppId = appID

	var netCollectedFee sdk.Int
	for _, data := range collectorData.AssetIdToFeeCollected {
		if data.AssetId == assetID {
			assetCollected.AssetId = assetID
			netCollectedFee = data.NetFeesCollected.Sub(amount)
			if netCollectedFee.IsNegative() {
				return types.ErrorNetFeesCanNotBeNegative
			}
			assetCollected.NetFeesCollected = netCollectedFee
			netCollected.AssetIdToFeeCollected = append(netCollected.AssetIdToFeeCollected, assetCollected)
		}
	}

	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.NetFeeCollectedDataKey(appID)
		value = k.cdc.MustMarshal(&netCollected)
	)

	store.Set(key, value)

	return nil
}

// UpdateCollector update collector store.
func (k Keeper) UpdateCollector(ctx sdk.Context, appID, assetID uint64, collectedStabilityFee, collectedClosingFee, collectedOpeningFee, liquidationRewardsCollected sdk.Int) error {
	if !k.HasAsset(ctx, assetID) {
		return types.ErrorAssetDoesNotExist
	}

	collectorData, found := k.GetAppidToAssetCollectorMapping(ctx, appID)
	if !found {
		//create a new instance of appID To AssetCollectorMapping
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
				newCollector.LiquidationRewardsCollected = data.Collector.LiquidationRewardsCollected.Add(liquidationRewardsCollected)
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
			var assetIDCollect types.AssetIdCollectorMapping
			assetIDCollect.AssetId = assetID
			var newCollector types.CollectorData

			newCollector.CollectedClosingFee = collectedClosingFee
			newCollector.CollectedOpeningFee = collectedOpeningFee
			newCollector.CollectedStabilityFee = collectedStabilityFee
			newCollector.LiquidationRewardsCollected = liquidationRewardsCollected
			assetIDCollect.Collector = newCollector

			collectorData.AssetCollector = append(collectorData.AssetCollector, assetIDCollect)

			k.SetAppidToAssetCollectorMapping(ctx, collectorData)
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

// SetAppidToAssetCollectorMapping update collector with app_id and asset.
func (k Keeper) SetAppidToAssetCollectorMapping(ctx sdk.Context, appAssetCollectorData types.AppIdToAssetCollectorMapping) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.AppidToAssetCollectorMappingKey(appAssetCollectorData.AppId)
		value = k.cdc.MustMarshal(&appAssetCollectorData)
	)
	store.Set(key, value)
}

// GetAppidToAssetCollectorMapping returns app_id to asset mapping for collector.
func (k Keeper) GetAppidToAssetCollectorMapping(ctx sdk.Context, appID uint64) (appAssetCollectorData types.AppIdToAssetCollectorMapping, found bool) {
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

func (k Keeper) GetAllAppidToAssetCollectorMapping(ctx sdk.Context) (appIdToAssetCollectorMapping []types.AppIdToAssetCollectorMapping) {
	var (
		store = ctx.KVStore(k.storeKey)
		iter  = sdk.KVStorePrefixIterator(store, types.AppIDToAssetCollectorMappingPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var fee types.AppIdToAssetCollectorMapping
		k.cdc.MustUnmarshal(iter.Value(), &fee)
		appIdToAssetCollectorMapping = append(appIdToAssetCollectorMapping, fee)
	}
	return appIdToAssetCollectorMapping
}

// GetCollectorDataForAppIDAssetID returns app_id to asset mapping for collector.
func (k Keeper) GetCollectorDataForAppIDAssetID(ctx sdk.Context, appID uint64, assetID uint64) (collectorData types.CollectorData, found bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.AppidToAssetCollectorMappingKey(appID)
		value = store.Get(key)
	)
	var appAssetCollectorData types.AppIdToAssetCollectorMapping
	if value == nil {
		return collectorData, false
	}
	k.cdc.MustUnmarshal(value, &appAssetCollectorData)

	for _, data := range appAssetCollectorData.AssetCollector {
		if data.AssetId == assetID {
			collectorData = data.Collector
			return collectorData, true
		}
	}
	return collectorData, false
}

// SetCollectorLookupTable updates the collector lookup store.
func (k Keeper) SetCollectorLookupTable(ctx sdk.Context, records ...types.CollectorLookupTable) error {
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
			//initialize the mapping
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

func (k Keeper) SetCollectorLookupTableForWasm(ctx sdk.Context, records ...types.CollectorLookupTable) error {
	for _, msg := range records {
		accmLookup, _ := k.GetCollectorLookupTable(ctx, msg.AppId)
		accmLookup.AppId = msg.AppId
		aa := accmLookup.AssetRateInfo
		for j, v := range aa {
			if v.CollectorAssetId == msg.CollectorAssetId {
				v.LockerSavingRate = msg.LockerSavingRate
				v.BidFactor = msg.BidFactor
				v.DebtLotSize = msg.DebtLotSize
				v.DebtThreshold = msg.DebtThreshold
				v.LotSize = msg.LotSize
				v.SurplusThreshold = msg.SurplusThreshold
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

// GetCollectorLookupTable returns collector lookup table.
func (k Keeper) GetCollectorLookupTable(ctx sdk.Context, appID uint64) (collectorLookup types.CollectorLookup, found bool) {
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

func (k Keeper) GetAllCollectorLookupTable(ctx sdk.Context) (collectorLookup []types.CollectorLookup) {
	var (
		store = ctx.KVStore(k.storeKey)
		iter  = sdk.KVStorePrefixIterator(store, types.AddCollectorLookupKey)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var fee types.CollectorLookup
		k.cdc.MustUnmarshal(iter.Value(), &fee)
		collectorLookup = append(collectorLookup, fee)
	}
	return collectorLookup
}

// GetCollectorLookupByAsset return collector lookup data queried on asset.
func (k Keeper) GetCollectorLookupByAsset(ctx sdk.Context, appID, assetID uint64) (collectorLookupTable types.CollectorLookupTable, found bool) {
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

// SetAppToDenomsMapping set denoms for appId in Collector LookupTable.
func (k Keeper) SetAppToDenomsMapping(ctx sdk.Context, appID uint64, appToDenom types.AppToDenomsMapping) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.CollectorForDenomKey(appID)
		value = k.cdc.MustMarshal(&appToDenom)
	)

	store.Set(key, value)
}

// GetAppToDenomsMapping get denoms for appId in Collector LookupTable.
func (k Keeper) GetAppToDenomsMapping(ctx sdk.Context, appID uint64) (appToDenom types.AppToDenomsMapping, found bool) {
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

func (k Keeper) GetAllAppToDenomsMapping(ctx sdk.Context) (appToDenomsMapping []types.AppToDenomsMapping) {
	var (
		store = ctx.KVStore(k.storeKey)
		iter  = sdk.KVStorePrefixIterator(store, types.CollectorForDenomKeyPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var fee types.AppToDenomsMapping
		k.cdc.MustUnmarshal(iter.Value(), &fee)
		appToDenomsMapping = append(appToDenomsMapping, fee)
	}
	return appToDenomsMapping
}

// SetAuctionMappingForApp sets auction map data for app/product.
func (k Keeper) SetAuctionMappingForApp(ctx sdk.Context, records ...types.CollectorAuctionLookupTable) error {

	for _, msg := range records {
		_, found := k.GetApp(ctx, msg.AppId)
		if !found {
			return types.ErrorAppDoesNotExist
		}
		_, found1 := k.GetAuctionParams(ctx, msg.AppId)
		if !found1 {
			return types.ErrorAuctionParamsNotSet
		}

		var collectorAuctionLookup types.CollectorAuctionLookupTable
		collectorAuctionLookup.AppId = msg.AppId
		var assetIDToAuctionLookups []types.AssetIdToAuctionLookupTable
		result1, found := k.GetAuctionMappingForApp(ctx, msg.AppId)

		if found {
			assetIDToAuctionLookups = result1.AssetIdToAuctionLookup
		}

		for _, data := range msg.AssetIdToAuctionLookup {
			_, found := k.GetAsset(ctx, data.AssetId)
			if !found {
				return types.ErrorAssetDoesNotExist
			}
			if data.IsSurplusAuction && data.IsDistributor {
				return types.ErrorSurplusDistributerCantbeTrue
			}
			if data.IsSurplusAuction && data.IsDebtAuction {
				return types.ErrorSurplusDebtrCantbeTrueSameTime
			}
			duplicate, index := k.DuplicateCheck(ctx, msg.AppId, data.AssetId)
			if duplicate {
				assetIDToAuctionLookups = append(assetIDToAuctionLookups[:index], assetIDToAuctionLookups[index+1:]...)
				var assetToAuctionUpdate types.AssetIdToAuctionLookupTable
				assetToAuctionUpdate.AssetId = data.AssetId
				assetToAuctionUpdate.IsSurplusAuction = data.IsSurplusAuction
				assetToAuctionUpdate.IsDebtAuction = data.IsDebtAuction
				assetToAuctionUpdate.IsDistributor = data.IsDistributor
				assetToAuctionUpdate.IsAuctionActive = data.IsAuctionActive
				assetToAuctionUpdate.AssetOutOraclePrice = data.AssetOutOraclePrice
				assetToAuctionUpdate.AssetOutPrice = data.AssetOutPrice
				assetIDToAuctionLookups = append(assetIDToAuctionLookups, assetToAuctionUpdate)
				continue
			}
			assetIDToAuctionLookup := types.AssetIdToAuctionLookupTable{

				AssetId:             data.AssetId,
				IsSurplusAuction:    data.IsSurplusAuction,
				IsDebtAuction:       data.IsDebtAuction,
				IsDistributor:       data.IsDistributor,
				IsAuctionActive:     data.IsAuctionActive,
				AssetOutOraclePrice: data.AssetOutOraclePrice,
				AssetOutPrice:       data.AssetOutPrice,
			}
			assetIDToAuctionLookups = append(assetIDToAuctionLookups, assetIDToAuctionLookup)
		}
		collectorAuctionLookup.AssetIdToAuctionLookup = assetIDToAuctionLookups
		var (
			store = ctx.KVStore(k.storeKey)
			key   = types.AppIDToAuctionMappingKey(msg.AppId)
			value = k.cdc.MustMarshal(&collectorAuctionLookup)
		)

		store.Set(key, value)
	}
	return nil
}

func (k Keeper) DuplicateCheck(ctx sdk.Context, appID, assetID uint64) (found bool, index int) {
	result, found := k.GetAuctionMappingForApp(ctx, appID)
	if !found {
		return false, 0
	}
	for i, data := range result.AssetIdToAuctionLookup {
		if data.AssetId == assetID {
			return true, i
		}
	}

	return false, 0
}

// GetAuctionMappingForApp gets auction map data for app/product.
func (k Keeper) GetAuctionMappingForApp(ctx sdk.Context, appID uint64) (collectorAuctionLookupTable types.CollectorAuctionLookupTable, found bool) {
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

func (k Keeper) GetAllAuctionMappingForApp(ctx sdk.Context) (collectorAuctionLookupTable []types.CollectorAuctionLookupTable, found bool) {
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

func (k Keeper) SetNetFeeCollectedData(ctx sdk.Context, appID, assetID uint64, fee sdk.Int) error {
	if fee.IsNegative() {
		return types.ErrorNetFeesCanNotBeNegative
	}
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

// GetNetFeeCollectedData sets net fees collected.
func (k Keeper) GetNetFeeCollectedData(ctx sdk.Context, appID uint64) (netFeeData types.NetFeeCollectedData, found bool) {
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

func (k Keeper) GetAllNetFeeCollectedData(ctx sdk.Context) (netFeeCollectedData []types.NetFeeCollectedData) {
	var (
		store = ctx.KVStore(k.storeKey)
		iter  = sdk.KVStorePrefixIterator(store, types.NetFeeCollectedDataPrefix)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var fee types.NetFeeCollectedData
		k.cdc.MustUnmarshal(iter.Value(), &fee)
		netFeeCollectedData = append(netFeeCollectedData, fee)
	}
	return netFeeCollectedData
}

func (k Keeper) WasmSetCollectorLookupTable(ctx sdk.Context, collectorBindings *bindings.MsgSetCollectorLookupTable) error {
	if !k.HasAsset(ctx, collectorBindings.CollectorAssetID) {
		return types.ErrorAssetDoesNotExist
	}
	if !k.HasAsset(ctx, collectorBindings.SecondaryAssetID) {
		return types.ErrorAssetDoesNotExist
	}
	if collectorBindings.CollectorAssetID == collectorBindings.SecondaryAssetID {
		return types.ErrorDuplicateAssetDenoms
	}
	appDenom, found := k.GetAppToDenomsMapping(ctx, collectorBindings.AppID)
	if found {
		//check if asset denom already exists
		var check = 0
		for _, data := range appDenom.AssetIds {
			if data == collectorBindings.CollectorAssetID {
				check++
			}
		}
		if check > 0 {
			return types.ErrorDuplicateCollectorDenomForApp
		}
		// if denom is new then append
		appDenom.AssetIds = append(appDenom.AssetIds, collectorBindings.CollectorAssetID)
		k.SetAppToDenomsMapping(ctx, collectorBindings.AppID, appDenom)
	} else {
		//initialize the mapping
		var appDenomNew types.AppToDenomsMapping
		appDenomNew.AppId = collectorBindings.AppID
		appDenomNew.AssetIds = append(appDenomNew.AssetIds, collectorBindings.CollectorAssetID)
		k.SetAppToDenomsMapping(ctx, collectorBindings.AppID, appDenomNew)
	}

	var Collector = types.CollectorLookupTable{
		AppId:            collectorBindings.AppID,
		CollectorAssetId: collectorBindings.CollectorAssetID,
		SecondaryAssetId: collectorBindings.SecondaryAssetID,
		SurplusThreshold: collectorBindings.SurplusThreshold,
		DebtThreshold:    collectorBindings.DebtThreshold,
		LockerSavingRate: collectorBindings.LockerSavingRate,
		LotSize:          collectorBindings.LotSize,
		BidFactor:        collectorBindings.BidFactor,
		DebtLotSize:      collectorBindings.DebtLotSize,
	}
	accmLookup, _ := k.GetCollectorLookupTable(ctx, collectorBindings.AppID)
	accmLookup.AppId = collectorBindings.AppID
	accmLookup.AssetRateInfo = append(accmLookup.AssetRateInfo, Collector)

	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.CollectorLookupTableMappingKey(collectorBindings.AppID)
		value = k.cdc.MustMarshal(&accmLookup)
	)

	store.Set(key, value)
	return nil
}

func (k Keeper) WasmSetCollectorLookupTableQuery(ctx sdk.Context, appID, collectorAssetID, secondaryAssetID uint64) (bool, string) {
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

func (k Keeper) WasmSetAuctionMappingForApp(ctx sdk.Context, auctionMappingBinding *bindings.MsgSetAuctionMappingForApp) error {
	result1, found := k.GetAuctionMappingForApp(ctx, auctionMappingBinding.AppID)
	var collectorAuctionLookup types.CollectorAuctionLookupTable
	var assetIDToAuctionLookups []types.AssetIdToAuctionLookupTable
	collectorAuctionLookup.AppId = auctionMappingBinding.AppID
	if found {
		assetIDToAuctionLookups = result1.AssetIdToAuctionLookup
	}
	for i := range auctionMappingBinding.AssetIDs {
		if auctionMappingBinding.IsSurplusAuctions[i] && auctionMappingBinding.IsDistributor[i] {
			return types.ErrorSurplusDistributerCantbeTrue
		}
		if auctionMappingBinding.IsSurplusAuctions[i] && auctionMappingBinding.IsDebtAuctions[i] {
			return types.ErrorSurplusDebtrCantbeTrueSameTime
		}
		duplicate, index := k.DuplicateCheck(ctx, auctionMappingBinding.AppID, auctionMappingBinding.AssetIDs[i])
		if duplicate {
			assetIDToAuctionLookups = append(assetIDToAuctionLookups[:index], assetIDToAuctionLookups[index+1:]...)
			var assetToAuctionUpdate types.AssetIdToAuctionLookupTable
			assetToAuctionUpdate.AssetId = auctionMappingBinding.AssetIDs[i]
			assetToAuctionUpdate.IsSurplusAuction = auctionMappingBinding.IsSurplusAuctions[i]
			assetToAuctionUpdate.IsDebtAuction = auctionMappingBinding.IsDebtAuctions[i]
			assetToAuctionUpdate.IsDistributor = auctionMappingBinding.IsDistributor[i]
			assetToAuctionUpdate.IsAuctionActive = false
			assetToAuctionUpdate.AssetOutOraclePrice = auctionMappingBinding.AssetOutOraclePrices[i]
			assetToAuctionUpdate.AssetOutPrice = auctionMappingBinding.AssetOutPrices[i]
			assetIDToAuctionLookups = append(assetIDToAuctionLookups, assetToAuctionUpdate)
			continue
		}
		assetIDToAuctionLookup := types.AssetIdToAuctionLookupTable{
			AssetId:             auctionMappingBinding.AssetIDs[i],
			IsSurplusAuction:    auctionMappingBinding.IsSurplusAuctions[i],
			IsDebtAuction:       auctionMappingBinding.IsDebtAuctions[i],
			IsDistributor:       auctionMappingBinding.IsDistributor[i],
			IsAuctionActive:     false,
			AssetOutOraclePrice: auctionMappingBinding.AssetOutOraclePrices[i],
			AssetOutPrice:       auctionMappingBinding.AssetOutPrices[i],
		}
		assetIDToAuctionLookups = append(assetIDToAuctionLookups, assetIDToAuctionLookup)
	}
	collectorAuctionLookup.AssetIdToAuctionLookup = assetIDToAuctionLookups
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.AppIDToAuctionMappingKey(auctionMappingBinding.AppID)
		value = k.cdc.MustMarshal(&collectorAuctionLookup)
	)

	store.Set(key, value)

	return nil
}

func (k Keeper) WasmSetAuctionMappingForAppQuery(ctx sdk.Context, appID uint64) (bool, string) {
	_, _ = k.GetAppidToAssetCollectorMapping(ctx, appID)

	return true, ""
}

func (k Keeper) WasmUpdateCollectorLookupTable(ctx sdk.Context, updateColBinding *bindings.MsgUpdateCollectorLookupTable) error {
	var Collector types.CollectorLookupTable
	accmLookup, _ := k.GetCollectorLookupTable(ctx, updateColBinding.AppID)

	for _, data := range accmLookup.AssetRateInfo {
		if data.CollectorAssetId == updateColBinding.AssetID {
			Collector.CollectorAssetId = updateColBinding.AssetID
			Collector.AppId = data.AppId
			Collector.BidFactor = updateColBinding.BidFactor
			Collector.DebtThreshold = updateColBinding.DebtThreshold
			Collector.SurplusThreshold = updateColBinding.SurplusThreshold
			Collector.LockerSavingRate = updateColBinding.LSR
			Collector.LotSize = updateColBinding.LotSize
			Collector.SecondaryAssetId = data.SecondaryAssetId
			Collector.DebtLotSize = updateColBinding.DebtLotSize
		}
	}
	err := k.SetCollectorLookupTableForWasm(ctx, Collector)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) WasmUpdateCollectorLookupTableQuery(ctx sdk.Context, appID, assetID uint64) (bool, string) {
	_, found := k.GetCollectorLookupByAsset(ctx, appID, assetID)
	if !found {
		return false, types.ErrorDataDoesNotExists.Error()
	}
	return true, ""
}
