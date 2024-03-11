package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// RegisterLegacyAminoCodec registers the necessary x/gasless interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateGasTank{}, "comdex/gasless/MsgCreateGasTank", nil)
	cdc.RegisterConcrete(&MsgAuthorizeActors{}, "comdex/gasless/MsgAuthorizeActors", nil)
	cdc.RegisterConcrete(&MsgUpdateGasTankStatus{}, "comdex/gasless/MsgUpdateGasTankStatus", nil)
	cdc.RegisterConcrete(&MsgUpdateGasTankConfig{}, "comdex/gasless/MsgUpdateGasTankConfig", nil)
	cdc.RegisterConcrete(&MsgBlockConsumer{}, "comdex/gasless/MsgBlockConsumer", nil)
	cdc.RegisterConcrete(&MsgUnblockConsumer{}, "comdex/gasless/MsgUnblockConsumer", nil)
	cdc.RegisterConcrete(&MsgUpdateGasConsumerLimit{}, "comdex/gasless/MsgUpdateGasConsumerLimit", nil)
	cdc.RegisterConcrete(&Params{}, "comdex/gasless/Params", nil)
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "comdex/gasless/MsgUpdateParams")
}

// RegisterInterfaces registers the x/gasless interfaces types with the
// interface registry.
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
	)

	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateGasTank{},
		&MsgAuthorizeActors{},
		&MsgUpdateGasTankStatus{},
		&MsgUpdateGasTankConfig{},
		&MsgBlockConsumer{},
		&MsgUnblockConsumer{},
		&MsgUpdateGasConsumerLimit{},
		&MsgUpdateParams{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino = codec.NewLegacyAmino()

	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	// sdk.RegisterLegacyAminoCodec(amino)
	amino.Seal()
}
