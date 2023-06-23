package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/comdex-official/comdex/x/liquidation/expected"
	"github.com/comdex-official/comdex/x/liquidation/types"
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
	esm        expected.EsmKeeper
	rewards    expected.RewardsKeeper
	lend       expected.LendKeeper
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
	esm expected.EsmKeeper,
	rewards expected.RewardsKeeper,
	lend expected.LendKeeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
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
		esm:        esm,
		rewards:    rewards,
		lend:       lend,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}

// Wasm tx and query binding functions
func (k Keeper) WasmWhitelistAppIDLiquidation(ctx sdk.Context, appID uint64) error {
	_, found := k.GetAppIDByAppForLiquidation(ctx, appID)
	if found {
		return types.ErrAppIDDoesNotExists
	}

	k.SetAppIDForLiquidation(ctx, appID)
	return nil
}

func (k Keeper) WasmWhitelistAppIDLiquidationQuery(ctx sdk.Context, appID uint64) (bool, string) {
	_, found := k.GetAppIDByAppForLiquidation(ctx, appID)
	if found {
		return false, types.ErrAppIDExists.Error()
	}
	return true, ""
}

func (k Keeper) WasmRemoveWhitelistAppIDLiquidation(ctx sdk.Context, appID uint64) error {
	_, found := k.GetAppIDByAppForLiquidation(ctx, appID)
	if !found {
		return types.ErrAppIDDoesNotExists
	}
	k.DeleteAppID(ctx, appID)

	return nil
}

func (k Keeper) WasmRemoveWhitelistAppIDLiquidationQuery(ctx sdk.Context, appID uint64) (bool, string) {
	_, found := k.GetAppIDByAppForLiquidation(ctx, appID)
	if !found {
		return false, types.ErrAppIDDoesNotExists.Error()
	}
	return true, ""
}
