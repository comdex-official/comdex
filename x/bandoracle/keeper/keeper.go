package keeper

import (
	"fmt"

	assetkeeper "github.com/comdex-official/comdex/x/asset/keeper"
	"github.com/comdex-official/comdex/x/bandoracle/expected"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/comdex-official/comdex/x/bandoracle/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/spm/ibckeeper"
)

type (
	Keeper struct {
		cdc           codec.BinaryCodec
		storeKey      sdk.StoreKey
		memKey        sdk.StoreKey
		paramstore    paramtypes.Subspace
		oracle        expected.OracleKeeper
		assetKeeper   assetkeeper.Keeper
		channelKeeper expected.ChannelKeeper
		portKeeper    expected.PortKeeper
		scopedKeeper  ibckeeper.ScopedKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,
	channelKeeper expected.ChannelKeeper,
	portKeeper expected.PortKeeper,
	scopedKeeper ibckeeper.ScopedKeeper,
	oracle expected.OracleKeeper,
	assetKeeper assetkeeper.Keeper,

) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		paramstore:    ps,
		oracle:        oracle,
		assetKeeper:   assetKeeper,
		channelKeeper: channelKeeper,
		portKeeper:    portKeeper,
		scopedKeeper:  scopedKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
