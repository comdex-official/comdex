package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateRequest{},
		&MsgDepositRequest{},
		&MsgWithdrawRequest{},
		&MsgDrawRequest{},
		&MsgRepayRequest{},
		&MsgLiquidateRequest{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_MsgService_serviceDesc)
}

var (
	ModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
)
