package keeper

import (
	"github.com/comdex-official/comdex/app/wasm/bindings"
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	"github.com/comdex-official/comdex/x/collector/types"
	lockertypes "github.com/comdex-official/comdex/x/locker/types"
	rewardstypes "github.com/comdex-official/comdex/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetAmountFromCollector returns amount from the collector.
func (k Keeper) GetAmountFromCollector(ctx sdk.Context, appID, assetID uint64, amount sdk.Int) (sdk.Int, error) {
	netFeeData, found := k.GetNetFeeCollectedData(ctx, appID, assetID)
	var returnedFee sdk.Int
	if !found {
		return returnedFee, types.ErrorDataDoesNotExists
	}
	if amount.IsNegative() {
		return returnedFee, types.ErrorAmountCanNotBeNegative
	}

	if !(netFeeData.NetFeesCollected.Sub(amount).GT(sdk.ZeroInt())) {
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

	returnedFee = amount
	return returnedFee, nil
}

func (k Keeper) DecreaseNetFeeCollectedData(ctx sdk.Context, appID, assetID uint64, amount sdk.Int, netCollected types.AppAssetIdToFeeCollectedData) error {
	netCollected.AppId = appID

	var netCollectedFee sdk.Int

	netCollected.AssetId = assetID
	netCollectedFee = netCollected.NetFeesCollected.Sub(amount)
	if netCollectedFee.IsNegative() {
		return types.ErrorNetFeesCanNotBeNegative
	}
	netCollected.NetFeesCollected = netCollectedFee

	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.NetFeeCollectedDataKey(appID, assetID)
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

	collectorData, found := k.GetAppidToAssetCollectorMapping(ctx, appID, assetID)
	if !found {
		//create a new instance of appID To AssetCollectorMapping
		var collectorNewData types.AppToAssetIdCollectorMapping
		collectorNewData.AppId = appID
		collectorNewData.AssetId = assetID

		var newCollector types.CollectorData
		newCollector.CollectedClosingFee = collectedClosingFee
		newCollector.CollectedOpeningFee = collectedOpeningFee
		newCollector.CollectedStabilityFee = collectedStabilityFee
		newCollector.LiquidationRewardsCollected = liquidationRewardsCollected

		collectorNewData.Collector = newCollector

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

		var collectorNewData types.AppToAssetIdCollectorMapping
		collectorNewData.AppId = appID
		collectorNewData.AssetId = assetID

		var newCollector types.CollectorData
		newCollector.CollectedClosingFee = collectorData.Collector.CollectedClosingFee.Add(collectedClosingFee)
		newCollector.CollectedOpeningFee = collectorData.Collector.CollectedOpeningFee.Add(collectedOpeningFee)
		newCollector.CollectedStabilityFee = collectorData.Collector.CollectedStabilityFee.Add(collectedStabilityFee)
		newCollector.LiquidationRewardsCollected = collectorData.Collector.LiquidationRewardsCollected.Add(liquidationRewardsCollected)
		collectorNewData.Collector = newCollector

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
	return nil
}

// SetAppidToAssetCollectorMapping update collector with app_id and asset.
func (k Keeper) SetAppidToAssetCollectorMapping(ctx sdk.Context, appAssetCollectorData types.AppToAssetIdCollectorMapping) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.AppidToAssetCollectorMappingKey(appAssetCollectorData.AppId, appAssetCollectorData.AssetId)
		value = k.cdc.MustMarshal(&appAssetCollectorData)
	)
	store.Set(key, value)
}

// GetAppidToAssetCollectorMapping returns app_id to asset mapping for collector.
func (k Keeper) GetAppidToAssetCollectorMapping(ctx sdk.Context, appID, assetID uint64) (appAssetCollectorData types.AppToAssetIdCollectorMapping, found bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.AppidToAssetCollectorMappingKey(appID, assetID)
		value = store.Get(key)
	)

	if value == nil {
		return appAssetCollectorData, false
	}

	k.cdc.MustUnmarshal(value, &appAssetCollectorData)
	return appAssetCollectorData, true
}

func (k Keeper) GetAllAppidToAssetCollectorMapping(ctx sdk.Context) (appIDToAssetCollectorMapping []types.AppToAssetIdCollectorMapping) {
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
		var fee types.AppToAssetIdCollectorMapping
		k.cdc.MustUnmarshal(iter.Value(), &fee)
		appIDToAssetCollectorMapping = append(appIDToAssetCollectorMapping, fee)
	}
	return appIDToAssetCollectorMapping
}

// GetCollectorDataForAppIDAssetID returns app_id to asset mapping for collector.
func (k Keeper) GetCollectorDataForAppIDAssetID(ctx sdk.Context, appID uint64, assetID uint64) (collectorData types.CollectorData, found bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.AppidToAssetCollectorMappingKey(appID, assetID)
		value = store.Get(key)
	)
	var appAssetCollectorData types.AppToAssetIdCollectorMapping
	if value == nil {
		return collectorData, false
	}
	k.cdc.MustUnmarshal(value, &appAssetCollectorData)
	collectorData = appAssetCollectorData.Collector

	return collectorData, false
}

// SetCollectorLookupTable updates the collector lookup store.
func (k Keeper) SetCollectorLookupTable(ctx sdk.Context, records types.CollectorLookupTableData) error {

	if !k.HasAsset(ctx, records.CollectorAssetId) {
		return types.ErrorAssetDoesNotExist
	}
	if !k.HasAsset(ctx, records.SecondaryAssetId) {
		return types.ErrorAssetDoesNotExist
	}
	if records.CollectorAssetId == records.SecondaryAssetId {
		return types.ErrorDuplicateAssetDenoms
	}
	_, found := k.GetMintGenesisTokenData(ctx, records.AppId, records.SecondaryAssetId)
	if !found {
		return types.ErrorAssetNotAddedForGenesisMinting
	}
	appDenom, found := k.GetAppToDenomsMapping(ctx, records.AppId)
	if found {
		//check if assetdenom already exists
		var check = 0
		for _, data := range appDenom.AssetIds {
			if data == records.CollectorAssetId {
				check++
			}
		}
		if check > 0 {
			return types.ErrorDuplicateCollectorDenomForApp
		}
		// if denom is new then append
		appDenom.AssetIds = append(appDenom.AssetIds, records.CollectorAssetId)
		k.SetAppToDenomsMapping(ctx, records.AppId, appDenom)
	} else {
		//initialize the mapping
		var appDenomNew types.AppToDenomsMapping
		appDenomNew.AppId = records.AppId
		appDenomNew.AssetIds = append(appDenomNew.AssetIds, records.CollectorAssetId)
		k.SetAppToDenomsMapping(ctx, records.AppId, appDenomNew)
	}

	var Collector = types.CollectorLookupTableData{
		AppId:            records.AppId,
		CollectorAssetId: records.CollectorAssetId,
		SecondaryAssetId: records.SecondaryAssetId,
		SurplusThreshold: records.SurplusThreshold,
		DebtThreshold:    records.DebtThreshold,
		LockerSavingRate: records.LockerSavingRate,
		LotSize:          records.LotSize,
		BidFactor:        records.BidFactor,
		DebtLotSize:      records.DebtLotSize,
	}

	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.CollectorLookupTableMappingKey(records.AppId, records.CollectorAssetId)
		value = k.cdc.MustMarshal(&Collector)
	)

	store.Set(key, value)
	return nil
}

// GetCollectorLookupTable returns collector lookup table.
func (k Keeper) GetCollectorLookupTable(ctx sdk.Context, appID, assetID uint64) (collectorLookup types.CollectorLookupTableData, found bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.CollectorLookupTableMappingKey(appID, assetID)
		value = store.Get(key)
	)

	if value == nil {
		return collectorLookup, false
	}

	k.cdc.MustUnmarshal(value, &collectorLookup)
	return collectorLookup, true
}

func (k Keeper) GetCollectorLookupTableByApp(ctx sdk.Context, appID uint64) (collectorLookup []types.CollectorLookupTableData, found bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.CollectorLookupTableMappingByAppKey(appID)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var table types.CollectorLookupTableData
		k.cdc.MustUnmarshal(iter.Value(), &table)
		collectorLookup = append(collectorLookup, table)
	}
	if collectorLookup == nil {
		return nil, false
	}
	return collectorLookup, true
}

func (k Keeper) GetAllCollectorLookupTable(ctx sdk.Context) (collectorLookup []types.CollectorLookupTableData) {
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
		var fee types.CollectorLookupTableData
		k.cdc.MustUnmarshal(iter.Value(), &fee)
		collectorLookup = append(collectorLookup, fee)
	}
	return collectorLookup
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
func (k Keeper) SetAuctionMappingForApp(ctx sdk.Context, record types.AppAssetIdToAuctionLookupTable) error {

	_, found := k.GetApp(ctx, record.AppId)
	if !found {
		return types.ErrorAppDoesNotExist
	}
	_, found1 := k.GetAuctionParams(ctx, record.AppId)
	if !found1 {
		return types.ErrorAuctionParamsNotSet
	}
	_, found2 := k.GetAsset(ctx, record.AssetId)
	if !found2 {
		return types.ErrorAssetDoesNotExist
	}
	if record.IsSurplusAuction && record.IsDistributor {
		return types.ErrorSurplusDistributerCantbeTrue
	}
	if record.IsSurplusAuction && record.IsDebtAuction {
		return types.ErrorSurplusDebtrCantbeTrueSameTime
	}

	var appAssetToAuction types.AppAssetIdToAuctionLookupTable
	result, found3 := k.GetAuctionMappingForApp(ctx, record.AppId, record.AssetId)
	if !found3 {
		appAssetToAuction.AppId = record.AppId
		appAssetToAuction.AssetId = record.AssetId
		appAssetToAuction.AssetOutOraclePrice = record.AssetOutOraclePrice
		appAssetToAuction.AssetOutPrice = record.AssetOutPrice
		appAssetToAuction.IsDebtAuction = record.IsDebtAuction
		appAssetToAuction.IsDistributor = record.IsDistributor
		appAssetToAuction.IsSurplusAuction = record.IsSurplusAuction
		appAssetToAuction.IsAuctionActive = record.IsAuctionActive
		result = appAssetToAuction
	} else {
		result.AssetOutOraclePrice = record.AssetOutOraclePrice
		result.AssetOutPrice = record.AssetOutPrice
		result.IsDebtAuction = record.IsDebtAuction
		result.IsDistributor = record.IsDistributor
		result.IsSurplusAuction = record.IsSurplusAuction
		result.IsAuctionActive = record.IsAuctionActive
	}
	if result.IsSurplusAuction && result.IsDistributor {
		return types.ErrorSurplusDistributerCantbeTrue
	}
	if result.IsSurplusAuction && result.IsDebtAuction {
		return types.ErrorSurplusDebtrCantbeTrueSameTime
	}

	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.AppIDToAuctionMappingKey(record.AppId, record.AssetId)
		value = k.cdc.MustMarshal(&result)
	)

	store.Set(key, value)
	return nil
}

// GetAuctionMappingForApp gets auction map data for app/product.
func (k Keeper) GetAuctionMappingForApp(ctx sdk.Context, appID, assetID uint64) (collectorAuctionLookupTable types.AppAssetIdToAuctionLookupTable, found bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.AppIDToAuctionMappingKey(appID, assetID)
		value = store.Get(key)
	)

	if value == nil {
		return collectorAuctionLookupTable, false
	}

	k.cdc.MustUnmarshal(value, &collectorAuctionLookupTable)
	return collectorAuctionLookupTable, true
}

func (k Keeper) GetAllAuctionMappingForApp(ctx sdk.Context) (collectorAuctionLookupTable []types.AppAssetIdToAuctionLookupTable, found bool) {
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
		var table types.AppAssetIdToAuctionLookupTable
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
	collectorData, found := k.GetNetFeeCollectedData(ctx, appID, assetID)
	if !found {
		var netCollected types.AppAssetIdToFeeCollectedData
		netCollected.AppId = appID
		netCollected.AssetId = assetID
		netCollected.NetFeesCollected = fee
		var (
			store = ctx.KVStore(k.storeKey)
			key   = types.NetFeeCollectedDataKey(appID, assetID)
			value = k.cdc.MustMarshal(&netCollected)
		)

		store.Set(key, value)
	} else {
		var netCollected types.AppAssetIdToFeeCollectedData
		netCollected.AppId = appID

		var netCollectedFee sdk.Int

		netCollected.AssetId = assetID
		netCollectedFee = collectorData.NetFeesCollected.Add(fee)

		netCollected.NetFeesCollected = netCollectedFee
		var (
			store = ctx.KVStore(k.storeKey)
			key   = types.NetFeeCollectedDataKey(appID, assetID)
			value = k.cdc.MustMarshal(&netCollected)
		)

		store.Set(key, value)
	}

	return nil
}

// GetNetFeeCollectedData sets net fees collected.
func (k Keeper) GetNetFeeCollectedData(ctx sdk.Context, appID, assetID uint64) (netFeeData types.AppAssetIdToFeeCollectedData, found bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.NetFeeCollectedDataKey(appID, assetID)
		value = store.Get(key)
	)

	if value == nil {
		return netFeeData, false
	}

	k.cdc.MustUnmarshal(value, &netFeeData)
	return netFeeData, true
}

func (k Keeper) GetAppNetFeeCollectedData(ctx sdk.Context, appID uint64) (netFeeData []types.AppAssetIdToFeeCollectedData, found bool) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.AppNetFeeCollectedDataKey(appID)
		iter  = sdk.KVStorePrefixIterator(store, key)
	)

	defer func(iter sdk.Iterator) {
		err := iter.Close()
		if err != nil {
			return
		}
	}(iter)

	for ; iter.Valid(); iter.Next() {
		var mapData types.AppAssetIdToFeeCollectedData
		k.cdc.MustUnmarshal(iter.Value(), &mapData)
		netFeeData = append(netFeeData, mapData)
	}
	if netFeeData == nil {
		return nil, false
	}
	return netFeeData, true
}

func (k Keeper) GetAllNetFeeCollectedData(ctx sdk.Context) (netFeeCollectedData []types.AppAssetIdToFeeCollectedData) {
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
		var fee types.AppAssetIdToFeeCollectedData
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
	blockHeight := ctx.BlockHeight()

	if collectorBindings.LockerSavingRate.IsZero() {
		blockHeight = 0
	}

	var Collector = types.CollectorLookupTableData{
		AppId:            collectorBindings.AppID,
		CollectorAssetId: collectorBindings.CollectorAssetID,
		SecondaryAssetId: collectorBindings.SecondaryAssetID,
		SurplusThreshold: collectorBindings.SurplusThreshold,
		DebtThreshold:    collectorBindings.DebtThreshold,
		LockerSavingRate: collectorBindings.LockerSavingRate,
		LotSize:          collectorBindings.LotSize,
		BidFactor:        collectorBindings.BidFactor,
		DebtLotSize:      collectorBindings.DebtLotSize,
		BlockHeight:      blockHeight,
		BlockTime:        ctx.BlockTime(),
	}

	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.CollectorLookupTableMappingKey(collectorBindings.AppID, collectorBindings.CollectorAssetID)
		value = k.cdc.MustMarshal(&Collector)
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
	result, found := k.GetAuctionMappingForApp(ctx, auctionMappingBinding.AppID, auctionMappingBinding.AssetIDs)
	var assetToAuctionUpdate types.AppAssetIdToAuctionLookupTable
	if auctionMappingBinding.IsSurplusAuctions && auctionMappingBinding.IsDistributor {
		return types.ErrorSurplusDistributerCantbeTrue
	}
	if auctionMappingBinding.IsSurplusAuctions && auctionMappingBinding.IsDebtAuctions {
		return types.ErrorSurplusDebtrCantbeTrueSameTime
	}
	if !found {
		assetToAuctionUpdate.AppId = auctionMappingBinding.AppID
		assetToAuctionUpdate.AssetId = auctionMappingBinding.AssetIDs
		assetToAuctionUpdate.IsSurplusAuction = auctionMappingBinding.IsSurplusAuctions
		assetToAuctionUpdate.IsDebtAuction = auctionMappingBinding.IsDebtAuctions
		assetToAuctionUpdate.IsDistributor = auctionMappingBinding.IsDistributor
		assetToAuctionUpdate.IsAuctionActive = false
		assetToAuctionUpdate.AssetOutOraclePrice = auctionMappingBinding.AssetOutOraclePrices
		assetToAuctionUpdate.AssetOutPrice = auctionMappingBinding.AssetOutPrices
		result = assetToAuctionUpdate
	} else {
		result.AssetOutOraclePrice = auctionMappingBinding.AssetOutOraclePrices
		result.AssetOutPrice = auctionMappingBinding.AssetOutPrices
		result.IsDebtAuction = auctionMappingBinding.IsDebtAuctions
		result.IsDistributor = auctionMappingBinding.IsDistributor
		result.IsSurplusAuction = auctionMappingBinding.IsSurplusAuctions
	}

	if result.IsSurplusAuction && result.IsDistributor {
		return types.ErrorSurplusDistributerCantbeTrue
	}
	if result.IsSurplusAuction && result.IsDebtAuction {
		return types.ErrorSurplusDebtrCantbeTrueSameTime
	}
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.AppIDToAuctionMappingKey(auctionMappingBinding.AppID, auctionMappingBinding.AssetIDs)
		value = k.cdc.MustMarshal(&result)
	)

	store.Set(key, value)

	return nil
}

func (k Keeper) WasmSetAuctionMappingForAppQuery(ctx sdk.Context, appID uint64) (bool, string) {
	_, found := k.GetApp(ctx, appID)
	if found {
		return true, ""
	}
	return false, types.ErrorAppDoesNotExist.Error()
}

func (k Keeper) WasmUpdateCollectorLookupTable(ctx sdk.Context, updateColBinding *bindings.MsgUpdateCollectorLookupTable) error {
	Collector, _ := k.GetCollectorLookupTable(ctx, updateColBinding.AppID, updateColBinding.AssetID)
	_, found := k.GetReward(ctx, Collector.AppId, Collector.CollectorAssetId)
	if found {
		if Collector.LockerSavingRate != updateColBinding.LSR {
			if updateColBinding.LSR.IsZero() {

				// run script to distrubyte reward
				k.LockerIterateRewards(ctx, Collector.LockerSavingRate, Collector.BlockHeight, Collector.BlockTime.Unix(), updateColBinding.AppID, updateColBinding.AssetID, false)
				Collector.BlockTime = ctx.BlockTime()
				Collector.BlockHeight = 0

			} else if Collector.LockerSavingRate.IsZero() {
				// do nothing
				Collector.BlockHeight = ctx.BlockHeight()
				Collector.BlockTime = ctx.BlockTime()
			} else if Collector.LockerSavingRate.GT(sdk.ZeroDec()) && updateColBinding.LSR.GT(sdk.ZeroDec()) {
				// run script to distribute
				k.LockerIterateRewards(ctx, Collector.LockerSavingRate, Collector.BlockHeight, Collector.BlockTime.Unix(), updateColBinding.AppID, updateColBinding.AssetID, true)
				Collector.BlockHeight = ctx.BlockHeight()
				Collector.BlockTime = ctx.BlockTime()

			}
		}
	}

	Collector.BidFactor = updateColBinding.BidFactor
	Collector.DebtThreshold = updateColBinding.DebtThreshold
	Collector.SurplusThreshold = updateColBinding.SurplusThreshold
	Collector.LockerSavingRate = updateColBinding.LSR
	Collector.LotSize = updateColBinding.LotSize
	Collector.DebtLotSize = updateColBinding.DebtLotSize

	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.CollectorLookupTableMappingKey(updateColBinding.AppID, updateColBinding.AssetID)
		value = k.cdc.MustMarshal(&Collector)
	)
	store.Set(key, value)
	return nil
}

func (k Keeper) LockerIterateRewards(ctx sdk.Context, collectorLsr sdk.Dec, collectorBh, collectorBt int64, appID, assetID uint64, changeTypes bool) {
	lockers, _ := k.GetLockerLookupTable(ctx, appID, assetID)
	for _, lockID := range lockers.LockerIds {
		lockerData, _ := k.GetLocker(ctx, lockID)
		rewards := sdk.ZeroDec()
		var err error
		if lockerData.BlockHeight == 0 {
			rewards, err = k.CalculationOfRewards(ctx, lockerData.NetBalance, collectorLsr, collectorBt)
			if err != nil {
				return
			}
		} else {
			rewards, err = k.CalculationOfRewards(ctx, lockerData.NetBalance, collectorLsr, lockerData.BlockTime.Unix())
			if err != nil {
				return
			}
		}

		lockerRewardsTracker, found := k.GetLockerRewardTracker(ctx, lockerData.LockerId, appID)
		if !found {
			lockerRewardsTracker = rewardstypes.LockerRewardsTracker{
				LockerId:           lockerData.LockerId,
				AppMappingId:       appID,
				RewardsAccumulated: sdk.ZeroDec(),
			}
		}

		lockerRewardsTracker.RewardsAccumulated = lockerRewardsTracker.RewardsAccumulated.Add(rewards)
		newReward := sdk.ZeroInt()
		if lockerRewardsTracker.RewardsAccumulated.GTE(sdk.OneDec()) {
			newReward = lockerRewardsTracker.RewardsAccumulated.TruncateInt()
			newRewardDec := sdk.NewDec(newReward.Int64())
			lockerRewardsTracker.RewardsAccumulated = lockerRewardsTracker.RewardsAccumulated.Sub(newRewardDec)
		}
		k.SetLockerRewardTracker(ctx, lockerRewardsTracker)

		netFeeCollectedData, found := k.GetNetFeeCollectedData(ctx, appID, lockerData.AssetDepositId)
		if !found {
			continue
		}
		err = k.DecreaseNetFeeCollectedData(ctx, appID, lockerData.AssetDepositId, newReward, netFeeCollectedData)
		if err != nil {
			continue
		}

		assetData, _ := k.GetAsset(ctx, assetID)
		newrewards := rewards.TruncateInt()
		if newrewards.GT(sdk.ZeroInt()) {
			err = k.SendCoinFromModuleToModule(ctx, types.ModuleName, lockertypes.ModuleName, sdk.NewCoins(sdk.NewCoin(assetData.Denom, newrewards)))
			if err != nil {
				continue
			}
		}

		lockerRewardsMapping, found := k.GetLockerTotalRewardsByAssetAppWise(ctx, appID, lockerData.AssetDepositId)
		if !found {
			var lockerReward lockertypes.LockerTotalRewardsByAssetAppWise
			lockerReward.AppId = appID
			lockerReward.AssetId = lockerData.AssetDepositId
			lockerReward.TotalRewards = sdk.ZeroInt().Add(newReward)
			err = k.SetLockerTotalRewardsByAssetAppWise(ctx, lockerReward)
			if err != nil {
				continue
			}
		} else {
			lockerRewardsMapping.TotalRewards = lockerRewardsMapping.TotalRewards.Add(newReward)

			err = k.SetLockerTotalRewardsByAssetAppWise(ctx, lockerRewardsMapping)
			if err != nil {
				continue
			}
		}

		// updating user rewards data
		lockerData.BlockTime = ctx.BlockTime()
		if changeTypes {
			lockerData.BlockHeight = ctx.BlockHeight()
		} else {
			lockerData.BlockHeight = 0
		}

		lockerData.NetBalance = lockerData.NetBalance.Add(newrewards)
		lockerData.ReturnsAccumulated = lockerData.ReturnsAccumulated.Add(newrewards)
		k.SetLocker(ctx, lockerData)
		lockers.DepositedAmount = lockers.DepositedAmount.Add(newrewards)
		k.SetLockerLookupTable(ctx, lockers)

	}

}

func (k Keeper) WasmUpdateCollectorLookupTableQuery(ctx sdk.Context, appID, assetID uint64) (bool, string) {
	_, found := k.GetCollectorLookupTable(ctx, appID, assetID)
	if !found {
		return false, types.ErrorDataDoesNotExists.Error()
	}
	return true, ""
}

func (k Keeper) WasmCheckSurplusRewardQuery(ctx sdk.Context, appID, assetID uint64) sdk.Coin {
	asset, _ := k.GetAsset(ctx, assetID)
	netFeeCollectedData, _ := k.GetNetFeeCollectedData(ctx, appID, assetID)
	auctionMapping, _ := k.GetAuctionMappingForApp(ctx, appID, assetID)
	collectorLookup, _ := k.GetCollectorLookupTable(ctx, appID, assetID)
	netAmount := collectorLookup.SurplusThreshold + collectorLookup.LotSize
	if auctionMapping.IsDistributor && netFeeCollectedData.NetFeesCollected.GT(sdk.NewInt(int64(netAmount))) {
		finalAmount := netFeeCollectedData.NetFeesCollected.Sub(sdk.NewInt(int64(collectorLookup.SurplusThreshold)))
		return sdk.NewCoin(asset.Denom, finalAmount)
	}
	return sdk.NewCoin(asset.Denom, sdk.NewInt(0))
}

func (k Keeper) WasmMsgGetSurplusFund(ctx sdk.Context, appID, assetID uint64, addr sdk.AccAddress, amount sdk.Coin) error {
	err := k.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.NewCoins(amount))
	if err != nil {
		return err
	}
	err = k.DecreaseNetFeeCollectedData(ctx, appID, assetID, amount.Amount, types.AppAssetIdToFeeCollectedData{})
	if err != nil {
		return err
	}
	return nil
}
