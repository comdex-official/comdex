package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgPlaceSurplusBidRequest{}, "comdex/auction/bid-surplus", nil)
	cdc.RegisterConcrete(&MsgPlaceDebtBidRequest{}, "comdex/auction/bid-debt", nil)
	cdc.RegisterConcrete(&MsgPlaceDutchBidRequest{}, "comdex/auction/bid-dutch", nil)
}

func RegisterCodec(cdc *codec.LegacyAmino) {
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgPlaceSurplusBidRequest{},
		&MsgPlaceDebtBidRequest{},
		&MsgPlaceDutchBidRequest{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	amino.Seal()
}
