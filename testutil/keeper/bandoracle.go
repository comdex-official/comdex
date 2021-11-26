package keeper

import (
	"testing"

	"github.com/comdex-official/comdex/x/bandoracle/keeper"
	"github.com/comdex-official/comdex/x/bandoracle/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	ibckeeper "github.com/cosmos/ibc-go/modules/core/keeper"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

func BandoracleKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	logger := log.NewNopLogger()

	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	appCodec := codec.NewProtoCodec(registry)
	capabilityKeeper := capabilitykeeper.NewKeeper(appCodec, storeKey, memStoreKey)

	amino := codec.NewLegacyAmino()
	ss := typesparams.NewSubspace(appCodec,
		amino,
		storeKey,
		memStoreKey,
		"BandoracleSubSpace",
	)
	IBCKeeper := ibckeeper.NewKeeper(
		appCodec,
		storeKey,
		ss,
		nil,
        nil,
		capabilityKeeper.ScopeToModule("BandoracleIBCKeeper"),
	)

	k := keeper.NewKeeper(codec.NewProtoCodec(registry), storeKey, memStoreKey,, IBCKeeper.ChannelKeeper, &IBCKeeper.PortKeeper, capabilityKeeper.ScopeToModule("BandoracleScopedKeeper"))

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, logger)
	return k, ctx
}
