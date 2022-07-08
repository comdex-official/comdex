package keeper_test

import (
	"github.com/comdex-official/comdex/x/liquidity/amm"
	"github.com/comdex-official/comdex/x/liquidity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestCreatePool() {

	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appOne")
	appID2 := s.CreateNewApp("appTwo")

	asset1 := s.CreateNewAsset("ASSET1", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSET2", "uasset2", 2000000)
	asset3 := s.CreateNewAsset("ASSET3", "uasset3", 3000000)

	params, err := s.keeper.GetGenericParams(s.ctx, appID1)
	s.Require().NoError(err)

	s.fundAddr(addr1, params.PairCreationFee)

	app1pairMsg := types.NewMsgCreatePair(appID1, addr1, asset1.Denom, asset2.Denom)
	app1pair, err := s.keeper.CreatePair(s.ctx, app1pairMsg, false)
	s.Require().NoError(err)
	s.Require().IsType(types.Pair{}, app1pair)

	s.fundAddr(addr1, params.PairCreationFee)

	app2pairMsg := types.NewMsgCreatePair(appID2, addr1, asset1.Denom, asset2.Denom)
	app2pair, err := s.keeper.CreatePair(s.ctx, app2pairMsg, false)
	s.Require().NoError(err)
	s.Require().IsType(types.Pair{}, app2pair)

	s.fundAddr(addr1, params.PairCreationFee)

	dummyPair1Msg := types.NewMsgCreatePair(appID1, addr1, asset2.Denom, asset1.Denom)
	dummyPair1, err := s.keeper.CreatePair(s.ctx, dummyPair1Msg, false)
	s.Require().NoError(err)
	s.Require().IsType(types.Pair{}, dummyPair1)

	s.fundAddr(addr1, params.PairCreationFee)

	dummyPair2Msg := types.NewMsgCreatePair(appID1, addr1, asset2.Denom, asset3.Denom)
	dummyPair2, err := s.keeper.CreatePair(s.ctx, dummyPair2Msg, false)
	s.Require().NoError(err)
	s.Require().IsType(types.Pair{}, dummyPair2)
	s.fundAddr(addr1, sdk.NewCoins(sdk.NewCoin(dummyPair2.BaseCoinDenom, sdk.NewInt(1000000000000)), sdk.NewCoin(dummyPair2.QuoteCoinDenom, sdk.NewInt(1000000000000))))

	// s.fundAddr(addr1, params.PairCreationFee)

	// pair2Msg := types.NewMsgCreatePair(appID1, addr1, asset1.Denom, asset2.Denom)
	// pair2, err := s.keeper.CreatePair(s.ctx, pairMsg, false)
	// s.Require().NoError(err)
	// s.Require().IsType(types.Pair{}, pair)

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

}
