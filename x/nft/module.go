package nft

import (
	"context"
	// "encoding/json"
	// "github.com/gogo/protobuf/grpc"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/comdex-official/comdex/x/nft/client/cli"
	"github.com/comdex-official/comdex/x/nft/client/rest"
	"github.com/comdex-official/comdex/x/nft/keeper"
	"github.com/comdex-official/comdex/x/nft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	// "github.com/cosmos/cosmos-sdk/types/module"
)

var (
	// _ module.AppModule      = AppModule{}
	// _ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct {
	cdc codec.Codec
}

func (AppModuleBasic) Name() string { return types.ModuleName }

func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

// func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
// 	return cdc.MustMarshalJSON(DefaultGenesisState())
// }

// func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
// 	var data types.GenesisState
// 	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
// 		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
// 	}

// 	return types.ValidateGenesis(data)
// }

// func (AppModuleBasic) RegisterGRPCRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
// 	types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
// }

func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

type AppModule struct {
	AppModuleBasic
	keeper keeper.Keeper
}

// func (am AppModule) RegisterQueryService(server grpc.Server) {
// 	types.RegisterQueryServer(server, am.keeper)
// }

func NewAppModule(cdc codec.Codec, keeper keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},
		keeper:         keeper,
	}
}

func (AppModule) Name() string { return types.ModuleName }

func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
}

func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(types.RouterKey, NewHandler(am.keeper))
}

func (AppModule) QuerierRoute() string { return types.RouterKey }

// func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
// 	return keeper.NewQuerier(am.keeper, legacyQuerierCdc)
// }

func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr *mux.Router) {
	rest.RegisterHandlers(clientCtx, rtr, types.RouterKey)
}

func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
}
// func (am AppModule) RegisterServices(cfg module.Configurator) {
// 	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)
// }

// func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
// 	// var genesisState types.GenesisState

// 	cdc.MustUnmarshalJSON(data, &genesisState)

// 	// InitGenesis(ctx, am.keeper, genesisState)
// 	return []abci.ValidatorUpdate{}
// }

// func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
// 	// gs := ExportGenesis(ctx, am.keeper)
// 	return cdc.MustMarshalJSON(gs)
// }
func (AppModule) ConsensusVersion() uint64 { return 1 }

func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
