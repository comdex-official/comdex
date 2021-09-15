package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/capability/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
	tmdb "github.com/tendermint/tm-db"
	"testing"
)

func setupKeeper(t testing.TB) *Keeper {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	registry := codectypes.NewInterfaceRegistry()
	paramspace := paramtypes.NewSubspace(codec.NewProtoCodec(registry), nil, storeKey, memStoreKey, "cdp_key")
	keeper := NewKeeper(codec.NewProtoCodec(registry), storeKey, memStoreKey, nil, nil,paramspace)
	return keeper
}


func setupctx(t testing.TB) sdk.Context {
	sdk.NewKVStoreKey(types.StoreKey)
	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	ctx := sdk.NewContext(stateStore, tmproto.Header{Height: 1, Time: tmtime.Now()}, true, log.NewNopLogger())
	return ctx
}