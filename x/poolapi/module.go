package liquidity

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

	"github.com/comdex-official/comdex/x/poolapi/keeper"
	"github.com/comdex-official/comdex/x/poolapi/types"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

type AppModuleBasic struct{}

func (a AppModuleBasic) Name() string {
	return types.ModuleName
}

func (a AppModuleBasic) RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {}

func (a AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
}

func (a AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return nil
}

func (a AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, message json.RawMessage) error {
	return nil
}

func (a AppModuleBasic) RegisterRESTRoutes(_ client.Context, _ *mux.Router) {}

func (a AppModuleBasic) RegisterGRPCGatewayRoutes(ctx client.Context, mux *runtime.ServeMux) {
	_ = types.RegisterQueryServiceHandlerClient(context.Background(), mux, types.NewQueryServiceClient(ctx))
}

func (a AppModuleBasic) GetTxCmd() *cobra.Command {
	return nil
}

func (a AppModuleBasic) GetQueryCmd() *cobra.Command {
	return nil
}

type AppModule struct {
	AppModuleBasic
	cdc codec.Codec
	k   keeper.Keeper
}

func NewAppModule(cdc codec.Codec, keeper keeper.Keeper) AppModule {
	return AppModule{
		cdc: cdc,
		k:   keeper,
	}
}

func (a AppModule) ConsensusVersion() uint64 {
	return 1
}

func (a AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, message json.RawMessage) []abcitypes.ValidatorUpdate {
	return nil
}

func (a AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	return nil
}

func (a AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

func (a AppModule) Route() sdk.Route {
	return sdk.NewRoute(types.RouterKey, NewHandler(a.k))
}

func (a AppModule) QuerierRoute() string {
	return types.ModuleName
}

func (a AppModule) LegacyQuerierHandler(_ *codec.LegacyAmino) sdk.Querier { return nil }

func (a AppModule) RegisterServices(configurator module.Configurator) {
	types.RegisterQueryServiceServer(configurator.QueryServer(), keeper.NewQueryServiceServer(a.k))
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
