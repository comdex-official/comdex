package keeper_test

import (
	"time"

	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/liquidity"
	"github.com/comdex-official/comdex/x/liquidity/amm"
	"github.com/comdex-official/comdex/x/liquidity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestLimitOrder() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appOne")

	asset1 := s.CreateNewAsset("ASSET1", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSET2", "uasset2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)

	params, err := s.keeper.GetGenericParams(s.ctx, appID1)
	s.Require().NoError(err)

	testCases := []struct {
		Name    string
		Msg     types.MsgLimitOrder
		ExpErr  error
		ExpResp *types.Order
	}{
		{
			Name: "error invalid app id",
			Msg: *types.NewMsgLimitOrder(
				69,
				addr1,
				pair.Id,
				types.OrderDirectionBuy,
				utils.ParseCoin("10030000uasset2"),
				asset1.Denom,
				newDec(1),
				newInt(10000000),
				time.Second*10,
			),
			ExpErr:  sdkerrors.Wrap(sdkerrors.Wrapf(types.ErrInvalidAppID, "app id %d not found", 69), "params retreval failed"),
			ExpResp: &types.Order{},
		},
		{
			Name: "error max order life span",
			Msg: *types.NewMsgLimitOrder(
				appID1,
				addr1,
				pair.Id,
				types.OrderDirectionBuy,
				utils.ParseCoin("10030000uasset2"),
				asset1.Denom,
				newDec(1),
				newInt(10000000),
				time.Hour*48,
			),
			ExpErr:  sdkerrors.Wrapf(types.ErrTooLongOrderLifespan, "%s is longer than %s", time.Hour*48, params.MaxOrderLifespan),
			ExpResp: &types.Order{},
		},
		{
			Name: "error invalid pair id",
			Msg: *types.NewMsgLimitOrder(
				appID1,
				addr1,
				69,
				types.OrderDirectionBuy,
				utils.ParseCoin("10030000uasset2"),
				asset1.Denom,
				newDec(1),
				newInt(10000000),
				time.Second*10,
			),
			ExpErr:  sdkerrors.Wrapf(sdkerrors.ErrNotFound, "pair %d not found", 69),
			ExpResp: &types.Order{},
		},
		{
			Name: "error price higher than upper limit",
			Msg: *types.NewMsgLimitOrder(
				appID1,
				addr1,
				pair.Id,
				types.OrderDirectionBuy,
				utils.ParseCoin("10030000uasset2"),
				asset1.Denom,
				amm.HighestTick(int(params.TickPrecision+1)),
				newInt(10000000),
				time.Second*10,
			),
			ExpErr:  sdkerrors.Wrapf(types.ErrPriceOutOfRange, "%s is higher than %s", amm.HighestTick(int(params.TickPrecision+1)), amm.HighestTick(int(params.TickPrecision))),
			ExpResp: &types.Order{},
		},
		{
			Name: "error price lower than lower limit",
			Msg: *types.NewMsgLimitOrder(
				appID1,
				addr1,
				pair.Id,
				types.OrderDirectionBuy,
				utils.ParseCoin("10030000uasset2"),
				asset1.Denom,
				amm.LowestTick(int(params.TickPrecision-1)),
				newInt(10000000),
				time.Second*10,
			),
			ExpErr:  sdkerrors.Wrapf(types.ErrPriceOutOfRange, "%s is lower than %s", amm.LowestTick(int(params.TickPrecision-1)), amm.LowestTick(int(params.TickPrecision))),
			ExpResp: &types.Order{},
		},
		{
			Name: "error invalid denom pair buy direction",
			Msg: *types.NewMsgLimitOrder(
				appID1,
				addr1,
				pair.Id,
				types.OrderDirectionBuy,
				utils.ParseCoin("10030000uasset1"),
				asset2.Denom,
				newDec(1),
				newInt(10000000),
				time.Second*10,
			),
			ExpErr:  sdkerrors.Wrapf(types.ErrWrongPair, "denom pair (%s, %s) != (%s, %s)", asset2.Denom, asset1.Denom, pair.BaseCoinDenom, pair.QuoteCoinDenom),
			ExpResp: &types.Order{},
		},
		{
			Name: "error invalid denom pair sell direction",
			Msg: *types.NewMsgLimitOrder(
				appID1,
				addr1,
				pair.Id,
				types.OrderDirectionSell,
				utils.ParseCoin("10030000uasset2"),
				asset1.Denom,
				newDec(1),
				newInt(10000000),
				time.Second*10,
			),
			ExpErr:  sdkerrors.Wrapf(types.ErrWrongPair, "denom pair (%s, %s) != (%s, %s)", asset2.Denom, asset1.Denom, pair.BaseCoinDenom, pair.QuoteCoinDenom),
			ExpResp: &types.Order{},
		},
		{
			Name: "error insufficient offer coin buy direction",
			Msg: *types.NewMsgLimitOrder(
				appID1,
				addr1,
				pair.Id,
				types.OrderDirectionBuy,
				utils.ParseCoin("10000000uasset2"), // swap fee excluded
				asset1.Denom,
				newDec(1),
				newInt(10000000),
				time.Second*10,
			),
			ExpErr:  sdkerrors.Wrapf(types.ErrInsufficientOfferCoin, "10000000uasset2 is smaller than 10030000uasset2"),
			ExpResp: &types.Order{},
		},
		{
			Name: "error insufficient offer coin sell direction",
			Msg: *types.NewMsgLimitOrder(
				appID1,
				addr1,
				pair.Id,
				types.OrderDirectionSell,
				utils.ParseCoin("10000000uasset1"), // swap fee excluded
				asset2.Denom,
				newDec(1),
				newInt(10000000),
				time.Second*10,
			),
			ExpErr:  sdkerrors.Wrapf(types.ErrInsufficientOfferCoin, "10000000uasset1 is smaller than 10030000uasset1"),
			ExpResp: &types.Order{},
		},
		{
			Name: "error too small order buy direction",
			Msg: *types.NewMsgLimitOrder(
				appID1,
				addr1,
				pair.Id,
				types.OrderDirectionBuy,
				utils.ParseCoin("130uasset2"),
				asset1.Denom,
				newDec(1),
				newInt(99),
				time.Second*10,
			),
			ExpErr:  types.ErrTooSmallOrder,
			ExpResp: &types.Order{},
		},
		{
			Name: "error too small order sell direction",
			Msg: *types.NewMsgLimitOrder(
				appID1,
				addr1,
				pair.Id,
				types.OrderDirectionSell,
				utils.ParseCoin("130uasset1"),
				asset2.Denom,
				newDec(1),
				newInt(99),
				time.Second*10,
			),
			ExpErr:  types.ErrTooSmallOrder,
			ExpResp: &types.Order{},
		},
		{
			Name: "error insufficient funds",
			Msg: *types.NewMsgLimitOrder(
				appID1,
				addr1,
				pair.Id,
				types.OrderDirectionBuy,
				utils.ParseCoin("1003000uasset2"),
				asset1.Denom,
				newDec(1),
				newInt(1000000),
				time.Second*10,
			),
			ExpErr:  sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "0uasset2 is smaller than 1003000uasset2"),
			ExpResp: &types.Order{},
		},
		{
			Name: "success valid case buy direction",
			Msg: *types.NewMsgLimitOrder(
				appID1,
				addr1,
				pair.Id,
				types.OrderDirectionBuy,
				utils.ParseCoin("1003000uasset2"),
				asset1.Denom,
				newDec(1),
				newInt(1000000),
				time.Second*10,
			),
			ExpErr: nil,
			ExpResp: &types.Order{
				Id:                 1,
				PairId:             1,
				MsgHeight:          0,
				Orderer:            addr1.String(),
				Direction:          types.OrderDirectionBuy,
				OfferCoin:          utils.ParseCoin("1000000uasset2"),
				RemainingOfferCoin: utils.ParseCoin("1000000uasset2"),
				ReceivedCoin:       utils.ParseCoin("0uasset1"),
				Price:              newDec(1),
				Amount:             newInt(1000000),
				OpenAmount:         newInt(1000000),
				BatchId:            1,
				ExpireAt:           s.ctx.BlockTime().Add(time.Second * 10),
				Status:             types.OrderStatusNotExecuted,
				AppId:              appID1,
			},
		},
		{
			Name: "success valid case sell direction",
			Msg: *types.NewMsgLimitOrder(
				appID1,
				addr1,
				pair.Id,
				types.OrderDirectionSell,
				utils.ParseCoin("1003000uasset1"),
				asset2.Denom,
				newDec(1),
				newInt(1000000),
				time.Second*10,
			),
			ExpErr: nil,
			ExpResp: &types.Order{
				Id:                 2,
				PairId:             1,
				MsgHeight:          0,
				Orderer:            addr1.String(),
				Direction:          types.OrderDirectionSell,
				OfferCoin:          utils.ParseCoin("1000000uasset1"),
				RemainingOfferCoin: utils.ParseCoin("1000000uasset1"),
				ReceivedCoin:       utils.ParseCoin("0uasset2"),
				Price:              newDec(1),
				Amount:             newInt(1000000),
				OpenAmount:         newInt(1000000),
				BatchId:            1,
				ExpireAt:           s.ctx.BlockTime().Add(time.Second * 10),
				Status:             types.OrderStatusNotExecuted,
				AppId:              appID1,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.Name, func() {
			if tc.ExpErr == nil {
				s.fundAddr(tc.Msg.GetOrderer(), sdk.NewCoins(tc.Msg.OfferCoin))
			}
			order, err := s.keeper.LimitOrder(s.ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				s.Require().Equal(tc.ExpResp, &order)
			} else {
				s.Require().NoError(err)
				s.Require().IsType(tc.ExpResp, &order)
				s.Require().Equal(tc.ExpResp, &order)

				qorder, found := s.keeper.GetOrder(s.ctx, tc.Msg.AppId, tc.Msg.PairId, order.Id)
				s.Require().True(found)
				s.Require().Equal(qorder, order)
			}
		})
	}
}

func (s *KeeperTestSuite) TestLimitOrderRefund() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appOne")

	asset1 := s.CreateNewAsset("ASSET1", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSET2", "denom2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)

	s.fundAddr(addr1, utils.ParseCoins("1000000000denom1,1000000000denom2"))

	testCases := []struct {
		Name         string
		Msg          *types.MsgLimitOrder
		RefundedCoin sdk.Coin
	}{
		{
			Name: "refund tc 1",
			Msg: types.NewMsgLimitOrder(
				appID1, addr1, pair.Id, types.OrderDirectionBuy, utils.ParseCoin("1003000denom2"), "denom1",
				utils.ParseDec("1.0"), newInt(1000000), 0),
			RefundedCoin: utils.ParseCoin("0denom2"),
		},
		{
			Name: "refund tc 2",
			Msg: types.NewMsgLimitOrder(
				appID1, addr1, pair.Id, types.OrderDirectionBuy, utils.ParseCoin("1000000denom2"), "denom1",
				utils.ParseDec("1.0"), newInt(10000), 0,
			),
			RefundedCoin: utils.ParseCoin("989970denom2"),
		},
		{
			Name: "refund tc 3",
			Msg: types.NewMsgLimitOrder(
				appID1, addr1, pair.Id, types.OrderDirectionBuy, utils.ParseCoin("1003denom2"), "denom1",
				utils.ParseDec("0.9999"), newInt(1000), 0,
			),
			RefundedCoin: utils.ParseCoin("0denom2"),
		},
		{
			Name: "refund tc 4",
			Msg: types.NewMsgLimitOrder(
				appID1, addr1, pair.Id, types.OrderDirectionBuy, utils.ParseCoin("102denom2"), "denom1",
				utils.ParseDec("1.001"), newInt(100), 0,
			),
			RefundedCoin: utils.ParseCoin("1denom2"),
		},
		{
			Name: "refund tc 5",
			Msg: types.NewMsgLimitOrder(
				appID1, addr1, pair.Id, types.OrderDirectionSell, utils.ParseCoin("1003denom1"), "denom2",
				utils.ParseDec("1.100"), newInt(1000), 0,
			),
			RefundedCoin: utils.ParseCoin("0denom1"),
		},
		{
			Name: "refund tc 6",
			Msg: types.NewMsgLimitOrder(
				appID1, addr1, pair.Id, types.OrderDirectionSell, utils.ParseCoin("1000denom1"), "denom2",
				utils.ParseDec("1.100"), newInt(100), 0),
			RefundedCoin: utils.ParseCoin("900denom1"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			s.Require().NoError(tc.Msg.ValidateBasic())

			balanceBefore := s.getBalance(addr1, tc.Msg.OfferCoin.Denom)
			_, err := s.keeper.LimitOrder(s.ctx, tc.Msg)
			s.Require().NoError(err)

			balanceAfter := s.getBalance(addr1, tc.Msg.OfferCoin.Denom)

			refundedCoin := balanceAfter.Sub(balanceBefore.Sub(tc.Msg.OfferCoin))
			s.Require().True(tc.RefundedCoin.IsEqual(refundedCoin))
		})
	}
}

func (s *KeeperTestSuite) TestLimitOrderWithPoolSwap() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appOne")

	asset1 := s.CreateNewAsset("ASSET1", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSET2", "uasset2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	_ = s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000000uasset1,500000000000uasset2")

	currentTime := time.Now()
	s.ctx = s.ctx.WithBlockTime(currentTime)

	testCases := []struct {
		Name                  string
		Msg                   types.MsgLimitOrder
		ExpResp               types.Order
		ExpOrderStatus        types.OrderStatus
		ExpBalanceAfterExpire sdk.Coins
	}{
		{
			Name: "swap at pool price buy direction",
			Msg: *types.NewMsgLimitOrder(
				appID1,
				s.addr(2),
				pair.Id,
				types.OrderDirectionBuy,
				utils.ParseCoin("1003000uasset2"),
				asset1.Denom,
				sdk.MustNewDecFromStr("0.5"),
				newInt(2000000),
				time.Minute*1,
			),
			ExpResp: types.Order{
				Id:                 1,
				PairId:             1,
				MsgHeight:          0,
				Orderer:            s.addr(2).String(),
				Direction:          types.OrderDirectionBuy,
				OfferCoin:          utils.ParseCoin("1000000uasset2"),
				RemainingOfferCoin: utils.ParseCoin("1000000uasset2"),
				ReceivedCoin:       utils.ParseCoin("0uasset1"),
				Price:              sdk.MustNewDecFromStr("0.5"),
				Amount:             newInt(2000000),
				OpenAmount:         newInt(2000000),
				BatchId:            1,
				ExpireAt:           s.ctx.BlockTime().Add(time.Minute * 1),
				Status:             types.OrderStatusNotExecuted,
				AppId:              appID1,
			},
			ExpOrderStatus:        types.OrderStatusNotMatched,
			ExpBalanceAfterExpire: utils.ParseCoins("1003000uasset2"),
		},
		{
			Name: "swap at slight higher than pool price buy direction",
			Msg: *types.NewMsgLimitOrder(
				appID1,
				s.addr(3),
				pair.Id,
				types.OrderDirectionBuy,
				utils.ParseCoin("1005006uasset2"),
				asset1.Denom,
				sdk.MustNewDecFromStr("0.501"),
				newInt(2000000),
				time.Minute*1,
			),
			ExpResp: types.Order{
				Id:                 2,
				PairId:             1,
				MsgHeight:          0,
				Orderer:            s.addr(3).String(),
				Direction:          types.OrderDirectionBuy,
				OfferCoin:          utils.ParseCoin("1002000uasset2"),
				RemainingOfferCoin: utils.ParseCoin("1002000uasset2"),
				ReceivedCoin:       utils.ParseCoin("0uasset1"),
				Price:              sdk.MustNewDecFromStr("0.501"),
				Amount:             newInt(2000000),
				OpenAmount:         newInt(2000000),
				BatchId:            2,
				ExpireAt:           s.ctx.BlockTime().Add(time.Minute * 1),
				Status:             types.OrderStatusNotExecuted,
				AppId:              appID1,
			},
			ExpOrderStatus:        types.OrderStatusCompleted,
			ExpBalanceAfterExpire: utils.ParseCoins("2000000uasset1,2003uasset2"),
		},
		{
			Name: "swap at slight lower than pool price buy direction",
			Msg: *types.NewMsgLimitOrder(
				appID1,
				s.addr(4),
				pair.Id,
				types.OrderDirectionBuy,
				utils.ParseCoin("1000994uasset2"),
				asset1.Denom,
				sdk.MustNewDecFromStr("0.499"),
				newInt(2000000),
				time.Minute*1,
			),
			ExpResp: types.Order{
				Id:                 3,
				PairId:             1,
				MsgHeight:          0,
				Orderer:            s.addr(4).String(),
				Direction:          types.OrderDirectionBuy,
				OfferCoin:          utils.ParseCoin("998000uasset2"),
				RemainingOfferCoin: utils.ParseCoin("998000uasset2"),
				ReceivedCoin:       utils.ParseCoin("0uasset1"),
				Price:              sdk.MustNewDecFromStr("0.499"),
				Amount:             newInt(2000000),
				OpenAmount:         newInt(2000000),
				BatchId:            3,
				ExpireAt:           s.ctx.BlockTime().Add(time.Minute * 1),
				Status:             types.OrderStatusNotExecuted,
				AppId:              appID1,
			},
			ExpOrderStatus:        types.OrderStatusNotMatched,
			ExpBalanceAfterExpire: utils.ParseCoins("1000994uasset2"),
		},
		{
			Name: "swap at pool price sell direction",
			Msg: *types.NewMsgLimitOrder(
				appID1,
				s.addr(5),
				pair.Id,
				types.OrderDirectionSell,
				utils.ParseCoin("2006000uasset1"),
				asset2.Denom,
				sdk.MustNewDecFromStr("0.501"),
				newInt(2000000),
				time.Minute*1,
			),
			ExpResp: types.Order{
				Id:                 4,
				PairId:             1,
				MsgHeight:          0,
				Orderer:            s.addr(5).String(),
				Direction:          types.OrderDirectionSell,
				OfferCoin:          utils.ParseCoin("2000000uasset1"),
				RemainingOfferCoin: utils.ParseCoin("2000000uasset1"),
				ReceivedCoin:       utils.ParseCoin("0uasset2"),
				Price:              sdk.MustNewDecFromStr("0.501"),
				Amount:             newInt(2000000),
				OpenAmount:         newInt(2000000),
				BatchId:            4,
				ExpireAt:           s.ctx.BlockTime().Add(time.Minute * 1),
				Status:             types.OrderStatusNotExecuted,
				AppId:              appID1,
			},
			ExpOrderStatus:        types.OrderStatusNotMatched,
			ExpBalanceAfterExpire: utils.ParseCoins("2006000uasset1"),
		},
		{
			Name: "swap at slight higher than pool price sell direction",
			Msg: *types.NewMsgLimitOrder(
				appID1,
				s.addr(6),
				pair.Id,
				types.OrderDirectionSell,
				utils.ParseCoin("2006000uasset1"),
				asset2.Denom,
				sdk.MustNewDecFromStr("0.51"),
				newInt(2000000),
				time.Minute*1,
			),
			ExpResp: types.Order{
				Id:                 5,
				PairId:             1,
				MsgHeight:          0,
				Orderer:            s.addr(6).String(),
				Direction:          types.OrderDirectionSell,
				OfferCoin:          utils.ParseCoin("2000000uasset1"),
				RemainingOfferCoin: utils.ParseCoin("2000000uasset1"),
				ReceivedCoin:       utils.ParseCoin("0uasset2"),
				Price:              sdk.MustNewDecFromStr("0.51"),
				Amount:             newInt(2000000),
				OpenAmount:         newInt(2000000),
				BatchId:            5,
				ExpireAt:           s.ctx.BlockTime().Add(time.Minute * 1),
				Status:             types.OrderStatusNotExecuted,
				AppId:              appID1,
			},
			ExpOrderStatus:        types.OrderStatusNotMatched,
			ExpBalanceAfterExpire: utils.ParseCoins("2006000uasset1"),
		},
		{
			Name: "swap at slight lower than pool price sell direction",
			Msg: *types.NewMsgLimitOrder(
				appID1,
				s.addr(7),
				pair.Id,
				types.OrderDirectionSell,
				utils.ParseCoin("2006000uasset1"),
				asset2.Denom,
				sdk.MustNewDecFromStr("0.50"),
				newInt(2000000),
				time.Minute*1,
			),
			ExpResp: types.Order{
				Id:                 6,
				PairId:             1,
				MsgHeight:          0,
				Orderer:            s.addr(7).String(),
				Direction:          types.OrderDirectionSell,
				OfferCoin:          utils.ParseCoin("2000000uasset1"),
				RemainingOfferCoin: utils.ParseCoin("2000000uasset1"),
				ReceivedCoin:       utils.ParseCoin("0uasset2"),
				Price:              sdk.MustNewDecFromStr("0.50"),
				Amount:             newInt(2000000),
				OpenAmount:         newInt(2000000),
				BatchId:            6,
				ExpireAt:           s.ctx.BlockTime().Add(time.Minute * 1),
				Status:             types.OrderStatusNotExecuted,
				AppId:              appID1,
			},
			ExpOrderStatus:        types.OrderStatusCompleted,
			ExpBalanceAfterExpire: utils.ParseCoins("1000002uasset2"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			s.fundAddr(tc.Msg.GetOrderer(), sdk.NewCoins(tc.Msg.OfferCoin))

			// order placed
			order, err := s.keeper.LimitOrder(s.ctx, &tc.Msg)
			s.Require().NoError(err)
			s.Require().IsType(types.Order{}, order)
			s.Require().Equal(tc.ExpResp, order)

			// execute order request
			s.nextBlock()

			if tc.ExpOrderStatus != types.OrderStatusCompleted {
				order, found := s.keeper.GetOrder(s.ctx, appID1, pair.Id, order.Id)
				s.Require().True(found)
				s.Require().Equal(tc.ExpOrderStatus, order.Status)

				// make order expire
				s.ctx = s.ctx.WithBlockTime(tc.ExpResp.ExpireAt)
				s.nextBlock()
			}

			_, found := s.keeper.GetOrder(s.ctx, appID1, pair.Id, order.Id)
			s.Require().False(found)

			availableBalance := s.getBalances(tc.Msg.GetOrderer())
			s.Require().True(tc.ExpBalanceAfterExpire.IsEqual(availableBalance))

			// reset to default time
			s.ctx = s.ctx.WithBlockTime(currentTime)
		})
	}
}

func (s *KeeperTestSuite) TestLimitOrderWithoutPool() {
	addr1 := s.addr(1)
	dummyAcc := s.addr(696969)

	appID1 := s.CreateNewApp("appOne")

	asset1 := s.CreateNewAsset("ASSET1", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSET2", "uasset2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)

	currentTime := time.Now()
	s.ctx = s.ctx.WithBlockTime(currentTime)

	accumulatedSwapFee := s.getBalances(pair.GetSwapFeeCollectorAddress())

	testCases := []struct {
		Name               string
		BuyMsg             types.MsgLimitOrder
		BuyerExpBalance    sdk.Coins
		BuyExpOrderStatus  types.OrderStatus
		SellMsg            types.MsgLimitOrder
		SellerExpBalance   sdk.Coins
		SellExpOrderStatus types.OrderStatus
		CollectedSwapFee   sdk.Coins
	}{
		{
			Name: "buyer seller full order match",
			BuyMsg: *types.NewMsgLimitOrder(
				appID1,
				s.addr(2),
				pair.Id,
				types.OrderDirectionBuy,
				utils.ParseCoin("17552500uasset2"),
				asset1.Denom,
				sdk.MustNewDecFromStr("0.5"),
				newInt(35000000),
				time.Minute*1,
			),
			BuyerExpBalance:   utils.ParseCoins("35000000uasset1"),
			BuyExpOrderStatus: types.OrderStatusCompleted,

			SellMsg: *types.NewMsgLimitOrder(
				appID1,
				s.addr(3),
				pair.Id,
				types.OrderDirectionSell,
				utils.ParseCoin("35105000uasset1"),
				asset2.Denom,
				sdk.MustNewDecFromStr("0.5"),
				newInt(35000000),
				time.Minute*1,
			),
			SellerExpBalance:   utils.ParseCoins("17500000uasset2"),
			SellExpOrderStatus: types.OrderStatusCompleted,
			CollectedSwapFee:   utils.ParseCoins("105000uasset1,52500uasset2"),
		},
		{
			Name: "buyer partial order match seller full order match",
			BuyMsg: *types.NewMsgLimitOrder(
				appID1,
				s.addr(2),
				pair.Id,
				types.OrderDirectionBuy,
				utils.ParseCoin("17552500uasset2"),
				asset1.Denom,
				sdk.MustNewDecFromStr("0.5"),
				newInt(35000000),
				time.Minute*1,
			),
			BuyerExpBalance:   utils.ParseCoins("20000000uasset1,7522500uasset2"),
			BuyExpOrderStatus: types.OrderStatusPartiallyMatched,

			SellMsg: *types.NewMsgLimitOrder(
				appID1,
				s.addr(3),
				pair.Id,
				types.OrderDirectionSell,
				utils.ParseCoin("20060000uasset1"),
				asset2.Denom,
				sdk.MustNewDecFromStr("0.5"),
				newInt(20000000),
				time.Minute*1,
			),
			SellerExpBalance:   utils.ParseCoins("10000000uasset2"),
			SellExpOrderStatus: types.OrderStatusCompleted,
			CollectedSwapFee:   utils.ParseCoins("60000uasset1,30000uasset2"),
		},
		{
			Name: "buyer full order match seller partial order match",
			BuyMsg: *types.NewMsgLimitOrder(
				appID1,
				s.addr(2),
				pair.Id,
				types.OrderDirectionBuy,
				utils.ParseCoin("15045000uasset2"),
				asset1.Denom,
				sdk.MustNewDecFromStr("0.5"),
				newInt(30000000),
				time.Minute*1,
			),
			BuyerExpBalance:   utils.ParseCoins("30000000uasset1"),
			BuyExpOrderStatus: types.OrderStatusCompleted,

			SellMsg: *types.NewMsgLimitOrder(
				appID1,
				s.addr(3),
				pair.Id,
				types.OrderDirectionSell,
				utils.ParseCoin("70210000uasset1"),
				asset2.Denom,
				sdk.MustNewDecFromStr("0.5"),
				newInt(70000000),
				time.Minute*1,
			),
			SellerExpBalance:   utils.ParseCoins("15000000uasset2,40120000uasset1"),
			SellExpOrderStatus: types.OrderStatusPartiallyMatched,
			CollectedSwapFee:   utils.ParseCoins("90000uasset1,45000uasset2"),
		},
		{
			Name: "buyer seller no order match",
			BuyMsg: *types.NewMsgLimitOrder(
				appID1,
				s.addr(2),
				pair.Id,
				types.OrderDirectionBuy,
				utils.ParseCoin("15045000uasset2"),
				asset1.Denom,
				sdk.MustNewDecFromStr("0.45"),
				newInt(30000000),
				time.Minute*1,
			),
			BuyerExpBalance:   utils.ParseCoins("15045000uasset2"),
			BuyExpOrderStatus: types.OrderStatusNotMatched,

			SellMsg: *types.NewMsgLimitOrder(
				appID1,
				s.addr(3),
				pair.Id,
				types.OrderDirectionSell,
				utils.ParseCoin("70210000uasset1"),
				asset2.Denom,
				sdk.MustNewDecFromStr("0.55"),
				newInt(70000000),
				time.Minute*1,
			),
			SellerExpBalance:   utils.ParseCoins("70210000uasset1"),
			SellExpOrderStatus: types.OrderStatusNotMatched,
			CollectedSwapFee:   utils.ParseCoins("0uasset1,0uasset2"),
		},
		{
			Name: "buyer price high seller price low",
			BuyMsg: *types.NewMsgLimitOrder(
				appID1,
				s.addr(2),
				pair.Id,
				types.OrderDirectionBuy,
				utils.ParseCoin("16549500uasset2"),
				asset1.Denom,
				sdk.MustNewDecFromStr("0.55"),
				newInt(30000000),
				time.Minute*1,
			),
			BuyerExpBalance:   utils.ParseCoins("30000000uasset1,752250uasset2"),
			BuyExpOrderStatus: types.OrderStatusCompleted,

			SellMsg: *types.NewMsgLimitOrder(
				appID1,
				s.addr(3),
				pair.Id,
				types.OrderDirectionSell,
				utils.ParseCoin("30090000uasset1"),
				asset2.Denom,
				sdk.MustNewDecFromStr("0.50"),
				newInt(30000000),
				time.Minute*1,
			),
			SellerExpBalance:   utils.ParseCoins("15750000uasset2"),
			SellExpOrderStatus: types.OrderStatusCompleted,
			CollectedSwapFee:   utils.ParseCoins("90000uasset1,47250uasset2"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			s.fundAddr(tc.BuyMsg.GetOrderer(), sdk.NewCoins(tc.BuyMsg.OfferCoin))
			s.fundAddr(tc.SellMsg.GetOrderer(), sdk.NewCoins(tc.SellMsg.OfferCoin))

			// buy order placed
			buyOrder, err := s.keeper.LimitOrder(s.ctx, &tc.BuyMsg)
			s.Require().NoError(err)
			s.Require().IsType(types.Order{}, buyOrder)
			s.Require().Equal(types.OrderStatusNotExecuted, buyOrder.Status)

			// buy order placed
			sellOrder, err := s.keeper.LimitOrder(s.ctx, &tc.SellMsg)
			s.Require().NoError(err)
			s.Require().IsType(types.Order{}, sellOrder)
			s.Require().Equal(types.OrderStatusNotExecuted, sellOrder.Status)

			// execute order request
			s.nextBlock()

			bOrder, found := s.keeper.GetOrder(s.ctx, tc.BuyMsg.AppId, tc.BuyMsg.PairId, buyOrder.Id)
			if tc.BuyExpOrderStatus != types.OrderStatusCompleted {
				s.Require().True(found)
				s.Require().Equal(tc.BuyExpOrderStatus, bOrder.Status)
			}

			sOrder, found := s.keeper.GetOrder(s.ctx, tc.SellMsg.AppId, tc.SellMsg.PairId, sellOrder.Id)
			if tc.SellExpOrderStatus != types.OrderStatusCompleted {
				s.Require().True(found)
				s.Require().Equal(tc.SellExpOrderStatus, sOrder.Status)
			}

			// change blocktime, so order gets expired
			s.ctx = s.ctx.WithBlockTime(currentTime.Add(tc.BuyMsg.OrderLifespan))
			s.nextBlock()

			_, found = s.keeper.GetOrder(s.ctx, tc.BuyMsg.AppId, tc.BuyMsg.PairId, buyOrder.Id)
			s.Require().False(found)
			_, found = s.keeper.GetOrder(s.ctx, tc.SellMsg.AppId, tc.SellMsg.PairId, sellOrder.Id)
			s.Require().False(found)

			buyerAvailableBalance := s.getBalances(tc.BuyMsg.GetOrderer())
			s.Require().True(tc.BuyerExpBalance.IsEqual(buyerAvailableBalance))

			selllerAvailableBalance := s.getBalances(tc.SellMsg.GetOrderer())
			s.Require().True(tc.SellerExpBalance.IsEqual(selllerAvailableBalance))

			// verify swapfee coolection
			accumulatedSwapFee = accumulatedSwapFee.Add(tc.CollectedSwapFee...)
			availableSwapFees := s.getBalances(pair.GetSwapFeeCollectorAddress())
			s.Require().True(accumulatedSwapFee.IsEqual(availableSwapFees))

			// transfer all funds from testing account to dummy account
			// for reusing the accounts, leads to easy account balance calculation
			s.sendCoins(tc.BuyMsg.GetOrderer(), dummyAcc, s.getBalances(tc.BuyMsg.GetOrderer()))
			s.sendCoins(tc.SellMsg.GetOrderer(), dummyAcc, s.getBalances(tc.SellMsg.GetOrderer()))

			// reset to default time
			s.ctx = s.ctx.WithBlockTime(currentTime)
		})
	}
}

func (s *KeeperTestSuite) TestMarketOrder() {
	creator := s.addr(1)
	trader := s.addr(2)
	// escrow := s.addr(3)

	appID1 := s.CreateNewApp("appOne")

	asset1 := s.CreateNewAsset("ASSET1", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSET2", "uasset2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset2.Denom)

	params, err := s.keeper.GetGenericParams(s.ctx, appID1)
	s.Require().NoError(err)

	testCases := []struct {
		Name    string
		Msg     types.MsgMarketOrder
		ExpErr  error
		ExpResp *types.Order
	}{
		{
			Name: "error invalid app id",
			Msg: *types.NewMsgMarketOrder(
				69,
				trader,
				pair.Id,
				types.OrderDirectionBuy,
				utils.ParseCoin("10030000uasset2"),
				asset1.Denom,
				newInt(10000000),
				time.Second*10,
			),
			ExpErr:  sdkerrors.Wrap(sdkerrors.Wrapf(types.ErrInvalidAppID, "app id %d not found", 69), "params retreval failed"),
			ExpResp: &types.Order{},
		},
		{
			Name: "error max order life span",
			Msg: *types.NewMsgMarketOrder(
				appID1,
				trader,
				pair.Id,
				types.OrderDirectionBuy,
				utils.ParseCoin("10030000uasset2"),
				asset1.Denom,
				newInt(10000000),
				time.Hour*48,
			),
			ExpErr:  sdkerrors.Wrapf(types.ErrTooLongOrderLifespan, "%s is longer than %s", time.Hour*48, params.MaxOrderLifespan),
			ExpResp: &types.Order{},
		},
		{
			Name: "error invalid pair id",
			Msg: *types.NewMsgMarketOrder(
				appID1,
				trader,
				69,
				types.OrderDirectionBuy,
				utils.ParseCoin("10030000uasset2"),
				asset1.Denom,
				newInt(10000000),
				time.Second*10,
			),
			ExpErr:  sdkerrors.Wrapf(sdkerrors.ErrNotFound, "pair %d not found", 69),
			ExpResp: &types.Order{},
		},
		{
			Name: "error last price not available",
			Msg: *types.NewMsgMarketOrder(
				appID1,
				trader,
				pair.Id,
				types.OrderDirectionBuy,
				utils.ParseCoin("1003000uasset2"),
				asset1.Denom,
				newInt(1000000),
				time.Second*10,
			),
			ExpErr:  types.ErrNoLastPrice,
			ExpResp: &types.Order{},
		},
		{
			Name: "error invalid denom pair buy direction",
			Msg: *types.NewMsgMarketOrder(
				appID1,
				trader,
				pair.Id,
				types.OrderDirectionBuy,
				utils.ParseCoin("10030000uasset1"),
				asset2.Denom,
				newInt(10000000),
				time.Second*10,
			),
			ExpErr:  sdkerrors.Wrapf(types.ErrWrongPair, "denom pair (%s, %s) != (%s, %s)", asset2.Denom, asset1.Denom, pair.BaseCoinDenom, pair.QuoteCoinDenom),
			ExpResp: &types.Order{},
		},
		{
			Name: "error invalid denom pair sell direction",
			Msg: *types.NewMsgMarketOrder(
				appID1,
				trader,
				pair.Id,
				types.OrderDirectionSell,
				utils.ParseCoin("10030000uasset2"),
				asset1.Denom,
				newInt(10000000),
				time.Second*10,
			),
			ExpErr:  sdkerrors.Wrapf(types.ErrWrongPair, "denom pair (%s, %s) != (%s, %s)", asset2.Denom, asset1.Denom, pair.BaseCoinDenom, pair.QuoteCoinDenom),
			ExpResp: &types.Order{},
		},
		{
			Name: "error insufficient offer coin buy direction",
			Msg: *types.NewMsgMarketOrder(
				appID1,
				trader,
				pair.Id,
				types.OrderDirectionBuy,
				utils.ParseCoin("10000000uasset2"), // swap fee excluded, also at price = 1.1
				asset1.Denom,
				newInt(10000000),
				time.Second*10,
			),
			ExpErr:  sdkerrors.Wrapf(types.ErrInsufficientOfferCoin, "10000000uasset2 is smaller than 11033000uasset2"),
			ExpResp: &types.Order{},
		},
		{
			Name: "error insufficient offer coin sell direction",
			Msg: *types.NewMsgMarketOrder(
				appID1,
				trader,
				pair.Id,
				types.OrderDirectionSell,
				utils.ParseCoin("10000000uasset1"), // swap fee excluded
				asset2.Denom,
				newInt(10000000),
				time.Second*10,
			),
			ExpErr:  sdkerrors.Wrapf(types.ErrInsufficientOfferCoin, "10000000uasset1 is smaller than 10030000uasset1"),
			ExpResp: &types.Order{},
		},
		{
			Name: "error too small order buy direction",
			Msg: *types.NewMsgMarketOrder(
				appID1,
				trader,
				pair.Id,
				types.OrderDirectionBuy,
				utils.ParseCoin("130uasset2"),
				asset1.Denom,
				newInt(99),
				time.Second*10,
			),
			ExpErr:  types.ErrTooSmallOrder,
			ExpResp: &types.Order{},
		},
		{
			Name: "error too small order sell direction",
			Msg: *types.NewMsgMarketOrder(
				appID1,
				trader,
				pair.Id,
				types.OrderDirectionSell,
				utils.ParseCoin("130uasset1"),
				asset2.Denom,
				newInt(99),
				time.Second*10,
			),
			ExpErr:  types.ErrTooSmallOrder,
			ExpResp: &types.Order{},
		},
		{
			Name: "error insufficient funds",
			Msg: *types.NewMsgMarketOrder(
				appID1,
				trader,
				pair.Id,
				types.OrderDirectionBuy,
				utils.ParseCoin("1103300uasset2"),
				asset1.Denom,
				newInt(1000000),
				time.Second*10,
			),
			ExpErr:  sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "0uasset2 is smaller than 1103300uasset2"),
			ExpResp: &types.Order{},
		},
		{
			Name: "success valid case buy direction",
			Msg: *types.NewMsgMarketOrder(
				appID1,
				trader,
				pair.Id,
				types.OrderDirectionBuy,
				utils.ParseCoin("1103300uasset2"),
				asset1.Denom,
				newInt(1000000),
				time.Second*10,
			),
			ExpErr: nil,
			ExpResp: &types.Order{
				Id:                 3,
				PairId:             1,
				MsgHeight:          0,
				Orderer:            trader.String(),
				Direction:          types.OrderDirectionBuy,
				OfferCoin:          utils.ParseCoin("1100000uasset2"),
				RemainingOfferCoin: utils.ParseCoin("1100000uasset2"),
				ReceivedCoin:       utils.ParseCoin("0uasset1"),
				Price:              sdk.MustNewDecFromStr("1.1"),
				Amount:             newInt(1000000),
				OpenAmount:         newInt(1000000),
				BatchId:            2,
				ExpireAt:           s.ctx.BlockTime().Add(time.Second * 10),
				Status:             types.OrderStatusNotExecuted,
				AppId:              appID1,
			},
		},
		{
			Name: "success valid case sell direction",
			Msg: *types.NewMsgMarketOrder(
				appID1,
				trader,
				pair.Id,
				types.OrderDirectionSell,
				utils.ParseCoin("1003000uasset1"),
				asset2.Denom,
				newInt(1000000),
				time.Second*10,
			),
			ExpErr: nil,
			ExpResp: &types.Order{
				Id:                 4,
				PairId:             1,
				MsgHeight:          0,
				Orderer:            trader.String(),
				Direction:          types.OrderDirectionSell,
				OfferCoin:          utils.ParseCoin("1000000uasset1"),
				RemainingOfferCoin: utils.ParseCoin("1000000uasset1"),
				ReceivedCoin:       utils.ParseCoin("0uasset2"),
				Price:              sdk.MustNewDecFromStr("0.9"),
				Amount:             newInt(1000000),
				OpenAmount:         newInt(1000000),
				BatchId:            2,
				ExpireAt:           s.ctx.BlockTime().Add(time.Second * 10),
				Status:             types.OrderStatusNotExecuted,
				AppId:              appID1,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.Name, func() {
			if tc.ExpErr == nil {
				s.fundAddr(tc.Msg.GetOrderer(), sdk.NewCoins(tc.Msg.OfferCoin))
			}
			order, err := s.keeper.MarketOrder(s.ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				s.Require().Equal(tc.ExpResp, &order)
			} else {
				s.Require().NoError(err)
				s.Require().IsType(tc.ExpResp, &order)
				s.Require().Equal(tc.ExpResp, &order)

				qorder, found := s.keeper.GetOrder(s.ctx, tc.Msg.AppId, tc.Msg.PairId, order.Id)
				s.Require().True(found)
				s.Require().Equal(qorder, order)
			}

			// make limit order after this testcase
			if tc.ExpErr == types.ErrNoLastPrice {
				// When there is no last price in the pair, only limit orders can be made.
				// These two orders will be matched.
				s.LimitOrder(tc.Msg.AppId, creator, tc.Msg.PairId, types.OrderDirectionBuy, utils.ParseDec("1"), sdk.NewInt(10000), 0)
				s.LimitOrder(tc.Msg.AppId, creator, tc.Msg.PairId, types.OrderDirectionSell, utils.ParseDec("1"), sdk.NewInt(10000), 0)
				s.nextBlock()
			}
		})
	}
}

func (s *KeeperTestSuite) TestMarketOrderTwo() {
	appID1 := s.CreateNewApp("appOne")

	asset1 := s.CreateNewAsset("ASSET1", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSET2", "uasset2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)

	// When there is no last price in the pair, only limit orders can be made.
	// These two orders will be matched.
	s.LimitOrder(appID1, s.addr(1), pair.Id, types.OrderDirectionBuy, utils.ParseDec("1.0"), sdk.NewInt(10000), 0)
	s.LimitOrder(appID1, s.addr(2), pair.Id, types.OrderDirectionSell, utils.ParseDec("1.0"), sdk.NewInt(10000), 0)
	s.nextBlock()

	// Now users can make market orders.
	// In this case, addr(3) user's order takes higher priority than addr(4) user's,
	// because market buy orders have 10% higher price than the last price(1.0).
	s.MarketOrder(appID1, s.addr(3), pair.Id, types.OrderDirectionBuy, sdk.NewInt(10000), 0)
	s.LimitOrder(appID1, s.addr(4), pair.Id, types.OrderDirectionBuy, utils.ParseDec("1.08"), sdk.NewInt(10000), 0)
	s.LimitOrder(appID1, s.addr(5), pair.Id, types.OrderDirectionSell, utils.ParseDec("1.07"), sdk.NewInt(10000), 0)
	s.nextBlock()

	// Check the result.
	s.Require().True(utils.ParseCoin("10000uasset1").IsEqual(s.getBalance(s.addr(3), "uasset1")))
	s.Require().True(utils.ParseCoins("10832uasset2").IsEqual(s.getBalances(s.addr(4))))
}

func (s *KeeperTestSuite) TestMarketOrderRefund() {
	appID1 := s.CreateNewApp("appOne")

	asset1 := s.CreateNewAsset("ASSET1", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSET2", "denom2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)
	p := utils.ParseDec("1.0")
	pair.LastPrice = &p
	s.keeper.SetPair(s.ctx, pair)
	orderer := s.addr(1)
	s.fundAddr(orderer, utils.ParseCoins("1000000000denom1,1000000000denom2"))

	for _, tc := range []struct {
		msg          *types.MsgMarketOrder
		refundedCoin sdk.Coin
	}{
		{
			types.NewMsgMarketOrder(
				appID1, orderer, pair.Id, types.OrderDirectionBuy, utils.ParseCoin("1103300denom2"), "denom1",
				newInt(1000000), 0),
			utils.ParseCoin("0denom2"),
		},
		{
			types.NewMsgMarketOrder(
				appID1, orderer, pair.Id, types.OrderDirectionBuy, utils.ParseCoin("1000000denom2"), "denom1",
				newInt(10000), 0),
			utils.ParseCoin("988967denom2"),
		},
		{
			types.NewMsgMarketOrder(
				appID1, orderer, pair.Id, types.OrderDirectionSell, utils.ParseCoin("1000000denom1"), "denom2",
				newInt(10000), 0),
			utils.ParseCoin("989970denom1"),
		},
	} {
		s.Run("", func() {
			s.Require().NoError(tc.msg.ValidateBasic())

			balanceBefore := s.getBalance(orderer, tc.msg.OfferCoin.Denom)
			_, err := s.keeper.MarketOrder(s.ctx, tc.msg)
			s.Require().NoError(err)

			balanceAfter := s.getBalance(orderer, tc.msg.OfferCoin.Denom)

			refundedCoin := balanceAfter.Sub(balanceBefore.Sub(tc.msg.OfferCoin))
			s.Require().True(tc.refundedCoin.IsEqual(refundedCoin))
		})
	}
}

func (s *KeeperTestSuite) TestMarketOrderWithNoLastPrice() {
	appID1 := s.CreateNewApp("appOne")

	asset1 := s.CreateNewAsset("ASSET1", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSET2", "denom2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)

	s.Require().Nil(pair.LastPrice)
	offerCoin := utils.ParseCoin("10000denom2")
	s.fundAddr(s.addr(1), sdk.NewCoins(offerCoin))
	msg := types.NewMsgMarketOrder(
		appID1, s.addr(1), pair.Id, types.OrderDirectionBuy, offerCoin, "denom1", sdk.NewInt(10000), 0)
	_, err := s.keeper.MarketOrder(s.ctx, msg)
	s.Require().ErrorIs(err, types.ErrNoLastPrice)
}

func (s *KeeperTestSuite) TestSingleOrderNoMatch() {
	k, ctx := s.keeper, s.ctx

	appID1 := s.CreateNewApp("appOne")

	asset1 := s.CreateNewAsset("ASSET1", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSET2", "denom2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)

	order := s.LimitOrder(appID1, s.addr(1), pair.Id, types.OrderDirectionBuy, utils.ParseDec("1.0"), sdk.NewInt(1000000), 10*time.Second)
	// Execute matching
	liquidity.EndBlocker(ctx, k)

	order, found := k.GetOrder(ctx, appID1, order.PairId, order.Id)
	s.Require().True(found)
	s.Require().Equal(types.OrderStatusNotMatched, order.Status)

	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(10 * time.Second))
	// Expire the order, here BeginBlocker is not called to check
	// the request's changed status
	liquidity.EndBlocker(ctx, k)

	order, _ = k.GetOrder(ctx, appID1, order.PairId, order.Id)
	s.Require().Equal(types.OrderStatusExpired, order.Status)

	s.Require().True(utils.ParseCoins("1003000denom2").IsEqual(s.getBalances(s.addr(1))))
}

func (s *KeeperTestSuite) TestTwoOrderExactMatch() {
	k, ctx := s.keeper, s.ctx

	appID1 := s.CreateNewApp("appOne")

	asset1 := s.CreateNewAsset("ASSET1", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSET2", "denom2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)

	req1 := s.LimitOrder(appID1, s.addr(1), pair.Id, types.OrderDirectionBuy, utils.ParseDec("1.0"), newInt(10000), time.Hour)
	req2 := s.LimitOrder(appID1, s.addr(2), pair.Id, types.OrderDirectionSell, utils.ParseDec("1.0"), newInt(10000), time.Hour)
	liquidity.EndBlocker(ctx, k)

	req1, _ = k.GetOrder(ctx, appID1, req1.PairId, req1.Id)
	s.Require().Equal(types.OrderStatusCompleted, req1.Status)
	req2, _ = k.GetOrder(ctx, appID1, req2.PairId, req2.Id)
	s.Require().Equal(types.OrderStatusCompleted, req2.Status)

	s.Require().True(utils.ParseCoins("10000denom1").IsEqual(s.getBalances(s.addr(1))))
	s.Require().True(utils.ParseCoins("10000denom2").IsEqual(s.getBalances(s.addr(2))))

	pair, _ = k.GetPair(ctx, appID1, pair.Id)
	s.Require().NotNil(pair.LastPrice)
	s.Require().True(utils.ParseDec("1.0").Equal(*pair.LastPrice))
}

func (s *KeeperTestSuite) TestPartialMatch() {
	appID1 := s.CreateNewApp("appOne")

	asset1 := s.CreateNewAsset("ASSET1", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSET2", "denom2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)

	order := s.LimitOrder(appID1, s.addr(1), pair.Id, types.OrderDirectionBuy, utils.ParseDec("1.0"), sdk.NewInt(10000), time.Hour)
	s.LimitOrder(appID1, s.addr(2), pair.Id, types.OrderDirectionSell, utils.ParseDec("1.0"), sdk.NewInt(5000), 0)
	s.nextBlock()

	order, found := s.keeper.GetOrder(s.ctx, appID1, order.PairId, order.Id)
	s.Require().True(found)
	s.Require().Equal(types.OrderStatusPartiallyMatched, order.Status)
	s.Require().True(utils.ParseCoin("5000denom2").IsEqual(order.RemainingOfferCoin))
	s.Require().True(utils.ParseCoin("5000denom1").IsEqual(order.ReceivedCoin))
	s.Require().True(sdk.NewInt(5000).Equal(order.OpenAmount))

	s.MarketOrder(appID1, s.addr(3), pair.Id, types.OrderDirectionSell, sdk.NewInt(5000), 0)
	s.nextBlock()

	// Now completely matched.
	_, found = s.keeper.GetOrder(s.ctx, appID1, order.PairId, order.Id)
	s.Require().False(found)
}

func (s *KeeperTestSuite) TestMatchWithLowPricePool() {
	appID1 := s.CreateNewApp("appOne")

	asset1 := s.CreateNewAsset("ASSET1", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSET2", "denom2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)

	// Create a pool with very low price.
	s.CreateNewLiquidityPool(appID1, pair.Id, s.addr(0), "10000000000000000000000000000000000000000denom1,1000000denom2")
	order := s.LimitOrder(appID1, s.addr(1), pair.Id, types.OrderDirectionBuy, utils.ParseDec("0.000000000001000000"), sdk.NewInt(100000000000000000), 10*time.Second)
	liquidity.EndBlocker(s.ctx, s.keeper)
	order, found := s.keeper.GetOrder(s.ctx, appID1, order.PairId, order.Id)
	s.Require().True(found)
	s.Require().Equal(types.OrderStatusNotMatched, order.Status)
}
