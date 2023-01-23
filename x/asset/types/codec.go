package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&AddAssetsProposal{}, "comdex/asset/AddAssetsProposal", nil)
	cdc.RegisterConcrete(&AddMultipleAssetsProposal{}, "comdex/asset/AddMultipleAssetsProposal", nil)
	cdc.RegisterConcrete(&UpdateAssetProposal{}, "comdex/asset/UpdateAssetProposal", nil)
	cdc.RegisterConcrete(&AddPairsProposal{}, "comdex/asset/AddPairsProposal", nil)
	cdc.RegisterConcrete(&AddMultiplePairsProposal{}, "comdex/asset/AddMultiplePairsProposal", nil)
	cdc.RegisterConcrete(&UpdatePairProposal{}, "comdex/asset/UpdatePairProposal", nil)
	cdc.RegisterConcrete(&AddAppProposal{}, "comdex/asset/AddAppProposal", nil)
	cdc.RegisterConcrete(&AddAssetInAppProposal{}, "comdex/asset/AddAssetInAppProposal", nil)
	cdc.RegisterConcrete(&UpdateGovTimeInAppProposal{}, "comdex/asset/UpdateGovTimeInAppProposal", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&AddAssetsProposal{},
		&AddMultipleAssetsProposal{},
		&UpdateAssetProposal{},
		&AddPairsProposal{},
		&AddMultiplePairsProposal{},
		&UpdatePairProposal{},
		&AddAppProposal{},
		&AddAssetInAppProposal{},
		&UpdateGovTimeInAppProposal{},
	)

	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
	)
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
