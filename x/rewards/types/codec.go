package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// this line is used by starport scaffolding # 1
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&WhitelistAsset{}, "comdex/rewards/whitelistAsset", nil)
	cdc.RegisterConcrete(&RemoveWhitelistAsset{}, "comdex/rewards/removeWhitelistAsset", nil)
	cdc.RegisterConcrete(&WhitelistAppIdVault{}, "comdex/rewards/whitelistAppIdVault", nil)
	cdc.RegisterConcrete(&RemoveWhitelistAppIdVault{}, "comdex/rewards/removeWhitelistAppIdVault", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {

	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&WhitelistAsset{},
		&RemoveWhitelistAsset{},
		&WhitelistAppIdVault{},
		&RemoveWhitelistAppIdVault{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)

func init() {
	RegisterCodec(Amino)
	Amino.Seal()
}
