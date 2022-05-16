package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgAddMarketRequest{}, "comdex/market/add-market", nil)
	cdc.RegisterConcrete(&MsgUpdateMarketRequest{}, "comdex/market/update-market", nil)
	cdc.RegisterConcrete(&MsgRemoveMarketForAssetRequest{}, "comdex/market/remove-market-for-asset", nil)

}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&UpdateAdminProposal{},
	)

	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgAddMarketRequest{},
		&MsgUpdateMarketRequest{},
		&MsgRemoveMarketForAssetRequest{},
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
