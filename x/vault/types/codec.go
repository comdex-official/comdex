package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateRequest{}, "comdex/vault/create", nil)
	cdc.RegisterConcrete(&MsgCloseRequest{}, "comdex/vault/close", nil)
	cdc.RegisterConcrete(&MsgDepositRequest{}, "comdex/vault/deposit", nil)
	cdc.RegisterConcrete(&MsgRepayRequest{}, "comdex/vault/repay", nil)
	cdc.RegisterConcrete(&MsgWithdrawRequest{}, "comdex/vault/withdraw", nil)
	cdc.RegisterConcrete(&MsgDrawRequest{}, "comdex/vault/draw", nil)
	cdc.RegisterConcrete(&MsgCreateStableMintRequest{}, "comdex/vault/create-stable-mint", nil)
	cdc.RegisterConcrete(&MsgDepositStableMintRequest{}, "comdex/vault/deposit-stable-mint", nil)
	cdc.RegisterConcrete(&MsgWithdrawStableMintRequest{}, "comdex/vault/withdraw-stable-mint", nil)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateRequest{},
		&MsgDepositRequest{},
		&MsgWithdrawRequest{},
		&MsgDrawRequest{},
		&MsgRepayRequest{},
		&MsgCloseRequest{},
		&MsgCreateStableMintRequest{},
		&MsgDepositStableMintRequest{},
		&MsgWithdrawStableMintRequest{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_MsgService_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
)

func init() {
	RegisterLegacyAminoCodec(amino)
	amino.Seal()
}
