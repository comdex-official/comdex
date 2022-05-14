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
	cdc.RegisterConcrete(&AddWhitelistedAssetsProposal{}, "comdex/asset/add-whitelisted-assets", nil)
	cdc.RegisterConcrete(&UpdateWhitelistedAssetProposal{}, "comdex/asset/update-whitelisted-assets", nil)
	cdc.RegisterConcrete(&AddWhitelistedPairsProposal{}, "comdex/asset/add-whitelisted-pairs", nil)
	cdc.RegisterConcrete(&UpdateWhitelistedPairProposal{}, "comdex/asset/update-pairs", nil)
	cdc.RegisterConcrete(&AddAppMappingProposal{}, "comdex/asset/add-app-mapping", nil)
	cdc.RegisterConcrete(&AddExtendedPairsVaultProposal{}, "comdex/asset/add-Extended-Pairs-vault", nil)
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
		&AddAppMappingProposal{},
		&AddExtendedPairsVaultProposal{},
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
