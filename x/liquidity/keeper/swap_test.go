package keeper_test

import (
	"math/rand"
	"time"

	errorsmod "cosmossdk.io/errors"

	sdkmath "cosmossdk.io/math"
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
	dummyAcc := s.addr(1234567890)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)

	params, err := s.keeper.GetGenericParams(s.ctx, appID1)
	s.Require().NoError(err)

	testCases := []struct {
		Name         string
		Msg          types.MsgLimitOrder
		FundRequired bool
		ExpErr       error
		ExpResp      *types.Order
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
			FundRequired: true,
			ExpErr:       errorsmod.Wrap(errorsmod.Wrapf(types.ErrInvalidAppID, "app id %d not found", 69), "params retreval failed"),
			ExpResp:      &types.Order{},
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
			FundRequired: true,
			ExpErr:       errorsmod.Wrapf(types.ErrTooLongOrderLifespan, "%s is longer than %s", time.Hour*48, params.MaxOrderLifespan),
			ExpResp:      &types.Order{},
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
			FundRequired: true,
			ExpErr:       errorsmod.Wrapf(sdkerrors.ErrNotFound, "pair %d not found", 69),
			ExpResp:      &types.Order{},
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
			FundRequired: true,
			ExpErr:       errorsmod.Wrapf(types.ErrPriceOutOfRange, "%s is higher than %s", amm.HighestTick(int(params.TickPrecision+1)), amm.HighestTick(int(params.TickPrecision))),
			ExpResp:      &types.Order{},
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
			FundRequired: true,
			ExpErr:       errorsmod.Wrapf(types.ErrPriceOutOfRange, "%s is lower than %s", amm.LowestTick(int(params.TickPrecision-1)), amm.LowestTick(int(params.TickPrecision))),
			ExpResp:      &types.Order{},
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
			FundRequired: true,
			ExpErr:       errorsmod.Wrapf(types.ErrWrongPair, "denom pair (%s, %s) != (%s, %s)", asset2.Denom, asset1.Denom, pair.BaseCoinDenom, pair.QuoteCoinDenom),
			ExpResp:      &types.Order{},
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
			FundRequired: true,
			ExpErr:       errorsmod.Wrapf(types.ErrWrongPair, "denom pair (%s, %s) != (%s, %s)", asset2.Denom, asset1.Denom, pair.BaseCoinDenom, pair.QuoteCoinDenom),
			ExpResp:      &types.Order{},
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
			FundRequired: true,
			ExpErr:       errorsmod.Wrapf(types.ErrInsufficientOfferCoin, "10000000uasset2 is smaller than 10030000uasset2"),
			ExpResp:      &types.Order{},
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
			FundRequired: true,
			ExpErr:       errorsmod.Wrapf(types.ErrInsufficientOfferCoin, "10000000uasset1 is smaller than 10030000uasset1"),
			ExpResp:      &types.Order{},
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
			FundRequired: true,
			ExpErr:       types.ErrTooSmallOrder,
			ExpResp:      &types.Order{},
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
			FundRequired: true,
			ExpErr:       types.ErrTooSmallOrder,
			ExpResp:      &types.Order{},
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
			FundRequired: false,
			ExpErr:       errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds, "0uasset2 is smaller than 1003000uasset2"),
			ExpResp:      &types.Order{},
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
			FundRequired: true,
			ExpErr:       nil,
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
				Type:               types.OrderTypeLimit,
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
			FundRequired: true,
			ExpErr:       nil,
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
				Type:               types.OrderTypeLimit,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.Name, func() {
			if tc.FundRequired {
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
			s.sendCoins(tc.Msg.GetOrderer(), dummyAcc, s.getBalances(tc.Msg.GetOrderer()))
		})
	}
}

func (s *KeeperTestSuite) TestLimitOrderRefund() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 2000000)

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

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	_ = s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000000uasset1,500000000000uasset2")

	currentTime := s.ctx.BlockTime()
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
				sdkmath.LegacyMustNewDecFromStr("0.5"),
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
				Price:              sdkmath.LegacyMustNewDecFromStr("0.5"),
				Amount:             newInt(2000000),
				OpenAmount:         newInt(2000000),
				BatchId:            1,
				ExpireAt:           s.ctx.BlockTime().Add(time.Minute * 1),
				Status:             types.OrderStatusNotExecuted,
				AppId:              appID1,
				Type:               types.OrderTypeLimit,
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
				sdkmath.LegacyMustNewDecFromStr("0.501"),
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
				Price:              sdkmath.LegacyMustNewDecFromStr("0.501"),
				Amount:             newInt(2000000),
				OpenAmount:         newInt(2000000),
				BatchId:            3,
				ExpireAt:           s.ctx.BlockTime().Add(time.Minute * 1),
				Status:             types.OrderStatusNotExecuted,
				AppId:              appID1,
				Type:               types.OrderTypeLimit,
			},
			ExpOrderStatus:        types.OrderStatusCompleted,
			ExpBalanceAfterExpire: utils.ParseCoins("2000000uasset1,1986uasset2"),
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
				sdkmath.LegacyMustNewDecFromStr("0.499"),
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
				Price:              sdkmath.LegacyMustNewDecFromStr("0.499"),
				Amount:             newInt(2000000),
				OpenAmount:         newInt(2000000),
				BatchId:            4,
				ExpireAt:           s.ctx.BlockTime().Add(time.Minute * 1),
				Status:             types.OrderStatusNotExecuted,
				AppId:              appID1,
				Type:               types.OrderTypeLimit,
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
				sdkmath.LegacyMustNewDecFromStr("0.501"),
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
				Price:              sdkmath.LegacyMustNewDecFromStr("0.501"),
				Amount:             newInt(2000000),
				OpenAmount:         newInt(2000000),
				BatchId:            6,
				ExpireAt:           s.ctx.BlockTime().Add(time.Minute * 1),
				Status:             types.OrderStatusNotExecuted,
				AppId:              appID1,
				Type:               types.OrderTypeLimit,
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
				sdkmath.LegacyMustNewDecFromStr("0.51"),
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
				Price:              sdkmath.LegacyMustNewDecFromStr("0.51"),
				Amount:             newInt(2000000),
				OpenAmount:         newInt(2000000),
				BatchId:            8,
				ExpireAt:           s.ctx.BlockTime().Add(time.Minute * 1),
				Status:             types.OrderStatusNotExecuted,
				AppId:              appID1,
				Type:               types.OrderTypeLimit,
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
				sdkmath.LegacyMustNewDecFromStr("0.499980"),
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
				Price:              sdkmath.LegacyMustNewDecFromStr("0.499980"),
				Amount:             newInt(2000000),
				OpenAmount:         newInt(2000000),
				BatchId:            10,
				ExpireAt:           s.ctx.BlockTime().Add(time.Minute * 1),
				Status:             types.OrderStatusNotExecuted,
				AppId:              appID1,
				Type:               types.OrderTypeLimit,
			},
			ExpOrderStatus:        types.OrderStatusCompleted,
			ExpBalanceAfterExpire: utils.ParseCoins("1000000uasset2"),
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
			s.Require().True(tc.ExpBalanceAfterExpire.Equal(availableBalance))

			// reset to default time
			s.ctx = s.ctx.WithBlockTime(currentTime)
		})
	}
}

func (s *KeeperTestSuite) TestLimitOrderWithoutPool() {
	addr1 := s.addr(1)
	dummyAcc := s.addr(696969)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)

	currentTime := s.ctx.BlockTime()
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
				sdkmath.LegacyMustNewDecFromStr("0.5"),
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
				sdkmath.LegacyMustNewDecFromStr("0.5"),
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
				sdkmath.LegacyMustNewDecFromStr("0.5"),
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
				sdkmath.LegacyMustNewDecFromStr("0.5"),
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
				sdkmath.LegacyMustNewDecFromStr("0.5"),
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
				sdkmath.LegacyMustNewDecFromStr("0.5"),
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
				sdkmath.LegacyMustNewDecFromStr("0.45"),
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
				sdkmath.LegacyMustNewDecFromStr("0.55"),
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
				sdkmath.LegacyMustNewDecFromStr("0.55"),
				newInt(30000000),
				time.Minute*1,
			),
			BuyerExpBalance:   utils.ParseCoins("30000000uasset1,1504500uasset2"),
			BuyExpOrderStatus: types.OrderStatusCompleted,

			SellMsg: *types.NewMsgLimitOrder(
				appID1,
				s.addr(3),
				pair.Id,
				types.OrderDirectionSell,
				utils.ParseCoin("30090000uasset1"),
				asset2.Denom,
				sdkmath.LegacyMustNewDecFromStr("0.50"),
				newInt(30000000),
				time.Minute*1,
			),
			SellerExpBalance:   utils.ParseCoins("15000000uasset2"),
			SellExpOrderStatus: types.OrderStatusCompleted,
			CollectedSwapFee:   utils.ParseCoins("90000uasset1,45000uasset2"),
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
			s.Require().True(tc.BuyerExpBalance.Equal(buyerAvailableBalance))

			selllerAvailableBalance := s.getBalances(tc.SellMsg.GetOrderer())
			s.Require().True(tc.SellerExpBalance.Equal(selllerAvailableBalance))

			// verify swapfee coolection
			accumulatedSwapFee = accumulatedSwapFee.Add(tc.CollectedSwapFee...)
			availableSwapFees := s.getBalances(pair.GetSwapFeeCollectorAddress())
			s.Require().True(accumulatedSwapFee.Equal(availableSwapFees))

			// transfer all funds from testing account to dummy account
			// for reusing the accounts, leads to easy account balance calculation
			s.sendCoins(tc.BuyMsg.GetOrderer(), dummyAcc, s.getBalances(tc.BuyMsg.GetOrderer()))
			s.sendCoins(tc.SellMsg.GetOrderer(), dummyAcc, s.getBalances(tc.SellMsg.GetOrderer()))

			// reset to default time
			s.ctx = s.ctx.WithBlockTime(currentTime)
		})
	}
}

func (s *KeeperTestSuite) TestDustCollector() {
	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)

	s.LimitOrder(appID1, s.addr(1), pair.Id, types.OrderDirectionBuy, utils.ParseDec("0.9005"), sdkmath.NewInt(1000), 0)
	s.LimitOrder(appID1, s.addr(2), pair.Id, types.OrderDirectionSell, utils.ParseDec("0.9005"), sdkmath.NewInt(1000), 0)
	s.nextBlock()

	s.Require().True(coinsEq(utils.ParseCoins("1000denom1,1denom2"), s.getBalances(s.addr(1))))
	s.Require().True(coinsEq(utils.ParseCoins("900denom2"), s.getBalances(s.addr(2))))

	s.Require().True(coinsEq(sdk.Coins{}, s.getBalances(pair.GetEscrowAddress())))
	params, err := s.keeper.GetGenericParams(s.ctx, appID1)
	s.Require().NoError(err)
	s.Require().True(coinsEq(utils.ParseCoins("1denom2"), s.getBalances(sdk.MustAccAddressFromBech32(params.DustCollectorAddress))))
}

func (s *KeeperTestSuite) TestInsufficientRemainingOfferCoin() {
	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)

	order := s.LimitOrder(appID1, s.addr(1), pair.Id, types.OrderDirectionBuy, utils.ParseDec("0.5"), sdkmath.NewInt(10000), time.Minute)
	s.LimitOrder(appID1, s.addr(2), pair.Id, types.OrderDirectionSell, utils.ParseDec("0.5"), sdkmath.NewInt(1001), 0)
	s.nextBlock()

	s.LimitOrder(appID1, s.addr(2), pair.Id, types.OrderDirectionSell, utils.ParseDec("0.5"), sdkmath.NewInt(8999), 0)
	liquidity.EndBlocker(s.ctx, s.keeper, s.app.AssetKeeper)
	order, found := s.keeper.GetOrder(s.ctx, appID1, order.PairId, order.Id)
	s.Require().True(found)
	s.Require().Equal(types.OrderStatusExpired, order.Status)
	s.Require().True(intEq(sdkmath.OneInt(), order.OpenAmount))
}

func (s *KeeperTestSuite) TestNegativeOpenAmount() {
	s.ctx = s.ctx.WithBlockHeight(1).WithBlockTime(utils.ParseTime("2022-03-01T00:00:00Z"))

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)

	order := s.LimitOrder(appID1, s.addr(1), pair.Id, types.OrderDirectionBuy, utils.ParseDec("0.82"), sdkmath.NewInt(648744), 0)
	s.LimitOrder(appID1, s.addr(2), pair.Id, types.OrderDirectionSell, utils.ParseDec("0.82"), sdkmath.NewInt(648745), 0)
	liquidity.EndBlocker(s.ctx, s.keeper, s.app.AssetKeeper)

	order, found := s.keeper.GetOrder(s.ctx, appID1, order.PairId, order.Id)
	s.Require().True(found)
	s.Require().False(order.OpenAmount.IsNegative())

	genState := s.keeper.ExportGenesis(s.ctx)
	s.Require().NotPanics(func() {
		s.keeper.InitGenesis(s.ctx, *genState)
	})
}

func (s *KeeperTestSuite) TestExpireSmallOrders() {
	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)
	s.LimitOrder(appID1, s.addr(1), pair.Id, types.OrderDirectionBuy, utils.ParseDec("0.000018"), sdkmath.NewInt(10000000), 0)

	// This order should have 10000 open amount after matching.
	// If this order would be matched after that, then the orderer will receive
	// floor(10000*0.000018) demand coin, which is zero.
	// So the order must have been expired after matching.
	order := s.LimitOrder(appID1, s.addr(2), pair.Id, types.OrderDirectionSell, utils.ParseDec("0.000018"), sdkmath.NewInt(10010000), time.Minute)
	liquidity.EndBlocker(s.ctx, s.keeper, s.app.AssetKeeper)
	order, found := s.keeper.GetOrder(s.ctx, appID1, order.PairId, order.Id)
	s.Require().True(found)
	s.Require().Equal(types.OrderStatusExpired, order.Status)
	liquidity.BeginBlocker(s.ctx, s.keeper, s.app.AssetKeeper) // Delete outdated states.

	s.LimitOrder(appID1, s.addr(1), pair.Id, types.OrderDirectionBuy, utils.ParseDec("0.000019"), sdkmath.NewInt(100000000), time.Minute)
	s.LimitOrder(appID1, s.addr(3), pair.Id, types.OrderDirectionSell, utils.ParseDec("0.000019"), sdkmath.NewInt(100000000), time.Minute)
	liquidity.EndBlocker(s.ctx, s.keeper, s.app.AssetKeeper)
}

func (s *KeeperTestSuite) TestRejectSmallOrders() {
	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)

	s.fundAddr(s.addr(1), utils.ParseCoins("10000000denom1,10000000denom2"))

	// Too small offer coin amount.
	msg := types.NewMsgLimitOrder(
		appID1, s.addr(1), pair.Id, types.OrderDirectionBuy, utils.ParseCoin("99denom2"),
		"denom1", utils.ParseDec("0.1"), sdkmath.NewInt(990), 0)
	s.Require().EqualError(msg.ValidateBasic(), "offer coin 99denom2 is smaller than the min amount 100: invalid request")

	// Too small order amount.
	msg = types.NewMsgLimitOrder(
		appID1, s.addr(1), pair.Id, types.OrderDirectionBuy, utils.ParseCoin("990denom2"),
		"denom1", utils.ParseDec("10.0"), sdkmath.NewInt(99), 0)
	s.Require().EqualError(msg.ValidateBasic(), "order amount 99 is smaller than the min amount 100: invalid request")

	// Too small orders.
	msg = types.NewMsgLimitOrder(
		appID1, s.addr(1), pair.Id, types.OrderDirectionBuy, utils.ParseCoin("101denom2"),
		"denom1", utils.ParseDec("0.00010001"), sdkmath.NewInt(999000), 0)
	s.Require().NoError(msg.ValidateBasic())
	_, err := s.keeper.LimitOrder(s.ctx, msg)
	s.Require().ErrorIs(err, types.ErrTooSmallOrder)

	msg = types.NewMsgLimitOrder(
		appID1, s.addr(1), pair.Id, types.OrderDirectionSell, utils.ParseCoin("1002999denom1"),
		"denom2", utils.ParseDec("0.0001"), sdkmath.NewInt(999999), 0)
	s.Require().NoError(msg.ValidateBasic())
	_, err = s.keeper.LimitOrder(s.ctx, msg)
	s.Require().ErrorIs(err, types.ErrTooSmallOrder)

	// Too small offer coin amount.
	msg2 := types.NewMsgMarketOrder(
		appID1, s.addr(1), pair.Id, types.OrderDirectionSell, utils.ParseCoin("99denom1"),
		"denom2", sdkmath.NewInt(99), 0)
	s.Require().EqualError(msg2.ValidateBasic(), "offer coin 99denom1 is smaller than the min amount 100: invalid request")

	// Too small order amount.
	msg2 = types.NewMsgMarketOrder(
		appID1, s.addr(1), pair.Id, types.OrderDirectionSell, utils.ParseCoin("100denom1"),
		"denom2", sdkmath.NewInt(99), 0)
	s.Require().EqualError(msg2.ValidateBasic(), "order amount 99 is smaller than the min amount 100: invalid request")

	p := utils.ParseDec("0.0001")
	pair.LastPrice = &p
	s.keeper.SetPair(s.ctx, pair)

	// Too small orders.
	msg2 = types.NewMsgMarketOrder(
		appID1, s.addr(1), pair.Id, types.OrderDirectionBuy, utils.ParseCoin("100denom2"),
		"denom1", sdkmath.NewInt(909090), 0)
	s.Require().NoError(msg2.ValidateBasic())
	_, err = s.keeper.MarketOrder(s.ctx, msg2)
	s.Require().ErrorIs(err, types.ErrTooSmallOrder)

	msg2 = types.NewMsgMarketOrder(
		appID1, s.addr(1), pair.Id, types.OrderDirectionSell, utils.ParseCoin("1003denom1"),
		"denom2", sdkmath.NewInt(1000), 0)
	s.Require().NoError(msg2.ValidateBasic())
	_, err = s.keeper.MarketOrder(s.ctx, msg2)
	s.Require().ErrorIs(err, types.ErrTooSmallOrder)
}

func (s *KeeperTestSuite) TestMarketOrder() {
	creator := s.addr(1)
	trader := s.addr(2)
	dummyAcc := s.addr(1234567890)
	// escrow := s.addr(3)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset2.Denom)

	params, err := s.keeper.GetGenericParams(s.ctx, appID1)
	s.Require().NoError(err)

	testCases := []struct {
		Name         string
		Msg          types.MsgMarketOrder
		FundRequired bool
		ExpErr       error
		ExpResp      *types.Order
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
			FundRequired: true,
			ExpErr:       errorsmod.Wrap(errorsmod.Wrapf(types.ErrInvalidAppID, "app id %d not found", 69), "params retreval failed"),
			ExpResp:      &types.Order{},
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
			FundRequired: true,
			ExpErr:       errorsmod.Wrapf(types.ErrTooLongOrderLifespan, "%s is longer than %s", time.Hour*48, params.MaxOrderLifespan),
			ExpResp:      &types.Order{},
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
			FundRequired: true,
			ExpErr:       errorsmod.Wrapf(sdkerrors.ErrNotFound, "pair %d not found", 69),
			ExpResp:      &types.Order{},
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
			FundRequired: true,
			ExpErr:       types.ErrNoLastPrice,
			ExpResp:      &types.Order{},
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
			FundRequired: true,
			ExpErr:       errorsmod.Wrapf(types.ErrWrongPair, "denom pair (%s, %s) != (%s, %s)", asset2.Denom, asset1.Denom, pair.BaseCoinDenom, pair.QuoteCoinDenom),
			ExpResp:      &types.Order{},
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
			FundRequired: true,
			ExpErr:       errorsmod.Wrapf(types.ErrWrongPair, "denom pair (%s, %s) != (%s, %s)", asset2.Denom, asset1.Denom, pair.BaseCoinDenom, pair.QuoteCoinDenom),
			ExpResp:      &types.Order{},
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
			FundRequired: true,
			ExpErr:       errorsmod.Wrapf(types.ErrInsufficientOfferCoin, "10000000uasset2 is smaller than 11033000uasset2"),
			ExpResp:      &types.Order{},
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
			FundRequired: true,
			ExpErr:       errorsmod.Wrapf(types.ErrInsufficientOfferCoin, "10000000uasset1 is smaller than 10030000uasset1"),
			ExpResp:      &types.Order{},
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
			FundRequired: true,
			ExpErr:       types.ErrTooSmallOrder,
			ExpResp:      &types.Order{},
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
			FundRequired: true,
			ExpErr:       types.ErrTooSmallOrder,
			ExpResp:      &types.Order{},
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
			FundRequired: false,
			ExpErr:       errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds, "0uasset2 is smaller than 1103300uasset2"),
			ExpResp:      &types.Order{},
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
			FundRequired: true,
			ExpErr:       nil,
			ExpResp: &types.Order{
				Id:                 3,
				PairId:             1,
				MsgHeight:          0,
				Orderer:            trader.String(),
				Direction:          types.OrderDirectionBuy,
				OfferCoin:          utils.ParseCoin("1100000uasset2"),
				RemainingOfferCoin: utils.ParseCoin("1100000uasset2"),
				ReceivedCoin:       utils.ParseCoin("0uasset1"),
				Price:              sdkmath.LegacyMustNewDecFromStr("1.1"),
				Amount:             newInt(1000000),
				OpenAmount:         newInt(1000000),
				BatchId:            2,
				ExpireAt:           s.ctx.BlockTime().Add(time.Second * 10),
				Status:             types.OrderStatusNotExecuted,
				AppId:              appID1,
				Type:               types.OrderTypeMarket,
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
			FundRequired: true,
			ExpErr:       nil,
			ExpResp: &types.Order{
				Id:                 4,
				PairId:             1,
				MsgHeight:          0,
				Orderer:            trader.String(),
				Direction:          types.OrderDirectionSell,
				OfferCoin:          utils.ParseCoin("1000000uasset1"),
				RemainingOfferCoin: utils.ParseCoin("1000000uasset1"),
				ReceivedCoin:       utils.ParseCoin("0uasset2"),
				Price:              sdkmath.LegacyMustNewDecFromStr("0.9"),
				Amount:             newInt(1000000),
				OpenAmount:         newInt(1000000),
				BatchId:            2,
				ExpireAt:           s.ctx.BlockTime().Add(time.Second * 10),
				Status:             types.OrderStatusNotExecuted,
				AppId:              appID1,
				Type:               types.OrderTypeMarket,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.Name, func() {
			if tc.FundRequired {
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
				s.LimitOrder(tc.Msg.AppId, creator, tc.Msg.PairId, types.OrderDirectionBuy, utils.ParseDec("1"), sdkmath.NewInt(10000), 0)
				s.LimitOrder(tc.Msg.AppId, creator, tc.Msg.PairId, types.OrderDirectionSell, utils.ParseDec("1"), sdkmath.NewInt(10000), 0)
				s.nextBlock()
			}
			s.sendCoins(tc.Msg.GetOrderer(), dummyAcc, s.getBalances(tc.Msg.GetOrderer()))
		})
	}
}

func (s *KeeperTestSuite) TestMarketOrderTwo() {
	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uasset2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)

	// When there is no last price in the pair, only limit orders can be made.
	// These two orders will be matched.
	s.LimitOrder(appID1, s.addr(1), pair.Id, types.OrderDirectionBuy, utils.ParseDec("1.0"), sdkmath.NewInt(10000), 0)
	s.LimitOrder(appID1, s.addr(2), pair.Id, types.OrderDirectionSell, utils.ParseDec("1.0"), sdkmath.NewInt(10000), 0)
	s.nextBlock()

	// Now users can make market orders.
	// In this case, addr(3) user's order takes higher priority than addr(4) user's,
	// because market buy orders have 10% higher price than the last price(1.0).
	s.MarketOrder(appID1, s.addr(3), pair.Id, types.OrderDirectionBuy, sdkmath.NewInt(10000), 0)
	s.LimitOrder(appID1, s.addr(4), pair.Id, types.OrderDirectionBuy, utils.ParseDec("1.08"), sdkmath.NewInt(10000), 0)
	s.LimitOrder(appID1, s.addr(5), pair.Id, types.OrderDirectionSell, utils.ParseDec("1.07"), sdkmath.NewInt(10000), 0)
	s.nextBlock()

	// Check the result.
	s.Require().True(utils.ParseCoin("10000uasset1").IsEqual(s.getBalance(s.addr(3), "uasset1")))
	s.Require().True(utils.ParseCoins("10832uasset2").Equal(s.getBalances(s.addr(4))))
}

func (s *KeeperTestSuite) TestMarketOrderRefund() {
	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 2000000)

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
	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)

	s.Require().Nil(pair.LastPrice)
	offerCoin := utils.ParseCoin("10000denom2")
	s.fundAddr(s.addr(1), sdk.NewCoins(offerCoin))
	msg := types.NewMsgMarketOrder(
		appID1, s.addr(1), pair.Id, types.OrderDirectionBuy, offerCoin, "denom1", sdkmath.NewInt(10000), 0)
	_, err := s.keeper.MarketOrder(s.ctx, msg)
	s.Require().ErrorIs(err, types.ErrNoLastPrice)
}

func (s *KeeperTestSuite) TestSingleOrderNoMatch() {
	k, ctx := s.keeper, s.ctx

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)

	order := s.LimitOrder(appID1, s.addr(1), pair.Id, types.OrderDirectionBuy, utils.ParseDec("1.0"), sdkmath.NewInt(1000000), 10*time.Second)
	// Execute matching
	liquidity.EndBlocker(ctx, k, s.app.AssetKeeper)

	order, found := k.GetOrder(ctx, appID1, order.PairId, order.Id)
	s.Require().True(found)
	s.Require().Equal(types.OrderStatusNotMatched, order.Status)

	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(10 * time.Second))
	// Expire the order, here BeginBlocker is not called to check
	// the request's changed status
	liquidity.EndBlocker(ctx, k, s.app.AssetKeeper)

	order, _ = k.GetOrder(ctx, appID1, order.PairId, order.Id)
	s.Require().Equal(types.OrderStatusExpired, order.Status)

	s.Require().True(utils.ParseCoins("1003000denom2").Equal(s.getBalances(s.addr(1))))
}

func (s *KeeperTestSuite) TestTwoOrderExactMatch() {
	k, ctx := s.keeper, s.ctx

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)

	req1 := s.LimitOrder(appID1, s.addr(1), pair.Id, types.OrderDirectionBuy, utils.ParseDec("1.0"), newInt(10000), time.Hour)
	req2 := s.LimitOrder(appID1, s.addr(2), pair.Id, types.OrderDirectionSell, utils.ParseDec("1.0"), newInt(10000), time.Hour)
	liquidity.EndBlocker(ctx, k, s.app.AssetKeeper)

	req1, _ = k.GetOrder(ctx, appID1, req1.PairId, req1.Id)
	s.Require().Equal(types.OrderStatusCompleted, req1.Status)
	req2, _ = k.GetOrder(ctx, appID1, req2.PairId, req2.Id)
	s.Require().Equal(types.OrderStatusCompleted, req2.Status)

	s.Require().True(utils.ParseCoins("10000denom1").Equal(s.getBalances(s.addr(1))))
	s.Require().True(utils.ParseCoins("10000denom2").Equal(s.getBalances(s.addr(2))))

	pair, _ = k.GetPair(ctx, appID1, pair.Id)
	s.Require().NotNil(pair.LastPrice)
	s.Require().True(utils.ParseDec("1.0").Equal(*pair.LastPrice))
}

func (s *KeeperTestSuite) TestPartialMatch() {
	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)

	order := s.LimitOrder(appID1, s.addr(1), pair.Id, types.OrderDirectionBuy, utils.ParseDec("1.0"), sdkmath.NewInt(10000), time.Hour)
	s.LimitOrder(appID1, s.addr(2), pair.Id, types.OrderDirectionSell, utils.ParseDec("1.0"), sdkmath.NewInt(5000), 0)
	s.nextBlock()

	order, found := s.keeper.GetOrder(s.ctx, appID1, order.PairId, order.Id)
	s.Require().True(found)
	s.Require().Equal(types.OrderStatusPartiallyMatched, order.Status)
	s.Require().True(utils.ParseCoin("5000denom2").IsEqual(order.RemainingOfferCoin))
	s.Require().True(utils.ParseCoin("5000denom1").IsEqual(order.ReceivedCoin))
	s.Require().True(sdkmath.NewInt(5000).Equal(order.OpenAmount))

	s.MarketOrder(appID1, s.addr(3), pair.Id, types.OrderDirectionSell, sdkmath.NewInt(5000), 0)
	s.nextBlock()

	// Now completely matched.
	_, found = s.keeper.GetOrder(s.ctx, appID1, order.PairId, order.Id)
	s.Require().False(found)
}

func (s *KeeperTestSuite) TestMatchWithLowPricePool() {
	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)

	// Create a pool with very low price.
	s.CreateNewLiquidityPool(appID1, pair.Id, s.addr(0), "1000000000000000000000denom1,1000000denom2")
	order := s.LimitOrder(appID1, s.addr(1), pair.Id, types.OrderDirectionBuy, utils.ParseDec("0.000000000001000000"), sdkmath.NewInt(100000000000000000), 10*time.Second)
	liquidity.EndBlocker(s.ctx, s.keeper, s.app.AssetKeeper)
	order, found := s.keeper.GetOrder(s.ctx, appID1, order.PairId, order.Id)
	s.Require().True(found)
	s.Require().Equal(types.OrderStatusNotMatched, order.Status)
}

func (s *KeeperTestSuite) TestMMOrder() {
	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)
	pair.LastPrice = utils.ParseDecP("1.0")
	s.keeper.SetPair(s.ctx, pair)

	params, err := s.keeper.GetGenericParams(s.ctx, appID1)
	s.Require().NoError(err)

	orders := s.MarketMakingOrder(
		s.addr(1), appID1, pair.Id,
		utils.ParseDec("1.1"), utils.ParseDec("1.03"), sdkmath.NewInt(1000_000000),
		utils.ParseDec("0.97"), utils.ParseDec("0.9"), sdkmath.NewInt(1000_000000),
		10*time.Second, true)
	maxNumTicks := int(params.MaxNumMarketMakingOrderTicks)
	s.Require().Len(orders, 2*maxNumTicks)

	pair, _ = s.keeper.GetPair(s.ctx, appID1, pair.Id)
	s.Require().EqualValues(2*maxNumTicks, pair.LastOrderId)

	// failed cancel on same batch
	cancelMsg := types.MsgCancelMMOrder{
		Orderer: s.addr(1).String(),
		PairId:  pair.Id,
		AppId:   appID1,
	}
	_, err = s.keeper.CancelMMOrder(s.ctx, &cancelMsg)
	s.Require().ErrorIs(err, types.ErrSameBatch)

	// successful cancel on next block
	s.nextBlock()
	ids, err := s.keeper.CancelMMOrder(s.ctx, &cancelMsg)
	s.Require().NoError(err)
	for i := range ids {
		s.Require().Equal(ids[i], orders[i].Id)
	}
	_, found := s.keeper.GetMMOrderIndex(s.ctx, s.addr(1), appID1, pair.Id)
	s.Require().False(found)

}

func (s *KeeperTestSuite) TestMMOrderCancelPreviousOrders() {
	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)
	pair.LastPrice = utils.ParseDecP("1.0")
	s.keeper.SetPair(s.ctx, pair)

	s.MarketMakingOrder(
		s.addr(1), appID1, pair.Id,
		utils.ParseDec("1.1"), utils.ParseDec("1.03"), sdkmath.NewInt(1000_000000),
		utils.ParseDec("0.97"), utils.ParseDec("0.9"), sdkmath.NewInt(1000_000000),
		10*time.Second, true)

	// Cannot place MM orders again because it's not allowed to cancel previous
	// orders within same batch.
	s.fundAddr(s.addr(1), utils.ParseCoins("1000000000denom1,1000000000denom2"))
	_, err := s.keeper.MMOrder(s.ctx, types.NewMsgMMOrder(
		appID1, s.addr(1), pair.Id,
		utils.ParseDec("1.1"), utils.ParseDec("1.03"), sdkmath.NewInt(1000_000000),
		utils.ParseDec("0.97"), utils.ParseDec("0.9"), sdkmath.NewInt(1000_000000),
		10*time.Second))
	s.Require().ErrorIs(err, types.ErrSameBatch)

	// Now it's OK to cancel previous orders and place new orders.
	s.nextBlock()
	s.MarketMakingOrder(
		s.addr(1), appID1, pair.Id,
		utils.ParseDec("1.1"), utils.ParseDec("1.03"), sdkmath.NewInt(1000_000000),
		utils.ParseDec("0.97"), utils.ParseDec("0.9"), sdkmath.NewInt(1000_000000),
		10*time.Second, true)

	params, err := s.keeper.GetGenericParams(s.ctx, appID1)
	s.Require().NoError(err)
	orders := s.keeper.GetAllOrders(s.ctx, appID1)
	maxNumTicks := int(params.MaxNumMarketMakingOrderTicks)
	s.Require().Len(orders, 2*2*maxNumTicks) // canceled previous orders + new orders

	s.nextBlock()

	orders = s.keeper.GetAllOrders(s.ctx, appID1)
	s.Require().Len(orders, 2*maxNumTicks) // new orders
	// Check order ids.
	for _, order := range orders {
		s.Require().EqualValues(2, order.BatchId)
		s.Require().GreaterOrEqual(order.Id, uint64(2*maxNumTicks+1))
	}
}

func (s *KeeperTestSuite) TestCancelOrder() {
	creator := s.addr(0)
	dummy := s.addr(1)
	trader := s.addr(2)
	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset2.Denom)

	cancelledOrder := s.LimitOrder(appID1, creator, pair.Id, types.OrderDirectionSell, utils.ParseDec("1.1"), newInt(10000000), time.Second*10)
	s.nextBlock()
	msg := types.NewMsgCancelOrder(appID1, creator, pair.Id, cancelledOrder.Id)
	err := s.keeper.CancelOrder(s.ctx, msg)
	s.Require().NoError(err)

	order := s.LimitOrder(appID1, trader, pair.Id, types.OrderDirectionSell, utils.ParseDec("1.1"), newInt(1000000), time.Second*10)

	testCases := []struct {
		Name   string
		Msg    types.MsgCancelOrder
		ExpErr error
	}{
		{
			Name:   "error app id invalid",
			Msg:    *types.NewMsgCancelOrder(69, creator, pair.Id, order.Id),
			ExpErr: errorsmod.Wrapf(types.ErrInvalidAppID, "app id %d not found", 69),
		},
		{
			Name:   "error order id invalid",
			Msg:    *types.NewMsgCancelOrder(appID1, creator, pair.Id, 69),
			ExpErr: errorsmod.Wrapf(sdkerrors.ErrNotFound, "order %d not found in pair %d", 69, pair.Id),
		},
		{
			Name:   "error invalid orderer",
			Msg:    *types.NewMsgCancelOrder(appID1, dummy, pair.Id, order.Id),
			ExpErr: errorsmod.Wrap(sdkerrors.ErrUnauthorized, "mismatching orderer"),
		},
		{
			Name:   "error order already cancelled",
			Msg:    *types.NewMsgCancelOrder(appID1, creator, pair.Id, cancelledOrder.Id),
			ExpErr: types.ErrAlreadyCanceled,
		},
		{
			Name:   "error same batch",
			Msg:    *types.NewMsgCancelOrder(appID1, trader, pair.Id, order.Id),
			ExpErr: types.ErrSameBatch,
		},
		{
			Name:   "success valid case",
			Msg:    *types.NewMsgCancelOrder(appID1, trader, pair.Id, order.Id),
			ExpErr: nil,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			if tc.ExpErr == nil {
				// triggering new batch by going to next block
				s.nextBlock()
			}
			err := s.keeper.CancelOrder(s.ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(tc.ExpErr, err.Error())
			} else {
				s.Require().NoError(err)

				order, found := s.keeper.GetOrder(s.ctx, tc.Msg.AppId, tc.Msg.PairId, tc.Msg.OrderId)
				s.Require().True(found)
				s.Require().Equal(types.OrderStatusCanceled, order.Status)

				s.Require().True(utils.ParseCoins("1003000denom1").Equal(s.getBalances(tc.Msg.GetOrderer())))

				s.nextBlock()
				_, found = s.keeper.GetOrder(s.ctx, tc.Msg.AppId, tc.Msg.PairId, tc.Msg.OrderId)
				s.Require().False(found)

			}
		})
	}
}

func (s *KeeperTestSuite) TestCancelOrderTwo() {
	k, ctx := s.keeper, s.ctx

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)

	order := s.LimitOrder(appID1, s.addr(1), pair.Id, types.OrderDirectionBuy, utils.ParseDec("1.0"), newInt(10000), types.DefaultMaxOrderLifespan)

	// Cannot cancel an order within a same batch
	err := k.CancelOrder(ctx, types.NewMsgCancelOrder(appID1, s.addr(1), order.PairId, order.Id))
	s.Require().ErrorIs(err, types.ErrSameBatch)

	s.nextBlock()

	// Now an order can be canceled
	err = k.CancelOrder(ctx, types.NewMsgCancelOrder(appID1, s.addr(1), order.PairId, order.Id))
	s.Require().NoError(err)

	order, found := k.GetOrder(ctx, appID1, order.PairId, order.Id)
	s.Require().True(found)
	s.Require().Equal(types.OrderStatusCanceled, order.Status)

	// Coins are refunded
	s.Require().True(utils.ParseCoins("10030denom2").Equal(s.getBalances(s.addr(1))))

	s.nextBlock()

	// Order is deleted
	_, found = k.GetOrder(ctx, appID1, order.PairId, order.Id)
	s.Require().False(found)
}

func (s *KeeperTestSuite) TestCancelAllOrders() {
	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 2000000)
	asset3 := s.CreateNewAsset("ASSETHREE", "denom3", 3000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)

	order := s.LimitOrder(appID1, s.addr(1), pair.Id, types.OrderDirectionBuy, utils.ParseDec("1.0"), sdkmath.NewInt(10000), time.Hour)
	msg := types.NewMsgCancelAllOrders(appID1, s.addr(1), nil)
	s.keeper.CancelAllOrders(s.ctx, msg) // CancelAllOrders doesn't cancel orders within in same batch
	s.nextBlock()

	// The order is still alive.
	_, found := s.keeper.GetOrder(s.ctx, appID1, order.PairId, order.Id)
	s.Require().True(found)

	s.keeper.CancelAllOrders(s.ctx, msg) // This time, it cancels the order.
	order, found = s.keeper.GetOrder(s.ctx, appID1, order.PairId, order.Id)
	// Canceling an order doesn't delete the order immediately.
	s.Require().True(found)
	// Instead, the order becomes canceled.
	s.Require().Equal(types.OrderStatusCanceled, order.Status)

	// The order won't be matched with this market order, since the order is
	// already canceled.
	s.LimitOrder(appID1, s.addr(3), pair.Id, types.OrderDirectionSell, utils.ParseDec("1.0"), sdkmath.NewInt(10000), 0)
	s.nextBlock()
	s.Require().True(utils.ParseCoins("10030denom2").Equal(s.getBalances(s.addr(1))))

	pair2 := s.CreateNewLiquidityPair(appID1, s.addr(0), asset2.Denom, asset3.Denom)
	s.LimitOrder(appID1, s.addr(2), pair.Id, types.OrderDirectionBuy, utils.ParseDec("1.0"), sdkmath.NewInt(10000), time.Hour)
	s.LimitOrder(appID1, s.addr(2), pair.Id, types.OrderDirectionSell, utils.ParseDec("1.5"), sdkmath.NewInt(10000), time.Hour)
	s.LimitOrder(appID1, s.addr(2), pair2.Id, types.OrderDirectionSell, utils.ParseDec("1.0"), sdkmath.NewInt(10000), time.Hour)
	s.nextBlock()

	msg = types.NewMsgCancelAllOrders(appID1, s.addr(2), []uint64{pair.Id})
	// CancelAllOrders can cancel orders in specific pairs.
	s.keeper.CancelAllOrders(s.ctx, msg)
	// Coins from first two orders are refunded, but not from the last order.
	s.Require().True(utils.ParseCoins("10030denom2,10030denom1").Equal(s.getBalances(s.addr(2))))
}

func (s *KeeperTestSuite) TestSwapFeeCollectionWithoutPool() {
	creator := s.addr(0)
	buyer := s.addr(1)
	seller := s.addr(2)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset2.Denom)

	buyOrder := s.LimitOrder(appID1, buyer, pair.Id, types.OrderDirectionBuy, utils.ParseDec("1"), newInt(52000000), time.Second*10)
	sellOrder := s.LimitOrder(appID1, seller, pair.Id, types.OrderDirectionSell, utils.ParseDec("1"), newInt(52000000), time.Second*10)

	s.nextBlock()

	buyOrder, found := s.keeper.GetOrder(s.ctx, appID1, pair.Id, buyOrder.Id)
	s.Require().False(found)

	sellOrder, found = s.keeper.GetOrder(s.ctx, appID1, pair.Id, sellOrder.Id)
	s.Require().False(found)

	collectedSwapFee := s.getBalances(pair.GetSwapFeeCollectorAddress())
	s.Require().True(utils.ParseCoins("156000denom2,156000denom1").Equal(collectedSwapFee))
}

func (s *KeeperTestSuite) TestSwapFeeCollectionWithPool() {
	creator := s.addr(0)
	buyer := s.addr(1)
	seller := s.addr(2)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset2.Denom)
	_ = s.CreateNewLiquidityPool(appID1, pair.Id, creator, "100000000000denom1,100000000000denom2")

	buyOrder1 := s.LimitOrder(appID1, buyer, pair.Id, types.OrderDirectionBuy, utils.ParseDec("1.001"), newInt(1000000), time.Second*10)
	buyOrder2 := s.LimitOrder(appID1, buyer, pair.Id, types.OrderDirectionBuy, utils.ParseDec("1.001"), newInt(1000000), time.Second*10)

	s.nextBlock()

	buyOrder1, found := s.keeper.GetOrder(s.ctx, appID1, pair.Id, buyOrder1.Id)
	s.Require().False(found)

	buyOrder2, found = s.keeper.GetOrder(s.ctx, appID1, pair.Id, buyOrder2.Id)
	s.Require().False(found)

	sellOrder1 := s.LimitOrder(appID1, seller, pair.Id, types.OrderDirectionSell, utils.ParseDec("0.99"), newInt(1000000), time.Second*10)
	sellOrder2 := s.LimitOrder(appID1, seller, pair.Id, types.OrderDirectionSell, utils.ParseDec("0.99"), newInt(1000000), time.Second*10)

	s.nextBlock()

	sellOrder1, found = s.keeper.GetOrder(s.ctx, appID1, pair.Id, sellOrder1.Id)
	s.Require().False(found)

	sellOrder2, found = s.keeper.GetOrder(s.ctx, appID1, pair.Id, sellOrder2.Id)
	s.Require().False(found)

	collectedSwapFee := s.getBalances(pair.GetSwapFeeCollectorAddress())
	s.Require().True(utils.ParseCoins("6000denom2,6000denom1").Equal(collectedSwapFee))
}

func (s *KeeperTestSuite) TestSwapFeeCollectionMarketOrder() {
	creator := s.addr(0)
	buyer := s.addr(1)
	trader1 := s.addr(2)
	trader2 := s.addr(3)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset2.Denom)
	_ = s.CreateNewLiquidityPool(appID1, pair.Id, creator, "100000000000denom1,100000000000denom2")

	buyOrder1 := s.LimitOrder(appID1, buyer, pair.Id, types.OrderDirectionBuy, utils.ParseDec("1.001"), newInt(1000000), time.Second*10)
	s.nextBlock()
	buyOrder1, found := s.keeper.GetOrder(s.ctx, appID1, pair.Id, buyOrder1.Id)
	s.Require().False(found)

	sellMarketOrder := s.MarketOrder(appID1, trader1, pair.Id, types.OrderDirectionSell, newInt(100000000), time.Second*10)
	s.nextBlock()
	_, found = s.keeper.GetOrder(s.ctx, appID1, pair.Id, sellMarketOrder.Id)
	s.Require().False(found)
	s.Require().True(utils.ParseCoins("99902053denom2").Equal(s.getBalances(trader1)))

	buyMarketOrder := s.MarketOrder(appID1, trader2, pair.Id, types.OrderDirectionBuy, newInt(100000000), time.Second*10)
	s.nextBlock()
	_, found = s.keeper.GetOrder(s.ctx, appID1, pair.Id, buyMarketOrder.Id)
	s.Require().False(found)
	s.Require().True(utils.ParseCoins("100000000denom1,9908395denom2").Equal(s.getBalances(trader2)))

	accumulatedSwapFee := s.getBalances(pair.GetSwapFeeCollectorAddress())
	s.Require().True(utils.ParseCoins("300000denom1,302707denom2").Equal(accumulatedSwapFee))

	s.nextBlock()
}

func (s *KeeperTestSuite) TestAccumulatedSwapFeeConversion() {
	creator := s.addr(0)
	trader1 := s.addr(1)
	trader2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "ucmdx", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uharbor", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, creator, asset1.Denom, asset2.Denom)

	// accumulateSwapFee in swapeecollector address by placing orders
	s.LimitOrder(appID1, trader1, pair.Id, types.OrderDirectionBuy, utils.ParseDec("1"), newInt(52000000), time.Second*10)
	s.LimitOrder(appID1, trader2, pair.Id, types.OrderDirectionSell, utils.ParseDec("1"), newInt(52000000), time.Second*10)
	s.LimitOrder(appID1, trader1, pair.Id, types.OrderDirectionBuy, utils.ParseDec("1"), newInt(52000000), time.Second*10)
	s.LimitOrder(appID1, trader2, pair.Id, types.OrderDirectionSell, utils.ParseDec("1"), newInt(52000000), time.Second*10)
	s.LimitOrder(appID1, trader1, pair.Id, types.OrderDirectionBuy, utils.ParseDec("1"), newInt(52000000), time.Second*10)
	s.LimitOrder(appID1, trader2, pair.Id, types.OrderDirectionSell, utils.ParseDec("1"), newInt(52000000), time.Second*10)
	s.LimitOrder(appID1, trader1, pair.Id, types.OrderDirectionBuy, utils.ParseDec("1"), newInt(52000000), time.Second*10)
	s.LimitOrder(appID1, trader2, pair.Id, types.OrderDirectionSell, utils.ParseDec("1"), newInt(52000000), time.Second*10)
	s.LimitOrder(appID1, trader1, pair.Id, types.OrderDirectionBuy, utils.ParseDec("1"), newInt(52000000), time.Second*10)
	s.LimitOrder(appID1, trader2, pair.Id, types.OrderDirectionSell, utils.ParseDec("1"), newInt(52000000), time.Second*10)
	s.LimitOrder(appID1, trader1, pair.Id, types.OrderDirectionBuy, utils.ParseDec("1"), newInt(52000000), time.Second*10)
	s.LimitOrder(appID1, trader2, pair.Id, types.OrderDirectionSell, utils.ParseDec("1"), newInt(52000000), time.Second*10)

	// execute orders and try to convert, conversion will not take since there are no pool
	s.nextBlock()
	accumulatedSwapFee := s.getBalances(pair.GetSwapFeeCollectorAddress())
	s.Require().True(utils.ParseCoins("936000ucmdx,936000uharbor").Equal(accumulatedSwapFee))

	// try to convert again, conversion will not take place since there are no pool
	s.nextBlock()
	accumulatedSwapFee = s.getBalances(pair.GetSwapFeeCollectorAddress())
	s.Require().True(utils.ParseCoins("936000ucmdx,936000uharbor").Equal(accumulatedSwapFee))

	// now create pool, so that token conversion can go through this
	_ = s.CreateNewLiquidityPool(appID1, pair.Id, creator, "1000000000000ucmdx,1000000000000uharbor")

	// NOTE
	// Eery conversion internally is an limit order based on the pool path

	// retry to convert by moving to next block, this time conversin/swap will take place since pool exists
	s.nextBlock()
	accumulatedSwapFee = s.getBalances(pair.GetSwapFeeCollectorAddress())
	// here order is placed for swap, hence harbor tokens are reduced and this will get executed in next block
	s.Require().True(utils.ParseCoins("936000ucmdx,5643uharbor").Equal(accumulatedSwapFee))

	// now execute the order placed in above block, swap order for 9 uharbor placed again in next block
	s.nextBlock()
	accumulatedSwapFee = s.getBalances(pair.GetSwapFeeCollectorAddress())
	s.Require().True(utils.ParseCoins("1779250ucmdx,556uharbor").Equal(accumulatedSwapFee))

	// now execute the order placed in above block, this block will execute the order for 9 harbor placed above
	s.nextBlock()
	accumulatedSwapFee = s.getBalances(pair.GetSwapFeeCollectorAddress())
	s.Require().True(utils.ParseCoins("1862727ucmdx,54uharbor").Equal(accumulatedSwapFee))

	// now execute the order placed in above block, here 1uharbor is refunded back since it is very small amount for swap order.
	// here all harbor tokens are converted into cmdx, since cmdx is the default distribution token for rewards
	s.nextBlock()
	accumulatedSwapFee = s.getBalances(pair.GetSwapFeeCollectorAddress())
	s.Require().True(utils.ParseCoins("1870997ucmdx,5uharbor").Equal(accumulatedSwapFee))
}

func (s *KeeperTestSuite) TestConvertAccumulatedSwapFeesWithSwapDistrToken_1() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")
	asset1 := s.CreateNewAsset("ASSETONE", "ucmdx", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uharbor", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pair.LastPrice = utils.ParseDecP("1.01")
	s.keeper.SetPair(s.ctx, pair)
	_ = s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "100000000ucmdx,100000000uharbor")
	// the given deposit amount is not the actuall deposit amout, it will be decided by the protocal based on the min,max and initial prices of ranged pools
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "100000000ucmdx,100000000uharbor", sdkmath.LegacyMustNewDecFromStr("0.95"), sdkmath.LegacyMustNewDecFromStr("1.05"), sdkmath.LegacyMustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "100000000ucmdx,100000000uharbor", sdkmath.LegacyMustNewDecFromStr("0.90"), sdkmath.LegacyMustNewDecFromStr("1.1"), sdkmath.LegacyMustNewDecFromStr("1"))
	_ = s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "100000000ucmdx,100000000uharbor", sdkmath.LegacyMustNewDecFromStr("0.99"), sdkmath.LegacyMustNewDecFromStr("1.01"), sdkmath.LegacyMustNewDecFromStr("1"))

	s.MarketOrder(appID1, s.addr(2), pair.Id, types.OrderDirectionBuy, sdkmath.NewInt(30_000000), 0)
	s.MarketOrder(appID1, s.addr(3), pair.Id, types.OrderDirectionBuy, sdkmath.NewInt(54_000000), 0)
	s.MarketOrder(appID1, s.addr(4), pair.Id, types.OrderDirectionBuy, sdkmath.NewInt(17_000000), 0)

	s.MarketOrder(appID1, s.addr(5), pair.Id, types.OrderDirectionSell, sdkmath.NewInt(12_000000), 0)
	s.MarketOrder(appID1, s.addr(6), pair.Id, types.OrderDirectionSell, sdkmath.NewInt(19_000000), 0)
	s.MarketOrder(appID1, s.addr(7), pair.Id, types.OrderDirectionSell, sdkmath.NewInt(23_000000), 0)

	_, err := s.keeper.GetGenericParams(s.ctx, appID1)
	s.Require().NoError(err)

	s.ctx = s.ctx.WithBlockHeight(15) // to avoid swapfee conversion via beign blocker which occurs at block height 150
	s.nextBlock()
	orders := s.keeper.GetAllOrders(s.ctx, appID1)
	s.Require().Empty(orders)

	accumulatedSwapFee := s.getBalances(pair.GetSwapFeeCollectorAddress())
	s.Require().True(coinsEq(utils.ParseCoins("162000ucmdx,306030uharbor"), accumulatedSwapFee))

	s.keeper.ConvertAccumulatedSwapFeesWithSwapDistrToken(s.ctx, appID1)
	// swap order is placed in the above block, it will get executed in the next block
	s.Require().True(coinsEq(utils.ParseCoins("162000ucmdx,10uharbor"), s.getBalances(pair.GetSwapFeeCollectorAddress())))
	s.nextBlock()
	// previous order was completed and rest of the tokens are returned back again
	s.Require().True(coinsEq(utils.ParseCoins("436622ucmdx,28661uharbor"), s.getBalances(pair.GetSwapFeeCollectorAddress())))
}

func (s *KeeperTestSuite) TestConvertAccumulatedSwapFeesWithSwapDistrToken_2() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appone")
	asset1 := s.CreateNewAsset("ASSETONE", "ucmdx", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uharbor", 1000000)
	asset3 := s.CreateNewAsset("ASSETTHREE", "uatom", 1000000)
	asset4 := s.CreateNewAsset("ASSETFOUR", "stake", 1000000)

	pair1 := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pair2 := s.CreateNewLiquidityPair(appID1, addr1, asset2.Denom, asset3.Denom)
	pair3 := s.CreateNewLiquidityPair(appID1, addr1, asset3.Denom, asset4.Denom)
	pair4 := s.CreateNewLiquidityPair(appID1, addr1, asset4.Denom, asset1.Denom)

	pair1.LastPrice = utils.ParseDecP("1.01")
	pair2.LastPrice = utils.ParseDecP("1.01")
	pair3.LastPrice = utils.ParseDecP("1.01")
	pair4.LastPrice = utils.ParseDecP("1.01")
	s.keeper.SetPair(s.ctx, pair1)
	s.keeper.SetPair(s.ctx, pair2)
	s.keeper.SetPair(s.ctx, pair3)
	s.keeper.SetPair(s.ctx, pair4)

	_ = s.CreateNewLiquidityPool(appID1, pair1.Id, addr1, "100000000ucmdx,100000000uharbor")
	_ = s.CreateNewLiquidityPool(appID1, pair2.Id, addr1, "100000000uharbor,100000000uatom")
	_ = s.CreateNewLiquidityPool(appID1, pair3.Id, addr1, "100000000uatom,100000000stake")
	_ = s.CreateNewLiquidityPool(appID1, pair4.Id, addr1, "100000000stake,100000000ucmdx")
	// the given deposit amount is not the actuall deposit amout, it will be decided by the protocal based on the min,max and initial prices of ranged pools

	s.MarketOrder(appID1, s.addr(2), pair1.Id, types.OrderDirectionBuy, sdkmath.NewInt(30_000000), 0)
	s.MarketOrder(appID1, s.addr(3), pair2.Id, types.OrderDirectionBuy, sdkmath.NewInt(54_000000), 0)
	s.MarketOrder(appID1, s.addr(4), pair3.Id, types.OrderDirectionBuy, sdkmath.NewInt(17_000000), 0)
	s.MarketOrder(appID1, s.addr(4), pair4.Id, types.OrderDirectionBuy, sdkmath.NewInt(17_000000), 0)

	s.MarketOrder(appID1, s.addr(6), pair1.Id, types.OrderDirectionSell, sdkmath.NewInt(19_000000), 0)
	s.MarketOrder(appID1, s.addr(7), pair2.Id, types.OrderDirectionSell, sdkmath.NewInt(23_000000), 0)
	s.MarketOrder(appID1, s.addr(5), pair3.Id, types.OrderDirectionSell, sdkmath.NewInt(12_000000), 0)
	s.MarketOrder(appID1, s.addr(5), pair4.Id, types.OrderDirectionSell, sdkmath.NewInt(32_000000), 0)

	_, err := s.keeper.GetGenericParams(s.ctx, appID1)
	s.Require().NoError(err)

	s.ctx = s.ctx.WithBlockHeight(15) // to avoid swapfee conversion via beign blocker which occurs at block height 150
	s.nextBlock()
	orders := s.keeper.GetAllOrders(s.ctx, appID1)
	s.Require().Empty(orders)

	s.Require().True(coinsEq(utils.ParseCoins("57000ucmdx,73592uharbor"), s.getBalances(pair1.GetSwapFeeCollectorAddress())))
	s.Require().True(coinsEq(utils.ParseCoins("85712uatom,69000uharbor"), s.getBalances(pair2.GetSwapFeeCollectorAddress())))
	s.Require().True(coinsEq(utils.ParseCoins("52157stake,36000uatom"), s.getBalances(pair3.GetSwapFeeCollectorAddress())))
	s.Require().True(coinsEq(utils.ParseCoins("65880stake,51510ucmdx"), s.getBalances(pair4.GetSwapFeeCollectorAddress())))

	// doing s.nextBlock() or calling the method ConvertAccumulatedSwapFeesWithSwapDistrToken manually is same thing
	// swap order is placed in the above block, it will get executed in the next block
	s.keeper.ConvertAccumulatedSwapFeesWithSwapDistrToken(s.ctx, appID1)
	s.nextBlock()
	s.keeper.ConvertAccumulatedSwapFeesWithSwapDistrToken(s.ctx, appID1)
	s.nextBlock()

	s.Require().True(coinsEq(utils.ParseCoins("122747ucmdx,983uharbor"), s.getBalances(pair1.GetSwapFeeCollectorAddress())))
	s.Require().True(coinsEq(utils.ParseCoins("826uatom,117890ucmdx,13950uharbor"), s.getBalances(pair2.GetSwapFeeCollectorAddress())))
	s.Require().True(coinsEq(utils.ParseCoins("301uatom,70988ucmdx,5363uharbor"), s.getBalances(pair3.GetSwapFeeCollectorAddress())))
	s.Require().True(coinsEq(utils.ParseCoins("6stake,111354ucmdx"), s.getBalances(pair4.GetSwapFeeCollectorAddress())))
}

func (s *KeeperTestSuite) TestPoolPreserveK() {
	r := rand.New(rand.NewSource(0))

	addr1 := s.addr(0)
	addr2 := s.addr(2)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "ucmdx", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "uharbor", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	params, err := s.keeper.GetGenericParams(s.ctx, appID1)
	s.Require().NoError(err)
	for i := 0; i < 10; i++ {
		minPrice := amm.RandomTick(r, utils.ParseDec("0.001"), utils.ParseDec("10.0"), int(params.TickPrecision))
		maxPrice := amm.RandomTick(r, minPrice.Mul(utils.ParseDec("1.01")), utils.ParseDec("100.0"), int(params.TickPrecision))
		initialPrice := amm.RandomTick(r, minPrice, maxPrice, int(params.TickPrecision))
		p := s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr2, "1_000000000000ucmdx,1_000000000000uharbor",
			minPrice, maxPrice, initialPrice,
		)
		s.Require().IsType(types.Pool{}, p)
	}
	pools := s.keeper.GetAllPools(s.ctx, appID1)

	ks := map[uint64]sdkmath.LegacyDec{}
	for _, pool := range pools {
		rx, ry := s.keeper.GetPoolBalances(s.ctx, pool)
		ammPool := pool.AMMPool(rx.Amount, ry.Amount, sdkmath.Int{}).(*amm.RangedPool)
		transX, transY := ammPool.Translation()
		ks[pool.Id] = sdkmath.LegacyNewDec(rx.Amount.Int64()).Add(transX).Mul(sdkmath.LegacyNewDec(ry.Amount.Int64()).Add(transY))
	}

	for i := 0; i < 20; i++ {
		pair, _ = s.keeper.GetPair(s.ctx, appID1, pair.Id)
		for j := 0; j < 50; j++ {
			var price sdkmath.LegacyDec
			if pair.LastPrice == nil {
				price = utils.RandomDec(r, utils.ParseDec("0.001"), utils.ParseDec("100.0"))
			} else {
				price = utils.RandomDec(r, utils.ParseDec("0.91"), utils.ParseDec("1.09")).Mul(*pair.LastPrice)
			}
			amt := utils.RandomInt(r, sdkmath.NewInt(10000), sdkmath.NewInt(1000000))
			lifespan := time.Duration(r.Intn(60)) * time.Second
			if r.Intn(2) == 0 {
				s.LimitOrder(appID1, s.addr(j+2), pair.Id, types.OrderDirectionBuy, price, amt, lifespan)
			} else {
				s.LimitOrder(appID1, s.addr(j+2), pair.Id, types.OrderDirectionBuy, price, amt, lifespan)
			}
		}

		s.nextBlock()
		s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(3 * time.Second))
		s.nextBlock()

		for _, pool := range pools {
			rx, ry := s.keeper.GetPoolBalances(s.ctx, pool)
			ammPool := pool.AMMPool(rx.Amount, ry.Amount, sdkmath.Int{}).(*amm.RangedPool)
			transX, transY := ammPool.Translation()
			k := sdkmath.LegacyNewDec(rx.Amount.Int64()).Add(transX).Mul(sdkmath.LegacyNewDec(ry.Amount.Int64()).Add(transY))
			s.Require().True(k.GTE(ks[pool.Id].Mul(utils.ParseDec("0.99999")))) // there may be a small error
			ks[pool.Id] = k
		}
	}
}

func (s *KeeperTestSuite) TestPoolOrderOverflow() {
	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)
	s.CreateNewLiquidityPool(appID1, pair.Id, s.addr(0), "1000000denom1,10000000000000000000000000denom2")

	s.LimitOrder(appID1, s.addr(1), pair.Id, types.OrderDirectionSell, utils.ParseDec("0.000000000000010000"), sdkmath.NewInt(1e17), 0)
	s.Require().NotPanics(func() {
		liquidity.EndBlocker(s.ctx, s.keeper, s.app.AssetKeeper)
	})
}

func (s *KeeperTestSuite) TestRangedLiquidity() {
	orderPrice := utils.ParseDec("1.05")
	orderAmt := sdkmath.NewInt(100000)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 1000000)
	asset3 := s.CreateNewAsset("ASSETTHREE", "denom3", 1000000)
	asset4 := s.CreateNewAsset("ASSETFOUR", "denom4", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)
	pair.LastPrice = utils.ParseDecP("1.0")
	s.keeper.SetPair(s.ctx, pair)

	s.CreateNewLiquidityPool(appID1, pair.Id, s.addr(1), "1000000denom1,1000000denom2")

	order := s.LimitOrder(appID1, s.addr(2), pair.Id, types.OrderDirectionBuy, orderPrice, orderAmt, 0)
	liquidity.EndBlocker(s.ctx, s.keeper, s.app.AssetKeeper)
	order, _ = s.keeper.GetOrder(s.ctx, appID1, order.PairId, order.Id)
	paid := order.OfferCoin.Sub(order.RemainingOfferCoin).Amount
	received := order.ReceivedCoin.Amount
	s.Require().True(received.LT(orderAmt))
	s.Require().True(sdkmath.LegacyNewDec(paid.Int64()).QuoInt(received).LTE(orderPrice))
	liquidity.BeginBlocker(s.ctx, s.keeper, s.app.AssetKeeper)

	pair = s.CreateNewLiquidityPair(appID1, s.addr(0), asset3.Denom, asset4.Denom)
	pair.LastPrice = utils.ParseDecP("1.0")
	s.keeper.SetPair(s.ctx, pair)

	s.CreateNewLiquidityRangedPool(appID1, pair.Id, s.addr(1), "1000000denom3,1000000denom4", utils.ParseDec("0.8"), utils.ParseDec("1.3"), utils.ParseDec("1.0"))
	order = s.LimitOrder(appID1, s.addr(2), pair.Id, types.OrderDirectionBuy, orderPrice, orderAmt, 0)
	liquidity.EndBlocker(s.ctx, s.keeper, s.app.AssetKeeper)
	order, _ = s.keeper.GetOrder(s.ctx, appID1, order.PairId, order.Id)
	paid = order.OfferCoin.Sub(order.RemainingOfferCoin).Amount
	received = order.ReceivedCoin.Amount
	s.Require().True(intEq(orderAmt, received))
	s.Require().True(sdkmath.LegacyNewDec(paid.Int64()).QuoInt(received).LTE(orderPrice))
}

func (s *KeeperTestSuite) TestOneSidedRangedPool() {
	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)
	pair.LastPrice = utils.ParseDecP("1.0")
	s.keeper.SetPair(s.ctx, pair)

	pool := s.CreateNewLiquidityRangedPool(appID1, pair.Id, s.addr(1), "1000000denom1,1000000denom2", utils.ParseDec("1.0"), utils.ParseDec("1.2"), utils.ParseDec("1.0"))
	rx, ry := s.keeper.GetPoolBalances(s.ctx, pool)
	ammPool := pool.AMMPool(rx.Amount, ry.Amount, sdkmath.Int{})
	s.Require().True(utils.DecApproxEqual(utils.ParseDec("1.0"), ammPool.Price()))
	s.Require().True(intEq(sdkmath.ZeroInt(), rx.Amount))
	s.Require().True(intEq(sdkmath.NewInt(1000000), ry.Amount))

	orderPrice := utils.ParseDec("1.1")
	orderAmt := sdkmath.NewInt(100000)
	order := s.LimitOrder(appID1, s.addr(2), pair.Id, types.OrderDirectionBuy, utils.ParseDec("1.1"), sdkmath.NewInt(100000), 0)
	liquidity.EndBlocker(s.ctx, s.keeper, s.app.AssetKeeper)
	order, _ = s.keeper.GetOrder(s.ctx, appID1, order.PairId, order.Id)
	paid := order.OfferCoin.Sub(order.RemainingOfferCoin).Amount
	received := order.ReceivedCoin.Amount
	s.Require().True(intEq(orderAmt, received))
	s.Require().True(sdkmath.LegacyNewDec(paid.Int64()).QuoInt(received).LTE(orderPrice))

	rx, _ = s.keeper.GetPoolBalances(s.ctx, pool)
	s.Require().True(rx.IsPositive())
}

func (s *KeeperTestSuite) TestExhaustRangedPool() {
	r := rand.New(rand.NewSource(0))

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, s.addr(0), asset1.Denom, asset2.Denom)

	minPrice, maxPrice := utils.ParseDec("0.5"), utils.ParseDec("2.0")
	initialPrice := utils.ParseDec("1.0")
	pool := s.CreateNewLiquidityRangedPool(appID1, pair.Id, s.addr(1), "1000000denom1,1000000denom2", minPrice, maxPrice, initialPrice)

	orderer := s.addr(2)
	s.fundAddr(orderer, utils.ParseCoins("10000000denom1,10000000denom2"))

	// Buy
	for {
		rx, ry := s.keeper.GetPoolBalances(s.ctx, pool)
		ammPool := pool.AMMPool(rx.Amount, ry.Amount, sdkmath.Int{})
		poolPrice := ammPool.Price()
		if ry.Amount.LT(sdkmath.NewInt(100)) {
			s.Require().True(utils.DecApproxEqual(maxPrice, poolPrice))
			break
		}
		orderPrice := utils.RandomDec(r, poolPrice, poolPrice.Mul(sdkmath.LegacyNewDecWithPrec(105, 2)))
		amt := utils.RandomInt(r, sdkmath.NewInt(5000), sdkmath.NewInt(15000))
		s.LimitOrder(appID1, orderer, pair.Id, types.OrderDirectionBuy, orderPrice, amt, 0)
		s.nextBlock()
	}

	// Sell
	for {
		rx, ry := s.keeper.GetPoolBalances(s.ctx, pool)
		ammPool := pool.AMMPool(rx.Amount, ry.Amount, sdkmath.Int{})
		poolPrice := ammPool.Price()
		if rx.Amount.LT(sdkmath.NewInt(100)) {
			s.Require().True(utils.DecApproxEqual(minPrice, poolPrice))
			break
		}
		orderPrice := utils.RandomDec(r, poolPrice.Mul(sdkmath.LegacyNewDecWithPrec(95, 2)), poolPrice)
		amt := utils.RandomInt(r, sdkmath.NewInt(5000), sdkmath.NewInt(15000))
		s.LimitOrder(appID1, orderer, pair.Id, types.OrderDirectionSell, orderPrice, amt, 0)
		s.nextBlock()
	}

	// Buy again
	for {
		rx, ry := s.keeper.GetPoolBalances(s.ctx, pool)
		ammPool := pool.AMMPool(rx.Amount, ry.Amount, sdkmath.Int{})
		poolPrice := ammPool.Price()
		if poolPrice.GTE(initialPrice) {
			break
		}
		orderPrice := utils.RandomDec(r, poolPrice, poolPrice.Mul(sdkmath.LegacyNewDecWithPrec(105, 2)))
		amt := utils.RandomInt(r, sdkmath.NewInt(5000), sdkmath.NewInt(15000))
		s.LimitOrder(appID1, orderer, pair.Id, types.OrderDirectionBuy, orderPrice, amt, 0)
		s.nextBlock()
	}

	params, err := s.keeper.GetGenericParams(s.ctx, appID1)
	s.Require().NoError(err)

	rx, ry := s.keeper.GetPoolBalances(s.ctx, pool)
	ammPool := pool.AMMPool(rx.Amount, ry.Amount, sdkmath.Int{})
	s.Require().True(coinEq(rx, utils.ParseCoin("997231denom2")))
	s.Require().True(coinEq(ry, utils.ParseCoin("984671denom1")))
	s.Require().True(decEq(ammPool.Price(), utils.ParseDec("1.003719250732340753")))

	s.Require().True(coinsEq(utils.ParseCoins("31534denom2"), s.getBalances(sdk.MustAccAddressFromBech32(params.DustCollectorAddress))))
	s.Require().True(coinsEq(utils.ParseCoins("12546884denom1,12666562denom2"), s.getBalances(orderer)))
}

func (s *KeeperTestSuite) TestOrderBooks_edgecase1() {
	addr1 := s.addr(0)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pair.LastPrice = utils.ParseDecP("0.57472")
	s.keeper.SetPair(s.ctx, pair)

	s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "991883358661denom2,620800303846denom1")
	s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "155025981873denom2,4703143223denom1", utils.ParseDec("1.15"), utils.ParseDec("1.55"), utils.ParseDec("1.5308"))
	s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "223122824634denom2,26528571912denom1", utils.ParseDec("1.25"), utils.ParseDec("1.45"), utils.ParseDec("1.4199"))

	resp, err := s.querier.OrderBooks(sdk.WrapSDKContext(s.ctx), &types.QueryOrderBooksRequest{
		AppId:    appID1,
		PairIds:  []uint64{pair.Id},
		NumTicks: 10,
	})
	s.Require().NoError(err)
	s.Require().Len(resp.Pairs, 1)
	s.Require().Len(resp.Pairs[0].OrderBooks, 3)

	s.Require().Len(resp.Pairs[0].OrderBooks[0].Buys, 2)
	s.Require().True(decEq(utils.ParseDec("0.63219"), resp.Pairs[0].OrderBooks[0].Buys[0].Price))
	s.Require().True(intEq(sdkmath.NewInt(1178846737645), resp.Pairs[0].OrderBooks[0].Buys[0].UserOrderAmount))
	s.Require().True(decEq(utils.ParseDec("0.5187"), resp.Pairs[0].OrderBooks[0].Buys[1].Price))
	s.Require().True(intEq(sdkmath.NewInt(13340086), resp.Pairs[0].OrderBooks[0].Buys[1].UserOrderAmount))
	s.Require().Len(resp.Pairs[0].OrderBooks[0].Sells, 0)
}

func (s *KeeperTestSuite) TestSwap_edgecase1() {
	addr1 := s.addr(0)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	s.LimitOrder(appID1, s.addr(2), pair.Id, types.OrderDirectionSell, utils.ParseDec("0.102"), sdkmath.NewInt(10000), 0)
	s.LimitOrder(appID1, s.addr(3), pair.Id, types.OrderDirectionSell, utils.ParseDec("0.101"), sdkmath.NewInt(9995), 0)
	s.LimitOrder(appID1, s.addr(4), pair.Id, types.OrderDirectionBuy, utils.ParseDec("0.102"), sdkmath.NewInt(10000), 0)

	s.nextBlock()
	pair, _ = s.keeper.GetPair(s.ctx, appID1, pair.Id)
	s.Require().True(decEq(utils.ParseDec("0.102"), *pair.LastPrice))

	s.LimitOrder(appID1, s.addr(2), pair.Id, types.OrderDirectionSell, utils.ParseDec("0.102"), sdkmath.NewInt(10000), 0)
	s.LimitOrder(appID1, s.addr(3), pair.Id, types.OrderDirectionSell, utils.ParseDec("0.101"), sdkmath.NewInt(9995), 0)
	s.LimitOrder(appID1, s.addr(4), pair.Id, types.OrderDirectionBuy, utils.ParseDec("0.102"), sdkmath.NewInt(10000), 0)
	s.nextBlock()
}

func (s *KeeperTestSuite) TestSwap_edgecase2() {
	addr1 := s.addr(0)
	addr2 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pair.LastPrice = utils.ParseDecP("1.6724")
	s.keeper.SetPair(s.ctx, pair)

	s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1005184935980denom2,601040339855denom1")
	s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "17335058855denom2", utils.ParseDec("1.15"), utils.ParseDec("1.55"), utils.ParseDec("1.55"))
	s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "217771046279denom2", utils.ParseDec("1.25"), utils.ParseDec("1.45"), utils.ParseDec("1.45"))

	s.MarketOrder(appID1, addr2, pair.Id, types.OrderDirectionSell, sdkmath.NewInt(4336_000000), 0)
	s.nextBlock()

	pair, _ = s.keeper.GetPair(s.ctx, appID1, pair.Id)
	s.Require().True(decEq(utils.ParseDec("1.6484"), *pair.LastPrice))

	s.nextBlock()
	pair, _ = s.keeper.GetPair(s.ctx, appID1, pair.Id)
	s.Require().True(decEq(utils.ParseDec("1.6484"), *pair.LastPrice))

	s.MarketOrder(appID1, addr2, pair.Id, types.OrderDirectionSell, sdkmath.NewInt(4450_000000), 0)
	s.nextBlock()

	pair, _ = s.keeper.GetPair(s.ctx, appID1, pair.Id)
	s.Require().True(decEq(utils.ParseDec("1.6248"), *pair.LastPrice))
}

func (s *KeeperTestSuite) TestSwap_edgecase3() {
	addr1 := s.addr(0)
	addr2 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pair.LastPrice = utils.ParseDecP("0.99992")
	s.keeper.SetPair(s.ctx, pair)

	s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "110001546090denom2,110013588106denom1")
	s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "140913832254denom2,130634675302denom1", utils.ParseDec("0.92"), utils.ParseDec("1.08"), utils.ParseDec("0.99989"))

	s.MarketOrder(appID1, addr2, pair.Id, types.OrderDirectionBuy, sdkmath.NewInt(30_000000), 0)
	s.nextBlock()

	pair, _ = s.keeper.GetPair(s.ctx, appID1, pair.Id)
	s.Require().True(decEq(utils.ParseDec("0.99992"), *pair.LastPrice))
}

func (s *KeeperTestSuite) TestSwap_edgecase4() {
	addr1 := s.addr(0)
	addr2 := s.addr(1)

	appID1 := s.CreateNewApp("appone")

	asset1 := s.CreateNewAsset("ASSETONE", "denom1", 1000000)
	asset2 := s.CreateNewAsset("ASSETTWO", "denom2", 1000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	pair.LastPrice = utils.ParseDecP("0.99999")
	s.keeper.SetPair(s.ctx, pair)

	s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000_000000denom1,100_000000denom2")
	s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000_000000denom1,1000_000000denom2", utils.ParseDec("0.95"), utils.ParseDec("1.05"), utils.ParseDec("1.02"))
	s.CreateNewLiquidityRangedPool(appID1, pair.Id, addr1, "1000_000000denom1,1000_000000denom2", utils.ParseDec("0.9"), utils.ParseDec("1.2"), utils.ParseDec("0.98"))

	s.LimitOrder(appID1, addr2, pair.Id, types.OrderDirectionSell, utils.ParseDec("1.05"), sdkmath.NewInt(50_000000), 0)
	s.LimitOrder(appID1, addr2, pair.Id, types.OrderDirectionBuy, utils.ParseDec("0.97"), sdkmath.NewInt(100_000000), 0)
	s.nextBlock()
	s.Require().True(utils.ParseCoins("50150000denom1,97291000denom2").Equal(s.getBalances(addr2)))

}
