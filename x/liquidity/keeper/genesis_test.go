package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/liquidity"
	"github.com/comdex-official/comdex/x/liquidity/types"
)

func (s *KeeperTestSuite) TestDefaultGenesis() {
	genState := *types.DefaultGenesis()

	s.keeper.InitGenesis(s.ctx, genState)
	got := s.keeper.ExportGenesis(s.ctx)
	s.Require().Equal(genState, *got)
}

func (s *KeeperTestSuite) TestImportExportGenesis() {
	s.ctx = s.ctx.WithBlockHeight(1).WithBlockTime(utils.ParseTime("2022-01-01T00:00:00Z"))
	k, ctx := s.keeper, s.ctx

	appID1 := s.CreateNewApp("appone")
	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 2000000)
	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, s.addr(0), "1000000denom1,1000000denom2")
	s.Deposit(appID1, pool.Id, s.addr(1), "1000000denom1,1000000denom2")
	s.nextBlock()

	poolCoin := s.getBalance(s.addr(1), pool.PoolCoinDenom)
	poolCoin.Amount = poolCoin.Amount.QuoRaw(2)
	s.Withdraw(appID1, pool.Id, s.addr(1), poolCoin)
	s.nextBlock()

	s.Farm(appID1, pool.Id, s.addr(3), "3330000pool1-1")
	s.ctx = s.ctx.WithBlockTime(ctx.BlockTime().Add(time.Hour * 25))
	s.nextBlock()

	s.Farm(appID1, pool.Id, s.addr(4), "4440000pool1-1")

	s.LimitOrder(appID1, s.addr(2), pair.Id, types.OrderDirectionBuy, utils.ParseDec("1.0"), sdk.NewInt(10000), 0)
	s.nextBlock()

	depositReq := s.Deposit(appID1, pool.Id, s.addr(3), "1000000denom1,1000000denom2")
	withdrawReq := s.Withdraw(appID1, pool.Id, s.addr(1), poolCoin)
	order := s.LimitOrder(appID1, s.addr(3), pair.Id, types.OrderDirectionSell, utils.ParseDec("1.0"), sdk.NewInt(10000), 0)

	pair, _ = k.GetPair(ctx, pair.AppId, pair.Id)
	pool, _ = k.GetPool(ctx, pool.AppId, pool.Id)
	allActiveFarmers := k.GetAllActiveFarmers(ctx, appID1, pool.Id)
	allQueuedFarmers := k.GetAllQueuedFarmers(ctx, appID1, pool.Id)

	genState := k.ExportGenesis(ctx)

	bz := s.app.AppCodec().MustMarshal(genState)

	s.SetupTest()
	s.ctx = s.ctx.WithBlockHeight(1).WithBlockTime(utils.ParseTime("2022-01-01T00:00:00Z"))
	k, ctx = s.keeper, s.ctx

	var genState2 types.GenesisState
	s.app.AppCodec().MustUnmarshal(bz, &genState2)
	k.InitGenesis(ctx, genState2)

	s.Require().Equal(genState.Params, genState2.Params)
	s.Require().Equal(len(genState.AppGenesisState), len(genState2.AppGenesisState))

	importedPair, found := k.GetPair(ctx, pair.AppId, pair.Id)
	s.Require().True(found)
	s.Require().Equal(pair, importedPair)

	importedPool, found := k.GetPool(ctx, pool.AppId, pool.Id)
	s.Require().True(found)
	s.Require().Equal(pool, importedPool)

	depositReq2, found := k.GetDepositRequest(ctx, depositReq.AppId, depositReq.PoolId, depositReq.Id)
	s.Require().True(found)
	s.Require().Equal(depositReq, depositReq2)
	withdrawReq2, found := k.GetWithdrawRequest(ctx, withdrawReq.AppId, withdrawReq.PoolId, withdrawReq.Id)
	s.Require().True(found)
	s.Require().Equal(withdrawReq, withdrawReq2)
	order2, found := k.GetOrder(ctx, order.AppId, order.PairId, order.Id)
	s.Require().True(found)
	s.Require().Equal(order, order2)

	importedAllActiveFarmers := k.GetAllActiveFarmers(ctx, appID1, pool.Id)
	s.Require().Equal(len(allActiveFarmers), len(importedAllActiveFarmers))
	s.Require().Equal(allActiveFarmers, importedAllActiveFarmers)

	importedAllQueuedFarmers := k.GetAllQueuedFarmers(ctx, appID1, pool.Id)
	s.Require().Equal(len(allQueuedFarmers), len(importedAllQueuedFarmers))
	s.Require().Equal(allQueuedFarmers, importedAllQueuedFarmers)
}

func (s *KeeperTestSuite) TestImportExportGenesisEmpty() {
	k, ctx := s.keeper, s.ctx
	genState := k.ExportGenesis(ctx)

	var genState2 types.GenesisState
	bz := s.app.AppCodec().MustMarshal(genState)
	s.app.AppCodec().MustUnmarshal(bz, &genState2)
	k.InitGenesis(ctx, genState2)

	genState3 := k.ExportGenesis(ctx)
	s.Require().Equal(genState.Params, genState2.Params)
	s.Require().Equal(len(genState.AppGenesisState), len(genState2.AppGenesisState))
	s.Require().Equal(genState2.Params, genState3.Params)
	s.Require().Equal(len(genState2.AppGenesisState), len(genState3.AppGenesisState))
}

func (s *KeeperTestSuite) TestIndexesAfterImport() {
	s.ctx = s.ctx.WithBlockHeight(1).WithBlockTime(utils.ParseTime("2022-03-01T00:00:00Z"))

	appID1 := s.CreateNewApp("appone")
	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 2000000)
	asset3 := s.CreateNewAsset("ASSETTWO", "denom3", 3000000)

	pair1 := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)
	pair2 := s.CreateNewLiquidityPair(appID1, s.addr(1), asset2.Denom, asset3.Denom)

	pool1 := s.CreateNewLiquidityPool(appID1, pair1.Id, s.addr(2), "1000000denom1,1000000denom2")
	pool2 := s.CreateNewLiquidityPool(appID1, pair2.Id, s.addr(3), "1000000denom2,1000000denom3")

	s.Deposit(appID1, pool1.Id, s.addr(4), "1000000denom1,1000000denom2")
	s.Deposit(appID1, pool2.Id, s.addr(5), "1000000denom2,1000000denom3")

	liquidity.EndBlocker(s.ctx, s.keeper)
	liquidity.BeginBlocker(s.ctx, s.keeper)

	depositReq1 := s.Deposit(appID1, pool1.Id, s.addr(4), "1000000denom1,1000000denom2")
	depositReq2 := s.Deposit(appID1, pool2.Id, s.addr(5), "1000000denom2,1000000denom3")

	withdrawReq1 := s.Withdraw(appID1, pool1.Id, s.addr(4), utils.ParseCoin("1000000pool1-1"))
	withdrawReq2 := s.Withdraw(appID1, pool2.Id, s.addr(5), utils.ParseCoin("1000000pool1-2"))

	order1 := s.LimitOrder(appID1, s.addr(6), pair1.Id, types.OrderDirectionBuy, utils.ParseDec("1.0"), sdk.NewInt(10000), time.Minute)
	order2 := s.LimitOrder(appID1, s.addr(7), pair2.Id, types.OrderDirectionSell, utils.ParseDec("1.0"), sdk.NewInt(10000), time.Minute)

	liquidity.EndBlocker(s.ctx, s.keeper)

	genState := s.keeper.ExportGenesis(s.ctx)
	s.SetupTest()
	s.ctx = s.ctx.WithBlockHeight(1).WithBlockTime(utils.ParseTime("2022-03-02T00:00:00Z"))
	s.keeper.InitGenesis(s.ctx, *genState)

	// Check pair indexes.
	pair, found := s.keeper.GetPairByDenoms(s.ctx, appID1, "denom1", "denom2")
	s.Require().True(found)
	s.Require().Equal(pair1.Id, pair.Id)

	resp1, err := s.querier.Pairs(sdk.WrapSDKContext(s.ctx), &types.QueryPairsRequest{
		Denoms: []string{"denom2", "denom1"},
		AppId:  appID1,
	})
	s.Require().NoError(err)
	s.Require().Len(resp1.Pairs, 1)
	s.Require().Equal(pair1.Id, resp1.Pairs[0].Id)

	resp2, err := s.querier.Pairs(sdk.WrapSDKContext(s.ctx), &types.QueryPairsRequest{
		Denoms: []string{"denom2", "denom3"},
		AppId:  appID1,
	})
	s.Require().NoError(err)
	s.Require().Len(resp2.Pairs, 1)
	s.Require().Equal(pair2.Id, resp2.Pairs[0].Id)

	// Check pool indexes.
	pools := s.keeper.GetPoolsByPair(s.ctx, appID1, pair2.Id)
	s.Require().Len(pools, 1)
	s.Require().Equal(pool2.Id, pools[0].Id)

	pool, found := s.keeper.GetPoolByReserveAddress(s.ctx, appID1, pool1.GetReserveAddress())
	s.Require().True(found)
	s.Require().Equal(pool1.Id, pool.Id)

	// Check deposit request indexes.
	depositReqs := s.keeper.GetDepositRequestsByDepositor(s.ctx, appID1, s.addr(4))
	s.Require().Len(depositReqs, 1)
	s.Require().Equal(depositReq1.Id, depositReqs[0].Id)

	depositReqs = s.keeper.GetDepositRequestsByDepositor(s.ctx, appID1, s.addr(5))
	s.Require().Len(depositReqs, 1)
	s.Require().Equal(depositReq2.Id, depositReqs[0].Id)

	// Check withdraw request indexes
	withdrawReqs := s.keeper.GetWithdrawRequestsByWithdrawer(s.ctx, appID1, s.addr(4))
	s.Require().Len(withdrawReqs, 1)
	s.Require().Equal(withdrawReq1.Id, withdrawReqs[0].Id)

	withdrawReqs = s.keeper.GetWithdrawRequestsByWithdrawer(s.ctx, appID1, s.addr(5))
	s.Require().Len(withdrawReqs, 1)
	s.Require().Equal(withdrawReq2.Id, withdrawReqs[0].Id)

	// Check order indexes
	orders := s.keeper.GetOrdersByOrderer(s.ctx, appID1, s.addr(6))
	s.Require().Len(orders, 1)
	s.Require().Equal(order1.Id, orders[0].Id)

	orders = s.keeper.GetOrdersByOrderer(s.ctx, appID1, s.addr(7))
	s.Require().Len(orders, 1)
	s.Require().Equal(order2.Id, orders[0].Id)
}
