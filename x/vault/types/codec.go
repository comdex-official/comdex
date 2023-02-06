package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateRequest{}, "comdex/vault/MsgCreateRequest", nil)
	cdc.RegisterConcrete(&MsgCloseRequest{}, "comdex/vault/MsgCloseRequest", nil)
	cdc.RegisterConcrete(&MsgDepositRequest{}, "comdex/vault/MsgDepositRequest", nil)
	cdc.RegisterConcrete(&MsgRepayRequest{}, "comdex/vault/MsgRepayRequest", nil)
	cdc.RegisterConcrete(&MsgWithdrawRequest{}, "comdex/vault/MsgWithdrawRequest", nil)
	cdc.RegisterConcrete(&MsgDrawRequest{}, "comdex/vault/MsgDrawRequest", nil)
	cdc.RegisterConcrete(&MsgDepositAndDrawRequest{}, "comdex/vault/MsgDepositAndDrawRequest", nil)
	cdc.RegisterConcrete(&MsgCreateStableMintRequest{}, "comdex/vault/MsgCreateStableMintRequest", nil)
	cdc.RegisterConcrete(&MsgDepositStableMintRequest{}, "comdex/vault/MsgDepositStableMintRequest", nil)
	cdc.RegisterConcrete(&MsgWithdrawStableMintRequest{}, "comdex/vault/MsgWithdrawStableMintRequest", nil)
	cdc.RegisterConcrete(&MsgVaultInterestCalcRequest{}, "comdex/vault/MsgVaultInterestCalcRequest", nil)
	cdc.RegisterConcrete(&MsgCorrectStabilityFeesRequest{}, "comdex/vault/MsgCorrectStabilityFeesRequest", nil) // need to remove later
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
		&MsgCorrectStabilityFeesRequest{}, // need to remove later
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
