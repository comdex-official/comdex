package app

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
)

type EncodingConfig struct {
	Amino             *codec.LegacyAmino
	InterfaceRegistry codectypes.InterfaceRegistry
	Marshaler         codec.ProtoCodecMarshaler
	TxConfig          client.TxConfig
}

func NewEncodingConfig() EncodingConfig {
	var (
		amino             = codec.NewLegacyAmino()
		interfaceRegistry = codectypes.NewInterfaceRegistry()
		marshaler         = codec.NewProtoCodec(interfaceRegistry)
		txConfig          = tx.NewTxConfig(marshaler, tx.DefaultSignModes)
	)

	return EncodingConfig{
		Amino:             amino,
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          txConfig,
	}
}

func MakeEncodingConfig() EncodingConfig {
	config := NewEncodingConfig()

	std.RegisterLegacyAminoCodec(config.Amino)
	std.RegisterInterfaces(config.InterfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(config.Amino)
	ModuleBasics.RegisterInterfaces(config.InterfaceRegistry)

	return config
}
