package keeper

import (
	"context"
	"fmt"

	"github.com/comdex-official/comdex/x/auctionsV2/expected"
	"github.com/comdex-official/comdex/x/auctionsV2/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"
)

type (
	Keeper struct {
		cdc            codec.BinaryCodec
		storeKey       sdk.StoreKey
		memKey         sdk.StoreKey
		paramstore     paramtypes.Subspace
		LiquidationsV2 expected.LiquidationsV2Keeper
		bankKeeper     types.BankKeeper
		market         expected.MarketKeeper
		asset          expected.AssetKeeper
	}
)

func (k Keeper) Params(ctx context.Context, request *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,
	LiquidationsV2Keeper expected.LiquidationsV2Keeper,
	bankKeeper types.BankKeeper,
	marketKeeper expected.MarketKeeper,
	assetKeeper expected.AssetKeeper,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		memKey:         memKey,
		paramstore:     ps,
		LiquidationsV2: LiquidationsV2Keeper,
		bankKeeper:     bankKeeper,
		market:         marketKeeper,
		asset:          assetKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
func (k Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}

func (k Keeper) SetAuctionParams(ctx sdk.Context, auctionParams types.AuctionParams) {
	var (
		store = k.Store(ctx)
		key   = types.AuctionParamsKey
		value = k.cdc.MustMarshal(&auctionParams)
	)

	store.Set(key, value)
}

func (k Keeper) GetAuctionParams(ctx sdk.Context) (auctionParams types.AuctionParams, found bool) {
	key := types.AuctionParamsKey

	var (
		store = k.Store(ctx)

		value = store.Get(key)
	)

	if value == nil {
		return auctionParams, false
	}

	k.cdc.MustUnmarshal(value, &auctionParams)
	return auctionParams, true
}
