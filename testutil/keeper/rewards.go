package keeper

//
//import (
//	"testing"
//
//	"github.com/comdex-official/comdex/x/rewards/keeper"
//	"github.com/comdex-official/comdex/x/rewards/types"
//	"github.com/cosmos/cosmos-sdk/codec"
//	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
//	"cosmossdk.io/store"
//	storetypes "cosmossdk.io/core/store"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
//	"github.com/stretchr/testify/require"
//	"cosmossdk.io/log"
//	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
//	tmdb "github.com/cometbft/cometbft-db"
//)
//
//func RewardsKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
//	storeKey := sdk.NewKVStoreKey(types.StoreKey)
//	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
//
//	db := tmdb.NewMemDB()
//	stateStore := store.NewCommitMultiStore(db)
//	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
//	stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)
//	require.NoError(t, stateStore.LoadLatestVersion())
//
//	registry := codectypes.NewInterfaceRegistry()
//	cdc := codec.NewProtoCodec(registry)
//
//	paramsSubspace := typesparams.NewSubspace(cdc,
//		types.Amino,
//		storeKey,
//		memStoreKey,
//		"RewardsParams",
//	)
//	k := keeper.NewKeeper(
//		cdc,
//		storeKey,
//		memStoreKey,
//		paramsSubspace,
//	)
//
//	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
//
//	// Initialize params
//	k.SetParams(ctx, types.DefaultParams())
//
//	return k, ctx
//}
