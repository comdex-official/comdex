package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	// this line is used by starport scaffolding # 1
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgLiquidateInternalKeeperRequest{}, "comdex/liquidation/MsgLiquidateInternalKeeperRequest", nil)
	cdc.RegisterConcrete(&MsgAppReserveFundsRequest{}, "comdex/liquidation/MsgAppReserveFundsRequest", nil)
	cdc.RegisterConcrete(&MsgLiquidateExternalKeeperRequest{}, "comdex/liquidation/MsgLiquidateExternalKeeperRequest", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&WhitelistLiquidationProposal{},
	)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgLiquidateInternalKeeperRequest{},
		&MsgAppReserveFundsRequest{},
		&MsgLiquidateExternalKeeperRequest{},
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
