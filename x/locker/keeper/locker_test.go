package keeper_test

import (
	"fmt"
	"github.com/petrichormoney/petri/app/wasm/bindings"
	utils "github.com/petrichormoney/petri/types"
	assetTypes "github.com/petrichormoney/petri/x/asset/types"
	collectortypes "github.com/petrichormoney/petri/x/collector/types"
	"github.com/petrichormoney/petri/x/locker/keeper"
	lockerTypes "github.com/petrichormoney/petri/x/locker/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *KeeperTestSuite) AddAppAsset() {
	assetKeeper, ctx := &s.assetKeeper, &s.ctx
	msg1 := assetTypes.AppData{
		Name:             "cswap",
		ShortName:        "cswap",
		MinGovDeposit:    sdk.NewIntFromUint64(10000000),
		GovTimeInSeconds: 900,
	}
	err := assetKeeper.AddAppRecords(*ctx, msg1)
	s.Require().NoError(err)

	msg2 := assetTypes.AppData{
		Name:             "commodo",
		ShortName:        "comdo",
		MinGovDeposit:    sdk.NewIntFromUint64(10000000),
		GovTimeInSeconds: 900,
	}
	err = assetKeeper.AddAppRecords(*ctx, msg2)
	s.Require().NoError(err)

	msg3 := assetTypes.Asset{
		Name:      "PETRI",
		Denom:     "upetri",
		Decimals:  sdk.NewInt(1000000),
		IsOnChain: true,
	}

	err = assetKeeper.AddAssetRecords(*ctx, msg3)
	s.Require().NoError(err)

	msg4 := assetTypes.Asset{
		Name:      "FUST",
		Denom:     "ufust",
		Decimals:  sdk.NewInt(1000000),
		IsOnChain: true,
	}
	err = assetKeeper.AddAssetRecords(*ctx, msg4)
	s.Require().NoError(err)

	msg5 := assetTypes.Asset{
		Name:      "HARBOR",
		Denom:     "uharbor",
		Decimals:  sdk.NewInt(1000000),
		IsOnChain: true,
	}
	err = assetKeeper.AddAssetRecords(*ctx, msg5)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) AddCollectorLookupTable() {
	collectorKeeper, ctx := &s.collector, &s.ctx
	msg1 := bindings.MsgSetCollectorLookupTable{
		AppID:            1,
		CollectorAssetID: 1,
		SecondaryAssetID: 3,
		SurplusThreshold: sdk.NewInt(10000000),
		DebtThreshold:    sdk.NewInt(5000000),
		LockerSavingRate: sdk.MustNewDecFromStr("0.1"),
		LotSize:          sdk.NewInt(2000000),
		BidFactor:        sdk.MustNewDecFromStr("0.01"),
		DebtLotSize:      sdk.NewInt(2000000),
	}
	err := collectorKeeper.WasmSetCollectorLookupTable(*ctx, &msg1)
	s.Require().NoError(err)

	msg2 := bindings.MsgSetCollectorLookupTable{
		AppID:            1,
		CollectorAssetID: 2,
		SecondaryAssetID: 3,
		SurplusThreshold: sdk.NewInt(10000000),
		DebtThreshold:    sdk.NewInt(5000000),
		LockerSavingRate: sdk.MustNewDecFromStr("0.1"),
		LotSize:          sdk.NewInt(2000000),
		BidFactor:        sdk.MustNewDecFromStr("0.01"),
		DebtLotSize:      sdk.NewInt(2000000),
	}
	err1 := collectorKeeper.WasmSetCollectorLookupTable(*ctx, &msg2)
	s.Require().NoError(err1)

	msg3 := bindings.MsgSetCollectorLookupTable{
		AppID:            2,
		CollectorAssetID: 1,
		SecondaryAssetID: 3,
		SurplusThreshold: sdk.NewInt(10000000),
		DebtThreshold:    sdk.NewInt(5000000),
		LockerSavingRate: sdk.MustNewDecFromStr("0.1"),
		LotSize:          sdk.NewInt(2000000),
		BidFactor:        sdk.MustNewDecFromStr("0.01"),
		DebtLotSize:      sdk.NewInt(2000000),
	}
	err2 := collectorKeeper.WasmSetCollectorLookupTable(*ctx, &msg3)
	s.Require().NoError(err2)
}

func (s *KeeperTestSuite) TestCreateLocker() {
	userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	lockerKeeper, ctx := &s.lockerKeeper, &s.ctx
	s.AddAppAsset()
	s.AddCollectorLookupTable()
	server := keeper.NewMsgServer(*lockerKeeper)

	// Add whitelisted App Asset combinations
	for _, tc := range []struct {
		name string
		msg  lockerTypes.MsgAddWhiteListedAssetRequest
	}{
		{
			"Whitelist : App1 Asset 1",
			lockerTypes.MsgAddWhiteListedAssetRequest{
				From:    userAddress,
				AppId:   1,
				AssetId: 1,
			},
		},
		{
			"Whitelist : App1 Asset 2",
			lockerTypes.MsgAddWhiteListedAssetRequest{
				From:    userAddress,
				AppId:   1,
				AssetId: 2,
			},
		},
		{
			"Whitelist : App2 Asset 1",
			lockerTypes.MsgAddWhiteListedAssetRequest{
				From:    userAddress,
				AppId:   2,
				AssetId: 1,
			},
		},
	} {
		s.Run(tc.name, func() {
			_, err := lockerKeeper.AddWhiteListedAsset(*ctx, &tc.msg)
			s.Require().NoError(err)
		})
	}

	// create lockers for App Asset combination , query locker and validate
	for _, tc := range []struct {
		name          string
		msg           lockerTypes.MsgCreateLockerRequest
		fundAmount    uint64
		expectedError bool
		query         lockerTypes.QueryLockerInfoRequest
		ExpErr        error
	}{
		{
			"CreateLocker : App1 Asset 1",
			lockerTypes.MsgCreateLockerRequest{
				Depositor: userAddress,
				Amount:    sdk.NewIntFromUint64(1000000),
				AssetId:   1,
				AppId:     1,
			},
			1000000,
			false,
			lockerTypes.QueryLockerInfoRequest{
				Id: 1,
			},
			nil,
		},
		{
			"CreateLocker : Duplicate locker App1 Asset1 should fail",
			lockerTypes.MsgCreateLockerRequest{
				Depositor: userAddress,
				Amount:    sdk.NewIntFromUint64(1000000),
				AssetId:   1,
				AppId:     1,
			},
			1000000,
			true,
			lockerTypes.QueryLockerInfoRequest{
				Id: 1,
			},
			lockerTypes.ErrorUserLockerAlreadyExists,
		},
		{
			"CreateLocker : ErrorAssetDoesNotExist",
			lockerTypes.MsgCreateLockerRequest{
				Depositor: userAddress,
				Amount:    sdk.NewIntFromUint64(1000000),
				AssetId:   10,
				AppId:     1,
			},
			1000000,
			true,
			lockerTypes.QueryLockerInfoRequest{
				Id: 1,
			},
			lockerTypes.ErrorAssetDoesNotExist,
		},
		{
			"CreateLocker : ErrorAppMappingDoesNotExist",
			lockerTypes.MsgCreateLockerRequest{
				Depositor: userAddress,
				Amount:    sdk.NewIntFromUint64(1000000),
				AssetId:   1,
				AppId:     10,
			},
			1000000,
			true,
			lockerTypes.QueryLockerInfoRequest{
				Id: 1,
			},
			lockerTypes.ErrorAppMappingDoesNotExist,
		},
		{
			"CreateLocker : ErrorAppMappingDoesNotExist",
			lockerTypes.MsgCreateLockerRequest{
				Depositor: userAddress,
				Amount:    sdk.NewIntFromUint64(1000000),
				AssetId:   2,
				AppId:     2,
			},
			1000000,
			true,
			lockerTypes.QueryLockerInfoRequest{
				Id: 1,
			},
			lockerTypes.ErrorCollectorLookupDoesNotExists,
		},
		{
			"CreateLocker : App1 Asset 2",
			lockerTypes.MsgCreateLockerRequest{
				Depositor: userAddress,
				Amount:    sdk.NewIntFromUint64(2000000),
				AssetId:   2,
				AppId:     1,
			},
			2000000,
			false,
			lockerTypes.QueryLockerInfoRequest{
				Id: 2,
			},
			nil,
		},
		{
			"CreateLocker : App2 Asset 1",
			lockerTypes.MsgCreateLockerRequest{
				Depositor: userAddress,
				Amount:    sdk.NewIntFromUint64(9900000),
				AssetId:   1,
				AppId:     2,
			},
			9900000,
			false,
			lockerTypes.QueryLockerInfoRequest{
				Id: 3,
			},
			nil,
		},
	} {
		s.Run(tc.name, func() {
			if tc.msg.AssetId == 1 {
				s.fundAddr(userAddress, sdk.NewCoin("upetri", sdk.NewIntFromUint64(tc.fundAmount)))
			} else {
				s.fundAddr(userAddress, sdk.NewCoin("ufust", sdk.NewIntFromUint64(tc.fundAmount)))
			}
			_, err := server.MsgCreateLocker(sdk.WrapSDKContext(*ctx), &tc.msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				lockerInfo, err := s.querier.QueryLockerInfo(sdk.WrapSDKContext(*ctx), &tc.query)
				s.Require().NoError(err)
				s.Require().Equal(lockerInfo.LockerInfo.Depositor, tc.msg.Depositor)
				s.Require().Equal(lockerInfo.LockerInfo.AppId, tc.msg.AppId)
				s.Require().Equal(lockerInfo.LockerInfo.AssetDepositId, tc.msg.AssetId)
				s.Require().Equal(lockerInfo.LockerInfo.NetBalance, tc.msg.Amount)
			}
		})
	}
}

func (s *KeeperTestSuite) TestDepositLocker() {
	userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	userAddress1 := "cosmos1fg240kge022yxh9yu5k5krhru9564u5cmrc57h"
	lockerKeeper, ctx := &s.lockerKeeper, &s.ctx
	// s.AddAppAsset()
	s.TestCreateLocker()
	server := keeper.NewMsgServer(*lockerKeeper)
	for _, tc := range []struct {
		name          string
		msg           lockerTypes.MsgDepositAssetRequest
		query         lockerTypes.QueryLockerInfoRequest
		fundAmount    uint64
		expectedError bool
		ExpErr        error
	}{
		{
			"DepositLocker : ErrorAssetDoesNotExist",
			lockerTypes.MsgDepositAssetRequest{
				Depositor: userAddress,
				LockerId:  1,
				Amount:    sdk.NewIntFromUint64(4000000),
				AssetId:   10,
				AppId:     1,
			},
			lockerTypes.QueryLockerInfoRequest{
				Id: 1,
			},
			9900000,
			false,
			lockerTypes.ErrorAssetDoesNotExist,
		},
		{
			"DepositLocker : ErrorAppMappingDoesNotExist",
			lockerTypes.MsgDepositAssetRequest{
				Depositor: userAddress,
				LockerId:  1,
				Amount:    sdk.NewIntFromUint64(4000000),
				AssetId:   1,
				AppId:     10,
			},
			lockerTypes.QueryLockerInfoRequest{
				Id: 1,
			},
			9900000,
			false,
			lockerTypes.ErrorAppMappingDoesNotExist,
		},
		{
			"DepositLocker : ErrorLockerDoesNotExists",
			lockerTypes.MsgDepositAssetRequest{
				Depositor: userAddress,
				LockerId:  10,
				Amount:    sdk.NewIntFromUint64(4000000),
				AssetId:   1,
				AppId:     1,
			},
			lockerTypes.QueryLockerInfoRequest{
				Id: 1,
			},
			9900000,
			false,
			lockerTypes.ErrorLockerDoesNotExists,
		},
		{
			"DepositLocker : ErrorInvalidAssetID",
			lockerTypes.MsgDepositAssetRequest{
				Depositor: userAddress,
				LockerId:  1,
				Amount:    sdk.NewIntFromUint64(4000000),
				AssetId:   2,
				AppId:     1,
			},
			lockerTypes.QueryLockerInfoRequest{
				Id: 1,
			},
			9900000,
			false,
			lockerTypes.ErrorInvalidAssetID,
		},
		{
			"DepositLocker : ErrorUnauthorized",
			lockerTypes.MsgDepositAssetRequest{
				Depositor: userAddress1,
				LockerId:  1,
				Amount:    sdk.NewIntFromUint64(4000000),
				AssetId:   1,
				AppId:     1,
			},
			lockerTypes.QueryLockerInfoRequest{
				Id: 1,
			},
			9900000,
			false,
			lockerTypes.ErrorUnauthorized,
		},
		{
			"DepositLocker : ErrorAppMappingDoesNotExist 2",
			lockerTypes.MsgDepositAssetRequest{
				Depositor: userAddress,
				LockerId:  1,
				Amount:    sdk.NewIntFromUint64(4000000),
				AssetId:   1,
				AppId:     2,
			},
			lockerTypes.QueryLockerInfoRequest{
				Id: 1,
			},
			9900000,
			false,
			lockerTypes.ErrorAppMappingDoesNotExist,
		},
		{
			"DepositLocker : App1 Asset 1",
			lockerTypes.MsgDepositAssetRequest{
				Depositor: userAddress,
				LockerId:  1,
				Amount:    sdk.NewIntFromUint64(4000000),
				AssetId:   1,
				AppId:     1,
			},
			lockerTypes.QueryLockerInfoRequest{
				Id: 1,
			},
			9900000,
			false,
			nil,
		},
		{
			"DepositLocker : App2 Asset 1",
			lockerTypes.MsgDepositAssetRequest{
				Depositor: userAddress,
				LockerId:  3,
				Amount:    sdk.NewIntFromUint64(4000000),
				AssetId:   1,
				AppId:     2,
			},
			lockerTypes.QueryLockerInfoRequest{
				Id: 3,
			},
			9900000,
			false,
			nil,
		},
		{
			"DepositLocker : App2 Asset 1",
			lockerTypes.MsgDepositAssetRequest{
				Depositor: userAddress,
				LockerId:  3,
				Amount:    sdk.NewIntFromUint64(9223372036854775807),
				AssetId:   1,
				AppId:     2,
			},
			lockerTypes.QueryLockerInfoRequest{
				Id: 3,
			},
			9223372036854775807,
			false,
			nil,
		},
	} {
		s.Run(tc.name, func() {
			if tc.msg.AssetId == 1 {
				s.fundAddr(userAddress, sdk.NewCoin("upetri", sdk.NewIntFromUint64(tc.fundAmount)))
			} else {
				s.fundAddr(userAddress, sdk.NewCoin("ufust", sdk.NewIntFromUint64(tc.fundAmount)))
			}
			lockerInfo, err := s.querier.QueryLockerInfo(sdk.WrapSDKContext(*ctx), &tc.query)
			s.Require().NoError(err)
			previousNetAmount := lockerInfo.LockerInfo.NetBalance
			_, err = server.MsgDepositAsset(sdk.WrapSDKContext(*ctx), &tc.msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				lockerInfo, err := s.querier.QueryLockerInfo(sdk.WrapSDKContext(*ctx), &tc.query)
				s.Require().NoError(err)
				s.Require().Equal(lockerInfo.LockerInfo.Depositor, tc.msg.Depositor)
				s.Require().Equal(lockerInfo.LockerInfo.AppId, tc.msg.AppId)
				s.Require().Equal(lockerInfo.LockerInfo.AssetDepositId, tc.msg.AssetId)
				s.Require().Equal(lockerInfo.LockerInfo.NetBalance, tc.msg.Amount.Add(previousNetAmount))
			}
		})
	}
}

func (s *KeeperTestSuite) TestWithdrawLocker() {
	userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	userAddress1 := "cosmos1fg240kge022yxh9yu5k5krhru9564u5cmrc57h"

	lockerKeeper, ctx := &s.lockerKeeper, &s.ctx
	// s.AddAppAsset()
	s.TestCreateLocker()
	server := keeper.NewMsgServer(*lockerKeeper)
	for _, tc := range []struct {
		name          string
		msg           lockerTypes.MsgWithdrawAssetRequest
		query         lockerTypes.QueryLockerInfoRequest
		expectedError bool
		partial       bool
		ExpErr        error
	}{
		{
			"WithdrawLocker : ErrorAssetDoesNotExist",
			lockerTypes.MsgWithdrawAssetRequest{
				Depositor: userAddress,
				LockerId:  1,
				Amount:    sdk.NewIntFromUint64(100000),
				AssetId:   10,
				AppId:     1,
			},
			lockerTypes.QueryLockerInfoRequest{
				Id: 1,
			},
			false,
			true,
			lockerTypes.ErrorAssetDoesNotExist,
		},
		{
			"WithdrawLocker : ErrorAppMappingDoesNotExist",
			lockerTypes.MsgWithdrawAssetRequest{
				Depositor: userAddress,
				LockerId:  1,
				Amount:    sdk.NewIntFromUint64(100000),
				AssetId:   1,
				AppId:     10,
			},
			lockerTypes.QueryLockerInfoRequest{
				Id: 1,
			},
			false,
			true,
			lockerTypes.ErrorAppMappingDoesNotExist,
		},
		{
			"WithdrawLocker : ErrorLockerDoesNotExists",
			lockerTypes.MsgWithdrawAssetRequest{
				Depositor: userAddress,
				LockerId:  10,
				Amount:    sdk.NewIntFromUint64(100000),
				AssetId:   1,
				AppId:     1,
			},
			lockerTypes.QueryLockerInfoRequest{
				Id: 1,
			},
			false,
			true,
			lockerTypes.ErrorLockerDoesNotExists,
		},
		{
			"WithdrawLocker : ErrorInvalidAssetID",
			lockerTypes.MsgWithdrawAssetRequest{
				Depositor: userAddress,
				LockerId:  1,
				Amount:    sdk.NewIntFromUint64(100000),
				AssetId:   2,
				AppId:     1,
			},
			lockerTypes.QueryLockerInfoRequest{
				Id: 1,
			},
			false,
			true,
			lockerTypes.ErrorInvalidAssetID,
		},
		{
			"WithdrawLocker : ErrorAppMappingDoesNotExist",
			lockerTypes.MsgWithdrawAssetRequest{
				Depositor: userAddress,
				LockerId:  1,
				Amount:    sdk.NewIntFromUint64(100000),
				AssetId:   1,
				AppId:     2,
			},
			lockerTypes.QueryLockerInfoRequest{
				Id: 1,
			},
			false,
			true,
			lockerTypes.ErrorAppMappingDoesNotExist,
		},
		{
			"WithdrawLocker : ErrorUnauthorized",
			lockerTypes.MsgWithdrawAssetRequest{
				Depositor: userAddress1,
				LockerId:  1,
				Amount:    sdk.NewIntFromUint64(100000),
				AssetId:   1,
				AppId:     1,
			},
			lockerTypes.QueryLockerInfoRequest{
				Id: 1,
			},
			false,
			true,
			lockerTypes.ErrorUnauthorized,
		},
		{
			"WithdrawLocker : ErrorRequestedAmountExceedsDepositAmount",
			lockerTypes.MsgWithdrawAssetRequest{
				Depositor: userAddress,
				LockerId:  1,
				Amount:    sdk.NewIntFromUint64(9223372036854775807),
				AssetId:   1,
				AppId:     1,
			},
			lockerTypes.QueryLockerInfoRequest{
				Id: 1,
			},
			false,
			true,
			lockerTypes.ErrorRequestedAmountExceedsDepositAmount,
		},
		{
			"WithdrawLocker : Partial withdraw App1 Asset 1",
			lockerTypes.MsgWithdrawAssetRequest{
				Depositor: userAddress,
				LockerId:  1,
				Amount:    sdk.NewIntFromUint64(100000),
				AssetId:   1,
				AppId:     1,
			},
			lockerTypes.QueryLockerInfoRequest{
				Id: 1,
			},
			false,
			true,
			nil,
		},
		{
			"WithdrawLocker : Full Withdraw App1 Asset 1",
			lockerTypes.MsgWithdrawAssetRequest{
				Depositor: userAddress,
				LockerId:  1,
				Amount:    sdk.NewIntFromUint64(900000),
				AssetId:   1,
				AppId:     1,
			},
			lockerTypes.QueryLockerInfoRequest{
				Id: 1,
			},
			false,
			false,
			nil,
		},
		{
			"WithdrawLocker : Full Withdraw App2 Asset 1",
			lockerTypes.MsgWithdrawAssetRequest{
				Depositor: userAddress,
				LockerId:  3,
				Amount:    sdk.NewIntFromUint64(9900000),
				AssetId:   1,
				AppId:     2,
			},
			lockerTypes.QueryLockerInfoRequest{
				Id: 3,
			},
			false,
			false,
			nil,
		},
	} {
		s.Run(tc.name, func() {
			lockerInfo, err := s.querier.QueryLockerInfo(sdk.WrapSDKContext(*ctx), &tc.query)
			s.Require().NoError(err)
			previousNetAmount := lockerInfo.LockerInfo.NetBalance
			_, err = server.MsgWithdrawAsset(sdk.WrapSDKContext(*ctx), &tc.msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				lockerInfo, err := s.querier.QueryLockerInfo(sdk.WrapSDKContext(*ctx), &tc.query)
				s.Require().NoError(err)
				s.Require().Equal(lockerInfo.LockerInfo.Depositor, tc.msg.Depositor)
				s.Require().Equal(lockerInfo.LockerInfo.AppId, tc.msg.AppId)
				s.Require().Equal(lockerInfo.LockerInfo.AssetDepositId, tc.msg.AssetId)
				if tc.partial {
					s.Require().Equal(lockerInfo.LockerInfo.NetBalance, previousNetAmount.Sub(tc.msg.Amount))
				} else {
					s.Require().Equal(lockerInfo.LockerInfo.NetBalance, sdk.NewIntFromUint64(0))
				}
			}
		})
	}
}

func (s *KeeperTestSuite) TestCloseLocker() {
	userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	userAddress1 := "cosmos1fg240kge022yxh9yu5k5krhru9564u5cmrc57h"

	lockerKeeper, ctx := &s.lockerKeeper, &s.ctx
	// s.AddAppAsset()
	s.TestCreateLocker()
	server := keeper.NewMsgServer(*lockerKeeper)
	for _, tc := range []struct {
		name          string
		msg           lockerTypes.MsgCloseLockerRequest
		expectedError bool
		partial       bool
		ExpErr        error
	}{
		{
			"WithdrawLocker : ErrorAssetDoesNotExist",
			lockerTypes.MsgCloseLockerRequest{
				Depositor: userAddress,
				LockerId:  1,
				AssetId:   10,
				AppId:     1,
			},
			false,
			true,
			lockerTypes.ErrorAssetDoesNotExist,
		},
		{
			"WithdrawLocker : ErrorAppMappingDoesNotExist",
			lockerTypes.MsgCloseLockerRequest{
				Depositor: userAddress,
				LockerId:  1,
				AssetId:   1,
				AppId:     10,
			},
			false,
			true,
			lockerTypes.ErrorAppMappingDoesNotExist,
		},
		{
			"WithdrawLocker : ErrorLockerDoesNotExists",
			lockerTypes.MsgCloseLockerRequest{
				Depositor: userAddress,
				LockerId:  10,
				AssetId:   1,
				AppId:     1,
			},
			false,
			true,
			lockerTypes.ErrorLockerDoesNotExists,
		},
		{
			"WithdrawLocker : ErrorInvalidAssetID",
			lockerTypes.MsgCloseLockerRequest{
				Depositor: userAddress,
				LockerId:  1,
				AssetId:   2,
				AppId:     1,
			},
			false,
			true,
			lockerTypes.ErrorInvalidAssetID,
		},
		{
			"WithdrawLocker : ErrorAppMappingDoesNotExist",
			lockerTypes.MsgCloseLockerRequest{
				Depositor: userAddress,
				LockerId:  1,
				AssetId:   1,
				AppId:     2,
			},
			false,
			true,
			lockerTypes.ErrorAppMappingDoesNotExist,
		},
		{
			"WithdrawLocker : ErrorUnauthorized",
			lockerTypes.MsgCloseLockerRequest{
				Depositor: userAddress1,
				LockerId:  1,
				AssetId:   1,
				AppId:     1,
			},
			false,
			true,
			lockerTypes.ErrorUnauthorized,
		},
		{
			"WithdrawLocker : success",
			lockerTypes.MsgCloseLockerRequest{
				Depositor: userAddress,
				LockerId:  1,
				AssetId:   1,
				AppId:     1,
			},
			false,
			true,
			nil,
		},
	} {
		s.Run(tc.name, func() {

			_, err := server.MsgCloseLocker(sdk.WrapSDKContext(*ctx), &tc.msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				//lockerInfo, err := s.querier.QueryLockerInfo(sdk.WrapSDKContext(*ctx), &tc.query)
				//s.Require().NoError(err)
				//s.Require().Equal(lockerInfo.LockerInfo.Depositor, tc.msg.Depositor)
				//s.Require().Equal(lockerInfo.LockerInfo.AppId, tc.msg.AppId)
				//s.Require().Equal(lockerInfo.LockerInfo.AssetDepositId, tc.msg.AssetId)
			}
		})
	}
}

//SetNetFeeCollectedData

func (s *KeeperTestSuite) TestSetNetFeeCollectedData() {
	collectorkeeper, ctx := &s.collector, &s.ctx
	err := collectorkeeper.SetNetFeeCollectedData(*ctx, 1, 1, sdk.NewInt(1000000000))
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestFundModule() {
	err := s.app.BankKeeper.MintCoins(s.ctx, lockerTypes.ModuleName, sdk.NewCoins(sdk.NewCoin("upetri", sdk.NewIntFromUint64(10000000000))))
	s.Require().NoError(err)
	err = s.app.BankKeeper.SendCoinsFromModuleToModule(s.ctx, lockerTypes.ModuleName, collectortypes.ModuleName, sdk.NewCoins(sdk.NewCoin("upetri", sdk.NewIntFromUint64(10000000000))))
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TestLockerRewardCalc() {
	userAddress := "cosmos1q7q90qsl9g0gl2zz0njxwv2a649yqrtyxtnv3v"
	lockerKeeper, ctx := &s.lockerKeeper, &s.ctx
	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2022-03-01T12:00:00Z"))
	s.ctx = s.ctx.WithBlockHeight(10)
	s.TestCreateLocker()
	s.TestSetNetFeeCollectedData()
	s.TestFundModule()
	s.ctx = s.ctx.WithBlockTime(utils.ParseTime("2022-03-02T12:00:00Z"))
	s.ctx = s.ctx.WithBlockHeight(15)
	//s.TestCreateIntRewards()
	rewardskeeper, ctx := &s.rewardsKeeper, &s.ctx
	err := rewardskeeper.WhitelistAssetForInternalRewards(*ctx, 1, 1)
	s.Require().NoError(err)
	server := keeper.NewMsgServer(*lockerKeeper)
	for _, tc := range []struct {
		name          string
		msg           lockerTypes.MsgLockerRewardCalcRequest
		expectedError bool
		partial       bool
		query         lockerTypes.QueryLockerInfoRequest
		ExpErr        error
	}{
		{
			"WithdrawLocker : success",
			lockerTypes.MsgLockerRewardCalcRequest{
				From:     userAddress,
				AppId:    1,
				LockerId: 1,
			},
			false,
			true,
			lockerTypes.QueryLockerInfoRequest{
				Id: 1,
			},
			nil,
		},
		{
			"WithdrawLocker : success",
			lockerTypes.MsgLockerRewardCalcRequest{
				From:     userAddress,
				AppId:    1,
				LockerId: 1,
			},
			false,
			true,
			lockerTypes.QueryLockerInfoRequest{
				Id: 1,
			},
			nil,
		},
	} {
		s.Run(tc.name, func() {

			_, err := server.MsgLockerRewardCalc(sdk.WrapSDKContext(*ctx), &tc.msg)
			if tc.ExpErr != nil {
				s.Require().Error(err)
				s.Require().EqualError(err, tc.ExpErr.Error())
			} else {
				s.Require().NoError(err)
				lockerInfo, err := s.querier.QueryLockerInfo(sdk.WrapSDKContext(*ctx), &tc.query)
				s.Require().NoError(err)
				s.Require().NotEqual(lockerInfo.LockerInfo.ReturnsAccumulated, sdk.ZeroInt())
				fmt.Println("ReturnsAccumulated", lockerInfo.LockerInfo.ReturnsAccumulated)
			}
		})
	}
}
