package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&AddAssetsProposal{}, "comdex/asset/add-assets", nil)
	cdc.RegisterConcrete(&UpdateAssetProposal{}, "comdex/asset/update-asset", nil)
	cdc.RegisterConcrete(&AddPairsProposal{}, "comdex/asset/add-pairs", nil)
	cdc.RegisterConcrete(&AddWhitelistedAssetsProposal{}, "comdex/lend/add-whitelisted-assets", nil)
	cdc.RegisterConcrete(&UpdateWhitelistedAssetProposal{}, "comdex/lend/update-whitelisted-assets", nil)
	cdc.RegisterConcrete(&AddWhitelistedPairsProposal{},"comdex/lend/add-pairs", nil)
	cdc.RegisterConcrete(&UpdateWhitelistedPairProposal{},"comdex/lend/update-pairs", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&AddAssetsProposal{},
		&UpdateAssetProposal{},
		&AddPairsProposal{},
		&AddWhitelistedAssetsProposal{},
		&UpdateWhitelistedAssetProposal{},
		&AddWhitelistedPairsProposal{},
		&UpdateWhitelistedPairProposal{},
	)

	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgAddAssetRequest{},
		&MsgUpdateAssetRequest{},
		&MsgAddPairRequest{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_MsgService_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	amino.Seal()
}
