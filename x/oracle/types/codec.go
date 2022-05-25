package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&UpdateAdminProposal{},
	)

	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgAddMarketRequest{},
		&MsgUpdateMarketRequest{},
		&MsgFetchPriceRequest{},
		&MsgRemoveMarketForAssetRequest{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_MsgService_serviceDesc)
}

var (
	ModuleCdc = codec.NewProtoCodec(types.NewInterfaceRegistry())
)
