package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgPlaceMarketBidRequest{}, "comdex/auctions/MsgPlaceMarketBidRequest", nil)
	cdc.RegisterConcrete(&MsgDepositLimitBidRequest{}, "comdex/auctions/MsgPlaceLimitBidRequest", nil)
	cdc.RegisterConcrete(&MsgCancelLimitBidRequest{}, "comdex/auctions/MsgCancelLimitBidRequest", nil)
	cdc.RegisterConcrete(&MsgWithdrawLimitBidRequest{}, "comdex/auctions/MsgWithdrawLimitBidRequest", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&DutchAutoBidParamsProposal{},
	)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgPlaceMarketBidRequest{},
		&MsgDepositLimitBidRequest{},
		&MsgCancelLimitBidRequest{},
		&MsgWithdrawLimitBidRequest{},
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
