package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
    cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	// this line is used by starport scaffolding # 1
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	// this line is used by starport scaffolding # 2
	cdc.RegisterConcrete(&MsgCreateLockerRequest{}, "comdex/locker/create-locker", nil)
	cdc.RegisterConcrete(&MsgDepositAssetRequest{}, "comdex/locker/deposit-locker", nil)
	cdc.RegisterConcrete(&MsgWithdrawAssetRequest{}, "comdex/locker/withdraw-locker", nil)
	cdc.RegisterConcrete(&MsgAddWhiteListedAssetRequest{}, "comdex/locker/whitelist-asset-locker", nil)
} 

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateLockerRequest{},
		&MsgWithdrawAssetRequest{},
		&MsgDepositAssetRequest{},
		&MsgAddWhiteListedAssetRequest{},
	)
	// msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
	msgservice.RegisterMsgServiceDesc(registry, &_MsgService_serviceDesc)
}

var (
	Amino = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
func init() {
	RegisterCodec(Amino)
	Amino.Seal()
}