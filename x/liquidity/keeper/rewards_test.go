package keeper_test

import (
	"time"

	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/liquidity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestSoftLockTokens() {
	creator := s.addr(0)

	appID1 := s.CreateNewApp("appOne")
	appID2 := s.CreateNewApp("appTwo")

	asset1 := s.CreateNewAsset("ASSET1", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSET2", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, creator, "1000000000000uasset1,1000000000000uasset2")

	pair2 := s.CreateNewLiquidityPair(appID2, creator, asset1.Denom, asset2.Denom)
	pool2 := s.CreateNewLiquidityPool(appID2, pair2.Id, creator, "1000000000000uasset1,1000000000000uasset2")

	liquidityProvider1 := s.addr(1)
	// app1 deposit
	s.Deposit(appID1, pool.Id, liquidityProvider1, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("10000000000pool1-1").IsEqual(s.getBalances(liquidityProvider1)))
	// app2 deposit
	s.Deposit(appID2, pool2.Id, liquidityProvider1, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("10000000000pool1-1,10000000000pool2-1").IsEqual(s.getBalances(liquidityProvider1)))

	liquidityProvider2 := s.addr(2)
	// app1 deposit
	s.Deposit(appID1, pool.Id, liquidityProvider2, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("9999999999pool1-1").IsEqual(s.getBalances(liquidityProvider2)))
	// app2 deposit
	s.Deposit(appID2, pool2.Id, liquidityProvider2, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("9999999999pool1-1,9999999999pool2-1").IsEqual(s.getBalances(liquidityProvider2)))

	currentTime := time.Now()
	s.ctx = s.ctx.WithBlockTime(currentTime)

	testCases := []struct {
		Name             string
		Msg              types.MsgTokensSoftLock
		ExpErr           error
		AvailableBalance sdk.Coins
		QueueLenght      uint64
	}{
		{
			Name:             "error app id invalid",
			Msg:              *types.NewMsgSoftLock(69, liquidityProvider1, pool.Id, utils.ParseCoin("699000000pool1-1")),
			ExpErr:           sdkerrors.Wrapf(types.ErrInvalidAppID, "app id %d not found", 69),
			AvailableBalance: utils.ParseCoins("10000000000pool1-1,10000000000pool2-1"),
		},
		{
			Name:             "error pool id invalid",
			Msg:              *types.NewMsgSoftLock(appID1, liquidityProvider1, 69, utils.ParseCoin("699000000pool1-1")),
			ExpErr:           sdkerrors.Wrapf(types.ErrInvalidPoolID, "no pool exists with id : %d", 69),
			AvailableBalance: utils.ParseCoins("10000000000pool1-1,10000000000pool2-1"),
		},
		{
			Name:             "error pool denom invalid",
			Msg:              *types.NewMsgSoftLock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("699000000pool1-2")),
			ExpErr:           sdkerrors.Wrapf(types.ErrWrongPoolCoinDenom, "expected pool coin denom %s, found pool1-2", pool.PoolCoinDenom),
			AvailableBalance: utils.ParseCoins("10000000000pool1-1,10000000000pool2-1"),
		},
		{
			Name:             "error insufficient pool denoms",
			Msg:              *types.NewMsgSoftLock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("100000000000pool1-1")),
			ExpErr:           sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "10000000000pool1-1 is smaller than 100000000000pool1-1"),
			AvailableBalance: utils.ParseCoins("10000000000pool1-1,10000000000pool2-1"),
		},
		{
			Name:             "success liquidity provider 1 app1",
			Msg:              *types.NewMsgSoftLock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("5252000000pool1-1")),
			ExpErr:           nil,
			AvailableBalance: utils.ParseCoins("4748000000pool1-1,10000000000pool2-1"),
			QueueLenght:      1,
		},
		{
			Name:             "success liquidity provider 2 app1",
			Msg:              *types.NewMsgSoftLock(appID1, liquidityProvider2, pool.Id, utils.ParseCoin("6934000000pool1-1")),
			ExpErr:           nil,
			AvailableBalance: utils.ParseCoins("3065999999pool1-1,9999999999pool2-1"),
			QueueLenght:      2,
		},
		{
			Name:             "success liquidity provider 1 app2",
			Msg:              *types.NewMsgSoftLock(appID2, liquidityProvider1, pool2.Id, utils.ParseCoin("123000000pool2-1")),
			ExpErr:           nil,
			AvailableBalance: utils.ParseCoins("4748000000pool1-1,9877000000pool2-1"),
			QueueLenght:      1,
		},
		{
			Name:             "success liquidity provider 2 app2",
			Msg:              *types.NewMsgSoftLock(appID2, liquidityProvider2, pool2.Id, utils.ParseCoin("546000000pool2-1")),
			ExpErr:           nil,
			AvailableBalance: utils.ParseCoins("3065999999pool1-1,9453999999pool2-1"),
			QueueLenght:      2,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			err := s.keeper.SoftLockTokens(s.ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				_, found := s.keeper.GetPoolLiquidityProvidersData(s.ctx, tc.Msg.AppId, tc.Msg.PoolId)
				s.Require().False(found)
			} else {
				s.Require().NoError(err)
				s.Require().True(tc.AvailableBalance.IsEqual(s.getBalances(sdk.MustAccAddressFromBech32(tc.Msg.Depositor))))

				lpData, found := s.keeper.GetPoolLiquidityProvidersData(s.ctx, tc.Msg.AppId, tc.Msg.PoolId)
				s.Require().True(found)

				s.Require().Equal(tc.Msg.AppId, lpData.AppId)
				s.Require().Equal(tc.Msg.PoolId, lpData.PoolId)
				s.Require().Len(lpData.QueuedLiquidityProviders, int(tc.QueueLenght))
				s.Require().Len(lpData.BondedLockIds, 0)
				s.Require().Len(lpData.LiquidityProviders, 0)

				s.Require().Equal(tc.Msg.Depositor, lpData.QueuedLiquidityProviders[tc.QueueLenght-1].Address)
				s.Require().Equal(&tc.Msg.SoftLockCoin, lpData.QueuedLiquidityProviders[tc.QueueLenght-1].SupplyProvided[0])
			}
		})
	}

	// increase time and check if deque works
	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(types.DefaultFarmingQueueDuration).Add(time.Second * 10))
	s.nextBlock()
	// app1 check
	lpData, found := s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool.Id)
	s.Require().True(found)
	s.Require().Equal(appID1, lpData.AppId)
	s.Require().Equal(pool.Id, lpData.PoolId)
	s.Require().Len(lpData.QueuedLiquidityProviders, 0)
	s.Require().Len(lpData.LiquidityProviders, 2)

	lp1, found := lpData.LiquidityProviders[liquidityProvider1.String()]
	s.Require().True(found)
	s.Require().True(lp1.Coins[0].IsEqual(utils.ParseCoin("5252000000pool1-1")))

	lp2, found := lpData.LiquidityProviders[liquidityProvider2.String()]
	s.Require().True(found)
	s.Require().True(lp2.Coins[0].IsEqual(utils.ParseCoin("6934000000pool1-1")))

	_, found = lpData.LiquidityProviders[creator.String()]
	s.Require().False(found)

	// app2 check
	lpData, found = s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID2, pool2.Id)
	s.Require().True(found)
	s.Require().Equal(appID2, lpData.AppId)
	s.Require().Equal(pool2.Id, lpData.PoolId)
	s.Require().Len(lpData.QueuedLiquidityProviders, 0)
	s.Require().Len(lpData.LiquidityProviders, 2)

	lp1, found = lpData.LiquidityProviders[liquidityProvider1.String()]
	s.Require().True(found)
	s.Require().True(lp1.Coins[0].IsEqual(utils.ParseCoin("123000000pool2-1")))

	lp2, found = lpData.LiquidityProviders[liquidityProvider2.String()]
	s.Require().True(found)
	s.Require().True(lp2.Coins[0].IsEqual(utils.ParseCoin("546000000pool2-1")))

	_, found = lpData.LiquidityProviders[creator.String()]
	s.Require().False(found)
}
