package keeper_test

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/comdex-official/comdex/x/gasless/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errors "github.com/cosmos/cosmos-sdk/types/errors"
	_ "github.com/stretchr/testify/suite"
	"golang.org/x/exp/slices"
)

func (s *KeeperTestSuite) TestCreateGasTank() {
	params := s.keeper.GetParams(s.ctx)

	provider1 := s.addr(1)
	provider1Tanks := []types.GasTank{}
	for i := 0; i < int(params.TankCreationLimit); i++ {
		tank := s.CreateNewGasTank(provider1, "ucmdx", sdkmath.NewInt(1000), 10, sdkmath.NewInt(1000000), []string{"/comdex.liquidity.v1beta1.MsgLimitOrder"}, []string{}, "100000000ucmdx")
		provider1Tanks = append(provider1Tanks, tank)
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
			Msg:    *types.NewMsgCreateGasTank(s.addr(2), "ucmdx", sdkmath.NewInt(123), 123, sdkmath.NewInt(1000000), []string{"/comdex.liquidity.v1beta1.MsgLimitOrder"}, []string{"comdex1qa4hswlcjmttulj0q9qa46jf64f93pecl6tydcsjldfe0hy5ju0s7r3hn3"}, sdk.NewCoin("ucmdx", sdk.NewInt(100000000))),
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

				for _, tx := range tc.Msg.TxsAllowed {
					txGtids, found := s.keeper.GetTxGTIDs(s.ctx, tx)
					s.Require().True(found)
					s.Require().IsType(types.TxGTIDs{}, txGtids)
					s.Require().IsType([]uint64{}, txGtids.GasTankIds)
					s.Require().Equal(txGtids.TxPathOrContractAddress, tx)
					s.Require().Equal(tank.Id, txGtids.GasTankIds[len(txGtids.GasTankIds)-1])
				}

				for _, c := range tc.Msg.ContractsAllowed {
					txGtids, found := s.keeper.GetTxGTIDs(s.ctx, c)
					s.Require().True(found)
					s.Require().IsType(types.TxGTIDs{}, txGtids)
					s.Require().IsType([]uint64{}, txGtids.GasTankIds)
					s.Require().Equal(txGtids.TxPathOrContractAddress, c)
					s.Require().Equal(tank.Id, txGtids.GasTankIds[len(txGtids.GasTankIds)-1])
				}
			}
		})
	}
}

func (s *KeeperTestSuite) TestAuthorizeActors() {
	provider1 := s.addr(1)
	tank1 := s.CreateNewGasTank(provider1, "ucmdx", sdkmath.NewInt(1000), 10, sdkmath.NewInt(1000000), []string{"/comdex.liquidity.v1beta1.MsgLimitOrder"}, []string{}, "100000000ucmdx")

	provider2 := s.addr(2)
	inactiveTank := s.CreateNewGasTank(provider2, "ucmdx", sdkmath.NewInt(1000), 10, sdkmath.NewInt(1000000), []string{"/comdex.liquidity.v1beta1.MsgLimitOrder"}, []string{}, "100000000ucmdx")
	inactiveTank.IsActive = false
	s.keeper.SetGasTank(s.ctx, inactiveTank)

	testCases := []struct {
		Name   string
		Msg    types.MsgAuthorizeActors
		ExpErr error
	}{
		{
			Name: "error invalid gas tank ID",
			Msg: *types.NewMsgAuthorizeActors(
				12, provider1, []sdk.AccAddress{s.addr(10), s.addr(11), s.addr(12)},
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrNotFound, "gas tank with id %d not found", 12),
		},
		{
			Name: "error unauthorized provider",
			Msg: *types.NewMsgAuthorizeActors(
				tank1.Id, provider2, []sdk.AccAddress{s.addr(10), s.addr(11), s.addr(12)},
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrUnauthorized, "unauthorized provider"),
		},
		{
			Name: "error inactive tank",
			Msg: *types.NewMsgAuthorizeActors(
				inactiveTank.Id, provider2, []sdk.AccAddress{s.addr(10), s.addr(11), s.addr(12)},
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrInvalidRequest, "gas tank inactive"),
		},
		{
			Name: "error max actor limit ",
			Msg: *types.NewMsgAuthorizeActors(
				tank1.Id, provider1, []sdk.AccAddress{s.addr(10), s.addr(11), s.addr(12), s.addr(13), s.addr(14), s.addr(15), s.addr(16)},
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrInvalidRequest, "maximum %d actors can be authorized", types.MaximumAuthorizedActorsLimit),
		},
		{
			Name: "success with duplicate actors",
			Msg: *types.NewMsgAuthorizeActors(
				tank1.Id, provider1, []sdk.AccAddress{s.addr(10), s.addr(10), s.addr(10), s.addr(10), s.addr(10), s.addr(10), s.addr(10)},
			),
			ExpErr: nil,
		},
		{
			Name: "success with unique actors",
			Msg: *types.NewMsgAuthorizeActors(
				tank1.Id, provider1, []sdk.AccAddress{s.addr(10), s.addr(11), s.addr(12), s.addr(13), s.addr(14)},
			),
			ExpErr: nil,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			tank, err := s.keeper.AuthorizeActors(s.ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(tank)

				s.Require().IsType(types.GasTank{}, tank)
				s.Require().Equal(len(types.RemoveDuplicates(tc.Msg.Actors)), len(tank.AuthorizedActors))
				slices.Sort(tc.Msg.Actors)
				slices.Sort(tank.AuthorizedActors)
				s.Require().Equal(tc.Msg.Actors, tank.AuthorizedActors)
			}
		})
	}

}

func (s *KeeperTestSuite) TestUpdateGasTankStatus() {
	provider1 := s.addr(1)
	tank1 := s.CreateNewGasTank(provider1, "ucmdx", sdkmath.NewInt(1000), 10, sdkmath.NewInt(1000000), []string{"/comdex.liquidity.v1beta1.MsgLimitOrder"}, []string{}, "100000000ucmdx")

	testCases := []struct {
		Name   string
		Msg    types.MsgUpdateGasTankStatus
		ExpErr error
	}{
		{
			Name: "error invalid gas tank ID",
			Msg: *types.NewMsgUpdateGasTankStatus(
				12, provider1,
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrNotFound, "gas tank with id %d not found", 12),
		},
		{
			Name: "error unauthorized provider",
			Msg: *types.NewMsgUpdateGasTankStatus(
				tank1.Id, s.addr(10),
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrUnauthorized, "unauthorized provider"),
		},
		{
			Name: "success update status to inactive",
			Msg: *types.NewMsgUpdateGasTankStatus(
				tank1.Id, provider1,
			),
			ExpErr: nil,
		},
		{
			Name: "success update status to active",
			Msg: *types.NewMsgUpdateGasTankStatus(
				tank1.Id, provider1,
			),
			ExpErr: nil,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			tank, _ := s.keeper.GetGasTank(s.ctx, tc.Msg.GasTankId)
			updatedTank, err := s.keeper.UpdateGasTankStatus(s.ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(tank)

				s.Require().IsType(types.GasTank{}, updatedTank)
				s.Require().Equal(tank.IsActive, !updatedTank.IsActive)
			}
		})
	}
}

func (s *KeeperTestSuite) TestUpdateGasTankConfig() {
	provider1 := s.addr(1)
	tank1 := s.CreateNewGasTank(provider1, "ucmdx", sdkmath.NewInt(1000), 10, sdkmath.NewInt(1000000), []string{"/comdex.liquidity.v1beta1.MsgLimitOrder"}, []string{}, "100000000ucmdx")

	provider2 := s.addr(2)
	inactiveTank := s.CreateNewGasTank(provider2, "ucmdx", sdkmath.NewInt(1000), 10, sdkmath.NewInt(1000000), []string{"/comdex.liquidity.v1beta1.MsgCreatePool"}, []string{}, "100000000ucmdx")
	inactiveTank.IsActive = false
	s.keeper.SetGasTank(s.ctx, inactiveTank)

	testCases := []struct {
		Name   string
		Msg    types.MsgUpdateGasTankConfig
		ExpErr error
	}{
		{
			Name: "error invalid gas tank ID",
			Msg: *types.NewMsgUpdateGasTankConfig(
				12, provider1, sdk.NewInt(1000), 10, sdk.NewInt(1000000),
				[]string{"/comdex.liquidity.v1beta1.MsgLimitOrder"},
				[]string{},
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrNotFound, "gas tank with id %d not found", 12),
		},
		{
			Name: "error unauthorized provider",
			Msg: *types.NewMsgUpdateGasTankConfig(
				tank1.Id, provider2, sdk.NewInt(1000), 10, sdk.NewInt(1000000),
				[]string{"/comdex.liquidity.v1beta1.MsgLimitOrder"},
				[]string{},
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrUnauthorized, "unauthorized provider"),
		},
		{
			Name: "error inactive tank",
			Msg: *types.NewMsgUpdateGasTankConfig(
				inactiveTank.Id, provider1, sdk.NewInt(1000), 10, sdk.NewInt(1000000),
				[]string{"/comdex.liquidity.v1beta1.MsgLimitOrder"},
				[]string{},
			),
			ExpErr: sdkerrors.Wrapf(errors.ErrInvalidRequest, "gas tank inactive"),
		},
		{
			Name: "error max tx count consumer is 0",
			Msg: *types.NewMsgUpdateGasTankConfig(
				tank1.Id, provider1, sdk.NewInt(1000), 0, sdk.NewInt(1000000),
				[]string{"/comdex.liquidity.v1beta1.MsgLimitOrder"},
				[]string{},
			),
			ExpErr: sdkerrors.Wrap(types.ErrorInvalidrequest, "max tx count per consumer must not be 0"),
		},
		{
			Name: "error max fee usage per tx should be positive",
			Msg: *types.NewMsgUpdateGasTankConfig(
				tank1.Id, provider1, sdk.ZeroInt(), 10, sdk.NewInt(1000000),
				[]string{"/comdex.liquidity.v1beta1.MsgLimitOrder"},
				[]string{},
			),
			ExpErr: sdkerrors.Wrapf(types.ErrorInvalidrequest, "max_fee_usage_per_tx should be positive"),
		},
		{
			Name: "error max fee usage per consumer should be positive",
			Msg: *types.NewMsgUpdateGasTankConfig(
				tank1.Id, provider1, sdk.NewInt(1000), 10, sdk.ZeroInt(),
				[]string{"/comdex.liquidity.v1beta1.MsgLimitOrder"},
				[]string{},
			),
			ExpErr: sdkerrors.Wrapf(types.ErrorInvalidrequest, "max_fee_usage_per_consumer should be positive"),
		},
		{
			Name: "error atleast one txPath or contract is required",
			Msg: *types.NewMsgUpdateGasTankConfig(
				tank1.Id, provider1, sdk.NewInt(1000), 10, sdk.NewInt(1000000),
				[]string{},
				[]string{},
			),
			ExpErr: sdkerrors.Wrapf(types.ErrorInvalidrequest, "request should have atleast one tx path or contract address"),
		},
		{
			Name: "error invalid message type URL",
			Msg: *types.NewMsgUpdateGasTankConfig(
				tank1.Id, provider1, sdk.NewInt(1000), 10, sdk.NewInt(1000000),
				[]string{"random message type"},
				[]string{"contract address"},
			),
			ExpErr: sdkerrors.Wrapf(types.ErrorInvalidrequest, "invalid message - %s", "random message type"),
		},
		{
			Name: "error invalid contract address",
			Msg: *types.NewMsgUpdateGasTankConfig(
				tank1.Id, provider1, sdk.NewInt(1000), 10, sdk.NewInt(1000000),
				[]string{"/comdex.liquidity.v1beta1.MsgLimitOrder"},
				[]string{"invalid contract address"},
			),
			ExpErr: sdkerrors.Wrapf(types.ErrorInvalidrequest, "invalid contract address - %s", "invalid contract address"),
		},
		{
			Name: "success tank configs updated",
			Msg: *types.NewMsgUpdateGasTankConfig(
				tank1.Id, provider1, sdk.NewInt(25000), 120, sdk.NewInt(150000000),
				[]string{"/comdex.liquidity.v1beta1.MsgMarketOrder"},
				nil,
			),
			ExpErr: nil,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.Name, func() {
			resp, err := s.keeper.UpdateGasTankConfig(s.ctx, &tc.Msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(resp)

				s.Require().IsType(types.GasTank{}, resp)

				checkTank, _ := s.keeper.GetGasTank(s.ctx, tc.Msg.GasTankId)
				s.Require().Equal(tc.Msg.MaxFeeUsagePerTx, checkTank.MaxFeeUsagePerTx)
				s.Require().Equal(tc.Msg.MaxTxsCountPerConsumer, checkTank.MaxTxsCountPerConsumer)
				s.Require().Equal(tc.Msg.MaxFeeUsagePerConsumer, checkTank.MaxFeeUsagePerConsumer)
				slices.Sort(tc.Msg.TxsAllowed)
				slices.Sort(checkTank.TxsAllowed)
				slices.Sort(tc.Msg.ContractsAllowed)
				slices.Sort(checkTank.ContractsAllowed)
				s.Require().Equal(tc.Msg.TxsAllowed, checkTank.TxsAllowed)
				s.Require().Equal(tc.Msg.ContractsAllowed, checkTank.ContractsAllowed)

				// validate if new txs and contracts has been added to the index of TxGTIDs
				for _, tx := range tc.Msg.TxsAllowed {
					txGtids, found := s.keeper.GetTxGTIDs(s.ctx, tx)
					s.Require().True(found)
					s.Require().IsType(types.TxGTIDs{}, txGtids)
					s.Require().IsType([]uint64{}, txGtids.GasTankIds)
					s.Require().Equal(txGtids.TxPathOrContractAddress, tx)
					s.Require().Equal(resp.Id, txGtids.GasTankIds[len(txGtids.GasTankIds)-1])
				}

				for _, c := range tc.Msg.ContractsAllowed {
					txGtids, found := s.keeper.GetTxGTIDs(s.ctx, c)
					s.Require().True(found)
					s.Require().IsType(types.TxGTIDs{}, txGtids)
					s.Require().IsType([]uint64{}, txGtids.GasTankIds)
					s.Require().Equal(txGtids.TxPathOrContractAddress, c)
					s.Require().Equal(resp.Id, txGtids.GasTankIds[len(txGtids.GasTankIds)-1])
				}

				// validate if old txs and contracts has been removed from the index of TxGTIDs
				for _, tx := range tank1.TxsAllowed {
					_, found := s.keeper.GetTxGTIDs(s.ctx, tx)
					s.Require().False(found)
				}

				for _, c := range tank1.ContractsAllowed {
					_, found := s.keeper.GetTxGTIDs(s.ctx, c)
					s.Require().False(found)
				}
			}
		})
	}
}
