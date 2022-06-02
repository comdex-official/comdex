package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// this line is used by starport scaffolding # 1
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateGauge{}, "comdex/rewards/MsgCreateGauge", nil)
	cdc.RegisterConcrete(&WhitelistAsset{}, "comdex/rewards/whitelistAsset", nil)
	cdc.RegisterConcrete(&RemoveWhitelistAsset{}, "comdex/rewards/removeWhitelistAsset", nil)
	cdc.RegisterConcrete(&WhitelistAppIdVault{}, "comdex/rewards/whitelistAppIdVault", nil)
	cdc.RegisterConcrete(&RemoveWhitelistAppIdVault{}, "comdex/rewards/removeWhitelistAppIdVault", nil)
	cdc.RegisterConcrete(&ActivateExternalRewardsLockers{}, "comdex/rewards/activateExternalRewardsLockers", nil)
	cdc.RegisterConcrete(&ActivateExternalRewardsVault{}, "comdex/rewards/activateExternalRewardsVault", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {

	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateGauge{},
		&WhitelistAsset{},
		&RemoveWhitelistAsset{},
		&WhitelistAppIdVault{},
		&RemoveWhitelistAppIdVault{},
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
