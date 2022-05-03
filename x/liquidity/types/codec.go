package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary x/liquidity interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreatePair{}, "liquidity/MsgCreatePair", nil)
	cdc.RegisterConcrete(&MsgCreatePool{}, "liquidity/MsgCreatePool", nil)
	cdc.RegisterConcrete(&MsgDeposit{}, "liquidity/MsgDeposit", nil)
	cdc.RegisterConcrete(&MsgWithdraw{}, "liquidity/MsgWithdraw", nil)
	cdc.RegisterConcrete(&MsgLimitOrder{}, "liquidity/MsgLimitOrder", nil)
	cdc.RegisterConcrete(&MsgMarketOrder{}, "liquidity/MsgMarketOrder", nil)
	cdc.RegisterConcrete(&MsgCancelOrder{}, "liquidity/MsgCancelOrder", nil)
	cdc.RegisterConcrete(&MsgCancelAllOrders{}, "liquidity/MsgCancelAllOrders", nil)
	cdc.RegisterConcrete(&MsgBondPoolTokens{}, "comdex/liquidity/bond", nil)
	cdc.RegisterConcrete(&MsgUnbondPoolTokens{}, "comdex/liquidity/unbond", nil)
}

// RegisterInterfaces registers the x/liquidity interfaces types with the
// interface registry.
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreatePair{},
		&MsgCreatePool{},
		&MsgDeposit{},
		&MsgWithdraw{},
		&MsgLimitOrder{},
		&MsgMarketOrder{},
		&MsgCancelOrder{},
		&MsgCancelAllOrders{},
		&MsgBondPoolTokens{},
		&MsgUnbondPoolTokens{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino = codec.NewLegacyAmino()

	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
