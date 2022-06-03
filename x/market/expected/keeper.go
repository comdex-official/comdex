package expected

import (
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	"github.com/comdex-official/comdex/x/bandoracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v3/modules/core/exported"
)

type ChannelKeeper interface {
	ChanCloseInit(ctx sdk.Context, portID, channelID string, capability *capabilitytypes.Capability) error
	GetChannel(ctx sdk.Context, srcPort, srcChannel string) (ibcchanneltypes.Channel, bool)
	GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool)
	SendPacket(ctx sdk.Context, cap *capabilitytypes.Capability, packet ibcexported.PacketI) error
}

type PortKeeper interface {
	BindPort(ctx sdk.Context, id string) *capabilitytypes.Capability
}

type ScopedKeeper interface {
	AuthenticateCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) bool
	ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error
	GetCapability(ctx sdk.Context, name string) (*capabilitytypes.Capability, bool)
}

type AssetKeeper interface {
	GetAssetsForOracle(ctx sdk.Context, id uint64) (assettypes.Asset, bool)
	GetPair(ctx sdk.Context, id uint64) (assettypes.Pair, bool)
}

type BandOracleKeeper interface {
	GetFetchPriceResult(ctx sdk.Context, id types.OracleRequestID) (types.FetchPriceResult, error)
	GetLastFetchPriceID(ctx sdk.Context) int64
	GetLastBlockheight(ctx sdk.Context) int64
	GetFetchPriceMsg(ctx sdk.Context) types.MsgFetchPriceData
}
