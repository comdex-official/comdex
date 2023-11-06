package main

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
	// bump47: uncomment below 1
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/router/types"
)

func Migrate(appState types.AppMap, clientCtx client.Context) types.AppMap {
	// Migrate packetfowardmiddleware.
	if appState["packetfowardmiddleware"] != nil {
		// unmarshal relative source genesis application state
		var oldGovState packetforwardtypes.GenesisState
		clientCtx.Codec.MustUnmarshalJSON(appState["packetfowardmiddleware"], &oldGovState)

		// delete deprecated packetfowardmiddleware genesis state
		delete(appState, "packetfowardmiddleware")

		// Migrate relative source genesis application state and marshal it into
		// the respective key.
		appState["packetfowardmiddleware"] = clientCtx.Codec.MustMarshalJSON(MigrateJSON(&oldGovState))
	}

	return appState
}

func MigrateJSON(oldState *packetforwardtypes.GenesisState) *packetforwardtypes.GenesisState {
	return &packetforwardtypes.GenesisState{
		Params:          oldState.Params,
		InFlightPackets: oldState.InFlightPackets,
	}
}
