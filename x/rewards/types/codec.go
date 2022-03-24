package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&NewMintRewardsProposal{}, "comdex/NewMintRewardsProposal", nil)
	cdc.RegisterConcrete(&DisbaleMintRewardsProposal{}, "comdex/DisbaleMintRewardsProposal", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&NewMintRewardsProposal{},
		&DisbaleMintRewardsProposal{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_MsgService_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
