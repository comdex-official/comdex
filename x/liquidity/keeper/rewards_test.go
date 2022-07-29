package keeper_test

import (
	"time"

	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/liquidity/types"
	rewardtypes "github.com/comdex-official/comdex/x/rewards/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestSoftLockTokens() {
	creator := s.addr(0)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")

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

	currentTime := s.ctx.BlockTime()
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
			Name:             "success liquidity provider 1 app1 re-add",
			Msg:              *types.NewMsgSoftLock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("23000000pool1-1")),
			ExpErr:           nil,
			AvailableBalance: utils.ParseCoins("4725000000pool1-1,10000000000pool2-1"),
			QueueLenght:      3,
		},
		{
			Name:             "success liquidity provider 1 app2",
			Msg:              *types.NewMsgSoftLock(appID2, liquidityProvider1, pool2.Id, utils.ParseCoin("123000000pool2-1")),
			ExpErr:           nil,
			AvailableBalance: utils.ParseCoins("4725000000pool1-1,9877000000pool2-1"),
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
	s.Require().True(lp1.Coins[0].IsEqual(utils.ParseCoin("5275000000pool1-1")))

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

func (s *KeeperTestSuite) TestSoftUnlockTokens() {
	creator := s.addr(0)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSET1", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSET2", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, creator, "1000000000000uasset1,1000000000000uasset2")

	pair2 := s.CreateNewLiquidityPair(appID1, creator, asset2.Denom, asset1.Denom)
	pool2 := s.CreateNewLiquidityPool(appID1, pair2.Id, creator, "1000000000000uasset1,1000000000000uasset2")

	liquidityProvider1 := s.addr(1)
	s.Deposit(appID1, pool.Id, liquidityProvider1, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("10000000000pool1-1").IsEqual(s.getBalances(liquidityProvider1)))

	currentTime := s.ctx.BlockTime()
	s.ctx = s.ctx.WithBlockTime(currentTime)

	msg := types.NewMsgSoftLock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("10000000000pool1-1"))
	err := s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)
	lpData, found := s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool.Id)
	s.Require().True(found)
	s.Require().Equal(appID1, lpData.AppId)
	s.Require().Equal(pool.Id, lpData.PoolId)
	s.Require().Len(lpData.QueuedLiquidityProviders, 1)
	s.Require().Len(lpData.LiquidityProviders, 0)

	testCases := []struct {
		Name             string
		Msg              types.MsgTokensSoftUnlock
		ExpErr           error
		AvailableBalance sdk.Coins
		QueueLenght      uint64
	}{
		{
			Name:             "error app id invalid",
			Msg:              *types.NewMsgSoftUnlock(69, liquidityProvider1, pool.Id, utils.ParseCoin("699000000pool1-1")),
			ExpErr:           sdkerrors.Wrapf(types.ErrInvalidAppID, "app id %d not found", 69),
			AvailableBalance: sdk.Coins{},
		},
		{
			Name:             "error pool id invalid",
			Msg:              *types.NewMsgSoftUnlock(appID1, liquidityProvider1, 69, utils.ParseCoin("699000000pool1-1")),
			ExpErr:           sdkerrors.Wrapf(types.ErrInvalidPoolID, "no pool exists with id : %d", 69),
			AvailableBalance: sdk.Coins{},
		},
		{
			Name:             "error pool denom invalid",
			Msg:              *types.NewMsgSoftUnlock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("699000000pool1-2")),
			ExpErr:           sdkerrors.Wrapf(types.ErrWrongPoolCoinDenom, "expected pool coin denom %s, found pool1-2", pool.PoolCoinDenom),
			AvailableBalance: sdk.Coins{},
		},
		{
			Name:             "error no soft locks present for pool",
			Msg:              *types.NewMsgSoftUnlock(appID1, liquidityProvider1, pool2.Id, utils.ParseCoin("699000000pool1-2")),
			ExpErr:           sdkerrors.Wrapf(types.ErrNoSoftLockPresent, "no soft locks present for given pool id %d", pool2.Id),
			AvailableBalance: sdk.Coins{},
		},
		{
			Name:             "error insufficient farmed amounts",
			Msg:              *types.NewMsgSoftUnlock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("100000000000pool1-1")),
			ExpErr:           sdkerrors.Wrapf(types.ErrInvalidUnlockAmount, "available soft locked amount 10000000000pool1-1 smaller than requested amount 100000000000pool1-1"),
			AvailableBalance: sdk.Coins{},
		},
		{
			Name:             "success partial unlock",
			Msg:              *types.NewMsgSoftUnlock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("5000000000pool1-1")),
			ExpErr:           nil,
			AvailableBalance: utils.ParseCoins("5000000000pool1-1"),
		},
		{
			Name:             "success full unlock",
			Msg:              *types.NewMsgSoftUnlock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("5000000000pool1-1")),
			ExpErr:           nil,
			AvailableBalance: utils.ParseCoins("10000000000pool1-1"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			err := s.keeper.SoftUnlockTokens(s.ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().True(tc.AvailableBalance.IsEqual(s.getBalances(sdk.MustAccAddressFromBech32(tc.Msg.Depositor))))
			}
		})
	}
}

func (s *KeeperTestSuite) TestSoftUnlockTokensTwo() {
	currentTime := s.ctx.BlockTime()
	s.ctx = s.ctx.WithBlockTime(currentTime)

	creator := s.addr(0)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSET1", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSET2", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, creator, "1000000000000uasset1,1000000000000uasset2")

	liquidityProvider1 := s.addr(1)
	liquidityProvider2 := s.addr(2)

	s.Deposit(appID1, pool.Id, liquidityProvider1, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("10000000000pool1-1").IsEqual(s.getBalances(liquidityProvider1)))

	s.Deposit(appID1, pool.Id, liquidityProvider2, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("9999999999pool1-1").IsEqual(s.getBalances(liquidityProvider2)))

	// farm 1, queue size 1
	// SortedByTimeFarmQueue -> [10000000pool1-1]
	msg := types.NewMsgSoftLock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("10000000pool1-1"))
	err := s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(time.Hour * 1))
	// farm 2, queue size 2
	// SortedByTimeFarmQueue -> [20000000pool1-1, 10000000pool1-1]
	msg = types.NewMsgSoftLock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("20000000pool1-1"))
	err = s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(time.Hour * 1))
	// farm 3, queue size 3
	// SortedByTimeFarmQueue -> [30000000pool1-1, 20000000pool1-1, 10000000pool1-1]
	msg = types.NewMsgSoftLock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("30000000pool1-1"))
	err = s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(time.Hour * 1))
	// farm 4, queue size 4
	// SortedByTimeFarmQueue -> [40000000pool1-1, 30000000pool1-1, 20000000pool1-1, 10000000pool1-1]
	msg = types.NewMsgSoftLock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("40000000pool1-1"))
	err = s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(time.Hour * 1))
	// farm 5, queue size 5
	// SortedByTimeFarmQueue -> [50000000pool1-1, 40000000pool1-1, 30000000pool1-1, 20000000pool1-1, 10000000pool1-1]
	msg = types.NewMsgSoftLock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("50000000pool1-1"))
	err = s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)

	lpData, found := s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool.Id)
	s.Require().True(found)
	s.Require().Equal(appID1, lpData.AppId)
	s.Require().Equal(pool.Id, lpData.PoolId)
	s.Require().Len(lpData.QueuedLiquidityProviders, 5)
	s.Require().Len(lpData.LiquidityProviders, 0)

	// lp1 trying to unfarm/softUnlock more than farmed/softLocked
	msgUnlock := types.NewMsgSoftUnlock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("160000000pool1-1"))
	err = s.keeper.SoftUnlockTokens(s.ctx, msgUnlock)
	s.Require().Error(err)
	s.Require().EqualError(err, sdkerrors.Wrapf(types.ErrInvalidUnlockAmount, "available soft locked amount 150000000pool1-1 smaller than requested amount 160000000pool1-1").Error())

	// unfarming small portions, below unlock removes token from most recently added queue
	// unlock is done from a single latest object in a queue since this object itself can satisfy the unlock requirement,
	// Before - SortedByTimeFarmQueue -> [50000000pool1-1, 40000000pool1-1, 30000000pool1-1, 20000000pool1-1, 10000000pool1-1]
	// After  - SortedByTimeFarmQueue -> [45000000pool1-1, 40000000pool1-1, 30000000pool1-1, 20000000pool1-1, 10000000pool1-1]
	msgUnlock = types.NewMsgSoftUnlock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("5000000pool1-1"))
	err = s.keeper.SoftUnlockTokens(s.ctx, msgUnlock)
	s.Require().NoError(err)
	lpData, found = s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool.Id)
	s.Require().True(found)
	s.Require().Len(lpData.QueuedLiquidityProviders, 5)
	s.Require().Equal(utils.ParseCoin("45000000pool1-1").Denom, lpData.QueuedLiquidityProviders[0].SupplyProvided[0].Denom)
	s.Require().Equal(utils.ParseCoin("45000000pool1-1").Amount, lpData.QueuedLiquidityProviders[0].SupplyProvided[0].Amount)

	// unfarming small portions, below unlock removes token from most recently added queue
	// unlock is done from a single latest object in a queue since this object itself can satisfy the unlock requirement,
	// Before  - SortedByTimeFarmQueue -> [45000000pool1-1, 40000000pool1-1, 30000000pool1-1, 20000000pool1-1, 10000000pool1-1]
	// After   - SortedByTimeFarmQueue -> [34000000pool1-1, 40000000pool1-1, 30000000pool1-1, 20000000pool1-1, 10000000pool1-1]
	msgUnlock = types.NewMsgSoftUnlock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("11000000pool1-1"))
	err = s.keeper.SoftUnlockTokens(s.ctx, msgUnlock)
	s.Require().NoError(err)
	lpData, found = s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool.Id)
	s.Require().True(found)
	s.Require().Len(lpData.QueuedLiquidityProviders, 5)
	s.Require().Equal(utils.ParseCoin("34000000pool1-1").Denom, lpData.QueuedLiquidityProviders[0].SupplyProvided[0].Denom)
	s.Require().Equal(utils.ParseCoin("34000000pool1-1").Amount, lpData.QueuedLiquidityProviders[0].SupplyProvided[0].Amount)

	// below case will delete the most recent object from queue since it satisfies the required unlock condition
	// here the unlock is being satisfied from the two queue objects, most recent one gets deleted after it fullfills all
	// of its token for unlocking, and the remaining unlock tokens are fullfilled from 2nd most recent queue object
	// Before - SortedByTimeFarmQueue -> [34000000pool1-1, 40000000pool1-1, 30000000pool1-1, 20000000pool1-1, 10000000pool1-1]
	// After   - SortedByTimeFarmQueue -> [36000000pool1-1, 30000000pool1-1, 20000000pool1-1, 10000000pool1-1]
	msgUnlock = types.NewMsgSoftUnlock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("38000000pool1-1"))
	err = s.keeper.SoftUnlockTokens(s.ctx, msgUnlock)
	s.Require().NoError(err)
	lpData, found = s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool.Id)
	s.Require().True(found)
	s.Require().Len(lpData.QueuedLiquidityProviders, 4)
	s.Require().Equal(utils.ParseCoin("36000000pool1-1").Denom, lpData.QueuedLiquidityProviders[0].SupplyProvided[0].Denom)
	s.Require().Equal(utils.ParseCoin("36000000pool1-1").Amount, lpData.QueuedLiquidityProviders[0].SupplyProvided[0].Amount)

	// similarly below cases are followed as above
	// Before   - SortedByTimeFarmQueue -> [36000000pool1-1, 30000000pool1-1, 20000000pool1-1, 10000000pool1-1]
	// After    - SortedByTimeFarmQueue -> [30000000pool1-1, 20000000pool1-1, 10000000pool1-1]
	msgUnlock = types.NewMsgSoftUnlock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("36000000pool1-1"))
	err = s.keeper.SoftUnlockTokens(s.ctx, msgUnlock)
	s.Require().NoError(err)
	lpData, found = s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool.Id)
	s.Require().True(found)
	s.Require().Len(lpData.QueuedLiquidityProviders, 3)
	s.Require().Equal(utils.ParseCoin("30000000pool1-1").Denom, lpData.QueuedLiquidityProviders[0].SupplyProvided[0].Denom)
	s.Require().Equal(utils.ParseCoin("30000000pool1-1").Amount, lpData.QueuedLiquidityProviders[0].SupplyProvided[0].Amount)

	// Before    - SortedByTimeFarmQueue -> [30000000pool1-1, 20000000pool1-1, 10000000pool1-1]
	// After     - SortedByTimeFarmQueue -> [10000000pool1-1]
	msgUnlock = types.NewMsgSoftUnlock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("50000000pool1-1"))
	err = s.keeper.SoftUnlockTokens(s.ctx, msgUnlock)
	s.Require().NoError(err)
	lpData, found = s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool.Id)
	s.Require().True(found)
	s.Require().Len(lpData.QueuedLiquidityProviders, 1)
	s.Require().Equal(utils.ParseCoin("10000000pool1-1").Denom, lpData.QueuedLiquidityProviders[0].SupplyProvided[0].Denom)
	s.Require().Equal(utils.ParseCoin("10000000pool1-1").Amount, lpData.QueuedLiquidityProviders[0].SupplyProvided[0].Amount)

	// lp1 trying to unfarm/softUnlock more than farmed/softLocked
	msgUnlock = types.NewMsgSoftUnlock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("11000000pool1-1"))
	err = s.keeper.SoftUnlockTokens(s.ctx, msgUnlock)
	s.Require().Error(err)
	s.Require().EqualError(err, sdkerrors.Wrapf(types.ErrInvalidUnlockAmount, "available soft locked amount 10000000pool1-1 smaller than requested amount 11000000pool1-1").Error())

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(time.Hour * 1))
	// SortedByTimeFarmQueue -> [69000000pool1-1, 10000000pool1-1]
	msg = types.NewMsgSoftLock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("69000000pool1-1"))
	err = s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)

	// marking oldest farmed object as valid and dequing it, assuming queue duration is satisfied
	s.ctx = s.ctx.WithBlockTime(currentTime.Add(types.DefaultFarmingQueueDuration).Add(time.Second * 10))
	s.nextBlock()
	lpData, found = s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool.Id)
	s.Require().True(found)
	s.Require().Len(lpData.QueuedLiquidityProviders, 1)
	s.Require().Len(lpData.LiquidityProviders, 1)

	// now the data is something like this
	// SortedByTimeFarmQueue -> [69000000pool1-1]
	// ActiveFarmedTokens -> [10000000pool1-1]

	// unlocking more tokens, some tokens will be unlocked from queue and some from active farmed tokens
	// Before
	// SortedByTimeFarmQueue -> [69000000pool1-1]
	// ActiveFarmedTokens -> [10000000pool1-1]
	// After
	// SortedByTimeFarmQueue -> []
	// ActiveFarmedTokens -> [9000000pool1-1]
	msgUnlock = types.NewMsgSoftUnlock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("70000000pool1-1"))
	err = s.keeper.SoftUnlockTokens(s.ctx, msgUnlock)
	s.Require().NoError(err)
	lpData, found = s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool.Id)
	s.Require().True(found)
	s.Require().Len(lpData.QueuedLiquidityProviders, 0)
	s.Require().Len(lpData.LiquidityProviders, 1)
	s.Require().Equal(utils.ParseCoin("9000000pool1-1").Denom, lpData.LiquidityProviders[liquidityProvider1.String()].Coins[0].Denom)
	s.Require().Equal(utils.ParseCoin("9000000pool1-1").Amount, lpData.LiquidityProviders[liquidityProvider1.String()].Coins[0].Amount)

	// unlocking all farmed tokens
	// Before
	// SortedByTimeFarmQueue -> []
	// ActiveFarmedTokens -> [9000000pool1-1]
	// After
	// SortedByTimeFarmQueue -> []
	// ActiveFarmedTokens -> []
	msgUnlock = types.NewMsgSoftUnlock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("9000000pool1-1"))
	err = s.keeper.SoftUnlockTokens(s.ctx, msgUnlock)
	s.Require().NoError(err)
	lpData, found = s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool.Id)
	s.Require().True(found)
	s.Require().Len(lpData.QueuedLiquidityProviders, 0)
	s.Require().Len(lpData.LiquidityProviders, 1)
	s.Require().Equal(utils.ParseCoin("0pool1-1").Denom, lpData.LiquidityProviders[liquidityProvider1.String()].Coins[0].Denom)
	s.Require().Equal(utils.ParseCoin("0pool1-1").Amount, lpData.LiquidityProviders[liquidityProvider1.String()].Coins[0].Amount)

	s.Require().True(utils.ParseCoins("10000000000pool1-1").IsEqual(s.getBalances(liquidityProvider1)))

	s.ctx = s.ctx.WithBlockTime(currentTime)
	// SortedByTimeFarmQueue -> [11000000pool1-1]
	// ActiveFarmedTokens -> []
	msg = types.NewMsgSoftLock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("11000000pool1-1"))
	err = s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(time.Hour * 1))
	// SortedByTimeFarmQueue -> [12000000pool1-1, 11000000pool1-1]
	// ActiveFarmedTokens -> []
	msg = types.NewMsgSoftLock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("12000000pool1-1"))
	err = s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(time.Hour * 1))
	// SortedByTimeFarmQueue -> [13000000pool1-1, 12000000pool1-1, 11000000pool1-1]
	// ActiveFarmedTokens -> []
	msg = types.NewMsgSoftLock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("13000000pool1-1"))
	err = s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)

	// marking oldest farmed object as valid and dequing it, assuming queue duration is satisfied
	s.ctx = s.ctx.WithBlockTime(currentTime.Add(types.DefaultFarmingQueueDuration).Add(time.Second * 10))
	s.nextBlock()
	lpData, found = s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool.Id)
	s.Require().True(found)
	s.Require().Len(lpData.QueuedLiquidityProviders, 2)
	s.Require().Len(lpData.LiquidityProviders, 1)

	// now the data is something like this
	// SortedByTimeFarmQueue -> [13000000pool1-1, 12000000pool1-1]
	// ActiveFarmedTokens -> [11000000pool1-1]

	// unlocking all queue tokes and some active tokens
	// Before
	// SortedByTimeFarmQueue -> [13000000pool1-1, 12000000pool1-1]
	// ActiveFarmedTokens -> [11000000pool1-1]
	// After
	// SortedByTimeFarmQueue -> []
	// ActiveFarmedTokens -> [10000000pool1-1]
	msgUnlock = types.NewMsgSoftUnlock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("26000000pool1-1"))
	err = s.keeper.SoftUnlockTokens(s.ctx, msgUnlock)
	s.Require().NoError(err)
	lpData, found = s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool.Id)
	s.Require().True(found)
	s.Require().Len(lpData.QueuedLiquidityProviders, 0)
	s.Require().Len(lpData.LiquidityProviders, 1)
	s.Require().Equal(utils.ParseCoin("10000000pool1-1").Denom, lpData.LiquidityProviders[liquidityProvider1.String()].Coins[0].Denom)
	s.Require().Equal(utils.ParseCoin("10000000pool1-1").Amount, lpData.LiquidityProviders[liquidityProvider1.String()].Coins[0].Amount)

	// SortedByTimeFarmQueue -> [ (l2) 7000000pool1-1]
	// ActiveFarmedTokens -> [ (l1) 10000000pool1-1]
	msg = types.NewMsgSoftLock(appID1, liquidityProvider2, pool.Id, utils.ParseCoin("7000000pool1-1"))
	err = s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)

	lpData, found = s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool.Id)
	s.Require().True(found)
	s.Require().Len(lpData.QueuedLiquidityProviders, 1)
	s.Require().Len(lpData.LiquidityProviders, 1)

	// SortedByTimeFarmQueue -> [(l1) 9000000pool1-1, (l2) 7000000pool1-1]
	// ActiveFarmedTokens -> [ (l1) 10000000pool1-1]
	msg = types.NewMsgSoftLock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("9000000pool1-1"))
	err = s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)

	lpData, found = s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool.Id)
	s.Require().True(found)
	s.Require().Len(lpData.QueuedLiquidityProviders, 2)
	s.Require().Len(lpData.LiquidityProviders, 1)

	// Before
	// SortedByTimeFarmQueue -> [(l2) 7000000pool1-1, (l1) 9000000pool1-1]
	// ActiveFarmedTokens -> [ (l1) 10000000pool1-1]
	// After
	// SortedByTimeFarmQueue -> [(l2) 7000000pool1-1, (l1) 6000000pool1-1]
	// ActiveFarmedTokens -> [ (l1) 10000000pool1-1]
	msgUnlock = types.NewMsgSoftUnlock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("3000000pool1-1"))
	err = s.keeper.SoftUnlockTokens(s.ctx, msgUnlock)
	s.Require().NoError(err)
	lpData, found = s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool.Id)
	s.Require().True(found)
	s.Require().Len(lpData.QueuedLiquidityProviders, 2)
	s.Require().Len(lpData.LiquidityProviders, 1)
	s.Require().Equal(utils.ParseCoin("7000000pool1-1").Denom, lpData.QueuedLiquidityProviders[0].SupplyProvided[0].Denom)
	s.Require().Equal(utils.ParseCoin("7000000pool1-1").Amount, lpData.QueuedLiquidityProviders[0].SupplyProvided[0].Amount)
	s.Require().Equal(utils.ParseCoin("6000000pool1-1").Denom, lpData.QueuedLiquidityProviders[1].SupplyProvided[0].Denom)
	s.Require().Equal(utils.ParseCoin("6000000pool1-1").Amount, lpData.QueuedLiquidityProviders[1].SupplyProvided[0].Amount)
	s.Require().Equal(utils.ParseCoin("10000000pool1-1").Denom, lpData.LiquidityProviders[liquidityProvider1.String()].Coins[0].Denom)
	s.Require().Equal(utils.ParseCoin("10000000pool1-1").Amount, lpData.LiquidityProviders[liquidityProvider1.String()].Coins[0].Amount)

	// Before
	// SortedByTimeFarmQueue -> [(l2) 7000000pool1-1, (l1) 6000000pool1-1]
	// ActiveFarmedTokens -> [ (l1) 10000000pool1-1]
	// After
	// SortedByTimeFarmQueue -> [(l2) 7000000pool1-1]
	// ActiveFarmedTokens -> [ (l1) 8000000pool1-1]
	msgUnlock = types.NewMsgSoftUnlock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("8000000pool1-1"))
	err = s.keeper.SoftUnlockTokens(s.ctx, msgUnlock)
	s.Require().NoError(err)
	lpData, found = s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool.Id)
	s.Require().True(found)
	s.Require().Len(lpData.QueuedLiquidityProviders, 1)
	s.Require().Len(lpData.LiquidityProviders, 1)
	s.Require().Equal(utils.ParseCoin("7000000pool1-1").Denom, lpData.QueuedLiquidityProviders[0].SupplyProvided[0].Denom)
	s.Require().Equal(utils.ParseCoin("7000000pool1-1").Amount, lpData.QueuedLiquidityProviders[0].SupplyProvided[0].Amount)
	s.Require().Equal(utils.ParseCoin("8000000pool1-1").Denom, lpData.LiquidityProviders[liquidityProvider1.String()].Coins[0].Denom)
	s.Require().Equal(utils.ParseCoin("8000000pool1-1").Amount, lpData.LiquidityProviders[liquidityProvider1.String()].Coins[0].Amount)

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(types.DefaultFarmingQueueDuration).Add(time.Minute * 1))
	s.nextBlock()
	// Now
	// SortedByTimeFarmQueue -> []
	// ActiveFarmedTokens -> [ (l1) 8000000pool1-1, (l2) 7000000pool1-1]
	lpData, found = s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool.Id)
	s.Require().True(found)
	s.Require().Len(lpData.QueuedLiquidityProviders, 0)
	s.Require().Len(lpData.LiquidityProviders, 2)
	s.Require().Equal(utils.ParseCoin("8000000pool1-1").Denom, lpData.LiquidityProviders[liquidityProvider1.String()].Coins[0].Denom)
	s.Require().Equal(utils.ParseCoin("8000000pool1-1").Amount, lpData.LiquidityProviders[liquidityProvider1.String()].Coins[0].Amount)

	s.Require().Equal(utils.ParseCoin("7000000pool1-1").Denom, lpData.LiquidityProviders[liquidityProvider2.String()].Coins[0].Denom)
	s.Require().Equal(utils.ParseCoin("7000000pool1-1").Amount, lpData.LiquidityProviders[liquidityProvider2.String()].Coins[0].Amount)

	// total unlock - lp1
	msgUnlock = types.NewMsgSoftUnlock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("8000000pool1-1"))
	err = s.keeper.SoftUnlockTokens(s.ctx, msgUnlock)
	s.Require().NoError(err)

	// total unlock - lp2
	msgUnlock = types.NewMsgSoftUnlock(appID1, liquidityProvider2, pool.Id, utils.ParseCoin("7000000pool1-1"))
	err = s.keeper.SoftUnlockTokens(s.ctx, msgUnlock)
	s.Require().NoError(err)

	// SortedByTimeFarmQueue -> []
	// ActiveFarmedTokens -> [ (l1) 0pool1-1, (l2) 0pool1-1]
	lpData, found = s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool.Id)
	s.Require().True(found)
	s.Require().Len(lpData.QueuedLiquidityProviders, 0)
	s.Require().Len(lpData.LiquidityProviders, 2)
	s.Require().Equal(utils.ParseCoin("0pool1-1").Denom, lpData.LiquidityProviders[liquidityProvider1.String()].Coins[0].Denom)
	s.Require().Equal(utils.ParseCoin("0pool1-1").Amount, lpData.LiquidityProviders[liquidityProvider1.String()].Coins[0].Amount)

	s.Require().Equal(utils.ParseCoin("0pool1-1").Denom, lpData.LiquidityProviders[liquidityProvider2.String()].Coins[0].Denom)
	s.Require().Equal(utils.ParseCoin("0pool1-1").Amount, lpData.LiquidityProviders[liquidityProvider2.String()].Coins[0].Amount)

	s.Require().True(utils.ParseCoins("10000000000pool1-1").IsEqual(s.getBalances(liquidityProvider1)))
	s.Require().True(utils.ParseCoins("9999999999pool1-1").IsEqual(s.getBalances(liquidityProvider2)))
}

// liquidity provided in incrementel order
func (s *KeeperTestSuite) TestGetFarmingRewardsDataLinearLPs() {
	currentTime := s.ctx.BlockTime()
	s.ctx = s.ctx.WithBlockTime(currentTime)

	creator := s.addr(0)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSET1", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSET2", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, creator, "1000000000000uasset1,1000000000000uasset2")

	liquidityProvider1 := s.addr(1)
	s.Deposit(appID1, pool.Id, liquidityProvider1, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("10000000000pool1-1").IsEqual(s.getBalances(liquidityProvider1)))
	msg := types.NewMsgSoftLock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("10000000000pool1-1"))
	err := s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)

	liquidityProvider2 := s.addr(2)
	s.Deposit(appID1, pool.Id, liquidityProvider2, "2000000000uasset1,2000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("19999999999pool1-1").IsEqual(s.getBalances(liquidityProvider2)))
	msg = types.NewMsgSoftLock(appID1, liquidityProvider2, pool.Id, utils.ParseCoin("19999999999pool1-1"))
	err = s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)

	liquidityProvider3 := s.addr(3)
	s.Deposit(appID1, pool.Id, liquidityProvider3, "3000000000uasset1,3000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("29999999999pool1-1").IsEqual(s.getBalances(liquidityProvider3)))
	msg = types.NewMsgSoftLock(appID1, liquidityProvider3, pool.Id, utils.ParseCoin("29999999999pool1-1"))
	err = s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)

	liquidityProvider4 := s.addr(4)
	s.Deposit(appID1, pool.Id, liquidityProvider4, "4000000000uasset1,4000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("39999999999pool1-1").IsEqual(s.getBalances(liquidityProvider4)))
	msg = types.NewMsgSoftLock(appID1, liquidityProvider4, pool.Id, utils.ParseCoin("39999999999pool1-1"))
	err = s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)

	lpData, found := s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool.Id)
	s.Require().True(found)
	s.Require().Equal(appID1, lpData.AppId)
	s.Require().Equal(pool.Id, lpData.PoolId)
	s.Require().Len(lpData.QueuedLiquidityProviders, 4)
	s.Require().Len(lpData.LiquidityProviders, 0)

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(types.DefaultFarmingQueueDuration).Add(time.Minute * 10))
	s.nextBlock()

	lpData, found = s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool.Id)
	s.Require().True(found)
	s.Require().Len(lpData.QueuedLiquidityProviders, 0)
	s.Require().Len(lpData.LiquidityProviders, 4)

	liquidityGauge := rewardtypes.LiquidtyGaugeMetaData{
		PoolId:       pool.Id,
		IsMasterPool: true,
		ChildPoolIds: []uint64{},
	}

	rewardDistrData, err := s.keeper.GetFarmingRewardsData(s.ctx, appID1, sdk.NewCoin("ucmdx", newInt(10000000000)), liquidityGauge)
	s.Require().NoError(err)
	s.Require().IsType([]rewardtypes.RewardDistributionDataCollector{}, rewardDistrData)

	mappedResp := map[string]sdk.Coin{}
	for _, d := range rewardDistrData {
		mappedResp[d.RewardReceiver.String()] = d.RewardCoin
	}

	s.Require().True(utils.ParseCoin("1000000000ucmdx").IsEqual(mappedResp[liquidityProvider1.String()]))
	s.Require().True(utils.ParseCoin("1999999999ucmdx").IsEqual(mappedResp[liquidityProvider2.String()]))
	s.Require().True(utils.ParseCoin("2999999999ucmdx").IsEqual(mappedResp[liquidityProvider3.String()]))
	s.Require().True(utils.ParseCoin("4000000000ucmdx").IsEqual(mappedResp[liquidityProvider4.String()]))

	liquidityGauge = rewardtypes.LiquidtyGaugeMetaData{
		PoolId:       pool.Id,
		IsMasterPool: false,
		ChildPoolIds: []uint64{},
	}

	rewardDistrData, err = s.keeper.GetFarmingRewardsData(s.ctx, appID1, sdk.NewCoin("ucmdx", newInt(20000000000)), liquidityGauge)
	s.Require().NoError(err)
	s.Require().IsType([]rewardtypes.RewardDistributionDataCollector{}, rewardDistrData)

	mappedResp = map[string]sdk.Coin{}
	for _, d := range rewardDistrData {
		mappedResp[d.RewardReceiver.String()] = d.RewardCoin
	}

	s.Require().True(utils.ParseCoin("2000000000ucmdx").IsEqual(mappedResp[liquidityProvider1.String()]))
	s.Require().True(utils.ParseCoin("3999999999ucmdx").IsEqual(mappedResp[liquidityProvider2.String()]))
	s.Require().True(utils.ParseCoin("5999999999ucmdx").IsEqual(mappedResp[liquidityProvider3.String()]))
	s.Require().True(utils.ParseCoin("8000000000ucmdx").IsEqual(mappedResp[liquidityProvider4.String()]))
}

// Equal liquidity provided my by all liquidity providers
func (s *KeeperTestSuite) TestGetFarmingRewardsDataEqualLPs() {
	currentTime := s.ctx.BlockTime()
	s.ctx = s.ctx.WithBlockTime(currentTime)

	creator := s.addr(0)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSET1", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSET2", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, creator, "1000000000000uasset1,1000000000000uasset2")

	liquidityProvider1 := s.addr(1)
	s.Deposit(appID1, pool.Id, liquidityProvider1, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("10000000000pool1-1").IsEqual(s.getBalances(liquidityProvider1)))
	msg := types.NewMsgSoftLock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("10000000000pool1-1"))
	err := s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)

	liquidityProvider2 := s.addr(2)
	s.Deposit(appID1, pool.Id, liquidityProvider2, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("9999999999pool1-1").IsEqual(s.getBalances(liquidityProvider2)))
	msg = types.NewMsgSoftLock(appID1, liquidityProvider2, pool.Id, utils.ParseCoin("9999999999pool1-1"))
	err = s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)

	liquidityProvider3 := s.addr(3)
	s.Deposit(appID1, pool.Id, liquidityProvider3, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("9999999999pool1-1").IsEqual(s.getBalances(liquidityProvider3)))
	msg = types.NewMsgSoftLock(appID1, liquidityProvider3, pool.Id, utils.ParseCoin("9999999999pool1-1"))
	err = s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)

	liquidityProvider4 := s.addr(4)
	s.Deposit(appID1, pool.Id, liquidityProvider4, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("9999999999pool1-1").IsEqual(s.getBalances(liquidityProvider4)))
	msg = types.NewMsgSoftLock(appID1, liquidityProvider4, pool.Id, utils.ParseCoin("9999999999pool1-1"))
	err = s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)

	lpData, found := s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool.Id)
	s.Require().True(found)
	s.Require().Equal(appID1, lpData.AppId)
	s.Require().Equal(pool.Id, lpData.PoolId)
	s.Require().Len(lpData.QueuedLiquidityProviders, 4)
	s.Require().Len(lpData.LiquidityProviders, 0)

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(types.DefaultFarmingQueueDuration).Add(time.Minute * 10))
	s.nextBlock()

	lpData, found = s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool.Id)
	s.Require().True(found)
	s.Require().Len(lpData.QueuedLiquidityProviders, 0)
	s.Require().Len(lpData.LiquidityProviders, 4)

	liquidityGauge := rewardtypes.LiquidtyGaugeMetaData{
		PoolId:       pool.Id,
		IsMasterPool: true,
		ChildPoolIds: []uint64{},
	}

	rewardDistrData, err := s.keeper.GetFarmingRewardsData(s.ctx, appID1, sdk.NewCoin("ucmdx", newInt(10000000000)), liquidityGauge)
	s.Require().NoError(err)
	s.Require().IsType([]rewardtypes.RewardDistributionDataCollector{}, rewardDistrData)

	mappedResp := map[string]sdk.Coin{}
	for _, d := range rewardDistrData {
		mappedResp[d.RewardReceiver.String()] = d.RewardCoin
	}

	s.Require().True(utils.ParseCoin("2500000001ucmdx").IsEqual(mappedResp[liquidityProvider1.String()]))
	s.Require().True(utils.ParseCoin("2499999999ucmdx").IsEqual(mappedResp[liquidityProvider2.String()]))
	s.Require().True(utils.ParseCoin("2499999999ucmdx").IsEqual(mappedResp[liquidityProvider3.String()]))
	s.Require().True(utils.ParseCoin("2499999999ucmdx").IsEqual(mappedResp[liquidityProvider4.String()]))

	liquidityGauge = rewardtypes.LiquidtyGaugeMetaData{
		PoolId:       pool.Id,
		IsMasterPool: false,
		ChildPoolIds: []uint64{},
	}

	rewardDistrData, err = s.keeper.GetFarmingRewardsData(s.ctx, appID1, sdk.NewCoin("ucmdx", newInt(20000000000)), liquidityGauge)
	s.Require().NoError(err)
	s.Require().IsType([]rewardtypes.RewardDistributionDataCollector{}, rewardDistrData)

	mappedResp = map[string]sdk.Coin{}
	for _, d := range rewardDistrData {
		mappedResp[d.RewardReceiver.String()] = d.RewardCoin
	}

	s.Require().True(utils.ParseCoin("5000000003ucmdx").IsEqual(mappedResp[liquidityProvider1.String()]))
	s.Require().True(utils.ParseCoin("4999999998ucmdx").IsEqual(mappedResp[liquidityProvider2.String()]))
	s.Require().True(utils.ParseCoin("4999999998ucmdx").IsEqual(mappedResp[liquidityProvider3.String()]))
	s.Require().True(utils.ParseCoin("4999999998ucmdx").IsEqual(mappedResp[liquidityProvider4.String()]))
}

// no liquidity providers
func (s *KeeperTestSuite) TestGetFarmingRewardsDataNoLPs() {
	currentTime := s.ctx.BlockTime()
	s.ctx = s.ctx.WithBlockTime(currentTime)

	creator := s.addr(0)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSET1", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSET2", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, creator, "1000000000000uasset1,1000000000000uasset2")

	_, found := s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool.Id)
	s.Require().False(found)

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(types.DefaultFarmingQueueDuration).Add(time.Minute * 10))
	s.nextBlock()

	liquidityGauge := rewardtypes.LiquidtyGaugeMetaData{
		PoolId:       pool.Id,
		IsMasterPool: true,
		ChildPoolIds: []uint64{},
	}

	rewardDistrData, err := s.keeper.GetFarmingRewardsData(s.ctx, appID1, sdk.NewCoin("ucmdx", newInt(10000000000)), liquidityGauge)
	s.Require().NoError(err)
	s.Require().IsType([]rewardtypes.RewardDistributionDataCollector{}, rewardDistrData)
	s.Require().Len(rewardDistrData, 0)

	liquidityGauge = rewardtypes.LiquidtyGaugeMetaData{
		PoolId:       pool.Id,
		IsMasterPool: false,
		ChildPoolIds: []uint64{},
	}

	rewardDistrData, err = s.keeper.GetFarmingRewardsData(s.ctx, appID1, sdk.NewCoin("ucmdx", newInt(20000000000)), liquidityGauge)
	s.Require().NoError(err)
	s.Require().IsType([]rewardtypes.RewardDistributionDataCollector{}, rewardDistrData)
	s.Require().Len(rewardDistrData, 0)
}

// create 2 pools, one master and another child
func (s *KeeperTestSuite) TestGetFarmingRewardsDataEqualLPsWChildPool() {
	currentTime := s.ctx.BlockTime()
	s.ctx = s.ctx.WithBlockTime(currentTime)

	creator := s.addr(0)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSET1", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSET2", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, creator, "1000000000000uasset1,1000000000000uasset2")

	pair2 := s.CreateNewLiquidityPair(appID1, creator, asset2.Denom, asset1.Denom)
	pool2 := s.CreateNewLiquidityPool(appID1, pair2.Id, creator, "1000000000000uasset1,1000000000000uasset2")

	// lp1 - farming only in master pool, not child pool (not eligible for masterpool type reward)
	liquidityProvider1 := s.addr(1)
	s.Deposit(appID1, pool.Id, liquidityProvider1, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("10000000000pool1-1").IsEqual(s.getBalances(liquidityProvider1)))
	msg := types.NewMsgSoftLock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("10000000000pool1-1"))
	err := s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)

	// lp2 - farming in master pool and child pool (eligible for masterpool type reward)
	liquidityProvider2 := s.addr(2)
	s.Deposit(appID1, pool.Id, liquidityProvider2, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Deposit(appID1, pool2.Id, liquidityProvider2, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("9999999999pool1-1,10000000000pool1-2").IsEqual(s.getBalances(liquidityProvider2)))
	msg = types.NewMsgSoftLock(appID1, liquidityProvider2, pool.Id, utils.ParseCoin("9999999999pool1-1"))
	err = s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)
	msg = types.NewMsgSoftLock(appID1, liquidityProvider2, pool2.Id, utils.ParseCoin("10000000000pool1-2"))
	err = s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)

	// lp3 - farming only in master pool, not child pool (not eligible for masterpool type reward)
	liquidityProvider3 := s.addr(3)
	s.Deposit(appID1, pool.Id, liquidityProvider3, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("9999999999pool1-1").IsEqual(s.getBalances(liquidityProvider3)))
	msg = types.NewMsgSoftLock(appID1, liquidityProvider3, pool.Id, utils.ParseCoin("9999999999pool1-1"))
	err = s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)

	// lp4 - farming  in master pool and  child pool (eligible for masterpool type reward)
	liquidityProvider4 := s.addr(4)
	s.Deposit(appID1, pool.Id, liquidityProvider4, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Deposit(appID1, pool2.Id, liquidityProvider4, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("9999999999pool1-1,9999999999pool1-2").IsEqual(s.getBalances(liquidityProvider4)))
	msg = types.NewMsgSoftLock(appID1, liquidityProvider4, pool.Id, utils.ParseCoin("9999999999pool1-1"))
	err = s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)
	msg = types.NewMsgSoftLock(appID1, liquidityProvider4, pool2.Id, utils.ParseCoin("9999999999pool1-2"))
	err = s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)

	lpData, found := s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool.Id)
	s.Require().True(found)
	s.Require().Equal(appID1, lpData.AppId)
	s.Require().Equal(pool.Id, lpData.PoolId)
	s.Require().Len(lpData.QueuedLiquidityProviders, 4)
	s.Require().Len(lpData.LiquidityProviders, 0)

	lpData, found = s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool2.Id)
	s.Require().True(found)
	s.Require().Equal(appID1, lpData.AppId)
	s.Require().Equal(pool2.Id, lpData.PoolId)
	s.Require().Len(lpData.QueuedLiquidityProviders, 2)
	s.Require().Len(lpData.LiquidityProviders, 0)

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(types.DefaultFarmingQueueDuration).Add(time.Minute * 10))
	s.nextBlock()

	lpData, found = s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool.Id)
	s.Require().True(found)
	s.Require().Len(lpData.QueuedLiquidityProviders, 0)
	s.Require().Len(lpData.LiquidityProviders, 4)

	lpData, found = s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool2.Id)
	s.Require().True(found)
	s.Require().Len(lpData.QueuedLiquidityProviders, 0)
	s.Require().Len(lpData.LiquidityProviders, 2)

	liquidityGauge := rewardtypes.LiquidtyGaugeMetaData{
		PoolId:       pool.Id,
		IsMasterPool: true,
		ChildPoolIds: []uint64{},
	}

	rewardDistrData, err := s.keeper.GetFarmingRewardsData(s.ctx, appID1, sdk.NewCoin("ucmdx", newInt(10000000000)), liquidityGauge)
	s.Require().NoError(err)
	s.Require().IsType([]rewardtypes.RewardDistributionDataCollector{}, rewardDistrData)

	mappedResp := map[string]sdk.Coin{}
	for _, d := range rewardDistrData {
		mappedResp[d.RewardReceiver.String()] = d.RewardCoin
	}

	s.Require().True(utils.ParseCoin("5000000000ucmdx").IsEqual(mappedResp[liquidityProvider2.String()]))
	s.Require().True(utils.ParseCoin("5000000000ucmdx").IsEqual(mappedResp[liquidityProvider4.String()]))

	liquidityGauge = rewardtypes.LiquidtyGaugeMetaData{
		PoolId:       pool.Id,
		IsMasterPool: true,
		ChildPoolIds: []uint64{1, 2},
	}

	rewardDistrData, err = s.keeper.GetFarmingRewardsData(s.ctx, appID1, sdk.NewCoin("ucmdx", newInt(40000000000)), liquidityGauge)
	s.Require().NoError(err)
	s.Require().IsType([]rewardtypes.RewardDistributionDataCollector{}, rewardDistrData)

	mappedResp = map[string]sdk.Coin{}
	for _, d := range rewardDistrData {
		mappedResp[d.RewardReceiver.String()] = d.RewardCoin
	}

	s.Require().True(utils.ParseCoin("20000000000ucmdx").IsEqual(mappedResp[liquidityProvider2.String()]))
	s.Require().True(utils.ParseCoin("20000000000ucmdx").IsEqual(mappedResp[liquidityProvider4.String()]))

	liquidityGauge = rewardtypes.LiquidtyGaugeMetaData{
		PoolId:       pool.Id,
		IsMasterPool: false,
		ChildPoolIds: []uint64{},
	}

	rewardDistrData, err = s.keeper.GetFarmingRewardsData(s.ctx, appID1, sdk.NewCoin("ucmdx", newInt(20000000000)), liquidityGauge)
	s.Require().NoError(err)
	s.Require().IsType([]rewardtypes.RewardDistributionDataCollector{}, rewardDistrData)

	mappedResp = map[string]sdk.Coin{}
	for _, d := range rewardDistrData {
		mappedResp[d.RewardReceiver.String()] = d.RewardCoin
	}

	s.Require().True(utils.ParseCoin("5000000003ucmdx").IsEqual(mappedResp[liquidityProvider1.String()]))
	s.Require().True(utils.ParseCoin("4999999998ucmdx").IsEqual(mappedResp[liquidityProvider2.String()]))
	s.Require().True(utils.ParseCoin("4999999998ucmdx").IsEqual(mappedResp[liquidityProvider3.String()]))
	s.Require().True(utils.ParseCoin("4999999998ucmdx").IsEqual(mappedResp[liquidityProvider4.String()]))
}

// pool reserve sent to somewhere else, and maked pool as depleted and disabled
func (s *KeeperTestSuite) TestGetFarmingRewardsDataErrorHandellings() {
	currentTime := s.ctx.BlockTime()
	s.ctx = s.ctx.WithBlockTime(currentTime)

	creator := s.addr(0)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSET1", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSET2", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, creator, "1000000000000uasset1,1000000000000uasset2")

	// lp1 - farming only in master pool, not child pool (not eligible for masterpool type reward)
	liquidityProvider1 := s.addr(1)
	s.Deposit(appID1, pool.Id, liquidityProvider1, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("10000000000pool1-1").IsEqual(s.getBalances(liquidityProvider1)))
	msg := types.NewMsgSoftLock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("10000000000pool1-1"))
	err := s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)

	s.sendCoins(pool.GetReserveAddress(), creator, s.getBalances(pool.GetReserveAddress()))

	liquidityGauge := rewardtypes.LiquidtyGaugeMetaData{
		PoolId:       pool.Id,
		IsMasterPool: true,
		ChildPoolIds: []uint64{},
	}

	rewardDistrData, err := s.keeper.GetFarmingRewardsData(s.ctx, appID1, sdk.NewCoin("ucmdx", newInt(10000000000)), liquidityGauge)
	s.Require().Error(err)
	s.Require().EqualError(err, sdkerrors.Wrapf(types.ErrDepletedPool, "pool 1 is depleted").Error())
	s.Require().IsType([]rewardtypes.RewardDistributionDataCollector{}, rewardDistrData)

	s.Deposit(appID1, pool.Id, liquidityProvider1, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	rewardDistrData, err = s.keeper.GetFarmingRewardsData(s.ctx, appID1, sdk.NewCoin("ucmdx", newInt(10000000000)), liquidityGauge)
	s.Require().Error(err)
	s.Require().EqualError(err, sdkerrors.Wrapf(types.ErrDisabledPool, "pool 1 is disabled").Error())
	s.Require().IsType([]rewardtypes.RewardDistributionDataCollector{}, rewardDistrData)
}

// LP added =>  farmed => unfarmed
func (s *KeeperTestSuite) TestGetFarmingRewardsDataZeroLPs() {
	currentTime := s.ctx.BlockTime()
	s.ctx = s.ctx.WithBlockTime(currentTime)

	creator := s.addr(0)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSET1", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSET2", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, creator, "1000000000000uasset1,1000000000000uasset2")

	liquidityProvider1 := s.addr(1)
	s.Deposit(appID1, pool.Id, liquidityProvider1, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("10000000000pool1-1").IsEqual(s.getBalances(liquidityProvider1)))
	msg := types.NewMsgSoftLock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("10000000000pool1-1"))
	err := s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)

	liquidityProvider2 := s.addr(2)
	s.Deposit(appID1, pool.Id, liquidityProvider2, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("9999999999pool1-1").IsEqual(s.getBalances(liquidityProvider2)))
	msg = types.NewMsgSoftLock(appID1, liquidityProvider2, pool.Id, utils.ParseCoin("9999999999pool1-1"))
	err = s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)

	liquidityProvider3 := s.addr(3)
	s.Deposit(appID1, pool.Id, liquidityProvider3, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("9999999999pool1-1").IsEqual(s.getBalances(liquidityProvider3)))
	msg = types.NewMsgSoftLock(appID1, liquidityProvider3, pool.Id, utils.ParseCoin("9999999999pool1-1"))
	err = s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)

	liquidityProvider4 := s.addr(4)
	s.Deposit(appID1, pool.Id, liquidityProvider4, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("9999999999pool1-1").IsEqual(s.getBalances(liquidityProvider4)))
	msg = types.NewMsgSoftLock(appID1, liquidityProvider4, pool.Id, utils.ParseCoin("9999999999pool1-1"))
	err = s.keeper.SoftLockTokens(s.ctx, msg)
	s.Require().NoError(err)

	lpData, found := s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool.Id)
	s.Require().True(found)
	s.Require().Equal(appID1, lpData.AppId)
	s.Require().Equal(pool.Id, lpData.PoolId)
	s.Require().Len(lpData.QueuedLiquidityProviders, 4)
	s.Require().Len(lpData.LiquidityProviders, 0)

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(types.DefaultFarmingQueueDuration).Add(time.Minute * 10))
	s.nextBlock()

	lpData, found = s.keeper.GetPoolLiquidityProvidersData(s.ctx, appID1, pool.Id)
	s.Require().True(found)
	s.Require().Len(lpData.QueuedLiquidityProviders, 0)
	s.Require().Len(lpData.LiquidityProviders, 4)

	msgUnlock := types.NewMsgSoftUnlock(appID1, liquidityProvider1, pool.Id, utils.ParseCoin("10000000000pool1-1"))
	err = s.keeper.SoftUnlockTokens(s.ctx, msgUnlock)
	s.Require().NoError(err)
	msgUnlock = types.NewMsgSoftUnlock(appID1, liquidityProvider2, pool.Id, utils.ParseCoin("9999999999pool1-1"))
	err = s.keeper.SoftUnlockTokens(s.ctx, msgUnlock)
	s.Require().NoError(err)
	msgUnlock = types.NewMsgSoftUnlock(appID1, liquidityProvider3, pool.Id, utils.ParseCoin("9999999999pool1-1"))
	err = s.keeper.SoftUnlockTokens(s.ctx, msgUnlock)
	s.Require().NoError(err)
	msgUnlock = types.NewMsgSoftUnlock(appID1, liquidityProvider4, pool.Id, utils.ParseCoin("9999999999pool1-1"))
	err = s.keeper.SoftUnlockTokens(s.ctx, msgUnlock)
	s.Require().NoError(err)

	liquidityGauge := rewardtypes.LiquidtyGaugeMetaData{
		PoolId:       pool.Id,
		IsMasterPool: true,
		ChildPoolIds: []uint64{},
	}

	rewardDistrData, err := s.keeper.GetFarmingRewardsData(s.ctx, appID1, sdk.NewCoin("ucmdx", newInt(10000000000)), liquidityGauge)
	s.Require().NoError(err)
	s.Require().IsType([]rewardtypes.RewardDistributionDataCollector{}, rewardDistrData)
	s.Require().Len(rewardDistrData, 0)

	liquidityGauge = rewardtypes.LiquidtyGaugeMetaData{
		PoolId:       pool.Id,
		IsMasterPool: false,
		ChildPoolIds: []uint64{},
	}

	rewardDistrData, err = s.keeper.GetFarmingRewardsData(s.ctx, appID1, sdk.NewCoin("ucmdx", newInt(20000000000)), liquidityGauge)
	s.Require().NoError(err)
	s.Require().IsType([]rewardtypes.RewardDistributionDataCollector{}, rewardDistrData)
	s.Require().Len(rewardDistrData, 0)
}
