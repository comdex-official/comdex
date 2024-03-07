package keeper_test

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/comdex-official/comdex/x/gasless/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/stretchr/testify/suite"
)

func (s *KeeperTestSuite) TestCreateGasTank() {
	params := s.keeper.GetParams(s.ctx)

	provider1 := s.addr(1)
	provider1Tanks := []uint64{}
	for i := 0; i < int(params.TankCreationLimit); i++ {
		tankID := s.CreateNewGasTank(provider1, "ucmdx", sdkmath.NewInt(1000), 10, sdkmath.NewInt(1000000), []string{"/comdex.liquidity.v1beta1.MsgLimitOrder"}, []string{}, "100000000ucmdx")
		provider1Tanks = append(provider1Tanks, tankID)
	}

	testCases := []struct {
		Name   string
		Msg    types.MsgCreateGasTank
		ExpErr error
	}{
		{
			Name:   "error tank creation limit reached",
			Msg:    *types.NewMsgCreateGasTank(provider1, "ucmdx", sdkmath.NewInt(1000), 10, sdkmath.NewInt(1000000), []string{}, []string{"comdex1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3"}, sdk.NewCoin("ucmdx", sdk.NewInt(100000000))),
			ExpErr: sdkerrors.Wrapf(types.ErrorMaxLimitReachedByProvider, " %d gas tanks already created by the provider", params.TankCreationLimit),
		},
		{
			Name:   "error fee and deposit denom mismatch",
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "uatom", sdkmath.NewInt(1000), 10, sdkmath.NewInt(1000000), []string{}, []string{"comdex1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3"}, sdk.NewCoin("ucmdx", sdk.NewInt(100000000))),
			ExpErr: sdkerrors.Wrapf(types.ErrorInvalidrequest, " fee denom %s do not match gas depoit denom %s ", "uatom", "ucmdx"),
		},
		{
			Name:   "error max tx count consumer is 0",
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "ucmdx", sdkmath.NewInt(1000), 0, sdkmath.NewInt(1000000), []string{}, []string{"comdex1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3"}, sdk.NewCoin("ucmdx", sdk.NewInt(100000000))),
			ExpErr: sdkerrors.Wrap(types.ErrorInvalidrequest, "max tx count per consumer must not be 0"),
		},
		{
			Name:   "error max fee usage per tx should be positive",
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "ucmdx", sdkmath.NewInt(0), 123, sdkmath.NewInt(1000000), []string{}, []string{"comdex1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3"}, sdk.NewCoin("ucmdx", sdk.NewInt(100000000))),
			ExpErr: sdkerrors.Wrapf(types.ErrorInvalidrequest, "max_fee_usage_per_tx should be positive"),
		},
		{
			Name:   "error max fee usage per consumer should be positive",
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "ucmdx", sdkmath.NewInt(123), 123, sdkmath.NewInt(0), []string{}, []string{"comdex1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3"}, sdk.NewCoin("ucmdx", sdk.NewInt(100000000))),
			ExpErr: sdkerrors.Wrapf(types.ErrorInvalidrequest, "max_fee_usage_per_consumer should be positive"),
		},
		{
			Name:   "error atleast one txPath or contract is required",
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "ucmdx", sdkmath.NewInt(123), 123, sdkmath.NewInt(1000000), []string{}, []string{}, sdk.NewCoin("ucmdx", sdk.NewInt(100000000))),
			ExpErr: sdkerrors.Wrapf(types.ErrorInvalidrequest, "request should have atleast one tx path or contract address"),
		},
		{
			Name:   "error deposit samller than required min deposit",
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "ucmdx", sdkmath.NewInt(123), 123, sdkmath.NewInt(1000000), []string{"/comdex.liquidity.v1beta1.MsgLimitOrder"}, []string{}, sdk.NewCoin("ucmdx", sdk.NewInt(100))),
			ExpErr: sdkerrors.Wrapf(types.ErrorInvalidrequest, "minimum required deposit is %s", params.MinimumGasDeposit[0].String()),
		},
		{
			Name:   "error fee denom not allowed",
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "uatom", sdkmath.NewInt(123), 123, sdkmath.NewInt(1000000), []string{"/comdex.liquidity.v1beta1.MsgLimitOrder"}, []string{}, sdk.NewCoin("uatom", sdk.NewInt(100))),
			ExpErr: sdkerrors.Wrapf(types.ErrorInvalidrequest, " fee denom %s not allowed ", "uatom"),
		},
		{
			Name:   "error invalid message type URL",
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "ucmdx", sdkmath.NewInt(123), 123, sdkmath.NewInt(1000000), []string{"random message type"}, []string{""}, sdk.NewCoin("ucmdx", sdk.NewInt(100000000))),
			ExpErr: sdkerrors.Wrapf(types.ErrorInvalidrequest, "invalid message - %s", "random message type"),
		},
		{
			Name:   "error invalid contract address",
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "ucmdx", sdkmath.NewInt(123), 123, sdkmath.NewInt(1000000), []string{}, []string{"comdex1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3"}, sdk.NewCoin("ucmdx", sdk.NewInt(100000000))),
			ExpErr: sdkerrors.Wrapf(types.ErrorInvalidrequest, "invalid contract address - %s", "comdex1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3"),
		},
		{
			Name:   "success gas tank creation",
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "ucmdx", sdkmath.NewInt(123), 123, sdkmath.NewInt(1000000), []string{"/comdex.liquidity.v1beta1.MsgLimitOrder"}, []string{}, sdk.NewCoin("ucmdx", sdk.NewInt(100000000))),
			ExpErr: nil,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			// add funds to acount for valid case
			if tc.ExpErr == nil {
				s.fundAddr(sdk.MustAccAddressFromBech32(tc.Msg.Provider), sdk.NewCoins(tc.Msg.GasDeposit))
			}

			tank, err := s.keeper.CreateGasTank(s.ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(tank)

				s.Require().IsType(types.GasTank{}, tank)
				s.Require().Equal(tc.Msg.FeeDenom, tank.FeeDenom)
				s.Require().Equal(tc.Msg.MaxFeeUsagePerTx, tank.MaxFeeUsagePerTx)
				s.Require().Equal(tc.Msg.MaxTxsCountPerConsumer, tank.MaxTxsCountPerConsumer)
				s.Require().Equal(tc.Msg.MaxFeeUsagePerConsumer, tank.MaxFeeUsagePerConsumer)
				s.Require().Equal(tc.Msg.TxsAllowed, tank.TxsAllowed)
				s.Require().Equal(tc.Msg.ContractsAllowed, tank.ContractsAllowed)
				s.Require().Equal(tc.Msg.GasDeposit, s.getBalance(tank.GetGasTankReserveAddress(), tank.FeeDenom))
			}
		})
	}
}
