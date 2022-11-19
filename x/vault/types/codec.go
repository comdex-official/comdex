package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateRequest{}, "petri/vault/MsgCreateRequest", nil)
	cdc.RegisterConcrete(&MsgCloseRequest{}, "petri/vault/MsgCloseRequest", nil)
	cdc.RegisterConcrete(&MsgDepositRequest{}, "petri/vault/MsgDepositRequest", nil)
	cdc.RegisterConcrete(&MsgRepayRequest{}, "petri/vault/MsgRepayRequest", nil)
	cdc.RegisterConcrete(&MsgWithdrawRequest{}, "petri/vault/MsgWithdrawRequest", nil)
	cdc.RegisterConcrete(&MsgDrawRequest{}, "petri/vault/MsgDrawRequest", nil)
	cdc.RegisterConcrete(&MsgDepositAndDrawRequest{}, "petri/vault/MsgDepositAndDrawRequest", nil)
	cdc.RegisterConcrete(&MsgCreateStableMintRequest{}, "petri/vault/MsgCreateStableMintRequest", nil)
	cdc.RegisterConcrete(&MsgDepositStableMintRequest{}, "petri/vault/MsgDepositStableMintRequest", nil)
	cdc.RegisterConcrete(&MsgWithdrawStableMintRequest{}, "petri/vault/MsgWithdrawStableMintRequest", nil)
	cdc.RegisterConcrete(&MsgVaultInterestCalcRequest{}, "petri/vault/MsgVaultInterestCalcRequest", nil)
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
		&MsgDepositAndDrawRequest{},
		&MsgCreateStableMintRequest{},
		&MsgDepositStableMintRequest{},
		&MsgWithdrawStableMintRequest{},
		&MsgVaultInterestCalcRequest{},
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
	amino.Seal()
}
