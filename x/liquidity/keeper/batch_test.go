package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/liquidity"
	"github.com/comdex-official/comdex/x/liquidity/types"

	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestOrderExpiration() {
	creator := s.addr(0)

	appID1 := s.CreateNewApp("appOne")

	asset1 := s.CreateNewAsset("ASSET1", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSET2", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset2.Denom)

	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2022-03-01T12:00:00Z"))
	order := s.LimitOrder(appID1, s.addr(1), pair.Id, types.OrderDirectionSell, utils.ParseDec("1.0"), sdk.NewInt(10000), 10*time.Second)
	liquidity.EndBlocker(s.ctx, s.keeper)

	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2022-03-01T12:00:06Z"))
	liquidity.BeginBlocker(s.ctx, s.keeper)
	order, found := s.keeper.GetOrder(s.ctx, appID1, order.PairId, order.Id)
	s.Require().True(found) // The order is not yet deleted.
	// A buy order comes in.
	s.LimitOrder(appID1, s.addr(2), pair.Id, types.OrderDirectionBuy, utils.ParseDec("1.0"), sdk.NewInt(5000), 0)
	liquidity.EndBlocker(s.ctx, s.keeper)

	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2022-03-01T12:00:12Z"))
	liquidity.BeginBlocker(s.ctx, s.keeper)
	order, found = s.keeper.GetOrder(s.ctx, appID1, order.PairId, order.Id)
	s.Require().True(found)
	s.Require().Equal(types.OrderStatusPartiallyMatched, order.Status)
	// Another buy order comes in, but this time the first order has been expired,
	// so there is no match.
	s.LimitOrder(appID1, s.addr(3), pair.Id, types.OrderDirectionBuy, utils.ParseDec("1.0"), sdk.NewInt(5000), 0)
	liquidity.EndBlocker(s.ctx, s.keeper)
	order, _ = s.keeper.GetOrder(s.ctx, appID1, order.PairId, order.Id)
	s.Require().Equal(types.OrderStatusExpired, order.Status)
	s.Require().True(sdk.NewInt(5000).Equal(order.OpenAmount))

	liquidity.BeginBlocker(s.ctx, s.keeper)
	_, found = s.keeper.GetOrder(s.ctx, appID1, order.PairId, order.Id)
	s.Require().False(found) // The order is gone.
}
