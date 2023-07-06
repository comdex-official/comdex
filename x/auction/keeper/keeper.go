package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/comdex-official/comdex/x/auction/expected"
	"github.com/comdex-official/comdex/x/auction/types"
)

type (
	Keeper struct {
		cdc         codec.BinaryCodec
		storeKey    storetypes.StoreKey
		memKey      storetypes.StoreKey
		paramstore  paramtypes.Subspace
		account     expected.AccountKeeper
		bank        expected.BankKeeper
		market      expected.MarketKeeper
		liquidation expected.LiquidationKeeper
		asset       expected.AssetKeeper
		vault       expected.VaultKeeper
		collector   expected.CollectorKeeper
		tokenMint   expected.TokenMintKeeper
		esm         expected.EsmKeeper
		lend        expected.LendKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	account expected.AccountKeeper,
	bank expected.BankKeeper,
	market expected.MarketKeeper,
	liquidation expected.LiquidationKeeper,
	asset expected.AssetKeeper,
	vault expected.VaultKeeper,
	collector expected.CollectorKeeper,
	tokenMintKeeper expected.TokenMintKeeper,
	esm expected.EsmKeeper,
	lend expected.LendKeeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:         cdc,
		storeKey:    storeKey,
		memKey:      memKey,
		paramstore:  ps,
		account:     account,
		bank:        bank,
		market:      market,
		liquidation: liquidation,
		asset:       asset,
		vault:       vault,
		collector:   collector,
		tokenMint:   tokenMintKeeper,
		esm:         esm,
		lend:        lend,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}
