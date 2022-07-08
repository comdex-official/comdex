package keeper_test

import (
	"time"

	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/liquidity/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestLimitOrder() {
	addr1 := s.addr(1)

	appID1 := s.CreateNewApp("appOne")

	asset1 := s.CreateNewAsset("ASSET1", "uasset1", 1000000)
	asset2 := s.CreateNewAsset("ASSET2", "uasset2", 2000000)

	pair := s.CreateNewLiquidityPair(appID1, addr1, asset1.Denom, asset2.Denom)
	_ = s.CreateNewLiquidityPool(appID1, pair.Id, addr1, "1000000000000uasset1,500000000000uasset2")

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
				utils.ParseCoin("10030000uasset"),
				asset1.Denom,
				newDec(1),
				newInt(10000000),
				time.Second*10,
			),
			ExpErr:  sdkerrors.Wrap(sdkerrors.Wrapf(types.ErrInvalidAppID, "app id %d not found", 69), "params retreval failed"),
			ExpResp: &types.Order{},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.Name, func() {
			order, err := s.keeper.LimitOrder(s.ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
				s.Require().Equal(tc.ExpResp, &order)
			}

		})
	}
}
