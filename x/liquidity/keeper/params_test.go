package keeper_test

import (
	"fmt"
	"time"

	utils "github.com/comdex-official/comdex/types"
	"github.com/comdex-official/comdex/x/liquidity/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestUpdateGenericParams() {
	appID1 := s.CreateNewApp("appone")

	_, err := s.keeper.GetGenericParams(s.ctx, appID1)
	s.Require().NoError(err)

	testCases := []struct {
		Name   string
		appID  uint64
		Keys   []string
		Values []string
		ExpErr error
	}{
		{
			Name:   "error app id invalid",
			appID:  69,
			Keys:   []string{},
			Values: []string{},
			ExpErr: sdkerrors.Wrapf(types.ErrInvalidAppID, "app id 69 not found"),
		},
		{
			Name:   "error key-value length mismatch",
			appID:  appID1,
			Keys:   []string{"abc"},
			Values: []string{"def", "ghi"},
			ExpErr: fmt.Errorf("keys and values list length mismatch"),
		},
		{
			Name:   "error invalid key for update",
			appID:  appID1,
			Keys:   []string{"abc"},
			Values: []string{"def"},
			ExpErr: fmt.Errorf("invalid key for update: abc"),
		},
		{
			Name:   "error invalid value type",
			appID:  appID1,
			Keys:   []string{"BatchSize"},
			Values: []string{"abc"},
			ExpErr: fmt.Errorf("strconv.ParseUint: parsing \"abc\": invalid syntax"),
		},
		{
			Name:   "error invalid fee type",
			appID:  appID1,
			Keys:   []string{"PairCreationFee"},
			Values: []string{"helloworld"},
			ExpErr: fmt.Errorf("invalid decimal coin expression: helloworld"),
		},
		{
			Name:   "error invalid order lifespan type",
			appID:  appID1,
			Keys:   []string{"MaxOrderLifespan"},
			Values: []string{"10b"},
			ExpErr: fmt.Errorf("time: unknown unit \"b\" in duration \"10b\""),
		},
		{
			Name:   "success valid case 1",
			appID:  appID1,
			Keys:   []string{"BatchSize", "TickPrecision", "MinInitialPoolCoinSupply", "PairCreationFee", "PoolCreationFee", "MinInitialDepositAmount", "MaxPriceLimitRatio", "MaxOrderLifespan", "SwapFeeRate", "WithdrawFeeRate", "DepositExtraGas", "WithdrawExtraGas", "OrderExtraGas", "SwapFeeDistrDenom", "SwapFeeBurnRate"},
			Values: []string{"69", "699", "1000000000000000000", "10000000000dummy1", "10000000000dummy2", "10000", "0.1", "10h", "0.2", "0.4", "11", "12", "13", "loltoken", "0.8"},
			ExpErr: nil,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			err := s.keeper.UpdateGenericParams(s.ctx, tc.appID, tc.Keys, tc.Values)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				params, err := s.keeper.GetGenericParams(s.ctx, tc.appID)
				s.Require().NoError(err)

				s.Require().Equal(params.BatchSize, uint64(69))
				s.Require().Equal(params.TickPrecision, uint64(699))
				s.Require().Equal(params.MinInitialPoolCoinSupply, newInt(1000000000000000000))
				s.Require().Equal(params.PairCreationFee, utils.ParseCoins("10000000000dummy1"))
				s.Require().Equal(params.PoolCreationFee, utils.ParseCoins("10000000000dummy2"))
				s.Require().Equal(params.MinInitialDepositAmount, newInt(10000))
				s.Require().Equal(params.MaxPriceLimitRatio, utils.ParseDec("0.1"))
				s.Require().Equal(params.MaxOrderLifespan, time.Hour*10)
				s.Require().Equal(params.SwapFeeRate, utils.ParseDec("0.2"))
				s.Require().Equal(params.WithdrawFeeRate, utils.ParseDec("0.4"))
				s.Require().Equal(params.DepositExtraGas, uint64(11))
				s.Require().Equal(params.WithdrawExtraGas, uint64(12))
				s.Require().Equal(params.OrderExtraGas, uint64(13))
				s.Require().Equal(params.SwapFeeDistrDenom, "loltoken")
				s.Require().Equal(params.SwapFeeBurnRate, utils.ParseDec("0.8"))
			}
		})

	}
}
