package asset

import (
	"context"
	"encoding/json"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/simulation"
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
)

type AppModuleBasic struct{ cdc codec.BinaryCodec }

func NewAppModule(
	//nolint
	cdc codec.Codec,
	keeper keeper.Keeper,
) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(cdc),
		keeper:         keeper,
	}
}

func NewAppModuleBasic(cdc codec.BinaryCodec) AppModuleBasic {
	return AppModuleBasic{cdc: cdc}
}

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
	//nolint
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
