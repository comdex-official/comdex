package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateGauge{}, "comdex/rewards/MsgCreateGauge", nil)
	cdc.RegisterConcrete(&ActivateExternalRewardsLockers{}, "comdex/rewards/activateExternalRewardsLockers", nil)
	cdc.RegisterConcrete(&ActivateExternalRewardsVault{}, "comdex/rewards/activateExternalRewardsVault", nil)
	cdc.RegisterConcrete(&AddLendExternalRewardsProposal{}, "comdex/rewards/AddLendExternalRewardsProposal", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&AddLendExternalRewardsProposal{},
	)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateGauge{},
		&ActivateExternalRewardsLockers{},
		&ActivateExternalRewardsVault{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
