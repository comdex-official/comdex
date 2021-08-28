package asset

import (
	"context"
	"encoding/json"
	"math"
	"math/rand"

	bandpacket "github.com/bandprotocol/bandchain-packet/packet"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/simulation"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	ibcporttypes "github.com/cosmos/ibc-go/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/modules/core/exported"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	abcitypes "github.com/tendermint/tendermint/abci/types"

	"github.com/comdex-official/comdex/x/asset/client/cli"
	"github.com/comdex-official/comdex/x/asset/keeper"
	"github.com/comdex-official/comdex/x/asset/types"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
	_ ibcporttypes.IBCModule     = AppModule{}
)

type AppModuleBasic struct{}

func (a AppModuleBasic) Name() string {
	return types.ModuleName
}

func (a AppModuleBasic) RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {}

func (a AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

func (a AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

func (a AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, message json.RawMessage) error {
	var state types.GenesisState
	if err := cdc.UnmarshalJSON(message, &state); err != nil {
		return err
	}

	return types.ValidateGenesis(&state)
}

func (a AppModuleBasic) RegisterRESTRoutes(_ client.Context, _ *mux.Router) {}

func (a AppModuleBasic) RegisterGRPCGatewayRoutes(ctx client.Context, mux *runtime.ServeMux) {
	_ = types.RegisterQueryServiceHandlerClient(context.Background(), mux, types.NewQueryServiceClient(ctx))
}

func (a AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

func (a AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

type AppModule struct {
	AppModuleBasic
	cdc    codec.JSONCodec
	keeper keeper.Keeper
}

func (a AppModule) ConsensusVersion() uint64 {
	return 1
}

func (a AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, message json.RawMessage) []abcitypes.ValidatorUpdate {
	var state types.GenesisState
	cdc.MustUnmarshalJSON(message, &state)
	InitGenesis(ctx, a.keeper, &state)

	return nil
}

func (a AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(ExportGenesis(ctx, a.keeper))
}

func (a AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

func (a AppModule) Route() sdk.Route {
	return sdk.NewRoute(types.RouterKey, NewHandler(a.keeper))
}

func (a AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

func (a AppModule) LegacyQuerierHandler(_ *codec.LegacyAmino) sdk.Querier { return nil }

func (a AppModule) RegisterServices(configurator module.Configurator) {
	types.RegisterMsgServiceServer(configurator.MsgServer(), keeper.NewMsgServiceServer(a.keeper))
	types.RegisterQueryServiceServer(configurator.QueryServer(), keeper.NewQueryServiceServer(a.keeper))
}

func (a AppModule) BeginBlock(_ sdk.Context, _ abcitypes.RequestBeginBlock) {}

func (a AppModule) EndBlock(_ sdk.Context, _ abcitypes.RequestEndBlock) []abcitypes.ValidatorUpdate {
	return nil
}

func (a AppModule) GenerateGenesisState(_ *module.SimulationState) {}

func (a AppModule) ProposalContents(_ module.SimulationState) []simulation.WeightedProposalContent {
	return nil
}

func (a AppModule) RandomizedParams(_ *rand.Rand) []simulation.ParamChange {
	return nil
}

func (a AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

func (a AppModule) WeightedOperations(_ module.SimulationState) []simulation.WeightedOperation {
	return nil
}

func ValidateAssetChannelParams(
	ctx sdk.Context,
	keeper keeper.Keeper,
	order ibcchanneltypes.Order,
	portID, channelID, channelVersion string,
) error {
	version := keeper.IBCVersion(ctx)
	if channelVersion != version {
		return types.ErrorInvalidVersion
	}

	port := keeper.IBCPort(ctx)
	if portID != port {
		return ibcporttypes.ErrInvalidPort
	}

	sequence, err := ibcchanneltypes.ParseChannelSequence(channelID)
	if err != nil {
		return err
	}
	if sequence > uint64(math.MaxUint32) {
		return types.ErrorMaxAssetChannels
	}
	if order != ibcchanneltypes.UNORDERED {
		return ibcchanneltypes.ErrInvalidChannelOrdering
	}

	return nil
}

func (a AppModule) OnChanOpenInit(
	ctx sdk.Context,
	order ibcchanneltypes.Order,
	_ []string,
	portID, channelID string,
	capability *capabilitytypes.Capability,
	_ ibcchanneltypes.Counterparty,
	channelVersion string,
) error {
	if err := ValidateAssetChannelParams(ctx, a.keeper, order, portID, channelID, channelVersion); err != nil {
		return err
	}

	if err := a.keeper.ClaimCapability(ctx, capability, ibchost.ChannelCapabilityPath(portID, channelID)); err != nil {
		return err
	}

	return nil
}

func (a AppModule) OnChanOpenTry(
	ctx sdk.Context,
	order ibcchanneltypes.Order,
	_ []string,
	portID, channelID string,
	capability *capabilitytypes.Capability,
	_ ibcchanneltypes.Counterparty,
	channelVersion, counterpartyVersion string,
) error {
	if counterpartyVersion != a.keeper.IBCVersion(ctx) {
		return types.ErrorInvalidVersion
	}

	if err := ValidateAssetChannelParams(ctx, a.keeper, order, portID, channelID, channelVersion); err != nil {
		return err
	}

	if !a.keeper.AuthenticateCapability(ctx, capability, ibchost.ChannelCapabilityPath(portID, channelID)) {
		if err := a.keeper.ClaimCapability(ctx, capability, ibchost.ChannelCapabilityPath(portID, channelID)); err != nil {
			return err
		}
	}

	return nil
}

func (a AppModule) OnChanOpenAck(
	ctx sdk.Context,
	_, _, counterpartyVersion string,
) error {
	version := a.keeper.IBCVersion(ctx)
	if counterpartyVersion != version {
		return types.ErrorInvalidVersion
	}

	return nil
}

func (a AppModule) OnChanOpenConfirm(
	_ sdk.Context,
	_, _ string,
) error {
	return nil
}

func (a AppModule) OnChanCloseInit(
	_ sdk.Context,
	_, _ string,
) error {
	return errors.ErrInvalidRequest
}

func (a AppModule) OnChanCloseConfirm(
	_ sdk.Context,
	_, _ string,
) error {
	return nil
}

func (a AppModule) OnRecvPacket(
	ctx sdk.Context,
	packet ibcchanneltypes.Packet,
	_ sdk.AccAddress,
) ibcexported.Acknowledgement {
	var (
		res bandpacket.OracleResponsePacketData
		ack = ibcchanneltypes.NewResultAcknowledgement([]byte{0x01})
	)

	if err := a.cdc.UnmarshalJSON(packet.GetData(), &res); err != nil {
		ack = ibcchanneltypes.NewErrorAcknowledgement(err.Error())
	}

	if ack.Success() {
		if err := a.keeper.OnRecvPacket(ctx, res); err != nil {
			ack = ibcchanneltypes.NewErrorAcknowledgement(err.Error())
		}
	}

	return ack
}

func (a AppModule) OnAcknowledgementPacket(
	_ sdk.Context,
	_ ibcchanneltypes.Packet,
	_ []byte,
	_ sdk.AccAddress,
) (*sdk.Result, error) {
	return nil, nil
}

func (a AppModule) OnTimeoutPacket(
	_ sdk.Context,
	_ ibcchanneltypes.Packet,
	_ sdk.AccAddress,
) (*sdk.Result, error) {
	return nil, nil
}
