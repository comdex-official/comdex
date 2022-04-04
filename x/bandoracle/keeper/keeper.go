package keeper

import (
	"fmt"
	assetkeeper "github.com/comdex-official/comdex/x/asset/keeper"
	"github.com/comdex-official/comdex/x/bandoracle/expected"
	host "github.com/cosmos/ibc-go/v2/modules/core/24-host"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/comdex-official/comdex/x/bandoracle/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
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
		scopedKeeper  expected.ScopedKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,
	channelKeeper expected.ChannelKeeper,
	portKeeper expected.PortKeeper,
	scopedKeeper expected.ScopedKeeper,
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

// IsBound checks if the market module is already bound to the desired port
func (k Keeper) IsBound(ctx sdk.Context, portID string) bool {
	_, ok := k.scopedKeeper.GetCapability(ctx, host.PortPath(portID))
	return ok
}

// BindPort defines a wrapper function for the ort Keeper's function in
// order to expose it to module's InitGenesis function
func (k Keeper) BindPort(ctx sdk.Context, portID string) error {
	cap := k.portKeeper.BindPort(ctx, portID)
	return k.scopedKeeper.ClaimCapability(ctx, cap, host.PortPath(portID))
}

// GetPort returns the portID for the market module. Used in ExportGenesis
func (k Keeper) GetPort(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(types.PortKey))
}

// SetPort sets the portID for the market module. Used in InitGenesis
func (k Keeper) SetPort(ctx sdk.Context, portID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PortKey, []byte(portID))
}
