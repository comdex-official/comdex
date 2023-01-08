package keeper_test

import (
	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/liquidity"
	"github.com/comdex-official/comdex/x/liquidity/amm"
	"github.com/comdex-official/comdex/x/liquidity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestCreatePool() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	asset3 := s.CreateNewAsset("ASSETHREE", "uasset3", 3000000)

	app1pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	app2pair := s.CreateNewLiquidityPair(appID2, addr1, asset1.Denom, asset2.Denom)
	dummyPair1 := s.CreateNewLiquidityPair(appID1, addr1, asset2.Denom, asset1.Denom)
	dummyPair2 := s.CreateNewLiquidityPair(appID1, addr1, asset2.Denom, asset3.Denom)

	params, err := s.keeper.GetGenericParams(s.ctx, appID1)
	s.Require().NoError(err)
	s.fundAddr(addr1, sdk.NewCoins(sdk.NewCoin(dummyPair2.BaseCoinDenom, sdk.NewInt(1000000000000)), sdk.NewCoin(dummyPair2.QuoteCoinDenom, sdk.NewInt(1000000000000))))

	testCases := []struct {
		Name               string
		Msg                types.MsgCreatePool
		ExpErr             error
		ExpResp            *types.Pool
		QueryResponseIndex uint64
		QueryResponse      *types.Pool
		AvailableBalance   sdk.Coins
	}{
		{
			Name: "error app id invalid",
			Msg: *types.NewMsgCreatePool(
				69, addr1, app1pair.Id, sdk.NewCoins(sdk.NewCoin(app1pair.BaseCoinDenom, sdk.NewInt(1000000000000)), sdk.NewCoin(app1pair.QuoteCoinDenom, sdk.NewInt(1000000000000))),
			),
			ExpErr:             sdkerrors.Wrapf(types.ErrInvalidAppID, "app id %d not found", 69),
			ExpResp:            &types.Pool{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "error pair id invalid",
			Msg: *types.NewMsgCreatePool(
				appID1, addr1, 12, sdk.NewCoins(sdk.NewCoin(app1pair.BaseCoinDenom, sdk.NewInt(1000000000000)), sdk.NewCoin(app1pair.QuoteCoinDenom, sdk.NewInt(1000000000000))),
			),
			ExpErr:             sdkerrors.Wrapf(sdkerrors.ErrNotFound, "pair %d not found", 12),
			ExpResp:            &types.Pool{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "error invalid deposit coin denom 1",
			Msg: *types.NewMsgCreatePool(
				appID1, addr1, app1pair.Id, sdk.NewCoins(sdk.NewCoin("fakedenom1", sdk.NewInt(1000000000000)), sdk.NewCoin("fakedenom2", sdk.NewInt(1000000000000))),
			),
			ExpErr:             sdkerrors.Wrapf(types.ErrInvalidCoinDenom, "coin denom %s is not in the pair", "fakedenom1"),
			ExpResp:            &types.Pool{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "error invalid deposit coin denom 2",
			Msg: *types.NewMsgCreatePool(
				appID1, addr1, app1pair.Id, sdk.NewCoins(sdk.NewCoin(app1pair.BaseCoinDenom, sdk.NewInt(1000000000000)), sdk.NewCoin("fakedenom2", sdk.NewInt(1000000000000))),
			),
			ExpErr:             sdkerrors.Wrapf(types.ErrInvalidCoinDenom, "coin denom %s is not in the pair", "fakedenom2"),
			ExpResp:            &types.Pool{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "error invalid deposit coin denom 3",
			Msg: *types.NewMsgCreatePool(
				appID1, addr1, app1pair.Id, sdk.NewCoins(sdk.NewCoin("fakedenom1", sdk.NewInt(1000000000000)), sdk.NewCoin(app1pair.QuoteCoinDenom, sdk.NewInt(1000000000000))),
			),
			ExpErr:             sdkerrors.Wrapf(types.ErrInvalidCoinDenom, "coin denom %s is not in the pair", "fakedenom1"),
			ExpResp:            &types.Pool{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "error smaller than minimum deposit amount 1",
			Msg: *types.NewMsgCreatePool(
				appID1, addr1, app1pair.Id, sdk.NewCoins(sdk.NewCoin(app1pair.BaseCoinDenom, params.MinInitialDepositAmount.Sub(sdk.NewInt(1))), sdk.NewCoin(app1pair.QuoteCoinDenom, params.MinInitialDepositAmount.Sub(sdk.NewInt(1)))),
			),
			ExpErr:             sdkerrors.Wrapf(types.ErrInsufficientDepositAmount, "%s is smaller than %s", sdk.NewCoin(app1pair.BaseCoinDenom, params.MinInitialDepositAmount.Sub(sdk.NewInt(1))), sdk.NewCoin(app1pair.BaseCoinDenom, params.MinInitialDepositAmount)),
			ExpResp:            &types.Pool{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "error smaller than minimum deposit amount 2",
			Msg: *types.NewMsgCreatePool(
				appID1, addr1, app1pair.Id, sdk.NewCoins(sdk.NewCoin(app1pair.BaseCoinDenom, params.MinInitialDepositAmount), sdk.NewCoin(app1pair.QuoteCoinDenom, params.MinInitialDepositAmount.Sub(sdk.NewInt(1)))),
			),
			ExpErr:             sdkerrors.Wrapf(types.ErrInsufficientDepositAmount, "%s is smaller than %s", sdk.NewCoin(app1pair.QuoteCoinDenom, params.MinInitialDepositAmount.Sub(sdk.NewInt(1))), sdk.NewCoin(app1pair.QuoteCoinDenom, params.MinInitialDepositAmount)),
			ExpResp:            &types.Pool{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		// this case will create a pool even without deposit coins, since testcases run in non atomic
		// environment. This cannot happen in proper envirnoment with actual chain running
		{
			Name: "error insufficient deposit coins",
			Msg: *types.NewMsgCreatePool(
				appID1, addr1, dummyPair1.Id, sdk.NewCoins(sdk.NewCoin(dummyPair1.BaseCoinDenom, sdk.NewInt(1000000000000)), sdk.NewCoin(dummyPair1.QuoteCoinDenom, sdk.NewInt(1000000000000))),
			),
			ExpErr:             sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "0%s is smaller than 1000000000000%s", dummyPair1.QuoteCoinDenom, dummyPair1.QuoteCoinDenom),
			ExpResp:            &types.Pool{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		// this case will create a pool even without deposit coins, since testcases run in non atomic
		// environment. This cannot happen in proper envirnoment with actual chain running
		{
			Name: "error insufficient pool creation fees",
			Msg: *types.NewMsgCreatePool(
				appID1, addr1, dummyPair2.Id, sdk.NewCoins(sdk.NewCoin(dummyPair2.BaseCoinDenom, sdk.NewInt(1000000000000)), sdk.NewCoin(dummyPair2.QuoteCoinDenom, sdk.NewInt(1000000000000))),
			),
			ExpErr:             sdkerrors.Wrap(sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "0%s is smaller than %s", params.PoolCreationFee[0].Denom, params.PoolCreationFee[0].String()), "insufficient pool creation fee"),
			ExpResp:            &types.Pool{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "success valid case app1 pair1",
			Msg: *types.NewMsgCreatePool(
				appID1, addr1, app1pair.Id, sdk.NewCoins(sdk.NewCoin(app1pair.BaseCoinDenom, sdk.NewInt(1000000000000)), sdk.NewCoin(app1pair.QuoteCoinDenom, sdk.NewInt(1000000000000))),
			),
			ExpErr: nil,
			ExpResp: &types.Pool{
				Id:                    3,
				PairId:                1,
				ReserveAddress:        "cosmos1szmdem9z7zwgk8ug0gpcaalj7xdrwhczxfyen7yr025c58wha56svf8f6s",
				PoolCoinDenom:         "pool1-3",
				LastDepositRequestId:  0,
				LastWithdrawRequestId: 0,
				Disabled:              false,
				AppId:                 1,
				Type:                  types.PoolTypeBasic,
				Creator:               addr1.String(),
			},
			QueryResponseIndex: 2, // poolID 1 & 2 are taken by above two cases, since the test environment is non atomic.
			QueryResponse: &types.Pool{
				Id:                    3,
				PairId:                1,
				ReserveAddress:        "cosmos1szmdem9z7zwgk8ug0gpcaalj7xdrwhczxfyen7yr025c58wha56svf8f6s",
				PoolCoinDenom:         "pool1-3",
				LastDepositRequestId:  0,
				LastWithdrawRequestId: 0,
				Disabled:              false,
				AppId:                 1,
				Type:                  types.PoolTypeBasic,
				Creator:               addr1.String(),
			},
			AvailableBalance: sdk.NewCoins(sdk.NewCoin("pool1-3", amm.InitialPoolCoinSupply(sdk.NewInt(1000000000000), sdk.NewInt(1000000000000)))),
		},
		{
			Name: "error pool already exists",
			Msg: *types.NewMsgCreatePool(
				appID1, addr1, app1pair.Id, sdk.NewCoins(sdk.NewCoin(app1pair.BaseCoinDenom, sdk.NewInt(1000000000000)), sdk.NewCoin(app1pair.QuoteCoinDenom, sdk.NewInt(1000000000000))),
			),
			ExpErr:             types.ErrPoolAlreadyExists,
			ExpResp:            &types.Pool{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "success valid case app2 pair1",
			Msg: *types.NewMsgCreatePool(
				appID2, addr1, app2pair.Id, sdk.NewCoins(sdk.NewCoin(app2pair.BaseCoinDenom, sdk.NewInt(1000000000000)), sdk.NewCoin(app2pair.QuoteCoinDenom, sdk.NewInt(1000000000000))),
			),
			ExpErr: nil,
			ExpResp: &types.Pool{
				Id:                    1,
				PairId:                1,
				ReserveAddress:        "cosmos1khz4nd0duzvk4cm3glz3czncnq5ecp77gdh58k558k3wh460rn6qx4e3m0",
				PoolCoinDenom:         "pool2-1",
				LastDepositRequestId:  0,
				LastWithdrawRequestId: 0,
				Disabled:              false,
				AppId:                 2,
				Type:                  types.PoolTypeBasic,
				Creator:               addr1.String(),
			},
			QueryResponseIndex: 0,
			QueryResponse: &types.Pool{
				Id:                    1,
				PairId:                1,
				ReserveAddress:        "cosmos1khz4nd0duzvk4cm3glz3czncnq5ecp77gdh58k558k3wh460rn6qx4e3m0",
				PoolCoinDenom:         "pool2-1",
				LastDepositRequestId:  0,
				LastWithdrawRequestId: 0,
				Disabled:              false,
				AppId:                 2,
				Type:                  types.PoolTypeBasic,
				Creator:               addr1.String(),
			},
			AvailableBalance: sdk.NewCoins(sdk.NewCoin("pool1-3", amm.InitialPoolCoinSupply(sdk.NewInt(1000000000000), sdk.NewInt(1000000000000))), sdk.NewCoin("pool2-1", amm.InitialPoolCoinSupply(sdk.NewInt(1000000000000), sdk.NewInt(1000000000000)))),
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.Name, func() {
			// add funds to acount for valid case
			if tc.ExpErr == nil {
				params, err := s.keeper.GetGenericParams(s.ctx, tc.Msg.AppId)
				s.Require().NoError(err)
				s.fundAddr(sdk.MustAccAddressFromBech32(tc.Msg.Creator), params.PoolCreationFee)
				s.fundAddr(sdk.MustAccAddressFromBech32(tc.Msg.Creator), tc.Msg.DepositCoins)
			}

			resp, err := s.keeper.CreatePool(s.ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				s.Require().Equal(tc.ExpResp, &resp)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				s.Require().Equal(tc.ExpResp, &resp)

				availableBalances := s.getBalances(sdk.MustAccAddressFromBech32(tc.Msg.Creator))
				s.Require().True(tc.AvailableBalance.IsEqual(availableBalances))

				params, err := s.keeper.GetGenericParams(s.ctx, tc.Msg.AppId)
				s.Require().NoError(err)

				expectedPairCreationFeesCollected := sdk.NewCoin(params.PairCreationFee[0].Denom, params.PairCreationFee[0].Amount.Mul(sdk.NewInt(int64(len(s.keeper.GetAllPairs(s.ctx, tc.Msg.AppId))))))

				expectedPoolCreationFeesCollected := sdk.Coin{}
				if tc.Msg.AppId == appID1 {
					expectedPoolCreationFeesCollected = sdk.NewCoin(params.PoolCreationFee[0].Denom, params.PoolCreationFee[0].Amount.Mul(sdk.NewInt(int64(tc.QueryResponseIndex-1))))
				} else {
					expectedPoolCreationFeesCollected = sdk.NewCoin(params.PoolCreationFee[0].Denom, params.PoolCreationFee[0].Amount.Mul(sdk.NewInt(int64(tc.QueryResponseIndex+1))))
				}

				collectedPairPoolCreationFee := s.getBalances(sdk.MustAccAddressFromBech32(params.FeeCollectorAddress))
				s.Require().True(sdk.NewCoins(expectedPairCreationFeesCollected.Add(expectedPoolCreationFeesCollected)).IsEqual(collectedPairPoolCreationFee))

				pools := s.keeper.GetAllPools(s.ctx, tc.Msg.AppId)
				s.Require().Len(pools, int(tc.QueryResponseIndex)+1)
				s.Require().Equal(tc.QueryResponse.Id, pools[tc.QueryResponseIndex].Id)
				s.Require().Equal(tc.QueryResponse.PairId, pools[tc.QueryResponseIndex].PairId)
				s.Require().Equal(tc.QueryResponse.ReserveAddress, pools[tc.QueryResponseIndex].ReserveAddress)
				s.Require().Equal(tc.QueryResponse.PoolCoinDenom, pools[tc.QueryResponseIndex].PoolCoinDenom)
				s.Require().Equal(tc.QueryResponse.LastDepositRequestId, pools[tc.QueryResponseIndex].LastDepositRequestId)
				s.Require().Equal(tc.QueryResponse.LastWithdrawRequestId, pools[tc.QueryResponseIndex].LastWithdrawRequestId)
				s.Require().Equal(tc.QueryResponse.Disabled, pools[tc.QueryResponseIndex].Disabled)
				s.Require().Equal(tc.QueryResponse.AppId, pools[tc.QueryResponseIndex].AppId)
			}
		})
	}
	gauges := s.app.Rewardskeeper.GetAllGauges(s.ctx)
	s.Require().Len(gauges, 2)
}

func (s *KeeperTestSuite) TestDisabledPool() {
	// A disabled pool is:
	// 1. A pool with at least one side of its x/y coin's balance is 0.
	// 2. A pool with 0 pool coin supply(all investors has withdrawn their coins)

	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pair2 := s.CreateNewLiquidityPair(appID1, addr1, asset2.Denom, asset1.Denom)

	// Create a pool.
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2")
	// Send the pool's balances to somewhere else.
	s.sendCoins(pool.GetReserveAddress(), s.addr(2), s.getBalances(pool.GetReserveAddress()))

	// By now, the pool is not marked as disabled automatically.
	// When someone sends a deposit/withdraw request to the pool or
	// the pool tries to participate in matching, then the pool
	// is marked as disabled.
	pool, _ = s.keeper.GetPool(s.ctx, appID1, pool.Id)
	s.Require().False(pool.Disabled)

	// A depositor tries to deposit to the pool.
	s.Deposit(appID1, pool.Id, addr1, "10000000uasset1,10000000uasset2")
	s.nextBlock()

	// Now, the pool is disabled.
	pool, _ = s.keeper.GetPool(s.ctx, appID1, pool.Id)
	s.Require().True(pool.Disabled)

	// Here's the second example.
	// This time, the pool creator withdraws all his coins.
	pool = s.CreateNewLiquidityPool(appID1, pair2.Id, addr1, "1000000uasset1,1000000uasset2")
	s.Withdraw(appID1, pool.Id, addr1, s.getBalance(addr1, pool.PoolCoinDenom))
	s.nextBlock()

	// The pool is disabled again.
	pool, _ = s.keeper.GetPool(s.ctx, appID1, pool.Id)
	s.Require().True(pool.Disabled)
}

func (s *KeeperTestSuite) TestCreatePoolAfterDisabled() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)

	// Create a disabled pool.
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2")
	s.Withdraw(appID1, pool.Id, addr1, s.getBalance(addr1, pool.PoolCoinDenom))
	s.nextBlock()

	// Now a new pool can be created with same denom pair because
	// all pools with same denom pair are disabled.
	s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2")
}

func (s *KeeperTestSuite) TestPoolIndexes() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")
	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2")

	fetchedPool, found := s.keeper.GetPoolByReserveAddress(s.ctx, appID1, pool.GetReserveAddress())
	s.Require().True(found)
	s.Require().Equal(pool.Id, fetchedPool.Id)

	pools := s.keeper.GetPoolsByPair(s.ctx, appID1, pair.Id)
	s.Require().Len(pools, 1)
	s.Require().Equal(pool.Id, pools[0].Id)
}

func (s *KeeperTestSuite) TestDeposit() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)

	app1Pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	app1Pool := s.CreateNewLiquidityPool(appID1, app1Pair.Id, addr1, "1000000000000uasset1,1000000000000uasset2")

	app2Pair := s.CreateNewLiquidityPair(appID2, addr1, asset1.Denom, asset2.Denom)
	app2Pool := s.CreateNewLiquidityPool(appID2, app2Pair.Id, addr1, "1000000000000uasset1,1000000000000uasset2")

	addr1AvailableBalance := sdk.NewCoins(
		sdk.NewCoin(app1Pool.PoolCoinDenom, amm.InitialPoolCoinSupply(sdk.NewInt(1000000000000), sdk.NewInt(1000000000000))),
		sdk.NewCoin(app2Pool.PoolCoinDenom, amm.InitialPoolCoinSupply(sdk.NewInt(1000000000000), sdk.NewInt(1000000000000))),
	)

	testCases := []struct {
		Name               string
		Msg                types.MsgDeposit
		ExpErr             error
		ExpResp            *types.DepositRequest
		QueryResponseIndex uint64
		QueryResponse      *types.DepositRequest
		AvailableBalance   sdk.Coins
	}{
		{
			Name: "error app id invalid",
			Msg: *types.NewMsgDeposit(
				69, addr1, app1Pool.Id, sdk.NewCoins(sdk.NewCoin(app1Pair.BaseCoinDenom, sdk.NewInt(100000000)), sdk.NewCoin(app1Pair.QuoteCoinDenom, sdk.NewInt(100000000))),
			),
			ExpErr:             sdkerrors.Wrapf(types.ErrInvalidAppID, "app id %d not found", 69),
			ExpResp:            &types.DepositRequest{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "error pool id invalid",
			Msg: *types.NewMsgDeposit(
				appID1, addr1, 69, sdk.NewCoins(sdk.NewCoin(app1Pair.BaseCoinDenom, sdk.NewInt(100000000)), sdk.NewCoin(app1Pair.QuoteCoinDenom, sdk.NewInt(100000000))),
			),
			ExpErr:             sdkerrors.Wrapf(sdkerrors.ErrNotFound, "pool %d not found", 69),
			ExpResp:            &types.DepositRequest{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "error invalid deposit coin 1",
			Msg: *types.NewMsgDeposit(
				appID1, addr1, app1Pool.Id, sdk.NewCoins(sdk.NewCoin("fakedenom1", sdk.NewInt(100000000)), sdk.NewCoin("fakedenom2", sdk.NewInt(100000000))),
			),
			ExpErr:             sdkerrors.Wrapf(types.ErrInvalidCoinDenom, "coin denom fakedenom1 is not in the pair"),
			ExpResp:            &types.DepositRequest{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "error invalid deposit coin 2",
			Msg: *types.NewMsgDeposit(
				appID1, addr1, app1Pool.Id, sdk.NewCoins(sdk.NewCoin(app1Pair.BaseCoinDenom, sdk.NewInt(100000000)), sdk.NewCoin("fakedenom2", sdk.NewInt(100000000))),
			),
			ExpErr:             sdkerrors.Wrapf(types.ErrInvalidCoinDenom, "coin denom fakedenom2 is not in the pair"),
			ExpResp:            &types.DepositRequest{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "error invalid deposit coin 3",
			Msg: *types.NewMsgDeposit(
				appID1, addr1, app1Pool.Id, sdk.NewCoins(sdk.NewCoin("fakedenom1", sdk.NewInt(100000000)), sdk.NewCoin(app1Pair.QuoteCoinDenom, sdk.NewInt(100000000))),
			),
			ExpErr:             sdkerrors.Wrapf(types.ErrInvalidCoinDenom, "coin denom fakedenom1 is not in the pair"),
			ExpResp:            &types.DepositRequest{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "error insufficeint deposit coins",
			Msg: *types.NewMsgDeposit(
				appID1, addr1, app1Pool.Id, sdk.NewCoins(sdk.NewCoin(app1Pair.BaseCoinDenom, sdk.NewInt(100000000)), sdk.NewCoin(app1Pair.QuoteCoinDenom, sdk.NewInt(100000000))),
			),
			ExpErr:             sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "0%s is smaller than 100000000%s", app1Pair.BaseCoinDenom, app1Pair.BaseCoinDenom),
			ExpResp:            &types.DepositRequest{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "success valid case 1 app1 pool1",
			Msg: *types.NewMsgDeposit(
				appID1, addr1, app1Pool.Id, sdk.NewCoins(sdk.NewCoin(app1Pair.BaseCoinDenom, sdk.NewInt(100000000)), sdk.NewCoin(app1Pair.QuoteCoinDenom, sdk.NewInt(100000000))),
			),
			ExpErr: nil,
			ExpResp: &types.DepositRequest{
				Id:             1,
				PoolId:         1,
				MsgHeight:      0,
				Depositor:      addr1.String(),
				DepositCoins:   sdk.NewCoins(sdk.NewCoin(app1Pair.BaseCoinDenom, sdk.NewInt(100000000)), sdk.NewCoin(app1Pair.QuoteCoinDenom, sdk.NewInt(100000000))),
				AcceptedCoins:  nil,
				MintedPoolCoin: sdk.NewCoin(app1Pool.PoolCoinDenom, sdk.NewInt(0)),
				Status:         types.RequestStatusNotExecuted,
				AppId:          appID1,
			},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   addr1AvailableBalance.Add(sdk.NewCoin(app1Pool.PoolCoinDenom, sdk.NewInt(1000000000))),
		},
		{
			Name: "success valid case 2 app1 pool1",
			Msg: *types.NewMsgDeposit(
				appID1, addr1, app1Pool.Id, sdk.NewCoins(sdk.NewCoin(app1Pair.BaseCoinDenom, sdk.NewInt(300000000)), sdk.NewCoin(app1Pair.QuoteCoinDenom, sdk.NewInt(300000000))),
			),
			ExpErr: nil,
			ExpResp: &types.DepositRequest{
				Id:             2,
				PoolId:         1,
				MsgHeight:      0,
				Depositor:      addr1.String(),
				DepositCoins:   sdk.NewCoins(sdk.NewCoin(app1Pair.BaseCoinDenom, sdk.NewInt(300000000)), sdk.NewCoin(app1Pair.QuoteCoinDenom, sdk.NewInt(300000000))),
				AcceptedCoins:  nil,
				MintedPoolCoin: sdk.NewCoin(app1Pool.PoolCoinDenom, sdk.NewInt(0)),
				Status:         types.RequestStatusNotExecuted,
				AppId:          appID1,
			},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   addr1AvailableBalance.Add(sdk.NewCoin(app1Pool.PoolCoinDenom, sdk.NewInt(3999999999))),
		},
		{
			Name: "success valid case 3 app2 pool1",
			Msg: *types.NewMsgDeposit(
				appID2, addr1, app2Pool.Id, sdk.NewCoins(sdk.NewCoin(app2Pair.BaseCoinDenom, sdk.NewInt(100000000)), sdk.NewCoin(app2Pair.QuoteCoinDenom, sdk.NewInt(100000000))),
			),
			ExpErr: nil,
			ExpResp: &types.DepositRequest{
				Id:             1,
				PoolId:         1,
				MsgHeight:      0,
				Depositor:      addr1.String(),
				DepositCoins:   sdk.NewCoins(sdk.NewCoin(app2Pair.BaseCoinDenom, sdk.NewInt(100000000)), sdk.NewCoin(app2Pair.QuoteCoinDenom, sdk.NewInt(100000000))),
				AcceptedCoins:  nil,
				MintedPoolCoin: sdk.NewCoin(app2Pool.PoolCoinDenom, sdk.NewInt(0)),
				Status:         types.RequestStatusNotExecuted,
				AppId:          appID2,
			},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   addr1AvailableBalance.Add(sdk.NewCoin(app1Pool.PoolCoinDenom, sdk.NewInt(3999999999))).Add(sdk.NewCoin(app2Pool.PoolCoinDenom, sdk.NewInt(1000000000))),
		},
		{
			Name: "success valid case 4 app2 pool1",
			Msg: *types.NewMsgDeposit(
				appID2, addr1, app2Pool.Id, sdk.NewCoins(sdk.NewCoin(app2Pair.BaseCoinDenom, sdk.NewInt(700000000)), sdk.NewCoin(app2Pair.QuoteCoinDenom, sdk.NewInt(700000000))),
			),
			ExpErr: nil,
			ExpResp: &types.DepositRequest{
				Id:             2,
				PoolId:         1,
				MsgHeight:      0,
				Depositor:      addr1.String(),
				DepositCoins:   sdk.NewCoins(sdk.NewCoin(app2Pair.BaseCoinDenom, sdk.NewInt(700000000)), sdk.NewCoin(app2Pair.QuoteCoinDenom, sdk.NewInt(700000000))),
				AcceptedCoins:  nil,
				MintedPoolCoin: sdk.NewCoin(app2Pool.PoolCoinDenom, sdk.NewInt(0)),
				Status:         types.RequestStatusNotExecuted,
				AppId:          appID2,
			},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   addr1AvailableBalance.Add(sdk.NewCoin(app1Pool.PoolCoinDenom, sdk.NewInt(3999999999))).Add(sdk.NewCoin(app2Pool.PoolCoinDenom, sdk.NewInt(7999999999))),
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.Name, func() {
			// add funds to acount for valid case
			if tc.ExpErr == nil {
				s.fundAddr(sdk.MustAccAddressFromBech32(tc.Msg.Depositor), tc.Msg.DepositCoins)
			}

			depositReq, err := s.keeper.Deposit(s.ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				s.Require().Equal(tc.ExpResp, &depositReq)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(depositReq)
				s.Require().Equal(tc.ExpResp, &depositReq)

				err := s.keeper.ExecuteDepositRequest(s.ctx, depositReq)
				s.Require().NoError(err)

				availableBalances := s.getBalances(sdk.MustAccAddressFromBech32(tc.Msg.Depositor))
				s.Require().True(tc.AvailableBalance.IsEqual(availableBalances))
				resp, found := s.keeper.GetDepositRequest(s.ctx, tc.Msg.AppId, tc.Msg.PoolId, depositReq.Id)
				s.Require().True(found)
				s.Require().Equal(tc.ExpResp.Id, resp.Id)
				s.Require().Equal(tc.ExpResp.PoolId, resp.PoolId)
				s.Require().Equal(tc.ExpResp.MsgHeight, resp.MsgHeight)
				s.Require().Equal(tc.ExpResp.Depositor, resp.Depositor)
				s.Require().Equal(tc.ExpResp.DepositCoins, resp.DepositCoins)
				s.Require().Equal(tc.ExpResp.DepositCoins, resp.AcceptedCoins)
				s.Require().Equal(types.RequestStatusSucceeded, resp.Status)
				s.Require().Equal(tc.ExpResp.AppId, resp.AppId)
			}
		})
	}
}

func (s *KeeperTestSuite) TestDepositRefund() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000uasset1,1500000uasset2")

	req := s.Deposit(appID1, pool.Id, addr1, "20000uasset1,15000uasset2")
	liquidity.EndBlocker(s.ctx, s.keeper, s.app.AssetKeeper)
	req, _ = s.keeper.GetDepositRequest(s.ctx, appID1, req.PoolId, req.Id)
	s.Require().Equal(types.RequestStatusSucceeded, req.Status)

	s.Require().True(utils.ParseCoin("10000uasset1").IsEqual(s.getBalance(addr1, "uasset1")))
	s.Require().True(utils.ParseCoin("0uasset2").IsEqual(s.getBalance(addr1, "uasset2")))
	liquidity.BeginBlocker(s.ctx, s.keeper, s.app.AssetKeeper)

	pair = s.CreateNewLiquidityPair(appID1, addr1, asset2.Denom, asset1.Denom)
	pool = s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000uasset2,1000000000000000uasset1")

	req = s.Deposit(appID1, pool.Id, addr2, "1uasset1,1uasset2")
	liquidity.EndBlocker(s.ctx, s.keeper, s.app.AssetKeeper)
	req, _ = s.keeper.GetDepositRequest(s.ctx, appID1, req.PoolId, req.Id)
	s.Require().Equal(types.RequestStatusFailed, req.Status)
	s.Require().True(req.DepositCoins.IsEqual(s.getBalances(addr2)))
}

func (s *KeeperTestSuite) TestDepositRefundTooSmallMintedPoolCoin() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)
	addr3 := s.addr(3)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000uasset1,1500000uasset2")

	req := s.Deposit(appID1, pool.Id, addr2, "20000uasset1,15000uasset2")
	liquidity.EndBlocker(s.ctx, s.keeper, s.app.AssetKeeper)
	req, _ = s.keeper.GetDepositRequest(s.ctx, appID1, req.PoolId, req.Id)
	s.Require().Equal(types.RequestStatusSucceeded, req.Status)

	s.Require().True(utils.ParseCoin("10000uasset1").IsEqual(s.getBalance(addr2, "uasset1")))
	s.Require().True(utils.ParseCoin("0uasset2").IsEqual(s.getBalance(addr2, "uasset2")))
	liquidity.BeginBlocker(s.ctx, s.keeper, s.app.AssetKeeper)

	pair = s.CreateNewLiquidityPair(appID1, addr1, asset2.Denom, asset1.Denom)
	pool = s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000uasset2,1000000000000000uasset1")

	req = s.Deposit(appID1, pool.Id, addr3, "1uasset1,1uasset2")
	liquidity.EndBlocker(s.ctx, s.keeper, s.app.AssetKeeper)
	req, _ = s.keeper.GetDepositRequest(s.ctx, appID1, req.PoolId, req.Id)
	s.Require().Equal(types.RequestStatusFailed, req.Status)

	s.Require().True(req.DepositCoins.IsEqual(s.getBalances(addr3)))
}

func (s *KeeperTestSuite) TestDepositToDisabledPool() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)
	addr3 := s.addr(3)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)

	// Create a disabled pool by sending the pool's balances to somewhere else.
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2")
	s.sendCoins(pool.GetReserveAddress(), addr2, s.getBalances(pool.GetReserveAddress()))

	// The depositor deposits coins but this will fail because the pool
	// is treated as disabled.
	req := s.Deposit(appID1, pool.Id, addr3, "1000000uasset1,1000000uasset2")
	err := s.keeper.ExecuteDepositRequest(s.ctx, req)
	s.Require().NoError(err)
	req, _ = s.keeper.GetDepositRequest(s.ctx, appID1, pool.Id, req.Id)
	s.Require().Equal(types.RequestStatusFailed, req.Status)

	// Delete the previous request and refund coins to the depositor.
	liquidity.BeginBlocker(s.ctx, s.keeper, s.app.AssetKeeper)

	// Now any deposits will result in an error.
	_, err = s.keeper.Deposit(s.ctx, types.NewMsgDeposit(appID1, addr2, pool.Id, req.DepositCoins))
	s.Require().ErrorIs(err, types.ErrDisabledPool)
}

func (s *KeeperTestSuite) TestTooLargePool() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2")

	_, err := s.keeper.Deposit(s.ctx, types.NewMsgDeposit(appID1, addr1, pool.Id, utils.ParseCoins("10000000000000000000000000000000000000000uasset1,10000000000000000000000000000000000000000uasset2")))
	s.Require().ErrorIs(err, types.ErrTooLargePool)
}

func (s *KeeperTestSuite) TestWithdraw() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000uasset1,1000000000uasset2")

	availablePoolBalance := s.getBalance(addr1, pool.PoolCoinDenom)

	testCases := []struct {
		Name             string
		Msg              types.MsgWithdraw
		ExpErr           error
		ExpResp          *types.WithdrawRequest
		AvailableBalance sdk.Coins
	}{
		{
			Name: "error app id invalid",
			Msg: *types.NewMsgWithdraw(
				69, addr1, pool.Id, availablePoolBalance,
			),
			ExpErr:           sdkerrors.Wrapf(types.ErrInvalidAppID, "app id %d not found", 69),
			ExpResp:          &types.WithdrawRequest{},
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "error pool id invalid",
			Msg: *types.NewMsgWithdraw(
				appID1, addr1, 69, availablePoolBalance,
			),
			ExpErr:           sdkerrors.Wrapf(sdkerrors.ErrNotFound, "pool %d not found", 69),
			ExpResp:          &types.WithdrawRequest{},
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "error insufficeint pool coins",
			Msg: *types.NewMsgWithdraw(
				appID1, addr1, pool.Id, availablePoolBalance.Add(sdk.NewCoin(availablePoolBalance.Denom, availablePoolBalance.Amount.Add(newInt(1000)))),
			),
			ExpErr:           sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "%s is smaller than %s", availablePoolBalance.String(), availablePoolBalance.Add(sdk.NewCoin(availablePoolBalance.Denom, availablePoolBalance.Amount.Add(newInt(1000)))).String()),
			ExpResp:          &types.WithdrawRequest{},
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "error invalid pool coin denom",
			Msg: *types.NewMsgWithdraw(
				appID1, addr1, pool.Id, utils.ParseCoin("1000000pool1"),
			),
			ExpErr:           types.ErrWrongPoolCoinDenom,
			ExpResp:          &types.WithdrawRequest{},
			AvailableBalance: sdk.NewCoins(),
		},
		{
			Name: "success valid case",
			Msg: *types.NewMsgWithdraw(
				appID1, addr1, pool.Id, availablePoolBalance,
			),
			ExpErr: nil,
			ExpResp: &types.WithdrawRequest{
				Id:             1,
				PoolId:         1,
				MsgHeight:      0,
				Withdrawer:     addr1.String(),
				PoolCoin:       availablePoolBalance,
				WithdrawnCoins: nil,
				Status:         types.RequestStatusNotExecuted,
				AppId:          1,
			},
			AvailableBalance: sdk.NewCoins(),
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.Name, func() {
			withdrawReq, err := s.keeper.Withdraw(s.ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				s.Require().Equal(tc.ExpResp, &withdrawReq)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(withdrawReq)
				s.Require().Equal(tc.ExpResp, &withdrawReq)

				availableBalances := s.getBalances(sdk.MustAccAddressFromBech32(tc.Msg.Withdrawer))
				s.Require().True(tc.AvailableBalance.IsEqual(availableBalances))

				depositedCoins := s.getBalances(pool.GetReserveAddress())
				s.nextBlock()
				s.Require().True(depositedCoins.IsEqual(s.getBalances(sdk.MustAccAddressFromBech32(tc.Msg.Withdrawer))))
			}
		})
	}
}

func (s *KeeperTestSuite) TestWithdrawRefund() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)
	addr3 := s.addr(3)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000uasset1,1000000000uasset2")

	s.Deposit(appID1, pool.Id, addr2, "1000000uasset1,1000000uasset2")
	s.nextBlock()

	// Make the pool depleted.
	s.sendCoins(pool.GetReserveAddress(), addr3, s.getBalances(pool.GetReserveAddress()))

	poolCoin := s.getBalance(addr2, pool.PoolCoinDenom)
	s.Withdraw(appID1, pool.Id, addr2, poolCoin)
	s.nextBlock()

	s.Require().True(sdk.NewCoins(poolCoin).IsEqual(s.getBalances(addr2)))
}

func (s *KeeperTestSuite) TestWithdrawRefundTooSmallWithdrawCoins() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000uasset1,1000000000uasset2")

	s.Deposit(appID1, pool.Id, addr2, "1000000uasset1,1000000uasset2")
	s.nextBlock()
	poolCoin := s.getBalance(addr2, pool.PoolCoinDenom)

	// Withdrawing too small amount of pool coin.
	s.Withdraw(appID1, pool.Id, addr2, utils.ParseCoin("100pool1-1"))
	s.nextBlock()

	s.Require().True(sdk.NewCoins(poolCoin).IsEqual(s.getBalances(addr2)))
}

func (s *KeeperTestSuite) TestWithdrawFromDisabledPool() {
	addr1 := s.addr(1)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)

	// Create a disabled pool by sending the pool's balances to somewhere else.
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2")
	s.sendCoins(pool.GetReserveAddress(), addr2, s.getBalances(pool.GetReserveAddress()))

	// The pool creator tries to withdraw his coins, but this will fail.
	req := s.Withdraw(appID1, pool.Id, addr1, s.getBalance(addr1, pool.PoolCoinDenom))
	err := s.keeper.ExecuteWithdrawRequest(s.ctx, req)
	s.Require().NoError(err)
	req, _ = s.keeper.GetWithdrawRequest(s.ctx, appID1, pool.Id, req.Id)
	s.Require().Equal(types.RequestStatusFailed, req.Status)

	// Delete the previous request and refund coins to the withdrawer.
	liquidity.BeginBlocker(s.ctx, s.keeper, s.app.AssetKeeper)

	// Now any withdrawals will result in an error.
	_, err = s.keeper.Withdraw(s.ctx, types.NewMsgWithdraw(appID1, addr1, pool.Id, s.getBalance(addr1, pool.PoolCoinDenom)))
	s.Require().ErrorIs(err, types.ErrDisabledPool)
}

func (s *KeeperTestSuite) TestGetDepositRequestsByDepositor() {
	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, s.addr(0), "1000000denom1,1000000denom2")
	req1 := s.Deposit(appID1, pool.Id, s.addr(1), "1000000denom1,1000000denom2")
	req2 := s.Deposit(appID1, pool.Id, s.addr(1), "1000000denom1,1000000denom2")
	reqs := s.keeper.GetDepositRequestsByDepositor(s.ctx, appID1, s.addr(1))
	s.Require().Len(reqs, 2)
	s.Require().Equal(req1.PoolId, reqs[0].PoolId)
	s.Require().Equal(req1.Id, reqs[0].Id)
	s.Require().Equal(req2.PoolId, reqs[1].PoolId)
	s.Require().Equal(req2.Id, reqs[1].Id)
}

func (s *KeeperTestSuite) TestWithdrawRequestsByWithdrawer() {
	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityPool(appID1, pair.Id, s.addr(0), "1000000denom1,1000000denom2")
	s.Deposit(appID1, pool.Id, s.addr(1), "1000000denom1,1000000denom2")
	s.nextBlock()
	req1 := s.Withdraw(appID1, pool.Id, s.addr(1), utils.ParseCoin("10000pool1-1"))
	req2 := s.Withdraw(appID1, pool.Id, s.addr(1), utils.ParseCoin("10000pool1-1"))
	reqs := s.keeper.GetWithdrawRequestsByWithdrawer(s.ctx, appID1, s.addr(1))
	s.Require().Len(reqs, 2)
	s.Require().Equal(req1.PoolId, reqs[0].PoolId)
	s.Require().Equal(req1.Id, reqs[0].Id)
	s.Require().Equal(req2.PoolId, reqs[1].PoolId)
	s.Require().Equal(req2.Id, reqs[1].Id)
}

func (s *KeeperTestSuite) TestCreateRangedPool() {
	addr1 := s.addr(1)
	dummyAddr := s.addr(696969)

	appID1 := s.CreateNewApp("appone")
	appID2 := s.CreateNewApp("apptwo")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)
	asset3 := s.CreateNewAsset("ASSETHREE", "uasset3", 3000000)

	app1pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	app2pair := s.CreateNewLiquidityPair(appID2, addr1, asset1.Denom, asset2.Denom)
	dummyPair1 := s.CreateNewLiquidityPair(appID1, addr1, asset2.Denom, asset1.Denom)
	dummyPair2 := s.CreateNewLiquidityPair(appID1, addr1, asset2.Denom, asset3.Denom)

	params, err := s.keeper.GetGenericParams(s.ctx, appID1)
	s.Require().NoError(err)

	testCases := []struct {
		Name               string
		Msg                types.MsgCreateRangedPool
		ExpErr             error
		ExpResp            *types.Pool
		QueryResponseIndex uint64
		QueryResponse      *types.Pool
		AvailableBalance   sdk.Coins
	}{
		{
			Name: "error app id invalid",
			Msg: *types.NewMsgCreateRangedPool(
				69, addr1, app1pair.Id,
				sdk.NewCoins(sdk.NewCoin(app1pair.BaseCoinDenom, sdk.NewInt(1000000000000)), sdk.NewCoin(app1pair.QuoteCoinDenom, sdk.NewInt(1000000000000))),
				sdk.MustNewDecFromStr("0.99"), sdk.MustNewDecFromStr("1.01"), sdk.MustNewDecFromStr("1"),
			),
			ExpErr:             sdkerrors.Wrapf(types.ErrInvalidAppID, "app id %d not found", 69),
			ExpResp:            &types.Pool{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "error pair id invalid",
			Msg: *types.NewMsgCreateRangedPool(
				appID1, addr1, 12,
				sdk.NewCoins(sdk.NewCoin(app1pair.BaseCoinDenom, sdk.NewInt(1000000000000)), sdk.NewCoin(app1pair.QuoteCoinDenom, sdk.NewInt(1000000000000))),
				sdk.MustNewDecFromStr("0.99"), sdk.MustNewDecFromStr("1.01"), sdk.MustNewDecFromStr("1"),
			),
			ExpErr:             sdkerrors.Wrapf(sdkerrors.ErrNotFound, "pair %d not found", 12),
			ExpResp:            &types.Pool{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "error invalid deposit coin denom 1",
			Msg: *types.NewMsgCreateRangedPool(
				appID1, addr1, app1pair.Id,
				sdk.NewCoins(sdk.NewCoin("fakedenom1", sdk.NewInt(1000000000000)), sdk.NewCoin("fakedenom2", sdk.NewInt(1000000000000))),
				sdk.MustNewDecFromStr("0.99"), sdk.MustNewDecFromStr("1.01"), sdk.MustNewDecFromStr("1"),
			),
			ExpErr:             sdkerrors.Wrapf(types.ErrInvalidCoinDenom, "coin denom %s is not in the pair", "fakedenom1"),
			ExpResp:            &types.Pool{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "error invalid deposit coin denom 2",
			Msg: *types.NewMsgCreateRangedPool(
				appID1, addr1, app1pair.Id, sdk.NewCoins(sdk.NewCoin(app1pair.BaseCoinDenom, sdk.NewInt(1000000000000)), sdk.NewCoin("fakedenom2", sdk.NewInt(1000000000000))),
				sdk.MustNewDecFromStr("0.99"), sdk.MustNewDecFromStr("1.01"), sdk.MustNewDecFromStr("1"),
			),
			ExpErr:             sdkerrors.Wrapf(types.ErrInvalidCoinDenom, "coin denom %s is not in the pair", "fakedenom2"),
			ExpResp:            &types.Pool{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "error invalid deposit coin denom 3",
			Msg: *types.NewMsgCreateRangedPool(
				appID1, addr1, app1pair.Id, sdk.NewCoins(sdk.NewCoin("fakedenom1", sdk.NewInt(1000000000000)), sdk.NewCoin(app1pair.QuoteCoinDenom, sdk.NewInt(1000000000000))),
				sdk.MustNewDecFromStr("0.99"), sdk.MustNewDecFromStr("1.01"), sdk.MustNewDecFromStr("1"),
			),
			ExpErr:             sdkerrors.Wrapf(types.ErrInvalidCoinDenom, "coin denom %s is not in the pair", "fakedenom1"),
			ExpResp:            &types.Pool{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "error smaller than minimum deposit amount 1",
			Msg: *types.NewMsgCreateRangedPool(
				appID1, addr1, app1pair.Id, sdk.NewCoins(sdk.NewCoin(app1pair.BaseCoinDenom, params.MinInitialDepositAmount.Sub(sdk.NewInt(1))), sdk.NewCoin(app1pair.QuoteCoinDenom, params.MinInitialDepositAmount.Sub(sdk.NewInt(1)))),
				sdk.MustNewDecFromStr("0.99"), sdk.MustNewDecFromStr("1.01"), sdk.MustNewDecFromStr("1"),
			),
			ExpErr:             types.ErrInsufficientDepositAmount,
			ExpResp:            &types.Pool{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "error smaller than minimum deposit amount 2",
			Msg: *types.NewMsgCreateRangedPool(
				appID1, addr1, app1pair.Id, sdk.NewCoins(sdk.NewCoin(app1pair.BaseCoinDenom, params.MinInitialDepositAmount), sdk.NewCoin(app1pair.QuoteCoinDenom, params.MinInitialDepositAmount.Sub(sdk.NewInt(1)))),
				sdk.MustNewDecFromStr("0.99"), sdk.MustNewDecFromStr("1.01"), sdk.MustNewDecFromStr("1"),
			),
			ExpErr:             types.ErrInsufficientDepositAmount,
			ExpResp:            &types.Pool{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "error initial price lower than min price",
			Msg: *types.NewMsgCreateRangedPool(
				appID1, addr1, app1pair.Id, sdk.NewCoins(sdk.NewCoin(app1pair.BaseCoinDenom, params.MinInitialDepositAmount), sdk.NewCoin(app1pair.QuoteCoinDenom, params.MinInitialDepositAmount.Sub(sdk.NewInt(1)))),
				sdk.MustNewDecFromStr("0.99"), sdk.MustNewDecFromStr("1.01"), sdk.MustNewDecFromStr("0.98"),
			),
			ExpErr:             sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "initial price must not be lower than min price"),
			ExpResp:            &types.Pool{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "error initial price higher than max price",
			Msg: *types.NewMsgCreateRangedPool(
				appID1, addr1, app1pair.Id, sdk.NewCoins(sdk.NewCoin(app1pair.BaseCoinDenom, params.MinInitialDepositAmount), sdk.NewCoin(app1pair.QuoteCoinDenom, params.MinInitialDepositAmount.Sub(sdk.NewInt(1)))),
				sdk.MustNewDecFromStr("0.99"), sdk.MustNewDecFromStr("1.01"), sdk.MustNewDecFromStr("1.05"),
			),
			ExpErr:             sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "initial price must not be higher than max price"),
			ExpResp:            &types.Pool{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		// // this case will create a pool even without deposit coins, since testcases run in non atomic
		// // environment. This cannot happen in proper envirnoment with actual chain running
		{
			Name: "error insufficient deposit coins",
			Msg: *types.NewMsgCreateRangedPool(
				appID1, addr1, dummyPair1.Id, sdk.NewCoins(sdk.NewCoin(dummyPair1.BaseCoinDenom, sdk.NewInt(1000000000000)), sdk.NewCoin(dummyPair1.QuoteCoinDenom, sdk.NewInt(1000000000000))),
				sdk.MustNewDecFromStr("0.99"), sdk.MustNewDecFromStr("1.01"), sdk.MustNewDecFromStr("1"),
			),
			ExpErr:             sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "0%s is smaller than 1000000000000%s", dummyPair1.QuoteCoinDenom, dummyPair1.QuoteCoinDenom),
			ExpResp:            &types.Pool{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		// this case will create a pool even without deposit coins, since testcases run in non atomic
		// environment. This cannot happen in proper envirnoment with actual chain running
		{
			Name: "error insufficient pool creation fees",
			Msg: *types.NewMsgCreateRangedPool(
				appID1, addr1, dummyPair2.Id, sdk.NewCoins(sdk.NewCoin(dummyPair2.BaseCoinDenom, sdk.NewInt(1000000000000)), sdk.NewCoin(dummyPair2.QuoteCoinDenom, sdk.NewInt(1000000000000))),
				sdk.MustNewDecFromStr("0.99"), sdk.MustNewDecFromStr("1.01"), sdk.MustNewDecFromStr("1"),
			),
			ExpErr:             sdkerrors.Wrap(sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "0%s is smaller than %s", params.PoolCreationFee[0].Denom, params.PoolCreationFee[0].String()), "insufficient pool creation fee"),
			ExpResp:            &types.Pool{},
			QueryResponseIndex: 0,
			QueryResponse:      nil,
			AvailableBalance:   sdk.NewCoins(),
		},
		{
			Name: "success valid case app1 pair1 pool1",
			Msg: *types.NewMsgCreateRangedPool(
				appID1, addr1, app1pair.Id, sdk.NewCoins(sdk.NewCoin(app1pair.BaseCoinDenom, sdk.NewInt(1000000000000)), sdk.NewCoin(app1pair.QuoteCoinDenom, sdk.NewInt(1000000000000))),
				sdk.MustNewDecFromStr("0.99"), sdk.MustNewDecFromStr("1.01"), sdk.MustNewDecFromStr("1"),
			),
			ExpErr: nil,
			ExpResp: &types.Pool{
				Id:                    3,
				PairId:                1,
				ReserveAddress:        "cosmos1szmdem9z7zwgk8ug0gpcaalj7xdrwhczxfyen7yr025c58wha56svf8f6s",
				PoolCoinDenom:         "pool1-3",
				LastDepositRequestId:  0,
				LastWithdrawRequestId: 0,
				Disabled:              false,
				AppId:                 1,
				Type:                  types.PoolTypeRanged,
				Creator:               addr1.String(),
			},
			QueryResponseIndex: 2, // poolID 1 & 2 are taken by above two cases, since the test environment is non atomic.
			QueryResponse: &types.Pool{
				Id:                    3,
				PairId:                1,
				ReserveAddress:        "cosmos1szmdem9z7zwgk8ug0gpcaalj7xdrwhczxfyen7yr025c58wha56svf8f6s",
				PoolCoinDenom:         "pool1-3",
				LastDepositRequestId:  0,
				LastWithdrawRequestId: 0,
				Disabled:              false,
				AppId:                 1,
				Type:                  types.PoolTypeBasic,
				Creator:               addr1.String(),
			},
			AvailableBalance: sdk.NewCoins(
				sdk.NewCoin("pool1-3", amm.InitialPoolCoinSupply(sdk.NewInt(1000000000000), sdk.NewInt(1000000000000))),
				sdk.NewCoin(app1pair.BaseCoinDenom, sdk.NewInt(9925681617)),
			),
		},
		{
			Name: "success valid case app2 pair1",
			Msg: *types.NewMsgCreateRangedPool(
				appID2, addr1, app2pair.Id,
				sdk.NewCoins(sdk.NewCoin(app2pair.BaseCoinDenom, sdk.NewInt(1000000000000)), sdk.NewCoin(app2pair.QuoteCoinDenom, sdk.NewInt(1000000000000))),
				sdk.MustNewDecFromStr("0.99"), sdk.MustNewDecFromStr("1.01"), sdk.MustNewDecFromStr("1"),
			),
			ExpErr: nil,
			ExpResp: &types.Pool{
				Id:                    1,
				PairId:                1,
				ReserveAddress:        "cosmos1khz4nd0duzvk4cm3glz3czncnq5ecp77gdh58k558k3wh460rn6qx4e3m0",
				PoolCoinDenom:         "pool2-1",
				LastDepositRequestId:  0,
				LastWithdrawRequestId: 0,
				Disabled:              false,
				AppId:                 2,
				Type:                  types.PoolTypeRanged,
				Creator:               addr1.String(),
			},
			QueryResponseIndex: 0,
			QueryResponse: &types.Pool{
				Id:                    1,
				PairId:                1,
				ReserveAddress:        "cosmos1khz4nd0duzvk4cm3glz3czncnq5ecp77gdh58k558k3wh460rn6qx4e3m0",
				PoolCoinDenom:         "pool2-1",
				LastDepositRequestId:  0,
				LastWithdrawRequestId: 0,
				Disabled:              false,
				AppId:                 2,
				Type:                  types.PoolTypeRanged,
				Creator:               addr1.String(),
			},
			AvailableBalance: sdk.NewCoins(
				sdk.NewCoin("pool2-1", amm.InitialPoolCoinSupply(sdk.NewInt(1000000000000), sdk.NewInt(1000000000000))),
				sdk.NewCoin(app1pair.BaseCoinDenom, sdk.NewInt(9925681617)),
			),
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.Name, func() {
			s.sendCoins(sdk.MustAccAddressFromBech32(tc.Msg.Creator), dummyAddr, s.getBalances(addr1))

			if tc.Name == "error insufficient pool creation fees" {
				s.fundAddr(sdk.MustAccAddressFromBech32(tc.Msg.Creator), tc.Msg.DepositCoins)
			}
			// add funds to acount for valid case
			if tc.ExpErr == nil {
				params, err := s.keeper.GetGenericParams(s.ctx, tc.Msg.AppId)
				s.Require().NoError(err)
				s.fundAddr(sdk.MustAccAddressFromBech32(tc.Msg.Creator), params.PoolCreationFee)
				s.fundAddr(sdk.MustAccAddressFromBech32(tc.Msg.Creator), tc.Msg.DepositCoins)
			}

			resp, err := s.keeper.CreateRangedPool(s.ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				s.Require().Equal(tc.ExpResp, &resp)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)
				tc.ExpResp.MinPrice = &tc.Msg.MinPrice
				tc.ExpResp.MaxPrice = &tc.Msg.MaxPrice
				s.Require().Equal(tc.ExpResp, &resp)

				availableBalances := s.getBalances(sdk.MustAccAddressFromBech32(tc.Msg.Creator))
				s.Require().True(tc.AvailableBalance.IsEqual(availableBalances))

				params, err := s.keeper.GetGenericParams(s.ctx, tc.Msg.AppId)
				s.Require().NoError(err)

				expectedPairCreationFeesCollected := sdk.NewCoin(params.PairCreationFee[0].Denom, params.PairCreationFee[0].Amount.Mul(sdk.NewInt(int64(len(s.keeper.GetAllPairs(s.ctx, tc.Msg.AppId))))))

				expectedPoolCreationFeesCollected := sdk.Coin{}
				if tc.Msg.AppId == appID1 {
					expectedPoolCreationFeesCollected = sdk.NewCoin(params.PoolCreationFee[0].Denom, params.PoolCreationFee[0].Amount.Mul(sdk.NewInt(int64(tc.QueryResponseIndex-1))))
				} else {
					expectedPoolCreationFeesCollected = sdk.NewCoin(params.PoolCreationFee[0].Denom, params.PoolCreationFee[0].Amount.Mul(sdk.NewInt(int64(tc.QueryResponseIndex+1))))
				}

				collectedPairPoolCreationFee := s.getBalances(sdk.MustAccAddressFromBech32(params.FeeCollectorAddress))
				s.Require().True(sdk.NewCoins(expectedPairCreationFeesCollected.Add(expectedPoolCreationFeesCollected)).IsEqual(collectedPairPoolCreationFee))
				pools := s.keeper.GetAllPools(s.ctx, tc.Msg.AppId)
				s.Require().Len(pools, int(tc.QueryResponseIndex)+1)
				s.Require().Equal(tc.QueryResponse.Id, pools[tc.QueryResponseIndex].Id)
				s.Require().Equal(tc.QueryResponse.PairId, pools[tc.QueryResponseIndex].PairId)
				s.Require().Equal(tc.QueryResponse.ReserveAddress, pools[tc.QueryResponseIndex].ReserveAddress)
				s.Require().Equal(tc.QueryResponse.PoolCoinDenom, pools[tc.QueryResponseIndex].PoolCoinDenom)
				s.Require().Equal(tc.QueryResponse.LastDepositRequestId, pools[tc.QueryResponseIndex].LastDepositRequestId)
				s.Require().Equal(tc.QueryResponse.LastWithdrawRequestId, pools[tc.QueryResponseIndex].LastWithdrawRequestId)
				s.Require().Equal(tc.QueryResponse.Disabled, pools[tc.QueryResponseIndex].Disabled)
				s.Require().Equal(tc.QueryResponse.AppId, pools[tc.QueryResponseIndex].AppId)
				s.Require().Equal(&tc.Msg.MinPrice, pools[tc.QueryResponseIndex].MinPrice)
				s.Require().Equal(&tc.Msg.MaxPrice, pools[tc.QueryResponseIndex].MaxPrice)
			}
		})
	}
	gauges := s.app.Rewardskeeper.GetAllGauges(s.ctx)
	s.Require().Len(gauges, 2)
}

func (s *KeeperTestSuite) TestMultipleBasicPool() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")
	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	_ = s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2")

	params, err := s.keeper.GetGenericParams(s.ctx, appID1)
	s.Require().NoError(err)

	parsedDepositCoins := utils.ParseCoins("1000000uasset1,1000000uasset2")

	s.fundAddr(addr1, params.PoolCreationFee)
	s.fundAddr(addr1, parsedDepositCoins)
	msg := types.NewMsgCreatePool(appID1, addr1, pair.Id, parsedDepositCoins)
	_, err = s.keeper.CreatePool(s.ctx, msg)
	s.Require().EqualError(types.ErrPoolAlreadyExists, err.Error())
}

func (s *KeeperTestSuite) Test1MaximumRangePoolInPair() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")
	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	// maximum 20 pool can be created in each pair
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))

	params, err := s.keeper.GetGenericParams(s.ctx, appID1)
	s.Require().NoError(err)

	parsedDepositCoins := utils.ParseCoins("1000000uasset1,1000000uasset2")

	s.fundAddr(addr1, params.PoolCreationFee)
	s.fundAddr(addr1, parsedDepositCoins)
	msg := types.NewMsgCreateRangedPool(appID1, addr1, pair.Id, parsedDepositCoins, sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_, err = s.keeper.CreateRangedPool(s.ctx, msg)
	s.Require().EqualError(types.ErrTooManyPools, err.Error())
}

func (s *KeeperTestSuite) Test2MaximumRangePoolInPair() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")
	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	// maximum 20 pool can be created in each pair
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))

	_ = s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2")

	params, err := s.keeper.GetGenericParams(s.ctx, appID1)
	s.Require().NoError(err)

	parsedDepositCoins := utils.ParseCoins("1000000uasset1,1000000uasset2")

	s.fundAddr(addr1, params.PoolCreationFee)
	s.fundAddr(addr1, parsedDepositCoins)
	msg := types.NewMsgCreateRangedPool(appID1, addr1, pair.Id, parsedDepositCoins, sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_, err = s.keeper.CreateRangedPool(s.ctx, msg)
	s.Require().EqualError(types.ErrTooManyPools, err.Error())
}

func (s *KeeperTestSuite) Test3MaximumRangePoolInPair() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")
	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	// maximum 20 pool can be created in each pair
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000000uasset1,1000000uasset2", sdk.MustNewDecFromStr("0.95"), sdk.MustNewDecFromStr("1.05"), sdk.MustNewDecFromStr("1"))

	params, err := s.keeper.GetGenericParams(s.ctx, appID1)
	s.Require().NoError(err)

	parsedDepositCoins := utils.ParseCoins("1000000uasset1,1000000uasset2")

	s.fundAddr(addr1, params.PoolCreationFee)
	s.fundAddr(addr1, parsedDepositCoins)
	msg := types.NewMsgCreatePool(appID1, addr1, pair.Id, parsedDepositCoins)
	_, err = s.keeper.CreatePool(s.ctx, msg)
	s.Require().EqualError(types.ErrTooManyPools, err.Error())
}

func (s *KeeperTestSuite) TestRangedPoolDepositWithdraw() {
	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityRangedPool(appID1, pair.Id, s.addr(1), "1000000denom1,1000000denom2", utils.ParseDec("0.5"), utils.ParseDec("2.0"), utils.ParseDec("1.0"))

	rx, ry := s.keeper.GetPoolBalances(s.ctx, pool)
	ammPool := pool.AMMPool(rx.Amount, ry.Amount, sdk.Int{})
	s.Require().True(utils.DecApproxEqual(ammPool.Price(), utils.ParseDec("1.0")))

	s.Deposit(appID1, pool.Id, s.addr(2), "400000denom1,1000000denom2")
	liquidity.EndBlocker(s.ctx, s.keeper, s.app.AssetKeeper)
	rx, ry = s.keeper.GetPoolBalances(s.ctx, pool)
	ammPool = pool.AMMPool(rx.Amount, ry.Amount, sdk.Int{})
	s.Require().True(utils.DecApproxEqual(ammPool.Price(), utils.ParseDec("1.0")))

	poolCoin := s.getBalance(s.addr(2), pool.PoolCoinDenom)
	s.Withdraw(appID1, pool.Id, s.addr(2), poolCoin.SubAmount(poolCoin.Amount.QuoRaw(3))) // withdraw 2/3 pool coin
	liquidity.EndBlocker(s.ctx, s.keeper, s.app.AssetKeeper)
	rx, ry = s.keeper.GetPoolBalances(s.ctx, pool)
	ammPool = pool.AMMPool(rx.Amount, ry.Amount, sdk.Int{})
	s.Require().True(utils.DecApproxEqual(ammPool.Price(), utils.ParseDec("1.0")))
}

func (s *KeeperTestSuite) TestRangedPoolDepositWithdraw_single_side() {
	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityRangedPool(appID1, pair.Id, s.addr(1), "1000000denom1", utils.ParseDec("0.5"), utils.ParseDec("2.0"), utils.ParseDec("0.5"))

	rx, ry := s.keeper.GetPoolBalances(s.ctx, pool)
	s.Require().True(intEq(sdk.ZeroInt(), rx.Amount))
	s.Require().True(intEq(sdk.NewInt(1000000), ry.Amount))
	ps := s.keeper.GetPoolCoinSupply(s.ctx, pool)

	s.Deposit(appID1, pool.Id, s.addr(2), "50000denom1")
	s.nextBlock()

	pc := s.getBalance(s.addr(2), pool.PoolCoinDenom)

	rx, ry = s.keeper.GetPoolBalances(s.ctx, pool)
	s.Require().True(intEq(sdk.ZeroInt(), rx.Amount))
	s.Require().True(intEq(sdk.NewInt(1050000), ry.Amount))
	s.Require().True(intEq(ps.QuoRaw(20), pc.Amount))

	balanceBefore := s.getBalance(s.addr(2), "denom1")
	s.Withdraw(appID1, pool.Id, s.addr(2), sdk.NewCoin(pool.PoolCoinDenom, pc.Amount))
	s.nextBlock()
	balanceAfter := s.getBalance(s.addr(2), "denom1")

	s.Require().True(balanceAfter.Sub(balanceBefore).Amount.Sub(sdk.NewInt(50000)).LTE(sdk.OneInt()))

	s.Deposit(appID1, pool.Id, s.addr(3), "1000000denom1,1000000denom2")
	s.nextBlock()

	s.Require().True(intEq(sdk.ZeroInt(), s.getBalance(s.addr(3), "denom1").Amount))
	s.Require().True(intEq(sdk.NewInt(1000000), s.getBalance(s.addr(3), "denom2").Amount))
}

func (s *KeeperTestSuite) TestRangedPoolDepositWithdraw_single_side2() {
	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)
	pool := s.CreateNewLiquidityRangedPool(appID1, pair.Id, s.addr(1), "1000000denom2", utils.ParseDec("0.5"), utils.ParseDec("2.0"), utils.ParseDec("2.0"))

	rx, ry := s.keeper.GetPoolBalances(s.ctx, pool)
	s.Require().True(intEq(sdk.NewInt(1000000), rx.Amount))
	s.Require().True(intEq(sdk.ZeroInt(), ry.Amount))
	ps := s.keeper.GetPoolCoinSupply(s.ctx, pool)

	s.Deposit(appID1, pool.Id, s.addr(2), "50000denom2")
	s.nextBlock()

	pc := s.getBalance(s.addr(2), pool.PoolCoinDenom)

	rx, ry = s.keeper.GetPoolBalances(s.ctx, pool)
	s.Require().True(intEq(sdk.NewInt(1050000), rx.Amount))
	s.Require().True(intEq(sdk.ZeroInt(), ry.Amount))
	s.Require().True(intEq(ps.QuoRaw(20), pc.Amount))

	balanceBefore := s.getBalance(s.addr(2), "denom2")
	s.Withdraw(appID1, pool.Id, s.addr(2), sdk.NewCoin(pool.PoolCoinDenom, pc.Amount))
	s.nextBlock()
	balanceAfter := s.getBalance(s.addr(2), "denom2")

	s.Require().True(balanceAfter.Sub(balanceBefore).Amount.Sub(sdk.NewInt(50000)).LTE(sdk.OneInt()))

	s.Deposit(appID1, pool.Id, s.addr(3), "1000000denom1,1000000denom2")
	s.nextBlock()

	s.Require().True(intEq(sdk.ZeroInt(), s.getBalance(s.addr(3), "denom2").Amount))
	s.Require().True(intEq(sdk.NewInt(1000000), s.getBalance(s.addr(3), "denom1").Amount))
}
