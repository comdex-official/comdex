package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// RegisterLegacyAminoCodec registers the necessary x/liquidity interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreatePair{}, "comdex/liquidity/MsgCreatePair", nil)
	cdc.RegisterConcrete(&MsgCreatePool{}, "comdex/liquidity/MsgCreatePool", nil)
	cdc.RegisterConcrete(&MsgCreateRangedPool{}, "comdex/liquidity/MsgCreateRangedPool", nil)
	cdc.RegisterConcrete(&MsgDeposit{}, "comdex/liquidity/MsgDeposit", nil)
	cdc.RegisterConcrete(&MsgWithdraw{}, "comdex/liquidity/MsgWithdraw", nil)
	cdc.RegisterConcrete(&MsgLimitOrder{}, "comdex/liquidity/MsgLimitOrder", nil)
	cdc.RegisterConcrete(&MsgMarketOrder{}, "comdex/liquidity/MsgMarketOrder", nil)
	cdc.RegisterConcrete(&MsgMMOrder{}, "comdex/liquidity/MsgMMOrder", nil)
	cdc.RegisterConcrete(&MsgCancelOrder{}, "comdex/liquidity/MsgCancelOrder", nil)
	cdc.RegisterConcrete(&MsgCancelAllOrders{}, "comdex/liquidity/MsgCancelAllOrders", nil)
	cdc.RegisterConcrete(&MsgCancelMMOrder{}, "comdex/liquidity/MsgCancelMMOrder", nil)
	cdc.RegisterConcrete(&MsgFarm{}, "comdex/liquidity/MsgFarm", nil)
	cdc.RegisterConcrete(&MsgUnfarm{}, "comdex/liquidity/MsgUnfarm", nil)
	cdc.RegisterConcrete(&MsgDepositAndFarm{}, "comdex/liquidity/DepositAndFarm", nil)
	cdc.RegisterConcrete(&MsgUnfarmAndWithdraw{}, "comdex/liquidity/UnfarmAndWithdraw", nil)
	cdc.RegisterConcrete(&UpdateGenericParamsProposal{}, "comdex/liquidity/UpdateGenericParamsProposal", nil)
	cdc.RegisterConcrete(&CreateNewLiquidityPairProposal{}, "comdex/liquidity/CreateNewLiquidityPairProposal", nil)
}

// RegisterInterfaces registers the x/liquidity interfaces types with the
// interface registry.
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&UpdateGenericParamsProposal{},
		&CreateNewLiquidityPairProposal{},
	)

	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreatePair{},
		&MsgCreatePool{},
		&MsgCreateRangedPool{},
		&MsgDeposit{},
		&MsgWithdraw{},
		&MsgLimitOrder{},
		&MsgMarketOrder{},
		&MsgMMOrder{},
		&MsgCancelOrder{},
		&MsgCancelAllOrders{},
		&MsgCancelMMOrder{},
		&MsgFarm{},
		&MsgUnfarm{},
		&MsgDepositAndFarm{},
		&MsgUnfarmAndWithdraw{},
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
