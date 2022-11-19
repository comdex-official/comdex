package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgDepositESM{}, "petri/esm/deposit-esm", nil)
	cdc.RegisterConcrete(&MsgExecuteESM{}, "petri/esm/execute-esm", nil)
	cdc.RegisterConcrete(&MsgKillRequest{}, "petri/esm/stop-all-actions", nil)
	cdc.RegisterConcrete(&MsgCollateralRedemptionRequest{}, "petri/esm/redeem-collateral", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
	)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgDepositESM{},
		&MsgExecuteESM{},
		&MsgKillRequest{},
		&MsgCollateralRedemptionRequest{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
