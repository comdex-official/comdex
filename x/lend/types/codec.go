package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	// this line is used by starport scaffolding # 1
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgLend{}, "comdex/lend/lend", nil)
	cdc.RegisterConcrete(&MsgWithdraw{}, "comdex/lend/withdraw", nil)
	cdc.RegisterConcrete(&MsgBorrow{}, "comdex/lend/borrow", nil)
	cdc.RegisterConcrete(&MsgRepay{}, "comdex/lend/repay", nil)
	cdc.RegisterConcrete(&MsgFundModuleAccounts{}, "comdex/lend/fund-module", nil)
	cdc.RegisterConcrete(&AddWhitelistedAssetsProposal{}, "comdex/lend/add-whitelisted-assets", nil)
	cdc.RegisterConcrete(&UpdateWhitelistedAssetProposal{}, "comdex/lend/update-whitelisted-assets", nil)
	cdc.RegisterConcrete(&AddWhitelistedPairsProposal{},"comdex/lend/add-pairs", nil)
	cdc.RegisterConcrete(&UpdateWhitelistedPairProposal{},"comdex/lend/update-pairs", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&AddWhitelistedAssetsProposal{},
		&UpdateWhitelistedAssetProposal{},
		&AddWhitelistedPairsProposal{},
		&UpdateWhitelistedPairProposal{},
	)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgLend{},
		&MsgWithdraw{},
		&MsgBorrow{},
		&MsgRepay{},
		&MsgFundModuleAccounts{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
