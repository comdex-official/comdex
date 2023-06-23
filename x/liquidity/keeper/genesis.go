package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/liquidity/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	if err := genState.Validate(); err != nil {
		panic(err)
	}
	k.SetParams(ctx, genState.Params)

	for _, appState := range genState.AppGenesisState {
		k.SetGenericParams(ctx, appState.GenericParams)
		k.SetLastPairID(ctx, appState.AppId, appState.LastPairId)
		k.SetLastPoolID(ctx, appState.AppId, appState.LastPoolId)

		for _, pair := range appState.Pairs {
			k.SetPair(ctx, pair)
			k.SetPairIndex(ctx, appState.AppId, pair.BaseCoinDenom, pair.QuoteCoinDenom, pair.Id)
			k.SetPairLookupIndex(ctx, appState.AppId, pair.BaseCoinDenom, pair.QuoteCoinDenom, pair.Id)
			k.SetPairLookupIndex(ctx, appState.AppId, pair.QuoteCoinDenom, pair.BaseCoinDenom, pair.Id)
		}

		for _, pool := range appState.Pools {
			k.SetPool(ctx, pool)
			k.SetPoolByReserveIndex(ctx, pool)
			k.SetPoolsByPairIndex(ctx, pool)
		}

		for _, req := range appState.DepositRequests {
			k.SetDepositRequest(ctx, req)
			k.SetDepositRequestIndex(ctx, req)
		}

		for _, req := range appState.WithdrawRequests {
			k.SetWithdrawRequest(ctx, req)
			k.SetWithdrawRequestIndex(ctx, req)
		}

		for _, order := range appState.Orders {
			k.SetOrder(ctx, appState.AppId, order)
			k.SetOrderIndex(ctx, appState.AppId, order)
		}

		for _, mmOrderIndex := range appState.MarketMakingOrderIndexes {
			k.SetMMOrderIndex(ctx, appState.AppId, mmOrderIndex)
		}

		for _, activeFarmer := range appState.ActiveFarmers {
			k.SetActiveFarmer(ctx, activeFarmer)
		}

		for _, queuedFarmer := range appState.QueuedFarmers {
			k.SetQueuedFarmer(ctx, queuedFarmer)
		}
	}
}

func (k Keeper) GetActiveAndQueuedFarmersForGenesis(ctx sdk.Context, appID uint64) ([]types.ActiveFarmer, []types.QueuedFarmer) {
	allActiveFarmers := []types.ActiveFarmer{}
	allQueuedFarmers := []types.QueuedFarmer{}
	allPools := k.GetAllPools(ctx, appID)
	for _, pool := range allPools {
		activeFarmers := k.GetAllActiveFarmers(ctx, appID, pool.Id)
		queuedFarmers := k.GetAllQueuedFarmers(ctx, appID, pool.Id)

		for _, activeFarmer := range activeFarmers {
			allActiveFarmers = append(allActiveFarmers, types.ActiveFarmer{
				AppId:          appID,
				PoolId:         pool.Id,
				Farmer:         activeFarmer.Farmer,
				FarmedPoolCoin: activeFarmer.FarmedPoolCoin,
			})
		}

		for _, queuedFarmer := range queuedFarmers {
			allQueuedFarmers = append(allQueuedFarmers, types.QueuedFarmer{
				AppId:      appID,
				PoolId:     pool.Id,
				Farmer:     queuedFarmer.Farmer,
				QueudCoins: queuedFarmer.QueudCoins,
			})
		}
	}
	return allActiveFarmers, allQueuedFarmers
}

// ExportGenesis returns the capability module's exported genesis.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	allApps, found := k.assetKeeper.GetApps(ctx)

	appGenesisState := []types.AppGenesisState{}

	if found {
		for _, app := range allApps {
			genericParams, err := k.GetGenericParams(ctx, app.Id)
			if err != nil {
				genericParams = types.DefaultGenericParams(app.Id)
			}
			allActiveFarmers, allQueuedFarmers := k.GetActiveAndQueuedFarmersForGenesis(ctx, app.Id)
			appGenesisState = append(appGenesisState, types.AppGenesisState{
				AppId:                    app.Id,
				GenericParams:            genericParams,
				LastPairId:               k.GetLastPairID(ctx, app.Id),
				LastPoolId:               k.GetLastPoolID(ctx, app.Id),
				Pairs:                    k.GetAllPairs(ctx, app.Id),
				Pools:                    k.GetAllPools(ctx, app.Id),
				DepositRequests:          k.GetAllDepositRequests(ctx, app.Id),
				WithdrawRequests:         k.GetAllWithdrawRequests(ctx, app.Id),
				Orders:                   k.GetAllOrders(ctx, app.Id),
				ActiveFarmers:            allActiveFarmers,
				QueuedFarmers:            allQueuedFarmers,
				MarketMakingOrderIndexes: k.GetAllMMOrderIndexes(ctx, app.Id),
			})
		}
	}

	return &types.GenesisState{
		Params:          k.GetParams(ctx),
		AppGenesisState: appGenesisState,
	}
}
