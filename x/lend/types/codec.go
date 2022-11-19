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
	cdc.RegisterConcrete(&MsgLend{}, "petri/lend/MsgLend", nil)
	cdc.RegisterConcrete(&MsgWithdraw{}, "petri/lend/MsgWithdraw", nil)
	cdc.RegisterConcrete(&MsgDeposit{}, "petri/lend/MsgDeposit", nil)
	cdc.RegisterConcrete(&MsgCloseLend{}, "petri/lend/MsgCloseLend", nil)
	cdc.RegisterConcrete(&MsgBorrow{}, "petri/lend/MsgBorrow", nil)
	cdc.RegisterConcrete(&MsgDepositBorrow{}, "petri/lend/MsgDepositBorrow", nil)
	cdc.RegisterConcrete(&MsgDraw{}, "petri/lend/MsgDraw", nil)
	cdc.RegisterConcrete(&MsgCloseBorrow{}, "petri/lend/MsgCloseBorrow", nil)
	cdc.RegisterConcrete(&MsgRepay{}, "petri/lend/MsgRepay", nil)
	cdc.RegisterConcrete(&MsgBorrowAlternate{}, "petri/lend/MsgBorrowAlternate", nil)
	cdc.RegisterConcrete(&MsgFundModuleAccounts{}, "petri/lend/MsgFundModuleAccounts", nil)
	cdc.RegisterConcrete(&LendPairsProposal{}, "petri/lend/LendPairsProposal", nil)
	cdc.RegisterConcrete(&AddPoolsProposal{}, "petri/lend/AddPoolsProposal", nil)
	cdc.RegisterConcrete(&AddAssetToPairProposal{}, "petri/lend/AddAssetToPairProposal", nil)
	cdc.RegisterConcrete(&AddAssetRatesParams{}, "petri/lend/AddAssetRatesParams", nil)
	cdc.RegisterConcrete(&AddAuctionParamsProposal{}, "petri/lend/AddAuctionParamsProposal", nil)
	cdc.RegisterConcrete(&MsgCalculateInterestAndRewards{}, "petri/lend/MsgCalculateInterestAndRewards", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&LendPairsProposal{},
		&AddPoolsProposal{},
		&AddAssetToPairProposal{},
		&AddAssetRatesParams{},
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
		&MsgCalculateInterestAndRewards{},
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
