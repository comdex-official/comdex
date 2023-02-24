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

func (s *KeeperTestSuite) TestFarm() {
	creator := s.addr(0)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

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
		Msg              types.MsgFarm
		ExpErr           error
		AvailableBalance sdk.Coins
		QueueLenght      uint64
	}{
		{
			Name:             "error app id invalid",
			Msg:              *types.NewMsgFarm(69, pool.Id, liquidityProvider1, utils.ParseCoin("699000000pool1-1")),
			ExpErr:           sdkerrors.Wrapf(types.ErrInvalidAppID, "app id %d not found", 69),
			AvailableBalance: utils.ParseCoins("10000000000pool1-1,10000000000pool2-1"),
		},
		{
			Name:             "error pool id invalid",
			Msg:              *types.NewMsgFarm(appID1, 69, liquidityProvider1, utils.ParseCoin("699000000pool1-1")),
			ExpErr:           sdkerrors.Wrapf(types.ErrInvalidPoolID, "no pool exists with id : %d", 69),
			AvailableBalance: utils.ParseCoins("10000000000pool1-1,10000000000pool2-1"),
		},
		{
			Name:             "error pool denom invalid",
			Msg:              *types.NewMsgFarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("699000000pool1-2")),
			ExpErr:           sdkerrors.Wrapf(types.ErrWrongPoolCoinDenom, "expected pool coin denom %s, found pool1-2", pool.PoolCoinDenom),
			AvailableBalance: utils.ParseCoins("10000000000pool1-1,10000000000pool2-1"),
		},
		{
			Name:             "error insufficient pool denoms",
			Msg:              *types.NewMsgFarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("100000000000pool1-1")),
			ExpErr:           sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "10000000000pool1-1 is smaller than 100000000000pool1-1"),
			AvailableBalance: utils.ParseCoins("10000000000pool1-1,10000000000pool2-1"),
		},
		{
			Name:             "success liquidity provider 1 app1",
			Msg:              *types.NewMsgFarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("5252000000pool1-1")),
			ExpErr:           nil,
			AvailableBalance: utils.ParseCoins("4748000000pool1-1,10000000000pool2-1,5252000000farm1-1"),
			QueueLenght:      1,
		},
		{
			Name:             "success liquidity provider 2 app1",
			Msg:              *types.NewMsgFarm(appID1, pool.Id, liquidityProvider2, utils.ParseCoin("6934000000pool1-1")),
			ExpErr:           nil,
			AvailableBalance: utils.ParseCoins("3065999999pool1-1,9999999999pool2-1,6934000000farm1-1"),
			QueueLenght:      1,
		},
		{
			Name:             "success liquidity provider 1 app1 re-add",
			Msg:              *types.NewMsgFarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("23000000pool1-1")),
			ExpErr:           nil,
			AvailableBalance: utils.ParseCoins("4725000000pool1-1,10000000000pool2-1,5275000000farm1-1"),
			QueueLenght:      2,
		},
		{
			Name:             "success liquidity provider 1 app2",
			Msg:              *types.NewMsgFarm(appID2, pool2.Id, liquidityProvider1, utils.ParseCoin("123000000pool2-1")),
			ExpErr:           nil,
			AvailableBalance: utils.ParseCoins("4725000000pool1-1,9877000000pool2-1,5275000000farm1-1,123000000farm2-1"),
			QueueLenght:      1,
		},
		{
			Name:             "success liquidity provider 2 app2",
			Msg:              *types.NewMsgFarm(appID2, pool2.Id, liquidityProvider2, utils.ParseCoin("546000000pool2-1")),
			ExpErr:           nil,
			AvailableBalance: utils.ParseCoins("3065999999pool1-1,9453999999pool2-1,6934000000farm1-1,546000000farm2-1"),
			QueueLenght:      1,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			err := s.keeper.Farm(s.ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				_, found := s.keeper.GetQueuedFarmer(s.ctx, tc.Msg.AppId, tc.Msg.PoolId, tc.Msg.GetFarmer())
				s.Require().False(found)
				_, found = s.keeper.GetActiveFarmer(s.ctx, tc.Msg.AppId, tc.Msg.PoolId, tc.Msg.GetFarmer())
				s.Require().False(found)
			} else {
				s.Require().NoError(err)
				s.Require().True(tc.AvailableBalance.IsEqual(s.getBalances(sdk.MustAccAddressFromBech32(tc.Msg.Farmer))))

				queuedFarmer, found := s.keeper.GetQueuedFarmer(s.ctx, tc.Msg.AppId, tc.Msg.PoolId, tc.Msg.GetFarmer())
				s.Require().True(found)

				s.Require().Equal(tc.Msg.AppId, queuedFarmer.AppId)
				s.Require().Equal(tc.Msg.PoolId, queuedFarmer.PoolId)
				s.Require().Len(queuedFarmer.QueudCoins, int(tc.QueueLenght))

				s.Require().Equal(tc.Msg.Farmer, queuedFarmer.Farmer)
				s.Require().Equal(tc.Msg.FarmingPoolCoin, queuedFarmer.QueudCoins[tc.QueueLenght-1].FarmedPoolCoin)
			}
		})
	}

	// increase time and check if deque works
	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(types.DefaultFarmingQueueDuration).Add(time.Second * 10))
	s.nextBlock()
	// app1 check
	activeFarmers := s.keeper.GetAllActiveFarmers(s.ctx, appID1, pool.Id)
	queuedFarmers := s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool.Id)
	s.Require().Len(queuedFarmers[0].QueudCoins, 0)
	s.Require().Len(queuedFarmers[1].QueudCoins, 0)
	s.Require().Len(activeFarmers, 2)

	s.Require().True(activeFarmers[0].FarmedPoolCoin.IsEqual(utils.ParseCoin("5275000000pool1-1")))
	s.Require().True(activeFarmers[1].FarmedPoolCoin.IsEqual(utils.ParseCoin("6934000000pool1-1")))

	_, found := s.keeper.GetActiveFarmer(s.ctx, appID1, pool.Id, creator)
	s.Require().False(found)

	// app2 check
	activeFarmers = s.keeper.GetAllActiveFarmers(s.ctx, appID2, pool2.Id)
	queuedFarmers = s.keeper.GetAllQueuedFarmers(s.ctx, appID2, pool2.Id)
	s.Require().Len(queuedFarmers[0].QueudCoins, 0)
	s.Require().Len(queuedFarmers[1].QueudCoins, 0)
	s.Require().Len(activeFarmers, 2)

	s.Require().True(activeFarmers[0].FarmedPoolCoin.IsEqual(utils.ParseCoin("123000000pool2-1")))
	s.Require().True(activeFarmers[1].FarmedPoolCoin.IsEqual(utils.ParseCoin("546000000pool2-1")))

	_, found = s.keeper.GetActiveFarmer(s.ctx, appID2, pool2.Id, creator)
	s.Require().False(found)
}

func (s *KeeperTestSuite) TestUnfarm() {
	creator := s.addr(0)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

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

	msg := types.NewMsgFarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("10000000000pool1-1"))
	err := s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)
	s.Require().True(utils.ParseCoins("10000000000farm1-1").IsEqual(s.getBalances(liquidityProvider1)))
	queuedFarmers := s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool.Id)
	s.Require().Len(queuedFarmers, 1)

	testCases := []struct {
		Name             string
		Msg              types.MsgUnfarm
		ExpErr           error
		AvailableBalance sdk.Coins
		QueueLenght      uint64
	}{
		{
			Name:             "error app id invalid",
			Msg:              *types.NewMsgUnfarm(69, pool.Id, liquidityProvider1, utils.ParseCoin("699000000farm1-1")),
			ExpErr:           sdkerrors.Wrapf(types.ErrInvalidAppID, "app id %d not found", 69),
			AvailableBalance: sdk.Coins{},
		},
		{
			Name:             "error pool id invalid",
			Msg:              *types.NewMsgUnfarm(appID1, 69, liquidityProvider1, utils.ParseCoin("699000000farm1-1")),
			ExpErr:           sdkerrors.Wrapf(types.ErrInvalidPoolID, "no pool exists with id : %d", 69),
			AvailableBalance: sdk.Coins{},
		},
		{
			Name:             "error pool denom invalid",
			Msg:              *types.NewMsgUnfarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("699000000pool1-2")),
			ExpErr:           sdkerrors.Wrapf(types.ErrWrongPoolCoinDenom, "expected farm coin denom %s, found pool1-2", pool.FarmCoin.Denom),
			AvailableBalance: sdk.Coins{},
		},
		{
			Name:             "error farm not found",
			Msg:              *types.NewMsgUnfarm(appID1, pool2.Id, liquidityProvider1, utils.ParseCoin("699000000farm1-2")),
			ExpErr:           sdkerrors.Wrapf(types.ErrorFarmerNotFound, "no active farm found for given pool id %d", pool2.Id),
			AvailableBalance: sdk.Coins{},
		},
		{
			Name:             "error insufficient farmed amounts",
			Msg:              *types.NewMsgUnfarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("100000000000farm1-1")),
			ExpErr:           sdkerrors.Wrapf(types.ErrInvalidUnfarmAmount, "farmed pool coin amount 10000000000farm1-1 smaller than requested unfarming pool coin amount 100000000000farm1-1"),
			AvailableBalance: sdk.Coins{},
		},
		{
			Name:             "success partial unlock",
			Msg:              *types.NewMsgUnfarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("5000000000farm1-1")),
			ExpErr:           nil,
			AvailableBalance: utils.ParseCoins("5000000000pool1-1,5000000000farm1-1"),
		},
		{
			Name:             "success full unlock",
			Msg:              *types.NewMsgUnfarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("5000000000farm1-1")),
			ExpErr:           nil,
			AvailableBalance: utils.ParseCoins("10000000000pool1-1"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			err := s.keeper.Unfarm(s.ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().True(tc.AvailableBalance.IsEqual(s.getBalances(sdk.MustAccAddressFromBech32(tc.Msg.Farmer))))
			}
		})
	}
}

func (s *KeeperTestSuite) TestUnfarmTwo() {
	currentTime := s.ctx.BlockTime()
	s.ctx = s.ctx.WithBlockTime(currentTime)

	creator := s.addr(0)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

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
	msg := types.NewMsgFarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("10000000pool1-1"))
	err := s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)
	s.Require().True(utils.ParseCoins("9990000000pool1-1,10000000farm1-1").IsEqual(s.getBalances(liquidityProvider1)))

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(time.Hour * 1))
	// farm 2, queue size 2
	// SortedByTimeFarmQueue -> [20000000pool1-1, 10000000pool1-1]
	msg = types.NewMsgFarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("20000000pool1-1"))
	err = s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)
	s.Require().True(utils.ParseCoins("9970000000pool1-1,30000000farm1-1").IsEqual(s.getBalances(liquidityProvider1)))

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(time.Hour * 1))
	// farm 3, queue size 3
	// SortedByTimeFarmQueue -> [30000000pool1-1, 20000000pool1-1, 10000000pool1-1]
	msg = types.NewMsgFarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("30000000pool1-1"))
	err = s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)
	s.Require().True(utils.ParseCoins("9940000000pool1-1,60000000farm1-1").IsEqual(s.getBalances(liquidityProvider1)))

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(time.Hour * 1))
	// farm 4, queue size 4
	// SortedByTimeFarmQueue -> [40000000pool1-1, 30000000pool1-1, 20000000pool1-1, 10000000pool1-1]
	msg = types.NewMsgFarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("40000000pool1-1"))
	err = s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)
	s.Require().True(utils.ParseCoins("9900000000pool1-1,100000000farm1-1").IsEqual(s.getBalances(liquidityProvider1)))

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(time.Hour * 1))
	// farm 5, queue size 5
	// SortedByTimeFarmQueue -> [50000000pool1-1, 40000000pool1-1, 30000000pool1-1, 20000000pool1-1, 10000000pool1-1]
	msg = types.NewMsgFarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("50000000pool1-1"))
	err = s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)
	s.Require().True(utils.ParseCoins("9850000000pool1-1,150000000farm1-1").IsEqual(s.getBalances(liquidityProvider1)))

	queuedFarmers := s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool.Id)
	s.Require().Len(queuedFarmers, 1)
	s.Require().Len(queuedFarmers[0].QueudCoins, 5)

	// lp1 trying to unfarm more than farmed
	msgUnlock := types.NewMsgUnfarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("160000000farm1-1"))
	err = s.keeper.Unfarm(s.ctx, msgUnlock)
	s.Require().Error(err)
	s.Require().EqualError(err, sdkerrors.Wrapf(types.ErrInvalidUnfarmAmount, "farmed pool coin amount 150000000farm1-1 smaller than requested unfarming pool coin amount 160000000farm1-1").Error())

	// unfarming small portions, below unlock removes token from most recently added queue
	// unlock is done from a single latest object in a queue since this object itself can satisfy the unlock requirement,
	// Before - SortedByTimeFarmQueue -> [50000000pool1-1, 40000000pool1-1, 30000000pool1-1, 20000000pool1-1, 10000000pool1-1]
	// After  - SortedByTimeFarmQueue -> [45000000pool1-1, 40000000pool1-1, 30000000pool1-1, 20000000pool1-1, 10000000pool1-1]
	msgUnlock = types.NewMsgUnfarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("5000000farm1-1"))
	err = s.keeper.Unfarm(s.ctx, msgUnlock)
	s.Require().NoError(err)
	queuedFarmers = s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool.Id)
	s.Require().Len(queuedFarmers, 1)
	s.Require().Len(queuedFarmers[0].QueudCoins, 5)
	s.Require().Equal(utils.ParseCoin("45000000pool1-1").Denom, queuedFarmers[0].QueudCoins[4].FarmedPoolCoin.Denom)
	s.Require().Equal(utils.ParseCoin("45000000pool1-1").Amount, queuedFarmers[0].QueudCoins[4].FarmedPoolCoin.Amount)
	s.Require().True(utils.ParseCoins("9855000000pool1-1,145000000farm1-1").IsEqual(s.getBalances(liquidityProvider1)))

	// unfarming small portions, below unlock removes token from most recently added queue
	// unlock is done from a single latest object in a queue since this object itself can satisfy the unlock requirement,
	// Before  - SortedByTimeFarmQueue -> [45000000pool1-1, 40000000pool1-1, 30000000pool1-1, 20000000pool1-1, 10000000pool1-1]
	// After   - SortedByTimeFarmQueue -> [34000000pool1-1, 40000000pool1-1, 30000000pool1-1, 20000000pool1-1, 10000000pool1-1]
	msgUnlock = types.NewMsgUnfarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("11000000farm1-1"))
	err = s.keeper.Unfarm(s.ctx, msgUnlock)
	s.Require().NoError(err)
	queuedFarmers = s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool.Id)
	s.Require().Len(queuedFarmers, 1)
	s.Require().Len(queuedFarmers[0].QueudCoins, 5)
	s.Require().Equal(utils.ParseCoin("34000000pool1-1").Denom, queuedFarmers[0].QueudCoins[4].FarmedPoolCoin.Denom)
	s.Require().Equal(utils.ParseCoin("34000000pool1-1").Amount, queuedFarmers[0].QueudCoins[4].FarmedPoolCoin.Amount)
	s.Require().True(utils.ParseCoins("9866000000pool1-1,134000000farm1-1").IsEqual(s.getBalances(liquidityProvider1)))

	// below case will delete the most recent object from queue since it satisfies the required unlock condition
	// here the unlock is being satisfied from the two queue objects, most recent one gets deleted after it fullfills all
	// of its token for unlocking, and the remaining unlock tokens are fullfilled from 2nd most recent queue object
	// Before - SortedByTimeFarmQueue -> [34000000pool1-1, 40000000pool1-1, 30000000pool1-1, 20000000pool1-1, 10000000pool1-1]
	// After   - SortedByTimeFarmQueue -> [36000000pool1-1, 30000000pool1-1, 20000000pool1-1, 10000000pool1-1]
	msgUnlock = types.NewMsgUnfarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("38000000farm1-1"))
	err = s.keeper.Unfarm(s.ctx, msgUnlock)
	s.Require().NoError(err)
	queuedFarmers = s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool.Id)
	s.Require().Len(queuedFarmers, 1)
	s.Require().Len(queuedFarmers[0].QueudCoins, 4)
	s.Require().Equal(utils.ParseCoin("36000000pool1-1").Denom, queuedFarmers[0].QueudCoins[3].FarmedPoolCoin.Denom)
	s.Require().Equal(utils.ParseCoin("36000000pool1-1").Amount, queuedFarmers[0].QueudCoins[3].FarmedPoolCoin.Amount)
	s.Require().True(utils.ParseCoins("9904000000pool1-1,96000000farm1-1").IsEqual(s.getBalances(liquidityProvider1)))

	// similarly below cases are followed as above
	// Before   - SortedByTimeFarmQueue -> [36000000pool1-1, 30000000pool1-1, 20000000pool1-1, 10000000pool1-1]
	// After    - SortedByTimeFarmQueue -> [30000000pool1-1, 20000000pool1-1, 10000000pool1-1]
	msgUnlock = types.NewMsgUnfarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("36000000farm1-1"))
	err = s.keeper.Unfarm(s.ctx, msgUnlock)
	s.Require().NoError(err)
	queuedFarmers = s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool.Id)
	s.Require().Len(queuedFarmers, 1)
	s.Require().Len(queuedFarmers[0].QueudCoins, 3)
	s.Require().Equal(utils.ParseCoin("30000000pool1-1").Denom, queuedFarmers[0].QueudCoins[2].FarmedPoolCoin.Denom)
	s.Require().Equal(utils.ParseCoin("30000000pool1-1").Amount, queuedFarmers[0].QueudCoins[2].FarmedPoolCoin.Amount)
	s.Require().True(utils.ParseCoins("9940000000pool1-1,60000000farm1-1").IsEqual(s.getBalances(liquidityProvider1)))

	// Before    - SortedByTimeFarmQueue -> [30000000pool1-1, 20000000pool1-1, 10000000pool1-1]
	// After     - SortedByTimeFarmQueue -> [10000000pool1-1]
	msgUnlock = types.NewMsgUnfarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("50000000farm1-1"))
	err = s.keeper.Unfarm(s.ctx, msgUnlock)
	s.Require().NoError(err)
	queuedFarmers = s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool.Id)
	s.Require().Len(queuedFarmers, 1)
	s.Require().Len(queuedFarmers[0].QueudCoins, 1)
	s.Require().Equal(utils.ParseCoin("10000000pool1-1").Denom, queuedFarmers[0].QueudCoins[0].FarmedPoolCoin.Denom)
	s.Require().Equal(utils.ParseCoin("10000000pool1-1").Amount, queuedFarmers[0].QueudCoins[0].FarmedPoolCoin.Amount)
	s.Require().True(utils.ParseCoins("9990000000pool1-1,10000000farm1-1").IsEqual(s.getBalances(liquidityProvider1)))

	// lp1 trying to unfarm more than farmed
	msgUnlock = types.NewMsgUnfarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("11000000farm1-1"))
	err = s.keeper.Unfarm(s.ctx, msgUnlock)
	s.Require().Error(err)
	s.Require().EqualError(err, sdkerrors.Wrapf(types.ErrInvalidUnfarmAmount, "farmed pool coin amount 10000000farm1-1 smaller than requested unfarming pool coin amount 11000000farm1-1").Error())

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(time.Hour * 1))
	// SortedByTimeFarmQueue -> [69000000pool1-1, 10000000pool1-1]
	msg = types.NewMsgFarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("69000000pool1-1"))
	err = s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)
	s.Require().True(utils.ParseCoins("9921000000pool1-1,79000000farm1-1").IsEqual(s.getBalances(liquidityProvider1)))

	// marking oldest farmed object as valid and dequing it, assuming queue duration is satisfied
	s.ctx = s.ctx.WithBlockTime(currentTime.Add(types.DefaultFarmingQueueDuration).Add(time.Second * 10))
	s.nextBlock()
	afs := s.keeper.GetAllActiveFarmers(s.ctx, appID1, pool.Id)
	qfs := s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool.Id)
	s.Require().Len(qfs, 1)
	s.Require().Len(afs, 1)

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
	msgUnlock = types.NewMsgUnfarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("70000000farm1-1"))
	err = s.keeper.Unfarm(s.ctx, msgUnlock)
	s.Require().NoError(err)
	afs = s.keeper.GetAllActiveFarmers(s.ctx, appID1, pool.Id)
	qfs = s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool.Id)
	s.Require().Len(qfs[0].QueudCoins, 0)
	s.Require().Len(afs, 1)
	s.Require().Equal(utils.ParseCoin("9000000pool1-1").Denom, afs[0].FarmedPoolCoin.Denom)
	s.Require().Equal(utils.ParseCoin("9000000pool1-1").Amount, afs[0].FarmedPoolCoin.Amount)
	s.Require().True(utils.ParseCoins("9991000000pool1-1,9000000farm1-1").IsEqual(s.getBalances(liquidityProvider1)))

	// unlocking all farmed tokens
	// Before
	// SortedByTimeFarmQueue -> []
	// ActiveFarmedTokens -> [9000000pool1-1]
	// After
	// SortedByTimeFarmQueue -> []
	// ActiveFarmedTokens -> []
	msgUnlock = types.NewMsgUnfarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("9000000farm1-1"))
	err = s.keeper.Unfarm(s.ctx, msgUnlock)
	s.Require().NoError(err)
	afs = s.keeper.GetAllActiveFarmers(s.ctx, appID1, pool.Id)
	qfs = s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool.Id)
	s.Require().Len(qfs[0].QueudCoins, 0)
	s.Require().Len(afs, 0)

	s.Require().True(utils.ParseCoins("10000000000pool1-1").IsEqual(s.getBalances(liquidityProvider1)))

	s.ctx = s.ctx.WithBlockTime(currentTime)
	// SortedByTimeFarmQueue -> [11000000pool1-1]
	// ActiveFarmedTokens -> []
	msg = types.NewMsgFarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("11000000pool1-1"))
	err = s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)
	s.Require().True(utils.ParseCoins("9989000000pool1-1,11000000farm1-1").IsEqual(s.getBalances(liquidityProvider1)))

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(time.Hour * 1))
	// SortedByTimeFarmQueue -> [12000000pool1-1, 11000000pool1-1]
	// ActiveFarmedTokens -> []
	msg = types.NewMsgFarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("12000000pool1-1"))
	err = s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)
	s.Require().True(utils.ParseCoins("9977000000pool1-1,23000000farm1-1").IsEqual(s.getBalances(liquidityProvider1)))

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(time.Hour * 1))
	// SortedByTimeFarmQueue -> [13000000pool1-1, 12000000pool1-1, 11000000pool1-1]
	// ActiveFarmedTokens -> []
	msg = types.NewMsgFarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("13000000pool1-1"))
	err = s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)
	s.Require().True(utils.ParseCoins("9964000000pool1-1,36000000farm1-1").IsEqual(s.getBalances(liquidityProvider1)))

	// marking oldest farmed object as valid and dequing it, assuming queue duration is satisfied
	s.ctx = s.ctx.WithBlockTime(currentTime.Add(types.DefaultFarmingQueueDuration).Add(time.Second * 10))
	s.nextBlock()
	afs = s.keeper.GetAllActiveFarmers(s.ctx, appID1, pool.Id)
	qfs = s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool.Id)
	s.Require().Len(qfs[0].QueudCoins, 2)
	s.Require().Len(afs, 1)

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
	msgUnlock = types.NewMsgUnfarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("26000000farm1-1"))
	err = s.keeper.Unfarm(s.ctx, msgUnlock)
	s.Require().NoError(err)
	afs = s.keeper.GetAllActiveFarmers(s.ctx, appID1, pool.Id)
	qfs = s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool.Id)
	s.Require().Len(qfs[0].QueudCoins, 0)
	s.Require().Len(afs, 1)
	s.Require().Equal(utils.ParseCoin("10000000pool1-1").Denom, afs[0].FarmedPoolCoin.Denom)
	s.Require().Equal(utils.ParseCoin("10000000pool1-1").Amount, afs[0].FarmedPoolCoin.Amount)
	s.Require().True(utils.ParseCoins("9990000000pool1-1,10000000farm1-1").IsEqual(s.getBalances(liquidityProvider1)))

	// SortedByTimeFarmQueue -> [ (l2) 7000000pool1-1]
	// ActiveFarmedTokens -> [ (l1) 10000000pool1-1]
	msg = types.NewMsgFarm(appID1, pool.Id, liquidityProvider2, utils.ParseCoin("7000000pool1-1"))
	err = s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)
	s.Require().True(utils.ParseCoins("9992999999pool1-1,7000000farm1-1").IsEqual(s.getBalances(liquidityProvider2)))

	qf, found := s.keeper.GetQueuedFarmer(s.ctx, appID1, pool.Id, liquidityProvider2)
	afs = s.keeper.GetAllActiveFarmers(s.ctx, appID1, pool.Id)
	s.Require().True(found)
	s.Require().Len(qf.QueudCoins, 1)
	s.Require().Len(afs, 1)

	// SortedByTimeFarmQueue -> [(l1) 9000000pool1-1, (l2) 7000000pool1-1]
	// ActiveFarmedTokens -> [ (l1) 10000000pool1-1]
	msg = types.NewMsgFarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("9000000pool1-1"))
	err = s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)
	s.Require().True(utils.ParseCoins("9981000000pool1-1,19000000farm1-1").IsEqual(s.getBalances(liquidityProvider1)))

	afs = s.keeper.GetAllActiveFarmers(s.ctx, appID1, pool.Id)
	qfs = s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool.Id)
	s.Require().Len(qfs, 2)
	s.Require().Len(afs, 1)

	// Before
	// SortedByTimeFarmQueue -> [(l2) 7000000pool1-1, (l1) 9000000pool1-1]
	// ActiveFarmedTokens -> [ (l1) 10000000pool1-1]
	// After
	// SortedByTimeFarmQueue -> [(l2) 7000000pool1-1, (l1) 6000000pool1-1]
	// ActiveFarmedTokens -> [ (l1) 10000000pool1-1]
	msgUnlock = types.NewMsgUnfarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("3000000farm1-1"))
	err = s.keeper.Unfarm(s.ctx, msgUnlock)
	s.Require().NoError(err)
	s.Require().True(utils.ParseCoins("9984000000pool1-1,16000000farm1-1").IsEqual(s.getBalances(liquidityProvider1)))
	afs = s.keeper.GetAllActiveFarmers(s.ctx, appID1, pool.Id)
	qfs = s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool.Id)
	s.Require().Len(qfs, 2)
	s.Require().Len(afs, 1)
	s.Require().Equal(utils.ParseCoin("7000000pool1-1").Denom, qfs[1].QueudCoins[0].FarmedPoolCoin.Denom)
	s.Require().Equal(utils.ParseCoin("7000000pool1-1").Amount, qfs[1].QueudCoins[0].FarmedPoolCoin.Amount)
	s.Require().Equal(utils.ParseCoin("6000000pool1-1").Denom, qfs[0].QueudCoins[0].FarmedPoolCoin.Denom)
	s.Require().Equal(utils.ParseCoin("6000000pool1-1").Amount, qfs[0].QueudCoins[0].FarmedPoolCoin.Amount)
	s.Require().Equal(utils.ParseCoin("10000000pool1-1").Denom, afs[0].FarmedPoolCoin.Denom)
	s.Require().Equal(utils.ParseCoin("10000000pool1-1").Amount, afs[0].FarmedPoolCoin.Amount)

	// Before
	// SortedByTimeFarmQueue -> [(l2) 7000000pool1-1, (l1) 6000000pool1-1]
	// ActiveFarmedTokens -> [ (l1) 10000000pool1-1]
	// After
	// SortedByTimeFarmQueue -> [(l2) 7000000pool1-1]
	// ActiveFarmedTokens -> [ (l1) 8000000pool1-1]
	msgUnlock = types.NewMsgUnfarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("8000000farm1-1"))
	err = s.keeper.Unfarm(s.ctx, msgUnlock)
	s.Require().NoError(err)
	s.Require().True(utils.ParseCoins("9992000000pool1-1,8000000farm1-1").IsEqual(s.getBalances(liquidityProvider1)))
	afs = s.keeper.GetAllActiveFarmers(s.ctx, appID1, pool.Id)
	qfs = s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool.Id)
	s.Require().Len(qfs, 2)
	s.Require().Len(afs, 1)
	s.Require().Equal(utils.ParseCoin("7000000pool1-1").Denom, qfs[1].QueudCoins[0].FarmedPoolCoin.Denom)
	s.Require().Equal(utils.ParseCoin("7000000pool1-1").Amount, qfs[1].QueudCoins[0].FarmedPoolCoin.Amount)
	s.Require().Equal(utils.ParseCoin("8000000pool1-1").Denom, afs[0].FarmedPoolCoin.Denom)
	s.Require().Equal(utils.ParseCoin("8000000pool1-1").Amount, afs[0].FarmedPoolCoin.Amount)

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(types.DefaultFarmingQueueDuration).Add(time.Minute * 1))
	s.nextBlock()
	// Now
	// SortedByTimeFarmQueue -> []
	// ActiveFarmedTokens -> [ (l1) 8000000pool1-1, (l2) 7000000pool1-1]
	afs = s.keeper.GetAllActiveFarmers(s.ctx, appID1, pool.Id)
	qfs = s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool.Id)
	s.Require().Len(qfs, 2)
	s.Require().Len(qfs[1].QueudCoins, 0)
	s.Require().Len(qfs[0].QueudCoins, 0)
	s.Require().Len(afs, 2)
	s.Require().Equal(utils.ParseCoin("8000000pool1-1").Denom, afs[0].FarmedPoolCoin.Denom)
	s.Require().Equal(utils.ParseCoin("8000000pool1-1").Amount, afs[0].FarmedPoolCoin.Amount)

	s.Require().Equal(utils.ParseCoin("7000000pool1-1").Denom, afs[1].FarmedPoolCoin.Denom)
	s.Require().Equal(utils.ParseCoin("7000000pool1-1").Amount, afs[1].FarmedPoolCoin.Amount)

	// total unlock - lp1
	msgUnlock = types.NewMsgUnfarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("8000000farm1-1"))
	err = s.keeper.Unfarm(s.ctx, msgUnlock)
	s.Require().NoError(err)
	s.Require().True(utils.ParseCoins("10000000000pool1-1").IsEqual(s.getBalances(liquidityProvider1)))

	// total unlock - lp2
	msgUnlock = types.NewMsgUnfarm(appID1, pool.Id, liquidityProvider2, utils.ParseCoin("7000000farm1-1"))
	err = s.keeper.Unfarm(s.ctx, msgUnlock)
	s.Require().NoError(err)
	s.Require().True(utils.ParseCoins("9999999999pool1-1").IsEqual(s.getBalances(liquidityProvider2)))

	// SortedByTimeFarmQueue -> []
	// ActiveFarmedTokens -> [ (l1) 0pool1-1, (l2) 0pool1-1]
	afs = s.keeper.GetAllActiveFarmers(s.ctx, appID1, pool.Id)
	qfs = s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool.Id)
	s.Require().Len(qfs, 2)
	s.Require().Len(qfs[1].QueudCoins, 0)
	s.Require().Len(qfs[0].QueudCoins, 0)
	s.Require().Len(afs, 0)
	// s.Require().Equal(utils.ParseCoin("0pool1-1").Denom, afs[0].FarmedPoolCoin.Denom)
	// s.Require().Equal(utils.ParseCoin("0pool1-1").Amount, afs[0].FarmedPoolCoin.Amount)

	// s.Require().Equal(utils.ParseCoin("0pool1-1").Denom, afs[1].FarmedPoolCoin.Denom)
	// s.Require().Equal(utils.ParseCoin("0pool1-1").Amount, afs[1].FarmedPoolCoin.Amount)

	s.Require().True(utils.ParseCoins("10000000000pool1-1").IsEqual(s.getBalances(liquidityProvider1)))
	s.Require().True(utils.ParseCoins("9999999999pool1-1").IsEqual(s.getBalances(liquidityProvider2)))
}

func (s *KeeperTestSuite) TestFarmAndUnfarm() {
	creator := s.addr(0)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, creator, "1000000000000uasset1,1000000000000uasset2")

	liquidityProvider1 := s.addr(1)
	s.Deposit(appID1, pool.Id, liquidityProvider1, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("10000000000pool1-1").IsEqual(s.getBalances(liquidityProvider1)))

	currentTime := s.ctx.BlockTime()
	s.ctx = s.ctx.WithBlockTime(currentTime)

	msg := types.NewMsgFarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("10000000000pool1-1"))
	err := s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)
	s.Require().True(utils.ParseCoins("10000000000farm1-1").IsEqual(s.getBalances(liquidityProvider1)))
	queuedFarmers := s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool.Id)
	s.Require().Len(queuedFarmers, 1)

	modAddress := s.app.AccountKeeper.GetModuleAddress(types.ModuleName)
	s.Require().True(utils.ParseCoins("10000000000pool1-1").IsEqual(s.getBalances(modAddress)))

	// validate supply of farm coin after farming - mint check
	s.Require().True(utils.ParseCoin("10000000000farm1-1").IsEqual(s.app.BankKeeper.GetSupply(s.ctx, "farm1-1")))

	unfarmMsg := types.NewMsgUnfarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("11000000farm1-1"))
	err = s.keeper.Unfarm(s.ctx, unfarmMsg)
	s.Require().NoError(err)

	// validate supply of farm coin after farming - burn check
	s.Require().True(utils.ParseCoin("9989000000farm1-1").IsEqual(s.app.BankKeeper.GetSupply(s.ctx, "farm1-1")))

	s.Require().True(utils.ParseCoins("9989000000farm1-1,11000000pool1-1").IsEqual(s.getBalances(liquidityProvider1)))
	s.Require().True(utils.ParseCoins("9989000000pool1-1").IsEqual(s.getBalances(modAddress)))

	unfarmMsg = types.NewMsgUnfarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("9989000000farm1-1"))
	err = s.keeper.Unfarm(s.ctx, unfarmMsg)
	s.Require().NoError(err)

	// check after full cycle - 100% farm and unfarm
	s.Require().True(utils.ParseCoin("0farm1-1").IsEqual(s.app.BankKeeper.GetSupply(s.ctx, "farm1-1")))

	s.Require().True(utils.ParseCoins("10000000000pool1-1").IsEqual(s.getBalances(liquidityProvider1)))
	s.Require().True(utils.ParseCoins("").IsEqual(s.getBalances(modAddress)))

}

// liquidity provided in incrementel order
func (s *KeeperTestSuite) TestGetFarmingRewardsDataLinearLPs() {
	currentTime := s.ctx.BlockTime()
	s.ctx = s.ctx.WithBlockTime(currentTime)

	creator := s.addr(0)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, creator, "1000000000000uasset1,1000000000000uasset2")

	liquidityProvider1 := s.addr(1)
	s.Deposit(appID1, pool.Id, liquidityProvider1, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("10000000000pool1-1").IsEqual(s.getBalances(liquidityProvider1)))
	msg := types.NewMsgFarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("10000000000pool1-1"))
	err := s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)

	liquidityProvider2 := s.addr(2)
	s.Deposit(appID1, pool.Id, liquidityProvider2, "2000000000uasset1,2000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("19999999999pool1-1").IsEqual(s.getBalances(liquidityProvider2)))
	msg = types.NewMsgFarm(appID1, pool.Id, liquidityProvider2, utils.ParseCoin("19999999999pool1-1"))
	err = s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)

	liquidityProvider3 := s.addr(3)
	s.Deposit(appID1, pool.Id, liquidityProvider3, "3000000000uasset1,3000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("29999999999pool1-1").IsEqual(s.getBalances(liquidityProvider3)))
	msg = types.NewMsgFarm(appID1, pool.Id, liquidityProvider3, utils.ParseCoin("29999999999pool1-1"))
	err = s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)

	liquidityProvider4 := s.addr(4)
	s.Deposit(appID1, pool.Id, liquidityProvider4, "4000000000uasset1,4000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("39999999999pool1-1").IsEqual(s.getBalances(liquidityProvider4)))
	msg = types.NewMsgFarm(appID1, pool.Id, liquidityProvider4, utils.ParseCoin("39999999999pool1-1"))
	err = s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)

	afs := s.keeper.GetAllActiveFarmers(s.ctx, appID1, pool.Id)
	qfs := s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool.Id)
	s.Require().Len(qfs, 4)
	s.Require().Len(afs, 0)

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(types.DefaultFarmingQueueDuration).Add(time.Minute * 10))
	s.nextBlock()

	afs = s.keeper.GetAllActiveFarmers(s.ctx, appID1, pool.Id)
	qfs = s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool.Id)
	s.Require().Len(qfs, 4)
	s.Require().Len(qfs[0].QueudCoins, 0)
	s.Require().Len(qfs[1].QueudCoins, 0)
	s.Require().Len(qfs[2].QueudCoins, 0)
	s.Require().Len(qfs[3].QueudCoins, 0)
	s.Require().Len(afs, 4)

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

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, creator, "1000000000000uasset1,1000000000000uasset2")

	liquidityProvider1 := s.addr(1)
	s.Deposit(appID1, pool.Id, liquidityProvider1, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("10000000000pool1-1").IsEqual(s.getBalances(liquidityProvider1)))
	msg := types.NewMsgFarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("10000000000pool1-1"))
	err := s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)

	liquidityProvider2 := s.addr(2)
	s.Deposit(appID1, pool.Id, liquidityProvider2, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("9999999999pool1-1").IsEqual(s.getBalances(liquidityProvider2)))
	msg = types.NewMsgFarm(appID1, pool.Id, liquidityProvider2, utils.ParseCoin("9999999999pool1-1"))
	err = s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)

	liquidityProvider3 := s.addr(3)
	s.Deposit(appID1, pool.Id, liquidityProvider3, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("9999999999pool1-1").IsEqual(s.getBalances(liquidityProvider3)))
	msg = types.NewMsgFarm(appID1, pool.Id, liquidityProvider3, utils.ParseCoin("9999999999pool1-1"))
	err = s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)

	liquidityProvider4 := s.addr(4)
	s.Deposit(appID1, pool.Id, liquidityProvider4, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("9999999999pool1-1").IsEqual(s.getBalances(liquidityProvider4)))
	msg = types.NewMsgFarm(appID1, pool.Id, liquidityProvider4, utils.ParseCoin("9999999999pool1-1"))
	err = s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)

	afs := s.keeper.GetAllActiveFarmers(s.ctx, appID1, pool.Id)
	qfs := s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool.Id)
	s.Require().Len(qfs, 4)
	s.Require().Len(afs, 0)

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(types.DefaultFarmingQueueDuration).Add(time.Minute * 10))
	s.nextBlock()

	afs = s.keeper.GetAllActiveFarmers(s.ctx, appID1, pool.Id)
	qfs = s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool.Id)
	s.Require().Len(qfs, 4)
	s.Require().Len(qfs[0].QueudCoins, 0)
	s.Require().Len(qfs[1].QueudCoins, 0)
	s.Require().Len(qfs[2].QueudCoins, 0)
	s.Require().Len(qfs[3].QueudCoins, 0)
	s.Require().Len(afs, 4)

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

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, creator, "1000000000000uasset1,1000000000000uasset2")

	afs := s.keeper.GetAllActiveFarmers(s.ctx, appID1, pool.Id)
	qfs := s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool.Id)
	s.Require().Len(qfs, 0)
	s.Require().Len(afs, 0)

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

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, creator, "1000000000000uasset1,1000000000000uasset2")

	pair2 := s.CreateNewLiquidityPair(appID1, creator, asset2.Denom, asset1.Denom)
	pool2 := s.CreateNewLiquidityPool(appID1, pair2.Id, creator, "1000000000000uasset1,1000000000000uasset2")

	// lp1 - farming only in master pool, not child pool (not eligible for masterpool type reward)
	liquidityProvider1 := s.addr(1)
	s.Deposit(appID1, pool.Id, liquidityProvider1, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("10000000000pool1-1").IsEqual(s.getBalances(liquidityProvider1)))
	msg := types.NewMsgFarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("10000000000pool1-1"))
	err := s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)

	// lp2 - farming in master pool and child pool (eligible for masterpool type reward)
	liquidityProvider2 := s.addr(2)
	s.Deposit(appID1, pool.Id, liquidityProvider2, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Deposit(appID1, pool2.Id, liquidityProvider2, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("9999999999pool1-1,10000000000pool1-2").IsEqual(s.getBalances(liquidityProvider2)))
	msg = types.NewMsgFarm(appID1, pool.Id, liquidityProvider2, utils.ParseCoin("9999999999pool1-1"))
	err = s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)
	msg = types.NewMsgFarm(appID1, pool2.Id, liquidityProvider2, utils.ParseCoin("10000000000pool1-2"))
	err = s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)

	// lp3 - farming only in master pool, not child pool (not eligible for masterpool type reward)
	liquidityProvider3 := s.addr(3)
	s.Deposit(appID1, pool.Id, liquidityProvider3, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("9999999999pool1-1").IsEqual(s.getBalances(liquidityProvider3)))
	msg = types.NewMsgFarm(appID1, pool.Id, liquidityProvider3, utils.ParseCoin("9999999999pool1-1"))
	err = s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)

	// lp4 - farming  in master pool and  child pool (eligible for masterpool type reward)
	liquidityProvider4 := s.addr(4)
	s.Deposit(appID1, pool.Id, liquidityProvider4, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Deposit(appID1, pool2.Id, liquidityProvider4, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("9999999999pool1-1,9999999999pool1-2").IsEqual(s.getBalances(liquidityProvider4)))
	msg = types.NewMsgFarm(appID1, pool.Id, liquidityProvider4, utils.ParseCoin("9999999999pool1-1"))
	err = s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)
	msg = types.NewMsgFarm(appID1, pool2.Id, liquidityProvider4, utils.ParseCoin("9999999999pool1-2"))
	err = s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)

	afs := s.keeper.GetAllActiveFarmers(s.ctx, appID1, pool.Id)
	qfs := s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool.Id)
	s.Require().Len(qfs, 4)
	s.Require().Len(afs, 0)

	afs = s.keeper.GetAllActiveFarmers(s.ctx, appID1, pool2.Id)
	qfs = s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool2.Id)
	s.Require().Len(qfs, 2)
	s.Require().Len(afs, 0)

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(types.DefaultFarmingQueueDuration).Add(time.Minute * 10))
	s.nextBlock()

	afs = s.keeper.GetAllActiveFarmers(s.ctx, appID1, pool.Id)
	qfs = s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool.Id)
	s.Require().Len(qfs, 4)
	s.Require().Len(qfs[0].QueudCoins, 0)
	s.Require().Len(qfs[1].QueudCoins, 0)
	s.Require().Len(qfs[2].QueudCoins, 0)
	s.Require().Len(qfs[3].QueudCoins, 0)
	s.Require().Len(afs, 4)

	afs = s.keeper.GetAllActiveFarmers(s.ctx, appID1, pool2.Id)
	qfs = s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool2.Id)
	s.Require().Len(qfs, 2)
	s.Require().Len(qfs[0].QueudCoins, 0)
	s.Require().Len(qfs[1].QueudCoins, 0)
	s.Require().Len(afs, 2)

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

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, creator, "1000000000000uasset1,1000000000000uasset2")

	// lp1 - farming only in master pool, not child pool (not eligible for masterpool type reward)
	liquidityProvider1 := s.addr(1)
	s.Deposit(appID1, pool.Id, liquidityProvider1, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("10000000000pool1-1").IsEqual(s.getBalances(liquidityProvider1)))
	msg := types.NewMsgFarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("10000000000pool1-1"))
	err := s.keeper.Farm(s.ctx, msg)
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

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, creator, "1000000000000uasset1,1000000000000uasset2")

	liquidityProvider1 := s.addr(1)
	s.Deposit(appID1, pool.Id, liquidityProvider1, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("10000000000pool1-1").IsEqual(s.getBalances(liquidityProvider1)))
	msg := types.NewMsgFarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("10000000000pool1-1"))
	err := s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)

	liquidityProvider2 := s.addr(2)
	s.Deposit(appID1, pool.Id, liquidityProvider2, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("9999999999pool1-1").IsEqual(s.getBalances(liquidityProvider2)))
	msg = types.NewMsgFarm(appID1, pool.Id, liquidityProvider2, utils.ParseCoin("9999999999pool1-1"))
	err = s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)

	liquidityProvider3 := s.addr(3)
	s.Deposit(appID1, pool.Id, liquidityProvider3, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("9999999999pool1-1").IsEqual(s.getBalances(liquidityProvider3)))
	msg = types.NewMsgFarm(appID1, pool.Id, liquidityProvider3, utils.ParseCoin("9999999999pool1-1"))
	err = s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)

	liquidityProvider4 := s.addr(4)
	s.Deposit(appID1, pool.Id, liquidityProvider4, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("9999999999pool1-1").IsEqual(s.getBalances(liquidityProvider4)))
	msg = types.NewMsgFarm(appID1, pool.Id, liquidityProvider4, utils.ParseCoin("9999999999pool1-1"))
	err = s.keeper.Farm(s.ctx, msg)
	s.Require().NoError(err)

	afs := s.keeper.GetAllActiveFarmers(s.ctx, appID1, pool.Id)
	qfs := s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool.Id)
	s.Require().Len(qfs, 4)
	s.Require().Len(afs, 0)

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(types.DefaultFarmingQueueDuration).Add(time.Minute * 10))
	s.nextBlock()

	afs = s.keeper.GetAllActiveFarmers(s.ctx, appID1, pool.Id)
	qfs = s.keeper.GetAllQueuedFarmers(s.ctx, appID1, pool.Id)
	s.Require().Len(qfs, 4)
	s.Require().Len(qfs[0].QueudCoins, 0)
	s.Require().Len(qfs[1].QueudCoins, 0)
	s.Require().Len(qfs[2].QueudCoins, 0)
	s.Require().Len(qfs[3].QueudCoins, 0)
	s.Require().Len(afs, 4)

	msgUnlock := types.NewMsgUnfarm(appID1, pool.Id, liquidityProvider1, utils.ParseCoin("10000000000farm1-1"))
	err = s.keeper.Unfarm(s.ctx, msgUnlock)
	s.Require().NoError(err)
	msgUnlock = types.NewMsgUnfarm(appID1, pool.Id, liquidityProvider2, utils.ParseCoin("9999999999farm1-1"))
	err = s.keeper.Unfarm(s.ctx, msgUnlock)
	s.Require().NoError(err)
	msgUnlock = types.NewMsgUnfarm(appID1, pool.Id, liquidityProvider3, utils.ParseCoin("9999999999farm1-1"))
	err = s.keeper.Unfarm(s.ctx, msgUnlock)
	s.Require().NoError(err)
	msgUnlock = types.NewMsgUnfarm(appID1, pool.Id, liquidityProvider4, utils.ParseCoin("9999999999farm1-1"))
	err = s.keeper.Unfarm(s.ctx, msgUnlock)
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

func (s *KeeperTestSuite) TestGetAmountFarmedForAssetID() {
	creator := s.addr(0)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)
	asset3 := s.CreateNewAsset("ASSETTHREE", "uasset3", 1000000)
	asset4 := s.CreateNewAsset("ASSETFOUR", "uasset4", 1000000)

	pair1 := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset2.Denom)
	pair2 := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset3.Denom)
	pair3 := s.CreateNewLiquidityPair(appID1, creator, asset2.Denom, asset3.Denom)
	pair4 := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset4.Denom)

	pool1 := s.CreateNewLiquidityPool(appID1, pair1.Id, creator, "1000000000000uasset1,1000000000000uasset2")
	pool2 := s.CreateNewLiquidityPool(appID1, pair2.Id, creator, "1000000000000uasset1,1000000000000uasset3")
	pool3 := s.CreateNewLiquidityPool(appID1, pair3.Id, creator, "1000000000000uasset2,1000000000000uasset3")
	pool4 := s.CreateNewLiquidityPool(appID1, pair4.Id, creator, "1000000000000uasset1,1000000000000uasset4")

	liquidityProvider1 := s.addr(1)

	s.Deposit(appID1, pool1.Id, liquidityProvider1, "1000000000uasset1,1000000000uasset2")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("10000000000pool1-1").IsEqual(s.getBalances(liquidityProvider1)))

	s.Deposit(appID1, pool2.Id, liquidityProvider1, "1000000000uasset1,1000000000uasset3")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("10000000000pool1-1,10000000000pool1-2").IsEqual(s.getBalances(liquidityProvider1)))

	s.Deposit(appID1, pool3.Id, liquidityProvider1, "1000000000uasset2,1000000000uasset3")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("10000000000pool1-1,10000000000pool1-2,10000000000pool1-3").IsEqual(s.getBalances(liquidityProvider1)))

	s.Deposit(appID1, pool4.Id, liquidityProvider1, "1000000000uasset1,1000000000uasset4")
	s.nextBlock()
	s.Require().True(utils.ParseCoins("10000000000pool1-1,10000000000pool1-2,10000000000pool1-3,10000000000pool1-4").IsEqual(s.getBalances(liquidityProvider1)))

	currentTime := s.ctx.BlockTime()
	s.ctx = s.ctx.WithBlockTime(currentTime)

	msg1 := types.NewMsgFarm(appID1, pool1.Id, liquidityProvider1, utils.ParseCoin("10000000000pool1-1"))
	msg2 := types.NewMsgFarm(appID1, pool2.Id, liquidityProvider1, utils.ParseCoin("10000000000pool1-2"))
	msg3 := types.NewMsgFarm(appID1, pool3.Id, liquidityProvider1, utils.ParseCoin("10000000000pool1-3"))
	msg4 := types.NewMsgFarm(appID1, pool4.Id, liquidityProvider1, utils.ParseCoin("10000000000pool1-4"))

	err := s.app.LiquidityKeeper.Farm(s.ctx, msg1)
	s.Require().NoError(err)
	err = s.app.LiquidityKeeper.Farm(s.ctx, msg2)
	s.Require().NoError(err)
	err = s.app.LiquidityKeeper.Farm(s.ctx, msg3)
	s.Require().NoError(err)
	err = s.app.LiquidityKeeper.Farm(s.ctx, msg4)
	s.Require().NoError(err)

	queuedFarmers := s.app.LiquidityKeeper.GetAllQueuedFarmers(s.ctx, appID1, pool1.Id)
	s.Require().Len(queuedFarmers, 1)
	queuedFarmers = s.app.LiquidityKeeper.GetAllQueuedFarmers(s.ctx, appID1, pool2.Id)
	s.Require().Len(queuedFarmers, 1)
	queuedFarmers = s.app.LiquidityKeeper.GetAllQueuedFarmers(s.ctx, appID1, pool3.Id)
	s.Require().Len(queuedFarmers, 1)
	queuedFarmers = s.app.LiquidityKeeper.GetAllQueuedFarmers(s.ctx, appID1, pool4.Id)
	s.Require().Len(queuedFarmers, 1)

	activeFarmers := s.app.LiquidityKeeper.GetAllActiveFarmers(s.ctx, appID1, pool1.Id)
	s.Require().Len(activeFarmers, 0)
	activeFarmers = s.app.LiquidityKeeper.GetAllActiveFarmers(s.ctx, appID1, pool2.Id)
	s.Require().Len(activeFarmers, 0)
	activeFarmers = s.app.LiquidityKeeper.GetAllActiveFarmers(s.ctx, appID1, pool3.Id)
	s.Require().Len(activeFarmers, 0)
	activeFarmers = s.app.LiquidityKeeper.GetAllActiveFarmers(s.ctx, appID1, pool4.Id)
	s.Require().Len(activeFarmers, 0)

	s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(time.Hour * 25))
	s.nextBlock()

	activeFarmers = s.app.LiquidityKeeper.GetAllActiveFarmers(s.ctx, appID1, pool1.Id)
	s.Require().Len(activeFarmers, 1)
	activeFarmers = s.app.LiquidityKeeper.GetAllActiveFarmers(s.ctx, appID1, pool2.Id)
	s.Require().Len(activeFarmers, 1)
	activeFarmers = s.app.LiquidityKeeper.GetAllActiveFarmers(s.ctx, appID1, pool3.Id)
	s.Require().Len(activeFarmers, 1)
	activeFarmers = s.app.LiquidityKeeper.GetAllActiveFarmers(s.ctx, appID1, pool4.Id)
	s.Require().Len(activeFarmers, 1)

	activeFarmer, found := s.app.LiquidityKeeper.GetActiveFarmer(s.ctx, appID1, pool1.Id, liquidityProvider1)
	s.Require().True(found)
	s.Require().IsType(types.ActiveFarmer{}, activeFarmer)

	activeFarmer, found = s.app.LiquidityKeeper.GetActiveFarmer(s.ctx, appID1, pool2.Id, liquidityProvider1)
	s.Require().True(found)
	s.Require().IsType(types.ActiveFarmer{}, activeFarmer)

	activeFarmer, found = s.app.LiquidityKeeper.GetActiveFarmer(s.ctx, appID1, pool3.Id, liquidityProvider1)
	s.Require().True(found)
	s.Require().IsType(types.ActiveFarmer{}, activeFarmer)

	activeFarmer, found = s.app.LiquidityKeeper.GetActiveFarmer(s.ctx, appID1, pool4.Id, liquidityProvider1)
	s.Require().True(found)
	s.Require().IsType(types.ActiveFarmer{}, activeFarmer)

	quantityFarmed, err := s.keeper.GetAmountFarmedForAssetID(s.ctx, appID1, asset1.Id, liquidityProvider1)
	s.Require().NoError(err)
	s.Require().Equal(sdk.NewInt(2999999997), quantityFarmed)

	quantityFarmed, err = s.keeper.GetAmountFarmedForAssetID(s.ctx, appID1, asset2.Id, liquidityProvider1)
	s.Require().NoError(err)
	s.Require().Equal(sdk.NewInt(1999999998), quantityFarmed)

	quantityFarmed, err = s.keeper.GetAmountFarmedForAssetID(s.ctx, appID1, asset3.Id, liquidityProvider1)
	s.Require().NoError(err)
	s.Require().Equal(sdk.NewInt(1999999998), quantityFarmed)

	quantityFarmed, err = s.keeper.GetAmountFarmedForAssetID(s.ctx, appID1, asset4.Id, liquidityProvider1)
	s.Require().NoError(err)
	s.Require().Equal(sdk.NewInt(999999999), quantityFarmed)
}

func (s *KeeperTestSuite) TestOraclePricForRewardDistrbution() {
	currentTime := s.ctx.BlockTime()
	s.ctx = s.ctx.WithBlockTime(currentTime)

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 1000000)

	price, found, ass := s.keeper.OraclePrice(s.ctx, asset1.Denom)
	s.Require().Equal(uint64(1000000), price)
	s.Require().True(found)
	s.Require().Equal(asset1, ass)

	price, found, ass = s.keeper.OraclePrice(s.ctx, asset2.Denom)
	s.Require().Equal(uint64(1000000), price)
	s.Require().True(found)
	s.Require().Equal(asset2, ass)

	twa, _ := s.app.MarketKeeper.GetTwa(s.ctx, asset1.Id)
	twa.IsPriceActive = false
	s.app.MarketKeeper.SetTwa(s.ctx, twa)

	price, found, ass = s.keeper.OraclePrice(s.ctx, asset1.Denom)
	s.Require().Equal(uint64(1000000), price)
	s.Require().True(found)
	s.Require().Equal(asset1, ass)

	s.ctx = s.ctx.WithBlockHeight(10)

	twa.IsPriceActive = false
	twa.DiscardedHeightDiff = 10
	s.app.MarketKeeper.SetTwa(s.ctx, twa)

	price, found, ass = s.keeper.OraclePrice(s.ctx, asset1.Denom)
	s.Require().Equal(uint64(1000000), price)
	s.Require().True(found)
	s.Require().Equal(asset1, ass)

	s.ctx = s.ctx.WithBlockHeight(610)

	twa.IsPriceActive = false
	s.app.MarketKeeper.SetTwa(s.ctx, twa)

	price, found, ass = s.keeper.OraclePrice(s.ctx, asset1.Denom)
	s.Require().Equal(uint64(1000000), price)
	s.Require().True(found)

	twa.IsPriceActive = false
	twa.Twa = 0
	s.app.MarketKeeper.SetTwa(s.ctx, twa)

	price, found, _ = s.keeper.OraclePrice(s.ctx, asset1.Denom)
	s.Require().Equal(uint64(0), price)
	s.Require().False(found)
}
