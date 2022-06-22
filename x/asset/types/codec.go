package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&AddAssetsProposal{}, "comdex/asset/AddAssetsProposal", nil)
	cdc.RegisterConcrete(&UpdateAssetProposal{}, "comdex/asset/UpdateAssetProposal", nil)
	cdc.RegisterConcrete(&AddPairsProposal{}, "comdex/asset/AddPairsProposal", nil)
	cdc.RegisterConcrete(&AddWhitelistedAssetsProposal{}, "comdex/asset/AddWhitelistedAssetsProposal", nil)
	cdc.RegisterConcrete(&UpdateWhitelistedAssetProposal{}, "comdex/asset/UpdateWhitelistedAssetProposal", nil)
	cdc.RegisterConcrete(&AddWhitelistedPairsProposal{}, "comdex/asset/AddWhitelistedPairsProposal", nil)
	cdc.RegisterConcrete(&UpdateWhitelistedPairProposal{}, "comdex/asset/UpdateWhitelistedPairProposal", nil)
	cdc.RegisterConcrete(&AddAppMappingProposal{}, "comdex/asset/AddAppMappingProposal", nil)
	cdc.RegisterConcrete(&AddAssetMappingProposal{}, "comdex/asset/AddAssetMappingProposal", nil)
	cdc.RegisterConcrete(&AddExtendedPairsVaultProposal{}, "comdex/asset/AddExtendedPairsVaultProposal", nil)
	cdc.RegisterConcrete(&UpdateGovTimeInAppMappingProposal{}, "comdex/asset/UpdateGovTimeInAppMappingProposal", nil)
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
		&AddAssetMappingProposal{},
		&AddExtendedPairsVaultProposal{},
		&UpdateGovTimeInAppMappingProposal{},
	)

	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgAddAssetRequest{},
		&MsgUpdateAssetRequest{},
		&MsgAddPairRequest{},
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
