package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgPlaceSurplusBidRequest{}, "comdex/auction/MsgPlaceSurplusBidRequest", nil)
	cdc.RegisterConcrete(&MsgPlaceDebtBidRequest{}, "comdex/auction/MsgPlaceDebtBidRequest", nil)
	cdc.RegisterConcrete(&MsgPlaceDutchBidRequest{}, "comdex/auction/MsgPlaceDutchBidRequest", nil)
	cdc.RegisterConcrete(&MsgPlaceDutchLendBidRequest{}, "comdex/auction/MsgPlaceDutchLendBidRequest", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgPlaceSurplusBidRequest{},
		&MsgPlaceDebtBidRequest{},
		&MsgPlaceDutchBidRequest{},
		&MsgPlaceDutchLendBidRequest{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
