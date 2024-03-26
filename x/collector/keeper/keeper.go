package keeper

import (
	sdkmath "cosmossdk.io/math"
	"fmt"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	protobuftypes "github.com/cosmos/gogoproto/types"

	"github.com/comdex-official/comdex/x/collector/expected"
	"github.com/comdex-official/comdex/x/collector/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		asset      expected.AssetKeeper
		auction    expected.AuctionKeeper
		locker     expected.LockerKeeper
		rewards    expected.RewardsKeeper
		esm        expected.EsmKeeper
		paramStore paramtypes.Subspace
		bank       expected.BankKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	asset expected.AssetKeeper,
	auction expected.AuctionKeeper,
	locker expected.LockerKeeper,
	rewards expected.RewardsKeeper,
	esm expected.EsmKeeper,
	ps paramtypes.Subspace,
	bank expected.BankKeeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		asset:      asset,
		auction:    auction,
		locker:     locker,
		rewards:    rewards,
		esm:        esm,
		paramStore: ps,
		bank:       bank,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) ModuleBalance(ctx sdk.Context, moduleName string, denom string) sdk.Int {
	return k.bank.GetBalance(ctx, authtypes.NewModuleAddress(moduleName), denom).Amount
}

func (k Keeper) Deposit(ctx sdk.Context, Amount sdk.Coin, AppID uint64, addr string) error {

	asset, found := k.asset.GetAssetForDenom(ctx, Amount.Denom)
	if !found {
		return types.ErrorAssetDoesNotExist
	}
	// check if denom is cmst
	if asset.Denom != "ucmst" {
		return types.ErrorAssetDoesNotExist
	}
	// check if app id exists and app name is harbor
	app, found := k.asset.GetApp(ctx, AppID)
	if !found {
		return types.ErrorAppDoesNotExist
	}
	if app.Name != "harbor" {
		return types.ErrorAppDoesNotExist
	}

	address, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return err
	}

	err = k.bank.SendCoinsFromAccountToModule(ctx, address, types.ModuleName, sdk.NewCoins(Amount))
	if err != nil {
		return err
	}

	err = k.SetNetFeeCollectedData(ctx, AppID, asset.Id, Amount.Amount)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) UpdateDebtParams(ctx sdk.Context, appId, assetID, slots uint64, debtThreshold, lotSize, debtLotSize sdkmath.Int, isDebtAuction bool, addr string) error {

	// check if address is admin
	getAdmin := k.esm.GetParams(ctx).Admin

	// check if address is admin in getAdmin array
	if getAdmin[0] != addr {
		return esmtypes.ErrorUnauthorized
	}

	// check if app id exists and app name is harbor
	app, found := k.asset.GetApp(ctx, appId)
	if !found {
		return types.ErrorAppDoesNotExist
	}
	if app.Name != "harbor" {
		return types.ErrorAppDoesNotExist
	}

	asset, found := k.asset.GetAsset(ctx, assetID)
	if !found {
		return types.ErrorAssetDoesNotExist
	}

	if asset.Denom != "ucmst" {
		return types.ErrorAssetDoesNotExist
	}

	if slots == 0 {
		return types.ErrorAmountCanNotBeZero
	}

	if lotSize.IsZero() {
		return types.ErrorAmountCanNotBeZero
	}

	if debtLotSize.IsZero() {
		return types.ErrorAmountCanNotBeZero
	}

	// get CollectorLookupTableData
	collectorLookupTable, found := k.GetCollectorLookupTable(ctx, appId, assetID)
	if !found {
		return types.ErrorDataDoesNotExists
	}

	collectorLookupTable.DebtLotSize = debtLotSize
	collectorLookupTable.DebtThreshold = debtThreshold
	collectorLookupTable.LotSize = lotSize

	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.CollectorLookupTableMappingKey(appId, assetID)
		value = k.cdc.MustMarshal(&collectorLookupTable)
	)

	store.Set(key, value)

	// get AppAssetIdToAuctionLookupTable
	appAssetIdToAuctionLookupTable, found := k.GetAuctionMappingForApp(ctx, appId, assetID)
	if !found {
		return types.ErrorAuctionParamsNotSet
	}

	appAssetIdToAuctionLookupTable.IsDebtAuction = isDebtAuction

	// set SetAuctionMappingForApp
	err := k.SetAuctionMappingForApp(ctx, appAssetIdToAuctionLookupTable)
	if err != nil {
		return err
	}

	// set slots
	k.SetSlots(ctx, slots)

	return nil
}

func (k Keeper) SetSlots(ctx sdk.Context, slots uint64) {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.SlotsKeyPrefix
		value = k.cdc.MustMarshal(
			&protobuftypes.UInt64Value{
				Value: slots,
			},
		)
	)

	store.Set(key, value)
}

func (k Keeper) GetSlots(ctx sdk.Context) uint64 {
	var (
		store = ctx.KVStore(k.storeKey)
		key   = types.SlotsKeyPrefix
		value = store.Get(key)
	)

	if value == nil {
		return 1
	}

	var id protobuftypes.UInt64Value
	k.cdc.MustUnmarshal(value, &id)

	return id.GetValue()
}
