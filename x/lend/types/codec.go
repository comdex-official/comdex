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
	cdc.RegisterConcrete(&MsgLend{}, "comdex/lend/MsgLend", nil)
	cdc.RegisterConcrete(&MsgWithdraw{}, "comdex/lend/MsgWithdraw", nil)
	cdc.RegisterConcrete(&MsgDeposit{}, "comdex/lend/MsgDeposit", nil)
	cdc.RegisterConcrete(&MsgCloseLend{}, "comdex/lend/MsgCloseLend", nil)
	cdc.RegisterConcrete(&MsgBorrow{}, "comdex/lend/MsgBorrow", nil)
	cdc.RegisterConcrete(&MsgDepositBorrow{}, "comdex/lend/MsgDepositBorrow", nil)
	cdc.RegisterConcrete(&MsgDraw{}, "comdex/lend/MsgDraw", nil)
	cdc.RegisterConcrete(&MsgCloseBorrow{}, "comdex/lend/MsgCloseBorrow", nil)
	cdc.RegisterConcrete(&MsgRepay{}, "comdex/lend/MsgRepay", nil)
	cdc.RegisterConcrete(&MsgBorrowAlternate{}, "comdex/lend/MsgBorrowAlternate", nil)
	cdc.RegisterConcrete(&MsgFundModuleAccounts{}, "comdex/lend/MsgFundModuleAccounts", nil)
	cdc.RegisterConcrete(&LendPairsProposal{}, "comdex/lend/LendPairsProposal", nil)
	cdc.RegisterConcrete(&AddPoolsProposal{}, "comdex/lend/AddPoolsProposal", nil)
	cdc.RegisterConcrete(&AddAssetToPairProposal{}, "comdex/lend/AddAssetToPairProposal", nil)
	cdc.RegisterConcrete(&AddAssetRatesStats{}, "comdex/lend/AddAssetRatesStats", nil)
	cdc.RegisterConcrete(&AddAuctionParamsProposal{}, "comdex/lend/AddAuctionParamsProposal", nil)
	cdc.RegisterConcrete(&MsgCalculateBorrowInterest{}, "comdex/lend/MsgCalculateBorrowInterest", nil)
	cdc.RegisterConcrete(&MsgCalculateLendRewards{}, "comdex/lend/MsgCalculateLendRewards", nil)
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
		&MsgCalculateBorrowInterest{},
		&MsgCalculateLendRewards{},
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
