package keeper_test

import (
	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/liquidity/keeper"
	"github.com/comdex-official/comdex/x/liquidity/types"
	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestDepositCoinsEscrowInvariant() {
	creator := s.addr(0)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSET1", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSET2", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, creator, "1000000000000uasset1,1000000000000uasset2")

	req := s.Deposit(appID1, pool.Id, s.addr(1), "1000000uasset1,1000000uasset2")
	_, broken := keeper.DepositCoinsEscrowInvariant(s.keeper)(s.ctx)
	s.Require().False(broken)

	oldReq := req
	req.DepositCoins = utils.ParseCoins("2000000uasset1,2000000uasset2")
	s.keeper.SetDepositRequest(s.ctx, req)
	_, broken = keeper.DepositCoinsEscrowInvariant(s.keeper)(s.ctx)
	s.Require().True(broken)

	req = oldReq
	s.keeper.SetDepositRequest(s.ctx, req)
	s.nextBlock()
	_, broken = keeper.DepositCoinsEscrowInvariant(s.keeper)(s.ctx)
	s.Require().False(broken)
}

func (s *KeeperTestSuite) TestPoolCoinEscrowInvariant() {
	creator := s.addr(0)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSET1", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSET2", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, creator, "1000000000000uasset1,1000000000000uasset2")

	s.Deposit(appID1, pool.Id, s.addr(1), "1000000uasset1,1000000uasset2")
	s.nextBlock()

	req := s.Withdraw(appID1, pool.Id, s.addr(1), utils.ParseCoin("1000000pool1-1"))
	_, broken := keeper.PoolCoinEscrowInvariant(s.keeper)(s.ctx)
	s.Require().False(broken)

	oldReq := req
	req.PoolCoin = utils.ParseCoin("2000000pool1")
	s.keeper.SetWithdrawRequest(s.ctx, req)
	_, broken = keeper.PoolCoinEscrowInvariant(s.keeper)(s.ctx)
	s.Require().True(broken)

	req = oldReq
	s.keeper.SetWithdrawRequest(s.ctx, req)
	s.nextBlock()
	_, broken = keeper.PoolCoinEscrowInvariant(s.keeper)(s.ctx)
	s.Require().False(broken)
}

func (s *KeeperTestSuite) TestRemainingOfferCoinEscrowInvariant() {
	creator := s.addr(0)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSET1", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSET2", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset2.Denom)

	order := s.LimitOrder(appID1, s.addr(1), pair.Id, types.OrderDirectionBuy, utils.ParseDec("1.0"), newInt(1000000), 0)
	_, broken := keeper.RemainingOfferCoinEscrowInvariant(s.keeper)(s.ctx)
	s.Require().False(broken)

	oldOrder := order
	order.RemainingOfferCoin = utils.ParseCoin("2000000denom1")
	s.keeper.SetOrder(s.ctx, appID1, order)
	_, broken = keeper.RemainingOfferCoinEscrowInvariant(s.keeper)(s.ctx)
	s.Require().True(broken)

	order = oldOrder
	s.keeper.SetOrder(s.ctx, appID1, order)
	s.nextBlock()
	_, broken = keeper.RemainingOfferCoinEscrowInvariant(s.keeper)(s.ctx)
	s.Require().False(broken)
}

func (s *KeeperTestSuite) TestPoolStatusInvariant() {
	creator := s.addr(0)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSET1", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSET2", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, creator, "1000000000000uasset1,1000000000000uasset2")

	_, broken := keeper.PoolStatusInvariant(s.keeper)(s.ctx)
	s.Require().False(broken)

	s.Withdraw(appID1, pool.Id, creator, s.getBalance(s.addr(0), pool.PoolCoinDenom))
	s.nextBlock()

	_, broken = keeper.PoolStatusInvariant(s.keeper)(s.ctx)
	s.Require().False(broken)

	pool, _ = s.keeper.GetPool(s.ctx, appID1, pool.Id)
	pool.Disabled = false
	s.keeper.SetPool(s.ctx, pool)
	_, broken = keeper.PoolStatusInvariant(s.keeper)(s.ctx)
	s.Require().True(broken)
}
