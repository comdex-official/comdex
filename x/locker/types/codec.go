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
	// this line is used by starport scaffolding # 2
	cdc.RegisterConcrete(&MsgCreateLockerRequest{}, "comdex/locker/MsgCreateLockerRequest", nil)
	cdc.RegisterConcrete(&MsgDepositAssetRequest{}, "comdex/locker/MsgDepositAssetRequest", nil)
	cdc.RegisterConcrete(&MsgWithdrawAssetRequest{}, "comdex/locker/MsgWithdrawAssetRequest", nil)
	cdc.RegisterConcrete(&MsgAddWhiteListedAssetRequest{}, "comdex/locker/MsgAddWhiteListedAssetRequest", nil)
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
