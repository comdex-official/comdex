package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&UpdateLiquidationRatio{}, "comdex/UpdateLiquidationRatioProposal", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {

	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&UpdateLiquidationRatio{},
	)

	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgAddAssetRequest{},
		&MsgUpdateAssetRequest{},
		&MsgAddPairRequest{},
		&MsgUpdatePairRequest{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_MsgService_serviceDesc)
}
