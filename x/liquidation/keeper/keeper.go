package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/comdex-official/comdex/x/liquidation/expected"
	"github.com/comdex-official/comdex/x/liquidation/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   sdk.StoreKey
	memKey     sdk.StoreKey
	paramstore paramtypes.Subspace
	account    expected.AccountKeeper
	bank       expected.BankKeeper
	vault      expected.VaultKeeper
	asset      expected.AssetKeeper
	market     expected.MarketKeeper
	auction    expected.AuctionKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,
	account expected.AccountKeeper,
	bank expected.BankKeeper,
	asset expected.AssetKeeper,
	vault expected.VaultKeeper,
	market expected.MarketKeeper,
	auction expected.AuctionKeeper,

) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{

		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
		account:    account,
		bank:       bank,
		asset:      asset,
		vault:      vault,
		market:     market,
		auction:    auction,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}

func uint64InSlice(a uint64, list []uint64) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (k Keeper) WhitelistAppId(ctx sdk.Context, appMappingId uint64) error {
	found := uint64InSlice(appMappingId, k.GetAppIds(ctx).WhitelistedAppMappingIds)
	if found {
		return types.ErrAppIdExists
	}
	WhitelistedAppIds := append(k.GetAppIds(ctx).WhitelistedAppMappingIds, appMappingId)
	UpdatedWhitelistedAppIds := types.WhitelistedAppIds{
		WhitelistedAppMappingIds: WhitelistedAppIds,
	}
	k.SetAppId(ctx, UpdatedWhitelistedAppIds)
	return nil
}

func (k Keeper) RemoveWhitelistAsset(ctx sdk.Context, appMappingId uint64) error {
	WhitelistedAppIds := k.GetAppIds(ctx).WhitelistedAppMappingIds
	found := uint64InSlice(appMappingId, k.GetAppIds(ctx).WhitelistedAppMappingIds)
	if !found {
		return types.ErrAppIdDoesNotExists
	}
	var newAppIds []uint64
	for i := range WhitelistedAppIds {
		if appMappingId != WhitelistedAppIds[i] {
			newAppId := WhitelistedAppIds[i]
			newAppIds = append(newAppIds, newAppId)
		}
	}
	UpdatedWhitelistedAppIds := types.WhitelistedAppIds{
		WhitelistedAppMappingIds: newAppIds,
	}

	k.SetAppId(ctx, UpdatedWhitelistedAppIds)
	return nil
}

//Wasm tx and query binding functions

func (k Keeper) WasmWhitelistAppIdLiquidation(ctx sdk.Context, appMappingId uint64) error {
	WhitelistedAppIds := append(k.GetAppIds(ctx).WhitelistedAppMappingIds, appMappingId)
	UpdatedWhitelistedAppIds := types.WhitelistedAppIds{
		WhitelistedAppMappingIds: WhitelistedAppIds,
	}
	k.SetAppId(ctx, UpdatedWhitelistedAppIds)
	return nil
}

func (k Keeper) WasmWhitelistAppIdLiquidationQuery(ctx sdk.Context, appMappingId uint64) (bool, string) {
	found := uint64InSlice(appMappingId, k.GetAppIds(ctx).WhitelistedAppMappingIds)
	if found {
		return false, types.ErrAppIdExists.Error()
	}
	return true, ""
}

func (k Keeper) WasmRemoveWhitelistAppIdLiquidation(ctx sdk.Context, appMappingId uint64) error {
	WhitelistedAppIds := k.GetAppIds(ctx).WhitelistedAppMappingIds
	found := uint64InSlice(appMappingId, k.GetAppIds(ctx).WhitelistedAppMappingIds)
	if !found {
		return types.ErrAppIdDoesNotExists
	}
	var newAppIds []uint64
	for i := range WhitelistedAppIds {
		if appMappingId != WhitelistedAppIds[i] {
			newAppId := WhitelistedAppIds[i]
			newAppIds = append(newAppIds, newAppId)
		}
	}
	UpdatedWhitelistedAppIds := types.WhitelistedAppIds{
		WhitelistedAppMappingIds: newAppIds,
	}

	k.SetAppId(ctx, UpdatedWhitelistedAppIds)
	return nil
}

func (k Keeper) WasmRemoveWhitelistAppIdLiquidationQuery(ctx sdk.Context, appMappingId uint64) (bool, string) {
	found := uint64InSlice(appMappingId, k.GetAppIds(ctx).WhitelistedAppMappingIds)
	if !found {
		return false, types.ErrAppIdDoesNotExists.Error()
	}
	return true, ""
}
