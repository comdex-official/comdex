package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil))

	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&FetchPriceProposal{},
	)
}

var (
	Amino     = codec.NewLegacyAmino()
	// ModuleCdc = codec.NewAminoCodec(Amino)
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
