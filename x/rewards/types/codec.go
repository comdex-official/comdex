package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&NewMintRewardsProposal{}, "comdex/rewards/add-new-mint-rewards", nil)
	cdc.RegisterConcrete(&DisbaleMintRewardsProposal{}, "comdex/rewards/disable-mint-rewards", nil)

}

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&NewMintRewardsProposal{}, "comdex/rewards/add-new-mint-rewards", nil)
	cdc.RegisterConcrete(&DisbaleMintRewardsProposal{}, "comdex/rewards/disable-mint-rewards", nil)
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
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	amino.Seal()
}
