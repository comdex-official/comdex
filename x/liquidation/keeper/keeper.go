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

func (k Keeper) WhitelistAppID(ctx sdk.Context, appID uint64) error {
	found := uint64InSlice(appID, k.GetAppIds(ctx).WhitelistedAppIds)
	if found {
		return types.ErrAppIDExists
	}
	WhitelistedAppIds := append(k.GetAppIds(ctx).WhitelistedAppIds, appID)
	UpdatedWhitelistedAppIds := types.WhitelistedAppIds{
		WhitelistedAppIds: WhitelistedAppIds,
	}
	k.SetAppID(ctx, UpdatedWhitelistedAppIds)
	return nil
}

func (k Keeper) RemoveWhitelistAsset(ctx sdk.Context, appMappingID uint64) error {
	WhitelistedAppIds := k.GetAppIds(ctx).WhitelistedAppIds
	found := uint64InSlice(appMappingID, k.GetAppIds(ctx).WhitelistedAppIds)
	if !found {
		return types.ErrAppIDDoesNotExists
	}
	var newAppIds []uint64
	for i := range WhitelistedAppIds {
		if appMappingID != WhitelistedAppIds[i] {
			newAppID := WhitelistedAppIds[i]
			newAppIds = append(newAppIds, newAppID)
		}
	}
	UpdatedWhitelistedAppIds := types.WhitelistedAppIds{
		WhitelistedAppIds: newAppIds,
	}

	k.SetAppID(ctx, UpdatedWhitelistedAppIds)
	return nil
}

//Wasm tx and query binding functions

func (k Keeper) WasmWhitelistAppIDLiquidation(ctx sdk.Context, appID uint64) error {
	WhitelistedAppIds := append(k.GetAppIds(ctx).WhitelistedAppIds, appID)
	UpdatedWhitelistedAppIds := types.WhitelistedAppIds{
		WhitelistedAppIds: WhitelistedAppIds,
	}
	k.SetAppID(ctx, UpdatedWhitelistedAppIds)
	return nil
}

func (k Keeper) WasmWhitelistAppIDLiquidationQuery(ctx sdk.Context, appID uint64) (bool, string) {
	found := uint64InSlice(appID, k.GetAppIds(ctx).WhitelistedAppIds)
	if found {
		return false, types.ErrAppIDExists.Error()
	}
	return true, ""
}

func (k Keeper) WasmRemoveWhitelistAppIDLiquidation(ctx sdk.Context, appID uint64) error {
	WhitelistedAppIds := k.GetAppIds(ctx).WhitelistedAppIds
	found := uint64InSlice(appID, k.GetAppIds(ctx).WhitelistedAppIds)
	if !found {
		return types.ErrAppIDDoesNotExists
	}
	var newAppIds []uint64
	for i := range WhitelistedAppIds {
		if appID != WhitelistedAppIds[i] {
			newAppID := WhitelistedAppIds[i]
			newAppIds = append(newAppIds, newAppID)
		}
	}
	UpdatedWhitelistedAppIds := types.WhitelistedAppIds{
		WhitelistedAppIds: newAppIds,
	}

	k.SetAppID(ctx, UpdatedWhitelistedAppIds)
	return nil
}

func (k Keeper) WasmRemoveWhitelistAppIDLiquidationQuery(ctx sdk.Context, appID uint64) (bool, string) {
	found := uint64InSlice(appID, k.GetAppIds(ctx).WhitelistedAppIds)
	if !found {
		return false, types.ErrAppIDDoesNotExists.Error()
	}
	return true, ""
}
