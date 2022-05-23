package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgPlaceBidRequest{}, "comdex/auction/bid", nil)
	cdc.RegisterConcrete(&MsgPlaceDebtBidRequest{}, "comdex/auction/debt", nil)
	cdc.RegisterConcrete(&MsgPlaceDutchBidRequest{}, "comdex/auction/dutch", nil)
}

func RegisterCodec(cdc *codec.LegacyAmino) {
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgPlaceBidRequest{},
		&MsgPlaceDebtBidRequest{},
		&MsgPlaceDutchBidRequest{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_MsgService_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	amino.Seal()
}
