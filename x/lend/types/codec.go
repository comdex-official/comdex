package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
    cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// this line is used by starport scaffolding # 1
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgLend{}, "comdex/lend/lend", nil)
	cdc.RegisterConcrete(&MsgWithdraw{}, "comdex/lend/withdraw", nil)
	cdc.RegisterConcrete(&MsgBorrow{}, "comdex/lend/borrow", nil)
	cdc.RegisterConcrete(&MsgRepay{}, "comdex/lend/repay", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgLend{},
		&MsgWithdraw{},
		&MsgBorrow{},
		&MsgRepay{},
		)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
