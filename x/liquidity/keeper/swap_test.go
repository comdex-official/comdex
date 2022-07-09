package keeper_test

import (
	"fmt"
	"time"

	utils "github.com/comdex-official/comdex/types"
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
				fmt.Println(qorder)
				s.Require().True(found)
				s.Require().Equal(qorder, order)
			}
		})
	}
}

func (s *KeeperTestSuite) TestLimitOrderRefundWithPool() {
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
