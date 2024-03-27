package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authzcodec "github.com/cosmos/cosmos-sdk/x/authz/codec"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	// this line is used by starport scaffolding # 1
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgLiquidateInternalKeeperRequest{}, "comdex/liquidationsV2/MsgLiquidateInternalKeeperRequest", nil)
	cdc.RegisterConcrete(&MsgAppReserveFundsRequest{}, "comdex/liquidationsV2/MsgAppReserveFundsRequest", nil)
	cdc.RegisterConcrete(&MsgLiquidateExternalKeeperRequest{}, "comdex/liquidationsV2/MsgLiquidateExternalKeeperRequest", nil)
	cdc.RegisterConcrete(&Params{}, "comdex/liquidationsV2/Params", nil)
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "comdex/liquidationsV2/MsgUpdateParams")
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
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	RegisterLegacyAminoCodec(authzcodec.Amino)
	// sdk.RegisterLegacyAminoCodec(amino)
	amino.Seal()
}
