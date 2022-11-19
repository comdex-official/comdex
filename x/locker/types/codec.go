package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateLockerRequest{}, "petri/locker/MsgCreateLockerRequest", nil)
	cdc.RegisterConcrete(&MsgDepositAssetRequest{}, "petri/locker/MsgDepositAssetRequest", nil)
	cdc.RegisterConcrete(&MsgWithdrawAssetRequest{}, "petri/locker/MsgWithdrawAssetRequest", nil)
	cdc.RegisterConcrete(&MsgAddWhiteListedAssetRequest{}, "petri/locker/MsgAddWhiteListedAssetRequest", nil)
	cdc.RegisterConcrete(&MsgCloseLockerRequest{}, "petri/locker/MsgCloseLockerRequest", nil)
	cdc.RegisterConcrete(&MsgLockerRewardCalcRequest{}, "petri/locker/MsgLockerRewardCalcRequest", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateLockerRequest{},
		&MsgWithdrawAssetRequest{},
		&MsgDepositAssetRequest{},
		&MsgAddWhiteListedAssetRequest{},
		&MsgCloseLockerRequest{},
		&MsgLockerRewardCalcRequest{},
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
