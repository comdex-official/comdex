package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&FetchPriceProposal{}, "comdex/FetchPriceProposal", nil)
}
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgFetchPriceData{}, "bandoracle/FetchPriceData", nil)
}
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgFetchPriceData{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)

	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&FetchPriceProposal{},
		)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
