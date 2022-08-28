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
	cdc.RegisterConcrete(&MsgLend{}, "comdex/lend/lend", nil)
	cdc.RegisterConcrete(&MsgWithdraw{}, "comdex/lend/withdraw", nil)
	cdc.RegisterConcrete(&MsgDeposit{}, "comdex/lend/deposit", nil)
	cdc.RegisterConcrete(&MsgCloseLend{}, "comdex/lend/close-lend", nil)
	cdc.RegisterConcrete(&MsgBorrow{}, "comdex/lend/borrow", nil)
	cdc.RegisterConcrete(&MsgDepositBorrow{}, "comdex/lend/deposit-borrow", nil)
	cdc.RegisterConcrete(&MsgDraw{}, "comdex/lend/draw", nil)
	cdc.RegisterConcrete(&MsgCloseBorrow{}, "comdex/lend/close-borrow", nil)
	cdc.RegisterConcrete(&MsgRepay{}, "comdex/lend/repay", nil)
	cdc.RegisterConcrete(&MsgBorrowAlternate{}, "comdex/lend/borrow-alternate", nil)
	cdc.RegisterConcrete(&MsgFundModuleAccounts{}, "comdex/lend/fund-module", nil)
	cdc.RegisterConcrete(&LendPairsProposal{}, "comdex/lend/add-lend-pairs", nil)
	cdc.RegisterConcrete(&AddPoolsProposal{}, "comdex/lend/add-lend-pools", nil)
	cdc.RegisterConcrete(&AddAssetToPairProposal{}, "comdex/lend/add-asset-to-pair-mapping", nil)
	cdc.RegisterConcrete(&AddAssetRatesStats{}, "comdex/lend/add-asset-rates-stats", nil)
	cdc.RegisterConcrete(&AddAuctionParamsProposal{}, "comdex/lend/add-auction-params", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&LendPairsProposal{},
		&AddPoolsProposal{},
		&AddAssetToPairProposal{},
		&AddAssetRatesStats{},
		&AddAuctionParamsProposal{},
	)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgLend{},
		&MsgWithdraw{},
		&MsgDeposit{},
		&MsgCloseLend{},
		&MsgBorrow{},
		&MsgDepositBorrow{},
		&MsgDraw{},
		&MsgCloseBorrow{},
		&MsgRepay{},
		&MsgBorrowAlternate{},
		&MsgFundModuleAccounts{},
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
