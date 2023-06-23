package liquidity_test

import (
	"strings"
	"testing"
	"time"

	_ "github.com/stretchr/testify/suite"

	"github.com/comdex-official/comdex/app"
	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/liquidity"
	"github.com/comdex-official/comdex/x/liquidity/keeper"
	"github.com/comdex-official/comdex/x/liquidity/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestInvalidMsg(t *testing.T) {
	app := app.Setup(false)

	app.LiquidityKeeper = keeper.NewKeeper(
		app.AppCodec(),
		app.GetKey(types.StoreKey),
		app.GetSubspace(types.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		&app.AssetKeeper,
		&app.MarketKeeper,
		&app.Rewardskeeper,
		&app.TokenmintKeeper,
	)
	h := liquidity.NewHandler(app.LiquidityKeeper)

	res, err := h(sdk.NewContext(nil, tmproto.Header{}, false, nil), testdata.NewTestMsg())
	require.Error(t, err)
	require.Nil(t, res)

	_, _, log := sdkerrors.ABCIInfo(err, false)
	require.True(t, strings.Contains(log, "unrecognized liquidityV1 message type:"))
}

func (s *ModuleTestSuite) TestMsgCreatePair() {
	handler := liquidity.NewHandler(s.keeper)
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")
	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)

	params, err := s.keeper.GetGenericParams(s.ctx, appID1)
	s.Require().NoError(err)
	msg := types.NewMsgCreatePair(
		appID1, addr1, asset1.Denom, asset2.Denom,
	)
	s.fundAddr(addr1, params.PairCreationFee)

	_, err = handler(s.ctx, msg)
	s.Require().NoError(err)

	pairs := s.keeper.GetAllPairs(s.ctx, appID1)
	s.Require().Len(pairs, 1)
	s.Require().Equal(uint64(1), pairs[0].Id)
	s.Require().Equal(asset1.Denom, pairs[0].BaseCoinDenom)
	s.Require().Equal(asset2.Denom, pairs[0].QuoteCoinDenom)
	s.Require().Equal("cosmos1url34vfv5a5a7esm2aapklqelh2mzuwe34vvgc94eh752t8mcjeqla0n0v", pairs[0].EscrowAddress)
	s.Require().Equal(uint64(0), pairs[0].LastOrderId)
	s.Require().Equal(uint64(1), pairs[0].CurrentBatchId)
	s.Require().Equal("cosmos19a7w3ferywxjst035636dzktx94xyh22u64pwee3fl62sennw5hsw8erx3", pairs[0].SwapFeeCollectorAddress)
	s.Require().Equal(appID1, pairs[0].AppId)
}

func (s *ModuleTestSuite) TestMsgCreatePool() {
	handler := liquidity.NewHandler(s.keeper)
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)

	params, err := s.keeper.GetGenericParams(s.ctx, appID1)
	s.Require().NoError(err)
	msg := types.NewMsgCreatePool(
		appID1, addr1, pair.Id, utils.ParseCoins("1000000000000uasset1,1000000000000uasset2"),
	)
	s.fundAddr(addr1, params.PoolCreationFee)
	s.fundAddr(addr1, msg.DepositCoins)

	_, err = handler(s.ctx, msg)
	s.Require().NoError(err)

	pools := s.keeper.GetAllPools(s.ctx, appID1)
	s.Require().Len(pools, 1)
	s.Require().Equal(uint64(1), pools[0].Id)
	s.Require().Equal(pair.Id, pools[0].PairId)
	s.Require().Equal("cosmos1s83g2v83mtmc4y4wcf3mw8204flrt0wdlp2mmjnykn2csxcrss5ql5hh0m", pools[0].ReserveAddress)
	s.Require().Equal("pool1-1", pools[0].PoolCoinDenom)
	s.Require().Equal(uint64(0), pools[0].LastDepositRequestId)
	s.Require().Equal(uint64(0), pools[0].LastWithdrawRequestId)
	s.Require().Equal(false, pools[0].Disabled)
	s.Require().Equal(appID1, pools[0].AppId)

	s.Require().True(utils.ParseCoins("1000000000000uasset1,1000000000000uasset2").IsEqual(s.getBalances(pools[0].GetReserveAddress())))

	gauges := s.app.Rewardskeeper.GetAllGauges(s.ctx)
	s.Require().Len(gauges, 1)
	s.Require().True(gauges[0].ForSwapFee)
	s.Require().False(gauges[0].GetLiquidityMetaData().IsMasterPool)
}

func (s *ModuleTestSuite) TestMsgDeposit() {
	handler := liquidity.NewHandler(s.keeper)
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000000uasset1,1000000000000uasset2")

	msg := types.NewMsgDeposit(
		appID1, addr1, pair.Id, utils.ParseCoins("1000000000uasset1,1000000000uasset2"),
	)
	s.fundAddr(addr1, msg.DepositCoins)

	_, err := handler(s.ctx, msg)
	s.Require().NoError(err)

	depositReq := s.keeper.GetAllDepositRequests(s.ctx, appID1)
	s.Require().Equal(uint64(1), depositReq[0].Id)
	s.Require().Equal(pool.Id, depositReq[0].PoolId)
	s.Require().Equal(addr1.String(), depositReq[0].Depositor)
	s.Require().Equal(msg.DepositCoins, depositReq[0].DepositCoins)
	s.Require().Equal(utils.ParseCoin("0pool1-1"), depositReq[0].MintedPoolCoin)
	s.Require().Equal(types.RequestStatusNotExecuted, depositReq[0].Status)
	s.Require().Equal(appID1, depositReq[0].AppId)
}

func (s *ModuleTestSuite) TestMsgWithdraw() {
	handler := liquidity.NewHandler(s.keeper)
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000000uasset1,1000000000000uasset2")

	msg := types.NewMsgWithdraw(
		appID1, addr1, pair.Id, utils.ParseCoin("1000pool1-1"),
	)
	_, err := handler(s.ctx, msg)
	s.Require().NoError(err)

	withdrawReq := s.keeper.GetAllWithdrawRequests(s.ctx, appID1)
	s.Require().Equal(uint64(1), withdrawReq[0].Id)
	s.Require().Equal(pool.Id, withdrawReq[0].PoolId)
	s.Require().Equal(addr1.String(), withdrawReq[0].Withdrawer)
	s.Require().Equal(msg.PoolCoin, withdrawReq[0].PoolCoin)
	s.Require().Equal(types.RequestStatusNotExecuted, withdrawReq[0].Status)
	s.Require().Equal(appID1, withdrawReq[0].AppId)
}

func (s *ModuleTestSuite) TestMsgLimitOrder() {
	handler := liquidity.NewHandler(s.keeper)
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	_ = s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000000uasset1,1000000000000uasset2")

	msg := types.NewMsgLimitOrder(
		appID1,
		addr1,
		pair.Id,
		types.OrderDirectionBuy,
		utils.ParseCoin("1003000uasset2"),
		asset1.Denom,
		sdk.NewDec(1),
		sdk.NewInt(1000000),
		time.Second*10,
	)
	s.fundAddr(addr1, sdk.NewCoins(msg.OfferCoin))
	_, err := handler(s.ctx, msg)
	s.Require().NoError(err)

	orders := s.keeper.GetAllOrders(s.ctx, appID1)
	s.Require().Len(orders, 1)
	s.Require().Equal(uint64(1), orders[0].Id)
	s.Require().Equal(pair.Id, orders[0].PairId)
	s.Require().Equal(int64(0), orders[0].MsgHeight)
	s.Require().Equal(addr1.String(), orders[0].Orderer)
	s.Require().Equal(types.OrderDirectionBuy, orders[0].Direction)
	s.Require().Equal(utils.ParseCoin("1000000uasset2"), orders[0].OfferCoin)
	s.Require().Equal(utils.ParseCoin("1000000uasset2"), orders[0].RemainingOfferCoin)
	s.Require().Equal(utils.ParseCoin("0uasset1"), orders[0].ReceivedCoin)
	s.Require().Equal(sdk.NewDec(1), orders[0].Price)
	s.Require().Equal(sdk.NewInt(1000000), orders[0].Amount)
	s.Require().Equal(sdk.NewInt(1000000), orders[0].OpenAmount)
	s.Require().Equal(uint64(1), orders[0].BatchId)
	s.Require().Equal(s.ctx.BlockTime().Add(time.Second*10), orders[0].ExpireAt)
	s.Require().Equal(types.OrderStatusNotExecuted, orders[0].Status)
	s.Require().Equal(appID1, orders[0].AppId)
}

func (s *ModuleTestSuite) TestMsgMarketOrder() {
	handler := liquidity.NewHandler(s.keeper)
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	_ = s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000000uasset1,1000000000000uasset2")

	s.LimitOrder(appID1, addr1, pair.Id, types.OrderDirectionBuy, utils.ParseDec("1"), sdk.NewInt(10000), 0)
	s.LimitOrder(appID1, addr1, pair.Id, types.OrderDirectionSell, utils.ParseDec("1"), sdk.NewInt(10000), 0)
	s.nextBlock()

	msg := types.NewMsgMarketOrder(
		appID1,
		addr1,
		pair.Id,
		types.OrderDirectionBuy,
		utils.ParseCoin("1103300uasset2"),
		asset1.Denom,
		sdk.NewInt(1000000),
		time.Second*10,
	)
	s.fundAddr(addr1, sdk.NewCoins(msg.OfferCoin))
	_, err := handler(s.ctx, msg)
	s.Require().NoError(err)

	orders := s.keeper.GetAllOrders(s.ctx, appID1)
	s.Require().Len(orders, 1)
	s.Require().Equal(uint64(3), orders[0].Id)
	s.Require().Equal(pair.Id, orders[0].PairId)
	s.Require().Equal(int64(0), orders[0].MsgHeight)
	s.Require().Equal(addr1.String(), orders[0].Orderer)
	s.Require().Equal(types.OrderDirectionBuy, orders[0].Direction)
	s.Require().Equal(utils.ParseCoin("1100000uasset2"), orders[0].OfferCoin)
	s.Require().Equal(utils.ParseCoin("1100000uasset2"), orders[0].RemainingOfferCoin)
	s.Require().Equal(utils.ParseCoin("0uasset1"), orders[0].ReceivedCoin)
	s.Require().Equal(sdk.MustNewDecFromStr("1.1"), orders[0].Price)
	s.Require().Equal(sdk.NewInt(1000000), orders[0].Amount)
	s.Require().Equal(sdk.NewInt(1000000), orders[0].OpenAmount)
	s.Require().Equal(uint64(2), orders[0].BatchId)
	s.Require().Equal(s.ctx.BlockTime().Add(time.Second*10), orders[0].ExpireAt)
	s.Require().Equal(types.OrderStatusNotExecuted, orders[0].Status)
	s.Require().Equal(appID1, orders[0].AppId)
}

func (s *ModuleTestSuite) TestMsgCancelOrder() {
	handler := liquidity.NewHandler(s.keeper)
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	_ = s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000000uasset1,1000000000000uasset2")

	order := s.LimitOrder(appID1, addr2, pair.Id, types.OrderDirectionSell, utils.ParseDec("1.1"), sdk.NewInt(1000000), time.Second*10)

	s.nextBlock()

	msg := types.NewMsgCancelOrder(appID1, addr2, pair.Id, order.Id)
	_, err := handler(s.ctx, msg)
	s.Require().NoError(err)

	order, found := s.keeper.GetOrder(s.ctx, appID1, pair.Id, order.Id)
	s.Require().True(found)
	s.Require().Equal(types.OrderStatusCanceled, order.Status)

	s.Require().True(utils.ParseCoins("1003000uasset1").IsEqual(s.getBalances(addr2)))

	s.nextBlock()
	_, found = s.keeper.GetOrder(s.ctx, appID1, pair.Id, order.Id)
	s.Require().False(found)
}

func (s *ModuleTestSuite) TestMsgCancelAllOrders() {
	handler := liquidity.NewHandler(s.keeper)
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	_ = s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000000uasset1,1000000000000uasset2")

	order := s.LimitOrder(appID1, addr2, pair.Id, types.OrderDirectionSell, utils.ParseDec("1.1"), sdk.NewInt(1000000), time.Second*10)

	s.nextBlock()

	msg := types.NewMsgCancelAllOrders(appID1, addr2, []uint64{pair.Id})
	_, err := handler(s.ctx, msg)
	s.Require().NoError(err)

	order, found := s.keeper.GetOrder(s.ctx, appID1, pair.Id, order.Id)
	s.Require().True(found)
	s.Require().Equal(types.OrderStatusCanceled, order.Status)

	s.Require().True(utils.ParseCoins("1003000uasset1").IsEqual(s.getBalances(addr2)))

	s.nextBlock()
	_, found = s.keeper.GetOrder(s.ctx, appID1, pair.Id, order.Id)
	s.Require().False(found)
}

func (s *ModuleTestSuite) TestMsgFarm() {
	handler := liquidity.NewHandler(s.keeper)
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000000uasset1,1000000000000uasset2")

	liquidityProvider1 := s.addr(2)
	s.Deposit(appID1, pool.Id, liquidityProvider1, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("10000000000pool1-1").IsEqual(s.getBalances(liquidityProvider1)))

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime())
	msg := types.NewMsgFarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("5000000000pool1-1"))
	_, err := handler(s.ctx, msg)
	s.Require().NoError(err)

	queuedFarmer, found := s.keeper.GetQueuedFarmer(s.ctx, appID1, pool.Id, liquidityProvider1)
	s.Require().True(found)
	s.Require().Equal(queuedFarmer.QueudCoins[0].FarmedPoolCoin.Denom, "pool1-1")
	s.Require().Equal(queuedFarmer.QueudCoins[0].FarmedPoolCoin.Amount, sdk.NewInt(5000000000))
}

func (s *ModuleTestSuite) TestMsgUnfarm() {
	handler := liquidity.NewHandler(s.keeper)
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000000uasset1,1000000000000uasset2")

	liquidityProvider1 := s.addr(2)
	s.Deposit(appID1, pool.Id, liquidityProvider1, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("10000000000pool1-1").IsEqual(s.getBalances(liquidityProvider1)))

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime())
	msg := types.NewMsgFarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("5000000000pool1-1"))
	_, err := handler(s.ctx, msg)
	s.Require().NoError(err)

	queuedFarmer, found := s.keeper.GetQueuedFarmer(s.ctx, appID1, pool.Id, liquidityProvider1)
	s.Require().True(found)
	s.Require().Equal(queuedFarmer.QueudCoins[0].FarmedPoolCoin.Denom, "pool1-1")
	s.Require().Equal(queuedFarmer.QueudCoins[0].FarmedPoolCoin.Amount, sdk.NewInt(5000000000))

	msgUnlock := types.NewMsgUnfarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("5000000000pool1-1"))
	_, err = handler(s.ctx, msgUnlock)
	s.Require().NoError(err)

	queuedFarmer, found = s.keeper.GetQueuedFarmer(s.ctx, appID1, pool.Id, liquidityProvider1)
	s.Require().True(found)
	s.Require().Len(queuedFarmer.QueudCoins, 0)
	s.Require().True(utils.ParseCoins("10000000000pool1-1").IsEqual(s.getBalances(liquidityProvider1)))
}
